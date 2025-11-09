package models

import (
	"github.com/google/uuid"
)

// Branch represents organization branch
type Branch struct {
	BaseModelWithUser
	Code       string     `json:"code" gorm:"type:varchar(20);uniqueIndex;not null"`
	Name       string     `json:"name" gorm:"type:varchar(255);not null"`
	Type       string     `json:"type" gorm:"type:varchar(50);not null"` // headquarters, school, clinic, mosque, office
	Address    string     `json:"address" gorm:"type:text"`
	City       string     `json:"city" gorm:"type:varchar(100)"`
	Province   string     `json:"province" gorm:"type:varchar(100)"`
	PostalCode string     `json:"postal_code" gorm:"type:varchar(20)"`
	Phone      string     `json:"phone" gorm:"type:varchar(50)"`
	Email      string     `json:"email" gorm:"type:varchar(255)"`
	IsActive   bool       `json:"is_active" gorm:"default:true"`
	
	// Relationships
	Users      []User     `json:"-" gorm:"foreignKey:BranchID"`
}

// TableName specifies table name
func (Branch) TableName() string {
	return "branches"
}

// BranchType constants
const (
	BranchTypeHeadquarters = "headquarters"
	BranchTypeSchool       = "school"
	BranchTypeClinic       = "clinic"
	BranchTypeMosque       = "mosque"
	BranchTypeOffice       = "office"
)

// CreateBranchRequest for creating new branch
type CreateBranchRequest struct {
	Code       string `json:"code" binding:"required,max=20"`
	Name       string `json:"name" binding:"required,max=255"`
	Type       string `json:"type" binding:"required,oneof=headquarters school clinic mosque office"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
	Email      string `json:"email" binding:"omitempty,email"`
}

// UpdateBranchRequest for updating branch
type UpdateBranchRequest struct {
	Name       string `json:"name" binding:"required,max=255"`
	Type       string `json:"type" binding:"required,oneof=headquarters school clinic mosque office"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
	Email      string `json:"email" binding:"omitempty,email"`
	IsActive   *bool  `json:"is_active"`
}

// BranchResponse for API response
type BranchResponse struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Address    string    `json:"address,omitempty"`
	City       string    `json:"city,omitempty"`
	Province   string    `json:"province,omitempty"`
	PostalCode string    `json:"postal_code,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Email      string    `json:"email,omitempty"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

// ToBranchResponse converts Branch to BranchResponse
func (b *Branch) ToBranchResponse() *BranchResponse {
	return &BranchResponse{
		ID:         b.ID,
		Code:       b.Code,
		Name:       b.Name,
		Type:       b.Type,
		Address:    b.Address,
		City:       b.City,
		Province:   b.Province,
		PostalCode: b.PostalCode,
		Phone:      b.Phone,
		Email:      b.Email,
		IsActive:   b.IsActive,
		CreatedAt:  b.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  b.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
