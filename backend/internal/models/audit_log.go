package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditLog represents audit trail for all operations
type AuditLog struct {
	ID         uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     *uuid.UUID      `json:"user_id" gorm:"type:uuid"`
	BranchID   *uuid.UUID      `json:"branch_id" gorm:"type:uuid"`
	Action     string          `json:"action" gorm:"type:varchar(50);not null"` // create, update, delete, view
	EntityType string          `json:"entity_type" gorm:"type:varchar(100);not null"`
	EntityID   *uuid.UUID      `json:"entity_id" gorm:"type:uuid"`
	OldValues  json.RawMessage `json:"old_values,omitempty" gorm:"type:jsonb"`
	NewValues  json.RawMessage `json:"new_values,omitempty" gorm:"type:jsonb"`
	IPAddress  string          `json:"ip_address" gorm:"type:varchar(45)"`
	UserAgent  string          `json:"user_agent" gorm:"type:text"`
	CreatedAt  time.Time       `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Branch *Branch `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
}

// TableName specifies table name
func (AuditLog) TableName() string {
	return "audit_logs"
}

// Action constants
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionView   = "view"
	ActionLogin  = "login"
	ActionLogout = "logout"
)

// CreateAuditLogRequest for creating audit log
type CreateAuditLogRequest struct {
	UserID     *uuid.UUID      `json:"user_id"`
	BranchID   *uuid.UUID      `json:"branch_id"`
	Action     string          `json:"action"`
	EntityType string          `json:"entity_type"`
	EntityID   *uuid.UUID      `json:"entity_id"`
	OldValues  json.RawMessage `json:"old_values,omitempty"`
	NewValues  json.RawMessage `json:"new_values,omitempty"`
	IPAddress  string          `json:"ip_address"`
	UserAgent  string          `json:"user_agent"`
}

// AuditLogResponse for API response
type AuditLogResponse struct {
	ID         uuid.UUID       `json:"id"`
	UserID     *uuid.UUID      `json:"user_id"`
	BranchID   *uuid.UUID      `json:"branch_id"`
	Action     string          `json:"action"`
	EntityType string          `json:"entity_type"`
	EntityID   *uuid.UUID      `json:"entity_id"`
	OldValues  json.RawMessage `json:"old_values,omitempty"`
	NewValues  json.RawMessage `json:"new_values,omitempty"`
	IPAddress  string          `json:"ip_address"`
	UserName   string          `json:"user_name,omitempty"`
	BranchName string          `json:"branch_name,omitempty"`
	CreatedAt  string          `json:"created_at"`
}

// ToAuditLogResponse converts AuditLog to AuditLogResponse
func (a *AuditLog) ToAuditLogResponse() *AuditLogResponse {
	response := &AuditLogResponse{
		ID:         a.ID,
		UserID:     a.UserID,
		BranchID:   a.BranchID,
		Action:     a.Action,
		EntityType: a.EntityType,
		EntityID:   a.EntityID,
		OldValues:  a.OldValues,
		NewValues:  a.NewValues,
		IPAddress:  a.IPAddress,
		CreatedAt:  a.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if a.User != nil {
		response.UserName = a.User.FullName
	}

	if a.Branch != nil {
		response.BranchName = a.Branch.Name
	}

	return response
}
