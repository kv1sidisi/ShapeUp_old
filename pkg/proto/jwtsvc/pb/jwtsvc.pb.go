// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: jwtsvc.proto

package jwtsvc

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

type GenerateTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Uid           []byte                 `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Operation     string                 `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GenerateTokenRequest) Reset() {
	*x = GenerateTokenRequest{}
	mi := &file_jwtsvc_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenerateTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateTokenRequest) ProtoMessage() {}

func (x *GenerateTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_jwtsvc_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateTokenRequest.ProtoReflect.Descriptor instead.
func (*GenerateTokenRequest) Descriptor() ([]byte, []int) {
	return file_jwtsvc_proto_rawDescGZIP(), []int{0}
}

func (x *GenerateTokenRequest) GetUid() []byte {
	if x != nil {
		return x.Uid
	}
	return nil
}

func (x *GenerateTokenRequest) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

type GenerateTokenResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GenerateTokenResponse) Reset() {
	*x = GenerateTokenResponse{}
	mi := &file_jwtsvc_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenerateTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateTokenResponse) ProtoMessage() {}

func (x *GenerateTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_jwtsvc_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateTokenResponse.ProtoReflect.Descriptor instead.
func (*GenerateTokenResponse) Descriptor() ([]byte, []int) {
	return file_jwtsvc_proto_rawDescGZIP(), []int{1}
}

func (x *GenerateTokenResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ValidateTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ValidateTokenRequest) Reset() {
	*x = ValidateTokenRequest{}
	mi := &file_jwtsvc_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateTokenRequest) ProtoMessage() {}

func (x *ValidateTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_jwtsvc_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateTokenRequest.ProtoReflect.Descriptor instead.
func (*ValidateTokenRequest) Descriptor() ([]byte, []int) {
	return file_jwtsvc_proto_rawDescGZIP(), []int{2}
}

func (x *ValidateTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ValidateTokenResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Uid           []byte                 `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Operation     string                 `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ValidateTokenResponse) Reset() {
	*x = ValidateTokenResponse{}
	mi := &file_jwtsvc_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateTokenResponse) ProtoMessage() {}

func (x *ValidateTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_jwtsvc_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateTokenResponse.ProtoReflect.Descriptor instead.
func (*ValidateTokenResponse) Descriptor() ([]byte, []int) {
	return file_jwtsvc_proto_rawDescGZIP(), []int{3}
}

func (x *ValidateTokenResponse) GetUid() []byte {
	if x != nil {
		return x.Uid
	}
	return nil
}

func (x *ValidateTokenResponse) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

type GenerateLinkRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LinkBase      string                 `protobuf:"bytes,1,opt,name=link_base,json=linkBase,proto3" json:"link_base,omitempty"`
	Uid           []byte                 `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
	Operation     string                 `protobuf:"bytes,3,opt,name=operation,proto3" json:"operation,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GenerateLinkRequest) Reset() {
	*x = GenerateLinkRequest{}
	mi := &file_jwtsvc_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenerateLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateLinkRequest) ProtoMessage() {}

func (x *GenerateLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_jwtsvc_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateLinkRequest.ProtoReflect.Descriptor instead.
func (*GenerateLinkRequest) Descriptor() ([]byte, []int) {
	return file_jwtsvc_proto_rawDescGZIP(), []int{4}
}

func (x *GenerateLinkRequest) GetLinkBase() string {
	if x != nil {
		return x.LinkBase
	}
	return ""
}

func (x *GenerateLinkRequest) GetUid() []byte {
	if x != nil {
		return x.Uid
	}
	return nil
}

func (x *GenerateLinkRequest) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

type GenerateLinkResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Link          string                 `protobuf:"bytes,1,opt,name=link,proto3" json:"link,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GenerateLinkResponse) Reset() {
	*x = GenerateLinkResponse{}
	mi := &file_jwtsvc_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenerateLinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateLinkResponse) ProtoMessage() {}

func (x *GenerateLinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_jwtsvc_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateLinkResponse.ProtoReflect.Descriptor instead.
func (*GenerateLinkResponse) Descriptor() ([]byte, []int) {
	return file_jwtsvc_proto_rawDescGZIP(), []int{5}
}

func (x *GenerateLinkResponse) GetLink() string {
	if x != nil {
		return x.Link
	}
	return ""
}

var File_jwtsvc_proto protoreflect.FileDescriptor

var file_jwtsvc_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6a, 0x77, 0x74, 0x73, 0x76, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03,
	0x6a, 0x77, 0x74, 0x22, 0x46, 0x0a, 0x14, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x2d, 0x0a, 0x15, 0x47,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x2c, 0x0a, 0x14, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x47, 0x0a, 0x15, 0x56, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03,
	0x75, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x62, 0x0a, 0x13, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x69, 0x6e, 0x6b,
	0x5f, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x69, 0x6e,
	0x6b, 0x42, 0x61, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x2a, 0x0a, 0x14, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x69, 0x6e,
	0x6b, 0x32, 0xda, 0x01, 0x0a, 0x03, 0x4a, 0x57, 0x54, 0x12, 0x46, 0x0a, 0x0d, 0x47, 0x65, 0x6e,
	0x65, 0x72, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x19, 0x2e, 0x6a, 0x77, 0x74,
	0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6a, 0x77, 0x74, 0x2e, 0x47, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x46, 0x0a, 0x0d, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x12, 0x19, 0x2e, 0x6a, 0x77, 0x74, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e,
	0x6a, 0x77, 0x74, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0c, 0x47, 0x65, 0x6e,
	0x65, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x18, 0x2e, 0x6a, 0x77, 0x74, 0x2e,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6a, 0x77, 0x74, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x17,
	0x5a, 0x15, 0x73, 0x68, 0x61, 0x70, 0x65, 0x75, 0x70, 0x2e, 0x6a, 0x77, 0x74, 0x2e, 0x76, 0x31,
	0x3b, 0x6a, 0x77, 0x74, 0x73, 0x76, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_jwtsvc_proto_rawDescOnce sync.Once
	file_jwtsvc_proto_rawDescData = file_jwtsvc_proto_rawDesc
)

func file_jwtsvc_proto_rawDescGZIP() []byte {
	file_jwtsvc_proto_rawDescOnce.Do(func() {
		file_jwtsvc_proto_rawDescData = protoimpl.X.CompressGZIP(file_jwtsvc_proto_rawDescData)
	})
	return file_jwtsvc_proto_rawDescData
}

var file_jwtsvc_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_jwtsvc_proto_goTypes = []any{
	(*GenerateTokenRequest)(nil),  // 0: jwt.GenerateTokenRequest
	(*GenerateTokenResponse)(nil), // 1: jwt.GenerateTokenResponse
	(*ValidateTokenRequest)(nil),  // 2: jwt.ValidateTokenRequest
	(*ValidateTokenResponse)(nil), // 3: jwt.ValidateTokenResponse
	(*GenerateLinkRequest)(nil),   // 4: jwt.GenerateLinkRequest
	(*GenerateLinkResponse)(nil),  // 5: jwt.GenerateLinkResponse
}
var file_jwtsvc_proto_depIdxs = []int32{
	0, // 0: jwt.JWT.GenerateToken:input_type -> jwt.GenerateTokenRequest
	2, // 1: jwt.JWT.ValidateToken:input_type -> jwt.ValidateTokenRequest
	4, // 2: jwt.JWT.GenerateLink:input_type -> jwt.GenerateLinkRequest
	1, // 3: jwt.JWT.GenerateToken:output_type -> jwt.GenerateTokenResponse
	3, // 4: jwt.JWT.ValidateToken:output_type -> jwt.ValidateTokenResponse
	5, // 5: jwt.JWT.GenerateLink:output_type -> jwt.GenerateLinkResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_jwtsvc_proto_init() }
func file_jwtsvc_proto_init() {
	if File_jwtsvc_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_jwtsvc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_jwtsvc_proto_goTypes,
		DependencyIndexes: file_jwtsvc_proto_depIdxs,
		MessageInfos:      file_jwtsvc_proto_msgTypes,
	}.Build()
	File_jwtsvc_proto = out.File
	file_jwtsvc_proto_rawDesc = nil
	file_jwtsvc_proto_goTypes = nil
	file_jwtsvc_proto_depIdxs = nil
}
