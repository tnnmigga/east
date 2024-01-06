package userimpl

import (
	"east/core/conf"
	"east/core/eventbus"
	"east/core/infra/mongo"
	"east/core/log"
	"east/core/msgbus"
	"east/core/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *service) onEventUserMsg(event *eventbus.Event) {
	log.Infof("user event %v", util.String(event))
	time.Sleep(time.Second * 10)
	msgbus.Cast(conf.ServerID(), &mongo.MongoSave{
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
	msgbus.RPC(s, conf.ServerID(), &mongo.MongoLoad{
		DBName:   "test",
		CollName: "test",
		Filter:   bson.M{},
	}, func(res any, err error) {
		log.Info(res, err)
	})
}
