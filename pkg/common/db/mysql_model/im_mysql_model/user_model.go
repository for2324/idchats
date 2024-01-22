package im_mysql_model

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

func init() {
	for k, v := range config.Config.Manager.AppManagerUid {
		_, err := GetUserByUserID(v)
		if err != nil {

		} else {
			continue
		}
		var appMgr db.User
		appMgr.UserID = v
		if k == 0 {
			appMgr.Nickname = config.Config.Manager.AppSysNotificationName
		} else {
			if strings.EqualFold("0x8858Af738d3F7c33250d7cFd48c89196eA7Dc728", appMgr.UserID) {
				appMgr.Nickname = "IDCHATS_AI"
				appMgr.FaceURL = "ipfs://QmdsV7cijcBdyfCpCDAUx1jBMfwgy8qDTaBFL4SoKxdfVH"
			} else {
				appMgr.Nickname = "AppManager" + utils.IntToString(k+1)
			}

		}
		appMgr.AppMangerLevel = constant.AppAdmin
		err = UserRegister(appMgr)
		if err != nil {
			fmt.Println("AppManager insert error ", err.Error(), appMgr)
		} else {
			fmt.Println("AppManager insert ", appMgr)
		}
	}

	for _, userId := range config.Config.InitUser.UserId {
		_, err := GetUserByUserID(userId)
		if err != nil {
		} else {
			continue
		}
		var initUser db.User
		initUser.UserID = userId
		if strings.EqualFold("0x8858Af738d3F7c33250d7cFd48c89196eA7Dc728", userId) {
			initUser.Nickname = "IDCHATS_AI"
			initUser.FaceURL = "ipfs://QmdsV7cijcBdyfCpCDAUx1jBMfwgy8qDTaBFL4SoKxdfVH"
		} else {
			initUser.Nickname = userId
		}

		initUser.AppMangerLevel = constant.AppAdmin
		initUser.CreateTime = time.Now()
		initUser.Birth = time.Now()
		err = UserRegister(initUser)
		if err != nil {
			fmt.Println("InitUser insert error ", err.Error(), initUser)
		} else {
			fmt.Println("InitUser insert ", initUser)
		}
	}
}

func UserRegister(user db.User) error {
	user.CreateTime = time.Now()
	if user.AppMangerLevel == 0 {
		user.AppMangerLevel = constant.AppOrdinaryUsers
	}
	if user.Birth.Unix() < 0 {
		user.Birth = utils.UnixSecondToTime(0)
	}
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTokenAllUser() ([]db.User, error) {
	var userList []db.User
	// 判断 token_id 不为空 并且 face_url 不为 '1'
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("token_id != '' and face_url != '1'").Find(&userList).Error
	return userList, err
}

func GetUserByUserID(userID string) (*db.User, error) {
	if userID == "" {
		return nil, errors.New("dot input userid error")
	}
	var user db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", userID).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func GetUserEnsNameByChainID(userID string, chainid string) (string, error) {
	var resultstrinensname string
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserDomain{
		UserId:  userID,
		ChainID: chainid,
	}).Select("ens_domain").Where("user_id=? and chain_id=?", userID, chainid).First(&resultstrinensname).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	return resultstrinensname, nil
}

func GetUsersByUserIDList(userIDList []string) ([]*db.User, error) {
	var userList []*db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id in (?)", userIDList).Find(&userList).Error
	return userList, err
}

func GetUserNameByUserID(userID string) (string, error) {
	var user db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Select("name").Where("user_id=?", userID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Nickname, nil
}

func UpdateUserInfo(user db.User) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", user.UserID).Updates(&user).Error
}
func UpdateUserInfoWithMap(userID string, value map[string]interface{}) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", userID).Updates(&value).Error
}
func UpdateUserInfoWithMapping(userID string, mapValue map[string]interface{}) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", userID).Updates(mapValue).Error
}
func UpdateUserInfoByMap(user db.User, m map[string]interface{}) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", user.UserID).Updates(m).Error
	return err
}

