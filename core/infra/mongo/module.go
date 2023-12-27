package mongo

import (
	"context"
	"east/core/idef"
	"east/core/module"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Module struct {
	*module.Module
	mongocli *mongo.Client
}

func New(name, url string) idef.IModule {
	cli, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	return &Module{
		Module:   module.New(name, module.DefaultMQLen),
		mongocli: cli,
	}
}

func (m *Module) Run() {
	if err := m.mongocli.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(err)
	}
	m.Module.Run()
}
