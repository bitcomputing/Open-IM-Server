package main

import (
	"Open_IM/internal/rpc/user"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	pbuser "Open_IM/pkg/proto/user"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/prometheus"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := zrpc.RpcServerConf{
		ServiceConf: service.ServiceConf{
			Name: config.Config.ServerConfigs.User.Name,
			Prometheus: prometheus.Config{
				Host: config.Config.ServerConfigs.User.Prometheus.Host,
				Port: config.Config.ServerConfigs.User.Prometheus.Port,
				Path: config.Config.ServerConfigs.User.Prometheus.Path,
			},
		},
		ListenOn: fmt.Sprintf("0.0.0.0:%d", config.Config.ServerConfigs.User.Port),
		Etcd: discov.EtcdConf{
			Hosts: config.Config.ServerConfigs.User.Discovery.Hosts,
			Key:   config.Config.ServerConfigs.User.Discovery.Key,
		},
		Timeout: config.Config.ServerConfigs.User.Timeout,
		Middlewares: zrpc.ServerMiddlewaresConf{
			Trace:      config.Config.ServerConfigs.User.Middlewares.Trace,
			Recover:    config.Config.ServerConfigs.User.Middlewares.Recover,
			Stat:       config.Config.ServerConfigs.User.Middlewares.Stat,
			Prometheus: config.Config.ServerConfigs.User.Middlewares.Prometheus,
			Breaker:    config.Config.ServerConfigs.User.Middlewares.Breaker,
		},
	}

	defaultPorts := config.Config.RpcPort.OpenImUserPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	// prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.UserPrometheusPort[0], "userPrometheusPort default listen port")
	// flag.Parse()
	// fmt.Println("start user rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	// rpcServer := user.NewUserServer(*rpcPort)
	// go func() {
	// 	err := promePkg.StartPromeSrv(*prometheusPort)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	// rpcServer.Run()

	server := user.NewUserServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbuser.RegisterUserServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	s.Start()
}
