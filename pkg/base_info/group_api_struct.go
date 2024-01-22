package base_info

import (
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
)

type CommResp struct {
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type CommDataResp struct {
	CommResp
	Data []map[string]interface{} `json:"data"`
}

type KickGroupMemberReq struct {
	GroupID          string   `json:"groupID" binding:"required"`
	KickedUserIDList []string `json:"kickedUserIDList" binding:"required"`
	Reason           string   `json:"reason"`
	OperationID      string   `json:"operationID" binding:"required"`
}
type KickGroupMemberResp struct {
	CommResp
	UserIDResultList []*UserIDResult `json:"data"`
}

type GetGroupMembersInfoReq struct {
	GroupID     string   `json:"groupID" binding:"required"`
	MemberList  []string `json:"memberList" binding:"required"`
	OperationID string   `json:"operationID" binding:"required"`
}
type GetGroupMembersInfoResp struct {
	CommResp
	MemberList []*open_im_sdk.GroupMemberFullInfo `json:"-"`
	Data       []map[string]interface{}           `json:"data" swaggerignore:"true"`
}

type InviteUserToGroupReq struct {
	GroupID           string   `json:"groupID" binding:"required"`
	InvitedUserIDList []string `json:"invitedUserIDList" binding:"required"`
	Reason            string   `json:"reason"`
	OperationID       string   `json:"operationID" binding:"required"`
}
type InviteUserToGroupResp struct {
	CommResp
	UserIDResultList []*UserIDResult `json:"data"`
}

type GetJoinedGroupListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetJoinedGroupListResp struct {
	CommResp
	GroupInfoList []*open_im_sdk.GroupInfo `json:"-"`
	Data          []map[string]interface{} `json:"data" swaggerignore:"true"`
}

type GetGroupMemberListReq struct {
	GroupID     string `json:"groupID"`
	Filter      int32  `json:"filter"`
	NextSeq     int32  `json:"nextSeq"`
	OperationID string `json:"operationID"`
}
type GetGroupMemberListResp struct {
	CommResp
	NextSeq    int32                              `json:"nextSeq"`
	MemberList []*open_im_sdk.GroupMemberFullInfo `json:"-"`
	Data       []map[string]interface{}           `json:"data" swaggerignore:"true"`
}

type GetGroupAllMemberReq struct {
	GroupID     string `json:"groupID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
	Offset      int32  `json:"offset"`
	Count       int32  `json:"count"`
}
type GetGroupAllMemberResp struct {
	CommResp
	MemberList []*open_im_sdk.GroupMemberFullInfo `json:"-"`
	Data       []map[string]interface{}           `json:"data" swaggerignore:"true"`
}

//
//type GetGroupAllMemberListBySplitReq struct {
//	GroupID     string `json:"groupID" binding:"required"`
//	OperationID string `json:"operationID" binding:"required"`
//	Offset      int32  `json:"offset" binding:"required"`
//	Count       int32  `json:"count" binding:"required"`
//}
//type GetGroupAllMemberListBySplitResp struct {
//	CommResp
//	MemberList []*open_im_sdk.GroupMemberFullInfo `json:"-"`
//	Data       []map[string]interface{}           `json:"data" swaggerignore:"true"`
//}

type CreateGroupReq struct {
	MemberList   []*GroupAddMemberInfo `json:"memberList"`
	OwnerUserID  string                `json:"ownerUserID"`
	GroupType    int32                 `json:"groupType"`
	GroupName    string                `json:"groupName"`
	Notification string                `json:"notification"`
	Introduction string                `json:"introduction"`
	FaceURL      string                `json:"faceURL"`
	Ex           string                `json:"ex"`
	OperationID  string                `json:"operationID" binding:"required"`
	GroupID      string                `json:"groupID"`
}
type CreateGroupResp struct {
	CommResp
	GroupInfo open_im_sdk.GroupInfo  `json:"-"`
	Data      map[string]interface{} `json:"data" swaggerignore:"true"`
}

type GetGroupApplicationListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"` //作为管理员或群主收到的 进群申请
}
type GetGroupApplicationListResp struct {
	CommResp
	GroupRequestList []*open_im_sdk.GroupRequest `json:"-"`
	Data             []map[string]interface{}    `json:"data" swaggerignore:"true"`
}

type GetUserReqGroupApplicationListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}

