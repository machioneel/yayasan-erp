package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/utils"
)

// RequirePermission middleware checks if user has required permission
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is super admin (bypass permission check)
		isSuperAdmin, exists := c.Get("is_super_admin")
		if exists && isSuperAdmin.(bool) {
			c.Next()
			return
		}

		// Get user permissions from context
		permissionsInterface, exists := c.Get("permissions")
		if !exists {
			utils.ForbiddenResponse(c, "No permissions found")
			c.Abort()
			return
		}

		permissions, ok := permissionsInterface.([]string)
		if !ok {
			utils.ForbiddenResponse(c, "Invalid permissions format")
			c.Abort()
			return
		}

		// Check if user has the required permission
		hasPermission := false
		for _, p := range permissions {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.ForbiddenResponse(c, "You don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission checks if user has at least one of the required permissions
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is super admin (bypass permission check)
		isSuperAdmin, exists := c.Get("is_super_admin")
		if exists && isSuperAdmin.(bool) {
			c.Next()
			return
		}

		// Get user permissions from context
		userPermissionsInterface, exists := c.Get("permissions")
		if !exists {
			utils.ForbiddenResponse(c, "No permissions found")
			c.Abort()
			return
		}

		userPermissions, ok := userPermissionsInterface.([]string)
		if !ok {
			utils.ForbiddenResponse(c, "Invalid permissions format")
			c.Abort()
			return
		}

		// Check if user has any of the required permissions
		hasPermission := false
		for _, requiredPerm := range permissions {
			for _, userPerm := range userPermissions {
				if userPerm == requiredPerm {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			utils.ForbiddenResponse(c, "You don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSuperAdmin middleware checks if user is super admin
func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isSuperAdmin, exists := c.Get("is_super_admin")
		if !exists || !isSuperAdmin.(bool) {
			utils.ForbiddenResponse(c, "Super admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireBranch middleware checks if user belongs to specified branch
func RequireBranch(branchID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Super admin can access all branches
		isSuperAdmin, exists := c.Get("is_super_admin")
		if exists && isSuperAdmin.(bool) {
			c.Next()
			return
		}

		// Get user's branch ID
		userBranchID, exists := c.Get("branch_id")
		if !exists {
			utils.ForbiddenResponse(c, "Branch information not found")
			c.Abort()
			return
		}

		userBranchUUID, ok := userBranchID.(*uuid.UUID)
		if !ok || userBranchUUID == nil {
			utils.ForbiddenResponse(c, "Invalid branch information")
			c.Abort()
			return
		}

		// Check if user's branch matches required branch
		if *userBranchUUID != branchID {
			utils.ForbiddenResponse(c, "You don't have access to this branch")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOwnResource middleware checks if user is accessing their own resource
func RequireOwnResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Super admin can access all resources
		isSuperAdmin, exists := c.Get("is_super_admin")
		if exists && isSuperAdmin.(bool) {
			c.Next()
			return
		}

		// Get authenticated user ID
		authUserIDInterface, exists := c.Get("user_id")
		if !exists {
			utils.UnauthorizedResponse(c, "User not authenticated")
			c.Abort()
			return
		}

		authUserID, ok := authUserIDInterface.(uuid.UUID)
		if !ok {
			utils.UnauthorizedResponse(c, "Invalid user ID")
			c.Abort()
			return
		}

		// Get resource user ID from URL parameter
		resourceUserIDStr := c.Param("id")
		resourceUserID, err := uuid.Parse(resourceUserIDStr)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid user ID format")
			c.Abort()
			return
		}

		// Check if authenticated user matches resource user
		if authUserID != resourceUserID {
			utils.ForbiddenResponse(c, "You can only access your own resources")
			c.Abort()
			return
		}

		c.Next()
	}
}
