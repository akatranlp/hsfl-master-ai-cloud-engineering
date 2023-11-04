package httpproxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type RouteMapping struct {
	hostIndex int
	host      *regexp.Regexp
	path      *regexp.Regexp
	hosts     []*url.URL
}

type HTTPProxy struct {
	client   Client
	Mappings []*RouteMapping
}

func New(client Client) *HTTPProxy {
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

		r.Header.Set("X-Forwarded-For", r.RemoteAddr)
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Host = host.Host
		r.URL.Scheme = host.Scheme
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

func (p *HTTPProxy) Add(host string, path string, hosts []string) error {
	wildcardMatcher := regexp.MustCompile("(\\*)")
	wildcardHostMatches := wildcardMatcher.FindAllStringSubmatch(host, -1)
	wildcardPathMatches := wildcardMatcher.FindAllStringSubmatch(host, -1)

	if len(wildcardHostMatches) > 0 {
		host = wildcardMatcher.ReplaceAllLiteralString(host, "([^\\.]*)")
	}

	if len(wildcardPathMatches) > 0 {
		path = wildcardMatcher.ReplaceAllLiteralString(path, "([^/]*)")
	}

	hostPattern := regexp.MustCompile(host)
	pathPattern := regexp.MustCompile(path)

	if len(hosts) < 1 {
		return errors.New("there was no host provided")
	}

	var urls []*url.URL
	for _, hostAddr := range hosts {
		host, err := url.Parse(hostAddr)
		if err != nil {
			return errors.New("invalid origin server URL")
		}
		urls = append(urls, host)
	}

	p.Mappings = append(p.Mappings, &RouteMapping{0, hostPattern, pathPattern, urls})
	return nil
}
