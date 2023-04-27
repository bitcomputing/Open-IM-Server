package conversation

import (
	"context"

	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetRecvMsgOptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetRecvMsgOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetRecvMsgOptLogic {
	return &SetRecvMsgOptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetRecvMsgOptLogic) SetRecvMsgOpt(req *types.SetRecvMsgOptRequest) (resp *types.SetRecvMsgOptResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq user.SetRecvMsgOptReq
	if err := utils.CopyStructFields(&rpcReq, req); err != nil {
		logger.Debug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	rpcResp, err := l.svcCtx.UserClient.SetRecvMsgOpt(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "SetRecvMsgOpt rpc failed, ", rpcReq.String(), err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	resp.ErrMsg = rpcResp.CommonResp.ErrMsg
	resp.ErrCode = rpcResp.CommonResp.ErrCode

	return resp, nil
}
