package unisat_wallet

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/guonaihong/gout"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"k8s.io/utils/strings/slices"
	"strings"
	"time"
)

type UnisatWeb struct {
	LastScanTime int
	NetParam     string
	Ticker       string
}
type StakeInfoByAddresst struct {
	StakeAmount   string //无效的量
	StakingPeriod string //质押天数
	StartTime     int64
	InscriptionID string
}
type TTickerPriceFromOkxData struct {
	ConfirmedMinted      string `json:"confirmedMinted"`
	ConfirmedMinted1H    string `json:"confirmedMinted1h"`
	ConfirmedMinted24H   string `json:"confirmedMinted24h"`
	Creator              string `json:"creator"`
	Decimal              int    `json:"decimal"`
	DeployedTime         string `json:"deployedTime"`
	DiscordUrl           string `json:"discordUrl"`
	FloorPrice           string `json:"floorPrice"`
	Holders              string `json:"holders"`
	HoldersChangeRate24H string `json:"holdersChangeRate24H"`
	Image                string `json:"image"`
	InscriptionId        string `json:"inscriptionId"`
	InscriptionNumEnd    string `json:"inscriptionNumEnd"`
	InscriptionNumStart  string `json:"inscriptionNumStart"`
	IsCollect            int    `json:"isCollect"`
	Limit                string `json:"limit"`
	MarketCap            string `json:"marketCap"`
	PriceChangeRate24H   string `json:"priceChangeRate24H"`
	Supply               string `json:"supply"`
	Ticker               string `json:"ticker"`
	TotalMinted          string `json:"totalMinted"`
	Transactions         string `json:"transactions"`
	TwitterUrl           string `json:"twitterUrl"`
	UsdFloorPrice        string `json:"usdFloorPrice"`
	UsdMarketCap         string `json:"usdMarketCap"`
	UsdVolume            string `json:"usdVolume"`
	UsdVolumeIn24H       string `json:"usdVolumeIn24h"`
	Volume               string `json:"volume"`
	VolumeCurrencyUrl    string `json:"volumeCurrencyUrl"`
	VolumeIn24H          string `json:"volumeIn24h"`
}

func LoadFromRedisLastSyncTime() time.Time {
	return time.Now()
}
func (unsat *UnisatWeb) GetBlockHeight() (int64, error) {
	//获取当前的有效事件
	blockHeight := ""
	dataFrom := gout.GET(mempoolHost(unsat.NetParam) + "/blocks/tip/height").BindBody(&blockHeight)
	if !config.Config.IsPublicEnv {
		dataFrom = dataFrom.SetProxy("http://proxy.idchats.com:7890")
		dataFrom = dataFrom.Debug(true)
	}
	err := dataFrom.F().Retry().Attempt(3).Do()
	//err := gout.GET(mempoolHost(unsat.NetParam) + "/blocks/tip/height").BindBody(&blockHeight).F().Retry().Attempt(3).Do()
	if err != nil {
		return 0, err
	}
	return convertor.ToInt(blockHeight)
}

// 获取每个180，360，1080的质押量
func (unsat *UnisatWeb) GetTotalStake() (stakeInfoAmount string, err error) {
	var dbData []*db.ObbtPoolInfo
	err = db.DB.MysqlDB.DefaultGormDB().Table("obbt_pool_info").Where("pool_id=1").Find(&dbData).Error
	totalStake := decimal.Zero
	if err == nil {
		for _, value := range dbData {
			poolStakeInfo, _ := decimal.NewFromString(value.TotalStake)
			totalStake = totalStake.Add(poolStakeInfo)
		}
	}
	return totalStake.String(), err
}

func (unsat *UnisatWeb) getUserPreTotal(btcString string) (err error, result []*StakeInfoByAddresst) {
	var listUserStake []*db.ObbtStake
	err = db.DB.MysqlDB.DefaultGormDB().Table("obbt_stake").Where("sender_btc_address = ?", btcString).Find(&listUserStake).Error
	return
}

