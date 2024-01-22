package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

func InputSpaceArticleID(articleID string, articleType string, reprintedID string, userID string,
	created_time time.Time, updated_at time.Time, status int8) (newSpaceID int64, err error) {
	if articleID == "" || articleType == "" || reprintedID == "" || userID == "" {
		return 0, errors.New("can't empty data")
	}
	createData := &db.SpaceArticleList{
		CreatedAt:    created_time,
		UpdatedAt:    updated_at,
		ReprintedID:  reprintedID,
		ArticleID:    articleID,
		ArticleType:  articleType,
		CreatorID:    userID,
		ArticleIsPin: 0,
		Status:       status,
	}
	err = db.DB.MysqlDB.DefaultGormDB().Table("space_article_list").Create(createData).Error
	return createData.ID, err
}

func GetSpaceArticleByID(spaceArticleID string) (*db.SpaceArticleList, error) {
	var tempData db.SpaceArticleList
	err := db.DB.MysqlDB.DefaultGormDB().Table("space_article_list").Where("id=?", spaceArticleID).First(&tempData).Error
	return &tempData, err
}
func UpdateSpaceArticleByIDGlobal(spaceArticleID string) (*db.SpaceArticleList, error) {
	var tempData db.SpaceArticleList
	err := db.DB.MysqlDB.DefaultGormDB().Table("space_article_list").
		Where("id=?", spaceArticleID).Updates(map[string]interface{}{
		"status": 1,
	}).Error
	return &tempData, err
}
func GetSpaceArticleByArticleIdAndArticleType(articleID string, articleType string) (creator string, err error) {
	var tempData []*db.SpaceArticleList
	err = db.DB.MysqlDB.DefaultGormDB().Table("space_article_list").
		Where("article_id=? and article_type=?", articleID, articleType).Find(&tempData).Error
	if len(tempData) == 0 {
		return "", errors.New("暂时没有发布该文章")
	} else {
		return tempData[0].CreatorID, nil
	}
}
func GetPersonalSpaceArticleByArticleIdAndArticleType(opUserID, articleID string, articleType string) (
	creator string, err error) {
	var tempData []*db.PersonalSpaceArticleList
	err = db.DB.MysqlDB.DefaultGormDB().Table("personal_space_article_list").
		Where("user_id=? and article_id=? and article_type=?", opUserID, articleID, articleType).Find(&tempData).Error
	if len(tempData) == 0 {
		return "", errors.New("你确定接收到这个邮件么")
	} else {
		return tempData[0].CreatorID, nil
	}
}

func DelSpaceArticleID(ID int64, userID string) (err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		//如果当前的id 是ido 是不允许删除的
		var tempData db.SpaceArticleList
		err = tx.Table("space_article_list").Where("id=?", ID).First(&tempData).Error
		//如果当前的文章是可以删除的。
		if err != nil {
			return err
		}
		if !strings.EqualFold(userID, tempData.CreatorID) {
			return errors.New("can't delete other article")
		}
		if tempData.ArticleType == "ido" {
			return errors.New("con't delete ido")
		}
		err = tx.Table("space_article_list").Where("id=?", ID).Delete(&db.SpaceArticleList{}).Error
		if err != nil {
			return err
		}
		err = tx.Table("announcement_article").Where("article_id=?", tempData.ArticleID).
			Updates(map[string]interface{}{"status": 1}).Error
		return err
	})
	return err
}

