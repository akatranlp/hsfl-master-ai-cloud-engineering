package strategy

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer/target"
)

type LeastConnectionStrategy struct {
	targets []*target.Target
}

func NewLeastConnectionsStrategy(targets []*target.Target) *LeastConnectionStrategy {
	return &LeastConnectionStrategy{targets: targets}
}

func (s *LeastConnectionStrategy) NextTarget(r *http.Request) *target.Target {
	min := s.targets[0]
	for _, target := range s.targets {
		if target.CurrentRequests < min.CurrentRequests && target.Health == 0 {
			min = target
		}
	}
	return min
}

func (s *LeastConnectionStrategy) SetTargets(targets []*target.Target) {
	s.targets = targets
}
