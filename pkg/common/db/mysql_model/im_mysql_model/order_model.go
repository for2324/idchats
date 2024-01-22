package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

func CreatePayScanBlockTask(task *db.PayScanBlockTask) error {
	task.CreateTime = time.Now()
	return db.DB.MysqlDB.DefaultGormDB().Model(&db.PayScanBlockTask{}).Create(task).Error
}

type GroupTaskList struct {
	Type        string `gorm:"column:type"`
	Tag         string `gorm:"column:tag"`
	ChainId     int64  `gorm:"column:chain_id"`
	StartHeight uint64 `gorm:"column:start_block_number"`
}

func GetPayScanBlockGroupTaskList() ([]GroupTaskList, error) {
	var taskList []GroupTaskList
	err := db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Select("type,tag,chain_id,min(scan_block_number) as start_block_number").
		Where("status = ?", 0).
		Group("type,tag,chain_id").
		Find(&taskList).Error
	return taskList, err
}

func GetPayScanBlockTaskListByTaskTag(taskType string, taskTag string, start, end uint64) ([]db.PayScanBlockTask, error) {
	var taskList []db.PayScanBlockTask
	err := db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("status = ? and type = ? and tag = ? and scan_block_number >= ? and scan_block_number <= ?", 0, taskType, taskTag, start, end).
		Find(&taskList).Error
	return taskList, err
}

func MarkPayScanBlockTaskExpired(endBlockTime time.Time) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("status = ? and block_expire_time < ?", 0, endBlockTime).
		Updates(db.PayScanBlockTask{
			Status: constant.PayScanBlockTaskStatusExpired,
		}).Error
}

func UpdatePayScanBlockTaskStatusExpired(taskID uint64, txHash string) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("task_id = ? and status = ?", taskID, 0).
		Updates(db.PayScanBlockTask{
			Status:  constant.PayScanBlockTaskStatusExpired,
			TxnHash: txHash,
		}).Error
}

func UpdatePayScanBlockTaskStatusFinished(
	formAddress string, toAddress string, value string, ChainId int64, blockTime time.Time, txnHash string) (
	*db.PayScanBlockTask, error,
) {
	formAddress = strings.ToLower(formAddress)
	toAddress = strings.ToLower(toAddress)
	var orderInfo *db.PayScanBlockTask
	res := db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("status = ? and form_address = ? and to_address = ? and value = ? and chain_id = ? and block_start_time <= ? and block_expire_time >= ?", 0, formAddress, toAddress, value, ChainId, blockTime, blockTime).
		Updates(db.PayScanBlockTask{
			Status:       constant.PayScanBlockTaskStatusFinished,
			TxnHash:      txnHash,
			BlockPayTime: blockTime,
		})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("order not found")
	}
	err := db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("txn_hash = ?", txnHash).
		First(&orderInfo).Error
	return orderInfo, err
}
func ReplenishmentOrder(id uint64, txnHash string, remark string, BlockPayTime time.Time) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("id = ?", id).
		Updates(db.PayScanBlockTask{
			Status:       constant.PayScanBlockTaskStatusFinished,
			TxnHash:      txnHash,
			BlockPayTime: BlockPayTime,
			Ex:           remark,
		}).Error
}

func UpdatePayScanBlockTaskStatusComfirm(taskID uint64) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("id = ? and status = ?", taskID, 0).
		Updates(db.PayScanBlockTask{
			Status: constant.PayScanBlockTaskStatusConfirm,
		}).Error
}
func UpdatePayScanBlockTaskStatusNotified(taskID string) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("id = ? and status = ?", taskID, constant.PayScanBlockTaskStatusConfirm).
		Updates(db.PayScanBlockTask{
			Status: constant.PayScanBlockTaskStatusNotified,
		}).Error
}

func UpdatePayScanBlockGroupTaskProgress(progress uint64, taskType string, taskTag string, start, end uint64) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("status = ? and type = ? and tag = ? and scan_block_number >= ? and scan_block_number <= ?", 0, taskType, taskTag, start, end).
		Updates(db.PayScanBlockTask{
			ScanBlockNumber: progress,
		}).Error
}

func GetPayScanBlockTaskByOrderId(orderId string) (*db.PayScanBlockTask, error) {
	var task db.PayScanBlockTask
	err := db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("order_id = ?", orderId).
		First(&task).Error
	return &task, err
}

func GetPayScanBlockTaskById(id uint64) (*db.PayScanBlockTask, error) {
	var task db.PayScanBlockTask
	err := db.DB.MysqlDB.DefaultGormDB().
		Model(&db.PayScanBlockTask{}).
		Where("id = ?", id).
		First(&task).Error
	return &task, err
}

func CreateOrderPaidRecord(record *db.OrderPaidRecord) error {
	record.CreateTime = time.Now()
	record.FormAddress = strings.ToLower(record.FormAddress)
	record.ToAddress = strings.ToLower(record.ToAddress)
	return db.DB.MysqlDB.DefaultGormDB().Model(&db.OrderPaidRecord{}).Create(record).Error
}

func CreateNotifyRetried(record *db.NotifyRetried) error {
	record.CreateAt = time.Now()
	return db.DB.MysqlDB.DefaultGormDB().Model(&db.NotifyRetried{}).Create(record).Error
}

func GetNotifyRetriedList() ([]db.NotifyRetried, error) {
	var list []db.NotifyRetried
	err := db.DB.MysqlDB.DefaultGormDB().
		Model(&db.NotifyRetried{}).
		Where("retried_count < ? and status = 0", 10).
		Find(&list).Error
	return list, err
}

func UpdateNotifyRetriedStatusSuccess(id int) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.NotifyRetried{}).
		Where("id = ?", id).
		Updates(db.NotifyRetried{
			Status: 1,
		}).Error
}

func IncreaseNotifyRetriedCount(id int) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.NotifyRetried{}).
		Where("id = ?", id).
		Update("retried_count", gorm.Expr("retried_count + ?", 1)).Error
}
