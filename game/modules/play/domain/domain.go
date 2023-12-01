package domain

import (
	"east/core/domain"
	"east/game/modules/play/domain/api"
)

type Domain struct {
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

func (d *Domain) MsgCase() api.IMsgCase {
	return d.GetCase(MsgCaseIndex).(api.IMsgCase)
}
