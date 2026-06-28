package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanskarjain/authorization/database"
	"github.com/sanskarjain/authorization/models"
)

func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, role)
}

func GetRoles(c *gin.Context) {
	projectID := c.Query("project_id")
	var roles []models.Role
	
	query := database.DB
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	
	query.Find(&roles)
	c.JSON(http.StatusOK, roles)
}

func AssignRoleToUser(c *gin.Context) {
	var userRole models.UserRole
	if err := c.ShouldBindJSON(&userRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&userRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role to user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Role assigned successfully"})
}
