package base_info

// 可能没有token的传输，  以
type CreateCommunityRoleDataReq struct {
	GroupID    string `json:"groupID"`
	RoleTitle  string `json:"roleTitle"`
	RoleIcon   string `json:"roleIcon"`
	Contract   string `json:"contract"`
	ChainID    string `json:"chainID"`
	TokenID    string `json:"tokenID"`
	OperatorID string `json:"operatorID"`
}

type CreateCommunityRoleDataResp struct {
	CommResp
}

type CommunityUserRoleReq struct {
	GroupID    string `json:"groupID"`
	Contract   string `json:"contract"`
	ChainID    string `json:"chainID"`
	TokenID    string `json:"tokenID"`
	OperatorID string `json:"operatorID"`
}
type AddCommunityUserRoleRsp struct {
	CommResp
}
