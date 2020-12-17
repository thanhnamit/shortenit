// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// AliasProviderServiceClient is the client API for AliasProviderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AliasProviderServiceClient interface {
	GetNewAlias(ctx context.Context, in *GetNewAliasRequest, opts ...grpc.CallOption) (*GetNewAliasResponse, error)
	CheckAliasValidity(ctx context.Context, in *CheckAliasValidityRequest, opts ...grpc.CallOption) (*CheckAliasValidityResponse, error)
}

type aliasProviderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAliasProviderServiceClient(cc grpc.ClientConnInterface) AliasProviderServiceClient {
	return &aliasProviderServiceClient{cc}
}

func (c *aliasProviderServiceClient) GetNewAlias(ctx context.Context, in *GetNewAliasRequest, opts ...grpc.CallOption) (*GetNewAliasResponse, error) {
	out := new(GetNewAliasResponse)
	err := c.cc.Invoke(ctx, "/v1.AliasProviderService/GetNewAlias", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aliasProviderServiceClient) CheckAliasValidity(ctx context.Context, in *CheckAliasValidityRequest, opts ...grpc.CallOption) (*CheckAliasValidityResponse, error) {
	out := new(CheckAliasValidityResponse)
	err := c.cc.Invoke(ctx, "/v1.AliasProviderService/CheckAliasValidity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AliasProviderServiceServer is the server API for AliasProviderService service.
// All implementations must embed UnimplementedAliasProviderServiceServer
// for forward compatibility
type AliasProviderServiceServer interface {
	GetNewAlias(context.Context, *GetNewAliasRequest) (*GetNewAliasResponse, error)
	CheckAliasValidity(context.Context, *CheckAliasValidityRequest) (*CheckAliasValidityResponse, error)
	mustEmbedUnimplementedAliasProviderServiceServer()
}

// UnimplementedAliasProviderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAliasProviderServiceServer struct {
}

func (UnimplementedAliasProviderServiceServer) GetNewAlias(context.Context, *GetNewAliasRequest) (*GetNewAliasResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNewAlias not implemented")
}
func (UnimplementedAliasProviderServiceServer) CheckAliasValidity(context.Context, *CheckAliasValidityRequest) (*CheckAliasValidityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAliasValidity not implemented")
}
func (UnimplementedAliasProviderServiceServer) mustEmbedUnimplementedAliasProviderServiceServer() {}

// UnsafeAliasProviderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AliasProviderServiceServer will
// result in compilation errors.
type UnsafeAliasProviderServiceServer interface {
	mustEmbedUnimplementedAliasProviderServiceServer()
}

func RegisterAliasProviderServiceServer(s grpc.ServiceRegistrar, srv AliasProviderServiceServer) {
	s.RegisterService(&_AliasProviderService_serviceDesc, srv)
}

func _AliasProviderService_GetNewAlias_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNewAliasRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AliasProviderServiceServer).GetNewAlias(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.AliasProviderService/GetNewAlias",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AliasProviderServiceServer).GetNewAlias(ctx, req.(*GetNewAliasRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AliasProviderService_CheckAliasValidity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckAliasValidityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AliasProviderServiceServer).CheckAliasValidity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.AliasProviderService/CheckAliasValidity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AliasProviderServiceServer).CheckAliasValidity(ctx, req.(*CheckAliasValidityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AliasProviderService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.AliasProviderService",
	HandlerType: (*AliasProviderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNewAlias",
			Handler:    _AliasProviderService_GetNewAlias_Handler,
		},
		{
			MethodName: "CheckAliasValidity",
			Handler:    _AliasProviderService_CheckAliasValidity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/alias/v1/alias.proto",
}