package user

import (
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/middlewares"
)

type Route struct {
	logger         framework.Logger
	router         infrastructure.Router
	controller     *Controller
	authMiddleware middlewares.AuthMiddleware
}

func NewRoute(
	logger framework.Logger,
	router infrastructure.Router,
	controller *Controller,
	authMiddleware middlewares.AuthMiddleware,
) *Route {
	return &Route{
		logger:         logger,
		router:         router,
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}

func RegisterRoutes(r *Route) {
	routes := r.router.Group("/users")

	routes.POST("/register", r.controller.Register)
	routes.POST("/login", r.controller.Login)

	protected := routes.Group("", r.authMiddleware.HandleAuthWithRole())
	{
		protected.GET("/profile", r.controller.GetProfile)
		protected.PUT("/profile", r.controller.UpdateProfile)
	}
}
