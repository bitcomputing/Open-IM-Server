package main

import (
	"Open_IM/internal/rpc/auth"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	pbauth "Open_IM/pkg/proto/auth"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	// defaultPorts := config.Config.RpcPort.OpenImAuthPort
	// rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 10800")
	// prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.AuthPrometheusPort[0], "authPrometheusPort default listen port")
	// flag.Parse()
	// fmt.Println("start auth rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	// rpcServer := rpcAuth.NewRpcAuthServer(*rpcPort)
	// go func() {
	// 	err := promePkg.StartPromeSrv(*prometheusPort)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	// rpcServer.Run()

	cfg := zrpc.RpcServerConf{
		ServiceConf: service.ServiceConf{
			Name: config.Config.ServerConfigs.Auth.Name,
		},
		ListenOn: fmt.Sprintf("0.0.0.0:%d", config.Config.ServerConfigs.Auth.Port),
		Etcd: discov.EtcdConf{
			Hosts: config.Config.ServerConfigs.Auth.Discovery.Hosts,
			Key:   config.Config.ServerConfigs.Auth.Discovery.Key,
		},
		Timeout: config.Config.ServerConfigs.Auth.Timeout,
		Middlewares: zrpc.ServerMiddlewaresConf{
			Trace:      config.Config.ServerConfigs.Auth.Middlewares.Trace,
			Recover:    config.Config.ServerConfigs.Auth.Middlewares.Recover,
			Stat:       config.Config.ServerConfigs.Auth.Middlewares.Stat,
			Prometheus: config.Config.ServerConfigs.Auth.Middlewares.Prometheus,
			Breaker:    config.Config.ServerConfigs.Auth.Middlewares.Breaker,
		},
	}

	defaultPorts := config.Config.RpcPort.OpenImAuthPort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 10800")

	server := auth.NewRpcAuthServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbauth.RegisterAuthServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	s.Start()
}
