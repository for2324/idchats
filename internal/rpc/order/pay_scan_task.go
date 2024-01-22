package order

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbOrder "Open_IM/pkg/proto/order"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const SCAN_BLOCK_STEP uint64 = 100

var (
	ErrScanBlockNotNeed = errors.New("scan block not need")
	// 任务是否正在执行
	scanFlag bool = false
	mutex    sync.Mutex
)

// 设置变量的值
func SetFlag(value bool) {
	mutex.Lock()
	scanFlag = value
	mutex.Unlock()
}

// 获取变量的值
func GetFlag() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return scanFlag
}
func LoopScanBlock() {
	if GetFlag() {
		log.Info("scan block task is running")
		return
	}
	SetFlag(true)
	for {
		time.Sleep(time.Duration(config.Config.Pay.ScanInterval) * time.Second)
		if err := ScanBlock(); err != nil {
			if errors.Is(err, ErrScanBlockNotNeed) {
				log.Info("stop scan task scan block not need")
			} else {
				log.Error("scan block failed, err: ", err.Error())
			}
			break
		}
	}
	SetFlag(false)
}

func ScanBlock() error {
	// 根据 type + tag 分组去获取扫块任务
	groupList, err := imdb.GetPayScanBlockGroupTaskList()
	if err != nil {
		log.Error("get pay scan block group task list failed, err: ", err.Error())
		return err
	}
	if len(groupList) == 0 {
		return ErrScanBlockNotNeed
	}
	log.NewInfo("ScanBlock", groupList)
	// 获取分组信息（任务列表），执行扫块，根据信息去处理任务
	for _, group := range groupList {
		OperationID := utils.OperationIDGenerator()
		client, err := utils.GetEthClient(group.ChainId)
		if err != nil {
			log.Error(OperationID, "get eth client failed, err: ", err.Error())
			continue
		}
		defer client.Close()
		end, err := ScanGroup(OperationID, group.Type, group.Tag, group.ChainId, group.StartHeight)
		if err != nil {
			log.Error(OperationID, "scan group failed, err: ", err.Error())
			continue
		}
		header, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(end)))
		if err != nil {
			log.Error(OperationID, "get header by number failed, err: ", err.Error())
			continue
		}
		// 更新进度
		log.Debug(OperationID, "update pay scan block group task progress, start: ", group.StartHeight, ", end: ", end, ", type: ", group.Type, ", progress: ", header.Number.Uint64())
		imdb.UpdatePayScanBlockGroupTaskProgress(header.Number.Uint64(), group.Type, group.Tag, group.StartHeight, end)
		endBlockTime := time.Unix(int64(header.Time), 0)
		err = imdb.MarkPayScanBlockTaskExpired(endBlockTime)
		if err != nil {
			log.Error(OperationID, "mark pay scan block task expired failed, err: ", err.Error())
			continue
		}
		log.NewInfo(fmt.Sprintf("scan block success, type: %s, tag: %s, start: %d, end: %d", group.Type, group.Tag, group.StartHeight, end))
	}
	return nil
}

// var cacheBanlanceStartheight = make(map[string]*big.Int)
// var cacheNonceAtStartheight = make(map[string]uint64)

