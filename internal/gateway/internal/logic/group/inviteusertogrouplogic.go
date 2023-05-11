package group

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type InviteUserToGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInviteUserToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteUserToGroupLogic {
	return &InviteUserToGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InviteUserToGroupLogic) InviteUserToGroup(req *types.InviteUserToGroupRequest) (resp *types.InviteUserToGroupResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	if len(req.InvitedUserIDList) > constant.MaxNotificationNum {
		errMsg := req.OperationID + " too many members " + utils.IntToString(len(req.InvitedUserIDList))
		logger.Error(req.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}

	rpcReq := &group.InviteUserToGroupReq{}
	if err := utils.CopyStructFields(rpcReq, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(token, rpcReq.OperationID)
	if !ok {
		errMsg := rpcReq.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(rpcReq.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid

	rpcResp, err := l.svcCtx.GroupClient.InviteUserToGroup(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, "InviteUserToGroup failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	userIDResultList := []*types.UserIDResult{}
	for _, v := range rpcResp.Id2ResultList {
		userIDResultList = append(userIDResultList, &types.UserIDResult{
			UserID: v.UserID,
			Result: v.Result,
		})
	}

	return &types.InviteUserToGroupResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		UserIDResultList: userIDResultList,
	}, nil
}
