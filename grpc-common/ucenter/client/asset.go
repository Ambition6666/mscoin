package client

import (
	"github.com/zeromicro/go-zero/zrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"grpc-common/ucenter/types/asset"
)

type (
	AssetReq              = asset.AssetReq
	AssetResp             = asset.AssetResp
	AddressList           = asset.AddressList
	MemberWallet          = asset.MemberWallet
	Coin                  = asset.Coin
	MemberWalletList      = asset.MemberWalletList
	MemberTransactionList = asset.MemberTransactionList

	Asset interface {
		FindWalletBySymbol(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*MemberWallet, error)
		FindWallet(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*MemberWalletList, error)
		ResetWalletAddress(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*AssetResp, error)
		FindTransaction(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*MemberTransactionList, error)
		GetAddress(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*AddressList, error)
	}

	defaultAsset struct {
		cli zrpc.Client
	}
)

func NewAsset(cli zrpc.Client) Asset {
	return &defaultAsset{
		cli: cli,
	}
}

func (m *defaultAsset) FindWalletBySymbol(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*MemberWallet, error) {
	client := asset.NewAssetClient(m.cli.Conn())
	return client.FindWalletBySymbol(ctx, in, opts...)
}

func (m *defaultAsset) FindWallet(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*MemberWalletList, error) {
	client := asset.NewAssetClient(m.cli.Conn())
	return client.FindWallet(ctx, in, opts...)
}

func (m *defaultAsset) ResetWalletAddress(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*AssetResp, error) {
	client := asset.NewAssetClient(m.cli.Conn())
	return client.ResetWalletAddress(ctx, in, opts...)
}

func (m *defaultAsset) FindTransaction(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*MemberTransactionList, error) {
	client := asset.NewAssetClient(m.cli.Conn())
	return client.FindTransaction(ctx, in, opts...)
}

func (m *defaultAsset) GetAddress(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*AddressList, error) {
	client := asset.NewAssetClient(m.cli.Conn())
	return client.GetAddress(ctx, in, opts...)
}
