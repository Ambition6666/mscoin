package processor

import (
	"common/msdb"
	"common/tools"
	"context"
	"encoding/json"
	"exchange/database"
	"exchange/internal/domain"
	"exchange/internal/model"
	"sort"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	mclient "grpc-common/market/client"
	"grpc-common/market/types/market"
	"sync"
)

type CoinTradeFactory struct {
	tradeMap map[string]*CoinTrade
	mux      sync.RWMutex
}

func InitCoinTradeFactory() *CoinTradeFactory {
	return &CoinTradeFactory{
		tradeMap: make(map[string]*CoinTrade),
	}
}
func (f *CoinTradeFactory) GetCoinTrade(symbol string) *CoinTrade {
	f.mux.RLock()
	defer f.mux.RUnlock()

	return f.tradeMap[symbol]
}

func (f *CoinTradeFactory) AddCoinTrade(symbol string, trade *CoinTrade) {
	f.mux.Lock()
	defer f.mux.Unlock()
	_, ok := f.tradeMap[symbol]
	if !ok {
		f.tradeMap[symbol] = trade
	}
}

func (f *CoinTradeFactory) Init(marketRpc mclient.Market, client *database.KafkaClient, db *msdb.MsDB) {
	ctx := context.Background()
	exchangeCoinRes, err := marketRpc.FindVisibleExchangeCoins(ctx, &market.MarketReq{})
	if err != nil {
		logx.Error(err)
		return
	}

	for _, v := range exchangeCoinRes.List {
		f.AddCoinTrade(v.Symbol, NewCoinTrade(v.Symbol, client, db))
	}
}

type LimitPriceQueue struct {
	mux  sync.RWMutex
	list TradeQueue
}
type LimitPriceMap struct {
	price float64
	list  []*model.ExchangeOrder
}

// 降序的排序
type TradeQueue []*LimitPriceMap

