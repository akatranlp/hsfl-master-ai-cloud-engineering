package balancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/orchestrator"
)

type RoundRobinBalancer struct {
	idx     int
	healthLock *sync.Mutex
	healthcheckInterval time.Duration
	targets []*orchestrator.Target
	healthyTargets []*orchestrator.Target
	client http.Client
}
// round robin balancer
func NewRoundRobinBalancer(targets []*orchestrator.Target, healthcheckInterval time.Duration, client http.Client) *RoundRobinBalancer {
	for _, target := range targets {
		target.Handler = httputil.NewSingleHostReverseProxy(target.Url)
	}

	log.Println(targets[0])

	return &RoundRobinBalancer{
		idx:     0,
		targets: targets,
		healthcheckInterval: healthcheckInterval,
		healthyTargets: targets,
		healthLock: &sync.Mutex{},
		client: client,
		
	}
}

func (lb *RoundRobinBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.healthLock.Lock()
	lb.idx = (lb.idx + 1) % len(lb.healthyTargets)
	target := lb.healthyTargets[lb.idx]
	log.Println("Send message to:", target.Url)
	lb.healthLock.Unlock()
	target.Handler.ServeHTTP(w, r)
}

func (lb *RoundRobinBalancer) StartHealthCheck() {
	go func() {
		for {
			select {
				case <- time.After(lb.healthcheckInterval):
					lb.healthLock.Lock()
					lb.healthyTargets = make([]*orchestrator.Target, 0)
					for _, target := range lb.targets {
						log.Print("Health-check", target.Url)
						if GetHealth(lb.client, target.Url) {
							log.Println("is healthy")
							target.Health = 0
							lb.healthyTargets = append(lb.healthyTargets, target)
						} else {
							log.Println("is not healthy")
							target.Health++
						}
					}
					lb.healthLock.Unlock()
			}
		}
	}()
}