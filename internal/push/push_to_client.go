/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/5 14:31).
 */
package push

import (
	push "Open_IM/internal/push/providers"
	cacheclient "Open_IM/internal/rpc/cache/client"
	utils2 "Open_IM/internal/utils"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbPush "Open_IM/pkg/proto/push"
	pbRelay "Open_IM/pkg/proto/relay"
	pbRtc "Open_IM/pkg/proto/rtc"
	"Open_IM/pkg/utils"
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type OpenIMContent struct {
	SessionType int    `json:"sessionType"`
	From        string `json:"from"`
	To          string `json:"to"`
	Seq         uint32 `json:"seq"`
}

type AtContent struct {
	Text       string   `json:"text"`
	AtUserList []string `json:"atUserList"`
	IsAtSelf   bool     `json:"isAtSelf"`
}

func MsgToUser(ctx context.Context, pushMsg *pbPush.PushMsgReq) {
	logger := logx.WithContext(ctx).WithFields(logx.Field("op", pushMsg.OperationID))
	var wsResult []*pbRelay.SingelMsgToUserResultList
	isOfflinePush := utils.GetSwitchFromOptions(pushMsg.MsgData.Options, constant.IsOfflinePush)
	logger.Debugv(pushMsg)
	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), pushMsg.OperationID)

	var UIDList = []string{pushMsg.PushToUserID}
	callbackResp := callbackOnlinePush(ctx, pushMsg.OperationID, UIDList, pushMsg.MsgData)
	logger.Debug("OnlinePush callback Resp")
	if callbackResp.ErrCode != 0 {
		logger.Error("callbackOnlinePush result: ", callbackResp)
	}
	if callbackResp.ActionCode != constant.ActionAllow {
		logger.Debug("OnlinePush stop")
		return
	}

	//Online push message
	logger.Debug("len  grpc", len(grpcCons), "data", pushMsg.String())
	for _, v := range grpcCons {
		msgClient := pbRelay.NewRelayClient(v)
		reply, err := msgClient.SuperGroupOnlineBatchPushOneMsg(ctx, &pbRelay.OnlineBatchPushOneMsgReq{OperationID: pushMsg.OperationID, MsgData: pushMsg.MsgData, PushToUserIDList: []string{pushMsg.PushToUserID}})
		if err != nil {
			logger.Error("SuperGroupOnlineBatchPushOneMsg push data to client rpc err", err)
			continue
		}
		if reply != nil && reply.SinglePushResult != nil {
			wsResult = append(wsResult, reply.SinglePushResult...)
		}
	}
	logger.Info("push_result", wsResult, "sendData", pushMsg.MsgData, "isOfflinePush", isOfflinePush)

	if isOfflinePush && pushMsg.PushToUserID != pushMsg.MsgData.SendID {
		// save invitation info for offline push
		for _, v := range wsResult {
			if v.OnlinePush {
				return
			}
		}
		if pushMsg.MsgData.ContentType == constant.SignalingNotification {
			isSend, err := db.DB.HandleSignalInfo(pushMsg.OperationID, pushMsg.MsgData, pushMsg.PushToUserID)
			if err != nil {
				logger.Error(err.Error(), pushMsg.MsgData)
				return
			}
			if !isSend {
				return
			}
		}
		var title, detailContent string
		callbackResp := callbackOfflinePush(ctx, pushMsg.OperationID, UIDList, pushMsg.MsgData, &[]string{})
		logger.Debug("offline callback Resp")
		if callbackResp.ErrCode != 0 {
			logger.Error("callbackOfflinePush result: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			logger.Debug("offlinePush stop")
			return
		}
		if pushMsg.MsgData.OfflinePushInfo != nil {
			title = pushMsg.MsgData.OfflinePushInfo.Title
			detailContent = pushMsg.MsgData.OfflinePushInfo.Desc
		}

		if offlinePusher == nil {
			return
		}
		opts, err := GetOfflinePushOpts(pushMsg)
		if err != nil {
			logger.Error("GetOfflinePushOpts failed", pushMsg, err.Error())
		}
		logger.Info(UIDList, title, detailContent, "opts:", opts)
		if title == "" {
			switch pushMsg.MsgData.ContentType {
			case constant.Text:
				fallthrough
			case constant.Picture:
				fallthrough
			case constant.Voice:
				fallthrough
			case constant.Video:
				fallthrough
			case constant.File:
				title = constant.ContentType2PushContent[int64(pushMsg.MsgData.ContentType)]
			case constant.AtText:
				a := AtContent{}
				_ = utils.JsonStringToStruct(string(pushMsg.MsgData.Content), &a)
				if utils.IsContain(pushMsg.PushToUserID, a.AtUserList) {
					title = constant.ContentType2PushContent[constant.AtText] + constant.ContentType2PushContent[constant.Common]
				} else {
					title = constant.ContentType2PushContent[constant.GroupMsg]
				}
			case constant.SignalingNotification:
				title = constant.ContentType2PushContent[constant.SignalMsg]
			default:
				title = constant.ContentType2PushContent[constant.Common]

			}
			// detailContent = title
		}
		if detailContent == "" {
			detailContent = title
		}
		pushResult, err := offlinePusher.Push(ctx, UIDList, title, detailContent, pushMsg.OperationID, opts)
		if err != nil {
			// promePkg.PromeInc(promePkg.MsgOfflinePushFailedCounter)
			offlinePushFailureCounter.Inc()
			logger.Error("offline push error", pushMsg.String(), err.Error())
		} else {
			// promePkg.PromeInc(promePkg.MsgOfflinePushSuccessCounter)
			offlinePushSuccessCounter.Inc()
			logger.Debug("offline push return result is ", pushResult, pushMsg.MsgData)
		}
	}
}

