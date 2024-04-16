package conf

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	confs map[string]any = map[string]any{}
	fns   []func()
)

var errConfigNotFound error = errors.New("configs not found")

func LoadFromJSON(b []byte) {
	b = uncomment(b)
	tmp := map[string]any{}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		panic(err)
	}
	transform(tmp, confs, nil)
	afterLoad()
}

func LoadFromYAML(b []byte) error {
	tmp := map[string]any{}
	err := yaml.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	transform(tmp, confs, nil)
	afterLoad()
	return nil
}

func uncomment(b []byte) []byte {
	reg := regexp.MustCompile(`/\*{1,2}[\s\S]*?\*/`)
	b = reg.ReplaceAll(b, []byte("\n"))
	reg = regexp.MustCompile(`\s//[\s\S]*?\n`)
	return reg.ReplaceAll(b, []byte("\n"))
}

func RegInitFn(fn func()) {
	fns = append(fns, fn)
}

func transform(oc map[string]any, nc map[string]any, prefix []string) {
	for k, v := range oc {
		switch v := v.(type) {
		case map[string]any:
			transform(v, nc, append(prefix, k))
		default:
			nc[strings.Join(append(prefix, k), ".")] = oc[k]
		}
	}
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
		return v.(T), true // 允许不存在, 不允许类型错误
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
