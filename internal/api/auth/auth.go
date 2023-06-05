package apiAuth

import (
	authclient "Open_IM/internal/rpc/auth/client"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	rpc "Open_IM/pkg/proto/auth"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	authClient authclient.AuthClient
)

func init() {
	authClient = authclient.NewAuthClient(config.ConvertClientConfig(config.Config.ClientConfigs.Auth))
}

// @Summary 用户注册
// @Description 用户注册
// @Tags 鉴权认证
// @ID UserRegister
// @Accept json
// @Param req body api.UserRegisterReq true "secret为openIM密钥, 详细见服务端config.yaml secret字段 <br> platform为平台ID <br> ex为拓展字段 <br> gender为性别, 0为女, 1为男"
// @Produce json
// @Success 0 {object} api.UserRegisterResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /auth/user_register [post]
func UserRegister(c *gin.Context) {
	params := api.UserRegisterReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := " BindJSON failed " + err.Error()
		logx.Error(errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	logger := logx.WithContext(c.Request.Context()).WithFields(logx.Field("op", params.OperationID))
	if params.Secret != config.Config.Secret {
		errMsg := " params.Secret != config.Config.Secret "
		logger.Error(errMsg, params.Secret, config.Config.Secret)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": errMsg})
		return
	}
	req := &rpc.UserRegisterReq{UserInfo: &open_im_sdk.UserInfo{}}
	utils.CopyStructFields(req.UserInfo, &params)
	//copier.Copy(req.UserInfo, &params)
	req.OperationID = params.OperationID
	reply, err := authClient.UserRegister(c.Request.Context(), req)
	if err != nil {
		errMsg := req.OperationID + " " + "UserRegister failed " + err.Error() + req.String()
		logger.Error(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	if reply.CommonResp.ErrCode != 0 {
		errMsg := req.OperationID + " " + " UserRegister failed " + reply.CommonResp.ErrMsg + req.String()
		logger.Error(errMsg)
		if reply.CommonResp.ErrCode == constant.RegisterLimit {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterLimit, "errMsg": "用户注册被限制"})
		} else if reply.CommonResp.ErrCode == constant.InvitationError {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.InvitationError, "errMsg": "邀请码错误"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		}
		return
	}

	pbDataToken := &rpc.UserTokenReq{Platform: params.Platform, FromUserID: params.UserID, OperationID: params.OperationID}
	replyToken, err := authClient.UserToken(c.Request.Context(), pbDataToken)
	if err != nil {
		errMsg := req.OperationID + " " + " client.UserToken failed " + err.Error() + pbDataToken.String()
		logger.Error(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	resp := api.UserRegisterResp{CommResp: api.CommResp{ErrCode: replyToken.CommonResp.ErrCode, ErrMsg: replyToken.CommonResp.ErrMsg},
		UserToken: api.UserTokenInfo{UserID: req.UserInfo.UserID, Token: replyToken.Token, ExpiredTime: replyToken.ExpiredTime}}
	log.NewInfo(req.OperationID, "UserRegister return ", resp)
	c.JSON(http.StatusOK, resp)

}

// @Summary 用户登录
// @Description 获取用户的token
// @Tags 鉴权认证
// @ID UserToken
// @Accept json
// @Param req body api.UserTokenReq true "secret为openIM密钥, 详细见服务端config.yaml secret字段 <br> platform为平台ID"
// @Produce json
// @Success 0 {object} api.UserTokenResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /auth/user_token [post]
func UserToken(c *gin.Context) {
	params := api.UserTokenReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := " BindJSON failed " + err.Error()
		logx.Error(errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	logger := logx.WithContext(c.Request.Context()).WithFields(logx.Field("op", params.OperationID))
	if params.Secret != config.Config.Secret {
		errMsg := params.OperationID + " params.Secret != config.Config.Secret "
		logger.Error("params.Secret != config.Config.Secret", params.Secret, config.Config.Secret)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": errMsg})
		return
	}
	req := &rpc.UserTokenReq{Platform: params.Platform, FromUserID: params.UserID, OperationID: params.OperationID, LoginIp: params.LoginIp}
	reply, err := authClient.UserToken(c.Request.Context(), req)
	if err != nil {
		errMsg := req.OperationID + " UserToken failed " + err.Error() + " req: " + req.String()
		logger.Error(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	resp := api.UserTokenResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg},
		UserToken: api.UserTokenInfo{UserID: req.FromUserID, Token: reply.Token, ExpiredTime: reply.ExpiredTime}}
	c.JSON(http.StatusOK, resp)
}

// @Summary 解析当前用户token
// @Description 解析当前用户token(token在请求头中传入)
// @Tags 鉴权认证
// @ID ParseToken
// @Accept json
// @Param token header string true "im token"
// @Param req body api.ParseTokenReq true "secret为openIM密钥, 详细见服务端config.yaml secret字段<br>platform为平台ID"
// @Produce json
// @Success 0 {object} api.ParseTokenResp{Data=api.ExpireTime}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /auth/parse_token [post]
func ParseToken(c *gin.Context) {
	params := api.ParseTokenReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := " BindJSON failed " + err.Error()
		logx.Error(errMsg)
		c.JSON(http.StatusOK, gin.H{"errCode": 1001, "errMsg": errMsg})
		return
	}
	logger := logx.WithContext(c.Request.Context()).WithFields(logx.Field("op", params.OperationID))
	var ok bool
	var errInfo string
	var expireTime int64
	ok, _, errInfo, expireTime = token_verify.GetUserIDFromTokenExpireTime(c.Request.Context(), c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromTokenExpireTime failed " + errInfo
		logger.Error(errMsg)
		c.JSON(http.StatusOK, gin.H{"errCode": 1001, "errMsg": errMsg})
		return
	}

	resp := api.ParseTokenResp{CommResp: api.CommResp{ErrCode: 0, ErrMsg: ""}, ExpireTime: api.ExpireTime{ExpireTimeSeconds: uint32(expireTime)}}
	resp.Data = structs.Map(&resp.ExpireTime)
	c.JSON(http.StatusOK, resp)
}

// @Summary 强制登出
// @Description 对应的平台强制登出
// @Tags 鉴权认证
// @ID ForceLogout
// @Accept json
// @Param token header string true "im token"
// @Param req body api.ForceLogoutReq true "platform为平台ID <br> fromUserID为要执行强制登出的用户ID"
// @Produce json
// @Success 0 {object} api.ForceLogoutResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /auth/force_logout [post]
func ForceLogout(c *gin.Context) {
	params := api.ForceLogoutReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := " BindJSON failed " + err.Error()
		logx.Error(errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	logger := logx.WithContext(c.Request.Context()).WithFields(logx.Field("op", params.OperationID))

	req := &rpc.ForceLogoutReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Context(), c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		logger.Error(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reply, err := authClient.ForceLogout(c.Request.Context(), req)
	if err != nil {
		errMsg := req.OperationID + " UserToken failed " + err.Error() + req.String()
		logger.Error(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	resp := api.ForceLogoutResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	c.JSON(http.StatusOK, resp)
}
