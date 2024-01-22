package cronTask

import (
	"Open_IM/internal/api/brc20/services/unisat_wallet"
	"Open_IM/internal/contract"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"context"
	"errors"
	"fmt"
	"github.com/aizuyan/cron/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sourcegraph/conc"
	"math/big"
	"time"
)

// 质押池的Volume 的波动检测
func CheckPledgePoolVolume() {
	c := cron.New(cron.WithSeconds(), cron.WithLogger(cron.DefaultLogger))
	c.AddFuncWithTag("InsertNewBBtPledge", "56 59 23 * * *", func() {
		CheckChainBBTPledgePoolVolume()
	})
	c.AddFuncWithTag("InsertNewBlpPledge", "55 59 23 * * *", func() {
		CheckChainBLpPledgePoolVolume()
	})
	c.AddFuncWithTag("InsertNewBRC20BBTBlpPledge", "55 59 */1 * * *", func() {
		CheckBrc20PledgePoolVolume()
	})
	c.AddFuncWithTag("SyncBrc20TransferScan", "*/20 * * * * *", func() {
		SyncBrc20Transfer()
	})
	c.Start()
}
func SyncBrc20Transfer() {
	tempdata := new(unisat_wallet.UnisatWeb)
	if !config.Config.IsPublicEnv {
		tempdata.NetParam = "testnet"
		tempdata.LastScanTime = 1703217600
		tempdata.Ticker = "dead"
		tempdata.ScanUnisatWallet()
	} else {

		tempdata.NetParam = "mainnet"
		tempdata.LastScanTime = 1703217600
		tempdata.Ticker = "obbt"
		tempdata.ScanUnisatWallet()
	}

}

// 统计每日brc20 的质押
func CheckBrc20PledgePoolVolume() {
	for {
		netParam := "testnet"
		if config.Config.IsPublicEnv {
			netParam = "mainnet"
		}

		unisatPtr := &unisat_wallet.UnisatWeb{NetParam: netParam}
		totalStake, err1 := unisatPtr.GetTotalStake()
		blockHeight, err2 := unisatPtr.GetBlockHeight()
		if err1 == nil && err2 == nil {
			if config.Config.IsPublicEnv {
				InsertIntoBrc20PledgePoolDB("0", blockHeight, uint64(time.Now().Unix()), totalStake, "obbt")
			} else {
				InsertIntoBrc20PledgePoolDB("0", blockHeight, uint64(time.Now().Unix()), totalStake, "dead")
			}

			break
		} else {
			if err1 != nil {
				fmt.Println(err1.Error())
			}
			if err2 != nil {
				fmt.Println(err2.Error())
			}

		}
		time.Sleep(time.Second * 10)
	}

}

// 下面两份代码是一样的，
func CheckChainBLpPledgePoolVolume() {
	for {
		if config.Config.BBTPledge.BLpContractAddress == "" {
			log.NewError("BLpContractAddress", "当前BLpContractAddress的内容为空")
			break
		}
		var concgo conc.WaitGroup
		var chainStr string
		var blockHeight uint64
		var created_at uint64
		var totalLock string
		var err1 error
		var err2 error
		concgo.Go(func() {
			chainStr, blockHeight, created_at, err1 = getBlockChainHeight()
		})
		concgo.Go(func() {
			totalLock, err2 = getBLpPledgeTotalLock()
		})
		concgo.Wait()
		if err1 == nil && err2 == nil {
			InsertIntoBLpPledgePoolDB(chainStr, blockHeight, created_at, totalLock)
			break
		} else {
			if err1 != nil {
				fmt.Println(err1.Error())
			}
			if err2 != nil {
				fmt.Println(err2.Error())
			}
			time.Sleep(time.Second * 12)
		}
	}
}

