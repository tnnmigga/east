package player

import (
	"east/define"
	"east/pb"
	"time"

	"github.com/tnnmigga/nett/core"
	"github.com/tnnmigga/nett/eventbus"
	"github.com/tnnmigga/nett/msgbus"
	"github.com/tnnmigga/nett/util"
	"github.com/tnnmigga/nett/zlog"
)

func (c *useCase) onSayHelloReq(msg *pb.SayHelloReq) {
	zlog.Infof("client say hello %v", msg.Text)
	msgbus.Cast(&eventbus.Event{
		OwnerID: 1,
		Topic:   define.EventUserSayHello,
		Value:   1,
	})
	c.TimerCase().New(time.Second*2, &timerSayHello{
		UserID: 1,
		Text:   "hello client!",
	})
	zlog.Infof("http get 1 %f", util.NowSec())
	core.Async(c, func() ([]byte, error) {
		zlog.Infof("http get 2 %f", util.NowSec())
		res, err := util.HttpGet("https://www.baidu.com")
		zlog.Infof("http get 3 %f", util.NowSec())
		return res, err
	}, func(b []byte, err error) {
		if err != nil {
			zlog.Errorf("http get error %v", err)
			return
		}
		zlog.Infof("http get 4 %f", util.NowSec())
		zlog.Infof("http get res %v", string(b))
	})
}

func (c *useCase) onRPCTest(req *pb.TestRPC, resolve func(any), reject func(error)) {
	resolve(&pb.TestRPCRes{
		V: 22,
	})
}
