package web3pub

import (
	"context"
	"testing"
)

func TestFollowTwitter(t *testing.T) {

}

func Test_web3pubserver_InitCrawlingTwitterFollow(t *testing.T) {
	web3pubserve := NewWeb3PubServer(0)
	isFollow := web3pubserve.InitCrawlingTwitterFollow(context.Background(), "123123123","knightxv2")
	if !isFollow {
		t.FailNow()
		return
	}
}
