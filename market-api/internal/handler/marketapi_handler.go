package handler

import (
	"market-api/internal/svc"
)

func RegisterHandlers(r *Routers, serverCtx *svc.ServiceContext) {
	rate := NewExchangeRateHandler(serverCtx)
	rateGroup := r.Group()
	rateGroup.Post("/market/exchange-rate/usd/:unit", rate.GetUsdRate)

	market := NewMarketHandler(serverCtx)
	marketGroup := r.Group()
	marketGroup.Post("/market/symbol-thumb-trend", market.SymbolThumbTrend)
	marketGroup.Post("/market/symbol-thumb", market.SymbolThumb)
	marketGroup.Post("/market/symbol-info", market.SymbolInfo)
	marketGroup.Post("/market/coin-info", market.CoinInfo)
	marketGroup.Get("/market/history", market.History)
}

func RegisterWsHandlers(r *Routers, ctx *svc.ServiceContext) {
	group := r.Group()
	group.Get("/socket.io/", nil)
	group.Post("/socket.io/", nil)
}
