package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/http"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbOrder "Open_IM/pkg/proto/order"
	pbRelay "Open_IM/pkg/proto/relay"
	pbSdkWs "Open_IM/pkg/proto/sdk_ws"

	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

type OrderCallBackConsumerHandler struct {
	msgHandle                  map[string]fcb
	orderCallBackConsumerGroup *kfk.MConsumerGroup
}

func (pc *OrderCallBackConsumerHandler) Init() {
	pc.msgHandle = make(map[string]fcb)
	pc.msgHandle[config.Config.Kafka.MsgOrder.Topic] = pc.handleOrderCallBack
	pc.orderCallBackConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_8_1_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.MsgOrder.Topic},
		config.Config.Kafka.MsgOrder.Addr, config.Config.Kafka.ConsumerGroupID.MsgToOrderNotify)
	initRestartNotifyRetriedTask()
}

func (pc *OrderCallBackConsumerHandler) handleOrderCallBack(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession) {
	msg := cMsg.Value
	msgPushMQ := pbOrder.MsgDataToOrderByMQ{}
	err := proto.Unmarshal(msg, &msgPushMQ)
	if err != nil {
		log.NewError(msgPushMQ.OperationID, "msg_transfer Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	log.Debug(msgPushMQ.OperationID, "proto.Unmarshal MsgDataToMQ", msgPushMQ.String())

	// notify user
	NotifyUser(msgPushMQ.OperationID, &msgPushMQ)

	if msgPushMQ.NotifyUrl == "" {
		log.NewError(msgPushMQ.OperationID, "msg_transfer msgPushMQ.NotifyUrl is empty", "msg", string(msg))
		return
	}
	resp, err := HttpPostNotifyUrl(msgPushMQ.NotifyUrl, &msgPushMQ)
	if err != nil {
		log.NewError(msgPushMQ.OperationID, "msg_transfer HttpPostNotifyUrl err", "msg", string(msg), "err", err.Error())
		return
	}
	if resp.Code != "success" {
		// 加入重试列表
		record := &db.NotifyRetried{
			NotifyUrl:  msgPushMQ.NotifyUrl,
			NotifyBody: msg,
			Mark:       "order_notify",
		}
		err := imdb.CreateNotifyRetried(record)
		if err != nil {
			log.NewError(msgPushMQ.OperationID, "msg_transfer imdb.CreateNotifyRetried err", "msg", string(msg), "err", err.Error())
			return
		}
		go StartNotifyRetriedTask(record)
		log.NewError(msgPushMQ.OperationID, "msg_transfer HttpPostNotifyUrl err", "msg", string(msg), "resp", resp)
		return
	}
	// 修改 订单 状态 为已通知
	imdb.UpdatePayScanBlockTaskStatusNotified(msgPushMQ.ID)
}

func (*OrderCallBackConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*OrderCallBackConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (pc *OrderCallBackConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
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
		sess.MarkMessage(msg, "")
	}
	return nil
}

type NotifyContent struct {
	Mark    string `json:"mark"`
	Id      string `json:"id"`
	OrderId string `json:"orderId"`
}

func NotifyUser(OperationID string, msgPushMQ *pbOrder.MsgDataToOrderByMQ) {
	etcdConn := getcdv3.GetDefaultConn(
		config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImRelayName,
		OperationID,
	)
	if etcdConn == nil {
		errMsg := OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(OperationID, errMsg)
		return
	}
	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","), OperationID)

	//Online push message
	log.Debug(OperationID, "len  grpc", len(grpcCons), "data")
	notifyContent := &NotifyContent{
		Mark:    msgPushMQ.Mark,
		Id:      msgPushMQ.ID,
		OrderId: msgPushMQ.OrderId,
	}
	content, err := json.Marshal(notifyContent)
	if err != nil {
		log.NewError(OperationID, "json.Marshal failed ", err.Error())
		return
	}
	for _, v := range grpcCons {
		_, err := pbRelay.NewRelayClient(v).OnlinePushMsg(context.Background(), &pbRelay.OnlinePushMsgReq{
			OperationID: OperationID,
			MsgData: &pbSdkWs.MsgData{
				RecvID:      msgPushMQ.UserId,
				ContentType: constant.UserPayNotification,
				SessionType: constant.NotificationOnlinePushType,
				Content:     content,
			},
			PushToUserID: msgPushMQ.UserId,
		})
		if err != nil {
			log.NewError(OperationID, "OnlinePushMsg failed ", err.Error())
			return
		}
		log.NewInfo(OperationID, "OnlinePushMsg success", msgPushMQ.UserId)
	}

}

type HttpResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OrderSource struct {
	Mark       string `json:"mark"`
	CreateTime string `json:"createTime"`
	PayTime    string `json:"payTime"`
	TxnHash    string `json:"txnHash"`
	OrderId    string `json:"orderId"`
	TxnType    string `json:"txnType"`
	Attach     string `json:"attach"`
}

type PostData struct {
	Id          string `json:"id"`
	EventType   string `json:"eventType"`   // SUCCESS
	EncryptType string `json:"encryptType"` // AES RSA
	Source      string `json:"source"`
	Nonce       string `json:"nonce"`
}

// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_2.shtml
func HttpPostNotifyUrl(uri string, msg *pbOrder.MsgDataToOrderByMQ) (*HttpResponse, error) {
	log.NewDebug("", "HttpPostNotifyCallBackUrl", "uri", uri, "msg", msg.String())
	orderSource := OrderSource{
		Mark:       msg.Mark,
		PayTime:    msg.PayTime,
		TxnHash:    msg.TxnHash,
		OrderId:    msg.OrderId,
		TxnType:    msg.TxnType,
		Attach:     msg.Attach,
		CreateTime: msg.CreateTime,
	}
	orderSourceBytes, err := json.Marshal(orderSource)
	if err != nil {
		log.NewError("", "HttpPostNotifyCallBackUrl", "uri", uri, "msg", msg.String(), "err", err.Error())
		return nil, err
	}
	ciphertext, nonceStr, err := AESGCMEncrypt(orderSourceBytes, msg.NotifyEncryptKey)
	if err != nil {
		log.NewError("", "HttpPostNotifyCallBackUrl", "uri", uri, "msg", msg.String(), "err", err.Error())
		return nil, err
	}
	base64Source := base64.StdEncoding.EncodeToString(ciphertext)
	postData := PostData{
		Id:          msg.OperationID,
		EventType:   "SUCCESS",
		EncryptType: msg.NotifyEncryptType,
		Source:      base64Source,
		Nonce:       nonceStr,
	}
	httpResp, err := http.Post(uri, postData, 10)
	if err != nil {
		log.NewError("", "HttpPostNotifyCallBackUrl", "uri", uri, "msg", msg.String(), "err", err.Error())
		return nil, err
	}
	var httpHttpResponse *HttpResponse
	err = json.Unmarshal(httpResp, &httpHttpResponse)
	if err != nil {
		log.NewError("", "HttpPostNotifyCallBackUrl", "uri", uri, "msg", string(httpResp), "err", err.Error())
		return nil, err
	}
	log.NewInfo("", "HttpPostNotifyCallBackUrl", "uri", uri, "msg", msg.String(), "httpResp", *httpHttpResponse)
	return httpHttpResponse, nil
}

func StartNotifyRetriedTask(task *db.NotifyRetried) {
	retriedCount := task.RetriedCount + 1
	for retriedCount < 10 {
		<-time.After(time.Duration(retriedCount) * time.Second)

		msg := task.NotifyBody
		msgPushMQ := pbOrder.MsgDataToOrderByMQ{}
		err := proto.Unmarshal(msg, &msgPushMQ)
		if err != nil {
			log.NewError(msgPushMQ.OperationID, "msg_transfer Unmarshal msg err", "msg", string(msg), "err", err.Error())
			return
		}
		resp, err := HttpPostNotifyUrl(task.NotifyUrl, &msgPushMQ)
		if err == nil && resp.Code == "success" {
			imdb.UpdateNotifyRetriedStatusSuccess(task.ID)
		} else {
			imdb.IncreaseNotifyRetriedCount(task.ID)
		}
		retriedCount++
	}
}

func initRestartNotifyRetriedTask() {
	record, err := imdb.GetNotifyRetriedList()
	if err != nil {
		log.NewError("", "initRestartNotifyRetriedTask", "err", err.Error())
		return
	}
	for _, v := range record {
		go StartNotifyRetriedTask(&v)
	}
}

func AESGCMEncrypt(plaintext []byte, key string) (ciphertext []byte, nonceStr string, err error) {
	nonce := make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	nonceStr = base32.StdEncoding.EncodeToString(nonce)[:12]
	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		return
	}
	ciphertext = aesgcm.Seal(nil, []byte(nonceStr), plaintext, nil)
	return
}
func AESGCMDecrypt(ciphertext string, nonceStr string, key string) (plaintext string, err error) {
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		return
	}
	var openData []byte
	openData, err = aesgcm.Open(nil, []byte(nonceStr), []byte(ciphertext), nil)
	if err != nil {
		return
	}
	plaintext = string(openData)
	return
}
