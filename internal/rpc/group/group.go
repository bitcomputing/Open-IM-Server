package group

import (
	cacheclient "Open_IM/internal/rpc/cache/client"
	conversationclient "Open_IM/internal/rpc/conversation/client"
	chat "Open_IM/internal/rpc/msg"
	userclient "Open_IM/internal/rpc/user/client"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/token_verify"
	cp "Open_IM/pkg/common/utils"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbCache "Open_IM/pkg/proto/cache"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbGroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type groupServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	pbGroup.UnimplementedGroupServer
	conversationClient conversationclient.ConversationClient
	userClient         userclient.UserClient
	cacheClient        cacheclient.CacheClient
}

func NewGroupServer(port int) *groupServer {
	return &groupServer{
		rpcPort:            port,
		rpcRegisterName:    config.Config.RpcRegisterName.OpenImGroupName,
		etcdSchema:         config.Config.Etcd.EtcdSchema,
		etcdAddr:           config.Config.Etcd.EtcdAddr,
		conversationClient: conversationclient.NewConversationClient(config.ConvertClientConfig(config.Config.ClientConfigs.Conversation)),
		userClient:         userclient.NewUserClient(config.ConvertClientConfig(config.Config.ClientConfigs.User)),
		cacheClient:        cacheclient.NewCacheClient(config.ConvertClientConfig(config.Config.ClientConfigs.Cache)),
	}
}

