package booking

import (
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/middlewares"
)

type Router struct {
	controller     *Controller
	router         infrastructure.Router
	authMiddleware middlewares.AuthMiddleware
}

func NewRouter(
	controller *Controller,
	router infrastructure.Router,
	authMiddleware middlewares.AuthMiddleware,
) *Router {
	return &Router{
		controller:     controller,
		router:         router,
		authMiddleware: authMiddleware,
	}
}

func RegisterRoutes(r *Router) {
	api := r.router.Group("/api/v1")
	{
		// Public routes
		public := api.Group("/public")
		{
			public.POST("/book", r.controller.CreateBooking)
		}

		// Protected routes
		protected := api.Group("/bookings")
		protected.Use(r.authMiddleware.HandleAuthWithRole())
		{
			protected.GET("", r.controller.GetBookings)
		}
	}
}
