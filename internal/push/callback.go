package push

import (
	cbApi "Open_IM/pkg/call_back_struct"
	"Open_IM/pkg/common/callback"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	http2 "net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

func callbackOfflinePush(operationID string, userIDList []string, msg *commonPb.MsgData, offlinePushUserIDList *[]string) cbApi.CommonCallbackResp {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", operationID))
	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
	if !config.Config.Callback.CallbackOfflinePush.Enable {
		return callbackResp
	}
	req := cbApi.CallbackBeforePushReq{
		UserStatusBatchCallbackReq: cbApi.UserStatusBatchCallbackReq{
			UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackOfflinePushCommand,
				OperationID:     operationID,
				PlatformID:      msg.SenderPlatformID,
				Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
			},
			UserIDList: userIDList,
		},
		OfflinePushInfo: msg.OfflinePushInfo,
		ClientMsgID:     msg.ClientMsgID,
		SendID:          msg.SendID,
		GroupID:         msg.GroupID,
		ContentType:     msg.ContentType,
		SessionType:     msg.SessionType,
		AtUserIDList:    msg.AtUserIDList,
		Content:         callback.GetContent(msg),
	}
	resp := &cbApi.CallbackBeforePushResp{CommonCallbackResp: &callbackResp}
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackOfflinePushCommand, req, resp, config.Config.Callback.CallbackOfflinePush.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackOfflinePush.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	if resp.ErrCode == constant.CallbackHandleSuccess && resp.ActionCode == constant.ActionAllow {
		if len(resp.UserIDList) != 0 {
			*offlinePushUserIDList = resp.UserIDList
		}
		if resp.OfflinePushInfo != nil {
			msg.OfflinePushInfo = resp.OfflinePushInfo
		}
	}
	logger.Debug(offlinePushUserIDList, resp.UserIDList)
	return callbackResp
}

func callbackOnlinePush(operationID string, userIDList []string, msg *commonPb.MsgData) cbApi.CommonCallbackResp {
	// logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", operationID))
	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
	if !config.Config.Callback.CallbackOnlinePush.Enable || utils.IsContain(msg.SendID, userIDList) {
		return callbackResp
	}
	req := cbApi.CallbackBeforePushReq{
		UserStatusBatchCallbackReq: cbApi.UserStatusBatchCallbackReq{
			UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackOnlinePushCommand,
				OperationID:     operationID,
				PlatformID:      msg.SenderPlatformID,
				Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
			},
			UserIDList: userIDList,
		},
		//OfflinePushInfo: msg.OfflinePushInfo,
		ClientMsgID:  msg.ClientMsgID,
		SendID:       msg.SendID,
		GroupID:      msg.GroupID,
		ContentType:  msg.ContentType,
		SessionType:  msg.SessionType,
		AtUserIDList: msg.AtUserIDList,
		Content:      callback.GetContent(msg),
	}
	resp := &cbApi.CallbackBeforePushResp{CommonCallbackResp: &callbackResp}
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackOnlinePushCommand, req, resp, config.Config.Callback.CallbackOnlinePush.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackOnlinePush.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	// if resp.ErrCode == constant.CallbackHandleSuccess && resp.ActionCode == constant.ActionAllow {
	//if resp.OfflinePushInfo != nil {
	//	msg.OfflinePushInfo = resp.OfflinePushInfo
	//}
	// }
	return callbackResp
}

func callbackBeforeSuperGroupOnlinePush(operationID string, groupID string, msg *commonPb.MsgData, pushToUserList *[]string) cbApi.CommonCallbackResp {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", operationID))
	logger.Debug(groupID, msg.String(), pushToUserList)
	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
	if !config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.Enable {
		return callbackResp
	}
	req := cbApi.CallbackBeforeSuperGroupOnlinePushReq{
		UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
			CallbackCommand: constant.CallbackSuperGroupOnlinePushCommand,
			OperationID:     operationID,
			PlatformID:      msg.SenderPlatformID,
			Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
		},
		//OfflinePushInfo: msg.OfflinePushInfo,
		ClientMsgID:  msg.ClientMsgID,
		SendID:       msg.SendID,
		GroupID:      groupID,
		ContentType:  msg.ContentType,
		SessionType:  msg.SessionType,
		AtUserIDList: msg.AtUserIDList,
		Content:      callback.GetContent(msg),
		Seq:          msg.Seq,
	}
	resp := &cbApi.CallbackBeforeSuperGroupOnlinePushResp{CommonCallbackResp: &callbackResp}
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackSuperGroupOnlinePushCommand, req, resp, config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	if resp.ErrCode == constant.CallbackHandleSuccess && resp.ActionCode == constant.ActionAllow {
		if len(resp.UserIDList) != 0 {
			*pushToUserList = resp.UserIDList
		}
		//if resp.OfflinePushInfo != nil {
		//	msg.OfflinePushInfo = resp.OfflinePushInfo
		//}
	}
	logger.Debug(pushToUserList, resp.UserIDList)
	return callbackResp

}
