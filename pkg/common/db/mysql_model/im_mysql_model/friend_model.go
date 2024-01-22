package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func InsertToFriend(toInsertFollow *db.Friend) error {
	toInsertFollow.CreateTime = time.Now()
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Create(toInsertFollow).Error
	if err != nil {
		return err
	}
	return nil
}

func GetFriendRelationshipFromFriend(OwnerUserID, FriendUserID string) (*db.Friend, error) {
	var friend db.Friend
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Where("owner_user_id=? and friend_user_id=?", OwnerUserID, FriendUserID).
		Take(&friend).Error
	if err != nil {
		return nil, err
	}
	return &friend, err
}

func GetFriendListByUserID(OwnerUserID string) ([]db.Friend, error) {
	var friends []db.Friend
	var x db.Friend
	x.OwnerUserID = OwnerUserID
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Where("owner_user_id=?", OwnerUserID).Find(&friends).Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}
func GetFollowEachOtherUserId(OwnerUserID string) ([]string, error) {
	var friendIDList []string
	err := db.DB.MysqlDB.DefaultGormDB().Raw(" SELECT tb1.follow_user_id FROM user_follow tb1 INNER JOIN user_follow tb2 ON tb1.from_user_id=tb2.follow_user_id AND tb1.follow_user_id=tb2.from_user_id where tb1.from_user_id = ?", OwnerUserID).
		Pluck("follow_user_id", &friendIDList).Error
	if err != nil {

		return nil, err
	}
	return friendIDList, nil
}
func GetFriendIDListByUserID(OwnerUserID string) ([]string, error) {
	var friendIDList []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Where("owner_user_id=?", OwnerUserID).Pluck("friend_user_id", &friendIDList).Error
	if err != nil {
		return nil, err
	}
	return friendIDList, nil
}
func GetFollowListByFollowAndUserID(userid string, follow bool) (result []*db.User, err error) {
	if !follow { // 追随我的人
		err = db.DB.MysqlDB.DefaultGormDB().Raw(`
				select  users.* from user_follow a INNER  join users on 
				a.follow_user_id = users.user_id and  a.from_user_id=?
		`, userid).Find(&result).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return

	} else { // 我追随的人
		err = db.DB.MysqlDB.DefaultGormDB().Raw(`select  users.* from user_follow a INNER  join users on 
				a.from_user_id = users.user_id and  a.follow_user_id=?`, userid).Find(&result).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return
	}
	return

}
func CheckIsFriendCommentInsert(OwnerUserID, FriendUserID string, Remark string) error {
	_, err := GetFriendRelationshipFromFriend(OwnerUserID, FriendUserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return InsertToFriend(&db.Friend{
			OwnerUserID:    OwnerUserID,
			FriendUserID:   FriendUserID,
			Remark:         Remark,
			CreateTime:     time.Now(),
			AddSource:      0,
			OperatorUserID: OwnerUserID,
			Ex:             "",
		})
	} else if err != nil {
		return err
	} else {
		return UpdateFriendComment(OwnerUserID, FriendUserID, Remark)
	}
}
func UpdateFriendComment(OwnerUserID, FriendUserID, Remark string) error {
	return db.DB.MysqlDB.DefaultGormDB().Exec("update friends set remark=? where owner_user_id=? and friend_user_id=?", Remark, OwnerUserID, FriendUserID).Error
}

func DeleteSingleFriendInfo(OwnerUserID, FriendUserID string) error {
	return db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		err := tx.Where("owner_user_id=? and friend_user_id=?", OwnerUserID, FriendUserID).
			Delete(db.Friend{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("from_user_id=? and follow_user_id=?", OwnerUserID, FriendUserID).Delete(db.UserFollow{}).Error
		if err != nil {
			return err
		}
		return nil
	})
}

type FriendUser struct {
	db.Friend
	Nickname string `gorm:"column:name;size:255"`
}

func GetUserFriendsCMS(ownerUserID, friendUserName string, pageNumber, showNumber int32) (friendUserList []*FriendUser, count int64, err error) {
	db := db.DB.MysqlDB.DefaultGormDB().Table("friends").
		Select("friends.*, users.name").
		Where("friends.owner_user_id=?", ownerUserID).Limit(int(showNumber)).
		Joins("left join users on friends.friend_user_id = users.user_id").
		Offset(int(showNumber * (pageNumber - 1)))
	if friendUserName != "" {
		db = db.Where("users.name like ?", fmt.Sprintf("%%%s%%", friendUserName))
	}
	if err = db.Count(&count).Error; err != nil {
		return
	}
	err = db.Find(&friendUserList).Error
	return
}

func GetFriendByIDCMS(ownerUserID, friendUserID string) (friendUser *FriendUser, err error) {
	friendUser = &FriendUser{}
	err = db.DB.MysqlDB.DefaultGormDB().Table("friends").
		Select("friends.*, users.name").
		Where("friends.owner_user_id=? and friends.friend_user_id=?", ownerUserID, friendUserID).
		Joins("left join users on friends.friend_user_id = users.user_id").
		Take(friendUser).Error
	return friendUser, err
}
func CheckIsFollowEachOther(ownerAddress, toAddress string) bool {
	var users []*db.UserFollow
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_follow").Raw("select user_follow.* from user_follow where from_user_id=? and follow_user_id =? union "+
		"select user_follow.* from user_follow where from_user_id =? and follow_user_id=? ", ownerAddress, toAddress, toAddress, ownerAddress).Scan(&users).Error
	if err != nil {
		return false
	} else {
		if len(users) >= 2 {
			return true
		}
	}
	return false
}

func GetUserFollowingCount(userID string) (count int64, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_follow").Where("from_user_id=?", userID).Count(&count).Error
	return
}

func GetUserFollowsCount(userID string) (count int64, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_follow").Where("follow_user_id=?", userID).Count(&count).Error
	return
}

func GetUserFollowingList(userID string, pageNumber, showNumber int32) (userList []*db.User, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_follow").
		Select("users.*").
		Where("user_follow.from_user_id=?", userID).Limit(int(showNumber)).
		Joins("left join users on user_follow.follow_user_id = users.user_id").
		Offset(int(showNumber * pageNumber)).
		Find(&userList).Error
	return
}

func GetUserFollowedList(userID string, pageNumber, showNumber int32) (userList []*db.User, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_follow").
		Select("users.*").
		Where("user_follow.follow_user_id=?", userID).Limit(int(showNumber)).
		Joins("left join users on user_follow.from_user_id = users.user_id").
		Offset(int(showNumber * pageNumber)).
		Find(&userList).Error
	return
}
func GetAllOpenGlobalPushUsers(userID string, pageNumber, showNumber int32) (userList []*db.User, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Raw(`
		SELECT DISTINCT users.* from users
	JOIN user_follow ON users.user_id = user_follow.from_user_id
	where user_follow.follow_user_id = ?
	and users.user_id  <>? and users.open_announcement=0 
	UNION
	select users.* from users where users.open_announcement =1
	and users.user_id<>?
	limit ?,?`, userID, userID, userID, showNumber*pageNumber, showNumber).Find(&userList).Error
	//Where("open_announcement=1 and user_id<> ?", userID).Limit(int(showNumber)).
	//Offset(int(showNumber * pageNumber)).
	//Find(&userList).Error
	return
}

func IsFollowUser(fromUserID, followUserID string) (isFollow bool, err error) {
	var count int64
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_follow").
		Where("from_user_id=? and follow_user_id=?", fromUserID, followUserID).
		Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		isFollow = true
	}
	return
}
