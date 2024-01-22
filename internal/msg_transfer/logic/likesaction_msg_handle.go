package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	"context"
	"github.com/Shopify/sarama"
	"github.com/go-redsync/redsync/v4"
	"github.com/golang/protobuf/proto"
	"time"
)

type LikesActionConsumerHandler struct {
	msgHandle               map[string]fcb
	likeActionConsumerGroup *kfk.MConsumerGroup
}

func (pc *LikesActionConsumerHandler) Init() {
	pc.msgHandle = make(map[string]fcb)
	pc.msgHandle[config.Config.Kafka.LikesAction.Topic] = pc.handleLikeAction
	pc.likeActionConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_8_1_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.LikesAction.Topic},
		config.Config.Kafka.LikesAction.Addr, config.Config.Kafka.ConsumerGroupID.MsgToLike)

}

func (pc *LikesActionConsumerHandler) handleLikeAction(cMsg *sarama.ConsumerMessage, msgKey string, _ sarama.ConsumerGroupSession) {
	msg := cMsg.Value
	msgLikesActionMQ := pbMsg.MsgLikeMQ{}
	err := proto.Unmarshal(msg, &msgLikesActionMQ)
	if err != nil {
		log.NewError(msgLikesActionMQ.OperationID, "msg_transfer Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	log.Debug(msgLikesActionMQ.OperationID, "proto.Unmarshal MsgDataToMQ", msgLikesActionMQ.String())
	//消息处理函数
	switch msgLikesActionMQ.LikeReword.ContentType {
	case "nft":
		mutexname := "nftlike:" + msgLikesActionMQ.LikeReword.ArticleID
		rs := db.DB.Pool
		mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
		ctx := context.Background()
		if err := mutex.LockContext(ctx); err != nil {
			log.Error("===================handleLikeAction==========================")
			return
		}
		defer mutex.UnlockContext(ctx)
		if imdb.LikeActionNftCount(utils.StringToInt64(msgLikesActionMQ.LikeReword.ArticleID), msgLikesActionMQ.LikeReword.UserID, msgLikesActionMQ.LikeReword.Action) != nil {
			log.Error("===================handleLikeAction==========================")
		}
	case "announcement":
		mutexname := "announcement:" + msgLikesActionMQ.LikeReword.ArticleID
		rs := db.DB.Pool
		mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
		ctx := context.Background()
		if err := mutex.LockContext(ctx); err != nil {
			log.Error("===================handleLikeAction==========================")
			return
		}
		defer mutex.UnlockContext(ctx)
		if imdb.LikeActionAnnouncementCount(
			utils.StringToInt64(msgLikesActionMQ.LikeReword.ArticleID),
			msgLikesActionMQ.LikeReword.UserID, msgLikesActionMQ.LikeReword.Action) != nil {
			log.Error("===================handleLikeAction==========================")
		}
	}

}
func (*LikesActionConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*LikesActionConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (pc *LikesActionConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mysql", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			pc.msgHandle[msg.Topic](msg, string(msg.Key), sess)
		} else {
			log.Error("", "msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
