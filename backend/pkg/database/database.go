package database

import (
	"fmt"
	"log"
	"os"

	"dvorfs-repository-manager/internal/repository"
	"dvorfs-repository-manager/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")
}

func Migrate() {
	DB.AutoMigrate(&user.User{}, &user.Role{}, &repository.Repository{}, &repository.Artifact{}, &repository.CleanupPolicy{}, &repository.BlobStore{})
	log.Println("Database migration completed")

	var admin user.User
	if err := DB.Where("username = ?", "admin").First(&admin).Error; err != nil {
		if err := DB.Create(&user.User{
			Username:     "admin",
			PasswordHash: "admin",
			Email:        "admin@local",
		}).Error; err != nil {
			log.Println("Failed to seed default admin user:", err)
		} else {
			log.Println("Seeded default admin user")
		}
	}
}
