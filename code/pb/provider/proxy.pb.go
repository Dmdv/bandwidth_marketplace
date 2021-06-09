// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: pb/provider/proto/proxy.proto

package provider

import (
	magma "github.com/0chain/bandwidth_marketplace/code/pb/magma"
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

type NewSessionBillingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionID        string `protobuf:"bytes,1,opt,name=sessionID,proto3" json:"sessionID,omitempty"`
	UserID           string `protobuf:"bytes,2,opt,name=userID,proto3" json:"userID,omitempty"`
	ConsumerID       string `protobuf:"bytes,3,opt,name=consumerID,proto3" json:"consumerID,omitempty"`
	AccessPointID    string `protobuf:"bytes,4,opt,name=accessPointID,proto3" json:"accessPointID,omitempty"`
	AcknowledgmentID string `protobuf:"bytes,5,opt,name=acknowledgmentID,proto3" json:"acknowledgmentID,omitempty"`
}

func (x *NewSessionBillingRequest) Reset() {
	*x = NewSessionBillingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_provider_proto_proxy_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewSessionBillingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewSessionBillingRequest) ProtoMessage() {}

func (x *NewSessionBillingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_provider_proto_proxy_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewSessionBillingRequest.ProtoReflect.Descriptor instead.
func (*NewSessionBillingRequest) Descriptor() ([]byte, []int) {
	return file_pb_provider_proto_proxy_proto_rawDescGZIP(), []int{0}
}

func (x *NewSessionBillingRequest) GetSessionID() string {
	if x != nil {
		return x.SessionID
	}
	return ""
}

func (x *NewSessionBillingRequest) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *NewSessionBillingRequest) GetConsumerID() string {
	if x != nil {
		return x.ConsumerID
	}
	return ""
}

func (x *NewSessionBillingRequest) GetAccessPointID() string {
	if x != nil {
		return x.AccessPointID
	}
	return ""
}

func (x *NewSessionBillingRequest) GetAcknowledgmentID() string {
	if x != nil {
		return x.AcknowledgmentID
	}
	return ""
}

type NewSessionBillingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NewSessionBillingResponse) Reset() {
	*x = NewSessionBillingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_provider_proto_proxy_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewSessionBillingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewSessionBillingResponse) ProtoMessage() {}

func (x *NewSessionBillingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_provider_proto_proxy_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewSessionBillingResponse.ProtoReflect.Descriptor instead.
func (*NewSessionBillingResponse) Descriptor() ([]byte, []int) {
	return file_pb_provider_proto_proxy_proto_rawDescGZIP(), []int{1}
}

type ForwardUsageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ForwardUsageResponse) Reset() {
	*x = ForwardUsageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_provider_proto_proxy_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForwardUsageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForwardUsageResponse) ProtoMessage() {}

