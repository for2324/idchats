package brc20

import (
	"Open_IM/internal/api/brc20/services"
	"Open_IM/internal/api/brc20/services/unisat_wallet"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/algorithm"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/sourcegraph/conc"
	"github.com/sourcegraph/conc/stream"
	"net/http"
	"strings"
	"time"
)

// @Summary		预定质押Obbt
// @Description	Brc20
// @Tags		Brc20相关
// @ID			Brc20Pledge
// @Accept		json
// @Param			token	header	string				true	"im token"
// @Param		req		body	api.Brc20PledgeUtxo{}	true	"请求体RequestType :balance，deploy ,inscription"
// @Produce		json
// @Success		0	{object}	api.Brc20PledgeResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/brc20/brc20_pledge [post]
func Brc20Pledge(c *gin.Context) {
	//获取自己个人的内容
	params := api.Brc20PledgeUtxo{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, _, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "Brc20Pledge failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, utils.GetSelfFuncName(), "args ", convertor.ToString(params))
	if len(params.InscriptionId) == 0 {
		log.NewError("0", "InscriptionId is empty", "")
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "InscriptionId is empty"})
		return
	}
	signData := strings.Join(params.InscriptionId, ",")
	resultSignVerify, err := services.VerifyMessageWithAddr(context.Background(), &services.VerifyMessageWithAddrRequest{
		Msg:       []byte(signData),
		Signature: params.Sign,
		Addr:      params.Address,
		NetParam:  params.NetParams,
	})
	if err != nil || resultSignVerify.Valid == false {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "sign error"})
		return
	}
	//检查每一张
	_, AmountList, err := unisat_wallet.CheckIsHaveThisUtxo(params.Address, params.Ticker, params.InscriptionId, params.NetParams)
	if err == nil {
		//判断这些票据 是否出现过在质押的池里面的，  如果存在那么旧的旧要全部废     弃掉，
		getPledgeOrderListString, err := imdb.GetAllPrePledgeOnlyInscription(params.Address)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
			return
		}
		//返回交集
		sameInscriptionResultArray := slice.Intersection(getPledgeOrderListString, params.InscriptionId)
		var insertIntoDb []*db.ObbtPrePledge
		fmt.Println("相同的数据如下：", convertor.ToString(sameInscriptionResultArray))
		if len(sameInscriptionResultArray) == 0 {
			for key, value := range params.InscriptionId {
				insertIntoDb = append(insertIntoDb, &db.ObbtPrePledge{
					SenderBtcAddress: params.Address,
					InscriptionID:    value,
					Amount:           AmountList[key],
					IsUpToChain:      0,
					StakingPeriod:    params.StakingPeriod,
				})
			}
		} else {
			//删除这些票据
			err = imdb.DeletePrePledgeByInscriptionId(params.Address, sameInscriptionResultArray)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
				return
			}
			for key, value := range params.InscriptionId {
				insertIntoDb = append(insertIntoDb, &db.ObbtPrePledge{
					SenderBtcAddress: params.Address,
					InscriptionID:    value,
					Amount:           AmountList[key],
					IsUpToChain:      0,
					StakingPeriod:    params.StakingPeriod,
				})
			}
		}
		if len(insertIntoDb) > 0 {
			err = imdb.InsertIntoPrePledger(insertIntoDb)
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": err.Error()})
			return
		}
		log.NewInfo(params.OperationID, "请上传数据")
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": ""})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

}

