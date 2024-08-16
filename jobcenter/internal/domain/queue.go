package domain

import (
	"encoding/json"
	"jobcenter/internal/database"
	"jobcenter/internal/model"
)

const KLINE1M = "kline_1m"

type QueueDomain struct {
	cli *database.KafkaClient
}

func (d *QueueDomain) Sync1mKline(data []string, symbol string, period string) {
	kline := model.NewKline(data, period)
	bytes, _ := json.Marshal(kline)
	sendData := database.KafkaData{
		Topic: KLINE1M,
		Key:   []byte(symbol),
		Data:  bytes,
	}
	d.cli.Send(sendData)
}

func (d *QueueDomain) SendRecharge(value float64, address string, time int64) {
	data := make(map[string]any)
	data["value"] = value
	data["address"] = address
	data["time"] = time
	data["type"] = model.RECHARGE
	data["symbol"] = "BTC"
	marshal, _ := json.Marshal(data)
	msg := database.KafkaData{
		Topic: "BtcTransactionTopic",
		Data:  marshal,
		Key:   []byte(address),
	}
	d.cli.Send(msg)
}

func NewQueueDomain(cli *database.KafkaClient) *QueueDomain {
	return &QueueDomain{
		cli: cli,
	}
}
