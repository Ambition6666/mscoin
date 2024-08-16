package logic

import (
	"common/msdb"
	"common/msdb/tran"
	"context"
	"errors"
	"exchange/internal/domain"
	"exchange/internal/model"
	"exchange/internal/svc"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/asset"
	"grpc-common/ucenter/types/member"
)

type ExchangeOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	orderDomain *domain.ExchangeOrderDomain
	transaction tran.Transaction
	kafkaDomain *domain.KafkaDomain
}

func NewExchangeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeOrderLogic {
	orderDomain := domain.NewExchangeOrderDomain(svcCtx.Db)
	return &ExchangeOrderLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		orderDomain: orderDomain,
		transaction: tran.NewTransaction(svcCtx.Db.Conn),
		kafkaDomain: domain.NewKafkaDomain(svcCtx.Kafka, orderDomain),
	}
}

func (l *ExchangeOrderLogic) FindOrderHistory(in *order.OrderReq) (*order.OrderRes, error) {

	// 这里有自己做更改
	exchangeOrders, total, err := l.orderDomain.FindHistory(
		l.ctx,
		in.Symbol,
		in.Page,
		in.PageSize,
		in.UserId)
	if err != nil {
		return nil, err
	}
	var list []*order.ExchangeOrder
	copier.Copy(&list, exchangeOrders)
	return &order.OrderRes{
		List:  list,
		Total: total,
	}, nil
}

func (l *ExchangeOrderLogic) FindOrderCurrent(in *order.OrderReq) (*order.OrderRes, error) {
	exchangeOrders, total, err := l.orderDomain.FindTrading(
		l.ctx,
		in.Symbol,
		in.Page,
		in.PageSize,
		in.UserId)
	if err != nil {
		return nil, err
	}
	var list []*order.ExchangeOrder
	copier.Copy(&list, exchangeOrders)
	return &order.OrderRes{
		List:  list,
		Total: total,
	}, nil
}