func ScanGroup(OperationID, groupType, groupTag string, chainId int64, startHeight uint64) (uint64, error) {
	log.NewInfo(OperationID, utils.GetSelfFuncName(), "start scan group, type: ", groupType, ", tag: ", groupTag, ", start: ", startHeight)
	scanResolver := ScanFactory(groupType, groupTag, chainId)
	if scanResolver == nil {
		log.Error(OperationID, "scan type not found, type: ", groupType)
		return 0, errors.New("scan type not found")
	}
	// 根据 扫块最小进度 获取扫块时间
	// 获取当前最新区块高度和时间
	block, err := scanResolver.GetBlock(OperationID, context.TODO(), chainId, nil)
	if err != nil {
		log.NewError("Failed to get node block number: %v", err)
		return 0, err
	}
	end := block.BlockNumber
	start := startHeight
	if end == start {
		return end, nil
	}
	if end-start > config.Config.Pay.ScanStep {
		end = start + config.Config.Pay.ScanStep
	}
	taskList, err := imdb.GetPayScanBlockTaskListByTaskTag(groupType, groupTag, start, end)
	if err != nil {
		log.Error(OperationID, "get pay scan block task list by task tag failed, err: ", err.Error())
		return 0, err
	}
	if len(taskList) == 0 {
		return end, nil
	}
	log.NewInfo(OperationID, utils.GetSelfFuncName(), "start scan block, start: ", start, ", end: ", end, ", task count: ", len(taskList))
	// 执行扫块任务
	formAddresses := make([]string, 0)
	toAddreses := make([]string, 0)
	for _, task := range taskList {
		formAddresses = append(formAddresses, task.FormAddress)
		toAddreses = append(toAddreses, task.ToAddress)
	}
	formAddresses = RemoveRepeatedElement(formAddresses)
	toAddreses = RemoveRepeatedElement(toAddreses)
	// 如果前后两次扫块的余额一致，则不需要再次扫块
	for _, to := range toAddreses {
		same, err := scanResolver.CompareStartAndEndBanlance(OperationID, to, start, end)
		if err != nil {
			log.Error(OperationID, "compare start and end banlance failed, err: ", err.Error())
		}
		if !same {
			continue
		}
		toAddreses = RemoveFromArr(toAddreses, to)
	}
	if len(formAddresses) == 0 || len(toAddreses) == 0 {
		log.Info(OperationID, "no need to scan block, formAddresses: ", formAddresses, ", toAddreses: ", toAddreses, ", start: ", start, ", end: ", end)
		return end, nil
	}
	log.NewInfo(OperationID, "scan block, formAddresses: ", formAddresses, ", toAddreses: ", toAddreses, ", start: ", start, ", end: ", end)
	TransferEvent, err := scanResolver.Scan(formAddresses, toAddreses, start, end)
	if err != nil {
		log.Error(OperationID, "scan block failed, err: ", err.Error())
		return 0, err
	}
	log.NewInfo(OperationID, "scan block success, TransferEvent: ", TransferEvent)
	for _, event := range TransferEvent {
		ResolveTransferEvent(OperationID, scanResolver, chainId, groupType, event)
	}
	log.NewInfo(OperationID, "scan block success, start: ", start, ", end: ", end)
	return end, nil
}

func ResolveTransferEvent(OperationID string, scanResolver ScanResolver, chainId int64, groupType string, event *TransferEvent) (resolveErr error) {
	defer func() {
		resolveInfo := "success"
		if resolveErr != nil {
			resolveInfo = resolveErr.Error()
		}
		err := imdb.CreateOrderPaidRecord(&db.OrderPaidRecord{
			FormAddress: event.From,
			ToAddress:   event.To,
			Value:       fmt.Sprint(event.Value),
			ChainId:     chainId,
			TxnHash:     event.TxHash,
			PayType:     groupType,
			Ex:          resolveInfo,
		})
		if err != nil {
			log.NewError(OperationID, utils.GetSelfFuncName(), "CreateOrderPaidRecord db failed", err.Error())
		} else {
			log.NewInfo(OperationID, utils.GetSelfFuncName(), "CreateOrderPaidRecord db success")
		}
	}()
	block, err := scanResolver.GetBlock(OperationID, context.TODO(), chainId, big.NewInt(int64(event.BlockNumber)))
	if err != nil {
		log.Error(OperationID, "get block time failed, err: ", err.Error())
		return err
	}
	// 更新扫块任务状态为已支付
	orderInfo, err := imdb.UpdatePayScanBlockTaskStatusFinished(event.From, event.To, event.Value.String(), chainId, block.Time, event.TxHash)
	if err != nil {
		log.Error(OperationID, "update pay scan block task status failed, err: ", err.Error())
		return err
	}
	if orderInfo == nil {
		log.Error(OperationID, "update pay scan block task status failed, orderId is empty")
		return errors.New("update pay scan block task status failed, orderId is empty")
	}
	confInfo, ok := config.Config.Pay.TnxTypeConfMap[orderInfo.TxnType]
	if ok {
		NewChainTable(fmt.Sprintf("%s/%s", orderInfo.TxnType, orderInfo.ToAddress), confInfo.Accuracy).GiveBack(orderInfo.Ex)
	}
	key := fmt.Sprintf("%s:%s", orderInfo.Mark, orderInfo.FormAddress)
	pid, offset, err := producerToOrder.SendMessage(
		&pbOrder.MsgDataToOrderByMQ{
			OperationID:       OperationID,
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
		key, OperationID)
	if err != nil {
		log.Error(OperationID, utils.GetSelfFuncName(), "kafka send failed", "send data", orderInfo.Id, "pid", pid, "offset", offset, "err", err.Error(), "key", key)
		return err
	}
	if err := imdb.UpdatePayScanBlockTaskStatusComfirm(orderInfo.Id); err != nil {
		log.Error(OperationID, "update pay scan block task status failed, err: ", err.Error())
		return err
	}
	log.Info("confirm ens order has been paid success, orderId: ", orderInfo.OrderId)
	return nil
}

func CompareStartAndEndBanlance(OperationID string, chainId int64, account string, start, end uint64) (bool, error) {
	client, err := utils.GetEthClient(chainId)
	if err != nil {
		log.Error(OperationID, "get eth client failed, err: ", err.Error())
		return false, err
	}
	defer client.Close()
	// 如果收款账户在此区间内有主动发生过交易，则需要重新扫块
	// var startNonce uint64
	// nonceCacheKey := fmt.Sprintf("%d:%s", chainId, account)
	// if val, ok := cacheNonceAtStartheight[nonceCacheKey]; !ok {
	// 	startNonce, err = client.NonceAt(context.Background(), common.HexToAddress(account), big.NewInt(int64(start)))
	// 	if err != nil {
	// 		log.Error(OperationID, "get start balance failed, err: ", err.Error())
	// 		return false, err
	// 	}
	// 	cacheNonceAtStartheight[nonceCacheKey] = startNonce
	// } else {
	// 	startNonce = val
	// }
	// endNonce, err := client.NonceAt(context.Background(), common.HexToAddress(account), big.NewInt(int64(end)))
	// if err != nil {
	// 	log.Error(OperationID, "get end nonce failed, err: ", err.Error())
	// 	return false, err
	// }
	// if startNonce != endNonce {
	// 	return false, nil
	// }
	// 对比账户余额
	// var startVal *big.Int
	startVal, err := client.BalanceAt(context.Background(), common.HexToAddress(account), big.NewInt(int64(start)))
	if err != nil {
		// 太久远的区块，无法获取到余额，直接返回 false
		log.Error(OperationID, "get start balance failed, err: ", err.Error())
		return false, nil
	}
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(account), big.NewInt(int64(end)))
	if err != nil {
		log.Error(OperationID, "get end balance failed, err: ", err.Error())
		return false, err
	}
	return balance.Cmp(startVal) == 0, nil
}