func (s *groupServer) RegisterLegacyDiscovery() {
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

// func (s *groupServer) Run() {
// 	log.NewInfo("", "group rpc start ")
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
// 	log.NewInfo("", "listen network success, ", address, listener)
// 	defer listener.Close()
// 	//grpc server
// 	recvSize := 1024 * 1024 * constant.GroupRPCRecvSize
// 	sendSize := 1024 * 1024 * constant.GroupRPCSendSize
// 	var grpcOpts = []grpc.ServerOption{
// 		grpc.MaxRecvMsgSize(recvSize),
// 		grpc.MaxSendMsgSize(sendSize),
// 	}
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
// 	pbGroup.RegisterGroupServer(srv, s)

// 	rpcRegisterIP := config.Config.RpcRegisterIP
// 	if config.Config.RpcRegisterIP == "" {
// 		rpcRegisterIP, err = utils.GetLocalIP()
// 		if err != nil {
// 			log.Error("", "GetLocalIP failed ", err.Error())
// 		}
// 	}
// 	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
// 	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
// 	if err != nil {
// 		log.NewError("", "RegisterEtcd failed ", err.Error())
// 		panic(utils.Wrap(err, "register group module  rpc to etcd err"))

// 	}
// 	log.Info("", "RegisterEtcd ", s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
// 	err = srv.Serve(listener)
// 	if err != nil {
// 		log.NewError("", "Serve failed ", err.Error())
// 		return
// 	}
// 	log.NewInfo("", "group rpc success")
// }

func (s *groupServer) CreateGroup(ctx context.Context, req *pbGroup.CreateGroupReq) (*pbGroup.CreateGroupResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if !token_verify.CheckAccess(req.OpUserID, req.OwnerUserID) {
		logger.Error("CheckAccess false ", req.OpUserID, req.OwnerUserID)
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	callbackResp := callbackBeforeCreateGroup(ctx, req)
	if callbackResp.ErrCode != 0 {
		logger.Error("callbackBeforeSendSingleMsg resp: ", callbackResp)
	}
	if callbackResp.ActionCode != constant.ActionAllow {
		if callbackResp.ErrCode == 0 {
			callbackResp.ErrCode = 201
		}
		logger.Debug("callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
		return &pbGroup.CreateGroupResp{
			ErrCode: int32(callbackResp.ErrCode),
			ErrMsg:  callbackResp.ErrMsg,
		}, nil
	}

	groupId := req.GroupInfo.GroupID
	if groupId == "" {
		groupId = utils.Md5(req.OperationID + strconv.FormatInt(time.Now().UnixNano(), 10))
		bi := big.NewInt(0)
		bi.SetString(groupId[0:8], 16)
		groupId = bi.String()
	}
	//to group
	groupInfo := db.Group{}
	utils.CopyStructFields(&groupInfo, req.GroupInfo)
	groupInfo.CreatorUserID = req.OpUserID
	groupInfo.GroupID = groupId
	groupInfo.CreateTime = time.Now()
	if groupInfo.NotificationUpdateTime.Unix() < 0 {
		groupInfo.NotificationUpdateTime = utils.UnixSecondToTime(0)
	}
	err := imdb.InsertIntoGroup(ctx, groupInfo)
	if err != nil {
		logger.Error("InsertIntoGroup failed, ", err.Error(), groupInfo)
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}
	var okUserIDList []string
	resp := &pbGroup.CreateGroupResp{GroupInfo: &open_im_sdk.GroupInfo{}}
	groupMember := db.GroupMember{}
	us := &db.User{}
	if req.OwnerUserID != "" {
		var userIDList []string
		for _, v := range req.InitMemberList {
			userIDList = append(userIDList, v.UserID)
		}
		userIDList = append(userIDList, req.OwnerUserID)
		if err := s.DelGroupAndUserCache(ctx, req.OperationID, "", userIDList); err != nil {
			logger.Error("DelGroupAndUserCache failed, ", err.Error(), userIDList)
			return &pbGroup.CreateGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
		}

		us, err = imdb.GetUserByUserID(ctx, req.OwnerUserID)
		if err != nil {
			logger.Error("GetUserByUserID failed ", err.Error(), req.OwnerUserID)
			return &pbGroup.CreateGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
		}
		//to group member
		groupMember = db.GroupMember{GroupID: groupId, RoleLevel: constant.GroupOwner, OperatorUserID: req.OpUserID, JoinSource: constant.JoinByInvitation, InviterUserID: req.OpUserID}
		utils.CopyStructFields(&groupMember, us)
		callbackResp := CallbackBeforeMemberJoinGroup(ctx, req.OperationID, &groupMember, groupInfo.Ex)
		if callbackResp.ErrCode != 0 {
			logger.Error("CallbackBeforeMemberJoinGroup resp: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			if callbackResp.ErrCode == 0 {
				callbackResp.ErrCode = 201
			}
			logger.Debug("CallbackBeforeMemberJoinGroup result", "end rpc and return", callbackResp)
			return &pbGroup.CreateGroupResp{
				ErrCode: int32(callbackResp.ErrCode),
				ErrMsg:  callbackResp.ErrMsg,
			}, nil
		}

		err = imdb.InsertIntoGroupMember(ctx, groupMember)
		if err != nil {
			logger.Error("InsertIntoGroupMember failed ", err.Error(), groupMember)
			return &pbGroup.CreateGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
		}
	}

	if req.GroupInfo.GroupType != constant.SuperGroup {
		//to group member
		for _, user := range req.InitMemberList {
			us, err := rocksCache.GetUserInfoFromCache(ctx, user.UserID)
			if err != nil {
				logger.Error("GetUserByUserID failed ", err.Error(), user.UserID)
				continue
			}
			if user.RoleLevel == constant.GroupOwner {
				logger.Error("only one owner, failed ", user)
				continue
			}
			groupMember.RoleLevel = user.RoleLevel
			groupMember.JoinSource = constant.JoinByInvitation
			groupMember.InviterUserID = req.OpUserID
			utils.CopyStructFields(&groupMember, us)
			callbackResp := CallbackBeforeMemberJoinGroup(ctx, req.OperationID, &groupMember, groupInfo.Ex)
			if callbackResp.ErrCode != 0 {
				logger.Error(utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg resp: ", callbackResp)
			}
			if callbackResp.ActionCode != constant.ActionAllow {
				if callbackResp.ErrCode == 0 {
					callbackResp.ErrCode = 201
				}
				logger.Debug("callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
				continue
			}
			err = imdb.InsertIntoGroupMember(ctx, groupMember)
			if err != nil {
				logger.Error("InsertIntoGroupMember failed ", err.Error(), groupMember)
				continue
			}
			okUserIDList = append(okUserIDList, user.UserID)
		}
		group, err := rocksCache.GetGroupInfoFromCache(ctx, groupId)
		if err != nil {
			logger.Error("GetGroupInfoByGroupID failed ", err.Error(), groupId)
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = err.Error()
			return resp, nil
		}
		utils.CopyStructFields(resp.GroupInfo, group)
		memberCount, err := rocksCache.GetGroupMemberNumFromCache(ctx, groupId)
		resp.GroupInfo.MemberCount = uint32(memberCount)
		if err != nil {
			logger.Error("GetGroupMemberNumByGroupID failed ", err.Error(), groupId)
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = err.Error()
			return resp, nil
		}
		if req.OwnerUserID != "" {
			resp.GroupInfo.OwnerUserID = req.OwnerUserID
			okUserIDList = append(okUserIDList, req.OwnerUserID)
		}
		// superGroup stored in mongodb
	} else {
		for _, v := range req.InitMemberList {
			okUserIDList = append(okUserIDList, v.UserID)
		}
		if err := db.DB.CreateSuperGroup(ctx, groupId, okUserIDList, len(okUserIDList)); err != nil {
			logger.Error("GetGroupMemberNumByGroupID failed ", err.Error(), groupId)
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = err.Error() + ": CreateSuperGroup failed"
			return resp, nil
		}
	}

	if len(okUserIDList) != 0 {
		logger.Info("rpc CreateGroup return ", resp.String())
		if req.GroupInfo.GroupType != constant.SuperGroup {
			chat.GroupCreatedNotification(ctx, req.OperationID, req.OpUserID, groupId, okUserIDList)
		} else {
			for _, userID := range okUserIDList {
				if err := rocksCache.DelJoinedSuperGroupIDListFromCache(ctx, userID); err != nil {
					logger.Error(userID, err.Error())
				}
			}
			go func() {
				for _, v := range okUserIDList {
					chat.SuperGroupNotification(ctx, req.OperationID, v, v)
				}
			}()
		}
		return resp, nil
	} else {
		return resp, nil
	}
}

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (*pbGroup.GetJoinedGroupListResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if !token_verify.CheckAccess(req.OpUserID, req.FromUserID) {
		logger.Error("CheckAccess false: ", req.OpUserID, req.FromUserID)
		return &pbGroup.GetJoinedGroupListResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	joinedGroupList, err := rocksCache.GetJoinedGroupIDListFromCache(ctx, req.FromUserID)
	if err != nil {
		logger.Error("GetJoinedGroupIDListFromCache failed", err.Error(), req.FromUserID)
		return &pbGroup.GetJoinedGroupListResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}

	var resp pbGroup.GetJoinedGroupListResp
	for _, v := range joinedGroupList {
		var groupNode open_im_sdk.GroupInfo
		num, err := rocksCache.GetGroupMemberNumFromCache(ctx, v)
		if err != nil {
			logger.Error(err.Error(), v)
			continue
		}
		owner, err2 := imdb.GetGroupOwnerInfoByGroupID(ctx, v)
		if err2 != nil {
			logger.Error(err2.Error(), v)
			continue
		}
		group, err := rocksCache.GetGroupInfoFromCache(ctx, v)
		if err != nil {
			logger.Error(err.Error(), v)
			continue
		}
		if group.GroupType == constant.SuperGroup {
			continue
		}
		if group.Status == constant.GroupStatusDismissed {
			logger.Error("constant.GroupStatusDismissed ", group)
			continue
		}
		utils.CopyStructFields(&groupNode, group)
		groupNode.CreateTime = uint32(group.CreateTime.Unix())
		groupNode.NotificationUpdateTime = uint32(group.NotificationUpdateTime.Unix())
		if group.NotificationUpdateTime.Unix() < 0 {
			groupNode.NotificationUpdateTime = 0
		}

		groupNode.MemberCount = uint32(num)
		groupNode.OwnerUserID = owner.UserID
		resp.GroupList = append(resp.GroupList, &groupNode)
		logger.Debug("joinedGroup ", groupNode)
	}

	return &resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbGroup.InviteUserToGroupReq) (*pbGroup.InviteUserToGroupResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if !imdb.IsExistGroupMember(ctx, req.GroupID, req.OpUserID) && !token_verify.IsManagerUserID(req.OpUserID) {
		logger.Error("no permission InviteUserToGroup ", req.GroupID, req.OpUserID)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	groupInfo, err := imdb.GetGroupInfoByGroupID(ctx, req.GroupID)
	if err != nil {
		logger.Error("GetGroupInfoByGroupID failed ", req.GroupID, err)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}, nil
	}
	var resp pbGroup.InviteUserToGroupResp
	if groupInfo.NeedVerification == constant.AllNeedVerification &&
		!imdb.IsGroupOwnerAdmin(ctx, req.GroupID, req.OpUserID) && !token_verify.IsManagerUserID(req.OpUserID) {
		var resp pbGroup.InviteUserToGroupResp
		joinReq := pbGroup.JoinGroupReq{}
		for _, v := range req.InvitedUserIDList {
			if imdb.IsExistGroupMember(ctx, req.GroupID, v) {
				logger.Error("IsExistGroupMember ", req.GroupID, v)
				var resultNode pbGroup.Id2Result
				resultNode.Result = -1
				resultNode.UserID = v
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				continue
			}
			var groupRequest db.GroupRequest
			groupRequest.UserID = v
			groupRequest.GroupID = req.GroupID
			groupRequest.JoinSource = constant.JoinByInvitation
			groupRequest.InviterUserID = req.OpUserID
			err = imdb.InsertIntoGroupRequest(ctx, groupRequest)
			if err != nil {
				var resultNode pbGroup.Id2Result
				resultNode.Result = -1
				resultNode.UserID = v
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)

				continue
				// log.NewError(req.OperationID, "InsertIntoGroupRequest failed ", err.Error(), groupRequest)
				//	return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			} else {
				var resultNode pbGroup.Id2Result
				resultNode.Result = 0
				resultNode.UserID = v
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				joinReq.GroupID = req.GroupID
				joinReq.OperationID = req.OperationID
				joinReq.OpUserID = v
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				chat.JoinGroupApplicationNotification(ctx, &joinReq)
			}
		}

		return &resp, nil
	}
	if err := s.DelGroupAndUserCache(ctx, req.OperationID, req.GroupID, req.InvitedUserIDList); err != nil {
		logger.Error("DelGroupAndUserCache failed", err.Error())
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
	}
	//from User:  invite: applicant
	//to user:  invite: invited
	var okUserIDList []string
	if groupInfo.GroupType != constant.SuperGroup {
		for _, v := range req.InvitedUserIDList {
			var resultNode pbGroup.Id2Result
			resultNode.UserID = v
			resultNode.Result = 0
			toUserInfo, err := imdb.GetUserByUserID(ctx, v)
			if err != nil {
				logger.Error("GetUserByUserID failed ", err.Error(), v)
				resultNode.Result = -1
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				continue
			}

			if imdb.IsExistGroupMember(ctx, req.GroupID, v) {
				logger.Error("IsExistGroupMember ", req.GroupID, v)
				resultNode.Result = -1
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				continue
			}
			var toInsertInfo db.GroupMember
			utils.CopyStructFields(&toInsertInfo, toUserInfo)
			toInsertInfo.GroupID = req.GroupID
			toInsertInfo.RoleLevel = constant.GroupOrdinaryUsers
			toInsertInfo.OperatorUserID = req.OpUserID
			toInsertInfo.InviterUserID = req.OpUserID
			toInsertInfo.JoinSource = constant.JoinByInvitation
			callbackResp := CallbackBeforeMemberJoinGroup(ctx, req.OperationID, &toInsertInfo, groupInfo.Ex)
			if callbackResp.ErrCode != 0 {
				logger.Error("callbackBeforeSendSingleMsg resp: ", callbackResp)
			}
			if callbackResp.ActionCode != constant.ActionAllow {
				logger.Debug("callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
				continue
			}

			err = imdb.InsertIntoGroupMember(ctx, toInsertInfo)
			if err != nil {
				logger.Error("InsertIntoGroupMember failed ", req.GroupID, toUserInfo.UserID, toUserInfo.Nickname, toUserInfo.FaceURL)
				resultNode.Result = -1
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				continue
			}
			okUserIDList = append(okUserIDList, v)
			err = db.DB.AddGroupMember(ctx, req.GroupID, toUserInfo.UserID)
			if err != nil {
				logger.Error("AddGroupMember failed ", err.Error(), req.GroupID, toUserInfo.UserID)
			}
			resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
		}
	} else {
		for _, v := range req.InvitedUserIDList {
			if imdb.IsExistGroupMember(ctx, req.GroupID, v) {
				logger.Error("IsExistGroupMember ", req.GroupID, v)
				var resultNode pbGroup.Id2Result
				resultNode.Result = -1
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				continue
			} else {
				okUserIDList = append(okUserIDList, v)
			}
		}
		//okUserIDList = req.InvitedUserIDList
		if err := db.DB.AddUserToSuperGroup(ctx, req.GroupID, okUserIDList); err != nil {
			logger.Error("AddUserToSuperGroup failed ", req.GroupID, err)
			return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
		}
	}

	// set conversations
	var haveConUserID []string
	var sessionType int
	if groupInfo.GroupType == constant.NormalGroup {
		sessionType = constant.GroupChatType
	} else {
		sessionType = constant.SuperGroupChatType
	}
	conversations, err := imdb.GetConversationsByConversationIDMultipleOwner(ctx, okUserIDList, utils.GetConversationIDBySessionType(req.GroupID, sessionType))
	if err != nil {
		logger.Error("GetConversationsByConversationIDMultipleOwner failed ", err.Error(), req.GroupID, sessionType)
	}
	for _, v := range conversations {
		haveConUserID = append(haveConUserID, v.OwnerUserID)
	}
	var reqPb pbUser.SetConversationReq
	var c pbConversation.Conversation
	for _, v := range conversations {
		reqPb.OperationID = req.OperationID
		c.OwnerUserID = v.OwnerUserID
		c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, sessionType)
		c.RecvMsgOpt = v.RecvMsgOpt
		c.ConversationType = int32(sessionType)
		c.GroupID = req.GroupID
		c.IsPinned = v.IsPinned
		c.AttachedInfo = v.AttachedInfo
		c.IsPrivateChat = v.IsPrivateChat
		c.GroupAtType = v.GroupAtType
		c.IsNotInGroup = false
		c.Ex = v.Ex
		reqPb.Conversation = &c
		respPb, err := s.userClient.SetConversation(ctx, &reqPb)
		if err != nil {
			logger.Error("SetConversation rpc failed, ", reqPb.String(), err.Error(), v.OwnerUserID)
		} else {
			logger.Debug("SetConversation success", respPb.String(), v.OwnerUserID)
		}
	}
	for _, v := range utils.DifferenceString(haveConUserID, okUserIDList) {
		reqPb.OperationID = req.OperationID
		c.OwnerUserID = v
		c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, sessionType)
		c.ConversationType = int32(sessionType)
		c.GroupID = req.GroupID
		c.IsNotInGroup = false
		c.UpdateUnreadCountTime = utils.GetCurrentTimestampByMill()
		reqPb.Conversation = &c
		respPb, err := s.userClient.SetConversation(ctx, &reqPb)
		if err != nil {
			logger.Error("SetConversation rpc failed, ", reqPb.String(), err.Error(), v)
		} else {
			logger.Debug("SetConversation success", respPb.String(), v)
		}
	}

	if groupInfo.GroupType != constant.SuperGroup {
		chat.MemberInvitedNotification(ctx, req.OperationID, req.GroupID, req.OpUserID, req.Reason, okUserIDList)
	} else {
		for _, v := range req.InvitedUserIDList {
			if err := rocksCache.DelJoinedSuperGroupIDListFromCache(ctx, v); err != nil {
				logger.Error(err.Error())
			}
		}
		for _, v := range req.InvitedUserIDList {
			chat.SuperGroupNotification(ctx, req.OperationID, v, v)
		}
	}

	return &resp, nil
}

func (s *groupServer) InviteUserToGroups(ctx context.Context, req *pbGroup.InviteUserToGroupsReq) (*pbGroup.InviteUserToGroupsResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if !token_verify.IsManagerUserID(req.OpUserID) {
		logger.Error("no permission InviteUserToGroup ", req.String())
		return &pbGroup.InviteUserToGroupsResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	for _, v := range req.GroupIDList {
		groupInfo, err := imdb.GetGroupInfoByGroupID(ctx, v)
		if err != nil {
			logger.Error("GetGroupInfoByGroupID failed ", v, err)
			return &pbGroup.InviteUserToGroupsResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error() + v}, nil
		}
		if groupInfo.Status == constant.GroupStatusDismissed {
			errMsg := " group status is dismissed " + v
			return &pbGroup.InviteUserToGroupsResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}, nil
		}
	}
	if err := db.DB.AddUserToSuperGroups(ctx, req.GroupIDList, req.InvitedUserID); err != nil {
		logger.Error("AddUserToSuperGroups failed ", err.Error())
		return &pbGroup.InviteUserToGroupsResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
	}
	if err := rocksCache.DelJoinedSuperGroupIDListFromCache(ctx, req.InvitedUserID); err != nil {
		logger.Error(err)
	}
	chat.SuperGroupNotification(ctx, req.OperationID, req.InvitedUserID, req.InvitedUserID)

	return nil, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	var resp pbGroup.GetGroupAllMemberResp
	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		logger.Error(err.Error(), req.GroupID)
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = constant.ErrDB.ErrMsg
		return &resp, nil
	}
	if groupInfo.GroupType != constant.SuperGroup {
		memberList, err := rocksCache.GetGroupMembersInfoFromCache(ctx, req.Count, req.Offset, req.GroupID)
		if err != nil {
			logger.Error(err.Error(), req.GroupID)
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = constant.ErrDB.ErrMsg
			return &resp, nil
		}
		for _, v := range memberList {
			var node open_im_sdk.GroupMemberFullInfo
			cp.GroupMemberDBCopyOpenIM(ctx, &node, v)
			resp.MemberList = append(resp.MemberList, &node)
		}
	}

	return &resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	var resp pbGroup.GetGroupMemberListResp
	memberList, err := imdb.GetGroupMemberByGroupID(ctx, req.GroupID, req.Filter, req.NextSeq, 30)
	if err != nil {
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = constant.ErrDB.ErrMsg
		logger.Error("GetGroupMemberByGroupId failed,", req.GroupID, req.Filter, req.NextSeq, 30)
		return &resp, nil
	}

	for _, v := range memberList {
		var node open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(&node, &v)
		resp.MemberList = append(resp.MemberList, &node)
	}
	//db operate  get db sorted by join time
	if int32(len(memberList)) < 30 {
		resp.NextSeq = 0
	} else {
		resp.NextSeq = req.NextSeq + int32(len(memberList))
	}
	resp.ErrCode = 0

	return &resp, nil
}

func (s *groupServer) getGroupUserLevel(ctx context.Context, groupID, userID string) (int, error) {
	opFlag := 0
	if !token_verify.IsManagerUserID(userID) {
		opInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(ctx, groupID, userID)
		if err != nil {
			return opFlag, utils.Wrap(err, "")
		}
		if opInfo.RoleLevel == constant.GroupOrdinaryUsers {
			opFlag = 0
		} else if opInfo.RoleLevel == constant.GroupOwner {
			opFlag = 2 //owner
		} else {
			opFlag = 3 //admin
		}
	} else {
		opFlag = 1 //app manager
	}
	return opFlag, nil
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbGroup.KickGroupMemberReq) (*pbGroup.KickGroupMemberResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		logger.Error("GetGroupInfoByGroupID", req.GroupID, err.Error())
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}
	var okUserIDList []string
	var resp pbGroup.KickGroupMemberResp
	if groupInfo.GroupType != constant.SuperGroup {
		opFlag := 0
		if !token_verify.IsManagerUserID(req.OpUserID) {
			opInfo, err := rocksCache.GetGroupMemberInfoFromCache(ctx, req.GroupID, req.OpUserID)
			if err != nil {
				errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID  failed " + err.Error() + req.GroupID + req.OpUserID
				logger.Error(errMsg)
				return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
			}
			if opInfo.RoleLevel == constant.GroupOrdinaryUsers {
				errMsg := req.OperationID + " opInfo.RoleLevel == constant.GroupOrdinaryUsers " + opInfo.UserID + opInfo.GroupID
				logger.Error(errMsg)
				return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
			} else if opInfo.RoleLevel == constant.GroupOwner {
				opFlag = 2 //owner
			} else {
				opFlag = 3 //admin
			}
		} else {
			opFlag = 1 //app manager
		}

		//op is group owner?
		if len(req.KickedUserIDList) == 0 {
			logger.Error("failed, kick list 0")
			return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}, nil
		}
		if err := s.DelGroupAndUserCache(ctx, req.OperationID, req.GroupID, req.KickedUserIDList); err != nil {
			logger.Error("DelGroupAndUserCache failed", err.Error())
			return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
		}
		//remove
		for _, v := range req.KickedUserIDList {
			kickedInfo, err := rocksCache.GetGroupMemberInfoFromCache(ctx, req.GroupID, v)
			if err != nil {
				logger.Error(" GetGroupMemberInfoByGroupIDAndUserID failed ", req.GroupID, v, err.Error())
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				continue
			}

			if kickedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
				logger.Error("is constant.GroupAdmin, can't kicked ", v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				continue
			}
			if kickedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
				logger.Debug("is constant.GroupOwner, can't kicked ", v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				continue
			}

			err = imdb.DeleteGroupMemberByGroupIDAndUserID(ctx, req.GroupID, v)
			if err != nil {
				logger.Error("RemoveGroupMember failed ", err.Error(), req.GroupID, v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
			} else {
				logger.Debug("kicked ", v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: 0})
				okUserIDList = append(okUserIDList, v)
			}
		}
		var reqPb pbUser.SetConversationReq
		var c pbConversation.Conversation
		for _, v := range okUserIDList {
			reqPb.OperationID = req.OperationID
			c.OwnerUserID = v
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsNotInGroup = true
			reqPb.Conversation = &c
			respPb, err := s.userClient.SetConversation(ctx, &reqPb)
			if err != nil {
				logger.Error("SetConversation rpc failed, ", reqPb.String(), err.Error(), v)
			} else {
				logger.Debug("SetConversation success", respPb.String(), v)
			}
		}
	} else {
		okUserIDList = req.KickedUserIDList
		if err := db.DB.RemoverUserFromSuperGroup(ctx, req.GroupID, okUserIDList); err != nil {
			logger.Error(req.GroupID, req.KickedUserIDList, err.Error())
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = constant.ErrDB.ErrMsg
			return &resp, nil
		}
		if err := rocksCache.DelGroupMemberListHashFromCache(ctx, req.GroupID); err != nil {
			logger.Error(req.GroupID, err.Error())
		}
		if err := rocksCache.DelGroupMemberIDListFromCache(ctx, req.GroupID); err != nil {
			logger.Error(err.Error(), req.GroupID)
		}
		reqPb := pbConversation.ModifyConversationFieldReq{Conversation: &pbConversation.Conversation{}}
		reqPb.OperationID = req.OperationID
		reqPb.UserIDList = okUserIDList
		reqPb.FieldType = constant.FieldUnread
		reqPb.Conversation.GroupID = req.GroupID
		reqPb.Conversation.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.SuperGroupChatType)
		reqPb.Conversation.ConversationType = int32(constant.SuperGroupChatType)
		reqPb.Conversation.UpdateUnreadCountTime = utils.GetCurrentTimestampByMill()
		respPb, err := s.conversationClient.ModifyConversationField(ctx, &reqPb)
		if err != nil {
			logger.Error("ModifyConversationField rpc failed, ", reqPb.String(), err.Error())
		} else {
			logger.Debug("ModifyConversationField success", respPb.String())
		}

	}

	if groupInfo.GroupType != constant.SuperGroup {
		for _, userID := range okUserIDList {
			if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, userID); err != nil {
				logger.Error(err.Error(), req.GroupID)
			}
		}
		chat.MemberKickedNotification(ctx, req, okUserIDList)
	} else {
		for _, userID := range okUserIDList {
			err = rocksCache.DelJoinedSuperGroupIDListFromCache(ctx, userID)
			if err != nil {
				logger.Error(err.Error(), userID)
			}
		}
		go func() {
			for _, v := range req.KickedUserIDList {
				chat.SuperGroupNotification(ctx, req.OperationID, v, v)
			}
		}()

	}

	return &resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (*pbGroup.GetGroupMembersInfoResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	var resp pbGroup.GetGroupMembersInfoResp
	resp.MemberList = []*open_im_sdk.GroupMemberFullInfo{}
	for _, userID := range req.MemberList {
		var (
			groupMember *db.GroupMember
			err         error
		)
		if req.NoCache {
			groupMember, err = imdb.GetGroupMemberInfoByGroupIDAndUserID(ctx, req.GroupID, userID)
		} else {
			groupMember, err = rocksCache.GetGroupMemberInfoFromCache(ctx, req.GroupID, userID)
		}
		if err != nil {
			logger.Error(req.GroupID, userID, err.Error())
			continue
		}
		var memberNode open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(&memberNode, groupMember)
		memberNode.JoinTime = int32(groupMember.JoinTime.Unix())
		resp.MemberList = append(resp.MemberList, &memberNode)
	}

	return &resp, nil
}

