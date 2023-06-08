package message

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullMsgBySeqListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPullMsgBySeqListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullMsgBySeqListLogic {
	return &PullMsgBySeqListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PullMsgBySeqListLogic) PullMsgBySeqList(req *types.PullMsgBySeqListRequest) (resp *types.PullMsgBySeqListResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	if ok, err := token_verify.VerifyToken(l.ctx, token, req.SendID); !ok {
		if err != nil {
			logger.Error(req.OperationID, utils.GetSelfFuncName(), err.Error(), token, req.SendID)
		}
		return nil, errors.BadRequest.WriteMessage(err.Error())
	}

	rpcReq := sdk.PullMessageBySeqListReq{}
	rpcReq.UserID = req.SendID
	rpcReq.OperationID = req.OperationID
	rpcReq.SeqList = req.SeqList

	rpcResp, err := l.svcCtx.MessageClient.PullMessageBySeqList(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(rpcReq.OperationID, "PullMessageBySeqList error", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	data, err := convertMsgDataRPC2API(rpcResp.List)
	if err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.PullMsgBySeqListResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		ReqIdentifier: req.ReqIdentifier,
		Data:          data,
	}, nil
}

func convertMsgDataRPC2API(list []*sdk.MsgData) ([]types.MsgData, error) {
	data := []types.MsgData{}

	for _, v := range list {
		item := &types.MsgData{}
		if err := utils.CopyStructFields(item, v); err != nil {
			return nil, err
		}
		data = append(data, *item)
	}
	return data, nil
}
