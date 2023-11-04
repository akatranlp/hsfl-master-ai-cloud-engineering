package httpproxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type HTTPProxy struct {
	client   Client
	Mappings []*RouteMapping
}

func NewHTTPProxy(client Client) *HTTPProxy {
	return &HTTPProxy{client: client}
}

func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

		r.Header.Set("X-Forwarded-For", strings.Split(r.RemoteAddr, ":")[0])
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.URL.Host = host.Host
		r.URL.Scheme = host.Scheme
		r.URL.Path = host.Path + r.URL.Path
		r.RequestURI = ""

		originServerResponse, err := p.client.Do(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			mapping.hostIndex = (mapping.hostIndex + 1) % len(mapping.hosts)
			log.Println(err.Error())
			return
		}

		for headerKey, headerValue := range originServerResponse.Header {
			for _, value := range headerValue {
				w.Header().Add(headerKey, value)
			}
		}

		w.WriteHeader(originServerResponse.StatusCode)

		io.Copy(w, originServerResponse.Body)

		// Round-robin
		mapping.hostIndex = (mapping.hostIndex + 1) % len(mapping.hosts)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
}

func (p *HTTPProxy) Append(mapping *RouteMapping) {
	p.Mappings = append(p.Mappings, mapping)
}
