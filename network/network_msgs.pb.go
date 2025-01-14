// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/bloxapp/ssv/network/network_msgs.proto

package network

import (
	fmt "fmt"
	proto1 "github.com/bloxapp/ssv/ibft/proto"
	proto "github.com/golang/protobuf/proto"
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

type NetworkMsg int32

const (
	// IBFTType are all iBFT related messages
	NetworkMsg_IBFTType NetworkMsg = 0
	// DecidedType is an iBFT specific message for broadcasting post consensus decided message with signatures
	NetworkMsg_DecidedType NetworkMsg = 1
	// SignatureType is an SSV node specific message for broadcasting post consensus signatures on eth2 duties
	NetworkMsg_SignatureType NetworkMsg = 2
	// SyncType is an SSV iBFT specific message that a node uses to sync up with other nodes
	NetworkMsg_SyncType NetworkMsg = 3
)

var NetworkMsg_name = map[int32]string{
	0: "IBFTType",
	1: "DecidedType",
	2: "SignatureType",
	3: "SyncType",
}

var NetworkMsg_value = map[string]int32{
	"IBFTType":      0,
	"DecidedType":   1,
	"SignatureType": 2,
	"SyncType":      3,
}

func (x NetworkMsg) String() string {
	return proto.EnumName(NetworkMsg_name, int32(x))
}

func (NetworkMsg) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3bcf522045b4e980, []int{0}
}

type Sync int32

const (
	// GetHighestType is a request from peers to return the highest decided/ prepared instance they know of
	Sync_GetHighestType Sync = 0
	// GetInstanceRange is a request from peers to return instances and their decided/ prepared justifications
	Sync_GetInstanceRange Sync = 1
)

var Sync_name = map[int32]string{
	0: "GetHighestType",
	1: "GetInstanceRange",
}

var Sync_value = map[string]int32{
	"GetHighestType":   0,
	"GetInstanceRange": 1,
}

func (x Sync) String() string {
	return proto.EnumName(Sync_name, int32(x))
}

func (Sync) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3bcf522045b4e980, []int{1}
}

