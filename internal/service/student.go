package service

import (
	"go-backend/internal/models"
	"go-backend/internal/repository"
	"strings"
)

// StudentService содержит бизнес-логику доступа пользователей и операций со студентами.
type StudentService struct {
	repo repository.UserRepositoryInterface
}

// NewStudentService создаёт новый сервис студентов.
func NewStudentService(repo repository.UserRepositoryInterface) *StudentService {
	if repo == nil {
		repo = &repository.UserRepository{}
	}
	return &StudentService{repo: repo}
}

func (s *StudentService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *StudentService) Create(student *models.User) error {
	return s.repo.Create(student)
}

func (s *StudentService) Update(id uint, updated models.User) (models.User, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return models.User{}, err
	}

	existing.Fio = updated.Fio
	existing.Group = updated.Group
	existing.PhoneNumber = updated.PhoneNumber

	if err := s.repo.Update(&existing); err != nil {
		return models.User{}, err
	}

	return existing, nil
}

func (s *StudentService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *StudentService) GetByID(id uint) (models.User, error) {
	return s.repo.GetByID(id)
}

func (s *StudentService) FilterByGroup(group string) ([]models.User, error) {
	return s.repo.FilterByGroup(group)
}

func (s *StudentService) FilterByGroupOptional(group string) ([]models.User, error) {
	if group == "" {
		return s.repo.GetAll()
	}
	return s.repo.FilterByGroup(group)
}

func (s *StudentService) GetPaginated(page, limit int) ([]models.User, int64, error) {
	return s.repo.GetPaginated(page, limit)
}

func (s *StudentService) Search(group, name string) ([]models.User, error) {
	return s.repo.Search(group, name)
}

// CanAccessUser проверяет, может ли текущий пользователь просматривать или редактировать целевого пользователя.
func (s *StudentService) CanAccessUser(currentUser, targetUser models.User) bool {
	if currentUser.Role == "admin" {
		return true
	}
	if currentUser.ID != 0 && targetUser.ID != 0 && currentUser.ID == targetUser.ID {
		return true
	}
	if currentUser.Role == "student" {
		return false
	}

	if currentUser.Role != "teacher" {
		return false
	}

	return hasAssignedGroup(currentUser.Group, targetUser.Group)
}

// FilterVisibleUsers оставляет только тех пользователей, которых преподаватель может видеть по закреплённым группам.
func (s *StudentService) FilterVisibleUsers(currentUser models.User, users []models.User) []models.User {
	if currentUser.Role == "admin" {
		return users
	}
	if currentUser.Role != "teacher" {
		return nil
	}

	filtered := make([]models.User, 0, len(users))
	for _, user := range users {
		if hasAssignedGroup(currentUser.Group, user.Group) {
			filtered = append(filtered, user)
		}
	}
	return filtered
}

func hasAssignedGroup(assignedGroups, targetGroup string) bool {
	for _, group := range strings.Split(assignedGroups, ",") {
		group = strings.TrimSpace(group)
		if group != "" && group == strings.TrimSpace(targetGroup) {
			return true
		}
	}
	return false
}
