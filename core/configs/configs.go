package configs

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

var (
	confs map[string]any
	fns   []func()
)

func LoadFromJSON(b []byte) error {
	err := json.Unmarshal(b, &confs)
	if err != nil {
		return err
	}
	afterLoad()
	return nil
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

func Int64(name string) int64 {
	return int64(Float64(name))
}

func Int32(name string) int32 {
	return int32(Float64(name))

}

func UInt64(name string) uint64 {
	return uint64(Float64(name))
}

func UInt32(name string) uint32 {
	return uint32(Float64(name))

}

func String(name string) string {
	return confs[name].(string)
}

func Float64(name string) float64 {
	return confs[name].(float64)
}
