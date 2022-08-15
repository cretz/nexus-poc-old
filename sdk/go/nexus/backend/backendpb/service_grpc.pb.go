// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: nexus/backend/v1/service.proto

package backendpb

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

// BackendServiceClient is the client API for BackendService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BackendServiceClient interface {
	UpdateServices(ctx context.Context, in *UpdateServicesRequest, opts ...grpc.CallOption) (*UpdateServicesResponse, error)
	GetServices(ctx context.Context, in *GetServicesRequest, opts ...grpc.CallOption) (*GetServicesResponse, error)
	Call(ctx context.Context, in *CallRequest, opts ...grpc.CallOption) (*CallResponse, error)
	// TODO(cretz): Should we do poll here instead?
	StreamTasks(ctx context.Context, opts ...grpc.CallOption) (BackendService_StreamTasksClient, error)
}

type backendServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBackendServiceClient(cc grpc.ClientConnInterface) BackendServiceClient {
	return &backendServiceClient{cc}
}

func (c *backendServiceClient) UpdateServices(ctx context.Context, in *UpdateServicesRequest, opts ...grpc.CallOption) (*UpdateServicesResponse, error) {
	out := new(UpdateServicesResponse)
	err := c.cc.Invoke(ctx, "/nexus.backend.v1.BackendService/UpdateServices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendServiceClient) GetServices(ctx context.Context, in *GetServicesRequest, opts ...grpc.CallOption) (*GetServicesResponse, error) {
	out := new(GetServicesResponse)
	err := c.cc.Invoke(ctx, "/nexus.backend.v1.BackendService/GetServices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendServiceClient) Call(ctx context.Context, in *CallRequest, opts ...grpc.CallOption) (*CallResponse, error) {
	out := new(CallResponse)
	err := c.cc.Invoke(ctx, "/nexus.backend.v1.BackendService/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendServiceClient) StreamTasks(ctx context.Context, opts ...grpc.CallOption) (BackendService_StreamTasksClient, error) {
	stream, err := c.cc.NewStream(ctx, &BackendService_ServiceDesc.Streams[0], "/nexus.backend.v1.BackendService/StreamTasks", opts...)
	if err != nil {
		return nil, err
	}
	x := &backendServiceStreamTasksClient{stream}
	return x, nil
}

type BackendService_StreamTasksClient interface {
	Send(*StreamTasksRequest) error
	Recv() (*StreamTasksResponse, error)
	grpc.ClientStream
}

type backendServiceStreamTasksClient struct {
	grpc.ClientStream
}

func (x *backendServiceStreamTasksClient) Send(m *StreamTasksRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *backendServiceStreamTasksClient) Recv() (*StreamTasksResponse, error) {
	m := new(StreamTasksResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BackendServiceServer is the server API for BackendService service.
// All implementations must embed UnimplementedBackendServiceServer
// for forward compatibility
type BackendServiceServer interface {
	UpdateServices(context.Context, *UpdateServicesRequest) (*UpdateServicesResponse, error)
	GetServices(context.Context, *GetServicesRequest) (*GetServicesResponse, error)
	Call(context.Context, *CallRequest) (*CallResponse, error)
	// TODO(cretz): Should we do poll here instead?
	StreamTasks(BackendService_StreamTasksServer) error
	mustEmbedUnimplementedBackendServiceServer()
}

// UnimplementedBackendServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBackendServiceServer struct {
}

func (UnimplementedBackendServiceServer) UpdateServices(context.Context, *UpdateServicesRequest) (*UpdateServicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateServices not implemented")
}
func (UnimplementedBackendServiceServer) GetServices(context.Context, *GetServicesRequest) (*GetServicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServices not implemented")
}
func (UnimplementedBackendServiceServer) Call(context.Context, *CallRequest) (*CallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func (UnimplementedBackendServiceServer) StreamTasks(BackendService_StreamTasksServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamTasks not implemented")
}
func (UnimplementedBackendServiceServer) mustEmbedUnimplementedBackendServiceServer() {}

// UnsafeBackendServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BackendServiceServer will
// result in compilation errors.
type UnsafeBackendServiceServer interface {
	mustEmbedUnimplementedBackendServiceServer()
}

func RegisterBackendServiceServer(s grpc.ServiceRegistrar, srv BackendServiceServer) {
	s.RegisterService(&BackendService_ServiceDesc, srv)
}

func _BackendService_UpdateServices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateServicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServiceServer).UpdateServices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nexus.backend.v1.BackendService/UpdateServices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServiceServer).UpdateServices(ctx, req.(*UpdateServicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BackendService_GetServices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServiceServer).GetServices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nexus.backend.v1.BackendService/GetServices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServiceServer).GetServices(ctx, req.(*GetServicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BackendService_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServiceServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nexus.backend.v1.BackendService/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServiceServer).Call(ctx, req.(*CallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BackendService_StreamTasks_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BackendServiceServer).StreamTasks(&backendServiceStreamTasksServer{stream})
}

type BackendService_StreamTasksServer interface {
	Send(*StreamTasksResponse) error
	Recv() (*StreamTasksRequest, error)
	grpc.ServerStream
}

type backendServiceStreamTasksServer struct {
	grpc.ServerStream
}

func (x *backendServiceStreamTasksServer) Send(m *StreamTasksResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *backendServiceStreamTasksServer) Recv() (*StreamTasksRequest, error) {
	m := new(StreamTasksRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BackendService_ServiceDesc is the grpc.ServiceDesc for BackendService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BackendService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nexus.backend.v1.BackendService",
	HandlerType: (*BackendServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateServices",
			Handler:    _BackendService_UpdateServices_Handler,
		},
		{
			MethodName: "GetServices",
			Handler:    _BackendService_GetServices_Handler,
		},
		{
			MethodName: "Call",
			Handler:    _BackendService_Call_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamTasks",
			Handler:       _BackendService_StreamTasks_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "nexus/backend/v1/service.proto",
}