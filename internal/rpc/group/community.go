package group

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_msg_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbGroup "Open_IM/pkg/proto/group"
	pbChat "Open_IM/pkg/proto/msg"
	pbOrder "Open_IM/pkg/proto/order"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func (s *groupServer) GetHotCommunityBanner(ctx context.Context, req *pbGroup.GetHotCommunityReq) (*pbGroup.GetHotSpaceBannerResp, error) {
	resp := pbGroup.GetHotSpaceBannerResp{
		ErrCode:       0,
		ErrMsg:        "",
		HotBannerInfo: make([]*pbGroup.HotBannerInfo, 0),
	}
	gethotbannerl, err := imdb.GetBannelGroupInfo()
	if err == nil {
		for _, value := range gethotbannerl {
			tempdata := new(pbGroup.HotBannerInfo)
			utils.CopyStructFields(tempdata, value)
			resp.HotBannerInfo = append(resp.HotBannerInfo, tempdata)
		}
	}
	return &resp, nil
}

func (s *groupServer) GetCommunityRoleDetail(
	ctx context.Context, detail *pbGroup.GetCommunityRoleReqDetail,
) (*pbGroup.GetCommunityRoleRespDetail, error) {

	resultResp := new(pbGroup.GetCommunityRoleRespDetail)
	resultResp.CommonResp = new(pbGroup.CommonResp)
	_, err := rocksCache.GetGroupInfoFromCache(detail.GroupID) //群不存在 或者说已经解散掉了
	if err != nil {
		resultResp.CommonResp.ErrCode = 10001
		resultResp.CommonResp.ErrMsg = err.Error()
		return resultResp, nil
	}
	if detail == nil || detail.GroupID == "" {
		resultResp.CommonResp.ErrCode = 10001
		resultResp.CommonResp.ErrMsg = "detail is nil"
		return resultResp, nil
	}
	//不需要判断是还否是群成圆
	//if groupinfo.OpUserID != detail.OpUserID {
	//	resultResp.CommonResp.ErrCode = 10001
	//	resultResp.CommonResp.ErrMsg = "not group creator can't create role" + groupinfo.OpUserID + ">" + detail.OpUserID
	//	return resultResp, nil
	//}
	// 查询标签分配的列表
	loadDataFromDB, err := imdb.GetGroupMemberTagList(detail.GroupID, detail.TokenID)
	if err == nil {
		for _, value := range loadDataFromDB {
			userInfodata, err := rocksCache.GetUserBaseInfoFromCache(value.UserID)
			if err != nil {
				//用户xx 不存在在系统之中
				//	log.NewError(detail.OperationID, utils.GetSelfFuncName(), detail.GroupID, value.UserID, err.Error())
				//		continue
				memberNodeUserInfo := new(open_im_sdk.UserInfo)
				memberNodeUserInfo.UserID = value.UserID
				memberNodeUserInfo.FaceURL = ""
				memberNodeUserInfo.Nickname = value.UserID
				resultResp.CommunityRoleInfoDetail = append(resultResp.CommunityRoleInfoDetail, &pbGroup.CommunityRoleInfoDetail{
					UserInfo: memberNodeUserInfo,
					Amount:   value.Amount,
				})
			} else {
				memberNodeUserInfo := new(open_im_sdk.UserInfo)
				utils.CopyStructFields(memberNodeUserInfo, userInfodata)
				resultResp.CommunityRoleInfoDetail = append(resultResp.CommunityRoleInfoDetail, &pbGroup.CommunityRoleInfoDetail{
					UserInfo: memberNodeUserInfo,
					Amount:   value.Amount,
				})
			}

		}
		return resultResp, nil
	} else {
		resultResp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resultResp.CommonResp.ErrMsg = "data error:" + err.Error()
		return resultResp, nil
	}
}