// 获取可以用来转账的票据：
// @Summary		获取可以用来转账的票据
// @Description	Brc20
// @Tags		Brc20相关
// @ID			GetBrc20TransferableList
// @Accept		json
// @Param			token	header	string				true	"im token"
// @Param		req		body	api.Brc20PledgeReq{}	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.Brc20TransferableListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/brc20/brc20_transferable_list [post]
func GetBrc20TransferableList(c *gin.Context) {
	//获取自己个人的内容
	params := api.Brc20PledgeReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, _, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetBrc20TransferableList failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, utils.GetSelfFuncName(), "args ", convertor.ToString(params))
	if params.Ticker == "" {
		log.NewError(params.OperationID, params.Address+"transferable ticker is empty", "")
		c.JSON(http.StatusOK, gin.H{"errCode": 1001, "errMsg": params.Address + " transferable ticker is empty"})
		return
	}
	beginIndex := int64(0)
	pageSize := int64(100)
	var rebackResul []*services.TransferAbleInscript
	for {
		serverTransferABleInscriptSub, err := unisat_wallet.NetGetAddressBrc20TransableListResult(params.Address,
			params.Ticker, params.NetParams, beginIndex, pageSize)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		if len(serverTransferABleInscriptSub) > 0 {
			rebackResul = append(rebackResul, serverTransferABleInscriptSub...)
		}
		if len(serverTransferABleInscriptSub) == 0 {
			break
		}
		if int64(len(serverTransferABleInscriptSub)) == pageSize {
			beginIndex = beginIndex + pageSize
			continue
		} else {
			break
		}
	}
	if len(rebackResul) == 0 {
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": []string{}})

		return
	} else {
		c.JSON(http.StatusOK, &api.Brc20TransferableListResp{
			Data: rebackResul,
		})
		return
	}
}

// 获取个人obbt质押的信息
// @Summary		获取个人obbt质押的信息
// @Description	Brc20
// @Tags		Brc20相关
// @ID			Brc20PledgePersonalInfoes
// @Accept		json
// @Param			token	header	string				true	"im token"
// @Param		req		body	api.Brc20PledgeReq{}	true	"请求体只有一个address有用"
// @Produce		json
// @Success		0	{object}	api.PersonalBrc20PledgeResonse
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/brc20/brc20_personal_pledge_infoes [post]
func Brc20PledgePersonalInfoes(c *gin.Context) {
	//获取自己个人的内容
	params := api.Brc20PledgeReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, _, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "Brc20PledgePersonalInfoes failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, utils.GetSelfFuncName(), "args ", convertor.ToString(params))

	//获取自己个人的内容
	var dbdata []*db.ObbtStake
	err := db.DB.MysqlDB.DefaultGormDB().Table("obbt_stake").Where("sender_btc_address=?", params.Address).
		Find(&dbdata).Error
	rebackData := new(api.PersonalBrc20PledgeResonse)
	dbListRewardHistory, err := getStakeRewardRate()
	if err == nil && len(dbdata) > 0 {
		streamSPtr := stream.New()
		for _, value := range dbdata {
			tempValueData := value
			streamSPtr.Go(func() stream.Callback {
				tempValue := &api.PersonalBrc20Pledge{
					StartTime:     tempValueData.StartTime.Unix(),
					StakedAmount:  tempValueData.StakeAmount,
					Reward:        "",
					StakingPeriod: tempValueData.StakingPeriod,
				}
				rewardValue := GetPending(tempValue, dbListRewardHistory)
				tempValue.Reward = rewardValue
				return func() {
					rebackData.Data = append(rebackData.Data, tempValue)
				}
			})
		}
		streamSPtr.Wait()
		log.NewError(params.OperationID, "reward 暂时没有计算")
		c.JSON(http.StatusOK, rebackData)
	} else {
		c.JSON(http.StatusOK, rebackData)
	}
}

type rewardRateComparator struct{}

func (c *rewardRateComparator) Compare(v1 any, v2 any) int {
	val1, _ := v1.(*api.Brc20PoolInfo)
	val2, _ := v2.(*api.Brc20PoolInfo)
	float64i, _ := convertor.ToFloat(val1.RewardsRate)
	float64j, _ := convertor.ToFloat(val2.RewardsRate)
	if float64i >= float64j {
		return 1
	}
	return -1
}

