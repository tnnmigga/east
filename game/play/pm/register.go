package pm

import (
	"reflect"

	"github.com/tnnmigga/core/codec"
	"github.com/tnnmigga/core/idef"
	"github.com/tnnmigga/core/utils"
)

var msgHandler = map[reflect.Type]func(p *Player, msg any){}

func RegMsgHandler[T any](m idef.IModule, handler func(p *Player, msg T)) {
	mType := reflect.TypeOf(utils.New[T]())
	if _, ok := msgHandler[mType]; ok {
		panic("msg handler already exists")
	}
	codec.Register[T]()
	// 注册消息处理函数
	msgHandler[reflect.TypeOf(utils.New[T]())] = func(p *Player, msg any) {
		handler(p, msg.(T))
	}
}
