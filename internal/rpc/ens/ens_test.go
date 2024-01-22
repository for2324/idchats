package ens

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/contracts/erc20"
	pbEns "Open_IM/pkg/proto/ens"
	"Open_IM/pkg/utils"
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestCreateRegisterEnsOrder(t *testing.T) {
	s := &EnsServer{}
	reps, err := s.CreateRegisterEnsOrder(context.Background(), &pbEns.CreateRegisterEnsOrderReq{
		OperationID: "test CreateRegisterEnsOrder",
		EnsName:     "te3trr",
		UserId:      "0xCBD033Ea3C05Dc9504610061C86C7aE191C5c913",
		TxnType:     "ETH",
	})
	if err != nil {
		t.Errorf("CreateRegisterEnsOrder error, err=%s", err.Error())
		return
	}
	t.Logf("CreateRegisterEnsOrder success, resp=%s", reps)
}

func TestGetEnsOrder(t *testing.T) {
	s := &EnsServer{}
	reps, err := s.GetEnsOrderInfo(context.Background(), &pbEns.GetEnsOrderInfoReq{
		OperationID: "test CreateRegisterEnsOrder",
		OrderId:     6,
		UserId:      "0xCBD033Ea3C05Dc9504610061C86C7aE191C5c913",
	})
	if err != nil {
		t.Errorf("GetEnsOrderInfo error, err=%s", err.Error())
		return
	}
	t.Logf("GetEnsOrderInfo success, resp=%s", reps)
}

func TestConfirmEnsOrderHasBeenPaid(t *testing.T) {
	s := &EnsServer{}
	reps, err := s.ConfirmEnsOrderHasBeenPaid(context.Background(), &pbEns.ConfirmEnsOrderHasBeenPaidReq{
		OperationID: "test CreateRegisterEnsOrder",
		OrderId:     "2",
		TxnHash:     "123",
	})
	if err != nil {
		t.Errorf("CreateRegisterEnsOrder error, err=%s", err.Error())
		return
	}
	t.Logf("CreateRegisterEnsOrder success, resp=%s", reps)
}

func TestCommitAndRegister(t *testing.T) {
	// Wait for transaction to be mined
	OperationID := "TestCommitAndRegister"
	var ChainId int64 = 80001
	client, err := GetEthClient(OperationID, ChainId)
	if err != nil {
		t.Errorf(OperationID, "Failed to connect to Ethereum network: %v", err)
		return
	}
	secretStr := "es77PCh0aZFpv/PyxuhgEpVGTz13C/sr7NHpml+z76E="
	secret, _ := base64.StdEncoding.DecodeString(secretStr)
	var Secret [32]byte
	copy(Secret[:], secret)
	ensDomain := "kest153765"
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
	bytes, err := arguments.Pack(ensDomainHash, common.HexToAddress("0x20B17db9a65D8a48f14338Efd829a5C77Ef8302d"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sig := "setAddr(bytes32,address)"
	id := crypto.Keccak256([]byte(sig))[:4]
	byteArray := append(id, bytes...)
	registerData := IETHRegistrarControllerRegisterData{
		Name:                 ensDomain,
		Owner:                common.HexToAddress("0xCBD033Ea3C05Dc9504610061C86C7aE191C5c913"),
		Secret:               Secret,
		Duration:             GetUint256Max(),
		Resolver:             common.HexToAddress(config.Config.Ens.Resolver),
		Data:                 [][]byte{byteArray},
		ReverseRecord:        true, // 反向解析记录
		OwnerControlledFuses: 1,    //只有owner 控制
		RebateName:           "",
	}

	// 注册上链
	tx, err := RegisterEnsChainUp(OperationID, ChainId, registerData, true)
	if err != nil {
		t.Errorf(OperationID, utils.GetSelfFuncName(), "GetEnsInstant failed", err.Error())
		return
	}
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Errorf(OperationID, "Failed to mine transaction: %v", err)
		return
	}
	if receipt.Status == types.ReceiptStatusFailed {
		return
	}
	t.Logf(OperationID, "RegisterEnsChainUp success")
}

func TestGUint256Max(t *testing.T) {
	fmt.Printf("max=%v\n", GetUint256Max())
}

func TestRestartRegisterMind(t *testing.T) {
	RestartRegisterMind()
}

func TestApprove(t *testing.T) {
	OperationID := "TestApprove"
	var chainId int64 = 80001
	client, err := GetEthClient(OperationID, chainId)
	if err != nil {
		log.NewError(OperationID, "Failed to connect to Ethereum network: %v", err)
		return
	}
	privateKey := "de5e623371122fa538c578265927129ef630561b6b9fc8b8641049ea6257f548"
	// b715c521c32503d125258aed0fd95ab1f7e224cbb269a45eab5220b656a8ca2a
	// privateKey := config.Config.Ens.EnsOwnerPrivateKeyHex
	token := "0xC0AC5CCF66c08a3115B0dd984aab7D07587D5f76" // 某个币种合约地址
	// fromAddress := "0xCBD033Ea3C05Dc9504610061C86C7aE191C5c913"
	toContractAddress := "0x3f71043213ae8aC3931309d04636C62dD09E459e"
	MaxApproveMoney := GetUint256Max()
	senderPrivateKey, err := crypto.HexToECDSA(string(privateKey))
	if err != nil {
		log.NewError(OperationID, "Failed to parse private key: %v", err)
		return
	}
	senderAuth, err := bind.NewKeyedTransactorWithChainID(senderPrivateKey, big.NewInt(chainId))
	if err != nil {
		log.NewError(OperationID, "Failed to create authorized transactor: %v", err)
		return
	}
	erc20Address := common.HexToAddress(token)
	caller, err := erc20.NewChainTransactor(erc20Address, client)
	if err != nil {
		log.NewError(OperationID, "Failed to connect to Ethereum network: %v", err)
		return
	}
	tx, err := caller.Approve(senderAuth, common.HexToAddress(toContractAddress), MaxApproveMoney)
	if err != nil {
		log.NewError(OperationID, "Failed to connect to Ethereum network: %v", err)
		return
	}
	fmt.Println("Approve Tx: " + tx.Hash().String())
	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.NewError(OperationID, "Failed to mine transaction: %v", err)
		return
	}
	// 交易被回滚了
	if receipt.Status == types.ReceiptStatusFailed {
		log.NewError(OperationID, "Failed to mine transaction: %v", err)
	} else {
		log.NewInfo(OperationID, "Approve success")
	}

}

