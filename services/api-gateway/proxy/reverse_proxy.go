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

	originalDirector := client.Director
	client.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	client.ModifyResponse = modifyResponse()
	client.ErrorHandler = errorHandler()
	return &ReverseProxy{client: client}, nil
}

func (reverseProxy *ReverseProxy) Map(sourcePath string, destinationPath string) {
	reverseProxy.routes[sourcePath] = destinationPath
}

func modifyRequest(req *http.Request) {

}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return nil
}

func modifyResponse() func(*http.Response) error {
	return nil
}

func (reverseProxy *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reverseProxy.client.ServeHTTP(w, r)
}
