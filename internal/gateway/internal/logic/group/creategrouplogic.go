package group

import (
	"context"
	"net/http"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	internalutils "Open_IM/internal/utils"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/group"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGroupLogic) CreateGroup(req *types.CreateGroupRequest) (*types.CreateGroupResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	if len(req.MemberList) > constant.MaxNotificationNum {
		errMsg := req.OperationID + " too many members " + utils.IntToString(len(req.MemberList))
		logger.Error(req.OperationID, errMsg)
		return nil, errors.Error{
			HttpStatusCode: http.StatusOK,
			Code:           400,
			Message:        errMsg,
		}
	}

	rpcReq := &group.CreateGroupReq{GroupInfo: &sdk.GroupInfo{}}
	if err := utils.CopyStructFields(rpcReq.GroupInfo, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	for _, v := range req.MemberList {
		rpcReq.InitMemberList = append(rpcReq.InitMemberList, &group.GroupAddMemberInfo{UserID: v.UserID, RoleLevel: v.RoleLevel})
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
	rpcReq.OwnerUserID = req.OwnerUserID
	rpcReq.OperationID = req.OperationID

	rpcResp, err := l.svcCtx.GroupClient.CreateGroup(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, "CreateGroup failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.CreateGroupResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		Data: internalutils.JsonDataOne(rpcResp.GroupInfo),
	}, nil
}
