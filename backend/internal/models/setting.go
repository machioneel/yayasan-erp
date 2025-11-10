package models

// Setting represents system settings
type Setting struct {
	BaseModel
	SettingKey   string `json:"setting_key" gorm:"type:varchar(100);uniqueIndex;not null"`
	SettingValue string `json:"setting_value" gorm:"type:text"`
	SettingType  string `json:"setting_type" gorm:"type:varchar(50)"` // string, boolean, integer, json
	Category     string `json:"category" gorm:"type:varchar(50)"`     // company, system, notification, etc.
	Description  string `json:"description" gorm:"type:text"`
	IsPublic     bool   `json:"is_public" gorm:"default:false"` // Can be accessed without auth
}

// TableName specifies table name
func (Setting) TableName() string {
	return "settings"
}

// SettingResponse represents setting response
type SettingResponse struct {
	ID           string `json:"id"`
	SettingKey   string `json:"setting_key"`
	SettingValue string `json:"setting_value"`
	SettingType  string `json:"setting_type"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	IsPublic     bool   `json:"is_public"`
}

// ToSettingResponse converts Setting to SettingResponse
func ToSettingResponse(s Setting) SettingResponse {
	return SettingResponse{
		ID:           s.ID.String(),
		SettingKey:   s.SettingKey,
		SettingValue: s.SettingValue,
		SettingType:  s.SettingType,
		Category:     s.Category,
		Description:  s.Description,
		IsPublic:     s.IsPublic,
	}
}

// CreateSettingRequest represents request to create setting
type CreateSettingRequest struct {
	SettingKey   string `json:"setting_key" binding:"required"`
	SettingValue string `json:"setting_value" binding:"required"`
	SettingType  string `json:"setting_type"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	IsPublic     bool   `json:"is_public"`
}

// UpdateSettingRequest represents request to update setting
type UpdateSettingRequest struct {
	SettingValue string `json:"setting_value"`
	Description  string `json:"description"`
	IsPublic     bool   `json:"is_public"`
}
