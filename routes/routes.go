package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sanskarjain/authorization/controllers"
	"github.com/sanskarjain/authorization/middleware"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Projects (Unprotected)
		api.POST("/projects", controllers.CreateProject)
		api.GET("/projects", controllers.GetProjects)
		api.POST("/projects/:id/api-key", controllers.GenerateAPIKey)

		// Protected Routes
		protected := api.Group("/")
		protected.Use(middleware.RequireAPIKey())
		{
			// Auth
			protected.POST("/auth/register", controllers.Register)
			protected.POST("/auth/login", controllers.Login)

			// Roles
			protected.POST("/roles", controllers.CreateRole)
			protected.GET("/roles", controllers.GetRoles)
			protected.POST("/roles/assign", controllers.AssignRoleToUser)

			// Permissions
			protected.POST("/permissions", controllers.CreatePermission)
			protected.POST("/permissions/assign", controllers.AssignPermissionToRole)
		}
	}
}
