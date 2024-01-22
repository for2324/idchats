package main

import (
	"Open_IM/internal/rpc/web3pub"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImWeb3JsPort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 11400")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.Web3PrometheusPort[0], "web3PrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start web3 rpc server, port: ", *rpcPort, "OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := web3pub.NewWeb3PubServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
