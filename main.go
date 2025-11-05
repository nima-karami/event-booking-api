package main

import (
	"net/http"

	"example.com/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

func main() {
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
	events := models.GetAllEvents()
	c.JSON(http.StatusOK, events)
}

func createEventHandler(c *gin.Context) {
	event := models.Event{}
	err := c.ShouldBindJSON(&event)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	event.ID = 1     // Placeholder for ID assignment
	event.UserID = 1 // Placeholder for UserID assignment

	event.Save()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}
