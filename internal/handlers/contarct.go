package handlers

import "go-backend/internal/models"

type UserINterface interface {
	Create(name, email, password string) (models.Student, error)
	GetById(id int) (models.Student, error)
	List() ([]models.Student, error)
	Update(id int, name, email, password string)
	Delete(id int) error
}
