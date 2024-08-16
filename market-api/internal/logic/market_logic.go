package logic

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market-api/internal/svc"
	"market-api/internal/types"
	"time"
)

type Market struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Market {
	return &Market{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Market) SymbolThumbTrend(req *types.MarketReq) (resp []*types.CoinThumbResp, err error) {
	processor := l.svcCtx.Processor
	thumb := processor.GetThumb()
	var list []*market.CoinThumb
	isCache := false
	if thumb != nil {
		m := thumb.(map[string]*market.CoinThumb)
		if len(m) > 0 {
			list = make([]*market.CoinThumb, len(m))
			i := 0
			for _, v := range m {
				list[i] = v
				i++
			}
			isCache = true
		}
	}
	if !isCache {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		thumbResp, err := l.svcCtx.MarketRPC.FindSymbolThumbTrend(ctx, &market.MarketReq{
			Ip: req.Ip,
		})
		if err != nil {
			return nil, err
		}
		list = thumbResp.List
		processor.PutThumb(list)
	}
	if err := copier.Copy(&resp, list); err != nil {
		return nil, errors.New("数据格式有误")
	}
	for _, v := range resp {
		if v.Trend == nil {
			v.Trend = []float64{}
		}
	}
	return
}

func (l *Market) SymbolThumb(req *types.MarketReq) (resp []*types.CoinThumbResp, err error) {
	processor := l.svcCtx.Processor
	thumb := processor.GetThumb()
	var list []*market.CoinThumb
	isCache := false
	if thumb != nil {
		m := thumb.(map[string]*market.CoinThumb)
		if len(m) > 0 {
			list = make([]*market.CoinThumb, len(m))
			i := 0
			for _, v := range m {
				list[i] = v
				i++
			}
			isCache = true
		}
	}
	if !isCache {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		thumbResp, err := l.svcCtx.MarketRPC.FindSymbolThumb(ctx, &market.MarketReq{
			Ip: req.Ip,
		})
		if err != nil {
			return nil, err
		}
		list = thumbResp.List
		processor.PutThumb(list)
	}
	if err := copier.Copy(&resp, list); err != nil {
		return nil, errors.New("数据格式有误")
	}
	for _, v := range resp {
		if v.Trend == nil {
			v.Trend = []float64{}
		}
	}
	return
}

func (l *Market) SymbolInfo(req *types.MarketReq) (*types.ExchangeCoinResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coin, err := l.svcCtx.MarketRPC.FindSymbolInfo(ctx, &market.MarketReq{
		Symbol: req.Symbol,
	})
	if err != nil {
		return nil, err
	}
	ec := &types.ExchangeCoinResp{}
	if err := copier.Copy(&ec, coin); err != nil {
		return nil, errors.New("数据格式有误")
	}
	ec.CurrentTime = time.Now().UnixMilli()
	return ec, nil
}

func (l *Market) CoinInfo(req *types.MarketReq) (*types.Coin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coin, err := l.svcCtx.MarketRPC.FindCoinInfo(ctx, &market.MarketReq{
		Unit: req.Unit,
	})
	if err != nil {
		return nil, err
	}
	ec := &types.Coin{}
	if err := copier.Copy(&ec, coin); err != nil {
		return nil, errors.New("数据格式有误")
	}
	return ec, nil
}

func (l *Market) History(req *types.MarketReq) (*types.HistoryKline, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	historyKline, err := l.svcCtx.MarketRPC.HistoryKline(ctx, &market.MarketReq{
		Symbol:     req.Symbol,
		From:       req.From,
		To:         req.To,
		Resolution: req.Resolution,
	})
	if err != nil {
		return nil, err
	}
	histories := historyKline.List
	var list = make([][]any, len(histories))
	for i, v := range histories {
		content := make([]any, 6)
		content[0] = v.Time
		content[1] = v.Open
		content[2] = v.High
		content[3] = v.Low
		content[4] = v.Close
		content[5] = v.Volume
		list[i] = content
	}
	return &types.HistoryKline{
		List: list,
	}, nil
}
