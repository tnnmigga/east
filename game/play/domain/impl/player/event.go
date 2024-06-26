package player

import (
	"east/define"
	"east/pb"

	"github.com/tnnmigga/core/infra/eventbus"
	"github.com/tnnmigga/core/infra/zlog"
	"github.com/tnnmigga/core/mods/mongo"
	"github.com/tnnmigga/core/msgbus"
	"github.com/tnnmigga/core/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func (c *useCase) onEventUserMsg(event *eventbus.Event) {
	zlog.Infof("user event %v", utils.String(event))
	b, _ := bson.Marshal(bson.M{
		"test": 11,
	})
	msgbus.Cast(&mongo.MongoSaveMulti{
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
	msgbus.RPC(c, msgbus.Local(), &mongo.MongoLoadMulti{
		CollName: "test",
		Filter:   bson.M{},
	}, func(raws []bson.Raw, err error) {
		ms := make([]bson.M, len(raws))
		for i, raw := range raws {
			bson.Unmarshal(raw, &ms[i])
		}
		zlog.Info("db load", ms, err)
	})
	msgbus.RPC(c, msgbus.ServerType(define.ServGateway), &pb.TestRPC{}, func(resp *pb.TestRPCRes, err error) {
		zlog.Info("remote call", resp, err)
	})
	msgbus.RPC(c, msgbus.Local(), &pb.TestRPC{}, func(resp *pb.TestRPCRes, err error) {
		zlog.Info("local call", resp, err)
	})
}
