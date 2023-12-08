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
		Pool:    domain.NewPool(MaxCaseIndex),
	}
	return d
}

const (
	MsgCaseIndex = iota
	UserCaseIndex
	MaxCaseIndex
)

func (d *Domain) MsgCase() api.IMsg {
	return d.GetCase(MsgCaseIndex).(api.IMsg)
}
