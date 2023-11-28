package strategy

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer/target"
)

type Strategy interface {
	NextTarget(*http.Request) *target.Target
	SetTargets([]*target.Target)
}
