package conversation

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbConversation "Open_IM/pkg/proto/conversation"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"Open_IM/pkg/common/config"
)

type rpcConversation struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	pbConversation.UnimplementedConversationServer
}

func (rpc *rpcConversation) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) (*pbConversation.ModifyConversationFieldResp, error) {
	logger := logx.WithContext(ctx).WithFields(logx.Field("op", req.OperationID))
	resp := &pbConversation.ModifyConversationFieldResp{}
	var err error
	isSyncConversation := true
	if req.Conversation.ConversationType == constant.GroupChatType {
		groupInfo, err := imdb.GetGroupInfoByGroupID(ctx, req.Conversation.GroupID)
		if err != nil {
			logger.Error(req.Conversation.GroupID, err)
			resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}
		if groupInfo.Status == constant.GroupStatusDismissed && !req.Conversation.IsNotInGroup && req.FieldType != constant.FieldUnread {
			errMsg := "group status is dismissed"
			resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}
			return resp, nil
		}
	}
	var conversation db.Conversation
	if err := utils.CopyStructFields(&conversation, req.Conversation); err != nil {
		logger.Debug(req.Conversation, err)
	}
	haveUserID, _ := imdb.GetExistConversationUserIDList(ctx, req.UserIDList, req.Conversation.ConversationID)
	switch req.FieldType {
	case constant.FieldRecvMsgOpt:
		for _, v := range req.UserIDList {
			if err = db.DB.SetSingleConversationRecvMsgOpt(ctx, v, req.Conversation.ConversationID, req.Conversation.RecvMsgOpt); err != nil {
				logger.Error(err)
				resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
				return resp, nil
			}
			if req.Conversation.ConversationType == constant.SuperGroupChatType {
				if req.Conversation.RecvMsgOpt == constant.ReceiveNotNotifyMessage {
					if err = db.DB.SetSuperGroupUserReceiveNotNotifyMessage(ctx, req.Conversation.GroupID, v); err != nil {
						logger.Error(err, req.Conversation.GroupID, v)
						resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
						return resp, nil
					}
				} else {
					if err = db.DB.SetSuperGroupUserReceiveNotifyMessage(ctx, req.Conversation.GroupID, v); err != nil {
						logger.Error(err, req.Conversation.GroupID, v)
						resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
						return resp, nil
					}
				}
			}
		}
		err = imdb.UpdateColumnsConversations(ctx, haveUserID, req.Conversation.ConversationID, map[string]interface{}{"recv_msg_opt": conversation.RecvMsgOpt})
	case constant.FieldGroupAtType:
		err = imdb.UpdateColumnsConversations(ctx, haveUserID, req.Conversation.ConversationID, map[string]interface{}{"group_at_type": conversation.GroupAtType})
	case constant.FieldIsNotInGroup:
		err = imdb.UpdateColumnsConversations(ctx, haveUserID, req.Conversation.ConversationID, map[string]interface{}{"is_not_in_group": conversation.IsNotInGroup})
	case constant.FieldIsPinned:
		err = imdb.UpdateColumnsConversations(ctx, haveUserID, req.Conversation.ConversationID, map[string]interface{}{"is_pinned": conversation.IsPinned})
	case constant.FieldIsPrivateChat:
		err = imdb.UpdateColumnsConversations(ctx, haveUserID, req.Conversation.ConversationID, map[string]interface{}{"is_private_chat": conversation.IsPrivateChat})
	case constant.FieldEx:
		err = imdb.UpdateColumnsConversations(ctx, haveUserID, req.Conversation.ConversationID, map[string]interface{}{"ex": conversation.Ex})
	case constant.FieldAttachedInfo:
		err = imdb.UpdateColumnsConversations(ctx, haveUserID, req.Conversation.ConversationID, map[string]interface{}{"attached_info": conversation.AttachedInfo})
	case constant.FieldUnread:
		isSyncConversation = false
		err = imdb.UpdateColumnsConversations(ctx, haveUserID, req.Conversation.ConversationID, map[string]interface{}{"update_unread_count_time": conversation.UpdateUnreadCountTime})
	case constant.FieldBurnDuration:
		err = imdb.UpdateColumnsConversations(ctx, haveUserID, req.Conversation.ConversationID, map[string]interface{}{"burn_duration": conversation.BurnDuration})
	}
	if err != nil {
		logger.Error(err)
		resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	for _, v := range utils.DifferenceString(haveUserID, req.UserIDList) {
		conversation.OwnerUserID = v
		err = rocksCache.DelUserConversationIDListFromCache(ctx, v)
		if err != nil {
			logger.Error(v, req.Conversation.ConversationID, err.Error())
		}
		err := imdb.SetOneConversation(ctx, conversation)
		if err != nil {
			logger.Error("SetConversation error", err.Error())
			resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}
	}

	// notification
	if req.Conversation.ConversationType == constant.SingleChatType && req.FieldType == constant.FieldIsPrivateChat {
		//sync peer user conversation if conversation is singleChatType
		if err := syncPeerUserConversation(ctx, req.Conversation, req.OperationID); err != nil {
			logger.Error("syncPeerUserConversation", err.Error())
			resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}

	} else {
		if isSyncConversation {
			for _, v := range req.UserIDList {
				if err = rocksCache.DelConversationFromCache(ctx, v, req.Conversation.ConversationID); err != nil {
					logger.Error(v, req.Conversation.ConversationID, err.Error())
				}
				chat.ConversationChangeNotification(ctx, req.OperationID, v)
			}
		} else {
			for _, v := range req.UserIDList {
				if err = rocksCache.DelConversationFromCache(ctx, v, req.Conversation.ConversationID); err != nil {
					logger.Error(v, req.Conversation.ConversationID, err.Error())
				}
				chat.ConversationUnreadChangeNotification(ctx, req.OperationID, v, req.Conversation.ConversationID, conversation.UpdateUnreadCountTime)
			}
		}

	}

	resp.CommonResp = &pbConversation.CommonResp{}
	return resp, nil
}

