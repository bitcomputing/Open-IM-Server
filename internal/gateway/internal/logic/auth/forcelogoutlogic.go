package auth

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/auth"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForceLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForceLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForceLogoutLogic {
	return &ForceLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForceLogoutLogic) ForceLogout(req *types.ForceLogoutRequest) (*types.ForceLogoutResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	forceLogoutReq := &auth.ForceLogoutReq{}
	if err := utils.CopyStructFields(forceLogoutReq, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, forceLogoutReq.OperationID)
	if !ok {
		logger.Error(errInfo)
		return nil, errors.BadRequest.WriteMessage(errInfo)
	}
	forceLogoutReq.OpUserID = opuid

	reply, err := l.svcCtx.AuthClient.ForceLogout(l.ctx, forceLogoutReq)
	if err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.ForceLogoutResponse{
		CommResp: types.CommResp{
			ErrCode: reply.CommonResp.ErrCode,
			ErrMsg:  reply.CommonResp.ErrMsg,
		},
	}, nil
}
