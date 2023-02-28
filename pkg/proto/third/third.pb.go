// Code generated by protoc-gen-go. DO NOT EDIT.
// source: third/third.proto

package third // import "OpenIM/pkg/proto/third"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import sdkws "OpenIM/pkg/proto/sdkws"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ApplySpaceReq struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Size                 int64    `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
	Hash                 string   `protobuf:"bytes,3,opt,name=hash" json:"hash,omitempty"`
	Purpose              uint32   `protobuf:"varint,4,opt,name=purpose" json:"purpose,omitempty"`
	ContentType          string   `protobuf:"bytes,5,opt,name=contentType" json:"contentType,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ApplySpaceReq) Reset()         { *m = ApplySpaceReq{} }
func (m *ApplySpaceReq) String() string { return proto.CompactTextString(m) }
func (*ApplySpaceReq) ProtoMessage()    {}
func (*ApplySpaceReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{0}
}
func (m *ApplySpaceReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApplySpaceReq.Unmarshal(m, b)
}
func (m *ApplySpaceReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApplySpaceReq.Marshal(b, m, deterministic)
}
func (dst *ApplySpaceReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApplySpaceReq.Merge(dst, src)
}
func (m *ApplySpaceReq) XXX_Size() int {
	return xxx_messageInfo_ApplySpaceReq.Size(m)
}
func (m *ApplySpaceReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ApplySpaceReq.DiscardUnknown(m)
}

var xxx_messageInfo_ApplySpaceReq proto.InternalMessageInfo

func (m *ApplySpaceReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ApplySpaceReq) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *ApplySpaceReq) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *ApplySpaceReq) GetPurpose() uint32 {
	if m != nil {
		return m.Purpose
	}
	return 0
}

func (m *ApplySpaceReq) GetContentType() string {
	if m != nil {
		return m.ContentType
	}
	return ""
}

type ApplySpaceResp struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	Size                 int64    `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
	Put                  []string `protobuf:"bytes,3,rep,name=put" json:"put,omitempty"`
	ConfirmID            string   `protobuf:"bytes,4,opt,name=confirmID" json:"confirmID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ApplySpaceResp) Reset()         { *m = ApplySpaceResp{} }
func (m *ApplySpaceResp) String() string { return proto.CompactTextString(m) }
func (*ApplySpaceResp) ProtoMessage()    {}
func (*ApplySpaceResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{1}
}
func (m *ApplySpaceResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApplySpaceResp.Unmarshal(m, b)
}
func (m *ApplySpaceResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApplySpaceResp.Marshal(b, m, deterministic)
}
func (dst *ApplySpaceResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApplySpaceResp.Merge(dst, src)
}
func (m *ApplySpaceResp) XXX_Size() int {
	return xxx_messageInfo_ApplySpaceResp.Size(m)
}
func (m *ApplySpaceResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ApplySpaceResp.DiscardUnknown(m)
}

var xxx_messageInfo_ApplySpaceResp proto.InternalMessageInfo

func (m *ApplySpaceResp) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *ApplySpaceResp) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *ApplySpaceResp) GetPut() []string {
	if m != nil {
		return m.Put
	}
	return nil
}

func (m *ApplySpaceResp) GetConfirmID() string {
	if m != nil {
		return m.ConfirmID
	}
	return ""
}

type ConfirmSpaceReq struct {
	ConfirmID            string   `protobuf:"bytes,1,opt,name=confirmID" json:"confirmID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfirmSpaceReq) Reset()         { *m = ConfirmSpaceReq{} }
func (m *ConfirmSpaceReq) String() string { return proto.CompactTextString(m) }
func (*ConfirmSpaceReq) ProtoMessage()    {}
func (*ConfirmSpaceReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{2}
}
func (m *ConfirmSpaceReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfirmSpaceReq.Unmarshal(m, b)
}
func (m *ConfirmSpaceReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfirmSpaceReq.Marshal(b, m, deterministic)
}
func (dst *ConfirmSpaceReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfirmSpaceReq.Merge(dst, src)
}
func (m *ConfirmSpaceReq) XXX_Size() int {
	return xxx_messageInfo_ConfirmSpaceReq.Size(m)
}
func (m *ConfirmSpaceReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfirmSpaceReq.DiscardUnknown(m)
}