func SelectAllUserID() ([]string, error) {
	var resultArr []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Pluck("user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func SelectSomeUserID(userIDList []string) ([]string, error) {
	var resultArr []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id IN (?) ", userIDList).Pluck("user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func GetUsers(showNumber, pageNumber int32) ([]db.User, error) {
	var users []db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, err
}

func AddUser(userID string, phoneNumber string, name string, email string, gender int32, faceURL string, birth string) error {
	_birth, err := utils.TimeStringToTime(birth)
	if err != nil {
		return err
	}
	user := db.User{
		UserID:      userID,
		Nickname:    name,
		FaceURL:     faceURL,
		Gender:      gender,
		PhoneNumber: phoneNumber,
		Birth:       _birth,
		Email:       email,
		Ex:          "",
		CreateTime:  time.Now(),
	}
	result := db.DB.MysqlDB.DefaultGormDB().Table("users").Create(&user)
	return result.Error
}

func UserIsBlock(userId string) (bool, error) {
	var user db.BlackList
	rows := db.DB.MysqlDB.DefaultGormDB().Table("black_lists").Where("uid=?", userId).First(&user).RowsAffected
	if rows >= 1 {
		return user.EndDisableTime.After(time.Now()), nil
	}
	return false, nil
}

func UsersIsBlock(userIDList []string) (inBlockUserIDList []string, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("black_lists").Where("uid in (?) and end_disable_time > now()", userIDList).Pluck("uid", &inBlockUserIDList).Error
	return inBlockUserIDList, err
}

func BlockUser(userID, endDisableTime string) error {
	user, err := GetUserByUserID(userID)
	if err != nil || user.UserID == "" {
		return err
	}
	end, err := time.Parse("2006-01-02 15:04:05", endDisableTime)
	if err != nil {
		return err
	}
	if end.Before(time.Now()) {
		return errors.New("endDisableTime is before now")
	}
	var blockUser db.BlackList
	db.DB.MysqlDB.DefaultGormDB().Table("black_lists").Where("uid=?", userID).First(&blockUser)
	if blockUser.UserId != "" {
		db.DB.MysqlDB.DefaultGormDB().Model(&blockUser).Where("uid=?", blockUser.UserId).Update("end_disable_time", end)
		return nil
	}
	blockUser = db.BlackList{
		UserId:           userID,
		BeginDisableTime: time.Now(),
		EndDisableTime:   end,
	}
	err = db.DB.MysqlDB.DefaultGormDB().Create(&blockUser).Error
	return err
}

func UnBlockUser(userID string) error {
	return db.DB.MysqlDB.DefaultGormDB().Where("uid=?", userID).Delete(&db.BlackList{}).Error
}

type BlockUserInfo struct {
	User             db.User
	BeginDisableTime time.Time
	EndDisableTime   time.Time
}

func GetBlockUserByID(userId string) (BlockUserInfo, error) {
	var blockUserInfo BlockUserInfo
	blockUser := db.BlackList{
		UserId: userId,
	}
	if err := db.DB.MysqlDB.DefaultGormDB().Table("black_lists").Where("uid=?", userId).Find(&blockUser).Error; err != nil {
		return blockUserInfo, err
	}
	user := db.User{
		UserID: blockUser.UserId,
	}
	if err := db.DB.MysqlDB.DefaultGormDB().Find(&user).Error; err != nil {
		return blockUserInfo, err
	}
	blockUserInfo.User.UserID = user.UserID
	blockUserInfo.User.FaceURL = user.FaceURL
	blockUserInfo.User.Nickname = user.Nickname
	blockUserInfo.User.Birth = user.Birth
	blockUserInfo.User.PhoneNumber = user.PhoneNumber
	blockUserInfo.User.Email = user.Email
	blockUserInfo.User.Gender = user.Gender
	blockUserInfo.BeginDisableTime = blockUser.BeginDisableTime
	blockUserInfo.EndDisableTime = blockUser.EndDisableTime
	return blockUserInfo, nil
}

func GetBlockUsers(showNumber, pageNumber int32) ([]BlockUserInfo, error) {
	var blockUserInfos []BlockUserInfo
	var blockUsers []db.BlackList
	if err := db.DB.MysqlDB.DefaultGormDB().Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&blockUsers).Error; err != nil {
		return blockUserInfos, err
	}
	for _, blockUser := range blockUsers {
		var user db.User
		if err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", blockUser.UserId).First(&user).Error; err == nil {
			blockUserInfos = append(blockUserInfos, BlockUserInfo{
				User: db.User{
					UserID:      user.UserID,
					Nickname:    user.Nickname,
					FaceURL:     user.FaceURL,
					Birth:       user.Birth,
					PhoneNumber: user.PhoneNumber,
					Email:       user.Email,
					Gender:      user.Gender,
				},
				BeginDisableTime: blockUser.BeginDisableTime,
				EndDisableTime:   blockUser.EndDisableTime,
			})
		}
	}
	return blockUserInfos, nil
}

