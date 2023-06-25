package message

import (
	"context"
	"net/http"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	message "Open_IM/pkg/proto/msg"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type ManagementSendMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManagementSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManagementSendMsgLogic {
	return &ManagementSendMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var validate *validator.Validate

func newUserSendMsgReq(logger logx.Logger, params *types.ManagementSendMsgRequest) *message.SendMsgReq {
	var newContent string
	options := make(map[string]bool, 5)
	var err error
	switch params.ContentType {
	case constant.Text:
		newContent = params.Content["text"].(string)
	case constant.Picture:
		fallthrough
	case constant.Custom:
		fallthrough
	case constant.Voice:
		fallthrough
	case constant.Video:
		fallthrough
	case constant.File:
		fallthrough
	case constant.CustomNotTriggerConversation:
		fallthrough
	case constant.CustomOnlineOnly:
		fallthrough
	case constant.AtText:
		fallthrough
	case constant.AdvancedRevoke:
		utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
		newContent = utils.StructToJsonString(params.Content)
	case constant.Revoke:
		utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
		newContent = params.Content["revokeMsgClientID"].(string)
	default:
	}
	if params.IsOnlineOnly {
		SetOptions(options, false)
	}
	if params.NotOfflinePush {
		utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	}
	if params.ContentType == constant.CustomOnlineOnly {
		SetOptions(options, false)
	} else if params.ContentType == constant.CustomNotTriggerConversation {
		utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	}

	pbData := message.SendMsgReq{
		OperationID: params.OperationID,
		MsgData: &sdk.MsgData{
			SendID:           params.SendID,
			GroupID:          params.GroupID,
			ClientMsgID:      utils.GetMsgID(params.SendID),
			SenderPlatformID: params.SenderPlatformID,
			SenderNickname:   params.SenderNickname,
			SenderFaceURL:    params.SenderFaceURL,
			SessionType:      params.SessionType,
			MsgFrom:          constant.SysMsgType,
			ContentType:      params.ContentType,
			Content:          []byte(newContent),
			RecvID:           params.RecvID,
			//	ForceList:        params.ForceList,
			CreateTime: utils.GetCurrentTimestampByMill(),
			Options:    options,
			OfflinePushInfo: &sdk.OfflinePushInfo{
				Title:         params.OfflinePushInfo.Title,
				Desc:          params.OfflinePushInfo.Desc,
				Ex:            params.OfflinePushInfo.Ex,
				IOSPushSound:  params.OfflinePushInfo.IOSPushSound,
				IOSBadgeCount: params.OfflinePushInfo.IOSBadgeCount,
			},
		},
	}
	if params.ContentType == constant.OANotification {
		var tips sdk.TipsComm
		tips.JsonDetail = utils.StructToJsonString(params.Content)
		pbData.MsgData.Content, err = proto.Marshal(&tips)
		if err != nil {
			logger.Error(params.OperationID, "Marshal failed ", err.Error(), tips.String())
		}
	}
	return &pbData
}

func SetOptions(options map[string]bool, value bool) {
	utils.SetSwitchFromOptions(options, constant.IsHistory, value)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, value)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, value)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, value)
}

func init() {
	validate = validator.New()
}

