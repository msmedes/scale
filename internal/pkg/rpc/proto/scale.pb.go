// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/scale.proto

package scale

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea80762c827c7836, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Success struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Success) Reset()         { *m = Success{} }
func (m *Success) String() string { return proto.CompactTextString(m) }
func (*Success) ProtoMessage()    {}
func (*Success) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea80762c827c7836, []int{1}
}

func (m *Success) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Success.Unmarshal(m, b)
}
func (m *Success) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Success.Marshal(b, m, deterministic)
}
func (m *Success) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Success.Merge(m, src)
}
func (m *Success) XXX_Size() int {
	return xxx_messageInfo_Success.Size(m)
}
func (m *Success) XXX_DiscardUnknown() {
	xxx_messageInfo_Success.DiscardUnknown(m)
}

var xxx_messageInfo_Success proto.InternalMessageInfo

type RemoteId struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoteId) Reset()         { *m = RemoteId{} }
func (m *RemoteId) String() string { return proto.CompactTextString(m) }
func (*RemoteId) ProtoMessage()    {}
func (*RemoteId) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea80762c827c7836, []int{2}
}

func (m *RemoteId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoteId.Unmarshal(m, b)
}
func (m *RemoteId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoteId.Marshal(b, m, deterministic)
}
func (m *RemoteId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoteId.Merge(m, src)
}
func (m *RemoteId) XXX_Size() int {
	return xxx_messageInfo_RemoteId.Size(m)
}
func (m *RemoteId) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoteId.DiscardUnknown(m)
}

var xxx_messageInfo_RemoteId proto.InternalMessageInfo

func (m *RemoteId) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

type RemoteQuery struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoteQuery) Reset()         { *m = RemoteQuery{} }
func (m *RemoteQuery) String() string { return proto.CompactTextString(m) }
func (*RemoteQuery) ProtoMessage()    {}
func (*RemoteQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea80762c827c7836, []int{3}
}

func (m *RemoteQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoteQuery.Unmarshal(m, b)
}
func (m *RemoteQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoteQuery.Marshal(b, m, deterministic)
}
func (m *RemoteQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoteQuery.Merge(m, src)
}
func (m *RemoteQuery) XXX_Size() int {
	return xxx_messageInfo_RemoteQuery.Size(m)
}
func (m *RemoteQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoteQuery.DiscardUnknown(m)
}

var xxx_messageInfo_RemoteQuery proto.InternalMessageInfo

func (m *RemoteQuery) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

type GetRequest struct {
	Key                  []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea80762c827c7836, []int{4}
}

