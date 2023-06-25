package errors

import (
	"net/http"
)

var (
	BadRequest    = Error{HttpStatusCode: http.StatusBadRequest, Code: 400, Message: ""}
	Unauthorized  = Error{HttpStatusCode: http.StatusUnauthorized, Code: 401, Message: ""}
	InternalError = Error{HttpStatusCode: http.StatusInternalServerError, Code: 500, Message: ""}
)
