package router

import (
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/controller"
)

func NewTemplateRouter(
	myController controller.Controller,
) *router.Router {
	templateRouter := router.New()
	templateRouter.GET("/api/health/templates", health.Check)

	templateRouter.GET("/api/templates", myController.GetAllTemplates)
	templateRouter.GET("/api/templates/:id", myController.GetTemplate)
	templateRouter.GET("/api/templates/:id/coalecing", myController.GetTemplateWithCoalecing)
	templateRouter.POST("/api/templates", myController.CreateTemplate)
	templateRouter.PUT("/api/templates/:id", myController.UpdateTemplate)
	templateRouter.DELETE("/api/templates/:id", myController.DeleteTemplate)

	return templateRouter
}
