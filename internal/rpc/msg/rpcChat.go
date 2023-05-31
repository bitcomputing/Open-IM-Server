package msg

import (
	pushclient "Open_IM/internal/push/client"
	cacheclient "Open_IM/internal/rpc/cache/client"
	conversationclient "Open_IM/internal/rpc/conversation/client"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/kafka"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"strings"

	"Open_IM/pkg/proto/msg"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metric"
	"google.golang.org/protobuf/proto"
)

type MessageWriter interface {
	SendMessage(m proto.Message, key string, operationID string) (int32, int64, error)
}

type rpcChat struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	messageWriter   MessageWriter
	//offlineProducer *kafka.Producer
	delMsgCh       chan deleteMsg
	dMessageLocker MessageLocker
	msg.UnimplementedMsgServer
	conversationClient                  conversationclient.ConversationClient
	pushClient                          pushclient.PushClient
	cacheClient                         cacheclient.CacheClient
	pullMessageRedisSuccess             metric.CounterVec
	pullMessageRedisFailure             metric.CounterVec
	pullMessageMongoSuccess             metric.CounterVec
	pullMessageMongoFailure             metric.CounterVec
	singleChatMsgRecvSuccess            metric.CounterVec
	groupChatMsgRecvSuccess             metric.CounterVec
	workSuperGroupChatMsgRecvSuccess    metric.CounterVec
	singleChatMsgProcessSuccess         metric.CounterVec
	singleChatMsgProcessFailure         metric.CounterVec
	groupChatMsgProcessSuccess          metric.CounterVec
	groupChatMsgProcessFailure          metric.CounterVec
	workSuperGroupChatMsgProcessSuccess metric.CounterVec
	workSuperGroupChatMsgProcessFailure metric.CounterVec
}

type deleteMsg struct {
	UserID      string
	OpUserID    string
	SeqList     []uint32
	OperationID string
}

func NewRpcChatServer(port int) *rpcChat {
	// log.NewPrivateLog(constant.LogFileName)
	rc := rpcChat{
		rpcPort:            port,
		rpcRegisterName:    config.Config.RpcRegisterName.OpenImMsgName,
		etcdSchema:         config.Config.Etcd.EtcdSchema,
		etcdAddr:           config.Config.Etcd.EtcdAddr,
		dMessageLocker:     NewLockerMessage(),
		conversationClient: conversationclient.NewConversationClient(config.ConvertClientConfig(config.Config.ClientConfigs.Conversation)),
		pushClient:         pushclient.NewPushClient(config.ConvertClientConfig(config.Config.ClientConfigs.Push)),
		cacheClient:        cacheclient.NewCacheClient(config.ConvertClientConfig(config.Config.ClientConfigs.Cache)),
		pullMessageRedisSuccess: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "msg_pull_from_redis_success",
			Help: "The number of successful pull msg from redis",
		}),
		pullMessageRedisFailure: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "msg_pull_from_redis_failed",
			Help: "The number of failed pull msg from redis",
		}),
		pullMessageMongoSuccess: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "msg_pull_from_mongo_success",
			Help: "The number of successful pull msg from mongo",
		}),
		pullMessageMongoFailure: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "msg_pull_from_mongo_failed",
			Help: "The number of failed pull msg from mongo",
		}),
		singleChatMsgRecvSuccess: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "single_chat_msg_recv_success",
			Help: "The number of single chat msg successful received ",
		}),
		groupChatMsgRecvSuccess: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "group_chat_msg_recv_success",
			Help: "The number of group chat msg successful received",
		}),
		workSuperGroupChatMsgRecvSuccess: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "work_super_group_chat_msg_recv_success",
			Help: "The number of work/super group chat msg successful received",
		}),
		singleChatMsgProcessSuccess: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "single_chat_msg_process_success",
			Help: "The number of single chat msg successful processed",
		}),
		singleChatMsgProcessFailure: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "single_chat_msg_process_failed",
			Help: "The number of single chat msg failed processed",
		}),
		groupChatMsgProcessSuccess: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "group_chat_msg_process_success",
			Help: "The number of group chat msg successful processed",
		}),
		groupChatMsgProcessFailure: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "group_chat_msg_process_failed",
			Help: "The number of group chat msg failed processed",
		}),
		workSuperGroupChatMsgProcessSuccess: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "work_super_group_chat_msg_process_success",
			Help: "The number of work/super group chat msg successful processed",
		}),
		workSuperGroupChatMsgProcessFailure: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "work_super_group_chat_msg_process_failed",
			Help: "The number of work/super group chat msg failed processed",
		}),
	}
	rc.messageWriter = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
	//rc.offlineProducer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschatOffline.Addr, config.Config.Kafka.Ws2mschatOffline.Topic)
	rc.delMsgCh = make(chan deleteMsg, 1000)
	return &rc
}

