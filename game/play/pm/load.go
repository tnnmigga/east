package pm

import (
	"fmt"
	"time"

	"github.com/tnnmigga/nett/conc"
	"github.com/tnnmigga/nett/infra/cluster"
	"github.com/tnnmigga/nett/infra/zlog"
	"github.com/tnnmigga/nett/mods/mongo"
	"github.com/tnnmigga/nett/msgbus"
	"github.com/tnnmigga/nett/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func LoadAsync(uid uint64, cb func(*Player, error)) {
	if p, ok := manager.cache[uid]; ok {
		cb(p, nil)
		return
	}
	manager.waiting[uid] = append(manager.waiting[uid], cb)
	cbs := manager.waiting[uid]
	if len(cbs) > 1 {
		return
	}
	conc.Async(manager, func() (*cluster.Lock, error) {
		lock := cluster.NewLock(fmt.Sprintf("userdata.%d", uid))
		err := lock.Wait(10 * time.Second) // or TryLock
		if err != nil {
			zlog.Errorf("wait userdata lock error %v", err)
			lock.Release()
			return nil, err
		}
		return lock, nil
	}, func(l *cluster.Lock, err error) {
		msgbus.RPC(manager, msgbus.Local(), &mongo.MongoLoadSingle{
			GroupKey: groupKey(uid),
			CollName: "player",
			Filter:   bson.M{"_id": uid},
		}, func(raw bson.Raw, err error) {
			p := &Player{}
			if err == nil {
				err = bson.Unmarshal(raw, p)
			}
			if err == nil {
				manager.cache[uid] = p
			}
			zlog.Debugf("load player %d %v", uid, err)
			cbs := manager.waiting[uid]
			for _, cb := range cbs {
				func() {
					defer utils.RecoverPanic()
					cb(p, err)
				}()
			}
		})
	})
}
