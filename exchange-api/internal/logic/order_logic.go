package logic

import (
	"common/pages"
	"context"
	"errors"
	"exchange-api/internal/svc"
	"exchange-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
)

type Order struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Order {
	return &Order{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Order) History(req *types.ExchangeReq) (*pages.PageResult, error) {
	value := l.ctx.Value("userId").(int64)
	history, err := l.svcCtx.OrderRPC.FindOrderHistory(l.ctx, &order.OrderReq{
		Symbol:   req.Symbol,
		Page:     req.PageNo,
		PageSize: req.PageSize,
		UserId:   value,
	})
	if err != nil {
		return nil, err
	}
	list := history.List
	b := make([]any, len(list))
	for i := range list {
		b[i] = list[i]
	}
	return pages.New(b, req.PageNo, req.PageSize, history.Total), nil
}

func (l *Order) Current(req *types.ExchangeReq) (*pages.PageResult, error) {
	value := l.ctx.Value("userId").(int64)
	history, err := l.svcCtx.OrderRPC.FindOrderCurrent(l.ctx, &order.OrderReq{
		Symbol:   req.Symbol,
		Page:     req.PageNo,
		PageSize: req.PageSize,
		UserId:   value,
	})
	if err != nil {
		return nil, err
	}
	list := history.List
	b := make([]any, len(list))
	for i := range list {
		b[i] = list[i]
	}
	return pages.New(b, req.PageNo, req.PageSize, history.Total), nil
}

func (l *Order) Add(req *types.ExchangeReq) (string, error) {
	value := l.ctx.Value("userId").(int64)
	if !req.OrderValid() {
		return "", errors.New("参数传递错误")
	}
	orderRes, err := l.svcCtx.OrderRPC.Add(l.ctx, &order.OrderReq{
		Symbol:    req.Symbol,
		UserId:    value,
		Direction: req.Direction,
		Type:      req.Type,
		Price:     req.Price,
		Amount:    req.Amount,
	})
	if err != nil {
		return "", err
	}
	return orderRes.OrderId, nil
}
