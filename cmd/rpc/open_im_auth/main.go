package main

import (
	"Open_IM/internal/rpc/auth"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	pbauth "Open_IM/pkg/proto/auth"
	"flag"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.ConvertServerConfig(config.Config.ServerConfigs.Auth)
	rpcPort := flag.Int("port", config.Config.ServerConfigs.Auth.Port, "RpcToken default listen port 10800")
	flag.Parse()

	server := auth.NewRpcAuthServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbauth.RegisterAuthServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	s.Start()
}
