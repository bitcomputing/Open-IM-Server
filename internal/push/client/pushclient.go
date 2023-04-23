package client

import (
	pbpush "Open_IM/pkg/proto/push"
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type PushClient interface {
	PushMsg(ctx context.Context, in *pbpush.PushMsgReq, opts ...grpc.CallOption) (*pbpush.PushMsgResp, error)
	DelUserPushToken(ctx context.Context, in *pbpush.DelUserPushTokenReq, opts ...grpc.CallOption) (*pbpush.DelUserPushTokenResp, error)
}

type pushClient struct {
	client zrpc.Client
}

func NewPushClient(config zrpc.RpcClientConf) PushClient {
	return &pushClient{
		client: zrpc.MustNewClient(config),
	}
}

func (c *pushClient) PushMsg(ctx context.Context, in *pbpush.PushMsgReq, opts ...grpc.CallOption) (*pbpush.PushMsgResp, error) {
	client := pbpush.NewPushMsgServiceClient(c.client.Conn())
	return client.PushMsg(ctx, in, opts...)
}

func (c *pushClient) DelUserPushToken(ctx context.Context, in *pbpush.DelUserPushTokenReq, opts ...grpc.CallOption) (*pbpush.DelUserPushTokenResp, error) {
	client := pbpush.NewPushMsgServiceClient(c.client.Conn())
	return client.DelUserPushToken(ctx, in, opts...)
}
