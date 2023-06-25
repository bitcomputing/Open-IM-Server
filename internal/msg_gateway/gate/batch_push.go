package gate

import (
	"Open_IM/pkg/common/db"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

var MaxPullMsgNum = 100

func (r *RPCServer) GenPullSeqList(ctx context.Context, currentSeq uint32, operationID string, userID string) ([]uint32, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	maxSeq, err := db.DB.GetUserMaxSeq(ctx, userID)
	if err != nil {
		logger.Errorf("GetUserMaxSeq failed: user_id: %s, err: %v", userID, err)
		return nil, utils.Wrap(err, "")
	}

	var seqList []uint32
	num := 0
	for i := currentSeq + 1; i < uint32(maxSeq); i++ {
		seqList = append(seqList, i)
		num++
		if num == MaxPullMsgNum {
			break
		}
	}

	return seqList, nil
}

func (r *RPCServer) GetSingleUserMsgForPushPlatforms(operationID string, msgData *sdk_ws.MsgData, pushToUserID string, platformIDList []int) map[int]*sdk_ws.MsgDataList {
	user2PushMsg := make(map[int]*sdk_ws.MsgDataList, 0)
	for _, v := range platformIDList {
		user2PushMsg[v] = r.GetSingleUserMsgForPush(operationID, msgData, pushToUserID, v)
		//log.Info(operationID, "GetSingleUserMsgForPush", msgData.Seq, pushToUserID, v, "len:", len(user2PushMsg[v]))
	}
	return user2PushMsg
}

func (r *RPCServer) GetSingleUserMsgForPush(operationID string, msgData *sdk_ws.MsgData, pushToUserID string, platformID int) *sdk_ws.MsgDataList {
	//msgData.MsgDataList = nil
	return &sdk_ws.MsgDataList{MsgDataList: []*sdk_ws.MsgData{msgData}}

	//userConn := ws.getUserConn(pushToUserID, platformID)
	//if userConn == nil {
	//	log.Debug(operationID, "userConn == nil")
	//	return []*sdk_ws.MsgData{msgData}
	//}
	//
	//if msgData.Seq <= userConn.PushedMaxSeq {
	//	log.Debug(operationID, "msgData.Seq <= userConn.PushedMaxSeq", msgData.Seq, userConn.PushedMaxSeq)
	//	return nil
	//}
	//
	//msgList := r.GetSingleUserMsg(operationID, msgData.Seq, pushToUserID)
	//if msgList == nil {
	//	log.Debug(operationID, "GetSingleUserMsg msgList == nil", msgData.Seq, userConn.PushedMaxSeq)
	//	userConn.PushedMaxSeq = msgData.Seq
	//	return []*sdk_ws.MsgData{msgData}
	//}
	//msgList = append(msgList, msgData)
	//
	//for _, v := range msgList {
	//	if v.Seq > userConn.PushedMaxSeq {
	//		userConn.PushedMaxSeq = v.Seq
	//	}
	//}
	//log.Debug(operationID, "GetSingleUserMsg msgList len ", len(msgList), userConn.PushedMaxSeq)
	//return msgList
}

func (r *RPCServer) GetSingleUserMsg(ctx context.Context, operationID string, currentMsgSeq uint32, userID string) []*sdk_ws.MsgData {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", operationID))
	logger := logx.WithContext(ctx)

	seqList, err := r.GenPullSeqList(ctx, currentMsgSeq, operationID, userID)
	if err != nil {
		logger.Errorf("GenPullSeqList failed: err: %v, current_msg_seq: %d, user_id: %s", err, currentMsgSeq, userID)
		return nil
	}
	if len(seqList) == 0 {
		logger.Errorf("GenPullSeqList len == 0, current_msg_seq: %d, user_id: %s", currentMsgSeq, userID)
		return nil
	}
	rpcReq := sdk_ws.PullMessageBySeqListReq{}
	rpcReq.SeqList = seqList
	rpcReq.UserID = userID
	rpcReq.OperationID = operationID

	reply, err := r.msgClient.PullMessageBySeqList(context.Background(), &rpcReq)
	if err != nil {
		logger.Errorf("PullMessageBySeqList failed: err: %v, resp: %v ", err.Error(), rpcReq)
		return nil
	}
	if len(reply.List) == 0 {
		return nil
	}
	return reply.List
}

//func (r *RPCServer) GetBatchUserMsgForPush(operationID string, msgData *sdk_ws.MsgData, pushToUserIDList []string, platformID int) map[string][]*sdk_ws.MsgData {
//	user2PushMsg := make(map[string][]*sdk_ws.MsgData, 0)
//	for _, v := range pushToUserIDList {
//		user2PushMsg[v] = r.GetSingleUserMsgForPush(operationID, msgData, v, platformID)
//	}
//	return user2PushMsg
//}