func (t TradeQueue) Len() int {
	return len(t)
}
func (t TradeQueue) Less(i, j int) bool {
	//降序
	return t[i].price > t[j].price
}
func (t TradeQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type TradeTimeQueue []*model.ExchangeOrder

func (t TradeTimeQueue) Len() int {
	return len(t)
}
func (t TradeTimeQueue) Less(i, j int) bool {
	//升序
	return t[i].Time < t[j].Time
}
func (t TradeTimeQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// CoinTrade 交易处理器
type CoinTrade struct {
	buyMarketQueue  TradeTimeQueue
	sellMarketQueue TradeTimeQueue
	buyLimitQueue   *LimitPriceQueue //从高到低
	sellLimitQueue  *LimitPriceQueue //从低到高
	symbol          string
	buyTradePlate   *TradePlate //买盘
	sellTradePlate  *TradePlate //卖盘
	kafkaClient     *database.KafkaClient
	orderDomain     *domain.ExchangeOrderDomain
}

func (t *CoinTrade) init() {
	t.buyTradePlate = NewTradePlate(t.symbol, model.BUY)
	t.sellTradePlate = NewTradePlate(t.symbol, model.SELL)
	t.buyLimitQueue = &LimitPriceQueue{}
	t.sellLimitQueue = &LimitPriceQueue{}
	t.initQueue()
}

// TradePlate 盘口信息
type TradePlate struct {
	Items     []*TradePlateItem `json:"items"`
	Symbol    string
	direction int
	maxDepth  int
	mux       sync.RWMutex
}

func (p *TradePlate) Add(order *model.ExchangeOrder) {
	if p.direction != order.Direction {
		logx.Error("买卖盘 direction not match，check code...")
		return
	}
	p.mux.Lock()
	defer p.mux.Unlock()
	if order.Type == model.MarketPrice {
		logx.Error("市价单 不加入买卖盘")
		return
	}
	size := len(p.Items)
	//检查items和当前的价格是否有一样的 有的话增加数量
	if size > 0 {
		for _, v := range p.Items {
			if (order.Direction == model.BUY && v.Price > order.Price) ||
				(order.Direction == model.SELL && v.Price < order.Price) {
				continue
			} else if v.Price == order.Price {
				v.Amount = tools.AddN(v.Amount, tools.SubFloor(order.Amount, order.TradedAmount, 5), 5)
				return
			} else {
				break
			}
		}
	}
	if size < p.maxDepth {
		tpi := &TradePlateItem{
			Amount: tools.SubFloor(order.Amount, order.TradedAmount, 5),
			Price:  order.Price,
		}
		p.Items = append(p.Items, tpi)
	}
}

type TradePlateResult struct {
	Direction    string            `json:"direction"`
	MaxAmount    float64           `json:"maxAmount"`
	MinAmount    float64           `json:"minAmount"`
	HighestPrice float64           `json:"highestPrice"`
	LowestPrice  float64           `json:"lowestPrice"`
	Symbol       string            `json:"symbol"`
	Items        []*TradePlateItem `json:"items"`
}

func (p *TradePlate) AllResult() *TradePlateResult {
	result := &TradePlateResult{}
	direction := model.DirectionMap.Value(p.direction)
	result.Direction = direction
	result.MaxAmount = p.getMaxAmount()
	result.MinAmount = p.getMinAmount()
	result.HighestPrice = p.getHighestPrice()
	result.LowestPrice = p.getLowestPrice()
	result.Symbol = p.Symbol
	result.Items = p.Items
	return result
}

func (p *TradePlate) Result(num int) *TradePlateResult {
	if num > len(p.Items) {
		num = len(p.Items)
	}
	result := &TradePlateResult{}
	direction := model.DirectionMap.Value(p.direction)
	result.Direction = direction
	result.MaxAmount = p.getMaxAmount()
	result.MinAmount = p.getMinAmount()
	result.HighestPrice = p.getHighestPrice()
	result.LowestPrice = p.getLowestPrice()
	result.Symbol = p.Symbol
	result.Items = p.Items[:num]
	return result
}

func (p *TradePlate) getMaxAmount() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var amount float64 = 0
	for _, v := range p.Items {
		if v.Amount > amount {
			amount = v.Amount
		}
	}
	return amount
}

func (p *TradePlate) getMinAmount() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var amount float64 = p.Items[0].Amount
	for _, v := range p.Items {
		if v.Amount < amount {
			amount = v.Amount
		}
	}
	return amount
}

func (p *TradePlate) getHighestPrice() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var price float64 = 0
	for _, v := range p.Items {
		if v.Price > price {
			price = v.Price
		}
	}
	return price
}
func (p *TradePlate) getLowestPrice() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var price float64 = p.Items[0].Price
	for _, v := range p.Items {
		if v.Price < price {
			price = v.Price
		}
	}
	return price
}

func (p *TradePlate) Remove(order *model.ExchangeOrder, amount float64) {
	for i, item := range p.Items {
		if item.Price == order.Price {
			item.Amount = tools.SubFloor(item.Amount, amount, 8)
			if item.Amount <= 0 {
				p.Items = append(p.Items[:i], p.Items[i+1:]...)
			}
			break
		}
	}
}

type TradePlateItem struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

func NewTradePlate(symbol string, direction int) *TradePlate {
	return &TradePlate{
		Symbol:    symbol,
		direction: direction,
		maxDepth:  100,
	}
}

func NewCoinTrade(symbol string, client *database.KafkaClient, db *msdb.MsDB) *CoinTrade {
	t := &CoinTrade{
		symbol:      symbol,
		kafkaClient: client,
		orderDomain: domain.NewExchangeOrderDomain(db),
	}
	t.init()
	return t
}

