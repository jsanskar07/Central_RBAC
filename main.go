package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sanskarjain/authorization/database"
	"github.com/sanskarjain/authorization/routes"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found or failed to load, falling back to environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=auth_db port=5432 sslmode=disable"
	}
	
	database.ConnectDB(dsn)
	database.Migrate()

	fmt.Println("Central Authorization System started.")

	// Set up Gin Router
	router := gin.Default()
	
	// Initialize Routes
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s...\n", port)
	router.Run(":" + port)
}
