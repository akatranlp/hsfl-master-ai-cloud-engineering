package httpproxy

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type HTTPUtilProxy struct {
	roundTripper http.RoundTripper
	Mappings     []*RouteMapping
}

func NewHTTPUtilProxy(roundTripper http.RoundTripper) *HTTPUtilProxy {
	return &HTTPUtilProxy{roundTripper: roundTripper}
}

func (p *HTTPUtilProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, mapping := range p.Mappings {
		hostMatches := mapping.host.FindAllStringSubmatch(r.Host, -1)

		if len(hostMatches) < 1 {
			continue
		}

		pathMatches := mapping.path.FindAllStringSubmatch(r.URL.Path, -1)

		if len(pathMatches) < 1 {
			continue
		}

		host := mapping.hosts[mapping.hostIndex]
		log.Printf("Got a connection from %s to %s: Redirect to %s\n", r.RemoteAddr, r.Host, host.Host)

		reverseProxy := httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetURL(host)
				r.Out.Host = r.In.Host
				r.SetXForwarded()
			},
			Transport: p.roundTripper,
		}

		mapping.hostIndex = (mapping.hostIndex + 1) % len(mapping.hosts)
		reverseProxy.ServeHTTP(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
}

func (p *HTTPUtilProxy) Append(mapping *RouteMapping) {
	p.Mappings = append(p.Mappings, mapping)
}
