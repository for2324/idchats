package sdk_struct

import "open_im_sdk/pkg/server_api_params"

////////////////////////// message/////////////////////////

type MessageReceipt struct {
	GroupID     string   `json:"groupID"`
	UserID      string   `json:"userID"`
	MsgIDList   []string `json:"msgIDList"`
	ReadTime    int64    `json:"readTime"`
	MsgFrom     int32    `json:"msgFrom"`
	ContentType int32    `json:"contentType"`
	SessionType int32    `json:"sessionType"`
}
type MessageRevoked struct {
	RevokerID                   string `json:"revokerID"`
	RevokerRole                 int32  `json:"revokerRole"`
	ClientMsgID                 string `json:"clientMsgID"`
	RevokerNickname             string `json:"revokerNickname"`
	RevokeTime                  int64  `json:"revokeTime"`
	SourceMessageSendTime       int64  `json:"sourceMessageSendTime"`
	SourceMessageSendID         string `json:"sourceMessageSendID"`
	SourceMessageSenderNickname string `json:"sourceMessageSenderNickname"`
	SessionType                 int32  `json:"sessionType"`
	Seq                         uint32 `json:"seq"`
}
type ImageInfo struct {
	Width  int32  `json:"x,omitempty"`
	Height int32  `json:"y,omitempty"`
	Type   string `json:"type,omitempty"`
	Size   int64  `json:"size,omitempty"`
}
type PictureBaseInfo struct {
	UUID   string `json:"uuid,omitempty"`
	Type   string `json:"type,omitempty"`
	Size   int64  `json:"size,omitempty"`
	Width  int32  `json:"width,omitempty"`
	Height int32  `json:"height,omitempty"`
	Url    string `json:"url,omitempty"`
}
type SoundBaseInfo struct {
	UUID      string `json:"uuid,omitempty"`
	SoundPath string `json:"soundPath,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	DataSize  int64  `json:"dataSize,omitempty"`
	Duration  int64  `json:"duration,omitempty"`
}
type VideoBaseInfo struct {
	VideoPath      string `json:"videoPath,omitempty"`
	VideoUUID      string `json:"videoUUID,omitempty"`
	VideoURL       string `json:"videoUrl,omitempty"`
	VideoType      string `json:"videoType,omitempty"`
	VideoSize      int64  `json:"videoSize,omitempty"`
	Duration       int64  `json:"duration,omitempty"`
	SnapshotPath   string `json:"snapshotPath,omitempty"`
	SnapshotUUID   string `json:"snapshotUUID,omitempty"`
	SnapshotSize   int64  `json:"snapshotSize,omitempty"`
	SnapshotURL    string `json:"snapshotUrl,omitempty"`
	SnapshotWidth  int32  `json:"snapshotWidth,omitempty"`
	SnapshotHeight int32  `json:"snapshotHeight,omitempty"`
}
type FileBaseInfo struct {
	FilePath  string `json:"filePath,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	FileName  string `json:"fileName,omitempty"`
	FileSize  int64  `json:"fileSize,omitempty"`
}

type MsgStruct struct {
	ClientMsgID       string                            `json:"clientMsgID,omitempty"`
	ServerMsgID       string                            `json:"serverMsgID,omitempty"`
	CreateTime        int64                             `json:"createTime"`
	SendTime          int64                             `json:"sendTime"`
	SessionType       int32                             `json:"sessionType"`
	SendID            string                            `json:"sendID,omitempty"`
	RecvID            string                            `json:"recvID,omitempty"`
	MsgFrom           int32                             `json:"msgFrom"`
	ContentType       int32                             `json:"contentType"`
	SenderPlatformID  int32                             `json:"platformID"`
	SenderNickname    string                            `json:"senderNickname,omitempty"`
	SenderFaceURL     string                            `json:"senderFaceUrl,omitempty"`
	GroupID           string                            `json:"groupID,omitempty"`
	ChannelID         string                            `json:"channelID,omitempty"`
	Content           string                            `json:"content,omitempty"`
	Seq               uint32                            `json:"seq"`
	IsRead            bool                              `json:"isRead"`
	Status            int32                             `json:"status"`
	OfflinePush       server_api_params.OfflinePushInfo `json:"offlinePush,omitempty"`
	AttachedInfo      string                            `json:"attachedInfo,omitempty"`
	Ex                string                            `json:"ex,omitempty"`
	PictureElem       *PictureElem                      `json:"pictureElem,omitempty"`
	SoundElem         *SoundElem                        `json:"soundElem,omitempty"`
	VideoElem         *VideoElem                        `json:"videoElem,omitempty"`
	FileElem          *FileElem                         `json:"fileElem,omitempty"`
	MergeElem         *MergeElem                        `json:"mergeElem,omitempty"`
	AtElem            *AtElem                           `json:"atElem,omitempty"`
	FaceElem          *FaceElem                         `json:"faceElem,omitempty"`
	LocationElem      *LocationElem                     `json:"locationElem,omitempty"`
	CustomElem        *CustomElem                       `json:"customElem,omitempty"`
	QuoteElem         *QuoteElem                        `json:"quoteElem,omitempty"`
	NotificationElem  *NotificationElem                 `json:"notificationElem,omitempty"`
	MessageEntityElem *MessageEntityElem                `json:"messageEntityElem,omitempty"`
	AttachedInfoElem  AttachedInfoElem                  `json:"attachedInfoElem,omitempty"`
	AnnouncementElem  *AnnouncementElem                 `json:"announcementElem,omitempty"`
}
type NotificationElem struct {
	Detail      string `json:"detail,omitempty"`
	DefaultTips string `json:"defaultTips,omitempty"`
}

