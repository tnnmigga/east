package domain

import (
	"east/game/play/domain/api"

	"github.com/tnnmigga/core/idef"
	"github.com/tnnmigga/core/infra/domain"
)

type Domain struct {
	idef.IModule
	domain.Root
}

func New(m idef.IModule) *Domain {
	d := &Domain{
		Root:    domain.New(m, MaxCaseIndex),
		IModule: m,
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
