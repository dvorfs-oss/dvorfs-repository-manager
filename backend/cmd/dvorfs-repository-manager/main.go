package main

import (
	"log"
	"net/http"

	"dvorfs-repository-manager/internal/auth"
	"dvorfs-repository-manager/internal/cleanup"
	"dvorfs-repository-manager/internal/repository"
	"dvorfs-repository-manager/internal/user"
	"dvorfs-repository-manager/pkg/api"
	"dvorfs-repository-manager/pkg/database"
)

func main() {
	database.Connect()
	database.Migrate()

	// Initialize services
	authService := auth.NewService()
	repoService := repository.NewService()
	userService := user.NewService()
	cleanupService := cleanup.NewService()

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	repoHandler := repository.NewHandler(repoService)
	userHandler := user.NewHandler(userService)
	cleanupHandler := cleanup.NewHandler(cleanupService)

	// Initialize router
	router := api.NewRouter(authHandler, repoHandler, userHandler, cleanupHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
