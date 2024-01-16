/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed BY
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/


// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.10
// source: mlops/scheduler/scheduler.proto

package scheduler

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

// SchedulerClient is the client API for Scheduler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SchedulerClient interface {
	ServerNotify(ctx context.Context, in *ServerNotifyRequest, opts ...grpc.CallOption) (*ServerNotifyResponse, error)
	LoadModel(ctx context.Context, in *LoadModelRequest, opts ...grpc.CallOption) (*LoadModelResponse, error)
	UnloadModel(ctx context.Context, in *UnloadModelRequest, opts ...grpc.CallOption) (*UnloadModelResponse, error)
	LoadPipeline(ctx context.Context, in *LoadPipelineRequest, opts ...grpc.CallOption) (*LoadPipelineResponse, error)
	UnloadPipeline(ctx context.Context, in *UnloadPipelineRequest, opts ...grpc.CallOption) (*UnloadPipelineResponse, error)
	StartExperiment(ctx context.Context, in *StartExperimentRequest, opts ...grpc.CallOption) (*StartExperimentResponse, error)
	StopExperiment(ctx context.Context, in *StopExperimentRequest, opts ...grpc.CallOption) (*StopExperimentResponse, error)
	ServerStatus(ctx context.Context, in *ServerStatusRequest, opts ...grpc.CallOption) (Scheduler_ServerStatusClient, error)
	ModelStatus(ctx context.Context, in *ModelStatusRequest, opts ...grpc.CallOption) (Scheduler_ModelStatusClient, error)
	PipelineStatus(ctx context.Context, in *PipelineStatusRequest, opts ...grpc.CallOption) (Scheduler_PipelineStatusClient, error)
	ExperimentStatus(ctx context.Context, in *ExperimentStatusRequest, opts ...grpc.CallOption) (Scheduler_ExperimentStatusClient, error)
	SchedulerStatus(ctx context.Context, in *SchedulerStatusRequest, opts ...grpc.CallOption) (*SchedulerStatusResponse, error)
	SubscribeServerStatus(ctx context.Context, in *ServerSubscriptionRequest, opts ...grpc.CallOption) (Scheduler_SubscribeServerStatusClient, error)
	SubscribeModelStatus(ctx context.Context, in *ModelSubscriptionRequest, opts ...grpc.CallOption) (Scheduler_SubscribeModelStatusClient, error)
	SubscribeExperimentStatus(ctx context.Context, in *ExperimentSubscriptionRequest, opts ...grpc.CallOption) (Scheduler_SubscribeExperimentStatusClient, error)
	SubscribePipelineStatus(ctx context.Context, in *PipelineSubscriptionRequest, opts ...grpc.CallOption) (Scheduler_SubscribePipelineStatusClient, error)
}

type schedulerClient struct {
	cc grpc.ClientConnInterface
}

func NewSchedulerClient(cc grpc.ClientConnInterface) SchedulerClient {
	return &schedulerClient{cc}
}

