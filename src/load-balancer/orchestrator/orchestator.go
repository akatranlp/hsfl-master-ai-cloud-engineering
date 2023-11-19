package orchestrator

import (
	"net/http"
	"net/url"
)

type Target struct {
	ContainerId string
	Handler     http.Handler
	Url         *url.URL
	Health      int
}

type Orchestrator interface {
	StartContainers(image string, replicas int) []string
	StopContainers(containers []string)
	StopAllContainers()
	GetContainerEndpoint(containers []string, networkName string) []*Target
	Close()
}
