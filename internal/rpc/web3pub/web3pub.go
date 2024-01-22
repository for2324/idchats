package web3pub

import (
	chat "Open_IM/internal/rpc/msg"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbTask "Open_IM/pkg/proto/task"
	pbweb3pb "Open_IM/pkg/proto/web3pub"
	"Open_IM/pkg/utils"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	twitterscraper "github.com/n0madic/twitter-scraper"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	erc20 "github.com/nattaponra/go-abi/erc20/contract"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type web3pubserver struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func (s *web3pubserver) GetFacebookTimeLine(ctx context.Context, req *pbweb3pb.ThirdPlatformTwitterReq) (*pbweb3pb.ThirdPlatformTwitterRsp, error) {
	//TODO implement me
	var userthird db.UserThird
	userid := req.Userid
	username := req.Username
	if count, err := imdb.IsExistThird("facebook", username); err != nil || count >= 1 {
		return &pbweb3pb.ThirdPlatformTwitterRsp{
			CommonResp: &pbweb3pb.CommonResp{
				ErrCode: 10001,
				ErrMsg:  "this facebook was bind",
			},
		}, nil
	}
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_third").Where("user_id=?", req.Userid).First(&userthird).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userthird.UserId = userid
		err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").Create(&userthird).Error
	} else if err == nil {
		userthird.UserId = userid
		err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").Where("user_id=?", req.Userid).
			Updates(&userthird).Error
	}
	return &pbweb3pb.ThirdPlatformTwitterRsp{
		CommonResp: &pbweb3pb.CommonResp{
			ErrCode: 0,
			ErrMsg:  "验证成功",
		},
	}, nil
}

func NewWeb3PubServer(port int) *web3pubserver {
	log.NewPrivateLog(constant.LogFileName)
	return &web3pubserver{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImWeb3Js,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

type TwitterUser struct {
	Name   string
	Tweets []*twitterscraper.TweetResult
}

func (s *web3pubserver) Run() {
	log.NewInfo("0", "web3pubserver run...")

	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)

	//listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo("0", "listen ok ", address)
	defer listener.Close()
	//grpc server
	var grpcOpts []grpc.ServerOption
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	//User friend related services register to etcd
	pbweb3pb.RegisterWeb3PubServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","),
		rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(),
			s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP,
			s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register web3rpc module  rpc to etcd err"))
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error(), listener)
		return
	}
}

func (s *web3pubserver) InitCrawlingTwitterFollow(ctx context.Context, Operation string, checkUserScreenName string) bool {
	return getRelationShip(Operation, config.Config.OfficialTwitter, checkUserScreenName)
}

// GetUserAuthorizedThirdPlatformList 获取用户已经授权的第三方平台列表
func (s *web3pubserver) GetUserAuthorizedThirdPlatformList(ctx context.Context, req *pbweb3pb.GetUserAuthorizedThirdPlatformListReq) (*pbweb3pb.GetUserAuthorizedThirdPlatformListRsp, error) {
	userId := req.Userid

	if userId == "" {
		return &pbweb3pb.GetUserAuthorizedThirdPlatformListRsp{
			CommonResp: &pbweb3pb.CommonResp{
				ErrCode: 10000,
				ErrMsg:  "userId is nil",
			},
		}, nil
	}

	userThird := db.UserThird{}
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_third").
		Where("user_id = ?", userId).First(&userThird).Error
	if err != nil {
		return &pbweb3pb.GetUserAuthorizedThirdPlatformListRsp{
			CommonResp: &pbweb3pb.CommonResp{
				ErrCode: 10000,
				ErrMsg:  err.Error(),
			},
		}, nil
	}

	platFormList := make([]*pbweb3pb.PlatFormRsp, 0)
	if userThird.Twitter != "" {
		model := pbweb3pb.PlatFormRsp{PlatForm: "twitter", UserId: userThird.Twitter}
		platFormList = append(platFormList, &model)
	}
	return &pbweb3pb.GetUserAuthorizedThirdPlatformListRsp{
		CommonResp: &pbweb3pb.CommonResp{
			ErrCode: 0,
			ErrMsg:  "",
		},
		PlatFormList: platFormList,
	}, nil
}

