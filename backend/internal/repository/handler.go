package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
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
	w.Header().Set("Content-Type", "application/json")
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
	var repo Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateRepository(&repo); err != nil {
		writeRepositoryError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(repo)
}

// @Summary Get a repository
// @Description Get a repository by name
// @Tags repositories
// @Produce  json
// @Param   name path string true "Repository name"
// @Success 200 {object} Repository
// @Router /repositories/{name} [get]
func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	repo, err := h.service.GetRepository(name)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
	name := mux.Vars(r)["name"]
	var repo Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	repo.Name = name

	if err := h.service.UpdateRepository(&repo); err != nil {
		writeRepositoryError(w, err)
		return
	}

	updated, err := h.service.GetRepository(name)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// @Summary Delete a repository
// @Description Delete a repository by name
// @Tags repositories
// @Param   name path string true "Repository name"
// @Success 200
// @Router /repositories/{name} [delete]
func (h *Handler) DeleteRepository(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if err := h.service.DeleteRepository(name); err != nil {
		writeRepositoryError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Handle artifact
// @Description Upload, download, or delete an artifact
// @Tags artifacts
// @Param   repository-name path string true "Repository name"
// @Param   path path string true "Artifact path"
// @Success 200
// @Router /repository/{repository-name}/{path} [put]
// @Router /repository/{repository-name}/{path} [get]
// @Router /repository/{repository-name}/{path} [delete]
func (h *Handler) HandleArtifact(w http.ResponseWriter, r *http.Request) {
	repoName, artifactPath, err := artifactRouteParts(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut, http.MethodPost:
		artifact, err := h.service.UploadArtifact(repoName, artifactPath, r.Header.Get("Content-Type"), r.Body)
		if err != nil {
			writeRepositoryError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(artifact)
	case http.MethodGet:
		file, artifact, err := h.service.OpenArtifact(repoName, artifactPath)
		if err != nil {
			writeRepositoryError(w, err)
			return
		}
		defer file.Close()

		contentType := strings.TrimSpace(artifact.ContentType)
		if contentType == "" {
			contentType = http.DetectContentType(nil)
			if contentType == "" {
				contentType = "application/octet-stream"
			}
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, pathBase(artifactPath)))
		if _, err := io.Copy(w, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		if err := h.service.DeleteArtifact(repoName, artifactPath); err != nil {
			writeRepositoryError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary Search artifacts
// @Description Search for artifacts
// @Tags search
// @Produce  json
// @Param   q query string false "Search query"
// @Success 200 {array} Artifact
// @Router /search/artifacts [get]
func (h *Handler) SearchArtifacts(w http.ResponseWriter, r *http.Request) {
	artifacts, err := h.service.SearchArtifacts(r.URL.Query().Get("q"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artifacts)
}

func artifactRouteParts(r *http.Request) (string, string, error) {
	repoName := mux.Vars(r)["repository-name"]
	if repoName == "" {
		return "", "", ErrInvalidRepositoryName
	}

	prefix := "/repository/" + repoName + "/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		return "", "", ErrInvalidArtifactPath
	}

	artifactPath := strings.TrimPrefix(r.URL.Path, prefix)
	if artifactPath == "" {
		return "", "", ErrInvalidArtifactPath
	}

	return repoName, artifactPath, nil
}

func writeRepositoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrRepositoryNotFound), errors.Is(err, ErrArtifactNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidRepositoryName), errors.Is(err, ErrInvalidArtifactPath), errors.Is(err, ErrUnsupportedRepositoryType):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func pathBase(p string) string {
	if idx := strings.LastIndexAny(p, `/\`); idx >= 0 {
		return p[idx+1:]
	}
	return p
}
