/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/3/4 11:18).
 */
package im_mysql_msg_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"errors"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func InsertMessageToChatLog(msg pbMsg.MsgDataToMQ) error {
	chatLog := new(db.ChatLog)
	copier.Copy(chatLog, msg.MsgData)
	switch msg.MsgData.SessionType {
	case constant.GroupChatType, constant.SuperGroupChatType:
		chatLog.RecvID = msg.MsgData.GroupID
	case constant.SingleChatType:
		chatLog.RecvID = msg.MsgData.RecvID
	}
	if msg.MsgData.ContentType >= constant.NotificationBegin && msg.MsgData.ContentType <= constant.NotificationEnd {
		var tips server_api_params.TipsComm
		_ = proto.Unmarshal(msg.MsgData.Content, &tips)
		marshaler := jsonpb.Marshaler{
			OrigName:     true,
			EnumsAsInts:  false,
			EmitDefaults: false,
		}
		chatLog.Content, _ = marshaler.MarshalToString(&tips)

	} else {
		chatLog.Content = string(msg.MsgData.Content)
	}
	chatLog.CreateTime = utils.UnixMillSecondToTime(msg.MsgData.CreateTime)
	chatLog.SendTime = utils.UnixMillSecondToTime(msg.MsgData.SendTime)
	log.NewDebug("test", "this is ", *chatLog)
	return db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Create(chatLog).Error
}

// 返回true 代表我发送消息的那个人 今天跟我聊天过
func CheckSomeBodyChatIsSendToMe(msg *pbMsg.MsgDataToMQ) bool {
	var chats db.ChatLog
	err := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Where("send_id=? and recv_id=? and session_type=1 and "+
		" send_time>=DATE_FORMAT(CURDATE(),'%Y-%m-%d %H:%i:%s') and content_type between 101 and 110 ", msg.MsgData.RecvID, msg.MsgData.SendID).
		Select("chat_logs.*").First(&chats).Error
	if err == nil {
		return true
	}
	return false
}
func CheckIsBindTwitterOrPhone(msg *pbMsg.MsgDataToMQ) bool {
	var chats db.EventLogs
	err := db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Where("user_id=? and event_typename='phone' ", msg.MsgData.SendID).First(&chats).Error
	if err == nil {
		var chats2 db.EventLogs
		err = db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Where("user_id=? and event_typename='phone' ", msg.MsgData.RecvID).First(&chats2).Error
		if err == nil {
			return true
		}
	}
	return false
}
func GetMessageFromMysqlDb(msg *db.ChatLog) (result *db.ChatLog, err error) {
	sqlDb := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs")
	if msg.ClientMsgID != "" {
		sqlDb = sqlDb.Where("client_msg_id =?", msg.ClientMsgID)
	}
	err = sqlDb.Take(result).Error
	return
}
func GetMessageFromMysqlDbByClientID(msgClientID string) (result *db.ChatLog, err error) {
	if msgClientID == "" {
		return nil, errors.New("dot select empty msgClient")
	}
	result = new(db.ChatLog)
	err = db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Where("client_msg_id=?", msgClientID).First(&result).Error
	return
}
func GetGroupMessageListNoTimeControllerFromMysqlDb(sourceID string, sessionType, count int, isReverse bool, channelId string) (resultmessage []*db.ChatLog, err error) {
	var condition, timeOrder string
	if isReverse {
		timeOrder = "send_time ASC"
	} else {
		timeOrder = "send_time DESC"
	}
	condition = "recv_id = ? AND status <=2 And session_type = ? "
	if channelId == "" {
		channelId = "1"
	}
	condition += " and channel_id=?"
	err = db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Where(condition,
		sourceID, sessionType, channelId).Order(timeOrder).Offset(0).Limit(count).Find(&resultmessage).Error

	return
}
func GetGroupMessageListControllerFromMysqlDb(sourceID string, sessionType, count int,
	startTime time.Time,
	isReverse bool, channelId string) (resultmessage []*db.ChatLog, err error) {

	var condition, timeOrder, timeSymbol string

	if isReverse {
		timeOrder = "send_time ASC"
		timeSymbol = ">"
	} else {
		timeOrder = "send_time DESC"
		timeSymbol = "<"
	}
	condition = "recv_id = ? AND status <=2 And session_type = ? And send_time " + timeSymbol + " ?"
	if channelId == "" {
		channelId = "1"
	}
	condition += " and channel_id=?"
	err = db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Where(condition,
		sourceID, sessionType, startTime, channelId).Order(timeOrder).Offset(0).Limit(count).Find(&resultmessage).Error

	return
}
func GetSingleChatMessageListNoTimeControllerFromMysqlDb(
	sendId, recvId string, count int, isReverse bool, channelId string,
) (resultmessage []*db.ChatLog, err error) {
	var condition, timeOrder string
	if isReverse {
		timeOrder = "send_time ASC"
	} else {
		timeOrder = "send_time DESC"
	}
	condition = "send_id = ? AND recv_id = ? AND status <=2 And session_type = ? "
	if channelId != "" {
		condition += gorm.Expr(" AND channel_id=?", channelId).SQL
	}
	condition += " and channel_id=?"
	err = db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Where(condition,
		sendId, recvId, constant.SingleChatType, channelId).Order(timeOrder).Offset(0).Limit(count).Find(&resultmessage).Error
	return
}
func GetSingleChatMessageListControllerFromMysqlDb(sendId, recvId string, count int,
	startTime time.Time,
	isReverse bool, channelId string,
) (resultmessage []*db.ChatLog, err error) {
	var condition, timeOrder, timeSymbol string
	if isReverse {
		timeOrder = "send_time ASC"
		timeSymbol = ">"
	} else {
		timeOrder = "send_time DESC"
		timeSymbol = "<"
	}
	condition = "send_id = ? AND recv_id = ? AND status <=2 And session_type = ? And send_time " + timeSymbol + " ?"
	if channelId != "" {
		condition += gorm.Expr(" AND channel_id=?", channelId).SQL
	}
	// condition += gorm.Expr(" AND content_type=?", constant.Text).SQL
	err = db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Where(condition,
		sendId, recvId, constant.SingleChatType, startTime).Order(timeOrder).Offset(0).Limit(count).Find(&resultmessage).Error

	return
}
