package gate

import (
	messageclient "Open_IM/internal/rpc/msg/client"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/discovery"
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"strconv"

	go_redis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	//"gopkg.in/errgo.v2/errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metric"
)

type UserConn struct {
	*websocket.Conn
	w            *sync.Mutex
	PlatformID   int32
	PushedMaxSeq uint32
	IsCompress   bool
	userID       string
	IsBackground bool
	token        string
	connID       string
}

type WServer struct {
	wsAddr                string
	wsMaxConnNum          int
	wsUpGrader            *websocket.Upgrader
	wsUserToConn          map[string]map[int][]*UserConn
	gatewayClient         *discovery.Client
	messageClient         messageclient.MsgClient
	msgRecvTotal          metric.CounterVec
	getNewestSeqTotal     metric.CounterVec
	pullMsgBySeqListTotal metric.CounterVec
	msgOnlinePushSuccess  metric.CounterVec
	onlineUser            metric.GaugeVec
}

func NewWebsocketServer(wsPort int) *WServer {
	gc, err := discovery.NewClient(config.ConvertClientConfig(config.Config.ClientConfigs.Gateway))
	if err != nil {
		panic(err)
	}
	ws = &WServer{
		wsAddr:       fmt.Sprintf(":%d", wsPort),
		wsMaxConnNum: config.Config.LongConnSvr.WebsocketMaxConnNum,
		wsUpGrader: &websocket.Upgrader{
			HandshakeTimeout: time.Duration(config.Config.LongConnSvr.WebsocketTimeOut) * time.Second,
			ReadBufferSize:   config.Config.LongConnSvr.WebsocketMaxMsgLen,
			CheckOrigin:      func(r *http.Request) bool { return true },
		},
		wsUserToConn:  make(map[string]map[int][]*UserConn),
		gatewayClient: gc,
		messageClient: messageclient.NewMsgClient(config.ConvertClientConfig(config.Config.ClientConfigs.Message)),
		msgRecvTotal: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "msg_recv_total",
			Help: "The number of msg received",
		}),
		getNewestSeqTotal: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "get_newest_seq_total",
			Help: "the number of get newest seq",
		}),
		pullMsgBySeqListTotal: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "pull_msg_by_seq_list_total",
			Help: "The number of pull msg by seq list",
		}),
		msgOnlinePushSuccess: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "msg_online_push_success",
			Help: "The number of msg successful online pushed",
		}),
		onlineUser: metric.NewGaugeVec(&metric.GaugeVecOpts{
			Name: "online_user_num",
			Help: "The number of online user num",
		}),
	}

	return ws
}

func (ws *WServer) Start() {
	http.HandleFunc("/", ws.wsHandler) //Get request from client to handle by wsHandler
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	err := http.ListenAndServe(ws.wsAddr, nil) //Start listening
	if err != nil {
		panic("Ws listening err:" + err.Error())
	}
}

func (ws *WServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	operationID := ""
	if len(query["operationID"]) != 0 {
		operationID = query["operationID"][0]
	} else {
		operationID = utils.OperationIDGenerator()
	}
	logger := logx.WithContext(r.Context()).WithFields(logx.Field("op", operationID))
	logger.Debug("args: ", query)
	if isPass, compression := ws.headerCheck(w, r, operationID); isPass {
		conn, err := ws.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
		if err != nil {
			logger.Errorf("upgrade http conn err: %v, query: %s", err, query)
			return
		} else {
			newConn := &UserConn{conn, new(sync.Mutex), utils.StringToInt32(query["platformID"][0]), 0, compression, query["sendID"][0], false, query["token"][0], utils.Md5(conn.RemoteAddr().String() + "_" + strconv.Itoa(int(utils.GetCurrentTimestampByMill())))}
			ws.addUserConn(query["sendID"][0], utils.StringToInt(query["platformID"][0]), newConn, query["token"][0], newConn.connID, operationID)
			go ws.readMsg(newConn)
		}
	} else {
		logger.Error("headerCheck failed: query: %v", query)
	}
}

