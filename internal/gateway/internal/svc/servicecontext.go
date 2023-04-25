package svc

import (
	"Open_IM/internal/gateway/internal/config"
	"Open_IM/internal/gateway/internal/middleware"
	auth "Open_IM/internal/rpc/auth/client"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config               config.Config
	CorsMiddleware       rest.Middleware
	BodyLoggerMiddleware rest.Middleware
	AuthClient           auth.AuthClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:               c,
		CorsMiddleware:       middleware.NewCorsMiddleware().Handle,
		BodyLoggerMiddleware: middleware.NewBodyLoggerMiddleware().Handle,
		AuthClient:           auth.NewAuthClient(c.AuthClient),
	}
}
