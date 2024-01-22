package im_mysql_model

import "Open_IM/pkg/common/db"

func SelectFriendsListByUserId(userId string) (list []db.Friend, err error) {
	model := db.Friend{}
	err = db.DB.MysqlDB.DefaultGormDB().Model(&model).Where("owner_user_id = ?", userId).Find(&list).Error
	return
}

type HotSpaceItem struct {
	UserId      string `gorm:"column:follow_user_id"`
	FollowCount int32  `gorm:"follow_count"`
	IsFollow    bool   `gorm:"is_follow"`
}

func GetHotSpace(userId string, pageIndex, pageSize int) (result []*HotSpaceItem, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_follow").
		Select("follow_user_id, count(*) as follow_count, SUM(from_user_id = ?) > 0 as is_follow", userId).Group("follow_user_id").
		Order("follow_count desc, follow_user_id").
		Offset(pageIndex * pageSize).
		Limit(pageSize).Find(&result).Error
	return
}

func GetMyFollowingSpace(userId string, pageIndex, pageSize int) (result []*HotSpaceItem, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_follow").
		Select("user_follow.follow_user_id, count(f.follow_user_id) as follow_count, SUM(from_user_id = ?) > 0 as is_follow", userId).
		Joins("inner join (select follow_user_id from user_follow where from_user_id = ? ) f ON user_follow.follow_user_id = f.follow_user_id", userId).
		Group("user_follow.follow_user_id").
		Order("follow_count desc, user_follow.follow_user_id").
		Offset(pageIndex * pageSize).
		Limit(pageSize).Find(&result).Error
	return
}
