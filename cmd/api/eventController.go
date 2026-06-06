package main

import (
	"log"
	"net/http"
	"rest_app_in_gin/internal/database"
	"strconv"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

// getEvent returns event by id
//
// @Summary Get event by ID
// @Description Returns a single event by its ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} database.Event
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/events/{id} [get]
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
// listEvents returns all events
//
// @Summary Returns all events
// @Description Returns all events
// @Tags Events
// @Accept json
// @Produce json
// @Success 200 {array} database.Event
// @Router /api/v1/events [get] 
func (app *application) listEventsHandler(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(200, events)
}
// updateEvent updates an existing event
//
// @Summary Update event
// @Description Updates an event owned by the authenticated user
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Param event body database.Event true "Updated event"
// @Success 200 {object} database.Event
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id} [put]
func (app *application) updateEventHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.GetUserFromContext(c)

	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
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
// deleteEvent deletes an event
//
// @Summary Delete event
// @Description Deletes an event owned by the authenticated user
// @Tags Events
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id} [delete]
func (app *application) deleteEventHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid event ID"})
		return
	}

		user := app.GetUserFromContext(c)

		existingEvent, err := app.models.Events.Get(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "Event not found"})
			return
		}
		if existingEvent == nil {
			c.JSON(404, gin.H{"error": "Event not found"})
			return
		}
		if existingEvent.OwnerId != user.Id {
			c.JSON(403, gin.H{"error": "You are not authorized to delete this event"})
			return
		}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(200, gin.H{"message": "Event deleted successfully"})
}
// createEvent creates a new event
//
// @Summary Create event
// @Description Creates a new event for the authenticated user
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event body database.Event true "Event data"
// @Success 201 {object} database.Event
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events [post]
func (app *application) createEventHandler(c *gin.Context) {
	var event database.Event
	
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	
	user := app.GetUserFromContext(c)
	event.OwnerId = user.Id

	if err := app.models.Events.Create(&event); err != nil {
	//	c.JSON(500, gin.H{"error": "Failed to create event"})
	log.Println("create event:", err)
	c.JSON(500, gin.H{"error": err.Error()})
	return
	}
	c.JSON(201, event)
}
// addAttendee adds attendee to event
//
// @Summary Add attendee
// @Description Adds a user as attendee to an event
// @Tags Attendees
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 201 {object} database.Attendee
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{user_id}/attendees [post]
func (app *application) addAttendeeHandler(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	event , err := app.models.Attendees.Get(user_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd, err := app.models.Users.Get(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.ID, userToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee already exists for this event and user"})
		return
	}

	attendee := database.Attendee{
		EventID: event.ID,
		UserID: userToAdd.Id,
	}

	_ , err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}
// addAttendeeToEvent adds a user to an event
//
// @Summary Add attendee to event
// @Description Adds a user to an event. Only the event owner can perform this action.
// @Tags Attendees
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Param userId path int true "User ID"
// @Success 201 {object} database.Attendee
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id}/attendees/{userId} [post]
func (app *application) addAttendeeToEventHandler(c *gin.Context) {
	
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user Id"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}

	userToAdd, err := app.models.Users.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive user"})
		return
	}

	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}

	user := app.GetUserFromContext(c)

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add an attendee"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.ID, userToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee already exists"})
		return
	}

	attendee := database.Attendee{
		EventID: event.ID,
		UserID:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add  attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)

}

func (app *application) deleteAttendeeFromEventHandler(c *gin.Context) {
	id , err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
	return
	}
	userId , err := strconv.Atoi(c.Param("userId"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
	return
	}

	event , err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	
	user := app.GetUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete attendees from this event"})
		return
	}

	err = app.models.Attendees.Delete(id, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee from event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendee deleted from event successfully"})
}
// getEventsByAttendee returns all events for attendee
//
// @Summary Get events by attendee
// @Description Returns all events that a user is attending
// @Tags Attendees
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} database.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/attendees/{id}/events [get]
func (app *application) getEventsByAttendeeHandler(c *gin.Context) {
	id , err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
	return
	}
	events, err := app.models.Attendees.GetEventsByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees for event"})
		return
	}

	c.JSON(http.StatusOK, events)
}
// getAttendeesForEvent returns all attendees for an event
//
// @Summary Get attendees for event
// @Description Returns all attendees registered for an event
// @Tags Attendees
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {array} database.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id}/attendees [get]
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id , err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
	return
	}
	attendees, err := app.models.Attendees.GetAttendeesByEventId(id) //у відповідь на запит GET /events/:id/attendees дає список користувачів, які є учасниками цього заходу
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees for event"})
		return
	}

	c.JSON(http.StatusOK, attendees)
}