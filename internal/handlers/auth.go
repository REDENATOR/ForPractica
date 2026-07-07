package handlers

import (
	"go-backend/internal/models"
	"go-backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler — конструктор, создаёт новый AuthHandler
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Register — регистрация нового пользователя по номеру телефона
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Fio         string `json:"fio" binding:"required"`
		PhoneNumber string `json:"phoneNumber" binding:"required"`
		Password    string `json:"password" binding:"required,min=6"`
	}

	// 1. Проверяем, что прислали JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Проверяем, есть ли уже пользователь с таким номером
	var existingUser models.User
	if err := h.db.Where("phone_number = ?", req.PhoneNumber).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	// 3. Хешируем пароль
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// 4. Создаём пользователя (роль student по умолчанию)
	user := models.User{
		Fio:         req.Fio,
		Group:       "",
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
		Role:        "student",
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// 5. Возвращаем ответ
	c.JSON(http.StatusCreated, gin.H{
		"message":     "user created",
		"user_id":     user.ID,
		"role":        user.Role,
		"phoneNumber": user.PhoneNumber,
	})
}

// Login — вход пользователя
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phoneNumber" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}

	// 1. Проверяем JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Ищем пользователя по номеру телефона
	var user models.User
	if err := h.db.Where("phone_number = ?", req.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 3. Проверяем пароль
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 4. Генерируем JWT
	token, err := utils.GenerateToken(user.ID, user.Role, user.Group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// 5. Возвращаем токен
	c.JSON(http.StatusOK, gin.H{
		"token":       token,
		"user_id":     user.ID,
		"role":        user.Role,
		"phoneNumber": user.PhoneNumber,
	})
}
