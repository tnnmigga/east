package domain

import (
	"east/core/domain"
	"east/core/mod"
	"east/game/modules/play/domain/api"
)

type Domain struct {
	*mod.Module
	*domain.Pool
}

func New(m *mod.Module) *Domain {
	d := &Domain{
		Module: m,
		Pool:   domain.NewPool(MaxImplIndex),
	}
	return d
}

const (
	MsgImplIndex = iota
	EventImplIndex
	TimerImplIndex
	UserImplIndex
	MaxImplIndex
)

func (d *Domain) MsgCase() api.IMsg {
	return d.GetImpl(MsgImplIndex).(api.IMsg)
}

func (d *Domain) EventImpl() api.IEvent {
	return d.GetImpl(EventImplIndex).(api.IEvent)
}

func (d *Domain) TimerImpl() api.ITimer {
	return d.GetImpl(TimerImplIndex).(api.ITimer)
}

func (d *Domain) UserImpl() api.IUser {
	return d.GetImpl(UserImplIndex).(api.IUser)
}
