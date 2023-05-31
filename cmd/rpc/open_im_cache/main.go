package main

import (
	"Open_IM/internal/rpc/cache"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	pbcache "Open_IM/pkg/proto/cache"
	"flag"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.ConvertServerConfig(config.Config.ServerConfigs.Cache)
	rpcPort := flag.Int("port", config.Config.ServerConfigs.Cache.Port, "RpcToken default listen port 10800")
	flag.Parse()
	server := cache.NewCacheServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbcache.RegisterCacheServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	s.Start()
}
