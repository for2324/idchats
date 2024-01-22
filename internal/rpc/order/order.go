package order

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/kafka"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"math/big"
	"time"

	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbOrder "Open_IM/pkg/proto/order"

	"Open_IM/pkg/utils"
	"context"
	"net"
	"strconv"
	"strings"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
)

var (
	producerToOrder *kafka.Producer
)

type OrderServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewOrderServer(port int) *OrderServer {
	log.NewPrivateLog(constant.LogFileName)
	return &OrderServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImOrder,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *OrderServer) Run() {
	log.NewInfo("0", "OrderServer run...")

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

	producerToOrder = kafka.NewKafkaProducer(config.Config.Kafka.MsgOrder.Addr, config.Config.Kafka.MsgOrder.Topic)

	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	//User friend related services register to etcd
	pbOrder.RegisterOrderServiceServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema,
		strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register friend module  rpc to etcd err"))
	}

	// run init task
	go LoopScanBlock()

	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error(), listener)
		return
	}
}

func (s *OrderServer) CreatePayScanBlockTask(ctx context.Context, req *pbOrder.CreatePayScanBlockTaskReq) (*pbOrder.CreatePayScanBlockTaskResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args", req)

	confInfo, ok := config.Config.Pay.TnxTypeConfMap[req.TxnType]
	if !ok {
		log.NewError(req.OperationID, "config.Config.Pay.TnxTypeConfMap[req.TxnType] is empty ", req.TxnType)
		return &pbOrder.CreatePayScanBlockTaskResp{
			CommonResp: &pbOrder.CommonResp{
				ErrCode: constant.ErrChainUp.ErrCode,
				ErrMsg:  "config.Config.Pay.TnxTypeConfMap[req.TxnType] is empty",
			},
		}, nil
	}
	// check confInfo
	if confInfo.Decimal == 0 || confInfo.Retention == 0 || confInfo.Accuracy == 0 || confInfo.Retention < confInfo.Accuracy {
		log.NewError(req.OperationID, "config.Config.Pay.TnxTypeConfMap[req.TxnType].Decimal or config.Config.Pay.TnxTypeConfMap[req.TxnType].ConcurrencyPrecision is empty ", req.TxnType)
		return &pbOrder.CreatePayScanBlockTaskResp{
			CommonResp: &pbOrder.CommonResp{
				ErrCode: constant.ErrChainUp.ErrCode,
				ErrMsg:  "config.Config.Pay.TnxTypeConfMap[req.TxnType].Decimal or config.Config.Pay.TnxTypeConfMap[req.TxnType].ConcurrencyPrecision is empty",
			},
		}, nil
	}
	if len(confInfo.ReceivedAddress) == 0 {
		log.NewError(req.OperationID, "config.Config.Pay.TnxTypeConfMap[req.TxnType].ReceivedAddress is empty ", req.TxnType)
		return &pbOrder.CreatePayScanBlockTaskResp{
			CommonResp: &pbOrder.CommonResp{
				ErrCode: constant.ErrChainUp.ErrCode,
				ErrMsg:  "config.Config.Pay.TnxTypeConfMap[req.TxnType].ReceivedAddress is empty",
			},
		}, nil
	}

	// task check
	taskType := req.TxnType
	tag := ""
	if isERC20(req.TxnType) {
		if confInfo.ContractAddress == "" {
			log.NewError(req.OperationID, "ContractAddress is empty ", req.TxnType)
			return &pbOrder.CreatePayScanBlockTaskResp{
				CommonResp: &pbOrder.CommonResp{
					ErrCode: constant.ErrChainUp.ErrCode,
					ErrMsg:  "ContractAddress is empty",
				},
			}, nil
		}
		tag = confInfo.ContractAddress
	}

	var valueNonce uint64 = 0
	var ReceivedAddress string
	var giveBackId string
	for _, v := range confInfo.ReceivedAddress {
		nonce, backId, err := NewChainTable(fmt.Sprintf("%s/%s", req.TxnType, v), confInfo.Accuracy).Next()
		if err != nil {
			continue
		}
		// 翻转 nonce 123 => 321
		valueNonce = uint64(nonce)
		ReceivedAddress = v
		giveBackId = backId
		break
	}
	if ReceivedAddress == "" {
		log.NewError(req.OperationID, "NewChainTable failed ", req.TxnType)
		return &pbOrder.CreatePayScanBlockTaskResp{
			CommonResp: &pbOrder.CommonResp{
				ErrCode: constant.ErrChainUp.ErrCode,
				ErrMsg:  "NewChainTable failed",
			},
		}, nil
	}
	// 1coin = ? usd
	var rate float64
	var err error
	resolver := ScanFactory(taskType, tag, confInfo.ChainId)
	if resolver == nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ScanFactory failed", req.TxnType)
		return &pbOrder.CreatePayScanBlockTaskResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrChainUp.ErrCode, ErrMsg: "ScanFactory failed"}}, nil
	}
	rate, err = resolver.GetCoinUSDPrice(req.OperationID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetCoinUSDPrice failed", err.Error())
		return &pbOrder.CreatePayScanBlockTaskResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrChainUp.ErrCode, ErrMsg: err.Error()}}, nil
	}
	log.Debug(req.OperationID, utils.GetSelfFuncName(), "rate", rate, "tnxType", req.TxnType)
	Decimal := confInfo.Decimal
	// 1usd = ? coin
	coinVal := math.Pow10(int(Decimal-6)) / rate
	coinVal = coinVal * float64(req.USD)
	round := (uint64(coinVal/math.Pow10(int(Decimal)-confInfo.Retention+confInfo.Accuracy))+1)*uint64(math.Pow10(confInfo.Accuracy)) + valueNonce
	power := new(big.Int)
	power.Exp(big.NewInt(10), big.NewInt(int64(int(Decimal)-confInfo.Retention)), nil)
	value := new(big.Int)
	value = value.SetBit(value, 256, 1)
	value = value.Mul(big.NewInt(int64(round)), power)
	// 获取当前最新区块高度和时间
	block, err := resolver.GetBlock(req.OperationID, ctx, confInfo.ChainId, nil)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "getBlock failed", err.Error())
		return &pbOrder.CreatePayScanBlockTaskResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrChainUp.ErrCode, ErrMsg: "get last block faild"}}, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "GetBlock block", block.BlockNumber, "blockTime", block.Time)
	blockStartNumber := block.BlockNumber
	blockStartTime := block.Time
	dbBlockTask := &db.PayScanBlockTask{
		OrderId:           req.OrderId,
		FormAddress:       strings.ToLower(req.FormAddress),
		ToAddress:         strings.ToLower(ReceivedAddress),
		ChainId:           confInfo.ChainId,
		Value:             value.String(),
		USDPrice:          req.USD,
		Decimal:           Decimal,
		Type:              taskType,
		Rate:              uint64(rate * math.Pow10(6)),
		Tag:               tag,
		TxnType:           req.TxnType,
		StartBlockNumber:  blockStartNumber,
		ScanBlockNumber:   blockStartNumber,
		BlockStartTime:    blockStartTime,
		BlockExpireTime:   blockStartTime.Add(time.Duration(config.Config.Pay.OrderExpireTime) * time.Minute),
		Mark:              req.Mark,
		Ex:                giveBackId,
		NotifyUrl:         req.NotifyUrl,
		Attach:            req.Attach,
		NotifyEncryptType: req.NotifyEncryptType,
		NotifyEncryptKey:  req.NotifyEncryptKey,
	}
	err = imdb.CreatePayScanBlockTask(dbBlockTask)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreatePayScanBlockTask db failed", err.Error())
		return &pbOrder.CreatePayScanBlockTaskResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	go LoopScanBlock()
	taskInfo := &pbOrder.ScanTaskInfo{}
	utils.CopyStructFields(&taskInfo, dbBlockTask)
	taskInfo.CreateTime = dbBlockTask.CreateTime.Format("2006-01-02 15:04:05")
	taskInfo.BlockStartTime = dbBlockTask.BlockStartTime.Format("2006-01-02 15:04:05")
	taskInfo.BlockExpireTime = dbBlockTask.BlockExpireTime.Format("2006-01-02 15:04:05")
	return &pbOrder.CreatePayScanBlockTaskResp{CommonResp: &pbOrder.CommonResp{}, ScanTaskInfo: taskInfo}, nil
}

