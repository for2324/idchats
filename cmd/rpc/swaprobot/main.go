package main

import (
	"Open_IM/internal/rpc/swaprobotrpc"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.SwapRobotPort
	rpcPort := flag.Int("port", defaultPorts[0], "get SwapRobotPort from cmd,default 12000 as port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.SwapRobotPrometheusPort[0],
		"SwapRobotPort default listen port")
	flag.Parse()
	fmt.Println("start friend rpc server, port: ", *rpcPort, "OpenIM version: ", constant.CurrentVersion)
	rpcServer := swaprobotrpc.NewSwapRobotServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
