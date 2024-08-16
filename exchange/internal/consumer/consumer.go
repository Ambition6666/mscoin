package consumer

import (
	"common/msdb"
	"encoding/json"
	"exchange/database"
	"exchange/internal/config"
	"exchange/internal/domain"
	"exchange/internal/model"
	"exchange/internal/processor"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
	"time"
)

type KafkaConsumer struct {
	c                config.KafkaConfig
	coinTradeFactory *processor.CoinTradeFactory
	orderDomain      *domain.ExchangeOrderDomain
}

func NewKafkaConsumer(
	c config.KafkaConfig,
	factory *processor.CoinTradeFactory,
	db *msdb.MsDB) *KafkaConsumer {
	return &KafkaConsumer{
		c:                c,
		coinTradeFactory: factory,
		orderDomain:      domain.NewExchangeOrderDomain(db),
	}
}

func (k *KafkaConsumer) Run() {
	k.orderTrading()
}
func (k *KafkaConsumer) orderTrading() {
	//topic exchange_order_trading
	client := database.NewKafkaClient(k.c)
	client = client.StartRead("exchange_order_trading")
	go k.readOrderTrading(client)

	client = client.StartRead("exchange_order_completed")
	go k.readOrderComplete(client)
}

func (k *KafkaConsumer) readOrderTrading(client *database.KafkaClient) {
	for {
		kafkaData, _ := client.Read()
		order := new(model.ExchangeOrder)
		json.Unmarshal(kafkaData.Data, &order)
		coinTrade := k.coinTradeFactory.GetCoinTrade(order.Symbol)
		coinTrade.Trade(order)
	}
}

func (k *KafkaConsumer) readOrderComplete(client *database.KafkaClient) {
	for {
		kafkaData, _ := client.Read()
		var order *model.ExchangeOrder
		json.Unmarshal(kafkaData.Data, &order)
		logx.Info("读取已完成订单数据成功:", order.OrderId)
		//更新订单
		err := k.orderDomain.UpdateOrderComplete(context.Background(), order)
		if err != nil {
			logx.Error(err)
			//如果失败重新放入  再次进行更新
			client.RPut(kafkaData)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		//继续放入MQ中 通知钱包服务更新
		for {
			kafkaData.Topic = "exchange_order_complete_update_success"
			err := client.SendSync(kafkaData)
			if err != nil {
				//确保一定发送成功
				logx.Error("钱包服务更新失败:", err)
				time.Sleep(250 * time.Millisecond)
				continue
			}
		}
	}
}
