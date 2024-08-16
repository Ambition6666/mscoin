package domain

import (
	"common/msdb"
	"context"
	"market/internal/dao"
	"market/internal/model"
	"market/internal/repo"
)

type CoinDomain struct {
	CoinRepo repo.CoinRepo
}

func (d *CoinDomain) FindCoinInfo(ctx context.Context, unit string) (*model.Coin, error) {
	coin, err := d.CoinRepo.FindByUnit(ctx, unit)
	coin.ColdWalletAddress = ""
	return coin, err
}

func (d *CoinDomain) FindCoinById(ctx context.Context, id int64) (*model.Coin, error) {
	coin, err := d.CoinRepo.FindById(ctx, id)
	coin.ColdWalletAddress = ""
	return coin, err
}

func (d *CoinDomain) FindAllCoin(ctx context.Context) ([]*model.Coin, error) {
	return d.CoinRepo.FindAll(ctx)
}

func NewConnDomain(db *msdb.MsDB) *CoinDomain {
	return &CoinDomain{
		CoinRepo: dao.NewCoinDao(db),
	}
}
