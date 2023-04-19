package gate

import (
	"sync"

	promePkg "Open_IM/pkg/common/prometheus"

	"github.com/go-playground/validator/v10"
)

var (
	rwLock   *sync.RWMutex
	validate *validator.Validate
	ws       WServer
	rpcSvr   RPCServer
)

func Init(rpcPort, wsPort int) {
	rwLock = new(sync.RWMutex)
	validate = validator.New()
	ws.onInit(wsPort)
	rpcSvr.onInit(rpcPort)
	initPrometheus()
}

func Run(promethuesPort int) {
	go ws.run()
	go rpcSvr.run()
	go func() {
		err := promePkg.StartPromeSrv(promethuesPort)
		if err != nil {
			panic(err)
		}
	}()
}
