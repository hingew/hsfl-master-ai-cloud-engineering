package router

import (
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templating-service/templates/controller"
)

type TemplateRouter struct {
	router http.Handler
}

func NewTemplateRouter(
	myController controller.IController,
) *TemplateRouter {
	myRouter := router.New()

	myRouter.GET("/templates", myController.GetAllTemplates)
	myRouter.GET("/templates/:id", myController.GetTemplate)
	myRouter.POST("/templates", myController.CreateTemplate)
	myRouter.PUT("/templates/:id", myController.UpdateTemplate)
	myRouter.DELETE("/templates/:id", myController.DeleteTemplate)

	return &TemplateRouter{myRouter}
}

func (router *TemplateRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
