package userimpl

import (
	"east/core/eventbus"
	"east/core/iconf"
	"east/core/infra/mongo"
	"east/core/log"
	"east/core/msgbus"
	"east/core/util"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *service) onEventUserMsg(event *eventbus.Event) {
	log.Infof("user event %v", util.String(event))
	msgbus.Cast(iconf.ServerID(), &mongo.MongoSave{
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
	msgbus.AsyncCall(s, iconf.ServerID(), &mongo.MongoLoad{
		DBName:   "test",
		CollName: "test",
		Filter:   bson.M{},
	}, func(res int, err error) {
		log.Info(res, err)
	})
}
