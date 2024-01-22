package base_info

import (
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"time"

	gogpt "github.com/sashabaranov/go-openai"
)

type GetUsersInfoReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}
type NftInfo struct {
	TokenID         string `json:"tokenID"`
	NftChainID      string `json:"nftChainID"`
	NftContract     string `json:"nftContract"`
	NftContractType string `json:"nftContractType"`
	NftTokenURL     string `json:"nftTokenURL"`
	LikesCount      int64  `json:"likesCount"`
	ArticleID       int64  `json:"ID"`
	IsLikes         int32  `json:"isLikes"`
}
type UpdateSelfUserHeadReq struct {
	OperationID string `json:"operationID" binding:"required"`
	NftInfo
}
type UpdateSelfUserHeadResp struct {
	CommResp
}

type UpdateSelfUserInfoReq struct {
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}
type SetGlobalRecvMessageOptReq struct {
	OperationID      string `json:"operationID" binding:"required"`
	GlobalRecvMsgOpt *int32 `json:"globalRecvMsgOpt" binding:"omitempty,oneof=0 1 2"`
	ShowBalance      *int32 `json:"showBalance" binding:"omitempty,oneof=0 1 2"`
	OpenAnnouncement *int32 `json:"openAnnouncement" binding:"omitempty,oneof=0 1 2"`
}
type SetGlobalRecvMessageOptResp struct {
	CommResp
}
type UpdateUserInfoResp struct {
	CommResp
}

type GetSelfUserInfoReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"userID"`
}
type ApiSelfUserInfo struct {
	open_im_sdk.UserInfo
	GroupInfo *open_im_sdk.GroupInfo `json:"group,omitempty"`
}

type GetSelfUserInfoResp struct {
	CommResp
	Data ApiSelfUserInfo `json:"data"`
}

type GetFriendIDListFromCacheReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type GetFriendIDListFromCacheResp struct {
	CommResp
	UserIDList []string `json:"userIDList" binding:"required"`
}

type GetBlackIDListFromCacheReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type GetBlackIDListFromCacheResp struct {
	CommResp
	UserIDList []string `json:"userIDList" binding:"required"`
}
type GetUsersThirdInfoReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}
type UserThirdPath struct {
	UserId    string `json:"userId"`
	Twitter   string `json:"twitter"`
	DnsDomain string `json:"dnsDomain"`
	EnsDomain string `json:"ensDomain"`
}
type GetUsersInfoResp struct {
	CommResp
	UserInfoList []*open_im_sdk.PublicUserInfo `json:"-"`
	Data         []map[string]interface{}      `json:"data" swaggerignore:"true"`
}

type GetUsersThirdInfoResp struct {
	CommResp
	UserThirdInfoList []*open_im_sdk.UserThirdInfo `json:"-"`
	Data              []map[string]interface{}     `json:"data" swaggerignore:"true"`
}

// 绑定第三方api
type BindUserSelfDomainReq struct {
	OperationID     string `json:"operationID" binding:"required"`
	EnsDomain       string `json:"ensDomain"`
	DnsDomain       string `json:"dnsDomain"`
	EmailAddress    string `json:"emailAddress"`
	EmailVerifyCode string `json:"emailVerifyCode"`
}

type BindUserSelfDomainResp struct {
	CommResp
}

type BindUserTelephoneCodeReq struct {
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	OperationID    string `json:"operationID" binding:"required"`
	UsedFor        int    `json:"usedFor"`
	AreaCode       string `json:"areaCode"` //區域 86  1 862
	InvitationCode string `json:"invitationCode"`
}
type BindUserTelephoneCodeResp struct {
	CommResp
}
type BindUserTelephoneReq struct {
	OperationID   string `json:"operationID" binding:"required"`
	TelephoneCode string `json:"code"`     //提交验证码去检验信息
	AreaCode      string `json:"areaCode"` //區域 86  1 862
	PhoneNumber   string `json:"phoneNumber"`
	Email         string `json:"email"`
	UpdateSecret  string `json:"updateSecret"`
}
type BindUserTelephoneResp struct {
	CommResp
}
type UpdateUserSignReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserProfile string `json:"userProfile"`
}
type UpdateUserSignResp struct {
	CommResp
}
type GetUserSignReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserId      string `json:"userId"`
}
type GetUserSignResp struct {
	CommResp
	UserInfo *open_im_sdk.UserInfo  `json:"-"`
	Data     map[string]interface{} `json:"data" swaggerignore:"true"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type TransferTokenToGroupReq struct {
	OperationID   string
	UserChatToken int64  //转入多少余额
	ToGroupID     string //注入的groupid
}
type TransferTokenToGroupRsp struct {
	CommResp
	ChatTokenCount map[string]int64 `json:"data"`
}

