package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"market/internal/model"
)

type KlineDao struct {
	db *mongo.Database
}

func (d *KlineDao) FindBySymbolTime(ctx context.Context, symbol, period string, from, to int64, sort string) (list []*model.Kline, err error) {
	collection := d.db.Collection("exchange_kline_" + symbol + "_" + period)
	//1是升序 -1 是降序
	s := -1
	if "asc" == sort {
		s = 1
	}
	cur, err := collection.Find(ctx, bson.D{{Key: "time", Value: bson.D{{"$gte", from}, {"$lte", to}}}}, &options.FindOptions{
		Sort: bson.D{{"time", s}},
	})
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, &list)
	if err != nil {
		return
	}
	return
}

func NewKlineDao(db *mongo.Database) *KlineDao {
	return &KlineDao{
		db: db,
	}
}

func (d *KlineDao) FindBySymbol(ctx context.Context, symbol, period string, count int64) (list []*model.Kline, err error) {
	collection := d.db.Collection("exchange_kline_" + symbol + "_" + period)
	//1是升序 -1 是降序
	cur, err := collection.Find(ctx, bson.D{{}}, &options.FindOptions{
		Limit: &count,
		Sort:  bson.D{{"time", -1}},
	})
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, &list)
	if err != nil {
		return
	}
	return
}
