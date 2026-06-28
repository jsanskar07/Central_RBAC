package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sanskarjain/authorization/database"
	"github.com/sanskarjain/authorization/models"
	"github.com/sanskarjain/authorization/utils"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// User doesn't exist globally, create them
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user = models.User{
			Email:        input.Email,
			PasswordHash: string(hashedPassword),
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Extract project_id from middleware context
	projectID := c.MustGet("project_id").(uint)

	// Check if already enrolled in project
	var projectUser models.ProjectUser
	if err := database.DB.Where("user_id = ? AND project_id = ?", user.ID, projectID).First(&projectUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already registered in this project"})
		return
	}

	// Enroll user in project
	projectUser = models.ProjectUser{
		UserID:    user.ID,
		ProjectID: projectID,
	}
	if err := database.DB.Create(&projectUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll user in project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": user.ID})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Extract project_id from middleware context
	projectID := c.MustGet("project_id").(uint)

	// Check if user is enrolled in this project
	var projectUser models.ProjectUser
	if err := database.DB.Where("user_id = ? AND project_id = ?", user.ID, projectID).First(&projectUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not belong to this project"})
		return
	}

	// In a real scenario, you'd fetch the user's roles for the specific ProjectID here
	var roles []string
	roles = append(roles, "user") // Default dummy role

	token, err := utils.GenerateToken(user.ID, user.Email, roles, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetPublicKey serves the RSA public key for verification by external services
func GetPublicKey(c *gin.Context) {
	pubData, err := os.ReadFile("public_key.pem")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read public key"})
		return
	}
	c.String(http.StatusOK, string(pubData))
}
