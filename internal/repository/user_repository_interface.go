package repository

import "go-backend/internal/models"

// UserRepositoryInterface определяет контракт для доступа к пользователям.
type UserRepositoryInterface interface {
	GetAll() ([]models.User, error)
	GetByID(id uint) (models.User, error)
	Create(student *models.User) error
	Update(student *models.User) error
	Delete(id uint) error
	FilterByGroup(group string) ([]models.User, error)
	GetPaginated(page, limit int) ([]models.User, int64, error)
	Search(group, name string) ([]models.User, error)
}

var _ UserRepositoryInterface = (*UserRepository)(nil)
