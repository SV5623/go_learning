package main

import 
(
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	{
		v1.POST("/auth/register", app.createUserHandler)
		
		v1.POST("/events", app.createEventHandler)
		v1.GET("/events", app.listEventsHandler)
		v1.GET("/events/:id", app.getEventHandler)
		v1.PUT("/events/:id", app.updateEventHandler)
		v1.DELETE("/events/:id", app.deleteEventHandler)
		// v1.POST("/events/:id/attendees", app.addAttendeeHandler)
	}
	
	return g
}

