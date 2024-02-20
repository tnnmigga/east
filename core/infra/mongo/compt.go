package mongo

import (
	"context"
	"east/core/compt"
	"east/core/conf"
	"east/core/idef"
	"east/core/infra"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type component struct {
	*compt.Module
	mongocli *mongo.Client
}

func New() idef.IComponent {
	com := &component{
		Module: compt.New(infra.ModTypMongo, compt.DefaultMQLen),
	}
	com.registerHandler()
	com.Before(idef.ServerStateRun, com.beforeRun)
	com.After(idef.ServerStateStop, com.afterStop)
	return com
}

func (com *component) beforeRun() (err error) {
	com.mongocli, err = mongo.Connect(context.Background(), options.Client().ApplyURI(conf.String("mongo-url", "mongodb://localhost")))
	if err != nil {
		return err
	}
	if err := com.mongocli.Ping(context.Background(), readpref.Primary()); err != nil {
		return err
	}
	return nil
}

func (com *component) afterStop() (err error) {
	com.mongocli.Disconnect(context.Background())
	return nil
}
