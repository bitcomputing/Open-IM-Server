package logic

import (
	pushclient "Open_IM/internal/push/client"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/kafka"
	promePkg "Open_IM/pkg/common/prometheus"
	"fmt"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

const OnlineTopicBusy = 1
const OnlineTopicVacancy = 0
const Msg = 2
const ConsumerMsgs = 3
const AggregationMessages = 4
const MongoMessages = 5
const ChannelNum = 100

var (
	persistentCH     PersistentConsumerHandler
	historyCH        OnlineHistoryRedisConsumerHandler
	historyMongoCH   OnlineHistoryMongoConsumerHandler
	modifyCH         ModifyMsgConsumerHandler
	producer         *kafka.Producer
	producerToModify *kafka.Producer
	producerToMongo  *kafka.Producer
	cmdCh            chan Cmd2Value
	// onlineTopicStatus int
	// w                 *sync.Mutex
	pushClient pushclient.PushClient
)

func Init() {
	cmdCh = make(chan Cmd2Value, 10000)
	// w = new(sync.Mutex)
	if config.Config.Prometheus.Enable {
		initPrometheus()
	}
	persistentCH.Init()   // ws2mschat save mysql
	historyCH.Init(cmdCh) //
	historyMongoCH.Init()
	modifyCH.Init()
	// onlineTopicStatus = OnlineTopicVacancy
	//offlineHistoryCH.Init(cmdCh)
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
	producerToModify = kafka.NewKafkaProducer(config.Config.Kafka.MsgToModify.Addr, config.Config.Kafka.MsgToModify.Topic)
	producerToMongo = kafka.NewKafkaProducer(config.Config.Kafka.MsgToMongo.Addr, config.Config.Kafka.MsgToMongo.Topic)
	pushClient = pushclient.NewPushClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: config.Config.ClientConfigs.Push.Disconvery.Hosts,
			Key:   config.Config.ClientConfigs.Push.Disconvery.Key,
		},
		Timeout:       config.Config.ClientConfigs.Push.Timeout,
		KeepaliveTime: 0,
		Middlewares: zrpc.ClientMiddlewaresConf{
			Trace:      config.Config.ClientConfigs.Push.Middlewares.Trace,
			Duration:   config.Config.ClientConfigs.Push.Middlewares.Duration,
			Prometheus: config.Config.ClientConfigs.Push.Middlewares.Prometheus,
			Breaker:    config.Config.ClientConfigs.Push.Middlewares.Breaker,
			Timeout:    config.Config.ClientConfigs.Push.Middlewares.Timeout,
		},
	})
}

func Run(promethuesPort int) {
	//register mysqlConsumerHandler to
	if config.Config.ChatPersistenceMysql {
		go persistentCH.persistentConsumerGroup.RegisterHandleAndConsumer(&persistentCH)
	} else {
		fmt.Println("not start mysql consumer")
	}
	go historyCH.historyConsumerGroup.RegisterHandleAndConsumer(&historyCH)
	go historyMongoCH.historyConsumerGroup.RegisterHandleAndConsumer(&historyMongoCH)
	go modifyCH.modifyMsgConsumerGroup.RegisterHandleAndConsumer(&modifyCH)
	//go offlineHistoryCH.historyConsumerGroup.RegisterHandleAndConsumer(&offlineHistoryCH)
	go func() {
		err := promePkg.StartPromeSrv(promethuesPort)
		if err != nil {
			panic(err)
		}
	}()
}

// func SetOnlineTopicStatus(status int) {
// 	w.Lock()
// 	defer w.Unlock()
// 	onlineTopicStatus = status
// }

// func GetOnlineTopicStatus() int {
// 	w.Lock()
// 	defer w.Unlock()
// 	return onlineTopicStatus
// }
