// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: office/office.proto

package office

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

// OfficeServiceClient is the client API for OfficeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OfficeServiceClient interface {
	GetUserTags(ctx context.Context, in *GetUserTagsReq, opts ...grpc.CallOption) (*GetUserTagsResp, error)
	CreateTag(ctx context.Context, in *CreateTagReq, opts ...grpc.CallOption) (*CreateTagResp, error)
	DeleteTag(ctx context.Context, in *DeleteTagReq, opts ...grpc.CallOption) (*DeleteTagResp, error)
	SetTag(ctx context.Context, in *SetTagReq, opts ...grpc.CallOption) (*SetTagResp, error)
	SendMsg2Tag(ctx context.Context, in *SendMsg2TagReq, opts ...grpc.CallOption) (*SendMsg2TagResp, error)
	GetTagSendLogs(ctx context.Context, in *GetTagSendLogsReq, opts ...grpc.CallOption) (*GetTagSendLogsResp, error)
	GetUserTagByID(ctx context.Context, in *GetUserTagByIDReq, opts ...grpc.CallOption) (*GetUserTagByIDResp, error)
	CreateOneWorkMoment(ctx context.Context, in *CreateOneWorkMomentReq, opts ...grpc.CallOption) (*CreateOneWorkMomentResp, error)
	DeleteOneWorkMoment(ctx context.Context, in *DeleteOneWorkMomentReq, opts ...grpc.CallOption) (*DeleteOneWorkMomentResp, error)
	LikeOneWorkMoment(ctx context.Context, in *LikeOneWorkMomentReq, opts ...grpc.CallOption) (*LikeOneWorkMomentResp, error)
	CommentOneWorkMoment(ctx context.Context, in *CommentOneWorkMomentReq, opts ...grpc.CallOption) (*CommentOneWorkMomentResp, error)
	DeleteComment(ctx context.Context, in *DeleteCommentReq, opts ...grpc.CallOption) (*DeleteCommentResp, error)
	GetWorkMomentByID(ctx context.Context, in *GetWorkMomentByIDReq, opts ...grpc.CallOption) (*GetWorkMomentByIDResp, error)
	ChangeWorkMomentPermission(ctx context.Context, in *ChangeWorkMomentPermissionReq, opts ...grpc.CallOption) (*ChangeWorkMomentPermissionResp, error)
	// / user self
	GetUserWorkMoments(ctx context.Context, in *GetUserWorkMomentsReq, opts ...grpc.CallOption) (*GetUserWorkMomentsResp, error)
	// / users friend
	GetUserFriendWorkMoments(ctx context.Context, in *GetUserFriendWorkMomentsReq, opts ...grpc.CallOption) (*GetUserFriendWorkMomentsResp, error)
	SetUserWorkMomentsLevel(ctx context.Context, in *SetUserWorkMomentsLevelReq, opts ...grpc.CallOption) (*SetUserWorkMomentsLevelResp, error)
}

type officeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOfficeServiceClient(cc grpc.ClientConnInterface) OfficeServiceClient {
	return &officeServiceClient{cc}
}

