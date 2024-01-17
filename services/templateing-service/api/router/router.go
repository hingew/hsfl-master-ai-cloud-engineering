package router

import (
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/middleware"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/controller"
)

func NewTemplateRouter(
	templateController controller.ControllerInterface,
	middleware middleware.AuthMiddleInterface,
) *router.Router {
	templateRouter := router.New()
	templateRouter.GET("/api/health/templates", health.Check)

	authRouter := templateRouter.USE("/api/templates*", middleware.AuthMiddleware)
	authRouter.GET("/api/templates", templateController.GetAllTemplates)
	authRouter.GET("/api/templates/:id", templateController.GetTemplate)
	authRouter.GET("/api/templates/:id/controller_coalescing", templateController.GetTemplateWithCoalecing)
	authRouter.POST("/api/templates", templateController.CreateTemplate)
	authRouter.PUT("/api/templates/:id", templateController.UpdateTemplate)
	authRouter.DELETE("/api/templates/:id", templateController.DeleteTemplate)

	return templateRouter
}