type GetGlobalUserProfileReq struct {
	OperationID string
	UserID      string `json:"userID" bind:"required"`
	ChainNet    string `json:"chainNet"` //eth-mainnet
}
type GetGlobalUserProfileResp struct {
	CommResp
	UserInfoProfile map[string]interface{} `json:"data"`
}
type TBtcWalletDetail struct {
	ConfirmAmount            string                       `json:"confirm_amount"`
	PendingAmount            string                       `json:"pending_amount"`
	Amount                   string                       `json:"amount"`
	ConfirmBtcAmount         string                       `json:"confirm_btc_amount"`
	PendingBtcAmount         string                       `json:"pending_btc_amount"`
	BtcAmount                string                       `json:"btc_amount"`
	ConfirmInscriptionAmount string                       `json:"confirm_inscription_amount"`
	PendingInscriptionAmount string                       `json:"pending_inscription_amount"`
	InscriptionAmount        string                       `json:"inscription_amount"`
	UsdValue                 string                       `json:"usd_value"`
	List                     []*TBtcWalletDetailTokenList `json:"list"`
	Total                    int                          `json:"total"`
}
type UserBtcBalanceInfoResp struct {
	CommResp
	Data *TBtcWalletDetail `json:"data"`
}
type TBtcWalletDetailTokenList struct {
	Ticker              string `json:"ticker"`
	OverallBalance      string `json:"overallBalance"`
	TransferableBalance string `json:"transferableBalance"`
	AvailableBalance    string `json:"availableBalance"`
}

type UserBalanceInfoResp struct {
	CommResp
	UserBalanceInfoDetail *UserBalanceInfoDetail `json:"data"`
	Error                 bool                   `json:"error"`
	ErrorMessage          interface{}            `json:"error_message"`
	ErrorCode             interface{}            `json:"error_code"`
}
type UserBalanceInfoDetail struct {
	Address       string                  `json:"address"`
	UpdatedAt     time.Time               `json:"updated_at"`
	NextUpdateAt  time.Time               `json:"next_update_at"`
	QuoteCurrency string                  `json:"quote_currency"`
	ChainID       int                     `json:"chain_id"`
	Items         []*UserBalanceInfoItems `json:"items"`
	Pagination    interface{}             `json:"pagination"`
}
type UserBalanceInfoItems struct {
	ContractDecimals     int         `json:"contract_decimals"`
	ContractName         string      `json:"contract_name"`
	ContractTickerSymbol string      `json:"contract_ticker_symbol"`
	ContractAddress      string      `json:"contract_address"`
	SupportsErc          []string    `json:"supports_erc"`
	LogoURL              string      `json:"logo_url"`
	LastTransferredAt    time.Time   `json:"last_transferred_at"`
	NativeToken          bool        `json:"native_token"`
	Type                 string      `json:"type"`
	Balance              string      `json:"balance"`
	Balance24H           string      `json:"balance_24h"`
	QuoteRate            float64     `json:"quote_rate"`
	QuoteRate24H         float64     `json:"quote_rate_24h"`
	Quote                float64     `json:"quote"`
	Quote24H             float64     `json:"quote_24h"`
	NftData              interface{} `json:"nft_data"`
}
type AskQuestionGptReq struct {
	OperationID       string                      `json:"operationID" binding:"required"`
	RecvID            string                      `json:"recvID" binding:"required"`
	ChannelID         string                      `json:"channelID"`
	ClientMsgID       string                      `json:"clientMsgID"`
	AskMode           string                      `json:"askMode"` //提问的方式
	CompletionRequest gogpt.ChatCompletionRequest `json:"completionRequest" binding:"required"`
}
type AskQuestionGptResp struct {
	CommResp
	CompletionResponse interface{} `json:"data"`
	UserChatToken      uint64
	UserChatCount      uint64
}

