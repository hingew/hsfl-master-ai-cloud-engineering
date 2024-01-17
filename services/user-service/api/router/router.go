package router

import (
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/controller"
)

func NewUserRouter(
	userController controller.ControllerInterface,
) *router.Router {
	userRouter := router.New()
	userRouter.GET("/api/health/user", health.Check)

	userRouter.POST("/auth/login", userController.Login)
	userRouter.POST("/auth/register", userController.Register)

	return userRouter
}
