package handlers

import (
	"go-backend/internal/models"
	"go-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminHandler обрабатывает админские операции со студентами.
type AdminHandler struct {
	svc *service.StudentService
}

// NewAdminHandler возвращает новый экземпляр AdminHandler.
func NewAdminHandler(svc *service.StudentService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) authorize(c *gin.Context) bool {
	roleValue, exists := c.Get("role")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}

	roleStr, ok := roleValue.(string)
	if !ok || roleStr != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: admin privileges required"})
		return false
	}
	return true
}

// GetAll returns list of all students (admin only).
func (h *AdminHandler) GetAll(c *gin.Context) {
	if !h.authorize(c) {
		return
	}
	students, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

// Create creates a new student from JSON payload (admin only).
func (h *AdminHandler) Create(c *gin.Context) {
	if !h.authorize(c) {
		return
	}

	var student models.User
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Create(&student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, student)
}

// Delete removes a student by ID (admin only).
func (h *AdminHandler) Delete(c *gin.Context) {
	if !h.authorize(c) {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted"})
}

// FilterByGroup returns students filtered by group (admin only).
func (h *AdminHandler) FilterByGroup(c *gin.Context) {
	if !h.authorize(c) {
		return
	}

	group := c.Query("group")

	if group == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Параметр group обязателен. Пример: ?group=VM",
		})
		return
	}

	students, err := h.svc.FilterByGroup(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

// FilterByGroupOptional returns all students when group is not provided (admin only).
func (h *AdminHandler) FilterByGroupOptional(c *gin.Context) {
	if !h.authorize(c) {
		return
	}

	group, _ := c.GetQuery("group")

	students, err := h.svc.FilterByGroupOptional(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

// GetPaginated returns a paginated student list (admin only).
func (h *AdminHandler) GetPaginated(c *gin.Context) {
	if !h.authorize(c) {
		return
	}

	pageStr := c.DefaultQuery("page", "0")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 0 {
		page = 0
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	students, total, err := h.svc.GetPaginated(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  students,
	})
}

// Search looks up students by group and name (admin only).
func (h *AdminHandler) Search(c *gin.Context) {
	if !h.authorize(c) {
		return
	}

	group := c.Query("group")
	name := c.Query("name")

	students, err := h.svc.Search(group, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}
