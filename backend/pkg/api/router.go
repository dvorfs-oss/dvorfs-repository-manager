package api

import (
	"dvorfs-repository-manager/internal/auth"
	"dvorfs-repository-manager/internal/cleanup"
	"dvorfs-repository-manager/internal/repository"
	"dvorfs-repository-manager/internal/user"
	"github.com/gorilla/mux"
)

func NewRouter(
	authHandler *auth.Handler,
	repoHandler *repository.Handler,
	userHandler *user.Handler,
	cleanupHandler *cleanup.Handler,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Authentication
	authRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	authRouter.HandleFunc("/me", authHandler.GetMe).Methods("GET")

	// Repositories
	repoRouter := router.PathPrefix("/api/v1/repositories").Subrouter()
	repoRouter.HandleFunc("/", repoHandler.GetAllRepositories).Methods("GET")
	repoRouter.HandleFunc("/", repoHandler.CreateRepository).Methods("POST")
	repoRouter.HandleFunc("/{name}", repoHandler.GetRepository).Methods("GET")
	repoRouter.HandleFunc("/{name}", repoHandler.UpdateRepository).Methods("PUT")
	repoRouter.HandleFunc("/{name}", repoHandler.DeleteRepository).Methods("DELETE")

	// Artifacts
	artifactRouter := router.PathPrefix("/repository").Subrouter()
	artifactRouter.PathPrefix("/{repository-name}").HandlerFunc(repoHandler.HandleArtifact)

	// Search
	searchRouter := router.PathPrefix("/api/v1/search").Subrouter()
	searchRouter.HandleFunc("/artifacts", repoHandler.SearchArtifacts).Methods("GET")

	// Security
	securityRouter := router.PathPrefix("/api/v1/security").Subrouter()
	securityRouter.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	securityRouter.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	securityRouter.HandleFunc("/users/{username}", userHandler.UpdateUser).Methods("PUT")
	securityRouter.HandleFunc("/users/{username}/password", userHandler.ChangeUserPassword).Methods("PUT")
	securityRouter.HandleFunc("/users/{username}", userHandler.DeleteUser).Methods("DELETE")
	securityRouter.HandleFunc("/roles", userHandler.GetAllRoles).Methods("GET")
	securityRouter.HandleFunc("/roles", userHandler.CreateRole).Methods("POST")
	securityRouter.HandleFunc("/roles/{roleId}", userHandler.UpdateRole).Methods("PUT")
	securityRouter.HandleFunc("/roles/{roleId}", userHandler.DeleteRole).Methods("DELETE")

	// Cleanup Policies
	cleanupRouter := router.PathPrefix("/api/v1/cleanup-policies").Subrouter()
	cleanupRouter.HandleFunc("/", cleanupHandler.GetAllCleanupPolicies).Methods("GET")
	cleanupRouter.HandleFunc("/", cleanupHandler.CreateCleanupPolicy).Methods("POST")
	cleanupRouter.HandleFunc("/{policyId}", cleanupHandler.UpdateCleanupPolicy).Methods("PUT")
	cleanupRouter.HandleFunc("/{policyId}", cleanupHandler.DeleteCleanupPolicy).Methods("DELETE")

	return router
}
