package group

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	rpc "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetJoinedSuperGroupList(c *gin.Context) {
	req := api.GetJoinedSuperGroupListReq{}
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, opUserID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPb := rpc.GetJoinedSuperGroupListReq{OperationID: req.OperationID, OpUserID: opUserID, UserID: req.FromUserID}
	rpcResp, err := groupClient.GetJoinedSuperGroupList(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, "InviteUserToGroup failed ", err.Error(), reqPb.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	GroupListResp := api.GetJoinedSuperGroupListResp{GetJoinedGroupListResp: api.GetJoinedGroupListResp{CommResp: api.CommResp{ErrCode: rpcResp.CommonResp.ErrCode, ErrMsg: rpcResp.CommonResp.ErrMsg}, GroupInfoList: rpcResp.GroupList}}
	GroupListResp.Data = jsonData.JsonDataList(GroupListResp.GroupInfoList)
	log.NewInfo(req.OperationID, "GetJoinedSuperGroupList api return ", GroupListResp)
	c.JSON(http.StatusOK, GroupListResp)
}

func GetSuperGroupsInfo(c *gin.Context) {
	req := api.GetSuperGroupsInfoReq{}
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, opUserID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPb := rpc.GetSuperGroupsInfoReq{OperationID: req.OperationID, OpUserID: opUserID, GroupIDList: req.GroupIDList}
	rpcResp, err := groupClient.GetSuperGroupsInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, "InviteUserToGroup failed ", err.Error(), reqPb.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.GetSuperGroupsInfoResp{GetGroupInfoResp: api.GetGroupInfoResp{CommResp: api.CommResp{ErrCode: rpcResp.CommonResp.ErrCode, ErrMsg: rpcResp.CommonResp.ErrMsg}, GroupInfoList: rpcResp.GroupInfoList}}
	resp.Data = jsonData.JsonDataList(resp.GroupInfoList)
	log.NewInfo(req.OperationID, "GetGroupsInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}
