package domain

import (
	"context"
	"jobcenter/internal/dao"
	"jobcenter/internal/database"
	"jobcenter/internal/model"
	"jobcenter/internal/repo"
	"log"
)

type KlineDomain struct {
	klineRepo repo.KlineRepo
}

func (d *KlineDomain) Save(data [][]string, symbol, period string) {
	klines := make([]*model.Kline, len(data))
	for i, v := range data {
		kline := model.NewKline(v, period)
		klines[i] = kline
	}
	//拿到获取的最后时间的数据，然后将大于此时间的数据删除
	kline := klines[len(klines)-1]
	dataTime := kline.Time
	err := d.klineRepo.DeleteGtTime(context.Background(), dataTime, symbol, period)
	if err != nil {
		log.Println(err)
		return
	}
	err = d.klineRepo.SaveBatch(context.Background(), klines, symbol, period)
	if err != nil {
		log.Println(err)
	}
}

func NewKlineDomain(cli *database.MongoClient) *KlineDomain {
	return &KlineDomain{
		klineRepo: dao.NewKlineDao(cli.Db),
	}
}
