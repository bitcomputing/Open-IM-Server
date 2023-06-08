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

type SetAppBadgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetAppBadgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetAppBadgeLogic {
	return &SetAppBadgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetAppBadgeLogic) SetAppBadge(req *types.SetAppBadgeRequest) (resp *types.SetAppBadgeResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}

	if !token_verify.CheckAccess(opuid, req.FromUserID) {
		logger.Error(req.OperationID, "CheckAccess false ", opuid, req.FromUserID)
		return nil, errors.BadRequest.WriteMessage("no permission")
	}

	err = db.DB.SetUserBadgeUnreadCountSum(req.FromUserID, int(req.AppUnreadCount))
	if err != nil {
		errMsg := req.OperationID + " " + "SetUserBadgeUnreadCountSum failed " + err.Error() + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.InternalError.WriteMessage(errMsg)
	}

	return &types.SetAppBadgeResponse{
		CommResp: types.CommResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
	}, nil
}
