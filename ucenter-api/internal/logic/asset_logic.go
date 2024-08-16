package logic

import (
	"common/pages"
	"context"
	"github.com/jinzhu/copier"
	"grpc-common/ucenter/types/asset"
	"time"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Asset struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Asset {
	return &Asset{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Asset) FindWalletBySymbol(req *types.AssetReq) (*types.MemberWallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value := l.ctx.Value("userId").(int64)
	memberWallet, err := l.svcCtx.UAssetRPC.FindWalletBySymbol(ctx, &asset.AssetReq{
		CoinName: req.CoinName,
		UserId:   value,
	})
	if err != nil {
		return nil, err
	}
	resp := &types.MemberWallet{}
	if err := copier.Copy(resp, memberWallet); err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *Asset) FindWallet(req *types.AssetReq) ([]*types.MemberWallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value := l.ctx.Value("userId").(int64)
	memberWallets, err := l.svcCtx.UAssetRPC.FindWallet(ctx, &asset.AssetReq{
		UserId: value,
	})
	if err != nil {
		return nil, err
	}
	var resp []*types.MemberWallet
	if err := copier.Copy(&resp, memberWallets.List); err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *Asset) ResetWalletAddress(req *types.AssetReq) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value := l.ctx.Value("userId").(int64)
	_, err := l.svcCtx.UAssetRPC.ResetWalletAddress(ctx, &asset.AssetReq{
		UserId:   value,
		CoinName: req.Unit,
	})
	if err != nil {
		return "", err
	}
	return "", nil
}

func (l *Asset) FindTransaction(req *types.AssetReq) (*pages.PageResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value := l.ctx.Value("userId").(int64)
	mts, err := l.svcCtx.UAssetRPC.FindTransaction(ctx, &asset.AssetReq{
		UserId:    value,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Type:      req.Type,
		Symbol:    req.Symbol,
		PageNo:    int64(req.PageNo),
		PageSize:  int64(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	var resp []*types.MemberTransaction
	if err := copier.Copy(&resp, mts.List); err != nil {
		return nil, err
	}
	if resp == nil {
		resp = []*types.MemberTransaction{}
	}
	b := make([]any, len(resp))
	for i := range resp {
		b[i] = resp[i]
	}
	return pages.New(b, int64(req.PageNo), int64(req.PageSize), mts.Total), nil
}
