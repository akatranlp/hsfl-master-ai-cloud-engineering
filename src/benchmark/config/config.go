package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RequestRamp struct {
	Duration  int `yaml:"duration"`
	TargetRPS int `yaml:"targetRPS"`
}

type LoadTestConfig struct {
	Users       int           `yaml:"users"`
	RequestRamp []RequestRamp `yaml:"requestRamp"`
	Targets     []string      `yaml:"targets"`
}

func FromFS(path string) (LoadTestConfig, error) {
	var config LoadTestConfig

	f, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