func CheckChainBBTPledgePoolVolume() {
	for {
		if config.Config.BBTPledge.ContractAddress == "" {
			log.NewError("CheckChainBBTPledgePoolVolume", "当前ContractAddress质押合约的内容为空")
			break

		}
		var concgo conc.WaitGroup
		var chainStr string
		var blockHeight uint64
		var created_at uint64
		var totalLock string
		var err1 error
		var err2 error
		concgo.Go(func() {
			chainStr, blockHeight, created_at, err1 = getBlockChainHeight()
		})
		concgo.Go(func() {
			totalLock, err2 = getBBTPledgeTotalLock()
		})
		concgo.Wait()
		if err1 == nil && err2 == nil {
			InsertIntoBBTPledgePoolDB(chainStr, blockHeight, created_at, totalLock)
			break
		} else {
			if err1 != nil {
				fmt.Println(err1.Error())
			}
			if err2 != nil {
				fmt.Println(err2.Error())
			}
			time.Sleep(time.Second * 12)
		}
	}
}
func InsertIntoBBTPledgePoolDB(chainStr string, blockHeight uint64, created_at uint64, totalLock string) error {
	// 插入数据
	timestamp := int64(created_at) // 替换为你要获取整点时间戳的时间戳
	t := time.Unix(timestamp, 0)
	roundedTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return db.DB.MysqlDB.DefaultGormDB().Create(&db.BbtPledgeLog{
		Chain:       chainStr,
		CreatedAt:   time.Now(),
		BlockDate:   roundedTime,
		Contract:    config.Config.BBTPledge.ContractAddress,
		TotalLock:   totalLock,
		BlockHeight: int64(blockHeight),
	}).Error
}
func InsertIntoBLpPledgePoolDB(chainStr string, blockHeight uint64, created_at uint64, totalLock string) error {
	// 插入数据
	timestamp := int64(created_at) // 替换为你要获取整点时间戳的时间戳
	t := time.Unix(timestamp, 0)
	roundedTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return db.DB.MysqlDB.DefaultGormDB().Create(&db.BLPPledgeLog{
		Chain:       chainStr,
		CreatedAt:   time.Now(),
		BlockDate:   roundedTime,
		Contract:    config.Config.BBTPledge.BLpContractAddress,
		TotalLock:   totalLock,
		BlockHeight: int64(blockHeight),
	}).Error
}

func InsertIntoBrc20PledgePoolDB(chainStr string, blockHeight int64, created_at uint64, totalLock string, ticker string) error {
	// 插入数据
	timestamp := int64(created_at) // 替换为你要获取整点时间戳的时间戳
	t := time.Unix(timestamp, 0)
	roundedTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return db.DB.MysqlDB.DefaultGormDB().Create(&db.ObbtPledgeLogDayReport{
		Chain:       chainStr,
		CreatedAt:   time.Now(),
		BlockDate:   roundedTime,
		Contract:    ticker,
		TotalLock:   totalLock,
		BlockHeight: blockHeight,
	}).Error
}

// 获取高度
func getBlockChainHeight() (chainStr string, blockHeight uint64, timestamp uint64, err error) {
	rpcClient := contract.GetRewardRpcClient()

	if rpcClient != nil {
		chainID, err := rpcClient.ChainID(context.Background())
		if err != nil {
			log.NewError("0", "Serve failed ", err.Error())
			return "", 0, 0, err
		}
		chainStr = config.GetChainName(chainID.Uint64())

		blockDetail, err := rpcClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.NewError("0", "Serve failed ", err.Error())
			return "", 0, 0, err
		}
		return chainStr, blockDetail.Number.Uint64(), blockDetail.Time, nil
	}
	return "", 0, 0, errors.New("未初始化 bbt rpc 连接")
}

// 获取当前totalsupply
func getBBTPledgeTotalLock() (string, error) {
	rpcClient := contract.GetRewardRpcClient()

	if rpcClient == nil || config.Config.BBTPledge.ContractAddress == "" {
		return "", errors.New("未初始化 bbt rpc 连接")
	}
	contractPtr, err := contract.NewBBTPledgePool(
		common.HexToAddress(config.Config.BBTPledge.ContractAddress),
		rpcClient)
	if err != nil {
		return "", errors.New("未初始化 config.Config.BBTPledge.ContractAddress  合约地址有错误")
	}
	totalLock, err := contractPtr.TotalLocked(nil)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return totalLock.String(), nil
}

// 获取当前totalsupply
func getBLpPledgeTotalLock() (string, error) {
	rpcClient := contract.GetRewardRpcClient()

	if rpcClient == nil || config.Config.BBTPledge.BLpContractAddress == "" {
		return "", errors.New("未初始化 bbt rpc 连接")
	}
	contractPtr, err := contract.NewBLPPledgePool(
		common.HexToAddress(config.Config.BBTPledge.BLpContractAddress),
		rpcClient)
	if err != nil {
		return "", errors.New("未初始化 config.Config.BBTPledge.ContractAddress  合约地址有错误")
	}
	poolInfo, err := contractPtr.PoolInfo(&bind.CallOpts{}, new(big.Int).SetInt64(0))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return poolInfo.TotalLocked.String(), nil
}
