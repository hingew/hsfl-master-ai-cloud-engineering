package coalecing

import "net/http"

type CoalecingServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
