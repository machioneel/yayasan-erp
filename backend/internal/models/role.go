package models

import (
	"time"
	
	"github.com/google/uuid"
)

// Role represents user role
type Role struct {
	BaseModel
	Code           string           `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Name           string           `json:"name" gorm:"type:varchar(100);not null"`
	Description    string           `json:"description" gorm:"type:text"`
	IsSystemRole   bool             `json:"is_system_role" gorm:"default:false"`
	
	// Relationships
	RolePermissions []RolePermission `json:"role_permissions,omitempty" gorm:"foreignKey:RoleID"`
	UserRoles       []UserRole       `json:"-" gorm:"foreignKey:RoleID"`
}

// TableName specifies table name
func (Role) TableName() string {
	return "roles"
}

// Permission represents system permission
type Permission struct {
	BaseModel
	Code        string           `json:"code" gorm:"type:varchar(100);uniqueIndex;not null"`
	Name        string           `json:"name" gorm:"type:varchar(100);not null"`
	Module      string           `json:"module" gorm:"type:varchar(50);not null"` // finance, inventory, sales, etc.
	Description string           `json:"description" gorm:"type:text"`
	
	// Relationships
	RolePermissions []RolePermission `json:"-" gorm:"foreignKey:PermissionID"`
}

// TableName specifies table name
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission junction table for Role and Permission
type RolePermission struct {
	RoleID       uuid.UUID   `json:"role_id" gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID   `json:"permission_id" gorm:"type:uuid;primaryKey"`
	
	// Relationships
	Role         *Role       `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Permission   *Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID"`
}

// TableName specifies table name
func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole junction table for User, Role, and Branch
type UserRole struct {
	UserID     uuid.UUID  `json:"user_id" gorm:"type:uuid;primaryKey"`
	RoleID     uuid.UUID  `json:"role_id" gorm:"type:uuid;primaryKey"`
	BranchID   uuid.UUID  `json:"branch_id" gorm:"type:uuid;primaryKey"`
	AssignedAt *time.Time `json:"assigned_at" gorm:"autoCreateTime"`
	AssignedBy *uuid.UUID `json:"assigned_by" gorm:"type:uuid"`
	
	// Relationships
	User   *User   `json:"-" gorm:"foreignKey:UserID"`
	Role   *Role   `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Branch *Branch `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
}

// TableName specifies table name
func (UserRole) TableName() string {
	return "user_roles"
}

// Module constants
const (
	ModuleFinance   = "finance"
	ModuleInventory = "inventory"
	ModuleSales     = "sales"
	ModuleCRM       = "crm"
	ModulePurchase  = "purchase"
	ModuleAssets    = "assets"
	ModuleReports   = "reports"
	ModuleSystem    = "system"
)

// CreateRoleRequest for creating new role
type CreateRoleRequest struct {
	Code          string      `json:"code" binding:"required,max=50"`
	Name          string      `json:"name" binding:"required,max=100"`
	Description   string      `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

// UpdateRoleRequest for updating role
type UpdateRoleRequest struct {
	Name          string      `json:"name" binding:"required,max=100"`
	Description   string      `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

// AssignRoleRequest for assigning role to user
type AssignRoleRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	RoleID   uuid.UUID `json:"role_id" binding:"required"`
	BranchID uuid.UUID `json:"branch_id" binding:"required"`
}

// RoleResponse for API response
type RoleResponse struct {
	ID          uuid.UUID            `json:"id"`
	Code        string               `json:"code"`
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	IsSystemRole bool                `json:"is_system_role"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   string               `json:"created_at"`
}

// PermissionResponse for API response
type PermissionResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Module      string    `json:"module"`
	Description string    `json:"description,omitempty"`
}

// ToRoleResponse converts Role to RoleResponse
func (r *Role) ToRoleResponse() *RoleResponse {
	response := &RoleResponse{
		ID:           r.ID,
		Code:         r.Code,
		Name:         r.Name,
		Description:  r.Description,
		IsSystemRole: r.IsSystemRole,
		CreatedAt:    r.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if len(r.RolePermissions) > 0 {
		permissions := make([]PermissionResponse, 0)
		for _, rp := range r.RolePermissions {
			if rp.Permission != nil {
				permissions = append(permissions, *rp.Permission.ToPermissionResponse())
			}
		}
		response.Permissions = permissions
	}

	return response
}

// ToPermissionResponse converts Permission to PermissionResponse
func (p *Permission) ToPermissionResponse() *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Module:      p.Module,
		Description: p.Description,
	}
}
