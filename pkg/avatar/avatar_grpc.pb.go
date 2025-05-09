// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.26.1
// source: api/avatar.proto

package avatar

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AvatarService_SetUserAvatar_FullMethodName        = "/AvatarService/SetUserAvatar"
	AvatarService_GetAllUserAvatars_FullMethodName    = "/AvatarService/GetAllUserAvatars"
	AvatarService_DeleteUserAvatar_FullMethodName     = "/AvatarService/DeleteUserAvatar"
	AvatarService_SetSocietyAvatar_FullMethodName     = "/AvatarService/SetSocietyAvatar"
	AvatarService_GetAllSocietyAvatars_FullMethodName = "/AvatarService/GetAllSocietyAvatars"
	AvatarService_DeleteSocietyAvatar_FullMethodName  = "/AvatarService/DeleteSocietyAvatar"
)

// AvatarServiceClient is the client API for AvatarService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AvatarServiceClient interface {
	SetUserAvatar(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[SetUserAvatarIn, SetUserAvatarOut], error)
	GetAllUserAvatars(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAllUserAvatarsOut, error)
	DeleteUserAvatar(ctx context.Context, in *DeleteUserAvatarIn, opts ...grpc.CallOption) (*Avatar, error)
	SetSocietyAvatar(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[SetSocietyAvatarIn, SetSocietyAvatarOut], error)
	GetAllSocietyAvatars(ctx context.Context, in *GetAllSocietyAvatarsIn, opts ...grpc.CallOption) (*GetAllSocietyAvatarsOut, error)
	DeleteSocietyAvatar(ctx context.Context, in *DeleteSocietyAvatarIn, opts ...grpc.CallOption) (*Avatar, error)
}

type avatarServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAvatarServiceClient(cc grpc.ClientConnInterface) AvatarServiceClient {
	return &avatarServiceClient{cc}
}

func (c *avatarServiceClient) SetUserAvatar(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[SetUserAvatarIn, SetUserAvatarOut], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &AvatarService_ServiceDesc.Streams[0], AvatarService_SetUserAvatar_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SetUserAvatarIn, SetUserAvatarOut]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type AvatarService_SetUserAvatarClient = grpc.ClientStreamingClient[SetUserAvatarIn, SetUserAvatarOut]

func (c *avatarServiceClient) GetAllUserAvatars(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAllUserAvatarsOut, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllUserAvatarsOut)
	err := c.cc.Invoke(ctx, AvatarService_GetAllUserAvatars_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *avatarServiceClient) DeleteUserAvatar(ctx context.Context, in *DeleteUserAvatarIn, opts ...grpc.CallOption) (*Avatar, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Avatar)
	err := c.cc.Invoke(ctx, AvatarService_DeleteUserAvatar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *avatarServiceClient) SetSocietyAvatar(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[SetSocietyAvatarIn, SetSocietyAvatarOut], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &AvatarService_ServiceDesc.Streams[1], AvatarService_SetSocietyAvatar_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SetSocietyAvatarIn, SetSocietyAvatarOut]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type AvatarService_SetSocietyAvatarClient = grpc.ClientStreamingClient[SetSocietyAvatarIn, SetSocietyAvatarOut]

