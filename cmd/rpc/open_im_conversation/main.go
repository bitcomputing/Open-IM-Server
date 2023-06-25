package main

import (
	"Open_IM/internal/rpc/conversation"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	pbconversation "Open_IM/pkg/proto/conversation"
	"flag"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.ConvertServerConfig(config.Config.ServerConfigs.Conversation)
	rpcPort := flag.Int("port", config.Config.ServerConfigs.Conversation.Port, "RpcConversation default listen port 11300")
	flag.Parse()

	server := conversation.NewRpcConversationServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbconversation.RegisterConversationServer(s, server)
	})
	defer s.Stop()

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	s.Start()
}
