package models

import "gorm.io/gorm"

// Student представляет запись студента в базе данных.
type Student struct {
	gorm.Model

	// Fio содержит полное имя студента.
	Fio string `json:"fio"`

	// Group соответствует группе студентов и сохраняется в колонку group_of_students.
	Group string `json:"group" gorm:"column:group_of_students"`

	// PhoneNumber содержит контактный номер студента.
	PhoneNumber string `json:"phoneNumber"`
}

// TableName возвращает имя таблицы для модели Student.
func (Student) TableName() string {
	return "students"
}
