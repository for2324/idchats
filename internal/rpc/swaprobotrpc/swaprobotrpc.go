package swaprobotrpc

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbSwapRobot "Open_IM/pkg/proto/swaprobot"
	"Open_IM/pkg/utils"
	"Open_IM/pkg/walletRdbService"
	client2 "Open_IM/pkg/walletRdbService/swaprobotservice"
	"context"
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/go-redsync/redsync/v4"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/metadata"
	"k8s.io/utils/strings/slices"
	"net"
	"strconv"
	"strings"
	"time"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type SwapRobotServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewSwapRobotServer(port int) *SwapRobotServer {
	log.NewPrivateLog(constant.LogFileName)
	return &SwapRobotServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.SwapRobotPort,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *SwapRobotServer) Run() {
	log.NewInfo("0", "SwapRobotServer run...")

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
	pbSwapRobot.RegisterSwaprobotServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register friend module  rpc to etcd err"))
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error(), listener)
		return
	}
}
func (s *SwapRobotServer) RecodeSwapStatus(ctx context.Context, req *pbSwapRobot.SwapRobotOrderStatusReq) (*pbSwapRobot.SwapRobotOrderStatusResp, error) {
	//TODO implement me
	panic("implement me")

}

func (s *SwapRobotServer) RecordSwapInfo(ctx context.Context, req *pbSwapRobot.SwapRecordInfoReq) (*pbSwapRobot.SwapRecordInfoResp, error) {
	resultPb := new(pbSwapRobot.SwapRecordInfoResp)
	resultPb.CommonResp = new(pbSwapRobot.CommonResp)

	resultPb.CommonResp.ErrCode = 0
	resultPb.CommonResp.ErrMsg = ""
	switch req.Method {
	case "createTask":
		err := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
			err := tx.Table("swap_robot_task").Create(&db.RoBotTask{
				ID:          0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UserID:      req.UserID,
				Address:     req.Address,
				OrdID:       req.OrdID,
				FromSymbol:  req.FromSymbol,
				ToSymbol:    req.ToSymbol,
				Amount:      req.Amount,
				Tp:          req.Tp,
				Sl:          req.Sl,
				OrderStatus: req.OrdStatus,
				MinimumOut:  req.MinimumOut,
				DeadlineDay: req.DeadlineDay,
			}).Error
			if err == nil {
				err = tx.Table("swap_robot_task_log").Create(&db.RoBotTaskLog{
					OrdID:       req.OrdID,
					OrderStatus: req.OrdStatus,
					Method:      req.Method,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}).Error
			}
			return err
		})
		if err == nil {
			resultPb.OrderID = req.OrdID
		}
		return resultPb, nil
	case "editTask":
		err := db.DB.MysqlDB.DefaultGormDB().Table("swap_robot_task").
			Where("ord_id=?", req.OrdID).
			Updates(map[string]interface{}{
				"updated_at":  time.Now(),
				"user_id":     req.UserID,
				"address":     req.Address,
				"ord_id":      req.OrdID,
				"from_symbol": req.FromSymbol,
				"to_symbol":   req.ToSymbol,
				"amount":      req.Amount,
				"tp":          req.Tp,
				"sl":          req.Sl,
			}).Error
		if err != nil {
			resultPb.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resultPb.CommonResp.ErrMsg = err.Error()
		}
	}
	return resultPb, nil
}
func (s *SwapRobotServer) BotOperation(ctx context.Context, req *pbSwapRobot.BotOperationReq) (*pbSwapRobot.BotOperationResp, error) {

	if req.UserID == "" && !slices.Contains([]string{"exploredPair", "tokenPrice",
		"quote", "newTokens", "tokensPrice"}, req.Method) {
		return &pbSwapRobot.BotOperationResp{CommonResp: &pbSwapRobot.CommonResp{
			ErrCode: constant.ErrInternal.ErrCode,
			ErrMsg:  "userId is empty",
		}}, nil
	}
	var robot *db.Robot
	var err error
	if req.UserID != "" {
		robot, err = imdb.GetUserRobotInfo(req.UserID)
		if err != nil {
			return &pbSwapRobot.BotOperationResp{
				CommonResp: &pbSwapRobot.CommonResp{
					ErrCode: constant.ErrInternal.ErrCode,
					ErrMsg:  "failed to get user robot info",
				},
			}, nil
		}
		if robot == nil {
			return &pbSwapRobot.BotOperationResp{CommonResp: &pbSwapRobot.CommonResp{
				ErrCode: constant.ErrInternal.ErrCode,
				ErrMsg:  "user robot is empty",
			}}, nil
		}
	} else {
		robot = new(db.Robot)
	}
	if config.Config.WalletService.OpenFlag && req.UserID != "" {
		client, err := walletRdbService.GetRdbService()
		if err != nil {
			return &pbSwapRobot.BotOperationResp{CommonResp: &pbSwapRobot.CommonResp{
				ErrCode: constant.ErrInternal.ErrCode,
				ErrMsg:  "wallet middleware con't connect",
			}}, nil
		}
		reqHeader := metadata.New(map[string]string{"userid": req.UserID})
		ctx2 := metadata.NewOutgoingContext(context.Background(), reqHeader)
		resultData, err := client.GetMnemonicFromMemory(ctx2, &client2.RequestUserID{
			UserID: req.UserID,
		})
		if err != nil {
			return &pbSwapRobot.BotOperationResp{CommonResp: &pbSwapRobot.CommonResp{
				ErrCode: constant.ErrInternal.ErrCode,
				ErrMsg:  "user robot is empty",
			}}, nil
		}
		if resultData.BaseResp.StatusCode != 0 {
			return &pbSwapRobot.BotOperationResp{CommonResp: &pbSwapRobot.CommonResp{
				ErrCode: constant.ErrInternal.ErrCode,
				ErrMsg:  "托管过期",
			}}, nil
		}
		robot.Mnemonic = resultData.Mnemonic
	}
	parentUserId, err := imdb.GetParentUserID(req.UserID)
	if err != nil {
		log.Info("查询当前的sql错误：:", err.Error())
	}
	var resp *RobotRunTaskResp
	postBody := req.Params
	resp, err = requestUriRobot2(req.OperatorID, robot, req.Method, postBody, req.BiBotKey, parentUserId)
	if err != nil {
		return &pbSwapRobot.BotOperationResp{CommonResp: &pbSwapRobot.CommonResp{
			ErrCode: constant.ErrInternal.ErrCode,
			ErrMsg:  err.Error(),
		}}, nil
	}
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return &pbSwapRobot.BotOperationResp{CommonResp: &pbSwapRobot.CommonResp{
			ErrCode: constant.ErrInternal.ErrCode,
			ErrMsg:  err.Error(),
		}}, nil
	}
	return &pbSwapRobot.BotOperationResp{CommonResp: &pbSwapRobot.CommonResp{}, Data: data}, nil
	// switch req.Method {
	// case "cancelTask":
	// 	//发送请求取消订单
	// 	if req.OrdID == "" {
	// 		return &pbSwapRobot.SwapRecordInfoResp{CommonResp: &pbSwapRobot.CommonResp{
	// 			ErrCode: constant.ErrInternal.ErrCode,
	// 			ErrMsg:  "取消的订单号不能为空",
	// 		}}, nil
	// 	}
	// 	if _, err := imdb.GetOrderInfo(req.UserID, req.OrdID); err != nil {
	// 		return &pbSwapRobot.SwapRecordInfoResp{CommonResp: &pbSwapRobot.CommonResp{
	// 			ErrCode: constant.ErrInternal.ErrCode,
	// 			ErrMsg:  "无法操作不存在的订单",
	// 		}}, nil

	// 	}
	// 	mapValue["ordId"] = req.OrdID
	// 	return cancelTask(req.UserID, req.OrdID, robot, req.Method, mapValue, req.OperationID)
	// case "createTask":
	// 	checkAmountStr, _ := decimal.NewFromString(req.Amount)
	// 	if checkAmountStr.LessThanOrEqual(decimal.NewFromInt(0)) {
	// 		return &pbSwapRobot.SwapRecordInfoResp{CommonResp: &pbSwapRobot.CommonResp{
	// 			ErrCode: constant.ErrInternal.ErrCode,
	// 			ErrMsg:  "金额不能低于0",
	// 		}}, nil
	// 	}
	// 	mapValue["fromSymbol"] = req.FromSymbol
	// 	mapValue["toSymbol"] = req.ToSymbol
	// 	mapValue["amount"], _ = checkAmountStr.Float64()
	// 	str, _ := decimal.NewFromString(req.Tp)
	// 	mapValue["tp"], _ = str.Float64()
	// 	str, _ = decimal.NewFromString(req.Sl)
	// 	mapValue["sl"], _ = str.Float64()
	// 	str, _ = decimal.NewFromString(req.MinimumOut)
	// 	mapValue["minimumOut"], _ = str.Float64()
	// 	mapValue["deadlineDay"] = utils.StringToInt64(req.DeadlineDay)
	// 	return createTask(req.UserID, robot, req.Method, mapValue, req.OperationID)
	// case "getTask":
	// 	if req.OrdID != "" {
	// 		mapValue["ordId"] = req.OrdID
	// 		mapValue["searchBy"] = "ordId"
	// 	} else {
	// 		mapValue["searchBy"] = "address"
	// 	}
	// 	return getTaskList(req.UserID, robot, req.Method, mapValue, req.OperationID)
	// default:
	// 	if robotReBack, err := requestUriRobot(robot, req.Method, mapValue, req.OperationID); err != nil {
	// 		return &pbSwapRobot.SwapRecordInfoResp{CommonResp: &pbSwapRobot.CommonResp{
	// 			ErrCode: constant.ErrInternal.ErrCode,
	// 			ErrMsg:  err.Error(),
	// 		}}, nil
	// 	} else {
	// 		log.Info(req.OperationID, "当前机器人返回数据:", utils.StructToJsonString(robotReBack))
	// 		return robotReBack, nil
	// 	}
	// }

}

