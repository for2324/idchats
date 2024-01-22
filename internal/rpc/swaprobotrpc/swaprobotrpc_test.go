package swaprobotrpc

import (
	pbSwapRobot "Open_IM/pkg/proto/swaprobot"
	"context"
	"encoding/json"
	"testing"
)

func TestSwao(t *testing.T) {
	s := NewSwapRobotServer(0)
	valMap := make(map[string]interface{})
	valMap["ordId"] = "1691391976072"
	valMap["searchBy"] = "ordId"

	valData, err := json.Marshal(valMap)
	if err != nil {
		t.Fatal(err)
		return
	}
	resp, err := s.BotOperation(context.TODO(), &pbSwapRobot.BotOperationReq{
		OperatorID: "123",
		Method:     "getTask",
		UserID:     "0x1c8ec996420db47c0859a5aaed0148fd23426bbf",
		Params:     valData,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(resp)
}
