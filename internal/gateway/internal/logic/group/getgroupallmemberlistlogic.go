package group

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	internalutils "Open_IM/internal/utils"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

type GetGroupAllMemberListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupAllMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupAllMemberListLogic {
	return &GetGroupAllMemberListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupAllMemberListLogic) GetGroupAllMemberList(req *types.GetGroupAllMemberListRequest) (resp *types.GetGroupAllMemberListResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	rpcReq := &group.GetGroupAllMemberReq{}
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

	maxSizeOption := grpc.MaxCallRecvMsgSize(1024 * 1024 * constant.GroupRPCRecvSize)
	rpcResp, err := l.svcCtx.GroupClient.GetGroupAllMember(l.ctx, rpcReq, maxSizeOption)
	if err != nil {
		logger.Error(rpcReq.OperationID, "GetGroupAllMember failed err", err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.GetGroupAllMemberListResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		Data: internalutils.JsonDataList(rpcResp.MemberList),
	}, nil
}
