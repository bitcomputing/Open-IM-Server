package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	utils2 "Open_IM/pkg/common/utils"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func OrganizationNotificationToAll(ctx context.Context, opUserID string, operationID string) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	err, userIDList := imdb.GetAllOrganizationUserID(ctx)
	if err != nil {
		logger.Error("GetAllOrganizationUserID failed ", err.Error())
		return
	}

	tips := open_im_sdk.OrganizationChangedTips{OpUser: &open_im_sdk.UserInfo{}}

	user, err := imdb.GetUserByUserID(ctx, opUserID)
	if err != nil {
		logger.Error("GetUserByUserID failed ", err.Error(), opUserID)
		return
	}
	utils2.UserDBCopyOpenIM(tips.OpUser, user)

	for _, v := range userIDList {
		logger.Debug("OrganizationNotification", opUserID, v, constant.OrganizationChangedNotification, &tips, operationID)
		OrganizationNotification(ctx, config.Config.Manager.AppManagerUid[0], v, constant.OrganizationChangedNotification, &tips, operationID)
	}
}

func OrganizationNotification(ctx context.Context, opUserID string, recvUserID string, contentType int32, m proto.Message, operationID string) {
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

	switch contentType {
	case constant.OrganizationChangedNotification:
		tips.DefaultTips = "OrganizationChangedNotification"

	default:
		logger.Error("contentType failed ", contentType)
		return
	}

	var n NotificationMsg
	n.SendID = opUserID
	n.RecvID = recvUserID
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
