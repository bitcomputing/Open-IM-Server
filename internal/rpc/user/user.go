package user

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"

	friendclient "Open_IM/internal/rpc/friend/client"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbFriend "Open_IM/pkg/proto/friend"
	sdkws "Open_IM/pkg/proto/sdk_ws"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"gorm.io/gorm"
)

type userServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	pbUser.UnimplementedUserServer
	friendClient friendclient.FriendClient
}

func NewUserServer(port int) *userServer {
	// log.NewPrivateLog(constant.LogFileName)
	return &userServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImUserName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
		friendClient:    friendclient.NewFriendClient(config.ConvertClientConfig(config.Config.ClientConfigs.Friend)),
	}
}

func (s *userServer) RegisterLegacyDiscovery() {
	rpcRegisterIP, err := utils.GetLocalIP()
	if err != nil {
		panic(fmt.Errorf("GetLocalIP failed: %w", err))
	}

	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		logx.Error("RegisterEtcd failed ", err.Error(),
			s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register auth module  rpc to etcd err"))

	}
}

// func (s *userServer) Run() {
// 	log.NewInfo("0", "rpc user start...")

// 	listenIP := ""
// 	if config.Config.ListenIP == "" {
// 		listenIP = "0.0.0.0"
// 	} else {
// 		listenIP = config.Config.ListenIP
// 	}
// 	address := listenIP + ":" + strconv.Itoa(s.rpcPort)

// 	//listener network
// 	listener, err := net.Listen("tcp", address)
// 	if err != nil {
// 		panic("listening err:" + err.Error() + s.rpcRegisterName)
// 	}
// 	log.NewInfo("0", "listen network success, address ", address, listener)
// 	defer listener.Close()
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
// 	//Service registers with etcd
// 	pbUser.RegisterUserServer(srv, s)
// 	rpcRegisterIP := config.Config.RpcRegisterIP
// 	if config.Config.RpcRegisterIP == "" {
// 		rpcRegisterIP, err = utils.GetLocalIP()
// 		if err != nil {
// 			log.Error("", "GetLocalIP failed ", err.Error())
// 		}
// 	}
// 	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
// 	if err != nil {
// 		log.NewError("0", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
// 		panic(utils.Wrap(err, "register user module  rpc to etcd err"))
// 	}
// 	err = srv.Serve(listener)
// 	if err != nil {
// 		log.NewError("0", "Serve failed ", err.Error())
// 		return
// 	}
// 	log.NewInfo("0", "rpc  user success")
// }

func syncPeerUserConversation(ctx context.Context, conversation *pbConversation.Conversation, operationID string) error {
	logger := logx.WithContext(ctx)
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
		logger.Error(err)
		return err
	}
	chat.ConversationSetPrivateNotification(operationID, conversation.OwnerUserID, conversation.UserID, conversation.IsPrivateChat)
	return nil
}

func (s *userServer) GetUserInfo(ctx context.Context, req *pbUser.GetUserInfoReq) (*pbUser.GetUserInfoResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)
	var userInfoList []*sdkws.UserInfo
	if len(req.UserIDList) > 0 {
		for _, userID := range req.UserIDList {
			var userInfo sdkws.UserInfo
			user, err := rocksCache.GetUserInfoFromCache(ctx, userID)
			if err != nil {
				logger.Error(err, userID)
				continue
			}
			utils.CopyStructFields(&userInfo, user)
			userInfo.BirthStr = utils.TimeToString(user.Birth)
			userInfoList = append(userInfoList, &userInfo)
		}
	} else {
		return &pbUser.GetUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
	}
	return &pbUser.GetUserInfoResp{CommonResp: &pbUser.CommonResp{}, UserInfoList: userInfoList}, nil
}

