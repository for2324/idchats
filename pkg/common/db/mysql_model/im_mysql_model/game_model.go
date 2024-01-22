package im_mysql_model

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

func GetLastStartTime(UserID string) (userDataGameStatus *db.UserGameScore, err error) {
	userDataGameStatus = new(db.UserGameScore)
	err = db.DB.MysqlDB.DefaultGormDB().Where("user_id=?", UserID).Order("created_at desc").Limit(1).Find(&userDataGameStatus).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("没有开始的记录")
	}
	return userDataGameStatus, err
}

type UserGameRank struct {
	UserID             string `gorm:"column:user_id"`
	Nickname           string `gorm:"column:name"`
	FaceURL            string `gorm:"column:face_url"`
	Score              int64  `gorm:"column:score"`
	Index              int32  `gorm:"column:index"`
	Reward             int64  `gorm:"column:reward"`
	TokenContractChain string `gorm:"column:token_contract_chain"`
}

func GetRankLink(gameID int32, selectUserID string) (result []*UserGameRank, NowIndex *UserGameRank, err error) {
	today := time.Now()
	nowMinute := today.Minute()
	if nowMinute > 30 {
		nowMinute = 30
	} else {
		nowMinute = 0
	}
	nowTody30MinUtc := time.Date(today.Year(), today.Month(), today.Day(), today.Hour(), nowMinute, 0, 0, today.Location())
	err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		err = tx.Table("user_game_score").
			Joins("left join users on user_game_score.user_id = users.user_id ").
			Where("game_id= ? and created_at=? and user_game_score.end_time<>'' ", gameID, nowTody30MinUtc).
			Select("user_game_score.user_id,users.name,users.face_url,users.token_contract_chain,user_game_score.score").
			Order("score desc").Limit(100).Find(&result).Error
		if err != nil {
			return err
		} else {
			isExist := false
			rankIndex := 0
			for key, value := range result {
				value.Index = int32(key) + 1
				value.Reward = config.GetRewordFromIndex(value.Index)
				if strings.EqualFold(value.UserID, selectUserID) {
					isExist = true
					rankIndex = key
				}
			}
			if len(result) > 0 {
				if !isExist {
					err = tx.Raw(`SELECT COUNT(score)+1 AS rank FROM user_game_score WHERE game_id= ?  and created_at =? and  score > (SELECT score FROM user_game_score WHERE game_id= ? and user_id = ? and created_at=?) `,
						gameID, nowTody30MinUtc, gameID, selectUserID, nowTody30MinUtc).
						Pluck("rank", &rankIndex).Error
					if err != nil {
						return err
					}
					NowIndex = new(UserGameRank)
					err = tx.Table("user_game_score").Joins("left join users on user_game_score.user_id = users.user_id ").
						Select("user_game_score.user_id,users.name,users.face_url,user_game_score.score,users.token_contract_chain").
						Where("user_game_score.game_id=? and user_game_score.user_id=? and user_game_score.created_at=?  and user_game_score.end_time<>''",
							gameID, selectUserID, nowTody30MinUtc).Find(NowIndex).Error
					if err == nil {
						if NowIndex.UserID == "" {
							NowIndex = nil
						} else {
							NowIndex.Index = int32(rankIndex)
						}

					} else {
						NowIndex = nil
					}
				} else {
					NowIndex = result[rankIndex]
				}
			}
			return err
		}
	})
	return
}

