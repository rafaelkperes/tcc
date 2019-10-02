// Code generated by protoc-gen-go. DO NOT EDIT.
// source: data.proto

package pbdata

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
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

type Ints struct {
	Ints                 []int64  `protobuf:"zigzag64,1,rep,packed,name=ints,proto3" json:"ints,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ints) Reset()         { *m = Ints{} }
func (m *Ints) String() string { return proto.CompactTextString(m) }
func (*Ints) ProtoMessage()    {}
func (*Ints) Descriptor() ([]byte, []int) {
	return fileDescriptor_871986018790d2fd, []int{0}
}

func (m *Ints) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ints.Unmarshal(m, b)
}
func (m *Ints) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ints.Marshal(b, m, deterministic)
}
func (m *Ints) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ints.Merge(m, src)
}
func (m *Ints) XXX_Size() int {
	return xxx_messageInfo_Ints.Size(m)
}
func (m *Ints) XXX_DiscardUnknown() {
	xxx_messageInfo_Ints.DiscardUnknown(m)
}

var xxx_messageInfo_Ints proto.InternalMessageInfo

func (m *Ints) GetInts() []int64 {
	if m != nil {
		return m.Ints
	}
	return nil
}

type Floats struct {
	Floats               []float64 `protobuf:"fixed64,1,rep,packed,name=floats,proto3" json:"floats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Floats) Reset()         { *m = Floats{} }
func (m *Floats) String() string { return proto.CompactTextString(m) }
func (*Floats) ProtoMessage()    {}
func (*Floats) Descriptor() ([]byte, []int) {
	return fileDescriptor_871986018790d2fd, []int{1}
}

func (m *Floats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Floats.Unmarshal(m, b)
}
func (m *Floats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Floats.Marshal(b, m, deterministic)
}
func (m *Floats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Floats.Merge(m, src)
}
func (m *Floats) XXX_Size() int {
	return xxx_messageInfo_Floats.Size(m)
}
func (m *Floats) XXX_DiscardUnknown() {
	xxx_messageInfo_Floats.DiscardUnknown(m)
}

var xxx_messageInfo_Floats proto.InternalMessageInfo

func (m *Floats) GetFloats() []float64 {
	if m != nil {
		return m.Floats
	}
	return nil
}

type Strings struct {
	Strings              []string `protobuf:"bytes,1,rep,name=strings,proto3" json:"strings,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Strings) Reset()         { *m = Strings{} }
func (m *Strings) String() string { return proto.CompactTextString(m) }
func (*Strings) ProtoMessage()    {}
func (*Strings) Descriptor() ([]byte, []int) {
	return fileDescriptor_871986018790d2fd, []int{2}
}

func (m *Strings) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Strings.Unmarshal(m, b)
}
func (m *Strings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Strings.Marshal(b, m, deterministic)
}
func (m *Strings) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Strings.Merge(m, src)
}
func (m *Strings) XXX_Size() int {
	return xxx_messageInfo_Strings.Size(m)
}
func (m *Strings) XXX_DiscardUnknown() {
	xxx_messageInfo_Strings.DiscardUnknown(m)
}

var xxx_messageInfo_Strings proto.InternalMessageInfo

func (m *Strings) GetStrings() []string {
	if m != nil {
		return m.Strings
	}
	return nil
}

type Objects struct {
	Objects              []*Objects_Object `protobuf:"bytes,1,rep,name=objects,proto3" json:"objects,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Objects) Reset()         { *m = Objects{} }
func (m *Objects) String() string { return proto.CompactTextString(m) }
func (*Objects) ProtoMessage()    {}
func (*Objects) Descriptor() ([]byte, []int) {
	return fileDescriptor_871986018790d2fd, []int{3}
}

func (m *Objects) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Objects.Unmarshal(m, b)
}
func (m *Objects) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Objects.Marshal(b, m, deterministic)
}
func (m *Objects) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Objects.Merge(m, src)
}
func (m *Objects) XXX_Size() int {
	return xxx_messageInfo_Objects.Size(m)
}
func (m *Objects) XXX_DiscardUnknown() {
	xxx_messageInfo_Objects.DiscardUnknown(m)
}

var xxx_messageInfo_Objects proto.InternalMessageInfo

func (m *Objects) GetObjects() []*Objects_Object {
	if m != nil {
		return m.Objects
	}
	return nil
}

