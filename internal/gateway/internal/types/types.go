// Code generated by goctl. DO NOT EDIT.
package types

type RegisterRequest struct {
	Secret   string `json:"secret" validate:"required,max=32"`
	Platform int32  `json:"platform" validate:"required,min=1,max=12"`
	ApiUserInfo
	OperationID string `json:"operationID" validate:"required"`
}

type ApiUserInfo struct {
	UserID      string `json:"userID" validate:"required,min=1,max=64"`
	Nickname    string `json:"nickname" validate:"omitempty,min=1,max=64"`
	FaceURL     string `json:"faceURL" validate:"omitempty,max=1024"`
	Gender      int32  `json:"gender" validate:"omitempty,oneof=0 1 2"`
	PhoneNumber string `json:"phoneNumber" validate:"omitempty,max=32"`
	Birth       uint32 `json:"birth" validate:"omitempty"`
	Email       string `json:"email" validate:"omitempty,max=64"`
	CreateTime  int64  `json:"createTime"`
	LoginLimit  int32  `json:"loginLimit" validate:"omitempty"`
	Ex          string `json:"ex" validate:"omitempty,max=1024"`
	BirthStr    string `json:"birthStr" validate:"omitempty"`
}

type RegisterResponse struct {
	CommResp
	UserToken UserTokenInfo `json:"data"`
}

type CommResp struct {
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type UserTokenInfo struct {
	UserID      string `json:"userID"`
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expiredTime"`
}

type LoginRequest struct {
	Secret      string `json:"secret" validate:"required,max=32"`
	Platform    int32  `json:"platform" validate:"required,min=1,max=12"`
	UserID      string `json:"userID" validate:"required,min=1,max=64"`
	LoginIp     string `json:"loginIp"`
	OperationID string `json:"operationID" validate:"required"`
}

type LoginResponse struct {
	CommResp
	UserToken UserTokenInfo `json:"data"`
}

type ParseTokenRequest struct {
	OperationID string `json:"operationID" validate:"required"`
}

type ParseTokenResponse struct {
	CommResp
	Data       map[string]interface{} `json:"data"`
	ExpireTime ExpireTime             `json:"-"`
}

type ExpireTime struct {
	ExpireTimeSeconds uint32 `json:"expireTimeSeconds"`
}

type ForceLogoutRequest struct {
	Platform    int32  `json:"platform" validate:"required,min=1,max=12"`
	FromUserID  string `json:"fromUserID" validate:"required,min=1,max=64"`
	OperationID string `json:"operationID" validate:"required"`
}

type ForceLogoutResponse struct {
	CommResp
}
