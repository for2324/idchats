package im_mysql_model

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"

	"time"

	"gorm.io/gorm"

	pbGroup "Open_IM/pkg/proto/group"
)

func GetFreeGroupCount(userid string) (error, int64) {
	var nNUmber int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("creator_user_id=? and is_fees=0 and status = 0", userid).Count(&nNUmber).Error
	if err != nil {
		return err, nNUmber
	}

	return nil, nNUmber
}
func InsertIntoGroup(groupInfo db.Group) error {
	if groupInfo.GroupName == "" {
		groupInfo.GroupName = "Group Chat"
	}
	groupInfo.CreateTime = time.Now()

	if groupInfo.NotificationUpdateTime.Unix() < 0 {
		groupInfo.NotificationUpdateTime = utils.UnixSecondToTime(0)
	}
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Create(&groupInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func GetGroupInfoByGroupID(groupID string) (*db.Group, error) {
	var groupInfo db.Group
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").
		Where("group_id=?", groupID).Take(&groupInfo).Error
	return &groupInfo, err
}
func GetOneGroupInfoByUserID(userid string) (*db.Group, error) {
	var groupInfo db.Group
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Joins("left join group_members on groups.group_id= group_members.group_id").
		Select("groups.* ,count(group_members.user_id) as member_count").
		Where("groups.creator_user_id=? and groups.status=0", userid).
		Group("groups.group_id").Order(" member_count desc").
		First(&groupInfo).Error
	return &groupInfo, err
}

func SetGroupInfo(groupInfo db.Group) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id=?", groupInfo.GroupID).Updates(&groupInfo).Error
}

type GroupWithNum struct {
	db.Group
	MemberCount int `gorm:"column:num"`
}

func GetGroupsByName(groupName string, pageNumber, showNumber int32) ([]GroupWithNum, int64, error) {
	var groups []GroupWithNum
	var count int64
	sql := db.DB.MysqlDB.DefaultGormDB().Table("groups").Select("groups.*, (select count(*) from group_members where group_members.group_id=groups.group_id) as num").
		Where(" name like ? and status != ?", fmt.Sprintf("%%%s%%", groupName), constant.GroupStatusDismissed)
	if err := sql.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := sql.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&groups).Error
	return groups, count, err
}

func GetGroups(pageNumber, showNumber int) ([]GroupWithNum, error) {
	var groups []GroupWithNum
	if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Select("groups.*, (select count(*) from group_members where group_members.group_id=groups.group_id) as num").
		Limit(showNumber).Offset(showNumber * (pageNumber - 1)).Find(&groups).Error; err != nil {
		return groups, err
	}
	return groups, nil
}

func OperateGroupStatus(groupId string, groupStatus int32) error {
	group := db.Group{
		GroupID: groupId,
		Status:  groupStatus,
	}
	if err := SetGroupInfo(group); err != nil {
		return err
	}
	return nil
}

func GetGroupsCountNum(group db.Group) (int32, error) {
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where(" name like ? ", fmt.Sprintf("%%%s%%", group.GroupName)).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

func UpdateGroupInfoDefaultZero(groupID string, args map[string]interface{}) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_id = ? ", groupID).Updates(args).Error
}

func GetGroupIDListByGroupType(groupType int) ([]string, error) {
	var groupIDList []string
	if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").Where("group_type = ? ", groupType).Pluck("group_id", &groupIDList).Error; err != nil {
		return nil, err
	}
	return groupIDList, nil
}

