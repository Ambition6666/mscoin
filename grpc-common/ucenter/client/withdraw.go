package client

import (
	"github.com/zeromicro/go-zero/zrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"grpc-common/ucenter/types/withdraw"
)

type (
	WithdrawReq       = withdraw.WithdrawReq
	AddressSimpleList = withdraw.AddressSimpleList
	RecordList        = withdraw.RecordList

	Withdraw interface {
		FindAddressByCoinId(ctx context.Context, in *WithdrawReq, opts ...grpc.CallOption) (*AddressSimpleList, error)
		SendCode(ctx context.Context, in *WithdrawReq, opts ...grpc.CallOption) (*withdraw.NoRes, error)
		WithdrawCode(ctx context.Context, in *WithdrawReq, opts ...grpc.CallOption) (*withdraw.NoRes, error)
		WithdrawRecord(ctx context.Context, in *WithdrawReq, opts ...grpc.CallOption) (*RecordList, error)
	}

	defaultWithdraw struct {
		cli zrpc.Client
	}
)

func NewWithdraw(cli zrpc.Client) Withdraw {
	return &defaultWithdraw{
		cli: cli,
	}
}

func (d *defaultWithdraw) FindAddressByCoinId(ctx context.Context, in *WithdrawReq, opts ...grpc.CallOption) (*AddressSimpleList, error) {
	l := withdraw.NewWithdrawClient(d.cli.Conn())
	return l.FindAddressByCoinId(ctx, in, opts...)
}
func (d *defaultWithdraw) SendCode(ctx context.Context, in *WithdrawReq, opts ...grpc.CallOption) (*withdraw.NoRes, error) {
	l := withdraw.NewWithdrawClient(d.cli.Conn())
	return l.SendCode(ctx, in, opts...)
}
func (d *defaultWithdraw) WithdrawCode(ctx context.Context, in *WithdrawReq, opts ...grpc.CallOption) (*withdraw.NoRes, error) {
	l := withdraw.NewWithdrawClient(d.cli.Conn())
	return l.WithdrawCode(ctx, in, opts...)
}
func (d *defaultWithdraw) WithdrawRecord(ctx context.Context, in *WithdrawReq, opts ...grpc.CallOption) (*RecordList, error) {
	l := withdraw.NewWithdrawClient(d.cli.Conn())
	return l.WithdrawRecord(ctx, in, opts...)
}
