package worker

import (
	"net/http"
	"sync"
)

type DefaultWorker struct {
	ID        int
	client    http.Client
	wg        *sync.WaitGroup
	terminate chan bool
	jobs      chan string
	results   chan bool
}

func NewDefaultWorker(id int, wg *sync.WaitGroup, client http.Client, jobs chan string, results chan bool, terminate chan bool) *DefaultWorker {
	return &DefaultWorker{
		ID:        id,
		client:    client,
		wg:        wg,
		jobs:      jobs,
		results:   results,
		terminate: terminate,
	}
}

func (w *DefaultWorker) Work() {
	defer w.wg.Done()
	for {
		select {
		case <-w.terminate:
			println("Terminating", w.ID)
			return
		case j := <-w.jobs:
			go func() {
				w.client.Get(j)
				w.results <- true
			}()
		}
	}
}
