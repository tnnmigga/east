package player

import (
	"east/pb"

	"github.com/tnnmigga/core/codec"
	"github.com/tnnmigga/core/infra/zlog"
	"github.com/tnnmigga/core/msgbus"
)

type timerSayHello struct {
	UserID uint64
	Text   string
}

func (c *useCase) onTimerSayHello(ctx *timerSayHello) {
	zlog.Infof("onTimerSayHello %v", ctx)
	msgbus.Cast(&pb.S2CPackage{
		UserID: ctx.UserID,
		Body: codec.Encode(&pb.SayHelloResp{
			Text: ctx.Text,
		}),
	}, msgbus.UseStream(), msgbus.ServerID(1888))
}
