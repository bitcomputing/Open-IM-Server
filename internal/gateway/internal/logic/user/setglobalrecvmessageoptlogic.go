package user

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetGlobalRecvMessageOptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetGlobalRecvMessageOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGlobalRecvMessageOptLogic {
	return &SetGlobalRecvMessageOptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetGlobalRecvMessageOptLogic) SetGlobalRecvMessageOpt(req *types.SetGlobalRecvMessageOptRequest) (resp *types.SetGlobalRecvMessageOptResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcReq := &user.SetGlobalRecvMessageOptReq{}
	if err := utils.CopyStructFields(rpcReq, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}
	rpcReq.OperationID = req.OperationID

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, userId, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, rpcReq.OperationID)
	if !ok {
		errMsg := rpcReq.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(rpcReq.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.UserID = userId

	rpcResp, err := l.svcCtx.UserClient.SetGlobalRecvMessageOpt(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, "SetGlobalRecvMessageOpt failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage("call  rpc server failed")
	}

	return &types.SetGlobalRecvMessageOptResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
	}, nil
}
