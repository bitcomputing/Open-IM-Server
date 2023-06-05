package discovery

import (
	"sync"
	"time"

	clientInterceptors "Open_IM/pkg/common/interceptors/client"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	subscriber         *discov.Subscriber
	conns              sync.Map
	clientInterceptors []grpc.UnaryClientInterceptor
}

func NewClient(cfg zrpc.RpcClientConf) (*Client, error) {
	sub, err := discov.NewSubscriber(cfg.Etcd.Hosts, cfg.Etcd.Key, discov.Exclusive())
	if err != nil {
		return nil, err
	}

	c := &Client{
		subscriber: sub,
	}

	c.buildInterceptors(cfg)

	c.subscriber.AddListener(c.Update)

	return c, nil
}

func (c *Client) ClientConns() []*grpc.ClientConn {
	clientConns := make([]*grpc.ClientConn, 0)
	c.conns.Range(func(key, value any) bool {
		conn := value.(*grpc.ClientConn)
		clientConns = append(clientConns, conn)
		return true
	})

	return clientConns
}

func (c *Client) Update() {
	endpoints := lo.Associate(c.subscriber.Values(), func(endpoint string) (string, struct{}) {
		return endpoint, struct{}{}
	})

	for k := range endpoints {
		_, ok := c.conns.Load(k)
		if !ok {
			conn, err := grpc.Dial(k, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(c.clientInterceptors...))
			if err != nil {
				logx.Error(err)
				continue
			}
			c.conns.Store(k, conn)
		}
	}

	c.conns.Range(func(key, value any) bool {
		k := key.(string)
		_, ok := endpoints[k]
		if !ok {
			c.conns.Delete(k)
		}

		return true
	})
}

func (c *Client) buildInterceptors(cfg zrpc.RpcClientConf) {
	if cfg.Middlewares.Breaker {
		c.clientInterceptors = append(c.clientInterceptors, clientInterceptors.BreakerInterceptor)
	}

	if cfg.Middlewares.Duration {
		c.clientInterceptors = append(c.clientInterceptors, clientInterceptors.DurationInterceptor)
	}

	if cfg.Middlewares.Prometheus {
		c.clientInterceptors = append(c.clientInterceptors, clientInterceptors.PrometheusInterceptor)
	}

	if cfg.Middlewares.Trace {
		c.clientInterceptors = append(c.clientInterceptors, clientInterceptors.UnaryTracingInterceptor)
	}

	if cfg.Middlewares.Timeout {
		c.clientInterceptors = append(c.clientInterceptors, clientInterceptors.TimeoutInterceptor(time.Duration(cfg.Timeout)*time.Millisecond))
	}
}