func (s *groupServer) GetGroupApplicationList(ctx context.Context, req *pbGroup.GetGroupApplicationListReq) (*pbGroup.GetGroupApplicationListResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	reply, err := imdb.GetGroupApplicationList(ctx, req.FromUserID)
	if err != nil {
		logger.Error("GetGroupApplicationList failed ", err.Error(), req.FromUserID)
		return &pbGroup.GetGroupApplicationListResp{ErrCode: 701, ErrMsg: "GetGroupApplicationList failed"}, nil
	}

	resp := pbGroup.GetGroupApplicationListResp{}
	for _, v := range reply {
		node := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		group, err := imdb.GetGroupInfoByGroupID(ctx, v.GroupID)
		if err != nil {
			logger.Error("GetGroupInfoByGroupID failed ", err.Error(), v.GroupID)
			continue
		}
		if group.Status == constant.GroupStatusDismissed {
			logger.Debug("group constant.GroupStatusDismissed  ", group.GroupID)
			continue
		}
		user, err := imdb.GetUserByUserID(ctx, v.UserID)
		if err != nil {
			logger.Error("GetUserByUserID failed ", err.Error(), v.UserID)
			continue
		}

		cp.GroupRequestDBCopyOpenIM(&node, &v)
		cp.UserDBCopyOpenIMPublicUser(node.UserInfo, user)
		cp.GroupDBCopyOpenIM(ctx, node.GroupInfo, group)
		resp.GroupRequestList = append(resp.GroupRequestList, &node)
	}

	return &resp, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	groupsInfoList := make([]*open_im_sdk.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		groupInfoFromRedis, err := rocksCache.GetGroupInfoFromCache(ctx, groupID)
		if err != nil {
			logger.Error("GetGroupInfoByGroupID failed ", err.Error(), groupID)
			continue
		}
		var groupInfo open_im_sdk.GroupInfo
		cp.GroupDBCopyOpenIM(ctx, &groupInfo, groupInfoFromRedis)
		//groupInfo.NeedVerification

		groupInfo.NeedVerification = groupInfoFromRedis.NeedVerification
		groupsInfoList = append(groupsInfoList, &groupInfo)
	}

	resp := pbGroup.GetGroupsInfoResp{GroupInfoList: groupsInfoList}

	return &resp, nil
}

