// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.1
// source: channel_rule.proto

package grpc

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

type ChannelsRuleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelID int64 `protobuf:"varint,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"`
	TagID     int64 `protobuf:"varint,2,opt,name=TagID,proto3" json:"TagID,omitempty"`
}

func (x *ChannelsRuleRequest) Reset() {
	*x = ChannelsRuleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_channel_rule_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChannelsRuleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChannelsRuleRequest) ProtoMessage() {}

func (x *ChannelsRuleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_channel_rule_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChannelsRuleRequest.ProtoReflect.Descriptor instead.
func (*ChannelsRuleRequest) Descriptor() ([]byte, []int) {
	return file_channel_rule_proto_rawDescGZIP(), []int{0}
}

func (x *ChannelsRuleRequest) GetChannelID() int64 {
	if x != nil {
		return x.ChannelID
	}
	return 0
}

func (x *ChannelsRuleRequest) GetTagID() int64 {
	if x != nil {
		return x.TagID
	}
	return 0
}

type ChannelRuleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelID int64 `protobuf:"varint,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"`
	TagID     int64 `protobuf:"varint,2,opt,name=TagID,proto3" json:"TagID,omitempty"`
}

func (x *ChannelRuleRequest) Reset() {
	*x = ChannelRuleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_channel_rule_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChannelRuleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChannelRuleRequest) ProtoMessage() {}

func (x *ChannelRuleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_channel_rule_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChannelRuleRequest.ProtoReflect.Descriptor instead.
func (*ChannelRuleRequest) Descriptor() ([]byte, []int) {
	return file_channel_rule_proto_rawDescGZIP(), []int{1}
}

func (x *ChannelRuleRequest) GetChannelID() int64 {
	if x != nil {
		return x.ChannelID
	}
	return 0
}

func (x *ChannelRuleRequest) GetTagID() int64 {
	if x != nil {
		return x.TagID
	}
	return 0
}

type ChannelRuleReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rules []*ChannelRule `protobuf:"bytes,1,rep,name=Rules,proto3" json:"Rules,omitempty"`
}

func (x *ChannelRuleReply) Reset() {
	*x = ChannelRuleReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_channel_rule_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChannelRuleReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChannelRuleReply) ProtoMessage() {}

func (x *ChannelRuleReply) ProtoReflect() protoreflect.Message {
	mi := &file_channel_rule_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChannelRuleReply.ProtoReflect.Descriptor instead.
func (*ChannelRuleReply) Descriptor() ([]byte, []int) {
	return file_channel_rule_proto_rawDescGZIP(), []int{2}
}

func (x *ChannelRuleReply) GetRules() []*ChannelRule {
	if x != nil {
		return x.Rules
	}
	return nil
}

type ChannelRule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelID int64   `protobuf:"varint,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"`
	TagID     int64   `protobuf:"varint,2,opt,name=TagID,proto3" json:"TagID,omitempty"`
	Price     float32 `protobuf:"fixed32,3,opt,name=Price,proto3" json:"Price,omitempty"`
	IsActive  bool    `protobuf:"varint,4,opt,name=IsActive,proto3" json:"IsActive,omitempty"`
}

func (x *ChannelRule) Reset() {
	*x = ChannelRule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_channel_rule_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChannelRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChannelRule) ProtoMessage() {}

func (x *ChannelRule) ProtoReflect() protoreflect.Message {
	mi := &file_channel_rule_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChannelRule.ProtoReflect.Descriptor instead.
func (*ChannelRule) Descriptor() ([]byte, []int) {
	return file_channel_rule_proto_rawDescGZIP(), []int{3}
}

func (x *ChannelRule) GetChannelID() int64 {
	if x != nil {
		return x.ChannelID
	}
	return 0
}

func (x *ChannelRule) GetTagID() int64 {
	if x != nil {
		return x.TagID
	}
	return 0
}

func (x *ChannelRule) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *ChannelRule) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

var File_channel_rule_proto protoreflect.FileDescriptor

var file_channel_rule_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x49, 0x0a, 0x13, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73,
	0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x61, 0x67,
	0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x54, 0x61, 0x67, 0x49, 0x44, 0x22,
	0x48, 0x0a, 0x12, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x61, 0x67, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x54, 0x61, 0x67, 0x49, 0x44, 0x22, 0x36, 0x0a, 0x10, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x22, 0x0a,
	0x05, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x52, 0x75, 0x6c, 0x65,
	0x73, 0x22, 0x73, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x75, 0x6c, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x44, 0x12, 0x14,
	0x0a, 0x05, 0x54, 0x61, 0x67, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x54,
	0x61, 0x67, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x73,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x49, 0x73,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x32, 0x87, 0x01, 0x0a, 0x12, 0x43, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3b, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x54, 0x61, 0x67, 0x73, 0x12,
	0x14, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52,
	0x75, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x0d, 0x47, 0x65,
	0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x54, 0x61, 0x67, 0x12, 0x13, 0x2e, 0x43, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0c, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x22, 0x00,
	0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x2e, 0x70, 0x61, 0x79, 0x67, 0x65, 0x61, 0x72, 0x2e,
	0x69, 0x72, 0x2f, 0x67, 0x69, 0x66, 0x74, 0x69, 0x6e, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_channel_rule_proto_rawDescOnce sync.Once
	file_channel_rule_proto_rawDescData = file_channel_rule_proto_rawDesc
)

func file_channel_rule_proto_rawDescGZIP() []byte {
	file_channel_rule_proto_rawDescOnce.Do(func() {
		file_channel_rule_proto_rawDescData = protoimpl.X.CompressGZIP(file_channel_rule_proto_rawDescData)
	})
	return file_channel_rule_proto_rawDescData
}

var file_channel_rule_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_channel_rule_proto_goTypes = []interface{}{
	(*ChannelsRuleRequest)(nil), // 0: ChannelsRuleRequest
	(*ChannelRuleRequest)(nil),  // 1: ChannelRuleRequest
	(*ChannelRuleReply)(nil),    // 2: ChannelRuleReply
	(*ChannelRule)(nil),         // 3: ChannelRule
}
var file_channel_rule_proto_depIdxs = []int32{
	3, // 0: ChannelRuleReply.Rules:type_name -> ChannelRule
	0, // 1: ChannelRuleService.GetChannelTags:input_type -> ChannelsRuleRequest
	1, // 2: ChannelRuleService.GetChannelTag:input_type -> ChannelRuleRequest
	2, // 3: ChannelRuleService.GetChannelTags:output_type -> ChannelRuleReply
	3, // 4: ChannelRuleService.GetChannelTag:output_type -> ChannelRule
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_channel_rule_proto_init() }
func file_channel_rule_proto_init() {
	if File_channel_rule_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_channel_rule_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChannelsRuleRequest); i {
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
		file_channel_rule_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChannelRuleRequest); i {
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
		file_channel_rule_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChannelRuleReply); i {
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
		file_channel_rule_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChannelRule); i {
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
			RawDescriptor: file_channel_rule_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_channel_rule_proto_goTypes,
		DependencyIndexes: file_channel_rule_proto_depIdxs,
		MessageInfos:      file_channel_rule_proto_msgTypes,
	}.Build()
	File_channel_rule_proto = out.File
	file_channel_rule_proto_rawDesc = nil
	file_channel_rule_proto_goTypes = nil
	file_channel_rule_proto_depIdxs = nil
}
