package user

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	kafkaMessage "Open_IM/pkg/proto/kafkamessage"
	pbRelay "Open_IM/pkg/proto/relay"
	sdkws "Open_IM/pkg/proto/sdk_ws"
	pbUser "Open_IM/pkg/proto/user"
	pbweb3pb "Open_IM/pkg/proto/web3pub"
	rpc "Open_IM/pkg/proto/web3pub"
	"Open_IM/pkg/utils"
	"Open_IM/pkg/xkafka"
	"Open_IM/pkg/xlog"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/syncx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/clause"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
)

// 本来要实现一种分布式的二阶段事务，目前先这么做，待优化
func (s *userServer) TransferChatTokenFromUserToGroup(ctx context.Context, req *pbUser.TransferChatTokenOperatorReq) (*pbUser.TransferChatTokenOperatorResp, error) {
	resp := &pbUser.TransferChatTokenOperatorResp{CommonResp: &pbUser.CommonResp{}}
	groupId := req.ToGroupID
	chatTokenCount := req.ChatTokenCount
	// 查询用户余额
	rs := db.DB.Pool
	mutexname := "user_chat_token:" + req.OpUserID
	mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
	if err := mutex.LockContext(ctx); err != nil {
		resp.CommonResp.ErrCode = 10001
		resp.CommonResp.ErrMsg = "正在操作，请稍后尝试"
		return resp, nil
	}
	defer mutex.UnlockContext(ctx)
	err := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		var userInfo db.User
		err := tx.Table("users").Where("user_id=?", req.OpUserID).Take(&userInfo).Error
		if err != nil {
			return err
		}
		timenow := time.Now()
		resultToken := int64(userInfo.ChatTokenCount) - chatTokenCount
		if resultToken < 0 || chatTokenCount <= 0 {
			return errors.New("用户余额不够 禁止转账")
		}
		//查询群id 是否存在：
		groupInfo, err := rocksCache.GetGroupInfoFromCache(groupId) //群不存在 或者说已经解散掉了
		if err != nil {
			return errors.New("群不存在 不可以转账")
		}
		// 入用户消费记录
		err = tx.Create(&db.UserChatTokenRecord{
			CreatedTime: timenow,
			UserID:      userInfo.UserID,
			TxID:        timenow.Format("20060102150405") + utils.Int64ToString(timenow.UnixMilli()%1000) + utils.Md5(userInfo.UserID),
			TxType:      "transfer",
			OldToken:    userInfo.ChatTokenCount,
			NewToken:    uint64(resultToken),
			ParamStr:    fmt.Sprintf(`{"transfer":%d,"groupid":%s"}`, req.ChatTokenCount, req.ToGroupID),
			ChainID:     "0",
		}).Error
		//扣款用户的余额
		err = tx.Table("users").Where("user_id=?", req.OpUserID).UpdateColumn("chat_token_count",
			gorm.Expr("chat_token_count - ?", chatTokenCount)).Error
		if err != nil {
			return err
		}
		//并发加锁group的余额
		mutexnamegroup := "user_chat_token_group:" + groupId
		mutexgroup := rs.NewMutex(mutexnamegroup, redsync.WithTries(3),
			redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
		if err = mutexgroup.LockContext(ctx); err != nil {
			return err
		}
		defer mutex.UnlockContext(ctx)
		//群余额记录增加
		err = tx.Create(&db.UserChatTokenRecord{
			CreatedTime: timenow,
			UserID:      groupId,
			TxID:        timenow.Format("20060102150405") + utils.Int64ToString(timenow.UnixMilli()%1000) + utils.Md5(userInfo.UserID),
			TxType:      "add",
			OldToken:    uint64(groupInfo.ChatTokenCount),
			NewToken:    uint64(groupInfo.ChatTokenCount + chatTokenCount),
			ParamStr:    fmt.Sprintf(`{"transfer":%d,"groupid":%s"}`, req.ChatTokenCount, req.ToGroupID),
		}).Error
		if err != nil {
			return errors.New("群不存在 不可以转账")
		}
		//群余额增加
		err = tx.Table("groups").Where("group_id=?", groupId).UpdateColumn("chat_token_count",
			gorm.Expr("chat_token_count + ?", chatTokenCount)).Error
		if err != nil {
			return err
		}
		resp.NowChatToken = uint64(resultToken)
		resp.GroupChatTokenCount = groupInfo.ChatTokenCount + chatTokenCount
		rocksCache.DelUserInfoFromCache(req.OpUserID)
		rocksCache.DelGroupInfoFromCache(groupId)
		return nil
	})
	if err != nil {
		resp.CommonResp.ErrCode = 10001
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	return resp, nil

}