func GetUserByName(userName string, showNumber, pageNumber int32) ([]db.User, error) {
	var users []db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where(" name like ?", fmt.Sprintf("%%%s%%", userName)).Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	return users, err
}

func GetUsersByNameAndID(content string, showNumber, pageNumber int32) ([]db.User, int64, error) {
	var users []db.User
	var count int64
	dbSql := db.DB.MysqlDB.DefaultGormDB().Table("users").Where(" name like ? or user_id = ? ",
		fmt.Sprintf("%%%s%%", content), content)
	if err := dbSql.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := dbSql.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	return users, count, err
}

func GetUserIDsByEmailAndID(phoneNumber, email string) ([]string, error) {
	dbSql := db.DB.MysqlDB.DefaultGormDB().Table("users")
	if phoneNumber == "" && email == "" {
		return nil, nil
	}
	if phoneNumber != "" {
		dbSql = dbSql.Where("phone_number = ? ", phoneNumber)
	}
	if email != "" {
		dbSql = dbSql.Where("email = ? ", email)
	}
	var userIDList []string
	err := dbSql.Pluck("user_id", &userIDList).Error
	return userIDList, err
}

func GetUsersCount(userName string) (int32, error) {
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where(" name like ? ", fmt.Sprintf("%%%s%%", userName)).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

func GetBlockUsersNumCount() (int32, error) {
	var count int64
	if err := db.DB.MysqlDB.DefaultGormDB().Model(&db.BlackList{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

type UserThirdPath struct {
	UserId          string
	Twitter         string
	DnsDomain       string
	DnsDomainVerify int32
	EnsDomain       string
	UserAddress     string
	ShowTwitter     bool
	ShowUserAddress bool
}

func GetUserAndThirdPath(userID string) (*db.YeWuUser, error) {
	if userID == "" {
		return nil, errors.New("dot input userid error")
	}
	var user db.YeWuUser
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").
		Joins("left join user_third on users.user_id=user_third.user_id").
		Select("users.*,user_third.dns_domain_verify").
		Where("users.user_id=?", userID).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func GetThirdUserInfoWithShowFlag(userid []string, chainID string) ([]*UserThirdPath, error) {
	if len(userid) == 0 {
		return nil, nil
	}
	var result []*UserThirdPath
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.User{}).
		Joins("left join user_third on users.user_id= user_third.user_id ").
		Joins("left join user_domains on users.user_id=user_domains.user_id and user_domains.chain_id = ? ", chainID).
		Where("users.user_id in (?)", userid).
		Select(`users.user_id as user_id ,case user_third.show_twitter  when 1 then user_third.twitter else "" end as twitter ,` +
			`user_third.dns_domain,user_third.weibo` +
			" ,user_third.show_facebook ,user_third.show_twitter, user_domains.ens_domain").Find(&result).Error

	fmt.Println("\n user.>>>>>>", utils.StructToJsonString(result))

	if err != nil {
		return nil, err
	}
	return result, err
}

func GetBindTwitterUserList() ([]*db.UserThird, error) {
	var result []*db.UserThird
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_third").Where("twitter != ''").Find(&result).Error
	return result, err
}

func HasBindTwitter(userId string) (bool, error) {
	twitterName := ""
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_third").Where("user_id = ?", userId).Pluck("twitter", &twitterName).Error
	if err != nil {
		return false, err
	}
	return twitterName != "", nil
}

func GetThirdUserInfoWithShowFlagWithOutDomain(userid []string) ([]*UserThirdPath, error) {
	if len(userid) == 0 {
		return nil, nil
	}
	var result []*UserThirdPath
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.User{}).
		Joins("left join user_third on users.user_id= user_third.user_id ").
		Where("users.user_id in (?)", userid).
		Select(`users.user_id as user_id ,case user_third.show_twitter  when 1 then user_third.twitter else "" end as twitter ,` +
			`user_third.dns_domain,user_third.dns_domain_verify,user_third.weibo` +
			` ,user_third.show_facebook ,user_third.show_twitter, case user_third.show_user_address when 1 then user_third.user_address else "" end as user_address  ,user_third.show_user_address`).
		Find(&result).Error

	fmt.Println("\n user.>>>>>>", utils.StructToJsonString(result))

	if err != nil {
		return nil, err
	}
	return result, err
}

func GetThirdUserInfoWithOutDomain(userid []string) ([]*UserThirdPath, error) {
	var result []*UserThirdPath
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.User{}).
		Joins("left join user_third on users.user_id= user_third.user_id ").
		Where("users.user_id in (?)", userid).Select(
		"users.user_id as user_id , user_third.twitter as twitter ," +
			"user_third.dns_domain,user_third.weibo ," +
			"user_third.show_twitter, user_third.dns_domain,user_third.dns_domain_verify, user_third.user_address,user_third.show_user_address").Find(&result).Error

	if err != nil {
		return nil, err
	}
	return result, err
}
func GetThirdUserInfo(userid []string, chainID string) ([]*UserThirdPath, error) {
	var result []*UserThirdPath
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.User{}).
		Joins("left join user_third on users.user_id= user_third.user_id ").
		Joins("left join user_domains on users.user_id=user_domains.user_id and user_domains.chain_id = ? ", chainID).
		Where("users.user_id in (?)", userid).Select(
		"users.user_id as user_id , user_third.twitter as twitter , user_third.dns_domain,user_third.dns_domain_verify,user_third.weibo ,user_third.show_twitter, user_domains.ens_domain").Find(&result).Error

	if err != nil {
		return nil, err
	}
	return result, err
}

func BindEnsDomain(userid string, ensdomain string, chainid string) error {
	var userdomain db.UserDomain
	db.DB.MysqlDB.DefaultGormDB().Where(&db.UserDomain{
		UserId:  userid,
		ChainID: chainid,
	}).Find(&userdomain)
	if userdomain.UserId != "" {
		return db.DB.MysqlDB.DefaultGormDB().Model(&db.UserDomain{
			UserId:  userid,
			ChainID: chainid,
		}).Updates(map[string]interface{}{
			"ens_domain": ensdomain,
			"chain_id":   chainid,
		}).Error

	} else {
		if ensdomain != "" {
			inseruserdomain := db.UserDomain{
				UserId:    userid,
				EnsDomain: ensdomain,
				ChainID:   chainid,
			}
			return db.DB.MysqlDB.DefaultGormDB().Create(&inseruserdomain).Error
		} else {
			return errors.New("无法设置空值")
		}

	}

}
func BindDomain(userid string, domain string) error {
	var user db.UserThird
	result := db.DB.MysqlDB.DefaultGormDB().Where(&db.UserThird{UserId: userid}).First(&user)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	if result.RowsAffected > 0 {
		// 如果找到了记录，则将其更新为新的值
		updates := map[string]interface{}{"dns_domain": domain}
		if err := db.DB.MysqlDB.DefaultGormDB().Model(&user).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	} else {
		// 如果未找到记录，则创建一个新的记录
		newUser := db.UserThird{UserId: userid, DnsDomain: domain}
		if err := db.DB.MysqlDB.DefaultGormDB().Create(&newUser).Error; err != nil {
			return err
		}
		return nil
	}
}
func BindUserEmail(userid string, emailAddress string) error {
	var user db.UserThird
	result := db.DB.MysqlDB.DefaultGormDB().Where(&db.UserThird{UserId: userid}).First(&user)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	if result.RowsAffected > 0 {
		// 如果找到了记录，则将其更新为新的值
		updates := map[string]interface{}{"user_address": emailAddress}
		if err := db.DB.MysqlDB.DefaultGormDB().Model(&user).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	} else {
		err := db.DB.MysqlDB.DefaultGormDB().Where(&db.UserThird{UserAddress: emailAddress}).First(&user).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("邮箱已经绑定了")
		}
		newUser := db.UserThird{UserId: userid, UserAddress: emailAddress}
		if err := db.DB.MysqlDB.DefaultGormDB().Create(&newUser).Error; err != nil {
			return err
		}
		return nil
	}
}

