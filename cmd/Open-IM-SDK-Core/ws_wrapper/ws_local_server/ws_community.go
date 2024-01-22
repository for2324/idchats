package ws_local_server

/*
func (wsRouter *WsFuncRouter) AddGroupChannelInfo(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupInfo", "groupID") {
		return
	}
	//(callback common.Base, groupInfo string, groupID string, operationID string)
	userWorker.Group().SetGroupInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupInfo"].(string), m["groupID"].(string), operationID)

	userWorker.Group().CreateGroup(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupBaseInfo"].(string), m["memberList"].(string), operationID)

}
func (wsRouter *WsFuncRouter) UpdateGroupChannelInfo(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupInfo", "groupID") {
		return
	}
	//(callback common.Base, groupInfo string, groupID string, operationID string)
	userWorker.Group().SetGroupInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupInfo"].(string), m["groupID"].(string), operationID)
}
func (wsRouter *WsFuncRouter) DeleteGroupChannelInfo(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupInfo", "groupID") {
		return
	}
	//(callback common.Base, groupInfo string, groupID string, operationID string)
	userWorker.Group().SetGroupInfo(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		m["groupInfo"].(string), m["groupID"].(string), operationID)
}
*/
