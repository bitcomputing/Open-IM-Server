package conversation

import (
	"context"

	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	errors "Open_IM/pkg/errors/api"
	convproto "Open_IM/pkg/proto/conversation"
	userproto "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetConversationLogic {
	return &SetConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetConversationLogic) SetConversation(req *types.SetConversationRequest) (resp *types.SetConversationResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq userproto.SetConversationReq
	rpcReq.Conversation = &convproto.Conversation{}
	if err := utils.CopyStructFields(&rpcReq, req); err != nil {
		logger.Debug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}
	if err = utils.CopyStructFields(rpcReq.Conversation, req.Conversation); err != nil {
		logger.Debug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	rpcResp, err := l.svcCtx.UserClient.SetConversation(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", rpcReq.String(), err.Error())
		return nil, errors.InternalError.WriteMessage("GetAllConversationMsgOpt rpc failed, " + err.Error())
	}

	resp.ErrMsg = rpcResp.CommonResp.ErrMsg
	resp.ErrCode = rpcResp.CommonResp.ErrCode

	return resp, nil
}