func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (m *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(m, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func (m *GetRequest) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

type GetResponse struct {
	Value                []byte   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}
func (*GetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea80762c827c7836, []int{5}
}

func (m *GetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetResponse.Unmarshal(m, b)
}
func (m *GetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetResponse.Marshal(b, m, deterministic)
}
func (m *GetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetResponse.Merge(m, src)
}
func (m *GetResponse) XXX_Size() int {
	return xxx_messageInfo_GetResponse.Size(m)
}
func (m *GetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetResponse proto.InternalMessageInfo

func (m *GetResponse) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type KeyTransferRequest struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeyTransferRequest) Reset()         { *m = KeyTransferRequest{} }
func (m *KeyTransferRequest) String() string { return proto.CompactTextString(m) }
func (*KeyTransferRequest) ProtoMessage()    {}
func (*KeyTransferRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea80762c827c7836, []int{6}
}

func (m *KeyTransferRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyTransferRequest.Unmarshal(m, b)
}
func (m *KeyTransferRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyTransferRequest.Marshal(b, m, deterministic)
}
func (m *KeyTransferRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyTransferRequest.Merge(m, src)
}
func (m *KeyTransferRequest) XXX_Size() int {
	return xxx_messageInfo_KeyTransferRequest.Size(m)
}
func (m *KeyTransferRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyTransferRequest.DiscardUnknown(m)
}

var xxx_messageInfo_KeyTransferRequest proto.InternalMessageInfo

func (m *KeyTransferRequest) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *KeyTransferRequest) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type SetRequest struct {
	Key                  []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetRequest) Reset()         { *m = SetRequest{} }
func (m *SetRequest) String() string { return proto.CompactTextString(m) }
func (*SetRequest) ProtoMessage()    {}
func (*SetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea80762c827c7836, []int{7}
}

func (m *SetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetRequest.Unmarshal(m, b)
}
func (m *SetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetRequest.Marshal(b, m, deterministic)
}
func (m *SetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetRequest.Merge(m, src)
}
func (m *SetRequest) XXX_Size() int {
	return xxx_messageInfo_SetRequest.Size(m)
}
func (m *SetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetRequest proto.InternalMessageInfo

func (m *SetRequest) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *SetRequest) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type NodeMetadata struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	PredecessorId        []byte   `protobuf:"bytes,3,opt,name=predecessorId,proto3" json:"predecessorId,omitempty"`
	PredecessorAddr      string   `protobuf:"bytes,4,opt,name=predecessorAddr,proto3" json:"predecessorAddr,omitempty"`
	SuccessorId          []byte   `protobuf:"bytes,5,opt,name=successorId,proto3" json:"successorId,omitempty"`
	SuccessorAddr        string   `protobuf:"bytes,6,opt,name=successorAddr,proto3" json:"successorAddr,omitempty"`
	FingerTable          [][]byte `protobuf:"bytes,7,rep,name=fingerTable,proto3" json:"fingerTable,omitempty"`
	Port                 string   `protobuf:"bytes,8,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeMetadata) Reset()         { *m = NodeMetadata{} }
func (m *NodeMetadata) String() string { return proto.CompactTextString(m) }
func (*NodeMetadata) ProtoMessage()    {}
func (*NodeMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea80762c827c7836, []int{8}
}

func (m *NodeMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeMetadata.Unmarshal(m, b)
}
func (m *NodeMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeMetadata.Marshal(b, m, deterministic)
}
func (m *NodeMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeMetadata.Merge(m, src)
}
func (m *NodeMetadata) XXX_Size() int {
	return xxx_messageInfo_NodeMetadata.Size(m)
}
func (m *NodeMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_NodeMetadata proto.InternalMessageInfo

func (m *NodeMetadata) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *NodeMetadata) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *NodeMetadata) GetPredecessorId() []byte {
	if m != nil {
		return m.PredecessorId
	}
	return nil
}

func (m *NodeMetadata) GetPredecessorAddr() string {
	if m != nil {
		return m.PredecessorAddr
	}
	return ""
}

func (m *NodeMetadata) GetSuccessorId() []byte {
	if m != nil {
		return m.SuccessorId
	}
	return nil
}

func (m *NodeMetadata) GetSuccessorAddr() string {
	if m != nil {
		return m.SuccessorAddr
	}
	return ""
}

func (m *NodeMetadata) GetFingerTable() [][]byte {
	if m != nil {
		return m.FingerTable
	}
	return nil
}

func (m *NodeMetadata) GetPort() string {
	if m != nil {
		return m.Port
	}
	return ""
}

type RemoteNode struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Present              bool     `protobuf:"varint,3,opt,name=present,proto3" json:"present,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoteNode) Reset()         { *m = RemoteNode{} }
func (m *RemoteNode) String() string { return proto.CompactTextString(m) }
func (*RemoteNode) ProtoMessage()    {}
func (*RemoteNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea80762c827c7836, []int{9}
}

func (m *RemoteNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoteNode.Unmarshal(m, b)
}
func (m *RemoteNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoteNode.Marshal(b, m, deterministic)
}
func (m *RemoteNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoteNode.Merge(m, src)
}
func (m *RemoteNode) XXX_Size() int {
	return xxx_messageInfo_RemoteNode.Size(m)
}
func (m *RemoteNode) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoteNode.DiscardUnknown(m)
}

var xxx_messageInfo_RemoteNode proto.InternalMessageInfo

func (m *RemoteNode) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *RemoteNode) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *RemoteNode) GetPresent() bool {
	if m != nil {
		return m.Present
	}
	return false
}

