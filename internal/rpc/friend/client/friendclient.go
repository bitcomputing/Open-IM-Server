package client

import (
	"context"

	pbfriend "Open_IM/pkg/proto/friend"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type FriendClient interface {
	AddFriend(context.Context, *pbfriend.AddFriendReq, ...grpc.CallOption) (*pbfriend.AddFriendResp, error)
	GetFriendApplyList(context.Context, *pbfriend.GetFriendApplyListReq, ...grpc.CallOption) (*pbfriend.GetFriendApplyListResp, error)
	GetSelfApplyList(context.Context, *pbfriend.GetSelfApplyListReq, ...grpc.CallOption) (*pbfriend.GetSelfApplyListResp, error)
	GetFriendList(context.Context, *pbfriend.GetFriendListReq, ...grpc.CallOption) (*pbfriend.GetFriendListResp, error)
	AddBlacklist(context.Context, *pbfriend.AddBlacklistReq, ...grpc.CallOption) (*pbfriend.AddBlacklistResp, error)
	RemoveBlacklist(context.Context, *pbfriend.RemoveBlacklistReq, ...grpc.CallOption) (*pbfriend.RemoveBlacklistResp, error)
	IsFriend(context.Context, *pbfriend.IsFriendReq, ...grpc.CallOption) (*pbfriend.IsFriendResp, error)
	IsInBlackList(context.Context, *pbfriend.IsInBlackListReq, ...grpc.CallOption) (*pbfriend.IsInBlackListResp, error)
	GetBlacklist(context.Context, *pbfriend.GetBlacklistReq, ...grpc.CallOption) (*pbfriend.GetBlacklistResp, error)
	DeleteFriend(context.Context, *pbfriend.DeleteFriendReq, ...grpc.CallOption) (*pbfriend.DeleteFriendResp, error)
	AddFriendResponse(context.Context, *pbfriend.AddFriendResponseReq, ...grpc.CallOption) (*pbfriend.AddFriendResponseResp, error)
	SetFriendRemark(context.Context, *pbfriend.SetFriendRemarkReq, ...grpc.CallOption) (*pbfriend.SetFriendRemarkResp, error)
	ImportFriend(context.Context, *pbfriend.ImportFriendReq, ...grpc.CallOption) (*pbfriend.ImportFriendResp, error)
}

type friendClient struct {
	client zrpc.Client
}

func NewFriendClient(config zrpc.RpcClientConf) FriendClient {
	return &friendClient{
		client: zrpc.MustNewClient(config),
	}
}

func (c *friendClient) AddFriend(ctx context.Context, in *pbfriend.AddFriendReq, opts ...grpc.CallOption) (*pbfriend.AddFriendResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.AddFriend(ctx, in, opts...)
}

func (c *friendClient) GetFriendApplyList(ctx context.Context, in *pbfriend.GetFriendApplyListReq, opts ...grpc.CallOption) (*pbfriend.GetFriendApplyListResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.GetFriendApplyList(ctx, in, opts...)
}

func (c *friendClient) GetSelfApplyList(ctx context.Context, in *pbfriend.GetSelfApplyListReq, opts ...grpc.CallOption) (*pbfriend.GetSelfApplyListResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.GetSelfApplyList(ctx, in, opts...)
}

func (c *friendClient) GetFriendList(ctx context.Context, in *pbfriend.GetFriendListReq, opts ...grpc.CallOption) (*pbfriend.GetFriendListResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.GetFriendList(ctx, in, opts...)
}

func (c *friendClient) AddBlacklist(ctx context.Context, in *pbfriend.AddBlacklistReq, opts ...grpc.CallOption) (*pbfriend.AddBlacklistResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.AddBlacklist(ctx, in, opts...)
}

func (c *friendClient) RemoveBlacklist(ctx context.Context, in *pbfriend.RemoveBlacklistReq, opts ...grpc.CallOption) (*pbfriend.RemoveBlacklistResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.RemoveBlacklist(ctx, in, opts...)
}

func (c *friendClient) IsFriend(ctx context.Context, in *pbfriend.IsFriendReq, opts ...grpc.CallOption) (*pbfriend.IsFriendResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.IsFriend(ctx, in, opts...)
}

func (c *friendClient) IsInBlackList(ctx context.Context, in *pbfriend.IsInBlackListReq, opts ...grpc.CallOption) (*pbfriend.IsInBlackListResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.IsInBlackList(ctx, in, opts...)
}

func (c *friendClient) GetBlacklist(ctx context.Context, in *pbfriend.GetBlacklistReq, opts ...grpc.CallOption) (*pbfriend.GetBlacklistResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.GetBlacklist(ctx, in, opts...)
}

func (c *friendClient) DeleteFriend(ctx context.Context, in *pbfriend.DeleteFriendReq, opts ...grpc.CallOption) (*pbfriend.DeleteFriendResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.DeleteFriend(ctx, in, opts...)
}

func (c *friendClient) AddFriendResponse(ctx context.Context, in *pbfriend.AddFriendResponseReq, opts ...grpc.CallOption) (*pbfriend.AddFriendResponseResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.AddFriendResponse(ctx, in, opts...)
}

func (c *friendClient) SetFriendRemark(ctx context.Context, in *pbfriend.SetFriendRemarkReq, opts ...grpc.CallOption) (*pbfriend.SetFriendRemarkResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.SetFriendRemark(ctx, in, opts...)
}

func (c *friendClient) ImportFriend(ctx context.Context, in *pbfriend.ImportFriendReq, opts ...grpc.CallOption) (*pbfriend.ImportFriendResp, error) {
	client := pbfriend.NewFriendClient(c.client.Conn())
	return client.ImportFriend(ctx, in, opts...)
}
