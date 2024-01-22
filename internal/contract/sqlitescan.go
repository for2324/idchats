package contract

import (
	"Open_IM/internal/contract/evmlisteninterface/eth"
	"Open_IM/pkg/common/config"
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/supermigo/xlog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

// 本地数据库 扫块的db 只记录到最大的处理块的db 避免数据过大
var (
	LocalDB *gorm.DB
)

func InitLocalDB(str string) {
	dbFileName := str
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		panic("can't save sqlite3 db")
		return
	}
	sqlDB, err := db.DB()
	sqlDB.SetConnMaxLifetime(time.Hour * 1)
	sqlDB.SetMaxOpenConns(3)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	_ = db.AutoMigrate(&eth.EthDefaultInfo{})
	LocalDB = db
	return
}

type EthInfoRepo struct {
}

func (*EthInfoRepo) Get(chainID string) (*eth.EthDefaultInfo, error) {
	var info eth.EthDefaultInfo
	err := LocalDB.First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}
func (*EthInfoRepo) Create(chainID string, info *eth.EthDefaultInfo) error {
	return LocalDB.Create(info).Error
}
func (*EthInfoRepo) Update(chainID string, info *eth.EthDefaultInfo) error {
	return LocalDB.Where("chain_id=?", chainID).Save(info).Error
}

// 建议用rpc 服务来做扫块 不要在api层面去做
func StartScanBlockFilterQuery() {
	ethInfo := &EthInfoRepo{}
	var createEthDefault *eth.EthDefaultInfo
	var err error
	//从第几块扫快
	beginBlock := config.Config.RewardScanBlock
	xlog.CInfo("从第几个快：", beginBlock)
	ethclientPtr := GetRewardRpcClient()
	if ethclientPtr == nil {
		xlog.CError("无法初始化 奖励池的 爬虫奖励")
		return
	}
	chainID, err := ethclientPtr.ChainID(context.Background())
	if err != nil {
		xlog.CError(err.Error())
		return
	} else {
		xlog.CInfo("当前的chain id 为：", chainID.String())
	}

	if createEthDefault, err = ethInfo.Get(chainID.String()); errors.Is(err, gorm.ErrRecordNotFound) {
		createEthDefault = &eth.EthDefaultInfo{
			ChainID:          chainID.String(),
			LastScannedBlock: beginBlock,
		}
		ethInfo.Create(chainID.String(), createEthDefault)
	}
	ethListen := eth.NewEthListenerByFilterQuery(
		[]*ethclient.Client{ethclientPtr},
		3, 3, chainID.String(),
		ethInfo)

	nt := NewRewardManagerContract(config.Config.RewardChainContractAddress, chainID.String(), ethclientPtr)
	ethListen.RegisterConsumer(nt)
	ethListen.Start(context.Background())
}
