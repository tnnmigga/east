package userimpl

import (
	"east/core/eventbus"
	"east/core/log"
	"east/define"
)

func (s *service) registerEventHandler() {
	s.EventImpl().RegisterHandler(define.EventUserSayHello, s.onEventUserMsg)
}

func (s *service) onEventUserMsg(event *eventbus.Event) {
	log.Infof("user event %v", event)
}
