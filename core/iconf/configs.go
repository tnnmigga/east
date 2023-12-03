package iconf

import (
	"east/core/util"
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v3"
)

var (
	confs map[string]any
	fns   []func()
)

var errConfigNotFound error = errors.New("configs not found")

func LoadFromJSON(b []byte) {
	b = util.Uncomment(b)
	err := json.Unmarshal(b, &confs)
	if err != nil {
		panic(err)
	}
	afterLoad()
}

func LoadFromYAML(b []byte) error {
	err := yaml.Unmarshal(b, &confs)
	if err != nil {
		return err
	}
	afterLoad()
	return nil
}

func RegInitFn(fn func()) {
	fns = append(fns, fn)
}

func afterLoad() {
	initServerConf()
	for _, fn := range fns {
		fn()
	}
}

type vType interface {
	float64 | bool | string
}

func Any[T vType](name string) (v T, ok bool) {
	if v, has := confs[name]; has {
		return v.(T), true // 允许不存在不允许类型错误
	}
	return v, false
}

func Int64(name string, defaultVal ...int64) int64 {
	v, ok := Any[float64](name)
	if ok {
		return int64(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	panic(errConfigNotFound)
}

func Int32(name string, defaultVal ...int32) int32 {
	v, ok := Any[float64](name)
	if ok {
		return int32(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	panic(errConfigNotFound)
}

func UInt64(name string, defaultVal ...uint64) uint64 {
	v, ok := Any[float64](name)
	if ok {
		return uint64(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	panic(errConfigNotFound)
}

func UInt32(name string, defaultVal ...uint32) uint32 {
	v, ok := Any[float64](name)
	if ok {
		return uint32(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	panic(errConfigNotFound)
}

func String(name string, defaultVal ...string) string {
	v, ok := Any[string](name)
	if ok {
		return string(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	panic(errConfigNotFound)
}

func Float64(name string, defaultVal ...float64) float64 {
	v, ok := Any[float64](name)
	if ok {
		return v
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	panic(errConfigNotFound)
}

func Bool(name string, defaultVal ...bool) bool {
	v, ok := Any[bool](name)
	if ok {
		return v
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	panic(errConfigNotFound)
}
