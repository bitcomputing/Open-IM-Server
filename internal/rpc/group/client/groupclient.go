package client

import (
	pbgroup "Open_IM/pkg/proto/group"
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type GroupClient interface {
	CreateGroup(ctx context.Context, in *pbgroup.CreateGroupReq, opts ...grpc.CallOption) (*pbgroup.CreateGroupResp, error)
	JoinGroup(ctx context.Context, in *pbgroup.JoinGroupReq, opts ...grpc.CallOption) (*pbgroup.JoinGroupResp, error)
	QuitGroup(ctx context.Context, in *pbgroup.QuitGroupReq, opts ...grpc.CallOption) (*pbgroup.QuitGroupResp, error)
	GetGroupsInfo(ctx context.Context, in *pbgroup.GetGroupsInfoReq, opts ...grpc.CallOption) (*pbgroup.GetGroupsInfoResp, error)
	SetGroupInfo(ctx context.Context, in *pbgroup.SetGroupInfoReq, opts ...grpc.CallOption) (*pbgroup.SetGroupInfoResp, error)
	GetGroupApplicationList(ctx context.Context, in *pbgroup.GetGroupApplicationListReq, opts ...grpc.CallOption) (*pbgroup.GetGroupApplicationListResp, error)
	GetUserReqApplicationList(ctx context.Context, in *pbgroup.GetUserReqApplicationListReq, opts ...grpc.CallOption) (*pbgroup.GetUserReqApplicationListResp, error)
	TransferGroupOwner(ctx context.Context, in *pbgroup.TransferGroupOwnerReq, opts ...grpc.CallOption) (*pbgroup.TransferGroupOwnerResp, error)
	GroupApplicationResponse(ctx context.Context, in *pbgroup.GroupApplicationResponseReq, opts ...grpc.CallOption) (*pbgroup.GroupApplicationResponseResp, error)
	GetGroupMemberList(ctx context.Context, in *pbgroup.GetGroupMemberListReq, opts ...grpc.CallOption) (*pbgroup.GetGroupMemberListResp, error)
	GetGroupMembersInfo(ctx context.Context, in *pbgroup.GetGroupMembersInfoReq, opts ...grpc.CallOption) (*pbgroup.GetGroupMembersInfoResp, error)
	KickGroupMember(ctx context.Context, in *pbgroup.KickGroupMemberReq, opts ...grpc.CallOption) (*pbgroup.KickGroupMemberResp, error)
	GetJoinedGroupList(ctx context.Context, in *pbgroup.GetJoinedGroupListReq, opts ...grpc.CallOption) (*pbgroup.GetJoinedGroupListResp, error)
	InviteUserToGroup(ctx context.Context, in *pbgroup.InviteUserToGroupReq, opts ...grpc.CallOption) (*pbgroup.InviteUserToGroupResp, error)
	InviteUserToGroups(ctx context.Context, in *pbgroup.InviteUserToGroupsReq, opts ...grpc.CallOption) (*pbgroup.InviteUserToGroupsResp, error)
	GetGroupAllMember(ctx context.Context, in *pbgroup.GetGroupAllMemberReq, opts ...grpc.CallOption) (*pbgroup.GetGroupAllMemberResp, error)
	GetGroups(ctx context.Context, in *pbgroup.GetGroupsReq, opts ...grpc.CallOption) (*pbgroup.GetGroupsResp, error)
	GetGroupMembersCMS(ctx context.Context, in *pbgroup.GetGroupMembersCMSReq, opts ...grpc.CallOption) (*pbgroup.GetGroupMembersCMSResp, error)
	DismissGroup(ctx context.Context, in *pbgroup.DismissGroupReq, opts ...grpc.CallOption) (*pbgroup.DismissGroupResp, error)
	MuteGroupMember(ctx context.Context, in *pbgroup.MuteGroupMemberReq, opts ...grpc.CallOption) (*pbgroup.MuteGroupMemberResp, error)
	CancelMuteGroupMember(ctx context.Context, in *pbgroup.CancelMuteGroupMemberReq, opts ...grpc.CallOption) (*pbgroup.CancelMuteGroupMemberResp, error)
	MuteGroup(ctx context.Context, in *pbgroup.MuteGroupReq, opts ...grpc.CallOption) (*pbgroup.MuteGroupResp, error)
	CancelMuteGroup(ctx context.Context, in *pbgroup.CancelMuteGroupReq, opts ...grpc.CallOption) (*pbgroup.CancelMuteGroupResp, error)
	SetGroupMemberNickname(ctx context.Context, in *pbgroup.SetGroupMemberNicknameReq, opts ...grpc.CallOption) (*pbgroup.SetGroupMemberNicknameResp, error)
	GetJoinedSuperGroupList(ctx context.Context, in *pbgroup.GetJoinedSuperGroupListReq, opts ...grpc.CallOption) (*pbgroup.GetJoinedSuperGroupListResp, error)
	GetSuperGroupsInfo(ctx context.Context, in *pbgroup.GetSuperGroupsInfoReq, opts ...grpc.CallOption) (*pbgroup.GetSuperGroupsInfoResp, error)
	SetGroupMemberInfo(ctx context.Context, in *pbgroup.SetGroupMemberInfoReq, opts ...grpc.CallOption) (*pbgroup.SetGroupMemberInfoResp, error)
	GetGroupAbstractInfo(ctx context.Context, in *pbgroup.GetGroupAbstractInfoReq, opts ...grpc.CallOption) (*pbgroup.GetGroupAbstractInfoResp, error)
	GroupIsExist(ctx context.Context, in *pbgroup.GroupIsExistReq, opts ...grpc.CallOption) (*pbgroup.GroupIsExistResp, error)
	UserIsInGroup(ctx context.Context, in *pbgroup.UserIsInGroupReq, opts ...grpc.CallOption) (*pbgroup.UserIsInGroupResp, error)
}

type groupClient struct {
	client zrpc.Client
}

func NewGroupClient(config zrpc.RpcClientConf) GroupClient {
	return &groupClient{
		client: zrpc.MustNewClient(config),
	}
}

func (c *groupClient) CreateGroup(ctx context.Context, in *pbgroup.CreateGroupReq, opts ...grpc.CallOption) (*pbgroup.CreateGroupResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.CreateGroup(ctx, in)
}

func (c *groupClient) JoinGroup(ctx context.Context, in *pbgroup.JoinGroupReq, opts ...grpc.CallOption) (*pbgroup.JoinGroupResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.JoinGroup(ctx, in, opts...)
}

func (c *groupClient) QuitGroup(ctx context.Context, in *pbgroup.QuitGroupReq, opts ...grpc.CallOption) (*pbgroup.QuitGroupResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.QuitGroup(ctx, in, opts...)
}

func (c *groupClient) GetGroupsInfo(ctx context.Context, in *pbgroup.GetGroupsInfoReq, opts ...grpc.CallOption) (*pbgroup.GetGroupsInfoResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetGroupsInfo(ctx, in, opts...)
}

func (c *groupClient) SetGroupInfo(ctx context.Context, in *pbgroup.SetGroupInfoReq, opts ...grpc.CallOption) (*pbgroup.SetGroupInfoResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.SetGroupInfo(ctx, in, opts...)
}

func (c *groupClient) GetGroupApplicationList(ctx context.Context, in *pbgroup.GetGroupApplicationListReq, opts ...grpc.CallOption) (*pbgroup.GetGroupApplicationListResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetGroupApplicationList(ctx, in, opts...)
}

func (c *groupClient) GetUserReqApplicationList(ctx context.Context, in *pbgroup.GetUserReqApplicationListReq, opts ...grpc.CallOption) (*pbgroup.GetUserReqApplicationListResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetUserReqApplicationList(ctx, in, opts...)
}

func (c *groupClient) TransferGroupOwner(ctx context.Context, in *pbgroup.TransferGroupOwnerReq, opts ...grpc.CallOption) (*pbgroup.TransferGroupOwnerResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.TransferGroupOwner(ctx, in, opts...)
}

func (c *groupClient) GroupApplicationResponse(ctx context.Context, in *pbgroup.GroupApplicationResponseReq, opts ...grpc.CallOption) (*pbgroup.GroupApplicationResponseResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GroupApplicationResponse(ctx, in, opts...)
}

func (c *groupClient) GetGroupMemberList(ctx context.Context, in *pbgroup.GetGroupMemberListReq, opts ...grpc.CallOption) (*pbgroup.GetGroupMemberListResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetGroupMemberList(ctx, in, opts...)
}

func (c *groupClient) GetGroupMembersInfo(ctx context.Context, in *pbgroup.GetGroupMembersInfoReq, opts ...grpc.CallOption) (*pbgroup.GetGroupMembersInfoResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetGroupMembersInfo(ctx, in, opts...)
}

func (c *groupClient) KickGroupMember(ctx context.Context, in *pbgroup.KickGroupMemberReq, opts ...grpc.CallOption) (*pbgroup.KickGroupMemberResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.KickGroupMember(ctx, in, opts...)
}

func (c *groupClient) GetJoinedGroupList(ctx context.Context, in *pbgroup.GetJoinedGroupListReq, opts ...grpc.CallOption) (*pbgroup.GetJoinedGroupListResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetJoinedGroupList(ctx, in, opts...)
}

func (c *groupClient) InviteUserToGroup(ctx context.Context, in *pbgroup.InviteUserToGroupReq, opts ...grpc.CallOption) (*pbgroup.InviteUserToGroupResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.InviteUserToGroup(ctx, in, opts...)
}

func (c *groupClient) InviteUserToGroups(ctx context.Context, in *pbgroup.InviteUserToGroupsReq, opts ...grpc.CallOption) (*pbgroup.InviteUserToGroupsResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.InviteUserToGroups(ctx, in, opts...)
}

func (c *groupClient) GetGroupAllMember(ctx context.Context, in *pbgroup.GetGroupAllMemberReq, opts ...grpc.CallOption) (*pbgroup.GetGroupAllMemberResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetGroupAllMember(ctx, in, opts...)
}

func (c *groupClient) GetGroups(ctx context.Context, in *pbgroup.GetGroupsReq, opts ...grpc.CallOption) (*pbgroup.GetGroupsResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetGroups(ctx, in, opts...)
}

func (c *groupClient) GetGroupMembersCMS(ctx context.Context, in *pbgroup.GetGroupMembersCMSReq, opts ...grpc.CallOption) (*pbgroup.GetGroupMembersCMSResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetGroupMembersCMS(ctx, in, opts...)
}

func (c *groupClient) DismissGroup(ctx context.Context, in *pbgroup.DismissGroupReq, opts ...grpc.CallOption) (*pbgroup.DismissGroupResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.DismissGroup(ctx, in, opts...)
}

func (c *groupClient) MuteGroupMember(ctx context.Context, in *pbgroup.MuteGroupMemberReq, opts ...grpc.CallOption) (*pbgroup.MuteGroupMemberResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.MuteGroupMember(ctx, in, opts...)
}

func (c *groupClient) CancelMuteGroupMember(ctx context.Context, in *pbgroup.CancelMuteGroupMemberReq, opts ...grpc.CallOption) (*pbgroup.CancelMuteGroupMemberResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.CancelMuteGroupMember(ctx, in, opts...)
}

func (c *groupClient) MuteGroup(ctx context.Context, in *pbgroup.MuteGroupReq, opts ...grpc.CallOption) (*pbgroup.MuteGroupResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.MuteGroup(ctx, in, opts...)
}

func (c *groupClient) CancelMuteGroup(ctx context.Context, in *pbgroup.CancelMuteGroupReq, opts ...grpc.CallOption) (*pbgroup.CancelMuteGroupResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.CancelMuteGroup(ctx, in, opts...)
}

func (c *groupClient) SetGroupMemberNickname(ctx context.Context, in *pbgroup.SetGroupMemberNicknameReq, opts ...grpc.CallOption) (*pbgroup.SetGroupMemberNicknameResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.SetGroupMemberNickname(ctx, in, opts...)
}

func (c *groupClient) GetJoinedSuperGroupList(ctx context.Context, in *pbgroup.GetJoinedSuperGroupListReq, opts ...grpc.CallOption) (*pbgroup.GetJoinedSuperGroupListResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetJoinedSuperGroupList(ctx, in, opts...)
}

func (c *groupClient) GetSuperGroupsInfo(ctx context.Context, in *pbgroup.GetSuperGroupsInfoReq, opts ...grpc.CallOption) (*pbgroup.GetSuperGroupsInfoResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetSuperGroupsInfo(ctx, in, opts...)
}

func (c *groupClient) SetGroupMemberInfo(ctx context.Context, in *pbgroup.SetGroupMemberInfoReq, opts ...grpc.CallOption) (*pbgroup.SetGroupMemberInfoResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.SetGroupMemberInfo(ctx, in, opts...)
}

func (c *groupClient) GetGroupAbstractInfo(ctx context.Context, in *pbgroup.GetGroupAbstractInfoReq, opts ...grpc.CallOption) (*pbgroup.GetGroupAbstractInfoResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GetGroupAbstractInfo(ctx, in, opts...)
}

func (c *groupClient) GroupIsExist(ctx context.Context, in *pbgroup.GroupIsExistReq, opts ...grpc.CallOption) (*pbgroup.GroupIsExistResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.GroupIsExist(ctx, in, opts...)
}

func (c *groupClient) UserIsInGroup(ctx context.Context, in *pbgroup.UserIsInGroupReq, opts ...grpc.CallOption) (*pbgroup.UserIsInGroupResp, error) {
	client := pbgroup.NewGroupClient(c.client.Conn())
	return client.UserIsInGroup(ctx, in, opts...)
}