// OperatorUserChatToken token 数量为0.03 换1000个chattoken
func (s *userServer) OperatorUserChatToken(ctx context.Context, req *pbUser.OperatorUserChatTokenReq) (*pbUser.OperatorUserChatTokenResp, error) {
	resp := &pbUser.OperatorUserChatTokenResp{CommonResp: &pbUser.CommonResp{}}
	rs := db.DB.Pool

	switch req.Operator {
	case "add":
		{
			//需要解析到转账用户的额token的信息在对其枷锁
			//如果添加的话 获取key的记录。
			etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
				strings.Join(config.Config.Etcd.EtcdAddr, ","),
				config.Config.RpcRegisterName.OpenImWeb3Js, req.OperationID)
			if etcdConn == nil {
				errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
				log.NewError(req.OperationID, errMsg)
				resp.CommonResp.ErrCode = 10001
				resp.CommonResp.ErrMsg = "正在操作，请稍后尝试1"
				return resp, nil
			}
			client := rpc.NewWeb3PubClient(etcdConn)
			reqNewTxid := &pbweb3pb.EthRpcTxIDReq{
				OperatorID: req.OperationID,
				TxID:       req.TxID,
				ChainID:    req.ChainID,
			}
			RpcResp, err := client.GetEthTxIDTaskRpc(context.Background(), reqNewTxid)
			if err != nil || RpcResp == nil || RpcResp.FromAddress == "" && RpcResp.ToAddress != config.Config.ReceiveTokenAddress {
				resp.CommonResp.ErrCode = 10001
				resp.CommonResp.ErrMsg = "正在操作，请稍后尝试2" + RpcResp.ToAddress + "&" + "RpcResp.FromAddress"
				log.NewInfo("now is data :", utils.StructToJsonString(RpcResp))
				return resp, nil
			}
			mutexname := "user_chat_token:" + RpcResp.FromAddress
			mutex := rs.NewMutex(mutexname, redsync.WithExpiry(time.Second*10))

			if err := mutex.LockContext(ctx); err != nil {
				resp.CommonResp.ErrCode = 10001
				resp.CommonResp.ErrMsg = "正在操作，请稍后尝试3"
				return resp, nil
			}
			defer mutex.UnlockContext(ctx)
			resp.CommonResp.ErrMsg = RpcResp.FromAddress
			err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
				var userInfo db.User
				err := tx.Table("users").Where("user_id=?", RpcResp.FromAddress).Take(&userInfo).Error
				if err != nil {
					return err
				}
				nowToken := userInfo.ChatTokenCount
				nowChatCount := userInfo.ChatCount
				var nValueCmp = 0
				var paramStr = ""
				if nValueCmp, err = utils.FloatCompare(RpcResp.Value, 0.005); err == nil && nValueCmp == 0 {
					nowChatCount = userInfo.ChatCount + 1
					paramStr = `{"recharge":"add count 1"}`

				} else {
					nowToken = userInfo.ChatTokenCount + uint64(RpcResp.Value*1000/0.003)
					paramStr = fmt.Sprintf(`{"recharge":"add token %d}`, uint64(RpcResp.Value*1000/0.003))
				}
				err = tx.Create(&db.UserChatTokenRecord{
					CreatedTime: time.Now(),
					UserID:      RpcResp.FromAddress,
					TxID:        RpcResp.TransactionHash,
					TxType:      "add",
					OldToken:    userInfo.ChatTokenCount,
					ChainID:     req.ChainID,
					NewToken:    nowToken,
					NowCount:    nowChatCount,
					ParamStr:    paramStr,
				}).Error
				if err != nil {
					return err
				}
				if nValueCmp != 0 {
					err = tx.Table("users").Where("user_id=?", RpcResp.FromAddress).
						UpdateColumn("chat_token_count", gorm.Expr("chat_token_count + ?", uint64(RpcResp.Value*1000/0.003))).Error
					if err != nil {
						return err
					}
				} else {
					err = tx.Table("users").Where("user_id=?", RpcResp.FromAddress).UpdateColumn("chat_count", gorm.Expr("chat_count + ?", 1)).Error
					if err != nil {
						return err
					}
				}

				return nil
			})
			if err != nil {
				resp.CommonResp.ErrCode = 10001
				resp.CommonResp.ErrMsg = "正在操作，请稍后尝试4"
				return resp, nil
			}
			if err2222 := rocksCache.DelUserInfoFromCache(strings.ToLower(RpcResp.FromAddress)); err2222 != nil {
				fmt.Println(err2222.Error())
			}
			return resp, nil
		}
	case "groupsub":
		mutexname := "user_chat_token_group:" + req.OpUserID
		mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
		if err := mutex.LockContext(ctx); err != nil {
			resp.CommonResp.ErrCode = 10001
			resp.CommonResp.ErrMsg = "正在操作，请稍后尝试"
			return resp, nil
		}
		defer mutex.UnlockContext(ctx)
		err := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
			var groupInfo db.Group
			err := tx.Table("groups").Where("group_id=?", req.OpUserID).Take(&groupInfo).Error
			if err != nil {
				return err
			}
			timenow := time.Now()
			needCost := req.Value
			resultToken := groupInfo.ChatTokenCount - req.Value
			if req.Value >= groupInfo.ChatTokenCount {
				needCost = groupInfo.ChatTokenCount
				resultToken = 0
			}
			err = tx.Create(&db.UserChatTokenRecord{
				CreatedTime: timenow,
				UserID:      groupInfo.GroupID,
				TxID:        timenow.Format("20060102150405") + utils.Int64ToString(timenow.UnixMilli()%1000) + utils.Md5(req.OpUserID),
				TxType:      "sub",
				OldToken:    uint64(groupInfo.ChatTokenCount),
				NewToken:    uint64(resultToken),
				ParamStr:    req.ParamStr,
				ChainID:     "0",
			}).Error

			err = tx.Table("groups").Where("group_id=?", req.OpUserID).UpdateColumn("chat_token_count",
				gorm.Expr("chat_token_count - ?", needCost)).Error
			if err != nil {
				return err
			}

			rocksCache.DelUserInfoFromCache(req.OpUserID)
			return nil
		})
		if err != nil {
			resp.CommonResp.ErrCode = 10001
			resp.CommonResp.ErrMsg = "正在操作，请稍后尝试"
			return resp, nil
		}
	case "sub":
		mutexname := "user_chat_token:" + req.OpUserID
		mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
		if err := mutex.LockContext(ctx); err != nil {
			resp.CommonResp.ErrCode = 10001
			resp.CommonResp.ErrMsg = "正在操作，请稍后尝试"
			return resp, nil
		}
		defer mutex.UnlockContext(ctx)
		err := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
			var userInfo db.User
			err := tx.Table("users").Where("user_id=?", req.OpUserID).Take(&userInfo).Error
			if err != nil {
				return err
			}
			timenow := time.Now()
			needCost := uint64(req.Value)
			resultToken := uint64(int64(userInfo.ChatTokenCount) - req.Value)
			if req.TxType == "times" {
				totalCount := userInfo.ChatCount - 1
				if totalCount < 0 {
					userInfo.ChatCount = 0
				}
				err = tx.Create(&db.UserChatTokenRecord{
					CreatedTime: timenow,
					UserID:      userInfo.UserID,
					TxID:        timenow.Format("20060102150405") + utils.Int64ToString(timenow.UnixMilli()%1000) + utils.Md5(userInfo.UserID),
					TxType:      "sub",
					OldToken:    userInfo.ChatTokenCount,
					NewToken:    userInfo.ChatTokenCount,
					ParamStr:    req.ParamStr,
					ChainID:     "0",
					NowCount:    userInfo.ChatCount,
				}).Error
			} else {
				if uint64(req.Value) >= userInfo.ChatTokenCount {
					needCost = userInfo.ChatTokenCount
					resultToken = 0
				}
				err = tx.Create(&db.UserChatTokenRecord{
					CreatedTime: timenow,
					UserID:      userInfo.UserID,
					TxID:        timenow.Format("20060102150405") + utils.Int64ToString(timenow.UnixMilli()%1000) + utils.Md5(userInfo.UserID),
					TxType:      "sub",
					OldToken:    userInfo.ChatTokenCount,
					NewToken:    resultToken,
					ParamStr:    req.ParamStr,
					ChainID:     "0",
				}).Error
			}

			if err != nil {
				return err
			}
			if req.TxType == "times" {
				err = tx.Table("users").Where("user_id=?", req.OpUserID).UpdateColumn("chat_count",
					gorm.Expr("chat_count - ?", 1)).Error
				if err != nil {
					return err
				}
			} else {
				err = tx.Table("users").Where("user_id=?", req.OpUserID).UpdateColumn("chat_token_count",
					gorm.Expr("chat_token_count - ?", needCost)).Error
				if err != nil {
					return err
				}
			}
			rocksCache.DelUserInfoFromCache(req.OpUserID)
			return nil
		})
		if err != nil {
			resp.CommonResp.ErrCode = 10001
			resp.CommonResp.ErrMsg = "正在操作，请稍后尝试"
			return resp, nil
		}
	}
	return resp, nil
}

