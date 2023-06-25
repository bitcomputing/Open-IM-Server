package friend

import (
	cbApi "Open_IM/pkg/call_back_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	pbFriend "Open_IM/pkg/proto/friend"
	"context"

	//"Open_IM/pkg/proto/msg"

	http2 "net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

func callbackBeforeAddFriend(ctx context.Context, req *pbFriend.AddFriendReq) cbApi.CommonCallbackResp {
	logger := logx.WithContext(ctx)
	callbackResp := cbApi.CommonCallbackResp{OperationID: req.CommID.OperationID}
	if !config.Config.Callback.CallbackBeforeAddFriend.Enable {
		return callbackResp
	}
	logger.Debugv(req)
	commonCallbackReq := &cbApi.CallbackBeforeAddFriendReq{
		CallbackCommand: constant.CallbackBeforeAddFriendCommand,
		FromUserID:      req.CommID.FromUserID,
		ToUserID:        req.CommID.ToUserID,
		ReqMsg:          req.ReqMsg,
		OperationID:     req.CommID.OperationID,
	}
	resp := &cbApi.CallbackBeforeAddFriendResp{
		CommonCallbackResp: &callbackResp,
	}
	//utils.CopyStructFields(req, msg.MsgData)
	defer logger.Debugv(commonCallbackReq)
	if err := http.CallBackPostReturnCtx(ctx, config.Config.Callback.CallbackUrl, constant.CallbackBeforeAddFriendCommand, commonCallbackReq, resp, config.Config.Callback.CallbackBeforeAddFriend.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackBeforeAddFriend.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	return callbackResp
}
