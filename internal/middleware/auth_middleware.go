package middleware

import (
	"net/http"
	"strings"

	"backend_go/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header must be Bearer token", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Check Blocklist
		var count int64
		err := db.Table("jwt_blocklists").Where("token = ?", tokenString).Count(&count).Error
		if err == nil && count > 0 {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token has been revoked", nil)
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		}

		c.Set("token", tokenString)
		// Set user info to context
		c.Set("userID", claims["sub"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		userRole := roleValue.(string)
		roleMatched := false
		for _, role := range roles {
			if role == userRole {
				roleMatched = true
				break
			}
		}

		if !roleMatched {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden: You don't have access to this resource", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
