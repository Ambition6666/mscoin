package dao

import (
	"common/msdb"
	"common/msdb/gorms"

	"golang.org/x/net/context"
	"ucenter/internal/model"
)

type WithdrawDao struct {
	conn *gorms.GormConn
}

func (m *WithdrawDao) UpdateTransactionNumber(ctx context.Context, wr model.WithdrawRecord) error {
	session := m.conn.Session(ctx)
	err := session.Model(&model.WithdrawRecord{}).
		Where("id=?", wr.Id).
		Updates(map[string]any{"transaction_number": wr.TransactionNumber, "status": wr.Status}).Error
	return err
}

func (m *WithdrawDao) Save(ctx context.Context, conn msdb.DbConn, record *model.WithdrawRecord) error {
	gormConn := conn.(*gorms.GormConn)
	tx := gormConn.Tx(ctx)
	err := tx.Create(&record).Error
	return err
}
func (m *WithdrawDao) FindByUserId(ctx context.Context, id int64, page int32, size int32) (list []*model.WithdrawRecord, total int64, err error) {
	session := m.conn.Session(ctx)
	db := session.Model(&model.WithdrawRecord{}).Where("member_id=?", id)

	offset := (page - 1) * size
	db.Count(&total)
	db.Order("create_time desc").Offset(int(offset)).Limit(int(size))
	err = db.Find(&list).Error
	return
}
func NewWithdrawDao(db *msdb.MsDB) *WithdrawDao {
	return &WithdrawDao{
		conn: gorms.New(db.Conn),
	}
}
