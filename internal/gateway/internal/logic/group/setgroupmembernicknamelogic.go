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

type SetGroupMemberNicknameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetGroupMemberNicknameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupMemberNicknameLogic {
	return &SetGroupMemberNicknameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetGroupMemberNicknameLogic) SetGroupMemberNickname(req *types.SetGroupMemberNicknameRequest) (*types.SetGroupMemberNicknameResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	rpcReq := &group.SetGroupMemberNicknameReq{}
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

	rpcResp, err := l.svcCtx.GroupClient.SetGroupMemberNickname(l.ctx, rpcReq)
	if err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.SetGroupMemberNicknameResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
	}, nil
}
