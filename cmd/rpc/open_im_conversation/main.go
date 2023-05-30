package main

import (
	"Open_IM/internal/rpc/conversation"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	pbconversation "Open_IM/pkg/proto/conversation"
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
			Name: config.Config.ServerConfigs.Conversation.Name,
			Prometheus: prometheus.Config{
				Host: config.Config.ServerConfigs.Conversation.Prometheus.Host,
				Port: config.Config.ServerConfigs.Conversation.Prometheus.Port,
				Path: config.Config.ServerConfigs.Conversation.Prometheus.Path,
			},
		},
		ListenOn: fmt.Sprintf("0.0.0.0:%d", config.Config.ServerConfigs.Conversation.Port),
		Etcd: discov.EtcdConf{
			Hosts: config.Config.ServerConfigs.Conversation.Discovery.Hosts,
			Key:   config.Config.ServerConfigs.Conversation.Discovery.Key,
		},
		Timeout: config.Config.ServerConfigs.Conversation.Timeout,
		Middlewares: zrpc.ServerMiddlewaresConf{
			Trace:      config.Config.ServerConfigs.Conversation.Middlewares.Trace,
			Recover:    config.Config.ServerConfigs.Conversation.Middlewares.Recover,
			Stat:       config.Config.ServerConfigs.Conversation.Middlewares.Stat,
			Prometheus: config.Config.ServerConfigs.Conversation.Middlewares.Prometheus,
			Breaker:    config.Config.ServerConfigs.Conversation.Middlewares.Breaker,
		},
	}

	defaultPorts := config.Config.RpcPort.OpenImConversationPort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcConversation default listen port 11300")
	// prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.ConversationPrometheusPort[0], "conversationPrometheusPort default listen port")
	// flag.Parse()
	// fmt.Println("start conversation rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion)
	// rpcServer := rpcConversation.NewRpcConversationServer(*rpcPort)
	// go func() {
	// 	err := promePkg.StartPromeSrv(*prometheusPort)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	// rpcServer.Run()

	server := conversation.NewRpcConversationServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbconversation.RegisterConversationServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	s.Start()
}
