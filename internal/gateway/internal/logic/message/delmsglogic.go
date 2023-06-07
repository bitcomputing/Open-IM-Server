package message

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelMsgLogic {
	return &DelMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelMsgLogic) DelMsg(req *types.DelMsgRequest) (resp *types.DelMsgResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq sdk.DelMsgListReq
	if err := utils.CopyStructFields(&rpcReq, &req); err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(err.Error())
	}
	rpcReq.OpUserID = opuid

	rpcResp, err := l.svcCtx.MessageClient.DelMsgList(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "DelMsgList failed", err.Error())
		return nil, errors.InternalError.WriteMessage(constant.ErrServer.ErrMsg + err.Error())
	}

	return &types.DelMsgResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
	}, nil
}