type GetUserRespGroupApplicationResp struct {
	CommResp
	GroupRequestList []*open_im_sdk.GroupRequest `json:"-"`
}

type GetGroupInfoReq struct {
	GroupIDList []string `json:"groupIDList" binding:"required"`
	OperationID string   `json:"operationID" binding:"required"`
}
type GetGroupInfoResp struct {
	CommResp
	GroupInfoList []*open_im_sdk.GroupInfo `json:"-"`
	Data          []map[string]interface{} `json:"data" swaggerignore:"true"`
}

//type GroupInfoAlias struct {
//	open_im_sdk.GroupInfo
//	NeedVerification int32 `protobuf:"bytes,13,opt,name=needVerification" json:"needVerification,omitempty"`
//}

//type GroupInfoAlias struct {
//	GroupID          string `protobuf:"bytes,1,opt,name=groupID" json:"groupID,omitempty"`
//	GroupName        string `protobuf:"bytes,2,opt,name=groupName" json:"groupName,omitempty"`
//	Notification     string `protobuf:"bytes,3,opt,name=notification" json:"notification,omitempty"`
//	Introduction     string `protobuf:"bytes,4,opt,name=introduction" json:"introduction,omitempty"`
//	FaceURL          string `protobuf:"bytes,5,opt,name=faceURL" json:"faceURL,omitempty"`
//	OwnerUserID      string `protobuf:"bytes,6,opt,name=ownerUserID" json:"ownerUserID,omitempty"`
//	CreateTime       uint32 `protobuf:"varint,7,opt,name=createTime" json:"createTime,omitempty"`
//	MemberCount      uint32 `protobuf:"varint,8,opt,name=memberCount" json:"memberCount,omitempty"`
//	Ex               string `protobuf:"bytes,9,opt,name=ex" json:"ex,omitempty"`
//	Status           int32  `protobuf:"varint,10,opt,name=status" json:"status,omitempty"`
//	OpUserID    string `protobuf:"bytes,11,opt,name=creatorUserID" json:"creatorUserID,omitempty"`
//	GroupType        int32  `protobuf:"varint,12,opt,name=groupType" json:"groupType,omitempty"`
//	NeedVerification int32  `protobuf:"bytes,13,opt,name=needVerification" json:"needVerification,omitempty"`
//}

type ApplicationGroupResponseReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	GroupID      string `json:"groupID" binding:"required"`
	FromUserID   string `json:"fromUserID" binding:"required"` //application from FromUserID
	HandledMsg   string `json:"handledMsg"`
	HandleResult int32  `json:"handleResult" binding:"required,oneof=-1 1"`
}
type ApplicationGroupResponseResp struct {
	CommResp
}

type JoinGroupReq struct {
	GroupID       string `json:"groupID" binding:"required"`
	ReqMessage    string `json:"reqMessage"`
	OperationID   string `json:"operationID" binding:"required"`
	JoinSource    int32  `json:"joinSource"`
	InviterUserID string `json:"inviterUserID"`
}

type JoinGroupResp struct {
	CommResp
}

