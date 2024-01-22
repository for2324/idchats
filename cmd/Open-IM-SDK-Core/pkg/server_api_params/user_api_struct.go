package server_api_params

type GetUsersInfoReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}
type GetUsersInfoResp struct {
	CommResp
	UserInfoList []*PublicUserInfo
	Data         []map[string]interface{} `json:"data"`
}

type GetUserScoreReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type UserScoreInfo struct {
	UserId string `json:"userId"`
	Score  int64  `json:"score"`
}
type ApiSpaceArticleMqData struct {
	SpaceID     string `json:"spaceID"`     // 全局文章id
	OpID        string `json:"opID"`        //运营人员 即发送文章者。
	CreateID    string `json:"createID"`    // 创建人与恩
	ArticleType string `json:"articleType"` //文章类型 包括文章和ido
	IsGlobal    string `json:"isGlobal"`    //是否全局推送
	CreatedAt   int64  `json:"createdAt"`   // 创建时间
}

type GetUserScoreResp struct {
	CommResp
	Data UserScoreInfo `json:"data"`
}

type UpdateSelfUserHeadReq struct {
	OperationID string `json:"operationID" binding:"required"`
	TokenID     string `json:"tokenID"`
	NftChainID  string `json:"nftChainID"`
	NftContract string `json:"nftContract"`
}
type UpdateSelfUserHeadResp struct {
	CommResp
}
type UpdateSelfUserInfoReq struct {
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}

type UpdateUserInfoResp struct {
	CommResp
}

type SetGlobalRecvMessageOptReq struct {
	OperationID      string `json:"operationID" binding:"required"`
	GlobalRecvMsgOpt *int32 `json:"globalRecvMsgOpt" binding:"omitempty,oneof=0 1 2"`
}
type SetGlobalRecvMessageOptResp struct {
	CommResp
}
type GetSelfUserInfoReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}
type GetSelfUserInfoResp struct {
	CommResp
	UserInfo *UserInfo              `json:"-"`
	Data     map[string]interface{} `json:"data"`
}
