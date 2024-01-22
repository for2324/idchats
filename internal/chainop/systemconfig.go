package chainop

import (
	"Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/token_verify"
	kafkaMessage "Open_IM/pkg/proto/kafkamessage"
	"Open_IM/pkg/utils"
	"Open_IM/pkg/xkafka"
	"Open_IM/pkg/xlog"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	erc1155 "github.com/nattaponra/go-abi/erc1155/contract"

	"github.com/Pallinder/go-randomdata"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	erc721 "github.com/nattaponra/go-abi/erc721/contract"
)

var ct *xkafka.Kafka

func CloseKafka() {
	if ct != nil && ct.Producer != nil {
		ct.Producer.Close()
	}
}
func InstallKafkaProduct() {
	optptr := xkafka.NewDefaultOptions()
	optptr.Name = "product"
	optptr.Addr = config.Config.Kafka.BusinessTop.Addr
	var err error
	ct, err = xkafka.New(optptr, nil, nil)
	if err != nil {
		xlog.CErrorf(err.Error())
		return
	}
}

func GetWhiteList(ctx *gin.Context) {
	if result, err := imdb.GetWhiteList(ctx); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"errCode": 40000, "errMsg": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": result})
	}
}

type PostEmailVeryCode struct {
	EmailAddress      string  `json:"emailAddress"`
	Code              string  `json:"code"`
	Password          *string `json:"password"`
	EncryptPrivateKey *string `json:"encryptPrivateKey"`
}

func GetEmailCode(c *gin.Context) {
	//让其生产一个邮件给发送
	var ask PostEmailVeryCode
	if err := c.ShouldBind(&ask); err != nil || ask.EmailAddress == "" ||
		!utils.VerifyEmailFormat(ask.EmailAddress) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parame id error"})
		return
	}
	_, requestRedisDbID, _ := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), utils.OperationIDGenerator())
	requestRedisDb := ask.EmailAddress
	if requestRedisDbID != "" {
		requestRedisDb += ":" + requestRedisDbID
	}
	requestRedisDb = db.EmailVerifyCodePrefix + "_" + requestRedisDb
	sendEmailCode(c, requestRedisDb, ask.EmailAddress)
}

func GetEmailCodeV2(c *gin.Context) {
	var ask PostEmailVeryCode
	if err := c.ShouldBind(&ask); err != nil || ask.EmailAddress == "" ||
		!utils.VerifyEmailFormat(ask.EmailAddress) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parame id error"})
		return
	}
	_, requestRedisDbID, _ := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), utils.OperationIDGenerator())
	requestRedisDb := ask.EmailAddress
	if requestRedisDbID != "" {
		requestRedisDb += ":" + requestRedisDbID
	}
	requestRedisDbKey := db.ChangeEmailVerifyCodePrefix + "_" + requestRedisDb
	sendEmailCode(c, requestRedisDbKey, ask.EmailAddress)
}

