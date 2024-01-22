package group

import (
	"context"
	"testing"

	pbGroup "Open_IM/pkg/proto/group"
)

func TestQuitGroup(t *testing.T) {
	s := &groupServer{}
	{
		groupId := "2636865895"
		userId := "0xcbd033ea3c05dc9504610061c86c7ae191c5c913"
		s.QuitGroup(context.TODO(), &pbGroup.QuitGroupReq{
			GroupID:  groupId,
			OpUserID: userId,
		})
	}
	// {
	// 	groupId := "2636865895"
	// 	userId := "0xd0821e32e0b3ccada132098a5eb27cce5d208730"
	// 	s.QuitGroup(context.TODO(), &pbGroup.QuitGroupReq{
	// 		GroupID:  groupId,
	// 		OpUserID: userId,
	// 	})
	// }
}

func TestJoinGroup(t *testing.T) {
	s := &groupServer{}
	{
		groupId := "2636865895"
		userId := "0xcbd033ea3c05dc9504610061c86c7ae191c5c913"
		s.JoinGroup(context.TODO(), &pbGroup.JoinGroupReq{
			GroupID:  groupId,
			OpUserID: userId,
		})
	}
}

func TestGetHotSpace(t *testing.T) {
	s := &groupServer{}
	{
		s.GetHotSpace(context.TODO(), &pbGroup.GetHotSpaceReq{
			UserId:    "0xcbd033ea3c05dc9504610061c86c7ae191c5c913",
			PageIndex: 0,
			PageSize:  10,
		})
	}
}
func TestGetMyFollowingSpace(t *testing.T) {
	s := &groupServer{}
	{
		s.GetMyFollowingSpace(context.TODO(), &pbGroup.GetHotSpaceReq{
			UserId:    "0xcbd033ea3c05dc9504610061c86c7ae191c5c913",
			PageIndex: 0,
			PageSize:  10,
		})
	}
}
func TestGlobalPushSpaceArticle(t *testing.T) {
	s := &groupServer{}
	{
		s.GlobalPushSpaceArticle(context.TODO(), &pbGroup.GlobalPushSpaceArticleReq{
			UserId:         "0xcbd033ea3c05dc9504610061c86c7ae191c5c913",
			SpaceArticleId: "1",
		})
	}
}
