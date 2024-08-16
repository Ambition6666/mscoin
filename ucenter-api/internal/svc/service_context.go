package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	mclient "grpc-common/market/client"
	"grpc-common/ucenter/client"
	"ucenter-api/internal/config"
)

type ServiceContext struct {
	Config       config.Config
	URegisterRPC client.Register
	ULoginRPC    client.Login
	UAssetRPC    client.Asset
	UMemberRPC   client.Member
	MarketRPC    mclient.Market
	UWithdrawRPC client.Withdraw
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		URegisterRPC: client.NewRegister(zrpc.MustNewClient(c.UCenterRPC)),
		ULoginRPC:    client.NewLogin(zrpc.MustNewClient(c.UCenterRPC)),
		UAssetRPC:    client.NewAsset(zrpc.MustNewClient(c.UCenterRPC)),
		UMemberRPC:   client.NewMember(zrpc.MustNewClient(c.UCenterRPC)),
		MarketRPC:    mclient.NewMarket(zrpc.MustNewClient(c.MarketRPC)),
		UWithdrawRPC: client.NewWithdraw(zrpc.MustNewClient(c.UCenterRPC)),
	}
}
