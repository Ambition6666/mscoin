package repo

import (
	"common/msdb"
	"context"
	"exchange/internal/model"
)

type ExchangeOrderRepo interface {
	FindBySymbolPage(
		ctx context.Context,
		symbol string,
		memId int64,
		status int,
		page int,
		pageSize int,
		isDesc bool) (list []*model.ExchangeOrder, total int64, err error)
	FindCount(
		ctx context.Context,
		memberId int64,
		symbol string,
		directionCode int,
		trading int) (int64, error)
	Save(ctx context.Context, conn msdb.DbConn, order *model.ExchangeOrder) error
	FindByOrderId(ctx context.Context, orderId string) (*model.ExchangeOrder, error)
	UpdateOrderStatusCancel(ctx context.Context, orderId string, status int, updateStatus int, cancelTime int64) error
	UpdateOrderStatus(ctx context.Context, orderId string, updateStatus int) error
	FindTradingOrders(ctx context.Context) (list []*model.ExchangeOrder, err error)
	UpdateOrderComplete(ctx context.Context, id string, amount float64, turnover float64, status int) error
}
