package base_info

import (
	"Open_IM/internal/api/brc20/services"
)

// 定义brc20 的结构体
type AllBrc20TokensReq struct {
	OperationID string `json:"operationID,required"`
	PageIndex   int    `json:"offset,required"`
	PageSize    int    `json:"limit,required"`
	Type        string `json:"type,required"`
	Sort        string `json:"sort"`               // deployTime,desc  mintRate,desc,...
	Timestamp   int    `json:"timestamp,required"` //请求的时间戳第一页的时间 会关系到分页面的行书数目
}

type AllBrc20TokensResp struct {
	CommResp
	Data *HitsBrc20Response `json:"data"`
}
type HitsBrc20Response struct {
	Total int             `json:"total"`
	Hits  []*ApiHitsBrc20 `json:"hits"`
}
type HitsBrc20 struct {
	Rank              int         `json:"rank"`
	LogoUrl           string      `json:"logoUrl"`
	Name              string      `json:"name"`
	DisplayName       string      `json:"displayName"`
	Supply            interface{} `json:"supply"`
	Minted            interface{} `json:"minted"`
	MintRate          float64     `json:"mintRate"`
	TransactionCount  int         `json:"transactionCount"`
	HolderCount       int         `json:"holderCount"`
	DeployTime        int         `json:"deployTime"`
	InscriptionNumber string      `json:"inscriptionNumber"`
	InscriptionId     string      `json:"inscriptionId"`
	Price             float64     `json:"price,omitempty"`
	TickId            string      `json:"tickId"`
	TokenType         string      `json:"tokenType"`
}
type ApiHitsBrc20 struct {
	Rank              int     `json:"rank"`
	LogoUrl           string  `json:"logoUrl"`
	Name              string  `json:"name"`
	DisplayName       string  `json:"displayName"`
	Supply            string  `json:"supply"`
	Minted            string  `json:"minted"`
	MintRate          float64 `json:"mintRate"`
	TransactionCount  int     `json:"transactionCount"`
	HolderCount       int     `json:"holderCount"`
	DeployTime        int     `json:"deployTime"`
	InscriptionNumber string  `json:"inscriptionNumber"`
	InscriptionId     string  `json:"inscriptionId"`
	Price             float64 `json:"price,omitempty"`
	TickId            string  `json:"tickId"`
	TokenType         string  `json:"tokenType"`
}
type PersonalAllBrc20TokenReq struct {
	OperationID string `json:"operationID,required"`
	PageIndex   int    `json:"offset,required"`
	PageSize    int    `json:"limit,required"`
	Type        string `json:"type,required"`
	Timestamp   int    `json:"timestamp,required"` //请求的时间戳第一页的时间 会关系到分页面的行书数目
	Address     string `json:"address,required"`
	RequestType string `json:"requestType,required"` //查询自己的余额balance，如果是自己的deploy ,inscription
}
type PersonalAllBrc20TokenRsp struct {
	CommResp
	Data *HitsPersonalBrc20Response `json:"data"`
}
type HitsPersonalBrc20Response struct {
	Extend struct {
		TokenHeldKind int `json:"tokenHeldKind"`
	} `json:"extend"`
	Hits  []*HitsPersonalBrc20 `json:"hits"`
	Total int64                `json:"total"`
}

type HitsPersonalBrc20 struct {
	Name                   string  `json:"name"`
	DisplayName            string  `json:"displayName"`
	LogoUrl                string  `json:"logoUrl"`
	Balance                float64 `json:"balance"`
	AvailableBalance       float64 `json:"availableBalance"`
	TransferableBalance    float64 `json:"transferableBalance"`
	InscriptionId          string  `json:"inscriptionId"`
	InscriptionNumber      string  `json:"inscriptionNumber"`
	TokenInscriptionNumber string  `json:"tokenInscriptionNumber"`
	InscriptionAmount      int     `json:"inscriptionAmount"`
	Price                  float64 `json:"price"`
	UsdValue               float64 `json:"usdValue"`
	Type                   string  `json:"type"`
	DeployTime             int64   `json:"deployTime"`
	Action                 string  `json:"action"`
	Owner                  string  `json:"owner"`
	Status                 string  `json:"status"`
	ErrorMsg               string  `json:"errorMsg"`
}
type Brc20PledgeReq struct {
	Address     string `json:"address,required"`
	Amount      string `json:"amount"`    //暂时不要考虑小数
	NetParams   string `json:"netParams"` //testnet或者mainnet
	Ticker      string `json:"ticker"`
	PageSize    string `json:"pageSize"`
	PageIndex   string `json:"pageIndex"`
	OperationID string `json:"operationID,required"`
}
type Brc20PledgeUtxo struct {
	Address       string   `json:"address,required"`
	InscriptionId []string `json:"inscriptionId,required"`
	Ticker        string   `json:"ticker,required"`
	NetParams     string   `json:"netParams,required"`
	StakingPeriod string   `json:"stakingPeriod,required"` //质押的年收益选项：
	OperationID   string   `json:"operationID,required"`
	Sign          string   `json:"sign,required"`
}
type Brc20PledgeResp struct {
	CommResp
	SignPsbtData string `json:"signData"` //签名的字段
}
type Brc20TransferableListResp struct {
	CommResp
	Data []*services.TransferAbleInscript `json:"data"`
}

type Brc20ObbtPledgePoolInfoes struct {
	CommResp
	Data []*Brc20PoolInfo `json:"data"`
}
type Brc20PoolInfo struct {
	PoolID        int    `json:"poolID"`
	StakingPeriod string `json:"stakingPeriod"` //180 .360 ,1080
	RewardsRate   string `json:"rewardsRate"`   //1000,2000,3000
	TotalStake    string `json:"totalStake"`
}
type PersonalBrc20Pledge struct {
	StartTime     int64  `json:"startTime"`
	StakedAmount  string `json:"stakedAmount"`
	Reward        string `json:"reward"`
	StakingPeriod string `json:"stakingPeriod"`
}
type PersonalBrc20PledgeResonse struct {
	CommResp
	Data []*PersonalBrc20Pledge `json:"data"`
}
type Brc20InfoForSearchResponse struct {
	CommResp
	Data []*Brc20InfoForSearch `json:"data"`
}
type Brc20InfoForSearch struct {
	Ticker        string `json:"ticker"`
	InscriptionID string `json:"inscriptionID"`
	Creator       string `json:"creator"`
	UsdFloorPrice string `json:"usdFloorPrice"`
	HolderCount   string `json:"holderCount"`
}
