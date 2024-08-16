package dao

import (
	"common/msdb"
	"common/msdb/gorms"
	"context"
	"market/internal/model"
)

type ExchangeCoinDao struct {
	conn *gorms.GormConn
}

func (d *ExchangeCoinDao) FindVisible(ctx context.Context) (list []*model.ExchangeCoin, err error) {
	session := d.conn.Session(ctx)
	err = session.Model(&model.ExchangeCoin{}).Where("visible=?", 1).Find(&list).Error
	return
}

func (d *ExchangeCoinDao) FindSymbol(ctx context.Context, symbol string) (*model.ExchangeCoin, error) {
	session := d.conn.Session(ctx)
	coin := &model.ExchangeCoin{}
	err := session.Model(&model.ExchangeCoin{}).Where("symbol=?", symbol).Take(coin).Error
	return coin, err
}

func NewExchangeCoinDao(db *msdb.MsDB) *ExchangeCoinDao {
	return &ExchangeCoinDao{
		conn: gorms.New(db.Conn),
	}
}
