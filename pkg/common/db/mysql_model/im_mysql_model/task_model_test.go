package im_mysql_model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinishDailyChatNFTHeadWithNewUserTask(t *testing.T) {
	userId := "0xcbd033ea3c05dc9504610061c86c7ae191c5c913"
	isFinishToday, err := IsFinishDailyChatNFTHeadWithNewUserTask("0xcbd033ea3c05dc9504610061c86c7ae191c5c913", "1")
	if err != nil {
		t.FailNow()
	}

	if !isFinishToday {
		FinishDailyChatNFTHeadWithNewUserTask(userId, "2")
	}
	isFinishToday, err = IsFinishDailyChatNFTHeadWithNewUserTask(userId, "1")
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, isFinishToday, true)
}

func TestFinishOfficialNFTHeadDailyChatWithNewUserTask(t *testing.T) {
	isFinishToday, err := IsFinishOfficialNFTHeadDailyChatWithNewUserTask("1", "2")
	if err != nil {
		t.FailNow()
	}

	if !isFinishToday {
		FinishOfficialNFTHeadDailyChatWithNewUserTask("1", "2")
	}
	isFinishToday, err = IsFinishOfficialNFTHeadDailyChatWithNewUserTask("1", "2")
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, isFinishToday, true)
}

func TestIsFinishDailyChatNFTHeadWithNewUserTask(t *testing.T) {
	isFinishToday, err := IsFinishDailyChatNFTHeadWithNewUserTask("0xcbd033ea3c05dc9504610061c86c7ae191c5c913", "1")
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, isFinishToday, false)
}

func TestFinishInviteBindTwitterTask(t *testing.T) {
	FinishInviteBindTwitterTask("0xcbd033ea3c05dc9504610061c86c7ae191c5c913", "0xcbd033ea3c05dc9504610061c86c7ae191c5c913")
}

func TestFinishUserDailyCheckInTask(t *testing.T) {
	FinishUserDailyCheckInTask("0xcbd033ea3c05dc9504610061c86c7ae191c5c913")
}
