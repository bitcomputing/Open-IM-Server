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

type CheckFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckFriendLogic {
	return &CheckFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckFriendLogic) CheckFriend(req *types.CheckFriendRequest) (*types.CheckFriendResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcReq := &friend.IsFriendReq{CommID: &friend.CommID{}}
	if err := utils.CopyStructFields(rpcReq.CommID, &req.ParamsCommFriend); err != nil {
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

	rpcResp, err := l.svcCtx.FriendClient.IsFriend(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.CommID.OperationID, "IsFriend failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.CheckFriendResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		Response: types.IsFriend{
			Friend: rpcResp.Response,
		},
	}, nil
}
