package availability

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
	api := r.router.Group("/api/v1")

	// Availabilities
	availabilities := api.Group("/availability").Use(r.authMiddleware.HandleAuthWithRole("doctor", "patient"))
	{
		availabilities.POST("", r.controller.CreateAvailability)         // Create availability slot
		availabilities.PUT("/:uuid", r.controller.UpdateAvailability)    // Update availability slot
		availabilities.DELETE("/:uuid", r.controller.DeleteAvailability) // Delete availability slot
		availabilities.GET("", r.controller.ListAvailability)            // List available slots
	}

	// Appointments
	appointments := api.Group("/appointments").Use(r.authMiddleware.HandleAuthWithRole("doctor", "patient"))
	{
		appointments.POST("", r.controller.BookAppointment)           // Book appointment
		appointments.DELETE("/:uuid", r.controller.CancelAppointment) // Cancel appointment
		appointments.GET("", r.controller.ListAppointments)           // List appointments
	}

}
