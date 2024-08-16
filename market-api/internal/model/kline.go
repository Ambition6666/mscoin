package model

import (
	"common/tools"
	"grpc-common/market/types/market"
)

type Kline struct {
	Period       string  `bson:"period,omitempty" json:"period"`
	OpenPrice    float64 `bson:"openPrice,omitempty" json:"openPrice"`
	HighestPrice float64 `bson:"highestPrice,omitempty" json:"highestPrice"`
	LowestPrice  float64 `bson:"lowestPrice,omitempty" json:"lowestPrice"`
	ClosePrice   float64 `bson:"closePrice,omitempty" json:"closePrice"`
	Time         int64   `bson:"time,omitempty" json:"time"`
	Count        float64 `bson:"count,omitempty" json:"count"`       //成交笔数
	Volume       float64 `bson:"volume,omitempty" json:"volume"`     //成交量
	Turnover     float64 `bson:"turnover,omitempty" json:"turnover"` //成交额
}

func (*Kline) Table(symbol, period string) string {
	return "exchange_kline_" + symbol + "_" + period
}
func (k *Kline) ToCoinThumb(symbol string, ct *market.CoinThumb) *market.CoinThumb {
	isSame := false
	if ct.Close == k.ClosePrice &&
		ct.Open == k.OpenPrice &&
		ct.High == k.HighestPrice &&
		ct.Low == k.LowestPrice {
		isSame = true
	}
	if !isSame {
		ct.Close = k.ClosePrice
		ct.Open = k.OpenPrice
		if ct.High < k.HighestPrice {
			ct.High = k.HighestPrice
		}
		ct.Low = k.LowestPrice
		if ct.Low > k.LowestPrice {
			ct.Low = k.LowestPrice
		}
		ct.Zone = 0
		ct.Volume = tools.AddN(k.Volume, ct.Volume, 8)
		ct.Turnover = tools.AddN(k.Turnover, ct.Turnover, 8)
		ct.Change = k.LowestPrice - ct.Close
		ct.Chg = tools.MulN(tools.DivN(ct.Change, ct.Close, 5), 100, 5)
		ct.UsdRate = k.ClosePrice
		ct.BaseUsdRate = 1
		ct.Trend = append(ct.Trend, k.ClosePrice)
	}
	return ct
}

func (k *Kline) InitCoinThumb(symbol string) *market.CoinThumb {
	ct := &market.CoinThumb{}
	ct.Symbol = symbol
	ct.Close = k.ClosePrice
	ct.Open = k.OpenPrice
	ct.High = k.HighestPrice
	ct.Volume = k.Volume
	ct.Turnover = k.Turnover
	ct.Low = k.LowestPrice
	ct.Zone = 0
	ct.Change = 0
	ct.Chg = 0.0
	ct.UsdRate = k.ClosePrice
	ct.BaseUsdRate = 1
	ct.Trend = make([]float64, 0)
	return ct
}
