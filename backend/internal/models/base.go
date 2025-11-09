package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BaseModelWithUser extends BaseModel with created_by and updated_by
type BaseModelWithUser struct {
	BaseModel
	CreatedBy *uuid.UUID `json:"created_by,omitempty" gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty" gorm:"type:uuid"`
}

// BeforeCreate sets UUID before creating
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}

// PaginationParams for pagination
type PaginationParams struct {
	Page      int    `form:"page" json:"page"`
	PageSize  int    `form:"page_size" json:"page_size"`
	Search    string `form:"search" json:"search"`
	SortBy    string `form:"sort_by" json:"sort_by"`
	SortDesc  bool   `form:"sort_desc" json:"sort_desc"`
	SortOrder string `form:"sort_order" json:"sort_order"` // asc or desc
}

// PaginationResult generic pagination result
type PaginationResult struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// PaginationResponse alias for PaginationResult
type PaginationResponse = PaginationResult

// Response standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}
