package router

import (
	"github.com/gin-gonic/gin"

	parkingH "parking/internal/http/handler/parking"
)

// Router is settings for http routs.
type Router struct {
	app *gin.Engine

	parking *parkingH.Handler
}

// New creates http router.
func New(
	app *gin.Engine,
	parking *parkingH.Handler,
) *Router {
	return &Router{
		app:     app,
		parking: parking,
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
}
