package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;index"`
	User      User       `gorm:"foreignKey:UserID"`
	Token     string     `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}