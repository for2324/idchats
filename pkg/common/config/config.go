package config

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

const ConfName = "openIMConf"

var Config config

type callBackConfig struct {
	Enable                 bool `yaml:"enable"`
	CallbackTimeOut        int  `yaml:"callbackTimeOut"`
	CallbackFailedContinue bool `yaml:"callbackFailedContinue"`
}

type config struct {
	ServerIP string `yaml:"serverip"`

	RpcRegisterIP string `yaml:"rpcRegisterIP"`
	ListenIP      string `yaml:"listenIP"`

	ServerVersion string `yaml:"serverversion"`
	Api           struct {
		GinPort  []int  `yaml:"openImApiPort"`
		ListenIP string `yaml:"listenIP"`
	}
	CmsApi struct {
		GinPort  []int  `yaml:"openImCmsApiPort"`
		ListenIP string `yaml:"listenIP"`
	}
	Sdk struct {
		WsPort  []int    `yaml:"openImSdkWsPort"`
		DataDir []string `yaml:"dataDir"`
	}
	Credential struct {
		Tencent struct {
			AppID     string `yaml:"appID"`
			Region    string `yaml:"region"`
			Bucket    string `yaml:"bucket"`
			SecretID  string `yaml:"secretID"`
			SecretKey string `yaml:"secretKey"`
		}
		Ali struct {
			RegionID           string `yaml:"regionID"`
			AccessKeyID        string `yaml:"accessKeyID"`
			AccessKeySecret    string `yaml:"accessKeySecret"`
			StsEndpoint        string `yaml:"stsEndpoint"`
			OssEndpoint        string `yaml:"ossEndpoint"`
			Bucket             string `yaml:"bucket"`
			FinalHost          string `yaml:"finalHost"`
			StsDurationSeconds int64  `yaml:"stsDurationSeconds"`
			OssRoleArn         string `yaml:"OssRoleArn"`
		}
		Minio struct {
			Bucket              string `yaml:"bucket"`
			AppBucket           string `yaml:"appBucket"`
			Location            string `yaml:"location"`
			Endpoint            string `yaml:"endpoint"`
			AccessKeyID         string `yaml:"accessKeyID"`
			SecretAccessKey     string `yaml:"secretAccessKey"`
			EndpointInner       string `yaml:"endpointInner"`
			EndpointInnerEnable bool   `yaml:"endpointInnerEnable"`
			StorageTime         int    `yaml:"storageTime"`
			IsDistributedMod    bool   `yaml:"isDistributedMod"`
		} `yaml:"minio"`
		Aws struct {
			AccessKeyID     string `yaml:"accessKeyID"`
			AccessKeySecret string `yaml:"accessKeySecret"`
			Region          string `yaml:"region"`
			Bucket          string `yaml:"bucket"`
			FinalHost       string `yaml:"finalHost"`
			RoleArn         string `yaml:"roleArn"`
			ExternalId      string `yaml:"externalId"`
			RoleSessionName string `yaml:"roleSessionName"`
		} `yaml:"aws"`
	}

	Dtm struct {
		ServerURL string `json:"serverURL"`
	}

	Mysql struct {
		DBAddress      []string `yaml:"dbMysqlAddress"`
		DBUserName     string   `yaml:"dbMysqlUserName"`
		DBPassword     string   `yaml:"dbMysqlPassword"`
		DBDatabaseName string   `yaml:"dbMysqlDatabaseName"`
		DBTableName    string   `yaml:"DBTableName"`
		DBMsgTableNum  int      `yaml:"dbMsgTableNum"`
		DBMaxOpenConns int      `yaml:"dbMaxOpenConns"`
		DBMaxIdleConns int      `yaml:"dbMaxIdleConns"`
		DBMaxLifeTime  int      `yaml:"dbMaxLifeTime"`
		LogLevel       int      `yaml:"logLevel"`
		SlowThreshold  int      `yaml:"slowThreshold"`
	}
	Mongo struct {
		DBUri                string   `yaml:"dbUri"`
		DBAddress            []string `yaml:"dbAddress"`
		DBDirect             bool     `yaml:"dbDirect"`
		DBTimeout            int      `yaml:"dbTimeout"`
		DBDatabase           string   `yaml:"dbDatabase"`
		DBSource             string   `yaml:"dbSource"`
		DBUserName           string   `yaml:"dbUserName"`
		DBPassword           string   `yaml:"dbPassword"`
		DBMaxPoolSize        int      `yaml:"dbMaxPoolSize"`
		DBRetainChatRecords  int      `yaml:"dbRetainChatRecords"`
		ChatRecordsClearTime string   `yaml:"chatRecordsClearTime"`
	}
	Redis struct {
		DBAddress     []string `yaml:"dbAddress"`
		DBMaxIdle     int      `yaml:"dbMaxIdle"`
		DBMaxActive   int      `yaml:"dbMaxActive"`
		DBIdleTimeout int      `yaml:"dbIdleTimeout"`
		DBUserName    string   `yaml:"dbUserName"`
		DBPassWord    string   `yaml:"dbPassWord"`
		EnableCluster bool     `yaml:"enableCluster"`
	}
	RpcPort struct {
		OpenImUserPort           []int `yaml:"openImUserPort"`
		OpenImFriendPort         []int `yaml:"openImFriendPort"`
		OpenImMessagePort        []int `yaml:"openImMessagePort"`
		OpenImMessageGatewayPort []int `yaml:"openImMessageGatewayPort"`
		OpenImGroupPort          []int `yaml:"openImGroupPort"`
		OpenImAuthPort           []int `yaml:"openImAuthPort"`
		OpenImPushPort           []int `yaml:"openImPushPort"`
		OpenImAdminCmsPort       []int `yaml:"openImAdminCmsPort"`
		OpenImOfficePort         []int `yaml:"openImOfficePort"`
		OpenImOrganizationPort   []int `yaml:"openImOrganizationPort"`
		OpenImConversationPort   []int `yaml:"openImConversationPort"`
		OpenImCachePort          []int `yaml:"openImCachePort"`
		OpenImRealTimeCommPort   []int `yaml:"openImRealTimeCommPort"`
		OpenImWeb3JsPort         []int `yaml:"openimweb3jsport"`
		OpenImTaskPort           []int `yaml:"openImTaskPort"`
		OpenImEnsPort            []int `yaml:"openImEnsPort"`
		OpenImOrderPort          []int `yaml:"openImOrderPort"`
		SwapRobotPort            []int `yaml:"swapRobotPort"`
	}
	RpcRegisterName struct {
		OpenImUserName   string `yaml:"openImUserName"`
		OpenImFriendName string `yaml:"openImFriendName"`
		//	OpenImOfflineMessageName     string `yaml:"openImOfflineMessageName"`
		OpenImMsgName          string `yaml:"openImMsgName"`
		OpenImPushName         string `yaml:"openImPushName"`
		OpenImRelayName        string `yaml:"openImRelayName"`
		OpenImGroupName        string `yaml:"openImGroupName"`
		OpenImAuthName         string `yaml:"openImAuthName"`
		OpenImAdminCMSName     string `yaml:"openImAdminCMSName"`
		OpenImOfficeName       string `yaml:"openImOfficeName"`
		OpenImOrganizationName string `yaml:"openImOrganizationName"`
		OpenImConversationName string `yaml:"openImConversationName"`
		OpenImCacheName        string `yaml:"openImCacheName"`
		OpenImRealTimeCommName string `yaml:"openImRealTimeCommName"`
		OpenImWeb3Js           string `yaml:"openimweb3js"`
		OpenImTask             string `yaml:"openImTaskName"`
		OpenImEns              string `yaml:"openImEnsName"`
		OpenImOrder            string `yaml:"openImOrderName"`
		UserScoreName          string `yaml:"userScoreName"`
		SwapRobotPort          string `yaml:"swapRobotPort"`
	}
	Etcd struct {
		EtcdSchema string   `yaml:"etcdSchema"`
		EtcdAddr   []string `yaml:"etcdAddr"`
		UserName   string   `yaml:"userName"`
		Password   string   `yaml:"password"`
		Secret     string   `yaml:"secret"`
	}
	Log struct {
		StorageLocation       string   `yaml:"storageLocation"`
		RotationTime          int      `yaml:"rotationTime"`
		RemainRotationCount   uint     `yaml:"remainRotationCount"`
		RemainLogLevel        uint     `yaml:"remainLogLevel"`
		ElasticSearchSwitch   bool     `yaml:"elasticSearchSwitch"`
		ElasticSearchAddr     []string `yaml:"elasticSearchAddr"`
		ElasticSearchUser     string   `yaml:"elasticSearchUser"`
		ElasticSearchPassword string   `yaml:"elasticSearchPassword"`
	}
	ModuleName struct {
		LongConnSvrName string `yaml:"longConnSvrName"`
		MsgTransferName string `yaml:"msgTransferName"`
		PushName        string `yaml:"pushName"`
	}
	Web3thirdpath struct {
		Twitterconsumerkey       string `yaml:"twitterconsumerkey"`
		Twitterconsumersecret    string `yaml:"twitterconsumersecret"`
		Twitteraccesstoken       string `yaml:"twitteraccesstoken"`
		Twitteraccesstokensecret string `yaml:"twitteraccesstokensecret"`
		TwitterBearToken         string `yaml:"twitterBearToken"`
	}
	LongConnSvr struct {
		WebsocketPort       []int `yaml:"openImWsPort"`
		WebsocketMaxConnNum int   `yaml:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int   `yaml:"websocketMaxMsgLen"`
		WebsocketTimeOut    int   `yaml:"websocketTimeOut"`
	}

	Push struct {
		Tpns struct {
			Ios struct {
				AccessID  string `yaml:"accessID"`
				SecretKey string `yaml:"secretKey"`
			}
			Android struct {
				AccessID  string `yaml:"accessID"`
				SecretKey string `yaml:"secretKey"`
			}
			Enable bool `yaml:"enable"`
		}
		Jpns struct {
			AppKey       string `yaml:"appKey"`
			MasterSecret string `yaml:"masterSecret"`
			PushUrl      string `yaml:"pushUrl"`
			PushIntent   string `yaml:"pushIntent"`
			Enable       bool   `yaml:"enable"`
		}
		Getui struct {
			PushUrl      string `yaml:"pushUrl"`
			AppKey       string `yaml:"appKey"`
			Enable       bool   `yaml:"enable"`
			Intent       string `yaml:"intent"`
			MasterSecret string `yaml:"masterSecret"`
			ChannelID    string `yaml:"channelID"`
			ChannelName  string `yaml:"channelName"`
		}
		Fcm struct {
			ServiceAccount string `yaml:"serviceAccount"`
			Enable         bool   `yaml:"enable"`
		}
		Mob struct {
			AppKey    string `yaml:"appKey"`
			PushUrl   string `yaml:"pushUrl"`
			Scheme    string `yaml:"scheme"`
			AppSecret string `yaml:"appSecret"`
			Enable    bool   `yaml:"enable"`
		}
	}
	Manager struct {
		AppManagerUid          []string `yaml:"appManagerUid"`
		Secrets                []string `yaml:"secrets"`
		AppSysNotificationName string   `yaml:"appSysNotificationName"`
	}
	InitUser struct {
		UserId []string `yaml:"userid"`
	}
	EmailSend struct {
		FromAddress   string `yaml:"from_address"`
		FromPassword  string `yaml:"from_password"`
		EmailSmtpHost string `yaml:"email_smtp_host"`
		EmailSmtpPort int    `yaml:"email_smtp_port"`
	} `yaml:"emailSend"`
	ReceiveTokenAddress string `yaml:"receiveTokenAddress"`
	ChatGptToken        string `yaml:"chatGptToken"`
	ChatGptMaxToken     int    `yaml:"chatGptMaxToken"`
	Kafka               struct {
		SASLUserName      string   `yaml:"SASLUserName"`
		SASLPassword      string   `yaml:"SASLPassword"`
		Addr              []string `yaml:"addr"`
		Partitions        int32    `yaml:"partitions"`  // 分区数
		ReplicationFactor int16    `yaml:"replication"` // 副本数
		BusinessTop       struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
			Group string   `yaml:"group"`
		} `yaml:"businessTop"`
		Ws2mschat struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		LikesAction struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		} `yaml:"likesAction"`
		AnnouncementAction struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		} `yaml:"announcementAction"`
		//Ws2mschatOffline struct {
		//	Addr  []string `yaml:"addr"`
		//	Topic string   `yaml:"topic"`
		//}
		MsgToMongo struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		Ms2pschat struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		MsgOrder struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		ConsumerGroupID struct {
			MsgToRedis        string `yaml:"msgToTransfer"`
			MsgToMongo        string `yaml:"msgToMongo"`
			MsgToMySql        string `yaml:"msgToMySql"`
			MsgToPush         string `yaml:"msgToPush"`
			MsgToLike         string `yaml:"msgToLike"`
			MsgToAnnounce     string `yaml:"msgToAnnounce"`
			MsgToOrderEns     string `yaml:"msgToOrderEns"`
			MsgToOrderArticle string `yaml:"msgToOrderArticle"`
			MsgToOrderNotify  string `yaml:"msgToOrderNotify"`
		}
	}
	Secret                            string `yaml:"secret"`
	MultiLoginPolicy                  int    `yaml:"multiloginpolicy"`
	ChatPersistenceMysql              bool   `yaml:"chatpersistencemysql"`
	ReliableStorage                   bool   `yaml:"reliablestorage"`
	MsgCacheTimeout                   int    `yaml:"msgCacheTimeout"`
	GroupMessageHasReadReceiptEnable  bool   `yaml:"groupMessageHasReadReceiptEnable"`
	SingleMessageHasReadReceiptEnable bool   `yaml:"singleMessageHasReadReceiptEnable"`
	TokenPolicy                       struct {
		AccessSecret string `yaml:"accessSecret"`
		AccessExpire int64  `yaml:"accessExpire"`
	}
	MessageVerify struct {
		FriendVerify *bool `yaml:"friendVerify"`
	}
	IOSPush struct {
		PushSound  string `yaml:"pushSound"`
		BadgeCount bool   `yaml:"badgeCount"`
		Production bool   `yaml:"production"`
	}
	Callback struct {
		CallbackUrl                        string         `yaml:"callbackUrl"`
		CallbackBeforeSendSingleMsg        callBackConfig `yaml:"callbackBeforeSendSingleMsg"`
		CallbackAfterSendSingleMsg         callBackConfig `yaml:"callbackAfterSendSingleMsg"`
		CallbackBeforeSendGroupMsg         callBackConfig `yaml:"callbackBeforeSendGroupMsg"`
		CallbackAfterSendGroupMsg          callBackConfig `yaml:"callbackAfterSendGroupMsg"`
		CallbackMsgModify                  callBackConfig `yaml:"callbackMsgModify"`
		CallbackUserOnline                 callBackConfig `yaml:"callbackUserOnline"`
		CallbackUserOffline                callBackConfig `yaml:"callbackUserOffline"`
		CallbackUserKickOff                callBackConfig `yaml:"callbackUserKickOff"`
		CallbackOfflinePush                callBackConfig `yaml:"callbackOfflinePush"`
		CallbackOnlinePush                 callBackConfig `yaml:"callbackOnlinePush"`
		CallbackBeforeSuperGroupOnlinePush callBackConfig `yaml:"callbackSuperGroupOnlinePush"`
		CallbackBeforeAddFriend            callBackConfig `yaml:"callbackBeforeAddFriend"`
		CallbackBeforeCreateGroup          callBackConfig `yaml:"callbackBeforeCreateGroup"`
	} `yaml:"callback"`
	Notification struct {
		///////////////////////group/////////////////////////////
		GroupCreated struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupCreated"`

		GroupInfoSet struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupInfoSet"`

		JoinGroupApplication struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"joinGroupApplication"`

		MemberQuit struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"memberQuit"`

		GroupApplicationAccepted struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupApplicationAccepted"`

		GroupApplicationRejected struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupApplicationRejected"`

		GroupOwnerTransferred struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupOwnerTransferred"`

		MemberKicked struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"memberKicked"`

		MemberInvited struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"memberInvited"`

		MemberEnter struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"memberEnter"`

		GroupDismissed struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupDismissed"`

		GroupMuted struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupMuted"`

		GroupCancelMuted struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupCancelMuted"`

		GroupMemberMuted struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupMemberMuted"`

		GroupMemberCancelMuted struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupMemberCancelMuted"`
		GroupMemberInfoSet struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupMemberInfoSet"`
		GroupMemberSetToAdmin struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupMemberSetToAdmin"`
		GroupMemberSetToOrdinary struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"groupMemberSetToOrdinaryUser"`
		OrganizationChanged struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"organizationChanged"`

		////////////////////////user///////////////////////
		UserInfoUpdated struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"userInfoUpdated"`

		//////////////////////friend///////////////////////
		FriendApplication struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"friendApplicationAdded"`
		FriendApplicationApproved struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"friendApplicationApproved"`

		FriendApplicationRejected struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"friendApplicationRejected"`

		FriendAdded struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"friendAdded"`
		FriendFollowApplication struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"friendFollowApplication"`
		FriendFollowDeleteApplication struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"friendFollowDeleteApplication"`

		FriendDeleted struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"friendDeleted"`
		FriendRemarkSet struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"friendRemarkSet"`
		BlackAdded struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"blackAdded"`
		BlackDeleted struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"blackDeleted"`
		ConversationOptUpdate struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"conversationOptUpdate"`
		ConversationSetPrivate struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  struct {
				OpenTips  string `yaml:"openTips"`
				CloseTips string `yaml:"closeTips"`
			} `yaml:"defaultTips"`
		} `yaml:"conversationSetPrivate"`
		WorkMomentsNotification struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"workMomentsNotification"`
		JoinDepartmentNotification struct {
			Conversation PConversation `yaml:"conversation"`
			OfflinePush  POfflinePush  `yaml:"offlinePush"`
			DefaultTips  PDefaultTips  `yaml:"defaultTips"`
		} `yaml:"joinDepartmentNotification"`
		Signal struct {
			OfflinePush struct {
				Title string `yaml:"title"`
			} `yaml:"offlinePush"`
		} `yaml:"signal"`
	}
	EnsSearch struct {
		Port     []int  `yaml:"openImEnsSearch"`
		ListenIP string `yaml:"listenIP"`
	}
	OpenNetProxy struct {
		OpenFlag bool   `yaml:"openFlag"`
		ProxyURL string `yaml:"proxyURL"`
	} `yaml:"openNetProxy"`
	IsPublicEnv bool `yaml:"ispublicenv"`

	//Contractaddress        string `yaml:"contractaddress"`
	OfficialGroupId        string `yaml:"officialGroupId"`
	OfficialSpaceId        string `yaml:"officialSpaceId"`
	UserGroupIdoBackGround string `yaml:"userGroupIdoBackGround"`
	UpdateSystemCountMine  int    `yaml:"updateSystemCountMine"`
	FreeGroupCount         int64  `yaml:"freeGroupCount"`
	OfficialTwitter        string `yaml:"officialTwitter"`

	SpaceArticle struct {
		PushUsdPrice uint64 `yaml:"pushUsdPrice"`
	} `yaml:"spaceArticle"`

	EnsPostCheck struct {
		Url string `yaml:"url"`
	} `yaml:"ensPostCheck"`
	IdoPostCheckUrl string `yaml:"idoPostCheckUrl"`
	WalletService   struct {
		Server   string `yaml:"server"`
		OpenFlag bool   `yaml:"openFlag"`
	} `yaml:"walletService"`
	Demo struct {
		Port       []int  `yaml:"openImDemoPort"`
		ListenIP   string `yaml:"listenIP"`
		Yunpiansms struct {
			Appid      string `yaml:"appid"`
			Templateid string `yaml:"templateid"`
		} `yaml:"yunpiansms"`
		AliSMSVerify struct {
			AccessKeyID                  string `yaml:"accessKeyId"`
			AccessKeySecret              string `yaml:"accessKeySecret"`
			SignName                     string `yaml:"signName"`
			VerificationCodeTemplateCode string `yaml:"verificationCodeTemplateCode"`
			Enable                       bool   `yaml:"enable"`
		}
		TencentSMS struct {
			AppID                        string `yaml:"appID"`
			Region                       string `yaml:"region"`
			SecretID                     string `yaml:"secretID"`
			SecretKey                    string `yaml:"secretKey"`
			SignName                     string `yaml:"signName"`
			VerificationCodeTemplateCode string `yaml:"verificationCodeTemplateCode"`
			Enable                       bool   `yaml:"enable"`
		}
		SuperCode    string `yaml:"superCode"`
		CodeTTL      int    `yaml:"codeTTL"`
		UseSuperCode bool   `yaml:"useSuperCode"`
		Mail         struct {
			Title                   string `yaml:"title"`
			SenderMail              string `yaml:"senderMail"`
			SenderAuthorizationCode string `yaml:"senderAuthorizationCode"`
			SmtpAddr                string `yaml:"smtpAddr"`
			SmtpPort                int    `yaml:"smtpPort"`
		}
		TestDepartMentID                        string   `yaml:"testDepartMentID"`
		ImAPIURL                                string   `yaml:"imAPIURL"`
		NeedInvitationCode                      bool     `yaml:"needInvitationCode"`
		OnboardProcess                          bool     `yaml:"onboardProcess"`
		JoinDepartmentIDList                    []string `yaml:"joinDepartmentIDList"`
		JoinDepartmentGroups                    bool     `yaml:"joinDepartmentGroups"`
		OaNotification                          bool     `yaml:"oaNotification"`
		CreateOrganizationUserAndJoinDepartment bool     `yaml:"createOrganizationUserAndJoinDepartment"`
	}
	WorkMoment struct {
		OnlyFriendCanSee bool `yaml:"onlyFriendCanSee"`
	} `yaml:"workMoment"`
	Rtc struct {
		SignalTimeout string `yaml:"signalTimeout"`
	} `yaml:"rtc"`

	Prometheus struct {
		Enable                        bool  `yaml:"enable"`
		UserPrometheusPort            []int `yaml:"userPrometheusPort"`
		FriendPrometheusPort          []int `yaml:"friendPrometheusPort"`
		MessagePrometheusPort         []int `yaml:"messagePrometheusPort"`
		MessageGatewayPrometheusPort  []int `yaml:"messageGatewayPrometheusPort"`
		GroupPrometheusPort           []int `yaml:"groupPrometheusPort"`
		AuthPrometheusPort            []int `yaml:"authPrometheusPort"`
		PushPrometheusPort            []int `yaml:"pushPrometheusPort"`
		AdminCmsPrometheusPort        []int `yaml:"adminCmsPrometheusPort"`
		OfficePrometheusPort          []int `yaml:"officePrometheusPort"`
		OrganizationPrometheusPort    []int `yaml:"organizationPrometheusPort"`
		ConversationPrometheusPort    []int `yaml:"conversationPrometheusPort"`
		CachePrometheusPort           []int `yaml:"cachePrometheusPort"`
		RealTimeCommPrometheusPort    []int `yaml:"realTimeCommPrometheusPort"`
		MessageTransferPrometheusPort []int `yaml:"messageTransferPrometheusPort"`
		Web3PrometheusPort            []int `yaml:"web3PrometheusPort"`
		TaskPrometheusPort            []int `yaml:"taskPrometheusPort"`
		EnsPrometheusPort             []int `yaml:"ensPrometheusPort"`
		OrderPrometheusPort           []int `yaml:"orderPrometheusPort"`
		SwapRobotPrometheusPort       []int `yaml:"swapRobotPrometheusPort"`
	} `yaml:"prometheus"`

	Ens struct {
		ChainId                   int64  `yaml:"chainId"`
		Contract                  string `yaml:"contract"`                  //ens 合约的地址
		Resolver                  string `yaml:"resolver"`                  //解析器的地址
		UniversalResolverContract string `yaml:"universalResolverContract"` //查询合约的地址
		EnsOwnerAddress           string `yaml:"ensOwnerAddress"`
		EnsOwnerPrivateKeyHex     string `yaml:"ensOwnerPrivateKeyHex"`
	} `yaml:"ens"`

	Pay struct {
		ScanStep        uint64                     `yaml:"scanStep"`
		ScanInterval    int64                      `yaml:"scanInterval"`
		FeeRate         uint64                     `yaml:"feeRate"`
		OrderExpireTime int64                      `yaml:"orderExpireTime"` // 分钟
		TnxTypeConfMap  map[string]TnxTypeConfInfo `yaml:"tnxTypeConfMap"`
	}
	ChainIdRpcMap  map[string][]string         `yaml:"chainIdRpcMap"`
	ChainIdHttpMap map[int64]ChainHttpEndpoint `yaml:"chainIdHttpMap"`
	UniswapRobot   struct {
		Uri              string                     `yaml:"uri"`
		BibotUri         string                     `yaml:"bibotUri"`
		ThorSwapEndpoint string                     `yaml:"thorSwapEndpoint"`
		FeeRateMap       map[string][]FeeRateConfig `yaml:"feeRateMap"`
	} `yaml:"uniswapRobot"`
	RewardTradeByScore         string `yaml:"rewardTradeByScore"`
	RewardChainRpc             string `yaml:"rewardChainRpc"` //rpc 链接， 奖励积分
	RewardChainContractAddress string `yaml:"rewardChainContractAddress"`
	RewardKey                  string `yaml:"rewardKey"`
	RewardScanBlock            int64  `yaml:"rewardScanBlock"`
	ReceiveBtcPledge           string `yaml:"receiveBtcPledge"`
	BBTPledge                  struct {
		ContractAddress    string `yaml:"contractAddress"`
		BLpContractAddress string `yaml:"BLpContractAddress"`
	} `yaml:"BBTPledge"`
	OklinkKey    string `yaml:"oklinkKey"`
	OklinkDomain string `yaml:"oklinkDomain"`
}
type FeeRateConfig struct {
	Method  string `yaml:"method"`
	FeeRate int    `yaml:"feeRate"`
}
type ChainHttpEndpoint struct {
	ApiKey   string `yaml:"apiKey"`
	EndPoint string `yaml:"endpoint"`
}
type TnxTypeConfInfo struct {
	ChainId         int64    `yaml:"chainId"`
	ChainName       string   `yaml:"chainName"`
	Decimal         uint32   `yaml:"decimal"`
	Accuracy        int      `yaml:"accuracy"`  // 精度位
	Retention       int      `yaml:"retention"` // 保留位
	ReceivedAddress []string `yaml:"receivedAddress"`
	ContractAddress string   `yaml:"contractAddress"`
}
type PConversation struct {
	ReliabilityLevel int  `yaml:"reliabilityLevel"`
	UnreadCount      bool `yaml:"unreadCount"`
}

type POfflinePush struct {
	PushSwitch bool   `yaml:"switch"`
	Title      string `yaml:"title"`
	Desc       string `yaml:"desc"`
	Ext        string `yaml:"ext"`
}
type PDefaultTips struct {
	Tips string `yaml:"tips"`
}

func unmarshalConfig(config interface{}, configName string) {
	var env string = "CONFIG_NAME"
	cfgName := os.Getenv(env)
	if len(cfgName) != 0 {
		bytes, err := ioutil.ReadFile(filepath.Join(cfgName, "config", configName))
		if err != nil {
			bytes, err = ioutil.ReadFile(filepath.Join(Root, "config", configName))
			if err != nil {
				panic(err.Error() + " config: " + filepath.Join(cfgName, "config", configName))
			}
		} else {
			Root = cfgName
		}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			panic(err.Error())
		}
	} else {
		bytes, err := ioutil.ReadFile(fmt.Sprintf("../config/%s", configName))
		if err != nil {
			bytes, err = ioutil.ReadFile(filepath.Join(Root, "config", configName))
			if err != nil {
				panic(err.Error() + " config: " + filepath.Join(cfgName, "config", configName))
			}
		}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			panic(err.Error())
		}
	}
}

func init() {
	unmarshalConfig(&Config, "config.yaml")

	//判断当前的 地方是否有.env 文件
	// 获取当前路径
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("无法获取当前路径：", err)
		return
	}
	currentDir = currentDir
	//// 拼接文件路径
	//envFilePath := currentDir + "/.env"
	//fmt.Println(envFilePath)
	//// 使用 os.Stat() 函数检查文件是否存在
	//_, err = os.Stat(envFilePath)
	//if os.IsNotExist(err) {
	//	// 用线上配置把本地配置覆盖，重新更新上去
	//	if conf, err := getEtcdConf(); err == nil {
	//		CopyMapToStruct(&Config, conf)
	//	}
	//}
	//// registerConf()
	//// 监听 etcd config 变化
	//go watchConf()
}

func watchConf() {
	key := GetPrefix(Config.Etcd.EtcdSchema, ConfName)
	etcdAddr := strings.Join(Config.Etcd.EtcdAddr, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 5 * time.Second},
	)
	if err != nil {
		panic(err)
	}
	watchChan := cli.Watch(context.Background(), key)
	for range watchChan {
		if conf, err := getEtcdConf(); err == nil {
			CopyMapToStruct(&Config, conf)
		}
	}
}

func GetRpcFromChainID(chainid string) []string {
	v, ok := Config.ChainIdRpcMap[chainid]
	if ok {
		return v
	}
	return []string{}
}
func GetRewordFromIndex(index int32) int64 {
	if index > 100 {
		return 0
	}
	//(2000000 - (index-1)*20000)
	return int64(40000 - (index-1)*400)

}
func GetSystemGroupInfo() string {
	return Config.OfficialGroupId
}

func GetConfigContractFromChainID(chainID string) string {
	if Config.IsPublicEnv {
		switch chainID {
		case "1":
			return "0xdAC17F958D2ee523a2206206994597C13D831ec7"
		case "10":
			return "0xdAC17F958D2ee523a2206206994597C13D831ec7"
		case "137":
			return "0xc2132D05D31c914a87C6611C10748AEb04B58e8F"
		case "42161":
			return "0xdAC17F958D2ee523a2206206994597C13D831ec7"
		case "56":
			return "0x55d398326f99059fF775485246999027B3197955"
		default:
			return "0xdAC17F958D2ee523a2206206994597C13D831ec7"
		}

	}
	switch chainID {
	case "1":
		return "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	case "10":
		return "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	case "137":
		return "0xc2132D05D31c914a87C6611C10748AEb04B58e8F"
	case "42161":
		return "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	case "5":
		return "0x15102793B94Bfe71b82acB884806E33b6AD5552A"
	case "56":
		return "0x55d398326f99059fF775485246999027B3197955"
	case "97":
		return "0x7E89c2b18B269864DE7caC7fCbCe64b2BF74b75D"
	case "80001":
		return "0xd989d103cc62ff24d58a127268aed1f3c99796f2"
	}
	return "0xdAC17F958D2ee523a2206206994597C13D831ec7"
}
func GetChainName(intValue uint64) string {
	switch intValue {
	case 1:
		return "eth"
	case 5:
		return "goerli"
	case 56:
		return "bsc"
	case 97:
		return "bsctest"
	case 137:
		return "matic"
	case 80001:
		return "mumbai"
	default:
		return "eth"

	}

}
func GetEnsTokenUrlServiceByChainID(chainID string) (string, string) {
	switch chainID {
	case "1":
		return "mainnet", "https://metadata.ens.domains"
	case "5":
		return "goerli", "https://metadata.ens.domains"
	case "80001":
		fallthrough
	case "137":
		return "", "https://metadata.biubiu.id/name"
	default:
		return "goerli", "https://metadata.ens.domains"
	}

}

func GetPrefix(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s/", schema, serviceName)
}

func getEtcdConf() (map[string]interface{}, error) {
	key := GetPrefix(Config.Etcd.EtcdSchema, ConfName)
	etcdAddr := strings.Join(Config.Etcd.EtcdAddr, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 5 * time.Second},
	)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := cli.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, errors.New("config not found")
	}
	var conf map[string]interface{}
	if err := yaml.Unmarshal(resp.Kvs[0].Value, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// func registerConf() {
// 	bytes, err := yaml.Marshal(Config)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	// secretMD5 := utils.Md5(config.Config.Etcd.Secret)
// 	// confBytes, err := utils.AesEncrypt(bytes, []byte(secretMD5[0:16]))
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	fmt.Println("start register", GetPrefix(Config.Etcd.EtcdSchema, ConfName))
// 	key := GetPrefix(Config.Etcd.EtcdSchema, ConfName)
// 	conf := string(bytes)
// 	etcdAddr := strings.Join(Config.Etcd.EtcdAddr, ",")
// 	cli, err := clientv3.New(clientv3.Config{
// 		Endpoints:   strings.Split(etcdAddr, ","),
// 		DialTimeout: 5 * time.Second},
// 	)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	//lease
// 	if _, err := cli.Put(context.Background(), key, conf); err != nil {
// 		fmt.Println("panic, params: ")
// 		panic(err.Error())
// 	}

// 	fmt.Println("etcd register conf ok")
// }

func CopyMapToStruct(s interface{}, m map[string]interface{}) error {
	// 创建一个与目标结构体相同类型的实例
	structType := reflect.TypeOf(s).Elem()
	structValue := reflect.ValueOf(s).Elem()
	// structValue := reflect.New(structType).Elem()

	// 将 map 转换为 JSON 字符串
	// jsonBytes, err := yaml.Marshal(m)
	// if err != nil {
	// 	return err
	// }
	// 遍历结构体的每个字段，如果对应的 map 中存在相应的键，则进行深度拷贝
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		key := field.Tag.Get("yaml")
		if key == "" {
			key = strings.ToLower(field.Name)
		}
		value, ok := m[key]
		if ok {
			fieldValue := structValue.Field(i)
			if fieldValue.CanSet() {
				// 如果字段是导出的，则进行深度拷贝
				jsonBytes, err := yaml.Marshal(value)
				if err != nil {
					return err
				}
				err = yaml.Unmarshal(jsonBytes, fieldValue.Addr().Interface())
				if err != nil {
					return err
				}
			}
		}

	}

	// 深度拷贝结构体
	// reflect.ValueOf(s).Elem().Set(structValue)
	return nil
}