func (s *groupServer) GetUserJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (*pbGroup.GetJoinedGroupListResp, error) {
	//TODO implement me
	log.NewInfo(req.OperationID, "GetUserJoinedGroupList, args ", req.String())
	//if !token_verify.CheckAccess(req.OpUserID, req.FromUserID) {
	//	log.NewError(req.OperationID, "CheckAccess false ", req.OpUserID, req.FromUserID)
	//	return &pbGroup.GetJoinedGroupListResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	//}

	joinedGroupList, err := rocksCache.GetJoinedGroupIDListFromCache(req.FromUserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserJoinedGroupIDListFromCache failed", err.Error(), req.FromUserID)
		return &pbGroup.GetJoinedGroupListResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "UserJoinedGroupList: ", joinedGroupList)
	var resp pbGroup.GetJoinedGroupListResp
	for _, v := range joinedGroupList {
		var groupNode open_im_sdk.GroupInfo
		num, err := rocksCache.GetGroupMemberNumFromCache(v)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), v)
			continue
		}
		owner, err2 := imdb.GetGroupOwnerInfoByGroupID(v)
		if err2 != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err2.Error(), v)
			continue
		}
		group, err := rocksCache.GetGroupInfoFromCache(v)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), v)
			continue
		}
		if group.GroupType == constant.SuperGroup {
			continue
		}
		if group.GroupType != constant.WorkingGroup {
			continue
		}
		if group.Status == constant.GroupStatusDismissed {
			log.NewError(req.OperationID, "constant.GroupStatusDismissed ", group)
			continue
		}
		utils.CopyStructFields(&groupNode, group)
		groupNode.CreateTime = uint32(group.CreateTime.Unix())
		groupNode.NotificationUpdateTime = uint32(group.NotificationUpdateTime.Unix())
		if group.NotificationUpdateTime.Unix() < 0 {
			groupNode.NotificationUpdateTime = 0
		}

		groupNode.MemberCount = uint32(num)
		groupNode.OwnerUserID = owner.UserID
		resp.GroupList = append(resp.GroupList, &groupNode)
		log.NewDebug(req.OperationID, "joinedGroup ", groupNode)
	}
	log.NewInfo(req.OperationID, "GetJoinedGroupList rpc return ", resp.String())
	return &resp, nil
}
func (s *groupServer) GetGroupHistoryMessageList(ctx context.Context, req *pbGroup.GetHistoryMessageListParamsReq) (*pbGroup.GetHistoryMessageListParamsResp, error) {
	t := time.Now()
	var sourceID string
	var channelID string
	var startTime time.Time
	sessionType := constant.SuperGroupChatType
	var err error
	var notStartTime bool
	if req.UserID == "" {
		sourceID = req.GroupID
		channelID = req.ChannelID
		if req.StartClientMsgID == "" {
			notStartTime = true
		} else {
			msg := db.ChatLog{
				ClientMsgID: req.StartClientMsgID,
			}
			resultmsg, err := im_mysql_msg_model.GetMessageFromMysqlDb(&msg)
			if err == nil {
				startTime = resultmsg.SendTime
			}
		}
	} else {
		return &pbGroup.GetHistoryMessageListParamsResp{
			CommonResp: &pbGroup.CommonResp{
				ErrCode: 50001,
				ErrMsg:  "only search community message",
			},
			Message: nil,
		}, nil
	}
	channelID = req.ChannelID
	log.Debug(req.OperationID, "Assembly parameters cost time", time.Since(t))
	t = time.Now()
	log.Info(req.OperationID, "sourceID:", sourceID, "startTime:", startTime, "count:", req.Count, "not start_time", notStartTime)

	sessionType = constant.SuperGroupChatType
	isReverse := req.IsReverse
	var list []*db.ChatLog
	if notStartTime {
		list, err = im_mysql_msg_model.GetGroupMessageListNoTimeControllerFromMysqlDb(sourceID, sessionType, int(req.Count), isReverse, channelID)
	} else {
		list, err = im_mysql_msg_model.GetGroupMessageListControllerFromMysqlDb(sourceID, sessionType, int(req.Count),
			startTime, isReverse, channelID)
	}
	fmt.Println("获取到的消息总数为：", len(list))
	log.Debug(req.OperationID, "db cost time", time.Since(t))
	t = time.Now()
	resultResp := new(pbGroup.GetHistoryMessageListParamsResp)
	resultResp.CommonResp = new(pbGroup.CommonResp)
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
		switch sessionType {
		case constant.GroupChatType:
			fallthrough
		case constant.SuperGroupChatType:
			temp.GroupID = temp.RecvID
			temp.ChannelID = v.ChannelID
		}
		resultResp.Message = append(resultResp.Message, &temp)
	}
	log.Debug(req.OperationID, "unmarshal cost time", time.Since(t))
	log.Debug(req.OperationID, "sort cost time", time.Since(t))
	return resultResp, err
}

func (s *groupServer) CreateCommunityRole(ctx context.Context, req *pbGroup.CreateCommunityRoleReq) (*pbGroup.CreateCommunityRoleResp, error) {
	//创建角色
	//构建数据库。
	resultResp := new(pbGroup.CreateCommunityRoleResp)
	resultResp.CommonResp = new(pbGroup.CommonResp)
	groupinfo, err := rocksCache.GetGroupInfoFromCache(req.GroupID) //群不存在 或者说已经解散掉了
	if err != nil || groupinfo.CreatorUserID != req.OpUserID {
		resultResp.CommonResp.ErrCode = 10001
		resultResp.CommonResp.ErrMsg = "not group creator can't create role" + groupinfo.CreatorUserID + ">" + req.OpUserID
		return resultResp, nil
	}

	roleid := utils.Md5(req.OperationID + strconv.FormatInt(time.Now().UnixNano(), 10))
	bi := big.NewInt(0)
	bi.SetString(roleid[0:8], 16)
	roleid = bi.String()
	if err = imdb.CreateGroupRoleInformation(&db.CommunityChannelRole{
		GroupID:   req.GroupID,
		RoleID:    roleid,
		RoleTitle: req.RoleTitle,
		RoleIPfs:  req.RoleIPfs,
	}); err != nil {
		resultResp.CommonResp.ErrCode = 10002
		resultResp.CommonResp.ErrMsg = err.Error()
		return resultResp, nil
	}

	resultResp.ReBackOrderID = roleid
	return resultResp, nil
}

