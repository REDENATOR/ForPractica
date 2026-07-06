package models

import "gorm.io/gorm"

// Student представляет запись студента в базе данных.
type User struct {
	gorm.Model

	// Fio содержит полное имя студента.
	Fio string `json:"fio"`

	// Group соответствует группе студентов и сохраняется в колонку group_of_students.
	Group string `json:"group" gorm:"column:group_of_students"`

	// PhoneNumber содержит контактный номер студента.
	PhoneNumber string `json:"phoneNumber"`

	Password string `json:"password"`
	// Не каждый пользователь студент, но каждый студент пользователь, поэтому добавим поле Role для различения ролей пользователей.
	Role string `gorm:"type:varchar(20);not null;default:'student'"`
}
