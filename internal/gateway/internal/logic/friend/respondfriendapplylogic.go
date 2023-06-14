package friend

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RespondFriendApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRespondFriendApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RespondFriendApplyLogic {
	return &RespondFriendApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RespondFriendApplyLogic) RespondFriendApply(req *types.RespondFriendApplyRequest) (*types.RespondFriendApplyResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcReq := &friend.AddFriendResponseReq{CommID: &friend.CommID{}}
	if err := utils.CopyStructFields(rpcReq.CommID, &req.ParamsCommFriend); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}
	rpcReq.HandleMsg = req.HandleMsg
	rpcReq.HandleResult = req.Flag

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

	if err := utils.CopyStructFields(rpcReq, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	rpcResp, err := l.svcCtx.FriendClient.AddFriendResponse(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.CommID.OperationID, "AddFriendResponse failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.RespondFriendApplyResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
	}, nil
}