func (s *MsgStruct) InitStruct() {
	s.PictureElem = new(PictureElem)
	s.SoundElem = new(SoundElem)
	s.VideoElem = new(VideoElem)
	s.FileElem = new(FileElem)
	s.MergeElem = new(MergeElem)
	s.AtElem = new(AtElem)
	s.FaceElem = new(FaceElem)
	s.LocationElem = new(LocationElem)
	s.QuoteElem = new(QuoteElem)
	s.CustomElem = new(CustomElem)
	s.MessageEntityElem = new(MessageEntityElem)
	//s.NotificationElem = new(NotificationElem)
	s.AnnouncementElem = new(AnnouncementElem)
}

type AnnouncementElem struct {
	Announcement  string   `json:"announcement"`   //在ipfs的公告的信息 url //文字、文字+视频、文字+图片、文字+图片+链接、文字+链接、图片、链接（单独一条外部链接）
	Text          string   `json:"text,omitempty"` //文字内容
	VideoPath     []string `json:"videoPath,omitempty"`
	PicturePath   []string `json:"videoPath,omitempty"`
	UrlPath       []string `json:"urlPath,omitempty""` //是否可以对url鉴权
	IsGlobalWorld bool     `json:"isGlobalWorld"`
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

type IMConfig struct {
	Platform      int32  `json:"platform"`
	ApiAddr       string `json:"api_addr"`
	WsAddr        string `json:"ws_addr"`
	DataDir       string `json:"data_dir"`
	LogLevel      uint32 `json:"log_level"`
	ObjectStorage string `json:"object_storage"` //"cos"(default)  "oss"
	EncryptionKey string `json:"encryption_key"`
}

var SvrConf IMConfig

type CmdNewMsgComeToConversation struct {
	MsgList       []*server_api_params.MsgData
	OperationID   string
	SyncFlag      int
	MaxSeqOnSvr   uint32
	MaxSeqOnLocal uint32
	CurrentMaxSeq uint32
	PullMsgOrder  int
}

type CmdPushMsgToMsgSync struct {
	Msg         *server_api_params.MsgData
	OperationID string
}

type CmdMaxSeqToMsgSync struct {
	MaxSeqOnSvr            uint32
	OperationID            string
	MinSeqOnSvr            uint32
	GroupID2MinMaxSeqOnSvr map[string]*server_api_params.MaxAndMinSeq
}

type CmdJoinedSuperGroup struct {
	OperationID string
}

type OANotificationElem struct {
	NotificationName    string `mapstructure:"notificationName" validate:"required"`
	NotificationFaceURL string `mapstructure:"notificationFaceURL" validate:"required"`
	NotificationType    int32  `mapstructure:"notificationType" validate:"required"`
	Text                string `mapstructure:"text" validate:"required"`
	Url                 string `mapstructure:"url"`
	MixType             int32  `mapstructure:"mixType"`
	Image               struct {
		SourceUrl   string `mapstructure:"sourceURL"`
		SnapshotUrl string `mapstructure:"snapshotURL"`
	} `mapstructure:"image"`
	Video struct {
		SourceUrl   string `mapstructure:"sourceURL"`
		SnapshotUrl string `mapstructure:"snapshotURL"`
		Duration    int64  `mapstructure:"duration"`
	} `mapstructure:"video"`
	File struct {
		SourceUrl string `mapstructure:"sourceURL"`
		FileName  string `mapstructure:"fileName"`
		FileSize  int64  `mapstructure:"fileSize"`
	} `mapstructure:"file"`
	Ex string `mapstructure:"ex"`
}
type MsgDeleteNotificationElem struct {
	GroupID     string   `json:"groupID"`
	IsAllDelete bool     `json:"isAllDelete"`
	SeqList     []string `json:"seqList"`
}
