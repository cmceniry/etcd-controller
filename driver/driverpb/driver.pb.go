// Code generated by protoc-gen-go. DO NOT EDIT.
// source: driver.proto

package driverpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StatusRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusRequest) Reset()         { *m = StatusRequest{} }
func (m *StatusRequest) String() string { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()    {}
func (*StatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{0}
}
func (m *StatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusRequest.Unmarshal(m, b)
}
func (m *StatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusRequest.Marshal(b, m, deterministic)
}
func (dst *StatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusRequest.Merge(dst, src)
}
func (m *StatusRequest) XXX_Size() int {
	return xxx_messageInfo_StatusRequest.Size(m)
}
func (m *StatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StatusRequest proto.InternalMessageInfo

type StatusResponse struct {
	State                int32    `protobuf:"varint,1,opt,name=state,proto3" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusResponse) Reset()         { *m = StatusResponse{} }
func (m *StatusResponse) String() string { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()    {}
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{1}
}
func (m *StatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusResponse.Unmarshal(m, b)
}
func (m *StatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusResponse.Marshal(b, m, deterministic)
}
func (dst *StatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusResponse.Merge(dst, src)
}
func (m *StatusResponse) XXX_Size() int {
	return xxx_messageInfo_StatusResponse.Size(m)
}
func (m *StatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StatusResponse proto.InternalMessageInfo

func (m *StatusResponse) GetState() int32 {
	if m != nil {
		return m.State
	}
	return 0
}

type InitClusterRequest struct {
	Force                bool     `protobuf:"varint,1,opt,name=force,proto3" json:"force,omitempty"`
	Snapshot             string   `protobuf:"bytes,2,opt,name=snapshot,proto3" json:"snapshot,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitClusterRequest) Reset()         { *m = InitClusterRequest{} }
func (m *InitClusterRequest) String() string { return proto.CompactTextString(m) }
func (*InitClusterRequest) ProtoMessage()    {}
func (*InitClusterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{2}
}
func (m *InitClusterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitClusterRequest.Unmarshal(m, b)
}
func (m *InitClusterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitClusterRequest.Marshal(b, m, deterministic)
}
func (dst *InitClusterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitClusterRequest.Merge(dst, src)
}
func (m *InitClusterRequest) XXX_Size() int {
	return xxx_messageInfo_InitClusterRequest.Size(m)
}
func (m *InitClusterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InitClusterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InitClusterRequest proto.InternalMessageInfo

func (m *InitClusterRequest) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

func (m *InitClusterRequest) GetSnapshot() string {
	if m != nil {
		return m.Snapshot
	}
	return ""
}

type InitClusterResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage         string   `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitClusterResponse) Reset()         { *m = InitClusterResponse{} }
func (m *InitClusterResponse) String() string { return proto.CompactTextString(m) }
func (*InitClusterResponse) ProtoMessage()    {}
func (*InitClusterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{3}
}
func (m *InitClusterResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitClusterResponse.Unmarshal(m, b)
}
func (m *InitClusterResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitClusterResponse.Marshal(b, m, deterministic)
}
func (dst *InitClusterResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitClusterResponse.Merge(dst, src)
}
func (m *InitClusterResponse) XXX_Size() int {
	return xxx_messageInfo_InitClusterResponse.Size(m)
}
func (m *InitClusterResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InitClusterResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InitClusterResponse proto.InternalMessageInfo

func (m *InitClusterResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *InitClusterResponse) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

type PeerInfo struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	URL                  string   `protobuf:"bytes,2,opt,name=URL,proto3" json:"URL,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PeerInfo) Reset()         { *m = PeerInfo{} }
func (m *PeerInfo) String() string { return proto.CompactTextString(m) }
func (*PeerInfo) ProtoMessage()    {}
func (*PeerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{4}
}
func (m *PeerInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerInfo.Unmarshal(m, b)
}
func (m *PeerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerInfo.Marshal(b, m, deterministic)
}
func (dst *PeerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerInfo.Merge(dst, src)
}
func (m *PeerInfo) XXX_Size() int {
	return xxx_messageInfo_PeerInfo.Size(m)
}
func (m *PeerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PeerInfo proto.InternalMessageInfo

func (m *PeerInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PeerInfo) GetURL() string {
	if m != nil {
		return m.URL
	}
	return ""
}

type JoinClusterRequest struct {
	Force                bool        `protobuf:"varint,1,opt,name=force,proto3" json:"force,omitempty"`
	Peers                []*PeerInfo `protobuf:"bytes,2,rep,name=peers,proto3" json:"peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *JoinClusterRequest) Reset()         { *m = JoinClusterRequest{} }
func (m *JoinClusterRequest) String() string { return proto.CompactTextString(m) }
func (*JoinClusterRequest) ProtoMessage()    {}
func (*JoinClusterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{5}
}
func (m *JoinClusterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinClusterRequest.Unmarshal(m, b)
}
func (m *JoinClusterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinClusterRequest.Marshal(b, m, deterministic)
}
func (dst *JoinClusterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinClusterRequest.Merge(dst, src)
}
func (m *JoinClusterRequest) XXX_Size() int {
	return xxx_messageInfo_JoinClusterRequest.Size(m)
}
func (m *JoinClusterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinClusterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_JoinClusterRequest proto.InternalMessageInfo

func (m *JoinClusterRequest) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

func (m *JoinClusterRequest) GetPeers() []*PeerInfo {
	if m != nil {
		return m.Peers
	}
	return nil
}

type JoinClusterResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage         string   `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JoinClusterResponse) Reset()         { *m = JoinClusterResponse{} }
func (m *JoinClusterResponse) String() string { return proto.CompactTextString(m) }
func (*JoinClusterResponse) ProtoMessage()    {}
func (*JoinClusterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{6}
}
func (m *JoinClusterResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinClusterResponse.Unmarshal(m, b)
}
func (m *JoinClusterResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinClusterResponse.Marshal(b, m, deterministic)
}
func (dst *JoinClusterResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinClusterResponse.Merge(dst, src)
}
func (m *JoinClusterResponse) XXX_Size() int {
	return xxx_messageInfo_JoinClusterResponse.Size(m)
}
func (m *JoinClusterResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinClusterResponse.DiscardUnknown(m)
}

var xxx_messageInfo_JoinClusterResponse proto.InternalMessageInfo

func (m *JoinClusterResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *JoinClusterResponse) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

type StopServerRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopServerRequest) Reset()         { *m = StopServerRequest{} }
func (m *StopServerRequest) String() string { return proto.CompactTextString(m) }
func (*StopServerRequest) ProtoMessage()    {}
func (*StopServerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{7}
}
func (m *StopServerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopServerRequest.Unmarshal(m, b)
}
func (m *StopServerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopServerRequest.Marshal(b, m, deterministic)
}
func (dst *StopServerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopServerRequest.Merge(dst, src)
}
func (m *StopServerRequest) XXX_Size() int {
	return xxx_messageInfo_StopServerRequest.Size(m)
}
func (m *StopServerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StopServerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StopServerRequest proto.InternalMessageInfo

type StopServerResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage         string   `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopServerResponse) Reset()         { *m = StopServerResponse{} }
func (m *StopServerResponse) String() string { return proto.CompactTextString(m) }
func (*StopServerResponse) ProtoMessage()    {}
func (*StopServerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{8}
}
func (m *StopServerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopServerResponse.Unmarshal(m, b)
}
func (m *StopServerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopServerResponse.Marshal(b, m, deterministic)
}
func (dst *StopServerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopServerResponse.Merge(dst, src)
}
func (m *StopServerResponse) XXX_Size() int {
	return xxx_messageInfo_StopServerResponse.Size(m)
}
func (m *StopServerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StopServerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StopServerResponse proto.InternalMessageInfo

func (m *StopServerResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *StopServerResponse) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

type StartServerRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartServerRequest) Reset()         { *m = StartServerRequest{} }
func (m *StartServerRequest) String() string { return proto.CompactTextString(m) }
func (*StartServerRequest) ProtoMessage()    {}
func (*StartServerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{9}
}
func (m *StartServerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartServerRequest.Unmarshal(m, b)
}
func (m *StartServerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartServerRequest.Marshal(b, m, deterministic)
}
func (dst *StartServerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartServerRequest.Merge(dst, src)
}
func (m *StartServerRequest) XXX_Size() int {
	return xxx_messageInfo_StartServerRequest.Size(m)
}
func (m *StartServerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StartServerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StartServerRequest proto.InternalMessageInfo

type StartServerResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage         string   `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartServerResponse) Reset()         { *m = StartServerResponse{} }
func (m *StartServerResponse) String() string { return proto.CompactTextString(m) }
func (*StartServerResponse) ProtoMessage()    {}
func (*StartServerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_driver_57b503804075e8b6, []int{10}
}
func (m *StartServerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartServerResponse.Unmarshal(m, b)
}
func (m *StartServerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartServerResponse.Marshal(b, m, deterministic)
}
func (dst *StartServerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartServerResponse.Merge(dst, src)
}
func (m *StartServerResponse) XXX_Size() int {
	return xxx_messageInfo_StartServerResponse.Size(m)
}
func (m *StartServerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StartServerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StartServerResponse proto.InternalMessageInfo

func (m *StartServerResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *StartServerResponse) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

func init() {
	proto.RegisterType((*StatusRequest)(nil), "driverpb.StatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "driverpb.StatusResponse")
	proto.RegisterType((*InitClusterRequest)(nil), "driverpb.InitClusterRequest")
	proto.RegisterType((*InitClusterResponse)(nil), "driverpb.InitClusterResponse")
	proto.RegisterType((*PeerInfo)(nil), "driverpb.PeerInfo")
	proto.RegisterType((*JoinClusterRequest)(nil), "driverpb.JoinClusterRequest")
	proto.RegisterType((*JoinClusterResponse)(nil), "driverpb.JoinClusterResponse")
	proto.RegisterType((*StopServerRequest)(nil), "driverpb.StopServerRequest")
	proto.RegisterType((*StopServerResponse)(nil), "driverpb.StopServerResponse")
	proto.RegisterType((*StartServerRequest)(nil), "driverpb.StartServerRequest")
	proto.RegisterType((*StartServerResponse)(nil), "driverpb.StartServerResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DriverClient is the client API for Driver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DriverClient interface {
	GetStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	InitializeCluster(ctx context.Context, in *InitClusterRequest, opts ...grpc.CallOption) (*InitClusterResponse, error)
	JoinCluster(ctx context.Context, in *JoinClusterRequest, opts ...grpc.CallOption) (*JoinClusterResponse, error)
	StopServer(ctx context.Context, in *StopServerRequest, opts ...grpc.CallOption) (*StopServerResponse, error)
	StartServer(ctx context.Context, in *StartServerRequest, opts ...grpc.CallOption) (*StartServerResponse, error)
}

type driverClient struct {
	cc *grpc.ClientConn
}

func NewDriverClient(cc *grpc.ClientConn) DriverClient {
	return &driverClient{cc}
}

func (c *driverClient) GetStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/driverpb.driver/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) InitializeCluster(ctx context.Context, in *InitClusterRequest, opts ...grpc.CallOption) (*InitClusterResponse, error) {
	out := new(InitClusterResponse)
	err := c.cc.Invoke(ctx, "/driverpb.driver/InitializeCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) JoinCluster(ctx context.Context, in *JoinClusterRequest, opts ...grpc.CallOption) (*JoinClusterResponse, error) {
	out := new(JoinClusterResponse)
	err := c.cc.Invoke(ctx, "/driverpb.driver/JoinCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) StopServer(ctx context.Context, in *StopServerRequest, opts ...grpc.CallOption) (*StopServerResponse, error) {
	out := new(StopServerResponse)
	err := c.cc.Invoke(ctx, "/driverpb.driver/StopServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) StartServer(ctx context.Context, in *StartServerRequest, opts ...grpc.CallOption) (*StartServerResponse, error) {
	out := new(StartServerResponse)
	err := c.cc.Invoke(ctx, "/driverpb.driver/StartServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DriverServer is the server API for Driver service.
type DriverServer interface {
	GetStatus(context.Context, *StatusRequest) (*StatusResponse, error)
	InitializeCluster(context.Context, *InitClusterRequest) (*InitClusterResponse, error)
	JoinCluster(context.Context, *JoinClusterRequest) (*JoinClusterResponse, error)
	StopServer(context.Context, *StopServerRequest) (*StopServerResponse, error)
	StartServer(context.Context, *StartServerRequest) (*StartServerResponse, error)
}

func RegisterDriverServer(s *grpc.Server, srv DriverServer) {
	s.RegisterService(&_Driver_serviceDesc, srv)
}

func _Driver_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverpb.driver/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).GetStatus(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_InitializeCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).InitializeCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverpb.driver/InitializeCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).InitializeCluster(ctx, req.(*InitClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_JoinCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).JoinCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverpb.driver/JoinCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).JoinCluster(ctx, req.(*JoinClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_StopServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).StopServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverpb.driver/StopServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).StopServer(ctx, req.(*StopServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_StartServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).StartServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverpb.driver/StartServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).StartServer(ctx, req.(*StartServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Driver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "driverpb.driver",
	HandlerType: (*DriverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatus",
			Handler:    _Driver_GetStatus_Handler,
		},
		{
			MethodName: "InitializeCluster",
			Handler:    _Driver_InitializeCluster_Handler,
		},
		{
			MethodName: "JoinCluster",
			Handler:    _Driver_JoinCluster_Handler,
		},
		{
			MethodName: "StopServer",
			Handler:    _Driver_StopServer_Handler,
		},
		{
			MethodName: "StartServer",
			Handler:    _Driver_StartServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "driver.proto",
}

func init() { proto.RegisterFile("driver.proto", fileDescriptor_driver_57b503804075e8b6) }

var fileDescriptor_driver_57b503804075e8b6 = []byte{
	// 396 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x5b, 0x8b, 0xda, 0x40,
	0x14, 0xf6, 0x52, 0x6d, 0x3c, 0x6a, 0x5b, 0x8f, 0x42, 0x43, 0x6a, 0x41, 0xa6, 0x50, 0xf2, 0x24,
	0xc5, 0xfe, 0x81, 0x42, 0xa1, 0xc5, 0x62, 0x97, 0x65, 0xe2, 0x3e, 0x2f, 0xd1, 0x3d, 0xee, 0x06,
	0x34, 0x93, 0x9d, 0x99, 0xf8, 0xb0, 0xef, 0xfb, 0xbf, 0x97, 0x38, 0x89, 0x4e, 0xd4, 0x85, 0x7d,
	0xf0, 0x6d, 0xce, 0xed, 0x3b, 0x97, 0xef, 0x63, 0xa0, 0x73, 0x27, 0xa3, 0x2d, 0xc9, 0x71, 0x22,
	0x85, 0x16, 0xe8, 0x18, 0x2b, 0x59, 0xb0, 0x8f, 0xd0, 0x0d, 0x74, 0xa8, 0x53, 0xc5, 0xe9, 0x31,
	0x25, 0xa5, 0xd9, 0x77, 0xf8, 0x50, 0x38, 0x54, 0x22, 0x62, 0x45, 0x38, 0x80, 0x86, 0xd2, 0xa1,
	0x26, 0xb7, 0x3a, 0xaa, 0xfa, 0x0d, 0x6e, 0x0c, 0xf6, 0x07, 0x70, 0x1a, 0x47, 0xfa, 0xf7, 0x3a,
	0x55, 0x9a, 0x64, 0x5e, 0x9d, 0xe5, 0xae, 0x84, 0x5c, 0x9a, 0x5c, 0x87, 0x1b, 0x03, 0x3d, 0x70,
	0x54, 0x1c, 0x26, 0xea, 0x41, 0x68, 0xb7, 0x36, 0xaa, 0xfa, 0x2d, 0xbe, 0xb7, 0xd9, 0x1c, 0xfa,
	0x25, 0x9c, 0xbc, 0xa9, 0x0b, 0xef, 0x55, 0xba, 0x5c, 0x92, 0x52, 0x39, 0x54, 0x61, 0xe2, 0x37,
	0xe8, 0x92, 0x94, 0x42, 0xde, 0x6e, 0x48, 0xa9, 0xf0, 0x9e, 0x72, 0xc4, 0xce, 0xce, 0xf9, 0xdf,
	0xf8, 0xd8, 0x0f, 0x70, 0xae, 0x89, 0xe4, 0x34, 0x5e, 0x09, 0x44, 0x78, 0x77, 0x15, 0x6e, 0xcc,
	0x48, 0x2d, 0xbe, 0x7b, 0xe3, 0x27, 0xa8, 0xdf, 0xf0, 0x59, 0x5e, 0x9a, 0x3d, 0xd9, 0x1c, 0xf0,
	0x9f, 0x88, 0xe2, 0x37, 0xed, 0xe3, 0x43, 0x23, 0x21, 0x92, 0xca, 0xad, 0x8d, 0xea, 0x7e, 0x7b,
	0x82, 0xe3, 0xe2, 0x9c, 0xe3, 0xa2, 0x29, 0x37, 0x09, 0xd9, 0x76, 0x25, 0xd4, 0xcb, 0x6c, 0xd7,
	0x87, 0x5e, 0xa0, 0x45, 0x12, 0x90, 0xdc, 0xee, 0x47, 0x65, 0x01, 0xa0, 0xed, 0xbc, 0x4c, 0xa7,
	0x41, 0x06, 0x1a, 0x4a, 0x5d, 0x6e, 0x35, 0x87, 0x7e, 0xc9, 0x7b, 0x91, 0x5e, 0x93, 0xe7, 0x3a,
	0x34, 0xcd, 0x21, 0xf1, 0x17, 0xb4, 0xfe, 0x92, 0x36, 0x3a, 0xc4, 0xcf, 0x87, 0xf3, 0x96, 0xa4,
	0xea, 0xb9, 0xa7, 0x01, 0x33, 0x09, 0xab, 0x20, 0x87, 0x5e, 0x26, 0xab, 0x28, 0x5c, 0x47, 0x4f,
	0x94, 0x9f, 0x1f, 0x87, 0x87, 0x82, 0x53, 0xed, 0x7a, 0x5f, 0x5f, 0x89, 0xee, 0x31, 0x67, 0xd0,
	0xb6, 0xc8, 0xb4, 0xd1, 0x4e, 0x95, 0x63, 0xa3, 0x9d, 0x51, 0x00, 0xab, 0xe0, 0x14, 0xe0, 0xc0,
	0x17, 0x7e, 0xb1, 0x77, 0x39, 0xa2, 0xd6, 0x1b, 0x9e, 0x0f, 0xda, 0x83, 0x59, 0x7c, 0x60, 0x29,
	0xfd, 0x98, 0x3c, 0x7b, 0xb0, 0x33, 0x24, 0xb2, 0xca, 0xa2, 0xb9, 0xfb, 0x23, 0x7e, 0xbe, 0x04,
	0x00, 0x00, 0xff, 0xff, 0x03, 0x84, 0xd9, 0xa9, 0x33, 0x04, 0x00, 0x00,
}
