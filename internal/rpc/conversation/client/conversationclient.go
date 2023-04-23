package client

import (
	"context"

	pbconversation "Open_IM/pkg/proto/conversation"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type ConversationClient interface{}

type conversationClient struct {
	client zrpc.Client
}

func NewConversationClient(config zrpc.RpcClientConf) ConversationClient {
	return &conversationClient{
		client: zrpc.MustNewClient(config),
	}
}

func (c *conversationClient) ModifyConversationField(ctx context.Context, in *pbconversation.ModifyConversationFieldReq, opts ...grpc.CallOption) (*pbconversation.ModifyConversationFieldResp, error) {
	client := pbconversation.NewConversationClient(c.client.Conn())
	return client.ModifyConversationField(ctx, in, opts...)
}