func (t *CoinTrade) Trade(order *model.ExchangeOrder) {

	var limitPriceList *LimitPriceQueue
	var marketPriceList TradeTimeQueue

	if order.Direction == model.BUY {
		limitPriceList = t.sellLimitQueue
		marketPriceList = t.sellMarketQueue
	} else {
		limitPriceList = t.buyLimitQueue
		marketPriceList = t.buyMarketQueue
	}

	if order.Type == model.MarketPrice {
		t.matchLimitPriceWithMP(marketPriceList, order)
	} else if order.Type == model.LimitPrice {
		//如果是限价单 先与限价单交易 在与市价单交易
		t.matchLimitPriceWithLP(limitPriceList, order)
		if order.Status == model.Trading {
			//证明还未交易完成 继续和市价单交易
			t.matchLimitPriceWithMP(marketPriceList, order)
		}
		if order.Status == model.Trading {
			t.addLimitPriceOrder(order)
			if order.Direction == model.BUY {
				t.sendTradePlateMsg(t.buyTradePlate)
			} else {
				t.sendTradePlateMsg(t.sellTradePlate)
			}
		}
	} else {
		//与限价单进行交易
		t.matchMarketPriceWithLP(limitPriceList, order)
	}
}

func (t *CoinTrade) GetTradePlate(direction int) *TradePlate {
	if direction == model.BUY {
		return t.buyTradePlate
	}
	return t.sellTradePlate
}

func (t *CoinTrade) sendTradePlateMsg(plate *TradePlate) {
	bytes, _ := json.Marshal(plate.Result(24))
	data := database.KafkaData{
		Topic: "exchange_order_trade_plate",
		Key:   []byte(plate.Symbol),
		Data:  bytes,
	}
	t.kafkaClient.SendSync(data)
}

func (t *CoinTrade) initQueue() {
	ctx := context.Background()
	list, err := t.orderDomain.FindTradingOrders(ctx)
	if err != nil {
		logx.Error(err)
		return
	}
	for _, v := range list {
		if v.Direction == model.BUY {
			//买
			if v.Type == model.MarketPrice {
				//市价买
				t.buyMarketQueue = append(t.buyMarketQueue, v)
			} else {
				isPut := false
				for _, bv := range t.buyLimitQueue.list {
					if bv.price == v.Price {
						bv.list = append(bv.list, v)
						isPut = true
						break
					}
				}
				if !isPut {
					plm := &LimitPriceMap{
						price: v.Price,
					}
					plm.list = append(plm.list, v)
					t.buyLimitQueue.list = append(t.buyLimitQueue.list, plm)
				}
				t.buyTradePlate.Add(v)
			}
		} else {
			//卖
			if v.Type == model.MarketPrice {
				//市价卖
				t.sellMarketQueue = append(t.sellMarketQueue, v)
			} else {
				isPut := false
				for _, bv := range t.sellLimitQueue.list {
					if bv.price == v.Price {
						bv.list = append(bv.list, v)
						isPut = true
						break
					}
				}
				if !isPut {
					plm := &LimitPriceMap{
						price: v.Price,
					}
					plm.list = append(plm.list, v)
					t.sellLimitQueue.list = append(t.sellLimitQueue.list, plm)
				}
				t.sellTradePlate.Add(v)
			}
		}
	}
	//排序
	sort.Sort(t.sellMarketQueue)
	sort.Sort(t.buyMarketQueue)
	sort.Sort(t.buyLimitQueue.list)
	sort.Sort(sort.Reverse(t.sellLimitQueue.list))
}

