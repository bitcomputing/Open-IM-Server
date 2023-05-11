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

type GetSuperGroupsInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSuperGroupsInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSuperGroupsInfoLogic {
	return &GetSuperGroupsInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSuperGroupsInfoLogic) GetSuperGroupsInfo(req *types.GetSuperGroupsInfoRequest) (resp *types.GetSuperGroupsInfoResponse, err error) {
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

	rpcReq := group.GetSuperGroupsInfoReq{OperationID: req.OperationID, OpUserID: opuid, GroupIDList: req.GroupIDList}
	rpcResp, err := l.svcCtx.GroupClient.GetSuperGroupsInfo(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, "InviteUserToGroup failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.GetSuperGroupsInfoResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
		Data: internalutils.JsonDataList(rpcResp.GroupInfoList),
	}, nil
}