// 定义 Transfer 事件结构体
type TransferEvent struct {
	From        string
	To          string
	Value       *big.Int
	BlockNumber uint64
	TxHash      string
}

type ChainBlock struct {
	BlockNumber uint64
	Time        time.Time
}

type ScanResolver interface {
	GetCoinUSDPrice(OperationID string) (float64, error)
	Scan(fromAddresses []string, toAddress []string, start uint64, end uint64) ([]*TransferEvent, error)
	CompareStartAndEndBanlance(OperationID string, account string, start, end uint64) (bool, error)
	GetBlock(OperationID string, ctx context.Context, chainId int64, number *big.Int) (*ChainBlock, error)
}

// ScanFactory 扫块任务工厂
func ScanFactory(scanType, tag string, chainId int64) ScanResolver {
	switch scanType {
	case constant.PayScanBlockTaskMaticErc20USDTType:
		return NewERC20ScanResolver(tag, chainId)
	case constant.PayScanBlockTaskETHType:
		return NewEthScanResolver(tag, chainId)
	case constant.PayScanBlockTaskBNBType:
		return NewBNBScanResolver(tag, chainId)
	case constant.PayScanBlockTaskMaticType:
		return NewMaticScanResolver(tag, chainId)
	case constant.PayScanBlockTaskOPType:
		return NewOPScanResolver(tag, chainId)
	case constant.PayScanBlockTaskARBType:
		return NewARBScanResolver(tag, chainId)
	default:
		return nil
	}
}

func isERC20(txnType string) bool {
	// 判断 txnType 是否已 ERC20 开头
	return strings.HasPrefix(txnType, "ERC20")
}

func GetTxBlockTime(chainId int64, txHash string) (*time.Time, error) {
	client, err := utils.GetEthClient(chainId)
	if err != nil {
		return nil, err
	}
	tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}
	block, err := client.BlockByHash(context.Background(), tx.Hash())
	if err != nil {
		return nil, err
	}
	blockTime := time.Unix(int64(block.Time()), 0)
	return &blockTime, nil
}

type CoinPriceReq struct {
	CoinName string `json:"coinName"  bind:"require"`
}
type ComResp struct {
	ErrCode int
	ErrMsg  string
	Data    interface{}
}

func PostGraph(url, postData interface{}) (*ComResp, error) {
	PostCheckDataByte, _ := json.Marshal(postData)
	resultByte, err := utils.HttpPost(
		fmt.Sprintf("%s%s", config.Config.EnsPostCheck.Url, url),
		"", map[string]string{"Content-Type": "application/json"}, PostCheckDataByte)
	if err != nil {
		return nil, err
	}
	var resultData ComResp
	err = json.Unmarshal(resultByte, &resultData)
	return &resultData, err
	// var resultData ComResp
	// resultData = ComResp{
	// 	Data: float64(1.0),
	// }
	// return &resultData, nil
}

// 数组去重
func RemoveRepeatedElement(arr []string) []string {
	result := make([]string, 0, len(arr))
	temp := map[string]struct{}{}
	for _, item := range arr {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveFromArr(arr []string, val string) []string {
	for i, v := range arr {
		if v == val {
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}
	return arr
}
