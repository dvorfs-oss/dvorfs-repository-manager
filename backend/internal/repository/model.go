package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Repository struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;"`
	Name            string         `gorm:"type:varchar(255);unique_index"`
	Format          string         `gorm:"type:varchar(50)"`
	Type            string         `gorm:"type:varchar(50)"`
	Attributes      datatypes.JSON `gorm:"type:jsonb"`
	CleanupPolicyID *uuid.UUID     `gorm:"type:uuid"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Artifacts       []Artifact
	CleanupPolicy   *CleanupPolicy
}

type Artifact struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;"`
	RepositoryID     uuid.UUID      `gorm:"type:uuid"`
	Path             string         `gorm:"type:text"`
	Size             int64          `gorm:"type:bigint"`
	ContentType      string         `gorm:"type:varchar(255)"`
	Checksums        datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt        time.Time
	LastDownloadedAt *time.Time
}

type CleanupPolicy struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;"`
	Name      string         `gorm:"type:varchar(255)"`
	Criteria  datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
