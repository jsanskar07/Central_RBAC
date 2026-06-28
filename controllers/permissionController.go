package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanskarjain/authorization/database"
	"github.com/sanskarjain/authorization/models"
)

func CreatePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permission"})
		return
	}

	c.JSON(http.StatusCreated, permission)
}

func AssignPermissionToRole(c *gin.Context) {
	var rolePerm models.RolePermission
	if err := c.ShouldBindJSON(&rolePerm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&rolePerm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign permission to role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Permission assigned to role successfully"})
}
