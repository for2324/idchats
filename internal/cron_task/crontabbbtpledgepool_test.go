package cronTask

import (
	"fmt"
	"testing"
)

func Test_getBLpPledgeTotalLock(t *testing.T) {
	data, err := getBLpPledgeTotalLock()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(data)
	}

}

func TestCheckBrc20PledgePoolVolume(t *testing.T) {

	CheckBrc20PledgePoolVolume()

}
