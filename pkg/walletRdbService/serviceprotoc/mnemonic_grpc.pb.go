// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: protoc/mnemonic.proto

package serviceprotoc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SwapRobotServiceClient is the client API for SwapRobotService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SwapRobotServiceClient interface {
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

type swapRobotServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSwapRobotServiceClient(cc grpc.ClientConnInterface) SwapRobotServiceClient {
	return &swapRobotServiceClient{cc}
}

func (c *swapRobotServiceClient) CreateUserRobot(ctx context.Context, in *CreateWalletMnemonicReq, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error) {
	out := new(CreateWalletMnemonicResp)
	err := c.cc.Invoke(ctx, "/mnemonicService.SwapRobotService/CreateUserRobot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swapRobotServiceClient) DeleteUserRobot(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*UserRobotResp, error) {
	out := new(UserRobotResp)
	err := c.cc.Invoke(ctx, "/mnemonicService.SwapRobotService/DeleteUserRobot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swapRobotServiceClient) CheckIsExpireTime(ctx context.Context, in *CheckIsNeedSign, opts ...grpc.CallOption) (*CheckIsNeedSignResp, error) {
	out := new(CheckIsNeedSignResp)
	err := c.cc.Invoke(ctx, "/mnemonicService.SwapRobotService/CheckIsExpireTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swapRobotServiceClient) GetMnemonic(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error) {
	out := new(CreateWalletMnemonicResp)
	err := c.cc.Invoke(ctx, "/mnemonicService.SwapRobotService/GetMnemonic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swapRobotServiceClient) GetMnemonicFromMemory(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error) {
	out := new(CreateWalletMnemonicResp)
	err := c.cc.Invoke(ctx, "/mnemonicService.SwapRobotService/GetMnemonicFromMemory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swapRobotServiceClient) ReloadMnemonic(ctx context.Context, in *RequestUserID, opts ...grpc.CallOption) (*UserRobotResp, error) {
	out := new(UserRobotResp)
	err := c.cc.Invoke(ctx, "/mnemonicService.SwapRobotService/ReloadMnemonic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swapRobotServiceClient) HelloWorldTest(ctx context.Context, in *Ceshi, opts ...grpc.CallOption) (*UserRobotResp, error) {
	out := new(UserRobotResp)
	err := c.cc.Invoke(ctx, "/mnemonicService.SwapRobotService/HelloWorldTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swapRobotServiceClient) ImportWallet(ctx context.Context, in *ImportWalletMnemonicReq, opts ...grpc.CallOption) (*CreateWalletMnemonicResp, error) {
	out := new(CreateWalletMnemonicResp)
	err := c.cc.Invoke(ctx, "/mnemonicService.SwapRobotService/ImportWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SwapRobotServiceServer is the server API for SwapRobotService service.
// All implementations must embed UnimplementedSwapRobotServiceServer
// for forward compatibility
type SwapRobotServiceServer interface {
	// 创建机器人 获取到助记词
	CreateUserRobot(context.Context, *CreateWalletMnemonicReq) (*CreateWalletMnemonicResp, error)
	// 删除某个机器人
	DeleteUserRobot(context.Context, *RequestUserID) (*UserRobotResp, error)
	// 检查内存保存的内容是否过期
	CheckIsExpireTime(context.Context, *CheckIsNeedSign) (*CheckIsNeedSignResp, error)
	// 某个用户请求助记词
	GetMnemonic(context.Context, *RequestUserID) (*CreateWalletMnemonicResp, error)
	// 某个用户过来助记词给
	GetMnemonicFromMemory(context.Context, *RequestUserID) (*CreateWalletMnemonicResp, error)
	// 重新生成载入助记词
	ReloadMnemonic(context.Context, *RequestUserID) (*UserRobotResp, error)
	// 测试内容
	HelloWorldTest(context.Context, *Ceshi) (*UserRobotResp, error)
	// 导入助记词
	ImportWallet(context.Context, *ImportWalletMnemonicReq) (*CreateWalletMnemonicResp, error)
	mustEmbedUnimplementedSwapRobotServiceServer()
}

// UnimplementedSwapRobotServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSwapRobotServiceServer struct {
}

func (UnimplementedSwapRobotServiceServer) CreateUserRobot(context.Context, *CreateWalletMnemonicReq) (*CreateWalletMnemonicResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserRobot not implemented")
}
func (UnimplementedSwapRobotServiceServer) DeleteUserRobot(context.Context, *RequestUserID) (*UserRobotResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserRobot not implemented")
}
func (UnimplementedSwapRobotServiceServer) CheckIsExpireTime(context.Context, *CheckIsNeedSign) (*CheckIsNeedSignResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckIsExpireTime not implemented")
}
func (UnimplementedSwapRobotServiceServer) GetMnemonic(context.Context, *RequestUserID) (*CreateWalletMnemonicResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMnemonic not implemented")
}
func (UnimplementedSwapRobotServiceServer) GetMnemonicFromMemory(context.Context, *RequestUserID) (*CreateWalletMnemonicResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMnemonicFromMemory not implemented")
}
func (UnimplementedSwapRobotServiceServer) ReloadMnemonic(context.Context, *RequestUserID) (*UserRobotResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReloadMnemonic not implemented")
}
func (UnimplementedSwapRobotServiceServer) HelloWorldTest(context.Context, *Ceshi) (*UserRobotResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HelloWorldTest not implemented")
}
func (UnimplementedSwapRobotServiceServer) ImportWallet(context.Context, *ImportWalletMnemonicReq) (*CreateWalletMnemonicResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportWallet not implemented")
}
func (UnimplementedSwapRobotServiceServer) mustEmbedUnimplementedSwapRobotServiceServer() {}

// UnsafeSwapRobotServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SwapRobotServiceServer will
// result in compilation errors.
type UnsafeSwapRobotServiceServer interface {
	mustEmbedUnimplementedSwapRobotServiceServer()
}

func RegisterSwapRobotServiceServer(s grpc.ServiceRegistrar, srv SwapRobotServiceServer) {
	s.RegisterService(&SwapRobotService_ServiceDesc, srv)
}

func _SwapRobotService_CreateUserRobot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWalletMnemonicReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwapRobotServiceServer).CreateUserRobot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemonicService.SwapRobotService/CreateUserRobot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwapRobotServiceServer).CreateUserRobot(ctx, req.(*CreateWalletMnemonicReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwapRobotService_DeleteUserRobot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwapRobotServiceServer).DeleteUserRobot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemonicService.SwapRobotService/DeleteUserRobot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwapRobotServiceServer).DeleteUserRobot(ctx, req.(*RequestUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwapRobotService_CheckIsExpireTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckIsNeedSign)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwapRobotServiceServer).CheckIsExpireTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemonicService.SwapRobotService/CheckIsExpireTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwapRobotServiceServer).CheckIsExpireTime(ctx, req.(*CheckIsNeedSign))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwapRobotService_GetMnemonic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwapRobotServiceServer).GetMnemonic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemonicService.SwapRobotService/GetMnemonic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwapRobotServiceServer).GetMnemonic(ctx, req.(*RequestUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwapRobotService_GetMnemonicFromMemory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwapRobotServiceServer).GetMnemonicFromMemory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemonicService.SwapRobotService/GetMnemonicFromMemory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwapRobotServiceServer).GetMnemonicFromMemory(ctx, req.(*RequestUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwapRobotService_ReloadMnemonic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwapRobotServiceServer).ReloadMnemonic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemonicService.SwapRobotService/ReloadMnemonic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwapRobotServiceServer).ReloadMnemonic(ctx, req.(*RequestUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwapRobotService_HelloWorldTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ceshi)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwapRobotServiceServer).HelloWorldTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemonicService.SwapRobotService/HelloWorldTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwapRobotServiceServer).HelloWorldTest(ctx, req.(*Ceshi))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwapRobotService_ImportWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportWalletMnemonicReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwapRobotServiceServer).ImportWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mnemonicService.SwapRobotService/ImportWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwapRobotServiceServer).ImportWallet(ctx, req.(*ImportWalletMnemonicReq))
	}
	return interceptor(ctx, in, info, handler)
}

// SwapRobotService_ServiceDesc is the grpc.ServiceDesc for SwapRobotService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SwapRobotService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mnemonicService.SwapRobotService",
	HandlerType: (*SwapRobotServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUserRobot",
			Handler:    _SwapRobotService_CreateUserRobot_Handler,
		},
		{
			MethodName: "DeleteUserRobot",
			Handler:    _SwapRobotService_DeleteUserRobot_Handler,
		},
		{
			MethodName: "CheckIsExpireTime",
			Handler:    _SwapRobotService_CheckIsExpireTime_Handler,
		},
		{
			MethodName: "GetMnemonic",
			Handler:    _SwapRobotService_GetMnemonic_Handler,
		},
		{
			MethodName: "GetMnemonicFromMemory",
			Handler:    _SwapRobotService_GetMnemonicFromMemory_Handler,
		},
		{
			MethodName: "ReloadMnemonic",
			Handler:    _SwapRobotService_ReloadMnemonic_Handler,
		},
		{
			MethodName: "HelloWorldTest",
			Handler:    _SwapRobotService_HelloWorldTest_Handler,
		},
		{
			MethodName: "ImportWallet",
			Handler:    _SwapRobotService_ImportWallet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protoc/mnemonic.proto",
}