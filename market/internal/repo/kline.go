package repo

import (
	"context"
	"market/internal/model"
)

type KlineRepo interface {
	FindBySymbol(ctx context.Context, symbol, period string, count int64) ([]*model.Kline, error)
	FindBySymbolTime(ctx context.Context, symbol, period string, from, to int64, sort string) (list []*model.Kline, err error)
}
