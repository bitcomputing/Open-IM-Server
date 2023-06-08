package message

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/utils"

	message "Open_IM/pkg/proto/msg"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearMsgLogic {
	return &ClearMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearMsgLogic) ClearMsg(req *types.ClearMsgRequest) (resp *types.ClearMsgResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcReq := &message.ClearMsgReq{}
	if err := utils.CopyStructFields(rpcReq, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, rpcReq.OperationID)
	if !ok {
		errMsg := rpcReq.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(rpcReq.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid

	rpcResp, err := l.svcCtx.MessageClient.ClearMsg(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, " CleanUpMsg failed ", err.Error(), rpcReq.String(), rpcResp.ErrMsg)
		return nil, errors.InternalError.WriteMessage(rpcResp.ErrMsg)
	}

	return &types.ClearMsgResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
	}, nil
}
