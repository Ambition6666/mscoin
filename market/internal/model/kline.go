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

func (k *Kline) ToCoinThumb(symbol string, end *Kline) *market.CoinThumb {
	ct := &market.CoinThumb{}
	ct.Symbol = symbol
	ct.Close = k.ClosePrice
	ct.Open = k.OpenPrice
	ct.Zone = 0
	ct.Change = k.ClosePrice - end.ClosePrice
	ct.Chg = tools.MulN(tools.DivN(ct.Change, ct.Close, 5), 100, 5)
	ct.UsdRate = k.ClosePrice
	ct.BaseUsdRate = 1
	return ct
}
