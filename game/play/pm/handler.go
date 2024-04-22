package pm

import (
	"east/pb"

	"github.com/tnnmigga/nett/conf"
	"github.com/tnnmigga/nett/infra/zlog"
	"github.com/tnnmigga/nett/mods/mongo"
	"github.com/tnnmigga/nett/msgbus"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *Manager) onCreatePlayerRPC(req *pb.CreatePlayerRPC, resolve func(any), reject func(err error)) {
	p := NewPlayer(req.UserID)
	b, _ := bson.Marshal(p)
	msgbus.Cast(&mongo.MongoSaveSingle{
		GroupKey: groupKey(req.UserID),
		CollName: "player",
		Op: &mongo.MongoSaveOp{
			Filter: bson.M{"_id": req.UserID},
			Value:  b,
		},
	})
	zlog.Infof("create new player: %d", req.UserID)
	resolve(&pb.CreatePlayerRPCRes{ServerID: conf.ServerID})
}
