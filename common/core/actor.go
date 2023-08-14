package core

import (
	"context"
	"eden/common/log"
	"eden/util"
	"reflect"

	"github.com/gogo/protobuf/proto"
)

type Actor struct {
	context.Context
	ID int64
	MQ chan proto.Message
}

var actors = map[int64]*Actor{}

func NewActor(ctx context.Context, ID int64, mqLen int32) *Actor {
	actor := &Actor{
		Context: ctx,
		ID:      ID,
		MQ:      make(chan proto.Message, mqLen),
	}
	go actor.run()
	return actor
}



func (actor *Actor) run() {
	select {
	case msg := <-actor.MQ:
		msgType := reflect.TypeOf(msg)
		cb, ok := msgCbs[msgType]
		if !ok {
			log.Errorf("msgCb not exist %v", msgType)
		}
		func() {
			defer util.PrintPanicStack()
			cb(msg)
		}()
	case <-actor.Done():
		return
	}
}
