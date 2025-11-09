package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yayasan/erp-backend/internal/utils"
)

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		tokenString, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid authorization header")
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			if err == utils.ErrExpiredToken {
				utils.UnauthorizedResponse(c, "Token has expired")
			} else {
				utils.UnauthorizedResponse(c, "Invalid token")
			}
			c.Abort()
			return
		}

		// Check token type
		if claims.TokenType != "access" {
			utils.UnauthorizedResponse(c, "Invalid token type")
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("branch_id", claims.BranchID)
		c.Set("is_super_admin", claims.IsSuperAdmin)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT token but doesn't require it
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Store claims in context if valid
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("branch_id", claims.BranchID)
		c.Set("is_super_admin", claims.IsSuperAdmin)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}
