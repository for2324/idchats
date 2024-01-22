package db

import (
	"gorm.io/gorm"
	"time"
)

type Register struct {
	Account        string    `gorm:"column:account;primary_key;type:char(255)" json:"account"`
	Password       string    `gorm:"column:password;type:varchar(255)" json:"password"`
	Ex             string    `gorm:"column:ex;size:1024" json:"ex"`
	UserID         string    `gorm:"column:user_id;type:varchar(255)" json:"userID"`
	AreaCode       string    `gorm:"column:area_code;type:varchar(255)"`
	InvitationCode string    `gorm:"column:invitation_code;type:varchar(255)"`
	RegisterIP     string    `gorm:"column:register_ip;type:varchar(255)"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

type Invitation struct {
	InvitationCode string    `gorm:"column:invitation_code;primary_key;type:varchar(64)"`
	CreateTime     time.Time `gorm:"column:create_time"`
	UserID         string    `gorm:"column:user_id;index:userID"`
	LastTime       time.Time `gorm:"column:last_time"`
	Status         int32     `gorm:"column:status"`
}

// message FriendInfo{
// string OwnerUserID = 1;
// string Remark = 2;
// int64 CreateTime = 3;
// UserInfo FriendUser = 4;
// int32 AddSource = 5;
// string OperatorUserID = 6;
// string Ex = 7;
// }
// open_im_sdk.FriendInfo(FriendUser) != imdb.Friend(FriendUserID)
type Friend struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	FriendUserID   string    `gorm:"column:friend_user_id;primary_key;size:64"`
	Remark         string    `gorm:"column:remark;size:255"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}
type UserFollow struct {
	ID            uint64    `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20) unsigned;not null" json:"id"`
	FromUserID    string    `gorm:"uniqueIndex:followindex;column:from_user_id;type:varchar(64);not null" json:"fromUserID"`
	FollowUserID  string    `gorm:"uniqueIndex:followindex;column:follow_user_id;type:varchar(64);not null" json:"followUserID"`
	HandleResult  int       `gorm:"column:handle_result;type:int(11);default:null" json:"handleResult"`
	Follow        int8      `gorm:"column:follow;type:tinyint(4);default:null" json:"follow"`
	Remark        string    `gorm:"column:remark;type:varchar(255);default:null" json:"remark"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime(3);default:null" json:"createTime"`
	HandlerUserID string    `gorm:"column:handler_user_id;type:varchar(64);default:null" json:"handlerUserId"`
	HandleMsg     string    `gorm:"column:handle_msg;type:varchar(255);default:null" json:"handleMsg"`
	HandleTime    time.Time `gorm:"column:handle_time;type:datetime(3);default:null" json:"handleTime"`
	Ex            string    `gorm:"column:ex;type:varchar(1024);default:null" json:"ex"`
}

func (UserFollow) TableName() string {
	return "user_follow"
}

// UserThird 用户绑定第三方的key
type UserThird struct {
	UserId          string `gorm:"column:user_id;type:varchar(255);primary_key" json:"user_id"`
	Twitter         string `gorm:"column:twitter;type:varchar(255)" json:"twitter"`
	ShowTwitter     int32  `gorm:"column:show_twitter;type:tinyint(1);default:1" json:"showTwitter"`
	DnsDomain       string `gorm:"column:dns_domain;type:varchar(255)" json:"dnsDomain"`
	DnsDomainVerify int32  `gorm:"column:dns_domain_verify;type:tinyint(1);default:0" json:"dnsDomainVerify"`
	UserAddress     string `gorm:"column:user_address;type:varchar(255)" json:"userAddress"`
	ShowUserAddress int32  `gorm:"column:show_user_address;type:tinyint(1)" json:"showUserAddress"`
}

// UserDomain 用户绑定第三方的key
type UserDomain struct {
	UserId    string `gorm:"column:user_id;type:varchar(255);primary_key" json:"userId"`
	ChainID   string `gorm:"column:chain_id;type:varchar(255)" json:"chainId"`
	EnsDomain string `gorm:"column:ens_domain;type:varchar(255)" json:"ensDomain"`
}

func (UserDomain) TableName() string {
	return "user_domains"
}
func (UserThird) TableName() string {
	return "user_third"
}

type EventUsers struct {
	Id                  int64     `gorm:"primary_key;AUTO_INCREMENT"`
	UserId              string    `gorm:"column:user_id;type:varchar(200);index:idx_user_id;" json:"user_id" form:"user_id"`
	UserEvent           int       `gorm:"column:user_event;type:int(10)" json:"user_event" form:"user_event"`
	UserEventCreateTime time.Time `gorm:"column:user_event_create_time" json:"user_event_create_time" form:"user_event_create_time"`
	UserEventParam      string    `gorm:"column:user_event_param;type:varchar(200)" json:"user_event_param" form:"user_event_param"`
}

func (EventUsers) TableName() string {
	return "event_users"
}

type CommunityChannelRole struct {
	GroupID        string `gorm:"primaryKey;column:group_id;type:varchar(64);not null" json:"-"`
	RoleID         string `gorm:"primaryKey;index:idx_role_id;column:role_id;type:varchar(20);not null" json:"-"`
	RoleTitle      string `gorm:"index:idx_title;column:role_title;type:varchar(20);default:null" json:"roleTitle"`
	RoleIPfs       string `gorm:"column:role_ipfs;type:varchar(255);default:null" json:"roleIpfs"`
	CreatorAddress string `gorm:"column:creator_address;type:varchar(255);default:null" json:"creatorAddress"`
	Contract       string `gorm:"column:contract;type:varchar(255);default:null" json:"contract"`
	ChainID        string `gorm:"column:chain_id;type:varchar(255);default:null" json:"chainId"`
	TokenID        string `gorm:"column:token_id;type:varchar(255);default:null" json:"tokenId"`
	TokenAmount    string `gorm:"column:token_amount;type:varchar(255);default:null" json:"tokenAmount"`
	TokenSub       string `gorm:"column:token_sub;type:varchar(255);default:null" json:"tokenSub"`
	Hash           string `gorm:"unique;column:hash;type:varchar(255);default:null" json:"hash"`
}

// TableName get sql table name.获取数据库表名
func (m *CommunityChannelRole) TableName() string {
	return "community_channel_role"
}

type CommunityRoleUserRelationship struct {
	ID              int64  `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"-"`
	GroupID         string `gorm:"column:group_id;type:varchar(60);not null" json:"-"`
	Contract        string `gorm:"index:idx_role_id;column:contract;type:varchar(90);default:null" json:"contract"`
	UserID          string `gorm:"index:idx_user_id;column:user_id;type:varchar(255);default:null" json:"userId"`
	TokenID         string `gorm:"column:token_id;type:varchar(20);default:null" json:"tokenId"`
	Amount          string `gorm:"column:amount;type:varchar(20);default:null" json:"amount"`
	LastBlockNumber int64  `gorm:"column:last_block_number;type:bigint(20);default:null" json:"lastBlockNumber"`
}

// TableName get sql table name.获取数据库表名
func (m *CommunityRoleUserRelationship) TableName() string {
	return "community_role_user_relationship"
}

// CommunityChannelRoleRelationship 某个群组下 频道关联角色  用户默认属于公共频道
type CommunityChannelRoleRelationship struct {
	GroupID   string `gorm:"column:group_id;type:varchar(60);index:idx_group_id;"`
	ChannelID string `gorm:"column:channel_id;type:varchar(20);index:idx_channel_id;"`
	RoleID    string `gorm:"column:role_id;type:varchar(20);index:idx_role_id;"`
}
type UserTokenGpt struct {
	ID           int64
	UserID       string
	TokenCount   int
	ParamsString string
}

func (CommunityChannelRoleRelationship) TableName() string {
	return "community_channel_role_relationship"
}

type EventBehaviour struct {
	EventBehaviourId   int    `gorm:"column:event_behaviour_id" json:"event_behaviour_id" form:"event_behaviour_id"`
	EventBehaviourName string `gorm:"column:event_behaviour_name" json:"event_behaviour_name" form:"event_behaviour_name"`
}

func (EventBehaviour) TableName() string {
	return "event_behaviour"
}

