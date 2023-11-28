package strategy

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer/target"
)

type RoundRobinStrategy struct {
	current int
	targets []*target.Target
}

func NewRoundRobinStrategy(targets []*target.Target) *RoundRobinStrategy {
	return &RoundRobinStrategy{
		current: 0,
		targets: targets,
	}
}

func (s *RoundRobinStrategy) NextTarget(*http.Request) *target.Target {
	s.current = (s.current + 1) % len(s.targets)
	target := s.targets[s.current]
	return target
}

func (s *RoundRobinStrategy) SetTargets(targets []*target.Target) {
	s.targets = targets
}