func (l *ManagementSendMsgLogic) ManagementSendMsg(req *types.ManagementSendMsgRequest) (*types.ManagementSendMsgResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var data any
	switch req.ContentType {
	case constant.Text:
		data = TextElem{}
	case constant.Picture:
		data = PictureElem{}
	case constant.Voice:
		data = SoundElem{}
	case constant.Video:
		data = VideoElem{}
	case constant.File:
		data = FileElem{}
	case constant.Custom:
		data = CustomElem{}
	case constant.Revoke:
		data = RevokeElem{}
	case constant.AdvancedRevoke:
		data = MessageRevoked{}
	case constant.OANotification:
		data = OANotificationElem{}
		req.SessionType = constant.NotificationChatType
	case constant.CustomNotTriggerConversation:
		data = CustomElem{}
	case constant.CustomOnlineOnly:
		data = CustomElem{}
	case constant.AtText:
		data = AtElem{}
	//case constant.HasReadReceipt:
	//case constant.Typing:
	//case constant.Quote:
	default:
		logger.Error("contentType err")
		return nil, errors.Error{
			HttpStatusCode: http.StatusOK,
			Code:           404,
			Message:        "contentType err",
		}
	}

	if err := mapstructure.WeakDecode(req.Content, &data); err != nil {
		logger.Error("content to Data struct  err", err.Error())
		return nil, errors.Error{
			HttpStatusCode: http.StatusOK,
			Code:           401,
			Message:        err.Error(),
		}
	} else if err := validate.Struct(data); err != nil {
		logger.Error("data args validate  err", err.Error())
		return nil, errors.Error{
			HttpStatusCode: http.StatusOK,
			Code:           403,
			Message:        err.Error(),
		}
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	claims, err := token_verify.ParseToken(l.ctx, token, req.OperationID)
	if err != nil {
		logger.Error(req.OperationID, "parse token failed", err.Error(), token)
		return nil, errors.Error{
			HttpStatusCode: http.StatusOK,
			Code:           400,
			Message:        "parse token failed",
		}
	}

	if !utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
		logger.Error(req.OperationID, "not authorized", token)
		return nil, errors.Error{
			HttpStatusCode: http.StatusOK,
			Code:           400,
			Message:        "not authorized",
		}
	}

	switch req.SessionType {
	case constant.SingleChatType:
		if len(req.RecvID) == 0 {
			logger.Error(req.OperationID, "recvID is a null string")
			return nil, errors.Error{
				HttpStatusCode: http.StatusOK,
				Code:           405,
				Message:        "recvID is a null string",
			}
		}
	case constant.GroupChatType, constant.SuperGroupChatType:
		if len(req.GroupID) == 0 {
			logger.Error(req.OperationID, "groupID is a null string")
			return nil, errors.Error{
				HttpStatusCode: http.StatusOK,
				Code:           405,
				Message:        "groupID is a null string",
			}
		}
	}

	pbData := newUserSendMsgReq(logger, req)

	var status int32
	rpcResp, err := l.svcCtx.MessageClient.SendMsg(l.ctx, pbData)
	if err != nil || (rpcResp != nil && rpcResp.ErrCode != 0) {
		status = constant.MsgSendFailed
	} else {
		status = constant.MsgSendSuccessed
	}

	respSetSendMsgStatus, err2 := l.svcCtx.MessageClient.SetSendMsgStatus(l.ctx, &message.SetSendMsgStatusReq{OperationID: req.OperationID, Status: status})
	if err2 != nil {
		logger.Error(err2.Error())
	}
	if respSetSendMsgStatus != nil && respSetSendMsgStatus.ErrCode != 0 {
		logger.Error(respSetSendMsgStatus.ErrCode, respSetSendMsgStatus.ErrMsg)
	}

	return &types.ManagementSendMsgResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.ErrCode,
			ErrMsg:  rpcResp.ErrMsg,
		},
		ResultList: types.UserSendMsgResp{
			ServerMsgID: rpcResp.ServerMsgID,
			ClientMsgID: rpcResp.ClientMsgID,
			SendTime:    rpcResp.SendTime,
			Ex:          "",
		},
	}, nil
}

