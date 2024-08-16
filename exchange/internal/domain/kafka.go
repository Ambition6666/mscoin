package domain

import (
	"encoding/json"
	"exchange/database"
	"exchange/internal/model"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
	"time"
)

type KafkaDomain struct {
	kafkaClient *database.KafkaClient
	orderDomain *ExchangeOrderDomain
}

func (d *KafkaDomain) Send(topic string, userId int64, orderId string, money float64, symbol string, direction int, baseSymbol string, coinSymbol string, status int) bool {
	m := make(map[string]any)
	m["userId"] = userId
	m["orderId"] = orderId
	m["money"] = money
	m["symbol"] = symbol
	m["direction"] = direction
	m["baseSymbol"] = baseSymbol
	m["coinSymbol"] = coinSymbol
	m["status"] = status
	marshal, _ := json.Marshal(m)
	data := database.KafkaData{
		Topic: topic,
		Key:   []byte(orderId),
		Data:  marshal,
	}
	logx.Info(string(marshal))
	err := d.kafkaClient.SendSync(data)
	if err != nil {
		logx.Error("kafka发送失败:", err)
		return false
	}
	return true
}

type AddOrderResult struct {
	UserId  int64  `json:"userId"`
	OrderId string `json:"orderId"`
}

func (d *KafkaDomain) WaitAddOrderResult(topic string) {
	cli := d.kafkaClient.StartRead(topic)
	for {
		data, _ := cli.Read()
		logx.Info("收到订单增加结果:" + string(data.Data))
		var result AddOrderResult
		json.Unmarshal(data.Data, &result)

		order, err := d.orderDomain.exchangeOrderRepo.FindByOrderId(context.Background(), result.OrderId)
		if err != nil {
			logx.Error(err)
			err := d.orderDomain.UpdateOrderStatus(context.Background(), result.OrderId, model.Canceled)
			if err != nil {
				logx.Error("更新状态失败:", err)
				cli.RPut(data)
				return
			}
			continue
		}
		err = d.orderDomain.UpdateOrderStatus(context.Background(), result.OrderId, model.Trading)
		if err != nil {
			logx.Error("更新状态失败:", err)
			cli.RPut(data)
			continue
		}
		if order.Status != model.Init {
			logx.Error("订单已经被处理过", order.Status)
			continue
		}
		//订单初始化完成 发送消息到kafka 等待撮合交易引擎进行交易撮合
		for {
			bytes, _ := json.Marshal(order)
			orderData := database.KafkaData{
				Topic: "exchange_order_trading",
				Key:   []byte(order.OrderId),
				Data:  bytes,
			}
			err := cli.SendSync(orderData)
			if err != nil {
				logx.Error("kafka发送失败:", err)
				time.Sleep(250 * time.Millisecond)
				continue
			}
			logx.Info("订单创建成功，发送创建成功消息:", order.OrderId)
			break
		}

	}
}

func NewKafkaDomain(kafkaClient *database.KafkaClient, orderDomain *ExchangeOrderDomain) *KafkaDomain {
	k := &KafkaDomain{
		kafkaClient: kafkaClient,
		orderDomain: orderDomain,
	}
	go k.WaitAddOrderResult("exchange_order_init_complete")
	return k
}