func (s *rpcChat) RegisterLegacyDiscovery() {
	rpcRegisterIP, err := utils.GetLocalIP()
	if err != nil {
		panic(fmt.Errorf("GetLocalIP failed: %w", err))
	}

	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		logx.Error("RegisterEtcd failed ", err.Error(),
			s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register auth module  rpc to etcd err"))
	}
}

// func (rpc *rpcChat) Run() {
// 	log.Info("", "rpcChat init...")
// 	listenIP := ""
// 	if config.Config.ListenIP == "" {
// 		listenIP = "0.0.0.0"
// 	} else {
// 		listenIP = config.Config.ListenIP
// 	}
// 	address := listenIP + ":" + strconv.Itoa(rpc.rpcPort)
// 	listener, err := net.Listen("tcp", address)
// 	if err != nil {
// 		panic("listening err:" + err.Error() + rpc.rpcRegisterName)
// 	}
// 	log.Info("", "listen network success, address ", address)
// 	recvSize := 1024 * 1024 * 30
// 	sendSize := 1024 * 1024 * 30
// 	var grpcOpts = []grpc.ServerOption{
// 		grpc.MaxRecvMsgSize(recvSize),
// 		grpc.MaxSendMsgSize(sendSize),
// 	}
// 	if config.Config.Prometheus.Enable {
// 		promePkg.NewGrpcRequestCounter()
// 		promePkg.NewGrpcRequestFailedCounter()
// 		promePkg.NewGrpcRequestSuccessCounter()
// 		grpcOpts = append(grpcOpts, []grpc.ServerOption{
// 			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
// 			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
// 			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
// 		}...)
// 	}
// 	srv := grpc.NewServer(grpcOpts...)
// 	defer srv.GracefulStop()

// 	rpcRegisterIP := config.Config.RpcRegisterIP
// 	msg.RegisterMsgServer(srv, rpc)
// 	if config.Config.RpcRegisterIP == "" {
// 		rpcRegisterIP, err = utils.GetLocalIP()
// 		if err != nil {
// 			log.Error("", "GetLocalIP failed ", err.Error())
// 		}
// 	}
// 	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
// 	if err != nil {
// 		log.Error("", "register rpcChat to etcd failed ", err.Error())
// 		panic(utils.Wrap(err, "register chat module  rpc to etcd err"))
// 	}
// 	go rpc.runCh()
// 	rpc.initPrometheus()
// 	err = srv.Serve(listener)
// 	if err != nil {
// 		log.Error("", "rpc rpcChat failed ", err.Error())
// 		return
// 	}
// 	log.Info("", "rpc rpcChat init success")
// }

func (rpc *rpcChat) RunCh() {
	ctx := context.Background()
	logx.Info("start del msg chan ")
	for {
		msg := <-rpc.delMsgCh
		logx.WithContext(ctx).WithFields(logx.Field("op", msg.OperationID)).Info("delmsgch recv new: ", msg)
		if len(msg.SeqList) > 0 {
			db.DB.DelMsgFromCache(ctx, msg.UserID, msg.SeqList, msg.OperationID)
			DeleteMessageNotification(ctx, msg.OpUserID, msg.UserID, msg.SeqList, msg.OperationID)
		}
	}
}
