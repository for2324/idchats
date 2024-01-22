package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/server_api_params"
)

// other user
type GetUsersInfoParam []string
type GetUsersInfoCallback []server_api_params.FullUserInfo

// type GetSelfUserInfoParam string
type GetSelfUserInfoCallback *model_struct.LocalUser
type GetSelfUserAndGroupInfoCallback struct {
	UserID           string                      `json:"userID" binding:"required,min=1,max=64"`
	Nickname         string                      `json:"nickname" binding:"omitempty,min=1,max=64"`
	FaceURL          string                      `json:"faceURL" binding:"omitempty,max=1024"`
	Gender           int32                       `json:"gender" binding:"omitempty,oneof=0 1 2"`
	PhoneNumber      string                      `json:"phoneNumber" binding:"omitempty,max=32"`
	Birth            uint32                      `json:"birth" binding:"omitempty"`
	Email            string                      `json:"email" binding:"omitempty,max=64"`
	UserIntroduction string                      `json:"userIntroduction" binding:"omitempty"`
	GlobalRecvMsgOpt int32                       `json:"globalRecvMsgOpt" binding:"omitempty,oneof=0 1 2"`
	Ex               string                      `json:"ex" binding:"omitempty,max=1024"`
	Group            server_api_params.GroupInfo `json:"group,omitempty"`
}

type SetSelfUserInfoParam server_api_params.ApiUserInfo

const SetSelfUserInfoCallback = constant.SuccessCallbackDefault
