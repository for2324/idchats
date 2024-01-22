package base_info

type RobotInfo struct {
	EthAddress string `json:"ethAddress"`
	BtcAddress string `json:"btcAddress"`
	BnbAddress string `json:"bnbAddress"`
}

type CreateRobotReq struct {
	OperationID string `json:"operationID"`
}

type CreateRobotResp struct {
	CommResp
	Data map[string]string
}
type GetRobotReq struct {
	OperationID string `json:"operationID"`
}
type GetRobotResp struct {
	CommResp
	Data map[string]string `json:"data"`
}
type DelegateCallReq struct {
	OperationID string                 `json:"operationID"`
	Method      string                 `json:"method"`
	Params      map[string]interface{} `json:"params"`
}
type DelegateCallResp struct {
	CommResp
	Data interface{} `json:"data"`
}

// chainid 从httpheader 里面获取
type RobotTransactionReq struct {
	OperationID string  `json:"operationID"`
	Address     string  `json:"address"`  //转出地址
	Contract    *string `json:"contract"` //提现母币不需要传输
	Amount      string  `json:"amount"`   //提现金额，小数位 尽量用6位吧。
}
type RobotTransactionResp struct {
	CommResp
	Data string `json:"data"` //data是一个hash值判断值是否成功就好
}

type TokenPriceReq struct {
	OperationID string   `json:"operationID"`
	Tokens      []string `json:"tokens"`
}
type TokenPriceResp struct {
	CommResp
	Data []interface{} `json:"data"`
}
type CreateRobotV2Req struct {
	OperationID string `json:"operationID"`
	Sign        string `json:"sign,required"`
}
type ImportWalletMnemonic struct {
	OperationID string `json:"operationID"`
	Sign        string `json:"sign,required"`
	S           string `json:"s"`
	I           string `json:"i"`
	C           string `json:"c"`
}