func (s *groupServer) GetCommunityRole(ctx context.Context, req *pbGroup.GetCommunityRoleReq) (
	*pbGroup.GetCommunityRoleResp, error) {
	resultResp := new(pbGroup.GetCommunityRoleResp)
	resultResp.CommonResp = new(pbGroup.CommonResp)
	groupinfo, err := rocksCache.GetGroupInfoFromCache(req.GroupID) //群不存在 或者说已经解散掉了

	resultData, err := imdb.GetGroupRoleInformation(groupinfo.GroupID)
	if err != nil {
		resultResp.CommonResp.ErrCode = 10002
		resultResp.CommonResp.ErrMsg = err.Error()
		return resultResp, nil
	}
	for _, value := range resultData {
		resuldattotla, err := imdb.GetTotalNft1155BurnAndTransfer(value)
		if err != nil {
			continue
		}
		tempdata := new(pbGroup.CommunityRoleInfo)
		dtAmount, _ := decimal.NewFromString(value.TokenAmount)
		dtSubBurn, _ := decimal.NewFromString(resuldattotla[0].TotalAmount)
		dtSub, _ := decimal.NewFromString(resuldattotla[1].TotalAmount)
		utils.CopyStructFields(tempdata, value)
		tempdata.TokenBurn = dtAmount.Sub(dtSubBurn).String()
		tempdata.TokenSub = dtSub.String()
		resultResp.CommunityRoleInfo = append(resultResp.CommunityRoleInfo, tempdata)
	}
	return resultResp, nil
}