func MsgToSuperGroupUser(ctx context.Context, cacheClient cacheclient.CacheClient, pushMsg *pbPush.PushMsgReq) {
	logger := logx.WithContext(ctx).WithFields(logx.Field("op", pushMsg.OperationID))
	var wsResult []*pbRelay.SingelMsgToUserResultList
	isOfflinePush := utils.GetSwitchFromOptions(pushMsg.MsgData.Options, constant.IsOfflinePush)
	logger.Debugv(pushMsg.String())
	var pushToUserIDList []string
	if config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.Enable {
		callbackResp := callbackBeforeSuperGroupOnlinePush(ctx, pushMsg.OperationID, pushMsg.PushToUserID, pushMsg.MsgData, &pushToUserIDList)
		logger.Debug("offline callback Resp")
		if callbackResp.ErrCode != 0 {
			logger.Error("callbackOfflinePush result: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			logger.Debug("onlinePush stop")
			return
		}
		logger.Debug("callback userIDList Resp", pushToUserIDList)
	}
	if len(pushToUserIDList) == 0 {
		userIDList, err := utils2.GetGroupMemberUserIDList(ctx, cacheClient, pushMsg.MsgData.GroupID, pushMsg.OperationID)
		if err != nil {
			logger.Error("GetGroupMemberUserIDList failed ", err.Error(), pushMsg.MsgData.GroupID)
			return
		}
		pushToUserIDList = userIDList
	}

	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), pushMsg.OperationID)

	//Online push message
	logger.Debug("len  grpc", len(grpcCons), "data", pushMsg.String())
	for _, v := range grpcCons {
		msgClient := pbRelay.NewRelayClient(v)
		reply, err := msgClient.SuperGroupOnlineBatchPushOneMsg(ctx, &pbRelay.OnlineBatchPushOneMsgReq{OperationID: pushMsg.OperationID, MsgData: pushMsg.MsgData, PushToUserIDList: pushToUserIDList})
		if err != nil {
			logger.Error("push data to client rpc err", err)
			continue
		}
		if reply != nil && reply.SinglePushResult != nil {
			wsResult = append(wsResult, reply.SinglePushResult...)
		}
	}
	logger.Debug("push_result", wsResult, "sendData", pushMsg.MsgData)

	if isOfflinePush {
		var onlineSuccessUserIDList []string
		var WebAndPcBackgroundUserIDList []string
		onlineSuccessUserIDList = append(onlineSuccessUserIDList, pushMsg.MsgData.SendID)
		for _, v := range wsResult {
			if v.OnlinePush && v.UserID != pushMsg.MsgData.SendID {
				onlineSuccessUserIDList = append(onlineSuccessUserIDList, v.UserID)
			}
			if !v.OnlinePush {
				if len(v.Resp) != 0 {
					for _, singleResult := range v.Resp {
						if singleResult.ResultCode == -2 {
							if constant.PlatformIDToClass(int(singleResult.RecvPlatFormID)) == constant.TerminalPC ||
								singleResult.RecvPlatFormID == constant.WebPlatformID {
								WebAndPcBackgroundUserIDList = append(WebAndPcBackgroundUserIDList, v.UserID)
							}
						}
					}
				}

			}
		}
		onlineFailedUserIDList := utils.DifferenceString(onlineSuccessUserIDList, pushToUserIDList)
		//Use offline push messaging
		var title, detailContent string
		if len(onlineFailedUserIDList) > 0 {
			var offlinePushUserIDList []string
			var needOfflinePushUserIDList []string
			callbackResp := callbackOfflinePush(ctx, pushMsg.OperationID, onlineFailedUserIDList, pushMsg.MsgData, &offlinePushUserIDList)
			logger.Debug("offline callback Resp")
			if callbackResp.ErrCode != 0 {
				logger.Error("callbackOfflinePush result: ", callbackResp)
			}
			if callbackResp.ActionCode != constant.ActionAllow {
				logger.Debug("offlinePush stop")
				return
			}
			if pushMsg.MsgData.OfflinePushInfo != nil {
				title = pushMsg.MsgData.OfflinePushInfo.Title
				detailContent = pushMsg.MsgData.OfflinePushInfo.Desc
			}
			if len(offlinePushUserIDList) > 0 {
				needOfflinePushUserIDList = offlinePushUserIDList
			} else {
				needOfflinePushUserIDList = onlineFailedUserIDList
			}
			if pushMsg.MsgData.ContentType != constant.SignalingNotification {
				notNotificationUserIDList, err := db.DB.GetSuperGroupUserReceiveNotNotifyMessageIDList(pushMsg.MsgData.GroupID)
				if err != nil {
					logger.Error("GetSuperGroupUserReceiveNotNotifyMessageIDList failed", pushMsg.MsgData.GroupID)
				} else {
					logger.Debug(notNotificationUserIDList)
				}
				needOfflinePushUserIDList = utils.RemoveFromSlice(notNotificationUserIDList, needOfflinePushUserIDList)
				logger.Debug(needOfflinePushUserIDList)

			}
			if offlinePusher == nil {
				return
			}
			opts, err := GetOfflinePushOpts(pushMsg)
			if err != nil {
				logger.Error("GetOfflinePushOpts failed", pushMsg, err.Error())
			}
			logger.Info(needOfflinePushUserIDList, title, detailContent, "opts:", opts)
			if title == "" {
				switch pushMsg.MsgData.ContentType {
				case constant.Text:
					fallthrough
				case constant.Picture:
					fallthrough
				case constant.Voice:
					fallthrough
				case constant.Video:
					fallthrough
				case constant.File:
					title = constant.ContentType2PushContent[int64(pushMsg.MsgData.ContentType)]
				case constant.AtText:
					a := AtContent{}
					_ = utils.JsonStringToStruct(string(pushMsg.MsgData.Content), &a)
					if utils.IsContain(pushMsg.PushToUserID, a.AtUserList) {
						title = constant.ContentType2PushContent[constant.AtText] + constant.ContentType2PushContent[constant.Common]
					} else {
						title = constant.ContentType2PushContent[constant.GroupMsg]
					}
				case constant.SignalingNotification:
					title = constant.ContentType2PushContent[constant.SignalMsg]
				default:
					title = constant.ContentType2PushContent[constant.Common]

				}
				detailContent = title
			}
			pushResult, err := offlinePusher.Push(ctx, needOfflinePushUserIDList, title, detailContent, pushMsg.OperationID, opts)
			if err != nil {
				// promePkg.PromeInc(promePkg.MsgOfflinePushFailedCounter)
				offlinePushFailureCounter.Inc()
				logger.Error("offline push error", pushMsg.String(), err.Error())
			} else {
				// promePkg.PromeInc(promePkg.MsgOfflinePushSuccessCounter)
				offlinePushSuccessCounter.Inc()
				logger.Debug("offline push return result is ", pushResult, pushMsg.MsgData)
			}
			needBackgroupPushUserID := utils.IntersectString(needOfflinePushUserIDList, WebAndPcBackgroundUserIDList)
			grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), pushMsg.OperationID)
			if len(needBackgroupPushUserID) > 0 {
				//Online push message
				logger.Debug("len  grpc", len(grpcCons), "data", pushMsg.String())
				for _, v := range grpcCons {
					msgClient := pbRelay.NewRelayClient(v)
					_, err := msgClient.SuperGroupBackgroundOnlinePush(ctx, &pbRelay.OnlineBatchPushOneMsgReq{OperationID: pushMsg.OperationID, MsgData: pushMsg.MsgData,
						PushToUserIDList: needBackgroupPushUserID})
					if err != nil {
						logger.Error("push data to client rpc err", err)
						continue
					}
				}
			}

		}
	}
}

func GetOfflinePushOpts(pushMsg *pbPush.PushMsgReq) (opts push.PushOpts, err error) {
	if pushMsg.MsgData.ContentType < constant.SignalingNotificationEnd && pushMsg.MsgData.ContentType > constant.SignalingNotificationBegin {
		req := &pbRtc.SignalReq{}
		if err := proto.Unmarshal(pushMsg.MsgData.Content, req); err != nil {
			return opts, utils.Wrap(err, "")
		}
		switch req.Payload.(type) {
		case *pbRtc.SignalReq_Invite, *pbRtc.SignalReq_InviteInGroup:
			opts.Signal.ClientMsgID = pushMsg.MsgData.ClientMsgID
		}
	}
	if pushMsg.MsgData.OfflinePushInfo != nil {
		opts.IOSBadgeCount = pushMsg.MsgData.OfflinePushInfo.IOSBadgeCount
		opts.IOSPushSound = pushMsg.MsgData.OfflinePushInfo.IOSPushSound
		opts.Data = pushMsg.MsgData.OfflinePushInfo.Ex
	}

	return opts, nil
}
