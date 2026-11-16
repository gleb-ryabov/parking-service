package router

import (
	"github.com/gin-gonic/gin"

	parkingH "parking/internal/http/handler/parking"
	scheduleH "parking/internal/http/handler/schedule"
)

// Router is settings for http routs.
type Router struct {
	app *gin.Engine

	parking  *parkingH.Handler
	schedule *scheduleH.Handler
}

// New creates http router.
func New(
	app *gin.Engine,
	parking *parkingH.Handler,
	schedule *scheduleH.Handler,
) *Router {
	return &Router{
		app:      app,
		parking:  parking,
		schedule: schedule,
	}
}

// Setup - configure router.
func (r *Router) Setup() {
	api := r.app.Group("/api")

	parkings := api.Group("/parking")
	{
		parkings.GET("", r.parking.GetParkings)
		parkings.POST("", r.parking.CreateParking)
		parkings.PUT("/:id", r.parking.UpdateParking)
		parkings.DELETE("/:id", r.parking.DeleteParking)
	}

	schedules := api.Group("/schedule")
	{
		schedules.GET("", r.schedule.GetSchedules)
		schedules.POST("", r.schedule.CreateSchedule)
		schedules.PUT("/:id", r.schedule.UpdateSchedule)
		schedules.DELETE("/:id", r.schedule.DeleteSchedule)
	}
}
