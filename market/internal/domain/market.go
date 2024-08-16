package domain

import (
	"common/tools"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market/database"
	"market/internal/dao"
	"market/internal/model"
	"market/internal/repo"
	"time"
)

type MarketDomain struct {
	klineRepo repo.KlineRepo
}

func (d *MarketDomain) SymbolThumbTrend(coins []*model.ExchangeCoin) []*market.CoinThumb {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	list := make([]*market.CoinThumb, len(coins))
	to := time.Now().UnixMilli()
	from := tools.ZeroTime()
	for i, v := range coins {
		klines, err := d.klineRepo.FindBySymbolTime(ctx, v.Symbol, "1H", from, to, "desc")
		if err != nil {
			logx.Error(err)
			continue
		}
		if len(klines) <= 0 {
			list[i] = &market.CoinThumb{
				Symbol: v.Symbol,
			}
			continue
		}
		trend := make([]float64, len(klines))
		var high float64 = 0
		var low float64 = 0
		var volumes float64 = 0
		var turnover float64 = 0
		for i, v := range klines {
			trend[i] = v.ClosePrice
			if v.HighestPrice > high {
				high = v.HighestPrice
			}
			if v.LowestPrice > low {
				low = v.LowestPrice
			}
			volumes += v.Volume
			turnover += v.Turnover
		}
		kline := klines[0]
		end := klines[len(klines)-1]
		ct := kline.ToCoinThumb(v.Symbol, end)
		ct.High = high
		ct.Low = low
		ct.Volume = volumes
		ct.Turnover = turnover
		ct.Trend = trend
		list[i] = ct
	}
	return list
}

func (d *MarketDomain) HistoryKline(ctx context.Context, symbol string, from int64, to int64, period string) ([]*market.History, error) {
	klines, err := d.klineRepo.FindBySymbolTime(ctx, symbol, period, from, to, "asc")
	if err != nil {
		return nil, err
	}
	list := make([]*market.History, len(klines))
	for i, v := range klines {
		h := &market.History{}
		h.Time = v.Time
		h.Open = v.OpenPrice
		h.High = v.HighestPrice
		h.Low = v.LowestPrice
		h.Volume = v.Volume
		h.Close = v.ClosePrice
		list[i] = h
	}
	return list, nil
}

func NewMarketDomain(cli *database.MongoClient) *MarketDomain {
	return &MarketDomain{
		klineRepo: dao.NewKlineDao(cli.Db),
	}
}
