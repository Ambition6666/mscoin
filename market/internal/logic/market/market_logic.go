package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/errgo.v2/errors"
	"grpc-common/market/types/market"
	"market/internal/domain"
	"market/internal/svc"
	"time"
)

type MarketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	marketDomain       *domain.MarketDomain
	exchangeCoinDomain *domain.ExchangeCoinDomain
	coinDomain         *domain.CoinDomain
}

func NewMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketLogic {
	return &MarketLogic{
		ctx:                ctx,
		svcCtx:             svcCtx,
		Logger:             logx.WithContext(ctx),
		marketDomain:       domain.NewMarketDomain(svcCtx.MongoClient),
		exchangeCoinDomain: domain.NewExchangeCoinDomain(svcCtx.Db),
		coinDomain:         domain.NewConnDomain(svcCtx.Db),
	}
}

func (l *MarketLogic) FindSymbolThumb(in *market.MarketReq) (*market.SymbolThumbRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exchangeCoins := l.exchangeCoinDomain.FindVisible(ctx)
	coinThumbs := make([]*market.CoinThumb, len(exchangeCoins))
	for i, v := range exchangeCoins {
		ct := &market.CoinThumb{}
		ct.Symbol = v.Symbol
		coinThumbs[i] = ct
	}
	return &market.SymbolThumbRes{
		List: coinThumbs,
	}, nil
}

func (l *MarketLogic) FindSymbolThumbTrend(in *market.MarketReq) (*market.SymbolThumbRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exchangeCoins := l.exchangeCoinDomain.FindVisible(ctx)
	coinThumbs := l.marketDomain.SymbolThumbTrend(exchangeCoins)
	return &market.SymbolThumbRes{
		List: coinThumbs,
	}, nil
}

func (l *MarketLogic) FindSymbolInfo(req *market.MarketReq) (*market.ExchangeCoin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exchangeCoin, err := l.exchangeCoinDomain.FindSymbol(ctx, req.Symbol)
	if err != nil {
		return nil, err
	}
	mc := &market.ExchangeCoin{}
	if err := copier.Copy(mc, exchangeCoin); err != nil {
		return nil, err
	}
	return mc, nil
}

func (l *MarketLogic) FindCoinInfo(req *market.MarketReq) (*market.Coin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coin, err := l.coinDomain.FindCoinInfo(ctx, req.Unit)
	if err != nil {
		return nil, err
	}
	mc := &market.Coin{}
	if err := copier.Copy(mc, coin); err != nil {
		return nil, err
	}
	return mc, nil
}

func (l *MarketLogic) FindCoinById(req *market.MarketReq) (*market.Coin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coin, err := l.coinDomain.FindCoinById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	mc := &market.Coin{}
	if err := copier.Copy(mc, coin); err != nil {
		return nil, err
	}
	return mc, nil
}

func (l *MarketLogic) HistoryKline(req *market.MarketReq) (*market.HistoryRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	period := "1H"
	if req.Resolution == "60" {
		period = "1H"
	} else if req.Resolution == "30" {
		period = "30m"
	} else if req.Resolution == "15" {
		period = "15m"
	} else if req.Resolution == "5" {
		period = "5m"
	} else if req.Resolution == "1" {
		period = "1m"
	} else if req.Resolution == "1D" {
		period = "1D"
	} else if req.Resolution == "1W" {
		period = "1W"
	} else if req.Resolution == "1M" {
		period = "1M"
	}
	histories, err := l.marketDomain.HistoryKline(ctx, req.Symbol, req.From, req.To, period)
	if err != nil {
		return nil, err
	}
	return &market.HistoryRes{
		List: histories,
	}, nil
}

func (l *MarketLogic) FindVisibleExchangeCoins() (*market.ExchangeCoinRes, error) {
	var list market.ExchangeCoinRes
	res := l.exchangeCoinDomain.FindVisible(context.Background())
	err := copier.Copy(&list.List, res)

	if err != nil {
		logx.Error("复制数据失败:", err)
		return nil, errors.New("exchangecoin复制数据失败")
	}

	return &list, nil
}
func (l *MarketLogic) FindAllCoin(req *market.MarketReq) (*market.CoinList, error) {
	ctx, cancel := context.WithTimeout(l.ctx, 5*time.Second)
	defer cancel()
	coinList, err := l.coinDomain.FindAllCoin(ctx)
	if err != nil {
		return nil, err
	}
	var list []*market.Coin
	copier.Copy(&list, coinList)
	return &market.CoinList{
		List: list,
	}, nil
}
