package pm

import (
	"east/pb"
	"reflect"

	"github.com/tnnmigga/core/codec"
	"github.com/tnnmigga/core/conf"
	"github.com/tnnmigga/core/infra/zlog"
	"github.com/tnnmigga/core/mods/mongo"
	"github.com/tnnmigga/core/msgbus"
	"github.com/tnnmigga/core/utils"
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

func (m *Manager) onC2SPackage(req *pb.C2SPackage) {
	pkg, err := codec.Decode(req.Body)
	if err != nil {
		zlog.Errorf("onC2SPackage, msg decode error %v", err)
		return
	}
	zlog.Infof("onC2SPackage, recv user %d msg %s", req.UserID, utils.String(pkg))
	uid := req.UserID
	h, ok := msgHandler[reflect.TypeOf(pkg)]
	if !ok {
		zlog.Errorf("onC2SPackage, %s handler not found", utils.TypeName(pkg))
		return
	}
	LoadAsync(uid, func(p *Player, err error) {
		if err != nil {
			zlog.Errorf("onC2SPackage, load player error: %v", err)
			return
		}
		h(p, pkg)
	})
}
