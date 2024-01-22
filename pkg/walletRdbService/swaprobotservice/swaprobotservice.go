// Code generated by goctl. DO NOT EDIT.
// Source: mnemonic.proto

package swaprobotservice

import (
	"context"

	__serviceprotoc "Open_IM/pkg/walletRdbService/serviceprotoc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	BaseResp                 = __serviceprotoc.BaseResp
	Ceshi                    = __serviceprotoc.Ceshi
	CheckIsNeedSign          = __serviceprotoc.CheckIsNeedSign
	CheckIsNeedSignResp      = __serviceprotoc.CheckIsNeedSignResp
	CreateWalletMnemonicReq  = __serviceprotoc.CreateWalletMnemonicReq
	CreateWalletMnemonicResp = __serviceprotoc.CreateWalletMnemonicResp
	ImportWalletMnemonicReq  = __serviceprotoc.ImportWalletMnemonicReq
	RequestUserID            = __serviceprotoc.RequestUserID
	UserRobotResp            = __serviceprotoc.UserRobotResp

	SwapRobotService interface {
		// 创建机器人 获取到助记词
		CreateUserRobot(ctx context.Context, in *CreateWalletMnemonicReq, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error)
		// 删除某个机器人
		DeleteUserRobot(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*UserRobotResp, error)
		// 检查内存保存的内容是否过期
		CheckIsExpireTime(ctx context.Context, in *CheckIsNeedSign, opts ...grpc.CallOption) (*CheckIsNeedSignResp, error)
		// 某个用户请求助记词
		GetMnemonic(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error)
		// 某个用户过来助记词给
		GetMnemonicFromMemory(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error)
		// 重新生成载入助记词
		ReloadMnemonic(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*UserRobotResp, error)
		// 测试内容
		HelloWorldTest(ctx context.Context, in *Ceshi, opts ...grpc.CallOption) (*UserRobotResp, error)
		// 导入助记词
		ImportWallet(ctx context.Context, in *ImportWalletMnemonicReq, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error)
	}

	defaultSwapRobotService struct {
		cli zrpc.Client
	}
)

func NewSwapRobotService(cli zrpc.Client) SwapRobotService {
	return &defaultSwapRobotService{
		cli: cli,
	}
}

// 创建机器人 获取到助记词
func (m *defaultSwapRobotService) CreateUserRobot(ctx context.Context, in *CreateWalletMnemonicReq, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error) {
	client := __serviceprotoc.NewSwapRobotServiceClient(m.cli.Conn())
	return client.CreateUserRobot(ctx, in, opts...)
}

// 删除某个机器人
func (m *defaultSwapRobotService) DeleteUserRobot(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*UserRobotResp, error) {
	client := __serviceprotoc.NewSwapRobotServiceClient(m.cli.Conn())
	return client.DeleteUserRobot(ctx, in, opts...)
}

// 检查内存保存的内容是否过期
func (m *defaultSwapRobotService) CheckIsExpireTime(ctx context.Context, in *CheckIsNeedSign, opts ...grpc.CallOption) (*CheckIsNeedSignResp, error) {
	client := __serviceprotoc.NewSwapRobotServiceClient(m.cli.Conn())
	return client.CheckIsExpireTime(ctx, in, opts...)
}

// 某个用户请求助记词
func (m *defaultSwapRobotService) GetMnemonic(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error) {
	client := __serviceprotoc.NewSwapRobotServiceClient(m.cli.Conn())
	return client.GetMnemonic(ctx, in, opts...)
}

// 某个用户过来助记词给
func (m *defaultSwapRobotService) GetMnemonicFromMemory(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error) {
	client := __serviceprotoc.NewSwapRobotServiceClient(m.cli.Conn())
	return client.GetMnemonicFromMemory(ctx, in, opts...)
}

// 重新生成载入助记词
func (m *defaultSwapRobotService) ReloadMnemonic(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*UserRobotResp, error) {
	client := __serviceprotoc.NewSwapRobotServiceClient(m.cli.Conn())
	return client.ReloadMnemonic(ctx, in, opts...)
}

// 测试内容
func (m *defaultSwapRobotService) HelloWorldTest(ctx context.Context, in *Ceshi, opts ...grpc.CallOption) (*UserRobotResp, error) {
	client := __serviceprotoc.NewSwapRobotServiceClient(m.cli.Conn())
	return client.HelloWorldTest(ctx, in, opts...)
}

// 导入助记词
func (m *defaultSwapRobotService) ImportWallet(ctx context.Context, in *ImportWalletMnemonicReq, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error) {
	client := __serviceprotoc.NewSwapRobotServiceClient(m.cli.Conn())
	return client.ImportWallet(ctx, in, opts...)
}
