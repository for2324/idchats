package brc20

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/db"
	"github.com/shopspring/decimal"
	"time"
)

// 模拟根据存款时间获取利率的逻辑
func GetInterestRates(depositTime int64, stakePeriod string, timeNow time.Time, interestRates []*db.ObbtPoolInfoHistory) []*db.ObbtPoolInfoHistory {
	// 找到所有符合条件的利率记录
	var applicableInterestRates []*db.ObbtPoolInfoHistory
	var PoolInfoRate map[int]struct{}
	PoolInfoRate = make(map[int]struct{}, 0)
	for i := range interestRates {
		if stakePeriod != interestRates[i].StakingPeriod {
			continue
		}
		dbStartTimeStamp := interestRates[i].StartTime.Unix()
		dbEndTimeStamp := interestRates[i].EndTime.Unix()
		if depositTime > dbStartTimeStamp && depositTime > dbEndTimeStamp {
			continue
		}
		//质押时间小于等于开始时间
		if depositTime <= dbStartTimeStamp && depositTime <= dbEndTimeStamp {
			if _, ok := PoolInfoRate[interestRates[i].ID]; !ok {
				//当前时间在这个范围内
				if timeNow.After(interestRates[i].StartTime) || timeNow.Equal(interestRates[i].StartTime) {
					applicableInterestRates = append(applicableInterestRates, interestRates[i])
				}
			}
			continue
		}
		//质押时间大于等于条件开始时间
		if depositTime >= dbStartTimeStamp && depositTime <= dbEndTimeStamp {
			if _, ok := PoolInfoRate[interestRates[i].ID]; !ok {
				applicableInterestRates = append(applicableInterestRates, interestRates[i])
			}
			continue
		}
	}
	return applicableInterestRates
}
func GetPending(stake *api.PersonalBrc20Pledge, rateList []*db.ObbtPoolInfoHistory) (pendingValue string) {
	timeNow := time.Now()
	timeNowStamp := timeNow.Unix()
	if timeNowStamp <= stake.StartTime {
		return "0"
	}
	// 找到所有符合条件的利率记录
	applicableInterestRates := GetInterestRates(stake.StartTime, stake.StakingPeriod, time.Now(), rateList)
	if len(applicableInterestRates) == 0 {
		return
	}
	// 计算利息
	interest := decimal.Zero

	stakeTimeStamp := stake.StartTime
	stakeAmountDecimal, _ := decimal.NewFromString(stake.StakedAmount)
	for i := range applicableInterestRates {
		if stake.StakedAmount == "0" {
			continue
		}
		if timeNowStamp < stakeTimeStamp {
			continue
		}
		if timeNow.Before(applicableInterestRates[i].StartTime) {
			continue
		}
		rewardRateDecimal, _ := decimal.NewFromString(applicableInterestRates[i].RewardsRate)
		startTimeStamp := applicableInterestRates[i].StartTime.Unix()
		endTimeStamp := applicableInterestRates[i].EndTime.Unix()
		if stakeTimeStamp <= startTimeStamp && stakeTimeStamp <= endTimeStamp {
			tempEndData := timeNowStamp
			if timeNowStamp >= endTimeStamp {
				tempEndData = endTimeStamp
			}
			interest = interest.Add(stakeAmountDecimal.Mul(rewardRateDecimal.Mul(decimal.NewFromInt(tempEndData - startTimeStamp)).
				Div(decimal.NewFromInt(360 * 86400))))
		} else if stakeTimeStamp >= startTimeStamp && stakeTimeStamp <= endTimeStamp {
			tempEndData := timeNowStamp
			if timeNowStamp >= endTimeStamp {
				tempEndData = endTimeStamp
			}
			interest = interest.Add(stakeAmountDecimal.Mul(rewardRateDecimal.Mul(decimal.NewFromInt(tempEndData - stakeTimeStamp)).Div(decimal.NewFromInt(360 * 86400))))
		} else if stakeTimeStamp >= startTimeStamp && startTimeStamp >= endTimeStamp {
			continue
		}
	}
	return interest.String()
}
func getStakeRewardRate() (rewardStartTime []*db.ObbtPoolInfoHistory, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("obbt_pool_info_history").
		Where("start_time < ? ", time.Now()).Order("start_time asc").Find(&rewardStartTime).Error
	return
}
