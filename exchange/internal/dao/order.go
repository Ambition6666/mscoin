package dao

import (
	"common/msdb"
	"common/msdb/gorms"
	"context"
	"exchange/internal/model"
	"fmt"
	"gorm.io/gorm/clause"
)

type ExchangeOrderDao struct {
	conn *gorms.GormConn
}

func (d *ExchangeOrderDao) FindCount(ctx context.Context, memberId int64, symbol string, directionCode int, trading int) (total int64, err error) {
	session := d.conn.Session(ctx)
	err = session.Model(&model.ExchangeOrder{}).Where("member_id = ? and symbol = ? and direction = ? and trading = ?", memberId, symbol, directionCode, trading).Count(&total).Error
	return
}

func (d *ExchangeOrderDao) FindBySymbolPage(
	ctx context.Context,
	symbol string,
	memId int64,
	status int,
	page int,
	pageSize int,
	isDesc bool) (list []*model.ExchangeOrder, total int64, err error) {
	session := d.conn.Session(ctx)
	index := (page - 1) * pageSize
	query := fmt.Sprintf("symbol = ?")
	params := []any{symbol}
	if memId != -1 {
		query += fmt.Sprintf(" and member_id = ?")
		params = append(params, memId)
	}
	if status != -1 {
		query += fmt.Sprintf(" and status = ?")
		params = append(params, status)
	}
	err = session.
		Model(&model.ExchangeOrder{}).
		Where(query, params...).
		Limit(pageSize).Offset(index).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "time"}, Desc: isDesc}).
		Find(&list).Error
	err = session.
		Model(&model.ExchangeOrder{}).
		Where(query, params...).
		Count(&total).Error
	return
}

func (d *ExchangeOrderDao) UpdateOrderStatusCancel(ctx context.Context, orderId string, status int, updateStatus int, cancelTime int64) error {
	session := d.conn.Session(ctx)
	err := session.Model(&model.ExchangeOrder{}).Where("order_id=? and status=?", orderId, updateStatus).Update("status", status).Update("canceled_time=?", cancelTime).Error
	return err
}

func (d *ExchangeOrderDao) UpdateOrderStatus(ctx context.Context, orderId string, updateStatus int) error {
	session := d.conn.Session(ctx)
	err := session.Model(&model.ExchangeOrder{}).Where("order_id=?", orderId).Update("status", updateStatus).Error
	return err
}

func (d *ExchangeOrderDao) FindByOrderId(ctx context.Context, orderId string) (o *model.ExchangeOrder, err error) {
	session := d.conn.Session(ctx)
	err = session.Model(&model.ExchangeOrder{}).Where("order_id=?", orderId).Take(&o).Error
	if err != nil {
		return nil, nil
	}
	return
}

func (d *ExchangeOrderDao) Save(ctx context.Context, conn msdb.DbConn, order *model.ExchangeOrder) error {
	d.conn = conn.(*gorms.GormConn)
	tx := d.conn.Tx(ctx)
	err := tx.Save(&order).Error
	return err
}

func NewExchangeOrderDao(db *msdb.MsDB) *ExchangeOrderDao {
	return &ExchangeOrderDao{
		conn: gorms.New(db.Conn),
	}
}

func (d *ExchangeOrderDao) FindTradingOrders(ctx context.Context) (list []*model.ExchangeOrder, err error) {
	session := d.conn.Session(ctx)
	err = session.Model(&model.ExchangeOrder{}).Where("status=?", model.Trading).Find(&list).Error
	return
}

func (e *ExchangeOrderDao) UpdateOrderComplete(
	ctx context.Context,
	orderId string,
	tradedAmount float64,
	turnover float64,
	status int) error {
	session := e.conn.Session(ctx)
	err := session.Model(&model.ExchangeOrder{}).Where("order_id= ? and status= ?", orderId, model.Trading).Update("traded_amount", tradedAmount).Update("turnover", turnover).Update("status", status).Error
	return err
}
