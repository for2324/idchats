package ens

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbEns "Open_IM/pkg/proto/ens"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary		预约
// @Description	预约
// @Tags			域名相关
// @ID				Appointment
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.AppointmentReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.AppointmentResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/ens/appointment [post]
func Appointment(c *gin.Context) {
	params := api.AppointmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	var UserId string
	ok, UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, "AppointmentReq args ", params)
	if params.Cancel {
		var err error
		if params.Ens != "" {
			err = imdb.CancelAppointmentEnsName(UserId, params.Ens)
		} else {
			err = imdb.CancelAppointment(UserId)
		}
		if err != nil {
			log.NewError(params.OperationID, "CancelAppointment failed ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
			return
		}
	} else {
		if len(params.Ens) >= 256 {
			c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": "domain can't more then 256 character"})
			return
		}
		err := imdb.Appointment(UserId, params.Ens)
		if err != nil {
			// if errors.Is(constant.ErrDoubleAppointment, err) {
			// 	log.NewError(params.OperationID, "Appointment failed ", err.Error())
			// 	c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": err.Error()})
			// 	return
			// }

			if errors.Is(constant.ErrHaveBeenAppointment, err) {
				log.NewError(params.OperationID, "Appointment failed ", err.Error())
				c.JSON(http.StatusOK, gin.H{"errCode": 409, "errMsg": constant.ErrHaveBeenAppointment})
				return
			}
			log.NewError(params.OperationID, "Appointment failed ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": constant.ErrHaveBeenAppointment})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "success"})
}

// @Summary		域名预约列表
// @Description	域名预约列表
// @Tags			域名相关
// @ID				AppointmentList
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.AppointmentListReq	true	"searchType:mine为我的预约列表 <br> pageIndex为页码（从0开始） <br> pageSize为每页数量"
// @Produce		json
// @Success		0	{object}	api.AppointmentListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/ens/appointment_list [post]
func AppointmentList(c *gin.Context) {
	params := api.AppointmentListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "AppointmentListReq args ", params)
	var ensList []api.AppointmentUserInfo
	var err error
	if params.SearchType == "mine" {
		var ok bool
		var errInfo string
		var UserId string
		ok, UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
		if !ok {
			errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
			log.NewError(params.OperationID, errMsg)
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
			return
		}
		ensList, err = imdb.MyAppointmentList(UserId, params.PageIndex, params.PageSize)
	} else {
		ensList, err = imdb.AppointmentList(params.PageIndex, params.PageSize)
	}
	for i := range ensList {
		ensList[i].CreateTime = ensList[i].CreatedAt.Unix()
	}

	if err != nil {
		log.NewError(params.OperationID, "AppointmentList failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "success", "data": ensList})
}

// @Summary		域名是否被预约
// @Description	域名是否被预约
// @Tags			域名相关
// @ID				HasAppointment
// @Accept			json
// @Param			req		body	api.HasAppointmentReq	true	"<br>"
// @Produce		json
// @Success		0	{object}	api.HasAppointmentResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/ens/has_appointment [post]
func HasAppointment(c *gin.Context) {
	params := api.HasAppointmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "HasAppointmentReq args ", params)
	hasAppointment, err := imdb.HasAppointment(params.Ens)
	if err != nil {
		log.NewError(params.OperationID, "HasAppointment failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "success", "data": hasAppointment})
}

// @Summary		创建注册ens订单
// @Description	创建注册ens订单
// @Tags			域名相关
// @ID				CreateRegisterEnsOrder
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.CreateRegisterEnsOrderReq	true	"创建注册ens订单"
// @Produce		json
// @Success		0	{object}	api.CreateRegisterEnsOrderResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/ens/create_register_ens_order [post]
func CreateRegisterEnsOrder(c *gin.Context) {
	params := api.CreateRegisterEnsOrderReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	// 校验 ensName 不能包含大写
	if strings.ToLower(params.EnsName) != params.EnsName {
		log.NewError(params.OperationID, "CreateRegisterEnsOrder failed ", "ensName must be lowercase")
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "ensName must be lowercase"})
		return
	}
	log.NewInfo(params.OperationID, utils.GetSelfFuncName(), " args ", params)
	var reqPb pbEns.CreateRegisterEnsOrderReq
	var ok bool
	var errInfo string
	ok, reqPb.UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	utils.CopyStructFields(&reqPb, params)
	etcdConn := getcdv3.GetDefaultConn(
		config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImEns,
		reqPb.OperationID,
	)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbEns.NewEnsServiceClient(etcdConn)
	respPb, err := client.CreateRegisterEnsOrder(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if respPb.CommonResp.ErrCode != 0 {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed ", respPb.CommonResp.ErrMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": respPb.CommonResp.ErrMsg})
		return
	}
	ensOrderInfo := api.EnsOrderInfo{}
	utils.CopyStructFields(&ensOrderInfo.Order, respPb.EnsOrderInfo)
	utils.CopyStructFields(&ensOrderInfo.PayInfo, respPb.ScanTaskInfo)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "success", "data": ensOrderInfo})
}

// @Summary		获取ens订单信息
// @Description	获取ens订单信息
// @Tags			域名相关
// @ID				GetEnsOrderInfo
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.GetEnsOrderInfoReq	true	"获取ens订单信息"
// @Produce		json
// @Success		0	{object}	api.GetEnsOrderInfoResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/ens/get_ens_order_info [post]
func GetEnsOrderInfo(c *gin.Context) {
	params := api.GetEnsOrderInfoReq{}
	apiResp := api.GetEnsOrderInfoResp{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, utils.GetSelfFuncName(), " args ", params)
	var reqPb pbEns.GetEnsOrderInfoReq
	var ok bool
	var errInfo string
	ok, reqPb.UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	utils.CopyStructFields(&reqPb, params)
	etcdConn := getcdv3.GetDefaultConn(
		config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImEns,
		reqPb.OperationID,
	)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbEns.NewEnsServiceClient(etcdConn)
	respPb, err := client.GetEnsOrderInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if respPb.CommonResp.ErrCode != 0 {
		errMsg := reqPb.OperationID + " GetEnsOrderInfo failed " + respPb.CommonResp.ErrMsg
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": errMsg})
		return
	}
	utils.CopyStructFields(&apiResp.Data.Order, respPb.EnsOrderInfo)
	utils.CopyStructFields(&apiResp.Data.PayInfo, respPb.ScanTaskInfo)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "success", "data": apiResp.Data})
}
