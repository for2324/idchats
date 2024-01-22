package user

import (
	"context"
	"testing"

	pbUser "Open_IM/pkg/proto/user"
)

func TestUploadOfficialNftHead(t *testing.T) {

	// 上传官方头像 模拟触发
	userId := "testUploadHeadUser"
	userId = userId
	var err error
	if err != nil {
		t.FailNow()
	}
	// 删除官方头像 模拟触发
	// err := imdb.CheckAndUpdateOrInsertIntoEventUser(userId, 1, time.Now(), true, "删除头像合约:")
	// if err != nil {
	// 	t.FailNow()
	// }
	// 重复上传头像

}

func TestGetUserInfo(t *testing.T) {
	userId := "0xcbd033ea3c05dc9504610061c86c7ae191c5c913"
	s := &userServer{}
	s.GetUserInfo(context.TODO(), &pbUser.GetUserInfoReq{
		OpUserID:   userId,
		UserIDList: []string{userId},
	})
}