func (s *OrderServer) GetPayScanBlockTaskByOrderId(ctx context.Context, req *pbOrder.GetPayScanBlockTaskByOrderIdReq) (*pbOrder.GetPayScanBlockTaskByOrderIdResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args", req)
	orderInfo, err := imdb.GetPayScanBlockTaskByOrderId(req.OrderId)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetPayScanBlockTaskByOrderId failed", err.Error())
		return &pbOrder.GetPayScanBlockTaskByOrderIdResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	taskInfo := &pbOrder.ScanTaskInfo{}
	utils.CopyStructFields(&taskInfo, orderInfo)
	taskInfo.CreateTime = orderInfo.CreateTime.Format("2006-01-02 15:04:05")
	taskInfo.BlockStartTime = orderInfo.BlockStartTime.Format("2006-01-02 15:04:05")
	taskInfo.BlockExpireTime = orderInfo.BlockExpireTime.Format("2006-01-02 15:04:05")
	return &pbOrder.GetPayScanBlockTaskByOrderIdResp{CommonResp: &pbOrder.CommonResp{}, ScanTaskInfo: taskInfo}, nil
}

func (s *OrderServer) GetPayScanBlockTaskById(ctx context.Context, req *pbOrder.GetPayScanBlockTaskByIdReq) (*pbOrder.GetPayScanBlockTaskByIdResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args", req)
	orderInfo, err := imdb.GetPayScanBlockTaskById(req.Id)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetPayScanBlockTaskById failed", err.Error())
		return &pbOrder.GetPayScanBlockTaskByIdResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	taskInfo := &pbOrder.ScanTaskInfo{}
	utils.CopyStructFields(&taskInfo, orderInfo)
	taskInfo.CreateTime = orderInfo.CreateTime.Format("2006-01-02 15:04:05")
	taskInfo.BlockStartTime = orderInfo.BlockStartTime.Format("2006-01-02 15:04:05")
	taskInfo.BlockExpireTime = orderInfo.BlockExpireTime.Format("2006-01-02 15:04:05")
	return &pbOrder.GetPayScanBlockTaskByIdResp{CommonResp: &pbOrder.CommonResp{}, ScanTaskInfo: taskInfo}, nil
}

