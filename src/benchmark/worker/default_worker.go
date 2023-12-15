package worker

import (
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/benchmark/client"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/benchmark/config"
)

type DefaultWorker struct {
	ID              int
	client          client.Client
	wg              *sync.WaitGroup
	ramp            config.Ramp
	targets         []*url.URL
	errors          chan error
	terminate       chan bool
	results         chan uint64
	currentDuration int
}

func NewDefaultWorker(
	id int,
	wg *sync.WaitGroup,
	client client.Client,
	ramp config.Ramp,
	targets []*url.URL,
	terminate chan bool,
) *DefaultWorker {
	return &DefaultWorker{
		ID:              id,
		wg:              wg,
		client:          client,
		ramp:            ramp,
		targets:         targets,
		errors:          make(chan error),
		terminate:       terminate,
		results:         make(chan uint64),
		currentDuration: 0,
	}
}

func (w *DefaultWorker) SendRequest(target *url.URL) {
	// fmt.Println("Sending request to", target.Host, target.Path)
	statusCode, err := w.client.Send(target.Host, target.Path)
	if err != nil {
		w.errors <- err
		return
	}
	w.results <- statusCode
}

func (w *DefaultWorker) Work() {
	defer w.wg.Done()

	constantLoadTicker := time.NewTicker(time.Second)
	for {
		select {
		case <-w.terminate:
			println("Terminating", w.ID)
			return
		case statusCode := <-w.results:
			if false {
				println("Worker", w.ID, "got status code", statusCode)
			}
		case err := <-w.errors:
			println("Worker", w.ID, "got error", err.Error())
		case <-constantLoadTicker.C:
			w.currentDuration++
			targetRPS := w.ramp.TargetRPS(w.currentDuration)
			if targetRPS == -1 {
				return
			}
			if targetRPS == 0 {
				continue
			}
			waitDuration := time.Duration(1/targetRPS) * time.Second
			for i := 0; i < targetRPS; i++ {
				beforeRequest := time.Now()
				go w.SendRequest(w.targets[rand.Intn(len(w.targets))])
				dt := time.Since(beforeRequest)
				time.Sleep(waitDuration - dt)
			}
		}
	}
}
