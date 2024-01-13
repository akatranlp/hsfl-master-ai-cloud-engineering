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
	ticksPerSecond  int
	ramp            config.Ramp
	targets         []*url.URL
	terminate       chan bool
	currentDuration int
	results         chan uint64
	errors          chan error
}

func NewDefaultWorker(
	id int,
	wg *sync.WaitGroup,
	client client.Client,
	ramp config.Ramp,
	targets []*url.URL,
	terminate chan bool,
	ticksPerSecond int,
) *DefaultWorker {
	return &DefaultWorker{
		ID:              id,
		wg:              wg,
		client:          client,
		ticksPerSecond:  ticksPerSecond,
		ramp:            ramp,
		targets:         targets,
		errors:          make(chan error),
		terminate:       terminate,
		results:         make(chan uint64),
		currentDuration: 0,
	}
}

func (w *DefaultWorker) SendRequest(target *url.URL) {
	statusCode, err := w.client.Send(target.Host, target.Path)
	w.results <- statusCode
	if err != nil {
		w.errors <- err
		return
	}
}

func (w *DefaultWorker) Work() {
	defer w.wg.Done()

	constantLoadTicker := time.NewTicker(time.Second / time.Duration(w.ticksPerSecond))
	for {
		select {
		case <-w.terminate:
			println("Terminating", w.ID)
			return
		case <-constantLoadTicker.C:
			w.currentDuration++
			targetRPS := w.ramp.TargetRPS(w.currentDuration)
			if targetRPS == -1 {
				return
			}
			if targetRPS == 0 {
				continue
			}
			waitDuration := time.Duration(1/targetRPS) * (time.Second / time.Duration(w.ticksPerSecond))
			for i := 0; i < targetRPS; i++ {
				beforeRequest := time.Now()
				go w.SendRequest(w.targets[rand.Intn(len(w.targets))])
				dt := time.Since(beforeRequest)
				time.Sleep(waitDuration - dt)
			}
		}
	}
}
