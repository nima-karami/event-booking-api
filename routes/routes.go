package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEventsHandler)
	server.GET("/events/:id", getEventHandler)
	server.PUT("/events/:id", updateEventHandler)
	server.DELETE("/events/:id", deleteEventHandler)
	server.POST("/events", createEventHandler)

	server.GET("/users", getUsersHandler)
	server.GET("/users/:id", getUserHandler)
	server.PUT("/users/:id", updateUserHandler)
	server.DELETE("/users/:id", deleteUserHandler)
	server.POST("/users/signup", userSignupHandler)
	server.POST("/users/login", userLoginHandler)
}
