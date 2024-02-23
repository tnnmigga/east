package mongo

import (
	"context"
	"east/core/basic"
	"east/core/conf"
	"east/core/idef"
	"east/core/infra"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type module struct {
	*basic.Module
	mongocli *mongo.Client
}

func New() idef.IModule {
	m := &module{
		Module: basic.New(infra.ModNameMongo, basic.DefaultMQLen),
	}
	m.registerHandler()
	m.Before(idef.ServerStateRun, m.beforeRun)
	m.After(idef.ServerStateStop, m.afterStop)
	return m
}

func (m *module) beforeRun() (err error) {
	m.mongocli, err = mongo.Connect(context.Background(), options.Client().ApplyURI(conf.String("mongo.url", "mongodb://localhost")))
	if err != nil {
		return err
	}
	if err := m.mongocli.Ping(context.Background(), readpref.Primary()); err != nil {
		return err
	}
	return nil
}

func (m *module) afterStop() (err error) {
	m.mongocli.Disconnect(context.Background())
	return nil
}
