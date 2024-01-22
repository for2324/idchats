package main

import (
	"Open_IM/internal/rpc/order"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImOrderPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcOrderPort from cmd,default 11700 as port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.OrderPrometheusPort[0], "ensPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start ens rpc server, port: ", *rpcPort, "OpenIM version: ", constant.CurrentVersion)
	rpcServer := order.NewOrderServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
