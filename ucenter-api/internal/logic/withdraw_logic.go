package logic

import (
	"common/pages"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/asset"
	"grpc-common/ucenter/types/member"
	"grpc-common/ucenter/types/withdraw"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
)

type Withdraw struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (w *Withdraw) QueryWithdrawCoin(req *types.WithdrawReq) ([]*types.WithdrawWalletInfo, error) {
	value := w.ctx.Value("userId").(int64)
	coinList, err := w.svcCtx.MarketRPC.FindAllCoin(w.ctx, &market.MarketReq{})
	if err != nil {
		return nil, err
	}
	coinMap := make(map[string]*market.Coin)
	for _, v := range coinList.List {
		coinMap[v.Unit] = v
	}
	//查询用户钱包
	walletList, err := w.svcCtx.UAssetRPC.FindWallet(w.ctx, &asset.AssetReq{
		UserId: value,
	})
	if err != nil {
		return nil, err
	}
	var list = make([]*types.WithdrawWalletInfo, len(walletList.List))
	for i, wallet := range walletList.List {
		ww := &types.WithdrawWalletInfo{}
		coin := coinMap[wallet.Coin.Unit]
		addressSimpleList, err := w.svcCtx.UWithdrawRPC.FindAddressByCoinId(w.ctx, &withdraw.WithdrawReq{
			UserId: value,
			CoinId: int64(coin.Id),
		})
		if err != nil {
			return nil, err
		}
		var addresses []types.AddressSimple
		copier.Copy(&addresses, addressSimpleList.List)
		ww.Addresses = addresses
		ww.Balance = wallet.Balance
		ww.WithdrawScale = int(coin.WithdrawScale)
		ww.MaxTxFee = coin.MaxTxFee
		ww.MinTxFee = coin.MinTxFee
		ww.MaxAmount = coin.MaxWithdrawAmount
		ww.MinAmount = coin.MinWithdrawAmount
		ww.Name = coin.GetName()
		ww.NameCn = coin.NameCn
		ww.Threshold = coin.WithdrawThreshold
		ww.Unit = coin.Unit
		ww.AccountType = int(coin.AccountType)
		if coin.CanAutoWithdraw == 0 {
			ww.CanAutoWithdraw = "true"
		} else {
			ww.CanAutoWithdraw = "false"
		}
		list[i] = ww
	}
	return list, nil
}

func (w *Withdraw) SendCode(req *types.WithdrawReq) (string, error) {
	value := w.ctx.Value("userId").(int64)
	memberRes, err := w.svcCtx.UMemberRPC.FindMemberById(w.ctx, &member.MemberReq{
		MemberId: value,
	})
	if err != nil {
		return "fail", err
	}
	phone := memberRes.MobilePhone
	_, err = w.svcCtx.UWithdrawRPC.SendCode(w.ctx, &withdraw.WithdrawReq{
		Phone: phone,
	})
	if err != nil {
		return "fail", err
	}
	return "success", nil
}

func (w *Withdraw) WithdrawCode(req *types.WithdrawReq) (string, error) {
	value := w.ctx.Value("userId").(int64)
	_, err := w.svcCtx.UWithdrawRPC.WithdrawCode(w.ctx, &withdraw.WithdrawReq{
		UserId:     value,
		Unit:       req.Unit,
		JyPassword: req.JyPassword,
		Code:       req.Code,
		Address:    req.Address,
		Amount:     req.Amount,
		Fee:        req.Fee,
	})
	if err != nil {
		return "fail", err
	}
	return "success", nil
}

func (w *Withdraw) Record(req *types.WithdrawReq) (*pages.PageResult, error) {
	value := w.ctx.Value("userId").(int64)
	records, err := w.svcCtx.UWithdrawRPC.WithdrawRecord(w.ctx, &withdraw.WithdrawReq{
		UserId:   value,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	list := records.List
	b := make([]any, len(list))
	for i := range list {
		b[i] = list[i]
	}
	return pages.New(b, int64(req.Page), int64(req.PageSize), records.Total), nil
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Withdraw {
	return &Withdraw{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
