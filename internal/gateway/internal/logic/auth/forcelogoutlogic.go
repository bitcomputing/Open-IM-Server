package auth

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	authproto "Open_IM/pkg/proto/auth"
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

func (l *ForceLogoutLogic) ForceLogout(req *types.ForceLogoutRequest) (resp *types.ForceLogoutResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	forceLogoutReq := &authproto.ForceLogoutReq{}
	utils.CopyStructFields(forceLogoutReq, &req)
	ok, opUserID, errInfo := token_verify.GetUserIDFromToken(token, forceLogoutReq.OperationID)
	if !ok {
		errMsg := forceLogoutReq.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(errMsg)
		return nil, errors.InternalError.WriteMessage(errMsg)
	}
	forceLogoutReq.OpUserID = opUserID

	reply, err := l.svcCtx.AuthClient.ForceLogout(l.ctx, forceLogoutReq)
	if err != nil {
		errMsg := forceLogoutReq.OperationID + " UserToken failed " + err.Error() + forceLogoutReq.String()
		logger.Error(errMsg)
		return nil, errors.InternalError.WriteMessage(errMsg)
	}

	return &types.ForceLogoutResponse{
		CommResp: types.CommResp{
			ErrCode: reply.CommonResp.ErrCode,
			ErrMsg:  reply.CommonResp.ErrMsg,
		},
	}, nil
}
