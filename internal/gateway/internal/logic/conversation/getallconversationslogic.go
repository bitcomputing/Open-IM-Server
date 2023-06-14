package conversation

import (
	"context"

	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllConversationsLogic {
	return &GetAllConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllConversationsLogic) GetAllConversations(req *types.GetAllConversationsRequest) (*types.GetAllConversationsResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq user.GetAllConversationsReq
	if err := utils.CopyStructFields(&rpcReq, req); err != nil {
		logger.Debug("CopyStructFields failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	rpcResp, err := l.svcCtx.UserClient.GetAllConversations(l.ctx, &rpcReq)
	if err != nil {
		logger.Error("SetConversation rpc failed, ", rpcReq.String(), err.Error())
		return nil, errors.InternalError.WriteMessage("GetAllConversationMsgOpt rpc failed, " + err.Error())
	}

	resp := new(types.GetAllConversationsResponse)
	if err := utils.CopyStructFields(&resp.Conversations, rpcResp.Conversations); err != nil {
		logger.Debug("CopyStructFields failed, ", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	resp.ErrMsg = rpcResp.CommonResp.ErrMsg
	resp.ErrCode = rpcResp.CommonResp.ErrCode

	return resp, nil
}
