package target

import (
	"net/http"
	"net/url"
)

type Target struct {
	proxy           http.Handler
	Url             *url.URL
	Health          int
	CurrentRequests int
}

func NewTarget(url *url.URL, handler http.Handler) *Target {
	return &Target{
		Url:             url,
		proxy:           handler,
		Health:          0,
		CurrentRequests: 0,
	}
}

func (t *Target) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.CurrentRequests++
	t.proxy.ServeHTTP(w, r)
	t.CurrentRequests--
}
