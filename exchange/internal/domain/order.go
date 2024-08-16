package domain

import (
	"common/msdb"
	"common/tools"
	"context"
	"errors"
	"exchange/internal/dao"
	"exchange/internal/model"
	"exchange/internal/repo"
	"fmt"
	mclient "grpc-common/market/client"
	uclient "grpc-common/ucenter/client"
	"time"
)

type ExchangeOrderDomain struct {
	exchangeOrderRepo repo.ExchangeOrderRepo
}

func (d *ExchangeOrderDomain) FindHistory(
	ctx context.Context,
	symbol string,
	page int64,
	pageSize int64,
	userId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.exchangeOrderRepo.FindBySymbolPage(
		ctx,
		symbol,
		userId,
		-1,
		int(page),
		int(pageSize),
		true)
	lv := make([]*model.ExchangeOrderVo, len(list))
	for i, v := range list {
		lv[i] = v.ToVo()
	}
	return lv, total, err
}

func NewExchangeOrderDomain(db *msdb.MsDB) *ExchangeOrderDomain {
	return &ExchangeOrderDomain{
		exchangeOrderRepo: dao.NewExchangeOrderDao(db),
	}
}

func (d *ExchangeOrderDomain) FindTrading(
	ctx context.Context,
	symbol string,
	page int64,
	pageSize int64,
	userId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.exchangeOrderRepo.FindBySymbolPage(
		ctx,
		symbol,
		userId,
		model.Trading,
		int(page),
		int(pageSize),
		true)
	lv := make([]*model.ExchangeOrderVo, len(list))
	for i, v := range list {
		lv[i] = v.ToVo()
	}
	return lv, total, err
}

func (d *ExchangeOrderDomain) AddOrder(ctx context.Context, conn msdb.DbConn, exchangeOrder *model.ExchangeOrder, coin *mclient.ExchangeCoin, wallet *uclient.MemberWallet, coinWallet *uclient.MemberWallet) (float64, error) {
	exchangeOrder.Status = model.Init
	exchangeOrder.TradedAmount = 0
	exchangeOrder.Time = time.Now().UnixMilli()
	exchangeOrder.OrderId = tools.Unique("E")
	var money float64
	if exchangeOrder.Direction == model.BUY {
		var turnover float64 = 0
		if exchangeOrder.Type == model.MarketPrice {
			turnover = exchangeOrder.Amount
		} else {
			turnover = tools.MulN(exchangeOrder.Amount, exchangeOrder.Price, 5)
		}
		//费率
		fee := tools.MulN(turnover, coin.Fee, 5)
		if wallet.Balance < turnover {
			return 0, errors.New("余额不足")
		}
		if wallet.Balance-turnover < fee {
			return 0, errors.New("手续费不足 需要:" + fmt.Sprintf("%f", fee))
		}
		//需要冻结的钱 turnover+fee
		money = tools.AddN(turnover, fee, 5)
	} else {
		fee := tools.MulN(exchangeOrder.Amount, coin.Fee, 5)
		if coinWallet.Balance < exchangeOrder.Amount {
			return 0, errors.New("余额不足")
		}
		if wallet.Balance-exchangeOrder.Amount < fee {
			return 0, errors.New("手续费不足 需要:" + fmt.Sprintf("%f", fee))
		}
		money = tools.AddN(exchangeOrder.Amount, fee, 5)
	}
	err := d.exchangeOrderRepo.Save(ctx, conn, exchangeOrder)
	if err != nil {
		return 0, err
	}
	return money, nil
}

func (d *ExchangeOrderDomain) FindByOrderId(ctx context.Context, orderId string) (*model.ExchangeOrder, error) {
	exchangeOrder, err := d.exchangeOrderRepo.FindByOrderId(ctx, orderId)
	if err == nil && exchangeOrder == nil {
		return nil, errors.New("订单号不存在")
	}
	return exchangeOrder, err
}

func (d *ExchangeOrderDomain) UpdateOrderStatusCancel(ctx context.Context, orderId string, updateStatus int) error {
	return d.exchangeOrderRepo.UpdateOrderStatusCancel(ctx, orderId, model.Canceled, updateStatus, time.Now().UnixMilli())
}

func (d *ExchangeOrderDomain) FindCurrentTradingCount(ctx context.Context, id int64, symbol string, direction string) (int64, error) {
	return d.exchangeOrderRepo.FindCount(ctx, id, symbol, model.DirectionMap.Code(direction), model.Trading)
}

func (d *ExchangeOrderDomain) UpdateOrderStatus(background context.Context, id string, status int) error {
	return d.exchangeOrderRepo.UpdateOrderStatus(background, id, status)
}

func (d *ExchangeOrderDomain) FindTradingOrders(ctx context.Context) ([]*model.ExchangeOrder, error) {
	return d.exchangeOrderRepo.FindTradingOrders(ctx)
}

func (d *ExchangeOrderDomain) UpdateOrderComplete(ctx context.Context, order *model.ExchangeOrder) error {
	return d.exchangeOrderRepo.UpdateOrderComplete(ctx, order.OrderId, order.TradedAmount, order.Turnover, order.Status)
}
