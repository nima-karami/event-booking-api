package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEventsHandler)
	server.GET("/events/:id", getEventHandler)
	server.POST("/events", createEventHandler)
}