// message FriendRequest{
// string  FromUserID = 1;
// string ToUserID = 2;
// int32 HandleResult = 3;
// string ReqMsg = 4;
// int64 CreateTime = 5;
// string HandlerUserID = 6;
// string HandleMsg = 7;
// int64 HandleTime = 8;
// string Ex = 9;
// }
// open_im_sdk.FriendRequest(nickname, farce url ...) != imdb.FriendRequest
type FriendRequest struct {
	FromUserID    string    `gorm:"column:from_user_id;primary_key;size:64"`
	ToUserID      string    `gorm:"column:to_user_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:255"`
	CreateTime    time.Time `gorm:"column:create_time"`
	HandlerUserID string    `gorm:"column:handler_user_id;size:64"`
	HandleMsg     string    `gorm:"column:handle_msg;size:255"`
	HandleTime    time.Time `gorm:"column:handle_time"`
	Ex            string    `gorm:"column:ex;size:1024"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}

//	message GroupInfo{
//		 string GroupID = 1;
//		 string GroupName = 2;
//		 string Notification = 3;
//		 string Introduction = 4;
//		 string FaceUrl = 5;
//		 string OwnerUserID = 6;
//		 uint32 MemberCount = 8;
//		 int64 CreateTime = 7;
//		 string Ex = 9;
//		 int32 Status = 10;
//		 string OpUserID = 11;
//		 int32 GroupType = 12;
//	}
//
// open_im_sdk.GroupInfo (OwnerUserID ,  MemberCount )> imdb.Group
type Group struct {
	GroupID                string    `gorm:"column:group_id;primary_key;size:64" json:"groupID" binding:"required"`
	GroupName              string    `gorm:"column:name;size:255" json:"groupName"`
	Notification           string    `gorm:"column:notification;size:255" json:"notification"`
	Introduction           string    `gorm:"column:introduction;size:255" json:"introduction"`
	FaceURL                string    `gorm:"column:face_url;size:255" json:"faceURL"`
	CreateTime             time.Time `gorm:"column:create_time;index:create_time"`
	Ex                     string    `gorm:"column:ex" json:"ex;size:1024" json:"ex"`
	Status                 int32     `gorm:"column:status"`
	CreatorUserID          string    `gorm:"column:creator_user_id;size:64"`
	GroupType              int32     `gorm:"column:group_type"`
	NeedVerification       int32     `gorm:"column:need_verification"`
	LookMemberInfo         int32     `gorm:"column:look_member_info" json:"lookMemberInfo"`
	ApplyMemberFriend      int32     `gorm:"column:apply_member_friend" json:"applyMemberFriend"`
	NotificationUpdateTime time.Time `gorm:"column:notification_update_time"`
	NotificationUserID     string    `gorm:"column:notification_user_id;size:64"`
	BlueVip                int32     `gorm:"column:blue_vip;type:int(2)" json:"blue_vip"` // 蓝v认证  0-否   1-是
	ChatTokenCount         int64     `gorm:"column:chat_token_count;type:bigint(20);default:0" json:"chatTokenCount"`
	IsFees                 int32     `gorm:"column:is_fees;type:tinyint(1);default:0" json:"isFees"` //是否是付费的群
}

// message GroupMemberFullInfo {
// string GroupID = 1 ;
// string UserID = 2 ;
// int32 roleLevel = 3;
// int64 JoinTime = 4;
// string NickName = 5;
// string FaceUrl = 6;
// int32 JoinSource = 8;
// string OperatorUserID = 9;
// string Ex = 10;
// int32 AppMangerLevel = 7; //if >0
// }  open_im_sdk.GroupMemberFullInfo(AppMangerLevel) > imdb.GroupMember
type GroupMember struct {
	GroupID        string    `gorm:"column:group_id;primary_key;size:64"`
	UserID         string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname       string    `gorm:"column:nickname;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci "`
	FaceURL        string    `gorm:"column:user_group_face_url;size:255"`
	RoleLevel      int32     `gorm:"column:role_level"`
	JoinTime       time.Time `gorm:"column:join_time"`
	JoinSource     int32     `gorm:"column:join_source"`
	InviterUserID  string    `gorm:"column:inviter_user_id;size:64"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	MuteEndTime    time.Time `gorm:"column:mute_end_time"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

// TableName get sql table name.获取数据库表名
func (m *GroupMember) TableName() string {
	return "group_members"
}

type GroupMemberData struct {
	GroupMember
	EnsDomain   string
	UserFaceURL string
}

// GroupChannel [...]
type GroupChannel struct {
	CreatedAt       time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updatedAt"`
	GroupID         string    `gorm:"primaryKey;column:group_id" json:"groupID"` // 社区id
	ChannelID       string    `gorm:"primaryKey;column:channel_id" json:"channelID"`
	ChannelName     string    `gorm:"column:channel_name" json:"channelName"`         // 频道名称
	ChannelStatus   int8      `gorm:"column:channel_status" json:"channelStatus"`     // 频道状态0为正常，1为删除
	ChannelType     string    `gorm:"column:channel_type" json:"channelType"`         // 频道类型
	ChannelDescript string    `gorm:"column:channel_descript" json:"channelDescript"` // 频道简介
	ChannelProfile  string    `gorm:"column:channel_profile" json:"channelProfile"`   // 允许什么角色进入频道
}

// TableName get sql table name.获取数据库表名
func (m *GroupChannel) TableName() string {
	return "group_channel"
}

// message GroupRequest{
// string UserID = 1;
// string GroupID = 2;
// string HandleResult = 3;
// string ReqMsg = 4;
// string  HandleMsg = 5;
// int64 ReqTime = 6;
// string HandleUserID = 7;
// int64 HandleTime = 8;
// string Ex = 9;
// }open_im_sdk.GroupRequest == imdb.GroupRequest
type GroupRequest struct {
	UserID        string    `gorm:"column:user_id;primary_key;size:64"`
	GroupID       string    `gorm:"column:group_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:1024"`
	HandledMsg    string    `gorm:"column:handle_msg;size:1024"`
	ReqTime       time.Time `gorm:"column:req_time"`
	HandleUserID  string    `gorm:"column:handle_user_id;size:64"`
	HandledTime   time.Time `gorm:"column:handle_time"`
	JoinSource    int32     `gorm:"column:join_source"`
	InviterUserID string    `gorm:"column:inviter_user_id;size:64"`
	Ex            string    `gorm:"column:ex;size:1024"`
}

// string UserID = 1;
// string Nickname = 2;
// string FaceUrl = 3;
// int32 Gender = 4;
// string PhoneNumber = 5;
// string Birth = 6;
// string Email = 7;
// string Ex = 8;
// string CreateIp = 9;
// int64 CreateTime = 10;
// int32 AppMangerLevel = 11;
// open_im_sdk.User == imdb.User
type User struct {
	UserID             string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname           string    `gorm:"column:name;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"`
	FaceURL            string    `gorm:"column:face_url;size:255"`
	Gender             int32     `gorm:"column:gender"`
	PhoneNumber        string    `gorm:"column:phone_number;size:32"`
	Birth              time.Time `gorm:"column:birth"`
	Email              string    `gorm:"column:email;size:64"`
	Ex                 string    `gorm:"column:ex;size:1024"`
	CreateTime         time.Time `gorm:"column:create_time;index:create_time"`
	AppMangerLevel     int32     `gorm:"column:app_manger_level"`
	Status             int32     `gorm:"column:status"`
	Chainid            int32     `gorm:"column:chainid;default:1" json:"chainid"`
	TokenId            string    `gorm:"column:token_id" json:"tokenId"`                                           //如果设置了官方的nft的地址 那么他的这个数值就是存在的
	UserProfile        string    `gorm:"column:user_profile" json:"userProfile"`                                   //用户的签名
	UserIntroduction   string    `gorm:"column:user_introduction" json:"userIntroduction"`                         //用户的简介
	TokenContractChain string    `gorm:"column:token_contract_chain;type:varchar(255);" json:"tokenContractChain"` //某条链条某个合约以&符号分割
	ChatTokenCount     uint64    `gorm:"column:chat_token_count;type:bigint(20);default:0" json:"chatTokenCount"`
	ChatCount          uint64    `gorm:"column:chat_count;type:bigint(20);default:0" json:"chatCount"`                //一次性消费 可以多次累加。
	GlobalMoneyCount   int64     `gorm:"column:global_money_count;type:bigint(20);default:0" json:"globalMoneyCount"` //全网广播可用次数
	GlobalRecvMsgOpt   int32     `gorm:"column:global_recv_msg_opt"`                                                  //0 正常不接收消息 1接收消息 全局推送开关
	ShowBalance        int32     `gorm:"column:show_balance;type:tinyint(1);default:1" json:"showBalance"`
	OpenAnnouncement   int32     `gorm:"column:open_announcement;type:tinyint(1);default:0" json:"openAnnouncement"`
}

// UserChatTokenRecord  用户token
type UserChatTokenRecord struct {
	ID          int64     `gorm:"primary_key;AUTO_INCREMENT;type:bigint(20)"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"createTime"`
	UserID      string    `gorm:"column:user_id;type:varchar(100);" json:"userID"`    //用户ID
	TxID        string    `gorm:"column:tx_id;type:varchar(100);unique;" json:"txID"` //交易号 唯一
	TxType      string    `gorm:"column:tx_type;type:varchar(10);" json:"txType"`     //交易类型 进账或者出账
	ParamStr    string    `gorm:"column:param_str;type:varchar(255);" json:"paramStr"`
	OldToken    uint64    `gorm:"column:old_token;type:bigint(20);default:0" json:"oldToken"`
	NewToken    uint64    `gorm:"column:new_token;type:bigint(20);default:0" json:"newToken"`
	ChainID     string    `gorm:"column:chain_id;type:varchar(10);" json:"chainID"`
	NowCount    uint64    `gorm:"column:now_count;type:bigint(20);default:0" json:"nowCount"`
}

func (UserChatTokenRecord) TableName() string {
	return "user_chat_token_record"
}

type UserIpRecord struct {
	UserID        string    `gorm:"column:user_id;primary_key;size:64"`
	CreateIp      string    `gorm:"column:create_ip;size:15"`
	LastLoginTime time.Time `gorm:"column:last_login_time"`
	LastLoginIp   string    `gorm:"column:last_login_ip;size:15"`
	LoginTimes    int32     `gorm:"column:login_times"`
}

// ip limit login
type IpLimit struct {
	Ip            string    `gorm:"column:ip;primary_key;size:15"`
	LimitRegister int32     `gorm:"column:limit_register;size:1"`
	LimitLogin    int32     `gorm:"column:limit_login;size:1"`
	CreateTime    time.Time `gorm:"column:create_time"`
	LimitTime     time.Time `gorm:"column:limit_time"`
}

// ip login
type UserIpLimit struct {
	UserID     string    `gorm:"column:user_id;primary_key;size:64"`
	Ip         string    `gorm:"column:ip;primary_key;size:15"`
	CreateTime time.Time `gorm:"column:create_time"`
}

// message BlackInfo{
// string OwnerUserID = 1;
// int64 CreateTime = 2;
// PublicUserInfo BlackUserInfo = 4;
// int32 AddSource = 5;
// string OperatorUserID = 6;
// string Ex = 7;
// }
// open_im_sdk.BlackInfo(BlackUserInfo) != imdb.Black (BlockUserID)
type Black struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	BlockUserID    string    `gorm:"column:block_user_id;primary_key;size:64"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

// AnnouncementArticleDraft 草稿箱子
type AnnouncementArticleDraft struct {
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime;default:null" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime;default:null" json:"updatedAt"`
	ArticleDraftID      int64     `gorm:"autoIncrement:true;primaryKey;column:article_draft_id;type:bigint(20);not null" json:"articleID"`
	CreatorUserID       string    `gorm:"column:creator_user_id;type:varchar(64)" json:"creatorUserID"`
	AnnouncementTitle   string    `gorm:"column:announcement_title;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci " json:"announcementTitle"`     //标题
	AnnouncementSummary string    `gorm:"column:announcement_summary;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci " json:"announcementSummary"` //摘要
	AnnouncementContent string    `gorm:"column:announcement_content;type:longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci " json:"announcementContent"`     //内容
	AnnouncementUrl     string    `gorm:"column:announcement_url;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci " json:"announcementUrl"`
	GroupID             string    `gorm:"column:group_id;type:varchar(255);not null" json:"groupID"`
	Status              int32     `gorm:"column:status;type:tinyint(1)" json:"status"`
}

func (AnnouncementArticleDraft) TableName() string {
	return "announcement_article_draft"
}

// AnnouncementArticle  已经发布的数据

