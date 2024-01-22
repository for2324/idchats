package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"errors"
	"github.com/Pallinder/go-randomdata"
	"github.com/duke-git/lancet/v2/convertor"
	"strings"
	"time"

	_ "gorm.io/gorm"
)

func GetRegister(account, areaCode, userID string) (*db.Register, error) {
	var r db.Register
	return &r, db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("user_id = ? and user_id != ? or account = ? or account =? and area_code=?",
		userID, "", account, account, areaCode).Take(&r).Error
}

func GetRegisterWallet(account, userID string) (*db.Register, error) {
	var r db.Register
	return &r, db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("user_id = ? and user_id != ? or account = ? or account =?",
		userID, "", account, account).Take(&r).Error
}

func GetRegisterInfo(account string) (*db.Register, error) {
	var r db.Register
	return &r, db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("account = ?", account).Take(&r).Error
}
func SetPasswordWithInvitationCode(account, password, ex, userID, areaCode, ip, invitationCode string) error {
	r := db.Register{
		Account:        account,
		Password:       password,
		Ex:             ex,
		UserID:         userID,
		RegisterIP:     ip,
		AreaCode:       areaCode,
		InvitationCode: invitationCode,
	}
	return db.DB.MysqlDB.DefaultGormDB().Table("registers").Create(&r).Error
}
func SetPassword(account, password, ex, userID, areaCode, ip string) error {
	r := db.Register{
		Account:    account,
		Password:   password,
		Ex:         ex,
		UserID:     userID,
		RegisterIP: ip,
		AreaCode:   areaCode,
	}
	return db.DB.MysqlDB.DefaultGormDB().Table("registers").Create(&r).Error
}

func ResetPassword(account, password string) error {
	r := db.Register{
		Password: password,
	}
	return db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("account = ?", account).Updates(&r).Error
}

func GetRegisterAddFriendList(showNumber, pageNumber int32) ([]string, error) {
	var IDList []string
	var err error
	model := db.DB.MysqlDB.DefaultGormDB().Model(&db.RegisterAddFriend{})
	if showNumber == 0 {
		err = model.Pluck("user_id", &IDList).Error
	} else {
		err = model.Limit(int(showNumber)).Offset(int(showNumber*(pageNumber-1))).Pluck("user_id", &IDList).Error
	}
	return IDList, err
}

func AddUserRegisterAddFriendIDList(userIDList ...string) error {
	var list []db.RegisterAddFriend
	for _, v := range userIDList {
		list = append(list, db.RegisterAddFriend{UserID: v})
	}
	result := db.DB.MysqlDB.DefaultGormDB().Create(list)
	if int(result.RowsAffected) < len(userIDList) {
		return errors.New("some line insert failed")
	}
	err := result.Error
	return err
}

func ReduceUserRegisterAddFriendIDList(userIDList ...string) error {
	var list []db.RegisterAddFriend
	for _, v := range userIDList {
		list = append(list, db.RegisterAddFriend{UserID: v})
	}
	err := db.DB.MysqlDB.DefaultGormDB().Delete(list).Error
	return err
}

func DeleteAllRegisterAddFriendIDList() error {
	err := db.DB.MysqlDB.DefaultGormDB().Where("1 = 1").Delete(&db.RegisterAddFriend{}).Error
	return err
}

func GetUserIPLimit(userID string) (db.UserIpLimit, error) {
	var limit db.UserIpLimit
	limit.UserID = userID
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserIpLimit{}).Take(&limit).Error
	return limit, err
}
func IsUserHaveApiKey(userID string) (bool, error) {
	var count int64
	result := db.DB.MysqlDB.DefaultGormDB().Table("user_robot_api").Where("user_id=?", userID).Count(&count)
	if err := result.Error; err != nil {
		return true, err
	}
	return count > 0, nil
}
func InsertUserApiKey(userID string) error {
	keyRand := strings.ToLower(randomdata.RandStringRunes(32))
	db.DB.MysqlDB.DefaultGormDB().Table("user_robot_api").Create(&db.UserRobotAPI{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
		UserID:      userID,
		APIKey:      strings.ToLower(randomdata.RandStringRunes(32)),
		APISecret:   "",
		TradeVolume: "",
		TradeFee:    convertor.ToString(0.0018),
		SniperFee:   convertor.ToString(0.0090),
		Status:      0,
		APIName:     keyRand,
	})
	return nil

}
