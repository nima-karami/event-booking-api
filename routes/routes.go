package routes

import (
	"example.com/event-booking-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	server.GET("/events", getEventsHandler)
	server.GET("/events/:id", getEventHandler)
	server.GET("/events/:id/registrations", getEventRegistrationsHandler)

	authenticated.POST("/events", createEventHandler)
	authenticated.PUT("/events/:id", updateEventHandler)
	authenticated.DELETE("/events/:id", deleteEventHandler)
	authenticated.POST("/events/:id/register", registerEventHandler)
	authenticated.DELETE("/events/:id/register", unregisterEventHandler)

	server.GET("/users", getUsersHandler)
	server.GET("/users/:id", getUserHandler)
	server.POST("/users/signup", userSignupHandler)
	server.POST("/users/login", userLoginHandler)

	authenticated.PUT("/users/:id", updateUserHandler)
	authenticated.DELETE("/users/:id", deleteUserHandler)
}
