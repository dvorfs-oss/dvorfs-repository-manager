package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"dvorfs-repository-manager/internal/repository"
	"dvorfs-repository-manager/internal/user"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	driver := os.Getenv("DB_DRIVER")
	switch driver {
	case "", "postgres":
		host := os.Getenv("DB_HOST")
		if host == "" && driver == "" {
			DB, err = connectSQLite()
		} else {
			DB, err = connectPostgres()
		}
	case "sqlite":
		DB, err = connectSQLite()
	default:
		log.Fatalf("Unsupported DB_DRIVER %q", driver)
	}

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

func connectPostgres() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func connectSQLite() (*gorm.DB, error) {
	dbPath := os.Getenv("SQLITE_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(".", "data", "dvorfs.db")
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, err
	}

	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}
