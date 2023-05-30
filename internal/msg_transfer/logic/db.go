package logic

import (
	"Open_IM/pkg/common/db"
	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

func saveUserChat(uid string, msg *pbMsg.MsgDataToMQ) error {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", msg.OperationID))
	time := utils.GetCurrentTimestampByMill()
	seq, err := db.DB.IncrUserSeq(uid)
	if err != nil {
		logger.Error("data insert to redis err", err.Error(), msg.String())
		return err
	}
	msg.MsgData.Seq = uint32(seq)
	pbSaveData := pbMsg.MsgDataToDB{}
	pbSaveData.MsgData = msg.MsgData
	logger.Info("IncrUserSeq cost time", utils.GetCurrentTimestampByMill()-time)
	return db.DB.SaveUserChatMongo2(uid, pbSaveData.MsgData.SendTime, &pbSaveData)
	//	return db.DB.SaveUserChatMongo2(uid, pbSaveData.MsgData.SendTime, &pbSaveData)
}

func saveUserChatList(userID string, msgList []*pbMsg.MsgDataToMQ, operationID string) (uint64, error) {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", operationID))
	logger.Info(userID, len(msgList))
	//return db.DB.BatchInsertChat(userID, msgList, operationID)
	return db.DB.BatchInsertChat2Cache(userID, msgList, operationID)
}
