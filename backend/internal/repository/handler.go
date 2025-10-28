package repository

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// @Summary Get all repositories
// @Description Get a list of all repositories
// @Tags repositories
// @Produce  json
// @Success 200 {array} Repository
// @Router /repositories [get]
func (h *Handler) GetAllRepositories(w http.ResponseWriter, r *http.Request) {
	repos, err := h.service.GetAllRepositories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(repos)
}

// @Summary Create a new repository
// @Description Create a new repository
// @Tags repositories
// @Accept  json
// @Produce  json
// @Param   repository body Repository true "Repository"
// @Success 201
// @Router /repositories [post]
func (h *Handler) CreateRepository(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusCreated)
}

// @Summary Get a repository
// @Description Get a repository by name
// @Tags repositories
// @Produce  json
// @Param   name path string true "Repository name"
// @Success 200 {object} Repository
// @Router /repositories/{name} [get]
func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	repo, err := h.service.GetRepository("test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(repo)
}

// @Summary Update a repository
// @Description Update a repository by name
// @Tags repositories
// @Accept  json
// @Produce  json
// @Param   name path string true "Repository name"
// @Param   repository body Repository true "Repository"
// @Success 200
// @Router /repositories/{name} [put]
func (h *Handler) UpdateRepository(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a repository
// @Description Delete a repository by name
// @Tags repositories
// @Param   name path string true "Repository name"
// @Success 200
// @Router /repositories/{name} [delete]
func (h *Handler) DeleteRepository(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

// @Summary Handle artifact
// @Description Upload or download an artifact
// @Tags artifacts
// @Param   repository-name path string true "Repository name"
// @Param   path path string true "Artifact path"
// @Success 200
// @Router /repository/{repository-name}/{path} [put]
// @Router /repository/{repository-name}/{path} [get]
func (h *Handler) HandleArtifact(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

// @Summary Search artifacts
// @Description Search for artifacts
// @Tags search
// @Produce  json
// @Param   q query string false "Search query"
// @Success 200 {array} Artifact
// @Router /search/artifacts [get]
func (h *Handler) SearchArtifacts(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	artifacts, err := h.service.SearchArtifacts("test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(artifacts)
}
