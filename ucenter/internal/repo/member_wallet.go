package repo

import (
	"common/msdb"
	"context"
	"ucenter/internal/model"
)

type MemberWalletRepo interface {
	Save(ctx context.Context, mw *model.MemberWallet) error
	FindByIdAndCoinName(ctx context.Context, memId int64, coinName string) (mw *model.MemberWallet, err error)
	UpdateFreeze(ctx context.Context, conn msdb.DbConn, memId int64, money float64, symbol string) error
	UpdateWallet(ctx context.Context, conn msdb.DbConn, id int64, balance float64, frozenBalance float64) error
	FindByMemId(ctx context.Context, memId int64) ([]*model.MemberWallet, error)
	UpdateAddress(ctx context.Context, mw *model.MemberWallet) error
	GetAddress(ctx context.Context, name string) ([]string, error)
	FindByAddress(background context.Context, address string) (*model.MemberWallet, error)
	FindByIdAndCoinId(ctx context.Context, id int64, id2 int64) (mw *model.MemberWallet, err error)
}
