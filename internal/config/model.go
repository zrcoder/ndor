package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Examples []Example
}

type Example struct {
	Name string
	Code string
}

var (
	//go:embed default.yaml
	defaultData []byte
	Default     *Config
)

func init() {
	err := yaml.Unmarshal(defaultData, &Default)
	if err != nil {
		panic(err)
	}
}
