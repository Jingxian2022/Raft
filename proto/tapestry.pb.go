//
//  Brown University, CS138, Spring 2023
//
//  Purpose: Defines the Tapestry RPC protocol using Google's Protocol Buffers
//  syntax. See https://developers.google.com/protocol-buffers for more details.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.12
// source: proto/tapestry.proto

package proto

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Ok struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *Ok) Reset() {
	*x = Ok{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ok) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ok) ProtoMessage() {}

func (x *Ok) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ok.ProtoReflect.Descriptor instead.
func (*Ok) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{0}
}

func (x *Ok) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type IdMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`        // ID of key to find or of local node
	Level int32  `protobuf:"varint,2,opt,name=level,proto3" json:"level,omitempty"` // Level in routing table to begin searching at
}

func (x *IdMsg) Reset() {
	*x = IdMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdMsg) ProtoMessage() {}

func (x *IdMsg) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdMsg.ProtoReflect.Descriptor instead.
func (*IdMsg) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{1}
}

func (x *IdMsg) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *IdMsg) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

type RootMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Next     string   `protobuf:"bytes,1,opt,name=next,proto3" json:"next,omitempty"`         // ID of root node
	ToRemove []string `protobuf:"bytes,2,rep,name=toRemove,proto3" json:"toRemove,omitempty"` // ID of nodes to remove
}

func (x *RootMsg) Reset() {
	*x = RootMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RootMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RootMsg) ProtoMessage() {}

func (x *RootMsg) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RootMsg.ProtoReflect.Descriptor instead.
func (*RootMsg) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{2}
}

func (x *RootMsg) GetNext() string {
	if x != nil {
		return x.Next
	}
	return ""
}

func (x *RootMsg) GetToRemove() []string {
	if x != nil {
		return x.ToRemove
	}
	return nil
}

type DataBlob struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"` // Addresses of replica group
	Key  string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`   // Replica group ID
}

func (x *DataBlob) Reset() {
	*x = DataBlob{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataBlob) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataBlob) ProtoMessage() {}

func (x *DataBlob) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataBlob.ProtoReflect.Descriptor instead.
func (*DataBlob) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{3}
}

func (x *DataBlob) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *DataBlob) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type TapestryKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"` // Replica group ID
}

func (x *TapestryKey) Reset() {
	*x = TapestryKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TapestryKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TapestryKey) ProtoMessage() {}

func (x *TapestryKey) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TapestryKey.ProtoReflect.Descriptor instead.
func (*TapestryKey) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{4}
}

func (x *TapestryKey) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type Registration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromNode string `protobuf:"bytes,1,opt,name=fromNode,proto3" json:"fromNode,omitempty"` // ID of node storing the blob for the key
	Key      string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`           // Replica group ID
}

func (x *Registration) Reset() {
	*x = Registration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Registration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Registration) ProtoMessage() {}

func (x *Registration) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Registration.ProtoReflect.Descriptor instead.
func (*Registration) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{5}
}

func (x *Registration) GetFromNode() string {
	if x != nil {
		return x.FromNode
	}
	return ""
}

func (x *Registration) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type NodeMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // ID of node
}

func (x *NodeMsg) Reset() {
	*x = NodeMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeMsg) ProtoMessage() {}

func (x *NodeMsg) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeMsg.ProtoReflect.Descriptor instead.
func (*NodeMsg) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{6}
}

func (x *NodeMsg) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type FetchedLocations struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsRoot bool     `protobuf:"varint,1,opt,name=isRoot,proto3" json:"isRoot,omitempty"` // Whether this node is the root node for the queried key
	Values []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`  // All the nodes storing the blob for the queried key
}

func (x *FetchedLocations) Reset() {
	*x = FetchedLocations{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchedLocations) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchedLocations) ProtoMessage() {}

func (x *FetchedLocations) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchedLocations.ProtoReflect.Descriptor instead.
func (*FetchedLocations) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{7}
}

func (x *FetchedLocations) GetIsRoot() bool {
	if x != nil {
		return x.IsRoot
	}
	return false
}

func (x *FetchedLocations) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type Neighbors struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Neighbors []string `protobuf:"bytes,1,rep,name=neighbors,proto3" json:"neighbors,omitempty"` // IDs of neighbors
}

