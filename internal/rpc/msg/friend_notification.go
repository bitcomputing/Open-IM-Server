package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	utils2 "Open_IM/pkg/common/utils"
	pbFriend "Open_IM/pkg/proto/friend"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func getFromToUserNickname(ctx context.Context, fromUserID, toUserID string) (string, string, error) {
	from, err := imdb.GetUserByUserID(ctx, fromUserID)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	to, err := imdb.GetUserByUserID(ctx, toUserID)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	return from.Nickname, to.Nickname, nil
}

func friendNotification(ctx context.Context, commID *pbFriend.CommID, contentType int32, m proto.Message) {
	log.Info(commID.OperationID, utils.GetSelfFuncName(), "args: ", commID, contentType)
	var err error
	var tips open_im_sdk.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.Error(commID.OperationID, "Marshal failed ", err.Error())
		return
	}

	marshaler := protojson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
	}
	buffer, _ := marshaler.Marshal(m)
	tips.JsonDetail = string(buffer)

	fromUserNickname, toUserNickname, err := getFromToUserNickname(ctx, commID.FromUserID, commID.ToUserID)
	if err != nil {
		log.Error(commID.OperationID, "getFromToUserNickname failed ", err.Error(), commID.FromUserID, commID.ToUserID)
		return
	}
	cn := config.Config.Notification
	switch contentType {
	case constant.FriendApplicationNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplication.DefaultTips.Tips
	case constant.FriendApplicationApprovedNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplicationApproved.DefaultTips.Tips
	case constant.FriendApplicationRejectedNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplicationRejected.DefaultTips.Tips
	case constant.FriendAddedNotification:
		tips.DefaultTips = cn.FriendAdded.DefaultTips.Tips
	case constant.FriendDeletedNotification:
		tips.DefaultTips = cn.FriendDeleted.DefaultTips.Tips + toUserNickname
	case constant.FriendRemarkSetNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendRemarkSet.DefaultTips.Tips
	case constant.BlackAddedNotification:
		tips.DefaultTips = cn.BlackAdded.DefaultTips.Tips
	case constant.BlackDeletedNotification:
		tips.DefaultTips = cn.BlackDeleted.DefaultTips.Tips + toUserNickname
	case constant.UserInfoUpdatedNotification:
		tips.DefaultTips = cn.UserInfoUpdated.DefaultTips.Tips
	case constant.FriendInfoUpdatedNotification:
		tips.DefaultTips = cn.FriendInfoUpdated.DefaultTips.Tips + toUserNickname
	default:
		log.Error(commID.OperationID, "contentType failed ", contentType)
		return
	}

	var n NotificationMsg
	n.SendID = commID.FromUserID
	n.RecvID = commID.ToUserID
	n.ContentType = contentType
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = commID.OperationID
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(commID.OperationID, "Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(&n)
}

func FriendApplicationNotification(ctx context.Context, req *pbFriend.AddFriendReq) {
	FriendApplicationTips := open_im_sdk.FriendApplicationTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	FriendApplicationTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationTips.FromToUserID.ToUserID = req.CommID.ToUserID
	friendNotification(ctx, req.CommID, constant.FriendApplicationNotification, &FriendApplicationTips)
}

func FriendApplicationApprovedNotification(ctx context.Context, req *pbFriend.AddFriendResponseReq) {
	FriendApplicationApprovedTips := open_im_sdk.FriendApplicationApprovedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	friendNotification(ctx, req.CommID, constant.FriendApplicationApprovedNotification, &FriendApplicationApprovedTips)
}

func FriendApplicationRejectedNotification(ctx context.Context, req *pbFriend.AddFriendResponseReq) {
	FriendApplicationApprovedTips := open_im_sdk.FriendApplicationApprovedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	friendNotification(ctx, req.CommID, constant.FriendApplicationRejectedNotification, &FriendApplicationApprovedTips)
}

func FriendAddedNotification(ctx context.Context, operationID, opUserID, fromUserID, toUserID string) {
	friendAddedTips := open_im_sdk.FriendAddedTips{Friend: &open_im_sdk.FriendInfo{}, OpUser: &open_im_sdk.PublicUserInfo{}}
	user, err := imdb.GetUserByUserID(ctx, opUserID)
	if err != nil {
		log.NewError(operationID, "GetUserByUserID failed ", err.Error(), opUserID)
		return
	}
	utils2.UserDBCopyOpenIMPublicUser(friendAddedTips.OpUser, user)
	friend, err := imdb.GetFriendRelationshipFromFriend(ctx, fromUserID, toUserID)
	if err != nil {
		log.NewError(operationID, "GetFriendRelationshipFromFriend failed ", err.Error(), fromUserID, toUserID)
		return
	}
	utils2.FriendDBCopyOpenIM(ctx, friendAddedTips.Friend, friend)
	commID := pbFriend.CommID{FromUserID: fromUserID, ToUserID: toUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(ctx, &commID, constant.FriendAddedNotification, &friendAddedTips)
}

func FriendDeletedNotification(ctx context.Context, req *pbFriend.DeleteFriendReq) {
	friendDeletedTips := open_im_sdk.FriendDeletedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	friendDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	friendDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	friendNotification(ctx, req.CommID, constant.FriendDeletedNotification, &friendDeletedTips)
}

func FriendRemarkSetNotification(ctx context.Context, operationID, opUserID, fromUserID, toUserID string) {
	friendInfoChangedTips := open_im_sdk.FriendInfoChangedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	friendInfoChangedTips.FromToUserID.FromUserID = fromUserID
	friendInfoChangedTips.FromToUserID.ToUserID = toUserID
	commID := pbFriend.CommID{FromUserID: fromUserID, ToUserID: toUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(ctx, &commID, constant.FriendRemarkSetNotification, &friendInfoChangedTips)
}

func BlackAddedNotification(ctx context.Context, req *pbFriend.AddBlacklistReq) {
	blackAddedTips := open_im_sdk.BlackAddedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	blackAddedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	blackAddedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	friendNotification(ctx, req.CommID, constant.BlackAddedNotification, &blackAddedTips)
}

func BlackDeletedNotification(ctx context.Context, req *pbFriend.RemoveBlacklistReq) {
	blackDeletedTips := open_im_sdk.BlackDeletedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	blackDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	blackDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	friendNotification(ctx, req.CommID, constant.BlackDeletedNotification, &blackDeletedTips)
}

// send to myself
func UserInfoUpdatedNotification(ctx context.Context, operationID, opUserID string, changedUserID string) {
	selfInfoUpdatedTips := open_im_sdk.UserInfoUpdatedTips{UserID: changedUserID}
	commID := pbFriend.CommID{FromUserID: opUserID, ToUserID: changedUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(ctx, &commID, constant.UserInfoUpdatedNotification, &selfInfoUpdatedTips)
}

func FriendInfoUpdatedNotification(ctx context.Context, operationID, changedUserID string, needNotifiedUserID string, opUserID string) {
	selfInfoUpdatedTips := open_im_sdk.UserInfoUpdatedTips{UserID: changedUserID}
	commID := pbFriend.CommID{FromUserID: opUserID, ToUserID: needNotifiedUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(ctx, &commID, constant.FriendInfoUpdatedNotification, &selfInfoUpdatedTips)
}
