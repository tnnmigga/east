package mongo

import (
	"context"
	"time"

	"github.com/tnnmigga/nett/basic"
	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/infra"

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	m.mongocli, err = mongo.Connect(ctx, options.Client().ApplyURI(conf.String("mongo.url", "mongodb://localhost")))
	if err != nil {
		return err
	}
	if err := m.mongocli.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}

func (m *module) afterStop() (err error) {
	m.mongocli.Disconnect(context.Background())
	return nil
}
