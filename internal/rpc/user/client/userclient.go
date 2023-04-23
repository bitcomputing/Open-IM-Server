package client

import (
	pbuser "Open_IM/pkg/proto/user"
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type UserClient interface {
	GetUserInfo(ctx context.Context, in *pbuser.GetUserInfoReq, opts ...grpc.CallOption) (*pbuser.GetUserInfoResp, error)
	UpdateUserInfo(ctx context.Context, in *pbuser.UpdateUserInfoReq, opts ...grpc.CallOption) (*pbuser.UpdateUserInfoResp, error)
	SetGlobalRecvMessageOpt(ctx context.Context, in *pbuser.SetGlobalRecvMessageOptReq, opts ...grpc.CallOption) (*pbuser.SetGlobalRecvMessageOptResp, error)
	GetAllUserID(ctx context.Context, in *pbuser.GetAllUserIDReq, opts ...grpc.CallOption) (*pbuser.GetAllUserIDResp, error)
	AccountCheck(ctx context.Context, in *pbuser.AccountCheckReq, opts ...grpc.CallOption) (*pbuser.AccountCheckResp, error)
	GetConversation(ctx context.Context, in *pbuser.GetConversationReq, opts ...grpc.CallOption) (*pbuser.GetConversationResp, error)
	GetAllConversations(ctx context.Context, in *pbuser.GetAllConversationsReq, opts ...grpc.CallOption) (*pbuser.GetAllConversationsResp, error)
	GetConversations(ctx context.Context, in *pbuser.GetConversationsReq, opts ...grpc.CallOption) (*pbuser.GetConversationsResp, error)
	BatchSetConversations(ctx context.Context, in *pbuser.BatchSetConversationsReq, opts ...grpc.CallOption) (*pbuser.BatchSetConversationsResp, error)
	SetConversation(ctx context.Context, in *pbuser.SetConversationReq, opts ...grpc.CallOption) (*pbuser.SetConversationResp, error)
	SetRecvMsgOpt(ctx context.Context, in *pbuser.SetRecvMsgOptReq, opts ...grpc.CallOption) (*pbuser.SetRecvMsgOptResp, error)
	GetUsers(ctx context.Context, in *pbuser.GetUsersReq, opts ...grpc.CallOption) (*pbuser.GetUsersResp, error)
	AddUser(ctx context.Context, in *pbuser.AddUserReq, opts ...grpc.CallOption) (*pbuser.AddUserResp, error)
	BlockUser(ctx context.Context, in *pbuser.BlockUserReq, opts ...grpc.CallOption) (*pbuser.BlockUserResp, error)
	UnBlockUser(ctx context.Context, in *pbuser.UnBlockUserReq, opts ...grpc.CallOption) (*pbuser.UnBlockUserResp, error)
	GetBlockUsers(ctx context.Context, in *pbuser.GetBlockUsersReq, opts ...grpc.CallOption) (*pbuser.GetBlockUsersResp, error)
}

type userClient struct {
	client zrpc.Client
}

func NewUserClient(config zrpc.RpcClientConf) UserClient {
	return &userClient{
		client: zrpc.MustNewClient(config),
	}
}

func (c *userClient) GetUserInfo(ctx context.Context, in *pbuser.GetUserInfoReq, opts ...grpc.CallOption) (*pbuser.GetUserInfoResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.GetUserInfo(ctx, in, opts...)
}

func (c *userClient) UpdateUserInfo(ctx context.Context, in *pbuser.UpdateUserInfoReq, opts ...grpc.CallOption) (*pbuser.UpdateUserInfoResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.UpdateUserInfo(ctx, in, opts...)
}

func (c *userClient) SetGlobalRecvMessageOpt(ctx context.Context, in *pbuser.SetGlobalRecvMessageOptReq, opts ...grpc.CallOption) (*pbuser.SetGlobalRecvMessageOptResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.SetGlobalRecvMessageOpt(ctx, in, opts...)
}

func (c *userClient) GetAllUserID(ctx context.Context, in *pbuser.GetAllUserIDReq, opts ...grpc.CallOption) (*pbuser.GetAllUserIDResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.GetAllUserID(ctx, in, opts...)
}

func (c *userClient) AccountCheck(ctx context.Context, in *pbuser.AccountCheckReq, opts ...grpc.CallOption) (*pbuser.AccountCheckResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.AccountCheck(ctx, in, opts...)
}

func (c *userClient) GetConversation(ctx context.Context, in *pbuser.GetConversationReq, opts ...grpc.CallOption) (*pbuser.GetConversationResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.GetConversation(ctx, in, opts...)
}

func (c *userClient) GetAllConversations(ctx context.Context, in *pbuser.GetAllConversationsReq, opts ...grpc.CallOption) (*pbuser.GetAllConversationsResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.GetAllConversations(ctx, in, opts...)
}

func (c *userClient) GetConversations(ctx context.Context, in *pbuser.GetConversationsReq, opts ...grpc.CallOption) (*pbuser.GetConversationsResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.GetConversations(ctx, in, opts...)
}

func (c *userClient) BatchSetConversations(ctx context.Context, in *pbuser.BatchSetConversationsReq, opts ...grpc.CallOption) (*pbuser.BatchSetConversationsResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.BatchSetConversations(ctx, in, opts...)
}

func (c *userClient) SetConversation(ctx context.Context, in *pbuser.SetConversationReq, opts ...grpc.CallOption) (*pbuser.SetConversationResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.SetConversation(ctx, in, opts...)
}

func (c *userClient) SetRecvMsgOpt(ctx context.Context, in *pbuser.SetRecvMsgOptReq, opts ...grpc.CallOption) (*pbuser.SetRecvMsgOptResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.SetRecvMsgOpt(ctx, in, opts...)
}

func (c *userClient) GetUsers(ctx context.Context, in *pbuser.GetUsersReq, opts ...grpc.CallOption) (*pbuser.GetUsersResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.GetUsers(ctx, in, opts...)
}

func (c *userClient) AddUser(ctx context.Context, in *pbuser.AddUserReq, opts ...grpc.CallOption) (*pbuser.AddUserResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.AddUser(ctx, in, opts...)
}

func (c *userClient) BlockUser(ctx context.Context, in *pbuser.BlockUserReq, opts ...grpc.CallOption) (*pbuser.BlockUserResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.BlockUser(ctx, in, opts...)
}

func (c *userClient) UnBlockUser(ctx context.Context, in *pbuser.UnBlockUserReq, opts ...grpc.CallOption) (*pbuser.UnBlockUserResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.UnBlockUser(ctx, in, opts...)
}

func (c *userClient) GetBlockUsers(ctx context.Context, in *pbuser.GetBlockUsersReq, opts ...grpc.CallOption) (*pbuser.GetBlockUsersResp, error) {
	client := pbuser.NewUserClient(c.client.Conn())
	return client.GetBlockUsers(ctx, in, opts...)
}
