package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlobService struct {
	db *gorm.DB
}

func NewBlobService(db *gorm.DB) *BlobService {
	return &BlobService{db: db}
}

func (s *BlobService) CreateBlobStore(blobStore *BlobStore) error {
	blobStore.ID = uuid.New()
	return s.db.Create(blobStore).Error
}

func (s *BlobService) GetBlobStores() ([]BlobStore, error) {
	var blobStores []BlobStore
	err := s.db.Find(&blobStores).Error
	return blobStores, err
}

func (s *BlobService) GetBlobStore(id uuid.UUID) (*BlobStore, error) {
	var blobStore BlobStore
	err := s.db.First(&blobStore, "id = ?", id).Error
	return &blobStore, err
}

func (s *BlobService) UpdateBlobStore(blobStore *BlobStore) error {
	return s.db.Save(blobStore).Error
}

func (s *BlobService) DeleteBlobStore(id uuid.UUID) error {
	return s.db.Delete(&BlobStore{}, "id = ?", id).Error
}
