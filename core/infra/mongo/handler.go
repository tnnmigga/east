package mongo

import (
	"context"
	"east/core/log"
	"east/core/msgbus"
	"east/core/sys"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (com *component) registerHandler() {
	msgbus.RegisterHandler(com, com.onMongoSave)
	msgbus.RegisterRPC(com, com.onMongoLoad)
}

func (com *component) onMongoSave(req *MongoSave) {
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
	sys.GoWithGroup(req.Key(), func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := com.mongocli.Database(req.DBName).Collection(req.CollName).BulkWrite(ctx, ms)
		cancel()
		if err != nil {
			log.Errorf("mongo save error %v", err)
		}
	})
}

func (com *component) onMongoLoad(req *MongoLoad, resolve func(any), reject func(error)) {
	sys.GoWithGroup(req.Key(), func() {
		cur, _ := com.mongocli.Database(req.DBName).Collection(req.CollName).Find(context.Background(), req.Filter)
		res := []bson.M{}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := cur.All(ctx, &res)
		cancel()
		if err != nil {
			reject(err)
		} else {
			resolve(res)
		}
	})
}
