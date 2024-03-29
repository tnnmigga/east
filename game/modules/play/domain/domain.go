package domain

import (
	"east/core/basic"
	"east/core/domain"
	"east/game/modules/play/domain/api"
)

type Domain struct {
	*basic.Module
	*domain.Pool
}

func New(m *basic.Module) *Domain {
	d := &Domain{
		Module: m,
		Pool:   domain.NewPool(MaxCaseIndex),
	}
	return d
}

const (
	MsgCaseIndex = iota
	EventCaseIndex
	TimerCaseIndex
	UserCaseIndex
	MaxCaseIndex
)

func (d *Domain) MsgCase() api.IMsg {
	return d.GetCase(MsgCaseIndex).(api.IMsg)
}

func (d *Domain) EventCase() api.IEvent {
	return d.GetCase(EventCaseIndex).(api.IEvent)
}

func (d *Domain) TimerCase() api.ITimer {
	return d.GetCase(TimerCaseIndex).(api.ITimer)
}

func (d *Domain) UserCase() api.IUser {
	return d.GetCase(UserCaseIndex).(api.IUser)
}
