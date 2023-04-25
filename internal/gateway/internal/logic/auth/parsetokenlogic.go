package auth

import (
	"context"
	"strings"

	"Open_IM/internal/gateway/internal/common/header"
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

	headerValue, err := header.GetHeaderValues(l.ctx)
	if err != nil {
		errMsg := strings.Join([]string{req.OperationID, "GetHeaderValues failed", err.Error()}, " ")
		logger.Error(errMsg)
		return nil, errors.InternalError.WriteMessage(errMsg)
	}
	token := headerValue.Token

	ok, _, errInfo, expireTime := token_verify.GetUserIDFromTokenExpireTime(token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromTokenExpireTime failed " + errInfo
		logger.Error(errMsg)
		return nil, errors.ParseTokenFailed.WriteMessage(errMsg)
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
