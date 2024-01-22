package walletRdbService

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	client "Open_IM/pkg/walletRdbService/swaprobotservice"
	"crypto/tls"
	"crypto/x509"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"io/ioutil"
	"time"
)

var GlobalSwapRobotService client.SwapRobotService

// tcp本身有保活动机制，所以只要建立一次链接即可， 2023.08.28
func GetRdbService() (client.SwapRobotService, error) {
	if GlobalSwapRobotService == nil {
		onceFunc := syncx.Once(func() {
			// 添加证书设置
			cred := GGetCreds()
			cred = cred
			clientConn, err := zrpc.NewClient(zrpc.RpcClientConf{
				Target: config.Config.WalletService.Server,
			}, zrpc.WithDialOption(grpc.WithTransportCredentials(cred)), zrpc.WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             30 * time.Second,
				PermitWithoutStream: false,
			})))
			if err != nil {
				log.NewInfo("err.e::::???????<>>>>>>>>>>", err.Error())
				GlobalSwapRobotService = nil
				return
			}
			GlobalSwapRobotService = client.NewSwapRobotService(clientConn)
		})
		onceFunc()
	}
	return GlobalSwapRobotService, nil
}

// getCreds 添加凭证
func GGetCreds() credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair("./tls/client.pem", "./tls/client.key")
	if err != nil {
		log.NewInfo("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	caFile, err := ioutil.ReadFile("./tls/ca.crt")
	if err != nil {
		log.NewInfo("加载 ca 失败!\n", err)
	}
	certPool.AppendCertsFromPEM(caFile)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
		RootCAs:      certPool,
	})

	return creds
}
