package models

import "github.com/google/uuid"

// Donor represents a donor
type Donor struct {
	BaseModel
	Code        string `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Name        string `gorm:"size:200;not null;index" json:"name"`
	Type        string `gorm:"size:20;not null" json:"type"` // individual, organization
	Email       string `gorm:"size:100" json:"email"`
	Phone       string `gorm:"size:20" json:"phone"`
	Address     string `gorm:"type:text" json:"address"`
	City        string `gorm:"size:100" json:"city"`
	Province    string `gorm:"size:100" json:"province"`
	PostalCode  string `gorm:"size:10" json:"postal_code"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Notes       string `gorm:"type:text" json:"notes"`
}

// TableName specifies table name
func (Donor) TableName() string {
	return "donors"
}

// DonorType constants
const (
	DonorTypeIndividual   = "individual"
	DonorTypeOrganization = "organization"
)

// DonorResponse for API responses
type DonorResponse struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Email      string    `json:"email,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	IsActive   bool      `json:"is_active"`
}

// ToDonorResponse converts Donor to DonorResponse
func (d *Donor) ToDonorResponse() *DonorResponse {
	return &DonorResponse{
		ID:       d.ID,
		Code:     d.Code,
		Name:     d.Name,
		Type:     d.Type,
		Email:    d.Email,
		Phone:    d.Phone,
		IsActive: d.IsActive,
	}
}
