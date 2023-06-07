package message

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	errors "Open_IM/pkg/errors/api"
	message "Open_IM/pkg/proto/msg"
	sdk "Open_IM/pkg/proto/sdk_ws"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMsgLogic) SendMsg(req *types.SendMsgRequest) (resp *types.SendMsgResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	rpcReq := convertSendMsgRequestAPI2RPC(token, req)
	rpcResp, err := l.svcCtx.MessageClient.SendMsg(l.ctx, rpcReq)
	if err != nil {
		logger.Error(req.OperationID, "SendMsg rpc failed, ", req, err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.SendMsgResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		Data: types.SendMsgData{
			ClientMsgID: rpcResp.ClientMsgID,
			ServerMsgID: rpcResp.ServerMsgID,
			SendTime:    rpcResp.SendTime,
		},
	}, nil
}

func convertSendMsgRequestAPI2RPC(token string, req *types.SendMsgRequest) *message.SendMsgReq {
	return &message.SendMsgReq{
		Token:       token,
		OperationID: req.OperationID,
		MsgData: &sdk.MsgData{
			SendID:           req.SendID,
			RecvID:           req.RecvID,
			GroupID:          req.GroupID,
			ClientMsgID:      req.ClientMsgID,
			SenderPlatformID: req.SenderPlatformID,
			SenderNickname:   req.SenderNickName,
			SenderFaceURL:    req.SenderFaceURL,
			SessionType:      req.SessionType,
			MsgFrom:          req.MsgFrom,
			ContentType:      req.ContentType,
			Content:          req.Content,
			CreateTime:       req.CreateTime,
			Options:          req.Options,
			OfflinePushInfo: &sdk.OfflinePushInfo{
				Title:         req.OfflineInfo.Title,
				Desc:          req.OfflineInfo.Desc,
				Ex:            req.OfflineInfo.Ex,
				IOSPushSound:  req.OfflineInfo.IOSPushSound,
				IOSBadgeCount: req.OfflineInfo.IOSBadgeCount,
			},
		},
	}
}
