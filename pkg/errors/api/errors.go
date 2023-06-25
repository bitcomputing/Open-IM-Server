package errors

import "fmt"

type Error struct {
	HttpStatusCode int    `json:"-"`
	Code           int32  `json:"errCode"`
	Message        string `json:"errMsg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("httpStatusCode: %d, errCode: %d, errMsg: %s", e.HttpStatusCode, e.Code, e.Message)
}

func (e *Error) WriteMessage(message string) Error {
	e.Message = message
	return *e
}
