package mongo

import (
	"context"
	"east/core/log"
	"east/core/msgbus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *module) registerHandler() {
	msgbus.RegisterHandler(m, m.onMongoSave)
	msgbus.RegisterRPC(m, m.onMongoLoad)
}

func (m *module) onMongoSave(msg *MongoSave) {
	ms := make([]mongo.WriteModel, 0, len(msg.Ops))
	for _, op := range msg.Ops {
		b, err := bson.Marshal(op.Value)
		if err != nil {
			log.Errorf("save %#v bson error %v", op.Value, err)
			continue
		}
		m := mongo.NewReplaceOneModel().SetFilter(op.Filter).SetReplacement(b).SetUpsert(true)
		ms = append(ms, m)
	}
	res, err := m.mongocli.Database(msg.DBName).Collection(msg.CollName).BulkWrite(context.Background(), ms)
	log.Info(res, err)
}

func (m *module) onMongoLoad(msg *MongoLoad, resolve func(any), reject func(error)) {
	cur, _ := m.mongocli.Database(msg.DBName).Collection(msg.CollName).Find(context.Background(), msg.Filter)
	var res any
	cur.Decode(res)
	resolve(res)
}
