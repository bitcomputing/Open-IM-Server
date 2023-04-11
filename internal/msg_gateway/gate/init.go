package gate

import (
	"sync"

	promePkg "Open_IM/pkg/common/prometheus"

	"github.com/go-playground/validator/v10"
)

var (
	rwLock              *sync.RWMutex
	validate            *validator.Validate
	ws                  WServer
	rpcSvr              RPCServer
	sendMsgAllCount     uint64
	sendMsgFailedCount  uint64
	sendMsgSuccessCount uint64
	userCount           uint64

	sendMsgAllCountLock sync.RWMutex
)

func Init(rpcPort, wsPort int) {
	rwLock = new(sync.RWMutex)
	validate = validator.New()

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
