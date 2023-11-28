package strategy

import (
	"hash/fnv"
	"net"
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer/target"
)

type IPHashStrategy struct {
	targets []*target.Target
}

func NewIPHashStrategy(targets []*target.Target) *IPHashStrategy {
	return &IPHashStrategy{targets: targets}
}

func (s *IPHashStrategy) NextTarget(r *http.Request) *target.Target {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	hash := fnv.New32a()
	hash.Write([]byte(ip))
	index := int(hash.Sum32()) % len(s.targets)

	target := s.targets[index]
	if target.Health == 0 {
		return target
	}
	return s.targets[0]
}

func (s *IPHashStrategy) SetTargets(targets []*target.Target) {
	s.targets = targets
}
