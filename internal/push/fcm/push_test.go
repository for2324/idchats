package fcm

import (
	"Open_IM/internal/push"
	"fmt"
	"testing"
)

func Test_Push(t *testing.T) {
	offlinePusher := NewFcm()
	resp, err := offlinePusher.Push([]string{"test_uid"}, "test", "test", "12321", push.PushOpts{})
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	fmt.Println(resp)
}