type AnnouncementArticle struct {
	ArticleID           int64     `gorm:"autoIncrement:true;primaryKey;column:article_id;type:bigint(20);not null" json:"articleID"`
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime(3);default:null" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime(3);default:null" json:"updatedAt"`
	DeletedAt           time.Time `gorm:"index:idx_announcement_article_deleted_at;column:deleted_at;type:datetime(3);default:null" json:"deletedAt"`
	GroupID             string    `gorm:"uniqueIndex:group_article;column:group_id;type:varchar(60);not null" json:"groupId"`
	GroupArticleID      int64     `gorm:"uniqueIndex:group_article;column:group_article_id;type:bigint(20);not null" json:"groupArticleId"`
	CreatorUserID       string    `gorm:"column:creator_user_id;type:char(64);default:null" json:"creatorUserId"`
	AnnouncementContent string    `gorm:"column:announcement_content;type:longtext  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ;default:null" json:"announcementContent"`
	AnnouncementUrl     string    `gorm:"column:announcement_url;type:longtext  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;default:null" json:"announcementUrl"`
	LikeCount           int64     `gorm:"column:like_count;type:bigint(20);default:null;default:0" json:"likeCount"`
	RewordCount         int64     `gorm:"column:reword_count;type:bigint(20);default:null;default:0" json:"rewordCount"`
	IsGlobal            int32     `gorm:"column:is_global;type:tinyint(1);default:null;default:0" json:"isGlobal"`
	OrderID             string    `gorm:"column:order_id;type:varchar(20);default:null" json:"orderId"`
	Status              int32     `gorm:"column:status;type:tinyint(1);default:null;default:0" json:"status"`
	AnnouncementTitle   string    `gorm:"column:announcement_title;type:longtext   CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;default:null" json:"announcementTitle"`
	AnnouncementSummary string    `gorm:"column:announcement_summary;type:longtext   CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;default:null" json:"announcementSummary"`
}

func (AnnouncementArticle) TableName() string {
	return "announcement_article"
}

type AnnouncementArticleLog struct {
	ID        int64     `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"ID"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:null" json:"updatedAt"`
	ArticleID int64     `gorm:"column:article_id;type:bigint(20);not null" json:"articleID"`
	UserID    string    `gorm:"column:user_id;type:char(64);default:null" json:"userId"`
	IsLikes   int32     `gorm:"column:is_likes;type:tinyint(1);default:0" json:"isLikes"`
	GroupID   string    `gorm:"column:group_id;type:varchar(60);default:null" json:"groupId"`
	Status    int32     `gorm:"column:status;type:tinyint(1);default:0" json:"status"`
	IsGlobal  int32     `gorm:"column:is_global;type:tinyint(1);default:0" json:"isGlobal"`
}

// TableName get sql table name.获取数据库表名
func (m *AnnouncementArticleLog) TableName() string {
	return "announcement_article_logs"
}

