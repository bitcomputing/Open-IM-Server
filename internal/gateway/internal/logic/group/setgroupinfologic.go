package group

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/group"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type SetGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupInfoLogic {
	return &SetGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetGroupInfoLogic) SetGroupInfo(req *types.SetGroupInfoRequest) (*types.SetGroupInfoResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	rpcReq := &group.SetGroupInfoReq{GroupInfoForSet: &sdk.GroupInfoForSet{}}
	if err := utils.CopyStructFields(rpcReq.GroupInfoForSet, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	rpcReq.OperationID = req.OperationID
	argsHandle(logger, req, rpcReq)
	ok, opuid, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, rpcReq.OperationID)
	if !ok {
		errMsg := rpcReq.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(rpcReq.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid

	rpcResp, err := l.svcCtx.GroupClient.SetGroupInfo(l.ctx, rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, "SetGroupInfo failed ", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.SetGroupInfoResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
	}, nil
}

func argsHandle(logger logx.Logger, params *types.SetGroupInfoRequest, req *group.SetGroupInfoReq) {
	if params.NeedVerification != nil {
		req.GroupInfoForSet.NeedVerification = &wrapperspb.Int32Value{Value: *params.NeedVerification}
		logger.Info(req.OperationID, "NeedVerification ", req.GroupInfoForSet.NeedVerification)
	}
	if params.LookMemberInfo != nil {
		req.GroupInfoForSet.LookMemberInfo = &wrapperspb.Int32Value{Value: *params.LookMemberInfo}
		logger.Info(req.OperationID, "LookMemberInfo ", req.GroupInfoForSet.LookMemberInfo)
	}
	if params.ApplyMemberFriend != nil {
		req.GroupInfoForSet.ApplyMemberFriend = &wrapperspb.Int32Value{Value: *params.ApplyMemberFriend}
		logger.Info(req.OperationID, "ApplyMemberFriend ", req.GroupInfoForSet.ApplyMemberFriend)
	}
}
