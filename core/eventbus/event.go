package eventbus

import (
	"east/core/idef"
	"east/core/log"
	"east/core/msgbus"
	"east/core/util"
	"strconv"
)

type ISubscriber interface {
	Name() string
	Topics() []string
	Handler(event *Event)
}

type EventBus struct {
	subs map[string][]ISubscriber
}

func (bus *EventBus) RegisterSubscriber(sub ISubscriber) {
	if bus.find(sub) {
		log.Errorf("%s has registered", sub.Name())
		return
	}
	for _, topic := range sub.Topics() {
		bus.subs[topic] = append(bus.subs[sub.Name()], sub)
	}
}

func (bus *EventBus) RegisterHandler(topic string, handler func(event *Event)) {
	h := &eventHandler{
		name:    util.FuncName(handler),
		topic:   topic,
		handler: handler,
	}
	bus.RegisterSubscriber(h)
}

func (bus *EventBus) UnregisterSubscriber(sub ISubscriber) {
	bus.removeSubscriber(sub.Name())
}

func (bus *EventBus) UnregisterHandler(topic string, handler func(event *Event)) {
	bus.removeSubscriber(util.FuncName(handler))
}

func (bus *EventBus) removeSubscriber(name string) {
	for topic, subs := range bus.subs {
		bus.subs[topic] = util.Filter(subs, func(sub ISubscriber) bool {
			return sub.Name() != name
		})
	}
}

func (bus *EventBus) find(sub ISubscriber) bool {
	subName := sub.Name()
	for _, sub := range bus.subs {
		for _, v := range sub {
			if v.Name() == subName {
				return true
			}
		}
	}
	return false
}

func (bus *EventBus) dispatch(event *Event) {
	subs := bus.subs[event.Topic]
	for _, sub := range subs {
		util.ExecAndRecover(func() {
			sub.Handler(event)
		})
	}
}

func New(m idef.IModule) *EventBus {
	bus := &EventBus{
		subs: map[string][]ISubscriber{},
	}
	msgbus.RegisterHandler(m, bus.dispatch)
	return bus
}

func (bus *EventBus) Cast(event *Event) {
	msgbus.CastLocal(event)
}

func (bus *EventBus) SyncCast(event *Event) {
	bus.dispatch(event)
}

type eventHandler struct {
	name    string
	topic   string
	handler func(*Event)
}

func (h *eventHandler) Name() string {
	return h.name
}

func (h *eventHandler) Topics() []string {
	return []string{h.topic}
}

func (h *eventHandler) Handler(event *Event) {
	h.handler(event)
}

func (e Event) Int64Arg(name string) (arg int64) {
	if e.Args != nil {
		if v, ok := e.Args[name]; ok {
			if n, err := strconv.Atoi(v); err != nil {
				return int64(n)
			}
		}
	}
	log.Errorf("event get param %s from %s faild", name, e.String())
	return arg
}

func (e Event) Int32Arg(name string) (arg int32) {
	return int32(e.Int64Arg(name))
}

func (e Event) StringArg(name string) (arg string) {
	if e.Args != nil {
		if v, ok := e.Args[name]; ok {
			return v
		}
	}
	log.Errorf("event get param %s from %s faild", name, e.String())
	return arg
}
