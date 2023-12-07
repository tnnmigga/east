package eventbus

import (
	"east/core/idef"
	"east/core/log"
	"east/core/message"
	"east/core/pb"
	"east/core/util"
)

type Subscriber interface {
	Name() string
	Topics() []string
	Handler(event *pb.Event)
}

type EventBus struct {
	subs map[string][]Subscriber
}

func (bus *EventBus) RegisterSubscriber(sub Subscriber) {
	if bus.find(sub) {
		log.Errorf("%s has registered", sub.Name())
		return
	}
	for _, topic := range sub.Topics() {
		bus.subs[topic] = append(bus.subs[sub.Name()], sub)
	}
}

func (bus *EventBus) RegisterHandler(topic string, handler func(event *pb.Event)) {
	h := &eventHandler{
		name:    util.FuncName(handler),
		topic:   topic,
		handler: handler,
	}
	bus.RegisterSubscriber(h)
}

func (bus *EventBus) UnregisterSubscriber(sub Subscriber) {
	bus.removeSubscriber(sub.Name())
}

func (bus *EventBus) UnregisterHandler(topic string, handler func(event *pb.Event)) {
	bus.removeSubscriber(util.FuncName(handler))
}

func (bus *EventBus) removeSubscriber(name string) {
	for topic, subs := range bus.subs {
		bus.subs[topic] = util.Filter(subs, func(sub Subscriber) bool {
			return sub.Name() != name
		})
	}
}

func (bus *EventBus) find(sub Subscriber) bool {
	subName := sub.Name()
	for _, m := range bus.subs {
		for _, v := range m {
			if v.Name() == subName {
				return true
			}
		}
	}
	return false
}

func (bus *EventBus) dispatch(event *pb.Event) {
	subs := bus.subs[event.Topic]
	for _, sub := range subs {
		util.ExecAndRecover(func() {
			sub.Handler(event)
		})
	}
}

func (bus *EventBus) Bind(m idef.IModule) {
	message.RegisterHandler(m, bus.dispatch)
}

type eventHandler struct {
	name    string
	topic   string
	handler func(*pb.Event)
}

func (h *eventHandler) Name() string {
	return h.name
}

func (h *eventHandler) Topics() []string {
	return []string{h.topic}
}

func (h *eventHandler) Handler(event *pb.Event) {
	h.handler(event)
}
