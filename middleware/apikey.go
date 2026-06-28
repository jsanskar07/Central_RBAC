package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanskarjain/authorization/database"
	"github.com/sanskarjain/authorization/models"
)

func RequireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing X-API-Key header"})
			return
		}

		var project models.Project
		if err := database.DB.Where("api_key = ?", apiKey).First(&project).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			return
		}

		// Inject the project ID securely into the context so controllers don't have to trust the payload
		c.Set("project_id", project.ID)
		
		c.Next()
	}
}
