package main

import (
	"net/http"

	"example.com/event-booking-api/db"
	"example.com/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	server.GET("/ping", pongHandler)
	server.GET("/events", getEventsHandler)
	server.POST("/events", createEventHandler)

	server.Run(":8080")
}

func pongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

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

func createEventHandler(c *gin.Context) {
	event := models.Event{}
	err := c.ShouldBindJSON(&event)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	event.ID = 1     // Placeholder for ID assignment
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