func (s *userServer) GetUserInfoWithoutToken(ctx context.Context, req *pbUser.GetUserInfoReq) (resultResp *pbUser.GetUserInfoWithProfileResp, err error) {
	log.NewInfo(req.OperationID, "GetUserInfo args ", req.String())
	resultResp = new(pbUser.GetUserInfoWithProfileResp)
	resultResp.CommonResp = new(pbUser.CommonResp)
	resultResp.CommonResp.ErrCode = constant.ErrArgs.ErrCode
	resultResp.CommonResp.ErrMsg = constant.ErrArgs.ErrMsg
	if len(req.UserIDList) > 0 {
		userID := req.UserIDList[0]
		userInfo := new(sdkws.UserInfo)
		user, err := rocksCache.GetUserInfoFromCacheByMerLin(userID, "1")
		if err == nil {
			utils.CopyStructFields(userInfo, user)
			userInfo.BirthStr = utils.TimeToString(user.Birth)
			resultResp.UserInfoList = userInfo
			userThirdInfo, err2 := imdb.GetUserTwitterWithFlagAll(userID)
			if err2 == nil {
				if req.OpUserID == req.UserIDList[0] {
					resultResp.Twitter = userThirdInfo.Twitter
					resultResp.DnsDomain = userThirdInfo.DnsDomain
					resultResp.DnsDomainVerify = userThirdInfo.DnsDomainVerify
					resultResp.EmailAddress = userThirdInfo.UserAddress
				} else {
					if userThirdInfo.ShowTwitter == 1 {
						resultResp.Twitter = userThirdInfo.Twitter
					} else {
						resultResp.Twitter = ""
					}
					resultResp.DnsDomain = userThirdInfo.DnsDomain
					resultResp.DnsDomainVerify = userThirdInfo.DnsDomainVerify
					if userThirdInfo.ShowUserAddress == 1 {
						resultResp.EmailAddress = userThirdInfo.UserAddress
					} else {
						resultResp.EmailAddress = ""
					}
				}
			}
			userLinkDb, err := imdb.GetUserLinkTree(userID, "")
			if err == nil {
				resultResp.LinkTree = make([]*sdkws.LinkTreeMsgReq, 0, len(userLinkDb))
				for key, _ := range userLinkDb {
					if userID != req.OpUserID {
						if userLinkDb[key].ShowStatus == 1 {
							resultResp.LinkTree = append(resultResp.LinkTree, &sdkws.LinkTreeMsgReq{
								LinkName:    userLinkDb[key].LinkName,
								Link:        userLinkDb[key].Link,
								FaceUrl:     userLinkDb[key].FaceURL,
								ShowStatus:  int32(userLinkDb[key].ShowStatus),
								UserID:      userID,
								DefaultIcon: userLinkDb[key].DefaultIcon,
								Des:         userLinkDb[key].Des,
								Bgc:         userLinkDb[key].Bgc,
								DefaultUrl:  userLinkDb[key].DefaultUrl,
								Type:        userLinkDb[key].Type,
							})
						}
					} else {
						resultResp.LinkTree = append(resultResp.LinkTree, &sdkws.LinkTreeMsgReq{
							LinkName:    userLinkDb[key].LinkName,
							Link:        userLinkDb[key].Link,
							FaceUrl:     userLinkDb[key].FaceURL,
							ShowStatus:  int32(userLinkDb[key].ShowStatus),
							UserID:      userID,
							DefaultIcon: userLinkDb[key].DefaultIcon,
							Des:         userLinkDb[key].Des,
							Bgc:         userLinkDb[key].Bgc,
							DefaultUrl:  userLinkDb[key].DefaultUrl,
							Type:        userLinkDb[key].Type,
						})
					}
				}
			}
			dbUserInfo, _ := imdb.GetUserByUserID(userID)
			if dbUserInfo != nil {
				resultResp.UserProfile = dbUserInfo.UserProfile
			}
			resultResp.CommonResp.ErrCode = 0
			resultResp.CommonResp.ErrMsg = ""
			return resultResp, nil
		} else {
			resultResp.CommonResp.ErrMsg = err.Error()
		}
	}
	return resultResp, nil
}
func (s *userServer) BindShowNft(ctx context.Context, req *pbUser.RPCBindShowNftReq) (*pbUser.RPCBindShowNftResp, error) {
	//存入数据库内容
	var InsertMap map[string]bool
	InsertMap = make(map[string]bool, 0)
	var inserInto []*db.UserNftConfig
	if len(req.NftInfo) > 0 {
		for _, value := range req.NftInfo {
			if value.NftTokenURL != "" {
				if _, ok := InsertMap[fmt.Sprintf("%d:%s:%s", value.NftChainID, value.NftContract, value.TokenID)]; !ok {
					inserInto = append(inserInto, &db.UserNftConfig{
						UserID:          req.UserID,
						NftChainID:      int(value.NftChainID),
						NftContract:     value.NftContract,
						TokenID:         value.TokenID,
						NftContractType: value.NftContractType,
						NftTokenURL:     value.NftTokenURL,
						IsShow:          1,
						Md5index:        utils.Md5(fmt.Sprintf("%s%d%s%s", req.UserID, value.NftChainID, value.NftContract, value.TokenID)),
					})
					InsertMap[fmt.Sprintf("%d:%s:%s", value.NftChainID, value.NftContract, value.TokenID)] = true
				}

			}
		}
		if len(inserInto) == 0 {
			return &pbUser.RPCBindShowNftResp{
				CommonResp: &pbUser.CommonResp{
					ErrCode: constant.ErrDB.ErrCode,
					ErrMsg:  "无法获取nft的描述地址",
				},
			}, nil
		}
	}

	err := imdb.InsertIntoUserNftConfig(req.UserID, inserInto)
	if err != nil {
		return &pbUser.RPCBindShowNftResp{
			CommonResp: &pbUser.CommonResp{
				ErrCode: constant.ErrDB.ErrCode,
				ErrMsg:  err.Error(),
			},
		}, nil
	}
	fmt.Println("完成设置nft 的树枝了")
	return &pbUser.RPCBindShowNftResp{
		CommonResp: &pbUser.CommonResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
	}, nil
}
func (s *userServer) GetBindShowNft(ctx context.Context, req *pbUser.RPCBindShowNftReq) (*pbUser.GetRPCBindShowNftResp, error) {
	//TODO implement me
	//存入数据库内容
	//var inserInto []*imdb.UserNftConfigInfo

	insertInto, err := imdb.GetUserNftConfig(req.UserID, req.OpUserID)
	if err != nil {
		return &pbUser.GetRPCBindShowNftResp{
			CommonResp: &pbUser.CommonResp{
				ErrCode: constant.ErrDB.ErrCode,
				ErrMsg:  err.Error(),
			},
		}, nil
	}

	pbResp := &pbUser.GetRPCBindShowNftResp{
		CommonResp: &pbUser.CommonResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
	}
	for _, value := range insertInto {
		pbResp.NftInfo = append(pbResp.NftInfo, &pbUser.NftInfo{
			NftChainID:      int32(value.NftChainID),
			NftContract:     value.NftContract,
			TokenID:         value.TokenID,
			NftContractType: value.NftContractType,
			NftTokenURL:     value.NftTokenURL,
			LikesCount:      value.LikeCount,
			ID:              value.ID,
			IsLikes:         value.IsLikes,
		})
	}
	return pbResp, nil
}
func (s *userServer) GetShowNftLikeStatus(ctx context.Context, req *pbUser.RpcLikeShowNftStatusReq) (*pbUser.RpcLikeShowNftStatusResp, error) {

	result := new(pbUser.RpcLikeShowNftStatusResp)
	result.NftLikeCount = imdb.GetUserNFtConfigLike(utils.StringToInt64(req.ArticleID))
	if req.UserID != "" {
		result.NftIsLike = int32(imdb.GetUserNFtConfigLikeWithUserID(utils.StringToInt64(req.ArticleID), req.UserID))
	}
	return result, nil

}

