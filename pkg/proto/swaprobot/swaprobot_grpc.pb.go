// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: swaprobot/swaprobot.proto

package swaprobot

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

// SwaprobotClient is the client API for Swaprobot service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SwaprobotClient interface {
	// 第三方通讯信息
	BotOperation(ctx context.Context, in *BotOperationReq, opts ...grpc.CallOption) (*BotOperationResp, error)
	FinishTaskToGetReword(ctx context.Context, in *BotSwapTradeReq, opts ...grpc.CallOption) (*BotSwapTradeResp, error)
}

type swaprobotClient struct {
	cc grpc.ClientConnInterface
}

func NewSwaprobotClient(cc grpc.ClientConnInterface) SwaprobotClient {
	return &swaprobotClient{cc}
}

func (c *swaprobotClient) BotOperation(ctx context.Context, in *BotOperationReq, opts ...grpc.CallOption) (*BotOperationResp, error) {
	out := new(BotOperationResp)
	err := c.cc.Invoke(ctx, "/swaprobot.swaprobot/BotOperation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaprobotClient) FinishTaskToGetReword(ctx context.Context, in *BotSwapTradeReq, opts ...grpc.CallOption) (*BotSwapTradeResp, error) {
	out := new(BotSwapTradeResp)
	err := c.cc.Invoke(ctx, "/swaprobot.swaprobot/FinishTaskToGetReword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SwaprobotServer is the server API for Swaprobot service.
// All implementations should embed UnimplementedSwaprobotServer
// for forward compatibility
type SwaprobotServer interface {
	// 第三方通讯信息
	BotOperation(context.Context, *BotOperationReq) (*BotOperationResp, error)
	FinishTaskToGetReword(context.Context, *BotSwapTradeReq) (*BotSwapTradeResp, error)
}

// UnimplementedSwaprobotServer should be embedded to have forward compatible implementations.
type UnimplementedSwaprobotServer struct {
}

func (UnimplementedSwaprobotServer) BotOperation(context.Context, *BotOperationReq) (*BotOperationResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BotOperation not implemented")
}
func (UnimplementedSwaprobotServer) FinishTaskToGetReword(context.Context, *BotSwapTradeReq) (*BotSwapTradeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FinishTaskToGetReword not implemented")
}

// UnsafeSwaprobotServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SwaprobotServer will
// result in compilation errors.
type UnsafeSwaprobotServer interface {
	mustEmbedUnimplementedSwaprobotServer()
}

func RegisterSwaprobotServer(s grpc.ServiceRegistrar, srv SwaprobotServer) {
	s.RegisterService(&Swaprobot_ServiceDesc, srv)
}

func _Swaprobot_BotOperation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BotOperationReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwaprobotServer).BotOperation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swaprobot.swaprobot/BotOperation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwaprobotServer).BotOperation(ctx, req.(*BotOperationReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Swaprobot_FinishTaskToGetReword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BotSwapTradeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwaprobotServer).FinishTaskToGetReword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swaprobot.swaprobot/FinishTaskToGetReword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwaprobotServer).FinishTaskToGetReword(ctx, req.(*BotSwapTradeReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Swaprobot_ServiceDesc is the grpc.ServiceDesc for Swaprobot service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Swaprobot_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "swaprobot.swaprobot",
	HandlerType: (*SwaprobotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BotOperation",
			Handler:    _Swaprobot_BotOperation_Handler,
		},
		{
			MethodName: "FinishTaskToGetReword",
			Handler:    _Swaprobot_FinishTaskToGetReword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "swaprobot/swaprobot.proto",
}