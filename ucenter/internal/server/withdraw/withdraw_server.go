package withdraw

import (
	"golang.org/x/net/context"
	"grpc-common/ucenter/types/withdraw"
	wlogic "ucenter/internal/logic/withdraw"
	"ucenter/internal/svc"
)

type WithdrawServer struct {
	svcCtx *svc.ServiceContext
	withdraw.UnimplementedWithdrawServer
}

func NewWithdrawServer(svcCtx *svc.ServiceContext) *WithdrawServer {
	return &WithdrawServer{
		svcCtx: svcCtx,
	}
}

func (s *WithdrawServer) FindAddressByCoinId(ctx context.Context, in *withdraw.WithdrawReq) (*withdraw.AddressSimpleList, error) {
	l := wlogic.NewWithdrawLogic(ctx, s.svcCtx)
	return l.FindAddressByCoinId(in)
}
func (s *WithdrawServer) SendCode(ctx context.Context, in *withdraw.WithdrawReq) (*withdraw.NoRes, error) {
	l := wlogic.NewWithdrawLogic(ctx, s.svcCtx)
	return l.SendCode(in)
}
func (s *WithdrawServer) WithdrawCode(ctx context.Context, in *withdraw.WithdrawReq) (*withdraw.NoRes, error) {
	l := wlogic.NewWithdrawLogic(ctx, s.svcCtx)
	return l.WithdrawCode(in)
}
func (s *WithdrawServer) WithdrawRecord(ctx context.Context, in *withdraw.WithdrawReq) (*withdraw.RecordList, error) {
	l := wlogic.NewWithdrawLogic(ctx, s.svcCtx)
	return l.WithdrawRecord(in)
}
