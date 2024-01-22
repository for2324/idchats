package logic

import (
	pbOrder "Open_IM/pkg/proto/order"
	"testing"
)

func TestNotifyUrl(t *testing.T) {
	notifyUrl := "http://127.0.0.1:10002/order/test_notify"
	msg := &pbOrder.MsgDataToOrderByMQ{
		OperationID:       "1234567890",
		NotifyUrl:         notifyUrl,
		USD:               1000000,
		OrderId:           "1234567890",
		UserId:            "1234567890",
		TxnType:           "MATIC",
		Attach:            "Attach",
		CreateTime:        "2021-09-01 00:00:00",
		NotifyEncryptType: "AES-GCM",
		NotifyEncryptKey:  "qqqqwwwweeeerrrr",
		ID:                "111",
	}
	HttpPostNotifyUrl(notifyUrl, msg)
}

func TestQueryOrder(t *testing.T) {

}
