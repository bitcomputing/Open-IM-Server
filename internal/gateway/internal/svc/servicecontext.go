package svc

import (
	"Open_IM/internal/gateway/internal/config"
	"Open_IM/internal/gateway/internal/middleware"
	auth "Open_IM/internal/rpc/auth/client"
	conversation "Open_IM/internal/rpc/conversation/client"
	friend "Open_IM/internal/rpc/friend/client"
	group "Open_IM/internal/rpc/group/client"
	user "Open_IM/internal/rpc/user/client"

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
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:               c,
		CorsMiddleware:       middleware.NewCorsMiddleware().Handle,
		BodyLoggerMiddleware: middleware.NewBodyLoggerMiddleware().Handle,
		AuthClient:           auth.NewAuthClient(c.AuthClient),
		UserClient:           user.NewUserClient(c.UserClient),
		ConversationClient:   conversation.NewConversationClient(c.ConversationClient),
		FriendClient:         friend.NewFriendClient(c.FriendClient),
		GroupClient:          group.NewGroupClient(c.GroupClient),
	}
}
