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
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type SetGroupMemberInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetGroupMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupMemberInfoLogic {
	return &SetGroupMemberInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetGroupMemberInfoLogic) SetGroupMemberInfo(req *types.SetGroupMemberInfoRequest) (resp *types.SetGroupMemberInfoResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}

	rpcReq := &group.SetGroupMemberInfoReq{
		GroupID:     req.GroupID,
		UserID:      req.UserID,
		OpUserID:    opuid,
		OperationID: req.OperationID,
	}
	if req.Nickname != nil {
		rpcReq.Nickname = &wrapperspb.StringValue{Value: *req.Nickname}
	}
	if req.FaceURL != nil {
		rpcReq.FaceURL = &wrapperspb.StringValue{Value: *req.FaceURL}
	}
	if req.Ex != nil {
		rpcReq.Ex = &wrapperspb.StringValue{Value: *req.Ex}
	}
	if req.RoleLevel != nil {
		rpcReq.RoleLevel = &wrapperspb.Int32Value{Value: *req.RoleLevel}
	}

	rpcResp, err := l.svcCtx.GroupClient.SetGroupMemberInfo(l.ctx, rpcReq)
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), " failed ", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.SetGroupMemberInfoResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
	}, nil
}
