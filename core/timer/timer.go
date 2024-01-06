package timer

import (
	"context"
	"east/core/algorithm"
	"east/core/conf"
	"east/core/idef"
	"east/core/msgbus"
	"east/core/sys"
	"east/core/util"
	"east/core/util/idgen"
	"fmt"
	"time"
)

type timerTrigger struct {
}

type timerCtx struct {
	ID   uint64
	Time time.Duration
	Ctx  any
}

func (t *timerCtx) String() string {
	return fmt.Sprintf("{ID: %d, Time: %d, Body: %s}", t.ID, t.Time, util.String(t.Ctx))
}

func (t *timerCtx) Key() uint64 {
	return t.ID
}

func (t *timerCtx) Value() time.Duration {
	return t.Time
}

type TimerHeap struct {
	algorithm.Heap[uint64, time.Duration, *timerCtx]
	component idef.IComponent
	timer     *time.Timer
}

func NewTimerHeap(com idef.IComponent) *TimerHeap {
	h := &TimerHeap{
		component: com,
	}
	msgbus.RegisterHandler(com, h.onTimerTrigger)
	sys.Go(h.tryNextTrigger)
	return h
}

func (h *TimerHeap) New(delay time.Duration, ctx any) uint64 {
	t := &timerCtx{
		ID:   idgen.NewUUID(),
		Time: util.NowNs() + delay,
		Ctx:  ctx,
	}
	top := h.Top()
	h.Push(t)
	if top != nil && top.Time <= t.Time {
		return t.ID
	}
	if h.timer != nil {
		h.timer.Stop()
		h.timer = nil
	}
	h.tryNextTrigger()
	return t.ID
}

func (h *TimerHeap) Stop(id uint64) bool {
	index := h.Find(id)
	if index == -1 {
		return false
	}
	h.RemoveByIndex(index)
	return true
}

func (h *TimerHeap) tryNextTrigger() {
	top := h.Top()
	if top == nil {
		return
	}
	nowNs := util.NowNs()
	if top.Time <= nowNs {
		h.trigger()
		return
	}
	h.timer = time.NewTimer(time.Duration(top.Time - util.NowNs()))
	sys.Go(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case <-h.timer.C:
			h.trigger()
		}
	})
}

func (h *TimerHeap) trigger() {
	msgbus.Cast(conf.ServerID(), &timerTrigger{}, msgbus.OneOfMods(h.component.Name()))
}