func LoginEmail(c *gin.Context) {
	var ask PostEmailVeryCode
	if err := c.ShouldBind(&ask); err != nil || ask.EmailAddress == "" || ask.Code != "" || *ask.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "param id error"})
		return
	}
	//请将数据家在到缓存里面。
	emailInfo, err := rocksCache.GetEmailUserInfo(ask.EmailAddress)
	if err == nil && utils.Md5(emailInfo.EmailPassword) == utils.Md5(*ask.Password) {
		c.JSON(http.StatusOK, gin.H{"errCode": 0,
			"errMsg": "ok"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 100009,
		"errMsg": "password error"})
	return
}
func RegisterEmailCode(c *gin.Context) {
	//让其生产一个邮件给发送
	var ask PostEmailVeryCode
	if err := c.ShouldBind(&ask); err != nil || ask.EmailAddress == "" || ask.Code == "" || *ask.Password == "" {

		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "param id error"})
		return
	}
	fmt.Println("ask.EmailAddress ", ask.EmailAddress, "ask.Code", ask.Code, "ask.Password", *ask.Password)
	if !db.DB.ExistVerifyCode(ask.EmailAddress) {
		c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "your email verify code is expire"})
		return
	}
	codestring, err := db.DB.GetEmailVerifyCode(ask.EmailAddress)
	fmt.Println("codestring", codestring, ask.Code)
	if err == nil && ask.Code == codestring {
		if !imdb.CheckIsHaveThisEmail(ask.EmailAddress) {
			imdb.InsertIntoEmail(ask.EmailAddress, *ask.Password)
			rocksCache.DelEmailUserInfo(ask.EmailAddress)
			c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "ok"})
			return
		} else {
			imdb.UpdateIntoEmail(ask.EmailAddress, *ask.Password)
			rocksCache.DelEmailUserInfo(ask.EmailAddress)
			c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "update password ok"})
			return
		}
	}
	db.DB.DeleteEmailVerifyCode(ask.EmailAddress)
	c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "your email verify code is expire"})
	return
}
func RegisterEmailCodeWithOutPassword(c *gin.Context) {
	//让其生产一个邮件给发送
	var ask PostEmailVeryCode
	if err := c.ShouldBind(&ask); err != nil || ask.EmailAddress == "" || ask.Code == "" {

		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "param id error"})
		return
	}
	fmt.Println("ask.EmailAddress ", ask.EmailAddress, "ask.Code", ask.Code)
	if !db.DB.ExistVerifyCode(ask.EmailAddress) {
		c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "your email verify code is expire"})
		return
	}
	codestring, err := db.DB.GetEmailVerifyCode(ask.EmailAddress)
	fmt.Println("codestring", codestring, ask.Code)
	if err == nil && ask.Code == codestring {
		if !imdb.CheckIsHaveThisEmail(ask.EmailAddress) {
			imdb.InsertIntoEmailWithPrivateKey(ask.EmailAddress, *ask.Password, *ask.EncryptPrivateKey)
			rocksCache.DelEmailUserInfo(ask.EmailAddress)
			c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "ok"})
			return
		} else {
			getdata, err := rocksCache.GetEmailUserInfo(ask.EmailAddress)
			if getdata == nil {
				c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "ok", "data": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "ok", "data": getdata.EncryptPrivateKey})
			}
			return
		}
	}
	db.DB.DeleteEmailVerifyCode(ask.EmailAddress)
	c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "your email verify code is expire"})
	return
}
func ChangeEmailPrivateWithOutPassword(c *gin.Context) {
	//让其生产一个邮件给发送
	var ask PostEmailVeryCode
	if err := c.ShouldBind(&ask); err != nil || ask.EmailAddress == "" || ask.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "param id error"})
		return
	}
	if ask.EncryptPrivateKey == nil && ask.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "至少要修改一种东西"})
		return
	}
	if ask.Password != nil && len(*ask.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "密码需要大于6位数"})
		return
	}
	_, requestRedisDbID, _ := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), utils.OperationIDGenerator())
	requestRedisDb := ask.EmailAddress
	if requestRedisDbID != "" {
		requestRedisDb += ":" + requestRedisDbID
	}
	requestRedisDbKey := db.ChangeEmailVerifyCodePrefix + "_" + requestRedisDb
	getUserInfoEmail, err := rocksCache.GetEmailUserInfo(ask.EmailAddress)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "your email not exist"})
		return
	}
	if nValueCode, _ := db.DB.GetEmailVerifyCodeV2(requestRedisDbKey); nValueCode == "" {
		c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "your email verify code is expire"})
		return
	} else {
		tempString := ""
		fmt.Println("codestring", nValueCode, ask.Code)
		db.DB.DeleteEmailVerifyCodeV2(requestRedisDbKey)
		if ask.Code == nValueCode {
			updateField := make(map[string]interface{}, 0)
			if ask.Password != nil {
				updateField["email_password"] = *ask.Password
				tempString += "<密码修改成功>"
			}
			if ask.EncryptPrivateKey != nil {
				updateField["encrypt_private_key"] = *ask.EncryptPrivateKey
				tempString += "<加密密钥修改成功>"
			}
			err = db.DB.MysqlDB.DefaultGormDB().Table("email_user_system").Where("email_address=?", getUserInfoEmail.EmailAddress).
				Updates(updateField).Error
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "db error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": tempString})
		return
	}

}

type SystemConfigInfoResp struct {
	SocialCount       int64  `json:"socialCount"`
	ImByMoneyCount    int64  `json:"imByMoneyCount"`
	NftImageHeadCount int64  `json:"nftImageHeadCount"`
	YaoQingHaoYou     int64  `json:"yaoQingHaoYou"`
	LastUpdateTime    string `json:"lastUpdateTime"`
	SheQuLingQu       int64  `json:"sheQuLingQu"`
}

