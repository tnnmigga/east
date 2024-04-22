package pm

import (
	"east/pb"
	"reflect"

	"github.com/tnnmigga/nett/codec"
	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/infra/zlog"
	"github.com/tnnmigga/nett/msgbus"
	"github.com/tnnmigga/nett/utils"
)

var manager *Manager

type Manager struct {
	idef.IModule
	cache   map[uint64]*Player
	waiting map[uint64][]func(*Player, error)
}

func Init(m idef.IModule) {
	manager = &Manager{
		IModule: m,
		cache:   map[uint64]*Player{},
		waiting: map[uint64][]func(*Player, error){},
	}
	msgbus.RegisterHandler(m, manager.onC2SPackage)
	msgbus.RegisterRPC(m, manager.onCreatePlayerRPC)
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
