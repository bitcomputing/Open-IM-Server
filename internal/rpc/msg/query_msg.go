package msg

import (
	commonDB "Open_IM/pkg/common/db"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/proto/msg"
	"context"

	go_redis "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

func (rpc *rpcChat) GetSuperGroupMsg(ctx context.Context, req *msg.GetSuperGroupMsgReq) (*msg.GetSuperGroupMsgResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := new(msg.GetSuperGroupMsgResp)
	redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(ctx, req.GroupID, []uint32{req.Seq}, req.OperationID)
	if err != nil {
		if err != go_redis.Nil {
			promePkg.PromeAdd(promePkg.MsgPullFromRedisFailedCounter, len(failedSeqList))
			logger.Error("get message from redis exception", err.Error(), failedSeqList)
		} else {
			logger.Debug("get message from redis is nil", failedSeqList)
		}
		msgList, _, err1 := commonDB.DB.GetSuperGroupMsgBySeqListMongo(ctx, req.GroupID, failedSeqList, req.OperationID)
		if err1 != nil {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoFailedCounter, len(failedSeqList))
			logger.Error("GetSuperGroupMsg data error", req.String(), err.Error())
			resp.ErrCode = 201
			resp.ErrMsg = err.Error()
			return resp, nil
		} else {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoSuccessCounter, len(msgList))
			redisMsgList = append(redisMsgList, msgList...)
			for _, m := range msgList {
				resp.MsgData = m
			}
		}
	} else {
		promePkg.PromeAdd(promePkg.MsgPullFromRedisSuccessCounter, len(redisMsgList))
		for _, m := range redisMsgList {
			resp.MsgData = m
		}
	}

	return resp, nil
}

func (rpc *rpcChat) GetWriteDiffMsg(context context.Context, req *msg.GetWriteDiffMsgReq) (*msg.GetWriteDiffMsgResp, error) {
	panic("implement me")
}
