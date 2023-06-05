package gate

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	rwLock   *sync.RWMutex
	validate *validator.Validate
	ws       *WServer
	server   *RPCServer
)

func init() {
	rwLock = new(sync.RWMutex)
	validate = validator.New()
}

// func Run(promethuesPort int) {
// 	go ws.run()
// 	go rpcSvr.run()
// 	go func() {
// 		err := promePkg.StartPromeSrv(promethuesPort)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()
// }
