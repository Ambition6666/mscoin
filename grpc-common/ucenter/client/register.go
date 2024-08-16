// Code generated by goctl. DO NOT EDIT.
// Source: asset.proto

package client

import (
	"context"
	register2 "grpc-common/ucenter/types/register"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CaptchaReq = register2.CaptchaReq
	RegReq     = register2.RegReq
	RegRes     = register2.RegRes
	CodeReq = register2.CodeReq
	NoRes = register2.NoRes

	Register interface {
		RegisterByPhone(ctx context.Context, in *RegReq, opts ...grpc.CallOption) (*RegRes, error)
		SendCode(ctx context.Context, in *CodeReq, opts ...grpc.CallOption) (*NoRes, error)
	}

	defaultRegister struct {
		cli zrpc.Client
	}
)

func NewRegister(cli zrpc.Client) Register {
	return &defaultRegister{
		cli: cli,
	}
}

func (m *defaultRegister) RegisterByPhone(ctx context.Context, in *RegReq, opts ...grpc.CallOption) (*RegRes, error) {
	client := register2.NewRegisterClient(m.cli.Conn())
	return client.RegisterByPhone(ctx, in, opts...)
}

func (m *defaultRegister) SendCode(ctx context.Context, in *CodeReq, opts ...grpc.CallOption) (*NoRes, error) {
	client := register2.NewRegisterClient(m.cli.Conn())
	return client.SendCode(ctx, in, opts...)
}