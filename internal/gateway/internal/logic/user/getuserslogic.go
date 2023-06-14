package user

import (
	"context"

	apiutils "Open_IM/internal/gateway/internal/common/utils"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	"Open_IM/pkg/common/token_verify"
	errors "Open_IM/pkg/errors/api"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersLogic {
	return &GetUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersLogic) GetUsers(req *types.GetUsersRequest) (*types.GetUsersResponse, error) {
	logger := l.Logger.WithFields(logx.Field("op", req.OperationID))

	var rpcReq user.GetUsersReq

	token, err := apiutils.GetTokenByContext(l.ctx, logger, req.OperationID)
	if err != nil {
		return nil, err
	}

	ok, _, errInfo := token_verify.GetUserIDFromToken(l.ctx, token, req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + token
		logger.Error(req.OperationID, errMsg)
		return nil, errors.BadRequest.WriteMessage(errMsg)
	}

	rpcReq.OperationID = req.OperationID
	rpcReq.UserID = req.UserID
	rpcReq.UserName = req.UserName
	rpcReq.Content = req.Content
	rpcReq.Pagination = &sdk.RequestPagination{ShowNumber: req.ShowNumber, PageNumber: req.PageNumber}

	rpcResp, err := l.svcCtx.UserClient.GetUsers(l.ctx, &rpcReq)
	if err != nil {
		logger.Error(err.Error(), rpcReq.String())
		return nil, errors.InternalError.WriteMessage(err.Error())
	}

	userList := []*types.CMSUser{}
	for _, v := range rpcResp.UserList {
		user := types.CMSUser{}
		if err := utils.CopyStructFields(&user, v.User); err != nil {
			logger.Error(err)
			return nil, errors.InternalError.WriteMessage(err.Error())
		}
		user.IsBlock = v.IsBlock
		userList = append(userList, &user)
	}

	return &types.GetUsersResponse{
		CommResp: types.CommResp{
			ErrCode: rpcResp.CommonResp.ErrCode,
			ErrMsg:  rpcResp.CommonResp.ErrMsg,
		},
		Data: types.GetUsersData{
			UserList:    userList,
			TotalNum:    rpcResp.TotalNums,
			CurrentPage: rpcResp.Pagination.CurrentPage,
			ShowNumber:  rpcResp.Pagination.ShowNumber,
		},
	}, nil
}
