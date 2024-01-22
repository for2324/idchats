package base_info

type BalanceOfCount struct {
	ErrCode   int    `json:"errCode"`
	ErrMsg    string `json:"errMsg"`
	BalanceOf int64  `json:"data"`
}
type CheckHaveNftReq struct {
	Address string
}

type RequestTokenIdReq struct {
	ContractAddress string `json:"contractAddress"`
	TokenID         string `json:"tokenID"`
	ChainID         string `json:"chainID"`
}
type RequestImageTokenIdReq struct {
	RequestTokenIdReq
	TokenImageModel string `json:"tokenImageModel"` //"ercContractAddress721 or erc1155"
	OwnerAddress    string `json:"ownerAddress"`
}
type RequestTokenIdResp struct {
	ErrCode           int    `json:"errCode"`
	ErrMsg            string `json:"errMsg"`
	TokenUrl          string `json:"tokenID"` //待优化， 后续修复
	TokenOwnerAddress string `json:"tokenOwnerAddress"`
}
type AppVersionResp struct {
	ErrCode            int                 `json:"errCode"`
	ErrMsg             string              `json:"errMsg"`
	AppVersionDataResp *AppVersionDataResp `json:"data"`
}
type AppVersionDataResp struct {
	Id          int32  `json:"id"`
	HasUpdate   bool   `json:"has_update"`
	IsIgnorable bool   `json:"is_ignorable"`
	VersionCode int32  ` json:"version_code"`
	VersionName string `json:"version_name"`
	UpdateLog   string ` json:"update_log"`
	UpdateLogEn string ` json:"update_log_en"`
	ApkUrl      string `json:"apk_url"`
	IosUrl      string `json:"ios_url"`
}
type TokenListResp struct {
}
type EnsApiReq struct {
	Address   string `json:"address" bind:"required"`
	EnsDomain string `json:"ensDomain" bind:"required"`
}
type EnsApiResp struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}
