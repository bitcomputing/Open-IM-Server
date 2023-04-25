package errors

import (
	"Open_IM/pkg/common/constant"
	"net/http"
)

var (
	BadRequest    = Error{HttpStatusCode: http.StatusBadRequest, Code: 400, Message: ""}
	Unauthorized  = Error{HttpStatusCode: http.StatusBadRequest, Code: 401, Message: ""} // todo change to matching HTTP status code
	InternalError = Error{HttpStatusCode: http.StatusInternalServerError, Code: 500, Message: ""}

	RegisterLimit    = Error{HttpStatusCode: http.StatusOK, Code: constant.RegisterLimit, Message: "用户注册被限制"}
	InvitationError  = Error{HttpStatusCode: http.StatusOK, Code: constant.InvitationError, Message: "邀请码错误"}
	ParseTokenFailed = Error{HttpStatusCode: http.StatusOK, Code: 1001, Message: ""}
)
