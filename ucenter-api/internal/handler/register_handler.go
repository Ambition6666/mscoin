package handler

import (
	"ucenter-api/internal/middlerware"
	"ucenter-api/internal/svc"
)

func RegisterHandlers(r *Routers, serverCtx *svc.ServiceContext) {
	register := NewRegisterHandler(serverCtx)
	registerGroup := r.Group()
	registerGroup.Post("/uc/register/phone", register.Register)
	registerGroup.Post("/uc/mobile/code", register.SendCode)

	login := NewLoginHandler(serverCtx)
	loginGroup := r.Group()
	loginGroup.Post("/uc/login", login.Login)
	loginGroup.Post("/uc/check/login", login.CheckLogin)

	asset := NewAssetHandler(serverCtx)
	assetGroup := r.Group()
	assetGroup.Use(middlerware.Auth(serverCtx.Config.JWT.AccessSecret))
	assetGroup.Post("/uc/asset/wallet/:coinName", asset.FindWalletBySymbol)

	assetGroup.Post("/uc/asset/wallet", asset.FindWallet)
	assetGroup.Post("/uc/asset/wallet/reset-address", asset.ResetWalletAddress)
	assetGroup.Post("/uc/asset/transaction/all", asset.FindTransaction)

	approveGroup := r.Group()
	approve := NewApproveHandler(serverCtx)
	approveGroup.Use(middlerware.Auth(serverCtx.Config.JWT.AccessSecret))
	approveGroup.Post("/uc/approve/security/setting", approve.SecuritySetting)

	withdrawGroup := r.Group()
	withdraw := NewWithdrawHandler(serverCtx)
	withdrawGroup.Use(middlerware.Auth(serverCtx.Config.JWT.AccessSecret))
	withdrawGroup.Post("/uc/withdraw/support/coin/info", withdraw.QueryWithdrawCoin)
	withdrawGroup.Post("/uc/withdraw/apply/code", withdraw.WithdrawCode)
	withdrawGroup.Post("/uc/mobile/withdraw/code", withdraw.SendCode)
	withdrawGroup.Post("/uc/withdraw/record", withdraw.Record)
}
