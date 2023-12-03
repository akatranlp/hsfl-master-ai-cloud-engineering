package config

import (
	"encoding/json"
	"os"
)

type LoadTestConfig struct {
	Users    int      `json:"users"`
	Rampup   int      `json:"rampup"`
	Duration int      `json:"duration"`
	Targets  []string `json:"targets"`
}

func FromFS(path string) (LoadTestConfig, error) {
	var config LoadTestConfig

	f, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&config)
	return config, err
}