type ChatLog struct {
	ServerMsgID      string    `gorm:"column:server_msg_id;primary_key;type:char(64)" json:"serverMsgID"`
	ClientMsgID      string    `gorm:"column:client_msg_id;type:char(64)" json:"clientMsgID"`
	SendID           string    `gorm:"column:send_id;type:char(64);index:send_id,priority:2" json:"sendID"`
	RecvID           string    `gorm:"column:recv_id;type:char(64);index:recv_id,priority:2" json:"recvID"`
	SenderPlatformID int32     `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string    `gorm:"column:sender_nick_name;type:varchar(255)  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci " json:"senderNickname"`
	SenderFaceURL    string    `gorm:"column:sender_face_url;tycpe:varchar(255);" json:"senderFaceURL"`
	SessionType      int32     `gorm:"column:session_type;index:session_type,priority:2;index:session_type_alone" json:"sessionType"`
	MsgFrom          int32     `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32     `gorm:"column:content_type;index:content_type,priority:2;index:content_type_alone" json:"contentType"`
	Content          string    `gorm:"column:content;type:varchar(3000)" json:"content"`
	Status           int32     `gorm:"column:status" json:"status"`
	SendTime         time.Time `gorm:"column:send_time;index:sendTime;index:content_type,priority:1;index:session_type,priority:1;index:recv_id,priority:1;index:send_id,priority:1" json:"sendTime"`
	CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`
	ChannelID        string    `gorm:"column:channel_id;index:index_recv_chanelid;type:char(128)" json:"channelID"`
	Ex               string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (ChatLog) TableName() string {
	return "chat_logs"
}

type BlackList struct {
	UserId           string    `gorm:"column:uid"`
	BeginDisableTime time.Time `gorm:"column:begin_disable_time"`
	EndDisableTime   time.Time `gorm:"column:end_disable_time"`
}
type Conversation struct {
	OwnerUserID           string `gorm:"column:owner_user_id;primary_key;type:char(128)" json:"OwnerUserID"`
	ConversationID        string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ConversationType      int32  `gorm:"column:conversation_type" json:"conversationType"`
	UserID                string `gorm:"column:user_id;type:char(64)" json:"userID"`
	GroupID               string `gorm:"column:group_id;type:char(128)" json:"groupID"`
	RecvMsgOpt            int32  `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
	UnreadCount           int32  `gorm:"column:unread_count" json:"unreadCount"`
	DraftTextTime         int64  `gorm:"column:draft_text_time" json:"draftTextTime"`
	IsPinned              bool   `gorm:"column:is_pinned" json:"isPinned"`
	IsPrivateChat         bool   `gorm:"column:is_private_chat" json:"isPrivateChat"`
	BurnDuration          int32  `gorm:"column:burn_duration;default:30" json:"burnDuration"`
	GroupAtType           int32  `gorm:"column:group_at_type" json:"groupAtType"`
	IsNotInGroup          bool   `gorm:"column:is_not_in_group" json:"isNotInGroup"`
	UpdateUnreadCountTime int64  `gorm:"column:update_unread_count_time" json:"updateUnreadCountTime"`
	AttachedInfo          string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex                    string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (Conversation) TableName() string {
	return "conversations"
}

type Department struct {
	DepartmentID   string    `gorm:"column:department_id;primary_key;size:64" json:"departmentID"`
	FaceURL        string    `gorm:"column:face_url;size:255" json:"faceURL"`
	Name           string    `gorm:"column:name;size:256" json:"name" binding:"required"`
	ParentID       string    `gorm:"column:parent_id;size:64" json:"parentID" binding:"required"` // "0" or Real parent id
	Order          int32     `gorm:"column:order" json:"order" `                                  // 1, 2, ...
	DepartmentType int32     `gorm:"column:department_type" json:"departmentType"`                //1, 2...
	RelatedGroupID string    `gorm:"column:related_group_id;size:64" json:"relatedGroupID"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
	Ex             string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (Department) TableName() string {
	return "departments"
}

type OrganizationUser struct {
	UserID      string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname    string    `gorm:"column:nickname;size:256"`
	EnglishName string    `gorm:"column:english_name;size:256"`
	FaceURL     string    `gorm:"column:face_url;size:256"`
	Gender      int32     `gorm:"column:gender"` //1 ,2
	Mobile      string    `gorm:"column:mobile;size:32"`
	Telephone   string    `gorm:"column:telephone;size:32"`
	Birth       time.Time `gorm:"column:birth"`
	Email       string    `gorm:"column:email;size:64"`
	CreateTime  time.Time `gorm:"column:create_time"`
	Ex          string    `gorm:"column:ex;size:1024"`
}

func (OrganizationUser) TableName() string {
	return "organization_users"
}

type DepartmentMember struct {
	UserID       string    `gorm:"column:user_id;primary_key;size:64"`
	DepartmentID string    `gorm:"column:department_id;primary_key;size:64"`
	Order        int32     `gorm:"column:order" json:"order"` //1,2
	Position     string    `gorm:"column:position;size:256" json:"position"`
	Leader       int32     `gorm:"column:leader" json:"leader"` //-1, 1
	Status       int32     `gorm:"column:status" json:"status"` //-1, 1
	CreateTime   time.Time `gorm:"column:create_time"`
	Ex           string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (DepartmentMember) TableName() string {
	return "department_members"
}

type AppVersion struct {
	Version     string `gorm:"column:version;size:64" json:"version"`
	Type        int    `gorm:"column:type;primary_key" json:"type"`
	UpdateTime  int    `gorm:"column:update_time" json:"update_time"`
	ForceUpdate bool   `gorm:"column:force_update" json:"force_update"`
	FileName    string `gorm:"column:file_name" json:"file_name"`
	YamlName    string `gorm:"column:yaml_name" json:"yaml_name"`
	UpdateLog   string `gorm:"column:update_log" json:"update_log"`
}

func (AppVersion) TableName() string {
	return "app_version"
}

type RegisterAddFriend struct {
	UserID string `gorm:"column:user_id;primary_key;size:64"`
}

func (RegisterAddFriend) TableName() string {
	return "register_add_friend"
}

type ClientInitConfig struct {
	DiscoverPageURL string `gorm:"column:discover_page_url;size:64" json:"version"`
}

func (ClientInitConfig) TableName() string {
	return "client_init_config"
}

// TwitterFollowsColumns get sql column name.获取数据库列名
var TwitterFollowsColumns = struct {
	ID              string
	TwitterUsername string
	TwitterUID      string
}{
	ID:              "id",
	TwitterUsername: "twitter_username",
	TwitterUID:      "twitter_uid",
}

// UserWhiteList [...]
type UserWhiteList struct {
	ID               int    `gorm:"primaryKey;column:id" json:"-"`
	WhiteListAddress string `gorm:"column:white_list_address" json:"whiteListAddress"`
}

// TableName get sql table name.获取数据库表名
func (m *UserWhiteList) TableName() string {
	return "user_white_list"
}

// UserWhiteListColumns get sql column name.获取数据库列名
var UserWhiteListColumns = struct {
	ID               string
	WhiteListAddress string
}{
	ID:               "id",
	WhiteListAddress: "white_list_address",
}

// EventLogs [...]
type EventLogs struct {
	ID              int64     `gorm:"primaryKey;column:id" json:"-"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"createdAt"`
	UserID          string    `gorm:"column:user_id" json:"userId"`
	EventID         string    `gorm:"column:event_id" json:"eventId"`
	EventTypename   string    `gorm:"column:event_typename" json:"eventTypename"`
	Chatwithaddress string    `gorm:"column:chatwithaddress" json:"chatwithaddress"`           // 和谁完成聊天
	UserTaskcount   int32     `gorm:"column:user_taskcount;type:int(10)" json:"userTaskcount"` // 完成度
	UserJSON        string    `gorm:"column:user_json" json:"userJson"`
	IsSync          int8      `gorm:"column:is_sync" json:"isSync"`
}

// TableName get sql table name.获取数据库表名
func (m *EventLogs) TableName() string {
	return "event_logs"
}

// EventLogsColumns get sql column name.获取数据库列名
var EventLogsColumns = struct {
	ID              string
	CreatedAt       string
	UserID          string
	EventID         string
	EventTypename   string
	Chatwithaddress string
	UserTaskcount   string
	UserJSON        string
	IsSync          string
}{
	ID:              "id",
	CreatedAt:       "created_at",
	UserID:          "user_id",
	EventID:         "event_id",
	EventTypename:   "event_typename",
	Chatwithaddress: "chatwithaddress",
	UserTaskcount:   "user_taskcount",
	UserJSON:        "user_json",
	IsSync:          "is_sync",
}

// EventTables [...]
type EventTables struct {
	EventID           int    `gorm:"column:event_id" json:"eventId"`
	EventTypename     string `gorm:"column:event_typename" json:"eventTypename"`
	EventName         string `gorm:"column:event_name" json:"eventName"`
	EventScore        int    `gorm:"column:event_score" json:"eventScore"`
	EventType         int    `gorm:"column:event_type" json:"eventType"`       // 1限制账号，2限制日期，3长期
	EventSubtype      int    `gorm:"column:event_subtype" json:"eventSubtype"` // 是否是指有一种
	Huxiangguanzhu    bool   `gorm:"column:huxiangguanzhu" json:"huxiangguanzhu"`
	Needcount         int    `gorm:"column:needcount" json:"needcount"`
	Shifoubangdingnft bool   `gorm:"column:shifoubangdingnft" json:"shifoubangdingnft"`
	IsOpen            bool   `gorm:"column:is_open" json:"isOpen"` // 是否开启
}

// TableName get sql table name.获取数据库表名
func (m *EventTables) TableName() string {
	return "event_tables"
}

// EventTablesColumns get sql column name.获取数据库列名
var EventTablesColumns = struct {
	EventID           string
	EventTypename     string
	EventName         string
	EventScore        string
	EventType         string
	EventSubtype      string
	Huxiangguanzhu    string
	Needcount         string
	Shifoubangdingnft string
	IsOpen            string
}{
	EventID:           "event_id",
	EventTypename:     "event_typename",
	EventName:         "event_name",
	EventScore:        "event_score",
	EventType:         "event_type",
	EventSubtype:      "event_subtype",
	Huxiangguanzhu:    "huxiangguanzhu",
	Needcount:         "needcount",
	Shifoubangdingnft: "shifoubangdingnft",
	IsOpen:            "is_open",
}

// AppVersionFlutter [...]
type AppVersionFlutter struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"-"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
	Platform    string    `gorm:"column:platform" json:"platform"`
	HasUpdate   int8      `gorm:"column:has_update" json:"hasUpdate"`
	IsIgnorable int8      `gorm:"column:is_ignorable" json:"isIgnorable"`
	VersionCode int       `gorm:"column:version_code" json:"versionCode"`
	VersionName string    `gorm:"column:version_name" json:"versionName"`
	UpdateLog   string    `gorm:"column:update_log" json:"updateLog"`
	UpdateLogEn string    `gorm:"column:update_log_en" json:"updateLogEn"`
	ApkURL      string    `gorm:"column:apk_url" json:"apkUrl"`
}

// TableName get sql table name.获取数据库表名
func (m *AppVersionFlutter) TableName() string {
	return "app_version_flutter"
}

// AppVersionFlutterColumns get sql column name.获取数据库列名
var AppVersionFlutterColumns = struct {
	ID          string
	Platform    string
	HasUpdate   string
	IsIgnorable string
	VersionCode string
	VersionName string
	UpdateLog   string
	ApkURL      string
	CreatedAt   string
	UpdatedAt   string
}{
	ID:          "id",
	Platform:    "platform",
	HasUpdate:   "has_update",
	IsIgnorable: "is_ignorable",
	VersionCode: "version_code",
	VersionName: "version_name",
	UpdateLog:   "update_log",
	ApkURL:      "apk_url",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// ChainToken [...]
type ChainToken struct {
	CoinChainid  int    `gorm:"primaryKey;column:coin_chainid" json:"coinChainid"`
	CoinToken    string `gorm:"primaryKey;column:coin_token" json:"coinToken"`
	CoinDecimals int    `gorm:"column:coin_decimals" json:"coinDecimals"`
	CoinName     string `gorm:"column:coin_name" json:"coinName"`
	CoinSymbol   string `gorm:"column:coin_symbol" json:"coinSymbol"`
	CoinType     string `gorm:"column:coin_type" json:"coinType"`
	CoinIsHot    int8   `gorm:"column:coin_is_hot" json:"coinIsHot"`
	CoinIcon     string `gorm:"column:coin_icon" json:"coinIcon"`
}

// TableName get sql table name.获取数据库表名
func (m *ChainToken) TableName() string {
	return "chain_token"
}

// ChainTokenColumns get sql column name.获取数据库列名
var ChainTokenColumns = struct {
	CoinChainid  string
	CoinToken    string
	CoinDecimals string
	CoinName     string
	CoinSymbol   string
	CoinType     string
	CoinIsHot    string
	CoinIcon     string
}{
	CoinChainid:  "coin_chainid",
	CoinToken:    "coin_token",
	CoinDecimals: "coin_decimals",
	CoinName:     "coin_name",
	CoinSymbol:   "coin_symbol",
	CoinType:     "coin_type",
	CoinIsHot:    "coin_is_hot",
	CoinIcon:     "coin_icon",
}

// GroupBanner [...]
type GroupBanner struct {
	BannerID    int    `gorm:"primaryKey;column:banner_id" json:"-"`   // id索引
	BannerSort  int    `gorm:"column:banner_sort" json:"bannerSort"`   // 排序
	BannerOpen  int    `gorm:"column:banner_open" json:"bannerOpen"`   // 是否打开
	BannerImage string `gorm:"column:banner_image" json:"bannerImage"` // 图片地址
	BannerUrl   string `gorm:"column:banner_url" json:"bannerUrl"`     // 跳转地址
}

// TableName get sql table name.获取数据库表名
func (m *GroupBanner) TableName() string {
	return "group_banner"
}

type EmailUserSystem struct {
	EmailAddress      string `gorm:"primaryKey;column:email_address;type:varchar(255)" json:"email_address"`
	EmailPassword     string `gorm:"column:email_password;type:varchar(255)" json:"email_password"`
	EncryptPrivateKey string `gorm:"column:encrypt_private_key;type:varchar(4096)" json:"encrypt_private_key"`
}

func (m *EmailUserSystem) TableName() string {
	return "email_user_system"
}

type GroupTagInfo struct {
	ID              int64  `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"-"`
	GroupID         string `gorm:"uniqueIndex:grouptagindex;column:group_id;type:varchar(64);not null" json:"groupId"`
	GroupTagTokenId int    `gorm:"uniqueIndex:grouptagindex;column:group_tag_tokenid;type:int(11);not null" json:"groupTagTokenId"`
	GroupTagName    string `gorm:"column:group_tag_name;type:varchar(255);default:null" json:"groupTagName"`
	GroupTagIpfsUrl string `gorm:"column:group_tag_ipfsurl;type:varchar(255);default:null" json:"groupTagIpfsUrl"`
	GroupTagCount   string `gorm:"column:group_tag_count;type:varchar(255);default:null" json:"groupTagCount"`
	Hash            string `gorm:"column:hash;type:varchar(255);default:null" json:"hash"`
	IsSync          int8   `gorm:"column:is_sync;type:tinyint(4);default:null;default:0" json:"isSync"`
}

// TableName get sql table name.获取数据库表名
func (m *GroupTagInfo) TableName() string {
	return "group_tag_info"
}

type UserNftConfig struct {
	ID              int64  `gorm:"autoIncrement:true;primaryKey;column:id;type:int(11);not null" json:"id"`
	UserID          string `gorm:"column:user_id;type:varchar(60);not null" json:"userId"`
	NftChainID      int    `gorm:"column:nft_chain_id;type:int(11);not null" json:"nftChainId"`
	NftContract     string `gorm:"column:nft_contract;type:varchar(255);not null" json:"nftContract"`
	TokenID         string `gorm:"column:token_id;type:varchar(255);not null" json:"tokenId"`
	NftContractType string `gorm:"column:nft_contract_type;type:varchar(20);not null;comment:''合约类型1155，721''" json:"nftContractType"` // '合约类型1155，721'
	NftTokenURL     string `gorm:"column:nft_token_url;type:varchar(255);not null" json:"nftTokenUrl"`
	IsShow          int    `gorm:"column:is_show;type:tinyint(1);default:null;default:1" json:"isShow"`
	Md5index        string `gorm:"unique;column:md5index;type:varchar(255);not null" json:"md5index"`
}

type UserNftConfigUserLikeLog struct {
	ID              int    `gorm:"autoIncrement:true;primaryKey;column:id;type:int(11);not null" json:"id"`
	UserID          string `gorm:"column:user_id;type:varchar(60);not null" json:"userId"`
	UserNftConfigID string `gorm:"column:user_nft_config_id;type:bigint(20);not null" json:"userNftConfigID"`
}

func (m *UserNftConfigUserLikeLog) TableName() string {
	return "user_nft_config_user_like_log"
}

// TableName get sql table name.获取数据库表名
func (m *UserNftConfig) TableName() string {
	return "user_nft_config"
}

type Task struct {
	Id              int32      `gorm:"autoIncrement:false;primary_key;column:task_id;type:int(10);not null" json:"id"` // id索引
	Name            string     `gorm:"column:name;type:varchar(255);not null" json:"name"`                             // 任务名称
	Head            string     `gorm:"column:head;type:varchar(255);not null" json:"head"`                             // 任务头像
	Type            string     `gorm:"column:type;type:varchar(20);not null" json:"type"`                              // 任务类型
	Classify        string     `gorm:"column:classify;type:varchar(20);not null" json:"classify"`                      // 任务分类
	Desc            string     `gorm:"column:desc;type:varchar(255);" json:"desc"`                                     // 任务描述
	Reward          uint64     `gorm:"column:reward;type:bigint(20);not null" json:"reward"`                           // 任务奖励
	CompletionCount int32      `gorm:"column:completionCount;type:int(10);not null" json:"completionCount"`            // 任务完成次数
	ClaimConditions string     `gorm:"column:claim_conditions;type:varchar(60);" json:"claimConditions"`
	Status          int8       `gorm:"column:status;type:tinyint(4);default:0" json:"status"`         // 0为正常，1为关闭
	CreatedAt       time.Time  `gorm:"column:created_at;type:datetime;default:null" json:"createdAt"` // 创建时间
	UpdatedAt       time.Time  `gorm:"column:updated_at;type:datetime;default:null" json:"updatedAt"` // 更新时间
	StartTime       int64      `gorm:"column:start_time;type:datetime;default:null" json:"startTime"` // 任务开始时间
	EndTime         int64      `gorm:"column:end_time;type:datetime;default:null" json:"endTime"`
	Ex              string     `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	EventType       string     `gorm:"column:event_type;type:varchar(60)" json:"eventType"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;type:datetime;default:null" json:"deletedAt"` // 删除时间
}

func (m *Task) TableName() string {
	return "task"
}

type UserTask struct {
	ID                string    `gorm:"primaryKey;column:id;varchar(255);not null" json:"id"`                     // user_task_id
	UserID            string    `gorm:"column:user_id;index:idx_user_id;type:varchar(60);not null" json:"userId"` // 用户id
	TaskID            int32     `gorm:"column:task_id;type:int(10);not null" json:"taskId"`                       // 任务id
	RewardEventLogsID string    `gorm:"column:reward_event_logs_id;type:varchar(255);" json:"rewardEventLogsId"`  // 奖励 event_logs id
	Status            int8      `gorm:"column:status;type:tinyint(4);default:0" json:"status"`                    // 0：未开始 1：进行中	2：已完成 3：已领取
	Progress          int32     `gorm:"column:progress;type:int(10);default:0" json:"progress"`                   // 任务进度
	StartTime         time.Time `gorm:"column:start_time;type:datetime;default:null" json:"startTime"`            // 任务开始时间
	EndTime           time.Time `gorm:"column:end_time;type:datetime;default:null" json:"endTime"`                // 任务结束时间
	Task              *Task     `gorm:"foreignKey:TaskID;references:Id" json:"task"`                              // 任务信息
	Ex                string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (m *UserTask) TableName() string {
	return "user_task"
}

type RewardEventLogs struct {
	ID         string    `gorm:"primaryKey;column:id;type:varchar(255);not null" json:"-"`            // id索引
	UserID     string    `gorm:"column:user_id;type:varchar(60);not null" json:"userId"`              // 用户id
	RewardType string    `gorm:"column:reward_type;type:varchar(50)" json:"rewardType"`               // 奖励类型
	Reward     uint64    `gorm:"column:reward;type:bigint(20);default:0" json:"reward"`               // 奖励数量
	UserJSON   string    `gorm:"column:user_json;type:varchar(1024)" json:"userJson"`                 // 用户信息
	IsSync     int8      `gorm:"column:is_sync;type:tinyint(4);default:null;default:0" json:"isSync"` // 是否同步 0未同步 1已同步
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;default:null" json:"createdAt"`
	SyncTime   time.Time `gorm:"column:sync_time;type:datetime;default:null" json:"syncTime"`
}

func (m *RewardEventLogs) TableName() string {
	return "reward_event_logs"
}

type UserGameScore struct {
	ID        int64          `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"-"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    string         `gorm:"column:user_id;type:varchar(60);not null" json:"userId"`
	GameID    int32          `gorm:"column:game_id;type:int(11);not null" json:"gameId"`
	GameName  string         `gorm:"column:game_name;type:varchar(255);not null" json:"gameName"`
	Score     int64          `gorm:"column:score;type:bigint(20);not null" json:"score"`
	IP        string         `gorm:"column:ip;type:varchar(255);not null" json:"ip"`
	UserAgent string         `gorm:"column:user_agent;type:text;not null" json:"userAgent"`
	StartTime int64          `gorm:"column:start_time;type:bigint(16);default:null" json:"startTime"`
	EndTime   int64          `gorm:"column:end_time;type:bigint(16);default:null" json:"endTime"`
	PlayNum   int            `gorm:"column:play_num;type:int(10);default:null;default:0" json:"playNum"`
}

// TableName get sql table name.获取数据库表名
func (m *UserGameScore) TableName() string {
	return "user_game_score"
}

type GameConfig struct {
	GameId               int32     `gorm:"column:game_id;primary_key;AUTO_INCREMENT;NOT NULL" json:"gameId"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"createdAt"`
	GameName             string    `gorm:"column:game_name" json:"gameName"`
	GameVerify           int8      `gorm:"column:game_verify;default:0" json:"gameVerify"`
	Status               int8      `gorm:"column:status;default:1" json:"status"`
	GameUrl              string    `gorm:"column:game_url" json:"gameUrl"`
	GameDesc             string    `gorm:"column:game_desc" json:"gameDesc"`
	GameCondition        string    `gorm:"column:game_condition" json:"gameCondition"` //游戏条件 json 方式
	GameCurrentPrizePool int64     `gorm:"column:game_current_prize_pool" json:"gameCurrentPrizePool"`
	GameMinPrizePool     int64     `gorm:"column:game_min_prize_pool" json:"gameMinPrizePool"`
	RewordCrontab        string    `gorm:"column:reword_crontab;type:varchar(255);'" json:"rewordCrontab"`
}
type GameConditionData struct {
	IsOfficialNft     bool `json:"isOfficialNft"`     //拥有官方NFT
	IsHeadNft         bool `json:"isHeadNft"`         // 头像设置了nft
	IsHeadOfficialNft bool `json:"isHeadOfficialNft"` //头像为官方nft头像ßß
}

func (g *GameConfig) TableName() string {
	return "game_config"
}

type GameScoreLog struct {
	Id                   int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	GameID               int       `gorm:"column:game_id;NOT NULL"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
	GameCurrentPrizePool int64     `gorm:"column:game_current_prize_pool"`
	GameAddPrizePool     int64     `gorm:"column:game_add_prize_pool"`
	JoinGameNumber       int       `gorm:"column:join_game_number"`
	OperatorIp           string    `gorm:"column:operator_ip"`
}

func (m *GameScoreLog) TableName() string {
	return "game_score_log"
}

type SystemNft struct {
	SystemNftId       int64  `gorm:"column:system_nft_id" db:"system_nft_id" json:"system_nft_id" form:"system_nft_id"`
	SystemNftChainId  int64  `gorm:"column:system_nft_chainid" db:"system_nft_chainid" json:"system_nft_chainid" form:"system_nft_chainid"`
	SystemNftContract string `gorm:"column:system_nft_contract" db:"system_nft_contract" json:"system_nft_contract" form:"system_nft_contract"`
	SystemNftClass    string `gorm:"column:system_nft_class" db:"system_nft_class" json:"system_nft_class" form:"system_nft_class"`
}

func (SystemNft) TableName() string {
	return "system_nft"
}

// 将ido 和space 合并起来
type SpaceArticleList struct {
	ID           int64      `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"ID"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:datetime;default:null" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:datetime;default:null" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"index:idx_announcement_article_deleted_at;column:deleted_at;type:datetime(3);default:null" json:"deletedAt"`
	ReprintedID  string     `gorm:"index:indexarticle;column:reprint_id;type:varchar(60);default:null" json:"reprintedID"`
	CreatorID    string     `gorm:"column:creator_id;type:varchar(60);default:null" json:"creatorId"`
	ArticleID    string     `gorm:"index:indexarticle;column:article_id;type:varchar(255);default:null" json:"articleId"`
	ArticleType  string     `gorm:"index:indexarticle;column:article_type;type:varchar(255);default:null" json:"articleType"`
	ArticleIsPin int8       `gorm:"column:article_is_pin;type:tinyint(4);default:0" json:"articleIsPin"`
	EffectEnd    time.Time  `gorm:"column:effect_end;type:datetime;default:null" json:"effectEnd"`
	IsGlobal     int8       `gorm:"column:is_global;type:tinyint(1);default:0;comment:''是否全局推送''" json:"isGlobal"`
	Status       int8       `gorm:"column:status;type:tinyint(1);default:0;comment:''是否已经支付''" json:"status"`
}

// TableName get sql table name.获取数据库表名
func (m *SpaceArticleList) TableName() string {
	return "space_article_list"
}

type PersonalSpaceArticleList struct {
	UserID string `gorm:"column:user_id;index:idx_user_id;type:varchar(60);not null" json:"userId"`
	SpaceArticleList
}

func (m *PersonalSpaceArticleList) TableName() string {
	return "personal_space_article_list"
}

type EnsAppointment struct {
	UserID    string    `gorm:"column:user_id;index:idx_user_id;type:varchar(60);not null" json:"userId"`
	EnsName   string    `gorm:"column:ens_name;unique;type:varchar(255);CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" json:"ensName"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:null" json:"createdAt"`
}

func (EnsAppointment) TableName() string {
	return "ens_appointment"
}

type EnsRegisterOrder struct {
	OrderId         uint64    `gorm:"autoIncrement:true;primaryKey;column:order_id;type:bigint(20);not null" json:"orderId"`
	EnsName         string    `gorm:"column:ens_name;type:varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null" json:"ensName"`
	Status          string    `gorm:"column:status;type:varchar(20);not null" json:"status"`
	TxnType         string    `gorm:"column:txn_type;type:varchar(20);not null" json:"txnType"`
	USDPrice        uint64    `gorm:"column:usd_price;type:bigint(20);not null" json:"usdPrice"`
	USDGasFee       uint64    `gorm:"column:usd_gas_fee;type:bigint(20);not null" json:"usdGasFee"` // 预估的 gas fee (coin)
	EnsInviter      string    `gorm:"column:ens_inviter;type:varchar(60);not null" json:"ensInviter"`
	ChainId         int64     `gorm:"column:chain_id;type:bigint(20);not null" json:"chainId"`
	TxnHash         string    `gorm:"column:txn_hash;type:varchar(255);" json:"txnHash"`
	UserId          string    `gorm:"column:txn_from_address;type:varchar(255);not null" json:"txnFromAddress"`
	RegisterTxnHash string    `gorm:"column:register_txn_hash;type:varchar(255);" json:"registerTxnHash"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;not null" json:"createTime"`
	PayTime         time.Time `gorm:"column:pay_time;type:datetime;default:null" json:"payTime"`
	ExpireTime      time.Time `gorm:"column:expire_time;type:datetime;not null" json:"expireTime"`
	Ex              string    `gorm:"column:ex;type:varchar(1024);" json:"ex"`
	// Value           uint64    `gorm:"column:value;type:bigint(20);not null" json:"value"`
	// Decimal         uint32    `gorm:"column:decimal;type:int(10);not null" json:"decimal"`
}

func (EnsRegisterOrder) TableName() string {
	return "ens_register_order"
}

type OrderPaidRecord struct {
	ID          uint64    `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"id"`
	FormAddress string    `gorm:"column:form_address;type:varchar(255);not null" json:"formAddress"`
	ToAddress   string    `gorm:"column:to_address;type:varchar(255);not null" json:"toAddress"`
	ChainId     int64     `gorm:"column:chain_id;type:bigint(20);not null" json:"chainId"`
	Value       string    `gorm:"column:value;type:varchar(255);not null" json:"value"`
	PayType     string    `gorm:"column:pay_type;type:varchar(20);" json:"payType"`
	TxnHash     string    `gorm:"column:txn_hash;type:varchar(255);not null" json:"txnHash"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;not null" json:"createTime"`
	Ex          string    `gorm:"column:ex;type:varchar(1024);" json:"ex"`
}

func (OrderPaidRecord) TableName() string {
	return "order_paid_record"
}

type PayScanBlockTask struct {
	Id                uint64    `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"id"`
	FormAddress       string    `gorm:"column:form_address;type:varchar(255);not null" json:"formAddress"`
	ToAddress         string    `gorm:"column:to_address;type:varchar(255);not null" json:"toAddress"`
	USDPrice          uint64    `gorm:"column:usd_price;type:bigint(20);not null" json:"usdPrice"`
	TxnType           string    `gorm:"column:txn_type;type:varchar(20);not null" json:"txnType"`
	Value             string    `gorm:"column:value;type:varchar(255);not null" json:"value"`
	Decimal           uint32    `gorm:"column:decimal;type:int(10);not null" json:"decimal"`
	StartBlockNumber  uint64    `gorm:"column:start_block_number;type:bigint(20);not null" json:"startBlockNumber"`
	ScanBlockNumber   uint64    `gorm:"index:idx_scan_block_number;column:scan_block_number;type:bigint(20);not null" json:"scanBlockNumber"`
	Rate              uint64    `gorm:"column:rate;type:bigint(20);not null" json:"rate"` // 1 coin = ? USDT
	Type              string    `gorm:"index:idx_type_tag_chain_id;column:type;type:varchar(20);not null" json:"type"`
	Status            int       `gorm:"column:status;type:int(10);default:0" json:"status"`
	Tag               string    `gorm:"index:idx_type_tag_chain_id;column:tag;type:varchar(50);" json:"tag"`
	ChainId           int64     `gorm:"index:idx_type_tag_chain_id;column:chain_id;type:bigint(20);not null" json:"chainId"`
	TxnHash           string    `gorm:"index:idx_txn_hash;column:txn_hash;type:varchar(255);" json:"txnHash"`
	CreateTime        time.Time `gorm:"column:create_time;type:datetime;not null" json:"createTime"`
	BlockStartTime    time.Time `gorm:"column:block_start_time;type:datetime;not null" json:"startTime"`
	BlockPayTime      time.Time `gorm:"column:block_pay_time;type:datetime;default:null" json:"payTime"`
	BlockExpireTime   time.Time `gorm:"column:block_expire_time;type:datetime;not null" json:"expireTime"`
	OrderId           string    `gorm:"column:order_id;type:varchar(255);not null" json:"orderId"`
	NotifyUrl         string    `gorm:"column:notify_url;type:varchar(255);" json:"notifyUrl"`
	Attach            string    `gorm:"column:attach;type:varchar(255);" json:"attach"`
	NotifyEncryptType string    `gorm:"column:notify_encrypt_type;type:varchar(255);" json:"notifyEncryptType"`
	NotifyEncryptKey  string    `gorm:"column:notify_encrypt_key;type:varchar(255);" json:"notifyEncryptKey"`
	Mark              string    `gorm:"column:mark;type:varchar(255);" json:"mark"`
	Ex                string    `gorm:"column:ex;type:varchar(1024);" json:"ex"`
}

func (PayScanBlockTask) TableName() string {
	return "pay_scan_block_task"
}

// UserLink 用户的链接树
type UserLink struct {
	ID          int        `gorm:"autoIncrement:true;primaryKey;column:id;type:int(11);not null" json:"-"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:datetime;default:null" json:"deletedAt"`
	UserID      string     `gorm:"column:user_id;type:varchar(60);not null" json:"userId"`
	LinkName    string     `gorm:"column:link_name;type:varchar(255)  CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci ;default:null" json:"linkName"`
	Link        string     `gorm:"column:link;type:varchar(255);" json:"link"`
	FaceURL     string     `gorm:"column:face_url;type:varchar(255)" json:"faceUrl"`
	ShowStatus  int8       `gorm:"column:show_status;type:tinyint(1);default:0" json:"showStatus"`
	DefaultIcon string     `gorm:"column:default_icon;type:varchar(255)" json:"defaultIcon"`
	Des         string     `gorm:"column:des;type:text  CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci" json:"des"`
	Bgc         string     `gorm:"column:bgc;type:varchar(255)" json:"bgc"`
	DefaultUrl  string     `gorm:"column:default_url;type:varchar(255)" json:"defaultUrl"`
	Type        string     `gorm:"column:type;type:varchar(255);default:'normal'" json:"type"`
}

// TableName get sql table name.获取数据库表名
func (m *UserLink) TableName() string {
	return "user_link"
}

type NotifyRetried struct {
	ID           int       `gorm:"autoIncrement:true;primaryKey;column:id;type:int(11);not null" json:"id"`
	NotifyUrl    string    `gorm:"column:notify_url;type:varchar(255);" json:"notifyUrl"`
	NotifyBody   []byte    `gorm:"column:notify_body;type:longblob;" json:"notifyBody"`
	RetriedCount int       `gorm:"column:retried_count;type:int(11);default:0" json:"retriedCount"`
	Mark         string    `gorm:"column:mark;type:varchar(255);" json:"mark"`
	Ex           string    `gorm:"column:ex;type:varchar(1024);" json:"ex"`
	Status       int       `gorm:"column:status;type:int(11);default:0" json:"status"`
	CreateAt     time.Time `gorm:"column:create_at;type:datetime;not null" json:"createAt"`
}

func (m *NotifyRetried) TableName() string {
	return "notify_retried"
}

type YeWuUser struct {
	User
	DnsDomainVerify int32
}

// 用户创建机器人
type Robot struct {
	UserID       string    `gorm:"column:user_id;primaryKey;type:varchar(60);not null" json:"userId"`
	EthAddress   string    `gorm:"column:eth_address;type:varchar(255);not null" json:"ethAddress"`
	BtcAddress   string    `gorm:"column:btc_address;type:varchar(255);not null" json:"btcAddress"`
	BnbAddress   string    `gorm:"column:bnb_address;type:varchar(255);not null" json:"bnbAddress"`
	TronAddress  string    `gorm:"column:tron_address;type:varchar(255);" json:"tronAddress"`
	Ex           string    `gorm:"column:ex;type:varchar(1024);" json:"ex"`
	Status       int       `gorm:"column:status;type:int(11);default:0" json:"status"`
	CreateAt     time.Time `gorm:"column:create_at;type:datetime;not null" json:"createAt"`
	FeeRate      int       `gorm:"column:fee_rate;type:int(4);default:2" json:"feeRate"`
	SnipeFeeRate int       `gorm:"column:snipe_fee_rate;type:int(4);default:10" json:"snipeFeeRate"`

	Mnemonic      string `gorm:"column:mnemonic;type:varchar(255)" json:"mnemonic"` //
	BnbPrivateKey string `gorm:"column:bnb_private_key;type:varchar(255);" json:"bnbPrivateKey"`
	EthPrivateKey string `gorm:"column:eth_private_key;type:varchar(255);" json:"ethPrivateKey"`
	BtcPrivateKey string `gorm:"column:btc_private_key;type:varchar(255);" json:"btcPrivateKey"`
}

func (m *Robot) TableName() string {
	return "uniswap_robot"
}

type RoBotTransaction struct {
	ID           uint64    `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20) unsigned;not null" json:"id"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`
	UserID       string    `gorm:"column:user_id;type:varchar(60);not null" json:"userId"`
	Address      string    `gorm:"column:address;type:varchar(255);not null;comment:''金融钱包公钥地址''" json:"address"`
	ChainID      string    `gorm:"column:chain_id;type:varchar(255);comment:''提现的网络''" json:"chainId"`
	DrawAddress  string    `gorm:"column:draw_address;type:varchar(255);not null;comment:''提到哪个地址''" json:"draw_address"`
	DrawContract string    `gorm:"column:draw_contract;type:varchar(255);comment:''提现的合约地址，空为母币''" json:"draw_contract"`
	OldAmount    string    `gorm:"column:oldAmount;type:varchar(255);not null;comment:''提现原先的金额''" json:"oldAmount"`
	Amount       string    `gorm:"column:amount;type:varchar(255);not null;comment:''提现金额''" json:"amount"`
	Hash         string    `gorm:"column:hash;type:varchar(255);not null;comment:''提现hash''" json:"hash"`
	Status       int       `gorm:"column:status;type:int(11);default:0;comment:''提现状态''" json:"status"`
	Ex           string    `gorm:"column:ex;type:varchar(1024);comment:''提现状态描述''" json:"ex"`
}

func (m *RoBotTransaction) TableName() string {
	//提现交易表格
	return "swap_robot_transaction"
}

type RoBotTask struct {
	ID          uint64    `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20) unsigned;not null" json:"id"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`
	UserID      string    `gorm:"uniqueIndex:order_id;column:user_id;type:varchar(60);not null" json:"userId"`
	Address     string    `gorm:"column:address;type:varchar(255);not null;comment:''金融钱包公钥地址''" json:"address"`
	OrdID       string    `gorm:"uniqueIndex:order_id;column:ord_id;type:varchar(255);not null" json:"ordID"`
	FromSymbol  string    `gorm:"column:from_symbol;type:varchar(255);not null" json:"fromSymbol"`
	ToSymbol    string    `gorm:"column:to_symbol;type:varchar(255);not null" json:"toSymbol"`
	Amount      string    `gorm:"column:amount;type:varchar(255);not null" json:"amount"`
	Tp          string    `gorm:"column:tp;type:varchar(255);not null" json:"tp"`
	Sl          string    `gorm:"column:sl;type:varchar(255);not null" json:"sl"`
	OrderStatus string    `gorm:"column:order_status;type:varchar(255);" json:"orderStatus"`
	MinimumOut  string    `gorm:"column:minimum_out;type:varchar(255);" json:"minimumOut"`
	DeadlineDay string    `gorm:"column:deadline_day;type:varchar(255);" json:"deadlineDay"`
	TxHash      string    `gorm:"column:tx_hash;type:varchar(255);" json:"txHash"`
}

func (m *RoBotTask) TableName() string {
	//提现交易表格
	return "swap_robot_task"
}

type RoBotTaskLog struct {
	ID          uint64    `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20) unsigned;not null" json:"id"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`
	OrdID       string    `gorm:"uniqueIndex:order_id;column:ord_id;type:varchar(255);not null" json:"ordID"`
	OrderStatus string    `gorm:"column:order_status;type:varchar(255);" json:"orderStatus"`
	Method      string    `gorm:"column:method;type:varchar(255);" json:"method"`
	TxHash      string    `gorm:"column:tx_hash;type:varchar(255);" json:"txHash"`
}

func (m *RoBotTaskLog) TableName() string {
	//提现交易表格
	return "swap_robot_task_log"
}

/******sql******
CREATE TABLE `user_robot_api` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` varchar(60) DEFAULT NULL,
  `api_key` varchar(255) DEFAULT NULL,
  `api_secrect` varchar(255) DEFAULT NULL,
  `trade_volume` varchar(255) DEFAULT NULL,
  `trade_fee` int(11) DEFAULT NULL,
  `sniper_fee` int(11) DEFAULT NULL,
  `status` int(1) DEFAULT '1',
  `api_name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id_index` (`user_id`) USING BTREE COMMENT '用户索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8
******sql******/
// UserRobotAPI [...]
type UserRobotAPI struct {
	ID          int64      `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"id"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;default:null" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp;default:null" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:timestamp;default:null" json:"deletedAt"`
	UserID      string     `gorm:"index:user_id_index;column:user_id;type:varchar(60);default:null" json:"userId"`
	APIKey      string     `gorm:"column:api_key;type:varchar(255);default:null" json:"apiKey"`
	APISecret   string     `gorm:"column:api_secret;type:varchar(1024);default:null" json:"apiSecret"`
	TradeVolume string     `gorm:"column:trade_volume;type:varchar(255);default:null" json:"tradeVolume"`
	TradeFee    string     `gorm:"column:trade_fee;type:varchar(255);default:null" json:"tradeFee"`
	SniperFee   string     `gorm:"column:sniper_fee;type:varchar(255);default:null" json:"sniperFee"`
	Status      int        `gorm:"column:status;type:int(1);default:null;default:1" json:"status"`
	APIName     string     `gorm:"column:api_name;type:varchar(255);default:null" json:"apiName"`
}

// TableName get sql table name.获取数据库表名
func (m *UserRobotAPI) TableName() string {
	return "user_robot_api"
}

// UserRobotAPIColumns get sql column name.获取数据库列名
var UserRobotAPIColumns = struct {
	ID          string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
	UserID      string
	APIKey      string
	APISecret   string
	TradeVolume string
	TradeFee    string
	SniperFee   string
	Status      string
	APIName     string
}{
	ID:          "id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	UserID:      "user_id",
	APIKey:      "api_key",
	APISecret:   "api_secret",
	TradeVolume: "trade_volume",
	TradeFee:    "trade_fee",
	SniperFee:   "sniper_fee",
	Status:      "status",
	APIName:     "api_name",
}

/******sql******
CREATE TABLE `user_history_reward` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(255) DEFAULT NULL,
  `usd_trade_volume` varchar(255) DEFAULT NULL,
  `task_id` varchar(11) DEFAULT NULL,
  `finish_time` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `task_type` varchar(255) DEFAULT NULL,
  `trade_no` varchar(255) DEFAULT NULL,
  `reward_score` varchar(255) DEFAULT NULL,
  `usd_fee_volume` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8
******sql******/
// UserHistoryReward [...]
type UserHistoryReward struct {
	ID             int64     `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"-"`
	UserID         string    `gorm:"column:user_id;type:varchar(255);default:null" json:"userId"`
	UsdTradeVolume string    `gorm:"column:usd_trade_volume;type:varchar(255);default:null" json:"usdTradeVolume"`
	TaskID         string    `gorm:"column:task_id;type:varchar(11);default:null" json:"taskId"`
	FinishTime     time.Time `gorm:"column:finish_time;type:datetime;default:null" json:"finishTime"`
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime;default:null" json:"createdAt"`
	TaskType       string    `gorm:"column:task_type;type:varchar(255);default:null" json:"taskType"`
	TradeNo        string    `gorm:"column:trade_no;type:varchar(255);default:null" json:"tradeNo"`
	RewardScore    string    `gorm:"column:reward_score;type:varchar(255);default:null" json:"rewardScore"`
	UsdFeeVolume   string    `gorm:"column:usd_fee_volume;type:varchar(255);default:null" json:"usdFeeVolume"`
}

// TableName get sql table name.获取数据库表名
func (m *UserHistoryReward) TableName() string {
	return "user_history_reward"
}

// UserHistoryRewardColumns get sql column name.获取数据库列名
var UserHistoryRewardColumns = struct {
	ID             string
	UserID         string
	UsdTradeVolume string
	TaskID         string
	FinishTime     string
	CreatedAt      string
	TaskType       string
	TradeNo        string
	RewardScore    string
	UsdFeeVolume   string
}{
	ID:             "id",
	UserID:         "user_id",
	UsdTradeVolume: "usd_trade_volume",
	TaskID:         "task_id",
	FinishTime:     "finish_time",
	CreatedAt:      "created_at",
	TaskType:       "task_type",
	TradeNo:        "trade_no",
	RewardScore:    "reward_score",
	UsdFeeVolume:   "usd_fee_volume",
}

/******sql******
CREATE TABLE `user_history_total` (
  `user_id` varchar(255) NOT NULL,
  `pending` varchar(255) DEFAULT NULL COMMENT '''待领取''',
  `claimed` varchar(255) DEFAULT NULL COMMENT '''已领取''',
  `on_chain_claimed` varchar(255) DEFAULT NULL COMMENT '''链上领取''',
  `total_trade_volume` varchar(255) DEFAULT NULL COMMENT '''个人总交易量''',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
******sql******/
// UserHistoryTotal [...]
type UserHistoryTotal struct {
	UserID              string `gorm:"primaryKey;column:user_id;type:varchar(255);not null" json:"-"`
	Pending             string `gorm:"column:pending;type:varchar(255);default:null;comment:''自己交易待领取''" json:"pending"`                             // '待领取'
	RakebackPending     string `gorm:"column:rakeback_pending;type:varchar(255);default:null;comment:''抽成待领取''" json:"rakeback_pending"`             // '待领取'
	Claimed             string `gorm:"column:claimed;type:varchar(255);default:null;comment:''已领取''" json:"claimed"`                                 // '已领取'
	RakebackClaimed     string `gorm:"column:rakeback_claimed;type:varchar(255);default:null;comment:''抽成已领取''" json:"rakebackClaimed"`              // '已领取'
	OnChainClaimed      string `gorm:"column:on_chain_claimed;type:varchar(255);default:null;comment:''链上领取''" json:"onChainClaimed"`                // '链上领取'
	TotalTradeVolume    string `gorm:"column:total_trade_volume;type:varchar(255);default:null;comment:''个人总交易量''" json:"totalTradeVolume"`          // '个人总交易量'
	SubTotalTradeVolume string `gorm:"column:sub_total_trade_volume;type:varchar(255);default:null;comment:''旗下用户总交易量''" json:"subTotalTradeVolume"` // '个人总交易量'
	CurrentNonce        string `gorm:"column:current_nonce;type:varchar(255);default:null;comment:''链上当前的nonce''" json:"currentNonce"`               // '当前链条上的最高的nonce
}

// TableName get sql table name.获取数据库表名
func (m *UserHistoryTotal) TableName() string {
	return "user_history_total"
}

// UserHistoryTotalColumns get sql column name.获取数据库列名
var UserHistoryTotalColumns = struct {
	UserID              string
	Pending             string
	RakebackPending     string
	Claimed             string
	OnChainClaimed      string
	TotalTradeVolume    string
	SubTotalTradeVolume string
	CurrentNonce        string
}{
	UserID:              "user_id",
	Pending:             "pending",
	RakebackPending:     "rakeback_pending",
	Claimed:             "claimed",
	OnChainClaimed:      "on_chain_claimed",
	TotalTradeVolume:    "total_trade_volume",
	SubTotalTradeVolume: "sub_total_trade_volume",
	CurrentNonce:        "current_nonce",
}

/******sql******
CREATE TABLE `bbt_pledge_log` (
  `id` bigint(20) NOT NULL,
  `chain` varchar(255) DEFAULT NULL,
  `created_at` date DEFAULT NULL,
  `contract` varchar(255) DEFAULT NULL,
  `total_lock` varchar(255) DEFAULT NULL,
  `block_height` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
******sql******/
/******sql******
CREATE TABLE `bbt_pledge_log` (
  `id` bigint(20) NOT NULL,
  `chain` varchar(255) NOT NULL,
  `created_at` date NOT NULL,
  `contract` varchar(255) NOT NULL,
  `total_lock` varchar(255) NOT NULL,
  `block_height` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
******sql******/
// BbtPledgeLog [...]
type BbtPledgeLog struct {
	ID          int64     `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"-"`
	Chain       string    `gorm:"column:chain;type:varchar(255);not null" json:"chain"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"createdAt"`
	BlockDate   time.Time `gorm:"column:block_date;type:date;not null" json:"blockDate"`
	Contract    string    `gorm:"column:contract;type:varchar(255);not null" json:"contract"`
	TotalLock   string    `gorm:"column:total_lock;type:varchar(255);not null" json:"totalLock"`
	BlockHeight int64     `gorm:"column:block_height;type:bigint(20);not null" json:"blockHeight"`
}

// TableName get sql table name.获取数据库表名
func (m *BbtPledgeLog) TableName() string {
	return "bbt_pledge_log"
}

// BbtPledgeLogColumns get sql column name.获取数据库列名
var BbtPledgeLogColumns = struct {
	ID          string
	Chain       string
	CreatedAt   string
	Contract    string
	TotalLock   string
	BlockHeight string
}{
	ID:          "id",
	Chain:       "chain",
	CreatedAt:   "created_at",
	Contract:    "contract",
	TotalLock:   "total_lock",
	BlockHeight: "block_height",
}

type BLPPledgeLog struct {
	ID          int64     `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"-"`
	Chain       string    `gorm:"column:chain;type:varchar(255);not null" json:"chain"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"createdAt"`
	BlockDate   time.Time `gorm:"column:block_date;type:date;not null" json:"blockDate"`
	Contract    string    `gorm:"column:contract;type:varchar(255);not null" json:"contract"`
	TotalLock   string    `gorm:"column:total_lock;type:varchar(255);not null" json:"totalLock"`
	BlockHeight int64     `gorm:"column:block_height;type:bigint(20);not null" json:"blockHeight"`
}

// TableName get sql table name.获取数据库表名
func (m *BLPPledgeLog) TableName() string {
	return "blp_pledge_log"
}

/******sql******
CREATE TABLE `obbt_pool_info` (
  `pool_id` int(11) NOT NULL AUTO_INCREMENT,
  `staking_period` varchar(255) NOT NULL COMMENT '180,360,1080',
  `pool_accobbt_pretime` datetime NOT NULL,
  `pool_last_reward_time` datetime NOT NULL,
  `rewards_rate` varchar(255) NOT NULL COMMENT '1,2,3',
  `total_stake` varchar(255) NOT NULL DEFAULT '0',
  PRIMARY KEY (`pool_id`,`staking_period`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Obbt质押的池子记录'
******sql******/
// ObbtPoolInfo Obbt质押的池子记录
type ObbtPoolInfo struct {
	PoolID             int       `gorm:"autoIncrement:true;primaryKey;column:pool_id;type:int(11);not null" json:"-"`
	StakingPeriod      string    `gorm:"primaryKey;column:staking_period;type:varchar(255);not null;comment:'180,360,1080'" json:"-"` // 180,360,1080
	PoolAccobbtPretime time.Time `gorm:"column:pool_accobbt_pretime;type:datetime;not null" json:"poolAccobbtPretime"`
	PoolLastRewardTime time.Time `gorm:"column:pool_last_reward_time;type:datetime;not null" json:"poolLastRewardTime"`
	RewardsRate        string    `gorm:"column:rewards_rate;type:varchar(255);not null;comment:'1,2,3'" json:"rewardsRate"` // 1,2,3
	TotalStake         string    `gorm:"column:total_stake;type:varchar(255);not null;default:0" json:"totalStake"`
}

// TableName get sql table name.获取数据库表名
func (m *ObbtPoolInfo) TableName() string {
	return "obbt_pool_info"
}

// ObbtPoolInfoColumns get sql column name.获取数据库列名
var ObbtPoolInfoColumns = struct {
	PoolID             string
	StakingPeriod      string
	PoolAccobbtPretime string
	PoolLastRewardTime string
	RewardsRate        string
	TotalStake         string
}{
	PoolID:             "pool_id",
	StakingPeriod:      "staking_period",
	PoolAccobbtPretime: "pool_accobbt_pretime",
	PoolLastRewardTime: "pool_last_reward_time",
	RewardsRate:        "rewards_rate",
	TotalStake:         "total_stake",
}

/******sql******
CREATE TABLE `obbt_pre_pledge` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sender_btc_address` varchar(255) NOT NULL,
  `inscription_id` varchar(255) NOT NULL,
  `amount` varchar(255) NOT NULL,
  `staking_period` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Btc用户打算用来质押的票据'
******sql******/
// ObbtPrePledge Btc用户打算用来质押的票据
type ObbtPrePledge struct {
	ID               int       `gorm:"autoIncrement:true;primaryKey;column:id;type:int(11);not null" json:"-"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	SenderBtcAddress string    `gorm:"index:sender_btc_address;column:sender_btc_address;type:varchar(255);not null" json:"senderBtcAddress"`
	InscriptionID    string    `gorm:"column:inscription_id;type:varchar(255);not null" json:"inscriptionId"`
	Amount           string    `gorm:"column:amount;type:varchar(255);not null" json:"amount"`
	IsUpToChain      uint8     `gorm:"column:is_up_to_chain;type:tinyint(1);not null;default:0" json:"isUpToChain"`
	StakingPeriod    string    `gorm:"column:staking_period;type:varchar(255);not null" json:"stakingPeriod"`
}

// TableName get sql table name.获取数据库表名
func (m *ObbtPrePledge) TableName() string {
	return "obbt_pre_pledge"
}

// ObbtPrePledgeColumns get sql column name.获取数据库列名
var ObbtPrePledgeColumns = struct {
	ID               string
	SenderBtcAddress string
	InscriptionID    string
	Amount           string
	StakingPeriod    string
	CreatedAt        string
	IsUpToChain      string
}{
	ID:               "id",
	SenderBtcAddress: "sender_btc_address",
	InscriptionID:    "inscription_id",
	Amount:           "amount",
	StakingPeriod:    "staking_period",
	CreatedAt:        "created_at",
	IsUpToChain:      "is_up_to_chain",
}

/******sql******
CREATE TABLE `obbt_stake` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `sender_btc_address` varchar(60) DEFAULT NULL,
  `pool_id` int(11) DEFAULT NULL,
  `stake_amount` varchar(255) DEFAULT NULL,
  `staking_period` varchar(255) DEFAULT NULL,
  `start_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `useridpool_id` (`sender_btc_address`,`pool_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8
******sql******/
// ObbtStake [...]
type ObbtStake struct {
	ID               int64     `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"-"`
	SenderBtcAddress string    `gorm:"uniqueIndex:useridpool_id;column:sender_btc_address;type:varchar(60);default:null" json:"senderBtcAddress"`
	PoolID           int       `gorm:"uniqueIndex:useridpool_id;column:pool_id;type:int(11);default:null" json:"poolId"`
	StakeAmount      string    `gorm:"column:stake_amount;type:varchar(255);default:null" json:"stakeAmount"`
	StakingPeriod    string    `gorm:"column:staking_period;type:varchar(255);default:null" json:"stakingPeriod"`
	StartTime        time.Time `gorm:"column:start_time;type:datetime;default:null" json:"startTime"`
	InscriptionID    string    `gorm:"column:inscription_id;type:varchar(255);not null" json:"inscriptionId"`
}

// TableName get sql table name.获取数据库表名
func (m *ObbtStake) TableName() string {
	return "obbt_stake"
}

// ObbtStakeColumns get sql column name.获取数据库列名
var ObbtStakeColumns = struct {
	ID               string
	SenderBtcAddress string
	PoolID           string
	StakeAmount      string
	StakingPeriod    string
	StartTime        string
}{
	ID:               "id",
	SenderBtcAddress: "sender_btc_address",
	PoolID:           "pool_id",
	StakeAmount:      "stake_amount",
	StakingPeriod:    "staking_period",
	StartTime:        "start_time",
}

/******sql******
CREATE TABLE `obbt_recive_from_api` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sender_btc_address` varchar(255) NOT NULL,
  `receive_btc_address` varchar(255) NOT NULL,
  `inscription_id` varchar(255) NOT NULL,
  `amount` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_up_to_chain` tinyint(1) NOT NULL DEFAULT '0',
  `blockhash` varchar(255) NOT NULL,
  `blocktime` int(10) NOT NULL,
  `height` double NOT NULL,
  `inscription_number` int(11) NOT NULL,
  `vout` int(10) NOT NULL,
  `valid` tinyint(4) NOT NULL,
  `txid` varchar(255) DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `satoshi` varchar(255) DEFAULT NULL,
  `txidx` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `sender_btc_address` (`sender_btc_address`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Btc用户打算用来质押的票据'
******sql******/
// ObbtReciveFromAPI Btc用户打算用来质押的票据
type ObbtReciveFromAPI struct {
	ID                int       `gorm:"autoIncrement:true;primaryKey;column:id;type:int(11);not null" json:"-"`
	SenderBtcAddress  string    `gorm:"index:sender_btc_address;column:sender_btc_address;type:varchar(255);not null" json:"senderBtcAddress"`
	ReceiveBtcAddress string    `gorm:"column:receive_btc_address;type:varchar(255);not null" json:"receiveBtcAddress"`
	InscriptionID     string    `gorm:"column:inscription_id;type:varchar(255);not null" json:"inscriptionId"`
	Amount            string    `gorm:"column:amount;type:varchar(255);not null" json:"amount"`
	CreatedAt         time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	IsUpToChain       bool      `gorm:"column:is_up_to_chain;type:tinyint(1);not null;default:0" json:"isUpToChain"`
	Blockhash         string    `gorm:"column:blockhash;type:varchar(255);not null" json:"blockhash"`
	Blocktime         int       `gorm:"column:blocktime;type:int(10);not null" json:"blocktime"`
	Height            int       `gorm:"column:height;type:int(11);not null" json:"height"`
	InscriptionNumber int       `gorm:"column:inscription_number;type:int(11);not null" json:"inscriptionNumber"`
	Vout              int       `gorm:"column:vout;type:int(10);not null" json:"vout"`
	Valid             int8      `gorm:"column:valid;type:tinyint(4);not null" json:"valid"`
	Txid              string    `gorm:"column:txid;type:varchar(255);default:null" json:"txid"`
	Type              string    `gorm:"column:type;type:varchar(255);default:null" json:"type"`
	Satoshi           string    `gorm:"column:satoshi;type:varchar(255);default:null" json:"satoshi"`
	Txidx             int       `gorm:"column:txidx;type:int(11);default:null" json:"txidx"`
	Ticker            string    `gorm:"column:ticker;type:varchar(255);default:null" json:"ticker"`
}

// TableName get sql table name.获取数据库表名
func (m *ObbtReciveFromAPI) TableName() string {
	return "obbt_recive_from_api"
}

/******sql******
CREATE TABLE `obbt_pool_info_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pool_id` int(11) NOT NULL,
  `ticker_name` varchar(255) NOT NULL,
  `staking_period` varchar(255) NOT NULL,
  `end_time` datetime NOT NULL DEFAULT '2999-12-31 23:59:59',
  `start_time` datetime NOT NULL,
  `reward_rate` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='记录每一种质押币的汇率更改的拉链表'
******sql******/
// ObbtPoolInfoHistory 记录每一种质押币的汇率更改的拉链表
type ObbtPoolInfoHistory struct {
	ID            int       `gorm:"autoIncrement:true;primaryKey;column:id;type:int(11);not null" json:"-"`
	PoolID        int       `gorm:"column:pool_id;type:int(11);not null" json:"poolId"`
	TickerName    string    `gorm:"column:ticker_name;type:varchar(255);not null" json:"tickerName"`
	StakingPeriod string    `gorm:"column:staking_period;type:varchar(255);not null" json:"stakingPeriod"`
	EndTime       time.Time `gorm:"column:end_time;type:datetime;not null;default:2999-12-31 23:59:59" json:"endTime"`
	StartTime     time.Time `gorm:"column:start_time;type:datetime;not null" json:"startTime"`
	RewardsRate   string    `gorm:"column:rewards_rate;type:varchar(255);default:null" json:"rewardsRate"`
}

// TableName get sql table name.获取数据库表名
func (m *ObbtPoolInfoHistory) TableName() string {
	return "obbt_pool_info_history"
}

/******sql******
CREATE TABLE `obbt_pledge_log_day_report` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `chain` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `block_date` date NOT NULL,
  `contract` varchar(255) NOT NULL,
  `total_lock` varchar(255) NOT NULL,
  `block_height` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8
******sql******/
// ObbtPledgeLogDayReport [...]
type ObbtPledgeLogDayReport struct {
	ID          int64     `gorm:"autoIncrement:true;primaryKey;column:id;type:bigint(20);not null" json:"-"`
	Chain       string    `gorm:"column:chain;type:varchar(255);not null" json:"chain"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	BlockDate   time.Time `gorm:"column:block_date;type:date;not null" json:"blockDate"`
	Contract    string    `gorm:"column:contract;type:varchar(255);not null" json:"contract"`
	TotalLock   string    `gorm:"column:total_lock;type:varchar(255);not null" json:"totalLock"`
	BlockHeight int64     `gorm:"column:block_height;type:bigint(20);not null" json:"blockHeight"`
}

// TableName get sql table name.获取数据库表名
func (m *ObbtPledgeLogDayReport) TableName() string {
	return "obbt_pledge_log_day_report"
}
