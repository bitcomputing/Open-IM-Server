package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/proto/msg"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

func (rpc *rpcChat) DelMsgList(ctx context.Context, req *commonPb.DelMsgListReq) (*commonPb.DelMsgListResp, error) {
	resp := &commonPb.DelMsgListResp{}
	select {
	case rpc.delMsgCh <- deleteMsg{
		UserID:      req.UserID,
		OpUserID:    req.OpUserID,
		SeqList:     req.SeqList,
		OperationID: req.OperationID,
	}:
	case <-time.After(1 * time.Second):
		resp.ErrCode = constant.ErrSendLimit.ErrCode
		resp.ErrMsg = constant.ErrSendLimit.ErrMsg
		return resp, nil
	}

	return resp, nil
}

func (rpc *rpcChat) DelSuperGroupMsg(ctx context.Context, req *msg.DelSuperGroupMsgReq) (*msg.DelSuperGroupMsgResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if !token_verify.CheckAccess(req.OpUserID, req.UserID) {
		logger.Error("CheckAccess false: ", req.OpUserID, req.UserID)
		return &msg.DelSuperGroupMsgResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	resp := &msg.DelSuperGroupMsgResp{}
	groupMaxSeq, err := db.DB.GetGroupMaxSeq(ctx, req.GroupID)
	if err != nil {
		logger.Error("GetGroupMaxSeq false: ", req.OpUserID, req.UserID, req.GroupID)
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
		return resp, nil
	}
	err = db.DB.SetGroupUserMinSeq(ctx, req.GroupID, req.UserID, groupMaxSeq)
	if err != nil {
		logger.Error("SetGroupUserMinSeq false ", req.OpUserID, req.UserID, req.GroupID)
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
		return resp, nil
	}

	return resp, nil
}
