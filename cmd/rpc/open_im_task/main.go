package main

import (
	"Open_IM/internal/rpc/task"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImTaskPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcTaskPort from cmd,default 11500 as port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.TaskPrometheusPort[0], "taskPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start task rpc server, port: ", *rpcPort, "OpenIM version: ", constant.CurrentVersion)
	rpcServer := task.NewTaskServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
