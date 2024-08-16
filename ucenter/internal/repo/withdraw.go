package repo

import (
	"common/msdb"
	"golang.org/x/net/context"
	"ucenter/internal/model"
)

type WithdrawRecordRepo interface {
	Save(ctx context.Context, conn msdb.DbConn, record *model.WithdrawRecord) error
	UpdateTransactionNumber(ctx context.Context, wr model.WithdrawRecord) error
	FindByUserId(ctx context.Context, id int64, page int32, size int32) ([]*model.WithdrawRecord, int64, error)
}
