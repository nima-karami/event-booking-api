package routes

import (
	"net/http"
	"strconv"

	"example.com/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

func getEventsHandler(c *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve events",
		})
		return
	}
	c.JSON(http.StatusOK, events)
}

func getEventHandler(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID",
		})
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve event",
		})
		return
	}

	c.JSON(http.StatusOK, event)
}

func createEventHandler(c *gin.Context) {
	event := models.Event{}
	err := c.ShouldBindJSON(&event)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	event.UserID = 1 // Placeholder for UserID assignment

	err = event.Save()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create event",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}
