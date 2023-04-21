package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func SetConversationNotification(operationID, sendID, recvID string, contentType int, m proto.Message, tips open_im_sdk.TipsComm) {
	log.NewInfo(operationID, "args: ", sendID, recvID, contentType, m, tips.String())
	var err error
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.NewError(operationID, "Marshal failed ", err.Error(), m)
		return
	}
	marshaler := protojson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
	}

	buffer, _ := marshaler.Marshal(m)
	tips.JsonDetail = string(buffer)
	var n NotificationMsg
	n.SendID = sendID
	n.RecvID = recvID
	n.ContentType = int32(contentType)
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(operationID, utils.GetSelfFuncName(), "Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(&n)
}

// SetPrivate调用
func ConversationSetPrivateNotification(operationID, sendID, recvID string, isPrivateChat bool) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	conversationSetPrivateTips := &open_im_sdk.ConversationSetPrivateTips{
		RecvID:    recvID,
		SendID:    sendID,
		IsPrivate: isPrivateChat,
	}
	var tips open_im_sdk.TipsComm
	var tipsMsg string
	if isPrivateChat {
		tipsMsg = config.Config.Notification.ConversationSetPrivate.DefaultTips.OpenTips
	} else {
		tipsMsg = config.Config.Notification.ConversationSetPrivate.DefaultTips.CloseTips
	}
	tips.DefaultTips = tipsMsg
	SetConversationNotification(operationID, sendID, recvID, constant.ConversationPrivateChatNotification, conversationSetPrivateTips, tips)
}

// 会话改变
func ConversationChangeNotification(operationID, userID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	ConversationChangedTips := &open_im_sdk.ConversationUpdateTips{
		UserID: userID,
	}
	var tips open_im_sdk.TipsComm
	tips.DefaultTips = config.Config.Notification.ConversationOptUpdate.DefaultTips.Tips
	SetConversationNotification(operationID, userID, userID, constant.ConversationOptChangeNotification, ConversationChangedTips, tips)
}

// 会话未读数同步
func ConversationUnreadChangeNotification(operationID, userID, conversationID string, updateUnreadCountTime int64) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	ConversationChangedTips := &open_im_sdk.ConversationUpdateTips{
		UserID:                userID,
		ConversationIDList:    []string{conversationID},
		UpdateUnreadCountTime: updateUnreadCountTime,
	}
	var tips open_im_sdk.TipsComm
	tips.DefaultTips = config.Config.Notification.ConversationOptUpdate.DefaultTips.Tips
	SetConversationNotification(operationID, userID, userID, constant.ConversationUnreadNotification, ConversationChangedTips, tips)
}
