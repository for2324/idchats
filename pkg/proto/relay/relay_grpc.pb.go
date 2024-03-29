// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: relay/relay.proto

package pbRelay

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

// RelayClient is the client API for Relay service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelayClient interface {
	OnlinePushMsg(ctx context.Context, in *OnlinePushMsgReq, opts ...grpc.CallOption) (*OnlinePushMsgResp, error)
	GetUsersOnlineStatus(ctx context.Context, in *GetUsersOnlineStatusReq, opts ...grpc.CallOption) (*GetUsersOnlineStatusResp, error)
	OnlineBatchPushOneMsg(ctx context.Context, in *OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*OnlineBatchPushOneMsgResp, error)
	OnlineBatchAnnouncementPushOneMsg(ctx context.Context, in *OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*OnlineBatchPushOneMsgResp, error)
	SuperGroupOnlineBatchPushOneMsg(ctx context.Context, in *OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*OnlineBatchPushOneMsgResp, error)
	KickUserOffline(ctx context.Context, in *KickUserOfflineReq, opts ...grpc.CallOption) (*KickUserOfflineResp, error)
	MultiTerminalLoginCheck(ctx context.Context, in *MultiTerminalLoginCheckReq, opts ...grpc.CallOption) (*MultiTerminalLoginCheckResp, error)
}

type relayClient struct {
	cc grpc.ClientConnInterface
}

func NewRelayClient(cc grpc.ClientConnInterface) RelayClient {
	return &relayClient{cc}
}