func (ws *WServer) readMsg(conn *UserConn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if messageType == websocket.PingMessage {
			logx.Info("this is a  pingMessage")
		}
		if err != nil {
			logx.Error("WS ReadMsg error:", " userIP: ", conn.RemoteAddr().String(), " error: ", err.Error())
			ws.delUserConn(conn)
			return
		}
		if messageType == websocket.CloseMessage {
			logx.Error("WS receive error: ", " userIP: ", conn.RemoteAddr().String(), " error: ", string(msg))
			ws.delUserConn(conn)
			return
		}

		if conn.IsCompress {
			buff := bytes.NewBuffer(msg)
			reader, err := gzip.NewReader(buff)
			if err != nil {
				logx.Error("ungzip read failed: ", err)
				continue
			}
			msg, err = ioutil.ReadAll(reader)
			if err != nil {
				logx.Error("ReadAll failed: ", err)
				continue
			}
			err = reader.Close()
			if err != nil {
				logx.Error("reader close failed: ", err)
			}
		}
		ws.msgParse(conn, msg)
	}
}

func (ws *WServer) SetWriteTimeout(conn *UserConn, timeout int) {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (ws *WServer) writeMsg(conn *UserConn, a int, msg []byte) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	if conn.IsCompress {
		var buffer bytes.Buffer
		gz := gzip.NewWriter(&buffer)
		if _, err := gz.Write(msg); err != nil {
			return utils.Wrap(err, "")
		}
		if err := gz.Close(); err != nil {
			return utils.Wrap(err, "")
		}
		msg = buffer.Bytes()
	}
	conn.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))
	return conn.WriteMessage(a, msg)
}

func (ws *WServer) SetWriteTimeoutWriteMsg(conn *UserConn, a int, msg []byte, timeout int) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	return conn.WriteMessage(a, msg)
}

func (ws *WServer) MultiTerminalLoginRemoteChecker(userID string, platformID int32, token string, operationID string) {
	ctx := logx.ContextWithFields(context.Background(), logx.Field("op", operationID))
	logger := logx.WithContext(ctx)
	grpcCons := ws.gatewayClient.ClientConns()
	logger.Info("args  grpcCons: ", userID, platformID, grpcCons)

	for i := 0; i < len(grpcCons); i++ {
		if grpcCons[i].Target() == server.target {
			logger.Debug(operationID, "Filter out this node ", server.target)
			continue
		}
		logger.Debug(operationID, "call this node ", grpcCons[i].Target(), server.target)
		logger.Info("current target:", grpcCons[i].Target())
		client := pbRelay.NewRelayClient(grpcCons[i])
		req := &pbRelay.MultiTerminalLoginCheckReq{OperationID: operationID, PlatformID: platformID, UserID: userID, Token: token}
		logger.Info(operationID, "MultiTerminalLoginCheckReq ", client, req.String())
		resp, err := client.MultiTerminalLoginCheck(ctx, req)
		if err != nil {
			logger.Error(operationID, "MultiTerminalLoginCheck failed ", err.Error())
			continue
		}
		if resp.ErrCode != 0 {
			logger.Error(operationID, "MultiTerminalLoginCheck errCode, errMsg: ", resp.ErrCode, resp.ErrMsg)
			continue
		}
		logger.Debug(operationID, "MultiTerminalLoginCheck resp ", resp.String())
	}
}