func (t *CoinTrade) matchMarketPriceWithLP(lpList *LimitPriceQueue, focusedOrder *model.ExchangeOrder) {
	lpList.mux.Lock()
	defer lpList.mux.Unlock()
	buyNotify := false
	sellNotify := false
	for _, v := range lpList.list {
		var delOrders []string
		for _, matchOrder := range v.list {
			if matchOrder.MemberId == focusedOrder.MemberId {
				//自己不与自己交易
				continue
			}
			//不管是买还是卖，如果匹配 那么 match从amount移动一部分到tradeAmount，同样focusOrder 从amount移动一部分到tradeAmount
			//focusedOrder是市价单 所以以matchOrder价格为主
			price := matchOrder.Price
			//可交易的数量
			matchAmount := tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
			focuseAmount := tools.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)
			if focusedOrder.Direction == model.BUY {
				//如果是市价买 amount是 USDT的数量 要计算买多少BTC 要根据match的price进行计算
				focuseAmount = tools.DivFloor(tools.SubFloor(focusedOrder.Amount, focusedOrder.Turnover, 8), price, 8)
			}
			if matchAmount >= focuseAmount {
				//能够进行匹配，直接完成即可
				matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, focuseAmount, 8)
				focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, focuseAmount, 8)
				to := tools.MulFloor(price, focuseAmount, 8)
				focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
				matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
				focusedOrder.Status = model.Completed
				if tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
					matchOrder.Status = model.Completed
					delOrders = append(delOrders, matchOrder.OrderId)
				}
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, focuseAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, focuseAmount)
					sellNotify = true
				}
				break
			} else {
				to := tools.MulFloor(price, matchAmount, 8)
				matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
				matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
				matchOrder.Status = model.Completed
				delOrders = append(delOrders, matchOrder.OrderId)
				focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
				focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
				//还得继续下一轮匹配
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, matchAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, matchAmount)
					sellNotify = true
				}
				continue
			}

		}
		for _, orderId := range delOrders {
			for index, order := range v.list {
				if order.OrderId == orderId {
					v.list = append(v.list[:index], v.list[index+1:]...)
					break
				}
			}
		}
	}
	//判断order完成了没，没完成 放入队列
	if focusedOrder.Status == model.Trading {
		t.addMarketPriceOrder(focusedOrder)
	}
	//通知买卖盘更新
	if buyNotify {
		t.sendTradePlateMsg(t.buyTradePlate)
	}
	if sellNotify {
		t.sendTradePlateMsg(t.sellTradePlate)
	}

}

func (t *CoinTrade) addMarketPriceOrder(order *model.ExchangeOrder) {
	if order.Type != model.MarketPrice {
		return
	}
	if order.Direction == model.BUY {
		t.buyMarketQueue = append(t.buyMarketQueue, order)
		sort.Sort(t.buyMarketQueue)
	} else {
		t.sellMarketQueue = append(t.sellMarketQueue, order)
		sort.Sort(t.sellMarketQueue)
	}
}

func (t *CoinTrade) addLimitPriceOrder(order *model.ExchangeOrder) {
	if order.Type != model.LimitPrice {
		return
	}
	if order.Direction == model.BUY {
		isPut := false
		for _, v := range t.buyLimitQueue.list {
			if v.price == order.Price {
				isPut = true
				v.list = append(v.list, order)
				break
			}
		}
		if !isPut {
			lm := &LimitPriceMap{
				price: order.Price,
				list:  []*model.ExchangeOrder{order},
			}
			t.buyLimitQueue.list = append(t.buyLimitQueue.list, lm)
		}
		sort.Sort(t.buyLimitQueue.list)
		t.buyTradePlate.Add(order)
	} else {
		isPut := false
		for _, v := range t.sellLimitQueue.list {
			if v.price == order.Price {
				isPut = true
				v.list = append(v.list, order)
				break
			}
		}
		if !isPut {
			lm := &LimitPriceMap{
				price: order.Price,
				list:  []*model.ExchangeOrder{order},
			}
			t.sellLimitQueue.list = append(t.sellLimitQueue.list, lm)
		}
		sort.Sort(sort.Reverse(t.sellLimitQueue.list))
		t.sellTradePlate.Add(order)
	}
}

func (t *CoinTrade) matchLimitPriceWithMP(mpList TradeTimeQueue, focusedOrder *model.ExchangeOrder) {
	//市价单时间是 从旧到新 先去匹配之前的单
	var delOrders []string
	for _, matchOrder := range mpList {
		if matchOrder.MemberId == focusedOrder.MemberId {
			//自己不与自己交易
			continue
		}
		price := focusedOrder.Price
		//可交易的数量
		matchAmount := tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
		focusedAmount := tools.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)
		if matchAmount >= focusedAmount {
			//能够进行匹配，直接完成即可
			matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, focusedAmount, 8)
			focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, focusedAmount, 8)
			to := tools.MulFloor(price, focusedAmount, 8)
			focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
			matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
			focusedOrder.Status = model.Completed
			if tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
				matchOrder.Status = model.Completed
				delOrders = append(delOrders, matchOrder.OrderId)
			}
			break
		} else {
			to := tools.MulFloor(price, matchAmount, 8)
			matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
			matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
			matchOrder.Status = model.Completed
			delOrders = append(delOrders, matchOrder.OrderId)
			focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
			focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
			//还得继续下一轮匹配
			continue
		}
	}
	for _, orderId := range delOrders {
		for index, order := range mpList {
			if order.OrderId == orderId {
				mpList = append(mpList[:index], mpList[index+1:]...)
				break
			}
		}
	}
}

