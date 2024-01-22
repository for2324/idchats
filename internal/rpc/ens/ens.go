package ens

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbEns "Open_IM/pkg/proto/ens"
	pbOrder "Open_IM/pkg/proto/order"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/protobuf/proto"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
)

type EnsServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewEnsServer(port int) *EnsServer {
	log.NewPrivateLog(constant.LogFileName)
	return &EnsServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImEns,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *EnsServer) Run() {
	log.NewInfo("0", "EnsServer run...")

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

	// init kafka consumer
	ConsumerGroup := kfk.NewMConsumerGroup(
		&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_8_1_0, OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false},
		[]string{config.Config.Kafka.MsgOrder.Topic},
		config.Config.Kafka.MsgOrder.Addr, config.Config.Kafka.ConsumerGroupID.MsgToOrderEns,
	)
	go ConsumerGroup.RegisterHandleAndConsumer(&EnsOrderConsumerGroupHandle{})

	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	//User friend related services register to etcd
	pbEns.RegisterEnsServiceServer(srv, s)
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

	go RestartRegisterMind()

	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error(), listener)
		return
	}
}

func (s *EnsServer) CreateRegisterEnsOrder(ctx context.Context, req *pbEns.CreateRegisterEnsOrderReq) (*pbEns.CreateRegisterEnsOrderResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args", req)
	// get USDPrice from ens contract
	chainId := config.Config.Ens.ChainId
	USDPrice, err := RentUSDPrice(req.OperationID, req.EnsName, chainId)
	if err != nil {
		log.NewError(req.OperationID, "RentUSDPrice failed ", err.Error())
		return &pbEns.CreateRegisterEnsOrderResp{
			CommonResp: &pbEns.CommonResp{
				ErrCode: constant.ErrChainUp.ErrCode,
				ErrMsg:  err.Error(),
			},
		}, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "RentUSDPrice:", USDPrice)

	//  预估调用合约 GasFee
	constractGasFee, err := SuggestContractGasFee(chainId, req.EnsName, req.EnsInviter, config.Config.Ens.Contract)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SuggestContractUsdGasFee failed", err.Error())
		return &pbEns.CreateRegisterEnsOrderResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrChainUp.ErrCode, ErrMsg: err.Error()}}, nil
	}
	// 1Coin = ? USD
	coinUsdPrice, err := utils.GetCoinUSDPrice(chainId)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetCoinUSDPrice failed", err.Error())
		return &pbEns.CreateRegisterEnsOrderResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrChainUp.ErrCode, ErrMsg: err.Error()}}, nil
	}
	// 预估调用合约花费的USD
	UsdGasPrice := float64(constractGasFee) / math.Pow(10, 18) * float64(coinUsdPrice)
	orderInfo := &db.EnsRegisterOrder{
		EnsName:    req.EnsName,
		TxnType:    req.TxnType,
		USDPrice:   USDPrice,
		USDGasFee:  uint64(UsdGasPrice),
		UserId:     req.UserId,
		ChainId:    chainId,
		EnsInviter: req.EnsInviter,
	}
	err = imdb.CreateEnsRegisterOrder(orderInfo)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateEnsRegisterOrder db failed", err.Error())
		return &pbEns.CreateRegisterEnsOrderResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	// 创建监听任务
	etcdConn := getcdv3.GetDefaultConn(
		config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImOrder,
		req.OperationID,
	)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		return &pbEns.CreateRegisterEnsOrderResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	client := pbOrder.NewOrderServiceClient(etcdConn)
	resp, err := client.CreatePayScanBlockTask(ctx, &pbOrder.CreatePayScanBlockTaskReq{
		USD:         USDPrice + uint64(UsdGasPrice),
		FormAddress: req.UserId,
		OperationID: req.OperationID,
		OrderId:     fmt.Sprint(orderInfo.OrderId),
		TxnType:     req.TxnType,
		Mark:        "ens",
	})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreatePayScanBlockTask failed", err.Error())
		return &pbEns.CreateRegisterEnsOrderResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	if resp.CommonResp.ErrCode != 0 {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreatePayScanBlockTask failed", resp.CommonResp.ErrMsg)
		return &pbEns.CreateRegisterEnsOrderResp{CommonResp: &pbEns.CommonResp{ErrCode: resp.CommonResp.ErrCode, ErrMsg: resp.CommonResp.ErrMsg}}, nil
	}
	// imdb.UpdateEnsRegisterOrder(resp.CommonResp)
	respOrderInfo := &pbEns.EnsOrderInfo{}
	respScanTaskInfo := &pbOrder.ScanTaskInfo{}
	utils.CopyStructFields(&respOrderInfo, orderInfo)
	respOrderInfo.CreateTime = orderInfo.CreateTime.Format("2006-01-02 15:04:05")
	respOrderInfo.ExpireTime = orderInfo.ExpireTime.Format("2006-01-02 15:04:05")
	if !orderInfo.PayTime.IsZero() {
		respOrderInfo.PayTime = orderInfo.PayTime.Format("2006-01-02 15:04:05")
	}
	utils.CopyStructFields(&respScanTaskInfo, resp.ScanTaskInfo)
	return &pbEns.CreateRegisterEnsOrderResp{
		CommonResp:   &pbEns.CommonResp{},
		EnsOrderInfo: respOrderInfo,
		ScanTaskInfo: respScanTaskInfo,
	}, nil
}