func (s *groupServer) BindCommunityChannelRole(ctx context.Context, req *pbGroup.OperatorCommunityChannelRoleReq) (*pbGroup.OperatorCommunityChannelRoleResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s *groupServer) GetUserRoleTagInfo(ctx context.Context, req *pbGroup.OperatorCommunityChannelRoleReq) (*pbGroup.UserRoleTagListRsp, error) {
	respdata := new(pbGroup.UserRoleTagListRsp)
	respdata.CommonResp = new(pbGroup.CommonResp)
	var err error
	//获取用户总共拥有的标签内容
	respdata.RoleIpfs, err = imdb.GetGroupMemberTagListByUserID(req.OpUserID, req.GroupID)
	if err != nil {
		respdata.CommonResp.ErrCode = constant.ErrDB.ErrCode
		respdata.CommonResp.ErrMsg = err.Error()
	}
	return respdata, nil
}
func (s *groupServer) PublishAnnounceMoment(ctx context.Context, req *pbGroup.PublishAnnouncementReq) (resp *pbGroup.PublishAnnouncementResp, err error) {
	resp = new(pbGroup.PublishAnnouncementResp)
	resp.CommonResp = new(pbGroup.CommonResp)
	var respSpaceArticleInfo *pbGroup.SpaceArticleIDResp
	if req.ArticleType == nil && req.ArticleID == nil {
		articleInfo, err := imdb.CreateAnnouncement(req.CreatorUserID,
			req.AnnouncementElem.AnnouncementUrl,
			req.AnnouncementTitle, req.AnnouncementSummary,
			req.OpUserID, req.IsGlobal)
		if err != nil {
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		log.NewInfo(req.OperationID, "当前的链路跟踪的operationID", req.OperationID)
		respSpaceArticleInfo, err = s.InputSpaceArticleList(ctx, &pbGroup.PutSpaceArticleID{
			OperationID:   req.OperationID,
			ArticleType:   "announce",
			ArticleID:     utils.Int64ToString(articleInfo.ArticleID),
			OpUserID:      req.OpUserID,
			CreatorUserID: req.CreatorUserID,
			EndTime:       nil,
			IsGlobal:      req.IsGlobal,
		})
		if err != nil {
			return &pbGroup.PublishAnnouncementResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	} else {
		//查询是否是自己的空间的文章， 如果不是那么要查询是否属于别人推送给我的文章
		var creatorIdFromMySpace, creatorIdFromMyPush string
		if creatorIdFromMySpace, err = imdb.GetSpaceArticleByArticleIdAndArticleType(utils.UInt64ToString(req.ArticleID.Value),
			req.ArticleType.Value); err != nil {
			log.NewInfo(req.OperationID, "无法从自己的文章中找到自己人1")
		}
		if creatorIdFromMySpace == "" {
			if creatorIdFromMyPush, err = imdb.GetPersonalSpaceArticleByArticleIdAndArticleType(req.OpUserID,
				utils.UInt64ToString(req.ArticleID.Value), req.ArticleType.Value); err != nil {
				log.NewInfo(req.OperationID, "无法从自己的文章中找到自己人2")
			}
		}
		var creatorUserId string
		if creatorIdFromMySpace == "" && creatorIdFromMyPush == "" {
			return &pbGroup.PublishAnnouncementResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode,
				ErrMsg: "你的文章可能还没发表"}}, nil
		} else {
			creatorUserId = creatorIdFromMySpace
			if creatorUserId == "" {
				creatorUserId = creatorIdFromMyPush
			}
		}
		respSpaceArticleInfo, err = s.InputSpaceArticleList(ctx, &pbGroup.PutSpaceArticleID{
			OperationID:   req.OperationID,
			ArticleType:   req.ArticleType.Value,
			ArticleID:     utils.UInt64ToString(req.ArticleID.Value),
			OpUserID:      req.OpUserID,
			CreatorUserID: creatorUserId,
			EndTime:       nil,
			IsGlobal:      req.IsGlobal,
		})
		if err != nil {
			return &pbGroup.PublishAnnouncementResp{CommonResp: &pbGroup.CommonResp{
				ErrCode: constant.ErrDB.ErrCode,
				ErrMsg:  err.Error()}}, nil
		}
	}
	resp.CommonResp.ErrCode = 0
	resp.CommonResp.ErrMsg = "ok"
	if req.IsGlobal == 1 {

		if config.Config.SpaceArticle.PushUsdPrice > 0 {
			//构建订单
			resultData, err := s.CreatePushSpaceArticelOrder(ctx, &pbGroup.CreatePushSpaceArticelOrderReq{
				OperationID:    req.OperationID,
				SpaceArticleId: strconv.FormatInt(respSpaceArticleInfo.NewSpaceArticleID, 10),
				UserId:         req.OpUserID,
				TxnType:        req.TxnType,
			})
			if err != nil {
				return &pbGroup.PublishAnnouncementResp{CommonResp: &pbGroup.CommonResp{
					ErrCode: constant.ErrDB.ErrCode,
					ErrMsg:  constant.ErrDB.ErrMsg}}, nil
			} else {
				resp.PayInfo = new(pbOrder.ScanTaskInfo)
				utils.CopyStructFields(resp.PayInfo, resultData.PayInfo)
			}
		}

	}
	return resp, nil

}
func (s *groupServer) GetPublishAnnounceMomentListWithIdo(ctx context.Context, req *pbGroup.GetPublishAnnouncementWithIdoReq) (
	resultData *pbGroup.GetPersonalPublishAnnouncementWithIdoResp, err error) {
	resultData = new(pbGroup.GetPersonalPublishAnnouncementWithIdoResp)
	resultData.CommonResp = new(pbGroup.CommonResp)
	resultData.TotalCount, _ = imdb.GetArticleListCount(req.OpUserID, req.ID, 0, 0)
	resultData.CurrentPage = int64(req.PageIndex)
	dataArray, err := imdb.GetArticleList(req.OpUserID, req.ID, int64(req.PageIndex), int64(req.PageSize))
	if err == nil && len(dataArray) > 0 {
		for _, value := range dataArray {
			tempValue := &pbGroup.PersonalGenerateArticle{
				GenerateArticle: &pbGroup.GenerateArticle{
					ArticleID:           value.ArticleID,
					CreatedAt:           value.CreatedAt.Format("2006-01-02 15:04:05"),
					UpdatedAt:           value.UpdatedAt.Format("2006-01-02 15:04:05"),
					DeletedAt:           value.DeletedAt.Format("2006-01-02 15:04:05"),
					OpUserID:            req.OpUserID, //获取自己的发送的文章
					GroupArticleId:      int32(value.GroupArticleID),
					CreatorUserId:       value.CreatorUserID,
					AnnouncementContent: value.AnnouncementContent,
					AnnouncementUrl:     value.AnnouncementUrl,
					LikeCount:           int32(value.LikeCount),
					RewordCount:         int32(value.RewordCount),
					IsGlobal:            value.IsGlobal,
					OrderId:             value.OrderID,
					Status:              value.Status,
					AnnouncementTitle:   value.AnnouncementTitle,
					AnnouncementSummary: value.AnnouncementSummary,
					ID:                  value.ID,
					ArticleType:         value.ArticleType,
					IsPin:               value.ArticleIsPin,
				},
				CreatorInfo: new(open_im_sdk.PublicUserInfo),
			}
			userInfo, err := rocksCache.GetUserBaseInfoFromCache(value.CreatorUserID)
			if err == nil {
				utils.CopyStructFields(tempValue.CreatorInfo, userInfo)
				linkTree, err := rocksCache.GetUserBaseInfoFromCacheUserLink(value.CreatorUserID)
				if err == nil {
					for _, value := range linkTree {
						tempData := new(open_im_sdk.LinkTreeMsgReq)
						utils.CopyStructFields(tempData, value)
						tempValue.CreatorInfo.LinkTree = append(tempValue.CreatorInfo.LinkTree, tempData)
					}
				}
			} else {

				tempValue.CreatorInfo = nil
			}

			resultData.PublishAnnounce = append(resultData.PublishAnnounce, tempValue)
		}
	}
	return resultData, err
}
func (s *groupServer) GetPublishAnnounceMomentList(ctx context.Context, req *pbGroup.GetPublishAnnouncementReq) (*pbGroup.GetPublishAnnouncementResp, error) {
	type AnnouncementArticleInfo struct {
		GroupID string
		Name    string
		FaceURL string
		IsLikes int32
		IsRead  int32
		db.AnnouncementArticle
	}

	if req.IsGlobal == 0 {
		var tempdata []*AnnouncementArticleInfo
		txdb := db.DB.MysqlDB.DefaultGormDB().Table("announcement_article ").
			Select(` groups.group_id,groups.name ,groups.face_url,COALESCE(announcement_article_logs.is_likes,0) as is_likes,COALESCE(announcement_article_logs.status,0) as is_read ,announcement_article.*`).
			Joins("LEFT JOIN groups on announcement_article.group_id= groups.group_id ").
			Joins("LEFT JOIN announcement_article_logs on announcement_article_logs.article_id = announcement_article.article_id and announcement_article_logs.user_id=?", req.CreatorUserID).
			Where("announcement_article.group_id=? and announcement_article.status=0", req.GroupID)
		if req.ArticleID != 0 {
			txdb = txdb.Where("announcement_article.article_id<=?", req.ArticleID)
		}
		err := txdb.Order("announcement_article.article_id desc").Limit(int(req.PageSize)).Offset(int(req.PageSize * req.PageIndex)).Find(&tempdata).Error
		resultpb := new(pbGroup.GetPublishAnnouncementResp)
		resultpb.CommonResp = new(pbGroup.CommonResp)
		log.NewInfo(req.OperationID, len(tempdata))
		if err != nil {
			resultpb.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resultpb.CommonResp.ErrMsg = err.Error()
			return resultpb, err
		} else {
			resultpb.CommonResp.ErrCode = 0
			resultpb.CommonResp.ErrMsg = ""
			resultpb.PublishAnnounce = make([]*pbGroup.AnnouncementInfo, 0)
			for _, value := range tempdata {
				resultpb.PublishAnnounce = append(resultpb.PublishAnnounce, &pbGroup.AnnouncementInfo{
					ArticleID:           value.ArticleID,
					CreatedAt:           value.CreatedAt.Unix(),
					UpdatedAt:           value.UpdatedAt.Unix(),
					CreatorUserID:       value.CreatorUserID,
					AnnouncementTitle:   value.AnnouncementTitle,
					AnnouncementSummary: value.AnnouncementSummary,
					AnnouncementContent: value.AnnouncementContent,
					AnnouncementUrl:     value.AnnouncementUrl,
					LikeCount:           value.LikeCount,
					RewordCount:         value.RewordCount,
					GroupID:             value.GroupID,
					GroupArticleID:      value.GroupArticleID,
					IsGlobal:            value.IsGlobal,
					GroupName:           value.Name,
					FaceUrl:             value.FaceURL,
					IsRead:              value.IsRead,
					IsLikes:             value.IsLikes,
				})
			}
		}
		return resultpb, nil
	} else if req.IsGlobal == 1 {
		var tempdata []*AnnouncementArticleInfo
		txdb := db.DB.MysqlDB.DefaultGormDB().Table("announcement_article ").
			Select(` groups.group_id,groups.name ,groups.face_url,COALESCE(announcement_article_logs.is_likes,0) as is_likes,COALESCE(announcement_article_logs.status,0) as is_read ,announcement_article.*`).
			Joins("LEFT JOIN groups on announcement_article.group_id=groups.group_id ").
			Joins("LEFT JOIN announcement_article_logs on announcement_article_logs.article_id = announcement_article.article_id and announcement_article_logs.user_id=?", req.CreatorUserID).
			Where("announcement_article.is_global=1  and announcement_article.status=0 and COALESCE(announcement_article_logs.status,0) < 2")
		if req.ArticleID != 0 {
			txdb = txdb.Where("announcement_article.article_id<=?", req.ArticleID)
		}
		err := txdb.Order("announcement_article.article_id desc").Limit(int(req.PageSize)).Offset(int(req.PageSize * req.PageIndex)).Find(&tempdata).Error
		resultpb := new(pbGroup.GetPublishAnnouncementResp)
		resultpb.CommonResp = new(pbGroup.CommonResp)
		log.NewInfo(req.OperationID, len(tempdata))
		if err != nil {
			resultpb.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resultpb.CommonResp.ErrMsg = err.Error()
			return resultpb, err
		} else {
			resultpb.CommonResp.ErrCode = 0
			resultpb.CommonResp.ErrMsg = ""
			resultpb.PublishAnnounce = make([]*pbGroup.AnnouncementInfo, 0)
			for _, value := range tempdata {
				resultpb.PublishAnnounce = append(resultpb.PublishAnnounce, &pbGroup.AnnouncementInfo{
					ArticleID:           value.ArticleID,
					CreatedAt:           value.CreatedAt.Unix(),
					UpdatedAt:           value.UpdatedAt.Unix(),
					CreatorUserID:       value.CreatorUserID,
					AnnouncementTitle:   value.AnnouncementTitle,
					AnnouncementSummary: value.AnnouncementSummary,
					AnnouncementContent: value.AnnouncementContent,
					AnnouncementUrl:     value.AnnouncementUrl,
					LikeCount:           value.LikeCount,
					RewordCount:         value.RewordCount,
					GroupID:             value.GroupID,
					GroupArticleID:      value.GroupArticleID,
					IsGlobal:            value.IsGlobal,
					GroupName:           value.Name,
					FaceUrl:             value.FaceURL,
					IsRead:              value.IsRead,
					IsLikes:             value.IsLikes,
				})
			}
		}
		return resultpb, nil
	}

	return nil, nil
}
func (s *groupServer) InputSpaceArticleList(ctx context.Context, reqid *pbGroup.PutSpaceArticleID) (result *pbGroup.SpaceArticleIDResp, err error) {
	// 将发布的内容展示到个人空间主页并提供status的服务

	log.NewInfo(reqid.OperationID, "当前的链路跟踪的operationID", reqid.OperationID)
	result = new(pbGroup.SpaceArticleIDResp)
	result.CommonResp = new(pbGroup.CommonResp)
	inPutNow := time.Now()
	if reqid.EndTime != nil {
		inPutNow = time.Unix(reqid.EndTime.Seconds, int64(reqid.EndTime.Nanos))
	}
	status := int8(1)
	if config.Config.SpaceArticle.PushUsdPrice > 0 {
		status = 0
	}
	newSpaceID, err := imdb.InputSpaceArticleID(reqid.ArticleID,
		reqid.ArticleType,
		reqid.OpUserID,
		reqid.CreatorUserID,
		inPutNow,
		inPutNow, status)
	if err != nil {
		result.CommonResp.ErrCode = constant.ErrDB.ErrCode
		result.CommonResp.ErrMsg = err.Error()
		return
	}
	if reqid.IsGlobal == 0 || config.Config.SpaceArticle.PushUsdPrice == 0 {
		pushMQ := &pbChat.NewPushActionMsgMq{
			OperationID: reqid.OperationID,
			PushMsg: &open_im_sdk.PushMessageToMailFromUserToFans{
				OperationID:       reqid.OperationID,
				ContentType:       reqid.ArticleType,
				ArticleID:         utils.StringToInt64(reqid.ArticleID),
				FromUserID:        reqid.OpUserID,
				FromArticleAuthor: reqid.CreatorUserID,
				IsGlobal:          reqid.IsGlobal,
			},
		}
		s.messagePushWrite.SendMessage(pushMQ, reqid.OpUserID, reqid.OperationID)
	}
	result.NewSpaceArticleID = newSpaceID
	return result, nil
}