func (c *avatarServiceClient) GetAllSocietyAvatars(ctx context.Context, in *GetAllSocietyAvatarsIn, opts ...grpc.CallOption) (*GetAllSocietyAvatarsOut, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllSocietyAvatarsOut)
	err := c.cc.Invoke(ctx, AvatarService_GetAllSocietyAvatars_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *avatarServiceClient) DeleteSocietyAvatar(ctx context.Context, in *DeleteSocietyAvatarIn, opts ...grpc.CallOption) (*Avatar, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Avatar)
	err := c.cc.Invoke(ctx, AvatarService_DeleteSocietyAvatar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AvatarServiceServer is the server API for AvatarService service.
// All implementations must embed UnimplementedAvatarServiceServer
// for forward compatibility.
type AvatarServiceServer interface {
	SetUserAvatar(grpc.ClientStreamingServer[SetUserAvatarIn, SetUserAvatarOut]) error
	GetAllUserAvatars(context.Context, *emptypb.Empty) (*GetAllUserAvatarsOut, error)
	DeleteUserAvatar(context.Context, *DeleteUserAvatarIn) (*Avatar, error)
	SetSocietyAvatar(grpc.ClientStreamingServer[SetSocietyAvatarIn, SetSocietyAvatarOut]) error
	GetAllSocietyAvatars(context.Context, *GetAllSocietyAvatarsIn) (*GetAllSocietyAvatarsOut, error)
	DeleteSocietyAvatar(context.Context, *DeleteSocietyAvatarIn) (*Avatar, error)
	mustEmbedUnimplementedAvatarServiceServer()
}

// UnimplementedAvatarServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAvatarServiceServer struct{}

func (UnimplementedAvatarServiceServer) SetUserAvatar(grpc.ClientStreamingServer[SetUserAvatarIn, SetUserAvatarOut]) error {
	return status.Errorf(codes.Unimplemented, "method SetUserAvatar not implemented")
}
func (UnimplementedAvatarServiceServer) GetAllUserAvatars(context.Context, *emptypb.Empty) (*GetAllUserAvatarsOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllUserAvatars not implemented")
}
func (UnimplementedAvatarServiceServer) DeleteUserAvatar(context.Context, *DeleteUserAvatarIn) (*Avatar, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserAvatar not implemented")
}
func (UnimplementedAvatarServiceServer) SetSocietyAvatar(grpc.ClientStreamingServer[SetSocietyAvatarIn, SetSocietyAvatarOut]) error {
	return status.Errorf(codes.Unimplemented, "method SetSocietyAvatar not implemented")
}
func (UnimplementedAvatarServiceServer) GetAllSocietyAvatars(context.Context, *GetAllSocietyAvatarsIn) (*GetAllSocietyAvatarsOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllSocietyAvatars not implemented")
}
func (UnimplementedAvatarServiceServer) DeleteSocietyAvatar(context.Context, *DeleteSocietyAvatarIn) (*Avatar, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSocietyAvatar not implemented")
}
func (UnimplementedAvatarServiceServer) mustEmbedUnimplementedAvatarServiceServer() {}
func (UnimplementedAvatarServiceServer) testEmbeddedByValue()                       {}

// UnsafeAvatarServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AvatarServiceServer will
// result in compilation errors.
type UnsafeAvatarServiceServer interface {
	mustEmbedUnimplementedAvatarServiceServer()
}

func RegisterAvatarServiceServer(s grpc.ServiceRegistrar, srv AvatarServiceServer) {
	// If the following call pancis, it indicates UnimplementedAvatarServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AvatarService_ServiceDesc, srv)
}

func _AvatarService_SetUserAvatar_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AvatarServiceServer).SetUserAvatar(&grpc.GenericServerStream[SetUserAvatarIn, SetUserAvatarOut]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type AvatarService_SetUserAvatarServer = grpc.ClientStreamingServer[SetUserAvatarIn, SetUserAvatarOut]

func _AvatarService_GetAllUserAvatars_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AvatarServiceServer).GetAllUserAvatars(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AvatarService_GetAllUserAvatars_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AvatarServiceServer).GetAllUserAvatars(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AvatarService_DeleteUserAvatar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserAvatarIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AvatarServiceServer).DeleteUserAvatar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AvatarService_DeleteUserAvatar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AvatarServiceServer).DeleteUserAvatar(ctx, req.(*DeleteUserAvatarIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _AvatarService_SetSocietyAvatar_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AvatarServiceServer).SetSocietyAvatar(&grpc.GenericServerStream[SetSocietyAvatarIn, SetSocietyAvatarOut]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type AvatarService_SetSocietyAvatarServer = grpc.ClientStreamingServer[SetSocietyAvatarIn, SetSocietyAvatarOut]

func _AvatarService_GetAllSocietyAvatars_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllSocietyAvatarsIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AvatarServiceServer).GetAllSocietyAvatars(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AvatarService_GetAllSocietyAvatars_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AvatarServiceServer).GetAllSocietyAvatars(ctx, req.(*GetAllSocietyAvatarsIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _AvatarService_DeleteSocietyAvatar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSocietyAvatarIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AvatarServiceServer).DeleteSocietyAvatar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AvatarService_DeleteSocietyAvatar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AvatarServiceServer).DeleteSocietyAvatar(ctx, req.(*DeleteSocietyAvatarIn))
	}
	return interceptor(ctx, in, info, handler)
}

// AvatarService_ServiceDesc is the grpc.ServiceDesc for AvatarService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AvatarService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AvatarService",
	HandlerType: (*AvatarServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllUserAvatars",
			Handler:    _AvatarService_GetAllUserAvatars_Handler,
		},
		{
			MethodName: "DeleteUserAvatar",
			Handler:    _AvatarService_DeleteUserAvatar_Handler,
		},
		{
			MethodName: "GetAllSocietyAvatars",
			Handler:    _AvatarService_GetAllSocietyAvatars_Handler,
		},
		{
			MethodName: "DeleteSocietyAvatar",
			Handler:    _AvatarService_DeleteSocietyAvatar_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SetUserAvatar",
			Handler:       _AvatarService_SetUserAvatar_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "SetSocietyAvatar",
			Handler:       _AvatarService_SetSocietyAvatar_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "api/avatar.proto",
}
