package conversation

import (
	"context"

	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/conversation"
	"Open_IM/pkg/proto/user"
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

func (l *SetConversationLogic) SetConversation(req *types.SetConversationRequest) (*types.SetConversationResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq user.SetConversationReq
	rpcReq.Conversation = &conversation.Conversation{}
	if err := utils.CopyStructFields(&rpcReq, req); err != nil {
		logger.Debug("CopyStructFields failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}
	if err := utils.CopyStructFields(rpcReq.Conversation, req.Conversation); err != nil {
		logger.Debug("CopyStructFields failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	rpcResp, err := l.svcCtx.UserClient.SetConversation(l.ctx, &rpcReq)
	if err != nil {
		logger.Error("SetConversation rpc failed, ", rpcReq.String(), err.Error())
		return nil, errors.InternalError.WriteMessage("GetAllConversationMsgOpt rpc failed, " + err.Error())
	}

	resp := new(types.SetConversationResponse)
	resp.ErrMsg = rpcResp.CommonResp.ErrMsg
	resp.ErrCode = rpcResp.CommonResp.ErrCode

	return resp, nil
}