type QuitGroupReq struct {
	GroupID     string `json:"groupID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type QuitGroupResp struct {
	CommResp
}

type SetGroupInfoReq struct {
	GroupID           string `json:"groupID" binding:"required"`
	GroupName         string `json:"groupName"`
	Notification      string `json:"notification"`
	Introduction      string `json:"introduction"`
	FaceURL           string `json:"faceURL"`
	Ex                string `json:"ex"`
	OperationID       string `json:"operationID" binding:"required"`
	NeedVerification  *int32 `json:"needVerification"`
	LookMemberInfo    *int32 `json:"lookMemberInfo"`
	ApplyMemberFriend *int32 `json:"applyMemberFriend"`
}

type SetGroupInfoResp struct {
	CommResp
}

type TransferGroupOwnerReq struct {
	GroupID        string `json:"groupID" binding:"required"`
	OldOwnerUserID string `json:"oldOwnerUserID" binding:"required"`
	NewOwnerUserID string `json:"newOwnerUserID" binding:"required"`
	OperationID    string `json:"operationID" binding:"required"`
}
type TransferGroupOwnerResp struct {
	CommResp
}

type DismissGroupReq struct {
	GroupID     string `json:"groupID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type DismissGroupResp struct {
	CommResp
}

type MuteGroupMemberReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	GroupID      string `json:"groupID" binding:"required"`
	UserID       string `json:"userID" binding:"required"`
	MutedSeconds uint32 `json:"mutedSeconds" binding:"required"`
}
type MuteGroupMemberResp struct {
	CommResp
}

type CancelMuteGroupMemberReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}
type CancelMuteGroupMemberResp struct {
	CommResp
}

type MuteGroupReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
}
type MuteGroupResp struct {
	CommResp
}

type CancelMuteGroupReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
}
type CancelMuteGroupResp struct {
	CommResp
}

type SetGroupMemberNicknameReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
	Nickname    string `json:"nickname"`
}

type SetGroupMemberNicknameResp struct {
	CommResp
}

type SetGroupMemberInfoReq struct {
	OperationID string  `json:"operationID" binding:"required"`
	GroupID     string  `json:"groupID" binding:"required"`
	UserID      string  `json:"userID" binding:"required"`
	Nickname    *string `json:"nickname"`
	FaceURL     *string `json:"userGroupFaceUrl"`
	RoleLevel   *int32  `json:"roleLevel" validate:"gte=1,lte=3"`
	Ex          *string `json:"ex"`
}

type SetGroupMemberInfoResp struct {
	CommResp
}

type GetGroupAbstractInfoReq struct {
	OperationID string `json:"operationID"`
	GroupID     string `json:"groupID"`
}

type GetGroupAbstractInfoResp struct {
	CommResp
	GroupMemberNumber   int32  `json:"groupMemberNumber"`
	GroupMemberListHash uint64 `json:"groupMemberListHash"`
}

type CreateSysUserGroupReq struct {
	OperationID string `json:"operationID"`
}
type CreateSysUserResp struct {
	CommResp
}
type BannerImage struct {
	BannerImage string `json:"bannerImage"`
	BannerSort  int32  `json:"bannerSort"`
	BannerUrl   string `json:"bannerUrl"`
}
type GroupInfoHotMessage struct {
	// 社区id
	GroupID string `json:"groupID"`
	// 社区名字
	GroupName string `json:"groupName"`
	// 社区公告
	Notification string `json:"notification"`
	// 社区介绍
	Introduction string `json:"introduction"`
	// 群图标
	FaceURL string `json:"faceURL"`
	// 群创建者
	OwnerUserID string `json:"ownerUserID"`
	//创建者的faceUrl
	CreatorFaceURL string `json:"creatorFaceURL"`
	//群总共人数
	MemberCount int64 `json:"memberCount"`
	//群总共人数
	FollowCount int64 `json:"followCount"`
	// 创建时间
	CreateTime    uint32                             `json:"createTime"`
	CreatorUserID string                             `json:"creatorUserID"`
	GroupType     uint32                             `json:"groupType"`
	IsJoinEd      bool                               `json:"isJoinEd"`
	MemberList    []*open_im_sdk.GroupMemberFullInfo `json:"memberList"`
}

