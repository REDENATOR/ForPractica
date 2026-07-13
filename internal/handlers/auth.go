// internal/handlers/auth.go
package handlers

import (
	"go-backend/internal/models"
	"go-backend/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db         *gorm.DB
	jwtManager *utils.JWTManager // ← ДОБАВИТЬ
}

// NewAuthHandler — ОБНОВИТЬ КОНСТРУКТОР
func NewAuthHandler(db *gorm.DB, jwtManager *utils.JWTManager) *AuthHandler {
	return &AuthHandler{
		db:         db,
		jwtManager: jwtManager,
	}
}

// Register — без изменений
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Fio         string `json:"fio" binding:"required"`
		PhoneNumber string `json:"phoneNumber" binding:"required"`
		Password    string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := h.db.Where("phone_number = ?", req.PhoneNumber).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

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

	c.JSON(http.StatusCreated, gin.H{
		"message":     "user created",
		"user_id":     user.ID,
		"role":        user.Role,
		"phoneNumber": user.PhoneNumber,
	})
}

// Login — ИЗМЕНИТЬ (выдавать пару токенов)
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phoneNumber" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("phone_number = ?", req.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Генерируем пару токенов
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Role, user.Group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// Сохраняем Refresh Token в БД
	refreshHash := utils.HashToken(refreshToken)
	rt := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(h.jwtManager.RefreshTTL),
		Revoked:   false,
	}

	if err := h.db.Create(&rt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_id":       user.ID,
		"role":          user.Role,
		"phoneNumber":   user.PhoneNumber,
		"expires_in":    int(h.jwtManager.AccessTTL.Seconds()),
	})
}

// Refresh — НОВЫЙ МЕТОД
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
		return
	}

	// 1. Валидируем Refresh Token
	claims, err := h.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 2. Ищем в БД
	refreshHash := utils.HashToken(req.RefreshToken)
	var storedRT models.RefreshToken
	if err := h.db.Where("token_hash = ? AND user_id = ?", refreshHash, claims.UserID).
		First(&storedRT).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	// 3. Проверяем статус
	if storedRT.Revoked {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token revoked"})
		return
	}

	if time.Now().After(storedRT.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
		return
	}

	// 4. РОТАЦИЯ: удаляем старый RT
	if err := h.db.Delete(&storedRT).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rotate token"})
		return
	}

	// 5. Получаем пользователя
	var user models.User
	if err := h.db.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	// 6. Генерируем новую пару
	newAccessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Role, user.Group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	newRefreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// 7. Сохраняем новый RT
	newRefreshHash := utils.HashToken(newRefreshToken)
	newRT := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: newRefreshHash,
		ExpiresAt: time.Now().Add(h.jwtManager.RefreshTTL),
		Revoked:   false,
	}

	if err := h.db.Create(&newRT).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save new refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
		"expires_in":    int(h.jwtManager.AccessTTL.Seconds()),
	})
}

// Logout — НОВЫЙ МЕТОД
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	c.ShouldBindJSON(&req)

	if req.RefreshToken != "" {
		// Удаляем конкретный RT
		refreshHash := utils.HashToken(req.RefreshToken)
		if err := h.db.Where("user_id = ? AND token_hash = ?", userID, refreshHash).
			Delete(&models.RefreshToken{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
			return
		}
	} else {
		// Удаляем все RT пользователя
		if err := h.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// LogoutAllDevices — НОВЫЙ МЕТОД
func (h *AuthHandler) LogoutAllDevices(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Помечаем все RT как revoked
	if err := h.db.Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout all devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out from all devices"})
}