// 泛型 可以存入多个值得ÏÏ
func GetArticleList(groupID string, fromID int64, offset int64, PageSize int64) (resultData []*ArticleInfo, err error) {
	//查询所有的推文包括ido的数据
	var ArticleInfoList []*ArticleInfo
	err = db.DB.MysqlDB.DefaultGormDB().Raw(`select space_article_list.id ,space_article_list.article_type,
       space_article_list.reprint_id,space_article_list.
        article_id,space_article_list.article_is_pin, 
        space_article_list.created_at,
		space_article_list.updated_at,
		announcement_article.group_article_id,
		announcement_article.deleted_at,
		space_article_list.creator_id as creator_user_id,
		announcement_article.announcement_content,
		announcement_article.announcement_url,
		announcement_article.like_count,
		announcement_article.reword_count,
		announcement_article.is_global,
		announcement_article.order_id,
		announcement_article.status,
		announcement_article.announcement_title,
		announcement_article.announcement_summary from space_article_list
		left join announcement_article on space_article_list.article_id= announcement_article.article_id  and space_article_list.article_type='announce'
		where space_article_list.reprint_id= ? and space_article_list.deleted_at is  null order by space_article_list.article_is_pin desc ,
		                                                                                         space_article_list.updated_at desc limit ?,?`,
		groupID, offset*PageSize, PageSize).Find(&ArticleInfoList).Error
	return ArticleInfoList, nil
}
func GetSpaceArticleListBanner() (resultData []*PersonalArticleInfo, err error) {
	//查询所有的推文包括ido的数据
	var ArticleInfoList []*PersonalArticleInfo
	err = db.DB.MysqlDB.DefaultGormDB().Raw("select " +
		"space_article_list.created_at," +
		"space_article_list.updated_at," +
		"announcement_article.group_article_id," +
		"announcement_article.deleted_at," +
		"space_article_list.creator_id as creator_user_id," +
		"announcement_article.announcement_content," +
		"announcement_article.announcement_url," +
		"announcement_article.like_count," +
		"announcement_article.reword_count," +
		"announcement_article.is_global," +
		"announcement_article.order_id," +
		"announcement_article.status," +
		"announcement_article.announcement_title," +
		"announcement_article.announcement_summary," +
		"space_article_list.reprint_id," +
		"space_article_list.article_type," +
		"space_article_list.id " +
		"from space_article_list  " +
		"left join announcement_article on  " +
		"space_article_list.article_id = announcement_article.article_id " +
		"and space_article_list.article_type='announce' where is_global=1 and status=1 order by space_article_list.id desc   limit 0,10").Find(&ArticleInfoList).Error
	return ArticleInfoList, nil
}

// 泛型 可以存入多个值得ÏÏ
func GetArticleListMyEmail(selfUserId string, fromID int64, offset int64, PageSize int64) (
	resultData []*PersonalArticleInfo, err error) {
	appendQuery := ""
	if fromID > 0 {
		appendQuery = "and personal_space_article_list.id<" + utils.Int64ToString(fromID)
	}
	//查询所有的推文包括ido的数据
	var ArticleInfoList []*PersonalArticleInfo
	err = db.DB.MysqlDB.DefaultGormDB().Raw("select "+
		"personal_space_article_list.created_at,"+
		"personal_space_article_list.updated_at,"+
		"announcement_article.group_article_id,"+
		"announcement_article.deleted_at,"+
		"personal_space_article_list.creator_id as creator_user_id,"+
		"announcement_article.announcement_content,"+
		"announcement_article.announcement_url,"+
		"announcement_article.like_count,"+
		"announcement_article.reword_count,"+
		"announcement_article.is_global,"+
		"announcement_article.order_id,"+
		"announcement_article.status,"+
		"announcement_article.announcement_title,"+
		"announcement_article.announcement_summary,"+
		"announcement_article.article_id, "+
		"personal_space_article_list.reprint_id,"+
		"personal_space_article_list.article_type,"+
		"personal_space_article_list.id from personal_space_article_list  "+
		"left join announcement_article on  personal_space_article_list.article_id = announcement_article.article_id and personal_space_article_list.article_type='announce' where personal_space_article_list.user_id = ? "+appendQuery+" order by personal_space_article_list.id desc   limit 0,?", selfUserId, PageSize).Find(&ArticleInfoList).Error
	return ArticleInfoList, nil
}

func GetArticleListCount(groupID string, fromID int64, offset int64, PageSize int64) (number int64, err error) {
	//获取总数
	err = db.DB.MysqlDB.DefaultGormDB().Raw(`select count(space_article_list.id)  as number from space_article_list
		left join announcement_article on space_article_list.article_id= announcement_article.article_id  and space_article_list.article_type='announce'
		where space_article_list.reprint_id= ? and space_article_list.deleted_at is  null`, groupID).Pluck("number", &number).Error
	return number, nil
}