func (s *userServer) RpcUserSettingInfo(ctx context.Context, req *pbUser.GetShowUserSettingReq) (resultResp *pbUser.GetShowUserSettingResp, err error) {
	resultResp = new(pbUser.GetShowUserSettingResp)
	//查看用户的信息:
	resultResp.CommonResp = new(pbUser.CommonResp)
	if req.UserID == "" {
		resultResp.CommonResp.ErrCode = constant.ErrInternal.ErrCode
		resultResp.CommonResp.ErrMsg = "user id is null"
		return
	}
	//查询个人信息
	userInfo, err := rocksCache.GetUserBaseInfoFromCache(req.UserID)
	if err != nil {
		resultResp.CommonResp.ErrCode = constant.ErrInternal.ErrCode
		resultResp.CommonResp.ErrMsg = "不存在该用户信息"
		return
	}
	//设置resultResp的基本用户信息
	resultResp.UserID = req.UserID
	resultResp.UserIntroduction = userInfo.UserIntroduction
	resultResp.UserProfile = userInfo.UserProfile
	resultResp.FaceURL = userInfo.FaceURL
	resultResp.Nickname = userInfo.Nickname
	resultResp.ShowBalance = userInfo.ShowBalance
	resultResp.OpenAnnouncement = userInfo.OpenAnnouncement
	if strList := strings.Split(userInfo.TokenContractChain, "&"); len(strList) >= 2 {
		resultResp.UserHeadTokenInfo = &pbUser.NftInfo{
			NftChainID:      utils.StringToInt32(strList[1]),
			NftContract:     strList[0],
			TokenID:         userInfo.TokenId,
			NftContractType: "erc721",
		}
	}

	var userThirdInfo []*imdb.UserThirdPath
	if req.OpUserID != req.UserID {
		//查看别人
		userThirdInfo, _ = imdb.GetThirdUserInfoWithShowFlagWithOutDomain([]string{req.UserID})
	} else {
		//查看自己
		userThirdInfo, _ = imdb.GetThirdUserInfoWithOutDomain([]string{req.UserID})
	}
	if len(userThirdInfo) > 0 {
		resultResp.DnsDomain = userThirdInfo[0].DnsDomain
		resultResp.UserTwitter = userThirdInfo[0].Twitter
		resultResp.IsShowTwitter = userThirdInfo[0].ShowTwitter
		resultResp.EmailAddress = userThirdInfo[0].UserAddress
		resultResp.ShowUserEmail = userThirdInfo[0].ShowUserAddress
		resultResp.DnsDomainVerify = userThirdInfo[0].DnsDomainVerify
	}
	insertInto, err := imdb.GetUserNftConfig(req.UserID, req.OpUserID)
	for _, value := range insertInto {
		resultResp.ShowNftList = append(resultResp.ShowNftList, &pbUser.NftInfo{
			NftChainID:      int32(value.NftChainID),
			NftContract:     value.NftContract,
			TokenID:         value.TokenID,
			NftContractType: value.NftContractType,
			NftTokenURL:     value.NftTokenURL,
			LikesCount:      value.LikeCount,
			ID:              value.ID,
			IsLikes:         value.IsLikes,
		})
	}
	resultResp.FollowsCount, _ = imdb.GetUserFollowsCount(req.UserID)
	resultResp.FollowingCount, _ = imdb.GetUserFollowingCount(req.UserID)
	userLinkDb, err := imdb.GetUserLinkTree(req.UserID, req.OpUserID)
	resultResp.LinkTree = make([]*sdkws.LinkTreeMsgReq, 0, len(userLinkDb))
	for key, _ := range userLinkDb {
		if req.UserID != req.OpUserID {
			if userLinkDb[key].ShowStatus == 1 {
				resultResp.LinkTree = append(resultResp.LinkTree, &sdkws.LinkTreeMsgReq{
					LinkName:    userLinkDb[key].LinkName,
					Link:        userLinkDb[key].Link,
					FaceUrl:     userLinkDb[key].FaceURL,
					ShowStatus:  int32(userLinkDb[key].ShowStatus),
					UserID:      req.UserID,
					DefaultIcon: userLinkDb[key].DefaultIcon,
					Des:         userLinkDb[key].Des,
					Bgc:         userLinkDb[key].Bgc,
					DefaultUrl:  userLinkDb[key].DefaultUrl,
					Type:        userLinkDb[key].Type,
				})
			}
		} else {
			resultResp.LinkTree = append(resultResp.LinkTree, &sdkws.LinkTreeMsgReq{
				LinkName:    userLinkDb[key].LinkName,
				Link:        userLinkDb[key].Link,
				FaceUrl:     userLinkDb[key].FaceURL,
				ShowStatus:  int32(userLinkDb[key].ShowStatus),
				UserID:      req.UserID,
				DefaultIcon: userLinkDb[key].DefaultIcon,
				Des:         userLinkDb[key].Des,
				Bgc:         userLinkDb[key].Bgc,
				DefaultUrl:  userLinkDb[key].DefaultUrl,
				Type:        userLinkDb[key].Type,
			})
		}
	}
	return resultResp, nil
}

