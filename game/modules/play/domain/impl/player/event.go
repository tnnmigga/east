package player

import (
	"east/core/conf"
	"east/core/eventbus"
	"east/core/infra/mongo"
	"east/core/log"
	"east/core/msgbus"
	"east/core/util"
	"east/pb"

	"go.mongodb.org/mongo-driver/bson"
)

func (c *useCase) onEventUserMsg(event *eventbus.Event) {
	log.Infof("user event %v", util.String(event))
	msgbus.CastLocal(&mongo.MongoSave{
		DBName:   "test",
		CollName: "test",
		Ops: []*mongo.MongoSaveOp{
			{
				Filter: bson.M{
					"_id": 1,
				},
				Value: bson.M{
					"test": 1,
				},
			},
		},
	})
	msgbus.RPC(c, conf.ServerID(), &mongo.MongoLoad{
		DBName:   "test",
		CollName: "test",
		Filter:   bson.M{},
		Data:     []bson.M{},
	}, func(res any, err error) {
		log.Info("db load", res, err)
	})
	msgbus.RPC(c, 1888, &pb.TestRPC{}, func(resp *pb.TestRPCRes, err error) {
		log.Info("remote call", resp, err)
	})
	msgbus.RPC(c, conf.ServerID(), &pb.TestRPC{}, func(resp *pb.TestRPCRes, err error) {
		log.Info("local call", resp, err)
	})
}
