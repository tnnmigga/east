package domain

import (
	"east/core/domain"
	"east/core/idef"
	"east/game/modules/play/domain/api"
)

type Domain struct {
	idef.IModule
	*domain.Pool
}

func New(m idef.IModule) *Domain {
	d := &Domain{
		IModule: m,
		Pool:    domain.NewPool(MaxImplIndex),
	}
	return d
}

const (
	MsgImplIndex = iota
	EventImplIndex
	UserImplIndex
	MaxImplIndex
)

func (d *Domain) MsgCase() api.IMsg {
	return d.GetImpl(MsgImplIndex).(api.IMsg)
}

func (d *Domain) EventImpl() api.IEvent {
	return d.GetImpl(EventImplIndex).(api.IEvent)
}

func (d *Domain) UserImpl() api.IUser {
	return d.GetImpl(UserImplIndex).(api.IUser)
}