type SyncMessage struct {
	SignedMessages       []*proto1.SignedMessage `protobuf:"bytes,1,rep,name=SignedMessages,proto3" json:"SignedMessages,omitempty"`
	FromPeerID           string                  `protobuf:"bytes,2,opt,name=FromPeerID,proto3" json:"FromPeerID,omitempty"`
	Params               []uint64                `protobuf:"varint,3,rep,packed,name=params,proto3" json:"params,omitempty"`
	Lambda               []byte                  `protobuf:"bytes,4,opt,name=Lambda,json=Lambda,proto3" json:"Lambda,omitempty"`
	Type                 Sync                    `protobuf:"varint,5,opt,name=Type,proto3,enum=network.Sync" json:"Type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SyncMessage) Reset()         { *m = SyncMessage{} }
func (m *SyncMessage) String() string { return proto.CompactTextString(m) }
func (*SyncMessage) ProtoMessage()    {}
func (*SyncMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_3bcf522045b4e980, []int{0}
}

func (m *SyncMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncMessage.Unmarshal(m, b)
}
func (m *SyncMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncMessage.Marshal(b, m, deterministic)
}
func (m *SyncMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncMessage.Merge(m, src)
}
func (m *SyncMessage) XXX_Size() int {
	return xxx_messageInfo_SyncMessage.Size(m)
}
func (m *SyncMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SyncMessage proto.InternalMessageInfo

func (m *SyncMessage) GetSignedMessages() []*proto1.SignedMessage {
	if m != nil {
		return m.SignedMessages
	}
	return nil
}

func (m *SyncMessage) GetFromPeerID() string {
	if m != nil {
		return m.FromPeerID
	}
	return ""
}

func (m *SyncMessage) GetParams() []uint64 {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *SyncMessage) GetLambda() []byte {
	if m != nil {
		return m.Lambda
	}
	return nil
}

func (m *SyncMessage) GetType() Sync {
	if m != nil {
		return m.Type
	}
	return Sync_GetHighestType
}

// Message is a wrapper struct for all network message types
type Message struct {
	Lambda               []byte                `protobuf:"bytes,1,opt,name=Lambda,proto3" json:"Lambda,omitempty"`
	SignedMessage        *proto1.SignedMessage `protobuf:"bytes,2,opt,name=SignedMessage,proto3" json:"SignedMessage,omitempty"`
	SyncMessage          *SyncMessage          `protobuf:"bytes,3,opt,name=SyncMessage,proto3" json:"SyncMessage,omitempty"`
	Type                 NetworkMsg            `protobuf:"varint,4,opt,name=Type,proto3,enum=network.NetworkMsg" json:"Type,omitempty"`
	Stream               SyncStream            `protobuf:"bytes,5,opt,name=Stream,proto3" json:"Stream,omitempty"` //TODO (proto is not updated with stream field) need to find better solution!!!
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_3bcf522045b4e980, []int{1}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetLambda() []byte {
	if m != nil {
		return m.Lambda
	}
	return nil
}

func (m *Message) GetSignedMessage() *proto1.SignedMessage {
	if m != nil {
		return m.SignedMessage
	}
	return nil
}

func (m *Message) GetSyncMessage() *SyncMessage {
	if m != nil {
		return m.SyncMessage
	}
	return nil
}

func (m *Message) GetType() NetworkMsg {
	if m != nil {
		return m.Type
	}
	return NetworkMsg_IBFTType
}

func init() {
	proto.RegisterEnum("network.NetworkMsg", NetworkMsg_name, NetworkMsg_value)
	proto.RegisterEnum("network.Sync", Sync_name, Sync_value)
	proto.RegisterType((*SyncMessage)(nil), "network.SyncMessage")
	proto.RegisterType((*Message)(nil), "network.Message")
}

func init() {
	proto.RegisterFile("github.com/bloxapp/ssv/network/network_msgs.proto", fileDescriptor_3bcf522045b4e980)
}

var fileDescriptor_3bcf522045b4e980 = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x4f, 0x6f, 0xda, 0x40,
	0x10, 0xc5, 0xbb, 0xd8, 0x85, 0x76, 0x0c, 0x94, 0x6e, 0x11, 0xb2, 0x7a, 0x40, 0x86, 0x4b, 0x2d,
	0x0e, 0xa6, 0xa5, 0x52, 0x0e, 0x51, 0x4e, 0x08, 0x41, 0x90, 0x20, 0x42, 0x86, 0x53, 0x2e, 0x68,
	0x6d, 0xaf, 0x8c, 0x05, 0xfe, 0x23, 0xef, 0x42, 0xc2, 0x97, 0xcb, 0x3d, 0xdf, 0x2a, 0xf2, 0x7a,
	0xf9, 0xe3, 0x48, 0xc9, 0x69, 0x35, 0xbf, 0x99, 0x37, 0x7e, 0x6f, 0x64, 0xf8, 0xe7, 0x07, 0x7c,
	0xb3, 0x77, 0x2c, 0x37, 0x0e, 0xfb, 0xce, 0x2e, 0x7e, 0x26, 0x49, 0xd2, 0x67, 0xec, 0xd0, 0x8f,
	0x28, 0x7f, 0x8a, 0xd3, 0xed, 0xe9, 0x5d, 0x87, 0xcc, 0x67, 0x56, 0x92, 0xc6, 0x3c, 0xc6, 0x15,
	0xc9, 0x7e, 0xc3, 0x05, 0x76, 0x5f, 0x11, 0x68, 0xcb, 0x63, 0xe4, 0xce, 0x29, 0x63, 0xc4, 0xa7,
	0xf8, 0x0e, 0xea, 0xcb, 0xc0, 0x8f, 0xa8, 0x27, 0x01, 0xd3, 0x91, 0xa1, 0x98, 0xda, 0xa0, 0x99,
	0xcf, 0x5b, 0x85, 0xa6, 0xfd, 0x6e, 0x16, 0xb7, 0x01, 0xc6, 0x69, 0x1c, 0x2e, 0x28, 0x4d, 0xa7,
	0x23, 0xbd, 0x64, 0x20, 0xf3, 0xbb, 0x7d, 0x45, 0x70, 0x0b, 0xca, 0x09, 0x49, 0x49, 0xc8, 0x74,
	0xc5, 0x50, 0x4c, 0xd5, 0x96, 0x15, 0xee, 0x40, 0xf5, 0x40, 0x76, 0x81, 0x47, 0x78, 0x9c, 0xae,
	0x93, 0xad, 0xae, 0x1a, 0xc8, 0xac, 0xda, 0xda, 0x99, 0x2d, 0xb6, 0xb8, 0x03, 0xea, 0xea, 0x98,
	0x50, 0xfd, 0xab, 0x81, 0xcc, 0xfa, 0xa0, 0x66, 0xc9, 0x30, 0x56, 0x66, 0xde, 0x16, 0xad, 0xee,
	0x0b, 0x82, 0xca, 0x29, 0x47, 0x0b, 0xca, 0x33, 0x12, 0x3a, 0x1e, 0xd1, 0x91, 0xd8, 0x25, 0x2b,
	0x7c, 0x0b, 0xb5, 0x82, 0x67, 0x61, 0xf2, 0xa3, 0x78, 0xc5, 0x51, 0x7c, 0x53, 0x38, 0x95, 0xae,
	0x48, 0xe5, 0xb5, 0x93, 0x93, 0xb2, 0x70, 0xd3, 0x3f, 0xd2, 0xba, 0x2a, 0xac, 0xff, 0x3a, 0x0b,
	0x1e, 0xf2, 0x77, 0xce, 0xfc, 0x3c, 0x40, 0x6f, 0x06, 0x70, 0x61, 0xb8, 0x0a, 0xdf, 0xa6, 0xc3,
	0xf1, 0x2a, 0xeb, 0x34, 0xbe, 0xe0, 0x1f, 0xa0, 0x8d, 0xa8, 0x1b, 0x78, 0xd4, 0x13, 0x00, 0xe1,
	0x9f, 0x79, 0x12, 0xc2, 0xf7, 0x29, 0x15, 0xa8, 0x94, 0x29, 0xb2, 0xef, 0x8a, 0x4a, 0xe9, 0xfd,
	0x05, 0x35, 0xab, 0x30, 0x86, 0xfa, 0x84, 0xf2, 0xfb, 0xc0, 0xdf, 0x50, 0xc6, 0xe5, 0xb6, 0x26,
	0x34, 0x26, 0x94, 0x4f, 0x23, 0xc6, 0x49, 0xe4, 0x52, 0x9b, 0x44, 0x3e, 0x6d, 0xa0, 0xa1, 0xf1,
	0xd8, 0xfe, 0xfc, 0xb7, 0x72, 0xca, 0xe2, 0x4c, 0xff, 0xdf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x91,
	0x78, 0x69, 0x2f, 0x7f, 0x02, 0x00, 0x00,
}
