package router

import (
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
)

func New(
	registerHandler http.HandlerFunc,
	loginHandler http.HandlerFunc,
) *router.Router {
	return router
}
