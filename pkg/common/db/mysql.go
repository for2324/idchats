package db

import (
	"Open_IM/pkg/common/config"
	"context"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mysqlDB struct {
	//sync.RWMutex
	db *gorm.DB
}

type Writer struct{}

func (w Writer) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func initMysqlDB() {
	fmt.Println("init mysqlDB start")
	//When there is no open IM database, connect to the mysql built-in database to create openIM database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=Local&parseTime=true",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], "mysql")
	var db *gorm.DB
	var err1 error
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		fmt.Println("Open failed ", err.Error(), dsn)
	}
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		db, err1 = gorm.Open(mysql.Open(dsn), nil)
		if err1 != nil {
			fmt.Println("Open failed ", err1.Error(), dsn)
			panic(err1.Error())
		}
	}

	//Check the database and table during initialization
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci;", config.Config.Mysql.DBDatabaseName)
	fmt.Println("exec sql: ", sql, " begin")
	err = db.Exec(sql).Error
	if err != nil {
		fmt.Println("Exec failed ", err.Error(), sql)
		panic(err.Error())
	}
	fmt.Println("exec sql: ", sql, " end")
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=Local&parseTime=true",
		config.Config.Mysql.DBUserName,
		config.Config.Mysql.DBPassword,
		config.Config.Mysql.DBAddress[0],
		config.Config.Mysql.DBDatabaseName)
	fmt.Println("db log  init===>", dsn)
	newLogger := logger.New(
		Writer{},
		logger.Config{
			SlowThreshold:             time.Duration(config.Config.Mysql.SlowThreshold) * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.LogLevel(config.Config.Mysql.LogLevel),                       // Log level
			IgnoreRecordNotFoundError: true,                                                                // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                                                // Disable color
		},
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	fmt.Println("db log  init", dsn)
	if err != nil {
		fmt.Println("Open failed ", err.Error(), dsn)
		panic(err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}

	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Config.Mysql.DBMaxLifeTime))
	sqlDB.SetMaxOpenConns(config.Config.Mysql.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Config.Mysql.DBMaxIdleConns)
	fmt.Println(" sql db init middle")
	//InitBeginProcessStart(db, ChangeSpaceArticleTableColumnsName)
	fmt.Println(" sql db init sync")
	db.AutoMigrate(
		&Register{},
		&Friend{},
		&FriendRequest{},
		&Group{},
		&GroupMember{},
		&GroupRequest{},
		&User{},
		&Black{}, &ChatLog{}, &Register{}, &Conversation{}, &AppVersion{}, &Department{}, &BlackList{}, &IpLimit{}, &UserIpLimit{}, &Invitation{}, &RegisterAddFriend{},
		&ClientInitConfig{}, &UserIpRecord{},
		&UserFollow{}, &UserThird{},
		&UserDomain{}, &EventLogs{}, &EventTables{},
		&UserWhiteList{}, &AppVersionFlutter{},
		&ChainToken{}, &GroupBanner{}, &EventUsers{},
		&CommunityChannelRoleRelationship{},
		&CommunityRoleUserRelationship{},
		&CommunityChannelRole{},
		&EmailUserSystem{}, &UserChatTokenRecord{}, &UserNftConfig{}, &UserNftConfigUserLikeLog{},
		&AnnouncementArticle{},
		&AnnouncementArticleLog{},
		&AnnouncementArticleDraft{}, &UserGameScore{},
		&RewardEventLogs{},
		&Task{},
		&UserTask{},
		&GameConfig{},
		&GameScoreLog{},
		&SystemNft{},
		&SpaceArticleList{},
		&GroupTagInfo{},
		&EnsRegisterOrder{},
		&OrderPaidRecord{},
		&PayScanBlockTask{},
		&UserLink{},
		&PersonalSpaceArticleList{},
		&NotifyRetried{},
		&Robot{},
		&RoBotTransaction{},
		&RoBotTask{},
		&UserRobotAPI{},
		&UserHistoryReward{},
		&UserHistoryTotal{},
		&BbtPledgeLog{},
		&BLPPledgeLog{},
		&ObbtPoolInfo{},
		&ObbtPrePledge{},
		&ObbtStake{},
		&ObbtReciveFromAPI{},
		&ObbtPledgeLogDayReport{},
		&ObbtPoolInfoHistory{},
	)
	db.Set("gorm:table_options", "CHARSET=utf8")
	db.Set("gorm:table_options", "collation=utf8_unicode_ci")

	if !db.Migrator().HasTable(&Friend{}) {
		fmt.Println("CreateTable Friend")
		db.Migrator().CreateTable(&Friend{})
	}
	if !db.Migrator().HasTable(&UserFollow{}) {
		fmt.Println("CreateTable UserFollow")
		db.Migrator().CreateTable(&UserFollow{})
	}

	if !db.Migrator().HasTable(&FriendRequest{}) {
		fmt.Println("CreateTable FriendRequest")
		db.Migrator().CreateTable(&FriendRequest{})
	}

	if !db.Migrator().HasTable(&Group{}) {
		fmt.Println("CreateTable Group")
		db.Migrator().CreateTable(&Group{})
	}

	if !db.Migrator().HasTable(&GroupMember{}) {
		fmt.Println("CreateTable GroupMember")
		db.Migrator().CreateTable(&GroupMember{})
	}
	if !db.Migrator().HasTable(&GroupRequest{}) {
		fmt.Println("CreateTable GroupRequest")
		db.Migrator().CreateTable(&GroupRequest{})
	}
	if !db.Migrator().HasTable(&User{}) {
		fmt.Println("CreateTable User")
		db.Migrator().CreateTable(&User{})
	}
	if !db.Migrator().HasTable(&Black{}) {
		fmt.Println("CreateTable Black")
		db.Migrator().CreateTable(&Black{})
	}
	if !db.Migrator().HasTable(&ChatLog{}) {
		fmt.Println("CreateTable ChatLog")
		db.Migrator().CreateTable(&ChatLog{})
	}
	if !db.Migrator().HasTable(&Register{}) {
		fmt.Println("CreateTable Register")
		db.Migrator().CreateTable(&Register{})
	}
	if !db.Migrator().HasTable(&Conversation{}) {
		fmt.Println("CreateTable Conversation")
		db.Migrator().CreateTable(&Conversation{})
	}

	if !db.Migrator().HasTable(&Department{}) {
		fmt.Println("CreateTable Department")
		db.Migrator().CreateTable(&Department{})
	}
	if !db.Migrator().HasTable(&OrganizationUser{}) {
		fmt.Println("CreateTable OrganizationUser")
		db.Migrator().CreateTable(&OrganizationUser{})
	}
	if !db.Migrator().HasTable(&DepartmentMember{}) {
		fmt.Println("CreateTable DepartmentMember")
		db.Migrator().CreateTable(&DepartmentMember{})
	}
	if !db.Migrator().HasTable(&AppVersion{}) {
		fmt.Println("CreateTable DepartmentMember")
		db.Migrator().CreateTable(&AppVersion{})
	}
	if !db.Migrator().HasTable(&BlackList{}) {
		fmt.Println("CreateTable BlackList")
		db.Migrator().CreateTable(&BlackList{})
	}
	if !db.Migrator().HasTable(&IpLimit{}) {
		fmt.Println("CreateTable IpLimit")
		db.Migrator().CreateTable(&IpLimit{})
	}
	if !db.Migrator().HasTable(&UserIpLimit{}) {
		fmt.Println("CreateTable UserIpLimit")
		db.Migrator().CreateTable(&UserIpLimit{})
	}

	if !db.Migrator().HasTable(&RegisterAddFriend{}) {
		fmt.Println("CreateTable RegisterAddFriend")
		db.Migrator().CreateTable(&RegisterAddFriend{})
	}
	if !db.Migrator().HasTable(&Invitation{}) {
		fmt.Println("CreateTable Invitation")
		db.Migrator().CreateTable(&Invitation{})
	}
	if !db.Migrator().HasTable(&UserDomain{}) {
		fmt.Println("CreateTable UserDomain")
		db.Migrator().CreateTable(&UserDomain{})
	}
	if !db.Migrator().HasTable(&EventLogs{}) {
		fmt.Println("CreateTable EventLogs")
		db.Migrator().CreateTable(&EventLogs{})
	}
	if !db.Migrator().HasTable(&EventTables{}) {
		fmt.Println("CreateTable EventTables")
		db.Migrator().CreateTable(&EventTables{})
	}
	if !db.Migrator().HasTable(&ClientInitConfig{}) {
		fmt.Println("CreateTable ClientInitConfig")
		db.Migrator().CreateTable(&ClientInitConfig{})
	}

	if !db.Migrator().HasTable(&UserIpRecord{}) {
		fmt.Println("CreateTable Friend")
		db.Migrator().CreateTable(&UserIpRecord{})
	}
	if !db.Migrator().HasTable(&UserThird{}) {
		fmt.Println("CreateTable UserThird")
		db.Migrator().CreateTable(&UserThird{})
	}
	if !db.Migrator().HasTable(&UserWhiteList{}) {
		fmt.Println("CreateTable UserWhiteList")
		db.Migrator().CreateTable(&UserWhiteList{})
	}
	if !db.Migrator().HasTable(&AppVersionFlutter{}) {
		fmt.Println("CreateTable AppVersionFlutter")
		db.Migrator().CreateTable(&AppVersionFlutter{})
	}
	if !db.Migrator().HasTable(&ChainToken{}) {
		fmt.Println("CreateTable ChainToken")
		db.Migrator().CreateTable(&ChainToken{})
	}
	if !db.Migrator().HasTable(&GroupBanner{}) {
		fmt.Println("CreateTable GroupBanner")
		db.Migrator().CreateTable(&GroupBanner{})
	}
	if !db.Migrator().HasTable(&GroupChannel{}) {
		fmt.Println("CreateTable GroupChannel")
		db.Migrator().CreateTable(&GroupChannel{})
	}
	if !db.Migrator().HasTable(&EventUsers{}) {
		fmt.Println("CreateTable EventUsers")
		db.Migrator().CreateTable(&EventUsers{})
	}
	if !db.Migrator().HasTable(&EventBehaviour{}) {
		fmt.Println("CreateTable EventBehaviour")
		db.Migrator().CreateTable(&EventBehaviour{})
	}
	if !db.Migrator().HasTable(&GroupTagInfo{}) {
		fmt.Println("CreateTable GroupTagInfo")
		db.Migrator().CreateTable(&GroupTagInfo{})
	}
	if !db.Migrator().HasTable(&RewardEventLogs{}) {
		fmt.Println("CreateTable RewardEventLogs")
		db.Migrator().CreateTable(&RewardEventLogs{})
	}
	if !db.Migrator().HasTable(&Task{}) {
		fmt.Println("CreateTable Task")
		db.Migrator().CreateTable(&Task{})
	}
	if !db.Migrator().HasTable(&UserTask{}) {
		fmt.Println("CreateTable UserTask")
		err := db.Migrator().CreateTable(&UserTask{})
		fmt.Println("CreateTable UserTask after", err)
	}

	if !db.Migrator().HasTable(&OrderPaidRecord{}) {
		fmt.Println("CreateTable OrderPaidRecord")
		err := db.Migrator().CreateTable(&OrderPaidRecord{})
		fmt.Println("CreateTable OrderPaidRecord after", err)
	}
	if !db.Migrator().HasTable(&EnsRegisterOrder{}) {
		fmt.Println("CreateTable EnsRegisterOrder")
		err := db.Migrator().CreateTable(&EnsRegisterOrder{})
		fmt.Println("CreateTable EnsRegisterOrder after", err)
	}
	if !db.Migrator().HasTable(&PayScanBlockTask{}) {
		fmt.Println("CreateTable PayScanBlockTask")
		err := db.Migrator().CreateTable(&PayScanBlockTask{})
		fmt.Println("CreateTable PayScanBlockTask after", err)
	}
	DB.MysqlDB.db = db
	return
}

