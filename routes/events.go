package routes

import (
	"net/http"
	"strconv"

	"example.com/event-booking-api/models"
	"example.com/event-booking-api/utils"
	"github.com/gin-gonic/gin"
)

func getEventsHandler(c *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		utils.Logger.Error("Failed to retrieve events", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve events",
		})
		return
	}
	utils.Logger.Debug("Retrieved events", "count", len(events))
	c.JSON(http.StatusOK, events)
}

func getEventHandler(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		utils.Logger.Warn("Invalid event ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID",
		})
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		utils.Logger.Error("Failed to retrieve event", "event_id", eventId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve event",
		})
		return
	}

	utils.Logger.Debug("Retrieved event", "event_id", eventId, "title", event.Title)
	c.JSON(http.StatusOK, event)
}

func createEventHandler(c *gin.Context) {
	event := models.Event{}
	err := c.ShouldBindJSON(&event)

	if err != nil {
		utils.Logger.Warn("Invalid event creation payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	userID := c.GetInt64("userID")
	event.UserID = userID

	err = event.Save()

	if err != nil {
		utils.Logger.Error("Failed to create event",
			"user_id", userID,
			"title", event.Title,
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create event",
		})
		return
	}

	utils.Logger.Info("Event created successfully",
		"event_id", event.ID,
		"title", event.Title,
		"user_id", userID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}

func updateEventHandler(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Warn("Invalid event ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID",
		})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		utils.Logger.Error("Failed to retrieve event for update", "event_id", eventId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve event",
		})
		return
	}

	userID := c.GetInt64("userID")
	role := c.GetString("role")

	if event.UserID != userID && role != "admin" {
		utils.Logger.Warn("Unauthorized event update attempt",
			"event_id", eventId,
			"event_owner", event.UserID,
			"user_id", userID,
			"role", role)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to update this event",
		})
		return
	}

	updatedEvent := models.Event{}
	err = c.ShouldBindJSON(&updatedEvent)
	if err != nil {
		utils.Logger.Warn("Invalid event update payload", "event_id", eventId, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	updatedEvent.ID = eventId

	err = updatedEvent.Update()
	if err != nil {
		utils.Logger.Error("Failed to update event",
			"event_id", eventId,
			"user_id", userID,
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update event",
		})
		return
	}

	utils.Logger.Info("Event updated successfully",
		"event_id", eventId,
		"title", updatedEvent.Title,
		"user_id", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully",
		"event":   updatedEvent,
	})
}

func deleteEventHandler(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Warn("Invalid event ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID",
		})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		utils.Logger.Error("Failed to retrieve event for deletion", "event_id", eventId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve event",
		})
		return
	}

	userID := c.GetInt64("userID")
	role := c.GetString("role")

	if event.UserID != userID && role != "admin" {
		utils.Logger.Warn("Unauthorized event deletion attempt",
			"event_id", eventId,
			"event_owner", event.UserID,
			"user_id", userID,
			"role", role)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to delete this event",
		})
		return
	}

	err = event.Delete()
	if err != nil {
		utils.Logger.Error("Failed to delete event",
			"event_id", eventId,
			"user_id", userID,
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete event",
		})
		return
	}

	utils.Logger.Info("Event deleted successfully",
		"event_id", eventId,
		"title", event.Title,
		"user_id", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
	})
}

func registerEventHandler(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Warn("Invalid event ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID",
		})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		utils.Logger.Error("Failed to retrieve event for registration", "event_id", eventId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve event",
		})
		return
	}

	userID := c.GetInt64("userID")
	err = event.Register(userID)
	if err != nil {
		utils.Logger.Error("Failed to register for event",
			"event_id", eventId,
			"user_id", userID,
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register for event",
		})
		return
	}

	utils.Logger.Info("User registered for event",
		"event_id", eventId,
		"event_title", event.Title,
		"user_id", userID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully registered for event",
	})
}

func unregisterEventHandler(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Warn("Invalid event ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID",
		})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		utils.Logger.Error("Failed to retrieve event for unregistration", "event_id", eventId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve event",
		})
		return
	}

	userID := c.GetInt64("userID")
	err = event.Unregister(userID)
	if err != nil {
		utils.Logger.Error("Failed to unregister from event",
			"event_id", eventId,
			"user_id", userID,
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to unregister from event",
		})
		return
	}

	utils.Logger.Info("User unregistered from event",
		"event_id", eventId,
		"event_title", event.Title,
		"user_id", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully unregistered from event",
	})
}

func getEventRegistrationsHandler(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Warn("Invalid event ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID",
		})
		return
	}

	registrations, err := models.GetRegistrationsByEventIDWithUsers(eventId)
	if err != nil {
		utils.Logger.Error("Failed to retrieve event registrations", "event_id", eventId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve registrations",
		})
		return
	}

	utils.Logger.Debug("Retrieved event registrations",
		"event_id", eventId,
		"count", len(registrations))
	c.JSON(http.StatusOK, registrations)
}