func (s *EnsServer) GetEnsOrderInfo(ctx context.Context, req *pbEns.GetEnsOrderInfoReq) (*pbEns.GetEnsOrderInfoResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args", req)
	orderInfo, err := imdb.GetEnsRegisterOrderInfo(req.OrderId)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetEnsRegisterOrderInfo db failed", err.Error())
		return &pbEns.GetEnsOrderInfoResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	respOrderInfo := &pbEns.EnsOrderInfo{}
	scanTaskInfo := &pbOrder.ScanTaskInfo{}
	utils.CopyStructFields(&respOrderInfo, orderInfo)
	respOrderInfo.CreateTime = orderInfo.CreateTime.Format("2006-01-02 15:04:05")
	respOrderInfo.ExpireTime = orderInfo.ExpireTime.Format("2006-01-02 15:04:05")
	if !orderInfo.PayTime.IsZero() {
		respOrderInfo.PayTime = orderInfo.PayTime.Format("2006-01-02 15:04:05")
	}
	task, _ := imdb.GetPayScanBlockTaskByOrderId(fmt.Sprint(orderInfo.OrderId))
	if task != nil {
		utils.CopyStructFields(&scanTaskInfo, task)
		scanTaskInfo.CreateTime = task.CreateTime.Format("2006-01-02 15:04:05")
		scanTaskInfo.BlockExpireTime = task.BlockExpireTime.Format("2006-01-02 15:04:05")
		scanTaskInfo.BlockStartTime = task.BlockStartTime.Format("2006-01-02 15:04:05")
	}
	return &pbEns.GetEnsOrderInfoResp{
		CommonResp:   &pbEns.CommonResp{},
		EnsOrderInfo: respOrderInfo,
		ScanTaskInfo: scanTaskInfo,
	}, nil
}

func (s *EnsServer) ConfirmEnsOrderHasBeenPaid(ctx context.Context, req *pbEns.ConfirmEnsOrderHasBeenPaidReq) (sResp *pbEns.ConfirmEnsOrderHasBeenPaidResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args", req)
	// 找到符合条件的订单
	orderId, err := strconv.ParseUint(req.OrderId, 10, 64)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ParseUint failed", err.Error())
		return &pbEns.ConfirmEnsOrderHasBeenPaidResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	res := imdb.UpdateEnsRegisterOrderPaid(orderId, req.TxnHash)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateEnsRegisterOrderPaid db failed", err.Error())
		return &pbEns.ConfirmEnsOrderHasBeenPaidResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	if res.RowsAffected == 0 {
		uErr := errors.New("UpdateEnsRegisterOrderPaid failed RowsAffected is 0")
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateEnsRegisterOrderPaid failed", "RowsAffected is 0")
		return &pbEns.ConfirmEnsOrderHasBeenPaidResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: uErr.Error()}}, nil
	}
	order, err := imdb.GetEnsRegisterOrderInfo(orderId)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetEnsRegisterOrderInfo db failed", err.Error())
		return &pbEns.ConfirmEnsOrderHasBeenPaidResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	var secret [32]byte
	_, _ = rand.Read(secret[:])
	ensDomain := order.EnsName
	ensDomainHash, _ := NameHash(ensDomain + ".biu")
	bytes32Type, _ := abi.NewType("bytes32", "bytes32", nil)
	addressType, _ := abi.NewType("address", "address", nil)
	arguments := abi.Arguments{
		{
			Type: bytes32Type,
		},
		{
			Type: addressType,
		},
	}
	bytes, err := arguments.Pack(ensDomainHash, common.HexToAddress(config.Config.Ens.Contract))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sig := "setAddr(bytes32,address)"
	id := crypto.Keccak256([]byte(sig))[:4]
	byteArray := append(id, bytes...)
	registerData := IETHRegistrarControllerRegisterData{
		Name:                 order.EnsName,
		Owner:                common.HexToAddress(order.UserId),
		Secret:               secret,
		Duration:             GetUint256Max(),
		Resolver:             common.HexToAddress(config.Config.Ens.Resolver),
		Data:                 [][]byte{byteArray},
		ReverseRecord:        true, // 反向解析记录
		OwnerControlledFuses: 1,    // 只有owner 控制
		RebateName:           order.EnsInviter,
	}
	isErc20 := isERC20(order.TxnType)
	tx, err := RegisterEnsChainUp(req.OperationID, order.ChainId, registerData, isErc20)
	if err != nil {
		if uerr := imdb.UpdateEnsOrderRegisterFailed(order.OrderId, err.Error()); uerr != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateEnsOrderRegisterFailed db failed", uerr.Error())
		}
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "RegisterEnsChainUp failed", err.Error())
		return &pbEns.ConfirmEnsOrderHasBeenPaidResp{CommonResp: &pbEns.CommonResp{ErrCode: constant.ErrChainUp.ErrCode, ErrMsg: err.Error()}}, nil
	}
	if uerr := imdb.UpdateEnsOrderRegisterConfirmed(order.OrderId, tx.Hash().String()); uerr != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateEnsOrderRegisterConfirmed db failed", uerr.Error())
	} else {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "UpdateEnsOrderRegisterConfirmed db success")
	}
	go WaitRegisterMined(req.OperationID, order, tx)
	return &pbEns.ConfirmEnsOrderHasBeenPaidResp{CommonResp: &pbEns.CommonResp{}}, nil
}

