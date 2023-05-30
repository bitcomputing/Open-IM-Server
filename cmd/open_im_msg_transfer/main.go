package main

import (
	"Open_IM/internal/msg_transfer/logic"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.MessageTransferPrometheusPort[0], "MessageTransferPrometheusPort default listen port")
	flag.Parse()
	logx.MustSetup(logx.LogConf{
		Mode:       "console",
		Encoding:   "json",
		TimeFormat: time.RFC3339,
		Level:      "debug",
	})
	// log.NewPrivateLog(constant.LogFileName)
	logic.Init()
	fmt.Println("start msg_transfer server ", ", OpenIM version: ", constant.CurrentVersion)
	logic.Run(*prometheusPort)
	wg.Wait()
}
