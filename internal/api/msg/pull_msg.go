package msg

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/proto/msg"
	pbMsg "Open_IM/pkg/proto/msg"
	rpc "Open_IM/pkg/proto/msg"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"

	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type paramsUserPullMsg struct {
	ReqIdentifier *int   `json:"reqIdentifier" binding:"required"`
	SendID        string `json:"sendID" binding:"required"`
	OperationID   string `json:"operationID" binding:"required"`
	Data          struct {
		SeqBegin *int64 `json:"seqBegin" binding:"required"`
		SeqEnd   *int64 `json:"seqEnd" binding:"required"`
	}
}

type paramsUserPullMsgBySeqList struct {
	ReqIdentifier int      `json:"reqIdentifier" binding:"required"`
	SendID        string   `json:"sendID" binding:"required"`
	OperationID   string   `json:"operationID" binding:"required"`
	SeqList       []uint32 `json:"seqList"`
}

func PullMsgBySeqList(c *gin.Context) {
	params := paramsUserPullMsgBySeqList{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	token := c.Request.Header.Get("token")
	if ok, err := token_verify.VerifyToken(token, params.SendID); !ok {
		if err != nil {
			log.NewError(params.OperationID, utils.GetSelfFuncName(), err.Error(), token, params.SendID)
		}
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	}
	pbData := open_im_sdk.PullMessageBySeqListReq{}
	pbData.UserID = params.SendID
	pbData.OperationID = params.OperationID
	pbData.SeqList = params.SeqList

	grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, pbData.OperationID)
	if grpcConn == nil {
		errMsg := pbData.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(pbData.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	msgClient := msg.NewMsgClient(grpcConn)
	reply, err := msgClient.PullMessageBySeqList(context.Background(), &pbData)
	if err != nil {
		log.Error(pbData.OperationID, "PullMessageBySeqList error", err.Error())
		return
	}
	log.NewInfo(pbData.OperationID, "rpc call success to PullMessageBySeqList", reply.String(), len(reply.List))
	c.JSON(http.StatusOK, gin.H{
		"errCode":       reply.ErrCode,
		"errMsg":        reply.ErrMsg,
		"reqIdentifier": params.ReqIdentifier,
		"data":          reply.List,
	})
}

// GetSingleChatHistoryMessageList
// @Summary		获取私聊聊天信息的历史消息
// @Description	获取私聊聊天信息的历史消息
// @Tags			消息相关
// @ID				GetSingleChatHistoryMessageList
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.GetSingleChatHistoryMessageListReq	true	"fromUserID为要获取的用户ID"
// @Produce		json
// @Success		0	{object}	api.GetHistoryMessageListResp{data=[]api.MsgStruct}
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/msg/get_single_chat_history_message_list [post]
func GetSingleChatHistoryMessageList(c *gin.Context) {
	params := api.GetSingleChatHistoryMessageListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbMsg.GetSingleChatHistoryMessageListReq{
		RecvUserId: params.UserID,
	}
	utils.CopyStructFields(req, params)

	var ok bool
	var errInfo string
	ok, req.SendUserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(
		config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImMsgName, req.OperationID,
	)
	if etcdConn == nil {
		errMsg := req.OperationID + " getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewMsgClient(etcdConn)
	RpcResp, err := client.GetSingleChatHistoryMessageList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, " SetMsgMinSeq failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if err != nil {
		log.NewError(req.OperationID, "GetGroupAllHistoryMessageList failed  ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	groupHistoryMessageListResp := api.GetHistoryMessageListResp{
		CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg},
	}
	groupHistoryMessageListResp.Data = make([]*api.MsgStruct, 0)
	for _, value := range RpcResp.Message {
		apimsgstruct := new(api.MsgStruct)
		utils.CopyStructFields(apimsgstruct, value)
		msgHandleByContentType(apimsgstruct)
		groupHistoryMessageListResp.Data = append(groupHistoryMessageListResp.Data, apimsgstruct)
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api return ", len(groupHistoryMessageListResp.Data))
	c.JSON(http.StatusOK, groupHistoryMessageListResp)
}

func msgHandleByContentType(msg *api.MsgStruct) (err error) {
	_ = utils.JsonStringToStruct(msg.AttachedInfo, &msg.AttachedInfoElem)
	switch msg.ContentType {
	case constant.Picture:
		msg.PictureElem = new(api.PictureElem)
		err = utils.JsonStringToStruct(msg.Content, msg.PictureElem)
	case constant.Voice:
		msg.SoundElem = new(api.SoundElem)
		err = utils.JsonStringToStruct(msg.Content, msg.SoundElem)
	case constant.Video:
		msg.VideoElem = new(api.VideoElem)
		err = utils.JsonStringToStruct(msg.Content, msg.VideoElem)
	case constant.File:
		msg.FileElem = new(api.FileElem)
		err = utils.JsonStringToStruct(msg.Content, msg.FileElem)
	case constant.AdvancedText:
		msg.MessageEntityElem = new(api.MessageEntityElem)
		err = utils.JsonStringToStruct(msg.Content, msg.MessageEntityElem)
	case constant.Location:
		msg.LocationElem = new(api.LocationElem)
		err = utils.JsonStringToStruct(msg.Content, &msg.LocationElem)
	case constant.Custom:
		msg.CustomElem = new(api.CustomElem)
		err = utils.JsonStringToStruct(msg.Content, &msg.CustomElem)
	case constant.Quote:
		msg.QuoteElem = new(api.QuoteElem)
		err = utils.JsonStringToStruct(msg.Content, msg.QuoteElem)
	case constant.Merger:
		msg.MergeElem = new(api.MergeElem)
		err = utils.JsonStringToStruct(msg.Content, &msg.MergeElem)
	}

	return utils.Wrap(err, "")
}
