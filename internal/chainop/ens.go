package chainop

import (
	api "Open_IM/pkg/base_info"
	biubiuens "Open_IM/pkg/biubiuens"
	"Open_IM/pkg/common/config"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// BindEnsCheck 例如查询某个人是否有权限设置某个域名信息，目前以biu 来做方案
// BindEnsCheck 查询ens 相关的信息
func BindEnsCheck(ctx *gin.Context) {
	var req api.EnsApiReq
	var responseData api.EnsApiResp
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	rpcChainID := "80001"
	if config.Config.IsPublicEnv {
		rpcChainID = "137"
	}

	rcplist := config.GetRpcFromChainID(rpcChainID)
	for _, valueUrl := range rcplist {
		tempValueUrl := valueUrl
		client, err := ethclient.Dial(tempValueUrl)
		if err != nil {
			continue
		}
		name, addr1, addr2, addr3, err := biubiuens.ReverseResolveEnsName(client, common.HexToAddress(config.Config.Ens.UniversalResolverContract), common.HexToAddress(req.Address))
		if err != nil {
			fmt.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
			return
		} else {
			if name == req.EnsDomain {
				ctx.JSON(http.StatusBadRequest, responseData)
				return
			} else {
				logrus.Info("当前配置的值得未：", name, addr1.Hex(), addr2.Hex(), addr3.Hex())
				if name == "" {
					ctx.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "当前的地址未配置主域名"})
					return
				} else {
					ctx.JSON(http.StatusBadRequest, gin.H{"errCode": 402, "errMsg": "当前的地址配置主域名:" + name})
					return
				}
			}
		}
	}
}
