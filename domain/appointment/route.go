package appointment

import (
	"github.com/mukezhz/appointment-booking/pkg/framework"
	"github.com/mukezhz/appointment-booking/pkg/infrastructure"
	"github.com/mukezhz/appointment-booking/pkg/middlewares"
)

// Route struct
type Route struct {
	logger         framework.Logger
	handler        infrastructure.Router
	controller     *Controller
	authMiddleware middlewares.AuthMiddleware
}

// NewRoute creates a new route
func NewRoute(
	logger framework.Logger,
	handler infrastructure.Router,
	controller *Controller,
	authMiddleware middlewares.AuthMiddleware,
) *Route {
	return &Route{
		handler:        handler,
		logger:         logger,
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes sets up appointment routes
func RegisterRoutes(r *Route) {
	r.logger.Info("Setting up appointment routes")

	api := r.handler.Group("/api")

	// Protected routes (require authentication)
	appointments := api.Group("/appointments", r.authMiddleware.HandleAuthWithRole())
	{
		// Availability routes
		appointments.POST("/availability", r.controller.CreateAvailability)
		appointments.GET("/availability", r.controller.GetAvailabilities)

		// Booking routes
		appointments.GET("/bookings", r.controller.GetBookings)
		appointments.GET("/bookings/:id", r.controller.GetBooking)
		appointments.PATCH("/bookings/:id/status", r.controller.UpdateBookingStatus)
	}

	// Public routes (no authentication required)
	public := api.Group("/public/appointments")
	{
		public.POST("/book", r.controller.CreateBooking)
	}
}
