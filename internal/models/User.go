package models

import "gorm.io/gorm"

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
