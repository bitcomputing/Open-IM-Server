package user

import (
	"context"
	"net/http"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	internalutils "Open_IM/internal/utils"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersPublicInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersPublicInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersPublicInfoLogic {
	return &GetUsersPublicInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersPublicInfoLogic) GetUsersPublicInfo(req *types.GetUsersPublicInfoRequest) (resp *types.GetUsersPublicInfoResponse, err error) {
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

	ok, userId, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, rpcReq.OperationID)
	if !ok {
		errMsg := rpcReq.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(rpcReq.OperationID, errMsg)
		return nil, errors.Error{
			HttpStatusCode: http.StatusOK,
			Code:           500,
			Message:        errMsg,
		}
	}
	rpcReq.OpUserID = userId

	rpcResp, err := l.svcCtx.UserClient.GetUserInfo(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, "GetUserInfo failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage("call  rpc server failed")
	}

	var publicUserInfoList []*sdk.PublicUserInfo
	for _, v := range rpcResp.UserInfoList {
		publicUserInfoList = append(publicUserInfoList,
			&sdk.PublicUserInfo{UserID: v.UserID, Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender, Ex: v.Ex})
	}

	return &types.GetUsersPublicInfoResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
		Data: internalutils.JsonDataList(publicUserInfoList),
	}, nil
}