func (s *groupServer) GroupApplicationResponse(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (*pbGroup.GroupApplicationResponseResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	groupRequest := db.GroupRequest{}
	utils.CopyStructFields(&groupRequest, req)
	groupRequest.UserID = req.FromUserID
	groupRequest.HandleUserID = req.OpUserID
	groupRequest.HandledTime = time.Now()
	if !token_verify.IsManagerUserID(req.OpUserID) && !imdb.IsGroupOwnerAdmin(ctx, req.GroupID, req.OpUserID) {
		logger.Error("IsManagerUserID IsGroupOwnerAdmin false ", req.GroupID, req.OpUserID)
		return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	err := imdb.UpdateGroupRequest(ctx, groupRequest)
	if err != nil {
		logger.Error("GroupApplicationResponse failed ", err.Error(), groupRequest)
		return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	request, err := imdb.GetGroupRequestByGroupIDAndUserID(ctx, req.GroupID, req.FromUserID)
	if err != nil {
		logger.Error("GroupApplicationResponse failed ", err.Error(), req.GroupID, req.FromUserID)
		return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		logger.Error("GetGroupInfoFromCache failed ", err.Error(), req.GroupID, req.FromUserID)
		return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if req.HandleResult == constant.GroupResponseAgree {
		user, err := imdb.GetUserByUserID(ctx, req.FromUserID)
		if err != nil {
			logger.Error("GroupApplicationResponse failed ", err.Error(), req.FromUserID)
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		if imdb.IsExistGroupMember(ctx, req.GroupID, req.FromUserID) {
			logger.Info("GroupApplicationResponse user in group", req.GroupID, req.FromUserID)
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{}}, nil
		}
		member := db.GroupMember{}
		member.GroupID = req.GroupID
		member.UserID = req.FromUserID
		member.RoleLevel = constant.GroupOrdinaryUsers
		member.OperatorUserID = req.OpUserID
		member.FaceURL = user.FaceURL
		member.Nickname = user.Nickname
		member.JoinSource = request.JoinSource
		member.InviterUserID = request.InviterUserID
		callbackResp := CallbackBeforeMemberJoinGroup(ctx, req.OperationID, &member, groupInfo.Ex)
		if callbackResp.ErrCode != 0 {
			logger.Error("callbackBeforeSendSingleMsg resp: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			if callbackResp.ErrCode == 0 {
				callbackResp.ErrCode = 201
			}
			logger.Debug("callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
			return &pbGroup.GroupApplicationResponseResp{
				CommonResp: &pbGroup.CommonResp{
					ErrCode: int32(callbackResp.ErrCode),
					ErrMsg:  callbackResp.ErrMsg,
				},
			}, nil
		}
		err = imdb.InsertIntoGroupMember(ctx, member)
		if err != nil {
			logger.Error("GroupApplicationResponse failed ", err.Error(), member)
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		var sessionType int
		if groupInfo.GroupType == constant.NormalGroup {
			sessionType = constant.GroupChatType
		} else {
			sessionType = constant.SuperGroupChatType
		}
		var reqPb pbUser.SetConversationReq
		reqPb.OperationID = req.OperationID
		var c pbConversation.Conversation
		conversation, err := imdb.GetConversation(ctx, req.FromUserID, utils.GetConversationIDBySessionType(req.GroupID, sessionType))
		if err != nil {
			c.OwnerUserID = req.FromUserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, sessionType)
			c.ConversationType = int32(sessionType)
			c.GroupID = req.GroupID
			c.IsNotInGroup = false
			c.UpdateUnreadCountTime = utils.GetCurrentTimestampByMill()
		} else {
			c.OwnerUserID = conversation.OwnerUserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, sessionType)
			c.RecvMsgOpt = conversation.RecvMsgOpt
			c.ConversationType = int32(sessionType)
			c.GroupID = req.GroupID
			c.IsPinned = conversation.IsPinned
			c.AttachedInfo = conversation.AttachedInfo
			c.IsPrivateChat = conversation.IsPrivateChat
			c.GroupAtType = conversation.GroupAtType
			c.IsNotInGroup = false
			c.Ex = conversation.Ex
		}
		reqPb.Conversation = &c
		respPb, err := s.userClient.SetConversation(ctx, &reqPb)
		if err != nil {
			logger.Error("SetConversation rpc failed, ", reqPb.String(), err.Error())
		} else {
			logger.Debug("SetConversation success", respPb.String())
		}

		cacheResp, err := s.cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{OperationID: req.OperationID, GroupID: req.GroupID})
		if err != nil {
			logger.Error("DelGroupMemberIDListFromCache rpc call failed ", err.Error())
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		if cacheResp.CommonResp.ErrCode != 0 {
			logger.Error("DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		if err := rocksCache.DelGroupMemberListHashFromCache(ctx, req.GroupID); err != nil {
			logger.Error(req.GroupID, err.Error())
		}
		if err := rocksCache.DelJoinedGroupIDListFromCache(ctx, req.FromUserID); err != nil {
			logger.Error(req.FromUserID, err.Error())
		}
		if err := rocksCache.DelGroupMemberNumFromCache(ctx, req.GroupID); err != nil {
			logger.Error(err.Error(), req.GroupID)
		}
		chat.GroupApplicationAcceptedNotification(ctx, req)
		chat.MemberEnterNotification(ctx, req)
	} else if req.HandleResult == constant.GroupResponseRefuse {
		chat.GroupApplicationRejectedNotification(ctx, req)
	} else {
		logger.Error("HandleResult failed ", req.HandleResult)
		return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
	}

	return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{}}, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.JoinGroupResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if imdb.IsExistGroupMember(ctx, req.GroupID, req.OpUserID) {
		logger.Info("IsExistGroupMember", req.GroupID, req.OpUserID)
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{}}, nil
	}
	_, err := imdb.GetUserByUserID(ctx, req.OpUserID)
	if err != nil {
		logger.Error("GetUserByUserID failed ", err.Error(), req.OpUserID)
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		logger.Error("GetGroupInfoByGroupID failed ", req.GroupID, err)
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}}, nil
	}

	if groupInfo.NeedVerification == constant.Directly {
		if groupInfo.GroupType != constant.SuperGroup {
			us, err := imdb.GetUserByUserID(ctx, req.OpUserID)
			if err != nil {
				logger.Error("GetUserByUserID failed ", err.Error(), req.OpUserID)
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			}
			//to group member
			groupMember := db.GroupMember{GroupID: req.GroupID, RoleLevel: constant.GroupOrdinaryUsers, OperatorUserID: req.OpUserID}
			utils.CopyStructFields(&groupMember, us)
			callbackResp := CallbackBeforeMemberJoinGroup(ctx, req.OperationID, &groupMember, groupInfo.Ex)
			if callbackResp.ErrCode != 0 {
				logger.Error("callbackBeforeSendSingleMsg resp: ", callbackResp)
			}
			if callbackResp.ActionCode != constant.ActionAllow {
				if callbackResp.ErrCode == 0 {
					callbackResp.ErrCode = 201
				}
				logger.Debug("callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
				return &pbGroup.JoinGroupResp{
					CommonResp: &pbGroup.CommonResp{
						ErrCode: int32(callbackResp.ErrCode),
						ErrMsg:  callbackResp.ErrMsg,
					},
				}, nil
			}

			if err := s.DelGroupAndUserCache(ctx, req.OperationID, req.GroupID, []string{req.OpUserID}); err != nil {
				logger.Error("DelGroupAndUserCache failed, ", err.Error(), req.OpUserID)
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
			}

			err = imdb.InsertIntoGroupMember(ctx, groupMember)
			if err != nil {
				logger.Error("InsertIntoGroupMember failed ", err.Error(), groupMember)
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			}
			//}

			var sessionType int
			if groupInfo.GroupType == constant.NormalGroup {
				sessionType = constant.GroupChatType
			} else {
				sessionType = constant.SuperGroupChatType
			}
			var reqPb pbUser.SetConversationReq
			var c pbConversation.Conversation
			reqPb.OperationID = req.OperationID
			c.OwnerUserID = req.OpUserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, sessionType)
			c.ConversationType = int32(sessionType)
			c.GroupID = req.GroupID
			c.IsNotInGroup = false
			c.UpdateUnreadCountTime = utils.GetCurrentTimestampByMill()
			reqPb.Conversation = &c
			respPb, err := s.userClient.SetConversation(ctx, &reqPb)
			if err != nil {
				logger.Error("SetConversation rpc failed, ", reqPb.String(), err.Error())
			} else {
				logger.Debug("SetConversation success", respPb.String())
			}

			chat.MemberEnterDirectlyNotification(ctx, req.GroupID, req.OpUserID, req.OperationID)

			return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
		} else {
			logger.Error("JoinGroup rpc failed, group type:  ", groupInfo.GroupType, "not support directly")
			return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
		}
	}
	var groupRequest db.GroupRequest
	groupRequest.UserID = req.OpUserID
	groupRequest.ReqMsg = req.ReqMessage
	groupRequest.GroupID = req.GroupID
	groupRequest.JoinSource = req.JoinSource
	err = imdb.InsertIntoGroupRequest(ctx, groupRequest)
	if err != nil {
		logger.Error("InsertIntoGroupRequest failed ", err.Error(), groupRequest)
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	//_, err = imdb.GetGroupMemberListByGroupIDAndRoleLevel(req.GroupID, constant.GroupOwner)
	//if err != nil {
	//	log.NewError(req.OperationID, "GetGroupMemberListByGroupIDAndRoleLevel failed ", err.Error(), req.GroupID, constant.GroupOwner)
	//	return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil

	chat.JoinGroupApplicationNotification(ctx, req)

	return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.QuitGroupResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	groupInfo, err := imdb.GetGroupInfoByGroupID(ctx, req.GroupID)
	if err != nil {
		logger.Error("ReduceGroupMemberFromCache rpc call failed ", err.Error())
		return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.GroupType != constant.SuperGroup {
		_, err = imdb.GetGroupMemberInfoByGroupIDAndUserID(ctx, req.GroupID, req.OpUserID)
		if err != nil {
			logger.Error("GetGroupMemberInfoByGroupIDAndUserID failed ", err.Error(), req.GroupID, req.OpUserID)
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}

		if err := s.DelGroupAndUserCache(ctx, req.OperationID, req.GroupID, []string{req.OpUserID}); err != nil {
			logger.Error("DelGroupAndUserCache failed, ", err.Error(), req.OpUserID)
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
		}

		err = imdb.DeleteGroupMemberByGroupIDAndUserID(ctx, req.GroupID, req.OpUserID)
		if err != nil {
			logger.Error("DeleteGroupMemberByGroupIdAndUserId failed ", err.Error(), req.GroupID, req.OpUserID)
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}

		err = db.DB.DelGroupMember(ctx, req.GroupID, req.OpUserID)
		if err != nil {
			logger.Error("DelGroupMember failed ", req.GroupID, req.OpUserID)
			//	return &pbGroup.CommonResp{ErrorCode: constant.ErrQuitGroup.ErrCode, ErrorMsg: constant.ErrQuitGroup.ErrMsg}, nil
		}
		//modify quitter conversation info
		var reqPb pbUser.SetConversationReq
		var c pbConversation.Conversation
		reqPb.OperationID = req.OperationID
		c.OwnerUserID = req.OpUserID
		c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
		c.ConversationType = constant.GroupChatType
		c.GroupID = req.GroupID
		c.IsNotInGroup = true
		reqPb.Conversation = &c
		respPb, err := s.userClient.SetConversation(ctx, &reqPb)
		if err != nil {
			logger.Error("SetConversation rpc failed, ", reqPb.String(), err.Error())
		} else {
			logger.Debug("SetConversation success", respPb.String())
		}
	} else {
		okUserIDList := []string{req.OpUserID}
		if err := db.DB.RemoverUserFromSuperGroup(ctx, req.GroupID, okUserIDList); err != nil {
			logger.Error(req.GroupID, okUserIDList, err.Error())
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}

	if groupInfo.GroupType != constant.SuperGroup {
		if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.OpUserID); err != nil {
			logger.Error(err.Error(), req.GroupID)
		}
		chat.MemberQuitNotification(ctx, req)
	} else {
		if err := rocksCache.DelJoinedSuperGroupIDListFromCache(ctx, req.OpUserID); err != nil {
			logger.Error(err.Error(), req.OpUserID)
		}
		if err := rocksCache.DelGroupMemberListHashFromCache(ctx, req.GroupID); err != nil {
			logger.Error(req.GroupID, err.Error())
		}
		chat.SuperGroupNotification(ctx, req.OperationID, req.OpUserID, req.OpUserID)
	}

	return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{}}, nil
}

func hasAccess(ctx context.Context, req *pbGroup.SetGroupInfoReq) bool {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if utils.IsContain(req.OpUserID, config.Config.Manager.AppManagerUid) {
		return true
	}
	groupUserInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(ctx, req.GroupInfoForSet.GroupID, req.OpUserID)
	if err != nil {
		logger.Error("GetGroupMemberInfoByGroupIDAndUserID failed, ", err.Error(), req.GroupInfoForSet.GroupID, req.OpUserID)
		return false

	}
	if groupUserInfo.RoleLevel == constant.GroupOwner || groupUserInfo.RoleLevel == constant.GroupAdmin {
		return true
	}
	return false
}

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbGroup.SetGroupInfoReq) (*pbGroup.SetGroupInfoResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if !hasAccess(ctx, req) {
		logger.Error("no access ", req)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	group, err := imdb.GetGroupInfoByGroupID(ctx, req.GroupInfoForSet.GroupID)
	if err != nil {
		logger.Error("GetGroupInfoByGroupID failed ", err.Error(), req.GroupInfoForSet.GroupID)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	if group.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}}, nil
	}

	////bitwise operators: 0001:groupName; 0010:Notification  0100:Introduction; 1000:FaceUrl; 10000:owner
	var changedType int32
	groupName := ""
	notification := ""
	introduction := ""
	faceURL := ""
	if group.GroupName != req.GroupInfoForSet.GroupName && req.GroupInfoForSet.GroupName != "" {
		changedType = 1
		groupName = req.GroupInfoForSet.GroupName
	}
	if group.Notification != req.GroupInfoForSet.Notification && req.GroupInfoForSet.Notification != "" {
		changedType = changedType | (1 << 1)
		notification = req.GroupInfoForSet.Notification
	}
	if group.Introduction != req.GroupInfoForSet.Introduction && req.GroupInfoForSet.Introduction != "" {
		changedType = changedType | (1 << 2)
		introduction = req.GroupInfoForSet.Introduction
	}
	if group.FaceURL != req.GroupInfoForSet.FaceURL && req.GroupInfoForSet.FaceURL != "" {
		changedType = changedType | (1 << 3)
		faceURL = req.GroupInfoForSet.FaceURL
	}

	if req.GroupInfoForSet.NeedVerification != nil {
		changedType = changedType | (1 << 4)
		m := make(map[string]interface{})
		m["need_verification"] = req.GroupInfoForSet.NeedVerification.Value
		if err := imdb.UpdateGroupInfoDefaultZero(ctx, req.GroupInfoForSet.GroupID, m); err != nil {
			logger.Error("UpdateGroupInfoDefaultZero failed ", err.Error(), m)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}
	if req.GroupInfoForSet.LookMemberInfo != nil {
		changedType = changedType | (1 << 5)
		m := make(map[string]interface{})
		m["look_member_info"] = req.GroupInfoForSet.LookMemberInfo.Value
		if err := imdb.UpdateGroupInfoDefaultZero(ctx, req.GroupInfoForSet.GroupID, m); err != nil {
			logger.Error("UpdateGroupInfoDefaultZero failed ", err.Error(), m)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}
	if req.GroupInfoForSet.ApplyMemberFriend != nil {
		changedType = changedType | (1 << 6)
		m := make(map[string]interface{})
		m["apply_member_friend"] = req.GroupInfoForSet.ApplyMemberFriend.Value
		if err := imdb.UpdateGroupInfoDefaultZero(ctx, req.GroupInfoForSet.GroupID, m); err != nil {
			logger.Error("UpdateGroupInfoDefaultZero failed ", err.Error(), m)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}
	//
	//if req.RoleLevel != nil {
	//
	//}
	//only administrators can set group information
	var groupInfo db.Group
	utils.CopyStructFields(&groupInfo, req.GroupInfoForSet)
	if req.GroupInfoForSet.Notification != "" {
		groupInfo.NotificationUserID = req.OpUserID
		groupInfo.NotificationUpdateTime = time.Now()
	}
	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupInfoForSet.GroupID); err != nil {
		logger.Error("DelGroupInfoFromCache failed ", err.Error(), req.GroupInfoForSet.GroupID)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	err = imdb.SetGroupInfo(ctx, groupInfo)
	if err != nil {
		logger.Error("SetGroupInfo failed ", err.Error(), groupInfo)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	if changedType != 0 {
		chat.GroupInfoSetNotification(ctx, req.OperationID, req.OpUserID, req.GroupInfoForSet.GroupID, groupName, notification,
			introduction, faceURL, req.GroupInfoForSet.NeedVerification, req.GroupInfoForSet.ApplyMemberFriend, req.GroupInfoForSet.LookMemberInfo)
	}
	if req.GroupInfoForSet.Notification != "" {
		//get group member user id
		getGroupMemberIDListFromCacheReq := &pbCache.GetGroupMemberIDListFromCacheReq{OperationID: req.OperationID, GroupID: req.GroupInfoForSet.GroupID}

		cacheResp, err := s.cacheClient.GetGroupMemberIDListFromCache(context.Background(), getGroupMemberIDListFromCacheReq)
		if err != nil {
			logger.Error("GetGroupMemberIDListFromCache rpc call failed ", err.Error())
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}}, nil
		}
		if cacheResp.CommonResp.ErrCode != 0 {
			logger.Error("GetGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}}, nil
		}
		var conversationReq pbConversation.ModifyConversationFieldReq

		conversation := pbConversation.Conversation{
			OwnerUserID:      req.OpUserID,
			ConversationID:   utils.GetConversationIDBySessionType(req.GroupInfoForSet.GroupID, constant.GroupChatType),
			ConversationType: constant.GroupChatType,
			GroupID:          req.GroupInfoForSet.GroupID,
		}
		conversationReq.Conversation = &conversation
		conversationReq.OperationID = req.OperationID
		conversationReq.FieldType = constant.FieldGroupAtType
		conversation.GroupAtType = constant.GroupNotification
		conversationReq.UserIDList = cacheResp.UserIDList
		conversationReply, err := s.conversationClient.ModifyConversationField(ctx, &conversationReq)
		if err != nil {
			logger.Error("ModifyConversationField rpc failed, ", conversationReq.String(), err.Error())
		} else if conversationReply.CommonResp.ErrCode != 0 {
			logger.Error("ModifyConversationField rpc failed, ", conversationReq.String(), conversationReply.String())
		}
	}
	return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}}, nil
}

func (s *groupServer) TransferGroupOwner(ctx context.Context, req *pbGroup.TransferGroupOwnerReq) (*pbGroup.TransferGroupOwnerResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	groupInfo, err := imdb.GetGroupInfoByGroupID(ctx, req.GroupID)
	if err != nil {
		logger.Error("GetGroupInfoByGroupID failed ", req.GroupID, err)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}}, nil
	}

	if req.OldOwnerUserID == req.NewOwnerUserID {
		logger.Error("same owner ", req.OldOwnerUserID, req.NewOwnerUserID)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
	}
	err = rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.NewOwnerUserID)
	if err != nil {
		logger.Error("DelGroupMemberInfoFromCache failed ", req.GroupID, req.NewOwnerUserID)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil

	}
	err = rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.OldOwnerUserID)
	if err != nil {
		logger.Error("DelGroupMemberInfoFromCache failed ", req.GroupID, req.OldOwnerUserID)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	groupMemberInfo := db.GroupMember{GroupID: req.GroupID, UserID: req.OldOwnerUserID, RoleLevel: constant.GroupOrdinaryUsers}
	err = imdb.UpdateGroupMemberInfo(ctx, groupMemberInfo)
	if err != nil {
		logger.Error("UpdateGroupMemberInfo failed ", groupMemberInfo)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	groupMemberInfo = db.GroupMember{GroupID: req.GroupID, UserID: req.NewOwnerUserID, RoleLevel: constant.GroupOwner}
	err = imdb.UpdateGroupMemberInfo(ctx, groupMemberInfo)
	if err != nil {
		logger.Error("UpdateGroupMemberInfo failed ", groupMemberInfo)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	chat.GroupOwnerTransferredNotification(ctx, req)
	return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil

}

func (s *groupServer) GetGroups(ctx context.Context, req *pbGroup.GetGroupsReq) (*pbGroup.GetGroupsResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbGroup.GetGroupsResp{
		CommonResp: &pbGroup.CommonResp{},
		CMSGroups:  []*pbGroup.CMSGroup{},
		Pagination: &open_im_sdk.ResponsePagination{CurrentPage: req.Pagination.PageNumber, ShowNumber: req.Pagination.ShowNumber},
	}
	if req.GroupID != "" {
		groupInfoDB, err := imdb.GetGroupInfoByGroupID(ctx, req.GroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, nil
			}
			logger.Error(err.Error(), req.GroupID)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		resp.GroupNum = 1
		groupInfo := &open_im_sdk.GroupInfo{}
		utils.CopyStructFields(groupInfo, groupInfoDB)
		groupMember, err := imdb.GetGroupOwnerInfoByGroupID(ctx, req.GroupID)
		if err != nil {
			logger.Error(err.Error(), req.GroupID)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		memberNum, err := imdb.GetGroupMembersCount(ctx, req.GroupID, "")
		if err != nil {
			logger.Error(err.Error(), req.GroupID)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		groupInfo.MemberCount = uint32(memberNum)
		groupInfo.CreateTime = uint32(groupInfoDB.CreateTime.Unix())
		resp.CMSGroups = append(resp.CMSGroups, &pbGroup.CMSGroup{GroupInfo: groupInfo, GroupOwnerUserName: groupMember.Nickname, GroupOwnerUserID: groupMember.UserID})
	} else {
		groups, count, err := imdb.GetGroupsByName(ctx, req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
		if err != nil {
			logger.Error("GetGroupsByName error", req.String(), req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
		}
		for _, v := range groups {
			group := &pbGroup.CMSGroup{GroupInfo: &open_im_sdk.GroupInfo{}}
			utils.CopyStructFields(group.GroupInfo, v)
			groupMember, err := imdb.GetGroupOwnerInfoByGroupID(ctx, v.GroupID)
			if err != nil {
				logger.Error("GetGroupOwnerInfoByGroupID failed", err.Error(), v)
				continue
			}
			group.GroupInfo.CreateTime = uint32(v.CreateTime.Unix())
			group.GroupOwnerUserID = groupMember.UserID
			group.GroupOwnerUserName = groupMember.Nickname
			resp.CMSGroups = append(resp.CMSGroups, group)
		}
		resp.GroupNum = int32(count)
	}

	return resp, nil
}

func (s *groupServer) GetGroupMembersCMS(ctx context.Context, req *pbGroup.GetGroupMembersCMSReq) (*pbGroup.GetGroupMembersCMSResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbGroup.GetGroupMembersCMSResp{CommonResp: &pbGroup.CommonResp{}}
	groupMembers, err := imdb.GetGroupMembersByGroupIdCMS(ctx, req.GroupID, req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		logger.Error("GetGroupMembersByGroupIdCMS Error", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	groupMembersCount, err := imdb.GetGroupMembersCount(ctx, req.GroupID, req.UserName)
	if err != nil {
		logger.Error("GetGroupMembersCMS Error", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}

	resp.MemberNums = int32(groupMembersCount)
	for _, groupMember := range groupMembers {
		member := open_im_sdk.GroupMemberFullInfo{}
		utils.CopyStructFields(&member, groupMember)
		member.JoinTime = int32(groupMember.JoinTime.Unix())
		member.MuteEndTime = uint32(groupMember.MuteEndTime.Unix())
		resp.Members = append(resp.Members, &member)
	}
	resp.Pagination = &open_im_sdk.ResponsePagination{
		CurrentPage: req.Pagination.PageNumber,
		ShowNumber:  req.Pagination.ShowNumber,
	}

	return resp, nil
}

func (s *groupServer) GetUserReqApplicationList(ctx context.Context, req *pbGroup.GetUserReqApplicationListReq) (*pbGroup.GetUserReqApplicationListResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbGroup.GetUserReqApplicationListResp{}
	groupRequests, err := imdb.GetUserReqGroupByUserID(ctx, req.UserID)
	if err != nil {
		logger.Error("GetUserReqGroupByUserID failed ", err.Error())
		resp.CommonResp = &pbGroup.CommonResp{
			ErrCode: constant.ErrDB.ErrCode,
			ErrMsg:  constant.ErrDB.ErrMsg,
		}
		return resp, nil
	}
	for _, groupReq := range groupRequests {
		node := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		group, err := imdb.GetGroupInfoByGroupID(ctx, groupReq.GroupID)
		if err != nil {
			logger.Error("GetGroupInfoByGroupID failed ", err.Error(), groupReq.GroupID)
			continue
		}
		user, err := imdb.GetUserByUserID(ctx, groupReq.UserID)
		if err != nil {
			logger.Error("GetUserByUserID failed ", err.Error(), groupReq.UserID)
			continue
		}
		cp.GroupRequestDBCopyOpenIM(&node, &groupReq)
		cp.UserDBCopyOpenIMPublicUser(node.UserInfo, user)
		cp.GroupDBCopyOpenIM(ctx, node.GroupInfo, group)
		resp.GroupRequestList = append(resp.GroupRequestList, &node)
	}
	resp.CommonResp = &pbGroup.CommonResp{
		ErrCode: 0,
		ErrMsg:  "",
	}
	return resp, nil
}

func (s *groupServer) DismissGroup(ctx context.Context, req *pbGroup.DismissGroupReq) (*pbGroup.DismissGroupResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if !token_verify.IsManagerUserID(req.OpUserID) && !imdb.IsGroupOwnerAdmin(ctx, req.GroupID, req.OpUserID) {
		logger.Error("verify failed ", req.OpUserID, req.GroupID)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupID); err != nil {
		logger.Error(err.Error(), req.GroupID)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	if err := s.DelGroupAndUserCache(ctx, req.OperationID, req.GroupID, nil); err != nil {
		logger.Error(err.Error(), req.GroupID)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	err := imdb.OperateGroupStatus(ctx, req.GroupID, constant.GroupStatusDismissed)
	if err != nil {
		logger.Error("OperateGroupStatus failed ", req.GroupID, constant.GroupStatusDismissed)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	groupInfo, err := imdb.GetGroupInfoByGroupID(ctx, req.GroupID)
	if err != nil {
		logger.Error(err.Error(), req.GroupID)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.GroupType != constant.SuperGroup {
		memberList, err := imdb.GetGroupMemberListByGroupID(ctx, req.GroupID)
		if err != nil {
			logger.Error("GetGroupMemberListByGroupID failed,", err.Error(), req.GroupID)
		}
		//modify quitter conversation info
		var reqPb pbUser.SetConversationReq
		var c pbConversation.Conversation
		for _, v := range memberList {
			reqPb.OperationID = req.OperationID
			c.OwnerUserID = v.UserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsNotInGroup = true
			reqPb.Conversation = &c
			respPb, err := s.userClient.SetConversation(ctx, &reqPb)
			if err != nil {
				logger.Error("SetConversation rpc failed, ", reqPb.String(), err.Error(), v.UserID)
			} else {
				logger.Debug("SetConversation success", respPb.String(), v.UserID)
			}
		}
		err = imdb.DeleteGroupMemberByGroupID(req.GroupID)
		if err != nil {
			logger.Error("DeleteGroupMemberByGroupID failed ", req.GroupID)
			return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		chat.GroupDismissedNotification(ctx, req)
	} else {
		err = db.DB.DeleteSuperGroup(ctx, req.GroupID)
		if err != nil {
			logger.Error("DeleteGroupMemberByGroupID failed ", req.GroupID)
			return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}

	return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

//  rpc MuteGroupMember(MuteGroupMemberReq) returns(MuteGroupMemberResp);
//  rpc CancelMuteGroupMember(CancelMuteGroupMemberReq) returns(CancelMuteGroupMemberResp);
//  rpc MuteGroup(MuteGroupReq) returns(MuteGroupResp);
//  rpc CancelMuteGroup(CancelMuteGroupReq) returns(CancelMuteGroupResp);

func (s *groupServer) MuteGroupMember(ctx context.Context, req *pbGroup.MuteGroupMemberReq) (*pbGroup.MuteGroupMemberResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	opFlag, err := s.getGroupUserLevel(ctx, req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		logger.Error(errMsg)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		logger.Error(errMsg)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	mutedInfo, err := rocksCache.GetGroupMemberInfoFromCache(ctx, req.GroupID, req.UserID)
	if err != nil {
		errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.UserID + err.Error()
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.UserID
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.UserID
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.UserID); err != nil {
		logger.Error(err.Error(), req.GroupID, req.UserID)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	groupMemberInfo := db.GroupMember{GroupID: req.GroupID, UserID: req.UserID}
	groupMemberInfo.MuteEndTime = time.Unix(int64(time.Now().Second())+int64(req.MutedSeconds), time.Now().UnixNano())
	err = imdb.UpdateGroupMemberInfo(ctx, groupMemberInfo)
	if err != nil {
		logger.Error("UpdateGroupMemberInfo failed ", err.Error(), groupMemberInfo)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	chat.GroupMemberMutedNotification(ctx, req.OperationID, req.OpUserID, req.GroupID, req.UserID, req.MutedSeconds)

	return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) CancelMuteGroupMember(ctx context.Context, req *pbGroup.CancelMuteGroupMemberReq) (*pbGroup.CancelMuteGroupMemberResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	opFlag, err := s.getGroupUserLevel(ctx, req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		logger.Error(errMsg)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		logger.Error(errMsg)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(ctx, req.GroupID, req.UserID)
	if err != nil {
		errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.UserID + err.Error()
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.UserID
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.UserID
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.UserID); err != nil {
		logger.Error(err.Error(), req.GroupID)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	groupMemberInfo := db.GroupMember{GroupID: req.GroupID, UserID: req.UserID}
	groupMemberInfo.MuteEndTime = time.Unix(0, 0)
	err = imdb.UpdateGroupMemberInfo(ctx, groupMemberInfo)
	if err != nil {
		logger.Error("UpdateGroupMemberInfo failed ", err.Error(), groupMemberInfo)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	chat.GroupMemberCancelMutedNotification(ctx, req.OperationID, req.OpUserID, req.GroupID, req.UserID)

	return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) MuteGroup(ctx context.Context, req *pbGroup.MuteGroupReq) (*pbGroup.MuteGroupResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	opFlag, err := s.getGroupUserLevel(ctx, req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		logger.Error(errMsg)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		logger.Error(errMsg)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	//mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.UserID)
	//if err != nil {
	//	errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupID); err != nil {
		logger.Error(err.Error(), req.GroupID)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	err = imdb.OperateGroupStatus(ctx, req.GroupID, constant.GroupStatusMuted)
	if err != nil {
		logger.Error("OperateGroupStatus failed ", err.Error(), req.GroupID, constant.GroupStatusMuted)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	chat.GroupMutedNotification(ctx, req.OperationID, req.OpUserID, req.GroupID)

	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) CancelMuteGroup(ctx context.Context, req *pbGroup.CancelMuteGroupReq) (*pbGroup.CancelMuteGroupResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	opFlag, err := s.getGroupUserLevel(ctx, req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		logger.Error(errMsg)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		logger.Error(errMsg)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	//mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.)
	//if err != nil {
	//	errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	logger.Debug("UpdateGroupInfoDefaultZero ", req.GroupID, map[string]interface{}{"status": constant.GroupOk})
	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupID); err != nil {
		logger.Error(err.Error(), req.GroupID)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	err = imdb.UpdateGroupInfoDefaultZero(ctx, req.GroupID, map[string]interface{}{"status": constant.GroupOk})
	if err != nil {
		logger.Error("UpdateGroupInfoDefaultZero failed ", err.Error(), req.GroupID)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	chat.GroupCancelMutedNotification(ctx, req.OperationID, req.OpUserID, req.GroupID)

	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) SetGroupMemberNickname(ctx context.Context, req *pbGroup.SetGroupMemberNicknameReq) (*pbGroup.SetGroupMemberNicknameResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " verify failed " + req.OpUserID + req.GroupID
		logger.Error(errMsg)
		return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	cbReq := &pbGroup.SetGroupMemberInfoReq{
		GroupID:     req.GroupID,
		UserID:      req.UserID,
		OperationID: req.OperationID,
		OpUserID:    req.OpUserID,
		Nickname:    &wrapperspb.StringValue{Value: req.Nickname},
	}
	callbackResp := CallbackBeforeSetGroupMemberInfo(ctx, cbReq)
	if callbackResp.ErrCode != 0 {
		logger.Error("CallbackBeforeMemberJoinGroup resp: ", callbackResp)
	}
	if callbackResp.ActionCode != constant.ActionAllow {
		if callbackResp.ErrCode == 0 {
			callbackResp.ErrCode = 201
		}
		logger.Debug("CallbackBeforeMemberJoinGroup result", "end rpc and return", callbackResp)
		return &pbGroup.SetGroupMemberNicknameResp{
			CommonResp: &pbGroup.CommonResp{
				ErrCode: int32(callbackResp.ErrCode),
				ErrMsg:  callbackResp.ErrMsg,
			},
		}, nil
	}

	nickName := cbReq.Nickname.Value
	groupMemberInfo := db.GroupMember{}
	groupMemberInfo.UserID = req.UserID
	groupMemberInfo.GroupID = req.GroupID
	if nickName == "" {
		userNickname, err := imdb.GetUserNameByUserID(ctx, groupMemberInfo.UserID)
		if err != nil {
			errMsg := req.OperationID + " GetUserNameByUserID failed " + err.Error()
			logger.Error(errMsg)
			return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		groupMemberInfo.Nickname = userNickname
	} else {
		groupMemberInfo.Nickname = nickName
	}

	if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.UserID); err != nil {
		logger.Error(err.Error(), req.GroupID)
		return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	err := imdb.UpdateGroupMemberInfo(ctx, groupMemberInfo)
	if err != nil {
		errMsg := req.OperationID + " UpdateGroupMemberInfo failed " + err.Error()
		logger.Error(errMsg)
		return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	chat.GroupMemberInfoSetNotification(ctx, req.OperationID, req.OpUserID, req.GroupID, req.UserID)

	return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) SetGroupMemberInfo(ctx context.Context, req *pbGroup.SetGroupMemberInfoReq) (resp *pbGroup.SetGroupMemberInfoResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbGroup.SetGroupMemberInfoResp{CommonResp: &pbGroup.CommonResp{}}
	if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.UserID); err != nil {
		logger.Error(err.Error(), req.GroupID)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	callbackResp := CallbackBeforeSetGroupMemberInfo(ctx, req)
	if callbackResp.ErrCode != 0 {
		logger.Error("CallbackBeforeMemberJoinGroup resp: ", callbackResp)
	}
	if callbackResp.ActionCode != constant.ActionAllow {
		if callbackResp.ErrCode == 0 {
			callbackResp.ErrCode = 201
		}
		logger.Debug("CallbackBeforeMemberJoinGroup result", "end rpc and return", callbackResp)
		return &pbGroup.SetGroupMemberInfoResp{
			CommonResp: &pbGroup.CommonResp{
				ErrCode: int32(callbackResp.ErrCode),
				ErrMsg:  callbackResp.ErrMsg,
			},
		}, nil
	}

	groupMember := db.GroupMember{
		GroupID: req.GroupID,
		UserID:  req.UserID,
	}
	m := make(map[string]interface{})
	if req.RoleLevel != nil {
		m["role_level"] = req.RoleLevel.Value
	}
	if req.FaceURL != nil {
		m["user_group_face_url"] = req.FaceURL.Value
	}
	if req.Nickname != nil {
		m["nickname"] = req.Nickname.Value
	}
	if req.Ex != nil {
		m["ex"] = req.Ex.Value
	} else {
		m["ex"] = nil
	}
	err = imdb.UpdateGroupMemberInfoByMap(ctx, groupMember, m)
	if err != nil {
		logger.Error("SetGroupMemberInfo failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + ":" + err.Error()
		return resp, nil
	}
	if req.RoleLevel != nil {
		switch req.RoleLevel.Value {
		case constant.GroupOrdinaryUsers:
			//msg.GroupMemberRoleLevelChangeNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, constant.GroupMemberSetToOrdinaryUserNotification)
			chat.GroupMemberInfoSetNotification(ctx, req.OperationID, req.OpUserID, req.GroupID, req.UserID)
		case constant.GroupAdmin, constant.GroupOwner:
			//msg.GroupMemberRoleLevelChangeNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, constant.GroupMemberSetToAdminNotification)
			chat.GroupMemberInfoSetNotification(ctx, req.OperationID, req.OpUserID, req.GroupID, req.UserID)
		}
	} else {
		chat.GroupMemberInfoSetNotification(ctx, req.OperationID, req.OpUserID, req.GroupID, req.UserID)
	}

	return resp, nil
}

func (s *groupServer) GetGroupAbstractInfo(ctx context.Context, req *pbGroup.GetGroupAbstractInfoReq) (*pbGroup.GetGroupAbstractInfoResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbGroup.GetGroupAbstractInfoResp{CommonResp: &pbGroup.CommonResp{}}
	hashCode, err := rocksCache.GetGroupMemberListHashFromCache(ctx, req.GroupID)
	if err != nil {
		logger.Error("GetGroupMemberListHashFromCache failed", req.GroupID, err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	resp.GroupMemberListHash = hashCode
	num, err := rocksCache.GetGroupMemberNumFromCache(ctx, req.GroupID)
	if err != nil {
		logger.Error("GetGroupMemberNumByGroupID failed", req.GroupID, err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	resp.GroupMemberNumber = int32(num)

	return resp, nil
}

func (s *groupServer) DelGroupAndUserCache(ctx context.Context, operationID, groupID string, userIDList []string) error {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	if groupID != "" {
		cacheResp, err := s.cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{
			GroupID:     groupID,
			OperationID: operationID,
		})
		if err != nil {
			logger.Error("DelGroupMemberIDListFromCache rpc call failed ", err.Error())
			return utils.Wrap(err, "")
		}
		if cacheResp.CommonResp.ErrCode != 0 {
			logger.Error("DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
			return errors.New(fmt.Sprintf("errMsg is %s, errCode is %d", cacheResp.CommonResp.ErrMsg, cacheResp.CommonResp.ErrCode))
		}
		err = rocksCache.DelGroupMemberListHashFromCache(ctx, groupID)
		if err != nil {
			logger.Error(groupID, err.Error())
			return utils.Wrap(err, "")
		}
		err = rocksCache.DelGroupMemberNumFromCache(ctx, groupID)
		if err != nil {
			logger.Error(err.Error(), groupID)
			return utils.Wrap(err, "")
		}
	}
	if userIDList != nil {
		for _, userID := range userIDList {
			err := rocksCache.DelJoinedGroupIDListFromCache(ctx, userID)
			if err != nil {
				logger.Error(err.Error())
				return utils.Wrap(err, "")
			}
		}
	}
	return nil
}

func (s *groupServer) GroupIsExist(ctx context.Context, req *pbGroup.GroupIsExistReq) (*pbGroup.GroupIsExistResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbGroup.GroupIsExistResp{CommonResp: &pbGroup.CommonResp{}}
	groups, err := imdb.GetGroupInfoByGroupIDList(ctx, req.GroupIDList)
	if err != nil {
		logger.Error(err.Error())
		resp.CommonResp.ErrMsg = err.Error()
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	var m = make(map[string]bool)
	for _, groupID := range req.GroupIDList {
		m[groupID] = false
		for _, group := range groups {
			if groupID == group.GroupID {
				m[groupID] = true
				break
			}
		}
	}
	resp.IsExistMap = m

	return resp, nil
}

func (s *groupServer) UserIsInGroup(ctx context.Context, req *pbGroup.UserIsInGroupReq) (*pbGroup.UserIsInGroupResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbGroup.UserIsInGroupResp{}
	groupMemberList, err := imdb.GetGroupMemberByUserIDList(ctx, req.GroupID, req.UserIDList)
	if err != nil {
		logger.Error(err)
		resp.CommonResp.ErrMsg = err.Error()
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	var m = make(map[string]bool)
	for _, userID := range req.UserIDList {
		m[userID] = false
		for _, user := range groupMemberList {
			if userID == user.UserID {
				m[userID] = true
				break
			}
		}
	}
	resp.IsExistMap = m

	return resp, nil
}
