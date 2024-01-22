package main

import (
	"Open_IM/internal/chainop"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/utils"
	"Open_IM/pkg/xlog"
	"flag"
	"fmt"
	"io"
	olog "log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	xlog.InitLevel(xlog.DefaultOption())
	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("../logs/enssearch.log")
	olog.SetOutput(new(redirector))
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.Default()
	r.Use(utils.CorsHandler())
	r.Use(gin.Recovery())
	if config.Config.Prometheus.Enable {
		r.GET("/metrics", promePkg.PrometheusHandler())
	}
	chainop.InstallKafkaProduct()
	g := r.Group("/graph")
	{
		g.POST("/ens", chainop.GraphSqlEnsDomain)
		g.POST("/enscheck", chainop.GraphSqlEnsDomainByName)
		g.GET("/testens", chainop.GraphSqlEnsDomainTest)
		g.GET("/systemmsg", chainop.SystemConfigInfo)
		g.GET("/whitelist", chainop.GetWhiteList)
		g.POST("/checkAddress", chainop.CheckIsHaveGuanFangNft)
		g.POST("/requesttokenuri", chainop.TokenUrl)
		g.POST("/tokenOwnerAddress", chainop.TokenOwnerByTokeID)
		g.POST("/tokenOwnerAddressContractChainID", chainop.TokenOwnerByTokeIDChainIDContract)
		g.POST("/gettokenurl", chainop.GetTokenUriByTokenID)
		g.GET("/appversion/:platform", chainop.GetAppVersion)
		g.GET("/coin", chainop.GetCoin)
		g.GET("/tokenlist/:chainid", chainop.GetTokenList)
		g.GET("/postnewtoken/:chainid/token/:tokenaddress", chainop.PostNewToken)
		g.POST("/getEmailCode", chainop.GetEmailCode)
		g.POST("/registerEmailCode", chainop.RegisterEmailCode)
		g.POST("/registerEmailCodeV2", chainop.RegisterEmailCodeWithOutPassword)

		g.POST("/getPrivateEmailCodeV2", chainop.GetEmailCodeV2)
		g.POST("/changeEmailCodeV2", chainop.ChangeEmailPrivateWithOutPassword)

		g.POST("/loginEmail", chainop.LoginEmail)
		g.POST("/parse_url", chainop.FetchLink)
		g.GET("/redirecturl", chainop.RedirectUrlLink)
		g.GET("/redirecturlv2", chainop.RedirectUrlLinkV2)
		g.POST("/coinprice", chainop.CoinPrice)
		g.POST("/bindEnsCheck", chainop.BindEnsCheck)
	}
	go func() {
		//30分钟刷一次时间
		wait.Until(func() {
			chainop.UpdateSystemInfo()
		}, time.Duration(config.Config.UpdateSystemCountMine)*time.Minute, wait.NeverStop)

	}()

	defaultPorts := config.Config.EnsSearch.Port
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10030 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	address = config.Config.EnsSearch.ListenIP + ":" + strconv.Itoa(*ginPort)
	fmt.Println("start demo api server address: ", address, "OpenIM version: ", constant.CurrentVersion, "\n")
	err := r.Run(address)
	chainop.CloseKafka()
	if err != nil {
		log.Error("", "run failed ", *ginPort, err.Error())
		fmt.Println(err.Error())
	}

}

type redirector struct{}

func (r *redirector) Write(p []byte) (n int, err error) {
	log.Info("redirector", string(p))
	return len(p), nil
}
