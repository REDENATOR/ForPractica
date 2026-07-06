package handlers

import (
	"go-backend/internal/models"
	"go-backend/internal/repository"
	"go-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *StudentHandler) authorizeSelfOrAdmin(c *gin.Context, targetUserID uint) bool {
	currentUserID, ok := getCurrentUserID(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}

	roleValue, exists := c.Get("role")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}

	roleStr, _ := roleValue.(string)
	if roleStr == "admin" || currentUserID == targetUserID {
		return true
	}

	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
	return false
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

// GetByID returns a student by ID.
func (h *StudentHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID должен быть числом",
		})
		return
	}

	if !h.authorizeSelfOrAdmin(c, uint(id)) {
		return
	}

	student, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Студент не найден",
		})
		return
	}

	c.JSON(http.StatusOK, student)
}

// Update updates an existing student by ID.
func (h *StudentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if !h.authorizeSelfOrAdmin(c, uint(id)) {
		return
	}

	var updatedStudent models.User
	if err := c.ShouldBindJSON(&updatedStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := h.svc.Update(uint(id), updatedStudent)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}
