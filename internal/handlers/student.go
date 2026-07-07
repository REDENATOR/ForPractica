package handlers

import (
	"go-backend/internal/models"
	"go-backend/internal/repository"
	"go-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StudentHandler обрабатывает запросы, связанные со студентами.
type StudentHandler struct {
	svc *service.StudentService
}

// NewStudentHandler возвращает новый экземпляр StudentHandler.
func NewStudentHandler(svc *service.StudentService) *StudentHandler {
	if svc == nil {
		svc = service.NewStudentService(&repository.UserRepository{})
	}
	return &StudentHandler{svc: svc}
}

func (h *StudentHandler) authorizeAccess(c *gin.Context, targetUser models.User) bool {
	currentUser, ok := getCurrentUser(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}

	if h.svc.CanAccessUser(currentUser, targetUser) {
		return true
	}

	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
	return false
}

func getCurrentUser(c *gin.Context) (models.User, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return models.User{}, false
	}

	userID := uint(0)
	switch v := userIDValue.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case int64:
		userID = uint(v)
	case float64:
		userID = uint(v)
	default:
		return models.User{}, false
	}

	roleValue, exists := c.Get("role")
	if !exists {
		return models.User{}, false
	}

	roleStr, _ := roleValue.(string)
	groupValue, _ := c.Get("group")
	groupStr, _ := groupValue.(string)

	return models.User{Model: gorm.Model{ID: userID}, Role: roleStr, Group: groupStr}, true
}

func getCurrentUserID(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	switch v := userIDValue.(type) {
	case uint:
		return v, true
	case int:
		return uint(v), true
	case int64:
		return uint(v), true
	case float64:
		return uint(v), true
	default:
		return 0, false
	}
}

// GetByID возвращает студента по ID.
func (h *StudentHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID должен быть числом",
		})
		return
	}

	student, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Студент не найден",
		})
		return
	}

	if !h.authorizeAccess(c, student) {
		return
	}

	c.JSON(http.StatusOK, student)
}

// Update обновляет существующего студента по ID.
func (h *StudentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	student, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	if !h.authorizeAccess(c, student) {
		return
	}

	var updatedStudent models.User
	if err := c.ShouldBindJSON(&updatedStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err = h.svc.Update(uint(id), updatedStudent)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}
