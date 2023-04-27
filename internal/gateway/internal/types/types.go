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
	Nickname    string `json:"nickname,optional" validate:"omitempty,min=1,max=64"`
	FaceURL     string `json:"faceURL,optional" validate:"omitempty,max=1024"`
	Gender      int32  `json:"gender,optional" validate:"omitempty,oneof=0 1 2"`
	PhoneNumber string `json:"phoneNumber,optional" validate:"omitempty,max=32"`
	Birth       uint32 `json:"birth,optional" validate:"omitempty"`
	Email       string `json:"email,optional" validate:"omitempty,max=64"`
	CreateTime  int64  `json:"createTime"`
	LoginLimit  int32  `json:"loginLimit,optional" validate:"omitempty"`
	Ex          string `json:"ex,optional" validate:"omitempty,max=1024"`
	BirthStr    string `json:"birthStr,optional" validate:"omitempty"`
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

type SetClientConfigRequest struct {
	OperationID     string  `json:"operationID"  validate:"required"`
	DiscoverPageURL *string `json:"discoverPageURL"`
}

type SetClientConfigResponse struct {
	CommResp
}

type GetClientConfigRequest struct {
	OperationID string `json:"operationID"  validate:"required"`
}

type GetClientConfigResponse struct {
	CommResp
	Data DiscoverPageURL `json:"data"`
}

type DiscoverPageURL struct {
	DiscoverPageURL string `json:"discoverPageURL"`
}

type GetAllConversationsRequest struct {
	OwnerUserID string `json:"ownerUserID" validate:"required"`
	OperationID string `json:"operationID" validate:"required"`
}

type GetAllConversationsResponse struct {
	CommResp
	Conversations []Conversation `json:"data"`
}

type Conversation struct {
	OwnerUserID           string `json:"ownerUserID" validate:"required"`
	ConversationID        string `json:"conversationID" validate:"required"`
	ConversationType      int32  `json:"conversationType" validate:"required"`
	UserID                string `json:"userID"`
	GroupID               string `json:"groupID"`
	RecvMsgOpt            int32  `json:"recvMsgOpt,optional"  validate:"omitempty,oneof=0 1 2"`
	UnreadCount           int32  `json:"unreadCount,optional"  validate:"omitempty"`
	DraftTextTime         int64  `json:"draftTextTime"`
	IsPinned              bool   `json:"isPinned,optional" validate:"omitempty"`
	IsPrivateChat         bool   `json:"isPrivateChat"`
	BurnDuration          int32  `json:"burnDuration"`
	GroupAtType           int32  `json:"groupAtType"`
	IsNotInGroup          bool   `json:"isNotInGroup"`
	UpdateUnreadCountTime int64  `json:"updateUnreadCountTime"`
	AttachedInfo          string `json:"attachedInfo"`
	Ex                    string `json:"ex"`
}

type GetConversationRequest struct {
	ConversationID string `json:"conversationID" validate:"required"`
	OwnerUserID    string `json:"ownerUserID" validate:"required"`
	OperationID    string `json:"operationID" validate:"required"`
}

type GetConversationResponse struct {
	CommResp
	Conversation Conversation `json:"data"`
}

type GetConversationsRequest struct {
	ConversationIDs []string `json:"conversationIDs" validate:"required"`
	OwnerUserID     string   `json:"ownerUserID" validate:"required"`
	OperationID     string   `json:"operationID" validate:"required"`
}

type GetConversationsResponse struct {
	CommResp
	Conversations []Conversation `json:"data"`
}

type SetConversationRequest struct {
	Conversation
	NotificationType int32  `json:"notificationType"`
	OperationID      string `json:"operationID" validate:"required"`
}

type SetConversationResponse struct {
	CommResp
}

type BatchSetConversationsRequest struct {
	Conversations    []Conversation `json:"conversations" validate:"required"`
	NotificationType int32          `json:"notificationType"`
	OwnerUserID      string         `json:"ownerUserID" validate:"required"`
	OperationID      string         `json:"operationID" validate:"required"`
}