func GetArticleListCountMyEmail(selfUserID string, fromID int64, offset int64, PageSize int64) (number int64, err error) {
	//获取总数
	err = db.DB.MysqlDB.DefaultGormDB().Raw(`select count(id) as number from personal_space_article_list where user_id=?`, selfUserID).Pluck("number", &number).Error
	return number, nil
}
func PinSpaceArticleID(ID, userID string, isPin int32) (err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		//如果当前的id 是ido 是不允许删除的
		var tempData db.SpaceArticleList
		err = tx.Table("space_article_list").Where("id=?", ID).First(&tempData).Error
		//如果当前的文章是可以删除的。
		if err != nil {
			return err
		}
		if !strings.EqualFold(userID, tempData.CreatorID) {
			return errors.New("can't delete other article")
		}
		if isPin > 0 {
			isPin = 1

			err = tx.Table("space_article_list").Where("id=?", ID).Updates(map[string]interface{}{
				"article_is_pin": isPin,
				"updated_at":     time.Now(),
			}).Error
		} else {
			isPin = 0

			err = tx.Table("space_article_list").Where("id=?", ID).Exec("update space_article_list set article_is_pin= 0,updated_at=created_at where  id=?", ID).Error
		}

		return err
	})
	return err
}

type ArticleInfo struct {
	db.AnnouncementArticle
	ID           int64 //sapceid
	ArticleType  string
	ArticleIsPin int32
}
type PersonalArticleInfo struct {
	db.AnnouncementArticle
	ID          int64 //personalSpaceId
	ArticleType string
	ReprintId   string
}
type ArticleInfoWithUser struct {
	db.AnnouncementArticle
	PublicUserInfo       *sdk_ws.PublicUserInfo //文章运营者的信息
	PublicUserInfoWriter *sdk_ws.PublicUserInfo //文章作者信息
}

type TIdoStruct struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    IdoData `json:"data"`
}
type OutIdoStruct struct {
	ArticleID          int64  `json:"articleID"`
	ArticleType string `json:"articleType"`
	IsPin       int32  `json:"isPin"`
	IdoData
}
type IdoData struct {
	Id         string        `json:"id"`
	Write_time string        `json:"time"`
	Address    string        `json:"address"`
	Num        string        `json:"num"`
	Person     string        `json:"person"`
	Owner      string        `json:"owner"`
	Chainid    string        `json:"chainid"`
	BaseInfo   string        `json:"baseInfo"`
	Groupid    string        `json:"groupid"`
	IdoDetail  IdoDetailData `json:"data"`
	BlockTime  string        `json:"blockTime"`
}
type IdoDetailData struct {
	TokenA           string `json:"tokenA"`
	TokenB           string `json:"tokenB"`
	ProjectText      string `json:"projectText"`
	ProjectType      string `json:"projectType"`
	GroupID          string `json:"groupID"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
	InTokenCapacity  string `json:"inTokenCapacity"`
	InTokenAmount    string `json:"inTokenAmount"`
	OutTokenCapacity string `json:"outTokenCapacity"`
	MaxExchange      string `json:"maxExchange"`
	Exchange         string `json:"exchange"`
	DecimalA         string `json:"decimalA"`
	DecimalB         string `json:"decimalB"`
	LockNum          string `json:"lockNum"`
	TimeList         []int  `json:"timeList"`
	TokenNameA       string `json:"tokenNameA"`
	SymbolA          string `json:"symbolA"`
	DecimalsA        string `json:"decimalsA"`
	TotalSupplyA     string `json:"totalSupplyA"`
	TokenNameB       string `json:"tokenNameB"`
	SymbolB          string `json:"symbolB"`
	DecimalsB        string `json:"decimalsB"`
	TotalSupplyB     string `json:"totalSupplyB"`
}
