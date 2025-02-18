// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v6.30.0--rc1
// source: core/pb/defs.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// the reader obtains associated timestamps rmin, rmax,rconf and auxiliary data Ctx, which may evolve
type Transaction struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The minimum round.
	RMin float64 `protobuf:"fixed64,1,opt,name=r_min,json=rMin,proto3" json:"r_min,omitempty"`
	// The maximum round.
	RMax float64 `protobuf:"fixed64,2,opt,name=r_max,json=rMax,proto3" json:"r_max,omitempty"`
	// Undefined confirmed round.
	RConf float64 `protobuf:"fixed64,3,opt,name=r_conf,json=rConf,proto3" json:"r_conf,omitempty"`
	// Transaction data.
	Ctx           []byte `protobuf:"bytes,4,opt,name=ctx,proto3" json:"ctx,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	mi := &file_core_pb_defs_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_core_pb_defs_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_core_pb_defs_proto_rawDescGZIP(), []int{0}
}

func (x *Transaction) GetRMin() float64 {
	if x != nil {
		return x.RMin
	}
	return 0
}

func (x *Transaction) GetRMax() float64 {
	if x != nil {
		return x.RMax
	}
	return 0
}

func (x *Transaction) GetRConf() float64 {
	if x != nil {
		return x.RConf
	}
	return 0
}

func (x *Transaction) GetCtx() []byte {
	if x != nil {
		return x.Ctx
	}
	return nil
}

// A vote is a tuple (tx,ts,sn,σ,R), where tx is a trans- action, ts is a timestamp, sn is a sequence number, σ is a signature, and R is a replica. A vote is valid if σ is a valid signature on message m = (tx, ts, sn) with respect to the public key pkR of replica R.
type Vote struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Transaction being voted on.
	Tx *Transaction `protobuf:"bytes,1,opt,name=tx,proto3" json:"tx,omitempty"`
	// Timestamp
	Ts float64 `protobuf:"fixed64,2,opt,name=ts,proto3" json:"ts,omitempty"`
	// Sequence number
	Sn int64 `protobuf:"varint,3,opt,name=sn,proto3" json:"sn,omitempty"`
	// Signature
	Sig []byte `protobuf:"bytes,4,opt,name=sig,proto3" json:"sig,omitempty"`
	// Is special heartbeat tx
	IsHeartbeat   bool `protobuf:"varint,5,opt,name=is_heartbeat,json=isHeartbeat,proto3" json:"is_heartbeat,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Vote) Reset() {
	*x = Vote{}
	mi := &file_core_pb_defs_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Vote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vote) ProtoMessage() {}

func (x *Vote) ProtoReflect() protoreflect.Message {
	mi := &file_core_pb_defs_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vote.ProtoReflect.Descriptor instead.
func (*Vote) Descriptor() ([]byte, []int) {
	return file_core_pb_defs_proto_rawDescGZIP(), []int{1}
}

func (x *Vote) GetTx() *Transaction {
	if x != nil {
		return x.Tx
	}
	return nil
}

func (x *Vote) GetTs() float64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

func (x *Vote) GetSn() int64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *Vote) GetSig() []byte {
	if x != nil {
		return x.Sig
	}
	return nil
}

func (x *Vote) GetIsHeartbeat() bool {
	if x != nil {
		return x.IsHeartbeat
	}
	return false
}

type GetVotesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SeqStart      int64                  `protobuf:"varint,1,opt,name=seqStart,proto3" json:"seqStart,omitempty"`
	SeqEnd        int64                  `protobuf:"varint,2,opt,name=seqEnd,proto3" json:"seqEnd,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetVotesRequest) Reset() {
	*x = GetVotesRequest{}
	mi := &file_core_pb_defs_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetVotesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVotesRequest) ProtoMessage() {}

func (x *GetVotesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_core_pb_defs_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVotesRequest.ProtoReflect.Descriptor instead.
func (*GetVotesRequest) Descriptor() ([]byte, []int) {
	return file_core_pb_defs_proto_rawDescGZIP(), []int{2}
}

func (x *GetVotesRequest) GetSeqStart() int64 {
	if x != nil {
		return x.SeqStart
	}
	return 0
}

func (x *GetVotesRequest) GetSeqEnd() int64 {
	if x != nil {
		return x.SeqEnd
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_core_pb_defs_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_core_pb_defs_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_core_pb_defs_proto_rawDescGZIP(), []int{3}
}

var File_core_pb_defs_proto protoreflect.FileDescriptor

var file_core_pb_defs_proto_rawDesc = string([]byte{
	0x0a, 0x12, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x62, 0x2f, 0x64, 0x65, 0x66, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x60, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x13, 0x0a, 0x05, 0x72, 0x5f, 0x6d, 0x69, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x72, 0x4d, 0x69, 0x6e, 0x12, 0x13, 0x0a, 0x05,
	0x72, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x72, 0x4d, 0x61,
	0x78, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x05, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x74, 0x78, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x63, 0x74, 0x78, 0x22, 0x7c, 0x0a, 0x04, 0x56, 0x6f,
	0x74, 0x65, 0x12, 0x1f, 0x0a, 0x02, 0x74, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x02, 0x74, 0x78, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x02, 0x74, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x73, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x03, 0x73, 0x69, 0x67, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x68, 0x65, 0x61, 0x72,
	0x74, 0x62, 0x65, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x48,
	0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x22, 0x45, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x56,
	0x6f, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73,
	0x65, 0x71, 0x53, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73,
	0x65, 0x71, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x71, 0x45, 0x6e,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x65, 0x71, 0x45, 0x6e, 0x64, 0x22,
	0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0xa8, 0x01, 0x0a, 0x0e, 0x52, 0x65, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x07, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x22, 0x0a, 0x05,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x08, 0x2e, 0x70, 0x62, 0x2e, 0x56, 0x6f, 0x74, 0x65,
	0x12, 0x24, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x56, 0x6f, 0x74, 0x65, 0x73, 0x12,
	0x09, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x08, 0x2e, 0x70, 0x62, 0x2e,
	0x56, 0x6f, 0x74, 0x65, 0x30, 0x01, 0x12, 0x2b, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x56, 0x6f, 0x74,
	0x65, 0x73, 0x12, 0x13, 0x2e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x08, 0x2e, 0x70, 0x62, 0x2e, 0x56, 0x6f, 0x74,
	0x65, 0x30, 0x01, 0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6c, 0x69, 0x61, 0x6d, 0x7a, 0x65, 0x62, 0x65, 0x64, 0x65, 0x65, 0x2f, 0x70, 0x6f,
	0x64, 0x2d, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_core_pb_defs_proto_rawDescOnce sync.Once
	file_core_pb_defs_proto_rawDescData []byte
)

func file_core_pb_defs_proto_rawDescGZIP() []byte {
	file_core_pb_defs_proto_rawDescOnce.Do(func() {
		file_core_pb_defs_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_core_pb_defs_proto_rawDesc), len(file_core_pb_defs_proto_rawDesc)))
	})
	return file_core_pb_defs_proto_rawDescData
}

var file_core_pb_defs_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_core_pb_defs_proto_goTypes = []any{
	(*Transaction)(nil),     // 0: pb.Transaction
	(*Vote)(nil),            // 1: pb.Vote
	(*GetVotesRequest)(nil), // 2: pb.GetVotesRequest
	(*Empty)(nil),           // 3: pb.Empty
}
var file_core_pb_defs_proto_depIdxs = []int32{
	0, // 0: pb.Vote.tx:type_name -> pb.Transaction
	3, // 1: pb.ReplicaService.Connect:input_type -> pb.Empty
	0, // 2: pb.ReplicaService.Write:input_type -> pb.Transaction
	3, // 3: pb.ReplicaService.StreamVotes:input_type -> pb.Empty
	2, // 4: pb.ReplicaService.GetVotes:input_type -> pb.GetVotesRequest
	3, // 5: pb.ReplicaService.Connect:output_type -> pb.Empty
	1, // 6: pb.ReplicaService.Write:output_type -> pb.Vote
	1, // 7: pb.ReplicaService.StreamVotes:output_type -> pb.Vote
	1, // 8: pb.ReplicaService.GetVotes:output_type -> pb.Vote
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_core_pb_defs_proto_init() }
func file_core_pb_defs_proto_init() {
	if File_core_pb_defs_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_core_pb_defs_proto_rawDesc), len(file_core_pb_defs_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_core_pb_defs_proto_goTypes,
		DependencyIndexes: file_core_pb_defs_proto_depIdxs,
		MessageInfos:      file_core_pb_defs_proto_msgTypes,
	}.Build()
	File_core_pb_defs_proto = out.File
	file_core_pb_defs_proto_goTypes = nil
	file_core_pb_defs_proto_depIdxs = nil
}
