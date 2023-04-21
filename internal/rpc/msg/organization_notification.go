package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	utils2 "Open_IM/pkg/common/utils"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func OrganizationNotificationToAll(ctx context.Context, opUserID string, operationID string) {
	err, userIDList := imdb.GetAllOrganizationUserID()
	if err != nil {
		log.Error(operationID, "GetAllOrganizationUserID failed ", err.Error())
		return
	}

	tips := open_im_sdk.OrganizationChangedTips{OpUser: &open_im_sdk.UserInfo{}}

	user, err := imdb.GetUserByUserID(ctx, opUserID)
	if err != nil {
		log.NewError(operationID, "GetUserByUserID failed ", err.Error(), opUserID)
		return
	}
	utils2.UserDBCopyOpenIM(tips.OpUser, user)

	for _, v := range userIDList {
		log.Debug(operationID, "OrganizationNotification", opUserID, v, constant.OrganizationChangedNotification, &tips, operationID)
		OrganizationNotification(config.Config.Manager.AppManagerUid[0], v, constant.OrganizationChangedNotification, &tips, operationID)
	}
}

func OrganizationNotification(opUserID string, recvUserID string, contentType int32, m proto.Message, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "args: ", contentType, opUserID)
	var err error
	var tips open_im_sdk.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.Error(operationID, "Marshal failed ", err.Error(), m)
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
		log.Error(operationID, "contentType failed ", contentType)
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
		log.Error(operationID, "Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(&n)
}
