package main

import (
	"context"
	"testing"
)

func TestUploadUserScore(t *testing.T) {
	cront := crontabServer{}
	cront.EventTaskReward(context.TODO())
}

func TestEventInsertIntoRecord(t *testing.T) {
	cront := crontabServer{}
	cront.EventInsertIntoRecord(context.TODO())
}

func Test_crontabServer_RechargeTotalReword(t *testing.T) {
	cront := crontabServer{}
	cront.RechargeTotalReword(context.TODO())
}
