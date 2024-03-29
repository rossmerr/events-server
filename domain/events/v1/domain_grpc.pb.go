// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: events/v1/domain.proto

package eventsv1

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

const (
	DomainEventsService_Stream_FullMethodName          = "/events.v1.DomainEventsService/Stream"
	DomainEventsService_Publish_FullMethodName         = "/events.v1.DomainEventsService/Publish"
	DomainEventsService_Subscribe_FullMethodName       = "/events.v1.DomainEventsService/Subscribe"
	DomainEventsService_Acknowledgement_FullMethodName = "/events.v1.DomainEventsService/Acknowledgement"
)

// DomainEventsServiceClient is the client API for DomainEventsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DomainEventsServiceClient interface {
	Stream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (*StreamResponse, error)
	Publish(ctx context.Context, opts ...grpc.CallOption) (DomainEventsService_PublishClient, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (DomainEventsService_SubscribeClient, error)
	Acknowledgement(ctx context.Context, opts ...grpc.CallOption) (DomainEventsService_AcknowledgementClient, error)
}

type domainEventsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDomainEventsServiceClient(cc grpc.ClientConnInterface) DomainEventsServiceClient {
	return &domainEventsServiceClient{cc}
}

func (c *domainEventsServiceClient) Stream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (*StreamResponse, error) {
	out := new(StreamResponse)
	err := c.cc.Invoke(ctx, DomainEventsService_Stream_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *domainEventsServiceClient) Publish(ctx context.Context, opts ...grpc.CallOption) (DomainEventsService_PublishClient, error) {
	stream, err := c.cc.NewStream(ctx, &DomainEventsService_ServiceDesc.Streams[0], DomainEventsService_Publish_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &domainEventsServicePublishClient{stream}
	return x, nil
}

type DomainEventsService_PublishClient interface {
	Send(*PublishRequest) error
	CloseAndRecv() (*PublishResponse, error)
	grpc.ClientStream
}

type domainEventsServicePublishClient struct {
	grpc.ClientStream
}

func (x *domainEventsServicePublishClient) Send(m *PublishRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *domainEventsServicePublishClient) CloseAndRecv() (*PublishResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(PublishResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *domainEventsServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (DomainEventsService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &DomainEventsService_ServiceDesc.Streams[1], DomainEventsService_Subscribe_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &domainEventsServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DomainEventsService_SubscribeClient interface {
	Recv() (*SubscribeResponse, error)
	grpc.ClientStream
}

type domainEventsServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *domainEventsServiceSubscribeClient) Recv() (*SubscribeResponse, error) {
	m := new(SubscribeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *domainEventsServiceClient) Acknowledgement(ctx context.Context, opts ...grpc.CallOption) (DomainEventsService_AcknowledgementClient, error) {
	stream, err := c.cc.NewStream(ctx, &DomainEventsService_ServiceDesc.Streams[2], DomainEventsService_Acknowledgement_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &domainEventsServiceAcknowledgementClient{stream}
	return x, nil
}

type DomainEventsService_AcknowledgementClient interface {
	Send(*AcknowledgementRequest) error
	Recv() (*AcknowledgementResponse, error)
	grpc.ClientStream
}

type domainEventsServiceAcknowledgementClient struct {
	grpc.ClientStream
}

func (x *domainEventsServiceAcknowledgementClient) Send(m *AcknowledgementRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *domainEventsServiceAcknowledgementClient) Recv() (*AcknowledgementResponse, error) {
	m := new(AcknowledgementResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DomainEventsServiceServer is the server API for DomainEventsService service.
// All implementations must embed UnimplementedDomainEventsServiceServer
// for forward compatibility
type DomainEventsServiceServer interface {
	Stream(context.Context, *StreamRequest) (*StreamResponse, error)
	Publish(DomainEventsService_PublishServer) error
	Subscribe(*SubscribeRequest, DomainEventsService_SubscribeServer) error
	Acknowledgement(DomainEventsService_AcknowledgementServer) error
	mustEmbedUnimplementedDomainEventsServiceServer()
}

// UnimplementedDomainEventsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDomainEventsServiceServer struct {
}

func (UnimplementedDomainEventsServiceServer) Stream(context.Context, *StreamRequest) (*StreamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedDomainEventsServiceServer) Publish(DomainEventsService_PublishServer) error {
	return status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedDomainEventsServiceServer) Subscribe(*SubscribeRequest, DomainEventsService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedDomainEventsServiceServer) Acknowledgement(DomainEventsService_AcknowledgementServer) error {
	return status.Errorf(codes.Unimplemented, "method Acknowledgement not implemented")
}
func (UnimplementedDomainEventsServiceServer) mustEmbedUnimplementedDomainEventsServiceServer() {}

// UnsafeDomainEventsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DomainEventsServiceServer will
// result in compilation errors.
type UnsafeDomainEventsServiceServer interface {
	mustEmbedUnimplementedDomainEventsServiceServer()
}

func RegisterDomainEventsServiceServer(s grpc.ServiceRegistrar, srv DomainEventsServiceServer) {
	s.RegisterService(&DomainEventsService_ServiceDesc, srv)
}

func _DomainEventsService_Stream_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StreamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DomainEventsServiceServer).Stream(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DomainEventsService_Stream_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DomainEventsServiceServer).Stream(ctx, req.(*StreamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DomainEventsService_Publish_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DomainEventsServiceServer).Publish(&domainEventsServicePublishServer{stream})
}

type DomainEventsService_PublishServer interface {
	SendAndClose(*PublishResponse) error
	Recv() (*PublishRequest, error)
	grpc.ServerStream
}

type domainEventsServicePublishServer struct {
	grpc.ServerStream
}

func (x *domainEventsServicePublishServer) SendAndClose(m *PublishResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *domainEventsServicePublishServer) Recv() (*PublishRequest, error) {
	m := new(PublishRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DomainEventsService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DomainEventsServiceServer).Subscribe(m, &domainEventsServiceSubscribeServer{stream})
}

type DomainEventsService_SubscribeServer interface {
	Send(*SubscribeResponse) error
	grpc.ServerStream
}

type domainEventsServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *domainEventsServiceSubscribeServer) Send(m *SubscribeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _DomainEventsService_Acknowledgement_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DomainEventsServiceServer).Acknowledgement(&domainEventsServiceAcknowledgementServer{stream})
}

type DomainEventsService_AcknowledgementServer interface {
	Send(*AcknowledgementResponse) error
	Recv() (*AcknowledgementRequest, error)
	grpc.ServerStream
}

type domainEventsServiceAcknowledgementServer struct {
	grpc.ServerStream
}

func (x *domainEventsServiceAcknowledgementServer) Send(m *AcknowledgementResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *domainEventsServiceAcknowledgementServer) Recv() (*AcknowledgementRequest, error) {
	m := new(AcknowledgementRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DomainEventsService_ServiceDesc is the grpc.ServiceDesc for DomainEventsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DomainEventsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "events.v1.DomainEventsService",
	HandlerType: (*DomainEventsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Stream",
			Handler:    _DomainEventsService_Stream_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Publish",
			Handler:       _DomainEventsService_Publish_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _DomainEventsService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Acknowledgement",
			Handler:       _DomainEventsService_Acknowledgement_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "events/v1/domain.proto",
}
