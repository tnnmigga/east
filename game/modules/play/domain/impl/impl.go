package impl

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/impl/event"
	"east/game/modules/play/domain/impl/msg"
	"east/game/modules/play/domain/impl/player"
	"east/game/modules/play/domain/impl/timer"
)

func Init(d *domain.Domain) {
	d.PutCase(domain.MsgCaseIndex, msg.New(d))
	d.PutCase(domain.EventCaseIndex, event.New(d))
	d.PutCase(domain.TimerCaseIndex, timer.New(d))
	d.PutCase(domain.UserCaseIndex, player.New(d))
}