type TestEnsRegisterOrder struct {
	CreateAt time.Time
	ExpireAt time.Time `json:"expireAt"`
}

func (e *TestEnsRegisterOrder) ExpireTime() string {
	return e.ExpireAt.Format("2006-01-02 15:04:05")
}

type TestApiEnsRegisterOrder struct {
	CreateTime string
	ExpireTime string `json:"expireTime"`
}

func (e *TestApiEnsRegisterOrder) CreateAt(cTime time.Time) {
	e.CreateTime = cTime.Format("2006-01-02 15:04:05")
}

func TestCopyTime(t *testing.T) {
	apiOrder := TestApiEnsRegisterOrder{}
	err := utils.CopyStructFields(&apiOrder, &TestEnsRegisterOrder{
		CreateAt: time.Now(),
		ExpireAt: time.Now().Add(time.Hour * 24 * 365),
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(apiOrder)
}

// const ENSResolverABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// // 查询 ENS 域名
// func GetENSName(client *ethclient.Client, tokenId *big.Int) (string, error) {
//     // ENS Resolver 合约地址
//     resolverAddress := common.HexToAddress("0x3CA935BFFb76789b7b87Fb983358c725CA7d68fF")

//     // 创建 ENS Resolver 合约实例
//     resolverInstance, err := bindENSResolver(client, resolverAddress)
//     if err != nil {
//         return "", err
//     }

//     // 查询域名
//     name, err := resolverInstance.Name(nil, common.BigToHash(tokenId))
//     if err != nil {
//         return "", err
//     }

//     return name, nil
// }

// // 创建 ENS Resolver 合约实例
// func bindENSResolver(client *ethclient.Client, contractAddress common.Address) (*ENSResolver, error) {
//     // 创建 ENS Resolver 合约实例
//     resolverInstance, err := NewENSResolver(contractAddress, client)
//     if err != nil {
//         return nil, err
//     }

//     // 返回绑定后的合约实例
//     return resolverInstance, nil
// }

// func TestQueryEnsName(t *testing.T) {
// 	tokenId := big.NewInt(123456)
//     name, err := GetENSName(client, tokenId)
// }
