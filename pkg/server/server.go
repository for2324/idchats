package server

import "context"

type Server interface {
	//Start 服务启动
	Start(ctx context.Context) error
	//Stop 服务关闭
	Stop(ctx context.Context) error
}
