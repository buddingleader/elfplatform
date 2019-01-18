// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common/metadata.proto

package common

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Metadata struct {
	Id                   int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	CreateAuthorId       int32                `protobuf:"varint,3,opt,name=create_author_id,json=createAuthorId,proto3" json:"create_author_id,omitempty"`
	CreateAuthorName     string               `protobuf:"bytes,4,opt,name=create_author_name,json=createAuthorName,proto3" json:"create_author_name,omitempty"`
	Created              *timestamp.Timestamp `protobuf:"bytes,5,opt,name=created,proto3" json:"created,omitempty"`
	UpdateAuthorId       int32                `protobuf:"varint,6,opt,name=update_author_id,json=updateAuthorId,proto3" json:"update_author_id,omitempty"`
	UpdateAuthorName     string               `protobuf:"bytes,7,opt,name=update_author_name,json=updateAuthorName,proto3" json:"update_author_name,omitempty"`
	LastUpdated          *timestamp.Timestamp `protobuf:"bytes,8,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Metadata) Reset()         { *m = Metadata{} }
func (m *Metadata) String() string { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()    {}
func (*Metadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_106914311898ec53, []int{0}
}

func (m *Metadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metadata.Unmarshal(m, b)
}
func (m *Metadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metadata.Marshal(b, m, deterministic)
}
func (m *Metadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metadata.Merge(m, src)
}
func (m *Metadata) XXX_Size() int {
	return xxx_messageInfo_Metadata.Size(m)
}
func (m *Metadata) XXX_DiscardUnknown() {
	xxx_messageInfo_Metadata.DiscardUnknown(m)
}

var xxx_messageInfo_Metadata proto.InternalMessageInfo

func (m *Metadata) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Metadata) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Metadata) GetCreateAuthorId() int32 {
	if m != nil {
		return m.CreateAuthorId
	}
	return 0
}

func (m *Metadata) GetCreateAuthorName() string {
	if m != nil {
		return m.CreateAuthorName
	}
	return ""
}

func (m *Metadata) GetCreated() *timestamp.Timestamp {
	if m != nil {
		return m.Created
	}
	return nil
}

func (m *Metadata) GetUpdateAuthorId() int32 {
	if m != nil {
		return m.UpdateAuthorId
	}
	return 0
}

func (m *Metadata) GetUpdateAuthorName() string {
	if m != nil {
		return m.UpdateAuthorName
	}
	return ""
}

func (m *Metadata) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

func init() {
	proto.RegisterType((*Metadata)(nil), "common.Metadata")
}

func init() { proto.RegisterFile("common/metadata.proto", fileDescriptor_106914311898ec53) }

var fileDescriptor_106914311898ec53 = []byte{
	// 274 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0x41, 0x4b, 0xf4, 0x30,
	0x10, 0x86, 0x69, 0xbf, 0xdd, 0xee, 0x7e, 0x59, 0x59, 0x96, 0x80, 0x50, 0xf6, 0x62, 0xf1, 0x54,
	0x50, 0x1b, 0x50, 0xaf, 0x1e, 0xf4, 0xe6, 0x41, 0x0f, 0x45, 0x2f, 0x5e, 0x4a, 0xda, 0xa4, 0xdd,
	0x40, 0xb3, 0x53, 0xd2, 0xe9, 0x2f, 0xf5, 0x0f, 0x49, 0x33, 0x16, 0xda, 0x93, 0xa7, 0x24, 0x6f,
	0x9e, 0x79, 0x79, 0x60, 0xd8, 0x65, 0x05, 0xd6, 0xc2, 0x59, 0x58, 0x8d, 0x52, 0x49, 0x94, 0x59,
	0xe7, 0x00, 0x81, 0x47, 0x14, 0x1f, 0xaf, 0x1a, 0x80, 0xa6, 0xd5, 0xc2, 0xa7, 0xe5, 0x50, 0x0b,
	0x34, 0x56, 0xf7, 0x28, 0x6d, 0x47, 0xe0, 0xf5, 0x77, 0xc8, 0xb6, 0x6f, 0xbf, 0xb3, 0x7c, 0xcf,
	0x42, 0xa3, 0xe2, 0x20, 0x09, 0xd2, 0x75, 0x1e, 0x1a, 0xc5, 0x39, 0x5b, 0x9d, 0xa5, 0xd5, 0x71,
	0x98, 0x04, 0xe9, 0xff, 0xdc, 0xdf, 0x79, 0xca, 0x0e, 0x95, 0xd3, 0x12, 0x75, 0x21, 0x07, 0x3c,
	0x81, 0x2b, 0x8c, 0x8a, 0xff, 0xf9, 0x89, 0x3d, 0xe5, 0xcf, 0x3e, 0x7e, 0x55, 0xfc, 0x96, 0xf1,
	0x25, 0xe9, 0xbb, 0x56, 0xbe, 0xeb, 0x30, 0x67, 0xdf, 0xc7, 0xde, 0x47, 0xb6, 0xa1, 0x4c, 0xc5,
	0xeb, 0x24, 0x48, 0x77, 0xf7, 0xc7, 0x8c, 0xdc, 0xb3, 0xc9, 0x3d, 0xfb, 0x98, 0xdc, 0xf3, 0x09,
	0x1d, 0x6d, 0x86, 0x4e, 0x2d, 0x6d, 0x22, 0xb2, 0xa1, 0x7c, 0x6e, 0xb3, 0x24, 0xbd, 0xcd, 0x86,
	0x6c, 0xe6, 0xac, 0xb7, 0x79, 0x62, 0x17, 0xad, 0xec, 0xb1, 0xa0, 0x0f, 0x15, 0x6f, 0xff, 0x54,
	0xda, 0x8d, 0xfc, 0x27, 0xe1, 0x2f, 0x77, 0x5f, 0x37, 0x8d, 0xc1, 0xd3, 0x50, 0x66, 0x15, 0x58,
	0xa1, 0xdb, 0x1a, 0x5c, 0x33, 0x1e, 0x5d, 0x2b, 0xb1, 0x06, 0x67, 0x69, 0x1f, 0xbd, 0xa0, 0x2d,
	0x95, 0x91, 0x7f, 0x3e, 0xfc, 0x04, 0x00, 0x00, 0xff, 0xff, 0x78, 0x7a, 0x60, 0x35, 0xcd, 0x01,
	0x00, 0x00,
}