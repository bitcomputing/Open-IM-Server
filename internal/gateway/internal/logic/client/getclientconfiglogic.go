package client

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetClientConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClientConfigLogic {
	return &GetClientConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetClientConfigLogic) GetClientConfig(req *types.GetClientConfigRequest) (resp *types.GetClientConfigResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	err, _ = token_verify.ParseTokenGetUserID(l.ctx, token, req.OperationID)
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.InternalError.WriteMessage(errMsg)
	}

	config, err := imdb.GetClientInitConfig()
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.GetClientConfigResponse{
		Data: types.DiscoverPageURL{
			DiscoverPageURL: config.DiscoverPageURL,
		},
	}, nil
}
