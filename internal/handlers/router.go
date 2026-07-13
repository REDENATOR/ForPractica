package handlers

import (
	"errors"
	"go-backend/internal/middlevare"
	"go-backend/internal/repository"
	"go-backend/internal/service"
	"go-backend/internal/utils" // ← ДОБАВИТЬ

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	engine *gin.Engine
	db     *gorm.DB
}

func NewRouter(engine *gin.Engine, db *gorm.DB) (*Router, error) {
	if engine == nil {
		return nil, errors.New("engine is nil")
	}
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &Router{engine: engine, db: db}, nil
}

func (r *Router) RegisterRoutes() {
	// ============================================
	// 1. ПУБЛИЧНЫЕ МАРШРУТЫ (без токена)
	// ============================================
	// ИНИЦИАЛИЗИРУЕМ JWTManager
	jwtManager := utils.NewJWTManager() // ← ДОБАВИТЬ

	// Передаем jwtManager в AuthHandler
	authHandler := NewAuthHandler(r.db, jwtManager) // ← ИЗМЕНИТЬ
	auth := r.engine.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh) // ← ДОБАВИТЬ
	}

	// ============================================
	// 2. ЗАЩИЩЁННЫЕ МАРШРУТЫ (с токеном)
	// ============================================
	api := r.engine.Group("/api")
	api.Use(middlevare.AuthMiddleware(jwtManager)) // ← ИЗМЕНИТЬ (передаем jwtManager)
	{
		repo := &repository.UserRepository{}
		studentService := service.NewStudentService(repo)
		studentHandler := NewStudentHandler(studentService)
		adminHandler := NewAdminHandler(studentService)

		// --- Роуты для работы со студентами ---
		api.GET("/students", middlevare.RoleMiddleware("admin", "teacher"), adminHandler.GetAll)
		api.POST("/students", middlevare.RoleMiddleware("admin"), adminHandler.Create)
		api.DELETE("/students/:id", middlevare.RoleMiddleware("admin"), adminHandler.Delete)
		api.GET("/students/filter", middlevare.RoleMiddleware("admin", "teacher"), adminHandler.FilterByGroup)
		api.GET("/students/filter-optional", middlevare.RoleMiddleware("admin", "teacher"), adminHandler.FilterByGroupOptional)
		api.GET("/students/paginated", middlevare.RoleMiddleware("admin", "teacher"), adminHandler.GetPaginated)
		api.GET("/students/search", middlevare.RoleMiddleware("admin", "teacher"), adminHandler.Search)

		api.GET("/students/:id", studentHandler.GetByID)
		api.PUT("/students/:id", studentHandler.Update)

		// --- ДОБАВИТЬ РОУТЫ ДЛЯ LOGOUT ---
		api.POST("/logout", authHandler.Logout)               // ← ДОБАВИТЬ
		api.POST("/logout-all", authHandler.LogoutAllDevices) // ← ДОБАВИТЬ
	}
}
