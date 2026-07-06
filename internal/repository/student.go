package repository

import (
	"go-backend/internal/config"
	"go-backend/internal/models"
)

// UserRepository предоставляет методы доступа к данным пользователей.
type UserRepository struct{}

// GetAll возвращает всех пользователей из базе данных.
func (r *UserRepository) GetAll() ([]models.User, error) {
	var students []models.User
	result := config.DB.Find(&students)
	return students, result.Error
}

// GetByID возвращает студента по его идентификатору.
func (r *UserRepository) GetByID(id uint) (models.User, error) {
	var student models.User
	result := config.DB.First(&student, id)
	return student, result.Error
}

// Create добавляет нового студента в базу данных.
func (r *UserRepository) Create(student *models.User) error {
	return config.DB.Create(student).Error
}

// Update сохраняет изменения студента в базе данных.
func (r *UserRepository) Update(student *models.User) error {
	return config.DB.Save(student).Error
}

// Delete удаляет студента по идентификатору.
func (r *UserRepository) Delete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}

// FilterByGroup возвращает студентов, относящихся к указанной группе.
func (r *UserRepository) FilterByGroup(group string) ([]models.User, error) {
	var students []models.User
	result := config.DB.Where("group_of_students = ?", group).Find(&students)
	return students, result.Error
}

func (r *UserRepository) GetPaginated(page, limit int) ([]models.User, int64, error) {
	var students []models.User
	var total int64
	if err := config.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := page * limit
	result := config.DB.Offset(offset).Limit(limit).Find(&students)
	return students, total, result.Error
}

func (r *UserRepository) Search(group, name string) ([]models.User, error) {
	var students []models.User
	query := config.DB

	if group != "" {
		query = query.Where("group_of_students = ?", group)
	}

	if name != "" {
		query = query.Where("fio LIKE ?", "%"+name+"%")
	}

	result := query.Find(&students)
	return students, result.Error
}
