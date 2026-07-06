package middlevare

import (
    "net/http"
    "strings"
	"go-backend/internal/utils"
    "github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT и кладёт user_id и role в контекст
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Забираем заголовок Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
            return
        }

        // 2. Проверяем формат "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
            return
        }

        tokenString := parts[1]

        // 3. Проверяем JWT
        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
            return
        }

        // 4. Кладём данные в контекст
        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)

        // 5. Передаём управление дальше
        c.Next()
    }
}

// RoleMiddleware проверяет, есть ли у пользователя нужная роль
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Достаём роль из контекста (установлена в AuthMiddleware)
        role, exists := c.Get("role")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }

        roleStr := role.(string)

        // 2. Проверяем, есть ли роль в списке разрешённых
        for _, allowed := range allowedRoles {
            if roleStr == allowed {
                c.Next()
                return
            }
        }

        // 3. Если нет — 403 Forbidden
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
    }
}