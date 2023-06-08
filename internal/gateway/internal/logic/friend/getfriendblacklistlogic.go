package friend

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	internalutils "Open_IM/internal/utils"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/friend"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendBlacklistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendBlacklistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendBlacklistLogic {
	return &GetFriendBlacklistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendBlacklistLogic) GetFriendBlacklist(req *types.GetFriendBlacklistRequest) (resp *types.GetFriendBlacklistResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcReq := &friend.GetBlacklistReq{CommID: &friend.CommID{}}
	if err := utils.CopyStructFields(rpcReq.CommID, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, rpcReq.CommID.OperationID)
	if !ok {
		errMsg := rpcReq.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(rpcReq.CommID.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.CommID.OpUserID = opuid

	rpcResp, err := l.svcCtx.FriendClient.GetBlacklist(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.CommID.OperationID, "GetBlacklist failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	blackUserInfoList := []*sdk.PublicUserInfo{}
	for _, v := range rpcResp.BlackUserInfoList {
		item := sdk.PublicUserInfo{}
		if err := utils.CopyStructFields(&item, v); err != nil {
			logger.Error(err)
			return nil, errors.InternalError.WriteMessage(err.Error())
		}
		blackUserInfoList = append(blackUserInfoList, &item)
	}

	return &types.GetFriendBlacklistResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		Data: internalutils.JsonDataList(blackUserInfoList),
	}, nil
}
