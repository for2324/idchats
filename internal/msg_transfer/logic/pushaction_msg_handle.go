package logic

import (
	"Open_IM/internal/rpc/user"
	"Open_IM/pkg/common/config"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"context"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type PushActionConsumerHandler struct {
	msgHandle               map[string]fcb
	pushActionConsumerGroup *kfk.MConsumerGroup
}

func (pc *PushActionConsumerHandler) Init() {
	pc.msgHandle = make(map[string]fcb)
	pc.msgHandle[config.Config.Kafka.AnnouncementAction.Topic] = pc.handlePushAction
	pc.pushActionConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_8_1_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.AnnouncementAction.Topic},
		config.Config.Kafka.AnnouncementAction.Addr, config.Config.Kafka.ConsumerGroupID.MsgToAnnounce)

}

// 消息为什么没有处理？？？
func (pc *PushActionConsumerHandler) handlePushAction(cMsg *sarama.ConsumerMessage, msgKey string, _ sarama.ConsumerGroupSession) {
	msg := cMsg.Value
	msgPushMQ := pbMsg.NewPushActionMsgMq{}
	err := proto.Unmarshal(msg, &msgPushMQ)
	if err != nil {
		log.NewError(msgPushMQ.OperationID, "PushActionConsumerHandler>>>>>msg_transfer Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	log.Debug(msgPushMQ.OperationID, "PushActionConsumerHandler>>>>>解析推送数据如下：", msgPushMQ.String())
	userService := user.NewUserServer(0)
	_, err = userService.RpcPushMessageToFollowsUser(context.Background(), &sdk_ws.PushMessageToMailFromUserToFans{
		OperationID:       msgPushMQ.OperationID,
		ContentType:       msgPushMQ.PushMsg.ContentType,
		ArticleID:         msgPushMQ.PushMsg.ArticleID,
		FromUserID:        msgPushMQ.PushMsg.FromUserID,
		FromArticleAuthor: msgPushMQ.PushMsg.FromArticleAuthor,
		IsGlobal:          msgPushMQ.PushMsg.IsGlobal,
	})
	if err != nil {
		log.NewError(msgPushMQ.OperationID, "推送信息返回错误")
		return
	}
}

func (*PushActionConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*PushActionConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (pc *PushActionConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to PushActionConsumerHandler", "msgTopic",
			msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			if _, ok := pc.msgHandle[msg.Topic]; ok {
				pc.msgHandle[msg.Topic](msg, string(msg.Key), sess)
			}

		} else {
			log.Error("", "msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