func SystemConfigInfo(ctx *gin.Context) {
	var ask SystemConfigInfoResp
	if !utils.FileExist("syteminfo.json") {
		ask = SystemConfigInfoResp{
			SocialCount:       500,
			ImByMoneyCount:    450,
			YaoQingHaoYou:     200,
			NftImageHeadCount: 1,
			SheQuLingQu:       10,
			LastUpdateTime:    time.Now().Format("2006-01-02 15:03:04"),
		}
		utils.SaveJson("syteminfo.json", ask)
	} else {
		utils.LoadJson("syteminfo.json", &ask)
	}
	ctx.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": ask})

}
func UpdateSystemInfo() {
	var ask SystemConfigInfoResp
	ask = SystemConfigInfoResp{
		SocialCount:       500,
		ImByMoneyCount:    450,
		YaoQingHaoYou:     200,
		NftImageHeadCount: 1,
		SheQuLingQu:       10,
		LastUpdateTime:    time.Now().Format("2006-01-02 15:03:04"),
	}
	utils.LoadJson("syteminfo.json", &ask)
	er, nowMinse := utils.SubDemo(ask.LastUpdateTime)
	if er == nil && nowMinse.Minutes() > float64(config.Config.UpdateSystemCountMine) {
		newAddRegist := randomdata.Number(2000, 2500)
		ask.SocialCount += int64(newAddRegist)
		ask.ImByMoneyCount += int64(randomdata.Number(1000, 1500))
		ask.YaoQingHaoYou += int64(newAddRegist / 4)
		ask.SheQuLingQu += int64(newAddRegist / 200)
		ask.LastUpdateTime = time.Now().Format("2006-01-02 15:03:04")
		utils.SaveJson("syteminfo.json", ask)
	}
}

func GetTokenUriByTokenID(c *gin.Context) {
	var ask base_info.RequestTokenIdReq
	if err := c.ShouldBind(&ask); err != nil || ask.ContractAddress == "" || ask.TokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parame id error"})
		return
	}
	isOK := false
	var resultdata base_info.RequestTokenIdResp
	rcplist := config.GetRpcFromChainID(ask.ChainID)
	for _, valueUrl := range rcplist {
		client2, err := ethclient.Dial(valueUrl)
		if err != nil {
			resultdata.ErrMsg = err.Error()
			continue
		}
		contractAddress := common.HexToAddress(ask.ContractAddress)
		contractptr, err := erc721.NewContract(contractAddress, client2)
		if err != nil {
			resultdata.ErrMsg = err.Error()
			continue
		}
		bigIntTokenID, errbool := big.NewInt(0).SetString(ask.TokenID, 10)
		if !errbool {
			resultdata.ErrMsg = "error request tokenid"
			continue
		}
		balanceOwnerTokenUrl, err := contractptr.TokenURI(&bind.CallOpts{}, bigIntTokenID)
		if strings.HasPrefix(balanceOwnerTokenUrl, "ipfs://") {
			if config.Config.IsPublicEnv {
				balanceOwnerTokenUrl = strings.ReplaceAll(balanceOwnerTokenUrl,
					"ipfs://", "https://ipfs.io/ipfs/")
			} else {
				balanceOwnerTokenUrl = strings.ReplaceAll(balanceOwnerTokenUrl, "ipfs://",
					"https://ipfs.io/ipfs/")
			}
		}
		resutlbyte, err := utils.HttpGetWithHeader(balanceOwnerTokenUrl, map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		})
		if err != nil {
			resultdata.ErrMsg = "获取nft头像的metadata 错误1"
			continue
		} else {
			var nftmetdata TNFTImageData
			json.Unmarshal(resutlbyte, &nftmetdata)
			if nftmetdata.Image != "" {
				resultdata.TokenUrl = nftmetdata.Image
			}
		}
		isOK = true
		break
	}
	if !isOK {
		resultdata.ErrCode = 50001
	} else {
		resultdata.ErrMsg = ""
	}
	c.JSON(http.StatusOK, resultdata)
}
func TokenOwnerByTokeID(c *gin.Context) {
	var ask base_info.RequestTokenIdReq
	if err := c.ShouldBind(&ask); err != nil || ask.ContractAddress == "" || ask.TokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parame id error"})
		return
	}
	fmt.Println(utils.StructToJsonString(ask))
	isOK := false
	var resultdata base_info.RequestTokenIdResp
	rcplist := config.GetRpcFromChainID(ask.ChainID)
	for _, valueUrl := range rcplist {
		client2, err := ethclient.Dial(valueUrl)
		if err != nil {
			resultdata.ErrMsg = err.Error()
			continue
		}
		contractAddress := common.HexToAddress(ask.ContractAddress)
		contractptr, err := erc721.NewContract(contractAddress, client2)
		if err != nil {
			resultdata.ErrMsg = err.Error()
			continue
		}
		bigIntTokenID, errbool := big.NewInt(0).SetString(ask.TokenID, 10)
		if !errbool {
			resultdata.ErrMsg = "error request tokenid"
			continue
		}
		balanceOwner, err := contractptr.OwnerOf(&bind.CallOpts{}, bigIntTokenID)
		if err != nil {
			resultdata.ErrMsg = err.Error()
			continue
		} else {
			resultdata.TokenOwnerAddress = balanceOwner.String()
		}
		isOK = true
		break
	}
	if !isOK {
		resultdata.ErrCode = 50001
	} else {
		resultdata.ErrMsg = ""
	}
	c.JSON(http.StatusOK, resultdata)
}

