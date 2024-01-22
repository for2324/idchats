package group

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

func (g *Group) MonitorGroupMessage(callback open_im_sdk_callback.Base, groupID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	return
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID)
		svrGroup, err := g.getGroupsInfoFromSvr([]string{groupID}, operationID)
		if err == nil || len(svrGroup) == 0 {
			groupInfo := svrGroup[0]
			conversationId := utils.GetConversationIDBySessionType(groupInfo.GroupID, constant.GroupChatType)
			conversationType := constant.GroupChatType
			if groupInfo.GroupType == constant.WorkingGroup {
				conversationId = utils.GetConversationIDBySessionType(groupInfo.GroupID, constant.SuperGroupChatType)
				conversationType = constant.SuperGroupChatType
			}
			var reqPb api.SetConversationReq
			var apiResp api.SetConversationResp
			reqPb.OperationID = operationID
			reqPb.OwnerUserID = g.loginUserID
			reqPb.ConversationID = conversationId
			reqPb.ConversationType = int32(conversationType)
			reqPb.GroupID = groupID
			reqPb.IsNotInGroup = false
			g.p.PostFatalCallback(callback, constant.MonitorConversationsRouter, reqPb, nil, operationID)
			callback.OnSuccess(utils.StructToJsonString(apiResp))
			log.NewInfo(operationID, fName, " callback: ")
		} else {
			common.CheckDBErrCallback(callback, err, operationID)
		}

	}()
}
