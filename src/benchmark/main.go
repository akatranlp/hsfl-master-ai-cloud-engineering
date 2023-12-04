package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

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
	fmt.Printf("Starting load test with %d users, rampup %ds, duration %ds\n", conf.Users, conf.Rampup, conf.Duration)

	var wg sync.WaitGroup
	wg.Add(conf.Users)

	startTime := time.Now()

	terminate := make(chan bool)
	results := make(chan bool, conf.Users*conf.Duration)
	jobs := make(chan string, conf.Users*conf.Duration)

	rampupDuration := time.Duration(conf.Rampup) * time.Second
	rampupRate := float64(conf.Users) / float64(rampupDuration.Seconds())

	workers := make([]worker.Worker, conf.Users)

	for i := 0; i < conf.Users; i++ {
		workers[i] = worker.NewDefaultWorker(i+1, &wg, http.Client{
			Timeout: time.Duration(100*time.Millisecond) * time.Second,
		}, jobs, results, terminate)
		go workers[i].Work()
	}

	reqCounter := 0

	currentRampUpUsers := rampupRate

	constantLoadTicker := time.NewTicker(time.Second)

	for elapsed := 0; elapsed < conf.Rampup; elapsed++ {
		<-constantLoadTicker.C
		fmt.Println(currentRampUpUsers)
		for j := 0; j < int(currentRampUpUsers); j++ {
			jobs <- conf.Targets[rand.Intn(len(conf.Targets))]
			reqCounter++
		}
		currentRampUpUsers += rampupRate
	}

	for elapsed := conf.Rampup; elapsed < conf.Duration; elapsed++ {
		<-constantLoadTicker.C
		fmt.Println(conf.Users)
		for j := 0; j < conf.Users; j++ {
			jobs <- conf.Targets[rand.Intn(len(conf.Targets))]
			reqCounter++
		}
	}

	<-constantLoadTicker.C
	wait := os.Getenv("WAIT_FOR_RESULTS")
	if wait != "" {
		for i := 0; i < reqCounter; i++ {
			<-results
		}
	}

	fmt.Println(reqCounter)

	close(terminate)

	wg.Wait()

	fmt.Printf("Load test completed in %s\n", time.Since(startTime))
}
