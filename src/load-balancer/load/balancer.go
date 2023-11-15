package load

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Balancer struct {
	idx     int
	targets []http.Handler
}

func NewBalancer(targetUrls []*url.URL) *Balancer {
	targets := make([]http.Handler, len(targetUrls))

	for i, targetUrl := range targetUrls {
		targets[i] = httputil.NewSingleHostReverseProxy(targetUrl)
	}

	return &Balancer{
		idx:     0,
		targets: targets,
	}
}

func (lb *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := lb.targets[lb.idx]
	target.ServeHTTP(w, r)
	lb.idx = (lb.idx + 1) % len(lb.targets)
}