// CreateSysUserGroup 创建官方群
// 须传入官方账号：userId，为他们创建一个群，随机一个人为群主，其他人为管理员
// 若userIDs，管理的群未满，不创建群
func CreateSysUserGroup(userIds []string) (err error) {
	tx := db.DB.MysqlDB.DefaultGormDB().Begin()
	defer func() {
		er := recover()
		if err != nil || er != nil {
			log.NewError("", "CreateSysUserGroup  Rollback", err.Error())
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if len(userIds) == 0 {
		return errors.New("user_id is nil")
	}

	row := tx.Table("groups").Where("creator_user_id in ?", userIds).Find(&[]db.Group{}).RowsAffected
	users := make([]db.User, 0, len(userIds))
	tx.Table("users").Where("user_id in ?", userIds).Find(&users)
	userMap := make(map[string]db.User)
	for _, v := range users {
		userMap[v.UserID] = v
	}
	// 随机群主
	rand.Seed(time.Now().UnixNano())
	creatorUserId := userIds[rand.Intn(len(userIds))]

	groupId := utils.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))
	bi := big.NewInt(0)
	bi.SetString(groupId[0:8], 16)
	groupId = bi.String()

	// 没有群组，初始化一个群
	if row == 0 {
		err = tx.Table("groups").Create(db.Group{GroupType: constant.NormalGroup, GroupName: constant.SysGroupName, CreatorUserID: creatorUserId, GroupID: groupId,
			CreateTime: time.Now(), NotificationUpdateTime: time.Now(), BlueVip: 1}).Error
		if err != nil {
			return err
		}

		// 添加管理员
		members := make([]db.GroupMember, 0)
		for _, uId := range userIds {
			if uId != creatorUserId {
				members = append(members, db.GroupMember{GroupID: groupId, UserID: uId, RoleLevel: constant.GroupAdmin, JoinTime: time.Now(),
					Nickname: "官方管理员", FaceURL: userMap[uId].FaceURL, MuteEndTime: time.Now()})
			} else {
				members = append(members, db.GroupMember{GroupID: groupId, UserID: uId, RoleLevel: constant.GroupOwner, JoinTime: time.Now(),
					Nickname: "群主", FaceURL: userMap[uId].FaceURL, MuteEndTime: time.Now()})
			}
		}
		err = tx.Table("group_members").Create(&members).Error
		if err != nil {
			return err
		}

		return nil
	}

	for _, userId := range userIds {
		groupId := ""
		tx.Table("groups").Select("group_id").Where("creator_user_id = ?", userId).
			Order("create_time desc").Limit(1).Take(&groupId)
		if groupId == "" {
			continue
		}
		memberNum := 0
		tx.Table("group_members").Select("count(*) as memberNum").Where("group_id = ?", groupId).
			Take(&memberNum)
		// 群未满，无需创建
		if memberNum <= 1000 {
			return nil
		}
	}

	// 群都满，创建群
	err = tx.Table("groups").Create(db.Group{GroupType: constant.NormalGroup, GroupName: constant.SysGroupName,
		CreatorUserID: creatorUserId, GroupID: groupId, CreateTime: time.Now(), NotificationUpdateTime: time.Now(), BlueVip: 1}).Error
	if err != nil {
		return err
	}

	// 添加管理员
	members := make([]db.GroupMember, 0)
	for _, uId := range userIds {
		if uId != creatorUserId {
			members = append(members, db.GroupMember{GroupID: groupId, UserID: uId, RoleLevel: constant.GroupAdmin, JoinTime: time.Now(),
				Nickname: "官方管理员", FaceURL: userMap[uId].FaceURL, MuteEndTime: time.Now()})
		} else {
			members = append(members, db.GroupMember{GroupID: groupId, UserID: uId, RoleLevel: constant.GroupOwner, JoinTime: time.Now(),
				Nickname: "群主", FaceURL: userMap[uId].FaceURL, MuteEndTime: time.Now()})
		}
	}
	err = tx.Table("group_members").Create(&members).Error
	if err != nil {
		return err
	}

	return nil
}
func InsertIntoGroupChannel(groupchannel *[]*db.GroupChannel) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("group_channel").Create(*groupchannel).Error
}
func IsGroupChannel(groupID, channelid string) bool {
	var dbdata db.GroupChannel
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_channel").Where("group_id=? and channel_id=?", groupID, channelid).First(&dbdata).Error
	if err == nil {
		return true
	}
	return false
}

func DissGroupChannel(groupID string, channelid string) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_channel").
		Where("group_id=? and channel_id=?", groupID, channelid).Updates(&db.GroupChannel{
		GroupID:       groupID,
		ChannelID:     channelid,
		ChannelStatus: constant.GroupStatusDismissed,
	}).Error
	return err
}

func UpdateGroupChannel(groupchannel *db.GroupChannel) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("group_channel").Updates(groupchannel).Error
}

