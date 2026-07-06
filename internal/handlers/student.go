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

// GetAll возвращает список всех студентов.
func (h *StudentHandler) GetAll(c *gin.Context) {
	students, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

// Create создает нового студента на основе JSON-запроса.
func (h *StudentHandler) Create(c *gin.Context) {
	var student models.Student
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

// Update обновляет данные существующего студента по ID.
func (h *StudentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedStudent models.Student
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

// Delete удаляет студента по указанному ID.
func (h *StudentHandler) Delete(c *gin.Context) {
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

	c.JSON(http.StatusOK, student)
}

// FilterByGroup возвращает студентов, отфильтрованных по группе.
func (h *StudentHandler) FilterByGroup(c *gin.Context) {
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

// FilterByGroupOptional возвращает всех студентов, если параметр group не задан.
func (h *StudentHandler) FilterByGroupOptional(c *gin.Context) {
	group, exists := c.GetQuery("group")

	students, err := h.svc.FilterByGroupOptional(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

// GetPaginated возвращает постраничный список студентов.
func (h *StudentHandler) GetPaginated(c *gin.Context) {
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

// Search ищет студентов по группе и имени.
func (h *StudentHandler) Search(c *gin.Context) {
	group := c.Query("group")
	name := c.Query("name")

	students, err := h.svc.Search(group, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}