func TokenOwnerByTokeIDChainIDContract(c *gin.Context) {
	var ask base_info.RequestImageTokenIdReq
	if err := c.ShouldBind(&ask); err != nil || ask.ContractAddress == "" || ask.TokenID == "" ||
		ask.TokenImageModel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parame id error"})
		return
	}
	isOK := false
	var resultdata base_info.RequestTokenIdResp
	var errmsg = ""
	var balanceOwnerTokenUrl = ""
	rcplist := config.GetRpcFromChainID(ask.ChainID)
	for _, valueUrl := range rcplist {
		fmt.Println("check rpclist:", valueUrl)
		tempValueUrl := valueUrl
		client2, err := ethclient.Dial(tempValueUrl)
		if err != nil {
			resultdata.ErrMsg = err.Error()
			continue
		}
		contractAddress := common.HexToAddress(ask.ContractAddress)
		switch ask.TokenImageModel {
		case "erc721":
			var balanceOwner common.Address
			ask2 := base_info.RequestTokenIdReq{
				ContractAddress: ask.ContractAddress,
				TokenID:         ask.TokenID,
				ChainID:         ask.ChainID,
			}
			if checkIsEnsToken(client2, ask2) {
				balanceOwnerTokenUrl, balanceOwner, err = getEnsTokenUrl(client2, ask2)
				fmt.Println(">>>>>>>>>>>>>>", balanceOwnerTokenUrl, balanceOwner)
			} else {
				balanceOwnerTokenUrl, balanceOwner, err = getErc721TokenUrl(client2, ask2)
			}
			if err != nil {
				resultdata.ErrMsg = err.Error()
				errmsg = err.Error()
				continue
			}
			resultdata.TokenOwnerAddress = balanceOwner.String()
		case "erc1155":
			fmt.Println(utils.StructToJsonString(ask))
			contractptr, err := erc1155.NewContract(contractAddress, client2)
			if err != nil {
				resultdata.ErrMsg = err.Error()
				errmsg = err.Error()
				continue
			}
			bigIntTokenID, errbool := big.NewInt(0).SetString(ask.TokenID, 10)
			if !errbool {
				resultdata.ErrMsg = "error request tokenid"
				errmsg = resultdata.ErrMsg
				continue
			}
			count, err := contractptr.BalanceOf(&bind.CallOpts{}, common.HexToAddress(ask.OwnerAddress), bigIntTokenID)
			if err != nil {
				resultdata.ErrMsg = err.Error()
				errmsg = err.Error()
				continue
			}
			balanceOwnerTokenUrl, err = contractptr.Uri(&bind.CallOpts{}, bigIntTokenID)
			if err != nil {
				resultdata.ErrMsg = err.Error()
				errmsg = err.Error()
				continue
			}
			if count.Int64() == 0 {
				resultdata.TokenOwnerAddress = ""
			} else {
				resultdata.TokenOwnerAddress = ask.OwnerAddress
			}
			fmt.Println(balanceOwnerTokenUrl)
			if strings.HasSuffix(balanceOwnerTokenUrl, "0x{id}") {
				balanceOwnerTokenUrl = strings.ReplaceAll(balanceOwnerTokenUrl, "0x{id}", ask.TokenID)
				_, balanceOwnerTokenUrl = checkIsBiuBiuMetaData(balanceOwnerTokenUrl, ask.TokenID)
			}
			fmt.Println(balanceOwnerTokenUrl)
		default:
			resultdata.ErrMsg = err.Error()
			errmsg = err.Error()
			isOK = false
			break

		}
		if err != nil {
			resultdata.ErrMsg = "获取nft头像的metadata 错误2"
			errmsg = err.Error() + balanceOwnerTokenUrl
			continue
		} else {
			resultdata.TokenUrl = balanceOwnerTokenUrl
		}
		isOK = true
		break
	}
	if !isOK {
		resultdata.ErrCode = 50001
		resultdata.ErrMsg = "chain rpc 的数据有问题" + errmsg
	} else {
		resultdata.ErrCode = 0
		resultdata.ErrMsg = ""
	}
	fmt.Println(utils.StructToJsonString(resultdata))
	c.JSON(http.StatusOK, resultdata)
}

