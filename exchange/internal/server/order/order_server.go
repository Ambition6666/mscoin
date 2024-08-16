package server

import (
	"context"

	orderlogic "exchange/internal/logic/order"
	"exchange/internal/svc"
	"grpc-common/exchange/types/order"
)

type OrderServer struct {
	svcCtx *svc.ServiceContext
	order.UnimplementedOrderServer
}

func NewOrderServer(svcCtx *svc.ServiceContext) *OrderServer {
	return &OrderServer{
		svcCtx: svcCtx,
	}
}

func (s *OrderServer) FindOrderHistory(ctx context.Context, in *order.OrderReq) (*order.OrderRes, error) {
	l := orderlogic.NewExchangeOrderLogic(ctx, s.svcCtx)
	return l.FindOrderHistory(in)
}

func (s *OrderServer) FindOrderCurrent(ctx context.Context, in *order.OrderReq) (*order.OrderRes, error) {
	l := orderlogic.NewExchangeOrderLogic(ctx, s.svcCtx)
	return l.FindOrderCurrent(in)
}

func (s *OrderServer) Add(ctx context.Context, in *order.OrderReq) (*order.AddOrderRes, error) {
	l := orderlogic.NewExchangeOrderLogic(ctx, s.svcCtx)
	return l.Add(in)
}

func (s *OrderServer) FindByOrderId(ctx context.Context, in *order.OrderReq) (*order.ExchangeOrderOrigin, error) {
	l := orderlogic.NewExchangeOrderLogic(ctx, s.svcCtx)
	return l.FindByOrderId(in)
}

func (s *OrderServer) CancelOrder(ctx context.Context, in *order.OrderReq) (*order.CancelOrderRes, error) {
	l := orderlogic.NewExchangeOrderLogic(ctx, s.svcCtx)
	return l.CancelOrder(in)
}