func UpdateUserPhoneNumber(userid string, phoneNumber string, isUpdate bool) error {
	dbtx := db.DB.MysqlDB.DefaultGormDB().Table("users")
	err := dbtx.Transaction(func(tx *gorm.DB) error {
		var tempPhone db.User
		err2 := tx.Model(&db.User{UserID: userid}).Find(&tempPhone).Error
		if tempPhone.PhoneNumber != "" && isUpdate == false {
			return errors.New("exist old phoneNumber cant change phoneNumber")
		}

		var tempUserCount int64
		err2 = tx.Where("phone_number =?", phoneNumber).Count(&tempUserCount).Error
		if err2 != nil {
			return err2
		}
		if !config.Config.OpenNetProxy.OpenFlag && tempUserCount > 0 {
			return errors.New("this phone was bind")
		}
		err2 = tx.Model(&db.User{UserID: userid}).Updates(map[string]interface{}{
			"phone_number": phoneNumber,
		}).Error
		return err2
	})
	return err
}
func ShowPlatformInfo(userid string, chanid string, platform string, showflag bool) error {
	switch platform {
	case "twitter":
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserThird{
			UserId: userid,
		}).Updates(map[string]interface{}{
			"show_twitter": showflag,
		}).Error
		return err
	case "facebook":
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserThird{
			UserId: userid,
		}).Updates(map[string]interface{}{
			"show_facebool": showflag,
		}).Error
		return err
		//case "balance":
		//	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserThird{
		//		UserId: userid,
		//	}).Updates(map[string]interface{}{
		//		"show_balance": showflag,
		//	}).Error
		//	return err
	}
	return nil
}
func DeletePlatformInfo(userid string, chanid string, platform string) error {
	switch platform {
	case "phone":
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.User{UserID: userid}).
			Updates(map[string]interface{}{
				"phone_number": "",
			}).Error
		return err
	case "twitter":
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserThird{
			UserId: userid,
		}).Updates(map[string]interface{}{
			"twitter": "",
		}).Error
		return err
	case "facebook":
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserThird{
			UserId: userid,
		}).Updates(map[string]interface{}{
			"facebook": "",
		}).Error
		return err
	case "ensDomain":
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserDomain{
			UserId:  userid,
			ChainID: chanid,
		}).Updates(map[string]interface{}{
			"ens_domain": "",
		}).Error
		if err == nil {
			err = db.DB.MysqlDB.DefaultGormDB().Model(&db.User{
				UserID: userid,
			}).Updates(map[string]interface{}{
				"ex": "",
			}).Error
			return err
		}
		return err
	case "faceURL":
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.User{
			UserID: userid,
		}).Updates(map[string]interface{}{
			"face_url": "",
			"token_id": "",
		}).Error
		return err
	case "dnsDomain":
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserThird{
			UserId: userid,
		}).Updates(map[string]interface{}{
			"dns_domain": "",
		}).Error
		return err
	case "email":
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserThird{
			UserId: userid,
		}).Updates(map[string]interface{}{
			"user_address": "",
		}).Error
		return err
	default:
		return errors.New("not system platform ")
	}
	return nil

}

