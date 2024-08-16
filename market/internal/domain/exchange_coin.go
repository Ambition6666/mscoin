package domain

import (
	"common/msdb"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"market/internal/dao"
	"market/internal/model"
	"market/internal/repo"
)

type ExchangeCoinDomain struct {
	ExchangeCoinRepo repo.ExchangeCoinRepo
}

func NewExchangeCoinDomain(db *msdb.MsDB) *ExchangeCoinDomain {
	return &ExchangeCoinDomain{
		ExchangeCoinRepo: dao.NewExchangeCoinDao(db),
	}
}

func (d *ExchangeCoinDomain) FindVisible(ctx context.Context) []*model.ExchangeCoin {
	list, err := d.ExchangeCoinRepo.FindVisible(ctx)
	if err != nil {
		logx.Error(err)
		return []*model.ExchangeCoin{}
	}
	return list
}

func (d *ExchangeCoinDomain) FindSymbol(ctx context.Context, symbol string) (*model.ExchangeCoin, error) {
	coin, err := d.ExchangeCoinRepo.FindSymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return coin, nil
}
