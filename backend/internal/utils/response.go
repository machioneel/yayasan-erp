package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yayasan/erp-backend/internal/models"
)

// SuccessResponse sends success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, models.Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, models.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	})
}

// ValidationErrorResponse sends validation error response
func ValidationErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, models.ErrorResponse{
		Success: false,
		Error:   "Validation failed: " + err.Error(),
		Code:    http.StatusBadRequest,
	})
}

// UnauthorizedResponse sends unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	c.JSON(http.StatusUnauthorized, models.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusUnauthorized,
	})
}

// ForbiddenResponse sends forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Access forbidden"
	}
	c.JSON(http.StatusForbidden, models.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusForbidden,
	})
}

// NotFoundResponse sends not found response
func NotFoundResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	c.JSON(http.StatusNotFound, models.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusNotFound,
	})
}

// InternalServerErrorResponse sends internal server error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	c.JSON(http.StatusInternalServerError, models.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusInternalServerError,
	})
}

// ConflictResponse sends conflict response
func ConflictResponse(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, models.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusConflict,
	})
}

// PaginatedResponse sends paginated response
func PaginatedResponse(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, models.PaginationResult{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetClientIP extracts client IP address from request
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	// Check X-Real-IP header
	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fallback to RemoteAddr
	return c.ClientIP()
}

// GetUserAgent extracts user agent from request
func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}
