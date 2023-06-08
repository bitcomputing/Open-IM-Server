package message

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	message "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListReactionExtensionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessageListReactionExtensionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListReactionExtensionsLogic {
	return &GetMessageListReactionExtensionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessageListReactionExtensionsLogic) GetMessageListReactionExtensions(req *types.GetMessageListReactionExtensionsRequest) (resp *types.GetMessageListReactionExtensionsResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq message.GetMessageListReactionExtensionsReq
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
		return nil, errors.InternalError.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid

	rpcResp, err := l.svcCtx.MessageClient.GetMessageListReactionExtensions(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "DelMsgList failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	data := []*types.SingleMessageExtensionResult{}
	for _, v := range rpcResp.SingleMessageResult {
		reactionExtensionList := map[string]*types.KeyValue{}
		for k, v := range v.ReactionExtensionList {
			reactionExtensionList[k] = &types.KeyValue{
				TypeKey:          v.TypeKey,
				Value:            v.Value,
				LatestUpdateTime: v.LatestUpdateTime,
			}
		}

		data = append(data, &types.SingleMessageExtensionResult{
			ErrCode:               v.ErrCode,
			ErrMsg:                v.ErrMsg,
			ReactionExtensionList: reactionExtensionList,
			ClientMsgID:           v.ClientMsgID,
		})
	}

	return &types.GetMessageListReactionExtensionsResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		Data: data,
	}, nil
}
