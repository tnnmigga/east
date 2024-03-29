package mongo

import (
	"context"
	"east/core/basic"
	"east/core/log"
	"east/core/msgbus"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *module) registerHandler() {
	msgbus.RegisterHandler(m, m.onMongoSave)
	msgbus.RegisterRPC(m, m.onMongoLoad)
}

func (m *module) onMongoSave(req *MongoSave) {
	ms := make([]mongo.WriteModel, 0, len(req.Ops))
	for _, op := range req.Ops {
		b, err := bson.Marshal(op.Value)
		if err != nil {
			log.Errorf("save %#v bson error %v", op.Value, err)
			continue
		}
		m := mongo.NewReplaceOneModel().SetFilter(op.Filter).SetReplacement(b).SetUpsert(true)
		ms = append(ms, m)
	}
	basic.GoWithGroup(req.Key(), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		_, err := m.mongocli.Database(req.DBName).Collection(req.CollName).BulkWrite(ctx, ms)
		cancel()
		if err != nil {
			log.Errorf("mongo save error %v", err)
		}
	})
}

func (m *module) onMongoLoad(req *MongoLoad, resolve func(any), reject func(error)) {
	basic.GoWithGroup(req.Key(), func() {
		cur, _ := m.mongocli.Database(req.DBName).Collection(req.CollName).Find(context.Background(), req.Filter)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		err := cur.All(ctx, &req.Data)
		cancel()
		if err != nil {
			reject(err)
		} else {
			resolve(req.Data)
		}
	})
}