type CreateCommunityReq struct {
	OwnerUserID  string `json:"ownerUserID"`
	GroupType    int32  `json:"groupType"`
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"` //官网的网址
	Introduction string `json:"introduction"`
	FaceURL      string `json:"faceURL"`
	Ex           string `json:"ex"`
	OperationID  string `json:"operationID" binding:"required"`
	GroupID      string `json:"groupID"`
	IsFees       int32  `json:"isFees"`
	FeeTx        string `json:"feeTx"`
}
type CreateCommunityResp struct {
	CommResp
	GroupInfo open_im_sdk.GroupInfo  `json:"-"`
	Data      map[string]interface{} `json:"data" swaggerignore:"true"`
}
type UpdateCommunityChannelReq struct {
	OwnerUserID string `json:"ownerUserID"`
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"` //群id
	OpInfo      string `json:"opInfo" binding:"required"`  //添加频道Add 或者删除频道Del 修改名称Update
	ChannelId   string `json:"channelId"`                  //添加的时候不需要修改 0和1 的ID不可以删除 其他可以
	ChannelName string `json:"channelName"`
}
type UpdateCommunityChannelResp struct {
	CommResp
	GroupChannelInfo open_im_sdk.GroupChannelInfo `json:"-"`
	Data             map[string]interface{}       `json:"data" swaggerignore:"true"`
}
type GetHotCommunityReq struct {
	OperationID string `json:"operationID" binding:"required"`
	SearchTitle string `json:"searchTitle"`
}
type GetHotSpaceReq struct {
	OperationID string `json:"operationID" binding:"required"`
	PageIndex   int32  `json:"pageIndex"`
	PageSize    int32  `json:"pageSize"`
	SearchType  string `json:"searchType" default:"all" commit:"all:全部,follow:关注的"`
}

type ApiSpaceUserInfo struct {
	open_im_sdk.UserInfo
	FollowCount int32                  `json:"followCount"`
	Group       *open_im_sdk.GroupInfo `json:"group"`
}

type GetHotSpaceResp struct {
	CommResp
	Data []*ApiSpaceUserInfo `json:"data"`
}

type GetHotSpaceMainInfoResp struct {
	CommResp
	BannerArrayImage []*BannerImage `json:"bannerArrayImage"`
}
type SearchCommunityResp struct {
	CommResp
	GroupInfoArray []*GroupInfoHotMessage `json:"groupInfoArray"`
}
type CommunityChannelAllListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
}
type CommunityChannelAllListResp struct {
	CommResp
	GroupChannelInfoList []*open_im_sdk.GroupChannelInfo `json:"-"`
	Data                 []map[string]interface{}        `json:"data"`
}
type CommunityChannelStatusReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
	ChannelID   string `json:"channelID" binding:"required"`
}
type CommunityChannelStatusResp struct {
	CommResp
	Data int32 `json:"data"`
}
type GetHistoryMessageListParams struct {
	UserID           string `json:"userID"`
	GroupID          string `json:"groupID"`
	ConversationID   string `json:"conversationID"`
	StartClientMsgID string `json:"startClientMsgID"`
	Count            int    `json:"count"`
	ChannelID        string `json:"channelID"`
	IsReverse        bool   `json:"isReverse"`
}

type GetSingleChatHistoryMessageListReq struct {
	UserID           string `json:"userID"`
	ConversationID   string `json:"conversationID"`
	StartClientMsgID string `json:"startClientMsgID"`
	Count            int    `json:"count"`
	ChannelID        string `json:"channelID"`
	IsReverse        bool   `json:"isReverse"`
}

