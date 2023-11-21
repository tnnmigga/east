package domain

import (
	"eden/core/domain"
	"eden/game/modules/play/domain/api"
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
	MsgCaseIndex domain.CaseIndex = iota
	MaxCaseIndex
)

func (d *Domain) MsgCase() api.IMsgCase {
	return d.GetCase(MsgCaseIndex).(api.IMsgCase)
}
