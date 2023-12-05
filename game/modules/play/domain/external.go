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

func NewDomain() *Domain {
	d := &Domain{
		Pool: domain.NewPool(MaxCaseIndex),
	}
	return d
}

const (
	MsgCaseIndex = iota
	MaxCaseIndex
)

func (d *Domain) MsgCase() api.IMsg {
	return d.GetCase(MsgCaseIndex).(api.IMsg)
}
