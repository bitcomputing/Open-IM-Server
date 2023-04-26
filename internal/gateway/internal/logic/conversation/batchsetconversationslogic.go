package conversation

import (
	"context"

	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/log"
	errors "Open_IM/pkg/errors/api"
	userproto "Open_IM/pkg/proto/user"
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

func (l *BatchSetConversationsLogic) BatchSetConversations(req *types.BatchSetConversationsRequest) (resp *types.BatchSetConversationsResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq userproto.BatchSetConversationsReq
	if err := utils.CopyStructFields(&rpcReq, req); err != nil {
		logger.Debug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	rpcResp, err := l.svcCtx.UserClient.BatchSetConversations(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", rpcReq.String(), err.Error())
		return nil, errors.InternalError.WriteMessage("GetAllConversationMsgOpt rpc failed, " + err.Error())
	}

	if err := utils.CopyStructFields(&resp.Data, rpcResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	resp.CommResp.ErrCode = rpcResp.CommonResp.ErrCode
	resp.CommResp.ErrMsg = rpcResp.CommonResp.ErrMsg

	return resp, nil
}
