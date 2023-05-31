package main

import (
	"Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	"flag"

	pbmsg "Open_IM/pkg/proto/msg"

	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.ConvertServerConfig(config.Config.ServerConfigs.Message)
	rpcPort := flag.Int("port", config.Config.ServerConfigs.Message.Port, "rpc listening port")
	flag.Parse()
	// prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.MessagePrometheusPort[0], "msgPrometheusPort default listen port")
	// flag.Parse()
	// fmt.Println("start msg rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	// rpcServer := msg.NewRpcChatServer(*rpcPort)
	// go func() {
	// 	err := promePkg.StartPromeSrv(*prometheusPort)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	// rpcServer.Run()
	server := msg.NewRpcChatServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbmsg.RegisterMsgServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	threading.GoSafe(func() {
		server.RunCh()
	})

	s.Start()
}
