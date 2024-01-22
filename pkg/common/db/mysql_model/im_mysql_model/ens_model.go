package im_mysql_model

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"time"

	"gorm.io/gorm"
)

func Appointment(userId string, ensName string) error {
	return db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		err := tx.
			Where("ens_name = ? AND created_at < ?", ensName, time.Now().Add(-10*time.Minute)).
			Delete(&db.EnsAppointment{}).Error
		if err != nil {
			return err
		}
		err = tx.Model(&db.EnsAppointment{}).Create(&db.EnsAppointment{
			UserID:    userId,
			EnsName:   ensName,
			CreatedAt: time.Now(),
		}).Error
		if err != nil {
			return constant.ErrHaveBeenAppointment
		}
		return nil
	})
}

func CancelAppointment(userId string) error {
	// delete from ens_appointment where user_id = ?
	return db.DB.MysqlDB.DefaultGormDB().Where("user_id = ?", userId).Delete(&db.EnsAppointment{}).Error
}
func CancelAppointmentEnsName(userId string, ensName string) error {
	// delete from ens_appointment where user_id = ?
	return db.DB.MysqlDB.DefaultGormDB().Where("user_id = ? and ens_name = ?", userId, ensName).Delete(&db.EnsAppointment{}).Error
}
func AppointmentList(pageIndex int, pageSize int) (ensList []api.AppointmentUserInfo, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().
		Model(&db.EnsAppointment{}).
		Select("ens_appointment.*, u.face_url").
		Joins("left join users u on ens_appointment.user_id = u.user_id").
		Where("ens_appointment.created_at > ?", time.Now().Add(-10*time.Minute)).
		Find(&ensList).
		Offset(pageIndex * pageSize).Limit(pageSize).Error
	return
}

func MyAppointmentList(userID string, pageIndex int, pageSize int) (ensList []api.AppointmentUserInfo, err error) {
	db.DB.MysqlDB.DefaultGormDB().
		Where("user_id = ? AND created_at < ?", userID, time.Now().Add(-10*time.Minute)).
		Delete(&db.EnsAppointment{})
	err = db.DB.MysqlDB.DefaultGormDB().
		Model(&db.EnsAppointment{}).
		Select("ens_appointment.*, u.face_url").
		Where("ens_appointment.user_id = ?", userID).
		Joins("left join users u on ens_appointment.user_id = u.user_id").
		Find(&ensList).Offset(pageIndex * pageSize).Limit(pageSize).Error
	return
}

func HasAppointment(ens string) (bool, error) {
	db.DB.MysqlDB.DefaultGormDB().
		Where("ens_name = ? AND created_at < ?", ens, time.Now().Add(-10*time.Minute)).
		Delete(&db.EnsAppointment{})
	var count int64
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.EnsAppointment{}).Where("ens_name = ?", ens).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateEnsRegisterOrder(orderInfo *db.EnsRegisterOrder) error {
	orderInfo.CreateTime = time.Now()
	orderInfo.ExpireTime = time.Now().Add(time.Minute * time.Duration(config.Config.Pay.OrderExpireTime))
	orderInfo.Status = constant.RegisterEnsOrderStatusUnpaid
	return db.DB.MysqlDB.DefaultGormDB().Model(&db.EnsRegisterOrder{}).Create(orderInfo).Error
}

func GetEnsRegisterOrderInfo(orderID uint64) (*db.EnsRegisterOrder, error) {
	var orderInfo db.EnsRegisterOrder
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.EnsRegisterOrder{}).Where("order_id = ?", orderID).First(&orderInfo).Error
	return &orderInfo, err
}

func GetEnsOrderRegisterNotMined() ([]db.EnsRegisterOrder, error) {
	var orderList []db.EnsRegisterOrder
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.EnsRegisterOrder{}).Where("status = ? and register_txn_hash <> ''", constant.RegisterEnsOrderStatusPaid).Find(&orderList).Error
	return orderList, err
}

func GetOrderInfoByRegisterTxnHash(RegisterTxnHash string) (*db.EnsRegisterOrder, error) {
	var orderInfo db.EnsRegisterOrder
	err := db.DB.MysqlDB.DefaultGormDB().
		Model(&db.EnsRegisterOrder{}).
		Where("register_txn_hash = ?", RegisterTxnHash).
		First(&orderInfo).Error
	return &orderInfo, err
}

func UpdateEnsRegisterOrderPaid(orderID uint64, TxnHash string) *gorm.DB {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.EnsRegisterOrder{}).
		Where("order_id = ? AND status = ?", orderID, constant.RegisterEnsOrderStatusUnpaid).
		Updates(db.EnsRegisterOrder{
			Status:  constant.RegisterEnsOrderStatusPaid,
			TxnHash: TxnHash,
			PayTime: time.Now(),
		})
}

func UpdateEnsOrderRegisterFailed(orderID uint64, errMsg string) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.EnsRegisterOrder{}).
		Where("order_id = ? AND status = ?", orderID, constant.RegisterEnsOrderStatusPaid).
		Updates(db.EnsRegisterOrder{
			Status: constant.RegisterEnsOrderStatusRegisterFailed,
			Ex:     errMsg,
		}).Error
}

func UpdateEnsOrderRegisterConfirmed(orderID uint64, txnHash string) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.EnsRegisterOrder{}).
		Where("order_id = ?", orderID).
		Updates(db.EnsRegisterOrder{
			RegisterTxnHash: txnHash,
		}).Error
}

func UpdateEnsOrderRegisterSuccess(orderID uint64) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.EnsRegisterOrder{}).
		Where("order_id = ? and status = ?", orderID, constant.RegisterEnsOrderStatusPaid).
		Updates(db.EnsRegisterOrder{
			Status: constant.RegisterEnsOrderStatusConfirmed,
		}).Error
}

func UpdateEnsOrderRegisterRefund(orderID uint64, errMsg string) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.EnsRegisterOrder{}).
		Where("order_id = ?", orderID).
		Updates(db.EnsRegisterOrder{
			Status: constant.RegisterEnsOrderStatusRefund,
			Ex:     errMsg,
		}).Error
}
