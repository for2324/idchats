package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/kafka"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/statistics"
	"fmt"
	"sync"
)

const OnlineTopicBusy = 1
const OnlineTopicVacancy = 0
const Msg = 2
const ConsumerMsgs = 3
const AggregationMessages = 4
const MongoMessages = 5
const ChannelNum = 100

var (
	persistentCH                    PersistentConsumerHandler
	historyCH                       OnlineHistoryRedisConsumerHandler
	historyMongoCH                  OnlineHistoryMongoConsumerHandler
	likeactionconsumberhandler      LikesActionConsumerHandler
	pushSpaceArticleConsumerHandler PushSpaceArticleConsumerHandler
	orderCallBackConsumerHandler    OrderCallBackConsumerHandler
	announcementHandler             PushActionConsumerHandler
	producer                        *kafka.Producer
	producerToMongo                 *kafka.Producer
	producerLikesAction             *kafka.Producer
	cmdCh                           chan Cmd2Value
	onlineTopicStatus               int
	w                               *sync.Mutex
	singleMsgSuccessCount           uint64
	groupMsgCount                   uint64
	singleMsgFailedCount            uint64
	singleMsgSuccessCountMutex      sync.Mutex
)

func Init() {
	cmdCh = make(chan Cmd2Value, 10000)
	w = new(sync.Mutex)
	if config.Config.Prometheus.Enable {
		initPrometheus()
	}
	persistentCH.Init()   // ws2mschat save mysql
	historyCH.Init(cmdCh) //  可能没什么用到的地方
	historyMongoCH.Init()
	likeactionconsumberhandler.Init()
	pushSpaceArticleConsumerHandler.Init()
	orderCallBackConsumerHandler.Init()
	announcementHandler.Init()
	onlineTopicStatus = OnlineTopicVacancy
	//offlineHistoryCH.Init(cmdCh)
	statistics.NewStatistics(&singleMsgSuccessCount, config.Config.ModuleName.MsgTransferName, fmt.Sprintf("%d second singleMsgCount insert to mongo", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	statistics.NewStatistics(&groupMsgCount, config.Config.ModuleName.MsgTransferName, fmt.Sprintf("%d second groupMsgCount insert to mongo", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
	producerToMongo = kafka.NewKafkaProducer(config.Config.Kafka.MsgToMongo.Addr, config.Config.Kafka.MsgToMongo.Topic)
	producerLikesAction = kafka.NewKafkaProducer(config.Config.Kafka.LikesAction.Addr, config.Config.Kafka.LikesAction.Topic)

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
	//go offlineHistoryCH.historyConsumerGroup.RegisterHandleAndConsumer(&offlineHistoryCH)
	go likeactionconsumberhandler.likeActionConsumerGroup.RegisterHandleAndConsumer(&likeactionconsumberhandler)
	go announcementHandler.pushActionConsumerGroup.RegisterHandleAndConsumer(&announcementHandler)
	go orderCallBackConsumerHandler.orderCallBackConsumerGroup.RegisterHandleAndConsumer(&orderCallBackConsumerHandler)
	go func() {
		err := promePkg.StartPromeSrv(promethuesPort)
		if err != nil {
			panic(err)
		}
	}()
}
func SetOnlineTopicStatus(status int) {
	w.Lock()
	defer w.Unlock()
	onlineTopicStatus = status
}
func GetOnlineTopicStatus() int {
	w.Lock()
	defer w.Unlock()
	return onlineTopicStatus
}

func initPrometheus() {
	promePkg.NewSeqGetSuccessCounter()
	promePkg.NewSeqGetFailedCounter()
	promePkg.NewSeqSetSuccessCounter()
	promePkg.NewSeqSetFailedCounter()
	promePkg.NewMsgInsertRedisSuccessCounter()
	promePkg.NewMsgInsertRedisFailedCounter()
	promePkg.NewMsgInsertMongoSuccessCounter()
	promePkg.NewMsgInsertMongoFailedCounter()
}
