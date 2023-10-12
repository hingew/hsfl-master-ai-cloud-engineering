package router

import "net/http"

type Router struct {
	mux *http.ServeMux
}

func New(
	registerHandler http.Handler,
	loginHandler http.Handler,
) *Router {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/auth/register", registerHandler)
	mux.Handle("/api/v1/auth/login", loginHandler)

	return &Router{mux}
}

func (handler *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.mux.ServeHTTP(w, r)
}
