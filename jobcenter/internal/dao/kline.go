package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"jobcenter/internal/model"
	"log"
)

type KlineDao struct {
	db *mongo.Database
}

func (d *KlineDao) DeleteGtTime(ctx context.Context, time int64, symbol string, period string) error {
	collection := d.db.Collection("exchange_kline_" + symbol + "_" + period)
	deleteResult, err := collection.DeleteMany(ctx, bson.D{{"time", bson.D{{"$gte", time}}}})
	log.Printf("删除表%s，数量为：%d \n", "exchange_kline_"+symbol+"_"+period, deleteResult.DeletedCount)
	return err
}

func NewKlineDao(db *mongo.Database) *KlineDao {
	return &KlineDao{
		db: db,
	}
}
func (d *KlineDao) SaveBatch(ctx context.Context, data []*model.Kline, symbol, period string) error {
	collection := d.db.Collection("exchange_kline_" + symbol + "_" + period)
	ds := make([]interface{}, len(data))
	for i, v := range data {
		ds[i] = v
	}
	_, err := collection.InsertMany(ctx, ds)
	return err
}
