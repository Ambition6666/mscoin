package client

import (
	"github.com/zeromicro/go-zero/zrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"grpc-common/ucenter/types/member"
)

type (
	MemberReq = member.MemberReq
	MemberRes = member.MemberRes

	Member interface {
		FindMemberById(ctx context.Context, in *MemberReq, opts ...grpc.CallOption) (*MemberRes, error)
	}

	defaultMember struct {
		cli zrpc.Client
	}
)

func NewMember(cli zrpc.Client) Member {
	return &defaultMember{
		cli: cli,
	}
}

func (m *defaultMember) FindMemberById(ctx context.Context, in *MemberReq, opts ...grpc.CallOption) (*MemberRes, error) {
	client := member.NewMemberClient(m.cli.Conn())
	return client.FindMemberById(ctx, in, opts...)
}