func GetEthClient(OperationID string, chainId int64) (*ethclient.Client, error) {
	chainKey := strconv.FormatInt(chainId, 10)
	rpcInfo, ok := config.Config.ChainIdRpcMap[chainKey]
	if !ok || len(rpcInfo) == 0 {
		log.NewError(OperationID, utils.GetSelfFuncName(), "GetEnsInstant failed", "chainId not found")
		return nil, errors.New("chainId not found")
	}

	endpoint := rpcInfo[0]
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.NewError(OperationID, "Failed to connect to Ethereum network: %v", err)
		return nil, err
	}
	return client, nil
}

func GetEnsInstant(OperationID string, chainId int64) (*Ens, error) {
	client, err := GetEthClient(OperationID, chainId)
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "GetEnsInstant failed", err.Error())
		return nil, err
	}
	contractAddress := config.Config.Ens.Contract
	ContractAddress := common.HexToAddress(contractAddress)
	return NewEns(ContractAddress, client)
}
func RegisterEnsChainUp(OperationID string, chainId int64, data IETHRegistrarControllerRegisterData, isERC20 bool) (
	*types.Transaction, error,
) {
	caller, err := GetEnsInstant(OperationID, chainId)
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "GetEnsInstant failed", err.Error())
		return nil, err
	}
	senderPrivateKey, err := crypto.HexToECDSA(config.Config.Ens.EnsOwnerPrivateKeyHex)
	if err != nil {
		log.NewError(OperationID, "Failed to parse private key: %v", err)
		return nil, err
	}
	senderAuth, err := bind.NewKeyedTransactorWithChainID(senderPrivateKey, big.NewInt(chainId))
	if err != nil {
		log.NewError(OperationID, "Failed to create authorized transactor: %v", err)
		return nil, err
	}
	return caller.Register(senderAuth, data, isERC20)
}

func RentUSDPrice(OperationID, ensName string, chainId int64) (uint64, error) {
	caller, err := GetEnsInstant(OperationID, chainId)
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "GetEnsInstant failed", err.Error())
		return 0, nil
	}
	price, err := caller.RentUSDPrice(nil, ensName, GetUint256Max())
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "RentPrice failed", err.Error())
		return 0, err
	}
	return price.Base.Uint64(), nil
}

func WaitRegisterMined(OperationID string, order *db.EnsRegisterOrder, tx *types.Transaction) {
	// Wait for the transaction to be mined
	client, err := GetEthClient(OperationID, order.ChainId)
	if err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "GetEthClient failed", err.Error())
		return
	}
	if tx == nil {
		tx, _, err = client.TransactionByHash(context.Background(), common.HexToHash(order.RegisterTxnHash))
		if err != nil {
			log.NewError(OperationID, utils.GetSelfFuncName(), "TransactionByHash failed", err.Error())
			return
		}
	}
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		imdb.UpdateEnsOrderRegisterFailed(order.OrderId, err.Error())
		log.NewError(OperationID, utils.GetSelfFuncName(), "WaitMined failed", err.Error())
		return
	}
	if receipt.Status == types.ReceiptStatusFailed {
		imdb.UpdateEnsOrderRegisterRefund(order.OrderId, "WaitMined failed")
		log.NewError(OperationID, utils.GetSelfFuncName(), "WaitMined failed", "receipt.Status == types.ReceiptStatusFailed")
		return
	}
	if err := imdb.UpdateEnsOrderRegisterSuccess(order.OrderId); err != nil {
		log.NewError(OperationID, utils.GetSelfFuncName(), "UpdateEnsOrderRegisterSuccess db failed", err.Error())
		return
	}
	log.NewInfo(OperationID, utils.GetSelfFuncName(), "UpdateEnsOrderRegisterSuccess db success")
}

