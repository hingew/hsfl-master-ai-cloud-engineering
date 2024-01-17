package middleware

import (
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
)

type AuthMiddleInterface interface {
	AuthMiddleware(http.ResponseWriter, *http.Request, router.Next)
}
