package repository

import (
	"go-backend/internal/config"
	"go-backend/internal/models"
)

// StudentRepository предоставляет методы доступа к данным студентов.
type StudentRepository struct{}

// GetAll возвращает всех студентов из базы данных.
func (r *StudentRepository) GetAll() ([]models.Student, error) {
	var students []models.Student
	result := config.DB.Find(&students)
	return students, result.Error
}

// GetByID возвращает студента по его идентификатору.
func (r *StudentRepository) GetByID(id uint) (models.Student, error) {
	var student models.Student
	result := config.DB.First(&student, id)
	return student, result.Error
}

// Create добавляет нового студента в базу данных.
func (r *StudentRepository) Create(student *models.Student) error {
	return config.DB.Create(student).Error
}

// Update сохраняет изменения студента в базе данных.
func (r *StudentRepository) Update(student *models.Student) error {
	return config.DB.Save(student).Error
}

// Delete удаляет студента по идентификатору.
func (r *StudentRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Student{}, id).Error
}

// FilterByGroup возвращает студентов, относящихся к указанной группе.
func (r *StudentRepository) FilterByGroup(group string) ([]models.Student, error) {
	var students []models.Student
	result := config.DB.Where("group_of_students = ?", group).Find(&students)
	return students, result.Error
}