func (s *web3pubserver) GetTwitterTimeLine(ctx context.Context, req *pbweb3pb.ThirdPlatformTwitterReq) (*pbweb3pb.ThirdPlatformTwitterRsp, error) {
	username := req.Username //twitter 的username
	userid := req.Userid     //用户当前目录的userid
	if count, err := imdb.IsExistThird("twitter", username); err != nil || count >= 1 {
		return &pbweb3pb.ThirdPlatformTwitterRsp{
			CommonResp: &pbweb3pb.CommonResp{
				ErrCode: 10001,
				ErrMsg:  "this twitter was bind",
			},
		}, nil
	}
	user := findUserTweets(username, 2)
	fmt.Println(utils.StructToJsonString(user))
	if len(user.Tweets) == 0 {
		return &pbweb3pb.ThirdPlatformTwitterRsp{
			CommonResp: &pbweb3pb.CommonResp{
				ErrCode: 10000,
				ErrMsg:  "推特没有推文：" + username,
			},
		}, nil
	} else {
		isOksign := false
		for _, value := range user.Tweets {
			compileRegex := regexp.MustCompile("codesign:(.*?)。")
			matchArr := compileRegex.FindStringSubmatch(value.Text)
			if len(matchArr) < 2 {
				continue
			}
			signtxt := matchArr[1]
			fmt.Println(">>>>signtxt:::::", signtxt)
			fmt.Println(">>>>signtxt:::::", req.Nonce)
			//验证签名
			if sigValid := utils.VerifySignature(strings.ToLower(req.Userid), signtxt, req.Nonce); sigValid {
				isOksign = true
				break
			}

		}
		if !isOksign {
			return &pbweb3pb.ThirdPlatformTwitterRsp{
				CommonResp: &pbweb3pb.CommonResp{
					ErrCode: 10001,
					ErrMsg:  "推文不正确",
				},
			}, nil
		} else {
			userBind, err := imdb.HasBindTwitter(userid)
			if err != nil {
				return &pbweb3pb.ThirdPlatformTwitterRsp{
					CommonResp: &pbweb3pb.CommonResp{
						ErrCode: constant.ErrCallback.ErrCode,
						ErrMsg:  err.Error(),
					},
				}, nil
			}
			if userBind {
				return &pbweb3pb.ThirdPlatformTwitterRsp{
					CommonResp: &pbweb3pb.CommonResp{
						ErrCode: 10002,
						ErrMsg:  "禁止重复绑定推特用户",
					},
				}, nil
			}
			//判断是否有其数据， 如果没有这个数据的情况 做插入的操作， 如果存在这个数据就做删除的操作。
			var userthird db.UserThird
			err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").Where("user_id=? ", req.Userid).First(&userthird).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				userthird.UserId = userid
				userthird.Twitter = username
				err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").Create(&userthird).Error
			} else if err == nil {
				userthird.UserId = userid
				userthird.Twitter = username
				err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").Where("user_id=?", req.Userid).
					Updates(&userthird).Error
			}
			err = imdb.UpdateUserInfo(db.User{
				UserID: req.Userid,
				Ex:     fmt.Sprintf(`{"twitter":"%s"}`, username),
			})
			rocksCache.DelUserInfoFromCache(req.Userid)
			chat.UserInfoUpdatedNotification(req.OperatorID, req.Userid, req.Userid)
			// 完成绑定推特任务
			FinishBindTwitterTask(req.OperatorID, req.Userid)
			return &pbweb3pb.ThirdPlatformTwitterRsp{
				CommonResp: &pbweb3pb.CommonResp{
					ErrCode: 0,
					ErrMsg:  "验证成功",
				},
			}, nil
		}
	}

}

func (s *web3pubserver) CheckIsFollowSystemTwitter(ctx context.Context, req *pbweb3pb.CheckUserIsFollowSystemTwitterReq) (*pbweb3pb.CheckUserIsFollowSystemTwitterRsp, error) {
	//查询是否领取任务
	result := &pbweb3pb.CheckUserIsFollowSystemTwitterRsp{CommonResp: &pbweb3pb.CommonResp{}}

	if finish, err := imdb.IsFinishFollowOfficialTwitterTask(req.UserId); err != nil {
		result.CommonResp.ErrCode = constant.ErrDB.ErrCode
		result.CommonResp.ErrMsg = err.Error()
		return result, nil
	} else {
		if finish {
			result.CommonResp.ErrCode = 50001
			result.CommonResp.ErrMsg = "your were followed biubiu twitter "
		}
	}

	//检查是否绑定的twitter
	twitterNameValue, err := imdb.GetUserTwitter(req.UserId)
	if err != nil {
		result.CommonResp.ErrCode = constant.ErrDB.ErrCode
		result.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return result, nil
	}
	if twitterNameValue == "" {
		result.CommonResp.ErrCode = 50003
		result.CommonResp.ErrMsg = "your not  bind twitter "
		return result, nil
	}
	//检查该twitterValue 是否是已经关注的：
	if !s.InitCrawlingTwitterFollow(ctx, req.OperatorID, twitterNameValue) {
		result.CommonResp.ErrCode = 50002
		result.CommonResp.ErrMsg = "暂未关注官方twitter"
		return result, nil
	} else {
		//已经有关注 旧需要去新增事件
		FinishFollowOfficialTwitterTask(req.OperatorID, req.UserId)
	}
	return result, nil
}

