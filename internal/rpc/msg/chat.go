package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/token_verify"
	pbChat "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

func (rpc *rpcChat) ClearMsg(ctx context.Context, req *pbChat.ClearMsgReq) (*pbChat.ClearMsgResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.UserID) {
		errMsg := "No permission" + req.OpUserID + req.UserID
		logger.Error(errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}

	err := db.DB.CleanUpOneUserAllMsgFromRedis(ctx, req.UserID, req.OperationID)
	if err != nil {
		errMsg := "CleanUpOneUserAllMsgFromRedis failed " + err.Error() + req.OperationID + req.UserID
		logger.Error(errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}

	err = db.DB.CleanUpUserMsgFromMongo(ctx, req.UserID, req.OperationID)
	if err != nil {
		errMsg := "CleanUpUserMsgFromMongo failed " + err.Error() + req.OperationID + req.UserID
		logger.Error(errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}

	resp := pbChat.ClearMsgResp{ErrCode: 0}

	return &resp, nil
}

func (rpc *rpcChat) SetMsgMinSeq(ctx context.Context, req *pbChat.SetMsgMinSeqReq) (*pbChat.SetMsgMinSeqResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.UserID) {
		errMsg := "No permission" + req.OpUserID + req.UserID
		logger.Error(errMsg)
		return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	if req.GroupID == "" {
		err := db.DB.SetUserMinSeq(ctx, req.UserID, req.MinSeq)
		if err != nil {
			errMsg := "SetUserMinSeq failed " + err.Error() + req.OperationID + req.UserID + utils.Uint32ToString(req.MinSeq)
			logger.Error(errMsg)
			return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
		}
		return &pbChat.SetMsgMinSeqResp{}, nil
	}
	err := db.DB.SetGroupUserMinSeq(ctx, req.GroupID, req.UserID, uint64(req.MinSeq))
	if err != nil {
		errMsg := "SetGroupUserMinSeq failed " + err.Error() + req.OperationID + req.GroupID + req.UserID + utils.Uint32ToString(req.MinSeq)
		logger.Error(errMsg)
		return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	return &pbChat.SetMsgMinSeqResp{}, nil
}
