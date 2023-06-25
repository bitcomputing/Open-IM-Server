package main

import (
	"Open_IM/internal/rpc/friend"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	pbfriend "Open_IM/pkg/proto/friend"
	"flag"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.ConvertServerConfig(config.Config.ServerConfigs.Friend)
	rpcPort := flag.Int("port", config.Config.ServerConfigs.Friend.Port, "get RpcFriendPort from cmd,default 12000 as port")
	flag.Parse()
	// prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.FriendPrometheusPort[0], "friendPrometheusPort default listen port")
	// flag.Parse()
	// fmt.Println("start friend rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	// rpcServer := friend.NewFriendServer(*rpcPort)
	// go func() {
	// 	err := promePkg.StartPromeSrv(*prometheusPort)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	// rpcServer.Run()
	server := friend.NewFriendServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbfriend.RegisterFriendServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	s.Start()
}