type MsgStruct struct {
	ClientMsgID      string       `json:"clientMsgID,omitempty"`
	ServerMsgID      string       `json:"serverMsgID,omitempty"`
	CreateTime       int64        `json:"createTime"`
	SendTime         int64        `json:"sendTime"`
	SessionType      int32        `json:"sessionType"`
	SendID           string       `json:"sendID,omitempty"`
	RecvID           string       `json:"recvID,omitempty"`
	MsgFrom          int32        `json:"msgFrom"`
	ContentType      int32        `json:"contentType"`
	SenderPlatformID int32        `json:"platformID"`
	SenderNickname   string       `json:"senderNickname,omitempty"`
	SenderFaceURL    string       `json:"senderFaceUrl,omitempty"`
	GroupID          string       `json:"groupID,omitempty"`
	ChannelID        string       `json:"channelID,omitempty"`
	Content          string       `json:"content,omitempty"`
	Seq              uint32       `json:"seq"`
	IsRead           bool         `json:"isRead"`
	Status           int32        `json:"status"`
	AttachedInfo     string       `json:"attachedInfo,omitempty"`
	Ex               string       `json:"ex,omitempty"`
	PictureElem      *PictureElem `json:"pictureElem,omitempty"`
	SoundElem        *SoundElem   `json:"soundElem,omitempty"`
	VideoElem        *VideoElem   `json:"videoElem,omitempty"`
	FileElem         *FileElem    `json:"fileElem,omitempty"`
	MergeElem        *MergeElem   `json:"mergeElem,omitempty"`
	AtElem           *AtElem      `json:"atElem,omitempty"`

	FaceElem          *FaceElem          `json:"faceElem,omitempty"`
	LocationElem      *LocationElem      `json:"locationElem,omitempty"`
	CustomElem        *CustomElem        `json:"customElem,omitempty"`
	QuoteElem         *QuoteElem         `json:"quoteElem,omitempty"`
	NotificationElem  *NotificationElem  `json:"notificationElem,omitempty"`
	MessageEntityElem *MessageEntityElem `json:"messageEntityElem,omitempty"`
	AttachedInfoElem  AttachedInfoElem   `json:"attachedInfoElem,omitempty"`
}
type NotificationElem struct {
	Detail      string `json:"detail,omitempty"`
	DefaultTips string `json:"defaultTips,omitempty"`
}

