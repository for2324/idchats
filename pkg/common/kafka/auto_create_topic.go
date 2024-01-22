package kafka

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"errors"
	"reflect"

	"github.com/Shopify/sarama"
)

func init() {
	log.NewInfo("AutoCreateTopic", "AutoCreateTopic init")
	AutoCreateTopic()
}

func AutoCreateTopic() {
	saConfig := sarama.NewConfig()
	saConfig.Version = sarama.V2_5_0_0

	OperationID := "auto_create_topic init"

	addr := config.Config.Kafka.Addr
	Partitions := config.Config.Kafka.Partitions
	if len(addr) == 0 {
		addr = []string{"127.0.0.1:9092"}
	}
	if Partitions == 0 {
		Partitions = 8
	}
	ReplicationFactor := config.Config.Kafka.ReplicationFactor
	if ReplicationFactor == 0 {
		ReplicationFactor = 1
	}

	admin, err := sarama.NewClusterAdmin(addr, saConfig)
	if err != nil {
		panic("Error creating cluster admin: " + err.Error())
	}
	defer admin.Close()

	// 使用反射获取 Kafka 结构体中的所有字段
	topics := []string{}
	kafkaConf := config.Config.Kafka
	kafkaValue := reflect.ValueOf(kafkaConf)
	for i := 0; i < kafkaValue.NumField(); i++ {
		field := kafkaValue.Field(i)
		if field.Type().Kind() == reflect.Struct {
			val := field.FieldByName("Topic")
			if val.IsValid() {
				topics = append(topics, val.String())
			}
		}
	}

	log.NewInfo(OperationID, "topics = %v", topics)
	detail := &sarama.TopicDetail{
		NumPartitions:     Partitions,
		ReplicationFactor: ReplicationFactor,
	}
	for _, topic := range topics {
		err = admin.CreateTopic(topic, detail, false)
		if err != nil {
			if !errors.Is(err, sarama.ErrTopicAlreadyExists) {
				log.Error("Error creating topic: %s", err)
				panic("Error creating topic: " + err.Error())
			}
		}
	}
}
