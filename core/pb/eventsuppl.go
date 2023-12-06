package pb

import (
	"east/core/log"
	"strconv"
)

func (e Event) Int64Arg(name string) (arg int64) {
	if e.Params != nil {
		if v, ok := e.Params[name]; ok {
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
	if e.Params != nil {
		if v, ok := e.Params[name]; ok {
			return v
		}
	}
	log.Errorf("event get param %s from %s faild", name, e.String())
	return arg
}
