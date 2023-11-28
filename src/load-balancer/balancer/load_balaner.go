package balancer

import (
	"net/http"
	"sync"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer/strategy"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer/target"
)

type LoadBalancer struct {
	healthLock          *sync.Mutex
	healthcheckInterval time.Duration
	targets             []*target.Target
	healthyTargets      []*target.Target
	client              *http.Client
	strategy            strategy.Strategy
}

func NewLoadBalancer(targets []*target.Target, healthcheckInterval time.Duration, client *http.Client, strategy strategy.Strategy) *LoadBalancer {
	return &LoadBalancer{
		healthLock:          &sync.Mutex{},
		healthcheckInterval: healthcheckInterval,
		targets:             targets,
		healthyTargets:      targets,
		client:              client,
		strategy:            strategy,
	}
}

func (lb *LoadBalancer) StartHealthCheck() {
	go func() {
		for {
			select {
			case <-time.After(lb.healthcheckInterval):
				lb.HealthCheck()
			}
		}
	}()
}

func (lb *LoadBalancer) HealthCheck() {
	lb.healthLock.Lock()
	lb.healthyTargets = make([]*target.Target, 0)
	for _, target := range lb.targets {
		if GetHealth(lb.client, target.Url) {
			target.Health = 0
			lb.healthyTargets = append(lb.healthyTargets, target)
		} else {
			target.Health++
		}
	}
	lb.strategy.SetTargets(lb.healthyTargets)
	lb.healthLock.Unlock()
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.strategy.NextTarget(r).ServeHTTP(w, r)
}
