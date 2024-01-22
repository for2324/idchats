package base_info

import "time"

type AppointmentReq struct {
	OperationID string `json:"operationID" binding:"required"`
	Ens         string `json:"ens" binding:"required"`
	Cancel      bool   `json:"cancel" default:"false"`
}

type AppointmentResp struct {
	CommResp
}
type HasAppointmentReq struct {
	OperationID string `json:"operationID"`
	Ens         string `json:"ens" binding:"required"`
}
type HasAppointmentResp struct {
	CommResp
	Data bool `json:"data"`
}

type AppointmentListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	SearchType  string `json:"searchType" enums:"all,mine" default:"all"`
	PageIndex   int    `json:"pageIndex" default:"0"`
	PageSize    int    `json:"pageSize" default:"20"`
}

type AppointmentUserInfo struct {
	UserId     string    `json:"userId"`
	EnsName    string    `json:"ensName"`
	FaceUrl    string    `json:"faceUrl"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;default:null" json:"-"`
	CreateTime int64     `json:"createTime"`
}

type AppointmentListResp struct {
	CommResp
	Data []AppointmentUserInfo `json:"data"`
}

type CreateRegisterEnsOrderReq struct {
	OperationID string `json:"operationID" binding:"required"`
	EnsName     string `json:"ensName" binding:"required"`
	TxnType     string `json:"txnType" binding:"required"`
	EnsInviter  string `json:"ensInviter"`
}

type EnsOrderInfo struct {
	Order   OrderInfo `json:"order"`
	PayInfo PayInfo   `json:"payInfo"`
}
type OrderInfo struct {
	OrderId         uint32 `json:"orderId"`
	EnsName         string `json:"ensName"`
	Status          string `json:"status"`
	TxnType         string `json:"txnType"`
	USDPrice        uint64 `json:"usdPrice"`
	USDGasFee       uint64 `json:"usdGasFee"`
	EnsInviter      string `json:"ensInviter"`
	TxnHash         string `json:"txnHash"`
	RegisterTxnHash string `json:"registerTxnHash"`
	CreateTime      string `json:"createTime"`
	PayTime         string `json:"payTime"`
	ExpireTime      string `json:"expireTime"`
}

type PayInfo struct {
	Id               uint64 `json:"id"`
	USDPrice         uint64 `json:"usdPrice"`
	FormAddress      string `json:"formAddress"`
	ToAddress        string `json:"toAddress"`
	TxnType          string `json:"txnType"`
	Value            string `json:"value"`
	Decimal          uint32 `json:"decimal"`
	StartBlockNumber uint64 `json:"startBlockNumber"`
	ScanBlockNumber  uint64 `json:"scanBlockNumber"`
	Rate             uint64 `json:"rate"`
	GasFee           uint64 `json:"gasFee"`
	Type             string `json:"type"`
	Status           int32  `json:"status"`
	Tag              string `json:"tag"`
	ChainId          int64  `json:"chainId"`
	TxnHash          string `json:"txnHash"`
	CreateTime       string `json:"createTime"`
	BlockStartTime   string `json:"blockStartTime"`
	BlockExpireTime  string `json:"blockExpireTime"`
	OrderId          string `json:"orderId"`
	Mark             string `json:"mark"`
	Ex               string `json:"ex"`
}

type CreateRegisterEnsOrderResp struct {
	CommResp
	Data EnsOrderInfo `json:"data"`
}

type GetEnsOrderInfoReq struct {
	OperationID string `json:"operationID" binding:"required"`
	OrderId     uint32 `json:"orderId" binding:"required"`
}

type GetEnsOrderInfoResp struct {
	CommResp
	Data EnsOrderInfo `json:"data"`
}

type GetSupportCoinListReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type GetSupportCoinListResp struct {
	CommResp
	Data []string `json:"data"`
}

type RequestPaymentReq struct {
	FormAddress       string `json:"formAddress" binding:"required"`                               // 付款人
	USDPrice          uint64 `json:"usdPrice" binding:"required"`                                  // 付款金额 美元
	TxnType           string `json:"txnType" binding:"required"`                                   // 交易类型
	OrderId           string `json:"orderId" binding:"required"`                                   // 订单ID
	NotifyUrl         string `json:"notifyUrl" binding:"required"`                                 // 回调地址
	Attach            string `json:"attach"`                                                       // 附加信息
	NotifyEncryptType string `json:"notifyEncryptType" enums:"AES-GCM" default:"AES-GCM"`          // 回调加密方式  枚举：AES-GCM
	NotifyEncryptKey  string `json:"notifyEncryptKey" binding:"required" validate:"len=16|len=32"` // 回调加密key 限制 16或32 位
}

type RequestPaymentResp struct {
	CommResp
	Data PayInfo `json:"data"`
}

type GetOrderInfoResp struct {
	CommResp
	Data PayInfo `json:"data"`
}

type TestNotifyReq struct {
	Id          string `json:"id"`
	EventType   string `json:"eventType"`   // SUCCESS
	EncryptType string `json:"encryptType"` // AES RSA
	Source      string `json:"source"`
	Nonce       string `json:"nonce"`
}

type TestNotifyReqResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
