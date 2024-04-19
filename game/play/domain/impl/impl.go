package impl

import (
	"east/game/play/domain"
	"east/game/play/domain/impl/event"
	"east/game/play/domain/impl/msg"
	"east/game/play/domain/impl/player"
	"east/game/play/domain/impl/timer"
)

func Init(d *domain.Domain) {
	d.PutCase(domain.MsgCaseIndex, msg.New(d))
	d.PutCase(domain.EventCaseIndex, event.New(d))
	d.PutCase(domain.TimerCaseIndex, timer.New(d))
	d.PutCase(domain.UserCaseIndex, player.New(d))
}
