package timerimpl

import (
	"east/core/timer"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"
)

type service struct {
	*domain.Domain
	*timer.TimerHeap
}

func New(d *domain.Domain) api.ITimer {
	return &service{
		Domain:    d,
		TimerHeap: timer.NewTimerHeap(d),
	}
}

func (s *service) Init() {
	// 加载定时器数据
}

func (s *service) Destroy() {
	// 定时器数据落地
}
