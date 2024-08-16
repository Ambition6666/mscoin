package domain

import (
	"common/msdb"
	"common/msdb/tran"
	"common/tools"
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/cache"
	mclient "grpc-common/market/client"
	"grpc-common/market/types/market"

	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"
)

type MemberWalletDomain struct {
	memberWalletRepo repo.MemberWalletRepo
	transaction      *tran.TransactionImpl
	cache            cache.Cache
}

func (d *MemberWalletDomain) FindWalletBySymbol(ctx context.Context, id int64, name string, coin *mclient.Coin) (*model.MemberWalletCoin, error) {
	mw, err := d.memberWalletRepo.FindByIdAndCoinName(ctx, id, name)
	if err != nil {
		return nil, err
	}
	if mw == nil {
		//新建并存储
		mw, walletCoin := model.NewMemberWallet(id, coin)
		err := d.memberWalletRepo.Save(ctx, mw)
		if err != nil {
			return nil, err
		}
		return walletCoin, nil
	}
	nwc := &model.MemberWalletCoin{}
	copier.Copy(nwc, mw)
	nwc.Coin = coin
	return nwc, nil
}

func (d *MemberWalletDomain) Freeze(ctx context.Context, conn msdb.DbConn, userId int64, money float64, symbol string) error {
	mw, err := d.memberWalletRepo.FindByIdAndCoinName(ctx, userId, symbol)
	if err != nil {
		return err
	}
	if mw.Balance < money {
		return errors.New("余额不足")
	}
	err = d.memberWalletRepo.UpdateFreeze(ctx, conn, userId, money, symbol)
	return err
}

func (d *MemberWalletDomain) UpdateWalletCoinAndBase(ctx context.Context, baseWallet *model.MemberWallet, coinWallet *model.MemberWallet) error {
	return d.transaction.Action(func(conn msdb.DbConn) error {
		err := d.memberWalletRepo.UpdateWallet(ctx, conn, baseWallet.Id, baseWallet.Balance, baseWallet.FrozenBalance)
		if err != nil {
			return err
		}
		err = d.memberWalletRepo.UpdateWallet(ctx, conn, coinWallet.Id, coinWallet.Balance, coinWallet.FrozenBalance)
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *MemberWalletDomain) FindWalletByMemIdAndCoinName(ctx context.Context, memberId int64, coinName string) (*model.MemberWallet, error) {
	mw, err := d.memberWalletRepo.FindByIdAndCoinName(ctx, memberId, coinName)
	if err != nil {
		return nil, err
	}
	return mw, nil
}

func (d *MemberWalletDomain) FindWalletByMemIdAndCoinId(ctx context.Context, memberId int64, coinId int64) (*model.MemberWallet, error) {
	mw, err := d.memberWalletRepo.FindByIdAndCoinId(ctx, memberId, coinId)
	if err != nil {
		return nil, err
	}
	return mw, nil
}

func (d *MemberWalletDomain) FindWalletByMemId(ctx context.Context, userId int64) ([]*model.MemberWallet, error) {
	memberWallets, err := d.memberWalletRepo.FindByMemId(ctx, userId)
	return memberWallets, err
}

func (d *MemberWalletDomain) Copy(memberWallet *model.MemberWallet, coinInfo *mclient.Coin) *model.MemberWalletCoin {
	mwc := &model.MemberWalletCoin{}
	copier.Copy(mwc, memberWallet)
	mwc.Coin = &market.Coin{}
	copier.Copy(mwc.Coin, coinInfo)
	var cnyRate string
	d.cache.Get("USDT::CNY::RATE", &cnyRate)
	if memberWallet.CoinName != "USDT" {
		//获取最新的汇率
		var usdRate string
		d.cache.Get(memberWallet.CoinName+"::USDT::RATE", &usdRate)
		if usdRate == "" {
			usdRate = "1"
		}
		mwc.Coin.UsdRate = tools.ToFloat64(usdRate)
		mwc.Coin.CnyRate = tools.MulN(tools.ToFloat64(usdRate), tools.ToFloat64(cnyRate), 10)
	} else {
		mwc.Coin.UsdRate = 1
		mwc.Coin.CnyRate = tools.ToFloat64(cnyRate)
	}
	return mwc
}

func (d *MemberWalletDomain) UpdateAddress(ctx context.Context, mw *model.MemberWallet) error {
	return d.memberWalletRepo.UpdateAddress(ctx, mw)
}

func (d *MemberWalletDomain) GetAddress(ctx context.Context, coin_name string) ([]string, error) {
	return d.memberWalletRepo.GetAddress(ctx, coin_name)
}
func (d *MemberWalletDomain) FindByAddress(address string) (*model.MemberWallet, error) {
	return d.memberWalletRepo.FindByAddress(context.Background(), address)
}
func NewMemberWalletDomain(db *msdb.MsDB, rcli cache.Cache) *MemberWalletDomain {
	return &MemberWalletDomain{
		dao.NewMemberWalletDao(db),
		tran.NewTransaction(db.Conn),
		rcli,
	}
}
