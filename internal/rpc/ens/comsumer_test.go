package ens

import (
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

func TestComsumerOrder(t *testing.T) {
	sConfig := sarama.NewConfig()
	sConfig.Consumer.Offsets.AutoCommit.Enable = true
	// sConfig.Consumer.Retry.Backoff = 2 * time.Second
	ConsumerGroup, err := kfk.NewConsumerGroup(
		"test_group_1",
		sConfig,
	)
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	for {
		err := ConsumerGroup.Consume(ctx, []string{"test_topic"}, &TestEnsOrderConsumerGroupHandle{})
		if err != nil {
			// 警告通知 error
			log.NewError("TestComsumerOrder", "ConsumerGroup.Consume", err.Error())
			time.Sleep(5 * time.Second)
		}
	}
}

type TestEnsOrderConsumerGroupHandle struct{}

func (*TestEnsOrderConsumerGroupHandle) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*TestEnsOrderConsumerGroupHandle) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (och *TestEnsOrderConsumerGroupHandle) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	GenerationID := fmt.Sprint(sess.GenerationID())
	log.NewDebug(GenerationID, "EnsOrderConsumerGroupHandle new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())

	for msg := range claim.Messages() {
		if len(msg.Value) != 0 {
			log.NewDebug(GenerationID, "kafka get info to kafka", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
			if string(msg.Value) == "failed" {
				// 操作失败的话，加入本地消息表进行重试  加入成功标记成功、加入失败（警告通知 error）
				// return localdb.Insert(msg)
				return fmt.Errorf("failed")
			}
			sess.MarkMessage(msg, "")
		} else {
			log.Error(GenerationID, "kafka msg get from kafka but is nil", msg.Key)
		}
	}
	return nil
}

type TestSkipEnsOrderConsumerGroupHandle struct{}

func (*TestSkipEnsOrderConsumerGroupHandle) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*TestSkipEnsOrderConsumerGroupHandle) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (och *TestSkipEnsOrderConsumerGroupHandle) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	GenerationID := fmt.Sprint(sess.GenerationID())
	log.NewDebug(GenerationID, "EnsOrderConsumerGroupHandle new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())

	for msg := range claim.Messages() {
		if len(msg.Value) != 0 {
			log.NewDebug(GenerationID, "kafka get info to kafka", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
			if string(msg.Value) == "failed" {
				sess.MarkMessage(msg, "")
				return nil
			}
		} else {
			log.Error(GenerationID, "kafka msg get from kafka but is nil", msg.Key)
		}
	}
	return nil
}

func TestSkipFaild(t *testing.T) {
	sConfig := sarama.NewConfig()
	sConfig.Consumer.Offsets.AutoCommit.Enable = true
	// sConfig.Consumer.Retry.Backoff = 2 * time.Second
	ConsumerGroup, _ := kfk.NewConsumerGroup(
		"test_group_1",
		sConfig,
	)
	ConsumerGroup.Pause(map[string][]int32{"test_topic": {0}})
	err := ConsumerGroup.Consume(context.Background(), []string{"test_topic"}, &TestSkipEnsOrderConsumerGroupHandle{})
	if err != nil {
		// 警告通知 error
		log.NewError("TestComsumerOrder", "ConsumerGroup.Consume", err.Error())
		return
	}
	ConsumerGroup.Resume(map[string][]int32{"test_topic": {0}})
}
