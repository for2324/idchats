package utils

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/eip4361"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const addressChecksumLen = 4

func VerifySignature(publicAddress, signature, message string) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Info(">>>>>>>>>>>>>>>>>VerifySignature")
		}
	}()
	messageBytes := accounts.TextHash([]byte(message))
	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false
	}

	if signatureBytes[crypto.RecoveryIDOffset] == 27 || signatureBytes[crypto.RecoveryIDOffset] == 28 {
		signatureBytes[crypto.RecoveryIDOffset] -= 27
	}
	recoveredPublicKey, err := crypto.SigToPub(messageBytes, signatureBytes)
	if err != nil {
		log.NewDebug("web3 signature verify fail. public_address:%s, message:%s, signature:%s", publicAddress, message, signature)
		return false
	}
	recoveredPublicAddress := crypto.PubkeyToAddress(*recoveredPublicKey).Hex()

	return strings.EqualFold(publicAddress, recoveredPublicAddress)
}
func VerifySignatureEip4361(publicAddress, signature, message string) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Info(">>>>>>>>>>>>>>>>>VerifySignature")
		}
	}()
	msg, eerrrr := eip4361.ParseMessageFromParse(message)
	if eerrrr != nil || !strings.EqualFold(msg.Address.String(), publicAddress) {
		log.Info(">>>>>>>>>>>>>>>>>VerifySignature1 ")
		return false
	}
	if pkey, err := msg.VerifyEIP191(signature); err != nil {
		log.Info("验证不通过")
		return false
	} else if !strings.EqualFold(crypto.PubkeyToAddress(*pkey).String(), publicAddress) {
		log.Info("地址不想等")
		return false
	}
	return true
}

/*
判断钱包地址是否有效
*/
func AddressIsValid(address []byte) bool {
	//base58解码
	address = Base58Decode(address)

	//从address截取rmd160
	rmd160 := address[1 : len(address)-addressChecksumLen]

	//从address截取checksum
	checksum := address[len(address)-addressChecksumLen:]

	//计算得出checksum
	checksumCal := getAddressChecksum(rmd160)

	return bytes.Compare(checksum, checksumCal) == 0
}

/*
getAddressChecksum
*/
func getAddressChecksum(ripemd160Hash []byte) []byte {
	//双hash256
	hash256 := sha256.Sum256(ripemd160Hash)
	hash256 = sha256.Sum256(hash256[:])

	//取前checksumLen
	return hash256[:addressChecksumLen]
}

/*
ripemd160(hash256(publicKey))
*/
func ripemd160Hash(publicKey []byte) []byte {
	//hash256 publicKey
	hash256 := sha256.Sum256(publicKey)

	//ripemd160
	rmd160 := ripemd160.New()
	rmd160.Write(hash256[:])
	return rmd160.Sum(nil)
}

var base58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

/*
字节数组反转
*/
func reverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

/*
*base58编码
 */
func Base58Encode(input []byte) (output []byte) {
	if len(input) == 0 {
		return
	}

	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64(len(base58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		output = append(output, base58Alphabet[mod.Int64()])
	}

	reverseBytes(output)

	for b := range input {
		if b != 0x00 {
			break
		}
		output = append([]byte{base58Alphabet[0]}, output...)
	}

	return
}

/*
base58解码
*/
func Base58Decode(input []byte) (output []byte) {
	if len(input) == 0 {
		return
	}

	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(base58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	output = result.Bytes()
	output = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), output...)

	return
}
func ToValidateAddress(address string) string {
	addrLowerStr := strings.ToLower(address)
	if strings.HasPrefix(addrLowerStr, "0x") {
		addrLowerStr = addrLowerStr[2:]
		address = address[2:]
	}
	var binaryStr string
	addrBytes := []byte(addrLowerStr)
	hash256 := Keccak256Hash([]byte(addrLowerStr)) //注意，这里是直接对字符串转换成byte切片然后哈希

	for i, e := range addrLowerStr {
		//如果是数字则跳过
		if e >= '0' && e <= '9' {
			continue
		} else {
			binaryStr = fmt.Sprintf("%08b", hash256[i/2]) //注意，这里一定要填充0
			if binaryStr[4*(i%2)] == '1' {
				addrBytes[i] -= 32
			}
		}
	}

	return "0x" + string(addrBytes)
}

func Keccak256Hash(data []byte) []byte {
	keccak256Hash2 := sha3.NewLegacyKeccak256()
	keccak256Hash2.Write(data)
	return keccak256Hash2.Sum(nil)
}

// 检查有大小写区别的以太坊地址是否合法
func CheckEthAddress(address string) bool {
	return strings.ToLower(ToValidateAddress(address)) == address
}

func GetEthClient(chainId int64) (*ethclient.Client, error) {
	chainKey := strconv.FormatInt(chainId, 10)
	rpcInfo, ok := config.Config.ChainIdRpcMap[chainKey]
	if !ok || len(rpcInfo) == 0 {
		return nil, errors.New("GetEnsInstant failed chainId not found")
	}
	endpoint := rpcInfo[0]
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type CoinUSDPriceResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Ethbtc          string `json:"ethbtc"`
		EthbtcTimestamp string `json:"ethbtc_


timestamp"`
		Ethusd          string `json:"ethusd"`
		EthusdTimestamp string `json:"ethusd_timestamp"`
	} `json:"result"`
}

func ChainHttpGet(uri string, v any) error {
	var resutlbyte []byte
	var err error
	if config.Config.OpenNetProxy.OpenFlag {
		proxyAddress, _ := url.Parse(config.Config.OpenNetProxy.ProxyURL)
		resutlbyte, err = HttpGetWithHeaderWithProxy(uri, map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		}, http.ProxyURL(proxyAddress))
	} else {
		resutlbyte, err = HttpGetWithHeader(uri, map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		})
	}
	if err != nil {
		log.NewError("Failed to send API request: %v", uri, err)
		return err
	}
	body := resutlbyte
	// Parse the JSON response
	if err := json.Unmarshal(body, &v); err != nil {
		log.NewError("Failed to parse API response: %v", uri, err)
		return err
	}
	return nil
}

func GetCoinUSDPrice(chainId int64) (uint64, error) {
	endPointConf := config.Config.ChainIdHttpMap[chainId]
	if endPointConf.EndPoint == "" {
		return 0, errors.New("chainId not support")
	}
	uri := fmt.Sprintf("%s/api?module=stats&action=ethprice&apikey=%s", endPointConf.EndPoint, endPointConf.ApiKey)
	var priceResult CoinUSDPriceResult
	err := ChainHttpGet(uri, &priceResult)
	if err != nil {
		return 0, err
	}
	if priceResult.Status != "1" {
		return 0, errors.New(priceResult.Message)
	}
	Ethusd := priceResult.Result.Ethusd
	priceRate, err := strconv.ParseFloat(Ethusd, 64)
	if err != nil {
		return 0, err
	}
	return uint64(priceRate * math.Pow10(6)), nil
}

// 预估消费的 gas fee
func SuggestGasFee(chainId int64) (uint64, error) {
	client, err := GetEthClient(chainId)
	if err != nil {
		return 0, err
	}
	// gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return 0, err
	}
	zeroAddr := common.HexToAddress("0x0000000")
	// 一笔转账预估的 gas used
	gasUsed, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To: &zeroAddr,
		// GasPrice: gasPrice,
		Value: big.NewInt(0),
		Data:  []byte{},
	})
	if err != nil {
		return 0, err
	}
	return gasPrice.Uint64() * gasUsed, nil
}