type BatchSetConversationsResponse struct {
	CommResp
	Data SuccessAndFailed `json:"data"`
}

type SuccessAndFailed struct {
	Success []string `json:"success"`
	Failed  []string `json:"failed"`
}

type SetRecvMsgOptRequest struct {
	OwnerUserID      string `json:"ownerUserID" validate:"required"`
	ConversationID   string `json:"conversationID"`
	RecvMsgOpt       int32  `json:"recvMsgOpt,optional"  validate:"omitempty,oneof=0 1 2"`
	OperationID      string `json:"operationID" validate:"required"`
	NotificationType int32  `json:"notificationType"`
}

type SetRecvMsgOptResponse struct {
	CommResp
}

type ModifyConversationFieldRequest struct {
	Conversation
	FieldType   int32    `json:"fieldType" validate:"required"`
	UserIDList  []string `json:"userIDList" validate:"required"`
	OperationID string   `json:"operationID" validate:"required"`
}

type ModifyConversationFieldResponse struct {
	CommResp
}

type AddFriendRequest struct {
	ParamsCommFriend
	ReqMsg string `json:"reqMsg"`
}

type ParamsCommFriend struct {
	OperationID string `json:"operationID" validate:"required"`
	ToUserID    string `json:"toUserID" validate:"required"`
	FromUserID  string `json:"fromUserID" validate:"required"`
}

type AddFriendResponse struct {
	CommResp
}

type DeleteFriendRequest struct {
	ParamsCommFriend
}

type DeleteFriendResponse struct {
	CommResp
}

type GetFriendApplyListRequest struct {
	OperationID string `json:"operationID" validate:"required"`
	FromUserID  string `json:"fromUserID" validate:"required"`
}

type GetFriendApplyListResponse struct {
	CommResp
	Data []map[string]interface{} `json:"data"`
}

type GetSelfFriendApplyListRequest struct {
	OperationID string `json:"operationID" validate:"required"`
	FromUserID  string `json:"fromUserID" validate:"required"`
}

type GetSelfFriendApplyListResponse struct {
	CommResp
	Data []map[string]interface{} `json:"data"`
}

type GetFriendListRequest struct {
	OperationID string `json:"operationID" validate:"required"`
	FromUserID  string `json:"fromUserID" validate:"required"`
}

type GetFriendListResponse struct {
	CommResp
	Data []map[string]interface{} `json:"data"`
}

type RespondFriendApplyRequest struct {
	ParamsCommFriend
	Flag      int32  `json:"flag" validate:"required,oneof=-1 0 1"`
	HandleMsg string `json:"handleMsg"`
}

type RespondFriendApplyResponse struct {
	CommResp
}

type SetFriendRemarkRequest struct {
	ParamsCommFriend
	Remark string `json:"remark"`
}

type SetFriendRemarkResponse struct {
	CommResp
}

type AddFriendBlacklistRequest struct {
	ParamsCommFriend
}

type AddFriendBlacklistResponse struct {
	CommResp
}

type GetFriendBlacklistRequest struct {
	OperationID string `json:"operationID" validate:"required"`
	FromUserID  string `json:"fromUserID" validate:"required"`
}

type GetFriendBlacklistResponse struct {
	CommResp
	Data []map[string]interface{} `json:"data"`
}

type RemoveFriendBlacklistRequest struct {
	ParamsCommFriend
}

type RemoveFriendBlacklistResponse struct {
	CommResp
}

type ImportFriendRequest struct {
	FriendUserIDList []string `json:"friendUserIDList" validate:"required"`
	OperationID      string   `json:"operationID" validate:"required"`
	FromUserID       string   `json:"fromUserID" validate:"required"`
}

type ImportFriendResponse struct {
	CommResp
	UserIDResultList []UserIDResult `json:"data"`
}

type UserIDResult struct {
	UserID string `json:"userID"`
	Result int32  `json:"result"`
}

type CheckFriendRequest struct {
	ParamsCommFriend
}

type CheckFriendResponse struct {
	CommResp
	Response IsFriend `json:"data"`
}

type IsFriend struct {
	Friend bool `json:"isFriend"`
}
