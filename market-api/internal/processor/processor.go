package processor

import (
	"encoding/json"
	"fmt"
	"grpc-common/market/types/market"
	"market-api/internal/database"
	"market-api/internal/model"
)

const KLINE1M = "kline_1m"
const KLINE = "kline"
const TRADE = "trade"
const TradePlateTopic = "exchange_order_trade_plate"
const TradePlate = "tradePlate"

type ProcessData struct {
	Type string //trade 交易 kline k线
	Key  []byte
	Data []byte
}
type MarketHandler interface {
	HandleTrade(symbol string, data []byte)
	HandleKLine(symbol string, kline *model.Kline, thumbMap map[string]*market.CoinThumb)
	HandlerTradePlate(symbol string, plate *model.TradePlateResult)
}
type Processor interface {
	PutThumb(any)
	GetThumb() any
	Process(data *ProcessData)
	AddHandler(h MarketHandler)
}

type DefaultProcessor struct {
	kafkaCli *database.KafkaClient
	handlers []MarketHandler
	thumbMap map[string]*market.CoinThumb
}

func (p *DefaultProcessor) PutThumb(data any) {
	switch data.(type) {
	case []*market.CoinThumb:
		list := data.([]*market.CoinThumb)
		for _, v := range list {
			p.thumbMap[v.Symbol] = v
		}
	}
}

func (p *DefaultProcessor) GetThumb() any {
	return p.thumbMap
}

func NewDefaultProcessor(kafkaCli *database.KafkaClient) *DefaultProcessor {
	return &DefaultProcessor{
		kafkaCli: kafkaCli,
		handlers: make([]MarketHandler, 0),
		thumbMap: make(map[string]*market.CoinThumb),
	}
}

func (p *DefaultProcessor) Init(wh *WebsocketHandler) {
	p.AddHandler(wh)
	//接收kline 1m的同步数据
	go p.startReadFromKafka(KLINE1M, KLINE)
	p.startReadTradePlate(TradePlateTopic)
}

func (p *DefaultProcessor) Process(data *ProcessData) {
	fmt.Println(data.Type)
	if data.Type == KLINE {
		kline := &model.Kline{}
		json.Unmarshal(data.Data, kline)
		for _, v := range p.handlers {
			v.HandleKLine(string(data.Key), kline, p.thumbMap)
		}
	} else if data.Type == TradePlate {
		tp := &model.TradePlateResult{}
		json.Unmarshal(data.Data, tp)
		for _, v := range p.handlers {
			v.HandlerTradePlate(string(data.Key), tp)
		}
	}
}
func (p *DefaultProcessor) AddHandler(h MarketHandler) {
	p.handlers = append(p.handlers, h)
}

func (p *DefaultProcessor) startReadFromKafka(topic string, tp string) {
	p.kafkaCli.StartRead(topic)
	go p.dealQueueData(p.kafkaCli, tp)
}

func (p *DefaultProcessor) dealQueueData(cli *database.KafkaClient, tp string) {
	for {
		kafkaData, _ := cli.Read()
		pd := &ProcessData{
			Type: tp,
			Key:  kafkaData.Key,
			Data: kafkaData.Data,
		}
		p.Process(pd)
	}
}

func (p *DefaultProcessor) startReadTradePlate(topic string) {
	cli := p.kafkaCli.StartReadNew(topic)
	go p.dealQueueData(cli, TradePlate)
}
