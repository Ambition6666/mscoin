// Code genemarketd by goctl. DO NOT EDIT.
// Source: asset.proto

package client

import (
	"context"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"grpc-common/market/types/market"
)

type (
	MarketReq       = market.MarketReq
	SymbolThumbRes  = market.SymbolThumbRes
	ExchangeCoin    = market.ExchangeCoin
	Coin            = market.Coin
	HistoryRes      = market.HistoryRes
	ExchangeCoinRes = market.ExchangeCoinRes
	CoinList        = market.CoinList

	Market interface {
		FindSymbolThumbTrend(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*SymbolThumbRes, error)
		FindSymbolThumb(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*SymbolThumbRes, error)
		FindSymbolInfo(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*ExchangeCoin, error)
		FindCoinInfo(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*Coin, error)
		HistoryKline(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*HistoryRes, error)
		FindVisibleExchangeCoins(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*ExchangeCoinRes, error)
		FindAllCoin(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*CoinList, error)
		FindCoinById(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*Coin, error)
	}

	defaultMarket struct {
		cli zrpc.Client
	}
)

func NewMarket(cli zrpc.Client) Market {
	return &defaultMarket{
		cli: cli,
	}
}

func (m *defaultMarket) FindSymbolThumbTrend(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*SymbolThumbRes, error) {
	client := market.NewMarketClient(m.cli.Conn())
	return client.FindSymbolThumbTrend(ctx, in, opts...)
}
func (m *defaultMarket) FindSymbolThumb(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*SymbolThumbRes, error) {
	client := market.NewMarketClient(m.cli.Conn())
	return client.FindSymbolThumb(ctx, in, opts...)
}
func (m *defaultMarket) FindSymbolInfo(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*ExchangeCoin, error) {
	client := market.NewMarketClient(m.cli.Conn())
	return client.FindSymbolInfo(ctx, in, opts...)
}

func (m *defaultMarket) FindCoinInfo(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*Coin, error) {
	client := market.NewMarketClient(m.cli.Conn())
	return client.FindCoinInfo(ctx, in, opts...)
}

func (m *defaultMarket) HistoryKline(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*HistoryRes, error) {
	client := market.NewMarketClient(m.cli.Conn())
	return client.HistoryKline(ctx, in, opts...)
}

func (m *defaultMarket) FindVisibleExchangeCoins(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*ExchangeCoinRes, error) {
	client := market.NewMarketClient(m.cli.Conn())
	return client.FindVisibleExchangeCoins(ctx, in, opts...)
}
func (m *defaultMarket) FindAllCoin(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*CoinList, error) {
	client := market.NewMarketClient(m.cli.Conn())
	return client.FindAllCoin(ctx, in, opts...)
}

func (m *defaultMarket) FindCoinById(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*Coin, error) {
	client := market.NewMarketClient(m.cli.Conn())
	return client.FindCoinById(ctx, in, opts...)
}
