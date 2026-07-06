package handlers

import "go-backend/internal/models"

type UserINterface interface {
	Create(name, email, password string) (models.User, error)
	GetById(id int) (models.User, error)
	List() ([]models.User, error)
	Update(id int, name, email, password string)
	Delete(id int) error
}
