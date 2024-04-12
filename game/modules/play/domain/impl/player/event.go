package player

import (
	"east/define"
	"east/pb"

	"github.com/tnnmigga/nett/eventbus"
	"github.com/tnnmigga/nett/infra/mongo"
	"github.com/tnnmigga/nett/msgbus"
	"github.com/tnnmigga/nett/util"
	"github.com/tnnmigga/nett/zlog"

	"go.mongodb.org/mongo-driver/bson"
)

func (c *useCase) onEventUserMsg(event *eventbus.Event) {
	zlog.Infof("user event %v", util.String(event))
	b, _ := bson.Marshal(bson.M{
		"test": 11,
	})
	msgbus.Cast(&mongo.MongoSave{
		DBName:   "test",
		CollName: "test",
		Ops: []*mongo.MongoSaveOp{
			{
				Filter: bson.M{
					"_id": 1,
				},
				Value: b,
			},
		},
	})
	msgbus.RPC(c, msgbus.Local(), &mongo.MongoLoad{
		DBName:   "test",
		CollName: "test",
		Filter:   bson.M{},
		Data:     []bson.M{},
	}, func(res any, err error) {
		zlog.Info("db load", res, err)
	})
	msgbus.RPC(c, msgbus.ServerType(define.ServGateway), &pb.TestRPC{}, func(resp *pb.TestRPCRes, err error) {
		zlog.Info("remote call", resp, err)
	})
	msgbus.RPC(c, msgbus.Local(), &pb.TestRPC{}, func(resp *pb.TestRPCRes, err error) {
		zlog.Info("local call", resp, err)
	})
}
