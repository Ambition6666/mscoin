package repo

import (
	"context"
	"market/internal/model"
)

type ExchangeCoinRepo interface {
	FindVisible(ctx context.Context) (list []*model.ExchangeCoin, err error)
	FindSymbol(ctx context.Context, symbol string) (*model.ExchangeCoin, error)
}
