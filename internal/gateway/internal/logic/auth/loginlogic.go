package auth

import (
	"context"

	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/config"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))
	if req.Secret != config.Config.Secret {
		errMsg := "params.Secret != config.Config.Secret"
		logger.Error("params.Secret != config.Config.Secret", req.Secret, config.Config.Secret)
		return nil, errors.Unauthorized.WriteMessage(errMsg)
	}

	userTokenReq := &auth.UserTokenReq{
		Platform:    req.Platform,
		FromUserID:  req.UserID,
		OpUserID:    "",
		OperationID: req.OperationID,
		LoginIp:     req.LoginIp,
	}

	reply, err := l.svcCtx.AuthClient.UserToken(l.ctx, userTokenReq)
	if err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.LoginResponse{
		CommResp: types.CommResp{
			ErrCode: reply.CommonResp.ErrCode,
			ErrMsg:  reply.CommonResp.ErrMsg,
		},
		UserToken: types.UserTokenInfo{
			UserID:      userTokenReq.FromUserID,
			Token:       reply.Token,
			ExpiredTime: reply.ExpiredTime,
		},
	}, nil
}
