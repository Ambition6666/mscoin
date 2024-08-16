package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/rate"
	"market-api/internal/svc"
	"market-api/internal/types"
	"time"
)

type ExchangeRate struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExchangeRateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeRate {
	return &ExchangeRate{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExchangeRate) GetUsdRate(req *types.RateRequest) (resp *types.RateResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usdRate, err := l.svcCtx.ExchangeRateRPC.UsdRate(ctx, &rate.RateReq{
		Unit: req.Unit,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.RateResponse{
		Rate: usdRate.Rate,
	}
	return
}
