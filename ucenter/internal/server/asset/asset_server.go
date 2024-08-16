package server

import (
	"context"

	"grpc-common/ucenter/types/asset"
	assetlogic "ucenter/internal/logic/asset"
	"ucenter/internal/svc"
)

type AssetServer struct {
	svcCtx *svc.ServiceContext
	asset.UnimplementedAssetServer
}

func NewAssetServer(svcCtx *svc.ServiceContext) *AssetServer {
	return &AssetServer{
		svcCtx: svcCtx,
	}
}

func (s *AssetServer) FindWalletBySymbol(ctx context.Context, in *asset.AssetReq) (*asset.MemberWallet, error) {
	l := assetlogic.NewAssetLogic(ctx, s.svcCtx)
	return l.FindWalletBySymbol(in)
}

func (s *AssetServer) FindWallet(ctx context.Context, in *asset.AssetReq) (*asset.MemberWalletList, error) {
	l := assetlogic.NewAssetLogic(ctx, s.svcCtx)
	return l.FindWallet(in)
}

func (s *AssetServer) ResetWalletAddress(ctx context.Context, in *asset.AssetReq) (*asset.AssetResp, error) {
	l := assetlogic.NewAssetLogic(ctx, s.svcCtx)
	return l.ResetWalletAddress(in)
}

func (s *AssetServer) FindTransaction(ctx context.Context, in *asset.AssetReq) (*asset.MemberTransactionList, error) {
	l := assetlogic.NewAssetLogic(ctx, s.svcCtx)
	return l.FindTransaction(in)
}

func (s *AssetServer) GetAddress(ctx context.Context, in *asset.AssetReq) (*asset.AddressList, error) {
	l := assetlogic.NewAssetLogic(ctx, s.svcCtx)
	return l.GetAddress(ctx, in.CoinName)
}
