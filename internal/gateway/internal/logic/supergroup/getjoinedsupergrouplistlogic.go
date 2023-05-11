package supergroup

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	internalutils "Open_IM/internal/utils"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/group"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJoinedSuperGroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJoinedSuperGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJoinedSuperGroupListLogic {
	return &GetJoinedSuperGroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJoinedSuperGroupListLogic) GetJoinedSuperGroupList(req *types.GetJoinedSuperGroupListRequest) (resp *types.GetJoinedSuperGroupListResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}

	rpcReq := group.GetJoinedSuperGroupListReq{OperationID: req.OperationID, OpUserID: opuid, UserID: req.FromUserID}
	rpcResp, err := l.svcCtx.GroupClient.GetJoinedSuperGroupList(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, "InviteUserToGroup failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.GetJoinedSuperGroupListResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
		Data: internalutils.JsonDataList(rpcResp.GroupList),
	}, nil
}
