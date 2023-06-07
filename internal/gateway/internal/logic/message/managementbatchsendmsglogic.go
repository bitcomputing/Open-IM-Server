package message

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	message "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"

	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
)

type ManagementBatchSendMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManagementBatchSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManagementBatchSendMsgLogic {
	return &ManagementBatchSendMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManagementBatchSendMsgLogic) ManagementBatchSendMsg(req *types.ManagementBatchSendMsgRequest) (resp *types.ManagementBatchSendMsgResponse, err error) {
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
	//case constant.AtText:
	//	data = AtElem{}
	//case constant.Merger:
	//	data =
	//case constant.Card:
	//case constant.Location:
	case constant.Custom:
		data = CustomElem{}
	case constant.Revoke:
		data = RevokeElem{}
	case constant.OANotification:
		data = OANotificationElem{}
		req.SessionType = constant.NotificationChatType
	case constant.CustomNotTriggerConversation:
		data = CustomElem{}
	case constant.CustomOnlineOnly:
		data = CustomElem{}
	//case constant.HasReadReceipt:
	//case constant.Typing:
	//case constant.Quote:
	default:
		logger.Error("contentType err")
		return &types.ManagementBatchSendMsgResponse{
			CommResp: types.CommResp{
				ErrCode: 404,
				ErrMsg:  "contentType err",
			},
			Data: types.ManagementBatchSendMsgData{},
		}, nil
	}

	if err := mapstructure.WeakDecode(req.Content, &data); err != nil {
		logger.Error("content to Data struct  err", err.Error())
		return &types.ManagementBatchSendMsgResponse{
			CommResp: types.CommResp{
				ErrCode: 401,
				ErrMsg:  err.Error(),
			},
			Data: types.ManagementBatchSendMsgData{},
		}, nil
	} else if err := validate.Struct(data); err != nil {
		logger.Error("data args validate  err", err.Error())
		return &types.ManagementBatchSendMsgResponse{
			CommResp: types.CommResp{
				ErrCode: 403,
				ErrMsg:  err.Error(),
			},
			Data: types.ManagementBatchSendMsgData{},
		}, nil
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	claims, err := token_verify.ParseToken(token, req.OperationID)
	if err != nil {
		logger.Error(req.OperationID, "parse token failed", err.Error())
		return nil, errors.BadRequest.WriteMessage(err.Error())
	}
	if !utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
		return nil, errors.BadRequest.WriteMessage("not authorized")
	}

	respSetSendMsgStatus, err := l.svcCtx.MessageClient.SetSendMsgStatus(l.ctx, &message.SetSendMsgStatusReq{OperationID: req.OperationID, Status: constant.MsgIsSending})
	if err != nil {
		logger.Error(req.OperationID, "call delete UserSendMsg rpc server failed", err.Error())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}
	if respSetSendMsgStatus.ErrCode != 0 {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), "rpc failed", respSetSendMsgStatus)
		return nil, errors.InternalError.WriteMessage(respSetSendMsgStatus.ErrMsg)
	}

	managementSendMsgReq := &types.ManagementSendMsgRequest{
		ManagementSendMsg: req.ManagementSendMsg,
	}
	pbData := newUserSendMsgReq(logger, managementSendMsgReq)
	var recvList []string
	if req.IsSendAll {
		recvList, err = im_mysql_model.SelectAllUserID()
		if err != nil {
			logger.Error(req.OperationID, utils.GetSelfFuncName(), err.Error())
			return nil, errors.InternalError.WriteMessage(err.Error())
		}
	} else {
		recvList = req.RecvIDList
	}

	var msgSendFailedFlag bool
	for _, recvID := range recvList {
		pbData.MsgData.RecvID = recvID
		// log.Info(req.OperationID, "", "api ManagementSendMsg call start..., ", pbData.String())
		rpcResp, err := l.svcCtx.MessageClient.SendMsg(l.ctx, pbData)
		if err != nil {
			logger.Error(req.OperationID, "call delete UserSendMsg rpc server failed", err.Error())
			// resp.Data.FailedIDList = append(resp.Data.FailedIDList, recvID)
			msgSendFailedFlag = true
			continue
		}
		if rpcResp.ErrCode != 0 {
			logger.Error(req.OperationID, utils.GetSelfFuncName(), "rpc failed", pbData, rpcResp)
			// resp.Data.FailedIDList = append(resp.Data.FailedIDList, recvID)
			msgSendFailedFlag = true
			continue
		}
		resp.Data.ResultList = append(resp.Data.ResultList, &types.SingleReturnResult{
			ServerMsgID: rpcResp.ServerMsgID,
			ClientMsgID: rpcResp.ClientMsgID,
			SendTime:    rpcResp.SendTime,
			RecvID:      recvID,
		})
	}

	var status int32
	if msgSendFailedFlag {
		status = constant.MsgSendFailed
	} else {
		status = constant.MsgSendSuccessed
	}
	respSetSendMsgStatus, err2 := l.svcCtx.MessageClient.SetSendMsgStatus(l.ctx, &message.SetSendMsgStatusReq{OperationID: req.OperationID, Status: status})
	if err2 != nil {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), err2.Error())
	}
	if respSetSendMsgStatus != nil && resp.ErrCode != 0 {
		logger.Error(req.OperationID, utils.GetSelfFuncName(), resp.ErrCode, resp.ErrMsg)
	}

	return resp, nil
}
