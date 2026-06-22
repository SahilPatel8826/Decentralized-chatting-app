package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	Username string `gorm:"uniqueIndex;not null" json:"username"`

	Email string `gorm:"uniqueIndex;not null" json:"email"`

	Password string `gorm:"not null" json:"-"`

	Role string `gorm:"default:user" json:"role"`

	Bio string `json:"bio"`

	AvatarURL string `json:"avatar_url"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
