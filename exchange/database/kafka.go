package database

import (
	"context"
	"errors"
	"exchange/internal/config"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"sync"
	"time"
)

type KafkaData struct {
	Topic string
	Key   []byte
	Data  []byte
}
type KafkaClient struct {
	w         *kafka.Writer
	r         *kafka.Reader
	rDataChan chan KafkaData
	data      chan KafkaData
	c         config.KafkaConfig
	closed    bool
	mutex     sync.Mutex
}

func NewKafkaClient(c config.KafkaConfig) *KafkaClient {
	return &KafkaClient{
		c: c,
	}
}

func (k *KafkaClient) StartWrite() {
	w := &kafka.Writer{
		Addr:     kafka.TCP(k.c.Addr),
		Balancer: &kafka.LeastBytes{},
	}
	k.w = w
	k.data = make(chan KafkaData, k.c.WriteCap)
	go k.sendKafka()
}

func (w *KafkaClient) Send(data KafkaData) {
	defer func() {
		if err := recover(); err != nil {
			w.closed = true
		}
	}()
	w.data <- data
	w.closed = false
}

func (k *KafkaClient) SendSync(data KafkaData) error {
	w := &kafka.Writer{
		Addr:     kafka.TCP(k.c.Addr),
		Balancer: &kafka.LeastBytes{},
	}
	k.w = w
	messages := []kafka.Message{
		{
			Topic: data.Topic,
			Key:   data.Key,
			Value: data.Data,
		},
	}
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = k.w.WriteMessages(ctx, messages...)
	return err
}

func (w *KafkaClient) Close() {
	if w.w != nil {
		w.w.Close()
		w.mutex.Lock()
		defer w.mutex.Unlock()
		if !w.closed {
			close(w.data)
			w.closed = true
		}
	}
	if w.r != nil {
		w.r.Close()
	}
}

func (w *KafkaClient) sendKafka() {
	for {
		select {
		case data := <-w.data:
			messages := []kafka.Message{
				{
					Topic: data.Topic,
					Key:   data.Key,
					Value: data.Data,
				},
			}
			var err error
			const retries = 3
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			success := false
			for i := 0; i < retries; i++ {
				// attempt to create topic prior to publishing the message
				err = w.w.WriteMessages(ctx, messages...)
				if err == nil {
					success = true
					break
				}
				if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
					time.Sleep(time.Millisecond * 250)
					success = false
					continue
				}
				if err != nil {
					success = false
					log.Printf("kafka send writemessage err %s \n", err.Error())
				}
			}
			if !success {
				//重新放进去等待消费
				w.Send(data)
			}
		}
	}

}

func (k *KafkaClient) StartRead(topic string) *KafkaClient {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.c.Addr},
		GroupID:  k.c.Group,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	client := NewKafkaClient(k.c)
	client.r = r
	client.rDataChan = make(chan KafkaData, k.c.ReadCap)
	go client.readMsg()
	return client
}

func (k *KafkaClient) readMsg() {
	for {
		m, err := k.r.ReadMessage(context.Background())
		if err != nil {
			logx.Error(err)
			continue
		}
		data := KafkaData{
			Key:  m.Key,
			Data: m.Value,
		}
		k.rDataChan <- data
	}
}
func (k *KafkaClient) Read() (KafkaData, error) {
	msg := <-k.rDataChan
	return msg, nil
}

func (k *KafkaClient) RPut(data KafkaData) {
	k.rDataChan <- data
}
