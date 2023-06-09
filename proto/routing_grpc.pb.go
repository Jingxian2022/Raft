// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// RouterClient is the client API for Router service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RouterClient interface {
	Lookup(ctx context.Context, in *RouteLookupRequest, opts ...grpc.CallOption) (*RouteLookupReply, error)
	Publish(ctx context.Context, in *RoutePublishRequest, opts ...grpc.CallOption) (*RoutePublishReply, error)
	Unpublish(ctx context.Context, in *RouteUnpublishRequest, opts ...grpc.CallOption) (*RouteUnpublishReply, error)
}

type routerClient struct {
	cc grpc.ClientConnInterface
}

func NewRouterClient(cc grpc.ClientConnInterface) RouterClient {
	return &routerClient{cc}
}

func (c *routerClient) Lookup(ctx context.Context, in *RouteLookupRequest, opts ...grpc.CallOption) (*RouteLookupReply, error) {
	out := new(RouteLookupReply)
	err := c.cc.Invoke(ctx, "/proto.Router/Lookup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routerClient) Publish(ctx context.Context, in *RoutePublishRequest, opts ...grpc.CallOption) (*RoutePublishReply, error) {
	out := new(RoutePublishReply)
	err := c.cc.Invoke(ctx, "/proto.Router/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routerClient) Unpublish(ctx context.Context, in *RouteUnpublishRequest, opts ...grpc.CallOption) (*RouteUnpublishReply, error) {
	out := new(RouteUnpublishReply)
	err := c.cc.Invoke(ctx, "/proto.Router/Unpublish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RouterServer is the server API for Router service.
// All implementations must embed UnimplementedRouterServer
// for forward compatibility
type RouterServer interface {
	Lookup(context.Context, *RouteLookupRequest) (*RouteLookupReply, error)
	Publish(context.Context, *RoutePublishRequest) (*RoutePublishReply, error)
	Unpublish(context.Context, *RouteUnpublishRequest) (*RouteUnpublishReply, error)
	mustEmbedUnimplementedRouterServer()
}

// UnimplementedRouterServer must be embedded to have forward compatible implementations.
type UnimplementedRouterServer struct {
}

func (UnimplementedRouterServer) Lookup(context.Context, *RouteLookupRequest) (*RouteLookupReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lookup not implemented")
}
func (UnimplementedRouterServer) Publish(context.Context, *RoutePublishRequest) (*RoutePublishReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedRouterServer) Unpublish(context.Context, *RouteUnpublishRequest) (*RouteUnpublishReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unpublish not implemented")
}
func (UnimplementedRouterServer) mustEmbedUnimplementedRouterServer() {}

// UnsafeRouterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RouterServer will
// result in compilation errors.
type UnsafeRouterServer interface {
	mustEmbedUnimplementedRouterServer()
}

func RegisterRouterServer(s grpc.ServiceRegistrar, srv RouterServer) {
	s.RegisterService(&Router_ServiceDesc, srv)
}

func _Router_Lookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RouteLookupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouterServer).Lookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Router/Lookup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouterServer).Lookup(ctx, req.(*RouteLookupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Router_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoutePublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouterServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Router/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouterServer).Publish(ctx, req.(*RoutePublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Router_Unpublish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RouteUnpublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouterServer).Unpublish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Router/Unpublish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouterServer).Unpublish(ctx, req.(*RouteUnpublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Router_ServiceDesc is the grpc.ServiceDesc for Router service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Router_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Router",
	HandlerType: (*RouterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Lookup",
			Handler:    _Router_Lookup_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _Router_Publish_Handler,
		},
		{
			MethodName: "Unpublish",
			Handler:    _Router_Unpublish_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/routing.proto",
}
