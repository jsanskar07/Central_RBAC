package controllers

import (
	"net/http"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanskarjain/authorization/database"
	"github.com/sanskarjain/authorization/models"
)

func CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if project.APIKey == "" {
		project.APIKey = GenerateSecureAPIKey()
	}

	if err := database.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func GetProjects(c *gin.Context) {
	var projects []models.Project
	database.DB.Find(&projects)
	c.JSON(http.StatusOK, projects)
}

func GenerateAPIKey(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Generate a new secure API Key
	project.APIKey = GenerateSecureAPIKey()

	// Save back to database
	if err := database.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API Key"})
		return
	}

	// Return the updated project including the new API Key
	c.JSON(http.StatusOK, gin.H{
		"message":      "API Key successfully generated",
		"project_id":   project.ID,
		"project_name": project.Name,
		"api_key":      project.APIKey,
		"updated_at":   project.UpdatedAt,
	})
}

// GenerateSecureAPIKey creates a standard 32-character secure token
func GenerateSecureAPIKey() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return "pk_live_" + string(b)
}
