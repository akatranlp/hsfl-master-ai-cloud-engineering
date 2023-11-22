package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/reverse-proxy/httpproxy"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type RouteMapping struct {
	Host  string   `yaml:"host"`
	Path  string   `yaml:"path"`
	Hosts []string `yaml:"hosts"`
}

type ApplicationConfig struct {
	Mappings []RouteMapping `yaml:"mappings"`
}

type EnvConfig struct {
	ConfigFilePath string `env:"CONFIG_FILE_PATH" envDefault:"config.yaml"`
	ConfigFile     string `env:"CONFIG_FILE"`
	Port           uint16 `env:"PORT" envDefault:"8080"`
}

func LoadConfigFromEnv(content string) (*ApplicationConfig, error) {
	f := strings.NewReader(content)

	var config ApplicationConfig
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	godotenv.Load()

	envConfig := EnvConfig{}
	if err := env.Parse(&envConfig); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	config, err := LoadConfigFromEnv(envConfig.ConfigFile)
	if err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	proxy := httpproxy.NewHTTPProxy(http.DefaultClient)
	httpUtilProxy := httpproxy.NewHTTPUtilProxy(http.DefaultTransport)

	for _, mapping := range config.Mappings {
		if err := httpproxy.AddToProxy(proxy, mapping.Host, mapping.Path, mapping.Hosts); err != nil {
			log.Fatalf("Could not parse application config #{err.Error()}")
		}
		if err := httpproxy.AddToProxy(httpUtilProxy, mapping.Host, mapping.Path, mapping.Hosts); err != nil {
			log.Fatalf("Could not parse application config #{err.Error()}")
		}
	}

	log.Println("Server Started!")
	addr := fmt.Sprintf("0.0.0.0:%d", envConfig.Port)
	if err := http.ListenAndServe(addr, httpUtilProxy); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