func (l *ExchangeOrderLogic) Add(req *order.OrderReq) (*order.AddOrderRes, error) {
	memberRes, err := l.svcCtx.MemberRPC.FindMemberById(l.ctx, &member.MemberReq{
		MemberId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	if memberRes.TransactionStatus == 0 {
		return nil, errors.New("此用户已经被禁止交易")
	}
	if req.Type == model.TypeMap[model.LimitPrice] && req.Price <= 0 {
		return nil, errors.New("限价模式下价格不能小于等于0")
	}
	if req.Amount <= 0 {
		return nil, errors.New("数量不能小于等于0")
	}
	exchangeCoin, err := l.svcCtx.MarketRPC.FindSymbolInfo(l.ctx, &market.MarketReq{
		Symbol: req.Symbol,
	})
	if err != nil {
		return nil, errors.New("nonsupport coin")
	}
	if exchangeCoin.Exchangeable != 1 && exchangeCoin.Enable != 1 {
		return nil, errors.New("coin forbidden")
	}
	//基准币
	baseSymbol := exchangeCoin.GetBaseSymbol()
	//交易币
	coinSymbol := exchangeCoin.GetCoinSymbol()
	cc := baseSymbol
	if req.Direction == model.DirectionMap[model.SELL] {
		//根据交易币查询
		cc = coinSymbol
	}
	coin, err := l.svcCtx.MarketRPC.FindCoinInfo(l.ctx, &market.MarketReq{
		Unit: cc,
	})
	if err != nil || coin == nil {
		return nil, errors.New("nonsupport coin")
	}
	if req.Type == model.TypeMap[model.MarketPrice] && req.Direction == model.DirectionMap[model.BUY] {
		if exchangeCoin.GetMinTurnover() > 0 && req.Amount < float64(exchangeCoin.GetMinTurnover()) {
			return nil, errors.New("成交额至少是" + fmt.Sprintf("%d", exchangeCoin.GetMinTurnover()))
		}
	} else {
		if exchangeCoin.GetMaxVolume() > 0 && exchangeCoin.GetMaxVolume() < req.Amount {
			return nil, errors.New("数量超出" + fmt.Sprintf("%f", exchangeCoin.GetMaxVolume()))
		}
		if exchangeCoin.GetMinVolume() > 0 && exchangeCoin.GetMinVolume() > req.Amount {
			return nil, errors.New("数量不能低于" + fmt.Sprintf("%f", exchangeCoin.GetMinVolume()))
		}
	}
	//查询用户钱包
	baseWallet, err := l.svcCtx.AssetRPC.FindWalletBySymbol(l.ctx, &asset.AssetReq{
		UserId:   req.UserId,
		CoinName: baseSymbol,
	})
	if err != nil {
		return nil, errors.New("no wallet")
	}
	exCoinWallet, err := l.svcCtx.AssetRPC.FindWalletBySymbol(l.ctx, &asset.AssetReq{
		UserId:   req.UserId,
		CoinName: coinSymbol,
	})
	if err != nil {
		return nil, errors.New("no wallet")
	}
	if baseWallet.IsLock == 1 || exCoinWallet.IsLock == 1 {
		return nil, errors.New("wallet locked")
	}
	if req.Direction == model.DirectionMap[model.SELL] && exchangeCoin.GetMinSellPrice() > 0 {
		if req.Price < exchangeCoin.GetMinSellPrice() || req.Type == model.TypeMap[model.MarketPrice] {
			return nil, errors.New("不能低于最低限价:" + fmt.Sprintf("%f", exchangeCoin.GetMinSellPrice()))
		}
	}
	if req.Direction == model.DirectionMap[model.BUY] && exchangeCoin.GetMaxBuyPrice() > 0 {
		if req.Price > exchangeCoin.GetMaxBuyPrice() || req.Type == model.TypeMap[model.MarketPrice] {
			return nil, errors.New("不能低于最高限价:" + fmt.Sprintf("%f", exchangeCoin.GetMaxBuyPrice()))
		}
	}
	//是否启用了市价买卖
	if req.Type == model.TypeMap[model.MarketPrice] {
		if req.Direction == model.DirectionMap[model.BUY] && exchangeCoin.EnableMarketBuy == 0 {
			return nil, errors.New("不支持市价购买")
		} else if req.Direction == model.DirectionMap[model.SELL] && exchangeCoin.EnableMarketSell == 0 {
			return nil, errors.New("不支持市价出售")
		}
	}
	//限制委托数量
	count, err := l.orderDomain.FindCurrentTradingCount(l.ctx, req.UserId, req.Symbol, req.Direction)
	if err != nil {
		return nil, err
	}
	if exchangeCoin.GetMaxTradingOrder() > 0 && count >= int64(exchangeCoin.GetMaxTradingOrder()) {
		return nil, errors.New("超过最大挂单数量 " + fmt.Sprintf("%d", exchangeCoin.GetMaxTradingOrder()))
	}
	//开始生成订单
	exchangeOrder := model.NewOrder()
	exchangeOrder.MemberId = req.UserId
	exchangeOrder.Symbol = req.Symbol
	exchangeOrder.BaseSymbol = baseSymbol
	exchangeOrder.CoinSymbol = coinSymbol
	typeCode := model.TypeMap.Code(req.Type)
	exchangeOrder.Type = typeCode
	directionCode := model.DirectionMap.Code(req.Direction)
	exchangeOrder.Direction = directionCode
	if exchangeOrder.Type == model.MarketPrice {
		exchangeOrder.Price = 0
	} else {
		exchangeOrder.Price = req.Price
	}
	exchangeOrder.UseDiscount = "0"
	exchangeOrder.Amount = req.Amount
	err = l.transaction.Action(func(conn msdb.DbConn) error {
		money, err := l.orderDomain.AddOrder(l.ctx, conn, exchangeOrder, exchangeCoin, baseWallet, exCoinWallet)
		if err != nil {
			return errors.New("订单提交失败")
		}
		//通过kafka发送订单消息，进行钱包货币扣除 同步发送 要保证发送成功
		ok := l.kafkaDomain.Send(
			"add-exchange-asset",
			req.UserId,
			exchangeOrder.OrderId,
			money,
			req.Symbol,
			exchangeOrder.Direction,
			baseSymbol,
			coinSymbol,
			model.Init)
		if !ok {
			return errors.New("消息队列出现故障，未能扣款")
		}
		logx.Info("发送成功，订单id:", exchangeOrder.OrderId)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &order.AddOrderRes{
		OrderId: exchangeOrder.OrderId,
	}, nil
}

func (l *ExchangeOrderLogic) FindByOrderId(req *order.OrderReq) (*order.ExchangeOrderOrigin, error) {
	orderId := req.OrderId
	exOrder, err := l.orderDomain.FindByOrderId(l.ctx, orderId)
	if err != nil {
		return nil, err
	}
	res := &order.ExchangeOrderOrigin{}
	copier.Copy(res, exOrder)
	return res, nil
}

func (l *ExchangeOrderLogic) CancelOrder(req *order.OrderReq) (*order.CancelOrderRes, error) {
	orderId := req.OrderId
	err := l.orderDomain.UpdateOrderStatusCancel(l.ctx, orderId, int(req.UpdateStatus))
	if err != nil {
		return nil, err
	}
	return &order.CancelOrderRes{OrderId: orderId}, nil
}
