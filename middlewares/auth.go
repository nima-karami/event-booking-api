package middlewares

import (
	"net/http"
	"strings"

	"example.com/event-booking-api/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "Authorization header missing",
			})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "Invalid token",
			})
		return
	}

	c.Set("userID", userID)
	c.Next()
}
