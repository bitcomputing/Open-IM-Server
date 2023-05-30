/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/22 15:33).
 */
package push

import (
	pusher "Open_IM/internal/push/providers"
	fcm "Open_IM/internal/push/providers/fcm"
	"Open_IM/internal/push/providers/getui"
	jpush "Open_IM/internal/push/providers/jpush"
	"Open_IM/internal/push/providers/mobpush"
	"Open_IM/pkg/common/config"

	"github.com/zeromicro/go-zero/core/metric"
)

var (
	pushCh PushConsumerHandler
	// producer                  *kafka.Producer
	offlinePusher             pusher.OfflinePusher
	offlinePushSuccessCounter metric.CounterVec
	offlinePushFailureCounter metric.CounterVec
)

func Init() {
	pushCh.Init()
	offlinePushSuccessCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Name: "msg_offline_push_success",
		Help: "The number of msg successful offline pushed",
	})
	offlinePushFailureCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Name: "msg_offline_push_failed",
		Help: "The number of msg failed offline pushed",
	})
}

func init() {
	// producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
	if *config.Config.Push.Getui.Enable {
		offlinePusher = getui.GetuiClient
	}
	if config.Config.Push.Jpns.Enable {
		offlinePusher = jpush.JPushClient
	}

	if config.Config.Push.Fcm.Enable {
		offlinePusher = fcm.NewFcm()
	}

	if config.Config.Push.Mob.Enable {
		offlinePusher = mobpush.MobPushClient
	}
}

func Run() {
	go pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&pushCh)
}
