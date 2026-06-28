package database

import (
	"fmt"
	"log"

	"github.com/sanskarjain/authorization/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to Database successfully!")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.ProjectUser{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.AuthKey{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database Migration completed!")
}
