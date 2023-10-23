package router

import (
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/controller"
)

type TemplateRouter struct {
	router http.Handler
}

func NewTemplateRouter(
	myController controller.Controller,
) *TemplateRouter {
	myRouter := router.New()

	myRouter.GET("/api/templates", myController.GetAllTemplates)
	myRouter.GET("/api/templates/:id", myController.GetTemplate)
	myRouter.POST("/api/templates", myController.CreateTemplate)
	myRouter.PUT("/api/templates/:id", myController.UpdateTemplate)
	myRouter.DELETE("/api/templates/:id", myController.DeleteTemplate)

	return &TemplateRouter{myRouter}
}

func (router *TemplateRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
