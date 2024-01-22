package order

import (
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/contracts/erc20"
	iutils "Open_IM/pkg/utils"
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func NewERC20ScanResolver(constract string, chainId int64) ScanResolver {
	return &ERC20ScanResolver{
		Constract: constract,
		ChainId:   chainId,
	}
}

type ERC20ScanResolver struct {
	Constract string
	ChainId   int64
}

func (e *ERC20ScanResolver) Scan(fromAddresses, toAddress []string, start, end uint64) ([]*TransferEvent, error) {
	client, err := iutils.GetEthClient(e.ChainId)
	if err != nil {
		log.Error("get eth client failed, err: ", err.Error())
		return nil, err
	}
	filterer, err := erc20.NewChainFilterer(common.HexToAddress(e.Constract), client)
	if err != nil {
		log.NewError("chain_event_listener NewChainFilterer faild", e.Constract, err)
		return nil, err
	}
	FromAddresses := []common.Address{}
	for _, address := range fromAddresses {
		FromAddresses = append(FromAddresses, common.HexToAddress(address))
	}
	ToAddresses := []common.Address{}
	for _, address := range toAddress {
		ToAddresses = append(ToAddresses, common.HexToAddress(address))
	}
	sub, err := filterer.FilterTransfer(&bind.FilterOpts{
		Start: start,
		End:   &end,
	}, FromAddresses, ToAddresses)
	if err != nil {
		log.Error("filter transfer failed, err: ", err.Error())
		return nil, err
	}
	transferEvent := make([]*TransferEvent, 0)
	for sub.Next() {
		transferEvent = append(transferEvent, &TransferEvent{
			From:        sub.Event.From.Hex(),
			To:          sub.Event.To.Hex(),
			Value:       sub.Event.Value,
			BlockNumber: sub.Event.Raw.BlockNumber,
			TxHash:      sub.Event.Raw.TxHash.Hex(),
		})
	}
	return transferEvent, nil
}

func (e *ERC20ScanResolver) CompareStartAndEndBanlance(OperationID string, account string, start, end uint64) (bool, error) {
	client, err := iutils.GetEthClient(e.ChainId)
	if err != nil {
		log.Error(OperationID, "get eth client failed, err: ", err.Error())
		return false, err
	}
	defer client.Close()
	// 如果收款账户在此区间内有主动发生过交易，则需要重新扫块
	// var startNonce uint64
	// nonceCacheKey := fmt.Sprintf("%d:%s:%s", e.ChainId, e.Constract, account)
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
	var startVal *big.Int
	erc20Cli, err := erc20.NewChain(common.HexToAddress(e.Constract), client)
	if err != nil {
		log.Error(OperationID, "get start balance failed, err: ", err.Error())
		return false, err
	}
	startVal, err = erc20Cli.BalanceOf(&bind.CallOpts{
		BlockNumber: big.NewInt(int64(start)),
	}, common.HexToAddress(account))
	if err != nil {
		log.Error(OperationID, "get start balance failed, err: ", err.Error())
		return false, err
	}
	balance, err := erc20Cli.BalanceOf(&bind.CallOpts{
		BlockNumber: big.NewInt(int64(end)),
	}, common.HexToAddress(account))
	if err != nil {
		log.Error(OperationID, "get end balance failed, err: ", err.Error())
		return false, err
	}
	return balance.Cmp(startVal) == 0, nil
}

func (e *ERC20ScanResolver) GetCoinUSDPrice(OperationID string) (float64, error) {
	// erc20CoinName := "usd"
	// resp, err := PostGraph("/graph/coinprice", CoinPriceReq{CoinName: erc20CoinName})
	// if err != nil {
	// 	return 0, err
	// }
	// if resp.ErrCode != 0 {
	// 	return 0, errors.New(resp.ErrMsg)
	// }
	// // 字币的兑u比例
	// coinPriceRate := resp.Data.(float64)
	// return coinPriceRate, nil
	return 1, nil
}

func (e *ERC20ScanResolver) GetBlock(OperationID string, ctx context.Context, chainId int64, number *big.Int) (*ChainBlock, error) {
	client, err := iutils.GetEthClient(chainId)
	if err != nil {
		return nil, err
	}

	block, err := client.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return &ChainBlock{
		BlockNumber: block.Number().Uint64(),
		Time:        time.Unix(int64(block.Time()), 0),
	}, nil
}
