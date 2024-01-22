package task

// func CreateTask(c *gin.Context) {
// 	var (
// 		req   cms_api_struct.CreateTaskReq
// 		resp  cms_api_struct.CreateTaskResp
// 		reqPb pbTask.CreateTaskReq
// 	)
// 	if err := c.BindJSON(&req); err != nil {
// 		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
// 		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
// 		return
// 	}
// 	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
// 	utils.CopyStructFields(&reqPb, req)
// 	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, reqPb.OperationID)
// 	if etcdConn == nil {
// 		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
// 		log.NewError(reqPb.OperationID, errMsg)
// 		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
// 		return
// 	}
// 	client := pbTask.NewTaskServiceClient(etcdConn)
// 	respPb, err := client.CreateTask(context.Background(), &reqPb)
// 	if err != nil {
// 		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "Create Task failed ", err.Error())
// 		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
// 		return
// 	}
// 	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
// 	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp})
// }
