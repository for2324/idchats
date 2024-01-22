package base_info

// GetUserThirdReq 第三方某个平台的发言段落
type GetUserThirdReq struct {
	// 操作ID
	OperationID string `json:"operationid" binding:"required"`
	//  第三方平台：[ twitter,weibo,facebook ]
	ThirdString string `json:"thirdstring" bindding:"required"`
}

// GetUserThirdRsp 返回 第三方某个平台的发言段落
type GetUserThirdRsp struct {
	CommResp
	//返回需要发到第三方帖子的内容
	Data string `json:"data"`
}

// VerifyThirdStringReq 验证发言在第三方平台上的内容，提交url给后端去获取短文内容 并验证签名
type VerifyThirdStringReq struct {
	PostUrl     string `json:"posturl"  binding:"required"`
	UserId      string `json:"userid"  binding:"required"`
	ThirdString string `json:"thirdstring"  binding:"required"`
	OperationID string `json:"operationid" binding:"required"`
}

// VerifyThirdStringRsp 返回验证发言在第三方平台上的内容
type VerifyThirdStringRsp struct {
	CommResp
}

// VerifyThirdStatusReq 查询已经认证的平台的内容
type VerifyThirdStatusReq struct {
}

// VerifyThirdStatusRsp 返回查询的结果
type VerifyThirdStatusRsp struct {
	CommResp
	//已经认证的接口
	ThirdString []string `json:"thirdstring"`
}

// GetUserAuthorizedThirdPlatformListReq 获取已经授权平台列表
type GetUserAuthorizedThirdPlatformListReq struct {
	OperationID string `json:"operationid" binding:"required"`
}
type DelThirdPlatformReq struct {
	OperationID  string `json:"operationid" binding:"required"`
	PlatformName string `json:"platformName"`
}
type DelThirdPlatformResp struct {
	CommResp
}
type ShowThirdPlatformReq struct {
	OperationID  string `json:"operationid" binding:"required"`
	ShowFlag     bool   `json:"showFlag"`
	PlatformName string `json:"platformName"`
}
type ShowThirdPlatformResp struct {
	CommResp
}
type CheckIsFinishTaskReq struct {
	OperationID string `json:"operationId"`
	TaskId      string `json:"taskId"`
	GroupID     string `json:"groupID"`
}
type CheckIsFinishTaskRsp struct {
	CommResp
}