func (c *schedulerClient) ServerNotify(ctx context.Context, in *ServerNotifyRequest, opts ...grpc.CallOption) (*ServerNotifyResponse, error) {
	out := new(ServerNotifyResponse)
	err := c.cc.Invoke(ctx, "/seldon.mlops.scheduler.Scheduler/ServerNotify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerClient) LoadModel(ctx context.Context, in *LoadModelRequest, opts ...grpc.CallOption) (*LoadModelResponse, error) {
	out := new(LoadModelResponse)
	err := c.cc.Invoke(ctx, "/seldon.mlops.scheduler.Scheduler/LoadModel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerClient) UnloadModel(ctx context.Context, in *UnloadModelRequest, opts ...grpc.CallOption) (*UnloadModelResponse, error) {
	out := new(UnloadModelResponse)
	err := c.cc.Invoke(ctx, "/seldon.mlops.scheduler.Scheduler/UnloadModel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerClient) LoadPipeline(ctx context.Context, in *LoadPipelineRequest, opts ...grpc.CallOption) (*LoadPipelineResponse, error) {
	out := new(LoadPipelineResponse)
	err := c.cc.Invoke(ctx, "/seldon.mlops.scheduler.Scheduler/LoadPipeline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerClient) UnloadPipeline(ctx context.Context, in *UnloadPipelineRequest, opts ...grpc.CallOption) (*UnloadPipelineResponse, error) {
	out := new(UnloadPipelineResponse)
	err := c.cc.Invoke(ctx, "/seldon.mlops.scheduler.Scheduler/UnloadPipeline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerClient) StartExperiment(ctx context.Context, in *StartExperimentRequest, opts ...grpc.CallOption) (*StartExperimentResponse, error) {
	out := new(StartExperimentResponse)
	err := c.cc.Invoke(ctx, "/seldon.mlops.scheduler.Scheduler/StartExperiment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerClient) StopExperiment(ctx context.Context, in *StopExperimentRequest, opts ...grpc.CallOption) (*StopExperimentResponse, error) {
	out := new(StopExperimentResponse)
	err := c.cc.Invoke(ctx, "/seldon.mlops.scheduler.Scheduler/StopExperiment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerClient) ServerStatus(ctx context.Context, in *ServerStatusRequest, opts ...grpc.CallOption) (Scheduler_ServerStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &Scheduler_ServiceDesc.Streams[0], "/seldon.mlops.scheduler.Scheduler/ServerStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &schedulerServerStatusClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Scheduler_ServerStatusClient interface {
	Recv() (*ServerStatusResponse, error)
	grpc.ClientStream
}

type schedulerServerStatusClient struct {
	grpc.ClientStream
}

func (x *schedulerServerStatusClient) Recv() (*ServerStatusResponse, error) {
	m := new(ServerStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *schedulerClient) ModelStatus(ctx context.Context, in *ModelStatusRequest, opts ...grpc.CallOption) (Scheduler_ModelStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &Scheduler_ServiceDesc.Streams[1], "/seldon.mlops.scheduler.Scheduler/ModelStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &schedulerModelStatusClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Scheduler_ModelStatusClient interface {
	Recv() (*ModelStatusResponse, error)
	grpc.ClientStream
}

type schedulerModelStatusClient struct {
	grpc.ClientStream
}

func (x *schedulerModelStatusClient) Recv() (*ModelStatusResponse, error) {
	m := new(ModelStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *schedulerClient) PipelineStatus(ctx context.Context, in *PipelineStatusRequest, opts ...grpc.CallOption) (Scheduler_PipelineStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &Scheduler_ServiceDesc.Streams[2], "/seldon.mlops.scheduler.Scheduler/PipelineStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &schedulerPipelineStatusClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Scheduler_PipelineStatusClient interface {
	Recv() (*PipelineStatusResponse, error)
	grpc.ClientStream
}

type schedulerPipelineStatusClient struct {
	grpc.ClientStream
}

func (x *schedulerPipelineStatusClient) Recv() (*PipelineStatusResponse, error) {
	m := new(PipelineStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *schedulerClient) ExperimentStatus(ctx context.Context, in *ExperimentStatusRequest, opts ...grpc.CallOption) (Scheduler_ExperimentStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &Scheduler_ServiceDesc.Streams[3], "/seldon.mlops.scheduler.Scheduler/ExperimentStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &schedulerExperimentStatusClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Scheduler_ExperimentStatusClient interface {
	Recv() (*ExperimentStatusResponse, error)
	grpc.ClientStream
}

type schedulerExperimentStatusClient struct {
	grpc.ClientStream
}

func (x *schedulerExperimentStatusClient) Recv() (*ExperimentStatusResponse, error) {
	m := new(ExperimentStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *schedulerClient) SchedulerStatus(ctx context.Context, in *SchedulerStatusRequest, opts ...grpc.CallOption) (*SchedulerStatusResponse, error) {
	out := new(SchedulerStatusResponse)
	err := c.cc.Invoke(ctx, "/seldon.mlops.scheduler.Scheduler/SchedulerStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerClient) SubscribeServerStatus(ctx context.Context, in *ServerSubscriptionRequest, opts ...grpc.CallOption) (Scheduler_SubscribeServerStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &Scheduler_ServiceDesc.Streams[4], "/seldon.mlops.scheduler.Scheduler/SubscribeServerStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &schedulerSubscribeServerStatusClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Scheduler_SubscribeServerStatusClient interface {
	Recv() (*ServerStatusResponse, error)
	grpc.ClientStream
}

type schedulerSubscribeServerStatusClient struct {
	grpc.ClientStream
}

func (x *schedulerSubscribeServerStatusClient) Recv() (*ServerStatusResponse, error) {
	m := new(ServerStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *schedulerClient) SubscribeModelStatus(ctx context.Context, in *ModelSubscriptionRequest, opts ...grpc.CallOption) (Scheduler_SubscribeModelStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &Scheduler_ServiceDesc.Streams[5], "/seldon.mlops.scheduler.Scheduler/SubscribeModelStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &schedulerSubscribeModelStatusClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Scheduler_SubscribeModelStatusClient interface {
	Recv() (*ModelStatusResponse, error)
	grpc.ClientStream
}

type schedulerSubscribeModelStatusClient struct {
	grpc.ClientStream
}

func (x *schedulerSubscribeModelStatusClient) Recv() (*ModelStatusResponse, error) {
	m := new(ModelStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *schedulerClient) SubscribeExperimentStatus(ctx context.Context, in *ExperimentSubscriptionRequest, opts ...grpc.CallOption) (Scheduler_SubscribeExperimentStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &Scheduler_ServiceDesc.Streams[6], "/seldon.mlops.scheduler.Scheduler/SubscribeExperimentStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &schedulerSubscribeExperimentStatusClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Scheduler_SubscribeExperimentStatusClient interface {
	Recv() (*ExperimentStatusResponse, error)
	grpc.ClientStream
}

type schedulerSubscribeExperimentStatusClient struct {
	grpc.ClientStream
}

func (x *schedulerSubscribeExperimentStatusClient) Recv() (*ExperimentStatusResponse, error) {
	m := new(ExperimentStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *schedulerClient) SubscribePipelineStatus(ctx context.Context, in *PipelineSubscriptionRequest, opts ...grpc.CallOption) (Scheduler_SubscribePipelineStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &Scheduler_ServiceDesc.Streams[7], "/seldon.mlops.scheduler.Scheduler/SubscribePipelineStatus", opts...)
	if err != nil {
		return nil, err
	}
	x := &schedulerSubscribePipelineStatusClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Scheduler_SubscribePipelineStatusClient interface {
	Recv() (*PipelineStatusResponse, error)
	grpc.ClientStream
}

type schedulerSubscribePipelineStatusClient struct {
	grpc.ClientStream
}

func (x *schedulerSubscribePipelineStatusClient) Recv() (*PipelineStatusResponse, error) {
	m := new(PipelineStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SchedulerServer is the server API for Scheduler service.
// All implementations must embed UnimplementedSchedulerServer
// for forward compatibility
type SchedulerServer interface {
	ServerNotify(context.Context, *ServerNotifyRequest) (*ServerNotifyResponse, error)
	LoadModel(context.Context, *LoadModelRequest) (*LoadModelResponse, error)
	UnloadModel(context.Context, *UnloadModelRequest) (*UnloadModelResponse, error)
	LoadPipeline(context.Context, *LoadPipelineRequest) (*LoadPipelineResponse, error)
	UnloadPipeline(context.Context, *UnloadPipelineRequest) (*UnloadPipelineResponse, error)
	StartExperiment(context.Context, *StartExperimentRequest) (*StartExperimentResponse, error)
	StopExperiment(context.Context, *StopExperimentRequest) (*StopExperimentResponse, error)
	ServerStatus(*ServerStatusRequest, Scheduler_ServerStatusServer) error
	ModelStatus(*ModelStatusRequest, Scheduler_ModelStatusServer) error
	PipelineStatus(*PipelineStatusRequest, Scheduler_PipelineStatusServer) error
	ExperimentStatus(*ExperimentStatusRequest, Scheduler_ExperimentStatusServer) error
	SchedulerStatus(context.Context, *SchedulerStatusRequest) (*SchedulerStatusResponse, error)
	SubscribeServerStatus(*ServerSubscriptionRequest, Scheduler_SubscribeServerStatusServer) error
	SubscribeModelStatus(*ModelSubscriptionRequest, Scheduler_SubscribeModelStatusServer) error
	SubscribeExperimentStatus(*ExperimentSubscriptionRequest, Scheduler_SubscribeExperimentStatusServer) error
	SubscribePipelineStatus(*PipelineSubscriptionRequest, Scheduler_SubscribePipelineStatusServer) error
	mustEmbedUnimplementedSchedulerServer()
}

// UnimplementedSchedulerServer must be embedded to have forward compatible implementations.
type UnimplementedSchedulerServer struct {
}

func (UnimplementedSchedulerServer) ServerNotify(context.Context, *ServerNotifyRequest) (*ServerNotifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServerNotify not implemented")
}
func (UnimplementedSchedulerServer) LoadModel(context.Context, *LoadModelRequest) (*LoadModelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadModel not implemented")
}
func (UnimplementedSchedulerServer) UnloadModel(context.Context, *UnloadModelRequest) (*UnloadModelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnloadModel not implemented")
}
func (UnimplementedSchedulerServer) LoadPipeline(context.Context, *LoadPipelineRequest) (*LoadPipelineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadPipeline not implemented")
}
func (UnimplementedSchedulerServer) UnloadPipeline(context.Context, *UnloadPipelineRequest) (*UnloadPipelineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnloadPipeline not implemented")
}
func (UnimplementedSchedulerServer) StartExperiment(context.Context, *StartExperimentRequest) (*StartExperimentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartExperiment not implemented")
}
func (UnimplementedSchedulerServer) StopExperiment(context.Context, *StopExperimentRequest) (*StopExperimentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopExperiment not implemented")
}
func (UnimplementedSchedulerServer) ServerStatus(*ServerStatusRequest, Scheduler_ServerStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method ServerStatus not implemented")
}
func (UnimplementedSchedulerServer) ModelStatus(*ModelStatusRequest, Scheduler_ModelStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method ModelStatus not implemented")
}
func (UnimplementedSchedulerServer) PipelineStatus(*PipelineStatusRequest, Scheduler_PipelineStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method PipelineStatus not implemented")
}
func (UnimplementedSchedulerServer) ExperimentStatus(*ExperimentStatusRequest, Scheduler_ExperimentStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method ExperimentStatus not implemented")
}
func (UnimplementedSchedulerServer) SchedulerStatus(context.Context, *SchedulerStatusRequest) (*SchedulerStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SchedulerStatus not implemented")
}
func (UnimplementedSchedulerServer) SubscribeServerStatus(*ServerSubscriptionRequest, Scheduler_SubscribeServerStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeServerStatus not implemented")
}
func (UnimplementedSchedulerServer) SubscribeModelStatus(*ModelSubscriptionRequest, Scheduler_SubscribeModelStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeModelStatus not implemented")
}
func (UnimplementedSchedulerServer) SubscribeExperimentStatus(*ExperimentSubscriptionRequest, Scheduler_SubscribeExperimentStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeExperimentStatus not implemented")
}
func (UnimplementedSchedulerServer) SubscribePipelineStatus(*PipelineSubscriptionRequest, Scheduler_SubscribePipelineStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribePipelineStatus not implemented")
}
func (UnimplementedSchedulerServer) mustEmbedUnimplementedSchedulerServer() {}

// UnsafeSchedulerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SchedulerServer will
// result in compilation errors.
type UnsafeSchedulerServer interface {
	mustEmbedUnimplementedSchedulerServer()
}

func RegisterSchedulerServer(s grpc.ServiceRegistrar, srv SchedulerServer) {
	s.RegisterService(&Scheduler_ServiceDesc, srv)
}

func _Scheduler_ServerNotify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerNotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServer).ServerNotify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seldon.mlops.scheduler.Scheduler/ServerNotify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServer).ServerNotify(ctx, req.(*ServerNotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduler_LoadModel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadModelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServer).LoadModel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seldon.mlops.scheduler.Scheduler/LoadModel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServer).LoadModel(ctx, req.(*LoadModelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduler_UnloadModel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnloadModelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServer).UnloadModel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seldon.mlops.scheduler.Scheduler/UnloadModel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServer).UnloadModel(ctx, req.(*UnloadModelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduler_LoadPipeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadPipelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServer).LoadPipeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seldon.mlops.scheduler.Scheduler/LoadPipeline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServer).LoadPipeline(ctx, req.(*LoadPipelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduler_UnloadPipeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnloadPipelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServer).UnloadPipeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seldon.mlops.scheduler.Scheduler/UnloadPipeline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServer).UnloadPipeline(ctx, req.(*UnloadPipelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduler_StartExperiment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartExperimentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServer).StartExperiment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seldon.mlops.scheduler.Scheduler/StartExperiment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServer).StartExperiment(ctx, req.(*StartExperimentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduler_StopExperiment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopExperimentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServer).StopExperiment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seldon.mlops.scheduler.Scheduler/StopExperiment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServer).StopExperiment(ctx, req.(*StopExperimentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduler_ServerStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ServerStatusRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SchedulerServer).ServerStatus(m, &schedulerServerStatusServer{stream})
}

type Scheduler_ServerStatusServer interface {
	Send(*ServerStatusResponse) error
	grpc.ServerStream
}

type schedulerServerStatusServer struct {
	grpc.ServerStream
}

func (x *schedulerServerStatusServer) Send(m *ServerStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Scheduler_ModelStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ModelStatusRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SchedulerServer).ModelStatus(m, &schedulerModelStatusServer{stream})
}

type Scheduler_ModelStatusServer interface {
	Send(*ModelStatusResponse) error
	grpc.ServerStream
}

type schedulerModelStatusServer struct {
	grpc.ServerStream
}

func (x *schedulerModelStatusServer) Send(m *ModelStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Scheduler_PipelineStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PipelineStatusRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SchedulerServer).PipelineStatus(m, &schedulerPipelineStatusServer{stream})
}

type Scheduler_PipelineStatusServer interface {
	Send(*PipelineStatusResponse) error
	grpc.ServerStream
}

type schedulerPipelineStatusServer struct {
	grpc.ServerStream
}

func (x *schedulerPipelineStatusServer) Send(m *PipelineStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Scheduler_ExperimentStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExperimentStatusRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SchedulerServer).ExperimentStatus(m, &schedulerExperimentStatusServer{stream})
}

type Scheduler_ExperimentStatusServer interface {
	Send(*ExperimentStatusResponse) error
	grpc.ServerStream
}

type schedulerExperimentStatusServer struct {
	grpc.ServerStream
}

func (x *schedulerExperimentStatusServer) Send(m *ExperimentStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Scheduler_SchedulerStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SchedulerStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulerServer).SchedulerStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seldon.mlops.scheduler.Scheduler/SchedulerStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulerServer).SchedulerStatus(ctx, req.(*SchedulerStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduler_SubscribeServerStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ServerSubscriptionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SchedulerServer).SubscribeServerStatus(m, &schedulerSubscribeServerStatusServer{stream})
}

type Scheduler_SubscribeServerStatusServer interface {
	Send(*ServerStatusResponse) error
	grpc.ServerStream
}

type schedulerSubscribeServerStatusServer struct {
	grpc.ServerStream
}

func (x *schedulerSubscribeServerStatusServer) Send(m *ServerStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Scheduler_SubscribeModelStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ModelSubscriptionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SchedulerServer).SubscribeModelStatus(m, &schedulerSubscribeModelStatusServer{stream})
}

type Scheduler_SubscribeModelStatusServer interface {
	Send(*ModelStatusResponse) error
	grpc.ServerStream
}

type schedulerSubscribeModelStatusServer struct {
	grpc.ServerStream
}

func (x *schedulerSubscribeModelStatusServer) Send(m *ModelStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Scheduler_SubscribeExperimentStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExperimentSubscriptionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SchedulerServer).SubscribeExperimentStatus(m, &schedulerSubscribeExperimentStatusServer{stream})
}

type Scheduler_SubscribeExperimentStatusServer interface {
	Send(*ExperimentStatusResponse) error
	grpc.ServerStream
}

type schedulerSubscribeExperimentStatusServer struct {
	grpc.ServerStream
}

func (x *schedulerSubscribeExperimentStatusServer) Send(m *ExperimentStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Scheduler_SubscribePipelineStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PipelineSubscriptionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SchedulerServer).SubscribePipelineStatus(m, &schedulerSubscribePipelineStatusServer{stream})
}

type Scheduler_SubscribePipelineStatusServer interface {
	Send(*PipelineStatusResponse) error
	grpc.ServerStream
}

type schedulerSubscribePipelineStatusServer struct {
	grpc.ServerStream
}

func (x *schedulerSubscribePipelineStatusServer) Send(m *PipelineStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Scheduler_ServiceDesc is the grpc.ServiceDesc for Scheduler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Scheduler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "seldon.mlops.scheduler.Scheduler",
	HandlerType: (*SchedulerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ServerNotify",
			Handler:    _Scheduler_ServerNotify_Handler,
		},
		{
			MethodName: "LoadModel",
			Handler:    _Scheduler_LoadModel_Handler,
		},
		{
			MethodName: "UnloadModel",
			Handler:    _Scheduler_UnloadModel_Handler,
		},
		{
			MethodName: "LoadPipeline",
			Handler:    _Scheduler_LoadPipeline_Handler,
		},
		{
			MethodName: "UnloadPipeline",
			Handler:    _Scheduler_UnloadPipeline_Handler,
		},
		{
			MethodName: "StartExperiment",
			Handler:    _Scheduler_StartExperiment_Handler,
		},
		{
			MethodName: "StopExperiment",
			Handler:    _Scheduler_StopExperiment_Handler,
		},
		{
			MethodName: "SchedulerStatus",
			Handler:    _Scheduler_SchedulerStatus_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ServerStatus",
			Handler:       _Scheduler_ServerStatus_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ModelStatus",
			Handler:       _Scheduler_ModelStatus_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PipelineStatus",
			Handler:       _Scheduler_PipelineStatus_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ExperimentStatus",
			Handler:       _Scheduler_ExperimentStatus_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeServerStatus",
			Handler:       _Scheduler_SubscribeServerStatus_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeModelStatus",
			Handler:       _Scheduler_SubscribeModelStatus_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeExperimentStatus",
			Handler:       _Scheduler_SubscribeExperimentStatus_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribePipelineStatus",
			Handler:       _Scheduler_SubscribePipelineStatus_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "mlops/scheduler/scheduler.proto",
}