type (
	TextElem struct {
		Text string `mapstructure:"text" validate:"required"`
	}

	PictureElem struct {
		SourcePath      string          `mapstructure:"sourcePath"`
		SourcePicture   PictureBaseInfo `mapstructure:"sourcePicture"`
		BigPicture      PictureBaseInfo `mapstructure:"bigPicture" `
		SnapshotPicture PictureBaseInfo `mapstructure:"snapshotPicture"`
	}

	PictureBaseInfo struct {
		UUID   string `mapstructure:"uuid"`
		Type   string `mapstructure:"type" `
		Size   int64  `mapstructure:"size" `
		Width  int32  `mapstructure:"width" `
		Height int32  `mapstructure:"height"`
		Url    string `mapstructure:"url" `
	}

	SoundElem struct {
		UUID      string `mapstructure:"uuid"`
		SoundPath string `mapstructure:"soundPath"`
		SourceURL string `mapstructure:"sourceUrl"`
		DataSize  int64  `mapstructure:"dataSize"`
		Duration  int64  `mapstructure:"duration"`
	}

	VideoElem struct {
		VideoPath      string `mapstructure:"videoPath"`
		VideoUUID      string `mapstructure:"videoUUID"`
		VideoURL       string `mapstructure:"videoUrl"`
		VideoType      string `mapstructure:"videoType"`
		VideoSize      int64  `mapstructure:"videoSize"`
		Duration       int64  `mapstructure:"duration"`
		SnapshotPath   string `mapstructure:"snapshotPath"`
		SnapshotUUID   string `mapstructure:"snapshotUUID"`
		SnapshotSize   int64  `mapstructure:"snapshotSize"`
		SnapshotURL    string `mapstructure:"snapshotUrl"`
		SnapshotWidth  int32  `mapstructure:"snapshotWidth"`
		SnapshotHeight int32  `mapstructure:"snapshotHeight"`
	}

	FileElem struct {
		FilePath  string `mapstructure:"filePath"`
		UUID      string `mapstructure:"uuid"`
		SourceURL string `mapstructure:"sourceUrl"`
		FileName  string `mapstructure:"fileName"`
		FileSize  int64  `mapstructure:"fileSize"`
	}

	CustomElem struct {
		Data        string `mapstructure:"data" validate:"required"`
		Description string `mapstructure:"description"`
		Extension   string `mapstructure:"extension"`
	}

	RevokeElem struct {
		RevokeMsgClientID string `mapstructure:"revokeMsgClientID" validate:"required"`
	}

	MessageRevoked struct {
		RevokerID       string `mapstructure:"revokerID" json:"revokerID" validate:"required"`
		RevokerRole     int32  `mapstructure:"revokerRole" json:"revokerRole" validate:"required"`
		ClientMsgID     string `mapstructure:"clientMsgID" json:"clientMsgID" validate:"required"`
		RevokerNickname string `mapstructure:"revokerNickname" json:"revokerNickname"`
		SessionType     int32  `mapstructure:"sessionType" json:"sessionType" validate:"required"`
		Seq             uint32 `mapstructure:"seq" json:"seq" validate:"required"`
	}

	OANotificationElem struct {
		NotificationName    string      `mapstructure:"notificationName" json:"notificationName" validate:"required"`
		NotificationFaceURL string      `mapstructure:"notificationFaceURL" json:"notificationFaceURL"`
		NotificationType    int32       `mapstructure:"notificationType" json:"notificationType" validate:"required"`
		Text                string      `mapstructure:"text" json:"text" validate:"required"`
		Url                 string      `mapstructure:"url" json:"url"`
		MixType             int32       `mapstructure:"mixType" json:"mixType"`
		PictureElem         PictureElem `mapstructure:"pictureElem" json:"pictureElem"`
		SoundElem           SoundElem   `mapstructure:"soundElem" json:"soundElem"`
		VideoElem           VideoElem   `mapstructure:"videoElem" json:"videoElem"`
		FileElem            FileElem    `mapstructure:"fileElem" json:"fileElem"`
		Ex                  string      `mapstructure:"ex" json:"ex"`
	}

	AtElem struct {
		Text        string   `mapstructure:"text"`
		AtUserList  []string `mapstructure:"atUserList"`
		AtUsersInfo []struct {
			AtUserID      string `json:"atUserID,omitempty"`
			GroupNickname string `json:"groupNickname,omitempty"`
		} `json:"atUsersInfo,omitempty"`
		IsAtSelf bool `mapstructure:"isAtSelf"`
	}
)