func GetGroupChannelAllInfo(groupid string) (resultdata []*db.GroupChannel, err error) {
	//正常状态
	err = db.DB.MysqlDB.DefaultGormDB().Table("group_channel").Where("group_id=? and channel_status=0", groupid).
		Find(&resultdata).Error
	return
}
func GetGroupChannelByGroupIDAndChannelID(groupID string, channelID string) (resultdata *db.GroupChannel, err error) {
	//正常状态
	err = db.DB.MysqlDB.DefaultGormDB().Table("group_channel").Where("group_id=? and channel_id =? ", groupID, channelID).
		Take(&resultdata).Error
	return
}
func GetGroupChannelInfo(groupid, channelid string) (*db.GroupChannel, error) {
	var groupChannelInfo db.GroupChannel //逃逸
	err := db.DB.MysqlDB.DefaultGormDB().Table("group_channel").Where("group_id=? and channel_id=?", groupid, channelid).
		Take(&groupChannelInfo).Error
	return &groupChannelInfo, err
}

func SearchCommunityInfo(title string) ([]*db.Group, error) {
	var groupInfo []*db.Group //逃逸
	err := db.DB.MysqlDB.DefaultGormDB().Table("groups").
		Where("name=? and status=0", title).
		Find(&groupInfo).Error
	return groupInfo, err
}

type GroupAndCount struct {
	db.Group
	Yescount   int
	Totalcount int
	Ffurl      string
}

func GetHotGroupInfo() (result []*GroupAndCount, err error) {
	count := 100
	//err = db.DB.MysqlDB.DefaultGormDB().Raw(`select b.yescount,b.totalcount,users.face_url as ffurl,groups.* from (
	//select a.group_id ,
	//		count( case   DATEDIFF(join_time,NOW()) WHEN -1  then 1 else null end ) as yescount,
	//		count(*) as totalcount
	//		from group_members a group by group_id order by yescount desc, totalcount desc  	 ) b
	//		INNER   JOIN    groups on b.group_id = groups.group_id
	//		left join  users on  groups.creator_user_id =users.user_id
	//		where groups.status=0
	//		order by yescount desc, totalcount desc limit ? `, count).Find(&result).Error
	var tempresult []*GroupAndCount
	err = db.DB.MysqlDB.DefaultGormDB().Raw(`	
 select b.totalcount,users.face_url as ffurl,groups.* from (
			select group_id , count(*) as totalcount from group_members group by group_id order by  count(*) desc  ) b left  JOIN    groups on b.group_id = groups.group_id
			left join  users on  groups.creator_user_id =users.user_id
			where groups.status=0	order by totalcount desc limit  ?`, count).Find(&tempresult).Error
	indexKey := -1
	groupID := config.GetSystemGroupInfo()
	for index, value := range tempresult {
		if value.GroupID == groupID {
			indexKey = index
		}
	}

	if indexKey >= 0 {
		result = append(result, tempresult[indexKey])
		result = append(result, tempresult[:indexKey]...)
		result = append(result, tempresult[indexKey+1:]...)
		return result, nil
	}
	tempdataresult := new(GroupAndCount)
	err = db.DB.MysqlDB.DefaultGormDB().Raw(`	select groups.* ,users.face_url  as ffurl from 
	groups 
	left join users on groups.creator_user_id =users.user_id
	where groups.status=0 and groups.group_id = ?  `, groupID).First(tempdataresult).Error
	if err == nil && tempdataresult.GroupID == groupID {
		int64Data := int64(1)
		db.DB.MysqlDB.DefaultGormDB().Table("group_members").Where("group_id=?", groupID).Count(&int64Data)
		tempdataresult.Totalcount = int(int64Data)
		result = append(result, tempdataresult)
		result = append(result, tempresult...)
		return
	}
	return tempresult, nil

}
func GetBannelGroupInfo() (result []*db.GroupBanner, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("group_banner").
		Where("banner_open=1").Find(&result).Error
	return
}
func GetGroupHaveNftMemberCount(groupId string) (count int32, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Raw(`select count(users.token_contract_chain) as totalnft from groups 
    left join group_members on groups.group_id = group_members.group_id and group_members.group_id=?
left join  users on group_members.user_id = users.user_id
where groups.group_id =? and users.token_contract_chain <>""`, groupId, groupId).Pluck("totalnft", &count).Error
	return
}
func GetGroupHaveNftMemberIDListCount(groupId string) (count []string, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Raw(`select users.user_id as totalnft from groups left join group_members on 
    groups.group_id = group_members.group_id and group_members.group_id=?
left join  users on group_members.user_id = users.user_id
where groups.group_id =? and users.token_contract_chain <>""`, groupId, groupId).Pluck("user_id", &count).Error
	return
}

