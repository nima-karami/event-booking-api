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
	userID, role, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "Invalid token",
			})
		return
	}

	c.Set("userID", userID)
	c.Set("role", role)
	c.Next()
}

func RequireRole(c *gin.Context, allowedRoles ...string) error {
	role, exists := c.Get("role")
	if !exists {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error": "Role not found in context",
			})
		return nil
	}

	roleStr, ok := role.(string)
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error": "Invalid role type",
			})
		return nil
	}

	for _, allowedRole := range allowedRoles {
		if roleStr == allowedRole {
			return nil
		}
	}

	c.AbortWithStatusJSON(
		http.StatusForbidden,
		gin.H{
			"error": "Insufficient permissions",
		})
	return nil
}

func AuthorizeAdmin(c *gin.Context) {
	err := RequireRole(c, "admin")
	if err != nil {
		return
	}
	c.Next()
}
