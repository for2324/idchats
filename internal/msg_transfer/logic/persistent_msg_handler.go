/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/11 15:37).
 */
package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_msg_model"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	pbMsg "Open_IM/pkg/proto/msg"
	pbTask "Open_IM/pkg/proto/task"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type PersistentConsumerHandler struct {
	msgHandle               map[string]fcb
	persistentConsumerGroup *kfk.MConsumerGroup
	rpcTaskClient           pbTask.TaskServiceClient
}

func (pc *PersistentConsumerHandler) Init() {

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, "PersistentConsumerHandler init etcd")
	if etcdConn == nil {
		panic("PersistentConsumerHandler init etcd failed")
	}
	// rpc
	pc.rpcTaskClient = pbTask.NewTaskServiceClient(etcdConn)

	pc.msgHandle = make(map[string]fcb)
	pc.msgHandle[config.Config.Kafka.Ws2mschat.Topic] = pc.handleChatWs2Mysql
	pc.persistentConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_8_1_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMySql)

}

func (pc *PersistentConsumerHandler) handleChatWs2Mysql(cMsg *sarama.ConsumerMessage, msgKey string, _ sarama.ConsumerGroupSession) {
	msg := cMsg.Value
	log.NewInfo("msg come here mysql!!!", "", "msg", string(msg), msgKey)
	var tag bool
	msgFromMQ := pbMsg.MsgDataToMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.NewError(msgFromMQ.OperationID, "msg_transfer Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	log.Debug(msgFromMQ.OperationID, "proto.Unmarshal MsgDataToMQ", msgFromMQ.String())
	//Control whether to store history messages (mysql)
	isPersist := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, constant.IsPersistent)
	//Only process receiver data
	if isPersist {
		switch msgFromMQ.MsgData.SessionType {
		case constant.SingleChatType, constant.NotificationChatType:
			if msgKey == msgFromMQ.MsgData.RecvID {
				tag = true
			}
		case constant.GroupChatType:
			if msgKey == msgFromMQ.MsgData.SendID {
				tag = true
			}
		case constant.SuperGroupChatType:
			tag = true
		}
		if tag {
			log.NewInfo(msgFromMQ.OperationID, "msg_transfer msg persisting", string(msg))
			if msgFromMQ.MsgData.SendID == "announce" {
				//通知消息不入库
				return
			}
			if err = im_mysql_msg_model.InsertMessageToChatLog(msgFromMQ); err != nil {
				log.NewError(msgFromMQ.OperationID, "Message insert failed", "err", err.Error(), "msg", msgFromMQ.String())
				return
			} else {
				log.NewInfo(msgFromMQ.OperationID, "Message insert success", "SessionType", msgFromMQ.MsgData.SessionType)
				if msgFromMQ.MsgData.SessionType == constant.SingleChatType && msgFromMQ.MsgData.SendID != msgFromMQ.MsgData.RecvID {
					//如果是聊天的信息
					if msgFromMQ.MsgData.ContentType >= constant.Text && msgFromMQ.MsgData.ContentType <= constant.Custom {
						log.NewInfo(msgFromMQ.OperationID, "CheckSomeBodyChatIsSendToMe", "SessionType", msgFromMQ.MsgData.SendID, "=>", msgFromMQ.MsgData.RecvID)
						// 如果是好友的话（不是陌生地址），就不需要去完成任务了
						if isFri, err := isFriend(msgFromMQ.OperationID, msgFromMQ.MsgData.SendID, msgFromMQ.MsgData.RecvID); err != nil || isFri {
							return
						}
						// 今天是否互相发送消息
						if im_mysql_msg_model.CheckSomeBodyChatIsSendToMe(&msgFromMQ) {
							log.NewInfo(msgFromMQ.OperationID, "todat is send to each", msgFromMQ.MsgData.SendID, "=>", msgFromMQ.MsgData.RecvID)
							// 完成 RecvID 的聊天任务
							{
								// 是否完成携带NFT头像的任务
								resp, err := pc.rpcTaskClient.IsFinishUploadNftHeadTask(context.Background(), &pbTask.IsFinishUploadNftHeadTaskReq{
									OperationID: msgFromMQ.OperationID,
									UserId:      msgFromMQ.MsgData.RecvID,
								})
								if err == nil {
									// 完成携带头像的任务才可以去完成互动任务
									if resp.IsFinish {
										resp, err := pc.rpcTaskClient.FinishDailyChatNFTHeadWithNewUserTask(context.Background(), &pbTask.FinishDailyChatNFTHeadWithNewUserTaskReq{
											OperationID: msgFromMQ.OperationID,
											UserId:      msgFromMQ.MsgData.RecvID,
											ChatUser:    msgFromMQ.MsgData.SendID,
										})
										if err != nil {
											log.NewError(msgFromMQ.OperationID, "FinishDailyChatNFTHeadWithNewUserTask", err.Error())
										} else {
											log.NewInfo(msgFromMQ.OperationID, "FinishDailyChatNFTHeadWithNewUserTask success", resp.CommonResp.ErrMsg)
										}
									}
								}
							}
							{
								// 是否完成携带官方NFT头像的任务
								resp, err := pc.rpcTaskClient.IsFinishOfficialNFTHeadTask(context.Background(), &pbTask.IsFinishOfficialNFTHeadTaskReq{
									OperationID: msgFromMQ.OperationID,
									UserId:      msgFromMQ.MsgData.RecvID,
								})
								if err == nil {
									// 完成携带头像的任务才可以去完成互动任务
									if resp.IsFinish {
										resp, err := pc.rpcTaskClient.FinishOfficialNFTHeadDailyChatWithNewUserTask(context.Background(), &pbTask.FinishOfficialNFTHeadDailyChatWithNewUserTaskReq{
											OperationID: msgFromMQ.OperationID,
											UserId:      msgFromMQ.MsgData.RecvID,
											ChatUser:    msgFromMQ.MsgData.SendID,
										})
										if err != nil {
											log.NewError(msgFromMQ.OperationID, "FinishOfficialNFTHeadDailyChatWithNewUserTask", err.Error())
										} else {
											log.NewInfo(msgFromMQ.OperationID, "FinishOfficialNFTHeadDailyChatWithNewUserTask success", resp.CommonResp.ErrMsg)
										}
									}
								}
							}
							// 完成 SendID 的聊天任务
							{
								// 是否完成携带NFT头像任务
								resp, err := pc.rpcTaskClient.IsFinishUploadNftHeadTask(context.Background(), &pbTask.IsFinishUploadNftHeadTaskReq{
									OperationID: msgFromMQ.OperationID,
									UserId:      msgFromMQ.MsgData.SendID,
								})
								if err == nil {
									// 完成携带头像的任务才可以去完成互动任务
									if resp.IsFinish {
										resp, err := pc.rpcTaskClient.FinishDailyChatNFTHeadWithNewUserTask(context.Background(), &pbTask.FinishDailyChatNFTHeadWithNewUserTaskReq{
											OperationID: msgFromMQ.OperationID,
											UserId:      msgFromMQ.MsgData.SendID,
											ChatUser:    msgFromMQ.MsgData.RecvID,
										})
										if err != nil {
											log.NewError(msgFromMQ.OperationID, "FinishDailyChatNFTHeadWithNewUserTask", err.Error())
										} else {
											log.NewInfo(msgFromMQ.OperationID, "FinishDailyChatNFTHeadWithNewUserTask success", resp.CommonResp.ErrMsg)
										}
									}
								}
							}
							{
								// 是否完成携带官方NFT头像任务
								resp, err := pc.rpcTaskClient.IsFinishOfficialNFTHeadTask(context.Background(), &pbTask.IsFinishOfficialNFTHeadTaskReq{
									OperationID: msgFromMQ.OperationID,
									UserId:      msgFromMQ.MsgData.SendID,
								})
								if err == nil {
									// 完成携带头像的任务才可以去完成互动任务
									if resp.IsFinish {
										resp, err := pc.rpcTaskClient.FinishOfficialNFTHeadDailyChatWithNewUserTask(context.Background(), &pbTask.FinishOfficialNFTHeadDailyChatWithNewUserTaskReq{
											OperationID: msgFromMQ.OperationID,
											UserId:      msgFromMQ.MsgData.SendID,
											ChatUser:    msgFromMQ.MsgData.RecvID,
										})
										if err != nil {
											log.NewError(msgFromMQ.OperationID, "FinishOfficialNFTHeadDailyChatWithNewUserTask", err.Error())
										} else {
											log.NewInfo(msgFromMQ.OperationID, "FinishOfficialNFTHeadDailyChatWithNewUserTask success", resp.CommonResp.ErrMsg)
										}
									}
								}
							}
						}
					}
				}
			}
		}

	}
}
func (*PersistentConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*PersistentConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (pc *PersistentConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
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

type EnsAddressReq struct {
	Address string
}
type EnsAddressRsp struct {
	ErrCode int32 `json:"errCode"`
	Data    int32 `json:"data"`
}

func isFriend(OperationID, fromUserID, toUserID string) (bool, error) {
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, OperationID)
	if etcdConn == nil {
		return false, errors.New("isFriend init etcd failed")
	}
	client := pbFriend.NewFriendClient(etcdConn)
	req := &pbFriend.IsFriendReq{CommID: &pbFriend.CommID{
		OpUserID:    fromUserID,
		OperationID: OperationID,
		ToUserID:    toUserID,
		FromUserID:  fromUserID,
	}}
	RpcResp, err := client.IsFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "IsFriend failed ", err.Error(), req.String())
		return false, err
	}
	return RpcResp.Response, nil
}