func (m *mysqlDB) DefaultGormDB() *gorm.DB {
	return DB.MysqlDB.db
}

type Callback func(*gorm.DB)

func InitBeginProcessStart(dbTemp *gorm.DB, callBack Callback) {
	// 初始化Etcd客户端
	etcdAddr := strings.Join(config.Config.Etcd.EtcdAddr, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 5 * time.Second,
	})
	defer cli.Close()

	// 创建分布式锁
	lockKey := "/db_init_lock"
	leaseResp, err := cli.Grant(context.Background(), 10)
	resp, err := cli.Txn(context.Background()).
		If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseResp.ID))).
		Commit()
	if err == nil && resp.Succeeded {
		if callBack != nil {
			callBack(dbTemp)
			fmt.Println("初始化数据库完毕")
		}
		cli.Put(context.Background(), lockKey, "1")
	}
}
func ChangeSpaceArticleTableColumnsName(dbTemp *gorm.DB) {
	hasTable := dbTemp.Migrator().HasTable(&SpaceArticleList{})
	if hasTable {
		hasColumn := dbTemp.Migrator().HasColumn(&SpaceArticleList{}, "group_id")
		if hasColumn {
			dbTemp.Table("space_article_list").Exec("ALTER TABLE space_article_list  CHANGE  group_id  reprint_id varchar(60)")
			dbTemp.Table("space_article_list").Exec("drop index indexarticle on space_article_list ")
		}
	}
}
