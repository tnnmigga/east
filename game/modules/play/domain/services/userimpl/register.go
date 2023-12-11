package userimpl

import (
	"east/core/msgbus"
	"east/define"
)

func (s *service) register() {
	s.registerHandler()
	s.registerEvent()
}

func (s *service) registerHandler() {
	msgbus.RegisterHandler(s, s.onSayHelloReq)
	// timer
	msgbus.RegisterHandler(s, s.onTimerSayHello)
}

func (s *service) registerEvent() {
	s.EventImpl().RegisterHandler(define.EventUserSayHello, s.onEventUserMsg)
}