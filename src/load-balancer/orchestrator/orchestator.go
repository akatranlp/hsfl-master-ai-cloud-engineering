package orchestrator

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer/target"

type Orchestrator interface {
	StartContainers(image string, replicas int) []string
	StopContainers(containers []string)
	StopAllContainers()
	GetContainerEndpoint(containers []string, networkName string) []*target.Target
	Close()
}
