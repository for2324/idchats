package order

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbOrder "Open_IM/pkg/proto/order"
	"Open_IM/pkg/utils"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary		获取支持的币种列表
// @Description	获取支持的币种列表
// @Tags			订单相关
// @ID				GetSupportCoinList
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.GetSupportCoinListReq	true	"获取ens订单信息"
// @Produce		json
// @Success		0	{object}	api.GetSupportCoinListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/order/get_support_coin_list [post]
func GetSupportCoinList(c *gin.Context) {
	params := api.GetSupportCoinListReq{}
	apiResp := api.GetSupportCoinListResp{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, utils.GetSelfFuncName(), " args ", params)
	confMap := config.Config.Pay.TnxTypeConfMap
	supportList := []string{}
	for k, _ := range confMap {
		supportList = append(supportList, k)
	}
	apiResp.Data = supportList
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "success", "data": apiResp.Data})
}

// @Summary		创建支付订单
// @Description	创建支付订单
// @Tags			订单相关
// @ID				RequestPayment
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.RequestPaymentReq	true	"获取ens订单信息"
// @Produce		json
// @Success		0	{object}	api.RequestPaymentResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/order/request_payment [post]
func RequestPayment(c *gin.Context) {
	params := api.RequestPaymentReq{}
	apiResp := api.RequestPaymentResp{}
	OperationID := utils.OperationIDGenerator()
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(OperationID, utils.GetSelfFuncName(), " args ", params)
	// 创建监听任务
	etcdConn := getcdv3.GetDefaultConn(
		config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImOrder,
		OperationID,
	)
	client := pbOrder.NewOrderServiceClient(etcdConn)
	resp, err := client.CreatePayScanBlockTask(c.Request.Context(), &pbOrder.CreatePayScanBlockTaskReq{
		OperationID:       OperationID,
		USD:               params.USDPrice,
		FormAddress:       params.FormAddress,
		OrderId:           params.OrderId,
		TxnType:           params.TxnType,
		NotifyUrl:         params.NotifyUrl,
		NotifyEncryptType: params.NotifyEncryptType,
		NotifyEncryptKey:  params.NotifyEncryptKey,
		Attach:            params.Attach,
		Mark:              "tenant",
	})
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "CreatePayScanBlockTask failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if resp.CommonResp.ErrCode != 0 {
		log.NewError(OperationID, utils.GetSelfFuncName(), "CreatePayScanBlockTask failed", resp.CommonResp.ErrMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": resp.CommonResp.ErrMsg})
		return
	}
	utils.CopyStructFields(&apiResp.Data, resp.ScanTaskInfo)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "success", "data": apiResp.Data})
}

type OrderSource struct {
	Mark       string `json:"mark"`
	CreateTime string `json:"createTime"`
	PayTime    string `json:"payTime"`
	TxnHash    string `json:"txnHash"`
	OrderId    string `json:"orderId"`
	TxnType    string `json:"txnType"`
	Attach     string `json:"attach"`
}

func TestNotify(c *gin.Context) {
	params := api.TestNotifyReq{}
	apiResp := api.TestNotifyReqResp{}
	OperationID := utils.OperationIDGenerator()
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(OperationID, utils.GetSelfFuncName(), " args ", params)

	// base 64 解码
	decodeBytes, err := base64.StdEncoding.DecodeString(params.Source)
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "base64 decode failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	plaintext, err := AESGCMDecrypt(decodeBytes, params.Nonce, "qqqqwwwweeeerrrr")
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "AESGCMDecrypt failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(OperationID, utils.GetSelfFuncName(), "plaintext ", plaintext)
	orderSource := OrderSource{}
	err = json.Unmarshal([]byte(plaintext), &orderSource)
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "json.Unmarshal failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(OperationID, utils.GetSelfFuncName(), "orderSource ", orderSource)
	c.JSON(http.StatusOK, apiResp)
}

func AESGCMDecrypt(ciphertext []byte, nonceStr string, key string) (plaintext string, err error) {
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		return
	}
	var openData []byte
	openData, err = aesgcm.Open(nil, []byte(nonceStr), ciphertext, nil)
	if err != nil {
		return
	}
	plaintext = string(openData)
	return
}

// @Summary		通过订单id获取支付订单
// @Description	通过订单id获取支付订单
// @Tags			订单相关
// @ID				GetOrderInfo
// @Accept			json
// @Produce		json
// @Success		0	{object}	api.GetOrderInfoResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/transactions/order/{id} [get]
func GetOrderInfo(c *gin.Context) {
	orderId := c.Params.ByName("id")
	apiResp := api.GetOrderInfoResp{}
	if orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "orderId is empty"})
		return
	}
	OperationID := utils.OperationIDGenerator()
	etcdConn := getcdv3.GetDefaultConn(
		config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImOrder,
		OperationID,
	)
	client := pbOrder.NewOrderServiceClient(etcdConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id, err := strconv.ParseUint(orderId, 10, 64)
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "strconv.ParseUint failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	resp, err := client.GetPayScanBlockTaskById(ctx, &pbOrder.GetPayScanBlockTaskByIdReq{
		Id: id,
	})
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "client.GetPayScanBlockTaskById failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	utils.CopyStructFields(&apiResp, resp.ScanTaskInfo)
	log.NewInfo(OperationID, utils.GetSelfFuncName(), "client.GetPayScanBlockTaskById resp ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		通过订单id获取支付订单
// @Description	通过订单id获取支付订单
// @Tags			订单相关
// @ID				GetOutTradeOrderInfo
// @Accept			json
// @Produce		json
// @Success		0	{object}	api.GetOrderInfoResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/transactions/out_trade_no/{orderId} [get]
func GetOutTradeOrderInfo(c *gin.Context) {
	orderId := c.Params.ByName("orderId")
	apiResp := api.GetOrderInfoResp{}
	if orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "orderId is empty"})
		return
	}
	OperationID := utils.OperationIDGenerator()
	etcdConn := getcdv3.GetDefaultConn(
		config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImOrder,
		OperationID,
	)
	client := pbOrder.NewOrderServiceClient(etcdConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := client.GetPayScanBlockTaskByOrderId(ctx, &pbOrder.GetPayScanBlockTaskByOrderIdReq{
		OperationID: OperationID,
		OrderId:     orderId,
	})
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "client.GetPayScanBlockTaskById failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	utils.CopyStructFields(&apiResp, resp.ScanTaskInfo)
	log.NewInfo(OperationID, utils.GetSelfFuncName(), "client.GetPayScanBlockTaskById resp ", resp)
	c.JSON(http.StatusOK, resp)
}
