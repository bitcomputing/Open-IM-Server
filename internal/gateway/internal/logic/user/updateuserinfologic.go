package user

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoRequest) (*types.UpdateUserInfoResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcReq := &user.UpdateUserInfoReq{UserInfo: &sdk.UserInfo{}}
	if err := utils.CopyStructFields(rpcReq.UserInfo, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}
	rpcReq.OperationID = req.OperationID

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, rpcReq.OperationID)
	if !ok {
		errMsg := rpcReq.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(rpcReq.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid

	rpcResp, err := l.svcCtx.UserClient.UpdateUserInfo(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, "UpdateUserInfo failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage("call  rpc server failed")
	}

	return &types.UpdateUserInfoResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
	}, nil
}
