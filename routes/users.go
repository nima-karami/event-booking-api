package routes

import (
	"net/http"
	"strconv"

	"example.com/event-booking-api/models"
	"example.com/event-booking-api/utils"
	"github.com/gin-gonic/gin"
)

func userSignupHandler(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.Logger.Warn("Invalid signup payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	_, err = models.GetUserByEmail(user.Email)
	if err == nil {
		utils.Logger.Warn("Duplicate user signup attempt", "email", user.Email)
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	err = user.Save()
	if err != nil {
		utils.Logger.Error("Failed to create user", "email", user.Email, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	utils.Logger.Info("User created successfully", "user_id", user.ID, "email", user.Email)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func userLoginHandler(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.Logger.Warn("Invalid login payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	err = user.Authenticate()
	if err != nil {
		utils.Logger.Warn("Authentication failed", "email", user.Email, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication failed",
		})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID, user.Role)
	if err != nil {
		utils.Logger.Error("Failed to generate token", "user_id", user.ID, "email", user.Email, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	utils.Logger.Info("User logged in successfully", "user_id", user.ID, "email", user.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"token":   token,
	})
}

func updateUserHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Warn("Invalid user ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	tokenUserID := c.GetInt64("userID")
	role := c.GetString("role")
	if tokenUserID != userID && role != "admin" {
		utils.Logger.Warn("Unauthorized user update attempt",
			"target_user_id", userID,
			"token_user_id", tokenUserID,
			"role", role)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to update this user",
		})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.Logger.Error("Failed to retrieve user for update", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	err = c.ShouldBindJSON(&user)
	if err != nil {
		utils.Logger.Warn("Invalid user update payload", "user_id", userID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	err = user.Update()
	if err != nil {
		utils.Logger.Error("Failed to update user", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	utils.Logger.Info("User updated successfully", "user_id", userID, "email", user.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

func getUsersHandler(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		utils.Logger.Error("Failed to retrieve users", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	sanitizedUsers := make([]*models.PublicUser, len(users))
	for i, user := range users {
		sanitizedUsers[i] = user.ToPublic()
	}

	utils.Logger.Debug("Retrieved users", "count", len(users))
	c.JSON(http.StatusOK, sanitizedUsers)
}

func getUserHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Warn("Invalid user ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	tokenUserID := c.GetInt64("userID")
	role := c.GetString("role")
	if tokenUserID != userID && role != "admin" {
		utils.Logger.Warn("Unauthorized user view attempt",
			"target_user_id", userID,
			"token_user_id", tokenUserID,
			"role", role)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to view this user",
		})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.Logger.Error("Failed to retrieve user", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	sanitizedUser := user.ToPublic()

	utils.Logger.Debug("Retrieved user", "user_id", userID, "email", user.Email)
	c.JSON(http.StatusOK, sanitizedUser)
}

func deleteUserHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Warn("Invalid user ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	role := c.GetString("role")
	tokenUserID := c.GetInt64("userID")
	if userID != tokenUserID && role != "admin" {
		utils.Logger.Warn("Unauthorized user deletion attempt",
			"target_user_id", userID,
			"token_user_id", tokenUserID,
			"role", role)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to delete this user",
		})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.Logger.Error("Failed to retrieve user for deletion", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	err = user.Delete()
	if err != nil {
		utils.Logger.Error("Failed to delete user", "user_id", userID, "email", user.Email, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	utils.Logger.Info("User deleted successfully", "user_id", userID, "email", user.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func updateUserRoleHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Warn("Invalid user ID parameter", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.Logger.Error("Failed to retrieve user for role update", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	var payload struct {
		Role string `json:"role"`
	}

	err = c.ShouldBindJSON(&payload)
	if err != nil {
		utils.Logger.Warn("Invalid role update payload", "user_id", userID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	oldRole := user.Role
	user.Role = payload.Role

	err = user.Update()
	if err != nil {
		utils.Logger.Error("Failed to update user role",
			"user_id", userID,
			"old_role", oldRole,
			"new_role", payload.Role,
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user role",
		})
		return
	}

	utils.Logger.Info("User role updated successfully",
		"user_id", userID,
		"email", user.Email,
		"old_role", oldRole,
		"new_role", user.Role)
	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"user":    user,
	})
}
