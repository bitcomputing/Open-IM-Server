package group

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	internalutils "Open_IM/internal/utils"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserReqGroupApplicationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserReqGroupApplicationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserReqGroupApplicationListLogic {
	return &GetUserReqGroupApplicationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserReqGroupApplicationListLogic) GetUserReqGroupApplicationList(req *types.GetUserReqGroupApplicationListRequest) (resp *types.GetUserReqGroupApplicationListResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	rpcReq := &group.GetUserReqApplicationListReq{}
	if err := utils.CopyStructFields(rpcReq, req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, rpcReq.OperationID)
	if !ok {
		errMsg := rpcReq.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(rpcReq.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid

	rpcResp, err := l.svcCtx.GroupClient.GetUserReqApplicationList(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, "GetGroupsInfo failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.GetUserReqGroupApplicationListResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
		Data: internalutils.JsonDataList(rpcResp.GroupRequestList),
	}, nil
}
