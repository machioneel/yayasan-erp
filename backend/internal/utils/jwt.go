package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims represents JWT claims
type Claims struct {
	UserID       uuid.UUID   `json:"user_id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	BranchID     *uuid.UUID  `json:"branch_id,omitempty"`
	IsSuperAdmin bool        `json:"is_super_admin"`
	Permissions  []string    `json:"permissions,omitempty"`
	TokenType    string      `json:"token_type"` // access or refresh
	jwt.RegisteredClaims
}

// GenerateToken generates JWT access token
func GenerateToken(userID uuid.UUID, username, email string, branchID *uuid.UUID, isSuperAdmin bool, permissions []string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(config.AppConfig.JWT.Expiry)

	claims := &Claims{
		UserID:       userID,
		Username:     username,
		Email:        email,
		BranchID:     branchID,
		IsSuperAdmin: isSuperAdmin,
		Permissions:  permissions,
		TokenType:    "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "yayasan-erp",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// GenerateRefreshToken generates JWT refresh token
func GenerateRefreshToken(userID uuid.UUID, username string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(config.AppConfig.JWT.RefreshExpiry)

	claims := &Claims{
		UserID:    userID,
		Username:  username,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "yayasan-erp",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// ValidateToken validates JWT token and returns claims
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshAccessToken generates new access token from refresh token
func RefreshAccessToken(refreshTokenString string, permissions []string) (string, error) {
	// Validate refresh token
	claims, err := ValidateToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	// Check if it's a refresh token
	if claims.TokenType != "refresh" {
		return "", ErrInvalidToken
	}

	// Generate new access token
	return GenerateToken(
		claims.UserID,
		claims.Username,
		claims.Email,
		claims.BranchID,
		claims.IsSuperAdmin,
		permissions,
	)
}

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Expected format: "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("invalid authorization header format")
	}

	return authHeader[len(bearerPrefix):], nil
}

// ExtractTokenFromHeader extracts JWT token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header format")
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}