func (t *CoinTrade) matchLimitPriceWithLP(lpList *LimitPriceQueue, focusedOrder *model.ExchangeOrder) {
	lpList.mux.Lock()
	defer lpList.mux.Unlock()
	buyNotify := false
	sellNotify := false
	var completeOrders []*model.ExchangeOrder
	for _, v := range lpList.list {
		var delOrders []string
		for _, matchOrder := range v.list {
			if matchOrder.MemberId == focusedOrder.MemberId {
				//自己不与自己交易
				continue
			}
			if focusedOrder.Direction == model.BUY {
				//买单 matchOrder为限价卖单 价格从低到高
				if matchOrder.Price > focusedOrder.Price {
					//最低卖价 比 买入价高 直接退出
					break
				}
			}
			if focusedOrder.Direction == model.SELL {
				if matchOrder.Price < focusedOrder.Price {
					//最高买价 比 卖价 低 直接退出
					break
				}
			}
			price := matchOrder.Price
			//可交易的数量
			matchAmount := tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
			focuseAmount := tools.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)
			if matchAmount <= 0 {
				//证明已经交易完成
				matchOrder.Status = model.Completed
				delOrders = append(delOrders, matchOrder.OrderId)
				completeOrders = append(completeOrders, matchOrder)
				continue
			}
			if matchAmount >= focuseAmount {
				//能够进行匹配，直接完成即可
				matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, focuseAmount, 8)
				focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, focuseAmount, 8)
				to := tools.MulFloor(price, focuseAmount, 8)
				focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
				matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
				focusedOrder.Status = model.Completed

				if tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
					//matchorder也完成了 需要从匹配列表中删除
					matchOrder.Status = model.Completed
					delOrders = append(delOrders, matchOrder.OrderId)
					completeOrders = append(completeOrders, matchOrder)
				}
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, focuseAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, focuseAmount)
					sellNotify = true
				}
				break
			} else {
				to := tools.MulFloor(price, matchAmount, 8)
				matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
				matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
				matchOrder.Status = model.Completed
				delOrders = append(delOrders, matchOrder.OrderId)
				completeOrders = append(completeOrders, matchOrder)
				focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
				focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
				//还得继续下一轮匹配
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, matchAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, matchAmount)
					sellNotify = true
				}
				continue
			}
		}
		for _, orderId := range delOrders {
			for index, order := range v.list {
				if order.OrderId == orderId {
					v.list = append(v.list[:index], v.list[index+1:]...)
					break
				}
			}
		}
	}
	//通知买卖盘更新
	if buyNotify {
		t.sendTradePlateMsg(t.buyTradePlate)
	}
	if sellNotify {
		t.sendTradePlateMsg(t.sellTradePlate)
	}
	t.onCompleteHandle(completeOrders)
}

func (t *CoinTrade) onCompleteHandle(orders []*model.ExchangeOrder) {
	if len(orders) <= 0 {
		return
	}
	for _, order := range orders {
		marshal, err := json.Marshal(order)
		if err != nil {
			logx.Error("封装已完成数据错误:", err)
			return
		}
		logx.Info("准备开始发已完成的数据")
		kafkaData := database.KafkaData{
			Topic: "exchange_order_completed",
			Key:   []byte(t.symbol),
			Data:  marshal,
		}
		for {
			//保证一定发成功
			err := t.kafkaClient.SendSync(kafkaData)
			if err != nil {
				logx.Error("更新订单错误:", err, "主题为exchange-asset-completed")
				time.Sleep(250 * time.Millisecond)
				continue
			}
		}
	}

}
