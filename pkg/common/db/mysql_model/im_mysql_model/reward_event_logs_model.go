package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"github.com/panjf2000/ants/v2"
	"gorm.io/gorm"
	"sync"
	"time"
)

func GetUnSyncUserTaskRewards(limit int) (resultdata []*db.RewardEventLogs, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("reward_event_logs").
		Where("is_sync=0").
		Limit(limit).
		Find(&resultdata).Error
	return
}

// 更新 reward_event_logs 表中 is_sync 为 1
func UpdateRewardEventLogSync(eventLogId string) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("reward_event_logs").Where("id=?", eventLogId).
		Update("is_sync", 1).Error
}
func EveryRechargeGameReword() {
	var gameConfig []*db.GameConfig
	err := db.DB.MysqlDB.DefaultGormDB().Table("game_config").Find(&gameConfig).Error
	operation := utils.OperationIDGenerator()
	if err != nil {
		log.NewInfo(operation, utils.GetSelfFuncName(), err.Error())
	} else {
		var wg sync.WaitGroup
		p, _ := ants.NewPoolWithFunc(5, func(msg interface{}) {
			msgNewMsgPoolDataData := (msg).(*db.GameConfig)
			if msgNewMsgPoolDataData.GameCurrentPrizePool < 5000000 {
				err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
					err := tx.Table("game_config").Where("game_id=?", msgNewMsgPoolDataData.GameId).Updates(map[string]interface{}{"game_current_prize_pool": 5000000}).Error
					if err != nil {
						return err
					}
					err = tx.Table("game_score_log").Create(&db.GameScoreLog{
						GameID:               int(msgNewMsgPoolDataData.GameId),
						CreatedAt:            time.Now(),
						UpdatedAt:            time.Now(),
						GameCurrentPrizePool: msgNewMsgPoolDataData.GameCurrentPrizePool,
						GameAddPrizePool:     5000000,
						JoinGameNumber:       0,
						OperatorIp:           "每日充值积分",
					}).Error
					return err
				})
			}
			wg.Done()
		})
		defer p.Release()
		for i := 0; i < len(gameConfig); i++ {
			wg.Add(1)
			msgData := gameConfig[i]
			_ = p.Invoke(msgData)
		}
		wg.Wait()
		log.NewInfo(operation, utils.GetSelfFuncName(), "每日充值积分")
	}
}