func (s *web3pubserver) CheckIsHaveNftRecvID(ctx context.Context, req *pbweb3pb.CheckIsHaveNftRecvIDReq) (*pbweb3pb.CheckIsHaveNftRecvIDResp, error) {
	haveNft := CheckIsHaveNftRecvID(req.UserId)
	return &pbweb3pb.CheckIsHaveNftRecvIDResp{
		CommonResp: &pbweb3pb.CommonResp{},
		HaveNft:    haveNft,
	}, nil
}

func (s *web3pubserver) CheckIsHaveGuanFangNftRecvID(ctx context.Context, req *pbweb3pb.CheckIsHaveGuanFangNftRecvIDReq) (*pbweb3pb.CheckIsHaveGuanFangNftRecvIDResp, error) {
	haveNft := CheckHeadIsOfficialNftContract(req.UserId)
	return &pbweb3pb.CheckIsHaveGuanFangNftRecvIDResp{
		CommonResp: &pbweb3pb.CommonResp{},
		HaveNft:    haveNft,
	}, nil
}

func FinishFollowOfficialTwitterTask(OperationID, userId string) {
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, OperationID)
	if etcdConn == nil {
		return
	}
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.FinishFollowOfficialTwitterTask(context.Background(), &pbTask.FinishFollowOfficialTwitterTaskReq{
		OperationID: OperationID,
		UserId:      userId,
	})
	if err != nil {
		log.NewError(OperationID, "FinishFollowOfficialTwitterTask failed ", err.Error())
		return
	}
	if respPb.CommonResp.ErrCode != 0 {
		log.NewDebug(OperationID, "FinishFollowOfficialTwitterTask failed ", respPb.CommonResp.ErrMsg)
		return
	}
	log.NewInfo(OperationID, "FinishFollowOfficialTwitterTask success ")
}

func FinishBindTwitterTask(OperationID, userId string) {
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, OperationID)
	if etcdConn == nil {
		return
	}
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.FinishBindTwitterTask(context.Background(), &pbTask.FinishBindTwitterTaskReq{
		OperationID: OperationID,
		UserId:      userId,
	})
	if err != nil {
		log.NewError(OperationID, "FinishBindTwitterTask failed ", err.Error())
		return
	}
	if respPb.CommonResp.ErrCode != 0 {
		log.NewDebug(OperationID, "FinishBindTwitterTask failed ", respPb.CommonResp.ErrMsg)
		return
	}
	log.NewInfo(OperationID, "FinishBindTwitterTask success ")
}

