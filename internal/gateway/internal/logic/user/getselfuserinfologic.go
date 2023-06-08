package user

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	internalutils "Open_IM/internal/utils"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelfUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSelfUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelfUserInfoLogic {
	return &GetSelfUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSelfUserInfoLogic) GetSelfUserInfo(req *types.GetSelfUserInfoRequest) (resp *types.GetSelfUserInfoResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcReq := &user.GetUserInfoReq{}
	if err := utils.CopyStructFields(rpcReq, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(token, rpcReq.OperationID)
	if !ok {
		errMsg := rpcReq.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid
	rpcReq.UserIDList = append(rpcReq.UserIDList, req.UserID)

	rpcResp, err := l.svcCtx.UserClient.GetUserInfo(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, "GetUserInfo failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage("call  rpc server failed")
	}

	if len(rpcResp.UserInfoList) == 1 {
		return &types.GetSelfUserInfoResponse{
			CommResp: types.CommResp{
				ErrCode: rpcResp.CommonResp.ErrCode,
				ErrMsg:  rpcResp.CommonResp.ErrMsg,
			},
			Data: internalutils.JsonDataOne(rpcResp.UserInfoList[0]),
		}, nil
	} else {
		return &types.GetSelfUserInfoResponse{
			CommResp: types.CommResp{
				ErrCode: constant.ErrDB.ErrCode,
				ErrMsg:  constant.ErrDB.ErrMsg,
			},
			Data: map[string]interface{}{},
		}, nil
	}
}
