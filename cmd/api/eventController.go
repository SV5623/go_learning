package main

import (
	"log"
	"net/http"
	"rest_app_in_gin/internal/database"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func (app *application) getEventHandler(c *gin.Context) {	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(200, event)
}

func (app *application) listEventsHandler(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(200, events)
}

func (app *application) updateEventHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	_, err = app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	event.ID = id

	if err := app.models.Events.Update(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (app *application) deleteEventHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(200, gin.H{"message": "Event deleted successfully"})
}

func (app *application) createEventHandler(c *gin.Context) {
	var newEvent database.Event
	
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	
	if err := app.models.Events.Create(&newEvent); err != nil {
	//	c.JSON(500, gin.H{"error": "Failed to create event"})
	log.Println("create event:", err)
	c.JSON(500, gin.H{"error": err.Error()})
	return
	}
	c.JSON(201, newEvent)
}

