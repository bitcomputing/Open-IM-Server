package main

import (
	"Open_IM/internal/push"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	pbPush "Open_IM/pkg/proto/push"
	"flag"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/prometheus"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := zrpc.RpcServerConf{
		ServiceConf: service.ServiceConf{
			Name: config.Config.ServerConfigs.Push.Name,
			Log: logx.LogConf{
				Mode:       "console",
				Encoding:   "json",
				TimeFormat: time.RFC3339,
				Level:      "debug",
			},
			Prometheus: prometheus.Config{
				Host: config.Config.ServerConfigs.Push.Prometheus.Host,
				Port: config.Config.ServerConfigs.Push.Prometheus.Port,
				Path: config.Config.ServerConfigs.Push.Prometheus.Path,
			},
		},
		ListenOn: fmt.Sprintf("0.0.0.0:%d", config.Config.ServerConfigs.Push.Port),
		Etcd: discov.EtcdConf{
			Hosts: config.Config.ServerConfigs.Push.Discovery.Hosts,
			Key:   config.Config.ServerConfigs.Push.Discovery.Key,
		},
		Timeout: config.Config.ServerConfigs.Push.Timeout,
		Middlewares: zrpc.ServerMiddlewaresConf{
			Trace:      config.Config.ServerConfigs.Auth.Middlewares.Trace,
			Recover:    config.Config.ServerConfigs.Auth.Middlewares.Recover,
			Stat:       config.Config.ServerConfigs.Auth.Middlewares.Stat,
			Prometheus: config.Config.ServerConfigs.Auth.Middlewares.Prometheus,
			Breaker:    config.Config.ServerConfigs.Auth.Middlewares.Breaker,
		},
	}

	defaultPorts := config.Config.RpcPort.OpenImPushPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")

	server := push.NewPushServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbPush.RegisterPushMsgServiceServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	push.Init()
	push.Run()

	s.Start()
}