func (x *ForwardUsageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_provider_proto_proxy_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForwardUsageResponse.ProtoReflect.Descriptor instead.
func (*ForwardUsageResponse) Descriptor() ([]byte, []int) {
	return file_pb_provider_proto_proxy_proto_rawDescGZIP(), []int{2}
}

var File_pb_provider_proto_proxy_proto protoreflect.FileDescriptor

var file_pb_provider_proto_proxy_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x12, 0x7a, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x1a, 0x1a, 0x70, 0x62, 0x2f, 0x6d, 0x61, 0x67, 0x6d, 0x61, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x61, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xc2, 0x01, 0x0a, 0x18, 0x4e, 0x65, 0x77, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x69,
	0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x49, 0x44,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72,
	0x49, 0x44, 0x12, 0x24, 0x0a, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x69, 0x6e,
	0x74, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x2a, 0x0a, 0x10, 0x61, 0x63, 0x6b, 0x6e,
	0x6f, 0x77, 0x6c, 0x65, 0x64, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x10, 0x61, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c, 0x65, 0x64, 0x67, 0x6d, 0x65,
	0x6e, 0x74, 0x49, 0x44, 0x22, 0x1b, 0x0a, 0x19, 0x4e, 0x65, 0x77, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x42, 0x69, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x16, 0x0a, 0x14, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x55, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xcf, 0x01, 0x0a, 0x05, 0x50, 0x72,
	0x6f, 0x78, 0x79, 0x12, 0x70, 0x0a, 0x11, 0x4e, 0x65, 0x77, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x42, 0x69, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x12, 0x2c, 0x2e, 0x7a, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x4e, 0x65,
	0x77, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x69, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x7a, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e,
	0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x4e, 0x65, 0x77, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x69, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x54, 0x0a, 0x0c, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x2e, 0x7a, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x70,
	0x62, 0x2e, 0x6d, 0x61, 0x67, 0x6d, 0x61, 0x2e, 0x55, 0x73, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x1a, 0x28, 0x2e, 0x7a, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x62, 0x2e, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x55, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3a, 0x5a, 0x38, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x30, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x2f, 0x62, 0x61, 0x6e, 0x64, 0x77, 0x69, 0x64, 0x74, 0x68, 0x5f, 0x6d, 0x61, 0x72, 0x6b, 0x65,
	0x74, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x62, 0x2f, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_provider_proto_proxy_proto_rawDescOnce sync.Once
	file_pb_provider_proto_proxy_proto_rawDescData = file_pb_provider_proto_proxy_proto_rawDesc
)

func file_pb_provider_proto_proxy_proto_rawDescGZIP() []byte {
	file_pb_provider_proto_proxy_proto_rawDescOnce.Do(func() {
		file_pb_provider_proto_proxy_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_provider_proto_proxy_proto_rawDescData)
	})
	return file_pb_provider_proto_proxy_proto_rawDescData
}

var file_pb_provider_proto_proxy_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pb_provider_proto_proxy_proto_goTypes = []interface{}{
	(*NewSessionBillingRequest)(nil),  // 0: zchain.pb.provider.NewSessionBillingRequest
	(*NewSessionBillingResponse)(nil), // 1: zchain.pb.provider.NewSessionBillingResponse
	(*ForwardUsageResponse)(nil),      // 2: zchain.pb.provider.ForwardUsageResponse
	(*magma.UsageData)(nil),           // 3: zchain.pb.magma.UsageData
}
var file_pb_provider_proto_proxy_proto_depIdxs = []int32{
	0, // 0: zchain.pb.provider.Proxy.NewSessionBilling:input_type -> zchain.pb.provider.NewSessionBillingRequest
	3, // 1: zchain.pb.provider.Proxy.ForwardUsage:input_type -> zchain.pb.magma.UsageData
	1, // 2: zchain.pb.provider.Proxy.NewSessionBilling:output_type -> zchain.pb.provider.NewSessionBillingResponse
	2, // 3: zchain.pb.provider.Proxy.ForwardUsage:output_type -> zchain.pb.provider.ForwardUsageResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_provider_proto_proxy_proto_init() }
func file_pb_provider_proto_proxy_proto_init() {
	if File_pb_provider_proto_proxy_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_provider_proto_proxy_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewSessionBillingRequest); i {
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
		file_pb_provider_proto_proxy_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewSessionBillingResponse); i {
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
		file_pb_provider_proto_proxy_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForwardUsageResponse); i {
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
			RawDescriptor: file_pb_provider_proto_proxy_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_provider_proto_proxy_proto_goTypes,
		DependencyIndexes: file_pb_provider_proto_proxy_proto_depIdxs,
		MessageInfos:      file_pb_provider_proto_proxy_proto_msgTypes,
	}.Build()
	File_pb_provider_proto_proxy_proto = out.File
	file_pb_provider_proto_proxy_proto_rawDesc = nil
	file_pb_provider_proto_proxy_proto_goTypes = nil
	file_pb_provider_proto_proxy_proto_depIdxs = nil
}