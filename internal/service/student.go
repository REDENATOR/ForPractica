package service

import (
	"go-backend/internal/models"
	"go-backend/internal/repository"
)

// StudentService содержит бизнес-логику для работы со студентами.
type StudentService struct {
	repo *repository.UserRepository
}

// NewStudentService создает новый сервис студентов.
func NewStudentService(repo *repository.UserRepository) *StudentService {
	if repo == nil {
		repo = &repository.UserRepository{}
	}
	return &StudentService{repo: repo}
}

func (s *StudentService) GetAll() ([]models.Student, error) {
	return s.repo.GetAll()
}

func (s *StudentService) Create(student *models.Student) error {
	return s.repo.Create(student)
}

func (s *StudentService) Update(id uint, updated models.Student) (models.Student, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return models.Student{}, err
	}

	existing.Fio = updated.Fio
	existing.Group = updated.Group
	existing.PhoneNumber = updated.PhoneNumber

	if err := s.repo.Update(&existing); err != nil {
		return models.Student{}, err
	}

	return existing, nil
}

func (s *StudentService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *StudentService) GetByID(id uint) (models.Student, error) {
	return s.repo.GetByID(id)
}

func (s *StudentService) FilterByGroup(group string) ([]models.Student, error) {
	return s.repo.FilterByGroup(group)
}

func (s *StudentService) FilterByGroupOptional(group string) ([]models.Student, error) {
	if group == "" {
		return s.repo.GetAll()
	}
	return s.repo.FilterByGroup(group)
}

func (s *StudentService) GetPaginated(page, limit int) ([]models.Student, int64, error) {
	return s.repo.GetPaginated(page, limit)
}

func (s *StudentService) Search(group, name string) ([]models.Student, error) {
	return s.repo.Search(group, name)
}
