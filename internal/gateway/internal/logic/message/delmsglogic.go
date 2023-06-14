package message

import (
	"context"
	"net/http"

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

func (l *DelMsgLogic) DelMsg(req *types.DelMsgRequest) (*types.DelMsgResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq sdk.DelMsgListReq
	if err := utils.CopyStructFields(&rpcReq, &req); err != nil {
		logger.Error("CopyStructFields", err.Error())
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
		return nil, errors.BadRequest.WriteMessage(err.Error())
	}
	rpcReq.OpUserID = opuid

	rpcResp, err := l.svcCtx.MessageClient.DelMsgList(l.ctx, &rpcReq)
	if err != nil {
		logger.Error("DelMsgList failed", err.Error(), rpcReq.String())
		return nil, errors.Error{
			HttpStatusCode: http.StatusOK,
			Code:           constant.ErrServer.ErrCode,
			Message:        constant.ErrServer.ErrMsg + err.Error(),
		}
	}

	return &types.DelMsgResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
	}, nil
}