// TokenUrl  需要兼容1155的标签的信息  待优化
func TokenUrl(c *gin.Context) {
	var ask base_info.RequestTokenIdReq
	if err := c.ShouldBind(&ask); err != nil || ask.ContractAddress == "" || ask.TokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parame id error"})
		return
	}
	isOK := false
	var resultdata base_info.RequestTokenIdResp
	var errmsg = ""
	rcplist := config.GetRpcFromChainID(ask.ChainID)
	for _, valueUrl := range rcplist {
		fmt.Println("check rpclist:", valueUrl)
		tempValueUrl := valueUrl
		client2, err := ethclient.Dial(tempValueUrl)
		if err != nil {
			resultdata.ErrMsg = err.Error()
			continue
		}
		balanceOwnerTokenUrl := ""
		var tempAddress common.Address
		if checkIsEnsToken(client2, ask) {
			balanceOwnerTokenUrl, tempAddress, err = getEnsTokenUrl(client2, ask)
			fmt.Println(balanceOwnerTokenUrl, tempAddress)
		} else {
			balanceOwnerTokenUrl, tempAddress, err = getErc721TokenUrl(client2, ask)
		}
		if err != nil {
			errmsg = err.Error()
			fmt.Println(err.Error())
			continue
		}
		resultdata.TokenOwnerAddress = strings.ToLower(tempAddress.String())
		if strings.HasPrefix(balanceOwnerTokenUrl, "ipfs://") {
			if config.Config.IsPublicEnv {
				balanceOwnerTokenUrl = strings.ReplaceAll(balanceOwnerTokenUrl, "ipfs://",
					"https://ipfs.io/ipfs/")
			} else {
				balanceOwnerTokenUrl = strings.ReplaceAll(balanceOwnerTokenUrl, "ipfs://",
					"https://ipfs.io/ipfs/")
			}
		}
		errorMsg, urlAddress := checkIsBiuBiuMetaData(balanceOwnerTokenUrl, ask.TokenID)
		if errorMsg != "" {
			errmsg = errorMsg
			continue
		} else {
			resultdata.TokenUrl = urlAddress
		}
		isOK = true
		break
	}
	if !isOK {
		resultdata.ErrCode = 50001
		resultdata.ErrMsg = "chain rpc 的数据有问题" + errmsg
	} else {
		resultdata.ErrCode = 0
		resultdata.ErrMsg = ""
	}
	c.JSON(http.StatusOK, resultdata)
}

type TNFTImageData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Dna         string `json:"dna"`
	Edition     int    `json:"edition"`
	Date        int64  `json:"date"`
	Attributes  []struct {
		TraitType string `json:"trait_type"`
		Value     string `json:"value"`
	} `json:"attributes"`
}

func CheckIsHaveGuanFangNft(c *gin.Context) {
	var ask EnsAddressReq
	if err := c.ShouldBind(&ask); err != nil || ask.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parame id error"})
		return
	}
	contractAddressList, err := rocksCache.GetOfficialNftFromCache()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parame id error"})
		return
	}
	var resultdata base_info.BalanceOfCount
	for _, valueContractaddress := range contractAddressList {
		isOK := false
		strinList := config.GetRpcFromChainID(utils.Int64ToString(valueContractaddress.SystemNftChainId))
		contractAddress := common.HexToAddress(valueContractaddress.SystemNftContract)
		for _, valueUrl := range strinList {
			client2, err := ethclient.Dial(valueUrl)
			if err != nil {
				resultdata.ErrMsg = err.Error()
				continue
			}
			holderAddress := common.HexToAddress(ask.Address)
			contractptr, err := erc721.NewContract(contractAddress, client2)
			fmt.Println("查询持有nft 的内容:", valueUrl, ask.Address, contractAddress)
			if err != nil {
				resultdata.ErrMsg = err.Error()
				continue
			}
			balance, err := contractptr.BalanceOf(&bind.CallOpts{}, holderAddress)
			if err != nil {
				resultdata.ErrMsg = err.Error()
				continue
			}
			resultdata.BalanceOf = balance.Int64()
			fmt.Println("resultdata.BalanceOf >>>>", resultdata.BalanceOf)
			if resultdata.BalanceOf == 0 {
				break
			}
			isOK = true
			break
		}
		if !isOK {
			resultdata.ErrCode = 50001
		} else {
			resultdata.ErrMsg = ""
			break
		}
	}
	c.JSON(http.StatusOK, resultdata)
}

var WriteFileCache sync.Mutex

