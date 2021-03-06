// Code generated by protoc-gen-go. DO NOT EDIT.
// source: storage_query.proto

package mgmt

import (
	fmt "fmt"
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

type BioHealthReq struct {
	DevUuid              string   `protobuf:"bytes,1,opt,name=dev_uuid,json=devUuid,proto3" json:"dev_uuid,omitempty"`
	TgtId                string   `protobuf:"bytes,2,opt,name=tgt_id,json=tgtId,proto3" json:"tgt_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BioHealthReq) Reset()         { *m = BioHealthReq{} }
func (m *BioHealthReq) String() string { return proto.CompactTextString(m) }
func (*BioHealthReq) ProtoMessage()    {}
func (*BioHealthReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_d87a8d20722a9416, []int{0}
}

func (m *BioHealthReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BioHealthReq.Unmarshal(m, b)
}
func (m *BioHealthReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BioHealthReq.Marshal(b, m, deterministic)
}
func (m *BioHealthReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BioHealthReq.Merge(m, src)
}
func (m *BioHealthReq) XXX_Size() int {
	return xxx_messageInfo_BioHealthReq.Size(m)
}
func (m *BioHealthReq) XXX_DiscardUnknown() {
	xxx_messageInfo_BioHealthReq.DiscardUnknown(m)
}

var xxx_messageInfo_BioHealthReq proto.InternalMessageInfo

func (m *BioHealthReq) GetDevUuid() string {
	if m != nil {
		return m.DevUuid
	}
	return ""
}

func (m *BioHealthReq) GetTgtId() string {
	if m != nil {
		return m.TgtId
	}
	return ""
}

type BioHealthResp struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	DevUuid              string   `protobuf:"bytes,2,opt,name=dev_uuid,json=devUuid,proto3" json:"dev_uuid,omitempty"`
	ErrorCount           uint64   `protobuf:"varint,3,opt,name=error_count,json=errorCount,proto3" json:"error_count,omitempty"`
	Temperature          uint32   `protobuf:"varint,4,opt,name=temperature,proto3" json:"temperature,omitempty"`
	MediaErrors          uint64   `protobuf:"varint,5,opt,name=media_errors,json=mediaErrors,proto3" json:"media_errors,omitempty"`
	ReadErrs             uint32   `protobuf:"varint,6,opt,name=read_errs,json=readErrs,proto3" json:"read_errs,omitempty"`
	WriteErrs            uint32   `protobuf:"varint,7,opt,name=write_errs,json=writeErrs,proto3" json:"write_errs,omitempty"`
	UnmapErrs            uint32   `protobuf:"varint,8,opt,name=unmap_errs,json=unmapErrs,proto3" json:"unmap_errs,omitempty"`
	ChecksumErrs         uint32   `protobuf:"varint,9,opt,name=checksum_errs,json=checksumErrs,proto3" json:"checksum_errs,omitempty"`
	Temp                 bool     `protobuf:"varint,10,opt,name=temp,proto3" json:"temp,omitempty"`
	Spare                bool     `protobuf:"varint,11,opt,name=spare,proto3" json:"spare,omitempty"`
	Readonly             bool     `protobuf:"varint,12,opt,name=readonly,proto3" json:"readonly,omitempty"`
	DeviceReliability    bool     `protobuf:"varint,13,opt,name=device_reliability,json=deviceReliability,proto3" json:"device_reliability,omitempty"`
	VolatileMemory       bool     `protobuf:"varint,14,opt,name=volatile_memory,json=volatileMemory,proto3" json:"volatile_memory,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BioHealthResp) Reset()         { *m = BioHealthResp{} }
func (m *BioHealthResp) String() string { return proto.CompactTextString(m) }
func (*BioHealthResp) ProtoMessage()    {}
func (*BioHealthResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_d87a8d20722a9416, []int{1}
}

func (m *BioHealthResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BioHealthResp.Unmarshal(m, b)
}
func (m *BioHealthResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BioHealthResp.Marshal(b, m, deterministic)
}
func (m *BioHealthResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BioHealthResp.Merge(m, src)
}
func (m *BioHealthResp) XXX_Size() int {
	return xxx_messageInfo_BioHealthResp.Size(m)
}
func (m *BioHealthResp) XXX_DiscardUnknown() {
	xxx_messageInfo_BioHealthResp.DiscardUnknown(m)
}

var xxx_messageInfo_BioHealthResp proto.InternalMessageInfo

func (m *BioHealthResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *BioHealthResp) GetDevUuid() string {
	if m != nil {
		return m.DevUuid
	}
	return ""
}

func (m *BioHealthResp) GetErrorCount() uint64 {
	if m != nil {
		return m.ErrorCount
	}
	return 0
}

func (m *BioHealthResp) GetTemperature() uint32 {
	if m != nil {
		return m.Temperature
	}
	return 0
}

func (m *BioHealthResp) GetMediaErrors() uint64 {
	if m != nil {
		return m.MediaErrors
	}
	return 0
}

func (m *BioHealthResp) GetReadErrs() uint32 {
	if m != nil {
		return m.ReadErrs
	}
	return 0
}

func (m *BioHealthResp) GetWriteErrs() uint32 {
	if m != nil {
		return m.WriteErrs
	}
	return 0
}

func (m *BioHealthResp) GetUnmapErrs() uint32 {
	if m != nil {
		return m.UnmapErrs
	}
	return 0
}

func (m *BioHealthResp) GetChecksumErrs() uint32 {
	if m != nil {
		return m.ChecksumErrs
	}
	return 0
}

func (m *BioHealthResp) GetTemp() bool {
	if m != nil {
		return m.Temp
	}
	return false
}

func (m *BioHealthResp) GetSpare() bool {
	if m != nil {
		return m.Spare
	}
	return false
}

func (m *BioHealthResp) GetReadonly() bool {
	if m != nil {
		return m.Readonly
	}
	return false
}

func (m *BioHealthResp) GetDeviceReliability() bool {
	if m != nil {
		return m.DeviceReliability
	}
	return false
}

func (m *BioHealthResp) GetVolatileMemory() bool {
	if m != nil {
		return m.VolatileMemory
	}
	return false
}

type SmdDevReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SmdDevReq) Reset()         { *m = SmdDevReq{} }
func (m *SmdDevReq) String() string { return proto.CompactTextString(m) }
func (*SmdDevReq) ProtoMessage()    {}
func (*SmdDevReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_d87a8d20722a9416, []int{2}
}

func (m *SmdDevReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SmdDevReq.Unmarshal(m, b)
}
func (m *SmdDevReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SmdDevReq.Marshal(b, m, deterministic)
}
func (m *SmdDevReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SmdDevReq.Merge(m, src)
}
func (m *SmdDevReq) XXX_Size() int {
	return xxx_messageInfo_SmdDevReq.Size(m)
}
func (m *SmdDevReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SmdDevReq.DiscardUnknown(m)
}

var xxx_messageInfo_SmdDevReq proto.InternalMessageInfo

type SmdDevResp struct {
	Status               int32                `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Devices              []*SmdDevResp_Device `protobuf:"bytes,2,rep,name=devices,proto3" json:"devices,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SmdDevResp) Reset()         { *m = SmdDevResp{} }
func (m *SmdDevResp) String() string { return proto.CompactTextString(m) }
func (*SmdDevResp) ProtoMessage()    {}
func (*SmdDevResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_d87a8d20722a9416, []int{3}
}

func (m *SmdDevResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SmdDevResp.Unmarshal(m, b)
}
func (m *SmdDevResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SmdDevResp.Marshal(b, m, deterministic)
}
func (m *SmdDevResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SmdDevResp.Merge(m, src)
}
func (m *SmdDevResp) XXX_Size() int {
	return xxx_messageInfo_SmdDevResp.Size(m)
}
func (m *SmdDevResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SmdDevResp.DiscardUnknown(m)
}

var xxx_messageInfo_SmdDevResp proto.InternalMessageInfo

func (m *SmdDevResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *SmdDevResp) GetDevices() []*SmdDevResp_Device {
	if m != nil {
		return m.Devices
	}
	return nil
}

type SmdDevResp_Device struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	TgtIds               []int32  `protobuf:"varint,2,rep,packed,name=tgt_ids,json=tgtIds,proto3" json:"tgt_ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SmdDevResp_Device) Reset()         { *m = SmdDevResp_Device{} }
func (m *SmdDevResp_Device) String() string { return proto.CompactTextString(m) }
func (*SmdDevResp_Device) ProtoMessage()    {}
func (*SmdDevResp_Device) Descriptor() ([]byte, []int) {
	return fileDescriptor_d87a8d20722a9416, []int{3, 0}
}

func (m *SmdDevResp_Device) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SmdDevResp_Device.Unmarshal(m, b)
}
func (m *SmdDevResp_Device) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SmdDevResp_Device.Marshal(b, m, deterministic)
}
func (m *SmdDevResp_Device) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SmdDevResp_Device.Merge(m, src)
}
func (m *SmdDevResp_Device) XXX_Size() int {
	return xxx_messageInfo_SmdDevResp_Device.Size(m)
}
func (m *SmdDevResp_Device) XXX_DiscardUnknown() {
	xxx_messageInfo_SmdDevResp_Device.DiscardUnknown(m)
}

var xxx_messageInfo_SmdDevResp_Device proto.InternalMessageInfo

func (m *SmdDevResp_Device) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *SmdDevResp_Device) GetTgtIds() []int32 {
	if m != nil {
		return m.TgtIds
	}
	return nil
}

type SmdPoolReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SmdPoolReq) Reset()         { *m = SmdPoolReq{} }
func (m *SmdPoolReq) String() string { return proto.CompactTextString(m) }
func (*SmdPoolReq) ProtoMessage()    {}
func (*SmdPoolReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_d87a8d20722a9416, []int{4}
}

func (m *SmdPoolReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SmdPoolReq.Unmarshal(m, b)
}
func (m *SmdPoolReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SmdPoolReq.Marshal(b, m, deterministic)
}
func (m *SmdPoolReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SmdPoolReq.Merge(m, src)
}
func (m *SmdPoolReq) XXX_Size() int {
	return xxx_messageInfo_SmdPoolReq.Size(m)
}
func (m *SmdPoolReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SmdPoolReq.DiscardUnknown(m)
}

var xxx_messageInfo_SmdPoolReq proto.InternalMessageInfo

type SmdPoolResp struct {
	Status               int32               `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Pools                []*SmdPoolResp_Pool `protobuf:"bytes,2,rep,name=pools,proto3" json:"pools,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *SmdPoolResp) Reset()         { *m = SmdPoolResp{} }
func (m *SmdPoolResp) String() string { return proto.CompactTextString(m) }
func (*SmdPoolResp) ProtoMessage()    {}
func (*SmdPoolResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_d87a8d20722a9416, []int{5}
}

func (m *SmdPoolResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SmdPoolResp.Unmarshal(m, b)
}
func (m *SmdPoolResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SmdPoolResp.Marshal(b, m, deterministic)
}
func (m *SmdPoolResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SmdPoolResp.Merge(m, src)
}
func (m *SmdPoolResp) XXX_Size() int {
	return xxx_messageInfo_SmdPoolResp.Size(m)
}
func (m *SmdPoolResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SmdPoolResp.DiscardUnknown(m)
}

var xxx_messageInfo_SmdPoolResp proto.InternalMessageInfo

func (m *SmdPoolResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *SmdPoolResp) GetPools() []*SmdPoolResp_Pool {
	if m != nil {
		return m.Pools
	}
	return nil
}

type SmdPoolResp_Pool struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	TgtIds               []int32  `protobuf:"varint,2,rep,packed,name=tgt_ids,json=tgtIds,proto3" json:"tgt_ids,omitempty"`
	Blobs                []uint64 `protobuf:"varint,3,rep,packed,name=blobs,proto3" json:"blobs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SmdPoolResp_Pool) Reset()         { *m = SmdPoolResp_Pool{} }
func (m *SmdPoolResp_Pool) String() string { return proto.CompactTextString(m) }
func (*SmdPoolResp_Pool) ProtoMessage()    {}
func (*SmdPoolResp_Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_d87a8d20722a9416, []int{5, 0}
}

func (m *SmdPoolResp_Pool) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SmdPoolResp_Pool.Unmarshal(m, b)
}
func (m *SmdPoolResp_Pool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SmdPoolResp_Pool.Marshal(b, m, deterministic)
}
func (m *SmdPoolResp_Pool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SmdPoolResp_Pool.Merge(m, src)
}
func (m *SmdPoolResp_Pool) XXX_Size() int {
	return xxx_messageInfo_SmdPoolResp_Pool.Size(m)
}
func (m *SmdPoolResp_Pool) XXX_DiscardUnknown() {
	xxx_messageInfo_SmdPoolResp_Pool.DiscardUnknown(m)
}

var xxx_messageInfo_SmdPoolResp_Pool proto.InternalMessageInfo

func (m *SmdPoolResp_Pool) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *SmdPoolResp_Pool) GetTgtIds() []int32 {
	if m != nil {
		return m.TgtIds
	}
	return nil
}