func (s *userServer) RpcUserSettingUpdate(ctx context.Context, req *pbUser.UpdateUserSettingReq) (resultResp *pbUser.UpdateUserSettingResp, err error) {
	resultResp = new(pbUser.UpdateUserSettingResp)
	resultResp.CommonResp = new(pbUser.CommonResp)
	if req.UserID == "" {
		resultResp.CommonResp.ErrCode = constant.ErrInternal.ErrCode
		resultResp.CommonResp.ErrMsg = "user id is null"
		return
	}
	//存入数据库内容
	var InsertMap map[string]bool
	InsertMap = make(map[string]bool, 0)
	var insertInto []*db.UserNftConfig
	if len(req.ShowNftList) > 0 {
		for _, value := range req.ShowNftList {
			if value.NftTokenURL != "" {
				if _, ok := InsertMap[fmt.Sprintf("%d:%s:%s", value.NftChainID, value.NftContract, value.TokenID)]; !ok {
					insertInto = append(insertInto, &db.UserNftConfig{
						UserID:          req.UserID,
						NftChainID:      int(value.NftChainID),
						NftContract:     value.NftContract,
						TokenID:         value.TokenID,
						NftContractType: value.NftContractType,
						NftTokenURL:     value.NftTokenURL,
						IsShow:          1,
						Md5index:        utils.Md5(fmt.Sprintf("%s%d%s%s", req.UserID, value.NftChainID, value.NftContract, value.TokenID)),
					})
					InsertMap[fmt.Sprintf("%d:%s:%s", value.NftChainID, value.NftContract, value.TokenID)] = true
				}

			}
		}
		if len(insertInto) == 0 {
			resultResp.CommonResp.ErrCode = constant.ErrInternal.ErrCode
			resultResp.CommonResp.ErrMsg = "你上传的nft 数据有问题"
			return
		}
	}

	userInfo, _ := rocksCache.GetUserBaseInfoFromCache(req.UserID)
	mapUpdateUserInfo := make(map[string]interface{}, 0)

	if req.Nickname != nil {
		mapUpdateUserInfo["name"] = req.Nickname.Value
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>:", req.Nickname.Value)

	}
	if req.FaceURL != nil {
		mapUpdateUserInfo["face_url"] = req.FaceURL.Value
	}
	if req.ShowBalance != nil {
		mapUpdateUserInfo["show_balance"] = req.ShowBalance.Value
	}
	if req.OpenAnnouncement != nil {
		mapUpdateUserInfo["open_announcement"] = req.OpenAnnouncement.Value
	}
	if req.UserIntroduction != nil {
		mapUpdateUserInfo["user_introduction"] = req.UserIntroduction.Value
	}
	if req.UserProfile != nil {
		mapUpdateUserInfo["user_profile"] = req.UserProfile.Value
	}
	if req.UserHeadTokenInfo != nil {
		mapUpdateUserInfo["token_id"] = req.UserHeadTokenInfo.NftChainID
		if req.UserHeadTokenInfo.NftContract == "" && req.UserHeadTokenInfo.NftChainID == 0 {
			mapUpdateUserInfo["token_contract_chain"] = ""
			mapUpdateUserInfo["face_url"] = ""
			mapUpdateUserInfo["token_id"] = ""
		} else {
			mapUpdateUserInfo["token_contract_chain"] = req.UserHeadTokenInfo.NftContract + "&" + utils.Int32ToString(req.UserHeadTokenInfo.NftChainID)
			mapUpdateUserInfo["face_url"] = req.UserHeadTokenInfo.NftTokenURL
			mapUpdateUserInfo["token_id"] = req.UserHeadTokenInfo.TokenID
		}

	}
	//关闭或者显示第三方信息
	mapUpdateUserThird := make(map[string]interface{}, 0)
	//在绑定的阶段 直接编辑
	//if req.EmailAddress != nil {
	//	mapUpdateUserThird["user_address"] = req.EmailAddress.Value
	//}
	//if req.UserTwitter != nil {
	//	mapUpdateUserThird["twitter"] = req.UserTwitter.Value
	//}
	if req.DnsDomain != nil {
		userThirdInfo, err := imdb.GetThirdUserInfoWithOutDomain([]string{req.UserID})

		if err == nil && len(userThirdInfo) > 0 {
			if !strings.EqualFold(req.DnsDomain.Value, userThirdInfo[0].DnsDomain) {
				mapUpdateUserThird["dns_domain"] = req.DnsDomain.Value
				mapUpdateUserThird["dns_domain_verify"] = 0
			}
			log.NewInfo(req.OperationID, "当前数据库的域名", userThirdInfo[0].DnsDomain, "数据库验证内容:", userThirdInfo[0].DnsDomainVerify)
			log.NewInfo(req.OperationID, "当前提交内容：", req.DnsDomain.Value)
		} else if err != nil {
			mapUpdateUserThird["dns_domain"] = req.DnsDomain.Value
			mapUpdateUserThird["dns_domain_verify"] = 0
		}
	}
	if req.IsShowTwitter != nil {
		mapUpdateUserThird["show_twitter"] = req.IsShowTwitter.Value
	}
	if req.IsShowUserEmail != nil {
		mapUpdateUserThird["show_user_address"] = req.IsShowUserEmail.Value
	}
	// 设置官方nft 的信息
	err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		if err = tx.Table("users").Where("user_id=?", req.UserID).Updates(mapUpdateUserInfo).Error; err != nil {
			return err
		}
		var tempUserThird db.UserThird
		if err := tx.Table("user_third").Where("user_id=?", req.UserID).First(&tempUserThird).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Table("user_third").Create(&db.UserThird{
				UserId: req.UserID,
			}).Error; err != nil {
				return err
			}
		}
		if err = tx.Table("user_third").Where("user_id=?", req.UserID).Updates(mapUpdateUserThird).Error; err != nil {
			return err
		}
		if req.ShowNftListCount != nil {
			if req.ShowNftListCount.Value == 0 {
				err = tx.Table("user_nft_config").Where("user_id=?", req.UserID).Updates(map[string]interface{}{"is_show": 0}).Error
				if err != nil {
					return err
				}

			} else {
				if req.ShowNftListCount != nil {
					if len(insertInto) == 0 {
						err = tx.Table("user_nft_config").Where("user_id=?", req.UserID).Updates(map[string]interface{}{"is_show": 0}).Error
						if err != nil {
							return err
						}
					}
					m := make(map[int]bool)
					for key, _ := range insertInto {
						m[key] = true
					}
					var inTableDataB []*db.UserNftConfig
					err = tx.Table("user_nft_config").Where("user_id=?", req.UserID).Find(&inTableDataB).Error
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return err
					}
					var showFlagID []int64
					var hideFlagID []int64
					//需要优化
					for _, valueInDb := range inTableDataB {
						inNewConfig := false
						for key2, value2 := range insertInto {
							if valueInDb.Md5index == value2.Md5index {
								if valueInDb.IsShow == 0 {
									showFlagID = append(showFlagID, valueInDb.ID)
								}
								if _, ok := m[key2]; ok {
									fmt.Println("存在重复数据 已经删除", key2, value2)
									delete(m, key2)
								}
								inNewConfig = true
								break
							}
						}
						if !inNewConfig && valueInDb.IsShow == 1 {
							hideFlagID = append(hideFlagID, valueInDb.ID)
						}
					}
					if len(showFlagID) > 0 {
						err = tx.Table("user_nft_config").Where("id in ?", showFlagID).Updates(map[string]interface{}{"is_show": 1}).Error
						if err != nil {
							return err
						}
					}
					if len(hideFlagID) > 0 {
						err = tx.Table("user_nft_config").Where("id in ?", hideFlagID).Updates(map[string]interface{}{"is_show": 0}).Error
						if err != nil {
							return err
						}
					}
					var newAppend []*db.UserNftConfig
					fmt.Println("\n", m)

					for key, _ := range m {
						newAppend = append(newAppend, insertInto[key])
					}
					if len(newAppend) > 0 {
						err = tx.Table("user_nft_config").Create(&newAppend).Error
					}
					if err != nil {
						return err
					}
				}

			}
		}
		if req.LinkTreeCount != nil {
			err = tx.Table("user_link").Where("user_id=?", req.UserID).Delete(&db.UserLink{}).Error
			if err == nil {
				if req.LinkTreeCount.Value > 0 {
					userLinkDB := make([]*db.UserLink, 0)
					for _, value := range req.LinkTree {
						tempValue := value
						showFlag := int8(0)
						if int8(tempValue.ShowStatus) > 0 {
							showFlag = 1
						}
						userLinkDB = append(userLinkDB, &db.UserLink{
							ID:          0,
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
							DeletedAt:   nil,
							UserID:      req.UserID,
							LinkName:    tempValue.LinkName,
							Link:        tempValue.Link,
							FaceURL:     tempValue.FaceUrl,
							ShowStatus:  showFlag,
							DefaultIcon: tempValue.DefaultIcon,
							Des:         tempValue.Des,
							Bgc:         tempValue.Bgc,
							DefaultUrl:  tempValue.DefaultUrl,
							Type:        tempValue.Type,
						})
					}
					err = tx.Table("user_link").Create(&userLinkDB).Error
				}
				if err != nil {
					return err
				}
			}
			_ = rocksCache.DeleteUserBaseInfoFromCacheUserLink(req.UserID)

		}
		return nil
	})
	if err != nil {
		resultResp.CommonResp.ErrCode = 500
		resultResp.CommonResp.ErrMsg = err.Error()
		return resultResp, nil
	}
	if len(mapUpdateUserInfo) > 0 {
		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
			log.NewError(req.OperationID, errMsg)
			resultResp.CommonResp.ErrCode = 500
			resultResp.CommonResp.ErrMsg = errMsg
			return
		}

		client := pbFriend.NewFriendClient(etcdConn)
		newReq := &pbFriend.GetFriendListReq{
			CommID: &pbFriend.CommID{OperationID: req.OperationID, FromUserID: req.UserID, OpUserID: req.UserID},
		}

		rpcResp, err := client.GetFriendList(context.Background(), newReq)
		if err != nil {
			log.NewError(req.OperationID, "GetFriendList failed ", err.Error(), newReq)
			resultResp.CommonResp.ErrCode = 500
			resultResp.CommonResp.ErrMsg = err.Error()
			return resultResp, nil
		}
		for _, v := range rpcResp.FriendInfoList {
			fmt.Println("\n需要提醒xxxxx 更新我的用户信息", req.OperationID, "UserInfoUpdatedNotification ", req.UserID, v.FriendUser.UserID)
			chat.UserInfoUpdatedNotification(req.OperationID, req.UserID, v.FriendUser.UserID)
		}
		if err = rocksCache.DelUserInfoFromCache(req.UserID); err != nil {
			log.NewError(req.OperationID, "userinfo", err.Error())
			resultResp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resultResp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
			return resultResp, nil
		}
		time.Sleep(1 * time.Second)
		chat.UserInfoUpdatedNotification(req.OperationID, req.UserID, req.UserID)
		fmt.Println("推送更新操作", req.UserID, req.UserID)
		log.Info(req.OperationID, "UserInfoUpdatedNotification ", req.UserID, req.UserID)

		if req.FaceURL != nil && req.FaceURL.Value != "" {
			s.SyncJoinedGroupMemberFaceURL(req.UserID, req.FaceURL.Value, req.OperationID, req.UserID)
		}
		if req.Nickname != nil && req.Nickname.Value != "" {
			s.SyncJoinedGroupMemberNickname(req.UserID, req.Nickname.Value, userInfo.Nickname, req.OperationID, req.UserID)
		}
	}
	return resultResp, nil
}

