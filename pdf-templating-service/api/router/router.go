package router

import (
	"context"
	"net/http"
	"regexp"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/pdf-templating-service/pdftemplates"
)

type route struct {
	method  string
	pattern *regexp.Regexp
	handler http.HandlerFunc
	params  []string
}

type Router struct {
	routes []route
}

func New(templateController pdftemplates.Controller) *Router {
	router := Router{}

	router.GET("api/v1/pdftemplates", templateController.GetTemplates)
	router.POST("api/v1/pdftemplates", templateController.PostTemplates)
	router.GET("api/v1/templates/:pdftemplateID", templateController.GetTemplate)
	router.PUT("api/v1/templates/:pdftemplateID", templateController.PutTemplates)
	router.DELETE("api/v1/templates/:pdftemplateID", templateController.DeleteTemplates)

	return &router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (router *Router) DELETE(pattern string, handler http.HandlerFunc) {
	router.addRoute(http.MethodDelete, pattern, handler)
}