func (ws *WServer) MultiTerminalLoginCheckerWithLock(uid string, platformID int, token string, operationID string) {
	ctx := logx.ContextWithFields(context.Background(), logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	rwLock.Lock()
	defer rwLock.Unlock()
	logger.Infof("user_id: %s, platform: %d, token: %s", uid, platformID, token)
	switch config.Config.MultiLoginPolicy {
	case constant.DefalutNotKick:
	case constant.PCAndOther:
		if constant.PlatformNameToClass(constant.PlatformIDToName(platformID)) == constant.TerminalPC {
			return
		}
		fallthrough
	case constant.AllLoginButSameTermKick:
		if oldConnMap, ok := ws.wsUserToConn[uid]; ok { // user->map[platform->conn]
			if oldConns, ok := oldConnMap[platformID]; ok {
				logger.Debugf("kick old conn: user_id: %s, platform: %d", uid, platformID)
				for _, conn := range oldConns {
					ws.sendKickMsg(conn, operationID)
				}
				m, err := db.DB.GetTokenMapByUidPid(context.Background(), uid, constant.PlatformIDToName(platformID))
				if err != nil && err != go_redis.Nil {
					logger.Error("get token from redis err: %v, user_id: %s, platform: %s", err, uid, constant.PlatformIDToName(platformID))
					return
				}
				if m == nil {
					logger.Errorf("get token from redis err: %v, user_id: %s, platform: %s, m is nil", uid, constant.PlatformIDToName(platformID))
					return
				}
				logger.Debugf("get token map is %v, user_id: %s, platform: %s", m, uid, constant.PlatformIDToName(platformID))

				for k := range m {
					if k != token {
						m[k] = constant.KickedToken
					}
				}
				logger.Debugf("set token map is %v, user_id: %s, platform: %s", m, uid, constant.PlatformIDToName(platformID))
				err = db.DB.SetTokenMapByUidPid(context.Background(), uid, platformID, m)
				if err != nil {
					logger.Errorf("SetTokenMapByUidPid err: %v, user_id: %s, platform: %d, m: %v", err, uid, platformID, m)
					return
				}

				delete(oldConnMap, platformID)
				ws.wsUserToConn[uid] = oldConnMap
				if len(oldConnMap) == 0 {
					delete(ws.wsUserToConn, uid)
				}
			} else {
				logger.Error("abnormal uid-conn: user_id: %s, platform: %s, conn_map: %v", uid, platformID, oldConnMap[platformID])
			}

		} else {
			logger.Debug("no other conn: user_to_conn: %v, user_id: %s, platform: %d", ws.wsUserToConn, uid, platformID)
		}
	case constant.SingleTerminalLogin:
	case constant.WebAndOther:
	}
}

func (ws *WServer) MultiTerminalLoginChecker(uid string, platformID int, newConn *UserConn, token string, operationID string) {
	ctx := logx.ContextWithFields(context.Background(), logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	switch config.Config.MultiLoginPolicy {
	case constant.DefalutNotKick:
	case constant.PCAndOther:
		if constant.PlatformNameToClass(constant.PlatformIDToName(platformID)) == constant.TerminalPC {
			return
		}
		fallthrough
	case constant.AllLoginButSameTermKick:
		if oldConnMap, ok := ws.wsUserToConn[uid]; ok { // user->map[platform->conn]
			if oldConns, ok := oldConnMap[platformID]; ok {
				logger.Debugf("kick old conn: user_id: %s, platform: %d", uid, platformID)
				for _, conn := range oldConns {
					ws.sendKickMsg(conn, operationID)
				}
				m, err := db.DB.GetTokenMapByUidPid(context.Background(), uid, constant.PlatformIDToName(platformID))
				if err != nil && err != go_redis.Nil {
					logger.Errorf("get token from redis err: %v, user_id: %s, platform: %s", err.Error(), uid, constant.PlatformIDToName(platformID))
					return
				}
				if m == nil {
					logger.Errorf("get token from redis err: %v, user_id: %s, platform: %s, m is nil", uid, constant.PlatformIDToName(platformID))
					return
				}
				logger.Debugf("get token map is %v, user_id: %s, platform: %s", m, uid, constant.PlatformIDToName(platformID))

				for k := range m {
					if k != token {
						m[k] = constant.KickedToken
					}
				}
				logger.Debugf("set token map is %v, user_id: %s, platform: %s", m, uid, constant.PlatformIDToName(platformID))
				err = db.DB.SetTokenMapByUidPid(context.Background(), uid, platformID, m)
				if err != nil {
					logger.Errorf("SetTokenMapByUidPid err: %v, user_id: %s, platform: %d, m: %v", err.Error(), uid, platformID, m)
					return
				}
				delete(oldConnMap, platformID)
				ws.wsUserToConn[uid] = oldConnMap
				if len(oldConnMap) == 0 {
					delete(ws.wsUserToConn, uid)
				}
				callbackResp := callbackUserKickOff(operationID, uid, platformID)
				if callbackResp.ErrCode != 0 {
					logger.Errorf("callbackUserOffline failed: resp: %v", callbackResp)
				}
			} else {
				logger.Debugf("normal uid-conn: user_id: %s, platform: %d, conn_map: %v", uid, platformID, oldConnMap[platformID])
			}

		} else {
			logger.Debug("no other conn: user_to_conn: %v, user_id: %s, platform: %d", ws.wsUserToConn, uid, platformID)
		}

	case constant.SingleTerminalLogin:
	case constant.WebAndOther:
	}
}

func (ws *WServer) sendKickMsg(oldConn *UserConn, operationID string) {
	ctx := logx.ContextWithFields(context.Background(), logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	mReply := Resp{
		ReqIdentifier: constant.WSKickOnlineMsg,
		ErrCode:       constant.ErrTokenInvalid.ErrCode,
		ErrMsg:        constant.ErrTokenInvalid.ErrMsg,
		OperationID:   operationID,
	}
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(mReply)
	if err != nil {
		logger.Errorf("req_identifier: %d, code: %d, msg: %s, err: %s, remote_addr: %s, err: %v", mReply.ReqIdentifier, mReply.ErrCode, mReply.ErrMsg, "Encode Msg error", oldConn.RemoteAddr().String(), err.Error())
		return
	}
	err = ws.writeMsg(oldConn, websocket.BinaryMessage, b.Bytes())
	if err != nil {
		logger.Errorf("req_identifier: %d, code: %d, msg: %s, err: %s, remote_addr: %s, err: %v", mReply.ReqIdentifier, mReply.ErrCode, mReply.ErrMsg, "sendKickMsg WS WriteMsg error", oldConn.RemoteAddr().String(), err.Error())
	}
	errClose := oldConn.Close()
	if errClose != nil {
		logger.Errorf("req_identifier: %d, code: %d, msg: %s, err: %s, remote_addr: %s, err: %v", mReply.ReqIdentifier, mReply.ErrCode, mReply.ErrMsg, "close old conn error", oldConn.RemoteAddr().String(), err.Error())

	}
}

func (ws *WServer) addUserConn(uid string, platformID int, conn *UserConn, token string, connID, operationID string) {
	ctx := logx.ContextWithFields(context.Background(), logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	rwLock.Lock()
	defer rwLock.Unlock()
	logger.Infof("user_id: %s, platform: %d, conn: %v, token: %s, remote_addr: %s", uid, platformID, conn, token, conn.RemoteAddr().String())
	callbackResp := callbackUserOnline(operationID, uid, platformID, token, false, connID)
	if callbackResp.ErrCode != 0 {
		logger.Error("callbackUserOnline resp: ", callbackResp)
	}
	go ws.MultiTerminalLoginRemoteChecker(uid, int32(platformID), token, operationID)
	ws.MultiTerminalLoginChecker(uid, platformID, conn, token, operationID)
	if oldConnMap, ok := ws.wsUserToConn[uid]; ok {
		if conns, ok := oldConnMap[platformID]; ok {
			conns = append(conns, conn)
			oldConnMap[platformID] = conns
		} else {
			var conns []*UserConn
			conns = append(conns, conn)
			oldConnMap[platformID] = conns
		}
		ws.wsUserToConn[uid] = oldConnMap
		logger.Debugf("user not first come in, add conn: user_id: %s, platform: %d, conn: %v, old_conn_map: %v", uid, platformID, conn, oldConnMap)
	} else {
		i := make(map[int][]*UserConn)
		var conns []*UserConn
		conns = append(conns, conn)
		i[platformID] = conns
		ws.wsUserToConn[uid] = i
		logger.Debugf("user first come in, new user, conn: user_id: %d, platform: %d, conn: %v, user_to_conn: %v", uid, platformID, conn, ws.wsUserToConn[uid])
	}
	count := 0
	for _, v := range ws.wsUserToConn {
		count = count + len(v)
	}
	ws.onlineUser.Inc()
	logger.Debug("WS Add operation", "", "wsUser added", ws.wsUserToConn, "connection_uid", uid, "connection_platform", constant.PlatformIDToName(platformID), "online_user_num", len(ws.wsUserToConn), "online_conn_num", count)
}

func (ws *WServer) delUserConn(conn *UserConn) {
	rwLock.Lock()
	defer rwLock.Unlock()
	operationID := utils.OperationIDGenerator()
	platform := int(conn.PlatformID)

	if oldConnMap, ok := ws.wsUserToConn[conn.userID]; ok { // only recycle self conn
		if oldconns, okMap := oldConnMap[platform]; okMap {

			var a []*UserConn

			for _, client := range oldconns {
				if client != conn {
					a = append(a, client)

				}
			}
			if len(a) != 0 {
				oldConnMap[platform] = a
			} else {
				delete(oldConnMap, platform)
			}

		}
		ws.wsUserToConn[conn.userID] = oldConnMap
		if len(oldConnMap) == 0 {
			delete(ws.wsUserToConn, conn.userID)
		}
		count := 0
		for _, v := range ws.wsUserToConn {
			count = count + len(v)
		}
		logx.Debug("WS delete operation", "", "wsUser deleted", ws.wsUserToConn, "disconnection_uid", conn.userID, "disconnection_platform", platform, "online_user_num", len(ws.wsUserToConn), "online_conn_num", count)
	}

	err := conn.Close()
	if err != nil {
		logx.Errorf("close err: user_id: %s, platform: %d", conn.userID, platform)
	}
	if conn.PlatformID == 0 || conn.connID == "" {
		logx.Errorf("PlatformID or connID is null: platform: %d, conn_id: %s", conn.PlatformID, conn.connID)
	}
	callbackResp := callbackUserOffline(operationID, conn.userID, int(conn.PlatformID), conn.connID)
	if callbackResp.ErrCode != 0 {
		logx.Error("callbackUserOffline failed: resp: ", callbackResp)
	}
	ws.onlineUser.Add(-1)

}

//	func (ws *WServer) getUserConn(uid string, platform int) *UserConn {
//		rwLock.RLock()
//		defer rwLock.RUnlock()
//		if connMap, ok := ws.wsUserToConn[uid]; ok {
//			if conn, flag := connMap[platform]; flag {
//				return conn
//			}
//		}
//		return nil
//	}
func (ws *WServer) getUserAllCons(uid string) map[int][]*UserConn {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if connMap, ok := ws.wsUserToConn[uid]; ok {
		newConnMap := make(map[int][]*UserConn)
		for k, v := range connMap {
			newConnMap[k] = v
		}
		return newConnMap
	}
	return nil
}

//	func (ws *WServer) getUserUid(conn *UserConn) (uid string, platform int) {
//		rwLock.RLock()
//		defer rwLock.RUnlock()
//
//		if stringMap, ok := ws.wsConnToUser[conn]; ok {
//			for k, v := range stringMap {
//				platform = k
//				uid = v
//			}
//			return uid, platform
//		}
//		return "", 0
//	}
func (ws *WServer) headerCheck(w http.ResponseWriter, r *http.Request, operationID string) (isPass, compression bool) {
	ctx := logx.ContextWithFields(r.Context(), logx.Field("op", operationID))
	logger := logx.WithContext(ctx)
	status := http.StatusUnauthorized
	query := r.URL.Query()
	if len(query["token"]) != 0 && len(query["sendID"]) != 0 && len(query["platformID"]) != 0 {
		if ok, err, msg := token_verify.WsVerifyToken(r.Context(), query["token"][0], query["sendID"][0], query["platformID"][0], operationID); !ok {
			if errors.Is(err, constant.ErrTokenExpired) {
				status = int(constant.ErrTokenExpired.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenInvalid) {
				status = int(constant.ErrTokenInvalid.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenMalformed) {
				status = int(constant.ErrTokenMalformed.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenNotValidYet) {
				status = int(constant.ErrTokenNotValidYet.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenUnknown) {
				status = int(constant.ErrTokenUnknown.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenKicked) {
				status = int(constant.ErrTokenKicked.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenDifferentPlatformID) {
				status = int(constant.ErrTokenDifferentPlatformID.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenDifferentUserID) {
				status = int(constant.ErrTokenDifferentUserID.ErrCode)
			}
			//switch errors.Cause(err) {
			//case constant.ErrTokenExpired:
			//	status = int(constant.ErrTokenExpired.ErrCode)
			//case constant.ErrTokenInvalid:
			//	status = int(constant.ErrTokenInvalid.ErrCode)
			//case constant.ErrTokenMalformed:
			//	status = int(constant.ErrTokenMalformed.ErrCode)
			//case constant.ErrTokenNotValidYet:
			//	status = int(constant.ErrTokenNotValidYet.ErrCode)
			//case constant.ErrTokenUnknown:
			//	status = int(constant.ErrTokenUnknown.ErrCode)
			//case constant.ErrTokenKicked:
			//	status = int(constant.ErrTokenKicked.ErrCode)
			//case constant.ErrTokenDifferentPlatformID:
			//	status = int(constant.ErrTokenDifferentPlatformID.ErrCode)
			//case constant.ErrTokenDifferentUserID:
			//	status = int(constant.ErrTokenDifferentUserID.ErrCode)
			//}

			logger.Errorf("Token verify failed: query: %v, msg: %s, err: %v, status: %d ", query, msg, err.Error(), status)
			w.Header().Set("Sec-Websocket-Version", "13")
			w.Header().Set("ws_err_msg", err.Error())
			http.Error(w, err.Error(), status)
			return false, false
		} else {
			if r.Header.Get("compression") == "gzip" {
				compression = true
			}
			if len(query["compression"]) != 0 && query["compression"][0] == "gzip" {
				compression = true
			}
			logger.Infof("Connection Authentication Success: token: %s, user_id: %s, platform: %s, compression: %v", query["token"][0], query["sendID"][0], query["platformID"][0], compression)
			return true, compression
		}
	} else {
		status = int(constant.ErrArgs.ErrCode)
		logger.Errorf("Args err: query: %v", query)
		w.Header().Set("Sec-Websocket-Version", "13")
		errMsg := "args err, need token, sendID, platformID"
		w.Header().Set("ws_err_msg", errMsg)
		http.Error(w, errMsg, status)
		return false, false
	}
}
