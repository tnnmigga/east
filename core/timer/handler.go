package timer

import (
	"east/core/conf"
	"east/core/msgbus"
	"east/core/util"
)

func (h *TimerHeap) onTimerTrigger(msg *timerTrigger) {
	defer h.tryNextTrigger()
	nowNs := util.NowNs()
	for top := h.Top(); top != nil && top.Time <= nowNs; top = h.Top() {
		h.Pop()
		msgbus.Cast(conf.ServerID(), top.Ctx)
	}
}
