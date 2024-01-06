package domain

import (
	"east/core/com"
	"east/core/domain"
	"east/game/modules/play/domain/api"
)

type Domain struct {
	*com.Component
	*domain.Pool
}

func New(m *com.Component) *Domain {
	d := &Domain{
		Component: m,
		Pool:      domain.NewPool(MaxImplIndex),
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
