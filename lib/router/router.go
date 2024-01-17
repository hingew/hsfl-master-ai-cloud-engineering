package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type route struct {
	method  string
	pattern *regexp.Regexp
	handler http.HandlerFunc
	params  []string
}

type middleware struct {
	pattern *regexp.Regexp
	handler MiddlewareHandler
	params  []string
}

type Router struct {
	routes      []route
	middlewares []middleware
}

type Next func(r *http.Request)
type MiddlewareHandler func(w http.ResponseWriter, r *http.Request, next Next)

func New() *Router {
	return &Router{}
}

// ServeHTTP implements the Http.Handler interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range router.middlewares {

		matches := middleware.pattern.FindStringSubmatch(r.URL.Path)

		if len(matches) > 0 {
			r = createRequestContext(r, middleware.params, matches[1:])

			called := false
			next := func(req *http.Request) {
				called = true
				r = req
			}

			middleware.handler(w, r, next)
			if !called {
				log.Println("Router Middleware: next not called")
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
	log.Println(fmt.Sprintf("Route not found: %s", r.URL.Path))
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
	matcher := regexp.MustCompile(":([a-zA-Z]+)")
	matches := matcher.FindAllStringSubmatch(pattern, -1)

	params := make([]string, len(matches))

	if len(matches) > 0 {
		pattern = matcher.ReplaceAllLiteralString(pattern, "([^/]+)")

		for i, match := range matches {
			params[i] = match[1]
		}
	}

	for _, route := range router.routes {
		if route.method == method && route.pattern.String() == ("^"+pattern+"$") {
			panic(fmt.Sprintf("The route %s with the method %s is already defined!", pattern, method))
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

func (router *Router) DELETE(pattern string, handler http.HandlerFunc) {
	router.addRoute(http.MethodDelete, pattern, handler)
}

func (router *Router) USE(pattern string, handler MiddlewareHandler) *Router {
	matcher := regexp.MustCompile(":([a-zA-Z]+)")
	matches := matcher.FindAllStringSubmatch(pattern, -1)

	// Match all routes
	pattern = pattern + "(.*)"

	params := make([]string, len(matches))

	if len(matches) > 0 {
		pattern = matcher.ReplaceAllLiteralString(pattern, "([^/]+)")

		for i, match := range matches {
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
