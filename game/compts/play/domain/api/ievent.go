package api

import (
	"east/core/eventbus"
)

type IEvent interface {
	Cast(event *eventbus.Event)
	SyncCast(event *eventbus.Event)
	RegisterSubscriber(sub eventbus.ISubscriber)
	RegisterHandler(topic string, handler func(event *eventbus.Event))
	UnregisterSubscriber(sub eventbus.ISubscriber)
	UnregisterHandler(topic string, handler func(event *eventbus.Event))
}
