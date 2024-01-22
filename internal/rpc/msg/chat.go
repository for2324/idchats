package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	pbChat "Open_IM/pkg/proto/msg"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"time"

	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_msg_model"
)

func (rpc *rpcChat) ClearMsg(_ context.Context, req *pbChat.ClearMsgReq) (*pbChat.ClearMsgResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc req: ", req.String())
	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.UserID) {
		errMsg := "No permission" + req.OpUserID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "CleanUpOneUserAllMsgFromRedis args", req.UserID)
	err := db.DB.CleanUpOneUserAllMsgFromRedis(req.UserID, req.OperationID)
	if err != nil {
		errMsg := "CleanUpOneUserAllMsgFromRedis failed " + err.Error() + req.OperationID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "CleanUpUserMsgFromMongo args", req.UserID)
	err = db.DB.CleanUpUserMsgFromMongo(req.UserID, req.OperationID)
	if err != nil {
		errMsg := "CleanUpUserMsgFromMongo failed " + err.Error() + req.OperationID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}

	resp := pbChat.ClearMsgResp{ErrCode: 0}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return &resp, nil
}

func (rpc *rpcChat) SetMsgMinSeq(_ context.Context, req *pbChat.SetMsgMinSeqReq) (*pbChat.SetMsgMinSeqResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc req: ", req.String())
	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.UserID) {
		errMsg := "No permission" + req.OpUserID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	if req.GroupID == "" {
		err := db.DB.SetUserMinSeq(req.UserID, req.MinSeq)
		if err != nil {
			errMsg := "SetUserMinSeq failed " + err.Error() + req.OperationID + req.UserID + utils.Uint32ToString(req.MinSeq)
			log.Error(req.OperationID, errMsg)
			return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
		}
		return &pbChat.SetMsgMinSeqResp{}, nil
	}
	err := db.DB.SetGroupUserMinSeq(req.GroupID, req.UserID, uint64(req.MinSeq))
	if err != nil {
		errMsg := "SetGroupUserMinSeq failed " + err.Error() + req.OperationID + req.GroupID + req.UserID + utils.Uint32ToString(req.MinSeq)
		log.Error(req.OperationID, errMsg)
		return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	return &pbChat.SetMsgMinSeqResp{}, nil
}

func (s *rpcChat) GetSingleChatHistoryMessageList(ctx context.Context, req *pbChat.GetSingleChatHistoryMessageListReq) (*pbChat.GetSingleChatHistoryMessageListResp, error) {
	t := time.Now()
	var sourceID string
	var channelID string
	var startTime time.Time
	var err error
	var notStartTime bool
	if req.StartClientMsgID == "" {
		notStartTime = true
	} else {
		msg := db.ChatLog{
			ClientMsgID: req.StartClientMsgID,
		}
		resultmsg, err := imdb.GetMessageFromMysqlDb(&msg)
		if err == nil {
			startTime = resultmsg.SendTime
		}
	}
	channelID = req.ChannelID
	log.Debug(req.OperationID, "Assembly parameters cost time", time.Since(t))
	t = time.Now()
	log.Info(req.OperationID, "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)

	isReverse := req.IsReverse
	var list []*db.ChatLog
	if notStartTime {
		list, err = imdb.GetSingleChatMessageListNoTimeControllerFromMysqlDb(
			req.SendUserId, req.RecvUserId, int(req.Count), isReverse, channelID,
		)
	} else {
		list, err = imdb.GetSingleChatMessageListControllerFromMysqlDb(
			req.SendUserId, req.RecvUserId, int(req.Count),
			startTime, isReverse, channelID,
		)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "get message list count:", len(list))
	log.Debug(req.OperationID, "db cost time", time.Since(t))
	t = time.Now()
	resultResp := new(pbChat.GetSingleChatHistoryMessageListResp)
	resultResp.CommonResp = new(pbChat.CommonResp)
	for _, v := range list {
		temp := open_im_sdk.MsgData{}
		temp.ClientMsgID = v.ClientMsgID
		temp.ServerMsgID = v.ServerMsgID
		temp.CreateTime = v.CreateTime.UnixMilli()
		temp.SendTime = v.SendTime.UnixMilli()
		temp.SessionType = v.SessionType
		temp.SendID = v.SendID
		temp.RecvID = v.RecvID
		temp.MsgFrom = v.MsgFrom
		temp.ContentType = v.ContentType
		temp.SenderPlatformID = v.SenderPlatformID
		temp.SenderNickname = v.SenderNickname
		temp.SenderFaceURL = v.SenderFaceURL
		temp.Content = utils.String2bytes(v.Content)
		temp.Status = v.Status
		temp.Ex = v.Ex
		resultResp.Message = append(resultResp.Message, &temp)
	}
	log.Debug(req.OperationID, "unmarshal cost time", time.Since(t))
	log.Debug(req.OperationID, "sort cost time", time.Since(t))
	return resultResp, err
}