func (s *web3pubserver) GetEthTxIDTaskRpc(ctx context.Context, req *pbweb3pb.EthRpcTxIDReq) (*pbweb3pb.EthRpcTxIDResp, error) {
	if config.Config.IsPublicEnv && (req.ChainID == "5" || req.ChainID == "80001" || req.ChainID == "97") {
		return nil, errors.New("禁止在正是服务器转测试")
	}
	rpcList := config.GetRpcFromChainID(req.ChainID)
	resp := new(pbweb3pb.EthRpcTxIDResp)
	resp.CommonResp = new(pbweb3pb.CommonResp)
	for _, rpcAddress := range rpcList {
		client2, err := ethclient.Dial(rpcAddress)
		if err != nil {
			continue
		}
		txHash := common.HexToHash(req.TxID)
		ethTxReceipt, err := client2.TransactionReceipt(context.Background(), txHash)
		if (err != nil && err != ethereum.NotFound) || ethTxReceipt == nil {
			continue
		}
		if ethTxReceipt.Status != 1 {
			break
		}
		ethRtx, _, err := client2.TransactionByHash(context.Background(), txHash)
		if err != nil || ethRtx == nil {
			continue
		}
		chainID, err := client2.NetworkID(context.Background())
		fromeaddress := ""
		contractAddress := ""
		from, err := types.Sender(types.LatestSignerForChainID(chainID), ethRtx)
		fromeaddress = from.String()
		toAddr := ethRtx.To()
		if toAddr != nil {
			contractAddress = toAddr.String()
		}
		if !strings.EqualFold(contractAddress, config.GetConfigContractFromChainID(req.ChainID)) {
			break
		}

		t := &pbweb3pb.EthRpcTxIDResp{
			CommonResp:      new(pbweb3pb.CommonResp),
			TransactionHash: strings.ToLower(ethRtx.Hash().Hex()),
			FromAddress:     fromeaddress,
			ToAddress:       "",
			Value:           0.0,
			Status:          int32(ethTxReceipt.Status),
			ContractAddress: strings.ToLower(ethRtx.To().String()),
		}
		if contractAddress != "" && len(ethRtx.Data()) > 0 {
			inputData := hexutil.Bytes(ethRtx.Data())
			inputStr := inputData.String()
			unpackData, method, err := UnpackInput(inputStr, erc20.ContractMetaData.ABI)
			if method == "transfer" && err == nil && len(unpackData) == 2 {
				p1, ok := unpackData[0].(common.Address)

				if ok {
					t.ToAddress = p1.Hex()
				}
				p2, ok := unpackData[1].(*big.Int)
				if ok {
					nValue, err := db.DB.Rc.Fetch("token:"+req.ChainID+":"+contractAddress, time.Hour*24*30, func() (string, error) {
						tokenInfo, _ := erc20.NewContract(*ethRtx.To(), client2)
						nValueTemp, err := tokenInfo.Decimals(nil)
						return strconv.Itoa(int(nValueTemp)), err
					})
					if err == nil {
						t.Decimals = uint32(utils.StringToInt32(nValue))
						amountStr := p2.String()
						amountFloat, _ := strconv.ParseFloat(amountStr, 64)
						amountFloat = amountFloat / math.Pow10(int(t.Decimals))
						t.Value = amountFloat
						return t, nil
					}
				}
			}
		}
	}
	resp.CommonResp.ErrCode = constant.ErrInternal.ErrCode
	resp.CommonResp.ErrMsg = "无法获取交易记录"
	return resp, nil
}

