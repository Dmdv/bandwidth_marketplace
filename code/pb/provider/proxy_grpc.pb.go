// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package provider

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

// ProxyClient is the client API for Proxy service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProxyClient interface {
	NewSessionBilling(ctx context.Context, in *NewSessionBillingRequest, opts ...grpc.CallOption) (*NewSessionBillingResponse, error)
	ForwardUsage(ctx context.Context, in *ForwardUsageRequest, opts ...grpc.CallOption) (*ForwardUsageResponse, error)
}

type proxyClient struct {
	cc grpc.ClientConnInterface
}

func NewProxyClient(cc grpc.ClientConnInterface) ProxyClient {
	return &proxyClient{cc}
}

func (c *proxyClient) NewSessionBilling(ctx context.Context, in *NewSessionBillingRequest, opts ...grpc.CallOption) (*NewSessionBillingResponse, error) {
	out := new(NewSessionBillingResponse)
	err := c.cc.Invoke(ctx, "/zchain.pb.provider.Proxy/NewSessionBilling", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyClient) ForwardUsage(ctx context.Context, in *ForwardUsageRequest, opts ...grpc.CallOption) (*ForwardUsageResponse, error) {
	out := new(ForwardUsageResponse)
	err := c.cc.Invoke(ctx, "/zchain.pb.provider.Proxy/ForwardUsage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProxyServer is the server API for Proxy service.
// All implementations must embed UnimplementedProxyServer
// for forward compatibility
type ProxyServer interface {
	NewSessionBilling(context.Context, *NewSessionBillingRequest) (*NewSessionBillingResponse, error)
	ForwardUsage(context.Context, *ForwardUsageRequest) (*ForwardUsageResponse, error)
	mustEmbedUnimplementedProxyServer()
}

// UnimplementedProxyServer must be embedded to have forward compatible implementations.
type UnimplementedProxyServer struct {
}

func (UnimplementedProxyServer) NewSessionBilling(context.Context, *NewSessionBillingRequest) (*NewSessionBillingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewSessionBilling not implemented")
}
func (UnimplementedProxyServer) ForwardUsage(context.Context, *ForwardUsageRequest) (*ForwardUsageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForwardUsage not implemented")
}
func (UnimplementedProxyServer) mustEmbedUnimplementedProxyServer() {}

// UnsafeProxyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProxyServer will
// result in compilation errors.
type UnsafeProxyServer interface {
	mustEmbedUnimplementedProxyServer()
}

func RegisterProxyServer(s grpc.ServiceRegistrar, srv ProxyServer) {
	s.RegisterService(&Proxy_ServiceDesc, srv)
}

func _Proxy_NewSessionBilling_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewSessionBillingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).NewSessionBilling(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zchain.pb.provider.Proxy/NewSessionBilling",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).NewSessionBilling(ctx, req.(*NewSessionBillingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Proxy_ForwardUsage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardUsageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).ForwardUsage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zchain.pb.provider.Proxy/ForwardUsage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).ForwardUsage(ctx, req.(*ForwardUsageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Proxy_ServiceDesc is the grpc.ServiceDesc for Proxy service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Proxy_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "zchain.pb.provider.Proxy",
	HandlerType: (*ProxyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewSessionBilling",
			Handler:    _Proxy_NewSessionBilling_Handler,
		},
		{
			MethodName: "ForwardUsage",
			Handler:    _Proxy_ForwardUsage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/provider/proto/proxy.proto",
}
