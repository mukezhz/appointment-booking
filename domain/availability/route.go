package availability

import (
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/middlewares"
)

type Route struct {
	controller     *Controller
	router         infrastructure.Router
	authMiddleware middlewares.AuthMiddleware
}

func NewRoute(
	controller *Controller,
	router infrastructure.Router,
	authMiddleware middlewares.AuthMiddleware,
) *Route {
	return &Route{
		controller:     controller,
		router:         router,
		authMiddleware: authMiddleware,
	}
}

func RegisterRoutes(r *Route) {
	auth := r.router.Group("/api/v1")
	{
		// Protected routes (require authentication)
		protected := auth.Group("/availability")
		protected.Use(r.authMiddleware.HandleAuthWithRole())
		{
			protected.POST("", r.controller.CreateAvailability)
			protected.GET("", r.controller.GetAvailability)
		}
	}
}
