package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

type MongoOp struct {
	Filter bson.M
	Value  any
}

type MongoSave struct {
	DBName   string
	CollName string
	Ops      []*MongoOp
}

type MongoLoad struct {
	DBName   string
	CollName string
	Filter   bson.M
	
}
