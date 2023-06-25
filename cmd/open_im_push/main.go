package main

import (
	"Open_IM/internal/push"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	pbPush "Open_IM/pkg/proto/push"
	"flag"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.ConvertServerConfig(config.Config.ServerConfigs.Push)
	rpcPort := flag.Int("port", config.Config.ServerConfigs.Push.Port, "rpc listening port")
	flag.Parse()

	server := push.NewPushServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbPush.RegisterPushMsgServiceServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	push.Init(server.CacheClient)
	push.Run()

	s.Start()
}