func GetAppVersion(platform string) (*db.AppVersionFlutter, error) {
	appVersion := new(db.AppVersionFlutter)
	if err := db.DB.MysqlDB.DefaultGormDB().Table("app_version_flutter").
		Where("platform = ?", platform).Order("id desc").First(appVersion).Error; err != nil {
		return nil, err
	}
	return appVersion, nil
}
func GetChatTokenHistory(fromid string, userid string, pagecount string) ([]*db.UserChatTokenRecord, error) {
	if pagecount == "" {
		pagecount = "100"
	}
	if fromid == "" {
		result := make([]*db.UserChatTokenRecord, 0)
		err := db.DB.MysqlDB.DefaultGormDB().Table("user_chat_token_record").Where("user_id=? ", userid).Order("id desc").Limit(utils.StringToInt(pagecount)).
			Find(&result).Error
		return result, err
	} else {
		result := make([]*db.UserChatTokenRecord, 0)
		err := db.DB.MysqlDB.DefaultGormDB().Table("user_chat_token_record").Where("user_id=? and id <? ", userid, fromid).Order("id desc").Limit(utils.StringToInt(pagecount)).
			Find(&result).Error
		return result, err
	}
}

type UserNft1155 struct {
	db.CommunityChannelRole
	UserID string
}

