package message

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	message "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageReactionExtensionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMessageReactionExtensionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageReactionExtensionsLogic {
	return &DeleteMessageReactionExtensionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMessageReactionExtensionsLogic) DeleteMessageReactionExtensions(req *types.DeleteMessageReactionExtensionsRequest) (resp *types.DeleteMessageReactionExtensionsResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq message.DeleteMessageListReactionExtensionsReq
	if err := utils.CopyStructFields(&rpcReq, &req); err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid

	rpcResp, err := l.svcCtx.MessageClient.DeleteMessageReactionExtensions(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "DelMsgList failed", err.Error())
		return nil, errors.InternalError.WriteMessage(constant.ErrServer.ErrMsg + err.Error())
	}

	return &types.DeleteMessageReactionExtensionsResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
	}, nil
}