var xxx_messageInfo_ConfirmSpaceReq proto.InternalMessageInfo

func (m *ConfirmSpaceReq) GetConfirmID() string {
	if m != nil {
		return m.ConfirmID
	}
	return ""
}

type ConfirmSpaceResp struct {
	ConfirmID            string   `protobuf:"bytes,1,opt,name=confirmID" json:"confirmID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfirmSpaceResp) Reset()         { *m = ConfirmSpaceResp{} }
func (m *ConfirmSpaceResp) String() string { return proto.CompactTextString(m) }
func (*ConfirmSpaceResp) ProtoMessage()    {}
func (*ConfirmSpaceResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{3}
}
func (m *ConfirmSpaceResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfirmSpaceResp.Unmarshal(m, b)
}
func (m *ConfirmSpaceResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfirmSpaceResp.Marshal(b, m, deterministic)
}
func (dst *ConfirmSpaceResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfirmSpaceResp.Merge(dst, src)
}
func (m *ConfirmSpaceResp) XXX_Size() int {
	return xxx_messageInfo_ConfirmSpaceResp.Size(m)
}
func (m *ConfirmSpaceResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfirmSpaceResp.DiscardUnknown(m)
}

var xxx_messageInfo_ConfirmSpaceResp proto.InternalMessageInfo

func (m *ConfirmSpaceResp) GetConfirmID() string {
	if m != nil {
		return m.ConfirmID
	}
	return ""
}

type GetRTCInvitationInfoReq struct {
	ClientMsgID          string   `protobuf:"bytes,1,opt,name=ClientMsgID" json:"ClientMsgID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRTCInvitationInfoReq) Reset()         { *m = GetRTCInvitationInfoReq{} }
func (m *GetRTCInvitationInfoReq) String() string { return proto.CompactTextString(m) }
func (*GetRTCInvitationInfoReq) ProtoMessage()    {}
func (*GetRTCInvitationInfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{4}
}
func (m *GetRTCInvitationInfoReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRTCInvitationInfoReq.Unmarshal(m, b)
}
func (m *GetRTCInvitationInfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRTCInvitationInfoReq.Marshal(b, m, deterministic)
}
func (dst *GetRTCInvitationInfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRTCInvitationInfoReq.Merge(dst, src)
}
func (m *GetRTCInvitationInfoReq) XXX_Size() int {
	return xxx_messageInfo_GetRTCInvitationInfoReq.Size(m)
}
func (m *GetRTCInvitationInfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRTCInvitationInfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetRTCInvitationInfoReq proto.InternalMessageInfo

func (m *GetRTCInvitationInfoReq) GetClientMsgID() string {
	if m != nil {
		return m.ClientMsgID
	}
	return ""
}

type GetRTCInvitationInfoResp struct {
	InvitationInfo       *sdkws.InvitationInfo  `protobuf:"bytes,1,opt,name=invitationInfo" json:"invitationInfo,omitempty"`
	OfflinePushInfo      *sdkws.OfflinePushInfo `protobuf:"bytes,2,opt,name=offlinePushInfo" json:"offlinePushInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *GetRTCInvitationInfoResp) Reset()         { *m = GetRTCInvitationInfoResp{} }
func (m *GetRTCInvitationInfoResp) String() string { return proto.CompactTextString(m) }
func (*GetRTCInvitationInfoResp) ProtoMessage()    {}
func (*GetRTCInvitationInfoResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{5}
}
func (m *GetRTCInvitationInfoResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRTCInvitationInfoResp.Unmarshal(m, b)
}
func (m *GetRTCInvitationInfoResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRTCInvitationInfoResp.Marshal(b, m, deterministic)
}
func (dst *GetRTCInvitationInfoResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRTCInvitationInfoResp.Merge(dst, src)
}
func (m *GetRTCInvitationInfoResp) XXX_Size() int {
	return xxx_messageInfo_GetRTCInvitationInfoResp.Size(m)
}
func (m *GetRTCInvitationInfoResp) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRTCInvitationInfoResp.DiscardUnknown(m)
}

var xxx_messageInfo_GetRTCInvitationInfoResp proto.InternalMessageInfo

func (m *GetRTCInvitationInfoResp) GetInvitationInfo() *sdkws.InvitationInfo {
	if m != nil {
		return m.InvitationInfo
	}
	return nil
}

func (m *GetRTCInvitationInfoResp) GetOfflinePushInfo() *sdkws.OfflinePushInfo {
	if m != nil {
		return m.OfflinePushInfo
	}
	return nil
}

type GetRTCInvitationInfoStartAppReq struct {
	UserID               string   `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRTCInvitationInfoStartAppReq) Reset()         { *m = GetRTCInvitationInfoStartAppReq{} }
func (m *GetRTCInvitationInfoStartAppReq) String() string { return proto.CompactTextString(m) }
func (*GetRTCInvitationInfoStartAppReq) ProtoMessage()    {}
func (*GetRTCInvitationInfoStartAppReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{6}
}
func (m *GetRTCInvitationInfoStartAppReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRTCInvitationInfoStartAppReq.Unmarshal(m, b)
}
func (m *GetRTCInvitationInfoStartAppReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRTCInvitationInfoStartAppReq.Marshal(b, m, deterministic)
}
func (dst *GetRTCInvitationInfoStartAppReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRTCInvitationInfoStartAppReq.Merge(dst, src)
}
func (m *GetRTCInvitationInfoStartAppReq) XXX_Size() int {
	return xxx_messageInfo_GetRTCInvitationInfoStartAppReq.Size(m)
}
func (m *GetRTCInvitationInfoStartAppReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRTCInvitationInfoStartAppReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetRTCInvitationInfoStartAppReq proto.InternalMessageInfo

func (m *GetRTCInvitationInfoStartAppReq) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type GetRTCInvitationInfoStartAppResp struct {
	InvitationInfo       *sdkws.InvitationInfo  `protobuf:"bytes,1,opt,name=invitationInfo" json:"invitationInfo,omitempty"`
	OfflinePushInfo      *sdkws.OfflinePushInfo `protobuf:"bytes,2,opt,name=offlinePushInfo" json:"offlinePushInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *GetRTCInvitationInfoStartAppResp) Reset()         { *m = GetRTCInvitationInfoStartAppResp{} }
func (m *GetRTCInvitationInfoStartAppResp) String() string { return proto.CompactTextString(m) }
func (*GetRTCInvitationInfoStartAppResp) ProtoMessage()    {}
func (*GetRTCInvitationInfoStartAppResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{7}
}
func (m *GetRTCInvitationInfoStartAppResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRTCInvitationInfoStartAppResp.Unmarshal(m, b)
}
func (m *GetRTCInvitationInfoStartAppResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRTCInvitationInfoStartAppResp.Marshal(b, m, deterministic)
}
func (dst *GetRTCInvitationInfoStartAppResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRTCInvitationInfoStartAppResp.Merge(dst, src)
}
func (m *GetRTCInvitationInfoStartAppResp) XXX_Size() int {
	return xxx_messageInfo_GetRTCInvitationInfoStartAppResp.Size(m)
}
func (m *GetRTCInvitationInfoStartAppResp) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRTCInvitationInfoStartAppResp.DiscardUnknown(m)
}

var xxx_messageInfo_GetRTCInvitationInfoStartAppResp proto.InternalMessageInfo

func (m *GetRTCInvitationInfoStartAppResp) GetInvitationInfo() *sdkws.InvitationInfo {
	if m != nil {
		return m.InvitationInfo
	}
	return nil
}

func (m *GetRTCInvitationInfoStartAppResp) GetOfflinePushInfo() *sdkws.OfflinePushInfo {
	if m != nil {
		return m.OfflinePushInfo
	}
	return nil
}

type FcmUpdateTokenReq struct {
	Platform             int32    `protobuf:"varint,1,opt,name=Platform" json:"Platform,omitempty"`
	FcmToken             string   `protobuf:"bytes,2,opt,name=FcmToken" json:"FcmToken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FcmUpdateTokenReq) Reset()         { *m = FcmUpdateTokenReq{} }
func (m *FcmUpdateTokenReq) String() string { return proto.CompactTextString(m) }
func (*FcmUpdateTokenReq) ProtoMessage()    {}
func (*FcmUpdateTokenReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{8}
}
func (m *FcmUpdateTokenReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FcmUpdateTokenReq.Unmarshal(m, b)
}
func (m *FcmUpdateTokenReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FcmUpdateTokenReq.Marshal(b, m, deterministic)
}
func (dst *FcmUpdateTokenReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FcmUpdateTokenReq.Merge(dst, src)
}
func (m *FcmUpdateTokenReq) XXX_Size() int {
	return xxx_messageInfo_FcmUpdateTokenReq.Size(m)
}
func (m *FcmUpdateTokenReq) XXX_DiscardUnknown() {
	xxx_messageInfo_FcmUpdateTokenReq.DiscardUnknown(m)
}

var xxx_messageInfo_FcmUpdateTokenReq proto.InternalMessageInfo

func (m *FcmUpdateTokenReq) GetPlatform() int32 {
	if m != nil {
		return m.Platform
	}
	return 0
}

func (m *FcmUpdateTokenReq) GetFcmToken() string {
	if m != nil {
		return m.FcmToken
	}
	return ""
}

type FcmUpdateTokenResp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FcmUpdateTokenResp) Reset()         { *m = FcmUpdateTokenResp{} }
func (m *FcmUpdateTokenResp) String() string { return proto.CompactTextString(m) }
func (*FcmUpdateTokenResp) ProtoMessage()    {}
func (*FcmUpdateTokenResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{9}
}
func (m *FcmUpdateTokenResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FcmUpdateTokenResp.Unmarshal(m, b)
}
func (m *FcmUpdateTokenResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FcmUpdateTokenResp.Marshal(b, m, deterministic)
}
func (dst *FcmUpdateTokenResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FcmUpdateTokenResp.Merge(dst, src)
}
func (m *FcmUpdateTokenResp) XXX_Size() int {
	return xxx_messageInfo_FcmUpdateTokenResp.Size(m)
}
func (m *FcmUpdateTokenResp) XXX_DiscardUnknown() {
	xxx_messageInfo_FcmUpdateTokenResp.DiscardUnknown(m)
}

var xxx_messageInfo_FcmUpdateTokenResp proto.InternalMessageInfo

type SetAppBadgeReq struct {
	FromUserID           string   `protobuf:"bytes,1,opt,name=FromUserID" json:"FromUserID,omitempty"`
	AppUnreadCount       int32    `protobuf:"varint,2,opt,name=AppUnreadCount" json:"AppUnreadCount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetAppBadgeReq) Reset()         { *m = SetAppBadgeReq{} }
func (m *SetAppBadgeReq) String() string { return proto.CompactTextString(m) }
func (*SetAppBadgeReq) ProtoMessage()    {}
func (*SetAppBadgeReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{10}
}
func (m *SetAppBadgeReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetAppBadgeReq.Unmarshal(m, b)
}
func (m *SetAppBadgeReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetAppBadgeReq.Marshal(b, m, deterministic)
}
func (dst *SetAppBadgeReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetAppBadgeReq.Merge(dst, src)
}
func (m *SetAppBadgeReq) XXX_Size() int {
	return xxx_messageInfo_SetAppBadgeReq.Size(m)
}
func (m *SetAppBadgeReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SetAppBadgeReq.DiscardUnknown(m)
}

var xxx_messageInfo_SetAppBadgeReq proto.InternalMessageInfo

func (m *SetAppBadgeReq) GetFromUserID() string {
	if m != nil {
		return m.FromUserID
	}
	return ""
}

func (m *SetAppBadgeReq) GetAppUnreadCount() int32 {
	if m != nil {
		return m.AppUnreadCount
	}
	return 0
}

type SetAppBadgeResp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetAppBadgeResp) Reset()         { *m = SetAppBadgeResp{} }
func (m *SetAppBadgeResp) String() string { return proto.CompactTextString(m) }
func (*SetAppBadgeResp) ProtoMessage()    {}
func (*SetAppBadgeResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_third_7f0ad743baded49b, []int{11}
}
func (m *SetAppBadgeResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetAppBadgeResp.Unmarshal(m, b)
}
func (m *SetAppBadgeResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetAppBadgeResp.Marshal(b, m, deterministic)
}
func (dst *SetAppBadgeResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetAppBadgeResp.Merge(dst, src)
}
func (m *SetAppBadgeResp) XXX_Size() int {
	return xxx_messageInfo_SetAppBadgeResp.Size(m)
}
func (m *SetAppBadgeResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SetAppBadgeResp.DiscardUnknown(m)
}

var xxx_messageInfo_SetAppBadgeResp proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ApplySpaceReq)(nil), "third.ApplySpaceReq")
	proto.RegisterType((*ApplySpaceResp)(nil), "third.ApplySpaceResp")
	proto.RegisterType((*ConfirmSpaceReq)(nil), "third.ConfirmSpaceReq")
	proto.RegisterType((*ConfirmSpaceResp)(nil), "third.ConfirmSpaceResp")
	proto.RegisterType((*GetRTCInvitationInfoReq)(nil), "third.GetRTCInvitationInfoReq")
	proto.RegisterType((*GetRTCInvitationInfoResp)(nil), "third.GetRTCInvitationInfoResp")
	proto.RegisterType((*GetRTCInvitationInfoStartAppReq)(nil), "third.GetRTCInvitationInfoStartAppReq")
	proto.RegisterType((*GetRTCInvitationInfoStartAppResp)(nil), "third.GetRTCInvitationInfoStartAppResp")
	proto.RegisterType((*FcmUpdateTokenReq)(nil), "third.FcmUpdateTokenReq")
	proto.RegisterType((*FcmUpdateTokenResp)(nil), "third.FcmUpdateTokenResp")
	proto.RegisterType((*SetAppBadgeReq)(nil), "third.SetAppBadgeReq")
	proto.RegisterType((*SetAppBadgeResp)(nil), "third.SetAppBadgeResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Third service

type ThirdClient interface {
	ApplySpace(ctx context.Context, in *ApplySpaceReq, opts ...grpc.CallOption) (*ApplySpaceResp, error)
	GetRTCInvitationInfo(ctx context.Context, in *GetRTCInvitationInfoReq, opts ...grpc.CallOption) (*GetRTCInvitationInfoResp, error)
	GetRTCInvitationInfoStartApp(ctx context.Context, in *GetRTCInvitationInfoStartAppReq, opts ...grpc.CallOption) (*GetRTCInvitationInfoStartAppResp, error)
	FcmUpdateToken(ctx context.Context, in *FcmUpdateTokenReq, opts ...grpc.CallOption) (*FcmUpdateTokenResp, error)
	SetAppBadge(ctx context.Context, in *SetAppBadgeReq, opts ...grpc.CallOption) (*SetAppBadgeResp, error)
}

type thirdClient struct {
	cc *grpc.ClientConn
}

func NewThirdClient(cc *grpc.ClientConn) ThirdClient {
	return &thirdClient{cc}
}

func (c *thirdClient) ApplySpace(ctx context.Context, in *ApplySpaceReq, opts ...grpc.CallOption) (*ApplySpaceResp, error) {
	out := new(ApplySpaceResp)
	err := grpc.Invoke(ctx, "/third.third/ApplySpace", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *thirdClient) GetRTCInvitationInfo(ctx context.Context, in *GetRTCInvitationInfoReq, opts ...grpc.CallOption) (*GetRTCInvitationInfoResp, error) {
	out := new(GetRTCInvitationInfoResp)
	err := grpc.Invoke(ctx, "/third.third/GetRTCInvitationInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *thirdClient) GetRTCInvitationInfoStartApp(ctx context.Context, in *GetRTCInvitationInfoStartAppReq, opts ...grpc.CallOption) (*GetRTCInvitationInfoStartAppResp, error) {
	out := new(GetRTCInvitationInfoStartAppResp)
	err := grpc.Invoke(ctx, "/third.third/GetRTCInvitationInfoStartApp", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *thirdClient) FcmUpdateToken(ctx context.Context, in *FcmUpdateTokenReq, opts ...grpc.CallOption) (*FcmUpdateTokenResp, error) {
	out := new(FcmUpdateTokenResp)
	err := grpc.Invoke(ctx, "/third.third/FcmUpdateToken", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *thirdClient) SetAppBadge(ctx context.Context, in *SetAppBadgeReq, opts ...grpc.CallOption) (*SetAppBadgeResp, error) {
	out := new(SetAppBadgeResp)
	err := grpc.Invoke(ctx, "/third.third/SetAppBadge", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Third service

type ThirdServer interface {
	ApplySpace(context.Context, *ApplySpaceReq) (*ApplySpaceResp, error)
	GetRTCInvitationInfo(context.Context, *GetRTCInvitationInfoReq) (*GetRTCInvitationInfoResp, error)
	GetRTCInvitationInfoStartApp(context.Context, *GetRTCInvitationInfoStartAppReq) (*GetRTCInvitationInfoStartAppResp, error)
	FcmUpdateToken(context.Context, *FcmUpdateTokenReq) (*FcmUpdateTokenResp, error)
	SetAppBadge(context.Context, *SetAppBadgeReq) (*SetAppBadgeResp, error)
}

func RegisterThirdServer(s *grpc.Server, srv ThirdServer) {
	s.RegisterService(&_Third_serviceDesc, srv)
}

func _Third_ApplySpace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplySpaceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThirdServer).ApplySpace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/third.third/ApplySpace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThirdServer).ApplySpace(ctx, req.(*ApplySpaceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Third_GetRTCInvitationInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRTCInvitationInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThirdServer).GetRTCInvitationInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/third.third/GetRTCInvitationInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThirdServer).GetRTCInvitationInfo(ctx, req.(*GetRTCInvitationInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Third_GetRTCInvitationInfoStartApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRTCInvitationInfoStartAppReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThirdServer).GetRTCInvitationInfoStartApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/third.third/GetRTCInvitationInfoStartApp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThirdServer).GetRTCInvitationInfoStartApp(ctx, req.(*GetRTCInvitationInfoStartAppReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Third_FcmUpdateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FcmUpdateTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThirdServer).FcmUpdateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/third.third/FcmUpdateToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThirdServer).FcmUpdateToken(ctx, req.(*FcmUpdateTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Third_SetAppBadge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAppBadgeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThirdServer).SetAppBadge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/third.third/SetAppBadge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThirdServer).SetAppBadge(ctx, req.(*SetAppBadgeReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Third_serviceDesc = grpc.ServiceDesc{
	ServiceName: "third.third",
	HandlerType: (*ThirdServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ApplySpace",
			Handler:    _Third_ApplySpace_Handler,
		},
		{
			MethodName: "GetRTCInvitationInfo",
			Handler:    _Third_GetRTCInvitationInfo_Handler,
		},
		{
			MethodName: "GetRTCInvitationInfoStartApp",
			Handler:    _Third_GetRTCInvitationInfoStartApp_Handler,
		},
		{
			MethodName: "FcmUpdateToken",
			Handler:    _Third_FcmUpdateToken_Handler,
		},
		{
			MethodName: "SetAppBadge",
			Handler:    _Third_SetAppBadge_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "third/third.proto",
}

func init() { proto.RegisterFile("third/third.proto", fileDescriptor_third_7f0ad743baded49b) }

var fileDescriptor_third_7f0ad743baded49b = []byte{
	// 588 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x54, 0x41, 0x4f, 0xdb, 0x30,
	0x14, 0x56, 0x08, 0x65, 0xeb, 0xab, 0x28, 0x60, 0x01, 0xcb, 0x22, 0x04, 0x51, 0x0e, 0xc0, 0x85,
	0x66, 0x62, 0x27, 0xc4, 0x26, 0x0d, 0xba, 0x75, 0xaa, 0xa6, 0x0a, 0x94, 0xb6, 0xda, 0xb4, 0x5b,
	0xd6, 0x3a, 0x6d, 0xd4, 0xc6, 0x7e, 0xb3, 0x1d, 0x10, 0xfb, 0x03, 0x3b, 0xef, 0xbc, 0xe3, 0xfe,
	0xe8, 0x14, 0xb7, 0x94, 0x24, 0x0b, 0xdd, 0x8e, 0xbb, 0x58, 0x7e, 0x9f, 0xbf, 0xef, 0xf9, 0x3d,
	0xbf, 0xe7, 0x07, 0x5b, 0x6a, 0x1c, 0x89, 0xa1, 0xa7, 0xd7, 0x06, 0x0a, 0xae, 0x38, 0xa9, 0x68,
	0xc3, 0x3e, 0xba, 0x42, 0xca, 0x4e, 0xda, 0x9d, 0x93, 0x2e, 0x15, 0x37, 0x54, 0x78, 0x38, 0x19,
	0x79, 0x9a, 0xe0, 0xc9, 0xe1, 0xe4, 0x56, 0x7a, 0xb7, 0x72, 0xc6, 0x77, 0xbf, 0x1b, 0xb0, 0x7e,
	0x81, 0x38, 0xbd, 0xeb, 0x62, 0x30, 0xa0, 0x3e, 0xfd, 0x4a, 0x08, 0xac, 0xb2, 0x20, 0xa6, 0x96,
	0xe1, 0x18, 0xc7, 0x55, 0x5f, 0xef, 0x53, 0x4c, 0x46, 0xdf, 0xa8, 0xb5, 0xe2, 0x18, 0xc7, 0xa6,
	0xaf, 0xf7, 0x29, 0x36, 0x0e, 0xe4, 0xd8, 0x32, 0x67, 0xbc, 0x74, 0x4f, 0x2c, 0x78, 0x82, 0x89,
	0x40, 0x2e, 0xa9, 0xb5, 0xea, 0x18, 0xc7, 0xeb, 0xfe, 0xbd, 0x49, 0x1c, 0xa8, 0x0d, 0x38, 0x53,
	0x94, 0xa9, 0xde, 0x1d, 0x52, 0xab, 0xa2, 0x45, 0x59, 0xc8, 0x0d, 0xa1, 0x9e, 0x0d, 0x44, 0x22,
	0xd9, 0x04, 0x33, 0x11, 0xd3, 0x79, 0x20, 0xe9, 0xb6, 0x34, 0x8e, 0x4d, 0x30, 0x31, 0x51, 0x96,
	0xe9, 0x98, 0x29, 0x0b, 0x13, 0x45, 0xf6, 0xa0, 0x3a, 0xe0, 0x2c, 0x8c, 0x44, 0xdc, 0x7e, 0xab,
	0xe3, 0xa8, 0xfa, 0x0f, 0x80, 0xeb, 0xc1, 0x46, 0x73, 0x66, 0x2c, 0x52, 0xce, 0x09, 0x8c, 0xa2,
	0xe0, 0x05, 0x6c, 0xe6, 0x05, 0x12, 0xff, 0xa2, 0x38, 0x87, 0x67, 0xef, 0xa9, 0xf2, 0x7b, 0xcd,
	0x36, 0xbb, 0x89, 0x54, 0xa0, 0x22, 0xce, 0xda, 0x2c, 0xe4, 0xe9, 0x55, 0x0e, 0xd4, 0x9a, 0xd3,
	0x88, 0x32, 0xd5, 0x91, 0xa3, 0x85, 0x34, 0x0b, 0xb9, 0x3f, 0x0d, 0xb0, 0xca, 0xd5, 0x12, 0xc9,
	0x6b, 0xa8, 0x47, 0x39, 0x54, 0x7b, 0xa8, 0x9d, 0xee, 0x34, 0x74, 0x5d, 0x1b, 0x05, 0x49, 0x81,
	0x4c, 0xde, 0xc0, 0x06, 0x0f, 0xc3, 0x69, 0xc4, 0xe8, 0x75, 0x22, 0xc7, 0x5a, 0xbf, 0xa2, 0xf5,
	0xbb, 0x73, 0xfd, 0x55, 0xfe, 0xd4, 0x2f, 0xd2, 0xdd, 0x33, 0x38, 0x28, 0x0b, 0xae, 0xab, 0x02,
	0xa1, 0x2e, 0x10, 0xd3, 0x14, 0x77, 0x61, 0x2d, 0x91, 0x54, 0x2c, 0xb2, 0x9b, 0x5b, 0xee, 0x2f,
	0x03, 0x9c, 0xe5, 0xda, 0xff, 0x21, 0xc1, 0x0f, 0xb0, 0xd5, 0x1a, 0xc4, 0x7d, 0x1c, 0x06, 0x8a,
	0xf6, 0xf8, 0x84, 0xb2, 0x34, 0x25, 0x1b, 0x9e, 0x5e, 0x4f, 0x03, 0x15, 0x72, 0x11, 0xeb, 0x78,
	0x2a, 0xfe, 0xc2, 0x4e, 0xcf, 0x5a, 0x83, 0x58, 0x53, 0xf5, 0x5d, 0x55, 0x7f, 0x61, 0xbb, 0xdb,
	0x40, 0x8a, 0xce, 0x24, 0xba, 0x9f, 0xa0, 0xde, 0xa5, 0x69, 0xc6, 0x97, 0xc1, 0x70, 0xa4, 0x1b,
	0x70, 0x1f, 0xa0, 0x25, 0x78, 0xdc, 0xcf, 0x3e, 0x5b, 0x06, 0x21, 0x87, 0xfa, 0x6f, 0xf4, 0x99,
	0xa0, 0xc1, 0xb0, 0xc9, 0x13, 0xa6, 0xf4, 0x4d, 0x15, 0xbf, 0x80, 0xba, 0x5b, 0xb0, 0x91, 0xf3,
	0x2c, 0xf1, 0xf4, 0x87, 0x09, 0xb3, 0x99, 0x40, 0xce, 0x00, 0x1e, 0x3e, 0x18, 0xd9, 0x6e, 0xcc,
	0xc6, 0x46, 0xee, 0xf3, 0xdb, 0x3b, 0x25, 0xa8, 0x44, 0xf2, 0x11, 0xb6, 0xcb, 0x2a, 0x47, 0xf6,
	0xe7, 0xf4, 0x47, 0xba, 0xdd, 0x3e, 0x58, 0x7a, 0x2e, 0x91, 0x70, 0xd8, 0x5b, 0xd6, 0x12, 0xe4,
	0x70, 0x89, 0x83, 0x4c, 0xcf, 0xd9, 0x47, 0xff, 0xc4, 0x93, 0x48, 0xde, 0x41, 0x3d, 0x5f, 0x11,
	0x62, 0xcd, 0xa5, 0x7f, 0x54, 0xdd, 0x7e, 0xfe, 0xc8, 0x89, 0x44, 0xf2, 0x0a, 0x6a, 0x99, 0x87,
	0x26, 0xf7, 0xcf, 0x96, 0x2f, 0xab, 0xbd, 0x5b, 0x06, 0x4b, 0xbc, 0xdc, 0xff, 0xbc, 0x97, 0xce,
	0xe7, 0x76, 0x27, 0x33, 0x97, 0x35, 0xf3, 0x5c, 0xaf, 0x5f, 0xd6, 0x34, 0xf4, 0xf2, 0x77, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x85, 0xa6, 0x11, 0xd8, 0xe0, 0x05, 0x00, 0x00,
}
