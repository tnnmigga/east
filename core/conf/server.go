package conf

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

var confs map[string]any

type ServerType string

var (
	serverID   uint64
	serverType string
)

func InitFromJSON(b []byte) error {
	err := json.Unmarshal(b, &confs)
	if err != nil {
		return err
	}
	return nil
}

func InitFromYAML(b []byte) error {
	err := yaml.Unmarshal(b, &confs)
	if err != nil {
		return err
	}
	return nil
}

func ServerID() uint64 {
	return serverID
}
