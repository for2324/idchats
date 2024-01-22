package eth

import (
	"Open_IM/internal/contract/evmlisteninterface/daemon"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"github.com/sourcegraph/conc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/supermigo/xlog"
	"math/big"
	"time"
)

type FilterEthListener struct {
	EthInfo          IEthInfoRepo
	Log              chan types.Log
	EthClient        []*ethclient.Client
	EventFilters     []ethereum.FilterQuery
	EventConsumerMap map[string]*EventConsumer
	TxMonitors       map[common.Address]ITxMonitor
	errC             chan error
	blockTime        uint64
	blockOffset      int64
	chainID          string
	IndexClient      int32
}

func NewEthListenerByFilterQuery(
	ethClient []*ethclient.Client,
	blockTime uint64,
	blockOffset int64,
	chainID string,
	info IEthInfoRepo,
) *FilterEthListener {
	return &FilterEthListener{
		EthInfo:          info,
		EthClient:        ethClient,
		EventConsumerMap: make(map[string]*EventConsumer),
		TxMonitors:       make(map[common.Address]ITxMonitor),
		Log:              make(chan types.Log),
		errC:             make(chan error),
		blockTime:        blockTime,
		blockOffset:      blockOffset,
		chainID:          chainID,
		IndexClient:      0,
	}
}

func (s *FilterEthListener) AddFilterQuery(query ethereum.FilterQuery) {
	s.EventFilters = append(s.EventFilters, query)
}

func (s *FilterEthListener) RegisterConsumer(consumer IEventConsumer) error {
	consumerHandler, err := consumer.GetConsumer()
	if err != nil {
		xlog.CError("[eth listener] Unable to get consumer", err.Error())
		return err
	}
	for i := 0; i < len(consumerHandler); i++ {
		s.EventConsumerMap[KeyFromBEConsumer(consumerHandler[i].Address.Hex(), consumerHandler[i].Topic.Hex())] = consumerHandler[i]
	}

	s.EventFilters = append(s.EventFilters, consumer.GetFilterQuery()...)
	return nil
}
func (s *FilterEthListener) GetEthClient() *ethclient.Client {
	if int(s.IndexClient) >= len(s.EthClient) {
		s.IndexClient = 0
	}
	result := s.EthClient[s.IndexClient]
	s.IndexClient++
	return result

}
func (s *FilterEthListener) RegisterTxMonitor(monitor ITxMonitor) error {
	if monitor == nil {
		err := fmt.Errorf("nil monitor")
		xlog.CError(err.Error(), "[eth listener] Register nil monitor")
		return err
	}
	address := monitor.MonitoredAddress()
	if _, ok := s.TxMonitors[address]; !ok {
		s.TxMonitors[address] = monitor
		return nil
	}
	return fmt.Errorf("Monitor for " + fmt.Sprintf("0x%x", address) + " already existed")
}

func (s *FilterEthListener) Start(ctx context.Context) {
	daemon.BootstrapDaemons(ctx, s.Handling, s.Scan)
}

func (s *FilterEthListener) Handling(parentContext context.Context) (fn daemon.Daemon, err error) {
	fn = func() {
		for {
			select {
			case err := <-s.errC:
				xlog.CError(err.Error(), "[eth_listener] Ethereum client scan block err")
			case vLog := <-s.Log:
				//	go func(vLog types.Log) {
				s.consumeEvent(vLog, vLog.BlockNumber)
			//	}(vLog)

			case <-parentContext.Done():
				return
			}
		}

	}
	return fn, nil
}

