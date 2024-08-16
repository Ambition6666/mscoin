package logic

import (
	"common/bc"
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/asset"
	"ucenter/internal/domain"
	"ucenter/internal/model"
	"ucenter/internal/svc"
)

type AssetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	memberDomain            *domain.MemberDomain
	memberWalletDomain      *domain.MemberWalletDomain
	memberTransactionDomain *domain.MemberTransactionDomain
}

func (l *AssetLogic) FindWalletBySymbol(req *asset.AssetReq) (*asset.MemberWallet, error) {
	ctx := context.Background()
	coinInfo, err := l.svcCtx.MarketRPC.FindCoinInfo(ctx, &market.MarketReq{
		Unit: req.CoinName,
	})
	if err != nil {
		return nil, err
	}
	memberWalletCoin, err := l.memberWalletDomain.FindWalletBySymbol(ctx, req.UserId, req.CoinName, coinInfo)
	if err != nil {
		return nil, err
	}
	resp := &asset.MemberWallet{}
	copier.Copy(resp, memberWalletCoin)
	return resp, nil
}

func (l *AssetLogic) FindWallet(req *asset.AssetReq) (*asset.MemberWalletList, error) {
	mws, err := l.memberWalletDomain.FindWalletByMemId(l.ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var list []*model.MemberWalletCoin
	for _, v := range mws {
		coinInfo, err := l.svcCtx.MarketRPC.FindCoinInfo(l.ctx, &market.MarketReq{
			Unit: v.CoinName,
		})
		if err != nil {
			return nil, err
		}
		list = append(list, l.memberWalletDomain.Copy(v, coinInfo))
	}
	var mwList []*asset.MemberWallet
	copier.Copy(&mwList, list)
	return &asset.MemberWalletList{
		List: mwList,
	}, nil
}

func (l *AssetLogic) ResetWalletAddress(req *asset.AssetReq) (*asset.AssetResp, error) {
	mw, err := l.memberWalletDomain.FindWalletByMemIdAndCoinName(l.ctx, req.UserId, req.CoinName)
	if err != nil {
		return nil, err
	}
	//BTC的钱包地址逻辑
	if mw.Address == "" && req.CoinName == "BTC" {
		wallet, err := bc.NewWallet()
		if err != nil {
			return nil, err
		}
		address := wallet.GetTestAddress()
		mw.Address = string(address)
		mw.AddressPrivateKey = wallet.GetPriKey()
		err = l.memberWalletDomain.UpdateAddress(l.ctx, mw)
		if err != nil {
			return nil, err
		}
	}
	return &asset.AssetResp{}, nil
}

func (l *AssetLogic) FindTransaction(req *asset.AssetReq) (*asset.MemberTransactionList, error) {
	mms, _, err := l.memberTransactionDomain.FindTransaction(
		l.ctx,
		req.PageNo,
		req.PageSize,
		req.UserId,
		req.Symbol,
		req.StartTime,
		req.EndTime,
		req.Type)
	if err != nil {
		return nil, err
	}
	var mts []*asset.MemberTransaction
	copier.Copy(&mts, mms)
	return &asset.MemberTransactionList{
		List: mts,
	}, nil
}

func (l *AssetLogic) GetAddress(ctx context.Context, coin_name string) (*asset.AddressList, error) {
	ars, err := l.memberWalletDomain.GetAddress(ctx, coin_name)
	if err != nil {
		return nil, err
	}

	return &asset.AddressList{
		List: ars,
	}, nil
}

func NewAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssetLogic {
	return &AssetLogic{
		ctx:                     ctx,
		svcCtx:                  svcCtx,
		Logger:                  logx.WithContext(ctx),
		memberDomain:            domain.NewMemberDomain(svcCtx.Db),
		memberWalletDomain:      domain.NewMemberWalletDomain(svcCtx.Db, svcCtx.Cache),
		memberTransactionDomain: domain.NewMemberTransactionDomain(svcCtx.Db),
	}
}
