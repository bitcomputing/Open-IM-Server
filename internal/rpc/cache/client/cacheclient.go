package client

import (
	"context"

	pbcache "Open_IM/pkg/proto/cache"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type CacheClient interface {
	// friendInfo
	GetFriendIDListFromCache(context.Context, *pbcache.GetFriendIDListFromCacheReq, ...grpc.CallOption) (*pbcache.GetFriendIDListFromCacheResp, error)
	// for dtm
	DelFriendIDListFromCache(context.Context, *pbcache.DelFriendIDListFromCacheReq, ...grpc.CallOption) (*pbcache.DelFriendIDListFromCacheResp, error)
	// blackList
	GetBlackIDListFromCache(context.Context, *pbcache.GetBlackIDListFromCacheReq, ...grpc.CallOption) (*pbcache.GetBlackIDListFromCacheResp, error)
	// for dtm
	DelBlackIDListFromCache(context.Context, *pbcache.DelBlackIDListFromCacheReq, ...grpc.CallOption) (*pbcache.DelBlackIDListFromCacheResp, error)
	// group
	GetGroupMemberIDListFromCache(context.Context, *pbcache.GetGroupMemberIDListFromCacheReq, ...grpc.CallOption) (*pbcache.GetGroupMemberIDListFromCacheResp, error)
	// for dtm
	DelGroupMemberIDListFromCache(context.Context, *pbcache.DelGroupMemberIDListFromCacheReq, ...grpc.CallOption) (*pbcache.DelGroupMemberIDListFromCacheResp, error)
}

type cacheClient struct {
	client zrpc.Client
}

func NewCacheClient(config zrpc.RpcClientConf) CacheClient {
	return &cacheClient{
		client: zrpc.MustNewClient(config),
	}
}

func (c *cacheClient) GetFriendIDListFromCache(ctx context.Context, in *pbcache.GetFriendIDListFromCacheReq, opts ...grpc.CallOption) (*pbcache.GetFriendIDListFromCacheResp, error) {
	client := pbcache.NewCacheClient(c.client.Conn())
	return client.GetFriendIDListFromCache(ctx, in, opts...)
}

func (c *cacheClient) DelFriendIDListFromCache(ctx context.Context, in *pbcache.DelFriendIDListFromCacheReq, opts ...grpc.CallOption) (*pbcache.DelFriendIDListFromCacheResp, error) {
	client := pbcache.NewCacheClient(c.client.Conn())
	return client.DelFriendIDListFromCache(ctx, in, opts...)
}

func (c *cacheClient) GetBlackIDListFromCache(ctx context.Context, in *pbcache.GetBlackIDListFromCacheReq, opts ...grpc.CallOption) (*pbcache.GetBlackIDListFromCacheResp, error) {
	client := pbcache.NewCacheClient(c.client.Conn())
	return client.GetBlackIDListFromCache(ctx, in, opts...)
}

func (c *cacheClient) DelBlackIDListFromCache(ctx context.Context, in *pbcache.DelBlackIDListFromCacheReq, opts ...grpc.CallOption) (*pbcache.DelBlackIDListFromCacheResp, error) {
	client := pbcache.NewCacheClient(c.client.Conn())
	return client.DelBlackIDListFromCache(ctx, in, opts...)
}

func (c *cacheClient) GetGroupMemberIDListFromCache(ctx context.Context, in *pbcache.GetGroupMemberIDListFromCacheReq, opts ...grpc.CallOption) (*pbcache.GetGroupMemberIDListFromCacheResp, error) {
	client := pbcache.NewCacheClient(c.client.Conn())
	return client.GetGroupMemberIDListFromCache(ctx, in, opts...)
}

func (c *cacheClient) DelGroupMemberIDListFromCache(ctx context.Context, in *pbcache.DelGroupMemberIDListFromCacheReq, opts ...grpc.CallOption) (*pbcache.DelGroupMemberIDListFromCacheResp, error) {
	client := pbcache.NewCacheClient(c.client.Conn())
	return client.DelGroupMemberIDListFromCache(ctx, in, opts...)
}
