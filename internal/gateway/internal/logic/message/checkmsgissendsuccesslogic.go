package message

import (
	"context"

	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	errors "Open_IM/pkg/errors/api"
	message "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckMsgIsSendSuccessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckMsgIsSendSuccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckMsgIsSendSuccessLogic {
	return &CheckMsgIsSendSuccessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckMsgIsSendSuccessLogic) CheckMsgIsSendSuccess(req *types.CheckMsgIsSendSuccessRequest) (resp *types.CheckMsgIsSendSuccessResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcResp, err := l.svcCtx.MessageClient.GetSendMsgStatus(l.ctx, &message.GetSendMsgStatusReq{OperationID: req.OperationID})
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), err.Error())
		return nil, errors.InternalError.WriteMessage("call GetSendMsgStatus  rpc server failed")
	}

	return &types.CheckMsgIsSendSuccessResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		Status: rpcResp.Status,
	}, nil
}
