package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

type RobotRunTaskResp struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}
type Param struct {
	FromSymbol  string  `json:"fromSymbol"`
	Amount      float64 `json:"amount"`
	ToSymbol    string  `json:"toSymbol"`
	Tp          string  `json:"tp"`
	Sl          string  `json:"sl"`
	OrdStatus   string  `json:"ordStatus"`
	MinimumOut  float64 `json:"minimumOut"`
	DeadlineDay string  `json:"deadlineDay"`
}
type SwapBotInfoParam struct {
	OrdId  string `json:"ordId"`
	Params Param  `json:"params"`
}

func TestMapTypStruct(t *testing.T) {
	respData := `{"code":200,"data":{"method":"createTask","msg":"OK","ordId":"1690253061205","ordStatus":"exPending","params":{"amount":0,"deadlineMin":"5","fromSymbol":"BUSD","minimumOut":1e-21,"sl":"0.2","toSymbol":"MATIC","tp":"1.0"},"privateKey":"819612ea784589631438d5a094520912a6b66cc306f4a622011a5588204283a8","taskStartTimestamp":"1690253061"},"msg":"OK"}`
	httpResp := RobotRunTaskResp{}
	err := json.Unmarshal([]byte(respData), &httpResp)
	if err != nil {

		return
	}
	if httpResp.Code != 200 {

		return
	}
	var tempResultData SwapBotInfoParam
	if err = Map2Struct(httpResp.Data, &tempResultData); err != nil {

		return
	}
	fmt.Println(StructToJsonString(tempResultData))
}