func (s *groupServer) PinSpaceArticleList(ctx context.Context, reqid *pbGroup.PinSpaceArticleID) (result *pbGroup.SpaceArticleIDResp, err error) {
	// 将发布的内容展示到个人空间主页并提供status的服务
	result = new(pbGroup.SpaceArticleIDResp)
	result.CommonResp = new(pbGroup.CommonResp)
	err = imdb.PinSpaceArticleID(reqid.ID, reqid.UserID, reqid.IsPin)
	if err != nil {
		result.CommonResp.ErrCode = constant.ErrDB.ErrCode
		result.CommonResp.ErrMsg = err.Error()
		return
	}
	return result, nil
}

func (s *groupServer) DelSpaceArticleList(ctx context.Context, reqid *pbGroup.DelSpaceArticleID) (result *pbGroup.SpaceArticleIDResp, err error) {
	//TODO implement me
	result = new(pbGroup.SpaceArticleIDResp)
	result.CommonResp = new(pbGroup.CommonResp)
	err = imdb.DelSpaceArticleID(utils.StringToInt64(reqid.ID), reqid.UserID)
	if err != nil {
		result.CommonResp.ErrCode = constant.ErrDB.ErrCode
		result.CommonResp.ErrMsg = err.Error()
		return
	}
	return result, nil
}

