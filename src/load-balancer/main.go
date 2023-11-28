package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer/strategy"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/orchestrator"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type ApplicationConfig struct {
	Port          int   `env:"PORT" envDefault:"8080"`
	HealthTimeout int64 `env:"HEALTH_TIMEOUT" envDefault:"200"`
}

func main() {
	godotenv.Load()

	envConfig := ApplicationConfig{}
	if err := env.Parse(&envConfig); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	image := flag.String("image", "akatranlp/web-service:latest", "")
	replicas := flag.Int("replicas", 1, "")
	network := flag.String("network", "bridge", "")
	strategyStr := flag.String("strategy", "round-robin", "")
	flag.Parse()

	orc := orchestrator.NewDefaultOrchestrator()
	defer orc.Close()

	containers := orc.StartContainers(*image, *replicas, *network)
	defer orc.StopContainers(containers)
	endpoints := orc.GetContainerEndpoints(containers, *network, envConfig.Port)

	client := &http.Client{
		Timeout: time.Duration(envConfig.HealthTimeout) * time.Millisecond,
	}

	var strategyImpl strategy.Strategy
	switch *strategyStr {
	case "round-robin":
		strategyImpl = strategy.NewRoundRobinStrategy(endpoints)
	case "ip-hash":
		strategyImpl = strategy.NewIPHashStrategy(endpoints)
	case "least-connections":
		strategyImpl = strategy.NewLeastConnectionsStrategy(endpoints)
	default:
		strategyImpl = strategy.NewRoundRobinStrategy(endpoints)
	}

	lb := balancer.NewLoadBalancer(endpoints, 10*time.Second, client, strategyImpl)

	lb.StartHealthCheck()

	addr := fmt.Sprintf(":%d", envConfig.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: lb,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		server.Shutdown(context.Background())
	}()

	server.ListenAndServe()
}
