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

func (h *Handler) GetAllRepositories(w http.ResponseWriter, r *http.Request) {
	repos, err := h.service.GetAllRepositories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(repos)
}

func (h *Handler) CreateRepository(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	repo, err := h.service.GetRepository("test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(repo)
}

func (h *Handler) UpdateRepository(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteRepository(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleArtifact(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SearchArtifacts(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	artifacts, err := h.service.SearchArtifacts("test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(artifacts)
}
