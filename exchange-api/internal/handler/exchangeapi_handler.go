package handler

import (
	"exchange-api/internal/middlerware"
	"exchange-api/internal/svc"
)

func RegisterHandlers(r *Routers, serverCtx *svc.ServiceContext) {
	order := NewOrderHandler(serverCtx)
	orderGroup := r.Group()
	orderGroup.Use(middlerware.Auth(serverCtx.Config.JWT.AccessSecret))
	orderGroup.Post("/exchange/asset/history", order.History)
	orderGroup.Post("/exchange/asset/current", order.Current)
	orderGroup.Post("/exchange/asset/add", order.Add)
}
