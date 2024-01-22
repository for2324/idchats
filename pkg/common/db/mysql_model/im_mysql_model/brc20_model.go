package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"errors"
	"gorm.io/gorm"
)

func GetReadToPledge(btcSenderAddress string) (rebackPreObbt []*db.ObbtPrePledge, err error) {
	// delete from ens_appointment where user_id = ?
	err = db.DB.MysqlDB.DefaultGormDB().Table("obbt_pre_pledge").
		Where("sender_btc_address = ?", btcSenderAddress).Find(&rebackPreObbt).Error
	return
}
func DelPrePledge(btcSenderAddress string) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("obbt_pre_pledge").Delete("sender_btc_address=?", btcSenderAddress).Error
	return err
}
func InsertIntoPrePledger(insertData []*db.ObbtPrePledge) error {
	err := db.DB.MysqlDB.DefaultGormDB().
		Table("obbt_pre_pledge").
		Create(&insertData).Error
	return err
}
func DeletePrePledgeByInscriptionId(btcAddress string, sameInsciptionList []string) (err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Exec(`delete FROM obbt_pre_pledge where sender_btc_address =? and staking_Period in ( select b.staking_period from
		(select staking_period from obbt_pre_pledge where sender_btc_address = ? and inscription_id in (?)) as b )`,
		btcAddress, btcAddress, sameInsciptionList).Error
	return
}

func GetAllPrePledge(btcAddress string) (result []*db.ObbtPrePledge, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("obbt_pre_pledge").
		Where("sender_btc_address = ?", btcAddress).
		Find(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}
func GetAllPrePledgeOnlyInscription(btcAddress string) (result []string, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("obbt_pre_pledge").
		Where("sender_btc_address = ?", btcAddress).Pluck("inscription_id", &result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}
