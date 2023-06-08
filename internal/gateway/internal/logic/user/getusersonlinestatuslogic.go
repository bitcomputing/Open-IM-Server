package user

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	relay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersOnlineStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersOnlineStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersOnlineStatusLogic {
	return &GetUsersOnlineStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersOnlineStatusLogic) GetUsersOnlineStatus(req *types.GetUsersOnlineStatusRequest) (resp *types.GetUsersOnlineStatusResponse, err error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	rpcReq := &relay.GetUsersOnlineStatusReq{}
	if err := utils.CopyStructFields(rpcReq, &req); err != nil {
		logger.Error(err)
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, opuid, errInfo := token_verify.GetUserIDFromToken(token, rpcReq.OperationID)
	if !ok {
		errMsg := rpcReq.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(rpcReq.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}
	rpcReq.OpUserID = opuid

	if len(config.Config.Manager.AppManagerUid) == 0 {
		logger.Error(rpcReq.OperationID, "Manager == 0")
		return nil, errors.InternalError.WriteMessage("Manager == 0")
	}
	rpcReq.OpUserID = config.Config.Manager.AppManagerUid[0]

	var wsResult []*relay.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*relay.GetUsersOnlineStatusResp_SuccessResult
	flag := false

	// todo
	// reply, err := l.svcCtx..GetUsersOnlineStatus(l.ctx, rpcReq)
	// 	if err != nil {
	// 		log.NewError(req.OperationID, "GetUsersOnlineStatus rpc  err", rpcReq.String(), err.Error())
	// 		continue
	// 	} else {
	// 		if reply.ErrCode == 0 {
	// 			wsResult = append(wsResult, reply.SuccessResult...)
	// 		}
	// 	}

	for _, v1 := range req.UserIDList {
		flag = false
		temp := new(relay.GetUsersOnlineStatusResp_SuccessResult)
		for _, v2 := range wsResult {
			if v2.UserID == v1 {
				flag = true
				temp.UserID = v1
				temp.Status = constant.OnlineStatus
				temp.DetailPlatformStatus = append(temp.DetailPlatformStatus, v2.DetailPlatformStatus...)
			}

		}
		if !flag {
			temp.UserID = v1
			temp.Status = constant.OfflineStatus
		}
		respResult = append(respResult, temp)
	}

	successResult := []*types.GetUsersOnlineStatusResp_SuccessResult{}
	for _, v := range respResult {
		detailPlatformStatus := []*types.GetUsersOnlineStatusResp_SuccessDetail{}
		for _, this := range v.DetailPlatformStatus {
			detailPlatformStatus = append(detailPlatformStatus, &types.GetUsersOnlineStatusResp_SuccessDetail{
				Platform:     this.Platform,
				Status:       this.Status,
				ConnID:       this.ConnID,
				IsBackground: this.IsBackground,
			})
		}
		successResult = append(successResult, &types.GetUsersOnlineStatusResp_SuccessResult{
			UserID:               v.UserID,
			Status:               v.Status,
			DetailPlatformStatus: detailPlatformStatus,
		})
	}

	return &types.GetUsersOnlineStatusResponse{
		CommResp: types.CommResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
		SuccessResult: successResult,
	}, nil
}
