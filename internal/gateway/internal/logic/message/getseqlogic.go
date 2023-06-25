package message

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	sdk "Open_IM/pkg/proto/sdk_ws"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSeqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSeqLogic {
	return &GetSeqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSeqLogic) GetSeq(req *types.GetSeqRequest) (*types.GetSeqResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	if ok, err := token_verify.VerifyToken(l.ctx, token, req.SendID); !ok {
		return nil, errors.BadRequest.WriteMessage(err.Error())
	}

	rpcReq := sdk.GetMaxAndMinSeqReq{
		GroupIDList: []string{},
		UserID:      req.SendID,
		OperationID: req.OperationID,
	}

	rpcResp, err := l.svcCtx.MessageClient.GetMaxAndMinSeq(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, "UserGetSeq rpc failed, ", req, err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	return &types.GetSeqResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		MsgIncr:       req.MsgIncr,
		ReqIdentifier: req.ReqIdentifier,
		Data: types.GetSeqData{
			MaxSeq: rpcResp.MaxSeq,
			MinSeq: rpcResp.MinSeq,
		},
	}, nil
}
