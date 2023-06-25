package msg

import (
	"Open_IM/pkg/common/constant"
	// "Open_IM/pkg/common/log"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func DeleteMessageNotification(ctx context.Context, opUserID, userID string, seqList []uint32, operationID string) {
	DeleteMessageTips := open_im_sdk.DeleteMessageTips{OpUserID: opUserID, UserID: userID, SeqList: seqList}
	MessageNotification(ctx, operationID, userID, userID, constant.DeleteMessageNotification, &DeleteMessageTips)
}

func MessageNotification(ctx context.Context, operationID, sendID, recvID string, contentType int32, m proto.Message) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	var err error
	var tips open_im_sdk.TipsComm
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
	n.ContentType = contentType
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
