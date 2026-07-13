package service

import (
	"errors"
	"testing"

	"go-backend/internal/models"
	"go-backend/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetByIDCallsRepositoryAndReturnsUser(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryInterface(t)
	student := models.User{Model: gorm.Model{ID: 42}, Fio: "Ivan Ivanov"}
	mockRepo.On("GetByID", uint(42)).Return(student, nil)

	svc := NewStudentService(mockRepo)
	result, err := svc.GetByID(42)

	assert.NoError(t, err)
	assert.Equal(t, student, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCallsRepositoryAndReturnsUpdatedUser(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryInterface(t)
	existing := models.User{Model: gorm.Model{ID: 10}, Fio: "Old", Group: "A", PhoneNumber: "123"}
	updated := models.User{Fio: "New", Group: "B", PhoneNumber: "456"}
	mockRepo.On("GetByID", uint(10)).Return(existing, nil)
	mockRepo.On("Update", mock.MatchedBy(func(user *models.User) bool {
		return user.ID == 10 && user.Fio == "New" && user.Group == "B" && user.PhoneNumber == "456"
	})).Return(nil)

	svc := NewStudentService(mockRepo)
	result, err := svc.Update(10, updated)

	assert.NoError(t, err)
	assert.Equal(t, uint(10), result.ID)
	assert.Equal(t, "New", result.Fio)
	assert.Equal(t, "B", result.Group)
	assert.Equal(t, "456", result.PhoneNumber)
	mockRepo.AssertExpectations(t)
}

func TestFilterByGroupOptionalEmptyGroupReturnsAll(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryInterface(t)
	users := []models.User{{Fio: "Ivan"}, {Fio: "Petr"}}
	mockRepo.On("GetAll").Return(users, nil)

	svc := NewStudentService(mockRepo)
	result, err := svc.FilterByGroupOptional("")

	assert.NoError(t, err)
	assert.Equal(t, users, result)
	mockRepo.AssertExpectations(t)
}

func TestFilterByGroupOptionalWithGroupReturnsFiltered(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryInterface(t)
	users := []models.User{{Fio: "Ivan", Group: "VM"}}
	mockRepo.On("FilterByGroup", "VM").Return(users, nil)

	svc := NewStudentService(mockRepo)
	result, err := svc.FilterByGroupOptional("VM")

	assert.NoError(t, err)
	assert.Equal(t, users, result)
	mockRepo.AssertExpectations(t)
}

func TestGetPaginatedCallsRepositoryAndReturnsPage(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryInterface(t)
	users := []models.User{{Fio: "Ivan"}}
	mockRepo.On("GetPaginated", 1, 10).Return(users, int64(1), nil)

	svc := NewStudentService(mockRepo)
	result, total, err := svc.GetPaginated(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, users, result)
	assert.Equal(t, int64(1), total)
	mockRepo.AssertExpectations(t)
}

func TestSearchCallsRepositoryAndReturnsResults(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryInterface(t)
	users := []models.User{{Fio: "Ivan", Group: "VM"}}
	mockRepo.On("Search", "VM", "Ivan").Return(users, nil)

	svc := NewStudentService(mockRepo)
	result, err := svc.Search("VM", "Ivan")

	assert.NoError(t, err)
	assert.Equal(t, users, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateCallsRepositoryAndReturnsErrorFromRepo(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryInterface(t)
	student := &models.User{Fio: "New"}
	errorExpected := errors.New("create failure")
	mockRepo.On("Create", student).Return(errorExpected)

	svc := NewStudentService(mockRepo)
	err := svc.Create(student)

	assert.ErrorIs(t, err, errorExpected)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCallsRepositoryAndReturnsErrorFromRepo(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryInterface(t)
	errorExpected := errors.New("delete failure")
	mockRepo.On("Delete", uint(5)).Return(errorExpected)

	svc := NewStudentService(mockRepo)
	err := svc.Delete(5)

	assert.ErrorIs(t, err, errorExpected)
	mockRepo.AssertExpectations(t)
}
