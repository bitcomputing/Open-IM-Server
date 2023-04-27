package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	AuthClient         zrpc.RpcClientConf
	UserClient         zrpc.RpcClientConf
	ConversationClient zrpc.RpcClientConf
	FriendClient       zrpc.RpcClientConf
}
