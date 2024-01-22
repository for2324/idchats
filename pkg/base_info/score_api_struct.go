package base_info

type GetUserScoreReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type UserScoreInfo struct {
	UserId string `json:"userId"`
	Score  int64  `json:"score"`
}

type GetUserScoreResp struct {
	CommResp
	Data UserScoreInfo `json:"data"`
}

type GetRewardEventLogsReq struct {
	OperationID string `json:"operationID" binding:"required"`
	RewardType  string `json:"rewardType"`
	PageIndex   int32  `json:"pageIndex"`
	PageSize    int32  `json:"pageSize" binding:"required" validate:"max=20"`
}

type GetRewardEventLogsResp struct {
	CommResp
	Data []RewardEventLogInfo `json:"data"`
}

type RewardEventLogInfo struct {
	Id         string `json:"id"`
	UserId     string `json:"userId"`
	RewardType string `json:"rewardType"`
	Reward     int64  `json:"reward"`
	Info       string `json:"info"`
	CreateTime int64  `json:"createTime"`
}

type WithdrawScoreReq struct {
	OperationID string `json:"operationID" binding:"required"`
	Score       int64  `json:"score" binding:"required"`
	Coin        string `json:"coin" commit:"BiuBiu" defalut:"BiuBiu"`
}

type WithdrawScoreRespInfo struct {
	Id int64 `json:"id"`
}
type WithdrawScoreResp struct {
	CommResp
	Data WithdrawScoreRespInfo `json:"data"`
}

type WithdrawInfo struct {
	Id         int32  `json:"id"`
	UserId     string `json:"userId"`
	Score      int64  `json:"score"`
	Amount     string `json:"amount"`
	Status     string `json:"status"`
	Remark     string `json:"remark"`
	Coin       string `json:"coin"`
	ChainId    int64  `json:"chainId"`
	TxHash     string `json:"txHash"`
	CreateTime int64  `json:"createTime"`
	Decimal    int32  `json:"decimal"`
}

type GetWithdrawScoreLogsReq struct {
	OperationID string `json:"operationID" binding:"required"`
	PageIndex   int32  `json:"pageIndex"`
	PageSize    int32  `json:"pageSize" binding:"required" validate:"max=20"`
}
type GetWithdrawScoreLogsResp struct {
	CommResp
	Data []WithdrawInfo `json:"data"`
}
