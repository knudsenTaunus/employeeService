// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: proto/user/v1/user_service.proto

package userpb

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

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	Register(ctx context.Context, in *UserServiceRegisterRequest, opts ...grpc.CallOption) (UserService_RegisterClient, error)
	Deregister(ctx context.Context, in *UserServiceDeregisterRequest, opts ...grpc.CallOption) (*UserServiceDeregisterResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Register(ctx context.Context, in *UserServiceRegisterRequest, opts ...grpc.CallOption) (UserService_RegisterClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserService_ServiceDesc.Streams[0], "/proto.user.v1.UserService/Register", opts...)
	if err != nil {
		return nil, err
	}
	x := &userServiceRegisterClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UserService_RegisterClient interface {
	Recv() (*UserServiceRegisterResponse, error)
	grpc.ClientStream
}

type userServiceRegisterClient struct {
	grpc.ClientStream
}

func (x *userServiceRegisterClient) Recv() (*UserServiceRegisterResponse, error) {
	m := new(UserServiceRegisterResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userServiceClient) Deregister(ctx context.Context, in *UserServiceDeregisterRequest, opts ...grpc.CallOption) (*UserServiceDeregisterResponse, error) {
	out := new(UserServiceDeregisterResponse)
	err := c.cc.Invoke(ctx, "/proto.user.v1.UserService/Deregister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	Register(*UserServiceRegisterRequest, UserService_RegisterServer) error
	Deregister(context.Context, *UserServiceDeregisterRequest) (*UserServiceDeregisterResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) Register(*UserServiceRegisterRequest, UserService_RegisterServer) error {
	return status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedUserServiceServer) Deregister(context.Context, *UserServiceDeregisterRequest) (*UserServiceDeregisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deregister not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_Register_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(UserServiceRegisterRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServiceServer).Register(m, &userServiceRegisterServer{stream})
}

type UserService_RegisterServer interface {
	Send(*UserServiceRegisterResponse) error
	grpc.ServerStream
}

type userServiceRegisterServer struct {
	grpc.ServerStream
}

func (x *userServiceRegisterServer) Send(m *UserServiceRegisterResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _UserService_Deregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserServiceDeregisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Deregister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.user.v1.UserService/Deregister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Deregister(ctx, req.(*UserServiceDeregisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.user.v1.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deregister",
			Handler:    _UserService_Deregister_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Register",
			Handler:       _UserService_Register_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/user/v1/user_service.proto",
}