// IsQunZhuHadGetRewordMemberCount ,判断是否领取过了奖励
func IsQunZhuHadGetRewordMemberCount(userid, event_id string, timeNow time.Time, groupid string) bool {
	timeparam := timeNow.Format("200601")
	if !config.Config.IsPublicEnv {
		timeparam = timeNow.Format("2006010215")
	}
	if event_id == "12" {
		//判断某群主是否在某个月已经领取了
		err := db.DB.MysqlDB.DefaultGormDB().Table("event_logs").
			Where("user_id=? and event_id=? and chatwithaddress=?", userid, event_id, timeparam).First(&db.EventLogs{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
	} else if event_id == "13" {
		//判断某群主是否在某个月并且某个群已经领取了
		timeparam += groupid
		err := db.DB.MysqlDB.DefaultGormDB().Table("event_logs").
			Where("user_id=? and event_id=? and chatwithaddress=?", userid, event_id, timeparam).First(&db.EventLogs{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
	}
	return false
}
func CreateGroupRoleInformation(dbdata *db.CommunityChannelRole) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("community_channel_role").
		Create(dbdata).Error
}

func GetGroupRoleInformation(groupid string) (dbdata []*db.CommunityChannelRole, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("community_channel_role").
		Where("group_id=? and hash <> ''", groupid).Find(&dbdata).Error
	return
}

type Nft1155BurnAndTranfer struct {
	FieldName   string
	TotalAmount string
}

func GetTotalNft1155BurnAndTransfer(dbdata *db.CommunityChannelRole) (result []*Nft1155BurnAndTranfer, err error) {

	err = db.DB.MysqlDB.DefaultGormDB().Raw(
		`select 'tokenburn' as fieldname ,sum(amount)  as total_amount from community_role_user_relationship where   contract=? and token_id =?
		union
		select 'tokentotal' as fieldname  ,ifnull(sum(amount),0)  as total_amount  from community_role_user_relationship where   contract=? and token_id =?  and user_id <>?`,
		dbdata.Contract, dbdata.TokenID, dbdata.Contract, dbdata.TokenID, dbdata.CreatorAddress).Find(&result).Error
	return
}
func GetGroupMemberTagList(groupid string, tokenid string) (dbdata []*db.CommunityRoleUserRelationship, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("community_role_user_relationship").
		Where("group_id=? and token_id=? and amount <> '0'", groupid, tokenid).Find(&dbdata).Error
	return
}
func GetGroupMemberTagListByUserID(userid string, groupid string) (dbdata []string, err error) {
	dbsql := db.DB.MysqlDB.DefaultGormDB().Table("community_role_user_relationship").
		Joins("left join community_channel_role  on community_role_user_relationship.contract=community_channel_role.contract").
		Where("user_id=? and amount <> '0'", userid)
	if groupid != "" {
		dbsql.Where("group_id=?", groupid)
	}
	err = dbsql.Pluck("community_channel_role.role_ipfs", &dbdata).Error
	return
}
func CreateAnnouncement(userid string, announcementurl, title, summary string, groupid string, isGlobal int32) (dbdata *db.AnnouncementArticle, err error) {
	var dbsumber int64
	db.DB.MysqlDB.DefaultGormDB().Table("announcement_article").Unscoped().Where("group_id=?", groupid).Count(&dbsumber)
	dbdata = &db.AnnouncementArticle{
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		GroupID:             groupid,
		CreatorUserID:       userid,
		AnnouncementContent: "",
		AnnouncementUrl:     announcementurl,
		LikeCount:           0,
		RewordCount:         0,
		IsGlobal:            isGlobal,
		OrderID:             "",
		Status:              0,
		AnnouncementTitle:   title,
		AnnouncementSummary: summary,
		GroupArticleID:      dbsumber + 1,
	}

	dbsqlError := db.DB.MysqlDB.DefaultGormDB().Table("announcement_article").Create(dbdata).Error
	return dbdata, dbsqlError
}

func CreateGlobalCreateAnnouncement(userid string, announcementurl, title, summary string, groupid string, isGlobal int32) (dbdata *db.AnnouncementArticle, err error) {

	err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		var dbsumber int64
		tx.Table("announcement_article").Unscoped().Where("group_id=?", groupid).Count(&dbsumber)
		dbdata = &db.AnnouncementArticle{
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			GroupID:             groupid,
			CreatorUserID:       userid,
			AnnouncementContent: "",
			AnnouncementUrl:     announcementurl,
			LikeCount:           0,
			RewordCount:         0,
			IsGlobal:            isGlobal,
			OrderID:             "",
			Status:              0,
			AnnouncementTitle:   title,
			AnnouncementSummary: summary,
			GroupArticleID:      dbsumber + 1,
		}
		err := db.DB.MysqlDB.DefaultGormDB().Table("announcement_article").Create(dbdata).Error
		if err != nil {
			return err
		}
		dbInfo, _ := GetUserByUserID(userid)
		err = tx.Table("users").Where("user_id=?", userid).Updates(map[string]interface{}{"global_money_count": gorm.Expr("global_money_count-", config.Config.SpaceArticle.PushUsdPrice)}).Error
		if err != nil {
			return err
		}
		paramStr := fmt.Sprintf(`{"transfer":%d,"articleID":%d"}`, config.Config.SpaceArticle.PushUsdPrice, dbdata.ArticleID)
		err = tx.Create(&db.UserChatTokenRecord{
			CreatedTime: time.Now(),
			UserID:      userid,
			TxID:        time.Now().Format("20060102150405") + utils.Int64ToString(time.Now().UnixMilli()%1000) + utils.Md5(userid),
			TxType:      "subglobal",
			OldToken:    uint64(dbInfo.GlobalMoneyCount),
			ChainID:     utils.Int32ToString(dbInfo.Chainid),
			NewToken:    uint64(dbInfo.GlobalMoneyCount - int64(config.Config.SpaceArticle.PushUsdPrice)),
			NowCount:    uint64(dbInfo.GlobalMoneyCount - int64(config.Config.SpaceArticle.PushUsdPrice)),
			ParamStr:    paramStr,
		}).Error
		return err
	})
	return dbdata, err
}

func AddUserGlobalMoneyCountByArticle(userid string, TxType string, articleID string, count uint64) error {
	return db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		dbInfo, err := GetUserByUserID(userid)
		if err != nil {
			return err
		}
		err = tx.Table("users").Where("user_id=?", userid).Updates(
			map[string]interface{}{"global_money_count": gorm.Expr("global_money_count+?", count)}).Error
		if err != nil {
			return err
		}
		paramStr := fmt.Sprintf(`{"transfer":%d,"articleID":%s"}`, config.Config.SpaceArticle.PushUsdPrice, articleID)
		err = tx.Create(&db.UserChatTokenRecord{
			CreatedTime: time.Now(),
			UserID:      userid,
			TxID:        TxType + ":" + articleID,
			TxType:      TxType,
			OldToken:    uint64(dbInfo.GlobalMoneyCount),
			ChainID:     utils.Int32ToString(dbInfo.Chainid),
			NewToken:    uint64(dbInfo.GlobalMoneyCount + int64(config.Config.SpaceArticle.PushUsdPrice)),
			NowCount:    uint64(dbInfo.GlobalMoneyCount + int64(config.Config.SpaceArticle.PushUsdPrice)),
			ParamStr:    paramStr,
		}).Error
		return err
	})
}