func GetCoin(c *gin.Context) {
	reqyest := c.Request.URL.RawQuery
	md5filekey := utils.Md5(reqyest)
	resultUpdateTime, _ := db.DB.GetFileUpdateTime(md5filekey)
	_, resulttmediff := utils.SubDemo(resultUpdateTime)
	if resultUpdateTime == "" || resulttmediff > 5*time.Minute {
		WriteFileCache.Lock()
		//缓存保存内容
		resultUpdateTime2, _ := db.DB.GetFileUpdateTime(md5filekey)
		if resultUpdateTime == resultUpdateTime2 {
			var resultdatabyte []byte
			var err error
			if config.Config.OpenNetProxy.OpenFlag {
				proxy := "http://proxy.idchats.com:7890"
				proxyAddress, _ := url.Parse(proxy)
				resultdatabyte, err = utils.HttpGetWithProxy("https://api.coingecko.com/api/v3/coins/markets?"+reqyest,
					http.ProxyURL(proxyAddress))
			} else {
				resultdatabyte, err = utils.HttpGet("https://api.coingecko.com/api/v3/coins/markets?" + reqyest)
			}
			utils.SaveString("currentcypt_"+md5filekey+".txt", &resultdatabyte)
			if err == nil {
				db.DB.SetFileUpdateTime(md5filekey)
			}
		}
		WriteFileCache.Unlock()
	}
	resultdatabytenew, err := utils.ReadString("currentcypt_" + md5filekey + ".txt")
	if err == nil {
		c.String(http.StatusOK, resultdatabytenew)

	} else {
		c.String(http.StatusBadRequest, "")
	}
}
func GetAppVersion(c *gin.Context) {
	platform := c.Param("platform")
	if platform == "" {
		c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": "not with platform"})
		return
	}
	appVersion, err := imdb.GetAppVersion(platform)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "version errror"})
		return
	}
	var resp = new(base_info.AppVersionDataResp)
	log.Printf("%+v", appVersion)
	if appVersion.HasUpdate != 0 {
		resp.HasUpdate = true
	}
	if appVersion.IsIgnorable != 0 {
		resp.IsIgnorable = true
	}
	resp.VersionCode = int32(appVersion.VersionCode)
	resp.VersionName = appVersion.VersionName
	resp.UpdateLog = appVersion.UpdateLog
	resp.UpdateLogEn = appVersion.UpdateLogEn
	resp.ApkUrl = appVersion.ApkURL
	c.JSON(http.StatusOK, &base_info.AppVersionResp{
		ErrCode:            0,
		ErrMsg:             "",
		AppVersionDataResp: resp,
	})
}

