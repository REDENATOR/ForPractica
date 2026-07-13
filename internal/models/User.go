package models

import (
	"time"

	"gorm.io/gorm"
)

// User представляет пользователя системы, например студента или преподавателя.
type User struct {
	gorm.Model

	Fio string `json:"fio"`

	// Group stores a single group for students and a comma-separated list of assigned groups for teachers.
	Group string `json:"group" gorm:"column:group_of_students"`

	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Role        string `gorm:"type:varchar(20);not null;default:'student'"`
}
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"unique;not null;size:64"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Revoked   bool      `gorm:"default:false;index"`
	User      User      `gorm:"foreignKey:UserID"`
}
