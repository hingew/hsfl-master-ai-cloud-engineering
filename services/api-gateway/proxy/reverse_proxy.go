package my_proxy

import "net/http"

type ReverseProxy interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Map(sourcePath string, destinationPath string, coalescing bool)
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}
