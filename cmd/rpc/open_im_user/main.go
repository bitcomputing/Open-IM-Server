package main

import (
	"Open_IM/internal/rpc/user"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	pbuser "Open_IM/pkg/proto/user"
	"flag"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.ConvertServerConfig(config.Config.ServerConfigs.User)

	rpcPort := flag.Int("port", config.Config.ServerConfigs.User.Port, "rpc listening port")
	flag.Parse()

	server := user.NewUserServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbuser.RegisterUserServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	s.Start()
}