func CheckIsHaveNftRecvID(userId string) bool {
	dbuser, _ := rocksCache.GetUserBaseInfoFromCache(userId)
	dbuserList := strings.Split(dbuser.TokenContractChain, "&")
	if len(dbuserList) < 2 {
		return false
	}
	if strings.EqualFold(dbuserList[0], "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee") {
		return false
	}
	PostCheckData, _ := json.Marshal(&api.RequestTokenIdReq{
		ChainID:         dbuserList[1],
		TokenID:         dbuser.TokenId,
		ContractAddress: dbuserList[0],
	})
	resultByte, err := utils.HttpPost(config.Config.EnsPostCheck.Url+"/graph/tokenOwnerAddress",
		"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, PostCheckData)
	if err != nil {
		return false
	}

	if err == nil {
		var resultData api.RequestTokenIdResp
		json.Unmarshal(resultByte, &resultData)
		if strings.ToLower(resultData.TokenOwnerAddress) != strings.ToLower(userId) {
			return false
		}
	}
	return true
}

func CheckHeadIsOfficialNftContract(address string) bool {
	dbuseradd, err := imdb.GetUserByUserID(address)
	if err == nil && dbuseradd.TokenContractChain != "" {
		stArray := strings.Split(dbuseradd.TokenContractChain, "&")
		if len(stArray) < 2 {
			return false
		}
		PostCheckData, _ := json.Marshal(&api.RequestTokenIdReq{
			ChainID:         stArray[1],
			TokenID:         dbuseradd.TokenId,
			ContractAddress: stArray[0],
		})
		resultByte, err := utils.HttpPost(config.Config.EnsPostCheck.Url+"/graph/tokenOwnerAddress",
			"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, PostCheckData)
		if err == nil {
			var resultData api.RequestTokenIdResp
			json.Unmarshal(resultByte, &resultData)
			if strings.EqualFold(resultData.TokenOwnerAddress, address) {
				return true
			}
		}
	}
	return false
}

// UnpackInput 解析输入
func UnpackInput(txInput string, abiJson string) ([]interface{}, string, error) {
	var data = make([]interface{}, 0)
	abi, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		return data, "", err
	}
	if len(txInput) > 10 {
		// decode txInput method signature
		decodeSign, err := hex.DecodeString(txInput[2:10])
		if err != nil {
			return data, "", err
		}
		// recover Method from signature and ABI
		method, err := abi.MethodById(decodeSign)
		if err != nil {
			return data, "", err
		}
		//decodeData, err := hex.DecodeString(txInput[2:])
		decodedData, err := hex.DecodeString(txInput[10:])
		if err != nil {
			return data, "", err
		}
		// unpack method inputs
		data, err = method.Inputs.Unpack(decodedData)
		return data, method.Name, err
	}
	return data, "", errors.New("数据：" + txInput + "解析失败")
}

func (s *web3pubserver) PostGamingStatus(ctx context.Context, req *pbweb3pb.UserGameReq) (result *pbweb3pb.UserGameResp, err error) {
	//检查user 是否含有nft的头像：
	result = new(pbweb3pb.UserGameResp)
	result.CommonResp = new(pbweb3pb.CommonResp)
	gameList, err := imdb.GetGameListFromDB(req.GameID)
	if err != nil && len(gameList) == 0 {
		result.CommonResp.ErrCode = constant.ErrInternal.ErrCode
		result.CommonResp.ErrMsg = "游戏类型不存在"
		return result, nil
	}
	jsonCondition := gameList[0].GameCondition
	gameName := gameList[0].GameName
	var dbGameCondition db.GameConditionData
	json.Unmarshal(utils.String2bytes(jsonCondition), &dbGameCondition)
	if dbGameCondition.IsHeadNft && !CheckIsHaveNftRecvID(req.UserID) {
		result.CommonResp.ErrCode = constant.ErrInternal.ErrCode
		result.CommonResp.ErrMsg = "head image not set nft address ,can't start game"
		return result, nil
	}
	if dbGameCondition.IsOfficialNft && !checkIsHaveGuangFangErc721Nft(req.UserID) {
		result.CommonResp.ErrCode = constant.ErrInternal.ErrCode
		result.CommonResp.ErrMsg = "have not official nft address ,can't start game"
		return result, nil
	}
	if req.Status == 1 { //游戏开始
		//游戏开始操作
		timeNow := time.Now().UnixMilli()
		err = imdb.UpdateGameStatus(req.UserID,
			gameName, utils.StringToInt32(req.GameID), req.Status, req.Ip, req.UserAgent, 0, timeNow)
		if err != nil {
			result.CommonResp.ErrCode = constant.ErrInternal.ErrCode
			result.CommonResp.ErrMsg = "未持有nft， 不能进行游戏"
			return result, nil
		} else {
			result.StartTime = timeNow
		}
		return result, nil

	} else if req.Status == 2 {
		var userDataGameStatus db.UserGameScore
		err := db.DB.MysqlDB.DefaultGormDB().
			Table("user_game_score").
			Where("user_id=?", req.UserID).
			Order("created_at desc").Limit(1).Find(&userDataGameStatus).Error
		if err != nil {
			result.CommonResp.ErrCode = constant.ErrInternal.ErrCode
			result.CommonResp.ErrMsg = "数据错误"
			return result, nil
		}
		inputscore := req.Score
		timeNow := time.Now().UnixMilli()
		if gameList[0].GameVerify == 1 {
			_, distance := CheckSumDistanceSource(timeNow - userDataGameStatus.StartTime)
			inputscore = distance
			if distance > req.Score {
				inputscore = req.Score
			}
		}
		err = imdb.UpdateGameStatus(req.UserID, gameName, utils.StringToInt32(req.GameID),
			req.Status, req.Ip, req.UserAgent, inputscore, timeNow)
		if err != nil {
			result.CommonResp.ErrCode = constant.ErrInternal.ErrCode
			result.CommonResp.ErrMsg = err.Error()
			return result, nil
		} else {
			result.CommonResp.ErrCode = 0
			result.CommonResp.ErrMsg = ""
			return result, nil
		}
	}
	return nil, nil
}

func (s *web3pubserver) GetGamingRankStatus(ctx context.Context, req *pbweb3pb.UserGameRankListReq) (*pbweb3pb.UserGameRankListResp, error) {
	result := new(pbweb3pb.UserGameRankListResp)
	result.CommonResp = new(pbweb3pb.CommonResp)
	dbList, UserRank, err := imdb.GetRankLink(req.GameID, req.UserID)
	if err != nil {
		result.CommonResp.ErrCode = constant.ErrDB.ErrCode
		result.CommonResp.ErrMsg = err.Error()
		return result, nil
	}
	for _, value := range dbList {
		result.UserRankInfo = append(result.UserRankInfo, &pbweb3pb.UserGameScore{
			RankIndex:          value.Index,
			UserID:             value.UserID,
			Nickname:           value.Nickname,
			FaceURL:            value.FaceURL,
			Score:              float64(value.Score),
			Reward:             value.Reward,
			TokenContractChain: value.TokenContractChain,
		})
	}
	fmt.Println(">>>>>>>>>>>>>>", utils.StructToJsonString(UserRank))
	if UserRank != nil {
		result.UserSelfRankInfo = &pbweb3pb.UserGameScore{
			RankIndex:          UserRank.Index,
			UserID:             UserRank.UserID,
			Nickname:           UserRank.Nickname,
			FaceURL:            UserRank.FaceURL,
			Score:              float64(UserRank.Score),
			Reward:             UserRank.Reward,
			TokenContractChain: UserRank.TokenContractChain,
		}
	}

	return result, nil
}
func checkIsHaveGuangFangErc721Nft(address string) bool {
	PostCheckData, _ := json.Marshal(&api.CheckHaveNftReq{
		Address: address,
	})
	resultByte, err := utils.HttpPost(config.Config.EnsPostCheck.Url+"/graph/checkAddress",
		"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, PostCheckData)
	if err == nil {
		var resultData api.BalanceOfCount
		json.Unmarshal(resultByte, &resultData)
		if resultData.BalanceOf > 0 {
			return true
		}
	}
	return false
}
func CheckSumDistanceSource(timeCost int64) (float64, float64) {
	if timeCost > 24*60*60*1000 { //24小时*60分钟*60秒
		log.NewError("存在错误数据", "作弊数据")
		return 0, 0
	}
	var initialVelocity float64 = 10
	//0.002
	acceleration := 2.0
	elapsedTime := float64(timeCost)
	distanceRan := 0.0
	msPerFrame := float64(1000 / 60)
	for t := float64(0); t < elapsedTime; t += msPerFrame {
		time := math.Min(msPerFrame, elapsedTime-t)
		distanceRan += initialVelocity * 0.0167
		initialVelocity += acceleration * time / 1000
		if initialVelocity >= 100 {
			initialVelocity = 100
		}
	}

	currentSpeed := initialVelocity
	return currentSpeed, distanceRan
}
func getRelationShip(Operation string, officialName string, username string) bool {
	var resutlbyte []byte
	var err error
	if config.Config.OpenNetProxy.OpenFlag == true {
		proxyAddress, _ := url.Parse("http://proxy.idchats.com:7890")
		resutlbyte, err = utils.HttpGetWithHeaderWithProxy(fmt.Sprintf("https://api.twitter.com/1.1/friendships/show.json?source_screen_name=%s&target_screen_name=%s",
			officialName, username), map[string]string{
			"User-Agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"authorization": fmt.Sprintf("Bearer %s", config.Config.Web3thirdpath.TwitterBearToken),
		}, http.ProxyURL(proxyAddress))

	} else {
		resutlbyte, err = utils.HttpGetWithHeader(fmt.Sprintf("https://api.twitter.com/1.1/friendships/show.json?source_screen_name=%s&target_screen_name=%s",
			officialName, username), map[string]string{
			"User-Agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"authorization": fmt.Sprintf("Bearer %s", config.Config.Web3thirdpath.TwitterBearToken),
		})
	}
	if err != nil {
		log.NewInfo(Operation, utils.GetSelfFuncName(), err.Error())
	} else {
		var TwitterRelationShipStruct TwitterRelationShip
		if err = json.Unmarshal(resutlbyte, &TwitterRelationShipStruct); err == nil {
			return TwitterRelationShipStruct.Relationship.Source.FollowedBy
		}
	}
	return false
}

type TwitterRelationShip struct {
	Relationship struct {
		Source struct {
			Id                   int64       `json:"id"`
			IdStr                string      `json:"id_str"`
			ScreenName           string      `json:"screen_name"`
			Following            bool        `json:"following"`
			FollowedBy           bool        `json:"followed_by"`
			LiveFollowing        bool        `json:"live_following"`
			FollowingReceived    interface{} `json:"following_received"`
			FollowingRequested   interface{} `json:"following_requested"`
			NotificationsEnabled interface{} `json:"notifications_enabled"`
			CanDm                bool        `json:"can_dm"`
			Blocking             interface{} `json:"blocking"`
			BlockedBy            interface{} `json:"blocked_by"`
			Muting               interface{} `json:"muting"`
			WantRetweets         interface{} `json:"want_retweets"`
			AllReplies           interface{} `json:"all_replies"`
			MarkedSpam           interface{} `json:"marked_spam"`
		} `json:"source"`
		Target struct {
			Id                 int64       `json:"id"`
			IdStr              string      `json:"id_str"`
			ScreenName         string      `json:"screen_name"`
			Following          bool        `json:"following"`
			FollowedBy         bool        `json:"followed_by"`
			FollowingReceived  interface{} `json:"following_received"`
			FollowingRequested interface{} `json:"following_requested"`
		} `json:"target"`
	} `json:"relationship"`
}

func findUserTweets(username string, count int) (userTweet *TwitterUser) {
	userTweet = new(TwitterUser)
	userTweet.Name = username
	done := make(chan bool)
	go func() {
		defer close(done) // 在子 goroutine 结束时关闭 done channel
		scraper := twitterscraper.New()
		if config.Config.OpenNetProxy.OpenFlag {
			err := scraper.SetProxy("http://proxy.idchats.com:7890")
			if err != nil {
				fmt.Println("无法翻墙")
			}
		}
		for tweet := range scraper.GetTweets(context.Background(), username, count) {
			if tweet.Error != nil {
				done <- true
			}
			userTweet.Tweets = append(userTweet.Tweets, tweet)
		}
	}()
	select {
	case <-done:
		fmt.Println("data end ")
		return
	case <-time.After(20 * time.Second):
		fmt.Println("time end ")
		close(done)
		return
	}
}

// 判断URL字符串是否包含HTTP或HTTPS前缀
func hasHTTPPrefix(urlStr string) bool {
	return len(urlStr) > 7 && (urlStr[:7] == "http://" || urlStr[:8] == "https://")
}
func (s *web3pubserver) CheckDnsDomainHadParseBiuBiuTxt(ctx context.Context, req *pbweb3pb.CheckDomainHadParseTxtReq) (*pbweb3pb.CheckDomainHadParseTxtResp, error) {
	result := new(pbweb3pb.CheckDomainHadParseTxtResp)
	result.CommonResp = new(pbweb3pb.CommonResp)
	// 在URL字符串前添加协议前缀
	checkDnsName := req.DnsDomain
	if !hasHTTPPrefix(req.DnsDomain) {
		checkDnsName = "http://" + req.DnsDomain
	}
	urldata, err := url.Parse(checkDnsName)
	if err != nil {
		result.CommonResp.ErrCode = 501
		result.CommonResp.ErrMsg = "无法解析域名" + checkDnsName
		return result, nil
	}
	domainName := urldata.Hostname()
	domainList := strings.Split(domainName, ".")
	if len(domainList) < 2 {
		result.CommonResp.ErrCode = 502
		result.CommonResp.ErrMsg = "无法解析域名" + checkDnsName
		return result, nil
	}
	needCheckName := domainList[len(domainList)-2] + "." + domainList[len(domainList)-1]
	resultData, _ := net.LookupTXT(needCheckName)
	for _, value := range resultData {
		if strings.EqualFold(value, req.UserId) {
			err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
				var tempUserThird db.UserThird
				if err := tx.Table("user_third").Where("user_id=?", req.UserId).First(&tempUserThird).Error; errors.Is(err, gorm.ErrRecordNotFound) {
					if err = tx.Table("user_third").Create(&db.UserThird{
						UserId: req.UserId,
					}).Error; err != nil {
						return err
					}
				}
				if err = tx.Table("user_third").Where("user_id=?", req.UserId).Updates(map[string]interface{}{
					"dns_domain":        req.DnsDomain,
					"dns_domain_verify": 1}).Error; err != nil {
					return err
				}
				return nil
			})
			result.CommonResp.ErrCode = 0
			result.CommonResp.ErrMsg = ""
			return result, nil
		}
	}
	result.CommonResp.ErrCode = 503
	result.CommonResp.ErrMsg = "你的解析不存在"
	return result, nil
}
