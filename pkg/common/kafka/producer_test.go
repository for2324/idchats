package kafka

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	"testing"
)

func TestSendMessage(t *testing.T) {
	producer := NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, "test2")
	OperationID := "TestSendMessage"
	pushToUserID := "0xcbd033ea3c05dc9504610061c86c7ae191c5c913"
	log.NewInfo(OperationID, "TestSendMessage, start", config.Config.Kafka.Ws2mschat.Addr, "test1")
	mqPushMsg := pbMsg.PushMsgDataToMQ{
		OperationID:  OperationID,
		MsgData:      nil,
		PushToUserID: pushToUserID}
	_, _, err := producer.SendMessage(&mqPushMsg, pushToUserID, "")
	if err != nil {
		log.NewError(OperationID, "TestSendMessage, err = %s", err.Error())
		return
	}
	log.NewInfo(OperationID, "TestSendMessage, success")
}

func TestAutoCreateTopic(t *testing.T) {
	AutoCreateTopic()
}
