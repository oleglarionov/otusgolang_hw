// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"

	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventServiceClient interface {
	Create(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*Event, error)
	Update(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*Event, error)
	Delete(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DayList(ctx context.Context, in *DayListRequest, opts ...grpc.CallOption) (*Events, error)
	WeekList(ctx context.Context, in *WeekListRequest, opts ...grpc.CallOption) (*Events, error)
	MonthList(ctx context.Context, in *MonthListRequest, opts ...grpc.CallOption) (*Events, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) Create(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/event.EventService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) Update(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/event.EventService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) Delete(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/event.EventService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) DayList(ctx context.Context, in *DayListRequest, opts ...grpc.CallOption) (*Events, error) {
	out := new(Events)
	err := c.cc.Invoke(ctx, "/event.EventService/DayList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) WeekList(ctx context.Context, in *WeekListRequest, opts ...grpc.CallOption) (*Events, error) {
	out := new(Events)
	err := c.cc.Invoke(ctx, "/event.EventService/WeekList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) MonthList(ctx context.Context, in *MonthListRequest, opts ...grpc.CallOption) (*Events, error) {
	out := new(Events)
	err := c.cc.Invoke(ctx, "/event.EventService/MonthList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceServer is the server API for EventService service.
// All implementations must embed UnimplementedEventServiceServer
// for forward compatibility
type EventServiceServer interface {
	Create(context.Context, *CreateEventRequest) (*Event, error)
	Update(context.Context, *UpdateEventRequest) (*Event, error)
	Delete(context.Context, *DeleteEventRequest) (*empty.Empty, error)
	DayList(context.Context, *DayListRequest) (*Events, error)
	WeekList(context.Context, *WeekListRequest) (*Events, error)
	MonthList(context.Context, *MonthListRequest) (*Events, error)
	mustEmbedUnimplementedEventServiceServer()
}

// UnimplementedEventServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEventServiceServer struct {
}

func (UnimplementedEventServiceServer) Create(context.Context, *CreateEventRequest) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedEventServiceServer) Update(context.Context, *UpdateEventRequest) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedEventServiceServer) Delete(context.Context, *DeleteEventRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedEventServiceServer) DayList(context.Context, *DayListRequest) (*Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DayList not implemented")
}
func (UnimplementedEventServiceServer) WeekList(context.Context, *WeekListRequest) (*Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WeekList not implemented")
}
func (UnimplementedEventServiceServer) MonthList(context.Context, *MonthListRequest) (*Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MonthList not implemented")
}
func (UnimplementedEventServiceServer) mustEmbedUnimplementedEventServiceServer() {}

// UnsafeEventServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServiceServer will
// result in compilation errors.
type UnsafeEventServiceServer interface {
	mustEmbedUnimplementedEventServiceServer()
}

func RegisterEventServiceServer(s grpc.ServiceRegistrar, srv EventServiceServer) {
	s.RegisterService(&_EventService_serviceDesc, srv)
}

func _EventService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).Create(ctx, req.(*CreateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).Update(ctx, req.(*UpdateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).Delete(ctx, req.(*DeleteEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_DayList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DayListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).DayList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/DayList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).DayList(ctx, req.(*DayListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_WeekList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WeekListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).WeekList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/WeekList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).WeekList(ctx, req.(*WeekListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_MonthList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MonthListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).MonthList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/MonthList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).MonthList(ctx, req.(*MonthListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EventService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "event.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _EventService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _EventService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _EventService_Delete_Handler,
		},
		{
			MethodName: "DayList",
			Handler:    _EventService_DayList_Handler,
		},
		{
			MethodName: "WeekList",
			Handler:    _EventService_WeekList_Handler,
		},
		{
			MethodName: "MonthList",
			Handler:    _EventService_MonthList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "EventService.proto",
}