type PictureBaseInfo struct {
	UUID   string `json:"uuid,omitempty"`
	Type   string `json:"type,omitempty"`
	Size   int64  `json:"size,omitempty"`
	Width  int32  `json:"width,omitempty"`
	Height int32  `json:"height,omitempty"`
	Url    string `json:"url,omitempty"`
}
type PictureElem struct {
	SourcePath      string          `json:"sourcePath,omitempty"`
	SourcePicture   PictureBaseInfo `json:"sourcePicture,omitempty"`
	BigPicture      PictureBaseInfo `json:"bigPicture,omitempty"`
	SnapshotPicture PictureBaseInfo `json:"snapshotPicture,omitempty"`
} //`json:"pictureElem,omitempty"`
type SoundElem struct {
	UUID      string `json:"uuid,omitempty"`
	SoundPath string `json:"soundPath,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	DataSize  int64  `json:"dataSize,omitempty"`
	Duration  int64  `json:"duration,omitempty"`
}

type VideoElem struct {
	VideoPath      string `json:"videoPath,omitempty"`
	VideoUUID      string `json:"videoUUID,omitempty"`
	VideoURL       string `json:"videoUrl,omitempty"`
	VideoType      string `json:"videoType,omitempty"`
	VideoSize      int64  `json:"videoSize,omitempty""`
	Duration       int64  `json:"duration,omitempty""`
	SnapshotPath   string `json:"snapshotPath,omitempty"`
	SnapshotUUID   string `json:"snapshotUUID,omitempty"`
	SnapshotSize   int64  `json:"snapshotSize,omitempty""`
	SnapshotURL    string `json:"snapshotUrl,omitempty"`
	SnapshotWidth  int32  `json:"snapshotWidth,omitempty""`
	SnapshotHeight int32  `json:"snapshotHeight,omitempty""`
}
type FileElem struct {
	FilePath  string `json:"filePath,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	FileName  string `json:"fileName,omitempty"`
	FileSize  int64  `json:"fileSize,omitempty"`
}
type MergeElem struct {
	Title             string           `json:"title,omitempty"`
	AbstractList      []string         `json:"abstractList,omitempty"`
	MultiMessage      []*MsgStruct     `json:"multiMessage,omitempty"`
	MessageEntityList []*MessageEntity `json:"messageEntityList,omitempty"`
}
type AtElem struct {
	Text         string     `json:"text,omitempty"`
	AtUserList   []string   `json:"atUserList,omitempty"`
	AtUsersInfo  []*AtInfo  `json:"atUsersInfo,omitempty"`
	QuoteMessage *MsgStruct `json:"quoteMessage,omitempty"`
	IsAtSelf     bool       `json:"isAtSelf"`
}
type FaceElem struct {
	Index int    `json:"index"`
	Data  string `json:"data,omitempty"`
}
type LocationElem struct {
	Description string  `json:"description,omitempty"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
}
type MessageEntityElem struct {
	Text              string           `json:"text,omitempty"`
	MessageEntityList []*MessageEntity `json:"messageEntityList,omitempty"`
}
type CustomElem struct {
	Data        string `json:"data,omitempty"`
	Description string `json:"description,omitempty"`
	Extension   string `json:"extension,omitempty"`
}
type QuoteElem struct {
	Text              string           `json:"text,omitempty"`
	QuoteMessage      *MsgStruct       `json:"quoteMessage,omitempty"`
	MessageEntityList []*MessageEntity `json:"messageEntityList,omitempty"`
}
type AtInfo struct {
	AtUserID      string `json:"atUserID,omitempty"`
	GroupNickname string `json:"groupNickname,omitempty"`
}
type AttachedInfoElem struct {
	GroupHasReadInfo          GroupHasReadInfo `json:"groupHasReadInfo,omitempty"`
	IsPrivateChat             bool             `json:"isPrivateChat,omitempty"`
	HasReadTime               int64            `json:"hasReadTime,omitempty"`
	NotSenderNotificationPush bool             `json:"notSenderNotificationPush,omitempty"`
	MessageEntityList         []*MessageEntity `json:"messageEntityList,omitempty"`
	IsEncryption              bool             `json:"isEncryption,omitempty"`
	InEncryptStatus           bool             `json:"inEncryptStatus,omitempty"`
}
type MessageEntity struct {
	Type   string `json:"type,omitempty"`
	Offset int32  `json:"offset"`
	Length int32  `json:"length"`
	Url    string `json:"url,omitempty"`
	Info   string `json:"info,omitempty"`
}

type GroupHasReadInfo struct {
	HasReadUserIDList []string `json:"hasReadUserIDList,omitempty"`
	HasReadCount      int32    `json:"hasReadCount"`
	GroupMemberCount  int32    `json:"groupMemberCount"`
}
type NewMsgList []*MsgStruct

// Len Implement the sort.Interface to get the number of elements method
func (n NewMsgList) Len() int {
	return len(n)
}

// Less Implement the sort.Interface  comparison element method
func (n NewMsgList) Less(i, j int) bool {
	return n[i].SendTime < n[j].SendTime
}

// Swap Implement the sort.Interface  exchange element method
func (n NewMsgList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

type GetHistoryMessageListResp struct {
	CommResp
	Data []*MsgStruct `json:"data"`
}

type GetCanRewordMemberCountRewordReq struct {
	MuteGroupReq
}
type GetCanRewordMemberCountRewordResp struct {
	CommResp
	Data int32 `json:"data"`
}

type ColligateSearchReq struct {
	OperationID string `json:"operationID" binding:"required"`
	SearchKey   string `json:"searchKey" binding:"required"`
}
type ColligateSearchResp struct {
	CommResp
	Data ColligateSearchBody `json:"data"`
}
type ColligateSearchGroupInfo struct {
	GroupID   string `json:"groupID"`
	GroupName string `json:"groupName"`
	FaceUrl   string `json:"faceUrl"`
}
type ColligateSearchUserInfo struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	FaceUrl  string `json:"faceUrl"`
}
type ColligateSearchBody struct {
	GroupList []*ColligateSearchGroupInfo `json:"groupList"`
	UserList  []*ColligateSearchUserInfo  `json:"userList"`
}