func (x *Neighbors) Reset() {
	*x = Neighbors{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Neighbors) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Neighbors) ProtoMessage() {}

func (x *Neighbors) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Neighbors.ProtoReflect.Descriptor instead.
func (*Neighbors) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{8}
}

func (x *Neighbors) GetNeighbors() []string {
	if x != nil {
		return x.Neighbors
	}
	return nil
}

type MulticastRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NewNode string `protobuf:"bytes,1,opt,name=newNode,proto3" json:"newNode,omitempty"` // ID of new node
	Level   int32  `protobuf:"varint,2,opt,name=level,proto3" json:"level,omitempty"`    // Row in routing table to propagate request to
}

func (x *MulticastRequest) Reset() {
	*x = MulticastRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MulticastRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MulticastRequest) ProtoMessage() {}

func (x *MulticastRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MulticastRequest.ProtoReflect.Descriptor instead.
func (*MulticastRequest) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{9}
}

func (x *MulticastRequest) GetNewNode() string {
	if x != nil {
		return x.NewNode
	}
	return ""
}

func (x *MulticastRequest) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

type TransferData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From string                `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`                                                                                         // ID of node data is being trasferred from
	Data map[string]*Neighbors `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // Data being transferred
}

func (x *TransferData) Reset() {
	*x = TransferData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferData) ProtoMessage() {}

func (x *TransferData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferData.ProtoReflect.Descriptor instead.
func (*TransferData) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{10}
}

func (x *TransferData) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *TransferData) GetData() map[string]*Neighbors {
	if x != nil {
		return x.Data
	}
	return nil
}

type BackpointerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From  string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`    // ID of from node
	Level int32  `protobuf:"varint,2,opt,name=level,proto3" json:"level,omitempty"` // Level to get backpointers at
}

func (x *BackpointerRequest) Reset() {
	*x = BackpointerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackpointerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackpointerRequest) ProtoMessage() {}

func (x *BackpointerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackpointerRequest.ProtoReflect.Descriptor instead.
func (*BackpointerRequest) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{11}
}

func (x *BackpointerRequest) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *BackpointerRequest) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

