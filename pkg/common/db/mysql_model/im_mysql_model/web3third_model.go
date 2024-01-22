package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GetUserTwitter(userid string) (twitterstring string, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").Where("user_id=? ", userid).Pluck("twitter", &twitterstring).Error
	return
}

func GetUserTwitterWithFlagAll(userid string) (userThirdInfo *db.UserThird, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").Where("user_id=? ", userid).First(&userThirdInfo).Error
	return
}
func IsExistThird(thirdname string, thirdValue string) (resultCount int64, err error) {
	switch thirdname {
	case "twitter":
		err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").
			Where("twitter=? ", thirdValue).Count(&resultCount).Error
	case "facebook":
		err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").
			Where("facebook=? ", thirdValue).Count(&resultCount).Error
	}
	return
}

func InsertIntoEventList(dbValue *db.EventLogs) bool {
	err := db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Create(dbValue).Error
	if err != nil {
		return false
	}
	return true
}
func CheckIsHavePhoneEvent(userid, taskid string) bool {
	//1次性任务
	var count int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Where("user_id=? and event_id=?", userid, taskid).Count(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func GetWhiteList(ctx context.Context) (result []string, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Table("user_white_list").Pluck("white_list_address", &result).Error
	return result, err
}

//	func GetTaskList() (result []*db.EventTables, err error) {
//		err = db.DB.MysqlDB.DefaultGormDB().Table("event_tables").Where("is_open=1").Find(&result).Error
//		return result, err
//	}
func GetTaskLogList(userid string, taskId int) (result []*db.EventLogs, err error) {
	if taskId >= 4 && taskId <= 7 {
		err = db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Where("user_id=? and event_id=? and created_at >= DATE_FORMAT(CURDATE(),'%Y-%m-%d %H:%i:%s')", userid, taskId).Find(&result).Error
		return
	}
	err = db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Where("user_id=? and event_id=?", userid, taskId).Find(&result).Error
	return
}
func GetTaskListByTaskId(taskid string) (result *db.EventTables, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("event_tables").Where("event_id=?", taskid).Find(&result).Error
	return result, err
}
func GetUserIsBindThisTwitter(twittername string) (userid string) {
	checkEventInfo := "bind twitter:" + twittername
	db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Where("user_json=?", checkEventInfo).Pluck("user_id", &userid)
	return
}
func FinishTaskId(userid, taskID string, exstring string) error {
	eventTask, _ := GetTaskListByTaskId(taskID)
	if eventTask.IsOpen {
		if utils.StringToInt32(taskID) <= 3 {
			// 1此行版定事件
			if !CheckIsHavePhoneEvent(userid, taskID) {
				insertEvent := &db.EventLogs{
					CreatedAt:     time.Now(),
					UserID:        userid,
					EventID:       taskID,
					EventTypename: eventTask.EventTypename,
					UserTaskcount: 1,
					UserJSON:      exstring,
					IsSync:        0,
				}
				InsertIntoEventList(insertEvent)
				return nil
			} else {
				return errors.New("你已经绑定这个事件")
			}
		} else { // 拉人头事件
			insertEvent := &db.EventLogs{
				CreatedAt:     time.Now(),
				UserID:        userid,
				EventID:       taskID,
				EventTypename: eventTask.EventTypename,
				UserTaskcount: 1,
				UserJSON:      exstring,
				IsSync:        0,
			}
			InsertIntoEventList(insertEvent)
			return nil
		}
	}
	return nil
}
func FinishTaskIdWithTalkAbout(userid, taskID string, exstring string, chatwithOne string) error {
	eventTask, _ := GetTaskListByTaskId(taskID)
	if eventTask.IsOpen {
		if utils.StringToInt32(taskID) >= 4 && utils.StringToInt32(taskID) <= 7 {
			insertEvent := &db.EventLogs{
				CreatedAt:       time.Now(),
				UserID:          userid,
				EventID:         taskID,
				EventTypename:   eventTask.EventTypename,
				Chatwithaddress: chatwithOne,
				UserTaskcount:   1,
				UserJSON:        exstring,
				IsSync:          0,
			}
			InsertIntoEventList(insertEvent)
			return nil
		}
	}
	return nil
}

type TodayFinishChatEvent struct {
	EventId string
	Total   int64
	UserId  string
}

func IsFinishChatToday(userid string) (isFinishChatEachOther, isFinishMoShengReng, isNftFinishChatEachOther, isNftFinishMoShengReng bool) {
	var resultTotal []*TodayFinishChatEvent
	err := db.DB.MysqlDB.DefaultGormDB().
		Raw(`select event_id ,count(event_id) as total ,user_id from event_logs where user_id = ? and created_at>=DATE_FORMAT(CURDATE(),'%Y-%m-%d %H:%i:%s')
				and  event_id between 4 and 7 GROUP BY event_id,user_id`, userid).
		Scan(&resultTotal).Error
	if err != nil {
		return false, false, false, false
	} else {
		for _, value := range resultTotal {
			switch value.EventId {
			case "4":
				if value.Total >= 1 {
					isFinishChatEachOther = true
				}
			case "5":
				if value.Total >= 1 {
					isFinishMoShengReng = true
				}
			case "6":
				if value.Total >= 1 {
					isNftFinishChatEachOther = true
				}
			case "7":
				if value.Total >= 2 {
					isNftFinishMoShengReng = true
				}
			}
		}
		return
	}
}
func IsChatSamePeople(userid string, chatwithOne string, eventId string) bool {
	var isExistThisChat db.EventLogs
	err := db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Where(`user_id = ? and chatwithaddress = ? and  event_id =? 
				and created_at>=DATE_FORMAT(CURDATE(),'%Y-%m-%d %H:%i:%s')`, userid, chatwithOne, eventId).
		Select("event_logs.*").First(&isExistThisChat).Error
	if err == nil {
		return true
	}
	return false
}
func GetUnSyncEventFromID(count int, fromID int64) (resultList []*db.EventLogs, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Where(`id>  ? and is_sync = 0 and event_id not between 4 and 7 `, fromID).
		Select("event_logs.*").Limit(count).Find(&resultList).Error
	return
}
func UpdateEventLogNotChat(id int64, status int) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Where(`id=?`, id).Updates(map[string]interface{}{
		"is_sync": status,
	}).Error
	return err
}

type TempEventLogs struct {
	db.EventLogs
	Dday time.Time
}

func GetUnSyncEventFromIDChatEvent(count int64) (resultList []*TempEventLogs, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Raw(`
select  CAST(DATE_FORMAT(event_logs.created_at,'%Y-%m-%d') as DATETIME) as dday, user_id ,event_id,event_typename,count(user_id) as number from event_logs where event_typename='chat_5'
		 and is_sync=0 and created_at >= CURRENT_DATE()
		GROUP BY dday,user_id ,event_id,event_typename  HAVING(number)>=1 
		UNION
		select  CAST(DATE_FORMAT(event_logs.created_at,'%Y-%m-%d') as DATETIME) as dday, user_id ,event_id,event_typename,count(user_id) as number from event_logs where event_typename='chat_10'
			and is_sync=0 and created_at >= CURRENT_DATE()
			GROUP BY dday,user_id,event_id,event_typename  HAVING(number)>=1
			UNION
			select  CAST(DATE_FORMAT(event_logs.created_at,'%Y-%m-%d') as DATETIME) as dday, user_id ,event_id,event_typename,count(user_id) as number from event_logs where event_typename='chat_1'
				 and is_sync=0 and created_at >= CURRENT_DATE()
				GROUP BY dday,user_id,event_id  HAVING(number)>=1 
				UNION
				select  CAST(DATE_FORMAT(event_logs.created_at,'%Y-%m-%d') as DATETIME) as dday,user_id ,event_id,event_typename,count(user_id) as number from event_logs where event_typename='chat_2'
					 and is_sync=0 and created_at >= CURRENT_DATE()
					GROUP BY dday,user_id,event_id,event_typename  HAVING(number)>=2   limit ?`, count).Find(&resultList).Error
	return
}

func UpdateEventLogUserWithChat(resultEvetnLog *db.EventLogs, status int, operatortime time.Time) error {
	//beginTime, endTime := utils.GetDateTimeBeginTimeAndEndTimeByInputTime(operatortime)
	err := db.DB.MysqlDB.DefaultGormDB().Table("event_logs").
		Where(`user_id = ? and event_typename=? and DATE_FORMAT(event_logs.created_at,'%Y-%m-%d')=? `,
			resultEvetnLog.UserID, resultEvetnLog.EventTypename, resultEvetnLog.CreatedAt.Format("2006-01-02")).Updates(map[string]interface{}{
		"is_sync": status,
	}).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
func CheckIsSyncPhoneEvent(address string) bool {
	var isSyncValue int
	db.DB.MysqlDB.DefaultGormDB().Table("event_logs").Where("user_id =? and event_typename='phone'", address).
		Pluck("is_sync", &isSyncValue)
	if isSyncValue == 1 || isSyncValue == 3 {
		return true
	}
	return false

}

// 获取某条链上的钱币
func GetChainTokenByChainID(chainid string) (result []*db.ChainToken, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("chain_token").
		Where("coin_chainid=?", chainid).
		Find(&result).Error
	return
}

// 插入新的token 钱币
func PostChainNewToken(result *[]*db.ChainToken) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("chain_token").Create(result).Error
}
func IsInsertCoinToken(chainid string, tokenaddress string) bool {
	var tempdata db.ChainToken
	err := db.DB.MysqlDB.DefaultGormDB().Table("chain_token").
		Where("coin_chainid=? and coin_token=?", chainid, tokenaddress).First(&tempdata).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	return false
}
func InsertIntoEmail(emailAddress string, emailPassword string) {
	db.DB.MysqlDB.DefaultGormDB().Table("email_user_system").
		Create(&db.EmailUserSystem{
			EmailAddress:  emailAddress,
			EmailPassword: emailPassword})
}
func InsertIntoEmailWithPrivateKey(emailAddress string, emailPassword string, privateKey string) {
	db.DB.MysqlDB.DefaultGormDB().Table("email_user_system").
		Create(&db.EmailUserSystem{
			EmailAddress:      emailAddress,
			EmailPassword:     emailPassword,
			EncryptPrivateKey: privateKey})
}
func GetEmailInfo(emailAddress string) (result *db.EmailUserSystem, err error) {

	resuldata := new(db.EmailUserSystem)
	err = db.DB.MysqlDB.DefaultGormDB().Table("email_user_system").
		Where("email_address=?", emailAddress).First(resuldata).Error
	if err != nil {
		return nil, err
	}
	return resuldata, nil
}
func UpdateIntoEmail(emailAddress string, emailPassword string) {
	db.DB.MysqlDB.DefaultGormDB().Table("email_user_system").
		Updates(&db.EmailUserSystem{EmailAddress: emailAddress, EmailPassword: emailPassword})
}
func CheckIsHaveThisEmail(emailAddress string) bool {
	var userEmail db.EmailUserSystem
	err := db.DB.MysqlDB.DefaultGormDB().Table("email_user_system").
		Where("email_address = ? ", emailAddress).First(&userEmail).Error
	if err == nil {
		return true
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return false
}
func GetSystemOfficialNftInfo() (result []*db.SystemNft, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("system_nft").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

/*select count(users.token_contract_chain) from groups left join group_members on groups.group_id = group_members.group_id and group_members.group_id='201303376'
left join  users on group_members.user_id = users.user_id
where groups.group_id =201303376 and users.token_contract_chain <>""*/