type CreatePushSpaceArticelOrderReq struct {
	OperationID    string `json:"operationID" binding:"required"`
	SpaceArticleId string `json:"spaceArticleId" binding:"required"`
	TxnType        string `json:"txnType" binding:"required"`
}
type CreatePushSpaceArticelOrderResp struct {
	CommResp
	Data PayInfo `json:"data"`
}

/* *****标签相关***************** */

// CreateRoleTagReq 创建标签
type CreateRoleTagReq struct {
	OperationID string `binding:"required"`
	RoleTitle   string `binding:"required"` //标签名称
	RoleIPfs    string `binding:"required"` //标签图片
	GroupID     string `binding:"required"`
}
type CreateRoleTagResp struct {
	CommResp
	GroupTagString string `json:"data"` //标签唯一值得
}
type GetCommunityRoleTagReq struct {
	OperationID string `binding:"required"`
	GroupID     string `json:"GroupID,omitempty"`
	Contract    string `json:"Contract,omitempty"`
	TokenID     string `json:"TokenID,omitempty"`
}
type GetCommunityRoleTagResp struct {
	CommResp
	CommitRoleTagReq []*CommitRoleTagReq `json:"data"`
}
type CommunityRoleUserInfoList struct {
	open_im_sdk.UserInfo
	Amount string `json:"amount"`
}
type GetCommunityRoleTagDetailResp struct {
	CommResp
	MemberList []*CommunityRoleUserInfoList `json:"data"`
}

// CommitRoleTagReq 确认块的方式：
type CommitRoleTagReq struct {
	OperationID string `json:"operationID,omitempty"`
	RoleID      string `json:"roleID,omitempty"`      //roleid
	GroupID     string `json:"groupID,omitempty"`     //哪个群
	TokenID     string `json:"tokenID,omitempty"`     //关联线上哪个token
	TokenAmount string `json:"tokenAmount,omitempty"` //发行总量
	TokenSub    string `json:"tokenSub,omitempty"`    //发行总量
	TokenBurn   string `json:"tokenBurn,omitempty"`   //发行总量
	RoleTitle   string `json:"roleTitle,omitempty"`
	RoleIPfs    string `json:"roleIPfs,omitempty"`
	Contract    string `json:"contract,omitempty"`
	ChainID     string `json:"chainID,omitempty"`
	Hash        string `json:"hash,omitempty"`
}
type CommitRoleTagResp struct {
	CommResp
}

type MintCommitRoleTagReq struct {
	Operator    string `json:"operator"` //mint  or burn
	Contract    string `json:"contract,omitempty"`
	TokenID     string `json:"tokenID,omitempty"`     //关联线上哪个token
	TokenAmount string `json:"tokenAmount,omitempty"` //发行总量
	UserAddress string `json:"userAddress,omitempty"`
	ChainID     string `json:"chainID,omitempty"`
	Hash        string `json:"hash,omitempty"`
	OperationID string `json:"operationID"`
}

// OperatorRoleUserReq 修改用关联的角色信息
type OperatorRoleUserReq struct {
	OperationID string
	UserID      []string
	Operator    string //行为操作
	RoleID      string
	groupID     string
}
type OperatorRoleUserResp struct {
	CommResp
}

// OperatorChannelRoleReq 修改channel
type OperatorChannelRoleReq struct {
	OperationID string
	ChannelID   string
	Operator    string //行为操作
	RoleID      string
	GroupID     string
}
type OperatorChannelRoleResp struct {
	CommResp
}
