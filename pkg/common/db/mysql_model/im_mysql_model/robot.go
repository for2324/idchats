package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"encoding/base64"
	"time"
)

func CreateRobot(robot *db.Robot) (err error) {
	robot.CreateAt = time.Now()
	return db.DB.MysqlDB.DefaultGormDB().Create(robot).Error
}

func GetUserRobot(userId string) (robot *db.Robot, err error) {
	robot = &db.Robot{}
	err = db.DB.MysqlDB.DefaultGormDB().Where("user_id = ?", userId).First(robot).Error
	return
}
func GetRobotAddressIsUserRobot(userId string, robotAddress string) (robot *db.Robot, err error) {
	robot = &db.Robot{}
	err = db.DB.MysqlDB.DefaultGormDB().Where("user_id = ? and eth_address=? ", userId, robotAddress).First(robot).Error
	return
}
func GetUserRobotInfo(userID string) (robot *db.Robot, err error) {
	robot = &db.Robot{}
	err = db.DB.MysqlDB.DefaultGormDB().Where("user_id = ?", userID).First(robot).Error
	if err == nil {
		privateKeyBytes, _ := base64.StdEncoding.DecodeString(robot.EthPrivateKey)
		privateKeyBytes, _ = utils.AesDecrypt(privateKeyBytes, []byte("U2FsdGVkX1+4xoFd+2jiqf+m16e3EdEQ"))
		robot.EthPrivateKey = utils.Bytes2string(privateKeyBytes)
		privateKeyBytes, _ = base64.StdEncoding.DecodeString(robot.BtcPrivateKey)
		privateKeyBytes, _ = utils.AesDecrypt(privateKeyBytes, []byte("U2FsdGVkX1+4xoFd+2jiqf+m16e3EdEQ"))
		robot.BtcPrivateKey = utils.Bytes2string(privateKeyBytes)
		return
	}
	return nil, err
}
func GetOrderInfo(userID string, OrderID string) (OrdInfo *db.RoBotTask, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("swap_robot_task").
		Where("user_id = ? and ord_id =?", userID, OrderID).First(OrdInfo).Error
	if err == nil {
		return
	}
	return nil, err
}
func GetUpdateOrderInfo(userID string, OrderID string) (OrdInfo *db.RoBotTask, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("swap_robot_task").
		Where("user_id = ? and ord_id =?", userID, OrderID).Updates(map[string]interface{}{"order_status": "cancel"}).First(OrdInfo).Error
	if err == nil {
		return
	}
	return nil, err
}
func ChangeOrderIDStatus(OrderID string, Status string, txHash string) (err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("swap_robot_task").
		Where("ord_id =?", OrderID).Updates(map[string]interface{}{
		"order_status": Status,
		"tx_hash":      txHash,
	}).Error
	return
}
func GetApiKeyUserID(apiKey string) (parentUserId string, err error) {
	if apiKey == "" {
		return "", nil
	}
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_robot_api").Where("api_key=?", apiKey).Pluck("user_id", &parentUserId).Error
	return
}
func GetApiKeyUserRobotApi(apiKey string) (api *db.UserRobotAPI, err error) {
	if apiKey == "" {
		return nil, nil
	}
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_robot_api").Where("api_key=?", apiKey).
		First(api).Error

	return
}
func GetParentUserID(userID string) (parentUserId string, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("user_id=?", userID).Pluck("invitation_code", &parentUserId).Error
	return
}
