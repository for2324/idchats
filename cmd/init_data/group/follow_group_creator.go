package group

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"log"

	pbFriend "Open_IM/pkg/proto/friend"
	pbGroup "Open_IM/pkg/proto/group"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"

	rpcFriend "Open_IM/internal/rpc/friend"
	rpcGroup "Open_IM/internal/rpc/group"

	"gorm.io/gorm"
)

type FollowItem struct {
	UserId        string `json:"userId"`
	CreatorUserId string `json:"creatorUserId"`
}

// INSERT INTO user_follow (from_user_id, follow_user_id, handle_result, follow, create_time,handle_time)
// SELECT DISTINCT *, 0, 1,NOW(), NOW() FROM (
//     SELECT gm.user_id, g.creator_user_id
//     FROM group_members
//     JOIN groups g ON g.group_id = gm.group_id
//     WHERE NOT EXISTS (
//     SELECT 1
//     FROM user_follow uf
//     WHERE uf.from_user_id = gm.user_id
//     AND uf.follow_user_id = g.creator_user_id
//     )
// ) im WHERE user_id <> creator_user_id

func FollowGroupCreator() {
	// 查询所有 group_member ,查看是否
	var followList []FollowItem
	// 添加好友
	err := db.DB.MysqlDB.DefaultGormDB().Raw(`SELECT DISTINCT user_id, creator_user_id FROM (
		SELECT gm.user_id, g.creator_user_id
		FROM group_members gm
		JOIN groups g ON g.group_id = gm.group_id
		WHERE NOT EXISTS (
		SELECT 1
		FROM user_follow uf
		WHERE uf.from_user_id = gm.user_id
		AND uf.follow_user_id = g.creator_user_id
		)
	) im WHERE user_id <> creator_user_id`).Find(&followList).Error
	if err != nil {
		panic(err)
	}
	for _, groupMember := range followList {
		if err := FollowAddFriend("follow_group_creator", groupMember.UserId, groupMember.CreatorUserId); err != nil {
			log.Println("follow_group_creator", "FollowAddFriend failed", err.Error())
		} else {
			log.Println("follow_group_creator", "FollowAddFriend success", groupMember.UserId, groupMember.CreatorUserId)
		}
	}
}

func FollowAddFriend(operationID, formUserId, toUserId string) error {
	req := &pbFriend.FollowAddFriendReq{
		CommID: &pbFriend.CommID{
			OpUserID:    formUserId,
			FromUserID:  formUserId,
			ToUserID:    toUserId,
			OperationID: operationID,
		},
		Follow: true,
	}
	// etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
	// 	strings.Join(config.Config.Etcd.EtcdAddr, ","),
	// 	config.Config.RpcRegisterName.OpenImFriendName,
	// 	req.CommID.OperationID,
	// )
	// if etcdConn == nil {
	// 	errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
	// 	log.Println(req.CommID.OperationID, errMsg)
	// 	return errors.New(errMsg)
	// }
	// client := pbFriend.NewFriendClient(etcdConn)

	client := rpcFriend.NewFriendServer(0)
	RpcResp, err := client.FollowAddFriend(context.Background(), req)
	if err != nil {
		log.Println(req.CommID.OperationID, "AddFriend failed ", err.Error(), req.String())
		return err
	}
	if RpcResp.CommonResp.ErrCode != 0 {
		return errors.New(RpcResp.CommonResp.ErrMsg)
	}
	log.Println(operationID, "FollowAddFriend success", formUserId, toUserId)
	return nil
}

type GroupCreatorItem struct {
	CreatorUserID string `gorm:"column:creator_user_id;size:64"`
}
type GroupItem struct {
	CreatorUserID string `gorm:"column:creator_user_id;size:64"`
	GroupId       string `gorm:"column:group_id;size:64"`
	MemberCount   int32  `gorm:"column:member_count"`
}