// GetTokenList 获取某条链条上的token 的信息
func GetTokenList(c *gin.Context) {
	platform := c.Param("chainid")
	if platform != "1" && platform != "56" && platform != "137" {
		c.JSON(http.StatusOK, gin.H{"errCode": 0,
			"errMsg": "not have this chain token", "data": nil})
		return
	}
	resultdata, err := imdb.GetChainTokenByChainID(platform)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"errCode": 0,
			"errMsg": "",
			"data":   resultdata})

	} else {
		c.JSON(http.StatusOK, gin.H{"errCode": 400,
			"errMsg": "not support this chain"})
	}

}
func PostNewToken(c *gin.Context) {
	platform := c.Param("chainid")
	if platform != "1" && platform != "56" && platform != "137" {
		c.JSON(http.StatusOK, gin.H{"errCode": 400,
			"errMsg": "not support this chain1"})
		return
	}
	tokenaddress := c.Param("tokenaddress")
	if tokenaddress == "" {
		c.JSON(http.StatusOK, gin.H{"errCode": 400,
			"errMsg": "not support this chain2"})
		return
	}
	if !imdb.IsInsertCoinToken(platform, tokenaddress) {
		c.JSON(http.StatusOK, gin.H{"errCode": 400,
			"errMsg": "has insert this token"})
		return
	}
	var tmpdata GenerateObjTokenCoin
	var resultByte []byte
	var errhttp error
	switch platform {
	case "1": //eth mainnet
		proxy := "http://proxy.idchats.com:7890"
		proxyAddress, _ := url.Parse(proxy)
		resultByte, errhttp = utils.HttpGetWithProxy(
			"https://pro-api.coingecko.com/api/v3/coins/binance-smart-chain/contract/"+tokenaddress+"?x_cg_pro_api_key=CG-ZRUJUuYGELw13WB13DptYbFc",
			http.ProxyURL(proxyAddress))
	case "56": //bsc mainnet
		proxy := "http://proxy.idchats.com:7890"
		proxyAddress, _ := url.Parse(proxy)
		resultByte, errhttp = utils.HttpGetWithProxy(
			"https://pro-api.coingecko.com/api/v3/coins/binance-smart-chain/contract/"+tokenaddress+"?x_cg_pro_api_key=CG-ZRUJUuYGELw13WB13DptYbFc",
			http.ProxyURL(proxyAddress))
	case "137":
		proxy := "http://proxy.idchats.com:7890"
		proxyAddress, _ := url.Parse(proxy)
		resultByte, errhttp = utils.HttpGetWithProxy(
			"https://pro-api.coingecko.com/api/v3/coins/polygon-pos/contract/"+tokenaddress+"?x_cg_pro_api_key=CG-ZRUJUuYGELw13WB13DptYbFc",
			http.ProxyURL(proxyAddress))
	}
	if errhttp != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": 400,
			"errMsg": "not support this chain3"})
		return
	}
	json.Unmarshal(resultByte, &tmpdata)
	if tmpdata.Error != "" {
		c.JSON(http.StatusOK, gin.H{"errCode": 400,
			"errMsg": tmpdata.Error})
		return
	}

	iconUrl := ""
	var insertIntoDb []*db.ChainToken

	if nData, ok := tmpdata.DetailPlatforms["ethereum"]; ok {
		//iconUrl = "https://raw.githubusercontent.com/yearn/yearn-assets/master/icons/multichain-tokens/1/" + nData.ContractAddress + `/logo-128.png`
		//iconUrl= "https://raw.githubusercontent.com/dappradar/tokens/master/binance-smart-chain/0xba2ae424d960c26247dd6c32edc70b295c744c43/logo.png
		iconUrl = "https://raw.githubusercontent.com/dappradar/tokens/master/ethereum/" + nData.ContractAddress + "/logo.png"
		insertIntoDb = append(insertIntoDb, &db.ChainToken{
			CoinChainid:  1,
			CoinToken:    nData.ContractAddress,
			CoinDecimals: nData.DecimalPlace,
			CoinName:     tmpdata.Name,
			CoinSymbol:   tmpdata.Symbol,
			CoinType:     "ERC20",
			CoinIsHot:    0,
			CoinIcon:     iconUrl,
		})
	}
	if nData, ok := tmpdata.DetailPlatforms["binance-smart-chain"]; ok {
		iconUrl = "https://raw.githubusercontent.com/dappradar/tokens/master/binance-smart-chain/" + nData.ContractAddress + "/logo.png"
		insertIntoDb = append(insertIntoDb, &db.ChainToken{
			CoinChainid:  56,
			CoinToken:    nData.ContractAddress,
			CoinDecimals: nData.DecimalPlace,
			CoinName:     tmpdata.Name,
			CoinSymbol:   tmpdata.Symbol,
			CoinType:     "BEP20",
			CoinIsHot:    0,
			CoinIcon:     iconUrl,
		})
	}
	if nData, ok := tmpdata.DetailPlatforms["polygon-pos"]; ok {
		iconUrl = "https://raw.githubusercontent.com/dappradar/tokens/master/polygon/" + nData.ContractAddress + "/logo.png"
		insertIntoDb = append(insertIntoDb, &db.ChainToken{
			CoinChainid:  137,
			CoinToken:    nData.ContractAddress,
			CoinDecimals: nData.DecimalPlace,
			CoinName:     tmpdata.Name,
			CoinSymbol:   tmpdata.Symbol,
			CoinType:     "ERC20",
			CoinIsHot:    0,
			CoinIcon:     iconUrl,
		})
	}
	if len(insertIntoDb) > 0 {
		errhttp = imdb.PostChainNewToken(&insertIntoDb)
	}
	if errhttp != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": 400,
			"errMsg": "not support this chain5"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"errCode": 0,
			"errMsg": ""})
		return
	}
}

type DetailPlatforms struct {
	DecimalPlace    int    `json:"decimal_place,omitempty"`
	ContractAddress string `json:"contract_address,omitempty"`
}
type GenerateObjTokenCoin struct {
	Error           string                     `json:"error"`
	Id              string                     `json:"id,omitempty"`
	Symbol          string                     `json:"symbol,omitempty"`
	Name            string                     `json:"name,omitempty"`
	AssetPlatformId string                     `json:"asset_platform_id,omitempty"`
	Platforms       map[string]string          `json:"platforms,omitempty"`
	DetailPlatforms map[string]DetailPlatforms `json:"detail_platforms,omitempty"`
}

