package impl

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/impl/event"
	"east/game/modules/play/domain/impl/msg"
	"east/game/modules/play/domain/impl/timer"
	"east/game/modules/play/domain/impl/user"
)

func Init(d *domain.Domain) {
	d.PutImpl(domain.MsgCaseIndex, msg.New(d))
	d.PutImpl(domain.EventCaseIndex, event.New(d))
	d.PutImpl(domain.TimerCaseIndex, timer.New(d))
	d.PutImpl(domain.UserCaseIndex, user.New(d))
}
