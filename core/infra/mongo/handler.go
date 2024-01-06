package mongo

import (
	"context"
	"east/core/log"
	"east/core/msgbus"

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
	res, err := com.mongocli.Database(req.DBName).Collection(req.CollName).BulkWrite(context.Background(), ms)
	log.Info(res, err)
}

func (com *component) onMongoLoad(msg *MongoLoad, resolve func(any), reject func(error)) {
	cur, _ := com.mongocli.Database(msg.DBName).Collection(msg.CollName).Find(context.Background(), msg.Filter)
	res := []bson.M{}
	err := cur.All(context.Background(), &res)
	if err != nil {
		reject(err)
		return
	}
	resolve(res)
}
