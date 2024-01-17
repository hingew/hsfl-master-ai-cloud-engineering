package router

import (
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/controller"
)

func NewTemplateRouter(
	templateController controller.ControllerInterface,
) *router.Router {
	templateRouter := router.New()
	templateRouter.GET("/api/health/templates", health.Check)

	templateRouter.GET("/api/templates", templateController.GetAllTemplates)
	templateRouter.GET("/api/templates/:id", templateController.GetTemplate)
	templateRouter.GET("/api/templates/:id/controller_coalescing", templateController.GetTemplateWithCoalecing)
	templateRouter.POST("/api/templates", templateController.CreateTemplate)
	templateRouter.PUT("/api/templates/:id", templateController.UpdateTemplate)
	templateRouter.DELETE("/api/templates/:id", templateController.DeleteTemplate)

	return templateRouter
}