func GetNft1155FromCommunityRoleDb(userid string) (UserNft1155Array []*UserNft1155, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Raw("select crur.user_id,ccr.group_id,ccr.token_id,ccr.contract,ccr.role_ipfs,ccr.role_title from "+
		"community_role_user_relationship crur left join community_channel_role  ccr on "+
		" crur.Contract = ccr.Contract and crur.token_id = ccr.token_id "+
		"where crur.user_id=? and amount <>'0' ", userid).Find(&UserNft1155Array).Error
	return UserNft1155Array, err
}
func InsertIntoUserNftConfig(userid string, configValueA []*db.UserNftConfig) error {
	if len(configValueA) == 0 {
		sqlErr := db.DB.MysqlDB.DefaultGormDB().Table("user_nft_config").Where("user_id=?", userid).Updates(map[string]interface{}{"is_show": 0}).Error
		return sqlErr
	}

	sqlErr := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		m := make(map[int]bool)
		for key, value := range configValueA {
			m[key] = true
			fmt.Println(value.NftContract)
		}
		var inTableDataB []*db.UserNftConfig
		err := tx.Table("user_nft_config").Where("user_id=?", userid).Find(&inTableDataB).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var showFlagID []int64
		var hideFlagID []int64
		for _, valueInDb := range inTableDataB {
			fmt.Println("check address 检查合约 :", valueInDb.NftContract)
			inNewConfig := false
			for key2, value2 := range configValueA {
				if valueInDb.Md5index == value2.Md5index {
					if valueInDb.IsShow == 0 {
						showFlagID = append(showFlagID, valueInDb.ID)
					}
					if _, ok := m[key2]; ok {
						fmt.Println("存在重复数据 已经删除", key2, value2)
						delete(m, key2)
					}
					inNewConfig = true
					break
				}
			}
			if !inNewConfig && valueInDb.IsShow == 1 {
				hideFlagID = append(hideFlagID, valueInDb.ID)
			}
		}
		if len(showFlagID) > 0 {
			err = tx.Table("user_nft_config").Where("id in ?", showFlagID).Updates(map[string]interface{}{"is_show": 1}).Error
			if err != nil {
				return err
			}
		}
		if len(hideFlagID) > 0 {
			err = tx.Table("user_nft_config").Where("id in ?", hideFlagID).Updates(map[string]interface{}{"is_show": 0}).Error
			if err != nil {
				return err
			}
		}
		var newAppend []*db.UserNftConfig
		fmt.Println("\n", m)

		for key := range m {
			newAppend = append(newAppend, configValueA[key])
		}
		if len(newAppend) > 0 {
			err = tx.Table("user_nft_config").Create(&newAppend).Error
		}
		return err
	})
	return sqlErr
}
func GetUserNFtConfigLike(id int64) (count int64) {
	//查询某个nft点赞的数量
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_nft_config_user_like_log").Where("user_nft_config_id=?", id).Count(&count).Error
	if err != nil {
		return 0
	}
	return
}
func GetUserNFtConfigLikeWithUserID(id int64, userid string) (count int64) {
	//查询某个nft点赞的数量
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_nft_config_user_like_log").Where("user_nft_config_id=? and user_id=?", id, userid).Count(&count).Error
	if err != nil {
		return 0
	}
	return
}
func LikeActionNftCount(id int64, userID string, isLike int32) error {
	if isLike == 1 {
		//关注
		err := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
			var user db.User
			var usernftlike db.UserNftConfig
			if tx.Table("users").Where("user_id=?", userID).Take(&user).Error == nil &&
				tx.Table("user_nft_config").Where("id=?", id).Take(&usernftlike).Error == nil {
				var count int64
				if tx.Table("user_nft_config_user_like_log").Where("user_nft_config_id=? and user_id=?", id, userID).Count(&count).Error == nil {
					if count == 0 {
						return tx.Table("user_nft_config_user_like_log").Create(&db.UserNftConfigUserLikeLog{
							UserID:          userID,
							UserNftConfigID: utils.Int64ToString(id),
						}).Error
					}
				}
			}
			return errors.New("文件不存在")
		})
		return err
	}
	if isLike == 2 {
		//关注
		err := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
			var user db.User
			var usernftlike db.UserNftConfig
			if tx.Table("users").Where("user_id=?", userID).Take(&user).Error == nil &&
				tx.Table("user_nft_config").Where("id=?", id).Take(&usernftlike).Error == nil {
				return tx.Table("user_nft_config_user_like_log").Where("user_nft_config_id =? and user_id =?", id, userID).Delete(&UserNftConfigInfo{}).Error
			} else {
				return errors.New("nft不存在")
			}
		})
		return err
	}
	return nil
}
func GetUserNftConfig(userid string, opUserid string) (configValue []*UserNftConfigInfo, err error) {
	if opUserid != "" {
		err = db.DB.MysqlDB.DefaultGormDB().Raw(`SELECT user_nft_config.*, COUNT(user_nft_config_user_like_log.user_id) AS like_count,
        (CASE WHEN COUNT(CASE WHEN user_nft_config_user_like_log.user_id = ? THEN 1 ELSE NULL END) > 0 THEN 1 ELSE 0 END) AS is_likes
		FROM user_nft_config
		LEFT JOIN user_nft_config_user_like_log ON user_nft_config.id = user_nft_config_user_like_log.user_nft_config_id
		where user_nft_config.user_id=?  and user_nft_config.is_show=1
		GROUP BY user_nft_config.id`, opUserid, userid).Find(&configValue).Error
	} else {
		err = db.DB.MysqlDB.DefaultGormDB().Raw(`SELECT user_nft_config.*, COUNT(user_nft_config_user_like_log.user_id) AS like_count,
        0 AS is_likes
		FROM user_nft_config
		LEFT JOIN user_nft_config_user_like_log ON user_nft_config.id = user_nft_config_user_like_log.user_nft_config_id
		where user_nft_config.user_id=?  and user_nft_config.is_show=1
		GROUP BY user_nft_config.id`, userid).Find(&configValue).Error
	}
	return
}
func GetUserLinkTree(userid string, opUserid string) (userLinkDb []*db.UserLink, err error) {
	if err = db.DB.MysqlDB.DefaultGormDB().Table("user_link").Where("user_id=?", userid).Find(&userLinkDb).Error; err != nil {
		return nil, err
	}
	return
}

