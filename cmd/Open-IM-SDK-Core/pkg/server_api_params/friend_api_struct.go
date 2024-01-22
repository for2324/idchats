package server_api_params

type ParamsCommFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	ToUserID    string `json:"toUserID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}

type AddBlacklistReq struct {
	ParamsCommFriend
}
type AddBlacklistResp struct {
	CommResp
}

type ImportFriendReq struct {
	FriendUserIDList []string `json:"friendUserIDList" binding:"required"`
	OperationID      string   `json:"operationID" binding:"required"`
	FromUserID       string   `json:"fromUserID" binding:"required"`
}
type UserIDResult struct {
	UserID string `json:"userID""`
	Result int32  `json:"result"`
}
type ImportFriendResp struct {
	CommResp
	UserIDResultList []UserIDResult `json:"data"`
}

type AddFriendReq struct {
	ParamsCommFriend
	ReqMsg string `json:"reqMsg"`
}
type AddFriendResp struct {
	CommResp
}
type FollowAddFriendReq struct {
	ParamsCommFriend
	Follow bool `json:"follow"`
}
type FollowAddFriendRsp struct {
	CommResp
	Follow         bool
	PublicUserInfo []*PublicUserInfo
	Data           []map[string]interface{} `json:"data"`
}

type AddFriendResponseReq struct {
	ParamsCommFriend
	Flag      int32  `json:"flag" binding:"required,oneof=-1 0 1"`
	HandleMsg string `json:"handleMsg"`
}
type AddFriendResponseResp struct {
	CommResp
}

type DeleteFriendReq struct {
	ParamsCommFriend
}
type DeleteFriendResp struct {
	CommResp
}

type GetBlackListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetBlackListResp struct {
	CommResp
	BlackUserInfoList []*PublicUserInfo
	Data              []map[string]interface{} `json:"data"`
}
type SetFriendRemarkReq struct {
	ParamsCommFriend
	Remark string `json:"remark" binding:"required"`
}
type SetFriendRemarkResp struct {
	CommResp
}

type RemoveBlackListReq struct {
	ParamsCommFriend
}
type RemoveBlackListResp struct {
	CommResp
}

type IsFriendReq struct {
	ParamsCommFriend
}
type Response struct {
	Friend bool `json:"isFriend"`
}
type IsFriendResp struct {
	CommResp
	Response Response `json:"data"`
}

type GetFriendsInfoReq struct {
	ParamsCommFriend
}
type GetFriendsInfoResp struct {
	CommResp
	FriendInfoList []*FriendInfo
	Data           []map[string]interface{} `json:"data"`
}

type GetFriendListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetFriendListResp struct {
	CommResp
	FriendInfoList []*FriendInfo
	Data           []map[string]interface{} `json:"data"`
}

type GetFollowEachOtherFriendListReq struct {
	OperationID string `json:"operationID" binding:"required"`
}
type GetFollowEachOtherFriendListResp struct {
	CommResp
	Data []string `json:"data"`
}

type GetFriendApplyListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetFriendApplyListResp struct {
	CommResp
	FriendRequestList []*FriendRequest
	Data              []map[string]interface{} `json:"data"`
}

type GetSelfFriendApplyListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetSelfFriendApplyListResp struct {
	CommResp
	FriendRequestList []*FriendRequest
	Data              []map[string]interface{} `json:"data"`
}
type GetFollowListReq struct {
	ParamsCommFriend
	IsFollow bool `json:"is_follow"` //true 跟随我的， false 跟随他人的
}
type FollowUserInfoList struct {
	UserID             string    `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
	Nickname           string    `protobuf:"bytes,2,opt,name=nickname,proto3" json:"nickname,omitempty"`
	FaceURL            string    `protobuf:"bytes,3,opt,name=faceURL,proto3" json:"faceURL,omitempty"`
	Gender             int32     `protobuf:"varint,4,opt,name=gender,proto3" json:"gender,omitempty"`
	Ex                 string    `protobuf:"bytes,5,opt,name=ex,proto3" json:"ex,omitempty"`
	Remark             string    `protobuf:"bytes,6,opt,name=remark,proto3" json:"remark,omitempty"`
	UserProfile        string    `protobuf:"bytes,7,opt,name=userProfile,proto3" json:"userProfile,omitempty"`
	TokenContractChain string    `protobuf:"bytes,8,opt,name=tokenContractChain,proto3" json:"tokenContractChain,omitempty"`
	Group              GroupInfo `json:"group,omitempty"`
}

type GetFollowListRsp struct {
	CommResp
	PublicUserInfo []*FollowUserInfoList
	Data           []map[string]interface{} `json:"data"`
}
