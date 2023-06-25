package msg

import (
	commonDB "Open_IM/pkg/common/db"
	promePkg "Open_IM/pkg/common/prometheus"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"context"

	go_redis "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

func (rpc *rpcChat) GetMaxAndMinSeq(ctx context.Context, in *open_im_sdk.GetMaxAndMinSeqReq) (*open_im_sdk.GetMaxAndMinSeqResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", in.OperationID))
	logger := logx.WithContext(ctx)

	resp := new(open_im_sdk.GetMaxAndMinSeqResp)
	m := make(map[string]*open_im_sdk.MaxAndMinSeq)
	var maxSeq, minSeq uint64
	var err1, err2 error
	maxSeq, err1 = commonDB.DB.GetUserMaxSeq(ctx, in.UserID)
	minSeq, err2 = commonDB.DB.GetUserMinSeq(ctx, in.UserID)
	if (err1 != nil && err1 != go_redis.Nil) || (err2 != nil && err2 != go_redis.Nil) {
		logger.Error("getMaxSeq from redis error", in.String())
		if err1 != nil {
			logger.Error(err1.Error())
		}
		if err2 != nil {
			logger.Error(err2.Error())
		}
		resp.ErrCode = 200
		resp.ErrMsg = "redis get err"
		return resp, nil
	}
	resp.MaxSeq = uint32(maxSeq)
	resp.MinSeq = uint32(minSeq)
	for _, groupID := range in.GroupIDList {
		x := new(open_im_sdk.MaxAndMinSeq)
		maxSeq, _ := commonDB.DB.GetGroupMaxSeq(ctx, groupID)
		minSeq, _ := commonDB.DB.GetGroupUserMinSeq(ctx, groupID, in.UserID)
		x.MaxSeq = uint32(maxSeq)
		x.MinSeq = uint32(minSeq)
		m[groupID] = x
	}
	resp.GroupMaxAndMinSeq = m
	return resp, nil
}

func (rpc *rpcChat) PullMessageBySeqList(ctx context.Context, in *open_im_sdk.PullMessageBySeqListReq) (*open_im_sdk.PullMessageBySeqListResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", in.OperationID))
	logger := logx.WithContext(ctx)

	resp := new(open_im_sdk.PullMessageBySeqListResp)
	m := make(map[string]*open_im_sdk.MsgDataList)
	redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(ctx, in.UserID, in.SeqList, in.OperationID)
	if err != nil {
		if err != go_redis.Nil {
			promePkg.PromeAdd(promePkg.MsgPullFromRedisFailedCounter, len(failedSeqList))
			logger.Error("get message from redis exception", err.Error(), failedSeqList)
		} else {
			logger.Debug("get message from redis is nil", failedSeqList)
		}
		msgList, err1 := commonDB.DB.GetMsgBySeqListMongo2(ctx, in.UserID, failedSeqList, in.OperationID)
		if err1 != nil {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoFailedCounter, len(failedSeqList))
			logger.Error("PullMessageBySeqList data error", in.String(), err1.Error())
			resp.ErrCode = 201
			resp.ErrMsg = err1.Error()
			return resp, nil
		} else {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoSuccessCounter, len(msgList))
			redisMsgList = append(redisMsgList, msgList...)
			resp.List = redisMsgList
		}
	} else {
		promePkg.PromeAdd(promePkg.MsgPullFromRedisSuccessCounter, len(redisMsgList))
		resp.List = redisMsgList
	}

	for k, v := range in.GroupSeqList {
		x := new(open_im_sdk.MsgDataList)
		redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(ctx, k, v.SeqList, in.OperationID)
		if err != nil {
			if err != go_redis.Nil {
				promePkg.PromeAdd(promePkg.MsgPullFromRedisFailedCounter, len(failedSeqList))
				logger.Error("get message from redis exception", err.Error(), failedSeqList)
			} else {
				logger.Debug("get message from redis is nil", failedSeqList)
			}
			msgList, _, err1 := commonDB.DB.GetSuperGroupMsgBySeqListMongo(ctx, k, failedSeqList, in.OperationID)
			if err1 != nil {
				promePkg.PromeAdd(promePkg.MsgPullFromMongoFailedCounter, len(failedSeqList))
				logger.Error("PullMessageBySeqList data error", in.String(), err1.Error())
				resp.ErrCode = 201
				resp.ErrMsg = err1.Error()
				return resp, nil
			} else {
				promePkg.PromeAdd(promePkg.MsgPullFromMongoSuccessCounter, len(msgList))
				redisMsgList = append(redisMsgList, msgList...)
				x.MsgDataList = redisMsgList
				m[k] = x
			}
		} else {
			promePkg.PromeAdd(promePkg.MsgPullFromRedisSuccessCounter, len(redisMsgList))
			x.MsgDataList = redisMsgList
			m[k] = x
		}
	}
	resp.GroupMsgDataList = m
	return resp, nil
}

type MsgFormats []*open_im_sdk.MsgData

// Implement the sort.Interface interface to get the number of elements method
func (s MsgFormats) Len() int {
	return len(s)
}

// Implement the sort.Interface interface comparison element method
func (s MsgFormats) Less(i, j int) bool {
	return s[i].SendTime < s[j].SendTime
}

// Implement the sort.Interface interface exchange element method
func (s MsgFormats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
