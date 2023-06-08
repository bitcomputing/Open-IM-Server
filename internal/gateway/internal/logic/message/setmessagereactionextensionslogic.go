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

type SetMessageReactionExtensionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetMessageReactionExtensionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetMessageReactionExtensionsLogic {
	return &SetMessageReactionExtensionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetMessageReactionExtensionsLogic) SetMessageReactionExtensions(req *types.SetMessageReactionExtensionsRequest) (resp *types.SetMessageReactionExtensionsResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq message.SetMessageReactionExtensionsReq
	if err := utils.CopyStructFields(&rpcReq, &req); err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo, oppid := token_verify.GetUserIDAndPlatformIDFromToken(l.ctx, token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid
	rpcReq.OpUserIDPlatformID = oppid

	rpcResp, err := l.svcCtx.MessageClient.SetMessageReactionExtensions(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "DelMsgList failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	resultKeyValue := []*types.KeyValueResp{}
	for _, v := range rpcResp.Result {
		resultKeyValue = append(resultKeyValue, &types.KeyValueResp{
			KeyValue: &types.KeyValue{
				TypeKey:          v.KeyValue.TypeKey,
				Value:            v.KeyValue.Value,
				LatestUpdateTime: v.KeyValue.LatestUpdateTime,
			},
			ErrCode: v.ErrCode,
			ErrMsg:  v.ErrMsg,
		})
	}

	return &types.SetMessageReactionExtensionsResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		Data: types.SetMessageReactionExtensionsData{
			ResultKeyValue:     resultKeyValue,
			MsgFirstModifyTime: rpcResp.MsgFirstModifyTime,
			IsReact:            rpcResp.IsReact,
		},
	}, nil
}
