package main

import (
	"example.com/event-booking-api/db"
	"example.com/event-booking-api/routes"
	"example.com/event-booking-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	utils.InitLogger()
	utils.Logger.Info("Starting Event Booking API")

	db.InitDB()
	server := gin.New()
	server.Use(gin.Recovery())

	routes.RegisterRoutes(server)

	utils.Logger.Info("Server starting", "port", 8080)
	if err := server.Run(":8080"); err != nil {
		utils.Logger.Error("Failed to start server", "error", err)
	}
}
