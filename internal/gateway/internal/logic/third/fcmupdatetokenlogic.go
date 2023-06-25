package third

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"

	"github.com/zeromicro/go-zero/core/logx"
)

type FCMUpdateTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFCMUpdateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FCMUpdateTokenLogic {
	return &FCMUpdateTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FCMUpdateTokenLogic) FCMUpdateToken(req *types.FCMUpdateTokenRequest) (*types.FCMUpdateTokenResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, userId, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.InternalError.WriteMessage(errMsg)
	}

	if err = db.DB.SetFcmToken(userId, req.Platform, req.FcmToken, 0); err != nil {
		errMsg := req.OperationID + " " + "SetFcmToken failed " + err.Error() + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.InternalError.WriteMessage(errMsg)
	}

	return &types.FCMUpdateTokenResponse{
		CommResp: types.CommResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
	}, nil
}
