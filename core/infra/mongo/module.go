package mongo

import (
	"context"
	"east/core/iconf"
	"east/core/idef"
	"east/core/infra"
	"east/core/module"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Module struct {
	*module.Module
	mongocli *mongo.Client
}

func New() idef.IModule {
	m := &Module{
		Module: module.New(infra.ModTypMongo, module.DefaultMQLen),
	}
	m.Before(idef.ServerStateRun, m.beforeRun)
	m.After(idef.ServerStateStop, m.afterStop)
	return m
}

func (m *Module) beforeRun() (err error) {
	m.mongocli, err = mongo.Connect(context.Background(), options.Client().ApplyURI(iconf.String("mongo-url", "mongodb://localhost")))
	if err != nil {
		return err
	}
	if err := m.mongocli.Ping(context.Background(), readpref.Primary()); err != nil {
		return err
	}
	return nil
}

func (m *Module) afterStop() (err error) {
	m.mongocli.Disconnect(context.Background())
	return nil
}
