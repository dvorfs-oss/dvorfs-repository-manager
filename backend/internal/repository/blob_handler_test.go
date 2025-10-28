package repository

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupBlobRouter() (*mux.Router, *BlobService) {
	router := mux.NewRouter()
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&BlobStore{}, &Repository{})
	service := NewBlobService(db)
	handler := NewBlobHandler(service)

	router.HandleFunc("/blob-stores", handler.CreateBlobStore).Methods("POST")
	router.HandleFunc("/blob-stores", handler.GetBlobStores).Methods("GET")
	router.HandleFunc("/blob-stores/{id}", handler.GetBlobStore).Methods("GET")
	router.HandleFunc("/blob-stores/{id}", handler.UpdateBlobStore).Methods("PUT")
	router.HandleFunc("/blob-stores/{id}", handler.DeleteBlobStore).Methods("DELETE")

	return router, service
}

func TestCreateBlobStore(t *testing.T) {
	router, _ := setupBlobRouter()

	blobStore := BlobStore{Name: "Test Blob Store", Type: "file"}
	jsonValue, _ := json.Marshal(blobStore)
	req, _ := http.NewRequest("POST", "/blob-stores", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdBlobStore BlobStore
	json.Unmarshal(w.Body.Bytes(), &createdBlobStore)
	assert.Equal(t, blobStore.Name, createdBlobStore.Name)
}

func TestGetBlobStores(t *testing.T) {
	router, service := setupBlobRouter()

	service.CreateBlobStore(&BlobStore{Name: "Test Blob Store 1", Type: "file"})
	service.CreateBlobStore(&BlobStore{Name: "Test Blob Store 2", Type: "file"})

	req, _ := http.NewRequest("GET", "/blob-stores", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var blobStores []BlobStore
	json.Unmarshal(w.Body.Bytes(), &blobStores)
	assert.Len(t, blobStores, 2)
}

func TestGetBlobStore(t *testing.T) {
	router, service := setupBlobRouter()

	blobStore := &BlobStore{Name: "Test Blob Store", Type: "file"}
	service.CreateBlobStore(blobStore)

	req, _ := http.NewRequest("GET", "/blob-stores/"+blobStore.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var foundBlobStore BlobStore
	json.Unmarshal(w.Body.Bytes(), &foundBlobStore)
	assert.Equal(t, blobStore.Name, foundBlobStore.Name)
}

func TestUpdateBlobStore(t *testing.T) {
	router, service := setupBlobRouter()

	blobStore := &BlobStore{Name: "Test Blob Store", Type: "file"}
	service.CreateBlobStore(blobStore)

	updatedBlobStore := BlobStore{Name: "Updated Blob Store", Type: "file"}
	jsonValue, _ := json.Marshal(updatedBlobStore)
	req, _ := http.NewRequest("PUT", "/blob-stores/"+blobStore.ID.String(), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result BlobStore
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, updatedBlobStore.Name, result.Name)
}

func TestDeleteBlobStore(t *testing.T) {
	router, service := setupBlobRouter()

	blobStore := &BlobStore{Name: "Test Blob Store", Type: "file"}
	service.CreateBlobStore(blobStore)

	req, _ := http.NewRequest("DELETE", "/blob-stores/"+blobStore.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
