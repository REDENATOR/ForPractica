package handlers

import (
	"go-backend/internal/models"
	"go-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TeacherHandler обрабатывает операции преподавателя, связанные со студентами.
type TeacherHandler struct {
	svc *service.StudentService
}

// NewTeacherHandler создаёт новый обработчик преподавателя.
func NewTeacherHandler(svc *service.StudentService) *TeacherHandler {
	if svc == nil {
		svc = service.NewStudentService(nil)
	}
	return &TeacherHandler{svc: svc}
}

func (h *TeacherHandler) authorize(c *gin.Context) bool {
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

func (h *TeacherHandler) GetAll(c *gin.Context) {
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

func (h *TeacherHandler) Create(c *gin.Context) {
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

func (h *TeacherHandler) Delete(c *gin.Context) {
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

func (h *TeacherHandler) FilterByGroup(c *gin.Context) {
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

func (h *TeacherHandler) FilterByGroupOptional(c *gin.Context) {
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

func (h *TeacherHandler) GetPaginated(c *gin.Context) {
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

func (h *TeacherHandler) Search(c *gin.Context) {
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
