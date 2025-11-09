package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents system user
type User struct {
	BaseModel
	BranchID          *uuid.UUID `json:"branch_id" gorm:"type:uuid"`
	Username          string     `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	Email             string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash      string     `json:"-" gorm:"type:varchar(255);not null"`
	FullName          string     `json:"full_name" gorm:"type:varchar(255);not null"`
	Phone             string     `json:"phone" gorm:"type:varchar(50)"`
	AvatarURL         string     `json:"avatar_url" gorm:"type:varchar(500)"`
	IsActive          bool       `json:"is_active" gorm:"default:true"`
	IsSuperAdmin      bool       `json:"is_super_admin" gorm:"default:false"`
	LastLoginAt       *time.Time `json:"last_login_at"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`

	// Relationships
	Branch    *Branch    `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	UserRoles []UserRole `json:"user_roles,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies table name
func (User) TableName() string {
	return "users"
}

// SetPassword hashes and sets the user password
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	now := time.Now()
	u.PasswordChangedAt = &now
	return nil
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// RegisterRequest for user registration
type RegisterRequest struct {
	BranchID *uuid.UUID `json:"branch_id"`
	Username string     `json:"username" binding:"required,min=3,max=100"`
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=8"`
	FullName string     `json:"full_name" binding:"required,max=255"`
	Phone    string     `json:"phone"`
}

// LoginRequest for user login
type LoginRequest struct {
	Login    string `json:"login" binding:"required"` // Can be username or email
	Password string `json:"password" binding:"required"`
}

// LoginResponse after successful login
type LoginResponse struct {
	Token        string        `json:"token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user"`
}

// ChangePasswordRequest for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UpdateProfileRequest for updating user profile
type UpdateProfileRequest struct {
	FullName  string `json:"full_name" binding:"required,max=255"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
}

// CreateUserRequest for admin creating user
type CreateUserRequest struct {
	BranchID     *uuid.UUID  `json:"branch_id"`
	Username     string      `json:"username" binding:"required,min=3,max=100"`
	Email        string      `json:"email" binding:"required,email"`
	Password     string      `json:"password" binding:"required,min=8"`
	FullName     string      `json:"full_name" binding:"required,max=255"`
	Phone        string      `json:"phone"`
	IsActive     *bool       `json:"is_active"`
	IsSuperAdmin *bool       `json:"is_super_admin"`
	RoleIDs      []uuid.UUID `json:"role_ids"`
}

// UpdateUserRequest for admin updating user
type UpdateUserRequest struct {
	BranchID     *uuid.UUID  `json:"branch_id"`
	FullName     string      `json:"full_name" binding:"required,max=255"`
	Phone        string      `json:"phone"`
	AvatarURL    string      `json:"avatar_url"`
	IsActive     *bool       `json:"is_active"`
	IsSuperAdmin *bool       `json:"is_super_admin"`
	RoleIDs      []uuid.UUID `json:"role_ids"`
}

// UserResponse for API response
type UserResponse struct {
	ID           uuid.UUID       `json:"id"`
	BranchID     *uuid.UUID      `json:"branch_id"`
	Username     string          `json:"username"`
	Email        string          `json:"email"`
	FullName     string          `json:"full_name"`
	Phone        string          `json:"phone,omitempty"`
	AvatarURL    string          `json:"avatar_url,omitempty"`
	IsActive     bool            `json:"is_active"`
	IsSuperAdmin bool            `json:"is_super_admin"`
	LastLoginAt  *string         `json:"last_login_at,omitempty"`
	Branch       *BranchResponse `json:"branch,omitempty"`
	Roles        []RoleResponse  `json:"roles,omitempty"`
	Permissions  []string        `json:"permissions,omitempty"`
	CreatedAt    string          `json:"created_at"`
}

// ToUserResponse converts User to UserResponse
func (u *User) ToUserResponse() *UserResponse {
	response := &UserResponse{
		ID:           u.ID,
		BranchID:     u.BranchID,
		Username:     u.Username,
		Email:        u.Email,
		FullName:     u.FullName,
		Phone:        u.Phone,
		AvatarURL:    u.AvatarURL,
		IsActive:     u.IsActive,
		IsSuperAdmin: u.IsSuperAdmin,
		CreatedAt:    u.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if u.LastLoginAt != nil {
		lastLogin := u.LastLoginAt.Format("2006-01-02 15:04:05")
		response.LastLoginAt = &lastLogin
	}

	if u.Branch != nil {
		response.Branch = u.Branch.ToBranchResponse()
	}

	// Convert roles
	if len(u.UserRoles) > 0 {
		roles := make([]RoleResponse, 0)
		permissionSet := make(map[string]bool)
		
		for _, ur := range u.UserRoles {
			if ur.Role != nil {
				roles = append(roles, *ur.Role.ToRoleResponse())
				
				// Collect unique permissions
				for _, rp := range ur.Role.RolePermissions {
					if rp.Permission != nil {
						permissionSet[rp.Permission.Code] = true
					}
				}
			}
		}
		
		response.Roles = roles
		
		// Convert permission set to slice
		permissions := make([]string, 0, len(permissionSet))
		for perm := range permissionSet {
			permissions = append(permissions, perm)
		}
		response.Permissions = permissions
	}

	return response
}

// UserListResponse for listing users
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}