func (m *SmdPoolResp_Pool) GetBlobs() []uint64 {
	if m != nil {
		return m.Blobs
	}
	return nil
}

type DevStateReq struct {
	DevUuid              string   `protobuf:"bytes,1,opt,name=dev_uuid,json=devUuid,proto3" json:"dev_uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DevStateReq) Reset()         { *m = DevStateReq{} }
func (m *DevStateReq) String() string { return proto.CompactTextString(m) }
func (*DevStateReq) ProtoMessage()    {}
func (*DevStateReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_d87a8d20722a9416, []int{6}
}

func (m *DevStateReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DevStateReq.Unmarshal(m, b)
}
func (m *DevStateReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DevStateReq.Marshal(b, m, deterministic)
}
func (m *DevStateReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DevStateReq.Merge(m, src)
}
func (m *DevStateReq) XXX_Size() int {
	return xxx_messageInfo_DevStateReq.Size(m)
}
func (m *DevStateReq) XXX_DiscardUnknown() {
	xxx_messageInfo_DevStateReq.DiscardUnknown(m)
}

var xxx_messageInfo_DevStateReq proto.InternalMessageInfo

func (m *DevStateReq) GetDevUuid() string {
	if m != nil {
		return m.DevUuid
	}
	return ""
}

type DevStateResp struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	DevUuid              string   `protobuf:"bytes,2,opt,name=dev_uuid,json=devUuid,proto3" json:"dev_uuid,omitempty"`
	DevState             string   `protobuf:"bytes,3,opt,name=dev_state,json=devState,proto3" json:"dev_state,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DevStateResp) Reset()         { *m = DevStateResp{} }
func (m *DevStateResp) String() string { return proto.CompactTextString(m) }
func (*DevStateResp) ProtoMessage()    {}
func (*DevStateResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_d87a8d20722a9416, []int{7}
}

func (m *DevStateResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DevStateResp.Unmarshal(m, b)
}
func (m *DevStateResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DevStateResp.Marshal(b, m, deterministic)
}
func (m *DevStateResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DevStateResp.Merge(m, src)
}
func (m *DevStateResp) XXX_Size() int {
	return xxx_messageInfo_DevStateResp.Size(m)
}
func (m *DevStateResp) XXX_DiscardUnknown() {
	xxx_messageInfo_DevStateResp.DiscardUnknown(m)
}

var xxx_messageInfo_DevStateResp proto.InternalMessageInfo

func (m *DevStateResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *DevStateResp) GetDevUuid() string {
	if m != nil {
		return m.DevUuid
	}
	return ""
}

func (m *DevStateResp) GetDevState() string {
	if m != nil {
		return m.DevState
	}
	return ""
}

func init() {
	proto.RegisterType((*BioHealthReq)(nil), "mgmt.BioHealthReq")
	proto.RegisterType((*BioHealthResp)(nil), "mgmt.BioHealthResp")
	proto.RegisterType((*SmdDevReq)(nil), "mgmt.SmdDevReq")
	proto.RegisterType((*SmdDevResp)(nil), "mgmt.SmdDevResp")
	proto.RegisterType((*SmdDevResp_Device)(nil), "mgmt.SmdDevResp.Device")
	proto.RegisterType((*SmdPoolReq)(nil), "mgmt.SmdPoolReq")
	proto.RegisterType((*SmdPoolResp)(nil), "mgmt.SmdPoolResp")
	proto.RegisterType((*SmdPoolResp_Pool)(nil), "mgmt.SmdPoolResp.Pool")
	proto.RegisterType((*DevStateReq)(nil), "mgmt.DevStateReq")
	proto.RegisterType((*DevStateResp)(nil), "mgmt.DevStateResp")
}

func init() {
	proto.RegisterFile("storage_query.proto", fileDescriptor_d87a8d20722a9416)
}

var fileDescriptor_d87a8d20722a9416 = []byte{
	// 514 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0xdd, 0x8e, 0x93, 0x40,
	0x14, 0x0e, 0x5b, 0xa0, 0xe5, 0x40, 0xd7, 0x38, 0xea, 0xee, 0xb8, 0x1b, 0x23, 0xe2, 0x85, 0x5c,
	0x68, 0x13, 0x35, 0xde, 0x1b, 0xed, 0x26, 0xee, 0x85, 0x89, 0x61, 0xe3, 0xad, 0x84, 0x96, 0x93,
	0x2e, 0x11, 0x3a, 0x74, 0x66, 0xc0, 0xf4, 0x25, 0x7c, 0x04, 0x1f, 0xc7, 0xe7, 0x32, 0x73, 0x46,
	0xda, 0xee, 0x4d, 0x13, 0xbd, 0x9b, 0xef, 0x67, 0x3e, 0xbe, 0x30, 0xe7, 0xc0, 0x03, 0xa5, 0x85,
	0x2c, 0x56, 0x98, 0x6f, 0x3a, 0x94, 0xdb, 0x59, 0x2b, 0x85, 0x16, 0xcc, 0x6d, 0x56, 0x8d, 0x4e,
	0xde, 0x43, 0xf4, 0xa1, 0x12, 0x9f, 0xb0, 0xa8, 0xf5, 0x6d, 0x86, 0x1b, 0xf6, 0x18, 0x26, 0x25,
	0xf6, 0x79, 0xd7, 0x55, 0x25, 0x77, 0x62, 0x27, 0x0d, 0xb2, 0x71, 0x89, 0xfd, 0xd7, 0xae, 0x2a,
	0xd9, 0x23, 0xf0, 0xf5, 0x4a, 0xe7, 0x55, 0xc9, 0x4f, 0x48, 0xf0, 0xf4, 0x4a, 0x5f, 0x97, 0xc9,
	0xef, 0x11, 0x4c, 0x0f, 0x22, 0x54, 0xcb, 0xce, 0xc0, 0x57, 0xba, 0xd0, 0x9d, 0xa2, 0x04, 0x2f,
	0xfb, 0x8b, 0xee, 0x64, 0x9f, 0xdc, 0xcd, 0x7e, 0x0a, 0x21, 0x4a, 0x29, 0x64, 0xbe, 0x14, 0xdd,
	0x5a, 0xf3, 0x51, 0xec, 0xa4, 0x6e, 0x06, 0x44, 0x7d, 0x34, 0x0c, 0x8b, 0x21, 0xd4, 0xd8, 0xb4,
	0x28, 0x0b, 0xdd, 0x49, 0xe4, 0x6e, 0xec, 0xa4, 0xd3, 0xec, 0x90, 0x62, 0xcf, 0x20, 0x6a, 0xb0,
	0xac, 0x8a, 0x9c, 0x6e, 0x29, 0xee, 0x51, 0x46, 0x48, 0xdc, 0x15, 0x51, 0xec, 0x12, 0x02, 0x89,
	0x45, 0x69, 0x1c, 0x8a, 0xfb, 0x14, 0x31, 0x31, 0xc4, 0x95, 0x94, 0x8a, 0x3d, 0x01, 0xf8, 0x21,
	0x2b, 0x8d, 0x56, 0x1d, 0x93, 0x1a, 0x10, 0x33, 0xc8, 0xdd, 0xba, 0x29, 0x5a, 0x2b, 0x4f, 0xac,
	0x4c, 0x0c, 0xc9, 0xcf, 0x61, 0xba, 0xbc, 0xc5, 0xe5, 0x77, 0xd5, 0x35, 0xd6, 0x11, 0x90, 0x23,
	0x1a, 0x48, 0x32, 0x31, 0x70, 0x4d, 0x63, 0x0e, 0xb1, 0x93, 0x4e, 0x32, 0x3a, 0xb3, 0x87, 0xe0,
	0xa9, 0xb6, 0x90, 0xc8, 0x43, 0x22, 0x2d, 0x60, 0x17, 0x40, 0xc5, 0xc4, 0xba, 0xde, 0xf2, 0x88,
	0x84, 0x1d, 0x66, 0xaf, 0x80, 0x95, 0xd8, 0x57, 0x4b, 0xcc, 0x25, 0xd6, 0x55, 0xb1, 0xa8, 0xea,
	0x4a, 0x6f, 0xf9, 0x94, 0x5c, 0xf7, 0xad, 0x92, 0xed, 0x05, 0xf6, 0x02, 0xee, 0xf5, 0xa2, 0x2e,
	0x74, 0x55, 0x63, 0xde, 0x60, 0x23, 0xe4, 0x96, 0x9f, 0x92, 0xf7, 0x74, 0xa0, 0x3f, 0x13, 0x9b,
	0x84, 0x10, 0xdc, 0x34, 0xe5, 0x1c, 0xfb, 0x0c, 0x37, 0xc9, 0x4f, 0x07, 0x60, 0x40, 0x47, 0x9e,
	0xf4, 0x35, 0x8c, 0xed, 0x17, 0x15, 0x3f, 0x89, 0x47, 0x69, 0xf8, 0xe6, 0x7c, 0x66, 0xc6, 0x6a,
	0xb6, 0xbf, 0x3a, 0x9b, 0xdb, 0x46, 0x83, 0xef, 0xe2, 0x1d, 0xf8, 0x96, 0x32, 0xbf, 0xe3, 0x60,
	0xce, 0xe8, 0xcc, 0xce, 0x61, 0x6c, 0x87, 0xcc, 0x06, 0x7a, 0x99, 0x4f, 0x53, 0xa6, 0x92, 0x88,
	0xfa, 0x7c, 0x11, 0xa2, 0x36, 0xf5, 0x7e, 0x39, 0x10, 0xee, 0xe0, 0x91, 0x7e, 0x2f, 0xc1, 0x6b,
	0x85, 0xa8, 0x87, 0x76, 0x67, 0xbb, 0x76, 0xc3, 0xcd, 0x19, 0x1d, 0xac, 0xe9, 0xe2, 0x1a, 0x5c,
	0x03, 0xff, 0xa9, 0x98, 0x79, 0xc0, 0x45, 0x2d, 0x16, 0x8a, 0x8f, 0xe2, 0x51, 0xea, 0x66, 0x16,
	0x24, 0x29, 0x84, 0x73, 0xec, 0x6f, 0x74, 0xa1, 0xf1, 0xf8, 0x5a, 0x25, 0xdf, 0x20, 0xda, 0x3b,
	0xff, 0x6f, 0x7b, 0x2e, 0x21, 0x30, 0x92, 0x31, 0x22, 0xed, 0x4e, 0x90, 0x19, 0x2f, 0x65, 0x2e,
	0x7c, 0x5a, 0xf7, 0xb7, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x67, 0xf4, 0xf7, 0xe5, 0x05, 0x04,
	0x00, 0x00,
}