func GetUint256Max() *big.Int {
	max := new(big.Int).SetUint64(1)
	max.Lsh(max, 256)
	max.Sub(max, big.NewInt(1))
	return max
}

func RestartRegisterMind() {
	orders, err := imdb.GetEnsOrderRegisterNotMined()
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "GetEnsOrderRegisterNotMined failed", err.Error())
		return
	}
	for _, order := range orders {
		WaitRegisterMined("init RestartRegisterMind", &order, nil)
	}
}

func isERC20(txnType string) bool {
	// 判断 txnType 是否已 ERC20 开头
	return strings.HasPrefix(txnType, "ERC20")
}

func SuggestContractGasFee(chainId int64, ensName, ensInviter string, constact string) (uint64, error) {
	client, err := utils.GetEthClient(chainId)
	if err != nil {
		return 0, err
	}
	// gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return 0, err
	}
	constractAddr := common.HexToAddress(constact)
	contract, _ := abi.JSON(strings.NewReader(EnsMetaData.ABI))
	var secret [32]byte
	_, _ = rand.Read(secret[:])
	ensDomain := ensName
	ensDomainHash, _ := NameHash(ensDomain + ".biu")
	bytes32Type, _ := abi.NewType("bytes32", "bytes32", nil)
	addressType, _ := abi.NewType("address", "address", nil)
	arguments := abi.Arguments{
		{
			Type: bytes32Type,
		},
		{
			Type: addressType,
		},
	}
	bytes, err := arguments.Pack(ensDomainHash, common.HexToAddress(config.Config.Ens.Contract))
	if err != nil {
		return 0, err
	}
	sig := "setAddr(bytes32,address)"
	id := crypto.Keccak256([]byte(sig))[:4]
	byteArray := append(id, bytes...)
	registerData := IETHRegistrarControllerRegisterData{
		Name:                 ensName,
		Owner:                common.HexToAddress(config.Config.Ens.EnsOwnerAddress),
		Secret:               secret,
		Duration:             GetUint256Max(),
		Resolver:             common.HexToAddress(config.Config.Ens.Resolver),
		Data:                 [][]byte{byteArray},
		ReverseRecord:        true, // 反向解析记录
		OwnerControlledFuses: 1,    // 只有owner 控制
		RebateName:           ensInviter,
	}
	inputData, err := contract.Pack("register", registerData, true)
	if err != nil {
		return 0, err
	}
	// 一笔转账预估的 gas used
	gasUsed, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     common.HexToAddress(config.Config.Ens.EnsOwnerAddress),
		To:       &constractAddr,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     inputData,
	})
	if err != nil {
		return 0, err
	}
	return gasPrice.Uint64() * gasUsed, nil
}

type EnsOrderConsumerGroupHandle struct{}

func (*EnsOrderConsumerGroupHandle) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*EnsOrderConsumerGroupHandle) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (och *EnsOrderConsumerGroupHandle) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	log.NewDebug("", "EnsOrderConsumerGroupHandle new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			och.EnsOrderHandler(msg, string(msg.Key))
		} else {
			log.Error("", "mongo msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
func (och *EnsOrderConsumerGroupHandle) EnsOrderHandler(msg *sarama.ConsumerMessage, key string) error {
	var order pbOrder.MsgDataToOrderByMQ
	err := proto.Unmarshal(msg.Value, &order)
	if err != nil {
		log.Error("", "EnsOrderConsumerGroupHandle EnsOrderHandler json.Unmarshal failed", err.Error())
		return err
	}
	if order.Mark != "ens" {
		return nil
	}
	log.NewInfo(order.OperationID, "EnsOrderConsumerGroupHandle EnsOrderHandler", "order", order.OrderId)
	// notify
	s := NewEnsServer(0)
	// notify order has been paid
	resp, err := s.ConfirmEnsOrderHasBeenPaid(context.Background(), &pbEns.ConfirmEnsOrderHasBeenPaidReq{
		OperationID: "ConfirmEnsOrderHasBeenPaid",
		OrderId:     order.OrderId,
		TxnHash:     order.TxnHash,
	})
	if err != nil {
		log.Error(order.OperationID, "confirm ens order has been paid failed, err: ", err.Error())
		return err
	}
	if resp.CommonResp.ErrCode != 0 {
		log.Error(order.OperationID, "confirm ens order has been paid failed, err: ", resp.CommonResp.ErrMsg)
		return errors.New(resp.CommonResp.ErrMsg)
	}
	return nil
}