func (s *groupServer) CreatePushSpaceArticelOrder(ctx context.Context, req *pbGroup.CreatePushSpaceArticelOrderReq) (
	resultData *pbGroup.CreatePushSpaceArticelOrderResp, err error) {

	etcdConn := getcdv3.GetDefaultConn(
		config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImOrder,
		req.OperationID,
	)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		return &pbGroup.CreatePushSpaceArticelOrderResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrServer.ErrCode, ErrMsg: errMsg}}, nil
	}
	client := pbOrder.NewOrderServiceClient(etcdConn)
	resp, err := client.CreatePayScanBlockTask(ctx, &pbOrder.CreatePayScanBlockTaskReq{
		USD:         config.Config.SpaceArticle.PushUsdPrice,
		FormAddress: req.UserId,
		OperationID: req.OperationID,
		OrderId:     req.SpaceArticleId,
		TxnType:     req.TxnType,
		Mark:        constant.PayMarkPushSpaceArticleType,
	})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreatePayScanBlockTask failed", err.Error())
		return &pbGroup.CreatePushSpaceArticelOrderResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrServer.ErrCode, ErrMsg: err.Error()}}, nil
	}
	if resp.CommonResp.ErrCode != 0 {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreatePayScanBlockTask failed", resp.CommonResp.ErrMsg)
		return &pbGroup.CreatePushSpaceArticelOrderResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrServer.ErrCode, ErrMsg: resp.CommonResp.ErrMsg}}, nil
	}
	resultData = &pbGroup.CreatePushSpaceArticelOrderResp{}
	resultData.PayInfo = resp.ScanTaskInfo
	return resultData, nil
}

