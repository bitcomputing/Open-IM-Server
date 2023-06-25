package logic

import (
	"Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	kfk "Open_IM/pkg/common/kafka"
	pbMsg "Open_IM/pkg/proto/msg"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/zeromicro/go-zero/core/logx"

	"google.golang.org/protobuf/proto"
)

type ModifyMsgConsumerHandler struct {
	msgHandle              map[string]fcb
	modifyMsgConsumerGroup *kfk.MConsumerGroup
}

func (mmc *ModifyMsgConsumerHandler) Init() {
	mmc.msgHandle = make(map[string]fcb)
	mmc.msgHandle[config.Config.Kafka.MsgToModify.Topic] = mmc.ModifyMsg
	mmc.modifyMsgConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.MsgToModify.Topic},
		config.Config.Kafka.MsgToModify.Addr, config.Config.Kafka.ConsumerGroupID.MsgToModify)
}

func (ModifyMsgConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ModifyMsgConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (mmc *ModifyMsgConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logx.Debug("kafka get info to mysql", "ModifyMsgConsumerHandler", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			mmc.msgHandle[msg.Topic](msg, string(msg.Key), sess)
		} else {
			logx.Error("msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (mmc *ModifyMsgConsumerHandler) ModifyMsg(cMsg *sarama.ConsumerMessage, msgKey string, _ sarama.ConsumerGroupSession) {
	ctx := context.Background()
	logx.Info("msg come here ModifyMsg!!!", "msg", string(cMsg.Value), msgKey)
	msgFromMQ := pbMsg.MsgDataToModifyByMQ{}
	err := proto.Unmarshal(cMsg.Value, &msgFromMQ)
	if err != nil {
		logx.Error(msgFromMQ.TriggerID, "msg_transfer Unmarshal msg err", "msg", string(cMsg.Value), "err", err.Error())
		return
	}

	logx.Debug(msgFromMQ.TriggerID, "proto.Unmarshal MsgDataToMQ", msgFromMQ.String())
	for _, msgDataToMQ := range msgFromMQ.MessageList {
		isReactionFromCache := utils.GetSwitchFromOptions(msgDataToMQ.MsgData.Options, constant.IsReactionFromCache)
		if !isReactionFromCache {
			continue
		}
		if msgDataToMQ.MsgData.ContentType == constant.ReactionMessageModifier {
			notification := &base_info.ReactionMessageModifierNotification{}
			if err := json.Unmarshal(msgDataToMQ.MsgData.Content, notification); err != nil {
				continue
			}
			if notification.IsExternalExtensions {
				logx.Info(msgDataToMQ.OperationID, "msg:", notification, "this is external extensions")
				continue
			}

			if msgDataToMQ.MsgData.SessionType == constant.SuperGroupChatType && utils.GetSwitchFromOptions(msgDataToMQ.MsgData.Options, constant.IsHistory) {
				if msgDataToMQ.MsgData.Seq == 0 {
					logx.Error(msgDataToMQ.OperationID, "seq==0, error msg", msgDataToMQ.MsgData)
					continue
				}
				msg, err := db.DB.GetMsgBySeqIndex(notification.SourceID, notification.Seq, msgDataToMQ.OperationID)
				if (msg != nil && msg.Seq != notification.Seq) || err != nil {
					if err != nil {
						logx.Error(msgDataToMQ.OperationID, "GetMsgBySeqIndex failed", notification, err.Error())
					}
					msgs, indexes, err := db.DB.GetSuperGroupMsgBySeqListMongo(ctx, notification.SourceID, []uint32{notification.Seq}, msgDataToMQ.OperationID)
					if err != nil {
						logx.Error(msgDataToMQ.OperationID, "GetSuperGroupMsgBySeqListMongo failed", notification, err.Error())
						continue
					}
					var index int
					if len(msgs) < 1 || len(indexes) < 1 {
						logx.Error(msgDataToMQ.OperationID, "GetSuperGroupMsgBySeqListMongo failed", notification, "len<1", msgs, indexes)
						continue
					} else {
						msg = msgs[0]
						index = indexes[msg.Seq]
					}
					msg.IsReact = true
					if err := db.DB.ReplaceMsgByIndex(notification.SourceID, msg, index, msgDataToMQ.OperationID); err != nil {
						logx.Error(msgDataToMQ.OperationID, "ReplaceMsgByIndex failed", notification.SourceID, msg)
					}

				} else {
					msg.IsReact = true
					if err = db.DB.ReplaceMsgBySeq(ctx, notification.SourceID, msg, msgDataToMQ.OperationID); err != nil {
						logx.Error(msgDataToMQ.OperationID, "ReplaceMsgBySeq failed", notification.SourceID, msg)
					}
				}
			}

			if !notification.IsReact {
				// first time to modify
				var reactionExtensionList = make(map[string]db.KeyValue)
				extendMsg := db.ExtendMsg{
					ReactionExtensionList: reactionExtensionList,
					ClientMsgID:           notification.ClientMsgID,
					MsgFirstModifyTime:    notification.MsgFirstModifyTime,
				}
				for _, v := range notification.SuccessReactionExtensionList {
					reactionExtensionList[v.TypeKey] = db.KeyValue{
						TypeKey:          v.TypeKey,
						Value:            v.Value,
						LatestUpdateTime: v.LatestUpdateTime,
					}
				}
				// modify old msg
				if err := db.DB.InsertExtendMsg(notification.SourceID, notification.SessionType, &extendMsg, msgDataToMQ.OperationID); err != nil {
					logx.Error(msgDataToMQ.OperationID, "MsgFirstModify InsertExtendMsg failed", notification.SourceID, notification.SessionType, extendMsg, err.Error())
					continue
				}

			} else {
				var reactionExtensionList = make(map[string]*server_api_params.KeyValue)
				for _, v := range notification.SuccessReactionExtensionList {
					reactionExtensionList[v.TypeKey] = &server_api_params.KeyValue{
						TypeKey:          v.TypeKey,
						Value:            v.Value,
						LatestUpdateTime: v.LatestUpdateTime,
					}
				}
				// is already modify
				if err := db.DB.InsertOrUpdateReactionExtendMsgSet(ctx, notification.SourceID, notification.SessionType, notification.ClientMsgID, notification.MsgFirstModifyTime, reactionExtensionList, msgDataToMQ.OperationID); err != nil {
					logx.Error(msgDataToMQ.OperationID, "InsertOrUpdateReactionExtendMsgSet failed")
				}
			}
		} else if msgDataToMQ.MsgData.ContentType == constant.ReactionMessageDeleter {
			notification := &base_info.ReactionMessageDeleteNotification{}
			if err := json.Unmarshal(msgDataToMQ.MsgData.Content, notification); err != nil {
				continue
			}
			if err := db.DB.DeleteReactionExtendMsgSet(ctx, notification.SourceID, notification.SessionType, notification.ClientMsgID, notification.MsgFirstModifyTime, notification.SuccessReactionExtensionList, msgDataToMQ.OperationID); err != nil {
				logx.Error(msgDataToMQ.OperationID, "InsertOrUpdateReactionExtendMsgSet failed")
			}
		}
	}

}

func UnMarshallSetReactionMsgContent(content []byte) (notification *base_info.ReactionMessageModifierNotification, err error) {
	return notification, nil
}
