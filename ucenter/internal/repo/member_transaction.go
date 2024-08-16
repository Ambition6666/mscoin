package repo

import (
	"context"
	"ucenter/internal/model"
)

type MemberTransactionRepo interface {
	FindTransaction(
		ctx context.Context,
		pageNo int,
		pageSize int,
		memberId int64,
		startTime string,
		endTime string,
		symbol string,
		transactionType string) (list []*model.MemberTransaction, total int64, err error)
	FindByAmountAndTime(ctx context.Context, amount float64, address string, time int64) (mt *model.MemberTransaction, err error)
	Save(ctx context.Context, transaction *model.MemberTransaction) error
}
