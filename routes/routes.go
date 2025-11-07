package routes

import (
	"example.com/event-booking-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", userSignupHandler)
	server.POST("/users/login", userLoginHandler)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	// Event routes
	authenticated.GET("/events", getEventsHandler)
	authenticated.POST("/events", createEventHandler)
	authenticated.GET("/events/:id", getEventHandler)
	authenticated.PUT("/events/:id", updateEventHandler)
	authenticated.DELETE("/events/:id", deleteEventHandler)

	// Event registration routes
	authenticated.POST("/events/:id/register", registerEventHandler)
	authenticated.DELETE("/events/:id/register", unregisterEventHandler)
	authenticated.GET("/events/:id/registrations", getEventRegistrationsHandler)

	// User routes
	authenticated.GET("/users/:id", getUserHandler)
	authenticated.PUT("/users/:id", updateUserHandler)
	authenticated.DELETE("/users/:id", deleteUserHandler)

	// Admin-only routes
	admin := authenticated.Group("/admin")
	admin.Use(middlewares.Authenticate, middlewares.AuthorizeAdmin)
	admin.GET("/users", getUsersHandler)
	admin.PUT("/users/:id/role", updateUserRoleHandler)
}
