package repository

import (
	"go-backend/internal/config"
	"go-backend/internal/models"
)

// UserRepository предоставляет методы доступа к данным пользователей.
type UserRepository struct{}

// GetAll возвращает всех пользователей из базы данных.
func (r *UserRepository) GetAll() ([]models.Student, error) {
	var students []models.Student
	result := config.DB.Find(&students)
	return students, result.Error
}

// GetByID возвращает студента по его идентификатору.
func (r *UserRepository) GetByID(id uint) (models.Student, error) {
	var student models.Student
	result := config.DB.First(&student, id)
	return student, result.Error
}

// Create добавляет нового студента в базу данных.
func (r *UserRepository) Create(student *models.Student) error {
	return config.DB.Create(student).Error
}

// Update сохраняет изменения студента в базе данных.
func (r *UserRepository) Update(student *models.Student) error {
	return config.DB.Save(student).Error
}

// Delete удаляет студента по идентификатору.
func (r *UserRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Student{}, id).Error
}

// FilterByGroup возвращает студентов, относящихся к указанной группе.
func (r *UserRepository) FilterByGroup(group string) ([]models.Student, error) {
	var students []models.Student
	result := config.DB.Where("group_of_students = ?", group).Find(&students)
	return students, result.Error
}
