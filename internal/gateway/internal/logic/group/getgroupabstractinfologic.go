package group

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupAbstractInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupAbstractInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupAbstractInfoLogic {
	return &GetGroupAbstractInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupAbstractInfoLogic) GetGroupAbstractInfo(req *types.GetGroupAbstractInfoRequest) (resp *types.GetGroupAbstractInfoResponse, err error) {
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

	rpcResp, err := l.svcCtx.GroupClient.GetGroupAbstractInfo(l.ctx, &group.GetGroupAbstractInfoReq{
		GroupID:     req.GroupID,
		OpUserID:    opuid,
		OperationID: req.OperationID,
	})
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), " failed ", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.GetGroupAbstractInfoResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
		GroupMemberNumber:   rpcResp.GroupMemberNumber,
		GroupMemberListHash: rpcResp.GroupMemberListHash,
	}, nil
}
