package game

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbweb3pb "Open_IM/pkg/proto/web3pub"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"net/http"
	"strings"
	"time"
)

// PostStartGame
// @Summary		游戏开始
// @Description	游戏开始
// @Tags			游戏相关
// @ID				PostStartGame
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.PostStartGameReq	true	"1为注册开始游戏 2 为注册结束游戏"
// @Produce		json
// @Success		0	{object}	api.PostStartGameResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/gameapi/start_game [post]
func PostStartGame(c *gin.Context) {
	var (
		req   api.PostStartGameReq
		resp  = new(api.PostStartGameResp)
		reqPb = new(pbweb3pb.UserGameReq)
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	utils.CopyStructFields(reqPb, &req)
	reqPb.Ip = GetRealIP(c)
	reqPb.UserAgent = c.Request.UserAgent()
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImWeb3Js, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbweb3pb.NewWeb3PubClient(etcdConn)
	RpcResp, err := client.PostGamingStatus(context.Background(), reqPb)
	if err == nil {
		resp.ErrMsg = RpcResp.CommonResp.ErrMsg
		resp.ErrCode = RpcResp.CommonResp.ErrCode
		resp.StartTime = RpcResp.StartTime
	} else {
		resp.ErrMsg = err.Error()
		resp.ErrCode = constant.ErrInternal.ErrCode
	}
	c.JSON(http.StatusOK, resp)
	return
}

func GetRealIP(c *gin.Context) string {
	// 从 X-Real-IP 头部中获取 IP 地址
	ip := c.Request.Header.Get("X-Real-IP")

	// 如果未能获取到，则从 c.Request.RemoteAddr 中获取 IP 地址
	if ip == "" {
		ip = c.Request.RemoteAddr
		if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
			ip = ip[:colonIndex]
		}
	}

	return ip
}

// GetGameRankList
// @Summary		游戏排行榜
// @Description	游戏排行榜
// @Tags			游戏相关
// @ID				GetGameRankList
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.PostGameRankListReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.PostGameRankListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/gameapi/game_rank_list [post]
func GetGameRankList(c *gin.Context) {
	var (
		req   api.PostGameRankListReq
		resp  = new(api.PostGameRankListResp)
		reqPb = new(pbweb3pb.UserGameRankListReq)
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	utils.CopyStructFields(reqPb, &req)
	reqPb.GameID = utils.StringToInt32(req.GameID)
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImWeb3Js, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbweb3pb.NewWeb3PubClient(etcdConn)
	RpcResp, err := client.GetGamingRankStatus(context.Background(), reqPb)
	if err != nil {
		resp.ErrMsg = err.Error()
		resp.ErrCode = constant.ErrInternal.ErrCode
		return
	}
	if len(RpcResp.UserRankInfo) > 0 {
		for _, value := range RpcResp.UserRankInfo {
			resp.RankLinkInfoInfo.UserScore = append(resp.RankLinkInfoInfo.UserScore, &api.UserGameScore{
				UserID:             value.UserID,
				Nickname:           value.Nickname,
				FaceURL:            value.FaceURL,
				Reward:             value.Reward,
				Score:              int64(value.Score),
				RankIndex:          value.RankIndex,
				TokenContractChain: value.TokenContractChain,
			})
		}
	}
	if RpcResp.UserSelfRankInfo != nil {
		resp.RankLinkInfoInfo.UserSelfScore = &api.UserGameScore{
			UserID:             RpcResp.UserSelfRankInfo.UserID,
			Nickname:           RpcResp.UserSelfRankInfo.Nickname,
			FaceURL:            RpcResp.UserSelfRankInfo.FaceURL,
			Reward:             RpcResp.UserSelfRankInfo.Reward,
			Score:              int64(RpcResp.UserSelfRankInfo.Score),
			RankIndex:          RpcResp.UserSelfRankInfo.RankIndex,
			TokenContractChain: RpcResp.UserSelfRankInfo.TokenContractChain,
		}
	}
	c.JSON(http.StatusOK, resp)
	return
}

// GetGameList
// @Summary		游戏列表
// @Description	游戏列表
// @Tags			游戏相关
// @ID				GetGameList
// @Accept			json
// @Param			token	header	string	true	"im token"
// @Param			req		body	api.PostGameRankListReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.PostGameListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/gameapi/game_list [post]
func GetGameList(c *gin.Context) {
	var (
		req  api.PostGameRankListReq
		resp = new(api.PostGameListResp)
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	//var ok bool
	//var errInfo string
	//var userID string
	//ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	//if !ok {
	//	errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
	//	log.NewError(req.OperationID, errMsg)
	//	c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
	//	return
	//}
	//if userID == "" {
	//	c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": " not  login"})
	//	return
	//}
	resultGame, err := imdb.GetGameListFromDB(req.GameID)
	if err == nil {
		if len(resultGame) > 0 {
			resp.GameConfigList = make([]*api.ApiGameConfig, 0)
			for _, value := range resultGame {
				//&db.GameConfig{
				//	GameId:               value.GameId,
				//	CreatedAt:            value.CreatedAt,
				//	GameName:             value.GameName,
				//	Status:               value.Status,
				//	GameUrl:              value.GameUrl,
				//	GameDesc:             value.GameDesc,
				//	GameCondition:        value.GameCondition,
				//	GameCurrentPrizePool: value.GameCurrentPrizePool,
				//	GameMinPrizePool:     value.GameMinPrizePool,
				//}
				tempApiGameConfig := new(api.ApiGameConfig)
				utils.CopyStructFields(tempApiGameConfig, value)
				// 解析Cron表达式
				specParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
				schedule, _ := specParser.Parse(value.RewordCrontab)
				now := time.Now()
				nextTime := schedule.Next(now)
				tempApiGameConfig.NextTime = nextTime.UTC().Unix()
				resp.GameConfigList = append(resp.GameConfigList, tempApiGameConfig)
			}
		}
		c.JSON(http.StatusOK, resp)
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
}