// 获取obbt质押池的信息：
// @Summary		获取obbt质押池的信息
// @Description	Brc20
// @Tags		Brc20相关
// @ID			GetBrc20PledgePoolInfoes
// @Accept		json
// @Produce		json
// @Success		0	{object}	api.Brc20TransferableListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/brc20/brc20_pledge_pool_infoes [get]
func GetBrc20PledgePoolInfoes(c *gin.Context) {
	//获取自己个人的内容
	var dbdata []*db.ObbtPoolInfo
	timeNowFormtValue := time.Now().Format(time.DateTime)
	err := db.DB.MysqlDB.DefaultGormDB().Raw(`select obbt_pool_info_history.*,obbt_pool_info.total_stake from obbt_pool_info left join obbt_pool_info_history on obbt_pool_info.pool_id = obbt_pool_info_history.pool_id
		and obbt_pool_info.staking_period = obbt_pool_info_history.staking_period
		where obbt_pool_info_history.start_time <=? and ?<=obbt_pool_info_history.end_time`, timeNowFormtValue, timeNowFormtValue).Find(&dbdata).Error
	if err == nil {
		rebackData := new(api.Brc20ObbtPledgePoolInfoes)
		for _, value := range dbdata {
			rebackData.Data = append(rebackData.Data, &api.Brc20PoolInfo{
				PoolID:        value.PoolID,
				StakingPeriod: value.StakingPeriod,
				RewardsRate:   value.RewardsRate,
				TotalStake:    value.TotalStake,
			})
		}
		if len(rebackData.Data) > 0 {
			algorithm.QuickSort(rebackData.Data, &rewardRateComparator{})
		}
		c.JSON(http.StatusOK, rebackData)
	} else {
		c.JSON(http.StatusOK, gin.H{"errCode": 1002, "errMsg": "pool not init "})
	}
}

// 查询brc20的信息
// @Summary		查询brc20的信息
// @Description	Brc20
// @Tags		Brc20相关
// @ID			Brc20SearchScan
// @Accept		json
// @Param		token	header	string				true	"im token"
// @Param		req		body	api.Brc20PledgeReq{}	true	"请求体只有一个address有用"
// @Produce		json
// @Success		0	{object}	api.Brc20InfoForSearchResponse
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/brc20/brc20_search [post]
func Brc20SearchScan(c *gin.Context) {
	//获取自己个人的内容
	params := api.Brc20PledgeReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, _, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "Brc20PledgePersonalInfoes failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, utils.GetSelfFuncName(), "args ", convertor.ToString(params))
	if params.Ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "scan ticker is empty"})
		return
	}
	reponseData := new(api.Brc20InfoForSearchResponse)
	returnList, err := getTickerListFromOkx(params.OperationID, params.Ticker)
	if err == nil {
		for _, value := range returnList {
			reponseData.Data = append(reponseData.Data, &api.Brc20InfoForSearch{
				Ticker:        value.Ticker,
				InscriptionID: value.InscriptionId,
				Creator:       value.Creator,
				UsdFloorPrice: value.UsdFloorPrice,
				HolderCount:   value.Holders,
			})
		}
	}
	c.JSON(http.StatusOK, reponseData)
}
func getTickerListFromOkx(operationID string, ticker string) ([]*TokSearchTickerInfo, error) {
	var tokenList TokSearchTickerResponse[TokSearchTickerInfo]
	code := 0

	requestUrl := fmt.Sprintf("https://www.okx.com/priapi/v1/nft/brc/tokens?t=%v&scope=1&page=1&size=20&sortBy=deployedTime&sort=asc&tokensLike=%v&timeType=1&walletAddress=", time.Now().UnixMilli(), ticker)
	fmt.Println("requestUrl:", requestUrl)
	dataForm := gout.GET(requestUrl).
		SetHeader(gout.H{"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"}).
		BindJSON(&tokenList).Code(&code)
	if !config.Config.IsPublicEnv {
		dataForm = dataForm.SetProxy("http://proxy.idchats.com:7890").Debug(true)
	}
	err := dataForm.Do()
	if err == nil && code == 200 && len(tokenList.Data.List) > 0 {
		var concWait conc.WaitGroup
		for key, _ := range tokenList.Data.List {
			fmt.Println("check ticker name :", tokenList.Data.List[key].Ticker)
			tempData := tokenList.Data.List[key]
			// 应该要使用有缓存池的方式去获取内容，但是因为项目赶工，那就让代码先上再说， 反正也就是20个携程
			concWait.Go(func() {
				tempDataInfo, err := getTokenPriceUsdPoolPrice(operationID, tempData.Ticker)
				if err == nil {
					tempData.UsdFloorPrice = tempDataInfo.Data.UsdFloorPrice
					tempData.InscriptionId = tempDataInfo.Data.InscriptionId
					tempData.Creator = tempDataInfo.Data.Creator
				} else {
					log.NewError(operationID, err.Error())
				}
			})
		}
		concWait.Wait()
		return tokenList.Data.List, nil
	}
	log.NewError(operationID, "无法获取该"+ticker+"的地板价格")
	return nil, errors.New("列表不对")

}