func (c *relayClient) OnlinePushMsg(ctx context.Context, in *OnlinePushMsgReq, opts ...grpc.CallOption) (*OnlinePushMsgResp, error) {
	out := new(OnlinePushMsgResp)
	err := c.cc.Invoke(ctx, "/relay.relay/OnlinePushMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayClient) GetUsersOnlineStatus(ctx context.Context, in *GetUsersOnlineStatusReq, opts ...grpc.CallOption) (*GetUsersOnlineStatusResp, error) {
	out := new(GetUsersOnlineStatusResp)
	err := c.cc.Invoke(ctx, "/relay.relay/GetUsersOnlineStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayClient) OnlineBatchPushOneMsg(ctx context.Context, in *OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*OnlineBatchPushOneMsgResp, error) {
	out := new(OnlineBatchPushOneMsgResp)
	err := c.cc.Invoke(ctx, "/relay.relay/OnlineBatchPushOneMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayClient) OnlineBatchAnnouncementPushOneMsg(ctx context.Context, in *OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*OnlineBatchPushOneMsgResp, error) {
	out := new(OnlineBatchPushOneMsgResp)
	err := c.cc.Invoke(ctx, "/relay.relay/OnlineBatchAnnouncementPushOneMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayClient) SuperGroupOnlineBatchPushOneMsg(ctx context.Context, in *OnlineBatchPushOneMsgReq, opts ...grpc.CallOption) (*OnlineBatchPushOneMsgResp, error) {
	out := new(OnlineBatchPushOneMsgResp)
	err := c.cc.Invoke(ctx, "/relay.relay/SuperGroupOnlineBatchPushOneMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayClient) KickUserOffline(ctx context.Context, in *KickUserOfflineReq, opts ...grpc.CallOption) (*KickUserOfflineResp, error) {
	out := new(KickUserOfflineResp)
	err := c.cc.Invoke(ctx, "/relay.relay/KickUserOffline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayClient) MultiTerminalLoginCheck(ctx context.Context, in *MultiTerminalLoginCheckReq, opts ...grpc.CallOption) (*MultiTerminalLoginCheckResp, error) {
	out := new(MultiTerminalLoginCheckResp)
	err := c.cc.Invoke(ctx, "/relay.relay/MultiTerminalLoginCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelayServer is the server API for Relay service.
// All implementations should embed UnimplementedRelayServer
// for forward compatibility
type RelayServer interface {
	OnlinePushMsg(context.Context, *OnlinePushMsgReq) (*OnlinePushMsgResp, error)
	GetUsersOnlineStatus(context.Context, *GetUsersOnlineStatusReq) (*GetUsersOnlineStatusResp, error)
	OnlineBatchPushOneMsg(context.Context, *OnlineBatchPushOneMsgReq) (*OnlineBatchPushOneMsgResp, error)
	OnlineBatchAnnouncementPushOneMsg(context.Context, *OnlineBatchPushOneMsgReq) (*OnlineBatchPushOneMsgResp, error)
	SuperGroupOnlineBatchPushOneMsg(context.Context, *OnlineBatchPushOneMsgReq) (*OnlineBatchPushOneMsgResp, error)
	KickUserOffline(context.Context, *KickUserOfflineReq) (*KickUserOfflineResp, error)
	MultiTerminalLoginCheck(context.Context, *MultiTerminalLoginCheckReq) (*MultiTerminalLoginCheckResp, error)
}

// UnimplementedRelayServer should be embedded to have forward compatible implementations.
type UnimplementedRelayServer struct {
}

func (UnimplementedRelayServer) OnlinePushMsg(context.Context, *OnlinePushMsgReq) (*OnlinePushMsgResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnlinePushMsg not implemented")
}
func (UnimplementedRelayServer) GetUsersOnlineStatus(context.Context, *GetUsersOnlineStatusReq) (*GetUsersOnlineStatusResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsersOnlineStatus not implemented")
}
func (UnimplementedRelayServer) OnlineBatchPushOneMsg(context.Context, *OnlineBatchPushOneMsgReq) (*OnlineBatchPushOneMsgResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnlineBatchPushOneMsg not implemented")
}
func (UnimplementedRelayServer) OnlineBatchAnnouncementPushOneMsg(context.Context, *OnlineBatchPushOneMsgReq) (*OnlineBatchPushOneMsgResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnlineBatchAnnouncementPushOneMsg not implemented")
}
func (UnimplementedRelayServer) SuperGroupOnlineBatchPushOneMsg(context.Context, *OnlineBatchPushOneMsgReq) (*OnlineBatchPushOneMsgResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuperGroupOnlineBatchPushOneMsg not implemented")
}
func (UnimplementedRelayServer) KickUserOffline(context.Context, *KickUserOfflineReq) (*KickUserOfflineResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KickUserOffline not implemented")
}
func (UnimplementedRelayServer) MultiTerminalLoginCheck(context.Context, *MultiTerminalLoginCheckReq) (*MultiTerminalLoginCheckResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiTerminalLoginCheck not implemented")
}

// UnsafeRelayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelayServer will
// result in compilation errors.
type UnsafeRelayServer interface {
	mustEmbedUnimplementedRelayServer()
}

func RegisterRelayServer(s grpc.ServiceRegistrar, srv RelayServer) {
	s.RegisterService(&Relay_ServiceDesc, srv)
}

func _Relay_OnlinePushMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnlinePushMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServer).OnlinePushMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relay.relay/OnlinePushMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServer).OnlinePushMsg(ctx, req.(*OnlinePushMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relay_GetUsersOnlineStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsersOnlineStatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServer).GetUsersOnlineStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relay.relay/GetUsersOnlineStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServer).GetUsersOnlineStatus(ctx, req.(*GetUsersOnlineStatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relay_OnlineBatchPushOneMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnlineBatchPushOneMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServer).OnlineBatchPushOneMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relay.relay/OnlineBatchPushOneMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServer).OnlineBatchPushOneMsg(ctx, req.(*OnlineBatchPushOneMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relay_OnlineBatchAnnouncementPushOneMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnlineBatchPushOneMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServer).OnlineBatchAnnouncementPushOneMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relay.relay/OnlineBatchAnnouncementPushOneMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServer).OnlineBatchAnnouncementPushOneMsg(ctx, req.(*OnlineBatchPushOneMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relay_SuperGroupOnlineBatchPushOneMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnlineBatchPushOneMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServer).SuperGroupOnlineBatchPushOneMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relay.relay/SuperGroupOnlineBatchPushOneMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServer).SuperGroupOnlineBatchPushOneMsg(ctx, req.(*OnlineBatchPushOneMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relay_KickUserOffline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KickUserOfflineReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServer).KickUserOffline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relay.relay/KickUserOffline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServer).KickUserOffline(ctx, req.(*KickUserOfflineReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relay_MultiTerminalLoginCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MultiTerminalLoginCheckReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServer).MultiTerminalLoginCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relay.relay/MultiTerminalLoginCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServer).MultiTerminalLoginCheck(ctx, req.(*MultiTerminalLoginCheckReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Relay_ServiceDesc is the grpc.ServiceDesc for Relay service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Relay_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "relay.relay",
	HandlerType: (*RelayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OnlinePushMsg",
			Handler:    _Relay_OnlinePushMsg_Handler,
		},
		{
			MethodName: "GetUsersOnlineStatus",
			Handler:    _Relay_GetUsersOnlineStatus_Handler,
		},
		{
			MethodName: "OnlineBatchPushOneMsg",
			Handler:    _Relay_OnlineBatchPushOneMsg_Handler,
		},
		{
			MethodName: "OnlineBatchAnnouncementPushOneMsg",
			Handler:    _Relay_OnlineBatchAnnouncementPushOneMsg_Handler,
		},
		{
			MethodName: "SuperGroupOnlineBatchPushOneMsg",
			Handler:    _Relay_SuperGroupOnlineBatchPushOneMsg_Handler,
		},
		{
			MethodName: "KickUserOffline",
			Handler:    _Relay_KickUserOffline_Handler,
		},
		{
			MethodName: "MultiTerminalLoginCheck",
			Handler:    _Relay_MultiTerminalLoginCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "relay/relay.proto",
}