func checkIsEnsToken(ethClient *ethclient.Client, ask base_info.RequestTokenIdReq) bool {
	fmt.Println("checkIsEnsToken", time.Now().UnixMilli())
	contractptr, _ := erc721.NewContract(common.HexToAddress(ask.ContractAddress), ethClient)
	dt, _ := big.NewInt(0).SetString(ask.TokenID, 10)
	_, err := contractptr.TokenURI(&bind.CallOpts{}, dt)
	if err != nil {
		fmt.Println("is not erc 721 token", time.Now().UnixMilli())
		return true
	}
	fmt.Println("is not erc 721 token", time.Now().UnixMilli())
	return false
}
func getErc721TokenUrl(ethClient *ethclient.Client, ask base_info.RequestTokenIdReq) (nftInfoUrl string, balanceOwner common.Address, err error) {
	contractAddress := common.HexToAddress(ask.ContractAddress)
	contractptr, err := erc721.NewContract(contractAddress, ethClient)
	if err != nil {
		return "", common.BytesToAddress([]byte("")), err
	}
	bigIntTokenID, errbool := big.NewInt(0).SetString(ask.TokenID, 10)
	if !errbool {
		return "", common.BytesToAddress([]byte("")), errors.New("error request tokenid")
	}
	balanceOwner, err = contractptr.OwnerOf(&bind.CallOpts{}, bigIntTokenID)
	if err != nil {
		return "", common.BytesToAddress([]byte("")), err
	}
	nftInfoUrl, err = contractptr.TokenURI(&bind.CallOpts{}, bigIntTokenID)
	if err != nil {
		return "", common.BytesToAddress([]byte("")), err
	}
	return
}
func getEnsTokenUrl(ethClient *ethclient.Client, ask base_info.RequestTokenIdReq) (nftInfoUrl string, balanceOwner common.Address, err error) {
	fmt.Println("getEsnTokenUrl begin", time.Now().UnixMilli())
	contractptr, _ := erc721.NewContract(common.HexToAddress(ask.ContractAddress),
		ethClient)
	dt, _ := big.NewInt(0).SetString(ask.TokenID, 10)
	balanceOwner, err = contractptr.OwnerOf(&bind.CallOpts{}, dt)
	chainName, chainImageUrl := config.GetEnsTokenUrlServiceByChainID(ask.ChainID)
	fmt.Println("getEsnTokenUrl end", time.Now().UnixMilli())
	if chainName == "" {
		return fmt.Sprintf("%s/%s", chainImageUrl, ask.TokenID), balanceOwner, nil
	} else {
		return fmt.Sprintf("%s/%s/%s/%s", chainImageUrl, chainName, ask.ContractAddress, ask.TokenID), balanceOwner, nil
	}
}
func checkIsBiuBiuMetaData(baseUrlPost string, tokenID string) (errormsg string, tokenUrl string) {
	balanceOwnerTokenUrl := baseUrlPost
	if strings.Contains(baseUrlPost, "0x{id}") {
		balanceOwnerTokenUrl = strings.ReplaceAll(baseUrlPost, "0x{id}", tokenID)
	}
	proxyAddress, _ := url.Parse("http://proxy.idchats.com:7890")
	var resutlbyte []byte
	var err error
	if config.Config.OpenNetProxy.OpenFlag {
		resutlbyte, err = utils.HttpGetWithHeaderWithProxy(balanceOwnerTokenUrl, map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		}, http.ProxyURL(proxyAddress))
	} else {
		resutlbyte, err = utils.HttpGetWithHeader(balanceOwnerTokenUrl, map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		})
	}

	if err != nil {
		return err.Error() + balanceOwnerTokenUrl, ""
	} else {
		var nftmetdata TNFTImageData
		_ = json.Unmarshal(resutlbyte, &nftmetdata)
		if nftmetdata.Image != "" {
			return "", nftmetdata.Image
		} else {
			return "合约上无设置url地址", ""
		}
	}
}
func sendEmailCode(c *gin.Context, requestRedisDbKey string, emailAddress string) {
	//检查是否带有token
	if nValue, _ := db.DB.GetEmailVerifyCodeV2(requestRedisDbKey); nValue != "" {
		c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "email verify code has send"})
		return
	}
	code := randomdata.Number(100000, 999999)
	levelTime := config.Config.Demo.CodeTTL * 2
	err := db.DB.SetEmailVerifyCodeV2(requestRedisDbKey, code, levelTime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr,
			"errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
		return
	}
	bodyString := "Your VerifyCode is:" + utils.IntToString(code)
	if strings.HasPrefix(requestRedisDbKey, db.ChangeEmailVerifyCodePrefix) {
		bodyString = "Your will change your wallet encrypt private key,\n now verifyCode is:" + utils.IntToString(code)
	}
	//发送kafka 时间到内容 然后让其发送email
	msgkafka := &kafkaMessage.KafkaMsg{
		MessageType: 1,
		EmailMsg: &kafkaMessage.EmailMessage{
			EmailType: 1,
			ToAddress: emailAddress,
			Subject:   "VerifyCode",
			Title:     "VerifyCode",
			Body:      bodyString,
		},
		SmsMsg: nil,
	}
	msgvalue, _ := proto.Marshal(msgkafka)
	msg := xkafka.ProducerMessage{
		Topic: config.Config.Kafka.BusinessTop.Topic,
		Value: msgvalue,
	}
	if ct != nil && ct.Producer != nil {
		ct.Producer.SendMessage(msg, func(metadata *xkafka.RecordMetadata, err error) {
			if err != nil {
				xlog.CInfo(err.Error())
				c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr,
					"errMsg": err.Error()})
				return
			}
		})
		c.JSON(http.StatusOK, gin.H{"errCode": 0,
			"errMsg": "yanzhengmaok"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"errCode": 100009,
			"errMsg": "email service not open"})
		return
	}
}
