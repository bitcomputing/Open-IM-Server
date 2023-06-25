package main

import (
	"Open_IM/internal/rpc/group"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	"flag"

	pbgroup "Open_IM/pkg/proto/group"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.ConvertServerConfig(config.Config.ServerConfigs.Group)
	rpcPort := flag.Int("port", config.Config.ServerConfigs.Group.Port, "get RpcGroupPort from cmd,default 16000 as port")
	flag.Parse()
	// prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.GroupPrometheusPort[0], "groupPrometheusPort default listen port")
	// flag.Parse()
	// fmt.Println("start group rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	// rpcServer := group.NewGroupServer(*rpcPort)
	// go func() {
	// 	err := promePkg.StartPromeSrv(*prometheusPort)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	// rpcServer.Run()

	server := group.NewGroupServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbgroup.RegisterGroupServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	s.Start()
}