func (c *officeServiceClient) GetUserTags(ctx context.Context, in *GetUserTagsReq, opts ...grpc.CallOption) (*GetUserTagsResp, error) {
	out := new(GetUserTagsResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/GetUserTags", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) CreateTag(ctx context.Context, in *CreateTagReq, opts ...grpc.CallOption) (*CreateTagResp, error) {
	out := new(CreateTagResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/CreateTag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) DeleteTag(ctx context.Context, in *DeleteTagReq, opts ...grpc.CallOption) (*DeleteTagResp, error) {
	out := new(DeleteTagResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/DeleteTag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) SetTag(ctx context.Context, in *SetTagReq, opts ...grpc.CallOption) (*SetTagResp, error) {
	out := new(SetTagResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/SetTag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) SendMsg2Tag(ctx context.Context, in *SendMsg2TagReq, opts ...grpc.CallOption) (*SendMsg2TagResp, error) {
	out := new(SendMsg2TagResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/SendMsg2Tag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) GetTagSendLogs(ctx context.Context, in *GetTagSendLogsReq, opts ...grpc.CallOption) (*GetTagSendLogsResp, error) {
	out := new(GetTagSendLogsResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/GetTagSendLogs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) GetUserTagByID(ctx context.Context, in *GetUserTagByIDReq, opts ...grpc.CallOption) (*GetUserTagByIDResp, error) {
	out := new(GetUserTagByIDResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/GetUserTagByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) CreateOneWorkMoment(ctx context.Context, in *CreateOneWorkMomentReq, opts ...grpc.CallOption) (*CreateOneWorkMomentResp, error) {
	out := new(CreateOneWorkMomentResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/CreateOneWorkMoment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) DeleteOneWorkMoment(ctx context.Context, in *DeleteOneWorkMomentReq, opts ...grpc.CallOption) (*DeleteOneWorkMomentResp, error) {
	out := new(DeleteOneWorkMomentResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/DeleteOneWorkMoment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) LikeOneWorkMoment(ctx context.Context, in *LikeOneWorkMomentReq, opts ...grpc.CallOption) (*LikeOneWorkMomentResp, error) {
	out := new(LikeOneWorkMomentResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/LikeOneWorkMoment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) CommentOneWorkMoment(ctx context.Context, in *CommentOneWorkMomentReq, opts ...grpc.CallOption) (*CommentOneWorkMomentResp, error) {
	out := new(CommentOneWorkMomentResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/CommentOneWorkMoment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) DeleteComment(ctx context.Context, in *DeleteCommentReq, opts ...grpc.CallOption) (*DeleteCommentResp, error) {
	out := new(DeleteCommentResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/DeleteComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) GetWorkMomentByID(ctx context.Context, in *GetWorkMomentByIDReq, opts ...grpc.CallOption) (*GetWorkMomentByIDResp, error) {
	out := new(GetWorkMomentByIDResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/GetWorkMomentByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) ChangeWorkMomentPermission(ctx context.Context, in *ChangeWorkMomentPermissionReq, opts ...grpc.CallOption) (*ChangeWorkMomentPermissionResp, error) {
	out := new(ChangeWorkMomentPermissionResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/ChangeWorkMomentPermission", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) GetUserWorkMoments(ctx context.Context, in *GetUserWorkMomentsReq, opts ...grpc.CallOption) (*GetUserWorkMomentsResp, error) {
	out := new(GetUserWorkMomentsResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/GetUserWorkMoments", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) GetUserFriendWorkMoments(ctx context.Context, in *GetUserFriendWorkMomentsReq, opts ...grpc.CallOption) (*GetUserFriendWorkMomentsResp, error) {
	out := new(GetUserFriendWorkMomentsResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/GetUserFriendWorkMoments", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *officeServiceClient) SetUserWorkMomentsLevel(ctx context.Context, in *SetUserWorkMomentsLevelReq, opts ...grpc.CallOption) (*SetUserWorkMomentsLevelResp, error) {
	out := new(SetUserWorkMomentsLevelResp)
	err := c.cc.Invoke(ctx, "/office.OfficeService/SetUserWorkMomentsLevel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OfficeServiceServer is the server API for OfficeService service.
// All implementations should embed UnimplementedOfficeServiceServer
// for forward compatibility
type OfficeServiceServer interface {
	GetUserTags(context.Context, *GetUserTagsReq) (*GetUserTagsResp, error)
	CreateTag(context.Context, *CreateTagReq) (*CreateTagResp, error)
	DeleteTag(context.Context, *DeleteTagReq) (*DeleteTagResp, error)
	SetTag(context.Context, *SetTagReq) (*SetTagResp, error)
	SendMsg2Tag(context.Context, *SendMsg2TagReq) (*SendMsg2TagResp, error)
	GetTagSendLogs(context.Context, *GetTagSendLogsReq) (*GetTagSendLogsResp, error)
	GetUserTagByID(context.Context, *GetUserTagByIDReq) (*GetUserTagByIDResp, error)
	CreateOneWorkMoment(context.Context, *CreateOneWorkMomentReq) (*CreateOneWorkMomentResp, error)
	DeleteOneWorkMoment(context.Context, *DeleteOneWorkMomentReq) (*DeleteOneWorkMomentResp, error)
	LikeOneWorkMoment(context.Context, *LikeOneWorkMomentReq) (*LikeOneWorkMomentResp, error)
	CommentOneWorkMoment(context.Context, *CommentOneWorkMomentReq) (*CommentOneWorkMomentResp, error)
	DeleteComment(context.Context, *DeleteCommentReq) (*DeleteCommentResp, error)
	GetWorkMomentByID(context.Context, *GetWorkMomentByIDReq) (*GetWorkMomentByIDResp, error)
	ChangeWorkMomentPermission(context.Context, *ChangeWorkMomentPermissionReq) (*ChangeWorkMomentPermissionResp, error)
	// / user self
	GetUserWorkMoments(context.Context, *GetUserWorkMomentsReq) (*GetUserWorkMomentsResp, error)
	// / users friend
	GetUserFriendWorkMoments(context.Context, *GetUserFriendWorkMomentsReq) (*GetUserFriendWorkMomentsResp, error)
	SetUserWorkMomentsLevel(context.Context, *SetUserWorkMomentsLevelReq) (*SetUserWorkMomentsLevelResp, error)
}

// UnimplementedOfficeServiceServer should be embedded to have forward compatible implementations.
type UnimplementedOfficeServiceServer struct {
}

func (UnimplementedOfficeServiceServer) GetUserTags(context.Context, *GetUserTagsReq) (*GetUserTagsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserTags not implemented")
}
func (UnimplementedOfficeServiceServer) CreateTag(context.Context, *CreateTagReq) (*CreateTagResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTag not implemented")
}
func (UnimplementedOfficeServiceServer) DeleteTag(context.Context, *DeleteTagReq) (*DeleteTagResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTag not implemented")
}
func (UnimplementedOfficeServiceServer) SetTag(context.Context, *SetTagReq) (*SetTagResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetTag not implemented")
}
func (UnimplementedOfficeServiceServer) SendMsg2Tag(context.Context, *SendMsg2TagReq) (*SendMsg2TagResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsg2Tag not implemented")
}
func (UnimplementedOfficeServiceServer) GetTagSendLogs(context.Context, *GetTagSendLogsReq) (*GetTagSendLogsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTagSendLogs not implemented")
}
func (UnimplementedOfficeServiceServer) GetUserTagByID(context.Context, *GetUserTagByIDReq) (*GetUserTagByIDResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserTagByID not implemented")
}
func (UnimplementedOfficeServiceServer) CreateOneWorkMoment(context.Context, *CreateOneWorkMomentReq) (*CreateOneWorkMomentResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOneWorkMoment not implemented")
}
func (UnimplementedOfficeServiceServer) DeleteOneWorkMoment(context.Context, *DeleteOneWorkMomentReq) (*DeleteOneWorkMomentResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteOneWorkMoment not implemented")
}
func (UnimplementedOfficeServiceServer) LikeOneWorkMoment(context.Context, *LikeOneWorkMomentReq) (*LikeOneWorkMomentResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikeOneWorkMoment not implemented")
}
func (UnimplementedOfficeServiceServer) CommentOneWorkMoment(context.Context, *CommentOneWorkMomentReq) (*CommentOneWorkMomentResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommentOneWorkMoment not implemented")
}
func (UnimplementedOfficeServiceServer) DeleteComment(context.Context, *DeleteCommentReq) (*DeleteCommentResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}
func (UnimplementedOfficeServiceServer) GetWorkMomentByID(context.Context, *GetWorkMomentByIDReq) (*GetWorkMomentByIDResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWorkMomentByID not implemented")
}
func (UnimplementedOfficeServiceServer) ChangeWorkMomentPermission(context.Context, *ChangeWorkMomentPermissionReq) (*ChangeWorkMomentPermissionResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeWorkMomentPermission not implemented")
}
func (UnimplementedOfficeServiceServer) GetUserWorkMoments(context.Context, *GetUserWorkMomentsReq) (*GetUserWorkMomentsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserWorkMoments not implemented")
}
func (UnimplementedOfficeServiceServer) GetUserFriendWorkMoments(context.Context, *GetUserFriendWorkMomentsReq) (*GetUserFriendWorkMomentsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserFriendWorkMoments not implemented")
}
func (UnimplementedOfficeServiceServer) SetUserWorkMomentsLevel(context.Context, *SetUserWorkMomentsLevelReq) (*SetUserWorkMomentsLevelResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUserWorkMomentsLevel not implemented")
}

// UnsafeOfficeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OfficeServiceServer will
// result in compilation errors.
type UnsafeOfficeServiceServer interface {
	mustEmbedUnimplementedOfficeServiceServer()
}

func RegisterOfficeServiceServer(s grpc.ServiceRegistrar, srv OfficeServiceServer) {
	s.RegisterService(&OfficeService_ServiceDesc, srv)
}

func _OfficeService_GetUserTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserTagsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).GetUserTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/GetUserTags",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).GetUserTags(ctx, req.(*GetUserTagsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_CreateTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTagReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).CreateTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/CreateTag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).CreateTag(ctx, req.(*CreateTagReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_DeleteTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTagReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).DeleteTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/DeleteTag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).DeleteTag(ctx, req.(*DeleteTagReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_SetTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetTagReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).SetTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/SetTag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).SetTag(ctx, req.(*SetTagReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_SendMsg2Tag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsg2TagReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).SendMsg2Tag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/SendMsg2Tag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).SendMsg2Tag(ctx, req.(*SendMsg2TagReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_GetTagSendLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTagSendLogsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).GetTagSendLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/GetTagSendLogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).GetTagSendLogs(ctx, req.(*GetTagSendLogsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_GetUserTagByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserTagByIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).GetUserTagByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/GetUserTagByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).GetUserTagByID(ctx, req.(*GetUserTagByIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_CreateOneWorkMoment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOneWorkMomentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).CreateOneWorkMoment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/CreateOneWorkMoment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).CreateOneWorkMoment(ctx, req.(*CreateOneWorkMomentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_DeleteOneWorkMoment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteOneWorkMomentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).DeleteOneWorkMoment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/DeleteOneWorkMoment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).DeleteOneWorkMoment(ctx, req.(*DeleteOneWorkMomentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_LikeOneWorkMoment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikeOneWorkMomentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).LikeOneWorkMoment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/LikeOneWorkMoment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).LikeOneWorkMoment(ctx, req.(*LikeOneWorkMomentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_CommentOneWorkMoment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommentOneWorkMomentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).CommentOneWorkMoment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/CommentOneWorkMoment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).CommentOneWorkMoment(ctx, req.(*CommentOneWorkMomentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_DeleteComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCommentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).DeleteComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/DeleteComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).DeleteComment(ctx, req.(*DeleteCommentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_GetWorkMomentByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWorkMomentByIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).GetWorkMomentByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/GetWorkMomentByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).GetWorkMomentByID(ctx, req.(*GetWorkMomentByIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_ChangeWorkMomentPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeWorkMomentPermissionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).ChangeWorkMomentPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/ChangeWorkMomentPermission",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).ChangeWorkMomentPermission(ctx, req.(*ChangeWorkMomentPermissionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_GetUserWorkMoments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserWorkMomentsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).GetUserWorkMoments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/GetUserWorkMoments",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).GetUserWorkMoments(ctx, req.(*GetUserWorkMomentsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_GetUserFriendWorkMoments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserFriendWorkMomentsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).GetUserFriendWorkMoments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/GetUserFriendWorkMoments",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).GetUserFriendWorkMoments(ctx, req.(*GetUserFriendWorkMomentsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OfficeService_SetUserWorkMomentsLevel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetUserWorkMomentsLevelReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfficeServiceServer).SetUserWorkMomentsLevel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/office.OfficeService/SetUserWorkMomentsLevel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfficeServiceServer).SetUserWorkMomentsLevel(ctx, req.(*SetUserWorkMomentsLevelReq))
	}
	return interceptor(ctx, in, info, handler)
}

// OfficeService_ServiceDesc is the grpc.ServiceDesc for OfficeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OfficeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "office.OfficeService",
	HandlerType: (*OfficeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserTags",
			Handler:    _OfficeService_GetUserTags_Handler,
		},
		{
			MethodName: "CreateTag",
			Handler:    _OfficeService_CreateTag_Handler,
		},
		{
			MethodName: "DeleteTag",
			Handler:    _OfficeService_DeleteTag_Handler,
		},
		{
			MethodName: "SetTag",
			Handler:    _OfficeService_SetTag_Handler,
		},
		{
			MethodName: "SendMsg2Tag",
			Handler:    _OfficeService_SendMsg2Tag_Handler,
		},
		{
			MethodName: "GetTagSendLogs",
			Handler:    _OfficeService_GetTagSendLogs_Handler,
		},
		{
			MethodName: "GetUserTagByID",
			Handler:    _OfficeService_GetUserTagByID_Handler,
		},
		{
			MethodName: "CreateOneWorkMoment",
			Handler:    _OfficeService_CreateOneWorkMoment_Handler,
		},
		{
			MethodName: "DeleteOneWorkMoment",
			Handler:    _OfficeService_DeleteOneWorkMoment_Handler,
		},
		{
			MethodName: "LikeOneWorkMoment",
			Handler:    _OfficeService_LikeOneWorkMoment_Handler,
		},
		{
			MethodName: "CommentOneWorkMoment",
			Handler:    _OfficeService_CommentOneWorkMoment_Handler,
		},
		{
			MethodName: "DeleteComment",
			Handler:    _OfficeService_DeleteComment_Handler,
		},
		{
			MethodName: "GetWorkMomentByID",
			Handler:    _OfficeService_GetWorkMomentByID_Handler,
		},
		{
			MethodName: "ChangeWorkMomentPermission",
			Handler:    _OfficeService_ChangeWorkMomentPermission_Handler,
		},
		{
			MethodName: "GetUserWorkMoments",
			Handler:    _OfficeService_GetUserWorkMoments_Handler,
		},
		{
			MethodName: "GetUserFriendWorkMoments",
			Handler:    _OfficeService_GetUserFriendWorkMoments_Handler,
		},
		{
			MethodName: "SetUserWorkMomentsLevel",
			Handler:    _OfficeService_SetUserWorkMomentsLevel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "office/office.proto",
}