func getTokenPriceUsdPoolPrice(operationID string, ticker string) (tokenPriceInfo *TTickerPriceInfoResponse, err error) {
	code := 0
	dataForm := gout.GET("https://www.okx.com/priapi/v1/nft/brc/tokens/" + ticker).
		SetHeader(gout.H{"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"}).
		BindJSON(&tokenPriceInfo).Code(&code)
	if !config.Config.IsPublicEnv {
		dataForm = dataForm.SetProxy("http://proxy.idchats.com:7890").Debug(true)
	}
	err = dataForm.Do()
	if err == nil && code == 200 && tokenPriceInfo.Code == 0 {
		return tokenPriceInfo, nil
	}
	log.NewError(operationID, "无法获取该"+ticker+"的地板价格")
	return nil, errors.New("无法获取该ticker 的地板价:" + ticker)
}

type TokSearchTickerInfo struct {
	ConfirmedMinted      string `json:"confirmedMinted"`
	ConfirmedMinted1H    string `json:"confirmedMinted1h"`
	ConfirmedMinted24H   string `json:"confirmedMinted24h"`
	DeployedTime         string `json:"deployedTime"`
	Holders              string `json:"holders"`
	HoldersChangeRate    string `json:"holdersChangeRate"`
	HoldersChangeRate24H string `json:"holdersChangeRate24H"`
	Image                string `json:"image"`
	IsCollect            int    `json:"isCollect"`
	LimitPerMint         string `json:"limitPerMint"`
	MarketCap            string `json:"marketCap"`
	Price                string `json:"price"`
	PriceChangeRate      string `json:"priceChangeRate"`
	PriceChangeRate24H   string `json:"priceChangeRate24H"`
	Supply               string `json:"supply"`
	Ticker               string `json:"ticker"`
	TickerId             string `json:"tickerId"`
	TotalMinted          string `json:"totalMinted"`
	Transactions         string `json:"transactions"`
	UsdMarketCap         string `json:"usdMarketCap"`
	UsdPrice             string `json:"usdPrice"`
	UsdVolume            string `json:"usdVolume"`
	Volume               string `json:"volume"`
	VolumeChangeRate     string `json:"volumeChangeRate"`
	VolumeCurrencyUrl    string `json:"volumeCurrencyUrl"`
	Creator              string `json:"creator"`
	Decimal              int    `json:"decimal"`
	DiscordUrl           string `json:"discordUrl"`
	FloorPrice           string `json:"floorPrice"`
	InscriptionId        string `json:"inscriptionId"`
	InscriptionNumEnd    string `json:"inscriptionNumEnd"`
	InscriptionNumStart  string `json:"inscriptionNumStart"`
	Limit                string `json:"limit"`
	TwitterUrl           string `json:"twitterUrl"`
	UsdFloorPrice        string `json:"usdFloorPrice"`
	UsdVolumeIn24H       string `json:"usdVolumeIn24h"`
	VolumeIn24H          string `json:"volumeIn24h"`
}
type TokSearchTickerResponse[T any] struct {
	Code int `json:"code"`
	Data struct {
		Data     []interface{} `json:"data"`
		List     []*T          `json:"list"`
		PageNum  int           `json:"pageNum"`
		PageSize int           `json:"pageSize"`
		Total    int           `json:"total"`
	} `json:"data"`
	DetailMsg    string `json:"detailMsg"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Msg          string `json:"msg"`
}
type TTickerPriceInfoResponse struct {
	Code         int                  `json:"code"`
	Data         *TokSearchTickerInfo `json:"data"`
	DetailMsg    string               `json:"detailMsg"`
	ErrorCode    string               `json:"error_code"`
	ErrorMessage string               `json:"error_message"`
	Msg          string               `json:"msg"`
}
