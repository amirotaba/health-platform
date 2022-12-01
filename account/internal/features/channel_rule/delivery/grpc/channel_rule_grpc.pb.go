// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.1
// source: channel_rule.proto

package grpc

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

// ChannelRuleServiceClient is the client API for ChannelRuleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChannelRuleServiceClient interface {
	GetChannelTags(ctx context.Context, in *ChannelsRuleRequest, opts ...grpc.CallOption) (*ChannelRuleReply, error)
	GetChannelTag(ctx context.Context, in *ChannelRuleRequest, opts ...grpc.CallOption) (*ChannelRule, error)
}

type channelRuleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChannelRuleServiceClient(cc grpc.ClientConnInterface) ChannelRuleServiceClient {
	return &channelRuleServiceClient{cc}
}

func (c *channelRuleServiceClient) GetChannelTags(ctx context.Context, in *ChannelsRuleRequest, opts ...grpc.CallOption) (*ChannelRuleReply, error) {
	out := new(ChannelRuleReply)
	err := c.cc.Invoke(ctx, "/ChannelRuleService/GetChannelTags", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelRuleServiceClient) GetChannelTag(ctx context.Context, in *ChannelRuleRequest, opts ...grpc.CallOption) (*ChannelRule, error) {
	out := new(ChannelRule)
	err := c.cc.Invoke(ctx, "/ChannelRuleService/GetChannelTag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChannelRuleServiceServer is the server API for ChannelRuleService service.
// All implementations must embed UnimplementedChannelRuleServiceServer
// for forward compatibility
type ChannelRuleServiceServer interface {
	GetChannelTags(context.Context, *ChannelsRuleRequest) (*ChannelRuleReply, error)
	GetChannelTag(context.Context, *ChannelRuleRequest) (*ChannelRule, error)
	mustEmbedUnimplementedChannelRuleServiceServer()
}

// UnimplementedChannelRuleServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChannelRuleServiceServer struct {
}

func (UnimplementedChannelRuleServiceServer) GetChannelTags(context.Context, *ChannelsRuleRequest) (*ChannelRuleReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChannelTags not implemented")
}
func (UnimplementedChannelRuleServiceServer) GetChannelTag(context.Context, *ChannelRuleRequest) (*ChannelRule, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChannelTag not implemented")
}
func (UnimplementedChannelRuleServiceServer) mustEmbedUnimplementedChannelRuleServiceServer() {}

// UnsafeChannelRuleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChannelRuleServiceServer will
// result in compilation errors.
type UnsafeChannelRuleServiceServer interface {
	mustEmbedUnimplementedChannelRuleServiceServer()
}

func RegisterChannelRuleServiceServer(s grpc.ServiceRegistrar, srv ChannelRuleServiceServer) {
	s.RegisterService(&ChannelRuleService_ServiceDesc, srv)
}

func _ChannelRuleService_GetChannelTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChannelsRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelRuleServiceServer).GetChannelTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ChannelRuleService/GetChannelTags",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelRuleServiceServer).GetChannelTags(ctx, req.(*ChannelsRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelRuleService_GetChannelTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChannelRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelRuleServiceServer).GetChannelTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ChannelRuleService/GetChannelTag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelRuleServiceServer).GetChannelTag(ctx, req.(*ChannelRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChannelRuleService_ServiceDesc is the grpc.ServiceDesc for ChannelRuleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChannelRuleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ChannelRuleService",
	HandlerType: (*ChannelRuleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetChannelTags",
			Handler:    _ChannelRuleService_GetChannelTags_Handler,
		},
		{
			MethodName: "GetChannelTag",
			Handler:    _ChannelRuleService_GetChannelTag_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "channel_rule.proto",
}