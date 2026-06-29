package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-backend/internal/config"
	"go-backend/internal/handlers"
	"go-backend/internal/models"
	"go-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDBs(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.Student{}))
	config.DB = db
}

func TestStudentRepository_CreateGetUpdateDeleteStudent(t *testing.T) {
	setupTestDB(t)

	repo := &repository.StudentRepository{}
	student := &models.Student{
		Fio:         "Ivan Ivanov",
		Group:       "VM",
		PhoneNumber: "1234567890",
	}

	require.NoError(t, repo.Create(student))
	require.NotZero(t, student.ID)

	students, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, students, 1)
	require.Equal(t, "Ivan Ivanov", students[0].Fio)
	require.Equal(t, "VM", students[0].Group)

	fetched, err := repo.GetByID(student.ID)
	require.NoError(t, err)
	require.Equal(t, "Ivan Ivanov", fetched.Fio)

	fetched.Group = "CS"
	fetched.PhoneNumber = "0987654321"
	require.NoError(t, repo.Update(&fetched))

	updated, err := repo.GetByID(student.ID)
	require.NoError(t, err)
	require.Equal(t, "CS", updated.Group)
	require.Equal(t, "0987654321", updated.PhoneNumber)

	require.NoError(t, repo.Delete(student.ID))
	_, err = repo.GetByID(student.ID)
	require.Error(t, err)
}

func TestStudentRepository_FilterByGroupStydent(t *testing.T) {
	setupTestDB(t)

	repo := &repository.StudentRepository{}
	require.NoError(t, repo.Create(&models.Student{Fio: "Anna", Group: "VM", PhoneNumber: "111111"}))
	require.NoError(t, repo.Create(&models.Student{Fio: "Boris", Group: "CS", PhoneNumber: "222222"}))
	require.NoError(t, repo.Create(&models.Student{Fio: "Carl", Group: "VM", PhoneNumber: "333333"}))

	students, err := repo.FilterByGroup("VM")
	require.NoError(t, err)
	require.Len(t, students, 2)
	for _, student := range students {
		require.Equal(t, "VM", student.Group)
	}
}

func TestStudentHandler_GetAllStudents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	repo := &repository.StudentRepository{}
	require.NoError(t, repo.Create(&models.Student{Fio: "Ivan Ivanov", Group: "VM", PhoneNumber: "1234567890"}))

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	handler := handlers.NewStudentHandler()
	handler.GetAll(c)

	require.Equal(t, http.StatusOK, recorder.Code)

	var students []models.Student
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &students))
	require.Len(t, students, 1)
	require.Equal(t, "Ivan Ivanov", students[0].Fio)
}

func TestStudentHandler_CreateStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	payload := models.Student{Fio: "Olga Petrovna", Group: "CS", PhoneNumber: "9876543210"}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/students", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := handlers.NewStudentHandler()
	handler.Create(c)

	require.Equal(t, http.StatusCreated, recorder.Code)

	var created models.Student
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &created))
	require.Equal(t, payload.Fio, created.Fio)
	require.Equal(t, payload.Group, created.Group)
	require.Equal(t, payload.PhoneNumber, created.PhoneNumber)
	require.NotZero(t, created.ID)
}

func TestStudentHandler_GetByID_NotFoundReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	handler := handlers.NewStudentHandler()
	handler.GetByID(c)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}
