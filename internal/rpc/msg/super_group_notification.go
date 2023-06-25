package msg

import (
	"Open_IM/pkg/common/constant"
	"context"
	//sdk "Open_IM/pkg/proto/sdk_ws"
)

func SuperGroupNotification(ctx context.Context, operationID, sendID, recvID string) {
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.SuperGroupUpdateNotification,
		SessionType: constant.SingleChatType,
		OperationID: operationID,
	}

	Notification(ctx, n)
}
