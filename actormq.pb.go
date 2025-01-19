// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.28.3
// source: actormq.proto

package actormq

import (
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

type PID struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Address       string                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	ID            string                 `protobuf:"bytes,2,opt,name=ID,proto3" json:"ID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PID) Reset() {
	*x = PID{}
	mi := &file_actormq_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PID) ProtoMessage() {}

func (x *PID) ProtoReflect() protoreflect.Message {
	mi := &file_actormq_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PID.ProtoReflect.Descriptor instead.
func (*PID) Descriptor() ([]byte, []int) {
	return file_actormq_proto_rawDescGZIP(), []int{0}
}

func (x *PID) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *PID) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

type RegisterNode struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterNode) Reset() {
	*x = RegisterNode{}
	mi := &file_actormq_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterNode) ProtoMessage() {}

func (x *RegisterNode) ProtoReflect() protoreflect.Message {
	mi := &file_actormq_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterNode.ProtoReflect.Descriptor instead.
func (*RegisterNode) Descriptor() ([]byte, []int) {
	return file_actormq_proto_rawDescGZIP(), []int{1}
}

type ActiveNodes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Nodes         []*PID                 `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ActiveNodes) Reset() {
	*x = ActiveNodes{}
	mi := &file_actormq_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ActiveNodes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActiveNodes) ProtoMessage() {}

func (x *ActiveNodes) ProtoReflect() protoreflect.Message {
	mi := &file_actormq_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActiveNodes.ProtoReflect.Descriptor instead.
func (*ActiveNodes) Descriptor() ([]byte, []int) {
	return file_actormq_proto_rawDescGZIP(), []int{2}
}

func (x *ActiveNodes) GetNodes() []*PID {
	if x != nil {
		return x.Nodes
	}
	return nil
}

var File_actormq_proto protoreflect.FileDescriptor

var file_actormq_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x6d, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x6d, 0x71, 0x22, 0x2f, 0x0a, 0x03, 0x50, 0x49, 0x44, 0x12,
	0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x22, 0x0e, 0x0a, 0x0c, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x22, 0x31, 0x0a, 0x0b, 0x41, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x22, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x6d,
	0x71, 0x2e, 0x50, 0x49, 0x44, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x42, 0x20, 0x5a, 0x1e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x72, 0x6f, 0x79, 0x67,
	0x69, 0x6c, 0x6d, 0x61, 0x6e, 0x30, 0x2f, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x6d, 0x71, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_actormq_proto_rawDescOnce sync.Once
	file_actormq_proto_rawDescData = file_actormq_proto_rawDesc
)

func file_actormq_proto_rawDescGZIP() []byte {
	file_actormq_proto_rawDescOnce.Do(func() {
		file_actormq_proto_rawDescData = protoimpl.X.CompressGZIP(file_actormq_proto_rawDescData)
	})
	return file_actormq_proto_rawDescData
}

var file_actormq_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_actormq_proto_goTypes = []any{
	(*PID)(nil),          // 0: actormq.PID
	(*RegisterNode)(nil), // 1: actormq.RegisterNode
	(*ActiveNodes)(nil),  // 2: actormq.ActiveNodes
}
var file_actormq_proto_depIdxs = []int32{
	0, // 0: actormq.ActiveNodes.nodes:type_name -> actormq.PID
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_actormq_proto_init() }
func file_actormq_proto_init() {
	if File_actormq_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_actormq_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_actormq_proto_goTypes,
		DependencyIndexes: file_actormq_proto_depIdxs,
		MessageInfos:      file_actormq_proto_msgTypes,
	}.Build()
	File_actormq_proto = out.File
	file_actormq_proto_rawDesc = nil
	file_actormq_proto_goTypes = nil
	file_actormq_proto_depIdxs = nil
}
