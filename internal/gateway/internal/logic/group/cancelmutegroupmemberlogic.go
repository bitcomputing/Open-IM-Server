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

type CancelMuteGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelMuteGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelMuteGroupMemberLogic {
	return &CancelMuteGroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelMuteGroupMemberLogic) CancelMuteGroupMember(req *types.CancelMuteGroupMemberRequest) (*types.CancelMuteGroupMemberResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	rpcReq := &group.CancelMuteGroupMemberReq{}
	if err := utils.CopyStructFields(rpcReq, &req); err != nil {
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

	rpcResp, err := l.svcCtx.GroupClient.CancelMuteGroupMember(l.ctx, rpcReq)
	if err != nil {
		logger.Error(" failed ", rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.CancelMuteGroupMemberResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
	}, nil
}