func syncPeerUserConversation(ctx context.Context, conversation *pbConversation.Conversation, operationID string) error {
	logger := logx.WithContext(ctx).WithFields(logx.Field("op", operationID))
	peerUserConversation := db.Conversation{
		OwnerUserID:      conversation.UserID,
		ConversationID:   utils.GetConversationIDBySessionType(conversation.OwnerUserID, constant.SingleChatType),
		ConversationType: constant.SingleChatType,
		UserID:           conversation.OwnerUserID,
		GroupID:          "",
		RecvMsgOpt:       0,
		UnreadCount:      0,
		DraftTextTime:    0,
		IsPinned:         false,
		IsPrivateChat:    conversation.IsPrivateChat,
		AttachedInfo:     "",
		Ex:               "",
	}
	err := imdb.PeerUserSetConversation(ctx, peerUserConversation)
	if err != nil {
		logger.Error("SetConversation error", err.Error())
		return err
	}
	err = rocksCache.DelUserConversationIDListFromCache(ctx, conversation.UserID)
	if err != nil {
		logger.Error("DelConversationFromCache failed", err.Error(), conversation.OwnerUserID, conversation.ConversationID)
	}
	err = rocksCache.DelConversationFromCache(ctx, conversation.UserID, utils.GetConversationIDBySessionType(conversation.OwnerUserID, constant.SingleChatType))
	if err != nil {
		logger.Error("DelConversationFromCache failed", err.Error(), conversation.OwnerUserID, conversation.ConversationID)
	}
	err = rocksCache.DelConversationFromCache(ctx, conversation.OwnerUserID, conversation.ConversationID)
	if err != nil {
		logger.Error("DelConversationFromCache failed", err.Error(), conversation.OwnerUserID, conversation.ConversationID)
	}
	chat.ConversationSetPrivateNotification(ctx, operationID, conversation.OwnerUserID, conversation.UserID, conversation.IsPrivateChat)
	return nil
}

func NewRpcConversationServer(port int) *rpcConversation {
	return &rpcConversation{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImConversationName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (r *rpcConversation) RegisterLegacyDiscovery() {
	rpcRegisterIP, err := utils.GetLocalIP()
	if err != nil {
		panic(fmt.Errorf("GetLocalIP failed: %w", err))
	}

	err = getcdv3.RegisterEtcd(r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName, 10)
	if err != nil {
		logx.Error("RegisterEtcd failed ", err.Error(),
			r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName)
		panic(utils.Wrap(err, "register auth module  rpc to etcd err"))
	}
}

// func (rpc *rpcConversation) Run() {
// 	log.NewInfo("0", "rpc conversation start...")

// 	listenIP := ""
// 	if config.Config.ListenIP == "" {
// 		listenIP = "0.0.0.0"
// 	} else {
// 		listenIP = config.Config.ListenIP
// 	}
// 	address := listenIP + ":" + strconv.Itoa(rpc.rpcPort)

// 	listener, err := net.Listen("tcp", address)
// 	if err != nil {
// 		panic("listening err:" + err.Error() + rpc.rpcRegisterName)
// 	}
// 	log.NewInfo("0", "listen network success, ", address, listener)
// 	//grpc server
// 	var grpcOpts []grpc.ServerOption
// 	if config.Config.Prometheus.Enable {
// 		promePkg.NewGrpcRequestCounter()
// 		promePkg.NewGrpcRequestFailedCounter()
// 		promePkg.NewGrpcRequestSuccessCounter()
// 		grpcOpts = append(grpcOpts, []grpc.ServerOption{
// 			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
// 			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
// 			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
// 		}...)
// 	}
// 	srv := grpc.NewServer(grpcOpts...)
// 	defer srv.GracefulStop()

// 	//service registers with etcd
// 	pbConversation.RegisterConversationServer(srv, rpc)
// 	rpcRegisterIP := config.Config.RpcRegisterIP
// 	if config.Config.RpcRegisterIP == "" {
// 		rpcRegisterIP, err = utils.GetLocalIP()
// 		if err != nil {
// 			log.Error("", "GetLocalIP failed ", err.Error())
// 		}
// 	}
// 	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
// 	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
// 	if err != nil {
// 		log.NewError("0", "RegisterEtcd failed ", err.Error(),
// 			rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
// 		panic(utils.Wrap(err, "register conversation module  rpc to etcd err"))
// 	}
// 	log.NewInfo("0", "RegisterConversationServer ok ", rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
// 	err = srv.Serve(listener)
// 	if err != nil {
// 		log.NewError("0", "Serve failed ", err.Error())
// 		return
// 	}
// 	log.NewInfo("0", "rpc conversation ok")
// }
