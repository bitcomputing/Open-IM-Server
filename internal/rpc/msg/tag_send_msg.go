package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	pbChat "Open_IM/pkg/proto/msg"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

func TagSendMessage(ctx context.Context, operationID string, user *db.User, recvID, content string, senderPlatformID int32) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	var req pbChat.SendMsgReq
	var msgData pbCommon.MsgData
	msgData.SendID = user.UserID
	msgData.RecvID = recvID
	msgData.ContentType = constant.Custom
	msgData.SessionType = constant.SingleChatType
	msgData.MsgFrom = constant.UserMsgType
	msgData.Content = []byte(content)
	msgData.SenderFaceURL = user.FaceURL
	msgData.SenderNickname = user.Nickname
	msgData.Options = map[string]bool{}
	msgData.Options[constant.IsSenderConversationUpdate] = false
	msgData.Options[constant.IsSenderNotificationPush] = false
	msgData.CreateTime = utils.GetCurrentTimestampByMill()
	msgData.ClientMsgID = utils.GetMsgID(user.UserID)
	msgData.SenderPlatformID = senderPlatformID
	req.MsgData = &msgData
	req.OperationID = operationID
	// etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, operationID)
	// if etcdConn == nil {
	// 	errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
	// 	logger.Error(errMsg)
	// 	return
	// }

	// client := pbChat.NewMsgClient(etcdConn)
	respPb, err := new(rpcChat).SendMsg(ctx, &req)
	if err != nil {
		logger.Error("send msg failed", err.Error())
		return
	}
	if respPb.ErrCode != 0 {
		logger.Error("send tag msg failed ", respPb)
	}
}
