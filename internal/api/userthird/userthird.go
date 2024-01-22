package userthird

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/web3pub"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
)

// ThirdBindPostFunc
// @Summary		获取sign
// @Description	获取sign
// @Tags		请求第三方授权授权
// @ID			ThirdBindPostFunc
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param		req		body	api.GetUserThirdReq{}	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetUserThirdReq
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/thirdSns/postSign [post]
func ThirdBindPostFunc(c *gin.Context) {
	params := api.GetUserThirdReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "ThirdBindPostFunc req: ", params)
	var ok bool
	var errInfo string
	var uid string
	ok, uid, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	response := new(api.GetUserThirdRsp)

	switch params.ThirdString {
	case "twitter":
		response.Data = fmt.Sprintf(`BiuBiuID wants you address %s to verify by Twitter ,@%s Now Code Is:%s`,
			uid, config.Config.OfficialTwitter, randomdata.RandStringRunes(8))
	}
	err := db.DB.SetAccountThirdPlatformString(uid, params.ThirdString, response.Data, config.Config.Demo.CodeTTL*5)
	if err != nil {
		log.NewError(params.OperationID, "set redis error", uid, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
		return
	}
	log.Debug(params.OperationID, "sign text is ", response.Data, uid)
	c.JSON(http.StatusOK, response)
}

// VerifyThirdPlatformSign
// @Summary		第三方平台验签接口
// @Description	第三方平台验签接口
// @Tags			第三方平台授权
// @ID				VerifyThirdPlatformSign
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.VerifyThirdStringReq{}	true	"请求体 </br>  ThirdString 第三方平台：[ twitter,weibo,facebook]"
// @Produce		json
// @Success		0	{object}	api.GetUsersInfoResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/thirdSns/verifyingSign [post]
func VerifyThirdPlatformSign(c *gin.Context) {
	params := api.VerifyThirdStringReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "VerifyThirdPlatformSign req: ", params)
	var ok bool
	var errInfo string
	var uid string
	ok, uid, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImWeb3Js, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	switch params.ThirdString {
	case "twitter":
		{
			req := new(rpc.ThirdPlatformTwitterReq)
			req.Userid = uid
			req.OperatorID = params.OperationID
			req.Username = params.UserId
			if params.PostUrl != "" {
				compileRegex := regexp.MustCompile(`\.com\/(.*)\/status`)
				matchArr := compileRegex.FindStringSubmatch(params.PostUrl)
				if len(matchArr) < 2 {
					c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrArgs.ErrCode, "errMsg": constant.ErrArgs.ErrMsg})
					return
				}
				req.Username = matchArr[1]
			}
			req.Nonce, _ = db.DB.GetAccountThirdPlatformString(uid, params.ThirdString)
			client := rpc.NewWeb3PubClient(etcdConn)
			RpcResp, err := client.GetTwitterTimeLine(context.Background(), req)
			if err != nil {
				log.NewError(params.OperationID, "VerifyThirdPlatformSign failed ", err.Error(), req.String())
				c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
				return
			}
			resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
			resp.Data = jsonData.JsonDataList(resp.UserInfoList)
			log.NewInfo(params.OperationID, "VerifyThirdPlatformSign api return ", resp)
			c.JSON(http.StatusOK, resp)
		}
	}
}

// GetUserAuthorizedThirdPlatformList
// @Summary		获取用户已经授权的平台列表
// @Description	获取用户已经授权的平台列表
// @Tags			第三方平台授权
// @ID				GetUserAuthorizedThirdPlatformList
// @Accept			json
// @Param			token	header	string										true	"im token"
// @Param			req		body	api.GetUserAuthorizedThirdPlatformListReq{}	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UserAuthorizedThirdPlatformListReq
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/thirdSns/get_user_auth_platform [post]
//func GetUserAuthorizedThirdPlatformList(c *gin.Context) {
//	params := api.GetUserAuthorizedThirdPlatformListReq{}
//	if err := c.BindJSON(&params); err != nil {
//		log.NewError("0", "BindJSON failed ", err.Error())
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
//		return
//	}
//	var ok bool
//	var errInfo string
//	var uid string
//	ok, uid, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
//	if !ok {
//		errMsg := "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
//		log.NewError(params.OperationID, errMsg)
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
//		return
//	}
//	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
//		strings.Join(config.Config.Etcd.EtcdAddr, ","),
//		config.Config.RpcRegisterName.OpenImWeb3Js, params.OperationID)
//	if etcdConn == nil {
//		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
//		log.NewError(params.OperationID, errMsg)
//		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
//		return
//	}
//
//	req := new(rpc.GetUserAuthorizedThirdPlatformListReq)
//	req.Userid = uid
//	client := rpc.NewWeb3PubClient(etcdConn)
//	RpcResp, err := client.GetUserAuthorizedThirdPlatformList(context.Background(), req)
//	if err != nil {
//		log.NewError(params.OperationID, "VerifyThirdPlatformSign failed ", err.Error(), req.String())
//		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
//		return
//	}
//	resp := api.UserAuthorizedThirdPlatformListReq{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
//	resp.PlatFormList = make([]api.PlatForm, 0)
//	for _, v := range RpcResp.PlatFormList {
//		model := new(api.PlatForm)
//		copier.Copy(&model, &v)
//		resp.PlatFormList = append(resp.PlatFormList, *model)
//	}
//
//	log.NewInfo(params.OperationID, "VerifyThirdPlatformSign api return ", resp)
//	c.JSON(http.StatusOK, resp)
//}
