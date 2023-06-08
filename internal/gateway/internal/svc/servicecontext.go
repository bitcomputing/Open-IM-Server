package svc

import (
	"Open_IM/internal/gateway/internal/config"
	"Open_IM/internal/gateway/internal/middleware"
	auth "Open_IM/internal/rpc/auth/client"
	cache "Open_IM/internal/rpc/cache/client"
	conversation "Open_IM/internal/rpc/conversation/client"
	friend "Open_IM/internal/rpc/friend/client"
	group "Open_IM/internal/rpc/group/client"
	message "Open_IM/internal/rpc/msg/client"
	user "Open_IM/internal/rpc/user/client"
	commoncfg "Open_IM/pkg/common/config"
	"Open_IM/pkg/discovery"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config               config.Config
	CorsMiddleware       rest.Middleware
	BodyLoggerMiddleware rest.Middleware
	AuthClient           auth.AuthClient
	UserClient           user.UserClient
	ConversationClient   conversation.ConversationClient
	FriendClient         friend.FriendClient
	GroupClient          group.GroupClient
	MessageClient        message.MsgClient
	CacheClient          cache.CacheClient
	GatewayClient        *discovery.Client
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	gatewayClient, err := discovery.NewClient(commoncfg.ConvertClientConfig(commoncfg.Config.ClientConfigs.Gateway))
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Config:               c,
		CorsMiddleware:       middleware.NewCorsMiddleware().Handle,
		BodyLoggerMiddleware: middleware.NewBodyLoggerMiddleware().Handle,
		AuthClient:           auth.NewAuthClient(c.AuthClient),
		UserClient:           user.NewUserClient(c.UserClient),
		ConversationClient:   conversation.NewConversationClient(c.ConversationClient),
		FriendClient:         friend.NewFriendClient(c.FriendClient),
		GroupClient:          group.NewGroupClient(c.GroupClient),
		MessageClient:        message.NewMsgClient(c.MessageClient),
		CacheClient:          cache.NewCacheClient(c.CacheClient),
		GatewayClient:        gatewayClient,
	}, nil
}
