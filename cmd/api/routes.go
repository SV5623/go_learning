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
		
		
		v1.GET("/events/:id", app.getEventHandler)
		
		v1.GET("/events", app.listEventsHandler)
	
		
		v1.POST("/events/:id/attendees", app.addAttendeeHandler)
		
		v1.GET("/events/:id/attendees", app.getAttendeesForEvent)
		
		v1.GET("/attendees/:id/events", app.getEventsByAttendeeHandler)

		v1.POST("/auth/login", app.loginUserHandler)
		v1.POST("/auth/register", app.createUserHandler)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.POST("/events", app.createEventHandler)
		authGroup.PUT("/events/:id", app.updateEventHandler)
		authGroup.DELETE("/events/:id", app.deleteEventHandler)
		authGroup.POST("/events/:id/attendees/:userId", app.addAttendeeToEventHandler)
		authGroup.DELETE("/events/:id/attendees/:userId", app.deleteAttendeeFromEventHandler)
	}
	
	

	return g
}