func (s *userServer) BatchSetConversations(ctx context.Context, req *pbUser.BatchSetConversationsReq) (*pbUser.BatchSetConversationsResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)
	if req.NotificationType == 0 {
		req.NotificationType = constant.ConversationOptChangeNotification
	}
	resp := &pbUser.BatchSetConversationsResp{}
	for _, v := range req.Conversations {
		conversation := db.Conversation{}
		if err := utils.CopyStructFields(&conversation, v); err != nil {
			logger.Debug(v.String(), "CopyStructFields failed", err.Error())
		}
		//redis op
		if err := db.DB.SetSingleConversationRecvMsgOpt(ctx, req.OwnerUserID, v.ConversationID, v.RecvMsgOpt); err != nil {
			logger.Error("cache failed, rpc return", err)
			resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}

		if v.ConversationType == constant.SuperGroupChatType {
			if v.RecvMsgOpt == constant.ReceiveNotNotifyMessage {
				if err := db.DB.SetSuperGroupUserReceiveNotNotifyMessage(ctx, v.GroupID, v.OwnerUserID); err != nil {
					logger.Error("cache failed, rpc return", err.Error(), v.GroupID, v.OwnerUserID)
					resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}
					return resp, nil
				}
			} else {
				if err := db.DB.SetSuperGroupUserReceiveNotifyMessage(ctx, v.GroupID, v.OwnerUserID); err != nil {
					logger.Error("cache failed, rpc return", err.Error(), v.GroupID, err.Error())
					resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}
					return resp, nil
				}
			}
		}

		isUpdate, err := imdb.SetConversation(ctx, conversation)
		if err != nil {
			logger.Error(err)
			resp.Failed = append(resp.Failed, v.ConversationID)
			continue
		}
		if isUpdate {
			err = rocksCache.DelConversationFromCache(ctx, v.OwnerUserID, v.ConversationID)
		} else {
			err = rocksCache.DelUserConversationIDListFromCache(ctx, v.OwnerUserID)
		}
		if err != nil {
			logger.Error(err, v.ConversationID, v.OwnerUserID)
		}
		resp.Success = append(resp.Success, v.ConversationID)
		// if is set private msg operation，then peer user need to sync and set tips\
		if v.ConversationType == constant.SingleChatType && req.NotificationType == constant.ConversationPrivateChatNotification {
			if err := syncPeerUserConversation(ctx, v, req.OperationID); err != nil {
				logger.Error(err)
			}
		}
	}
	chat.ConversationChangeNotification(req.OperationID, req.OwnerUserID)
	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) GetAllConversations(ctx context.Context, req *pbUser.GetAllConversationsReq) (*pbUser.GetAllConversationsResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)
	resp := &pbUser.GetAllConversationsResp{Conversations: []*pbConversation.Conversation{}}
	conversations, err := rocksCache.GetUserAllConversationList(ctx, req.OwnerUserID)
	logger.Debug(conversations)
	if err != nil {
		logger.Error(err)
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if err = utils.CopyStructFields(&resp.Conversations, conversations); err != nil {
		logger.Debug(err)
	}
	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) GetConversation(ctx context.Context, req *pbUser.GetConversationReq) (*pbUser.GetConversationResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbUser.GetConversationResp{Conversation: &pbConversation.Conversation{}}
	conversation, err := rocksCache.GetConversationFromCache(ctx, req.OwnerUserID, req.ConversationID)
	logger.Debug(conversation)
	if err != nil {
		logger.Error(err)
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if err := utils.CopyStructFields(resp.Conversation, &conversation); err != nil {
		logger.Debug(err)
	}

	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) GetConversations(ctx context.Context, req *pbUser.GetConversationsReq) (*pbUser.GetConversationsResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbUser.GetConversationsResp{Conversations: []*pbConversation.Conversation{}}
	conversations, err := rocksCache.GetConversationsFromCache(ctx, req.OwnerUserID, req.ConversationIDs)
	logger.Debug(conversations)
	if err != nil {
		logger.Error(err)
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if err := utils.CopyStructFields(&resp.Conversations, conversations); err != nil {
		logger.Debug(err)
	}

	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) SetConversation(ctx context.Context, req *pbUser.SetConversationReq) (*pbUser.SetConversationResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbUser.SetConversationResp{}
	if req.NotificationType == 0 {
		req.NotificationType = constant.ConversationOptChangeNotification
	}
	if req.Conversation.ConversationType == constant.GroupChatType || req.Conversation.ConversationType == constant.SuperGroupChatType {
		groupInfo, err := imdb.GetGroupInfoByGroupID(ctx, req.Conversation.GroupID)
		if err != nil {
			logger.Error(err)
			resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}
		if groupInfo.Status == constant.GroupStatusDismissed && !req.Conversation.IsNotInGroup {
			logger.Info("group status is dismissed", groupInfo)
			errMsg := "group status is dismissed"
			resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}
			return resp, nil
		}

		if req.Conversation.ConversationType == constant.SuperGroupChatType {
			if req.Conversation.RecvMsgOpt == constant.ReceiveNotNotifyMessage {
				if err = db.DB.SetSuperGroupUserReceiveNotNotifyMessage(ctx, req.Conversation.GroupID, req.Conversation.OwnerUserID); err != nil {
					logger.Error("cache failed, rpc return", err.Error(), req.Conversation.GroupID, req.Conversation.OwnerUserID)
					resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}
					return resp, nil
				}
			} else {
				if err = db.DB.SetSuperGroupUserReceiveNotifyMessage(ctx, req.Conversation.GroupID, req.Conversation.OwnerUserID); err != nil {
					logger.Error("cache failed, rpc return: ", err.Error(), req.Conversation.GroupID, err.Error())
					resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}
					return resp, nil
				}
			}
		}
	}

	var conversation db.Conversation
	if err := utils.CopyStructFields(&conversation, req.Conversation); err != nil {
		logger.Debug("CopyStructFields failed: ", *req.Conversation, err.Error())
	}
	if err := db.DB.SetSingleConversationRecvMsgOpt(ctx, req.Conversation.OwnerUserID, req.Conversation.ConversationID, req.Conversation.RecvMsgOpt); err != nil {
		logger.Error("cache failed, rpc return: ", err.Error())
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	isUpdate, err := imdb.SetConversation(ctx, conversation)
	if err != nil {
		logger.Error("SetConversation error: ", err.Error())
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if isUpdate {
		err = rocksCache.DelConversationFromCache(ctx, req.Conversation.OwnerUserID, req.Conversation.ConversationID)
	} else {
		err = rocksCache.DelUserConversationIDListFromCache(ctx, req.Conversation.OwnerUserID)
	}
	if err != nil {
		logger.Error(err.Error(), " ", req.Conversation.ConversationID, req.Conversation.OwnerUserID)
	}

	// notification
	if req.Conversation.ConversationType == constant.SingleChatType && req.NotificationType == constant.ConversationPrivateChatNotification {
		//sync peer user conversation if conversation is singleChatType
		if err := syncPeerUserConversation(ctx, req.Conversation, req.OperationID); err != nil {
			logger.Error("syncPeerUserConversation failed: ", err.Error())
			resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}
	} else {
		chat.ConversationChangeNotification(req.OperationID, req.Conversation.OwnerUserID)
	}

	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) SetRecvMsgOpt(ctx context.Context, req *pbUser.SetRecvMsgOptReq) (*pbUser.SetRecvMsgOptResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbUser.SetRecvMsgOptResp{}
	var conversation db.Conversation
	if err := utils.CopyStructFields(&conversation, req); err != nil {
		logger.Debug("CopyStructFields failed: ", *req, err.Error())
	}
	if err := db.DB.SetSingleConversationRecvMsgOpt(ctx, req.OwnerUserID, req.ConversationID, req.RecvMsgOpt); err != nil {
		logger.Error("cache failed, rpc return: ", err)
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	stringList := strings.Split(req.ConversationID, "_")
	if len(stringList) > 1 {
		switch stringList[0] {
		case "single":
			conversation.UserID = stringList[1]
			conversation.ConversationType = constant.SingleChatType
		case "group":
			conversation.GroupID = stringList[1]
			conversation.ConversationType = constant.GroupChatType
		case "super_group":
			conversation.GroupID = stringList[1]
			if req.RecvMsgOpt == constant.ReceiveNotNotifyMessage {
				if err := db.DB.SetSuperGroupUserReceiveNotNotifyMessage(ctx, conversation.GroupID, req.OwnerUserID); err != nil {
					logger.Error("cache failed, rpc return: ", err.Error(), conversation.GroupID, req.OwnerUserID)
					resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
					return resp, nil
				}
			} else {
				if err := db.DB.SetSuperGroupUserReceiveNotifyMessage(ctx, conversation.GroupID, req.OwnerUserID); err != nil {
					logger.Error("cache failed, rpc return: ", err.Error(), conversation.GroupID, req.OwnerUserID)
					resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
					return resp, nil
				}
			}
		}
	}
	isUpdate, err := imdb.SetRecvMsgOpt(ctx, conversation)
	if err != nil {
		logger.Error("SetConversation error: ", err.Error())
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}

	if isUpdate {
		err = rocksCache.DelConversationFromCache(ctx, conversation.OwnerUserID, conversation.ConversationID)
	} else {
		err = rocksCache.DelUserConversationIDListFromCache(ctx, conversation.OwnerUserID)
	}
	if err != nil {
		logger.Error(conversation.ConversationID, err.Error())
	}
	chat.ConversationChangeNotification(req.OperationID, req.OwnerUserID)

	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) GetAllUserID(ctx context.Context, req *pbUser.GetAllUserIDReq) (*pbUser.GetAllUserIDResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if !token_verify.IsManagerUserID(req.OpUserID) {
		logger.Error("IsManagerUserID false: ", req.OpUserID)
		return &pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	uidList, err := imdb.SelectAllUserID(ctx)
	if err != nil {
		logger.Error("SelectAllUserID false: ", err.Error())
		return &pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	return &pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{}, UserIDList: uidList}, nil
}

func (s *userServer) AccountCheck(ctx context.Context, req *pbUser.AccountCheckReq) (*pbUser.AccountCheckResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if !token_verify.IsManagerUserID(req.OpUserID) {
		logger.Error("IsManagerUserID false: ", req.OpUserID)
		return &pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	uidList, err := imdb.SelectSomeUserID(ctx, req.CheckUserIDList)

	if err != nil {
		logger.Error("SelectSomeUserID failed: ", err.Error(), req.CheckUserIDList)
		return &pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	var r []*pbUser.AccountCheckResp_SingleUserStatus
	for _, v := range req.CheckUserIDList {
		temp := new(pbUser.AccountCheckResp_SingleUserStatus)
		temp.UserID = v
		if utils.IsContain(v, uidList) {
			temp.AccountStatus = constant.Registered
		} else {
			temp.AccountStatus = constant.UnRegistered
		}
		r = append(r, temp)
	}
	resp := pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrCode: 0, ErrMsg: ""}, ResultList: r}

	return &resp, nil

}

func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbUser.UpdateUserInfoReq) (*pbUser.UpdateUserInfoResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if !token_verify.CheckAccess(req.OpUserID, req.UserInfo.UserID) {
		logger.Error("CheckAccess false: ", req.OpUserID, req.UserInfo.UserID)
		return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	oldNickname := ""
	if req.UserInfo.Nickname != "" {
		u, err := imdb.GetUserByUserID(ctx, req.UserInfo.UserID)
		if err == nil {
			oldNickname = u.Nickname
		}
	}
	var user db.User
	utils.CopyStructFields(&user, req.UserInfo)

	if req.UserInfo.BirthStr != "" {
		time, err := utils.TimeStringToTime(req.UserInfo.BirthStr)
		if err != nil {
			logger.Error("TimeStringToTime failed ", err.Error(), req.UserInfo.BirthStr)
			return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: "TimeStringToTime failed:" + err.Error()}}, nil
		}
		user.Birth = time
	}

	err := imdb.UpdateUserInfo(ctx, user)
	if err != nil {
		logger.Error("UpdateUserInfo failed ", err.Error(), user)
		return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	newReq := &pbFriend.GetFriendListReq{
		CommID: &pbFriend.CommID{OperationID: req.OperationID, FromUserID: req.UserInfo.UserID, OpUserID: req.OpUserID},
	}

	rpcResp, err := s.friendClient.GetFriendList(context.Background(), newReq)
	if err != nil {
		logger.Error("GetFriendList failed ", err.Error(), newReq)
		return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, nil
	}
	for _, v := range rpcResp.FriendInfoList {
		logger.Info(req.OperationID, "UserInfoUpdatedNotification ", req.UserInfo.UserID, v.FriendUser.UserID)
		//	chat.UserInfoUpdatedNotification(req.OperationID, req.UserInfo.UserID, v.FriendUser.UserID)
		chat.FriendInfoUpdatedNotification(ctx, req.OperationID, req.UserInfo.UserID, v.FriendUser.UserID, req.OpUserID)
	}
	if err := rocksCache.DelUserInfoFromCache(ctx, user.UserID); err != nil {
		logger.Error("GetFriendList failed ", err.Error(), newReq)
		return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	//chat.UserInfoUpdatedNotification(req.OperationID, req.OpUserID, req.UserInfo.UserID)
	chat.UserInfoUpdatedNotification(ctx, req.OperationID, req.OpUserID, req.UserInfo.UserID)
	logger.Info("UserInfoUpdatedNotification ", req.UserInfo.UserID, req.OpUserID)
	if req.UserInfo.FaceURL != "" {
		s.SyncJoinedGroupMemberFaceURL(ctx, req.UserInfo.UserID, req.UserInfo.FaceURL, req.OperationID, req.OpUserID)
	}
	if req.UserInfo.Nickname != "" {
		s.SyncJoinedGroupMemberNickname(ctx, req.UserInfo.UserID, req.UserInfo.Nickname, oldNickname, req.OperationID, req.OpUserID)
	}

	return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{}}, nil
}
func (s *userServer) SetGlobalRecvMessageOpt(ctx context.Context, req *pbUser.SetGlobalRecvMessageOptReq) (*pbUser.SetGlobalRecvMessageOptResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	var user db.User
	user.UserID = req.UserID
	m := make(map[string]interface{}, 1)

	m["global_recv_msg_opt"] = req.GlobalRecvMsgOpt
	err := db.DB.SetUserGlobalMsgRecvOpt(ctx, user.UserID, req.GlobalRecvMsgOpt)
	if err != nil {
		logger.Error("SetGlobalRecvMessageOpt failed: ", err.Error(), user)
		return &pbUser.SetGlobalRecvMessageOptResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	err = imdb.UpdateUserInfoByMap(ctx, user, m)
	if err != nil {
		logger.Error("SetGlobalRecvMessageOpt failed: ", err.Error(), user)
		return &pbUser.SetGlobalRecvMessageOptResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if err := rocksCache.DelUserInfoFromCache(ctx, user.UserID); err != nil {
		logger.Error("DelUserInfoFromCache failed: ", err.Error(), req.String())
		return &pbUser.SetGlobalRecvMessageOptResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	chat.UserInfoUpdatedNotification(ctx, req.OperationID, req.UserID, req.UserID)
	return &pbUser.SetGlobalRecvMessageOptResp{CommonResp: &pbUser.CommonResp{}}, nil
}

func (s *userServer) SyncJoinedGroupMemberFaceURL(ctx context.Context, userID string, faceURL string, operationID string, opUserID string) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	joinedGroupIDList, err := rocksCache.GetJoinedGroupIDListFromCache(ctx, userID)
	if err != nil {
		logger.Info("GetJoinedGroupIDListByUserID failed ", userID, err.Error())
		return
	}
	for _, groupID := range joinedGroupIDList {
		groupMemberInfo := db.GroupMember{UserID: userID, GroupID: groupID, FaceURL: faceURL}
		if err := imdb.UpdateGroupMemberInfo(ctx, groupMemberInfo); err != nil {
			logger.Error(err.Error(), groupMemberInfo)
			continue
		}
		//if err := rocksCache.DelAllGroupMembersInfoFromCache(groupID); err != nil {
		//	log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID)
		//	continue
		//}
		if err := rocksCache.DelGroupMemberInfoFromCache(ctx, groupID, userID); err != nil {
			logger.Error(err.Error(), groupID, userID)
			continue
		}
		chat.GroupMemberInfoSetNotification(ctx, operationID, opUserID, groupID, userID)
	}
}

func (s *userServer) SyncJoinedGroupMemberNickname(ctx context.Context, userID string, newNickname, oldNickname string, operationID string, opUserID string) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	joinedGroupIDList, err := imdb.GetJoinedGroupIDListByUserID(ctx, userID)
	if err != nil {
		logger.Info("GetJoinedGroupIDListByUserID failed ", userID, err.Error())
		return
	}
	for _, v := range joinedGroupIDList {
		member, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(ctx, v, userID)
		if err != nil {
			logger.Info("GetGroupMemberInfoByGroupIDAndUserID failed ", err.Error(), v, userID)
			continue
		}
		if member.Nickname == oldNickname {
			groupMemberInfo := db.GroupMember{UserID: userID, GroupID: v, Nickname: newNickname}
			if err := imdb.UpdateGroupMemberInfo(ctx, groupMemberInfo); err != nil {
				logger.Error(err.Error(), groupMemberInfo)
				continue
			}
			//if err := rocksCache.DelAllGroupMembersInfoFromCache(v); err != nil {
			//	log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), v)
			//	continue
			//}
			if err := rocksCache.DelGroupMemberInfoFromCache(ctx, v, userID); err != nil {
				logger.Error(err.Error(), v)
			}
			chat.GroupMemberInfoSetNotification(ctx, operationID, opUserID, v, userID)
		}
	}
}

func (s *userServer) GetUsers(ctx context.Context, req *pbUser.GetUsersReq) (*pbUser.GetUsersResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	var usersDB []db.User
	var err error
	resp := &pbUser.GetUsersResp{CommonResp: &pbUser.CommonResp{}, Pagination: &sdkws.ResponsePagination{CurrentPage: req.Pagination.PageNumber, ShowNumber: req.Pagination.ShowNumber}}
	if req.UserID != "" {
		userDB, err := imdb.GetUserByUserID(ctx, req.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, nil
			}
			logger.Error(req.UserID, err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
			return resp, nil
		}
		usersDB = append(usersDB, *userDB)
		resp.TotalNums = 1
	} else if req.UserName != "" {
		usersDB, err = imdb.GetUserByName(ctx, req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber)
		if err != nil {
			logger.Error(req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber, err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
			return resp, nil
		}
		resp.TotalNums, err = imdb.GetUsersCount(ctx, req.UserName)
		if err != nil {
			logger.Error(req.UserName, err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}

	} else if req.Content != "" {
		var count int64
		usersDB, count, err = imdb.GetUsersByNameAndID(ctx, req.Content, req.Pagination.ShowNumber, req.Pagination.PageNumber)
		if err != nil {
			logger.Error("GetUsers failed", req.Pagination.ShowNumber, req.Pagination.PageNumber, err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		resp.TotalNums = int32(count)
	} else {
		usersDB, err = imdb.GetUsers(ctx, req.Pagination.ShowNumber, req.Pagination.PageNumber)
		if err != nil {
			logger.Error("GetUsers failed", req.Pagination.ShowNumber, req.Pagination.PageNumber, err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		resp.TotalNums, err = imdb.GetTotalUserNum(ctx)
		if err != nil {
			logger.Error(err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
	}
	for _, userDB := range usersDB {
		var user sdkws.UserInfo
		utils.CopyStructFields(&user, userDB)
		user.CreateTime = uint32(userDB.CreateTime.Unix())
		user.BirthStr = utils.TimeToString(userDB.Birth)
		resp.UserList = append(resp.UserList, &pbUser.CmsUser{User: &user})
	}

	var userIDList []string
	for _, v := range resp.UserList {
		userIDList = append(userIDList, v.User.UserID)
	}
	isBlockUser, err := imdb.UsersIsBlock(ctx, userIDList)
	if err != nil {
		logger.Error(err.Error(), userIDList)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}

	for _, v := range resp.UserList {
		if utils.IsContain(v.User.UserID, isBlockUser) {
			v.IsBlock = true
		}
	}

	return resp, nil
}

func (s *userServer) AddUser(ctx context.Context, req *pbUser.AddUserReq) (*pbUser.AddUserResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbUser.AddUserResp{CommonResp: &pbUser.CommonResp{}}
	err := imdb.AddUser(ctx, req.UserInfo.UserID, req.UserInfo.PhoneNumber, req.UserInfo.Nickname, req.UserInfo.Email, req.UserInfo.Gender, req.UserInfo.FaceURL, req.UserInfo.BirthStr)
	if err != nil {
		logger.Error(err)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	return resp, nil
}

func (s *userServer) BlockUser(ctx context.Context, req *pbUser.BlockUserReq) (*pbUser.BlockUserResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbUser.BlockUserResp{CommonResp: &pbUser.CommonResp{}}
	err := imdb.BlockUser(ctx, req.UserID, req.EndDisableTime)
	if err != nil {
		logger.Error(err)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}

	return resp, nil
}

func (s *userServer) UnBlockUser(ctx context.Context, req *pbUser.UnBlockUserReq) (*pbUser.UnBlockUserResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbUser.UnBlockUserResp{CommonResp: &pbUser.CommonResp{}}
	err := imdb.UnBlockUser(ctx, req.UserID)
	if err != nil {
		logger.Error(err)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}

	return resp, nil
}

func (s *userServer) GetBlockUsers(ctx context.Context, req *pbUser.GetBlockUsersReq) (resp *pbUser.GetBlockUsersResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbUser.GetBlockUsersResp{CommonResp: &pbUser.CommonResp{}, Pagination: &sdkws.ResponsePagination{ShowNumber: req.Pagination.ShowNumber, CurrentPage: req.Pagination.PageNumber}}
	var blockUsers []imdb.BlockUserInfo
	if req.UserID != "" {
		blockUser, err := imdb.GetBlockUserByID(ctx, req.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, nil
			}
			logger.Error(err)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		blockUsers = append(blockUsers, blockUser)
		resp.UserNums = 1
	} else {
		blockUsers, err = imdb.GetBlockUsers(ctx, req.Pagination.ShowNumber, req.Pagination.PageNumber)
		if err != nil {
			logger.Error("GetBlockUsers", err.Error(), req.Pagination.ShowNumber, req.Pagination.PageNumber)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}

		nums, err := imdb.GetBlockUsersNumCount(ctx)
		if err != nil {
			logger.Error("GetBlockUsersNumCount failed", err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		resp.UserNums = nums
	}
	for _, v := range blockUsers {
		resp.BlockUsers = append(resp.BlockUsers, &pbUser.BlockUser{
			UserInfo: &sdkws.UserInfo{
				FaceURL:     v.User.FaceURL,
				Nickname:    v.User.Nickname,
				UserID:      v.User.UserID,
				PhoneNumber: v.User.PhoneNumber,
				Email:       v.User.Email,
				Gender:      v.User.Gender,
			},
			BeginDisableTime: (v.BeginDisableTime).String(),
			EndDisableTime:   (v.EndDisableTime).String(),
		})
	}

	return resp, nil
}
