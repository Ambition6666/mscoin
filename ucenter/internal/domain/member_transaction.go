package domain

import (
	"common/msdb"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"
)

type MemberTransactionDomain struct {
	memberTransactionRepo repo.MemberTransactionRepo
	memberWalletDomain    *MemberWalletDomain
}

func (d *MemberTransactionDomain) FindTransaction(
	ctx context.Context,
	pageNo int64,
	pageSize int64,
	userId int64,
	symbol string,
	startTime string,
	endTime string,
	t string) ([]*model.MemberTransactionVo, int64, error) {
	list, total, err := d.memberTransactionRepo.FindTransaction(ctx, int(pageNo), int(pageSize), userId, startTime, endTime, symbol, t)
	if err != nil {
		return nil, total, err
	}
	var voList = make([]*model.MemberTransactionVo, len(list))
	for i, v := range list {
		voList[i] = v.ToVo()
	}
	return voList, total, nil
}

func (d *MemberTransactionDomain) SaveRecharge(address string, value float64, time int64, t string, symbol string) error {
	time = time * 1000
	ctx := context.Background()
	memberTransaction, err := d.memberTransactionRepo.FindByAmountAndTime(ctx, value, address, time)
	if err != nil {
		return err
	}
	wallet, err := d.memberWalletDomain.FindByAddress(address)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("address not exist ")
	}
	if memberTransaction == nil {
		transactionType := model.TypeMap.Code(t)
		memberTransaction = &model.MemberTransaction{}
		memberTransaction.MemberId = wallet.MemberId
		memberTransaction.Address = address
		memberTransaction.Type = transactionType
		memberTransaction.CreateTime = time * 1000
		memberTransaction.Amount = value
		memberTransaction.Symbol = symbol
		err := d.memberTransactionRepo.Save(ctx, memberTransaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewMemberTransactionDomain(db *msdb.MsDB, rcli cache.Cache) *MemberTransactionDomain {
	return &MemberTransactionDomain{
		memberTransactionRepo: dao.NewMemberTransactionDao(db),
		memberWalletDomain:    NewMemberWalletDomain(db, rcli),
	}
}
