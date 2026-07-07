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

// GetAll возвращает список всех студентов (только для администратора).
func (h *AdminHandler) GetAll(c *gin.Context) {
	if !h.authorize(c) {
		return
	}
	currentUser, ok := getCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	students, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if currentUser.Role == "teacher" {
		students = h.svc.FilterVisibleUsers(currentUser, students)
	}
	c.JSON(http.StatusOK, students)
}

// Create создаёт нового студента из JSON-данных (только для администратора).
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

// Delete удаляет студента по ID (только для администратора).
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

// FilterByGroup возвращает студентов, отфильтрованных по группе (только для администратора).
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

	currentUser, ok := getCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	students, err := h.svc.FilterByGroup(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if currentUser.Role == "teacher" {
		students = h.svc.FilterVisibleUsers(currentUser, students)
	}

	c.JSON(http.StatusOK, students)
}

// FilterByGroupOptional возвращает всех студентов, если группа не указана (только для администратора).
func (h *AdminHandler) FilterByGroupOptional(c *gin.Context) {
	if !h.authorize(c) {
		return
	}

	group, _ := c.GetQuery("group")

	currentUser, ok := getCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	students, err := h.svc.FilterByGroupOptional(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if currentUser.Role == "teacher" {
		students = h.svc.FilterVisibleUsers(currentUser, students)
	}

	c.JSON(http.StatusOK, students)
}

// GetPaginated возвращает постраничный список студентов (только для администратора).
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

	currentUser, ok := getCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	students, total, err := h.svc.GetPaginated(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if currentUser.Role == "teacher" {
		students = h.svc.FilterVisibleUsers(currentUser, students)
		total = int64(len(students))
		if page*limit < len(students) {
			end := (page + 1) * limit
			if end > len(students) {
				end = len(students)
			}
			students = students[page*limit : end]
		} else {
			students = []models.User{}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  students,
	})
}

// Search ищет студентов по группе и имени (только для администратора).
func (h *AdminHandler) Search(c *gin.Context) {
	if !h.authorize(c) {
		return
	}

	group := c.Query("group")
	name := c.Query("name")

	currentUser, ok := getCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	students, err := h.svc.Search(group, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if currentUser.Role == "teacher" {
		students = h.svc.FilterVisibleUsers(currentUser, students)
	}

	c.JSON(http.StatusOK, students)
}
