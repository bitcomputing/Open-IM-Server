package client

import (
	pbauth "Open_IM/pkg/proto/auth"
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type AuthClient interface {
	UserRegister(context.Context, *pbauth.UserRegisterReq, ...grpc.CallOption) (*pbauth.UserRegisterResp, error)
	UserToken(context.Context, *pbauth.UserTokenReq, ...grpc.CallOption) (*pbauth.UserTokenResp, error)
	ForceLogout(context.Context, *pbauth.ForceLogoutReq, ...grpc.CallOption) (*pbauth.ForceLogoutResp, error)
	ParseToken(context.Context, *pbauth.ParseTokenReq, ...grpc.CallOption) (*pbauth.ParseTokenResp, error)
}

type authClient struct {
	client zrpc.Client
}

func NewAuthClient(config zrpc.RpcClientConf) AuthClient {
	return &authClient{
		client: zrpc.MustNewClient(config),
	}
}

func (c *authClient) UserRegister(ctx context.Context, in *pbauth.UserRegisterReq, opts ...grpc.CallOption) (*pbauth.UserRegisterResp, error) {
	client := pbauth.NewAuthClient(c.client.Conn())
	return client.UserRegister(ctx, in, opts...)
}

func (c *authClient) UserToken(ctx context.Context, in *pbauth.UserTokenReq, opts ...grpc.CallOption) (*pbauth.UserTokenResp, error) {
	client := pbauth.NewAuthClient(c.client.Conn())
	return client.UserToken(ctx, in, opts...)
}

func (c *authClient) ForceLogout(ctx context.Context, in *pbauth.ForceLogoutReq, opts ...grpc.CallOption) (*pbauth.ForceLogoutResp, error) {
	client := pbauth.NewAuthClient(c.client.Conn())
	return client.ForceLogout(ctx, in, opts...)
}

func (c *authClient) ParseToken(ctx context.Context, in *pbauth.ParseTokenReq, opts ...grpc.CallOption) (*pbauth.ParseTokenResp, error) {
	client := pbauth.NewAuthClient(c.client.Conn())
	return client.ParseToken(ctx, in, opts...)
}
