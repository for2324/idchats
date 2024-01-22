package base_info

type BusinessApiKeyInfo struct {
	TradeFee    string `json:"tradeFee"`
	SniperFee   string `json:"sniperFee"`
	Key         string `json:"key"`
	TradeVolume string `json:"tradeVolume"`
	CreatedAt   string `json:"createdAt"`
	ApiName     string `json:"apiName"`
}
type UserGetBusinessListResp struct {
	CommResp
	Data UserGetBusinessListResult `json:"Data"`
}
type UserGetBusinessListReq struct {
	PageIndex   int    `json:"pageIndex"`
	PageSize    int    `json:"pageSize" binding:"lte=50" `
	MerchantId  string `json:"merchantId" validate:"required"`
	OperationID string `json:"operationID"`
}
type UserGetBusinessListResult struct {
	Total      int64                 `json:"total"`
	ApiKeyInfo []*BusinessApiKeyInfo `json:"apiKeyInfo"`
}
type UserPersonalTotalVolumeListReq struct {
	PageIndex   int    `json:"pageIndex"`
	PageSize    int    `json:"pageSize" binding:"lte=50"`
	RequestTime int64  `json:"requestTime" validate:"required"` //第一页的请求时间
	OperationID string `json:"operationID"`
}

// 以后这种tmd的 全部改成用any 的设计 去做一个模版化的东西
type UserPersonalTotalVolumeListResp struct {
	CommResp
	Data *UserPersonalTotalVolumeListRespList `json:"data"`
}
type UserPersonalTotalVolumeListRespList struct {
	Total int64                         `json:"total"`
	Data  []*UserPersonalVolumeListData `json:"data"`
}

type UserPersonalVolumeListData struct {
	ID             int64  `json:"ID"`
	FinishTime     int64  `json:"finishDate"`
	Address        string `json:"address"`
	UsdTradeVolume string `json:"volume"`
}

// 前短提交给后端
type UserAddBusinessApiReq struct {
	Key         string `json:"key" validate:"required"`
	ApiName     string `json:"apiName" validate:"required"`
	Sign        string `json:"sign" validate:"required"`
	OperationID string `json:"operationID"`
}
type UserAddBusinessApiResp struct {
	CommResp
}

type UserUpdateBusinessApiReq struct {
	ApiName     string `json:"apiName" validate:"required"`
	Key         string `json:"key" validate:"required"`
	Method      string `json:"method" validate:"required"`
	Sign        string `json:"sign" validate:"required"`
	OperationID string `json:"operationID"`
}
type UserUpdateBusinessApiResp struct {
	CommResp
}
type GetUserTradeRewardScoreReq struct {
	OperationID string `json:"OperatorID"`
	RewardType  string `json:"rewardType" validate:"required"`
}
type GetUserTradeRewardScoreData struct {
	InviteCount         int64  `json:"inviteCount"` //邀请人
	Claim               string `json:"claim"`       //自己交易奖励领取
	Pending             string `json:"pending"`     //自己交易奖励未领取
	RewardFee           string `json:"rewardFee"`
	RakebackPending     string `json:"rakebackPending"` //抽成奖励
	RakebackClaim       string `json:"rakebackClaim"`   //抽成已领取
	PersonalTradeVolume string `json:"personalTradeVolume"`
	SubPersonalTradeFee string `json:"subPersonalTradeFee"`
}

type GetUserTradeRewardScoreResp struct {
	CommResp
	Data *GetUserTradeRewardScoreData `json:"data"`
}
type UserStakeRewardReq struct {
	//领取积分 需要获取是本人的信息
	Sign string `json:"sign,optional"`
}
type UserStakeRewardResp struct {
	CommResp
	Data *SignData `json:"data"` //发送到链条上广播的签名
}
type SignData struct {
	SignData  string `json:"signData"`
	Amount    string `json:"amount"`
	Recipient string `json:"recipient"`
	Custom    string `json:"custom"`
	Nonce     string `json:"nonce"`
}
