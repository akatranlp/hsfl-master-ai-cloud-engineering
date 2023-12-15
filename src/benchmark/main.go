package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/benchmark/client"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/benchmark/config"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/benchmark/worker"
)

func main() {
	confPath := os.Getenv("CONFIG_PATH")
	if confPath == "" {
		log.Fatal("There is no configPath provided")
	}
	conf, err := config.FromFS(confPath)
	if err != nil {
		log.Fatal("Conf couldn't pe parsed!", err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(conf.Users)

	startTime := time.Now()

	terminate := make(chan bool)

	workers := make([]worker.Worker, conf.Users)

	ramp := config.NewLinearRamp(conf.RequestRamp)
	targets := make([]*url.URL, len(conf.Targets))
	for i, t := range conf.Targets {
		targets[i], err = url.Parse(t)
		if err != nil {
			log.Fatal("Target couldn't be parsed!", err.Error())
		}
	}

	waitTime := time.Duration(1/conf.Users) * time.Second
	for i := 0; i < conf.Users; i++ {
		currentTime := time.Now()
		workers[i] = worker.NewDefaultWorker(i+1, &wg, client.NewTcpClient(), ramp, targets, terminate)
		go workers[i].Work()
		dt := time.Since(currentTime)
		time.Sleep(waitTime - dt)
	}

	wg.Wait()

	fmt.Printf("Load test completed in %s\n", time.Since(startTime))
}
