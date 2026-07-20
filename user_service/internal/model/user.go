package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"uniqueIndex;not null" json:"username"`

	Email string `gorm:"uniqueIndex;not null" json:"email"`

	Password string `gorm:"not null" json:"-"`

	Role string `gorm:"default:user" json:"role"`

	Bio string `json:"bio"`

	AvatarURL string `json:"avatar_url"`
}