type UserNftConfigInfo struct {
	db.UserNftConfig
	LikeCount int64
	IsLikes   int32
}

func DelUserAnnounceDraft(articleDraft *db.AnnouncementArticleDraft) error {
	if articleDraft.ArticleDraftID != 0 {
		return db.DB.MysqlDB.DefaultGormDB().Table("announcement_article_draft").Where(&db.AnnouncementArticleDraft{ArticleDraftID: articleDraft.ArticleDraftID}).
			Updates(map[string]interface{}{"status": 1}).Error
	} else {
		return errors.New("无效的参数，广播的id 不能为空")
	}
}

func GetTotalUserAnnounceDraft(userid string, groupID string) (result []*db.AnnouncementArticleDraft, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("announcement_article_draft").
		Where("creator_user_id=? and group_id=? and status=0", userid, groupID).Order("created_at desc").Find(&result).Error
	return
}
func GetTotalPublishGroupAnnounce(userid string, groupID string) (result []*db.AnnouncementArticle, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("announcement_article").
		Where("group_id=? and status=0", userid, groupID).Order("created_at desc").Find(&result).Error
	return
}

func DelUserAnnouncement(articleID int64) error {
	if articleID != 0 {
		return db.DB.MysqlDB.DefaultGormDB().Table("announcement_article").Where("article_id=?", articleID).
			Updates(map[string]interface{}{"status": 1}).Error
	} else {
		return errors.New("无效的参数，广播的id 不能为空")
	}
}
func CreateOrUpdateUserAnnounce(articleDraft *db.AnnouncementArticleDraft) error {
	articleDraft.UpdatedAt = time.Now()
	articleDraft.Status = 0
	if articleDraft.ArticleDraftID != 0 {
		//如果该值得不等于0 就需要判断是否是更新
		var userarticledraft db.AnnouncementArticleDraft
		result := db.DB.MysqlDB.DefaultGormDB().Table("announcement_article_draft").Where(&db.AnnouncementArticleDraft{ArticleDraftID: articleDraft.ArticleDraftID,
			CreatorUserID: articleDraft.CreatorUserID,
		}).First(&userarticledraft)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		if userarticledraft.Status == 1 {
			return errors.New("don't update has delete draft")
		}
		if result.RowsAffected > 0 {
			// 如果找到了记录，则将其更新为新的值
			updates := map[string]interface{}{
				"updated_at":           time.Now(),
				"article_draft_id":     articleDraft.ArticleDraftID,
				"creator_user_id":      articleDraft.CreatorUserID,
				"announcement_title":   articleDraft.AnnouncementTitle,
				"announcement_summary": articleDraft.AnnouncementSummary,
				"announcement_content": articleDraft.AnnouncementContent,
				"announcement_url":     articleDraft.AnnouncementUrl,
				"status":               0,
			}
			if err := db.DB.MysqlDB.DefaultGormDB().Table("announcement_article_draft").
				Where("article_draft_id =? and  creator_user_id=? and status=0 ", articleDraft.ArticleDraftID, articleDraft.CreatorUserID).Updates(updates).Error; err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("无法找更新的数据")
		}
	} else {
		var count int64
		if err := db.DB.MysqlDB.DefaultGormDB().Table("announcement_article_draft").Where("creator_user_id=? and status=0", articleDraft.CreatorUserID).Count(&count).Error; err != nil {
			return err
		}
		if count >= 5 {
			return errors.New("草稿箱已满")
		}
		articleDraft.CreatedAt = time.Now()
		if err := db.DB.MysqlDB.DefaultGormDB().Table("announcement_article_draft").Create(&articleDraft).Error; err != nil {
			return err
		}
		return nil
	}
}
