package userimpl

import (
	"east/core/msgbus"
	"east/define"
)

func (s *service) register() {
	s.registerHandler()
	s.registerEvent()
	s.registerTime()
}

func (s *service) registerHandler() {
	msgbus.RegisterHandler(s, s.onSayHelloReq)
}

func (s *service) registerEvent() {
	s.EventImpl().RegisterHandler(define.EventUserSayHello, s.onEventUserMsg)
}

func (s *service) registerTime() {
	msgbus.RegisterHandler(s, s.onTimerSayHello)
}
