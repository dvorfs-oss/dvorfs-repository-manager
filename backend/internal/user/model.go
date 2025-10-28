package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	Username     string    `gorm:"type:varchar(255);unique_index"`
	PasswordHash string    `gorm:"type:varchar(255)"`
	Email        string    `gorm:"type:varchar(255);unique_index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Roles        []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;"`
	Name       string         `gorm:"type:varchar(255);unique_index"`
	Privileges datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt  time.Time
}