func (s *userServer) RpcPushMessageToFollowsUser(ctx context.Context, req *sdkws.PushMessageToMailFromUserToFans) (*emptypb.Empty, error) {
	log.NewInfo(req.OperationID, "UpdateUserInfo args ", req.String())

	if req.IsGlobal == 1 {
		index := int32(0)
		for {
			followedList, err := imdb.GetAllOpenGlobalPushUsers(req.FromUserID, index, 1000)
			if err != nil {
				break
			}
			var dbPersonal []*db.PersonalSpaceArticleList
			var pushUserIDList []string
			for _, value := range followedList {
				pushUserIDList = append(pushUserIDList, value.UserID)
				dbPersonal = append(dbPersonal, &db.PersonalSpaceArticleList{
					value.UserID,
					db.SpaceArticleList{
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						DeletedAt:    nil,
						ReprintedID:  req.FromUserID,
						CreatorID:    req.FromArticleAuthor,
						ArticleID:    utils.Int64ToString(req.ArticleID),
						ArticleType:  req.ContentType,
						ArticleIsPin: 0,
						EffectEnd:    time.Now(),
						IsGlobal:     int8(req.IsGlobal),
						Status:       1,
					},
				})

			}
			if len(dbPersonal) > 0 {
				db.DB.MysqlDB.DefaultGormDB().Table("personal_space_article_list").Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "user_id"}, {Name: "article_id"}, {Name: "article_type"}},
					DoUpdates: clause.AssignmentColumns([]string{"created_at"}),
				}).Create(&dbPersonal)
				//
			}
			s.pushMailToUserList(pushUserIDList, req)
			//发送通知
			if len(followedList) < 1000 {
				break
			}
			index++
		}
	} else {
		index := int32(0)
		for {
			followedList, err := imdb.GetUserFollowedList(req.FromUserID, index, 1000)
			if err != nil {
				break
			}
			var dbPersonal []*db.PersonalSpaceArticleList
			var pushUserIDList []string
			for _, value := range followedList {
				pushUserIDList = append(pushUserIDList, value.UserID)
				dbPersonal = append(dbPersonal, &db.PersonalSpaceArticleList{
					value.UserID,
					db.SpaceArticleList{
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						DeletedAt:    nil,
						ReprintedID:  req.FromUserID,
						CreatorID:    req.FromArticleAuthor,
						ArticleID:    utils.Int64ToString(req.ArticleID),
						ArticleType:  req.ContentType,
						ArticleIsPin: 0,
						EffectEnd:    time.Now(),
						IsGlobal:     int8(req.IsGlobal),
						Status:       1,
					},
				})

			}
			if len(dbPersonal) > 0 {
				db.DB.MysqlDB.DefaultGormDB().Table("personal_space_article_list").Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "user_id"}, {Name: "article_id"}, {Name: "article_type"}},
					DoUpdates: clause.AssignmentColumns([]string{"created_at"}),
				}).Create(&dbPersonal)
				//

			}
			s.pushMailToUserList(pushUserIDList, req)
			//发送通知
			if len(followedList) < 1000 {
				break
			}
			index++
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *userServer) pushMailToUserList(pushUserIDList []string, req *sdkws.PushMessageToMailFromUserToFans) {

	s2 := &sdkws.PushMessageToMailFromUserToFans{
		OperationID:       req.OperationID,
		ContentType:       req.ContentType,
		ArticleID:         req.ArticleID,
		FromUserID:        req.FromUserID,
		FromArticleAuthor: req.FromArticleAuthor,
		IsGlobal:          req.IsGlobal,
	}
	strMsgEx, _ := json.Marshal(s2)
	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","), req.OperationID)
	for _, v := range grpcCons {
		client := pbRelay.NewRelayClient(v)
		_, err := client.OnlineBatchAnnouncementPushOneMsg(context.Background(), &pbRelay.OnlineBatchPushOneMsgReq{
			OperationID: req.OperationID,
			MsgData: &sdkws.MsgData{
				ContentType: constant.UserAnnouncement,
				SessionType: constant.NotificationOnlinePushType,
				Ex:          string(strMsgEx),
			},
			PushToUserIDList: pushUserIDList,
		})
		if err != nil {
			log.NewError(req.OperationID, "OnlinePushMsg failed ", err.Error())
		} else {
			log.NewInfo(req.OperationID, "成功推送给用户")
			return
		}
	}

}
func (s *userServer) RpcSendEmailToUserLinkEmail(ctx context.Context, req *pbUser.EmailContentReq) (*pbUser.CommonResp, error) {
	//TODO implement me
	linkTree, err := rocksCache.GetUserBaseInfoFromCacheUserLink(req.UserID)
	if err == nil {
		for _, value := range linkTree {
			if value.LinkName == "Email" && value.Link != "" && utils.IsEmailValid(value.Link) {
				//判断用户存在的邮箱信息
				//发送kafka 时间到内容 然后让其发送email
				msgkafka := &kafkaMessage.KafkaMsg{
					MessageType: 1,
					EmailMsg: &kafkaMessage.EmailMessage{
						EmailType: 1,
						ToAddress: value.Link,
						Subject:   "Robot",
						Title:     "Robot SwapValue",
						Body:      req.EmailContent,
					},
					SmsMsg: nil,
				}
				msgvalue, _ := proto.Marshal(msgkafka)
				msg := xkafka.ProducerMessage{
					Topic: config.Config.Kafka.BusinessTop.Topic,
					Value: msgvalue,
				}
				if ct := InstallKafkaProduct(); ct != nil && ct.Producer != nil {
					err := ct.Producer.SyncSendMessage(ctx, msg)
					if err != nil {
						return &pbUser.CommonResp{
							ErrCode: 401,
							ErrMsg:  "发送失败" + err.Error(),
						}, nil
					} else {
						return &pbUser.CommonResp{
							ErrCode: 0,
							ErrMsg:  "成功",
						}, nil
					}

				} else {
					return &pbUser.CommonResp{
						ErrCode: 402,
						ErrMsg:  "邮件服务未打开",
					}, nil
				}
			}
		}
		return &pbUser.CommonResp{
			ErrCode: 0,
			ErrMsg:  "用户未设置邮箱",
		}, nil
	} else {
		return &pbUser.CommonResp{
			ErrCode: 801,
			ErrMsg:  "用户没有设置邮箱",
		}, nil
	}
}

var ct *xkafka.Kafka

func InstallKafkaProduct() *xkafka.Kafka {
	if ct == nil {
		syncx.Once(func() {
			optptr := xkafka.NewDefaultOptions()
			optptr.Name = "product"
			optptr.Addr = config.Config.Kafka.BusinessTop.Addr
			var err error
			ct, err = xkafka.New(optptr, nil, nil)
			if err != nil {
				xlog.CErrorf(err.Error())
				return
			}
		})
	}
	return ct

}
