package main

import (
	"Open_IM/cmd/userscoreupload/rediscron"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	utils2 "Open_IM/pkg/utils"

	"gorm.io/gorm"

	pbScore "Open_IM/pkg/proto/score"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type crontabServer struct {
	tracer trace.Tracer
	cron   *rediscron.Cron
}

func NewServer() *crontabServer {
	tracer := otel.Tracer("crontab")
	c, err := rediscron.NewRedisMutexBuilder(rediscron.Config{
		Address: config.Config.Redis.DBAddress[0],
		Secret:  config.Config.Redis.DBPassWord,
		DB:      0,
	})

	if err != nil {
		panic(any(err))
	}
	var redisCron *rediscron.Cron
	redisCron, err = rediscron.New(rediscron.WithRedisMutexBuilder(c))
	if err != nil {
		panic(any(err))
	}

	redisCron.Entries()
	return &crontabServer{
		tracer: tracer,
		cron:   redisCron,
	}
}

func (s *crontabServer) Start(ctx context.Context) (err error) {
	// err = s.cron.AddJob(s.ExecutionSyncSource())
	// if err != nil {
	// 	return err
	// }
	// err = s.cron.AddJob(s.ExecutionSyncSourceWithChat())
	// if err != nil {
	// 	return err
	// }
	err = s.cron.AddJob(s.ExecutionSyncSourceWithTaskReward())
	if err != nil {
		return err
	}
	//每天充值500万积分
	err = s.cron.AddJob(s.ExecutionRechargeEveryGameTotalReword())
	if err != nil {
		return err
	}
	//计算游戏积分
	err = s.cron.AddJob(s.ExecutionGameSource())
	if err != nil {
		return err
	}

	s.cron.Start(ctx)
	<-ctx.Done()
	return
}

func (s *crontabServer) Stop() {
	s.cron.Stop()

}

func (s *crontabServer) ExecutionSyncSourceWithTaskReward() rediscron.Job {
	return rediscron.Job{
		Name:   "TaskRewardUpload",
		Rhythm: "*/10 * * * * ?",
		Func: func(ctx context.Context) error {
			err := s.EventTaskReward(ctx)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
func (s *crontabServer) ExecutionRechargeEveryGameTotalReword() rediscron.Job {
	return rediscron.Job{
		Name:   "RechargeTotalReword",
		Rhythm: "30 0 0 * * ?",
		Func: func(ctx context.Context) error {
			err := s.RechargeTotalReword(ctx)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

// func (s *crontabServer) ExecutionSyncSource() rediscron.Job {
// 	return rediscron.Job{
// 		Name:   "renwuguanzhu",
// 		Rhythm: "*/10 * * * * ?",
// 		Func: func(ctx context.Context) error {
// 			err := s.EventLogNotChat(ctx)
// 			if err != nil {
// 				return err
// 			}
// 			return nil
// 		},
// 	}
// }

func (s *crontabServer) ExecutionGameSource() rediscron.Job {
	pbjob := rediscron.Job{
		Name:   "gamescorelogic",
		Rhythm: "0 1,31 * * * ?",
		Func: func(ctx context.Context) error {
			err := s.EventInsertIntoRecord(ctx)
			if err != nil {
				return err
			}
			return nil
		},
	}
	if config.Config.IsPublicEnv {
		pbjob.Rhythm = "0 1,31 * * * ?"
	}
	return pbjob
}

func (s *crontabServer) EventTaskReward(ctx context.Context) error {
	resultList, err := imdb.GetUnSyncUserTaskRewards(50)
	if err != nil {
		return err
	}
	for _, event := range resultList {
		if err := s.UploadUserRewardEvent(event); err != nil {
			log.Error("UploadUserRewardEvent err", err)
		}
	}
	return nil
}
func (s *crontabServer) RechargeTotalReword(ctx context.Context) error {
	imdb.EveryRechargeGameReword()
	return nil
}
func (s *crontabServer) EventInsertIntoRecord(ctx context.Context) error {
	today := time.Now()
	nowMinute := today.Minute()
	oldHource := today.Add(-1 * time.Hour).Hour()
	if nowMinute > 30 {
		nowMinute = 0
		oldHource = today.Hour()
	} else {
		nowMinute = 30
	}
	yesterdayUtc := time.Date(today.Year(), today.Month(), today.Day(), oldHource, nowMinute, 0, 0, today.Location())
	operationID := utils2.OperationIDGenerator()
	dbgameList, _ := imdb.GetGameListFromDB("")
	for _, valueGameInfo := range dbgameList {
		if valueGameInfo.GameCurrentPrizePool >= valueGameInfo.GameMinPrizePool { //如果是超过这个积分 就需要分配了
			db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
				var db100data []*db.UserGameScore
				err := tx.Table("user_game_score").
					Where("game_id=? and created_at =? and score<>0", valueGameInfo.GameId, yesterdayUtc).
					Order("score desc").Limit(100).Find(&db100data).Error
				if err == nil {
					if len(db100data) > 0 {
						var inserData []*db.RewardEventLogs
						totalSumScore := uint64(0)
						joinCount := 0
						for key, value := range db100data {
							str := fmt.Sprintf(`{"type":"game","info":"在游戏%d里面跑了%d有效分,获得第%dming锁获得奖励"}`, valueGameInfo.GameId, value.Score, key+1)
							makeid := fmt.Sprintf("game:%d", valueGameInfo.GameId) + utils2.Md5(value.UserID+today.Format("2006-01-02 15:04:05"))
							if config.Config.IsPublicEnv {
								makeid = fmt.Sprintf("game:%d", valueGameInfo.GameId) + utils2.Md5(value.UserID+yesterdayUtc.Format("2006-01-02 15:04:05"))
							}
							newRewardEventLogs := &db.RewardEventLogs{
								ID:         makeid,
								UserID:     value.UserID,
								RewardType: fmt.Sprintf("game:%d", valueGameInfo.GameId),
								Reward:     uint64(config.GetRewordFromIndex(int32(key) + 1)),
								UserJSON:   str,
								IsSync:     0,
								CreatedAt:  time.Now(),
							}
							inserData = append(inserData, newRewardEventLogs)
							totalSumScore += newRewardEventLogs.Reward
							joinCount++
						}
						if joinCount > 0 {
							log.NewInfo(operationID, yesterdayUtc, "参与人数", joinCount, "发放奖励", totalSumScore)
							if err = tx.Table("reward_event_logs").Create(inserData).Error; err != nil {
								log.NewInfo(operationID, err, inserData)
								return err
							}
							err = tx.Table("game_config").Where("game_id=?", valueGameInfo.GameId).Update("game_current_prize_pool",
								gorm.Expr("game_current_prize_pool-?", totalSumScore)).Error
							if err != nil {
								return err
							}
							err = tx.Table("game_score_log").Create(&db.GameScoreLog{
								GameID:               int(valueGameInfo.GameId),
								CreatedAt:            time.Now(),
								UpdatedAt:            time.Now(),
								GameCurrentPrizePool: valueGameInfo.GameCurrentPrizePool,
								GameAddPrizePool:     -int64(totalSumScore),
								JoinGameNumber:       joinCount,
								OperatorIp:           "第" + yesterdayUtc.Format("2006-01-02 15:04:05") + "发放奖励",
							}).Error
						}

						return err
					}

				}
				return err
			})
		} else {
			log.NewInfo(operationID, "游戏id", valueGameInfo.GameId, "当前积分是：", valueGameInfo.GameCurrentPrizePool, "不够发放:")
		}
	}
	return nil
}

type PostEvent struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    interface{}
}

type UserJSON struct {
	Info string `json:"info"`
}

func (s *crontabServer) UploadUserRewardEvent(eventLog *db.RewardEventLogs) error {
	info := eventLog.UserJSON
	userInfo := UserJSON{}
	err := json.Unmarshal([]byte(info), &userInfo)
	if err == nil {
		info = userInfo.Info
	}
	reqPb := &pbScore.UploadUserRewardEventReq{
		OperationID: "UploadUserRewardEvent",
		Id:          eventLog.ID,
		UserId:      eventLog.UserID,
		RewardType:  eventLog.RewardType,
		Reward:      int64(eventLog.Reward),
		Info:        info,
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.UserScoreName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		return errors.New(errMsg)
	}
	client := pbScore.NewScoreServiceClient(etcdConn)
	_, err = client.UploadUserRewardEvent(context.Background(), reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, "UploadUserRewardEvent", err.Error(), reqPb.String())
		return err
	}
	return imdb.UpdateRewardEventLogSync(eventLog.ID)
}
