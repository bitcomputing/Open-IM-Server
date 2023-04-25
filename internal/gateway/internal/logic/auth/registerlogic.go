package auth

import (
	"context"

	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	errors "Open_IM/pkg/errors/api"
	authproto "Open_IM/pkg/proto/auth"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))
	if req.Secret != config.Config.Secret {
		errMsg := " params.Secret != config.Config.Secret "
		logger.Error(errMsg, req.Secret, config.Config.Secret)
		return nil, errors.Unauthorized.WriteMessage(errMsg)
	}

	userRegisterReq := &authproto.UserRegisterReq{
		UserInfo: &sdk.UserInfo{},
	}
	utils.CopyStructFields(userRegisterReq.UserInfo, &req)
	userRegisterReq.OperationID = req.OperationID

	reply, err := l.svcCtx.AuthClient.UserRegister(l.ctx, userRegisterReq)
	if err != nil {
		errMsg := userRegisterReq.OperationID + " " + "UserRegister failed " + err.Error() + userRegisterReq.String()
		logger.Error(errMsg)
		return nil, errors.InternalError.WriteMessage(errMsg)
	}
	if reply.CommonResp.ErrCode != 0 {
		errMsg := userRegisterReq.OperationID + " " + " UserRegister failed " + reply.CommonResp.ErrMsg + userRegisterReq.String()
		logger.Error(errMsg)
		if reply.CommonResp.ErrCode == constant.RegisterLimit {
			return nil, errors.RegisterLimit
		} else if reply.CommonResp.ErrCode == constant.InvitationError {
			return nil, errors.InternalError
		} else {
			return nil, errors.InternalError.WriteMessage(errMsg)
		}
	}

	userTokenReq := &authproto.UserTokenReq{
		Platform:    req.Platform,
		FromUserID:  req.UserID,
		OperationID: req.OperationID,
	}
	replyToken, err := l.svcCtx.AuthClient.UserToken(l.ctx, userTokenReq)
	if err != nil {
		errMsg := req.OperationID + " " + " client.UserToken failed " + err.Error() + userTokenReq.String()
		logger.Error(errMsg)
		return nil, errors.InternalError.WriteMessage(errMsg)
	}

	return &types.RegisterResponse{
		CommResp: types.CommResp{
			ErrCode: replyToken.CommonResp.ErrCode,
			ErrMsg:  replyToken.CommonResp.ErrMsg,
		},
		UserToken: types.UserTokenInfo{
			UserID:      userRegisterReq.UserInfo.UserID,
			Token:       replyToken.Token,
			ExpiredTime: replyToken.ExpiredTime,
		},
	}, nil
}
