package main

import (
	"Open_IM/internal/msg_gateway/gate"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/interceptors"
	"flag"

	pbgateway "Open_IM/pkg/proto/relay"

	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.ConvertServerConfig(config.Config.ServerConfigs.Gateway)
	// log.NewPrivateLog(constant.LogFileName)
	// defaultWsPorts := config.Config.LongConnSvr.WebsocketPort
	// defaultPromePorts := config.Config.Prometheus.MessageGatewayPrometheusPort
	rpcPort := flag.Int("rpc_port", config.Config.ServerConfigs.Gateway.Port, "rpc listening port")
	wsPort := flag.Int("ws_port", config.Config.LongConnSvr.WebsocketPort[0], "ws listening port")
	// prometheusPort := flag.Int("prometheus_port", defaultPromePorts[0], "PushrometheusPort default listen port")
	flag.Parse()
	// var wg sync.WaitGroup
	// wg.Add(1)
	// fmt.Println("start rpc/msg_gateway server, port: ", *rpcPort, *wsPort, *prometheusPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	// gate.Init(*rpcPort, *wsPort)
	// gate.Run(*prometheusPort)
	// wg.Wait()

	server := gate.NewRPCServer(*rpcPort)
	s := zrpc.MustNewServer(cfg, func(s *grpc.Server) {
		pbgateway.RegisterRelayServer(s, server)
	})
	defer s.Stop()

	ws := gate.NewWebsocketServer(*wsPort)

	server.RegisterLegacyDiscovery()

	s.AddUnaryInterceptors(interceptors.ResponseLogger)

	threading.GoSafe(func() {
		ws.Start()
	})

	s.Start()
}
