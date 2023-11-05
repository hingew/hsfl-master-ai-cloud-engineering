package my_proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	client ProxyClient
	routes map[string]string
}

func NewReverseProxy(targetHost string) (*ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	client := httputil.NewSingleHostReverseProxy(url)

	proxy := ReverseProxy{client: client, routes: make(map[string]string)}

	originalDirector := client.Director
	client.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req, proxy.routes)
	}

	client.ErrorHandler = errorHandler()

	proxy.client = client

	return &proxy, nil
}

func (reverseProxy *ReverseProxy) Map(sourcePath string, destinationPath string) {
	reverseProxy.routes[sourcePath] = destinationPath
}

func modifyRequest(req *http.Request, routes map[string]string) {
	destinationPath, ok := routes[req.URL.Path]
	if ok {
		req.URL.Path = destinationPath
	}
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "Error while computing request: "+err.Error(), http.StatusInternalServerError)
	}
}

func (reverseProxy *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reverseProxy.client.ServeHTTP(w, r)
}