// Scan  offset is the number of block to scan before current block to make sure event is confirmed
func (s *FilterEthListener) Scan(parentContext context.Context) (daemonfn daemon.Daemon, err error) {
	// scanners

	blockTransferScanner := func(from *big.Int, to *big.Int) {
		for i := from; i.Cmp(to) < 1; i = i.Add(i, big.NewInt(1)) {
			currBlock, err := s.GetEthClient().BlockByNumber(context.Background(), i)
			if err != nil {
				s.errC <- err
				continue
			}
			trans := currBlock.Transactions()
			for _, t := range trans {
				if t.To() == nil { // contract deployment, ignore
					continue
				}
				if len(t.Data()) > 0 {
					continue
				}

				if monitor, ok := s.TxMonitors[*t.To()]; ok {
					msg, err := core.TransactionToMessage(t, types.LatestSignerForChainID(t.ChainId()), nil)
					if err != nil {
						continue
					}
					from := msg.From.Hex()
					to := msg.To.Hex()
					amount := t.Value().String()
					_ = monitor.TxParse(t, from, to, monitor.MonitoredAddress().String(), amount, msg.Data)
				}
			}
		}
	}

	eventScanner := func(query ethereum.FilterQuery, from, to *big.Int) {
		query.FromBlock = from //scannedBlock
		query.ToBlock = to     //currBlock
		events, err := s.GetEthClient().FilterLogs(context.Background(), query)
		if err != nil {
			s.errC <- err
		}
		for _, event := range events {
			s.Log <- event
		}
	}

	// main scanning daemon
	daemonfn = func() {
		for {
			select {
			case <-parentContext.Done():
				return
			default:
				sysInfo, err := s.EthInfo.Get(s.chainID)
				if err != nil {
					xlog.CError("[eth_listener] can't get system info:", err.Error())
					continue
				}

				header, err := s.GetEthClient().HeaderByNumber(parentContext, nil)
				if err != nil {
					xlog.CError("[eth_listener] can't get head by number, possibly due to rpc node failure:", err.Error())
					continue
				} else {
					xlog.CInfo("当前高度:", header.Number.Int64())
				}
				currBlock := header.Number

				var scannedBlock *big.Int
				if sysInfo.LastScannedBlock <= 0 {
					// set first block
					scannedBlock = big.NewInt(s.blockOffset)
				} else {
					scannedBlock = big.NewInt(sysInfo.LastScannedBlock)
				}

				// scanned a offset - 1 block before to sure event confirmed
				// scannedBlock = scannedBlock.Sub(scannedBlock, big.NewInt(s.blockOffset-1))

				// if last scanned block is more than $BIGNUM blocks away just scan last $BIGNUM blocks
				diff := big.NewInt(0).Sub(currBlock, scannedBlock)
				if diff.Cmp(big.NewInt(600000)) > 0 {
					diff = big.NewInt(600000)
					scannedBlock = scannedBlock.Sub(currBlock, diff)
				} else {
					xlog.CInfo("ChainID:"+s.chainID+" 相差块：", diff.Int64())
				}

				if diff.Cmp(big.NewInt(100)) > 0 {
					// scan in 100-chunk
					for begin := scannedBlock; currBlock.Cmp(begin) > 0; begin = begin.Add(begin, big.NewInt(100)) {
						limit := big.NewInt(99)
						curDiff := big.NewInt(0)
						if curDiff.Sub(currBlock, begin).Cmp(limit) < 0 {
							limit = currBlock.Sub(currBlock, begin)
						}
						until := big.NewInt(0)
						until = until.Add(begin, limit)
						newBegin := begin
						newUtils := until
						var wg conc.WaitGroup
						wg.Go(func() {
							if len(s.TxMonitors) > 0 {
								blockTransferScanner(newBegin, newUtils)
							}
						})
						wg.Go(func() {
							if len(s.EventFilters) > 0 {
								xlog.CInfo("当前过滤事件:", len(s.EventFilters))
								for _, query := range s.EventFilters {
									eventScanner(query, newBegin, newUtils)
								}
							}
						})
						wg.Wait()
						// update last scan block
						sysInfo.LastScannedBlock = until.Int64()
						xlog.CInfo("更新当前最新块:", sysInfo.LastScannedBlock)
						_ = s.EthInfo.Update(s.chainID, sysInfo)
					}
				} else {
					xlog.CInfo("scannedBlock:", scannedBlock.Int64(), "currBlock:", currBlock.Int64(), currBlock.Int64()-scannedBlock.Int64())
					if (currBlock.Int64() - scannedBlock.Int64()) < 1 {
						xlog.CInfo("sleep 3s")
						daemon.SleepContext(parentContext, time.Second*time.Duration(s.blockTime))
						continue
					}
					scannedBlock = scannedBlock.Add(scannedBlock, big.NewInt(1))
					var wg conc.WaitGroup
					wg.Go(func() {
						if len(s.TxMonitors) > 0 {
							blockTransferScanner(scannedBlock, currBlock)

						}
					})
					wg.Go(func() {
						if len(s.EventFilters) > 0 {
							for _, query := range s.EventFilters {
								eventScanner(query, scannedBlock, currBlock)
							}
						}
					})
					wg.Wait()
					sysInfo.LastScannedBlock = currBlock.Int64()
					xlog.CInfo("更新当前最新块:", sysInfo.LastScannedBlock)
					_ = s.EthInfo.Update(s.chainID, sysInfo)
				}
				daemon.SleepContext(parentContext, time.Second*time.Duration(s.blockTime))
			}
		}
	}
	return daemonfn, nil
}

func (s *FilterEthListener) matchEvent(vLog types.Log) (*EventConsumer, bool) {
	key := KeyFromBEConsumer(vLog.Address.Hex(), vLog.Topics[0].Hex())
	consumer, isExisted := s.EventConsumerMap[key]
	return consumer, isExisted
}

func (s *FilterEthListener) consumeEvent(vLog types.Log, blockTime uint64) {
	consumer, isExisted := s.matchEvent(vLog)
	if isExisted {
		blockHeader, err := s.EthClient[0].HeaderByNumber(context.Background(), big.NewInt(int64(blockTime)))
		timstampBlock := uint64(time.Now().Unix())
		if blockHeader != nil {
			timstampBlock = blockHeader.Time
		}
		err = consumer.ParseEvent(vLog, timstampBlock)
		if err != nil {
			xlog.CError(err.Error(), "[eth_client] Consume event error")
		}
	}
}