type LeaveNotification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From        string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`               // ID of node that is leaving
	Replacement string `protobuf:"bytes,2,opt,name=replacement,proto3" json:"replacement,omitempty"` // ID of suitable alternative node for routing table
}

func (x *LeaveNotification) Reset() {
	*x = LeaveNotification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tapestry_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaveNotification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaveNotification) ProtoMessage() {}

func (x *LeaveNotification) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tapestry_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaveNotification.ProtoReflect.Descriptor instead.
func (*LeaveNotification) Descriptor() ([]byte, []int) {
	return file_proto_tapestry_proto_rawDescGZIP(), []int{12}
}

func (x *LeaveNotification) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *LeaveNotification) GetReplacement() string {
	if x != nil {
		return x.Replacement
	}
	return ""
}

var File_proto_tapestry_proto protoreflect.FileDescriptor

var file_proto_tapestry_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x61, 0x70, 0x65, 0x73, 0x74, 0x72, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x14, 0x0a,
	0x02, 0x4f, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x02, 0x6f, 0x6b, 0x22, 0x2d, 0x0a, 0x05, 0x49, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x65, 0x76,
	0x65, 0x6c, 0x22, 0x39, 0x0a, 0x07, 0x52, 0x6f, 0x6f, 0x74, 0x4d, 0x73, 0x67, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x65, 0x78,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x6f, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x08, 0x74, 0x6f, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x22, 0x30, 0x0a,
	0x08, 0x44, 0x61, 0x74, 0x61, 0x42, 0x6c, 0x6f, 0x62, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22,
	0x1f, 0x0a, 0x0b, 0x54, 0x61, 0x70, 0x65, 0x73, 0x74, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x22, 0x3c, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1a, 0x0a, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x4e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x19,
	0x0a, 0x07, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x73, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x42, 0x0a, 0x10, 0x46, 0x65, 0x74,
	0x63, 0x68, 0x65, 0x64, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x69, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69,
	0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x29, 0x0a,
	0x09, 0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x65,
	0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6e,
	0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x22, 0x42, 0x0a, 0x10, 0x4d, 0x75, 0x6c, 0x74,
	0x69, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x6e, 0x65, 0x77, 0x4e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e,
	0x65, 0x77, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x22, 0xa0, 0x01, 0x0a,
	0x0c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a,
	0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f,
	0x6d, 0x12, 0x31, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72,
	0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x1a, 0x49, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x26, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x65, 0x69, 0x67, 0x68,
	0x62, 0x6f, 0x72, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x3e, 0x0a, 0x12, 0x42, 0x61, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76,
	0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x22,
	0x49, 0x0a, 0x11, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x70, 0x6c,
	0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72,
	0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x32, 0x80, 0x05, 0x0a, 0x0b, 0x54,
	0x61, 0x70, 0x65, 0x73, 0x74, 0x72, 0x79, 0x52, 0x50, 0x43, 0x12, 0x2a, 0x0a, 0x08, 0x46, 0x69,
	0x6e, 0x64, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49,
	0x64, 0x4d, 0x73, 0x67, 0x1a, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x6f, 0x6f,
	0x74, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x09, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4f, 0x6b, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x05, 0x46, 0x65, 0x74, 0x63, 0x68, 0x12, 0x12, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61, 0x70, 0x65, 0x73, 0x74, 0x72, 0x79, 0x4b, 0x65,
	0x79, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x65,
	0x64, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x00, 0x12, 0x2d, 0x0a, 0x07,
	0x41, 0x64, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x73, 0x67, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x22, 0x00, 0x12, 0x2f, 0x0a, 0x0e, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x42, 0x61, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x10, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x1a,
	0x09, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x6b, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x10,
	0x41, 0x64, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74,
	0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x22, 0x00, 0x12, 0x2c, 0x0a,
	0x08, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x09,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x6b, 0x22, 0x00, 0x12, 0x2d, 0x0a, 0x0e, 0x41,
	0x64, 0x64, 0x42, 0x61, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x0e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x73, 0x67, 0x1a, 0x09, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x6b, 0x22, 0x00, 0x12, 0x30, 0x0a, 0x11, 0x52, 0x65,
	0x6d, 0x6f, 0x76, 0x65, 0x42, 0x61, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x12,
	0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x73, 0x67, 0x1a,
	0x09, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x6b, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x0f,
	0x47, 0x65, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x73, 0x12,
	0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x22, 0x00, 0x12, 0x34,
	0x0a, 0x0b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x12, 0x18, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x09, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4f, 0x6b, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x0e, 0x42, 0x6c, 0x6f, 0x62, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x46, 0x65, 0x74, 0x63, 0x68, 0x12, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54,
	0x61, 0x70, 0x65, 0x73, 0x74, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x42, 0x6c, 0x6f, 0x62, 0x22, 0x00, 0x42, 0x0e, 0x5a,
	0x0c, 0x6d, 0x6f, 0x64, 0x69, 0x73, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_tapestry_proto_rawDescOnce sync.Once
	file_proto_tapestry_proto_rawDescData = file_proto_tapestry_proto_rawDesc
)

func file_proto_tapestry_proto_rawDescGZIP() []byte {
	file_proto_tapestry_proto_rawDescOnce.Do(func() {
		file_proto_tapestry_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_tapestry_proto_rawDescData)
	})
	return file_proto_tapestry_proto_rawDescData
}

var file_proto_tapestry_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_proto_tapestry_proto_goTypes = []interface{}{
	(*Ok)(nil),                 // 0: proto.Ok
	(*IdMsg)(nil),              // 1: proto.IdMsg
	(*RootMsg)(nil),            // 2: proto.RootMsg
	(*DataBlob)(nil),           // 3: proto.DataBlob
	(*TapestryKey)(nil),        // 4: proto.TapestryKey
	(*Registration)(nil),       // 5: proto.Registration
	(*NodeMsg)(nil),            // 6: proto.NodeMsg
	(*FetchedLocations)(nil),   // 7: proto.FetchedLocations
	(*Neighbors)(nil),          // 8: proto.Neighbors
	(*MulticastRequest)(nil),   // 9: proto.MulticastRequest
	(*TransferData)(nil),       // 10: proto.TransferData
	(*BackpointerRequest)(nil), // 11: proto.BackpointerRequest
	(*LeaveNotification)(nil),  // 12: proto.LeaveNotification
	nil,                        // 13: proto.TransferData.DataEntry
}
var file_proto_tapestry_proto_depIdxs = []int32{
	13, // 0: proto.TransferData.data:type_name -> proto.TransferData.DataEntry
	8,  // 1: proto.TransferData.DataEntry.value:type_name -> proto.Neighbors
	1,  // 2: proto.TapestryRPC.FindRoot:input_type -> proto.IdMsg
	5,  // 3: proto.TapestryRPC.Register:input_type -> proto.Registration
	4,  // 4: proto.TapestryRPC.Fetch:input_type -> proto.TapestryKey
	6,  // 5: proto.TapestryRPC.AddNode:input_type -> proto.NodeMsg
	8,  // 6: proto.TapestryRPC.RemoveBadNodes:input_type -> proto.Neighbors
	9,  // 7: proto.TapestryRPC.AddNodeMulticast:input_type -> proto.MulticastRequest
	10, // 8: proto.TapestryRPC.Transfer:input_type -> proto.TransferData
	6,  // 9: proto.TapestryRPC.AddBackpointer:input_type -> proto.NodeMsg
	6,  // 10: proto.TapestryRPC.RemoveBackpointer:input_type -> proto.NodeMsg
	11, // 11: proto.TapestryRPC.GetBackpointers:input_type -> proto.BackpointerRequest
	12, // 12: proto.TapestryRPC.NotifyLeave:input_type -> proto.LeaveNotification
	4,  // 13: proto.TapestryRPC.BlobStoreFetch:input_type -> proto.TapestryKey
	2,  // 14: proto.TapestryRPC.FindRoot:output_type -> proto.RootMsg
	0,  // 15: proto.TapestryRPC.Register:output_type -> proto.Ok
	7,  // 16: proto.TapestryRPC.Fetch:output_type -> proto.FetchedLocations
	8,  // 17: proto.TapestryRPC.AddNode:output_type -> proto.Neighbors
	0,  // 18: proto.TapestryRPC.RemoveBadNodes:output_type -> proto.Ok
	8,  // 19: proto.TapestryRPC.AddNodeMulticast:output_type -> proto.Neighbors
	0,  // 20: proto.TapestryRPC.Transfer:output_type -> proto.Ok
	0,  // 21: proto.TapestryRPC.AddBackpointer:output_type -> proto.Ok
	0,  // 22: proto.TapestryRPC.RemoveBackpointer:output_type -> proto.Ok
	8,  // 23: proto.TapestryRPC.GetBackpointers:output_type -> proto.Neighbors
	0,  // 24: proto.TapestryRPC.NotifyLeave:output_type -> proto.Ok
	3,  // 25: proto.TapestryRPC.BlobStoreFetch:output_type -> proto.DataBlob
	14, // [14:26] is the sub-list for method output_type
	2,  // [2:14] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_proto_tapestry_proto_init() }
func file_proto_tapestry_proto_init() {
	if File_proto_tapestry_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_tapestry_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ok); i {
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
		file_proto_tapestry_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdMsg); i {
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
		file_proto_tapestry_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RootMsg); i {
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
		file_proto_tapestry_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataBlob); i {
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
		file_proto_tapestry_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TapestryKey); i {
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
		file_proto_tapestry_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Registration); i {
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
		file_proto_tapestry_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeMsg); i {
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
		file_proto_tapestry_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FetchedLocations); i {
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
		file_proto_tapestry_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Neighbors); i {
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
		file_proto_tapestry_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MulticastRequest); i {
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
		file_proto_tapestry_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferData); i {
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
		file_proto_tapestry_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackpointerRequest); i {
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
		file_proto_tapestry_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaveNotification); i {
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
			RawDescriptor: file_proto_tapestry_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_tapestry_proto_goTypes,
		DependencyIndexes: file_proto_tapestry_proto_depIdxs,
		MessageInfos:      file_proto_tapestry_proto_msgTypes,
	}.Build()
	File_proto_tapestry_proto = out.File
	file_proto_tapestry_proto_rawDesc = nil
	file_proto_tapestry_proto_goTypes = nil
	file_proto_tapestry_proto_depIdxs = nil
}