// 不翻页的请求
func GetGameListFromDB(gameID string) ([]*db.GameConfig, error) {
	var result []*db.GameConfig
	var err error
	if gameID != "" {
		err = db.DB.MysqlDB.DefaultGormDB().Table("game_config").Where("game_id=? and status=1", gameID).Order("game_id").Find(&result).Error
	} else {
		err = db.DB.MysqlDB.DefaultGormDB().Table("game_config").Where("status=1").Order("game_id").Find(&result).Error
	}

	return result, err
}
func UpdateGameStatus(UserID, gameName string, gameID, StartFlag int32, ip, useragent string, score float64, timeNow int64) error {
	today := time.Now()
	nowMinute := today.Minute()
	if nowMinute > 30 {
		nowMinute = 30
	} else {
		nowMinute = 0
	}

	nowTody30MinUtc := time.Date(today.Year(), today.Month(), today.Day(), today.Hour(), nowMinute, 0, 0, today.Location())
	var userDataGameStatus db.UserGameScore
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_game_score").
		Where("user_id=? and game_id=?", UserID, gameID).
		Order("created_at desc").First(&userDataGameStatus).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("没有开始的记录")
	} else if errors.Is(err, gorm.ErrRecordNotFound) && StartFlag == 2 {
		return errors.New("未登记开始的 无法结束游戏")
	}
	if StartFlag == 1 {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//插入数据
			return db.DB.MysqlDB.DefaultGormDB().Create(&db.UserGameScore{
				CreatedAt: nowTody30MinUtc,
				UpdatedAt: nowTody30MinUtc,
				UserID:    UserID,
				GameID:    gameID,
				GameName:  gameName,
				Score:     0,
				IP:        ip,
				UserAgent: useragent,
				StartTime: time.Now().UnixMilli(),
				EndTime:   0,
				PlayNum:   1,
			}).Error
		} else {
			//判断最后一天的时间是否是今天
			lastInsertDate := userDataGameStatus.CreatedAt
			daySub := nowTody30MinUtc.Sub(lastInsertDate).Minutes()
			if daySub < 30 && daySub >= 0 {
				err = db.DB.MysqlDB.DefaultGormDB().Table("user_game_score").Where("id= ?", userDataGameStatus.ID).Updates(
					map[string]interface{}{
						"start_time": timeNow,
						"ip":         ip,
						"user_agent": useragent,
						"play_num":   userDataGameStatus.PlayNum + 1}).Error

				//同一天的情况
			} else if daySub >= 30 {
				//大于24小时的情况 需要插入数据了
				err = db.DB.MysqlDB.DefaultGormDB().Create(&db.UserGameScore{
					CreatedAt: nowTody30MinUtc,
					UpdatedAt: nowTody30MinUtc,
					UserID:    UserID,
					GameID:    gameID,
					GameName:  gameName,
					Score:     0,
					IP:        ip,
					UserAgent: useragent,
					StartTime: timeNow,
					EndTime:   0,
					PlayNum:   1,
				}).Error
			}
			if err != nil {
				return err
			}
		}

	} else if StartFlag == 2 { //结束游戏的情况下：
		//判断最后一天的时间是否是今天
		lastInsertDate := userDataGameStatus.CreatedAt
		nowMinute := nowTody30MinUtc.Sub(lastInsertDate).Minutes()
		fmt.Println(userDataGameStatus.IP)
		fmt.Println(ip)
		fmt.Println(userDataGameStatus.UserAgent)
		fmt.Println(useragent)
		//if !strings.EqualFold(userDataGameStatus.IP, ip) || !strings.EqualFold(userDataGameStatus.UserAgent, useragent) {
		//	return errors.New("夸浏览器提交，作弊")
		//}

		if nowMinute < 30 && nowMinute >= 0 {
			if userDataGameStatus.Score > int64(score) {
				err = db.DB.MysqlDB.DefaultGormDB().Table("user_game_score").
					Where("id= ?", userDataGameStatus.ID).Updates(
					map[string]interface{}{
						"end_time": timeNow,
						"score":    userDataGameStatus.Score,
					}).Error
			} else {
				err = db.DB.MysqlDB.DefaultGormDB().Table("user_game_score").
					Where("id= ?", userDataGameStatus.ID).Updates(
					map[string]interface{}{
						"end_time": timeNow,
						"score":    score,
					}).Error
			}

		} else if nowMinute >= 24 {
			//大于24小时的情况 需要插入数据了
			err = db.DB.MysqlDB.DefaultGormDB().Create(&db.UserGameScore{
				CreatedAt: nowTody30MinUtc,
				UpdatedAt: nowTody30MinUtc,
				UserID:    UserID,
				GameID:    gameID,
				GameName:  gameName,
				Score:     int64(score),
				IP:        ip,
				UserAgent: useragent,
				StartTime: userDataGameStatus.StartTime,
				EndTime:   timeNow,
				PlayNum:   1,
			}).Error
		}
		return err
	}
	return nil
}