type Objects_Object struct {
	I                    int64    `protobuf:"zigzag64,1,opt,name=I,proto3" json:"I,omitempty"`
	F                    float64  `protobuf:"fixed64,2,opt,name=F,proto3" json:"F,omitempty"`
	T                    bool     `protobuf:"varint,3,opt,name=T,proto3" json:"T,omitempty"`
	S                    string   `protobuf:"bytes,4,opt,name=S,proto3" json:"S,omitempty"`
	B                    []byte   `protobuf:"bytes,5,opt,name=B,proto3" json:"B,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Objects_Object) Reset()         { *m = Objects_Object{} }
func (m *Objects_Object) String() string { return proto.CompactTextString(m) }
func (*Objects_Object) ProtoMessage()    {}
func (*Objects_Object) Descriptor() ([]byte, []int) {
	return fileDescriptor_871986018790d2fd, []int{3, 0}
}

func (m *Objects_Object) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Objects_Object.Unmarshal(m, b)
}
func (m *Objects_Object) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Objects_Object.Marshal(b, m, deterministic)
}
func (m *Objects_Object) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Objects_Object.Merge(m, src)
}
func (m *Objects_Object) XXX_Size() int {
	return xxx_messageInfo_Objects_Object.Size(m)
}
func (m *Objects_Object) XXX_DiscardUnknown() {
	xxx_messageInfo_Objects_Object.DiscardUnknown(m)
}

var xxx_messageInfo_Objects_Object proto.InternalMessageInfo

func (m *Objects_Object) GetI() int64 {
	if m != nil {
		return m.I
	}
	return 0
}

func (m *Objects_Object) GetF() float64 {
	if m != nil {
		return m.F
	}
	return 0
}

func (m *Objects_Object) GetT() bool {
	if m != nil {
		return m.T
	}
	return false
}

func (m *Objects_Object) GetS() string {
	if m != nil {
		return m.S
	}
	return ""
}

func (m *Objects_Object) GetB() []byte {
	if m != nil {
		return m.B
	}
	return nil
}

func init() {
	proto.RegisterType((*Ints)(nil), "data.Ints")
	proto.RegisterType((*Floats)(nil), "data.Floats")
	proto.RegisterType((*Strings)(nil), "data.Strings")
	proto.RegisterType((*Objects)(nil), "data.Objects")
	proto.RegisterType((*Objects_Object)(nil), "data.Objects.Object")
}

func init() { proto.RegisterFile("data.proto", fileDescriptor_871986018790d2fd) }

var fileDescriptor_871986018790d2fd = []byte{
	// 233 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8f, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x40, 0x75, 0x34, 0x38, 0xed, 0xd1, 0xe9, 0x84, 0x90, 0x95, 0x05, 0x2b, 0x2c, 0x9e, 0x52,
	0x09, 0xfe, 0x20, 0x43, 0xa4, 0x2c, 0x20, 0x5d, 0x3a, 0xb1, 0x25, 0x90, 0x46, 0x41, 0x6d, 0x1c,
	0xd5, 0xc7, 0x47, 0xf0, 0xd7, 0xc8, 0x4e, 0x3a, 0xf9, 0xbd, 0xbb, 0x27, 0xd9, 0x46, 0xfc, 0x6e,
	0xa5, 0x2d, 0xe6, 0xab, 0x13, 0x47, 0x49, 0xe0, 0xec, 0x79, 0x70, 0x6e, 0x38, 0xf7, 0x87, 0x38,
	0xeb, 0x7e, 0x4f, 0x07, 0x19, 0x2f, 0xbd, 0x97, 0xf6, 0x32, 0x2f, 0x59, 0x9e, 0x61, 0x52, 0x4f,
	0xe2, 0x89, 0x30, 0x19, 0x27, 0xf1, 0x1a, 0xcc, 0xc6, 0x12, 0x47, 0xce, 0x0d, 0xaa, 0xea, 0xec,
	0x5a, 0xf1, 0xf4, 0x84, 0xea, 0x14, 0x29, 0xee, 0x81, 0x57, 0xcb, 0x5f, 0x30, 0x6d, 0xe4, 0x3a,
	0x4e, 0x83, 0x27, 0x8d, 0xa9, 0x5f, 0x30, 0x36, 0x3b, 0xbe, 0x69, 0xfe, 0x07, 0x98, 0x7e, 0x74,
	0x3f, 0xfd, 0x97, 0x78, 0x2a, 0x30, 0x75, 0x0b, 0xc6, 0xea, 0xe1, 0xf5, 0xb1, 0x88, 0x6f, 0x5e,
	0xf7, 0xeb, 0xc9, 0xb7, 0x28, 0x7b, 0x47, 0xb5, 0x8c, 0x68, 0x8f, 0x50, 0x6b, 0x30, 0x60, 0x89,
	0xa1, 0x0e, 0x56, 0xe9, 0x3b, 0x03, 0x16, 0x18, 0xaa, 0x60, 0x47, 0xbd, 0x31, 0x60, 0xb7, 0x0c,
	0xc7, 0x60, 0x8d, 0x4e, 0x0c, 0xd8, 0x1d, 0x43, 0x13, 0xac, 0xd4, 0xf7, 0x06, 0xec, 0x9e, 0xa1,
	0x2c, 0xb7, 0x9f, 0x6a, 0xee, 0xc2, 0x8d, 0x9d, 0x8a, 0xff, 0x7f, 0xfb, 0x0f, 0x00, 0x00, 0xff,
	0xff, 0x3c, 0x3a, 0x9f, 0x9e, 0x34, 0x01, 0x00, 0x00,
}