type Param struct {
	FromSymbol  string  `json:"fromSymbol"`
	Amount      float64 `json:"amount"`
	ToSymbol    string  `json:"toSymbol"`
	Tp          float64 `json:"tp"`
	Sl          float64 `json:"sl"`
	OrdStatus   string  `json:"ordStatus"`
	MinimumOut  float64 `json:"minimumOut"`
	DeadlineDay int     `json:"deadlineDay"`
}
type SwapBotInfoParam struct {
	OrdId     string `json:"ordId"`
	Method    string `json:"method"`
	OrdStatus string `json:"ordStatus"`
	Params    Param  `json:"params"`
}
type RobotRunTaskPostReq struct {
	PrivateKey       string                 `json:"privateKey"`
	Address          string                 `json:"address"`
	SenderAddress    string                 `json:"senderAddress"`
	RecipientAddress string                 `json:"recipientAddress"`
	Phrase           string                 `json:"phrase"`
	Params           map[string]interface{} `json:"params"`
	TimeStamp        int64                  `json:"timeStamp"`
	Method           string                 `json:"method"`
	FeeRate          float64                `json:"feeRate"`
	FeeAddressMap    map[string]string      `json:"feeAddressMap"`
	UserID           string                 `json:"userID"`
	ApiKey           string                 `json:"apiKey"`
	ParentUserID     string                 `json:"parentUserID"`
}

type RobotRunTaskResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type JsonTaskDetail struct {
	Address            string      `json:"address"`
	Method             string      `json:"method"`
	OrdId              string      `json:"ordId"`
	OrdStatus          string      `json:"ordStatus"`
	Params             Param       `json:"params"`
	TaskStartTimestamp string      `json:"taskStartTimestamp"`
	TradeMsg           interface{} `json:"tradeMsg"`
}

func requestUriRobot2(
	operationID string,
	robot *db.Robot, method string, body []byte,
	apiKey, parentUserID string) (*RobotRunTaskResp, error) {
	params := map[string]interface{}{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		return nil, err
	}
	feeRate := 0.0
	if method == "snipeSwap" || method == "stableSwap" || method == "swap" {
		if apiKey == "" {
			tradeFee, sniperFee := rocksCache.GetCurrentFee(robot.UserID, "", nil)
			if method == "snipeSwap" {
				feeRate = sniperFee
			} else {
				feeRate = tradeFee
			}
		} else {
			tempUserID, _ := imdb.GetApiKeyUserRobotApi(apiKey)
			if tempUserID == nil {
				if method == "snipeSwap" {
					feeRate = 0.009
				} else {
					feeRate = 0.0018
				}
			} else {
				if method == "snipeSwap" {
					feeRate, _ = convertor.ToFloat(tempUserID.SniperFee)
				} else {
					feeRate, _ = convertor.ToFloat(tempUserID.TradeFee)
				}
			}
		}
	}

	postBody := RobotRunTaskPostReq{
		Phrase:       robot.Mnemonic,
		Params:       params,
		TimeStamp:    time.Now().Unix(),
		Method:       method, //某种类型的参数
		FeeRate:      feeRate,
		UserID:       robot.UserID,
		ApiKey:       apiKey,
		ParentUserID: parentUserID,
		FeeAddressMap: map[string]string{"btc": "bc1qyszn8mykfn52ejmqsw0qnthtrag5kxk9k7p2e9",
			"eth": "0xf8893d45bb5052fea90711fefe7c478167348c64"},
	}
	postBodyBytes, err := json.Marshal(postBody)
	if err != nil {
		return nil, err
	}
	if !config.Config.IsPublicEnv {
		log.NewInfo(operationID, "请求机器人参数：", string(postBodyBytes))
	}
	respData, err := utils.HttpPost(config.Config.UniswapRobot.BibotUri+"/api/carry_out", "", map[string]string{}, postBodyBytes)
	if err != nil {
		log.NewInfo(operationID, "请求机器人失败", config.Config.UniswapRobot.Uri, err.Error())
		return nil, err
	} else {
		log.NewInfo(operationID, "请求返回数据：", len(respData))
	}

	robotRunTaskRespData := new(RobotRunTaskResp)
	err = json.Unmarshal(respData, robotRunTaskRespData)
	if robotRunTaskRespData.Code != 200 {
		return nil, errors.New(robotRunTaskRespData.Msg)
	}
	return robotRunTaskRespData, err
}
func (s *SwapRobotServer) FinishTaskToGetReword(ctx context.Context, req *pbSwapRobot.BotSwapTradeReq) (*pbSwapRobot.BotSwapTradeResp, error) {
	//bibot 通知服务，告速本服务 已经完成的任务id 以及所获得积分，统计结果内容。
	//产生交易量 获得多少。

	log.NewInfo("内容数据", utils.StructToJsonString(req))
	parentUserId, err := imdb.GetParentUserID(req.UserID)
	mutexname := "trade_volume:" + req.UserID
	rs := db.DB.Pool
	mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))

	if err := mutex.LockContext(ctx); err != nil {
		return &pbSwapRobot.BotSwapTradeResp{CommonResp: &pbSwapRobot.CommonResp{
			ErrCode: 1,
			ErrMsg:  "正在交易",
		}}, nil
	}
	defer mutex.UnlockContext(ctx)
	if req.UserID != parentUserId && parentUserId != "" {
		mutexname2 := "trade_volume:" + parentUserId
		rs := db.DB.Pool
		mutex2 := rs.NewMutex(mutexname2, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
		if err := mutex2.LockContext(ctx); err != nil {
			return &pbSwapRobot.BotSwapTradeResp{CommonResp: &pbSwapRobot.CommonResp{
				ErrCode: 1,
				ErrMsg:  "正在交易",
			}}, nil
		}
		defer mutex2.UnlockContext(ctx)
	}

	//	apiKeyUserID, err := imdb.GetApiKeyUserID(req.ApiKey)
	err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		//自己的交易量：
		var dbUserHistoryTotal db.UserHistoryTotal
		err := tx.Table("user_history_total").Where("user_id=?", req.UserID).First(&dbUserHistoryTotal).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		isExist := true
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isExist = false
		}
		dbUserHistoryTotal.UserID = req.UserID

		oldPending, _ := decimal.NewFromString(dbUserHistoryTotal.Pending)
		tradeVolumeUsdDecimal, _ := decimal.NewFromString(req.TradeVolumeUsd)
		//新增的奖励为交易量的万分之8,换算后的资产比如1000000 = 1u ，那么存储的量是1 不是1000000
		addPendingScore := tradeVolumeUsdDecimal.Shift(-6).
			Mul(decimal.NewFromFloat(0.0008))
		newPending := oldPending.Add(addPendingScore)

		dbUserHistoryTotal.Pending = newPending.String()
		oldTotalTrade, _ := decimal.NewFromString(dbUserHistoryTotal.TotalTradeVolume)
		totalTrade, _ := decimal.NewFromString(req.TradeVolumeUsd)
		totalTrade = totalTrade.Add(oldTotalTrade)
		dbUserHistoryTotal.TotalTradeVolume = totalTrade.String()

		if !isExist {
			err = tx.Table("user_history_total").Where("user_id=?", req.UserID).Create(dbUserHistoryTotal).Error
			if err != nil {
				return err
			}
		} else {
			err = tx.Table("user_history_total").Where("user_id=?", req.UserID).Updates(map[string]interface{}{
				"pending":            newPending.String(),
				"total_trade_volume": totalTrade.String(),
			}).Error
			if err != nil {
				return err
			}
		}
		dbUserHistoryReward := &db.UserHistoryReward{
			UserID:         req.UserID,
			UsdTradeVolume: req.TradeVolumeUsd,
			TaskID:         req.TaskID,
			FinishTime:     utils.UnixSecondToTime(utils.ToInt64(req.FinishTime)),
			CreatedAt:      time.Now(),
			TaskType:       req.TradeType,
			TradeNo:        req.TradeNo,
			RewardScore:    addPendingScore.String(),
			UsdFeeVolume:   req.FeeUsdCost,
		}
		//上级抽成的交易量
		if req.UserID != parentUserId && parentUserId != "" {
			var dbUserHistoryTotal db.UserHistoryTotal
			err := tx.Table("user_history_total").Where("user_id=?", parentUserId).First(&dbUserHistoryTotal).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			isExist := true
			if errors.Is(err, gorm.ErrRecordNotFound) {
				isExist = false
			}
			//抽成奖励
			dbUserHistoryTotal.UserID = parentUserId
			newDecimal, _ := decimal.NewFromString(dbUserHistoryTotal.RakebackPending)
			newAddPending := tradeVolumeUsdDecimal.Shift(-6).Mul(decimal.NewFromFloat(0.0004))
			///Mul(
			//	decimal.NewFromFloat(
			//rocksCache.GetCurrentTradeParentRewardFee(dbUserHistoryTotal.SubTotalTradeVolume)))
			newDecimal = newDecimal.Add(newAddPending)
			dbUserHistoryTotal.RakebackPending = newDecimal.String()

			oldSubTotal, _ := decimal.NewFromString(dbUserHistoryTotal.SubTotalTradeVolume)
			oldSubTotal = tradeVolumeUsdDecimal.Add(oldSubTotal)
			dbUserHistoryTotal.SubTotalTradeVolume = oldSubTotal.String()
			if !isExist {
				err = tx.Table("user_history_total").Where("user_id=?", parentUserId).Create(&dbUserHistoryTotal).Error
				if err != nil {
					return err
				}
			} else {
				err = tx.Table("user_history_total").Where("user_id=?", parentUserId).Updates(map[string]interface{}{
					"rakeback_pending":       dbUserHistoryTotal.RakebackPending,
					"sub_total_trade_volume": dbUserHistoryTotal.SubTotalTradeVolume,
				}).Error
				if err != nil {
					return err
				}
			}
		}
		//if apiKeyUserID != "" {
		//	var dbUserHistoryTotal db.UserHistoryTotal
		//	err := tx.Table("user_history_total").Where("user_id=?", parentUserId).First(&dbUserHistoryTotal).Error
		//	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		//		return err
		//	}
		//	isExist := true
		//	if errors.Is(err, gorm.ErrRecordNotFound) {
		//		isExist = false
		//	}
		//	dbUserHistoryTotal.UserID = parentUserId
		//	oldTotalVolume, _ := decimal.NewFromString(dbUserHistoryTotal.TotalTradeVolume)
		//	addDecimal, _ := decimal.NewFromString(req.TradeVolumeUsd)
		//	oldTotalVolume = oldTotalVolume.Add(addDecimal)
		//	dbUserHistoryTotal.TotalTradeVolume = oldTotalVolume.String()
		//	if !isExist {
		//		err = tx.Table("user_history_total").Where("user_id=?", parentUserId).Create(dbUserHistoryTotal).Error
		//		if err != nil {
		//			return err
		//		}
		//	} else {
		//		err = tx.Table("user_history_total").Where("user_id=?", parentUserId).Updates(map[string]interface{}{
		//			"total_trade_volume": oldTotalVolume.String(),
		//		}).Error
		//		if err != nil {
		//			return err
		//		}
		//	}
		//}
		err = tx.Table("user_history_reward").Create(dbUserHistoryReward).Error
		return err

	})
	if err == nil {
		return &pbSwapRobot.BotSwapTradeResp{
			CommonResp: &pbSwapRobot.CommonResp{
				ErrCode: 0,
				ErrMsg:  "",
			},
		}, nil
	}
	return &pbSwapRobot.BotSwapTradeResp{
		CommonResp: &pbSwapRobot.CommonResp{
			ErrCode: 400,
			ErrMsg:  err.Error(),
		},
	}, nil
}