// 补单
func (s *OrderServer) ReplenishmentOrder(ctx context.Context, req *pbOrder.ReplenishmentOrderReq) (*pbOrder.ReplenishmentOrderResp, error) {
	orderInfo, err := imdb.GetPayScanBlockTaskById(req.Id)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetPayScanBlockTaskById failed", err.Error())
		return &pbOrder.ReplenishmentOrderResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	blockPayTime, err := GetTxBlockTime(orderInfo.ChainId, req.TxnHash)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTxBlockTime failed", err.Error())
		return &pbOrder.ReplenishmentOrderResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	err = imdb.ReplenishmentOrder(req.Id, req.TxnHash, req.Remark, *blockPayTime)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ReplenishmentOrder failed", err.Error())
		return &pbOrder.ReplenishmentOrderResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	confInfo, ok := config.Config.Pay.TnxTypeConfMap[orderInfo.TxnType]
	if ok {
		NewChainTable(fmt.Sprintf("%s/%s", orderInfo.TxnType, orderInfo.ToAddress), confInfo.Accuracy).GiveBack(orderInfo.Ex)
	}
	key := fmt.Sprintf("%s:%s", orderInfo.Mark, orderInfo.FormAddress)
	// 更优的方式应该是采用 dtm 的分布式事务（补偿机制）
	pid, offset, err := producerToOrder.SendMessage(
		&pbOrder.MsgDataToOrderByMQ{
			OperationID:       req.OperationID,
			ID:                fmt.Sprint(orderInfo.Id),
			USD:               orderInfo.USDPrice,
			OrderId:           orderInfo.OrderId,
			UserId:            orderInfo.FormAddress,
			TxnType:           orderInfo.TxnType,
			Mark:              orderInfo.Mark,
			Value:             orderInfo.Value,
			Decimal:           orderInfo.Decimal,
			TxnHash:           orderInfo.TxnHash,
			NotifyUrl:         orderInfo.NotifyUrl,
			CreateTime:        orderInfo.BlockStartTime.Format("2006-01-02 15:04:05"),
			PayTime:           orderInfo.BlockPayTime.Format("2006-01-02 15:04:05"),
			Attach:            orderInfo.Attach,
			NotifyEncryptKey:  orderInfo.NotifyEncryptKey,
			NotifyEncryptType: orderInfo.NotifyEncryptType,
		},
		key, req.OperationID)
	if err != nil {
		log.Error(req.OperationID, utils.GetSelfFuncName(), "kafka send failed", "send data", orderInfo.Id, "pid", pid, "offset", offset, "err", err.Error(), "key", key)
		return &pbOrder.ReplenishmentOrderResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrKafka.ErrCode, ErrMsg: err.Error()}}, nil
	}
	if err := imdb.UpdatePayScanBlockTaskStatusComfirm(orderInfo.Id); err != nil {
		log.Error(req.OperationID, "update pay scan block task status failed, err: ", err.Error())
		return &pbOrder.ReplenishmentOrderResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	return &pbOrder.ReplenishmentOrderResp{}, nil
}