// GetUserChatTokenReq 查询用户的chat token
type GetUserChatTokenReq struct {
	OperationID string `json:"operationID" binding:"required"`
}
type GetUserChatTokenRsp struct {
	CommResp
	ChatTokenCount   uint64 `json:"chat_token_count"`
	ChatCount        uint64 `json:"chat_count"`
	GlobalMoneyCount int64  `json:"global_money_count"`
}

// GetUserChatTokenReq 查询用户的chat token
type GetUserChatTokenHistoryReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromID      string `json:"fromID,omitempty"`
	PageCount   string `json:"pageCount,omitempty"`
	PageIndex   string `json:"pageIndex,omitempty"`
}

type GetUserChatTokenHistoryRsp struct {
	CommResp
	UserChatTokenHistory []*UserChatTokenHistory `json:"data"`
}
type UserChatTokenHistory struct {
	Action     string `json:"action,omitempty"`     //行为分析
	Param      string `json:"param,omitempty"`      //消费内容
	CreateTime string `json:"createTime,omitempty"` //消费时间
	ID         string `json:"ID"`
	ChainID    string `json:"chainID"`
	Value      string `json:"value"`
}

// UserChatTokenReq 只给内网调用
type UserChatTokenReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string //操作的用户对象。
	Operator    string //操作的方法 添加删除token
	Value       int64
	ChainID     string `json:"chainID"`
	TxID        string `json:"txID"`   // 交易序号 ，添加金币的时候需要这个参数来获取 做成记录事务的数据来做成单独的进程的初入帐。
	TxType      string `json:"txType"` //交易类型 充值时候需要
}
type UserChatTokenResp struct {
	CommResp
	UserChatToken uint64
}

type UserChatTokenOperatorReq struct {
	TxID string `json:"txID"` //
}
type BindShowNftReq struct {
	OperationID string     `json:"operationID" binding:"required"`
	NftListShow []*NftInfo `json:"nftListShow"`
}
type BindShowNftResp struct {
	CommResp
}
type GetShowNftResp struct {
	CommResp
	NftListShow []*NftInfo `json:"nftListShow"`
}
type LikeActionNftReq struct {
	OperationID string `json:"operationID" binding:"required"`
	Action      int32  `json:"action" binding:"required"`
	ArticleID   string `json:"articleID" binding:"required"`
}

// AnnouncementReq 读取谋篇文章
type AnnouncementReq struct {
	OperationID string `json:"operationID" binding:"required"`
	ArticleID   string `json:"articleID" binding:"required"`
}
type LikeActionNftResp struct {
	CommResp
	NCount int64 //暂时未使用
	IsLike int32 //是否已经点赞
}
type GetLikeShowNftCountReq struct {
	OperationID string `json:"operationID" binding:"required"`
	ArticleID   string //nft的id
}

// 获取个人设置页面
type UserSettingPageReq struct {
	UserID      string //获取个人页面的数据信息。
	OperationID string `json:"operationID" binding:"required"`
}
type UserSettingInfoWithFriend struct {
	UserSettingInfo
	FollowsCount   int64 `json:"followsCount"`
	FollowingCount int64 `json:"followingCount"`
}
type UserSettingInfo struct {
	FaceURL           *string                        `json:"faceURL"`
	Nickname          *string                        `json:"nickname"`
	UserID            *string                        `json:"userID"`
	UserIntroduction  *string                        `json:"userIntroduction"`
	UserProfile       *string                        `json:"userProfile"`
	UserHeadTokenInfo *NftInfo                       `json:"userHeadTokenInfo"`
	NftListShow       []*NftInfo                     `json:"nftListShow"` //展示的nft
	EmailAddress      *string                        `json:"emailAddress"`
	UserTwitter       *string                        `json:"userTwitter"`
	DnsDomain         *string                        `json:"dnsDomain"`
	IsShowTwitter     *bool                          `json:"isShowTwitter"`
	IsShowEmail       *bool                          `json:"isShowEmail"`
	IsShowBalance     *bool                          `json:"isShowBalance"`
	OpenAnnouncement  *int32                         `json:"openAnnouncement"`
	DnsDomainVerify   *int32                         `json:"dnsDomainVerify"` //不要去设置这个值得
	LinkTree          *[]*open_im_sdk.LinkTreeMsgReq `json:"linkTree"`
}
type UserSettingPageResp struct {
	CommResp
	UserSettingInfoWithFriend `json:"data"`
}
type UpdateUserSettingPageReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserSettingInfo
}
type UpdateUserSettingPageResp struct {
	CommResp
}