func AutoDismissGroup() {
	// 查看所有的 group 查看是否有重复的，有的话只保留一个

	creators := []GroupCreatorItem{}
	// 查看所有有创建群的用户
	err := db.DB.MysqlDB.DefaultGormDB().Raw(`SELECT DISTINCT creator_user_id from groups`).Find(&creators).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("AutoCreateGroup", "no group creator")
			return
		}
		panic(err)
	}
	// 查看他是否包含多个有效（未解散的）群
	for _, creator := range creators {
		groups := []GroupItem{}
		// 按照群成员人数进行排序
		err := db.DB.MysqlDB.DefaultGormDB().Raw(
			`SELECT g.group_id, g.creator_user_id, COUNT(m.user_id) AS member_count
			FROM groups g
			LEFT JOIN group_members m ON g.group_id = m.group_id
			WHERE g.creator_user_id = ? AND g.status = 0
			GROUP BY g.group_id
			ORDER BY member_count DESC`,
			creator.CreatorUserID,
		).Find(&groups).Error
		if err != nil {
			fmt.Println("", err.Error())
			continue
		}
		if len(groups) > 1 {
			// 只保留第一个（人数最高的），删除多余的群
			for i := 1; i < len(groups); i++ {
				log.Println("AutoDismissGroup", "dismiss group", groups[i].GroupId)
				err := dismissGroupRpc(creator.CreatorUserID, groups[i].GroupId)
				if err != nil {
					panic(err)
				}
			}
		}
	}

}

type UserItem struct {
	UserID string `gorm:"column:user_id;size:64"`
}

func AutoCraeteUserGroup() {
	// 查看没有建群的用户（status = 0）
	operationID := "init data AutoCraeteUserGroup"
	users := []UserItem{}
	err := db.DB.MysqlDB.DefaultGormDB().Raw(`SELECT u.user_id
	FROM users u
	LEFT JOIN (select * FROM groups WHERE status= 0) g ON u.user_id = g.creator_user_id WHERE g.group_id IS NULL`).Find(&users).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("AutoCreateGroup", "no user")
			return
		}
		panic(err)
	}
	for _, user := range users {
		// 创建群
		err := AutoCreateCommunity(operationID, user.UserID)
		if err != nil {
			log.Println(operationID, "AutoCreateCommunity failed", err.Error())
			panic(err)
		}
	}
}

func dismissGroupRpc(opUserId, groupId string) error {
	OperationID := "init data dismissGroupRpc"
	// etcdConn := getcdv3.GetDefaultConn(
	// 	config.Config.Etcd.EtcdSchema,
	// 	strings.Join(config.Config.Etcd.EtcdAddr, ","),
	// 	config.Config.RpcRegisterName.OpenImGroupName, OperationID)
	// if etcdConn == nil {
	// 	errMsg := OperationID + "getcdv3.GetDefaultConn == nil"
	// 	log.Println(OperationID, errMsg)
	// 	panic(errMsg)
	// }
	// client := pbGroup.NewGroupClient(etcdConn)
	client := rpcGroup.NewGroupServer(0)
	req := &pbGroup.DismissGroupReq{
		OpUserID:    opUserId,
		GroupID:     groupId,
		OperationID: OperationID,
	}
	reply, err := client.DismissGroup(context.Background(), req)
	if err != nil {
		log.Println(OperationID, utils.GetSelfFuncName(), " failed ", req.String())
		return err
	}
	if reply.CommonResp.ErrCode != 0 {
		log.Println(OperationID, utils.GetSelfFuncName(), " failed ", req.String())
		return errors.New(reply.CommonResp.ErrMsg)
	}
	log.Println(OperationID, utils.GetSelfFuncName(), " success ", req.String())
	return nil
}

func AutoCreateCommunity(operationID, userID string) error {
	req := pbGroup.CreateCommunityReq{
		OperationID: operationID,
		OpUserID:    userID,
		OwnerUserID: userID,
		GroupInfo: &sdk_ws.GroupInfo{
			GroupID:     userID,
			GroupName:   userID,
			OwnerUserID: userID,
			GroupType:   constant.WorkingGroup,
			FaceURL:     "",
		},
	}
	// etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
	// 	strings.Join(config.Config.Etcd.EtcdAddr, ","),
	// 	config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	// if etcdConn == nil {
	// 	panic("etcdConn is nil")
	// }
	// client := pbGroup.NewGroupClient(etcdConn)
	client := rpcGroup.NewGroupServer(0)
	RpcResp, err := client.CreateCommunity(context.Background(), &req)
	if err != nil {
		log.Println(req.OperationID, "CreateCommunity failed ", err.Error(), req.String())
		return err
	}
	if RpcResp.ErrCode != 0 {
		log.Println(req.OperationID, "CreateCommunity failed ", RpcResp.ErrMsg, req.String())
		return errors.New(RpcResp.ErrMsg)
	}
	log.Println(req.OperationID, "CreateCommunity ok ", req.String())
	return nil
}