// 强制补单（无 tnxHash）
func (s *OrderServer) ForcedReplenishmentOrder(ctx context.Context, req *pbOrder.ForcedReplenishmentOrderReq) (*pbOrder.ForcedReplenishmentOrderResp, error) {
	orderInfo, err := imdb.GetPayScanBlockTaskById(req.Id)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetPayScanBlockTaskById failed", err.Error())
		return &pbOrder.ForcedReplenishmentOrderResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	// 以补单时间作为支付时间
	blockPayTime := time.Now()
	err = imdb.ReplenishmentOrder(req.Id, "", req.Remark, blockPayTime)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ReplenishmentOrder failed", err.Error())
		return &pbOrder.ForcedReplenishmentOrderResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	confInfo, ok := config.Config.Pay.TnxTypeConfMap[orderInfo.TxnType]
	if ok {
		NewChainTable(fmt.Sprintf("%s/%s", orderInfo.TxnType, orderInfo.ToAddress), confInfo.Accuracy).GiveBack(orderInfo.Ex)
	}
	key := fmt.Sprintf("%s:%s", orderInfo.Mark, orderInfo.FormAddress)
	// 更优的方式应该是采用 dtm 的分布式事务（补偿机制）
	pid, offset, err := producerToOrder.SendMessage(
		&pbOrder.MsgDataToOrderByMQ{
			OperationID:       req.OperationID,
			ID:                fmt.Sprint(orderInfo.Id),
			USD:               orderInfo.USDPrice,
			OrderId:           orderInfo.OrderId,
			UserId:            orderInfo.FormAddress,
			TxnType:           orderInfo.TxnType,
			Mark:              orderInfo.Mark,
			Value:             orderInfo.Value,
			Decimal:           orderInfo.Decimal,
			TxnHash:           orderInfo.TxnHash,
			NotifyUrl:         orderInfo.NotifyUrl,
			CreateTime:        orderInfo.BlockStartTime.Format("2006-01-02 15:04:05"),
			PayTime:           orderInfo.BlockPayTime.Format("2006-01-02 15:04:05"),
			Attach:            orderInfo.Attach,
			NotifyEncryptKey:  orderInfo.NotifyEncryptKey,
			NotifyEncryptType: orderInfo.NotifyEncryptType,
		},
		key, req.OperationID)
	if err != nil {
		log.Error(req.OperationID, utils.GetSelfFuncName(), "kafka send failed", "send data", orderInfo.Id, "pid", pid, "offset", offset, "err", err.Error(), "key", key)
		return &pbOrder.ForcedReplenishmentOrderResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrKafka.ErrCode, ErrMsg: err.Error()}}, nil
	}
	if err := imdb.UpdatePayScanBlockTaskStatusComfirm(orderInfo.Id); err != nil {
		log.Error(req.OperationID, "update pay scan block task status failed, err: ", err.Error())
		return &pbOrder.ForcedReplenishmentOrderResp{CommonResp: &pbOrder.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	return &pbOrder.ForcedReplenishmentOrderResp{}, nil
}

// func SuggestERC20GasFee(chainId int64, toContractAddress string) (uint64, error) {
// 	client, err := utils.GetEthClient(chainId)
// 	if err != nil {
// 		return 0, err
// 	}
// 	// gas 价格
// 	gasPrice, err := client.SuggestGasPrice(context.Background())
// 	if err != nil {
// 		return 0, err
// 	}
// 	token := big.NewInt(0)
// 	// zeroAddr := common.HexToAddress("0x106E405a99C67258E915B12A50cb5fC33c217214")
// 	toContract := common.HexToAddress(toContractAddress)
// 	contract, _ := abi.JSON(strings.NewReader(erc20.ChainABI))
// 	inputData, err := contract.Pack("transfer", toContract, token)
// 	if err != nil {
// 		return 0, err
// 	}
// 	// 一笔转账预估的 gas used
// 	gasUsed, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
// 		From: toContract,
// 		To:   &toContract,
// 		// GasPrice: gasPrice,
// 		Value: big.NewInt(0),
// 		Data:  inputData,
// 	})
// 	if err != nil {
// 		return 0, err
// 	}
// 	return gasPrice.Uint64() * gasUsed, nil
// }

func AESGCMEncrypt(plaintext string, key string) (ciphertext string, nonceStr string, err error) {
	nonce := make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		return
	}
	ciphertextByte := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	ciphertext = string(ciphertextByte)
	nonceStr = string(nonce)
	return
}
func AESGCMDecrypt(ciphertext string, nonceStr string, key string) (plaintext string, err error) {
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		return
	}
	var openData []byte
	openData, err = aesgcm.Open(nil, []byte(nonceStr), []byte(ciphertext), nil)
	if err != nil {
		return
	}
	plaintext = string(openData)
	return
}