// 从unisat 上获取转账的内容
func (unsat *UnisatWeb) ScanUnisatWallet() {
	//https://api-testnet.unisat.io/query-v4/brc20/64656164/history?start=0&limit=20&type=transfer
	start := 0
	limit := 20
	for {
		var dbLastInsertData db.ObbtReciveFromAPI
		err := db.DB.MysqlDB.DefaultGormDB().Table("obbt_recive_from_api").Order("blocktime desc").First(&dbLastInsertData).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("无法查询sql obbt_recive_from_api 的数据 ")
			return
		}
		//查询出该事件含有的铭文的inscripttion_id
		var oldInscriptionID []string
		//代表相同时间内的不能重复插入的铭文id
		if dbLastInsertData.Blocktime != 0 {
			err := db.DB.MysqlDB.DefaultGormDB().Table("obbt_recive_from_api").Where("blocktime=?", dbLastInsertData.Blocktime).Pluck("inscription_id",
				&oldInscriptionID).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println("无法查询sql obbt_recive_from_api 的数据 ")
				return
			}
		}
		fmt.Println("请求服务器数据内容1")
		var tempData TUnisatApiQueryInsertInto
		apiHost, _ := unisatHostApi(unsat.NetParam)
		fmt.Println(apiHost + fmt.Sprintf("/v1/indexer/brc20/%v/history?type=transfer&start=%v&limit=%v",
			hex.EncodeToString([]byte(unsat.Ticker)), start, limit))
		dataForma := gout.GET(apiHost + fmt.Sprintf("/v1/indexer/brc20/%v/history?type=transfer&start=%v&limit=%v",
			hex.EncodeToString([]byte(unsat.Ticker)), start, limit)).SetHeader(gout.H{
			"Authorization": "Bearer 4afe26f1ae562b3b642d64a3ed448904b4e59d90ae435a948078bd9b80ed52a8",
		}).BindJSON(&tempData)
		if !config.Config.IsPublicEnv {
			dataForma = dataForma.Debug(true)
			dataForma = dataForma.SetProxy(config.Config.OpenNetProxy.ProxyURL)
		}
		err = dataForma.Do()
		fmt.Println("请求服务器数据内容2")
		if err == nil && tempData.Code == 0 && len(tempData.Data.Detail) > 0 {
			var tempListInscriptionID []string
			tempListInscriptionIDMap := make(map[string]int64, 0)
			//从数组中筛选出来要插入的数据
			var insertIntoObbtReciveFromApi []*db.ObbtReciveFromAPI
			for _, v := range tempData.Data.Detail {
				if v.Blocktime >= unsat.LastScanTime && //时间必须币数据库最新一条的时候大 或者想等，因为可能同一个时间段多条插入
					v.Blocktime >= dbLastInsertData.Blocktime && // 同时也要判断事件待遇数据库最新一条的时间大
					strings.EqualFold(v.To, config.Config.ReceiveBtcPledge) && //必须是给指定的钱包打钱
					v.Valid && !strings.EqualFold(v.From, v.To) { //必须是有效的铭文，且质押的钱包不能给自己转
					if !slices.Contains(oldInscriptionID, v.InscriptionId) { //事件想等的时候 不能拿同一条来插入
						tempListInscriptionID = append(tempListInscriptionID, v.InscriptionId)
						tempListInscriptionIDMap[v.InscriptionId] = int64(v.Blocktime)
						insertData := new(db.ObbtReciveFromAPI)
						insertData.Blockhash = v.Blockhash
						insertData.Blocktime = v.Blocktime
						insertData.Height = v.Height
						insertData.InscriptionID = v.InscriptionId
						insertData.InscriptionNumber = v.InscriptionNumber
						insertData.Txid = v.Txid
						insertData.Type = v.Type
						insertData.Vout = v.Vout
						insertData.Ticker = v.Ticker
						insertData.SenderBtcAddress = v.From
						insertData.ReceiveBtcAddress = v.To
						insertData.Satoshi = convertor.ToString(v.Satoshi)
						insertData.Amount = v.Amount
						insertData.Height = v.Height
						insertData.Txidx = v.Txidx
						insertData.Blockhash = v.Blockhash
						insertData.Blocktime = v.Blocktime
						insertIntoObbtReciveFromApi = append(insertIntoObbtReciveFromApi, insertData)
					}
				}
			}
			if len(insertIntoObbtReciveFromApi) > 0 { //如果有新增的数据的情况下：插入到数据库
				err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
					//将有效的转账写入到数据
					err = tx.Table("obbt_recive_from_api").Create(&insertIntoObbtReciveFromApi).Error
					if err != nil {
						return err
					}
					//查询质押池的信息
					var dbPoolInfo []*db.ObbtPoolInfo
					mapStakePeriod := make(map[string]*db.ObbtPoolInfo, 0)
					err = tx.Table("obbt_pool_info").Where("pool_id=1").Find(&dbPoolInfo).Error
					if err == nil {
						for _, value := range dbPoolInfo {
							mapStakePeriod[value.StakingPeriod] = value
						}
					}
					//循环有效票据检查是否是用户质押的票据
					for _, willInsertInscription := range insertIntoObbtReciveFromApi {
						var dbDataObbtPrePledge db.ObbtPrePledge
						err = tx.Table("obbt_pre_pledge").Where("inscription_id =? ", willInsertInscription.InscriptionID).
							First(&dbDataObbtPrePledge).Error
						if err == nil {
							//如果是用户预定质押的票据 那么旧统计到质押库里面
							err = tx.Table("obbt_stake").Create(&db.ObbtStake{
								SenderBtcAddress: willInsertInscription.SenderBtcAddress,
								PoolID:           1,
								StakeAmount:      willInsertInscription.Amount,
								StakingPeriod:    dbDataObbtPrePledge.StakingPeriod,
								StartTime:        utils.UnixSecondToTime(int64(willInsertInscription.Blocktime)),
								InscriptionID:    willInsertInscription.InscriptionID,
							}).Error
							//删除该预定质押的票据
							err = tx.Table("obbt_pre_pledge").Delete("inscription_id =? ", willInsertInscription.InscriptionID).Error
							if _, ok := mapStakePeriod[dbDataObbtPrePledge.StakingPeriod]; ok {
								oldTotalAmount, _ := decimal.NewFromString(mapStakePeriod[dbDataObbtPrePledge.StakingPeriod].TotalStake)
								updateTotalAmount, _ := decimal.NewFromString(willInsertInscription.Amount)
								mapStakePeriod[dbDataObbtPrePledge.StakingPeriod].TotalStake = oldTotalAmount.Add(updateTotalAmount).String()
							}
						}
						for _, value := range mapStakePeriod {
							err = db.DB.MysqlDB.DefaultGormDB().Table("obbt_pool_info").Where("pool_id=1 and ticker_name=? and staking_period=?",
								unsat.Ticker, value.StakingPeriod).Updates(map[string]interface{}{
								"total_stake": value.TotalStake,
							}).Error
							if err != nil {
								return err
							}
						}
					}
					return err
				})
				if err != nil {
					fmt.Println("无法插入到数据库")
					return
				}
			}
			if len(tempData.Data.Detail) < limit {
				break
			}
			//如果我最后一条的事件小雨最后一条插入的事件，那么也不需要统计了 「
			if len(tempData.Data.Detail) > 0 {
				if tempData.Data.Detail[len(tempData.Data.Detail)-1].Blocktime < dbLastInsertData.Blocktime ||
					tempData.Data.Detail[len(tempData.Data.Detail)-1].Blocktime < unsat.LastScanTime {
					break
				}
			}
			start = start + limit
			time.Sleep(time.Second * 2)
			continue
		} else if tempData.Code == 0 && len(tempData.Data.Detail) == 0 {
			break
		}
	}
	fmt.Println("同步转账信息成功")
}

type TUnisatApiQueryInsertInto struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Height int `json:"height"`
		Total  int `json:"total"`
		Start  int `json:"start"`
		Detail []struct {
			Ticker            string `json:"ticker"`
			Type              string `json:"type"`
			Valid             bool   `json:"valid"`
			Txid              string `json:"txid"`
			Idx               int    `json:"idx"`
			Vout              int    `json:"vout"`
			InscriptionNumber int    `json:"inscriptionNumber"`
			InscriptionId     string `json:"inscriptionId"`
			From              string `json:"from"`
			To                string `json:"to"`
			Satoshi           int    `json:"satoshi"`
			Amount            string `json:"amount"`
			OverallBalance    string `json:"overallBalance"`
			TransferBalance   string `json:"transferBalance"`
			AvailableBalance  string `json:"availableBalance"`
			Height            int    `json:"height"`
			Txidx             int    `json:"txidx"`
			Blockhash         string `json:"blockhash"`
			Blocktime         int    `json:"blocktime"`
		} `json:"detail"`
	} `json:"data"`
}
