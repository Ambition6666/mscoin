package server

import (
	"context"
	"grpc-common/market/types/market"
	"market/internal/logic/market"
	"market/internal/svc"
)

type MarketServer struct {
	svcCtx *svc.ServiceContext
	market.UnimplementedMarketServer
}

func NewMarketServer(svcCtx *svc.ServiceContext) *MarketServer {
	return &MarketServer{
		svcCtx: svcCtx,
	}
}

func (s *MarketServer) FindSymbolThumbTrend(ctx context.Context, in *market.MarketReq) (*market.SymbolThumbRes, error) {
	l := logic.NewMarketLogic(ctx, s.svcCtx)
	return l.FindSymbolThumbTrend(in)
}

func (s *MarketServer) FindSymbolThumb(ctx context.Context, in *market.MarketReq) (*market.SymbolThumbRes, error) {
	l := logic.NewMarketLogic(ctx, s.svcCtx)
	return l.FindSymbolThumb(in)
}

func (s *MarketServer) FindSymbolInfo(ctx context.Context, in *market.MarketReq) (*market.ExchangeCoin, error) {
	l := logic.NewMarketLogic(ctx, s.svcCtx)
	return l.FindSymbolInfo(in)
}

func (s *MarketServer) FindCoinInfo(ctx context.Context, in *market.MarketReq) (*market.Coin, error) {
	l := logic.NewMarketLogic(ctx, s.svcCtx)
	return l.FindCoinInfo(in)
}

func (s *MarketServer) HistoryKline(ctx context.Context, in *market.MarketReq) (*market.HistoryRes, error) {
	l := logic.NewMarketLogic(ctx, s.svcCtx)
	return l.HistoryKline(in)
}

func (s *MarketServer) FindVisibleExchangeCoins(ctx context.Context, in *market.MarketReq) (*market.ExchangeCoinRes, error) {
	l := logic.NewMarketLogic(ctx, s.svcCtx)
	return l.FindVisibleExchangeCoins()
}
func (s *MarketServer) FindAllCoin(ctx context.Context, in *market.MarketReq) (*market.CoinList, error) {
	l := logic.NewMarketLogic(ctx, s.svcCtx)
	return l.FindAllCoin(in)
}

func (s *MarketServer) FindCoinById(ctx context.Context, in *market.MarketReq) (*market.Coin, error) {
	l := logic.NewMarketLogic(ctx, s.svcCtx)
	return l.FindCoinById(in)
}
