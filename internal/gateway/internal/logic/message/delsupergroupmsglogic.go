package message

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	message "Open_IM/pkg/proto/msg"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type DelSuperGroupMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelSuperGroupMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelSuperGroupMsgLogic {
	return &DelSuperGroupMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelSuperGroupMsgLogic) DelSuperGroupMsg(req *types.DelSuperGroupMsgRequest) (resp *types.DelSuperGroupMsgResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcReq := &message.DelSuperGroupMsgReq{}
	if err := utils.CopyStructFields(rpcReq, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + rpcReq.OpUserID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid

	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	pbData := message.SendMsgReq{
		OperationID: req.OperationID,
		MsgData: &sdk.MsgData{
			SendID:      req.UserID,
			RecvID:      req.UserID,
			ClientMsgID: utils.GetMsgID(req.UserID),
			SessionType: constant.SingleChatType,
			MsgFrom:     constant.SysMsgType,
			ContentType: constant.MsgDeleteNotification,
			//	ForceList:        params.ForceList,
			CreateTime: utils.GetCurrentTimestampByMill(),
			Options:    options,
		},
	}

	var tips sdk.TipsComm
	deleteMsg := types.MsgDeleteNotificationElem{
		GroupID:     req.GroupID,
		IsAllDelete: req.IsAllDelete,
		SeqList:     req.SeqList,
	}
	tips.JsonDetail = utils.StructToJsonString(deleteMsg)

	content, err := proto.Marshal(&tips)
	if err != nil {
		logger.Error(req.OperationID, "Marshal failed ", err.Error(), tips.String())
		return &types.DelSuperGroupMsgResponse{
			CommResp: types.CommResp{
				ErrCode: 400,
				ErrMsg:  err.Error(),
			},
		}, nil
	}
	pbData.MsgData.Content = content

	if req.IsAllDelete {
		rpcResp, err := l.svcCtx.MessageClient.DelSuperGroupMsg(l.ctx, rpcReq)
		if err != nil {
			logger.Error(req.OperationID, "call delete DelSuperGroupMsg rpc server failed", err.Error())
			return nil, errors.InternalError.WriteMessage(err.Error())
		}
		return &types.DelSuperGroupMsgResponse{
			CommResp: types.CommResp{
				ErrCode: rpcResp.ErrCode,
				ErrMsg:  rpcResp.ErrMsg,
			},
		}, nil
	} else {
		rpcResp, err := l.svcCtx.MessageClient.SendMsg(l.ctx, &pbData)
		if err != nil {
			logger.Error(req.OperationID, "call delete UserSendMsg rpc server failed", err.Error())
			return nil, errors.InternalError.WriteMessage(err.Error())
		}
		return &types.DelSuperGroupMsgResponse{
			CommResp: types.CommResp{
				ErrCode: rpcResp.ErrCode,
				ErrMsg:  rpcResp.ErrMsg,
			},
		}, nil
	}
}
