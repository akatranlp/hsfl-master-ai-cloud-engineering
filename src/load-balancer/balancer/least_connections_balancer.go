package balancer

import (
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/orchestrator"
)

type LeastConnectionsBalancer struct {
	currentConnections  map[*orchestrator.Target]int32
	mutex               sync.Mutex
	healthLock          *sync.Mutex
	healthcheckInterval time.Duration
	targets             []*orchestrator.Target
	healthyTargets      []*orchestrator.Target
	client              http.Client
}

func NewLeastConnectionsBalancer(targets []*orchestrator.Target, healthcheckInterval time.Duration, client http.Client) *LeastConnectionsBalancer {
	currentConnections := make(map[*orchestrator.Target]int32)

	for _, target := range targets {
		target.Handler = httputil.NewSingleHostReverseProxy(target.Url)
		currentConnections[target] = 0
	}

	return &LeastConnectionsBalancer{targets: targets,
		healthcheckInterval: healthcheckInterval,
		healthyTargets:      targets,
		healthLock:          &sync.Mutex{},
		currentConnections:  currentConnections,
		client:              client,
	}
}

// Get the next server with least connections using mutex without atomic
func (lb *LeastConnectionsBalancer) NextServer() *orchestrator.Target {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	min := int32(0)

	var least *orchestrator.Target
	for _, target := range lb.healthyTargets {
		curr := lb.currentConnections[target]
		if min == 0 || curr < min {
			min = curr
			least = target
		}
	}

	lb.currentConnections[least]++
	return least

}

// Reduce the current connection on server when connection has ended
func (lb *LeastConnectionsBalancer) ReduceConnection(target *orchestrator.Target) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	lb.currentConnections[target]--
}

func (lb *LeastConnectionsBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.healthLock.Lock()
	target := lb.NextServer()
	target.Handler.ServeHTTP(w, r)
	lb.ReduceConnection(target)
	lb.healthLock.Unlock()
}

func (lb *LeastConnectionsBalancer) StartHealthCheck() {
	go func() {
		for {
			select {
			case <-time.After(lb.healthcheckInterval):
				lb.healthLock.Lock()
				lb.healthyTargets = make([]*orchestrator.Target, 0)
				for _, target := range lb.targets {

					if GetHealth(lb.client, target.Url) {
						target.Health = 0
						lb.healthyTargets = append(lb.healthyTargets, target)
					} else {
						target.Health++
					}
				}
				lb.healthLock.Unlock()
			}
		}
	}()
}
