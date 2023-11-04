package httpproxy

import (
	"errors"
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

type Proxy interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Append(*RouteMapping)
}

func AddToProxy(p Proxy, host string, path string, hosts []string) error {
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

	p.Append(&RouteMapping{0, hostPattern, pathPattern, urls})

	return nil
}
