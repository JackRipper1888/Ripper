// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: peer.proto

package models

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type HeartBeat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerId    []byte `protobuf:"bytes,1,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	TimeStamp int64  `protobuf:"varint,2,opt,name=timeStamp,proto3" json:"timeStamp,omitempty"`
}

func (x *HeartBeat) Reset() {
	*x = HeartBeat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartBeat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartBeat) ProtoMessage() {}

func (x *HeartBeat) ProtoReflect() protoreflect.Message {
	mi := &file_peer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartBeat.ProtoReflect.Descriptor instead.
func (*HeartBeat) Descriptor() ([]byte, []int) {
	return file_peer_proto_rawDescGZIP(), []int{0}
}

func (x *HeartBeat) GetPeerId() []byte {
	if x != nil {
		return x.PeerId
	}
	return nil
}

func (x *HeartBeat) GetTimeStamp() int64 {
	if x != nil {
		return x.TimeStamp
	}
	return 0
}

type FindNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerId    []byte `protobuf:"bytes,1,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	PeerCount int32  `protobuf:"varint,2,opt,name=peer_count,json=peerCount,proto3" json:"peer_count,omitempty"`
}

func (x *FindNode) Reset() {
	*x = FindNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindNode) ProtoMessage() {}

func (x *FindNode) ProtoReflect() protoreflect.Message {
	mi := &file_peer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindNode.ProtoReflect.Descriptor instead.
func (*FindNode) Descriptor() ([]byte, []int) {
	return file_peer_proto_rawDescGZIP(), []int{1}
}

func (x *FindNode) GetPeerId() []byte {
	if x != nil {
		return x.PeerId
	}
	return nil
}

func (x *FindNode) GetPeerCount() int32 {
	if x != nil {
		return x.PeerCount
	}
	return 0
}

type FindNodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerId   []byte      `protobuf:"bytes,1,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	Peerlist []*PeerInfo `protobuf:"bytes,2,rep,name=peerlist,proto3" json:"peerlist,omitempty"`
}

func (x *FindNodeResponse) Reset() {
	*x = FindNodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindNodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindNodeResponse) ProtoMessage() {}

func (x *FindNodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_peer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindNodeResponse.ProtoReflect.Descriptor instead.
func (*FindNodeResponse) Descriptor() ([]byte, []int) {
	return file_peer_proto_rawDescGZIP(), []int{2}
}

func (x *FindNodeResponse) GetPeerId() []byte {
	if x != nil {
		return x.PeerId
	}
	return nil
}

func (x *FindNodeResponse) GetPeerlist() []*PeerInfo {
	if x != nil {
		return x.Peerlist
	}
	return nil
}

type PeerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerId    []byte `protobuf:"bytes,1,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	Addr      string `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	TimeStamp int64  `protobuf:"varint,3,opt,name=timeStamp,proto3" json:"timeStamp,omitempty"`
}

func (x *PeerInfo) Reset() {
	*x = PeerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerInfo) ProtoMessage() {}

func (x *PeerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_peer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerInfo.ProtoReflect.Descriptor instead.
func (*PeerInfo) Descriptor() ([]byte, []int) {
	return file_peer_proto_rawDescGZIP(), []int{3}
}

func (x *PeerInfo) GetPeerId() []byte {
	if x != nil {
		return x.PeerId
	}
	return nil
}

func (x *PeerInfo) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *PeerInfo) GetTimeStamp() int64 {
	if x != nil {
		return x.TimeStamp
	}
	return 0
}

type FindValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key  []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Leve int32  `protobuf:"varint,2,opt,name=leve,proto3" json:"leve,omitempty"`
}

func (x *FindValue) Reset() {
	*x = FindValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindValue) ProtoMessage() {}

func (x *FindValue) ProtoReflect() protoreflect.Message {
	mi := &file_peer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindValue.ProtoReflect.Descriptor instead.
func (*FindValue) Descriptor() ([]byte, []int) {
	return file_peer_proto_rawDescGZIP(), []int{4}
}

func (x *FindValue) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *FindValue) GetLeve() int32 {
	if x != nil {
		return x.Leve
	}
	return 0
}

type FindValueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key          []byte      `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Leve         int32       `protobuf:"varint,2,opt,name=leve,proto3" json:"leve,omitempty"`
	Peerlist     []*PeerInfo `protobuf:"bytes,3,rep,name=peerlist,proto3" json:"peerlist,omitempty"`
	Nearpeerlist []*PeerInfo `protobuf:"bytes,4,rep,name=nearpeerlist,proto3" json:"nearpeerlist,omitempty"`
}

func (x *FindValueResponse) Reset() {
	*x = FindValueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peer_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindValueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindValueResponse) ProtoMessage() {}

func (x *FindValueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_peer_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindValueResponse.ProtoReflect.Descriptor instead.
func (*FindValueResponse) Descriptor() ([]byte, []int) {
	return file_peer_proto_rawDescGZIP(), []int{5}
}

func (x *FindValueResponse) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *FindValueResponse) GetLeve() int32 {
	if x != nil {
		return x.Leve
	}
	return 0
}

func (x *FindValueResponse) GetPeerlist() []*PeerInfo {
	if x != nil {
		return x.Peerlist
	}
	return nil
}

func (x *FindValueResponse) GetNearpeerlist() []*PeerInfo {
	if x != nil {
		return x.Nearpeerlist
	}
	return nil
}

var File_peer_proto protoreflect.FileDescriptor

var file_peer_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x70, 0x65, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x22, 0x42, 0x0a, 0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x42, 0x0a, 0x08, 0x46, 0x69, 0x6e, 0x64,
	0x4e, 0x6f, 0x64, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x70, 0x65, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x59, 0x0a, 0x10,
	0x46, 0x69, 0x6e, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x08, 0x70, 0x65, 0x65,
	0x72, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x70,
	0x65, 0x65, 0x72, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x55, 0x0a, 0x08, 0x50, 0x65, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x61, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x31,
	0x0a, 0x09, 0x46, 0x69, 0x6e, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a,
	0x04, 0x6c, 0x65, 0x76, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x6c, 0x65, 0x76,
	0x65, 0x22, 0x9d, 0x01, 0x0a, 0x11, 0x46, 0x69, 0x6e, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x65, 0x76,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x6c, 0x65, 0x76, 0x65, 0x12, 0x2c, 0x0a,
	0x08, 0x70, 0x65, 0x65, 0x72, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x08, 0x70, 0x65, 0x65, 0x72, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x0c, 0x6e,
	0x65, 0x61, 0x72, 0x70, 0x65, 0x65, 0x72, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x0c, 0x6e, 0x65, 0x61, 0x72, 0x70, 0x65, 0x65, 0x72, 0x6c, 0x69, 0x73,
	0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_peer_proto_rawDescOnce sync.Once
	file_peer_proto_rawDescData = file_peer_proto_rawDesc
)

func file_peer_proto_rawDescGZIP() []byte {
	file_peer_proto_rawDescOnce.Do(func() {
		file_peer_proto_rawDescData = protoimpl.X.CompressGZIP(file_peer_proto_rawDescData)
	})
	return file_peer_proto_rawDescData
}

var file_peer_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_peer_proto_goTypes = []interface{}{
	(*HeartBeat)(nil),         // 0: models.HeartBeat
	(*FindNode)(nil),          // 1: models.FindNode
	(*FindNodeResponse)(nil),  // 2: models.FindNodeResponse
	(*PeerInfo)(nil),          // 3: models.PeerInfo
	(*FindValue)(nil),         // 4: models.FindValue
	(*FindValueResponse)(nil), // 5: models.FindValueResponse
}
var file_peer_proto_depIdxs = []int32{
	3, // 0: models.FindNodeResponse.peerlist:type_name -> models.PeerInfo
	3, // 1: models.FindValueResponse.peerlist:type_name -> models.PeerInfo
	3, // 2: models.FindValueResponse.nearpeerlist:type_name -> models.PeerInfo
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_peer_proto_init() }
func file_peer_proto_init() {
	if File_peer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_peer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartBeat); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_peer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindNode); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_peer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindNodeResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_peer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeerInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_peer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_peer_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindValueResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_peer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_peer_proto_goTypes,
		DependencyIndexes: file_peer_proto_depIdxs,
		MessageInfos:      file_peer_proto_msgTypes,
	}.Build()
	File_peer_proto = out.File
	file_peer_proto_rawDesc = nil
	file_peer_proto_goTypes = nil
	file_peer_proto_depIdxs = nil
}