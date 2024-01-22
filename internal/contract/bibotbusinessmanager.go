package contract

import (
	"Open_IM/internal/contract/evmlisteninterface/eth"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/go-redsync/redsync/v4"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type NodeContract struct {
	ContractAddr    string
	NodeContractPtr *BBTTradeReward
	ChainID         string
}

func NewRewardManagerContract(addr string, ChainID string, ethClientPtr *ethclient.Client) *NodeContract {
	result := new(NodeContract)
	result.ContractAddr = addr
	result.ChainID = ChainID
	result.NodeContractPtr, _ = NewBBTTradeReward(common.HexToAddress(addr), ethClientPtr)
	return result
}
func (node *NodeContract) GetConsumer() ([]*eth.EventConsumer, error) {
	parsed, _ := abi.JSON(strings.NewReader(BBTTradeRewardMetaData.ABI))
	return []*eth.EventConsumer{
		{
			Address: common.HexToAddress(node.ContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(parsed.Events["ClaimReward"].Sig),
			),
			ParseEvent: node.ClaimRewardParse,
		},
	}, nil
}
func (node *NodeContract) GetFilterQuery() []ethereum.FilterQuery {
	parsed, _ := abi.JSON(strings.NewReader(BBTTradeRewardMetaData.ABI))
	topics := []common.Hash{
		crypto.Keccak256Hash([]byte(parsed.Events["ClaimReward"].Sig)),
	}
	return []ethereum.FilterQuery{{
		Addresses: []common.Address{common.HexToAddress(node.ContractAddr)},
		Topics:    [][]common.Hash{topics}},
	}
}
func (node *NodeContract) ClaimRewardParse(logValue types.Log, blockNumber uint64) error {
	ContractparseClaimRewardPtr, err := node.NodeContractPtr.ParseClaimReward(logValue)
	if err != nil {
		return err
	}

	mutexname := "trade_volume:" + ContractparseClaimRewardPtr.User.String()
	rs := db.DB.Pool
	mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
	if err := mutex.LockContext(context.Background()); err != nil {
		return err
	}
	defer mutex.UnlockContext(context.Background())

	db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		var oldData db.UserHistoryTotal
		tx.Table("user_history_total").Where("user_id=?", ContractparseClaimRewardPtr.User.String()).First(&oldData)
		if oldData.UserID == "" {
			return errors.New("无法查询到用户")
		}
		//如果存在数据 那么就更新update 下

		decimalDataPending, _ := decimal.NewFromString(oldData.Pending)
		decimalDataRakebackPending, _ := decimal.NewFromString(oldData.RakebackPending)
		customerDataArray := strings.Split(ContractparseClaimRewardPtr.CustomData, "&")
		customerDataType := "tradeReward"
		bbtPrice, _ := decimal.NewFromString(config.Config.RewardTradeByScore)
		if len(customerDataArray) >= 2 {
			customerDataType = customerDataArray[0]
			bbtPrice, _ = decimal.NewFromString(customerDataArray[1])

		}
		decimalOnChain := decimal.NewFromBigInt(ContractparseClaimRewardPtr.Amount, 0).Shift(-18).Mul(bbtPrice)
		switch customerDataType {
		case "tradeReward":
			decimalDataClaim, _ := decimal.NewFromString(oldData.Claimed)
			oldData.Pending = decimalDataPending.Sub(decimalOnChain).String()
			oldData.Claimed = decimalDataClaim.Add(decimalOnChain).String()
		case "rakebackReward":
			decimalDataClaim, _ := decimal.NewFromString(oldData.RakebackClaimed)
			oldData.RakebackPending = decimalDataRakebackPending.Sub(decimalOnChain).String()
			oldData.RakebackClaimed = decimalDataClaim.Add(decimalOnChain).String()
		}
		oldData.CurrentNonce = ""
		decimalOnChainClaim, _ := decimal.NewFromString(oldData.OnChainClaimed)
		oldData.OnChainClaimed = decimalOnChainClaim.Add(decimal.NewFromBigInt(ContractparseClaimRewardPtr.Amount, 0)).String()
		tx.Table("user_history_total").Save(oldData)
		return nil
	})
	return nil
}

type Eip712Transfer struct {
	Amount    *big.Int
	Recipient common.Address
	Nonce     *big.Int
	Custom    string
}

func SignTransfer(
	privateKey *ecdsa.PrivateKey, verifyingContract common.Address, verifyingContractChainID int64,
	reg *Eip712Transfer,
) ([]byte, []byte, error) {
	data := &apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": {
				{
					Name: "name",
					Type: "string",
				},
				{
					Name: "version",
					Type: "string",
				},
				{
					Name: "chainId",
					Type: "uint256",
				},
				{
					Name: "verifyingContract",
					Type: "address",
				},
			},
			"ClaimInfo": {
				{
					Name: "reward",
					Type: "uint256",
				},
				{
					Name: "recipient",
					Type: "address",
				},
				{
					Name: "nonce",
					Type: "uint256",
				},
				{
					Name: "custom",
					Type: "string",
				},
			},
		},
		Domain: apitypes.TypedDataDomain{
			Name:              "BiBot",
			Version:           "1.0.0",
			ChainId:           math.NewHexOrDecimal256(verifyingContractChainID),
			VerifyingContract: verifyingContract.String(),
		},
		PrimaryType: "ClaimInfo",
		Message: map[string]interface{}{
			"custom":    reg.Custom,
			"reward":    (*hexutil.Big)(reg.Amount).String(),
			"recipient": reg.Recipient.Hex(),
			"nonce":     (*hexutil.Big)(reg.Nonce).String(),
		},
	}
	typedData := apitypes.TypedData{
		Types:       data.Types,
		PrimaryType: data.PrimaryType,
		Domain:      data.Domain,
		Message:     data.Message,
	}
	encoded, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return nil, nil, err
	}
	log.NewInfo("encode is：", hex.EncodeToString(encoded))
	log.NewInfo("encode is>>>：", utils.StructToJsonString(typedData))
	sig, err := crypto.Sign(encoded, privateKey)
	if err != nil {
		return encoded, nil, err
	}
	return encoded, sig, nil
}
func GetRsv(signature []byte) (*big.Int, *big.Int, uint8) {
	return ecdsaSignatureToRSV(signature)
}
func ecdsaSignatureToRSV(sig []byte) (*big.Int, *big.Int, uint8) {
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	v := sig[64] + 27
	return r, s, v
}

func GetRewardRpcClient() (ethCli *ethclient.Client) {
	if !config.Config.IsPublicEnv {
		proxyURL := "http://proxy.idchats.com:7890" // 替换为你的代理URL
		urlProxy, _ := url.Parse(proxyURL)
		transport := &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
		}
		client := &http.Client{
			Transport: transport,
			Timeout:   15 * time.Second,
		}
		rpcClient, _ := rpc.DialOptions(context.Background(), config.Config.RewardChainRpc, rpc.WithHTTPClient(client))
		ethCli = ethclient.NewClient(rpcClient)
	} else {
		ethCli, _ = ethclient.Dial(config.Config.RewardChainRpc)
	}
	return ethCli
}
