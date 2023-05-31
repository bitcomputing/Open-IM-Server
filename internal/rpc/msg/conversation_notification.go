package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"

	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func SetConversationNotification(ctx context.Context, operationID, sendID, recvID string, contentType int, m proto.Message, tips open_im_sdk.TipsComm) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	var err error
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		logger.Error("Marshal failed ", err.Error(), m)
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
		logger.Error("Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(ctx, &n)
}

// SetPrivate调用
func ConversationSetPrivateNotification(ctx context.Context, operationID, sendID, recvID string, isPrivateChat bool) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))

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
	SetConversationNotification(ctx, operationID, sendID, recvID, constant.ConversationPrivateChatNotification, conversationSetPrivateTips, tips)
}

// 会话改变
func ConversationChangeNotification(ctx context.Context, operationID, userID string) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))

	ConversationChangedTips := &open_im_sdk.ConversationUpdateTips{
		UserID: userID,
	}
	var tips open_im_sdk.TipsComm
	tips.DefaultTips = config.Config.Notification.ConversationOptUpdate.DefaultTips.Tips
	SetConversationNotification(ctx, operationID, userID, userID, constant.ConversationOptChangeNotification, ConversationChangedTips, tips)
}

// 会话未读数同步
func ConversationUnreadChangeNotification(ctx context.Context, operationID, userID, conversationID string, updateUnreadCountTime int64) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))

	ConversationChangedTips := &open_im_sdk.ConversationUpdateTips{
		UserID:                userID,
		ConversationIDList:    []string{conversationID},
		UpdateUnreadCountTime: updateUnreadCountTime,
	}
	var tips open_im_sdk.TipsComm
	tips.DefaultTips = config.Config.Notification.ConversationOptUpdate.DefaultTips.Tips
	SetConversationNotification(ctx, operationID, userID, userID, constant.ConversationUnreadNotification, ConversationChangedTips, tips)
}
