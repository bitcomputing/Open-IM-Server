package client

import (
	pbmsg "Open_IM/pkg/proto/msg"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"

	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type MsgClient interface {
	GetMaxAndMinSeq(ctx context.Context, in *sdk_ws.GetMaxAndMinSeqReq, opts ...grpc.CallOption) (*sdk_ws.GetMaxAndMinSeqResp, error)
	PullMessageBySeqList(ctx context.Context, in *sdk_ws.PullMessageBySeqListReq, opts ...grpc.CallOption) (*sdk_ws.PullMessageBySeqListResp, error)
	SendMsg(ctx context.Context, in *pbmsg.SendMsgReq, opts ...grpc.CallOption) (*pbmsg.SendMsgResp, error)
	DelMsgList(ctx context.Context, in *sdk_ws.DelMsgListReq, opts ...grpc.CallOption) (*sdk_ws.DelMsgListResp, error)
	DelSuperGroupMsg(ctx context.Context, in *pbmsg.DelSuperGroupMsgReq, opts ...grpc.CallOption) (*pbmsg.DelSuperGroupMsgResp, error)
	ClearMsg(ctx context.Context, in *pbmsg.ClearMsgReq, opts ...grpc.CallOption) (*pbmsg.ClearMsgResp, error)
	SetMsgMinSeq(ctx context.Context, in *pbmsg.SetMsgMinSeqReq, opts ...grpc.CallOption) (*pbmsg.SetMsgMinSeqResp, error)
	SetSendMsgStatus(ctx context.Context, in *pbmsg.SetSendMsgStatusReq, opts ...grpc.CallOption) (*pbmsg.SetSendMsgStatusResp, error)
	GetSendMsgStatus(ctx context.Context, in *pbmsg.GetSendMsgStatusReq, opts ...grpc.CallOption) (*pbmsg.GetSendMsgStatusResp, error)
	GetSuperGroupMsg(ctx context.Context, in *pbmsg.GetSuperGroupMsgReq, opts ...grpc.CallOption) (*pbmsg.GetSuperGroupMsgResp, error)
	GetWriteDiffMsg(ctx context.Context, in *pbmsg.GetWriteDiffMsgReq, opts ...grpc.CallOption) (*pbmsg.GetWriteDiffMsgResp, error)
	// modify msg
	SetMessageReactionExtensions(ctx context.Context, in *pbmsg.SetMessageReactionExtensionsReq, opts ...grpc.CallOption) (*pbmsg.SetMessageReactionExtensionsResp, error)
	GetMessageListReactionExtensions(ctx context.Context, in *pbmsg.GetMessageListReactionExtensionsReq, opts ...grpc.CallOption) (*pbmsg.GetMessageListReactionExtensionsResp, error)
	AddMessageReactionExtensions(ctx context.Context, in *pbmsg.AddMessageReactionExtensionsReq, opts ...grpc.CallOption) (*pbmsg.AddMessageReactionExtensionsResp, error)
	DeleteMessageReactionExtensions(ctx context.Context, in *pbmsg.DeleteMessageListReactionExtensionsReq, opts ...grpc.CallOption) (*pbmsg.DeleteMessageListReactionExtensionsResp, error)
}

type msgClient struct {
	client zrpc.Client
}

func NewMsgClient(config zrpc.RpcClientConf) MsgClient {
	return &msgClient{
		client: zrpc.MustNewClient(config),
	}
}

func (c *msgClient) GetMaxAndMinSeq(ctx context.Context, in *sdk_ws.GetMaxAndMinSeqReq, opts ...grpc.CallOption) (*sdk_ws.GetMaxAndMinSeqResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.GetMaxAndMinSeq(ctx, in, opts...)
}

func (c *msgClient) PullMessageBySeqList(ctx context.Context, in *sdk_ws.PullMessageBySeqListReq, opts ...grpc.CallOption) (*sdk_ws.PullMessageBySeqListResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.PullMessageBySeqList(ctx, in, opts...)
}

func (c *msgClient) SendMsg(ctx context.Context, in *pbmsg.SendMsgReq, opts ...grpc.CallOption) (*pbmsg.SendMsgResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.SendMsg(ctx, in, opts...)
}

func (c *msgClient) DelMsgList(ctx context.Context, in *sdk_ws.DelMsgListReq, opts ...grpc.CallOption) (*sdk_ws.DelMsgListResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.DelMsgList(ctx, in, opts...)
}

func (c *msgClient) DelSuperGroupMsg(ctx context.Context, in *pbmsg.DelSuperGroupMsgReq, opts ...grpc.CallOption) (*pbmsg.DelSuperGroupMsgResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.DelSuperGroupMsg(ctx, in, opts...)
}

func (c *msgClient) ClearMsg(ctx context.Context, in *pbmsg.ClearMsgReq, opts ...grpc.CallOption) (*pbmsg.ClearMsgResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.ClearMsg(ctx, in, opts...)
}

func (c *msgClient) SetMsgMinSeq(ctx context.Context, in *pbmsg.SetMsgMinSeqReq, opts ...grpc.CallOption) (*pbmsg.SetMsgMinSeqResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.SetMsgMinSeq(ctx, in, opts...)
}

func (c *msgClient) SetSendMsgStatus(ctx context.Context, in *pbmsg.SetSendMsgStatusReq, opts ...grpc.CallOption) (*pbmsg.SetSendMsgStatusResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.SetSendMsgStatus(ctx, in, opts...)
}

func (c *msgClient) GetSendMsgStatus(ctx context.Context, in *pbmsg.GetSendMsgStatusReq, opts ...grpc.CallOption) (*pbmsg.GetSendMsgStatusResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.GetSendMsgStatus(ctx, in, opts...)
}

func (c *msgClient) GetSuperGroupMsg(ctx context.Context, in *pbmsg.GetSuperGroupMsgReq, opts ...grpc.CallOption) (*pbmsg.GetSuperGroupMsgResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.GetSuperGroupMsg(ctx, in, opts...)
}

func (c *msgClient) GetWriteDiffMsg(ctx context.Context, in *pbmsg.GetWriteDiffMsgReq, opts ...grpc.CallOption) (*pbmsg.GetWriteDiffMsgResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.GetWriteDiffMsg(ctx, in, opts...)
}

func (c *msgClient) SetMessageReactionExtensions(ctx context.Context, in *pbmsg.SetMessageReactionExtensionsReq, opts ...grpc.CallOption) (*pbmsg.SetMessageReactionExtensionsResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.SetMessageReactionExtensions(ctx, in, opts...)
}

func (c *msgClient) GetMessageListReactionExtensions(ctx context.Context, in *pbmsg.GetMessageListReactionExtensionsReq, opts ...grpc.CallOption) (*pbmsg.GetMessageListReactionExtensionsResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.GetMessageListReactionExtensions(ctx, in, opts...)
}

func (c *msgClient) AddMessageReactionExtensions(ctx context.Context, in *pbmsg.AddMessageReactionExtensionsReq, opts ...grpc.CallOption) (*pbmsg.AddMessageReactionExtensionsResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.AddMessageReactionExtensions(ctx, in, opts...)
}

func (c *msgClient) DeleteMessageReactionExtensions(ctx context.Context, in *pbmsg.DeleteMessageListReactionExtensionsReq, opts ...grpc.CallOption) (*pbmsg.DeleteMessageListReactionExtensionsResp, error) {
	client := pbmsg.NewMsgClient(c.client.Conn())
	return client.DeleteMessageReactionExtensions(ctx, in, opts...)
}
