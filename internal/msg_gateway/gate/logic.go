package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/msg"
	pbRtc "Open_IM/pkg/proto/rtc"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"bytes"
	"context"
	"encoding/gob"
	"runtime"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func (ws *WServer) msgParse(conn *UserConn, binaryMsg []byte) {
	b := bytes.NewBuffer(binaryMsg)
	m := Req{}
	dec := gob.NewDecoder(b)
	err := dec.Decode(&m)
	if err != nil {
		logx.Error("ws decode err: ", err.Error())
		err = conn.Close()
		if err != nil {
			logx.Error("ws close err", err.Error())
		}
		return
	}
	if err := validate.Struct(m); err != nil {
		logx.Error("ws args validate  err", err.Error())
		ws.sendErrMsg(conn, 201, err.Error(), m.ReqIdentifier, m.MsgIncr, m.OperationID)
		return
	}

	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", m.OperationID))
	if m.SendID != conn.userID {
		if err = conn.Close(); err != nil {
			logger.Errorf("close ws conn failed: user_id: %s, send_id: %s, err: %v", conn.userID, m.SendID)
			return
		}
	}
	switch m.ReqIdentifier {
	case constant.WSGetNewestSeq:
		logger.Infof("getSeqReq: send_id: %s, msg_incr: %s, req_identifier: %d", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.getSeqReq(conn, &m)
		// promePkg.PromeInc(promePkg.GetNewestSeqTotalCounter)
		ws.getNewestSeqTotal.Inc()
	case constant.WSSendMsg:
		logger.Infof("sendMsgReq: send_id: %s, msg_incr: %d, req_identifier: %d", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.sendMsgReq(conn, &m)
		// promePkg.PromeInc(promePkg.MsgRecvTotalCounter)
		ws.msgRecvTotal.Inc()
	case constant.WSSendSignalMsg:
		logger.Infof("sendSignalMsgReq: send_id: %s, msg_incr: %d, req_identifier: %d", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.sendSignalMsgReq(conn, &m)
	case constant.WSPullMsgBySeqList:
		logger.Infof("pullMsgBySeqListReq: send_id: %s, msg_incr: %d, req_identifier: %d ", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.pullMsgBySeqListReq(conn, &m)
		// promePkg.PromeInc(promePkg.PullMsgBySeqListTotalCounter)
		ws.pullMsgBySeqListTotal.Inc()
	case constant.WsLogoutMsg:
		logger.Infof("conn.Close(): send_id: %s, msg_incr: %d, req_identifier: %d", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.userLogoutReq(conn, &m)
	case constant.WsSetBackgroundStatus:
		logger.Infof("WsSetBackgroundStatus: send_id: %s, msg_incr: %d, req_identifier: %d", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.setUserDeviceBackground(conn, &m)
	default:
		logger.Errorf("ReqIdentifier failed: send_id: %s, msg_incr: %d, req_identifier: %d", m.SendID, m.MsgIncr, m.ReqIdentifier)
	}
	logger.Info("goroutine num is ", runtime.NumGoroutine())
}

func (ws *WServer) getSeqReq(conn *UserConn, m *Req) {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", m.OperationID))
	logger.Infof("Ws call success to getNewSeq: send_id: %s, msg_incr: %d, req_identifier: %d", m.SendID, m.MsgIncr, m.ReqIdentifier)
	nReply := new(sdk_ws.GetMaxAndMinSeqResp)
	isPass, errCode, errMsg, data := ws.argsValidate(m, constant.WSGetNewestSeq, m.OperationID)
	logger.Infof("argsValidate: is_pass: %v, code: %d, msg: %s", isPass, errCode, errMsg)
	if isPass {
		rpcReq := sdk_ws.GetMaxAndMinSeqReq{}
		rpcReq.GroupIDList = data.(sdk_ws.GetMaxAndMinSeqReq).GroupIDList
		rpcReq.UserID = m.SendID
		rpcReq.OperationID = m.OperationID
		logger.Debugf("Ws call success to getMaxAndMinSeq: send_id: %s, msg_incr: %d, req_identifier: %d", m.SendID, m.MsgIncr, m.ReqIdentifier, data.(sdk_ws.GetMaxAndMinSeqReq).GroupIDList)
		// if grpcConn == nil {
		// 	errMsg := rpcReq.OperationID + "getcdv3.GetDefaultConn == nil"
		// 	nReply.ErrCode = 500
		// 	nReply.ErrMsg = errMsg
		// 	logger.Error(errMsg)
		// 	ws.getSeqResp(conn, m, nReply)
		// 	return
		// }
		rpcReply, err := ws.messageClient.GetMaxAndMinSeq(context.Background(), &rpcReq)
		if err != nil {
			nReply.ErrCode = 500
			nReply.ErrMsg = err.Error()
			logger.Error("rpc call failed to GetMaxAndMinSeq: resp: ", nReply.String())
			ws.getSeqResp(conn, m, nReply)
		} else {
			logger.Info("rpc call success to getSeqReq: resp: ", rpcReply.String())
			ws.getSeqResp(conn, m, rpcReply)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		logger.Error("argsValidate failed send resp: ", nReply.String())
		ws.getSeqResp(conn, m, nReply)
	}
}

func (ws *WServer) getSeqResp(conn *UserConn, m *Req, pb *sdk_ws.GetMaxAndMinSeqResp) {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", m.OperationID))
	b, _ := proto.Marshal(pb)
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		ErrCode:       pb.GetErrCode(),
		ErrMsg:        pb.GetErrMsg(),
		OperationID:   m.OperationID,
		Data:          b,
	}
	logger.Debugf("getSeqResp come  here req: %v, send resp: msg_incr: %d, req_identifier: %d, code: %d, msg: %s", pb,
		mReply.ReqIdentifier, mReply.MsgIncr, mReply.ErrCode, mReply.ErrMsg)
	ws.sendMsg(conn, mReply)
}

func (ws *WServer) pullMsgBySeqListReq(conn *UserConn, m *Req) {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", m.OperationID))
	logger.Infof("Ws call success to pullMsgBySeqListReq start: send_id: %s, msg_incr: %d, req_identifier: %d, data: %s", m.SendID, m.MsgIncr, m.ReqIdentifier, string(m.Data))
	nReply := new(sdk_ws.PullMessageBySeqListResp)
	isPass, errCode, errMsg, data := ws.argsValidate(m, constant.WSPullMsgBySeqList, m.OperationID)
	if isPass {
		rpcReq := sdk_ws.PullMessageBySeqListReq{}
		rpcReq.SeqList = data.(sdk_ws.PullMessageBySeqListReq).SeqList
		rpcReq.UserID = m.SendID
		rpcReq.OperationID = m.OperationID
		rpcReq.GroupSeqList = data.(sdk_ws.PullMessageBySeqListReq).GroupSeqList
		logger.Infof("Ws call success to pullMsgBySeqListReq middle: send_id: %s, msg_incr: %d, req_identifier: %d, seq_list: %v", m.SendID, m.MsgIncr, m.ReqIdentifier, data.(sdk_ws.PullMessageBySeqListReq).SeqList)
		// grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, m.OperationID)
		// if grpcConn == nil {
		// 	errMsg := rpcReq.OperationID + "getcdv3.GetDefaultConn == nil"
		// 	nReply.ErrCode = 500
		// 	nReply.ErrMsg = errMsg
		// 	logger.Error(errMsg)
		// 	ws.pullMsgBySeqListResp(conn, m, nReply)
		// 	return
		// }
		// msgClient := pbChat.NewMsgClient(grpcConn)
		maxSizeOption := grpc.MaxCallRecvMsgSize(1024 * 1024 * 20)
		reply, err := ws.messageClient.PullMessageBySeqList(context.Background(), &rpcReq, maxSizeOption)
		if err != nil {
			logger.Error(err)
			nReply.ErrCode = 200
			nReply.ErrMsg = err.Error()
			ws.pullMsgBySeqListResp(conn, m, nReply)
		} else {
			logger.Infof("rpc call success to pullMsgBySeqListReq: resp: %v, msg_count: %v", reply.String(), len(reply.List))
			ws.pullMsgBySeqListResp(conn, m, reply)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		ws.pullMsgBySeqListResp(conn, m, nReply)
	}
}
func (ws *WServer) pullMsgBySeqListResp(conn *UserConn, m *Req, pb *sdk_ws.PullMessageBySeqListResp) {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", m.OperationID))
	logger.Info("pullMsgBySeqListResp come here ", pb.String())
	c, _ := proto.Marshal(pb)
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		ErrCode:       pb.GetErrCode(),
		ErrMsg:        pb.GetErrMsg(),
		OperationID:   m.OperationID,
		Data:          c,
	}
	logger.Infof("pullMsgBySeqListResp all data is: msg_incr: %d, req_identifier: %d, code: %d, msg: %s, msg_count: %d", mReply.MsgIncr, mReply.ReqIdentifier, mReply.ErrCode, mReply.ErrMsg,
		len(mReply.Data))
	ws.sendMsg(conn, mReply)
}

func (ws *WServer) userLogoutReq(conn *UserConn, m *Req) {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", m.OperationID))
	logger.Infof("Ws call success to userLogoutReq start: send_id: %s, msg_incr: %d, req_identifier: %d, data: %s", m.SendID, m.MsgIncr, m.ReqIdentifier, string(m.Data))

	// rpcReq := push.DelUserPushTokenReq{}
	// rpcReq.UserID = m.SendID
	// rpcReq.PlatformID = conn.PlatformID
	// rpcReq.OperationID = m.OperationID
	// grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImPushName, m.OperationID)
	// if grpcConn == nil {
	// 	errMsg := rpcReq.OperationID + "getcdv3.GetDefaultConn == nil"
	// 	log.NewError(rpcReq.OperationID, errMsg)
	// 	ws.userLogoutResp(conn, m)
	// 	return
	// }
	// msgClient := push.NewPushMsgServiceClient(grpcConn)
	// reply, err := msgClient.DelUserPushToken(context.Background(), &rpcReq)
	// if err != nil {
	// 	log.NewError(rpcReq.OperationID, "DelUserPushToken err", err.Error())

	// 	ws.userLogoutResp(conn, m)
	// } else {
	// 	log.NewInfo(rpcReq.OperationID, "rpc call success to DelUserPushToken", reply.String())
	// 	ws.userLogoutResp(conn, m)
	// }
	// ws.userLogoutResp(conn, m)
	if err := db.DB.DelFcmToken(context.Background(), m.SendID, int(conn.PlatformID)); err != nil {
		logger.Error("DelUserPushToken err: ", err.Error())
	}
	ws.userLogoutResp(conn, m)
}

func (ws *WServer) userLogoutResp(conn *UserConn, m *Req) {
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		OperationID:   m.OperationID,
	}
	ws.sendMsg(conn, mReply)
	_ = conn.Close()
}

func (ws *WServer) sendMsgReq(conn *UserConn, m *Req) {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", m.OperationID))
	logger.Infof("Ws call success to sendMsgReq start: send_id: %s, msg_incr: %d, req_identifier: %d", m.SendID, m.MsgIncr, m.ReqIdentifier)

	nReply := new(pbChat.SendMsgResp)
	isPass, errCode, errMsg, pData := ws.argsValidate(m, constant.WSSendMsg, m.OperationID)
	if isPass {
		data := pData.(sdk_ws.MsgData)
		pbData := pbChat.SendMsgReq{
			Token:       m.Token,
			OperationID: m.OperationID,
			MsgData:     &data,
		}
		logger.Infof("Ws call success to sendMsgReq middle: send_id: %s, msg_incr: %d, req_identifier: %d, data: %s", m.SendID, m.MsgIncr, m.ReqIdentifier, data.String())
		// etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, m.OperationID)
		// if etcdConn == nil {
		// 	errMsg := m.OperationID + "getcdv3.GetDefaultConn == nil"
		// 	nReply.ErrCode = 500
		// 	nReply.ErrMsg = errMsg
		// 	logger.Error(errMsg)
		// 	ws.sendMsgResp(conn, m, nReply)
		// 	return
		// }
		// client := pbChat.NewMsgClient(etcdConn)
		reply, err := ws.messageClient.SendMsg(context.Background(), &pbData)
		if err != nil {
			logger.Error(err)
			nReply.ErrCode = 200
			nReply.ErrMsg = err.Error()
			ws.sendMsgResp(conn, m, nReply)
		} else {
			logger.Info("rpc call success to sendMsgReq: resp: ", reply.String())
			ws.sendMsgResp(conn, m, reply)
		}

	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		ws.sendMsgResp(conn, m, nReply)
	}

}

func (ws *WServer) sendMsgResp(conn *UserConn, m *Req, pb *pbChat.SendMsgResp) {
	var mReplyData sdk_ws.UserSendMsgResp
	mReplyData.ClientMsgID = pb.GetClientMsgID()
	mReplyData.ServerMsgID = pb.GetServerMsgID()
	mReplyData.SendTime = pb.GetSendTime()
	mReplyData.Ex = pb.GetEx()
	b, _ := proto.Marshal(&mReplyData)
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		ErrCode:       pb.GetErrCode(),
		ErrMsg:        pb.GetErrMsg(),
		OperationID:   m.OperationID,
		Data:          b,
	}
	ws.sendMsg(conn, mReply)

}

func (ws *WServer) sendSignalMsgReq(conn *UserConn, m *Req) {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", m.OperationID))
	logger.Infof("Ws call success to sendSignalMsgReq start: send_id: %s, msg_incr: %d, req_identifier: %d, data: %s", m.SendID, m.MsgIncr, m.ReqIdentifier, string(m.Data))
	nReply := new(pbChat.SendMsgResp)
	isPass, errCode, errMsg, pData := ws.argsValidate(m, constant.WSSendSignalMsg, m.OperationID)
	if isPass {
		signalResp := pbRtc.SignalResp{}
		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImRealTimeCommName, m.OperationID)
		if etcdConn == nil {
			errMsg := m.OperationID + "getcdv3.GetDefaultConn == nil"
			logger.Error(errMsg)
			ws.sendSignalMsgResp(conn, 204, errMsg, m, &signalResp)
			return
		}
		rtcClient := pbRtc.NewRtcServiceClient(etcdConn)
		req := &pbRtc.SignalMessageAssembleReq{
			SignalReq:   pData.(*pbRtc.SignalReq),
			OperationID: m.OperationID,
		}
		respPb, err := rtcClient.SignalMessageAssemble(context.Background(), req)
		if err != nil {
			logger.Error("SignalMessageAssemble", err.Error(), config.Config.RpcRegisterName.OpenImRealTimeCommName)
			ws.sendSignalMsgResp(conn, 204, "grpc SignalMessageAssemble failed: "+err.Error(), m, &signalResp)
			return
		}
		signalResp.Payload = respPb.SignalResp.Payload
		msgData := sdk_ws.MsgData{}
		utils.CopyStructFields(&msgData, respPb.MsgData)
		logger.Info(respPb.String())
		if respPb.IsPass {
			pbData := pbChat.SendMsgReq{
				Token:       m.Token,
				OperationID: m.OperationID,
				MsgData:     &msgData,
			}
			logger.Info("pbData: ", pbData)
			logger.Info("Ws call success to sendSignalMsgReq middle", m.ReqIdentifier, m.SendID, m.MsgIncr, msgData)
			// etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, m.OperationID)
			// if etcdConn == nil {
			// 	errMsg := m.OperationID + "getcdv3.GetDefaultConn == nil"
			// 	logger.Error(errMsg)
			// 	ws.sendSignalMsgResp(conn, 200, errMsg, m, &signalResp)
			// 	return
			// }
			// client := pbChat.NewMsgClient(etcdConn)
			reply, err := ws.messageClient.SendMsg(context.Background(), &pbData)
			if err != nil {
				logger.Error("rpc sendMsg err", err.Error())
				nReply.ErrCode = 200
				nReply.ErrMsg = err.Error()
				ws.sendSignalMsgResp(conn, 200, err.Error(), m, &signalResp)
			} else {
				logger.Info("rpc call success to sendMsgReq", reply.String(), signalResp.String(), m)
				ws.sendSignalMsgResp(conn, 0, "", m, &signalResp)
			}
		} else {
			logger.Error(respPb.IsPass, respPb.CommonResp.ErrCode, respPb.CommonResp.ErrMsg)
			ws.sendSignalMsgResp(conn, respPb.CommonResp.ErrCode, respPb.CommonResp.ErrMsg, m, &signalResp)
		}
	} else {
		ws.sendSignalMsgResp(conn, errCode, errMsg, m, nil)
	}

}
func (ws *WServer) sendSignalMsgResp(conn *UserConn, errCode int32, errMsg string, m *Req, pb *pbRtc.SignalResp) {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", m.OperationID))
	// := make(map[string]interface{})
	logger.Debug("sendSignalMsgResp is", pb.String())
	b, _ := proto.Marshal(pb)
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		ErrCode:       errCode,
		ErrMsg:        errMsg,
		OperationID:   m.OperationID,
		Data:          b,
	}
	ws.sendMsg(conn, mReply)
}

func (ws *WServer) sendMsg(conn *UserConn, mReply interface{}) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(mReply)
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", mReply.(Resp).OperationID))
	if err != nil {
		//	uid, platform := ws.getUserUid(conn)
		logger.Errorf("req_identifier: %d, code: %d, msg: %s, remote_addr: %s, err: %v", mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg, conn.RemoteAddr().String(), err.Error())
		return
	}
	err = ws.writeMsg(conn, websocket.BinaryMessage, b.Bytes())
	if err != nil {
		//	uid, platform := ws.getUserUid(conn)
		logger.Errorf("req_identifier: %d, code: %d, msg: %s, remote_addr: %s, err: %v", mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg, conn.RemoteAddr().String(), err.Error())
	} else {
		logger.Debugf("req_identifier: %d, code: %d, msg: %s", mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg)
	}
}
func (ws *WServer) sendErrMsg(conn *UserConn, errCode int32, errMsg string, reqIdentifier int32, msgIncr string, operationID string) {
	mReply := Resp{
		ReqIdentifier: reqIdentifier,
		MsgIncr:       msgIncr,
		ErrCode:       errCode,
		ErrMsg:        errMsg,
		OperationID:   operationID,
	}
	ws.sendMsg(conn, mReply)
}

func SetTokenKicked(userID string, platformID int, operationID string) {
	ctx := logx.ContextWithFields(context.Background(), logx.Field("op", operationID))
	logger := logx.WithContext(ctx)
	m, err := db.DB.GetTokenMapByUidPid(ctx, userID, constant.PlatformIDToName(platformID))
	if err != nil {
		logger.Errorf("GetTokenMapByUidPid failed: err: %v, user_id: %s, platform: %s ", err.Error(), userID, constant.PlatformIDToName(platformID))
		return
	}
	for k := range m {
		m[k] = constant.KickedToken
	}
	err = db.DB.SetTokenMapByUidPid(ctx, userID, platformID, m)
	if err != nil {
		logger.Errorf("SetTokenMapByUidPid failed: err: %v, user_id: %s, platform: %s", err.Error(), userID, constant.PlatformIDToName(platformID))
		return
	}
}

func (ws *WServer) setUserDeviceBackground(conn *UserConn, m *Req) {
	logger := logx.WithContext(context.Background()).WithFields(logx.Field("op", m.OperationID))
	isPass, errCode, errMsg, pData := ws.argsValidate(m, constant.WsSetBackgroundStatus, m.OperationID)
	if isPass {
		req := pData.(*sdk_ws.SetAppBackgroundStatusReq)
		conn.IsBackground = req.IsBackground
		callbackResp := callbackUserOnline(m.OperationID, conn.userID, int(conn.PlatformID), conn.token, conn.IsBackground, conn.connID)
		if callbackResp.ErrCode != 0 {
			logger.Error("callbackUserOffline failed: resp: ", callbackResp)
		}
		logger.Infof("SetUserDeviceBackground success: conn: %v, background: %v", *conn, req.IsBackground)
	}
	ws.setUserDeviceBackgroundResp(conn, m, errCode, errMsg)
}

func (ws *WServer) setUserDeviceBackgroundResp(conn *UserConn, m *Req, errCode int32, errMsg string) {
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		OperationID:   m.OperationID,
		ErrCode:       errCode,
		ErrMsg:        errMsg,
	}
	ws.sendMsg(conn, mReply)
}
