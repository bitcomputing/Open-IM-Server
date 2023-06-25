package user

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccountCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAccountCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountCheckLogic {
	return &AccountCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountCheckLogic) AccountCheck(req *types.AccountCheckRequest) (*types.AccountCheckResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcReq := &user.AccountCheckReq{}
	if err := utils.CopyStructFields(rpcReq, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

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

	rpcResp, err := l.svcCtx.UserClient.AccountCheck(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, "call AccountCheck users rpc server failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	resultList := []*types.AccountCheckResp_SingleUserStatus{}
	for _, v := range rpcResp.ResultList {
		resultList = append(resultList, &types.AccountCheckResp_SingleUserStatus{
			UserID:        v.UserID,
			AccountStatus: v.AccountStatus,
		})
	}

	return &types.AccountCheckResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
		ResultList: resultList,
	}, nil
}
