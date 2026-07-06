package handlers

import (
	"errors"

	"go-backend/internal/repository"
	"go-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(engine *gin.Engine) (*Router, error) {
	if engine == nil {
		return nil, errors.New("engine is nil")
	}
	return &Router{engine: engine}, nil
}

func (r *Router) RegisterRoutes() {
	api := r.engine.Group("/api")
	{
		repo := &repository.UserRepository{}
		studentService := service.NewStudentService(repo)
		studentHandler := NewStudentHandler(studentService)

		// Роуты для работы со студентами.
		api.GET("/students", studentHandler.GetAll)                                // Получить всех студентов.
		api.POST("/students", studentHandler.Create)                               // Создать нового студента.
		api.PUT("/students/:id", studentHandler.Update)                            // Обновить данные студента по ID.
		api.DELETE("/students/:id", studentHandler.Delete)                         // Удалить студента по ID.
		api.GET("/students/:id", studentHandler.GetByID)                           // Получить студента по ID.
		api.GET("/students/filter", studentHandler.FilterByGroup)                  // Фильтрация студентов по группе.
		api.GET("/students/filter-optional", studentHandler.FilterByGroupOptional) // Фильтрация с необязательным параметром.
		api.GET("/students/paginated", studentHandler.GetPaginated)                // Получить студентов постранично.
		api.GET("/students/search", studentHandler.Search)                         // Поиск студентов.
	}
}
