package player

import (
	"east/define"
	"east/game/play/pm"
	"east/pb"
	"time"

	"github.com/tnnmigga/core/conc"
	"github.com/tnnmigga/core/infra/eventbus"
	"github.com/tnnmigga/core/infra/https"
	"github.com/tnnmigga/core/infra/zlog"
	"github.com/tnnmigga/core/msgbus"
	"github.com/tnnmigga/core/utils"
)

func (c *useCase) onSayHelloReq(p *pm.Player, req *pb.SayHelloReq) {
	zlog.Infof("client say hello %v", req.Text)
	msgbus.Cast(&eventbus.Event{
		OwnerID: 1,
		Topic:   define.EventUserSayHello,
		Value:   1,
	})
	c.TimerCase().New(time.Second*2, &timerSayHello{
		UserID: 1,
		Text:   "hello client!",
	})
	zlog.Infof("http get 1 %f", utils.NowSec())
	conc.Async(c, func() ([]byte, error) {
		zlog.Infof("http get 2 %f", utils.NowSec())
		res, err := https.HttpGet("https://www.baidu.com")
		zlog.Infof("http get 3 %f", utils.NowSec())
		return res, err
	}, func(b []byte, err error) {
		if err != nil {
			zlog.Errorf("http get error %v", err)
			return
		}
		zlog.Infof("http get 4 %f", utils.NowSec())
		zlog.Infof("http get res %v", string(b))
	})
}

func (c *useCase) onRPCTest(req *pb.TestRPC, resolve func(any), reject func(error)) {
	resolve(&pb.TestRPCRes{
		V: 22,
	})
}
