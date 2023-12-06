package eventbus

import (
	"east/core/log"
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
