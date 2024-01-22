package base_info

type DelMsgReq struct {
	UserID      string   `json:"userID,omitempty" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty" binding:"required"`
	OperationID string   `json:"operationID,omitempty" binding:"required"`
}

type DelMsgResp struct {
	CommResp
}

type CleanUpMsgReq struct {
	UserID      string `json:"userID"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}

type CleanUpMsgResp struct {
	CommResp
}
type DelSuperGroupMsgReq struct {
	UserID      string   `json:"userID" binding:"required"`
	GroupID     string   `json:"groupID" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty"`
	IsAllDelete bool     `json:"isAllDelete"`
	OperationID string   `json:"operationID" binding:"required"`
}

type DelSuperGroupMsgResp struct {
	CommResp
}
type MsgDeleteNotificationElem struct {
	GroupID     string   `json:"groupID"`
	IsAllDelete bool     `json:"isAllDelete"`
	SeqList     []uint32 `json:"seqList"`
}

// UserID               string   `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
//
//	GroupID              string   `protobuf:"bytes,2,opt,name=groupID" json:"groupID,omitempty"`
//	MinSeq               uint32   `protobuf:"varint,3,opt,name=minSeq" json:"minSeq,omitempty"`
//	OperationID          string   `protobuf:"bytes,4,opt,name=operationID" json:"operationID,omitempty"`
//	OpUserID             string   `protobuf:"bytes,5,opt,name=opUserID" json:"opUserID,omitempty"`
type SetMsgMinSeqReq struct {
	UserID      string `json:"userID"  binding:"required"`
	GroupID     string `json:"groupID"`
	MinSeq      uint32 `json:"minSeq"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}
type SetMsgMinSeqResp struct {
	CommResp
}
type SendMsgReqStructReq struct {
	SenderPlatformID int32              `json:"senderPlatformID" binding:"required"`
	SendID           string             `json:"sendID" binding:"required"`
	SenderNickName   string             `json:"senderNickName"`
	SenderFaceURL    string             `json:"senderFaceUrl"`
	OperationID      string             `json:"operationID" binding:"required"`
	Data             *SendMsgReqDataReq `json:"data"`
}
type SendMsgReqDataReq struct {
	SessionType int32           `json:"sessionType" binding:"required"`
	MsgFrom     int32           `json:"msgFrom" binding:"required"`
	ContentType int32           `json:"contentType" binding:"required"`
	RecvID      string          `json:"recvID" `
	GroupID     string          `json:"groupID" `
	ChannelID   string          `json:"channelID"`
	ForceList   []string        `json:"forceList"`
	Content     []byte          `json:"content" binding:"required"`
	Options     map[string]bool `json:"options" `
	ClientMsgID string          `json:"clientMsgID" binding:"required"`
	CreateTime  int64           `json:"createTime" binding:"required"`
}
