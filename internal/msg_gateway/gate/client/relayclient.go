package client

import (
	pbrelay "Open_IM/pkg/proto/relay"
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type RelayClient interface {
	OnlinePushMsg(ctx context.Context, in *pbrelay.OnlinePushMsgReq, opts ...grpc.CallOption) (*pbrelay.OnlinePushMsgResp, error)
	GetUsersOnlineStatus(ctx context.Context, in *pbrelay.GetUsersOnlineStatusReq, opts ...grpc.CallOption) (*pbrelay.GetUsersOnlineStatusResp, error)
	OnlineBatchPushOneMsg(ctx context.Context, in *pbrelay.OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*pbrelay.OnlineBatchPushOneMsgResp, error)
	SuperGroupOnlineBatchPushOneMsg(ctx context.Context, in *pbrelay.OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*pbrelay.OnlineBatchPushOneMsgResp, error)
	KickUserOffline(ctx context.Context, in *pbrelay.KickUserOfflineReq, opts ...grpc.CallOption) (*pbrelay.KickUserOfflineResp, error)
	MultiTerminalLoginCheck(ctx context.Context, in *pbrelay.MultiTerminalLoginCheckReq, opts ...grpc.CallOption) (*pbrelay.MultiTerminalLoginCheckResp, error)
	SuperGroupBackgroundOnlinePush(ctx context.Context, in *pbrelay.OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*pbrelay.OnlineBatchPushOneMsgResp, error)
}

type relayClient struct {
	client zrpc.Client
}

func NewRelayClient(config zrpc.RpcClientConf) RelayClient {
	return &relayClient{
		client: zrpc.MustNewClient(config),
	}
}

func (c *relayClient) OnlinePushMsg(ctx context.Context, in *pbrelay.OnlinePushMsgReq, opts ...grpc.CallOption) (*pbrelay.OnlinePushMsgResp, error) {
	client := pbrelay.NewRelayClient(c.client.Conn())
	return client.OnlinePushMsg(ctx, in, opts...)
}

func (c *relayClient) GetUsersOnlineStatus(ctx context.Context, in *pbrelay.GetUsersOnlineStatusReq, opts ...grpc.CallOption) (*pbrelay.GetUsersOnlineStatusResp, error) {
	client := pbrelay.NewRelayClient(c.client.Conn())
	return client.GetUsersOnlineStatus(ctx, in, opts...)
}

func (c *relayClient) OnlineBatchPushOneMsg(ctx context.Context, in *pbrelay.OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*pbrelay.OnlineBatchPushOneMsgResp, error) {
	client := pbrelay.NewRelayClient(c.client.Conn())
	return client.OnlineBatchPushOneMsg(ctx, in, opts...)
}

func (c *relayClient) SuperGroupOnlineBatchPushOneMsg(ctx context.Context, in *pbrelay.OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*pbrelay.OnlineBatchPushOneMsgResp, error) {
	client := pbrelay.NewRelayClient(c.client.Conn())
	return client.SuperGroupOnlineBatchPushOneMsg(ctx, in, opts...)
}

func (c *relayClient) KickUserOffline(ctx context.Context, in *pbrelay.KickUserOfflineReq, opts ...grpc.CallOption) (*pbrelay.KickUserOfflineResp, error) {
	client := pbrelay.NewRelayClient(c.client.Conn())
	return client.KickUserOffline(ctx, in, opts...)
}

func (c *relayClient) MultiTerminalLoginCheck(ctx context.Context, in *pbrelay.MultiTerminalLoginCheckReq, opts ...grpc.CallOption) (*pbrelay.MultiTerminalLoginCheckResp, error) {
	client := pbrelay.NewRelayClient(c.client.Conn())
	return client.MultiTerminalLoginCheck(ctx, in, opts...)
}

func (c *relayClient) SuperGroupBackgroundOnlinePush(ctx context.Context, in *pbrelay.OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*pbrelay.OnlineBatchPushOneMsgResp, error) {
	client := pbrelay.NewRelayClient(c.client.Conn())
	return client.SuperGroupBackgroundOnlinePush(ctx, in, opts...)
}
