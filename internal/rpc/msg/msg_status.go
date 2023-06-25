package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	pbMsg "Open_IM/pkg/proto/msg"
	"context"

	goRedis "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

func (rpc *rpcChat) SetSendMsgStatus(ctx context.Context, req *pbMsg.SetSendMsgStatusReq) (resp *pbMsg.SetSendMsgStatusResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbMsg.SetSendMsgStatusResp{}

	if err := db.DB.SetSendMsgStatus(ctx, req.Status, req.OperationID); err != nil {
		logger.Error(err)
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
		return resp, nil
	}

	return resp, nil
}

func (rpc *rpcChat) GetSendMsgStatus(ctx context.Context, req *pbMsg.GetSendMsgStatusReq) (resp *pbMsg.GetSendMsgStatusResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbMsg.GetSendMsgStatusResp{}
	status, err := db.DB.GetSendMsgStatus(ctx, req.OperationID)
	if err != nil {
		resp.Status = constant.MsgStatusNotExist
		if err == goRedis.Nil {
			logger.Info("not exist")
			return resp, nil
		} else {
			logger.Error(err)
			resp.ErrMsg = err.Error()
			resp.ErrCode = constant.ErrDB.ErrCode
			return resp, nil
		}
	}
	resp.Status = int32(status)

	return resp, nil
}
