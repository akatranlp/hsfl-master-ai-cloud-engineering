package balancer

import (
	"hash/fnv"
	"net"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/orchestrator"
)

type IPHashBalancer struct {
	healthLock          *sync.Mutex
	healthcheckInterval time.Duration
	targets             []*orchestrator.Target
	healthyTargets      []*orchestrator.Target
	client              http.Client
}

func NewIPHashBalancer(targets []*orchestrator.Target, healthcheckInterval time.Duration, client http.Client) *IPHashBalancer {
	for _, target := range targets {
		target.Handler = httputil.NewSingleHostReverseProxy(target.Url)
	}
	return &IPHashBalancer{targets: targets,
		healthcheckInterval: healthcheckInterval,
		healthyTargets:      targets,
		healthLock:          &sync.Mutex{},
		client:              client,
	}
}

func (lb *IPHashBalancer) getServerIndex(ip string) int {
	hash := fnv.New32a()
	hash.Write([]byte(ip))

	return int(hash.Sum32()) % len(lb.targets)
}

func (lb *IPHashBalancer) GetServer(ip string) http.Handler {
	index := lb.getServerIndex(ip)

	if lb.targets[index].Health > 0 {
		return lb.healthyTargets[0].Handler
	}

	return lb.targets[index].Handler
}

func (lb *IPHashBalancer) GetServerForRequest(req *http.Request) http.Handler {
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	return lb.GetServer(ip)
}

func (lb *IPHashBalancer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	lb.healthLock.Lock()
	lb.GetServerForRequest(req).ServeHTTP(rw, req)
	lb.healthLock.Unlock()
}

func (lb *IPHashBalancer) StartHealthCheck() {
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
