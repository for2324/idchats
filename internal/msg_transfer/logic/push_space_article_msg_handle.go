package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbGroup "Open_IM/pkg/proto/group"
	pbOrder "Open_IM/pkg/proto/order"
	"context"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type PushSpaceArticleConsumerHandler struct {
	msgHandle                     map[string]fcb
	pushSpaceArticleConsumerGroup *kfk.MConsumerGroup
}

func (pc *PushSpaceArticleConsumerHandler) Init() {
	pc.msgHandle = make(map[string]fcb)
	pc.msgHandle[config.Config.Kafka.MsgOrder.Topic] = pc.handlePushSpaceArticle
	pc.pushSpaceArticleConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_8_1_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.MsgOrder.Topic},
		config.Config.Kafka.MsgOrder.Addr, config.Config.Kafka.ConsumerGroupID.MsgToOrderArticle)

}

func (pc *PushSpaceArticleConsumerHandler) handlePushSpaceArticle(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession) {
	var rpcErr error
	defer func() {
		if rpcErr == nil {
			sess.MarkMessage(cMsg, "")
		}
	}()
	msg := cMsg.Value
	msgPushMQ := pbOrder.MsgDataToOrderByMQ{}
	err := proto.Unmarshal(msg, &msgPushMQ)
	if err != nil {
		log.NewError(msgPushMQ.OperationID, "msg_transfer Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	log.Debug(msgPushMQ.OperationID, "proto.Unmarshal MsgDataToMQ", msgPushMQ.String())
	if msgPushMQ.Mark != constant.PayMarkPushSpaceArticleType {
		return
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, msgPushMQ.OperationID)

	client := pbGroup.NewGroupClient(etcdConn)
	req := &pbGroup.GlobalPushSpaceArticleReq{
		OperationID:    msgPushMQ.OperationID,
		SpaceArticleId: msgPushMQ.OrderId,
		UserId:         msgPushMQ.UserId,
	}
	_, rpcErr = client.GlobalPushSpaceArticle(context.Background(), req)

}

func (*PushSpaceArticleConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*PushSpaceArticleConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (pc *PushSpaceArticleConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mysql", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			if _, ok := pc.msgHandle[msg.Topic]; ok {
				pc.msgHandle[msg.Topic](msg, string(msg.Key), sess)
			}
		} else {
			log.Error("", "msg get from kafka but is nil", msg.Key)
		}
	}
	return nil
}
