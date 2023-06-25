package conversation

import (
	"context"

	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/log"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchSetConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchSetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchSetConversationsLogic {
	return &BatchSetConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchSetConversationsLogic) BatchSetConversations(req *types.BatchSetConversationsRequest) (*types.BatchSetConversationsResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq user.BatchSetConversationsReq
	if err := utils.CopyStructFields(&rpcReq, req); err != nil {
		logger.Debug("CopyStructFields failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	rpcResp, err := l.svcCtx.UserClient.BatchSetConversations(l.ctx, &rpcReq)
	if err != nil {
		logger.Error("SetConversation rpc failed, ", rpcReq.String(), err.Error())
		return nil, errors.InternalError.WriteMessage("GetAllConversationMsgOpt rpc failed, " + err.Error())
	}

	resp := new(types.BatchSetConversationsResponse)
	if err := utils.CopyStructFields(&resp.Data, rpcResp); err != nil {
		log.NewDebug("CopyStructFields failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	resp.CommResp.ErrCode = rpcResp.CommonResp.ErrCode
	resp.CommResp.ErrMsg = rpcResp.CommonResp.ErrMsg

	return resp, nil
}
