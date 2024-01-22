package notifyservice

import (
	kafkaMessage "Open_IM/pkg/proto/kafkamessage"
	"Open_IM/pkg/xkafka"
	"Open_IM/pkg/xlog"

	"google.golang.org/protobuf/proto"
)

type Sender interface {
	Send(receivers []string, subject, content string) error
}

type ListenerSendServer struct {
	SenderServiceEmail Sender
	SenderServiceSms   Sender
}

func (lt *ListenerSendServer) Listen(message xkafka.ConsumerMessage, ackNoewLedgment *xkafka.Acknowledgment) {
	if len(message.Value) != 0 {
		xlog.CInfo("current offset is ", message.Offset)
		var kafkaValue kafkaMessage.KafkaMsg
		err := proto.Unmarshal(message.Value, &kafkaValue)
		if err == nil {
			switch kafkaValue.MessageType {
			case 1:
				err = lt.SenderServiceEmail.Send([]string{kafkaValue.EmailMsg.ToAddress}, kafkaValue.EmailMsg.Subject, kafkaValue.EmailMsg.Body)
				if err != nil {
					xlog.CError(err.Error())
				}
			case 2:
				lt.SenderServiceSms.Send([]string{kafkaValue.EmailMsg.ToAddress}, kafkaValue.EmailMsg.Subject, kafkaValue.EmailMsg.Body)
			}
			ackNoewLedgment.Acknowledge()
		}
	}
}