func LikeActionAnnouncementCount(articleID int64, userID string, isLike int32) error {
	switch isLike {
	case 1: //喜欢
		//关注
		err := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
			var user db.User
			var userAnnouncement db.AnnouncementArticleLog
			if tx.Table("users").Where("user_id=?", userID).Take(&user).Error == nil &&
				tx.Table("announcement_article").Where("article_id=?", articleID).Take(&userAnnouncement).Error == nil {
				var count int64
				err := tx.Table("announcement_article_logs").Where("user_id=? and article_id =?", userID, articleID).Count(&count).Error
				if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
					if count == 0 {
						fmt.Println("系统没有数据")
						err := tx.Table("announcement_article").Where("article_id=?", articleID).Updates(map[string]interface{}{"like_count": gorm.Expr("like_count+1")}).Error
						if err != nil {
							return err
						}
						return tx.Table("announcement_article_logs").Create(&db.AnnouncementArticleLog{
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
							UserID:    userID,
							ArticleID: articleID,
							IsLikes:   1,
							Status:    0,
						}).Error
					} else {
						fmt.Println("系统有数据")
						var isLikeValueDB int64
						err := tx.Table("announcement_article_logs").Where("user_id=? and article_id=?", userID, articleID).
							Pluck("is_likes", &isLikeValueDB).Error
						if isLikeValueDB == 0 {
							err = tx.Table("announcement_article").Where("article_id=?", articleID).
								Updates(map[string]interface{}{"like_count": gorm.Expr("like_count+1")}).Error
							if err != nil {
								return err
							}
							return tx.Table("announcement_article_logs").Where("user_id=? and article_id=?", userID, articleID).
								Updates(map[string]interface{}{"is_likes": 1}).Error
						}

					}
				}
			}
			return errors.New("userAnnouncement 不存在")
		})
		return err
	case 2: //取消喜欢
		err := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
			err := tx.Table("announcement_article").Where("article_id=?", articleID).Updates(map[string]interface{}{
				"like_count": gorm.Expr("like_count-1")}).Error
			if err != nil {
				return err
			}
			return tx.Table("announcement_article_logs").Where("user_id=? and article_id =?", userID, articleID).Updates(map[string]interface{}{"is_likes": 0}).Error
		})
		return err
	case 3: //个人全局消息里面删除
		err := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
			return tx.Table("announcement_article_logs").Where("user_id=? and article_id =?", userID, articleID).Updates(map[string]interface{}{"status": 2}).Error
		})
		return err
	default:
		return nil
	}
}
func GetPublishGroupAnnounceByArticleID(articleID string) (result *db.AnnouncementArticle, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("announcement_article").
		Where("article_id =? ", articleID).Order("created_at desc").First(&result).Error
	return
}

