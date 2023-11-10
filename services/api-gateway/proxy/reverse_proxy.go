package my_proxy

import "net/http"

type ReverseProxy interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Map(sourcePath string, destinationPath string)
}
