// internal/utils/jwt.go (полностью заменить)
package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Group  string `json:"group"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	accessSecret  string
	refreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func NewJWTManager() *JWTManager {
	accessTTL, _ := strconv.Atoi(getEnv("JWT_ACCESS_TTL", "15"))
	refreshTTL, _ := strconv.Atoi(getEnv("JWT_REFRESH_TTL", "43200"))

	return &JWTManager{
		accessSecret:  getEnv("JWT_ACCESS_SECRET", os.Getenv("JWT_SECRET")),
		refreshSecret: getEnv("JWT_REFRESH_SECRET", os.Getenv("JWT_SECRET")),
		AccessTTL:     time.Duration(accessTTL) * time.Minute,
		RefreshTTL:    time.Duration(refreshTTL) * time.Minute,
	}
}

func (m *JWTManager) GenerateAccessToken(userID uint, role, group string) (string, error) {
	claims := AccessClaims{
		UserID: userID,
		Role:   role,
		Group:  group,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        generateTokenID(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.accessSecret))
}

func (m *JWTManager) GenerateRefreshToken(userID uint) (string, error) {
	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.RefreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        generateTokenID(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.refreshSecret))
}

func (m *JWTManager) ValidateAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.accessSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token_expired")
		}
		return nil, errors.New("invalid_token")
	}
	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid_claims")
}

func (m *JWTManager) ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.refreshSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("refresh_token_expired")
		}
		return nil, errors.New("invalid_refresh_token")
	}
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid_refresh_claims")
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func generateTokenID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ОСТАВИТЬ СТАРЫЕ ФУНКЦИИ ДЛЯ СОВМЕСТИМОСТИ
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Group  string `json:"group"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role, group string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		Group:  group,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
