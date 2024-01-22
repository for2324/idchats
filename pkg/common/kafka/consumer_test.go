package kafka

import (
	"Open_IM/pkg/common/log"
	"testing"

	"github.com/Shopify/sarama"
)

type fcb func(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession)

type CgHandler struct{}

func TestConSumer(t *testing.T) {
	addr := []string{"127.0.0.1:9092"}
	// comsumer := NewKafkaConsumer([]string{"127.0.0.1:9092"}, "test")
	topic := "test2"
	groupId := "test2"
	cg := NewMConsumerGroup(
		&MConsumerGroupConfig{
			KafkaVersion:   sarama.V2_8_1_0,
			OffsetsInitial: sarama.OffsetNewest,
			IsReturnErr:    true,
		},
		[]string{topic},
		addr,
		groupId,
	)
	handler := &CgHandler{}
	cg.RegisterHandleAndConsumer(handler)
	live := make(chan int)
	<-live
}

func (*CgHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*CgHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (pc *CgHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		OperationID := "ConsumeClaim"
		log.NewDebug(OperationID, "kafka get info to mysql", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			log.NewInfo(OperationID, msg, string(msg.Key), sess)
		} else {
			log.NewError(OperationID, "msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
