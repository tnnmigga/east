package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoSaveOp struct {
	Filter bson.M
	Value  any
}

type MongoSave struct {
	DBName   string
	CollName string
	Ops      []*MongoSaveOp
}

func (m *MongoSave) Key() string {
	return fmt.Sprintf("%s-%s", m.DBName, m.CollName)
}

type MongoLoad struct {
	DBName   string
	CollName string
	Filter   bson.M
}

func (m *MongoLoad) Key() string {
	return fmt.Sprintf("%s-%s", m.DBName, m.CollName)
}
