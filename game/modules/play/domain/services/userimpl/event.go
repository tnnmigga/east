package userimpl

import (
	"east/core/eventbus"
	"east/core/log"
)

func (s *service) onEventUserMsg(event *eventbus.Event) {
	log.Infof("user event %v", event)
}
