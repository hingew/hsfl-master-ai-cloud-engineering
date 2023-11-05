package my_proxy

import "net/http"

type ProxyClient interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
