package auth

import (
	"context"
	"net/http"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"

	"github.com/fatih/structs"
	"github.com/zeromicro/go-zero/core/logx"
)

type ParseTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewParseTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseTokenLogic {
	return &ParseTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParseTokenLogic) ParseToken(req *types.ParseTokenRequest) (resp *types.ParseTokenResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, _, errInfo, expireTime := token_verify.GetUserIDFromTokenExpireTime(l.ctx, token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromTokenExpireTime failed " + errInfo
		logger.Error(errMsg)
		return nil, errors.Error{
			HttpStatusCode: http.StatusOK,
			Code:           1001,
			Message:        errMsg,
		}
	}

	expireTimeResp := types.ExpireTime{
		ExpireTimeSeconds: uint32(expireTime),
	}

	return &types.ParseTokenResponse{
		CommResp: types.CommResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
		Data:       structs.Map(&expireTimeResp),
		ExpireTime: expireTimeResp,
	}, nil
}
