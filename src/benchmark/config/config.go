package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type LoadTestConfig struct {
	Users    int      `yaml:"users"`
	Rampup   int      `yaml:"rampup"`
	Duration int      `yaml:"duration"`
	Targets  []string `yaml:"targets"`
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
