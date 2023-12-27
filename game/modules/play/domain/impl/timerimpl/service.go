package timerimpl

import (
	"east/core/idef"
	"east/core/timer"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"
	"errors"
)

type service struct {
	*domain.Domain
	*timer.TimerHeap
}

func New(d *domain.Domain) api.ITimer {
	s := &service{
		Domain:    d,
		TimerHeap: timer.NewTimerHeap(d),
	}
	s.After(idef.ServerStateInit, s.afterInit)
	s.Before(idef.ServerStateStop, s.beforeStop)
	return s
}

func (s *service) afterInit() error {
	// 加载数据
	return errors.New("timer error")
	return nil
}

func (s *service) beforeStop() error {
	// 数据落地
	return nil
}
