package router

import (
	"context"
	"net/http"
	"regexp"
)

type route struct {
	method  string
	pattern *regexp.Regexp
	handler http.HandlerFunc
	params  []string
}

type Next func(r *http.Request)
type MiddleWareFunc func(w http.ResponseWriter, r *http.Request, next Next)

type middleware struct {
	pattern *regexp.Regexp
	handler MiddleWareFunc
	params  []string
}

type Router struct {
	routes      []route
	middlewares []middleware
}

func New() *Router {
	return &Router{}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range router.middlewares {

		matches := middleware.pattern.FindStringSubmatch(r.URL.Path)

		if len(matches) > 0 {
			r = createRequestContext(r, middleware.params, matches[1:])

			nextWasCalled := false
			next := func(req *http.Request) {
				nextWasCalled = true
				r = req
			}

			middleware.handler(w, r, next)
			if !nextWasCalled {
				return
			}
		}
	}

	for _, route := range router.routes {
		if r.Method != route.method {
			continue
		}

		matches := route.pattern.FindStringSubmatch(r.URL.Path)

		if len(matches) > 0 {
			r = createRequestContext(r, route.params, matches[1:])
			route.handler(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func createRequestContext(r *http.Request, paramKeys []string, paramValues []string) *http.Request {
	if len(paramKeys) == 0 {
		return r
	}

	ctx := r.Context()
	for i := 0; i < len(paramKeys); i++ {
		ctx = context.WithValue(ctx, paramKeys[i], paramValues[i])
	}

	return r.WithContext(ctx)
}

func (router *Router) addRoute(method string, pattern string, handler http.HandlerFunc) {
	paramMatcher := regexp.MustCompile(":([a-zA-Z]+)")
	paramMatches := paramMatcher.FindAllStringSubmatch(pattern, -1)

	params := make([]string, len(paramMatches))

	if len(paramMatches) > 0 {
		pattern = paramMatcher.ReplaceAllLiteralString(pattern, "([^/]+)")

		for i, match := range paramMatches {
			params[i] = match[1]
		}
	}

	router.routes = append(router.routes, route{
		method:  method,
		pattern: regexp.MustCompile("^" + pattern + "$"),
		handler: handler,
		params:  params,
	})
}

func (router *Router) GET(pattern string, handler http.HandlerFunc) {
	router.addRoute(http.MethodGet, pattern, handler)
}

func (router *Router) POST(pattern string, handler http.HandlerFunc) {
	router.addRoute(http.MethodPost, pattern, handler)
}

func (router *Router) PUT(pattern string, handler http.HandlerFunc) {
	router.addRoute(http.MethodPut, pattern, handler)
}

func (router *Router) PATCH(pattern string, handler http.HandlerFunc) {
	router.addRoute(http.MethodPatch, pattern, handler)
}

func (router *Router) DELETE(pattern string, handler http.HandlerFunc) {
	router.addRoute(http.MethodDelete, pattern, handler)
}

func (router *Router) USE(pattern string, handler MiddleWareFunc) *Router {
	var params []string

	pattern = pattern + "(.*)"

	paramMatcher := regexp.MustCompile(":([a-zA-Z]+)")
	paramMatches := paramMatcher.FindAllStringSubmatch(pattern, -1)

	params = make([]string, len(paramMatches))

	if len(paramMatches) > 0 {
		pattern = paramMatcher.ReplaceAllLiteralString(pattern, "([^/]+)")

		for i, match := range paramMatches {
			params[i] = match[1]
		}
	}

	router.middlewares = append(router.middlewares, middleware{
		pattern: regexp.MustCompile("^" + pattern + "$"),
		handler: handler,
		params:  params,
	})

	return router
}
