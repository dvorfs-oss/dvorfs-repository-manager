package repository

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BlobHandler struct {
	service *BlobService
}

func NewBlobHandler(service *BlobService) *BlobHandler {
	return &BlobHandler{service: service}
}

// CreateBlobStore godoc
// @Summary Create a new BlobStore
// @Description Create a new BlobStore.
// @Tags BlobStore
// @Accept json
// @Produce json
// @Param blobStore body BlobStore true "BlobStore object"
// @Success 201 {object} BlobStore
// @Router /blob-stores [post]
func (h *BlobHandler) CreateBlobStore(w http.ResponseWriter, r *http.Request) {
	var blobStore BlobStore
	if err := json.NewDecoder(r.Body).Decode(&blobStore); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateBlobStore(&blobStore); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blobStore)
}

// GetBlobStores godoc
// @Summary Get all BlobStores
// @Description Get all BlobStores.
// @Tags BlobStore
// @Produce json
// @Success 200 {array} BlobStore
// @Router /blob-stores [get]
func (h *BlobHandler) GetBlobStores(w http.ResponseWriter, r *http.Request) {
	blobStores, err := h.service.GetBlobStores()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blobStores)
}

// GetBlobStore godoc
// @Summary Get a BlobStore by ID
// @Description Get a BlobStore by ID.
// @Tags BlobStore
// @Produce json
// @Param id path string true "BlobStore ID"
// @Success 200 {object} BlobStore
// @Router /blob-stores/{id} [get]
func (h *BlobHandler) GetBlobStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	blobStore, err := h.service.GetBlobStore(id)
	if err != nil {
		http.Error(w, "BlobStore not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blobStore)
}

// UpdateBlobStore godoc
// @Summary Update a BlobStore
// @Description Update a BlobStore.
// @Tags BlobStore
// @Accept json
// @Produce json
// @Param id path string true "BlobStore ID"
// @Param blobStore body BlobStore true "BlobStore object"
// @Success 200 {object} BlobStore
// @Router /blob-stores/{id} [put]
func (h *BlobHandler) UpdateBlobStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var blobStore BlobStore
	if err := json.NewDecoder(r.Body).Decode(&blobStore); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	blobStore.ID = id
	if err := h.service.UpdateBlobStore(&blobStore); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blobStore)
}

// DeleteBlobStore godoc
// @Summary Delete a BlobStore
// @Description Delete a BlobStore.
// @Tags BlobStore
// @Param id path string true "BlobStore ID"
// @Success 204
// @Router /blob-stores/{id} [delete]
func (h *BlobHandler) DeleteBlobStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteBlobStore(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
