package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/rate"
	"market/internal/domain"
	"market/internal/svc"
)

type ExchangeRateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	exchangeRateDomain *domain.ExchangeRateDomain
}

func NewExchangeRateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeRateLogic {
	return &ExchangeRateLogic{
		ctx:                ctx,
		svcCtx:             svcCtx,
		Logger:             logx.WithContext(ctx),
		exchangeRateDomain: domain.NewExchangeRateDomain(),
	}
}

func (l *ExchangeRateLogic) UsdRate(in *rate.RateReq) (*rate.RateRes, error) {
	usdtRate := l.exchangeRateDomain.GetUsdRate(in.GetUnit())
	return &rate.RateRes{
		Rate: usdtRate,
	}, nil
}