func ColligateSearch(searchKey string) (*pbGroup.ColligateSearchBody, error) {
	// 按照 search 查找 group
	// 如果是以 0x 开头的话
	searchBody := &pbGroup.ColligateSearchBody{}
	if strings.HasPrefix(searchKey, "0x") {
		// 查找的是用户，或者空间[group.groupName]
		var users []*db.User
		if err := db.DB.MysqlDB.DefaultGormDB().Table("users").
			Select("user_id, name, face_url").
			Where("name LIKE ? or user_id = ?", searchKey+"%", searchKey).
			Limit(50).
			Find(&users).Error; err != nil {
			return nil, err
		}
		for _, user := range users {
			searchBody.UserList = append(searchBody.UserList, &pbGroup.ColligateSearchUserInfo{
				UserId:   user.UserID,
				UserName: user.Nickname,
				FaceUrl:  user.FaceURL,
			})
		}
		// 按照 groups 进行查找，按照 name 进行模糊查找，只需要匹配后缀就可以，并且只需要 group_id 和 name 字段
		var groups []*db.Group
		if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").
			Select("group_id, name, face_url").
			Where("name LIKE ?", searchKey+"%").
			Where("status=0").
			Limit(50).
			Find(&groups).Error; err != nil {
			return nil, err
		}
		for _, group := range groups {
			searchBody.GroupList = append(searchBody.GroupList, &pbGroup.ColligateSearchGroupInfo{
				GroupID:   group.GroupID,
				GroupName: group.GroupName,
				FaceUrl:   group.FaceURL,
			})
		}
	}
	// 如果是纯数字的话
	if _, err := strconv.Atoi(searchKey); err == nil {
		// 按照 groups 进行查找，按照 group_id 查找，并且只需要 group_id 和 name 字段
		var groups []*db.Group
		if err := db.DB.MysqlDB.DefaultGormDB().Table("groups").
			Select("group_id, name, face_url").
			Where("group_id LIKE ?", searchKey+"%").
			Where("status=0").
			Limit(50).
			Find(&groups).Error; err != nil {
			return nil, err
		}
		for _, group := range groups {
			searchBody.GroupList = append(searchBody.GroupList, &pbGroup.ColligateSearchGroupInfo{
				GroupID:   group.GroupID,
				GroupName: group.GroupName,
				FaceUrl:   group.FaceURL,
			})
		}
	}
	return searchBody, nil
}
