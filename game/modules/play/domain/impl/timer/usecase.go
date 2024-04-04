package timer

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"

	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/timer"
)

type useCase struct {
	*domain.Domain
	*timer.TimerHeap
}

func New(d *domain.Domain) api.ITimer {
	s := &useCase{
		Domain:    d,
		TimerHeap: timer.NewTimerHeap(d),
	}
	s.After(idef.ServerStateInit, s.afterInit)
	s.Before(idef.ServerStateStop, s.beforeStop)
	return s
}

func (c *useCase) afterInit() error {
	// 加载数据
	return nil
}

func (c *useCase) beforeStop() error {
	// 数据落地
	return nil
}
