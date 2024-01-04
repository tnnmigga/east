package mongo

import (
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

type MongoLoad struct {
	DBName   string
	CollName string
	Filter   bson.M
}
