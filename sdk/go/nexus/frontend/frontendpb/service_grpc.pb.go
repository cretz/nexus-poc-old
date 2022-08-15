// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: nexus/frontend/v1/service.proto

package frontendpb

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

// FrontendManagementServiceClient is the client API for FrontendManagementService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FrontendManagementServiceClient interface {
	UpdateConfig(ctx context.Context, in *UpdateConfigRequest, opts ...grpc.CallOption) (*UpdateConfigResponse, error)
}

type frontendManagementServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFrontendManagementServiceClient(cc grpc.ClientConnInterface) FrontendManagementServiceClient {
	return &frontendManagementServiceClient{cc}
}

func (c *frontendManagementServiceClient) UpdateConfig(ctx context.Context, in *UpdateConfigRequest, opts ...grpc.CallOption) (*UpdateConfigResponse, error) {
	out := new(UpdateConfigResponse)
	err := c.cc.Invoke(ctx, "/nexus.frontend.v1.FrontendManagementService/UpdateConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FrontendManagementServiceServer is the server API for FrontendManagementService service.
// All implementations must embed UnimplementedFrontendManagementServiceServer
// for forward compatibility
type FrontendManagementServiceServer interface {
	UpdateConfig(context.Context, *UpdateConfigRequest) (*UpdateConfigResponse, error)
	mustEmbedUnimplementedFrontendManagementServiceServer()
}

// UnimplementedFrontendManagementServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFrontendManagementServiceServer struct {
}

func (UnimplementedFrontendManagementServiceServer) UpdateConfig(context.Context, *UpdateConfigRequest) (*UpdateConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateConfig not implemented")
}
func (UnimplementedFrontendManagementServiceServer) mustEmbedUnimplementedFrontendManagementServiceServer() {
}

// UnsafeFrontendManagementServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FrontendManagementServiceServer will
// result in compilation errors.
type UnsafeFrontendManagementServiceServer interface {
	mustEmbedUnimplementedFrontendManagementServiceServer()
}

func RegisterFrontendManagementServiceServer(s grpc.ServiceRegistrar, srv FrontendManagementServiceServer) {
	s.RegisterService(&FrontendManagementService_ServiceDesc, srv)
}

func _FrontendManagementService_UpdateConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontendManagementServiceServer).UpdateConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nexus.frontend.v1.FrontendManagementService/UpdateConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontendManagementServiceServer).UpdateConfig(ctx, req.(*UpdateConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FrontendManagementService_ServiceDesc is the grpc.ServiceDesc for FrontendManagementService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FrontendManagementService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nexus.frontend.v1.FrontendManagementService",
	HandlerType: (*FrontendManagementServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateConfig",
			Handler:    _FrontendManagementService_UpdateConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nexus/frontend/v1/service.proto",
}