func init() {
	proto.RegisterType((*Empty)(nil), "scale.Empty")
	proto.RegisterType((*Success)(nil), "scale.Success")
	proto.RegisterType((*RemoteId)(nil), "scale.RemoteId")
	proto.RegisterType((*RemoteQuery)(nil), "scale.RemoteQuery")
	proto.RegisterType((*GetRequest)(nil), "scale.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "scale.GetResponse")
	proto.RegisterType((*KeyTransferRequest)(nil), "scale.KeyTransferRequest")
	proto.RegisterType((*SetRequest)(nil), "scale.SetRequest")
	proto.RegisterType((*NodeMetadata)(nil), "scale.NodeMetadata")
	proto.RegisterType((*RemoteNode)(nil), "scale.RemoteNode")
}

func init() { proto.RegisterFile("proto/scale.proto", fileDescriptor_ea80762c827c7836) }

var fileDescriptor_ea80762c827c7836 = []byte{
	// 491 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x41, 0x6f, 0xd3, 0x4c,
	0x10, 0x4d, 0xe2, 0x38, 0x71, 0x27, 0x6e, 0xfa, 0x75, 0x3e, 0x0e, 0x26, 0x12, 0x28, 0x5a, 0x10,
	0x8a, 0x90, 0x48, 0x05, 0x45, 0xa8, 0x12, 0x27, 0x0e, 0x60, 0x95, 0x42, 0x55, 0xec, 0xfe, 0x01,
	0x37, 0x3b, 0xa9, 0x2c, 0x52, 0xdb, 0xec, 0x6e, 0x90, 0x7c, 0xe3, 0x5f, 0xf1, 0xf7, 0xd0, 0xae,
	0xed, 0xda, 0x8e, 0xab, 0x2a, 0xb7, 0x99, 0x37, 0x33, 0xef, 0x79, 0x26, 0x6f, 0x03, 0xc7, 0x99,
	0x48, 0x55, 0x7a, 0x22, 0x57, 0xd1, 0x86, 0x96, 0x26, 0x46, 0xdb, 0x24, 0x6c, 0x0c, 0xf6, 0xe7,
	0xbb, 0x4c, 0xe5, 0xec, 0x00, 0xc6, 0xe1, 0x76, 0xb5, 0x22, 0x29, 0xd9, 0x0c, 0x9c, 0x80, 0xee,
	0x52, 0x45, 0xe7, 0x1c, 0xa7, 0x30, 0x88, 0xb9, 0xd7, 0x9f, 0xf7, 0x17, 0x6e, 0x30, 0x88, 0x39,
	0x7b, 0x06, 0x93, 0xa2, 0xf6, 0x63, 0x4b, 0x22, 0xef, 0x94, 0x9f, 0x03, 0xf8, 0xa4, 0x02, 0xfa,
	0xb5, 0x25, 0xa9, 0xf0, 0x3f, 0xb0, 0x7e, 0x52, 0x5e, 0x96, 0x75, 0xc8, 0x5e, 0xc0, 0xc4, 0xd4,
	0x65, 0x96, 0x26, 0x92, 0xf0, 0x09, 0xd8, 0xbf, 0xa3, 0xcd, 0x96, 0xca, 0x96, 0x22, 0x61, 0x67,
	0x80, 0x17, 0x94, 0x5f, 0x8b, 0x28, 0x91, 0x6b, 0x12, 0x15, 0xd9, 0x8e, 0x14, 0x22, 0x0c, 0x23,
	0xce, 0x85, 0x37, 0x98, 0xf7, 0x17, 0x07, 0x81, 0x89, 0xd9, 0x7b, 0x80, 0xf0, 0x11, 0xf9, 0x5a,
	0x6f, 0xd0, 0xd4, 0xfb, 0x33, 0x00, 0xf7, 0x32, 0xe5, 0xf4, 0x9d, 0x54, 0xc4, 0x23, 0x15, 0xed,
	0x23, 0x85, 0x2f, 0xe1, 0x30, 0x13, 0xc4, 0x49, 0x5f, 0x2c, 0x15, 0xe7, 0xdc, 0xb3, 0x4c, 0x7b,
	0x1b, 0xc4, 0x05, 0x1c, 0x35, 0x80, 0x4f, 0x9a, 0x64, 0x68, 0x48, 0x76, 0x61, 0x9c, 0xc3, 0x44,
	0x16, 0xf7, 0x37, 0x6c, 0xb6, 0x61, 0x6b, 0x42, 0x5a, 0xf1, 0x3e, 0x35, 0x4c, 0x23, 0xc3, 0xd4,
	0x06, 0x35, 0xcf, 0x3a, 0x4e, 0x6e, 0x49, 0x5c, 0x47, 0x37, 0x1b, 0xf2, 0xc6, 0x73, 0x4b, 0xf3,
	0x34, 0x20, 0xbd, 0x4d, 0x96, 0x0a, 0xe5, 0x39, 0xc5, 0x36, 0x3a, 0x66, 0x5f, 0x01, 0x8a, 0x9f,
	0x55, 0xdf, 0x61, 0xaf, 0xfd, 0x3d, 0x18, 0x67, 0x82, 0x24, 0x25, 0xca, 0x6c, 0xee, 0x04, 0x55,
	0xfa, 0xee, 0xef, 0x10, 0xec, 0x50, 0x9b, 0x0b, 0xdf, 0x82, 0xeb, 0x93, 0x0a, 0xab, 0xef, 0x43,
	0x77, 0x59, 0x38, 0xd0, 0x38, 0x6e, 0x76, 0x5c, 0x66, 0xb5, 0x30, 0xeb, 0xe1, 0x29, 0x4c, 0x7d,
	0x52, 0x57, 0xf5, 0x71, 0xf6, 0x19, 0x3a, 0x83, 0xc3, 0x2f, 0x71, 0xc2, 0x6b, 0x21, 0x6c, 0x75,
	0x19, 0xab, 0x3e, 0x3c, 0xf9, 0x06, 0x46, 0x97, 0xa9, 0x8a, 0xd7, 0x39, 0x76, 0xcb, 0xb3, 0x69,
	0x09, 0x55, 0xef, 0x42, 0x7f, 0x9d, 0xe3, 0x93, 0xfa, 0x96, 0xae, 0xa2, 0xcd, 0xfd, 0x40, 0xed,
	0xf7, 0x19, 0x36, 0xa1, 0xc2, 0xe2, 0xac, 0x87, 0x27, 0xe0, 0x84, 0xbb, 0x43, 0xb5, 0x4b, 0x1f,
	0x50, 0x59, 0x82, 0xe5, 0x93, 0xda, 0x5f, 0xe0, 0x35, 0x58, 0x61, 0xa3, 0xff, 0x51, 0xee, 0x0f,
	0x70, 0xe4, 0x93, 0x6a, 0xb9, 0xbd, 0x7d, 0xe0, 0xff, 0xcb, 0xac, 0xd9, 0xc2, 0x7a, 0xf8, 0x0a,
	0x86, 0x57, 0x71, 0x72, 0xbb, 0xd3, 0xdc, 0xe5, 0xff, 0x08, 0x6e, 0xf5, 0x70, 0x2f, 0x28, 0x97,
	0xf8, 0xb4, 0xec, 0xe8, 0x3e, 0xe8, 0xee, 0xf0, 0xcd, 0xc8, 0xfc, 0x35, 0x9d, 0xfe, 0x0b, 0x00,
	0x00, 0xff, 0xff, 0xbf, 0xf8, 0x23, 0x83, 0xaf, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ScaleClient is the client API for Scale service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ScaleClient interface {
	GetSuccessor(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*RemoteNode, error)
	GetPredecessor(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*RemoteNode, error)
	FindSuccessor(ctx context.Context, in *RemoteQuery, opts ...grpc.CallOption) (*RemoteNode, error)
	Notify(ctx context.Context, in *RemoteNode, opts ...grpc.CallOption) (*Success, error)
	GetLocal(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	SetLocal(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*Success, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*Success, error)
	GetNodeMetadata(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NodeMetadata, error)
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Success, error)
	TransferKeys(ctx context.Context, in *KeyTransferRequest, opts ...grpc.CallOption) (*Success, error)
}

type scaleClient struct {
	cc *grpc.ClientConn
}

func NewScaleClient(cc *grpc.ClientConn) ScaleClient {
	return &scaleClient{cc}
}

func (c *scaleClient) GetSuccessor(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*RemoteNode, error) {
	out := new(RemoteNode)
	err := c.cc.Invoke(ctx, "/scale.Scale/GetSuccessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scaleClient) GetPredecessor(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*RemoteNode, error) {
	out := new(RemoteNode)
	err := c.cc.Invoke(ctx, "/scale.Scale/GetPredecessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scaleClient) FindSuccessor(ctx context.Context, in *RemoteQuery, opts ...grpc.CallOption) (*RemoteNode, error) {
	out := new(RemoteNode)
	err := c.cc.Invoke(ctx, "/scale.Scale/FindSuccessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scaleClient) Notify(ctx context.Context, in *RemoteNode, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/scale.Scale/Notify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scaleClient) GetLocal(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/scale.Scale/GetLocal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scaleClient) SetLocal(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/scale.Scale/SetLocal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scaleClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/scale.Scale/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scaleClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/scale.Scale/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scaleClient) GetNodeMetadata(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NodeMetadata, error) {
	out := new(NodeMetadata)
	err := c.cc.Invoke(ctx, "/scale.Scale/GetNodeMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scaleClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/scale.Scale/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scaleClient) TransferKeys(ctx context.Context, in *KeyTransferRequest, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/scale.Scale/TransferKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScaleServer is the server API for Scale service.
type ScaleServer interface {
	GetSuccessor(context.Context, *Empty) (*RemoteNode, error)
	GetPredecessor(context.Context, *Empty) (*RemoteNode, error)
	FindSuccessor(context.Context, *RemoteQuery) (*RemoteNode, error)
	Notify(context.Context, *RemoteNode) (*Success, error)
	GetLocal(context.Context, *GetRequest) (*GetResponse, error)
	SetLocal(context.Context, *SetRequest) (*Success, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Set(context.Context, *SetRequest) (*Success, error)
	GetNodeMetadata(context.Context, *Empty) (*NodeMetadata, error)
	Ping(context.Context, *Empty) (*Success, error)
	TransferKeys(context.Context, *KeyTransferRequest) (*Success, error)
}

// UnimplementedScaleServer can be embedded to have forward compatible implementations.
type UnimplementedScaleServer struct {
}

func (*UnimplementedScaleServer) GetSuccessor(ctx context.Context, req *Empty) (*RemoteNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSuccessor not implemented")
}
func (*UnimplementedScaleServer) GetPredecessor(ctx context.Context, req *Empty) (*RemoteNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPredecessor not implemented")
}
func (*UnimplementedScaleServer) FindSuccessor(ctx context.Context, req *RemoteQuery) (*RemoteNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindSuccessor not implemented")
}
func (*UnimplementedScaleServer) Notify(ctx context.Context, req *RemoteNode) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Notify not implemented")
}
func (*UnimplementedScaleServer) GetLocal(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLocal not implemented")
}
func (*UnimplementedScaleServer) SetLocal(ctx context.Context, req *SetRequest) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLocal not implemented")
}
func (*UnimplementedScaleServer) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedScaleServer) Set(ctx context.Context, req *SetRequest) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (*UnimplementedScaleServer) GetNodeMetadata(ctx context.Context, req *Empty) (*NodeMetadata, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodeMetadata not implemented")
}
func (*UnimplementedScaleServer) Ping(ctx context.Context, req *Empty) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedScaleServer) TransferKeys(ctx context.Context, req *KeyTransferRequest) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransferKeys not implemented")
}

func RegisterScaleServer(s *grpc.Server, srv ScaleServer) {
	s.RegisterService(&_Scale_serviceDesc, srv)
}

func _Scale_GetSuccessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).GetSuccessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/GetSuccessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).GetSuccessor(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scale_GetPredecessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).GetPredecessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/GetPredecessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).GetPredecessor(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scale_FindSuccessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoteQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).FindSuccessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/FindSuccessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).FindSuccessor(ctx, req.(*RemoteQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scale_Notify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoteNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).Notify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/Notify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).Notify(ctx, req.(*RemoteNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scale_GetLocal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).GetLocal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/GetLocal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).GetLocal(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scale_SetLocal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).SetLocal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/SetLocal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).SetLocal(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scale_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scale_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scale_GetNodeMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).GetNodeMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/GetNodeMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).GetNodeMetadata(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scale_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scale_TransferKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyTransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScaleServer).TransferKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scale.Scale/TransferKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScaleServer).TransferKeys(ctx, req.(*KeyTransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Scale_serviceDesc = grpc.ServiceDesc{
	ServiceName: "scale.Scale",
	HandlerType: (*ScaleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSuccessor",
			Handler:    _Scale_GetSuccessor_Handler,
		},
		{
			MethodName: "GetPredecessor",
			Handler:    _Scale_GetPredecessor_Handler,
		},
		{
			MethodName: "FindSuccessor",
			Handler:    _Scale_FindSuccessor_Handler,
		},
		{
			MethodName: "Notify",
			Handler:    _Scale_Notify_Handler,
		},
		{
			MethodName: "GetLocal",
			Handler:    _Scale_GetLocal_Handler,
		},
		{
			MethodName: "SetLocal",
			Handler:    _Scale_SetLocal_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Scale_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _Scale_Set_Handler,
		},
		{
			MethodName: "GetNodeMetadata",
			Handler:    _Scale_GetNodeMetadata_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Scale_Ping_Handler,
		},
		{
			MethodName: "TransferKeys",
			Handler:    _Scale_TransferKeys_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/scale.proto",
}
