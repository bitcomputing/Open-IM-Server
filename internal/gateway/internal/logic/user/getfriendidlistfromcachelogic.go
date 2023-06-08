package user

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/cache"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendIDListFromCacheLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendIDListFromCacheLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendIDListFromCacheLogic {
	return &GetFriendIDListFromCacheLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendIDListFromCacheLogic) GetFriendIDListFromCache(req *types.GetFriendIDListFromCacheRequest) (resp *types.GetFriendIDListFromCacheResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq cache.GetFriendIDListFromCacheReq

	rpcReq.OperationID = req.OperationID

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, userId, errInfo := token_verify.GetUserIDFromToken(token, req.OperationID)
	if !ok {
		errMsg := "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.InternalError.WriteMessage(errMsg)
	}
	rpcReq.UserID = userId

	rpcResp, err := l.svcCtx.CacheClient.GetFriendIDListFromCache(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "GetFriendIDListFromCache", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.GetFriendIDListFromCacheResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
		UserIDList: rpcResp.UserIDList,
	}, nil
}
