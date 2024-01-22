package friend

import (
	"context"
	"testing"

	pbFriend "Open_IM/pkg/proto/friend"
)

func TestGetFriendFollowList(t *testing.T) {
	s := &friendServer{}
	userId := "0xcbd033ea3c05dc9504610061c86c7ae191c5c913"
	s.GetFriendFollowList(context.TODO(), &pbFriend.GetFriendFollowListReq{
		IsFollow: true,
		CommID: &pbFriend.CommID{
			OpUserID:   userId,
			ToUserID:   userId,
			FromUserID: userId,
		},
	})
}

func TestFollowAddFriend(t *testing.T) {
	s := &friendServer{}
	userId := "0xcbd033ea3c05dc9504610061c86c7ae191c5c913"
	s.FollowAddFriend(context.TODO(), &pbFriend.FollowAddFriendReq{
		CommID: &pbFriend.CommID{
			OpUserID:   userId,
			FromUserID: userId,
			ToUserID:   "0x08175c7396f9491dca656f0e1b9ae66539382960",
		},
		Follow: false,
	})
}
