package task

import (
	"github.com/go-co-op/gocron"
	"jobcenter/internal/market"
	"jobcenter/internal/svc"
	"time"
)

type Task struct {
	s   *gocron.Scheduler
	ctx *svc.ServiceContext
}

func NewTask(ctx *svc.ServiceContext) *Task {
	return &Task{
		s:   gocron.NewScheduler(time.UTC),
		ctx: ctx,
	}
}

func (t *Task) Run() {
	kline := market.NewKline(t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Config.Okx, t.ctx.RedisClient) //开启kafka写
	rate := market.NewRate(t.ctx)
	t.ctx.KafkaClient.StartWrite()
	t.s.Every(60).Seconds().Do(func() {
		kline.Do("1m")
		kline.Do("1H")
		kline.Do("30m")
		kline.Do("15m")
		kline.Do("5m")
		kline.Do("1D")
		kline.Do("1W")
		kline.Do("1M")
		rate.Do()
	})
	t.s.Every(10).Minute().Do(func() {
		market.NewBitCoin(t.ctx.RedisClient.Cache, t.ctx.AssetRPC, t.ctx.MongoClient, t.ctx.KafkaClient).Do()
	})
}

func (t *Task) StartBlocking() {
	t.s.StartBlocking()
}

func (t *Task) Stop() {
	t.s.Stop()
}