func (s *groupServer) GlobalPushSpaceArticle(ctx context.Context, req *pbGroup.GlobalPushSpaceArticleReq) (
	_ *pbGroup.GlobalPushSpaceArticleResp, rpcErr error) {
	var err error
	defer func() {
		if err != nil {
			// 发送kafka失败、往用户账号加钱，如果失败还是失败就返回rpcErr
			if dbErr := imdb.AddUserGlobalMoneyCountByArticle(req.UserId, "pushSpaceArticleFailInc", req.SpaceArticleId, config.Config.SpaceArticle.PushUsdPrice); dbErr != nil {
				// 用户没找到,不管它
				if errors.Is(dbErr, gorm.ErrRecordNotFound) {
					log.NewWarn(req.OperationID, utils.GetSelfFuncName(), "AddUserGlobalMoneyCount failed", dbErr.Error())
					return
				}
				// 重复消费，并且重复添加余额
				DuplicateErrTip := fmt.Sprintf("Error 1062: Duplicate entry 'pushSpaceArticleFailInc:%s' for key 'tx_id'", req.SpaceArticleId)
				dbErrStr := dbErr.Error()
				if dbErrStr == DuplicateErrTip {
					log.NewWarn(req.OperationID, utils.GetSelfFuncName(), "AddUserGlobalMoneyCount failed", dbErr.Error())
					return
				}
				// 未知错误，不消费
				rpcErr = dbErr
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "AddUserGlobalMoneyCount failed", dbErr.Error())
			}
		}
	}()
	spaceArticleId := req.SpaceArticleId
	var spaceArticle *db.SpaceArticleList
	spaceArticle, err = imdb.GetSpaceArticleByID(spaceArticleId)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSpaceArticleByID failed", err.Error())
		return &pbGroup.GlobalPushSpaceArticleResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	imdb.UpdateSpaceArticleByIDGlobal(spaceArticleId)
	var articleId int64
	articleId, err = strconv.ParseInt(spaceArticle.ArticleID, 10, 64)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "strconv.ParseInt failed", err.Error())
		return &pbGroup.GlobalPushSpaceArticleResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrServer.ErrCode, ErrMsg: err.Error()}}, nil
	}
	// 全网广播某个 spaceArticle
	pushMQ := &pbChat.NewPushActionMsgMq{
		OperationID: req.OperationID,
		PushMsg: &open_im_sdk.PushMessageToMailFromUserToFans{
			OperationID:       req.OperationID,
			ContentType:       spaceArticle.ArticleType,
			ArticleID:         articleId,
			FromUserID:        spaceArticle.ReprintedID,
			FromArticleAuthor: spaceArticle.CreatorID,
			IsGlobal:          int32(spaceArticle.IsGlobal),
		},
	}
	_, _, err = s.messagePushWrite.SendMessage(pushMQ, spaceArticle.CreatorID, req.OperationID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SendMessage failed", err.Error())
		return &pbGroup.GlobalPushSpaceArticleResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrServer.ErrCode, ErrMsg: err.Error()}}, nil
	}
	return &pbGroup.GlobalPushSpaceArticleResp{}, nil
}
func (s *groupServer) GetPersonalPublishAnnounceMomentListWithIdo(ctx context.Context, req *pbGroup.GetPublishAnnouncementWithIdoReq) (
	resultData *pbGroup.GetPersonalPublishAnnouncementWithIdoResp, err error) {
	resultData = new(pbGroup.GetPersonalPublishAnnouncementWithIdoResp)
	resultData.CommonResp = new(pbGroup.CommonResp)
	resultData.TotalCount, _ = imdb.GetArticleListCountMyEmail(req.OpUserID, req.ID, 0, 0)
	resultData.CurrentPage = int64(req.PageIndex)
	dataArray, err := imdb.GetArticleListMyEmail(req.OpUserID, req.ID, int64(req.PageIndex), int64(req.PageSize))
	if err == nil && len(dataArray) > 0 {
		for _, value := range dataArray {
			wenzhang := &pbGroup.PersonalGenerateArticle{
				GenerateArticle: &pbGroup.GenerateArticle{
					ArticleID:           value.ArticleID,
					CreatedAt:           value.CreatedAt.Format("2006-01-02 15:04:05"),
					UpdatedAt:           value.UpdatedAt.Format("2006-01-02 15:04:05"),
					DeletedAt:           value.DeletedAt.Format("2006-01-02 15:04:05"),
					OpUserID:            value.ReprintId,
					GroupArticleId:      int32(value.GroupArticleID),
					CreatorUserId:       value.CreatorUserID,
					AnnouncementContent: value.AnnouncementContent,
					AnnouncementUrl:     value.AnnouncementUrl,
					LikeCount:           int32(value.LikeCount),
					RewordCount:         int32(value.RewordCount),
					IsGlobal:            value.IsGlobal,
					OrderId:             value.OrderID,
					Status:              value.Status,
					AnnouncementTitle:   value.AnnouncementTitle,
					AnnouncementSummary: value.AnnouncementSummary,
					ID:                  value.ID,
					ArticleType:         value.ArticleType,
				},
				OperatorInfo: new(open_im_sdk.PublicUserInfo),
				CreatorInfo:  new(open_im_sdk.PublicUserInfo),
			}
			userInfo, err := rocksCache.GetUserBaseInfoFromCache(value.CreatorUserID)
			if err == nil {
				utils.CopyStructFields(wenzhang.CreatorInfo, userInfo)
				linkTree, err := rocksCache.GetUserBaseInfoFromCacheUserLink(value.CreatorUserID)
				if err == nil {
					for _, value := range linkTree {
						tempData := new(open_im_sdk.LinkTreeMsgReq)
						utils.CopyStructFields(tempData, value)
						wenzhang.CreatorInfo.LinkTree = append(wenzhang.CreatorInfo.LinkTree, tempData)
					}
				}
			} else {
				wenzhang.CreatorInfo = nil
			}
			if value.ReprintId == value.CreatorUserID && userInfo != nil {
				utils.CopyStructFields(wenzhang.OperatorInfo, userInfo)
				for _, value := range wenzhang.CreatorInfo.LinkTree {
					tempData := value
					wenzhang.OperatorInfo.LinkTree = append(wenzhang.OperatorInfo.LinkTree, tempData)
				}
			} else {
				userInfo2, err := rocksCache.GetUserBaseInfoFromCache(value.ReprintId)
				if err == nil {
					utils.CopyStructFields(wenzhang.OperatorInfo, userInfo2)
					linkTree, err := rocksCache.GetUserBaseInfoFromCacheUserLink(value.ReprintId)
					if err == nil {
						for _, value := range linkTree {
							tempData := new(open_im_sdk.LinkTreeMsgReq)
							utils.CopyStructFields(tempData, value)
							wenzhang.OperatorInfo.LinkTree = append(wenzhang.OperatorInfo.LinkTree, tempData)
						}
					}
				} else {
					wenzhang.OperatorInfo = nil
				}
			}
			resultData.PublishAnnounce = append(resultData.PublishAnnounce, wenzhang)
		}
	}
	return resultData, err
}
func (s *groupServer) GetHotCommunityBannerAnnouncementList(ctx context.Context,
	req *pbGroup.GetHotCommunityReq) (
	resultData *pbGroup.GetPersonalPublishAnnouncementWithIdoResp, err error) {
	resultData = new(pbGroup.GetPersonalPublishAnnouncementWithIdoResp)
	resultData.CommonResp = new(pbGroup.CommonResp)
	resultData.TotalCount = 0
	dataArray, err := imdb.GetSpaceArticleListBanner()
	if err == nil && len(dataArray) > 0 {
		for _, value := range dataArray {
			wenzhang := &pbGroup.PersonalGenerateArticle{
				GenerateArticle: &pbGroup.GenerateArticle{
					ArticleID:           value.ArticleID,
					CreatedAt:           value.CreatedAt.Format(time.DateTime),
					UpdatedAt:           value.UpdatedAt.Format("2006-01-02 15:04:05"),
					DeletedAt:           value.DeletedAt.Format("2006-01-02 15:04:05"),
					OpUserID:            value.ReprintId,
					GroupArticleId:      int32(value.GroupArticleID),
					CreatorUserId:       value.CreatorUserID,
					AnnouncementContent: value.AnnouncementContent,
					AnnouncementUrl:     value.AnnouncementUrl,
					LikeCount:           int32(value.LikeCount),
					RewordCount:         int32(value.RewordCount),
					IsGlobal:            value.IsGlobal,
					OrderId:             value.OrderID,
					Status:              value.Status,
					AnnouncementTitle:   value.AnnouncementTitle,
					AnnouncementSummary: value.AnnouncementSummary,
					ID:                  value.ID,
					ArticleType:         value.ArticleType,
				},
				OperatorInfo: new(open_im_sdk.PublicUserInfo),
				CreatorInfo:  new(open_im_sdk.PublicUserInfo),
			}
			userInfo, err := rocksCache.GetUserBaseInfoFromCache(value.CreatorUserID)
			if err == nil {
				utils.CopyStructFields(wenzhang.CreatorInfo, userInfo)
			} else {
				wenzhang.CreatorInfo = nil
			}
			if value.ReprintId == value.CreatorUserID && userInfo != nil {
				utils.CopyStructFields(wenzhang.OperatorInfo, userInfo)
			} else {
				userInfo2, err := rocksCache.GetUserBaseInfoFromCache(value.ReprintId)
				if err == nil {
					utils.CopyStructFields(wenzhang.OperatorInfo, userInfo2)
				} else {
					wenzhang.OperatorInfo = nil
				}
			}
			resultData.PublishAnnounce = append(resultData.PublishAnnounce, wenzhang)
		}
	} else {

	}
	return resultData, err
}
