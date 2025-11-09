package models

import (
	"time"

	"github.com/google/uuid"
)

// Asset represents a fixed asset
type Asset struct {
	BaseModel
	AssetNumber      string     `gorm:"size:50;uniqueIndex;not null" json:"asset_number"`
	Name             string     `gorm:"size:200;not null;index" json:"name"`
	CategoryID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"category_id"`
	BranchID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	
	// Asset Details
	Description      string     `gorm:"type:text" json:"description,omitempty"`
	Brand            string     `gorm:"size:100" json:"brand,omitempty"`
	Model            string     `gorm:"size:100" json:"model,omitempty"`
	SerialNumber     string     `gorm:"size:100;unique" json:"serial_number,omitempty"`
	
	// Financial
	PurchaseDate     time.Time  `gorm:"not null" json:"purchase_date"`
	PurchasePrice    float64    `gorm:"type:decimal(15,2);not null" json:"purchase_price"`
	SupplierName     string     `gorm:"size:200" json:"supplier_name,omitempty"`
	InvoiceNumber    string     `gorm:"size:100" json:"invoice_number,omitempty"`
	
	// Depreciation
	DepreciationMethod string   `gorm:"size:50" json:"depreciation_method"` // straight_line, declining_balance
	UsefulLife       int        `gorm:"default:0" json:"useful_life"` // in years
	SalvageValue     float64    `gorm:"type:decimal(15,2);default:0" json:"salvage_value"`
	AccumulatedDepreciation float64 `gorm:"type:decimal(15,2);default:0" json:"accumulated_depreciation"`
	BookValue        float64    `gorm:"type:decimal(15,2)" json:"book_value"`
	
	// Location
	Location         string     `gorm:"size:200" json:"location,omitempty"`
	Room             string     `gorm:"size:100" json:"room,omitempty"`
	
	// Condition
	Condition        string     `gorm:"size:20;not null;default:'good'" json:"condition"` // excellent, good, fair, poor, damaged
	Status           string     `gorm:"size:20;not null;default:'active'" json:"status"` // active, disposed, lost, sold
	
	// Responsibility
	ResponsibleUser  *uuid.UUID `gorm:"type:uuid" json:"responsible_user_id,omitempty"`
	
	// Insurance
	InsuranceNumber  string     `gorm:"size:100" json:"insurance_number,omitempty"`
	InsuranceExpiry  *time.Time `json:"insurance_expiry,omitempty"`
	
	// Disposal
	DisposalDate     *time.Time `json:"disposal_date,omitempty"`
	DisposalValue    float64    `gorm:"type:decimal(15,2);default:0" json:"disposal_value"`
	DisposalReason   string     `gorm:"type:text" json:"disposal_reason,omitempty"`
	
	// Documents
	PhotoURL         string     `gorm:"type:text" json:"photo_url,omitempty"`
	DocumentURLs     string     `gorm:"type:text" json:"document_urls,omitempty"` // JSON array
	
	Notes            string     `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Category         AssetCategory      `gorm:"foreignKey:CategoryID" json:"category"`
	Branch           Branch             `gorm:"foreignKey:BranchID" json:"branch"`
	ResponsiblePerson *User             `gorm:"foreignKey:ResponsibleUser" json:"responsible_person,omitempty"`
	Depreciations    []AssetDepreciation `gorm:"foreignKey:AssetID" json:"depreciations,omitempty"`
	Maintenances     []AssetMaintenance  `gorm:"foreignKey:AssetID" json:"maintenances,omitempty"`
	Transfers        []AssetTransfer     `gorm:"foreignKey:AssetID" json:"transfers,omitempty"`
}

// TableName specifies table name
func (Asset) TableName() string {
	return "assets"
}

// AssetCategory represents asset category
type AssetCategory struct {
	BaseModel
	Code             string  `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Name             string  `gorm:"size:100;not null" json:"name"`
	Description      string  `gorm:"type:text" json:"description,omitempty"`
	AccountID        *uuid.UUID `gorm:"type:uuid" json:"account_id,omitempty"` // Link to COA
	DefaultUsefulLife int    `gorm:"default:5" json:"default_useful_life"`
	IsActive         bool    `gorm:"default:true" json:"is_active"`
}

// TableName specifies table name
func (AssetCategory) TableName() string {
	return "asset_categories"
}

// AssetDepreciation represents depreciation calculation
type AssetDepreciation struct {
	BaseModel
	AssetID          uuid.UUID `gorm:"type:uuid;not null;index" json:"asset_id"`
	Period           string    `gorm:"size:7;not null;index" json:"period"` // YYYY-MM
	DepreciationAmount float64 `gorm:"type:decimal(15,2);not null" json:"depreciation_amount"`
	AccumulatedAmount float64  `gorm:"type:decimal(15,2);not null" json:"accumulated_amount"`
	BookValue        float64   `gorm:"type:decimal(15,2);not null" json:"book_value"`
	
	// Journal Integration
	IsPosted         bool      `gorm:"default:false" json:"is_posted"`
	JournalID        *uuid.UUID `gorm:"type:uuid" json:"journal_id,omitempty"`
	
	Notes            string    `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Asset            Asset     `gorm:"foreignKey:AssetID" json:"asset"`
}

// TableName specifies table name
func (AssetDepreciation) TableName() string {
	return "asset_depreciations"
}

// AssetMaintenance represents maintenance record
type AssetMaintenance struct {
	BaseModel
	AssetID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"asset_id"`
	MaintenanceType  string     `gorm:"size:50;not null" json:"maintenance_type"` // routine, repair, emergency
	MaintenanceDate  time.Time  `gorm:"not null;index" json:"maintenance_date"`
	Description      string     `gorm:"type:text;not null" json:"description"`
	Cost             float64    `gorm:"type:decimal(15,2);default:0" json:"cost"`
	Vendor           string     `gorm:"size:200" json:"vendor,omitempty"`
	PerformedBy      string     `gorm:"size:200" json:"performed_by,omitempty"`
	NextScheduled    *time.Time `json:"next_scheduled,omitempty"`
	Status           string     `gorm:"size:20;not null;default:'completed'" json:"status"` // scheduled, completed, cancelled
	Notes            string     `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Asset            Asset      `gorm:"foreignKey:AssetID" json:"asset"`
}

// TableName specifies table name
func (AssetMaintenance) TableName() string {
	return "asset_maintenances"
}

// AssetTransfer represents asset transfer between branches
type AssetTransfer struct {
	BaseModel
	AssetID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"asset_id"`
	FromBranchID     uuid.UUID  `gorm:"type:uuid;not null" json:"from_branch_id"`
	ToBranchID       uuid.UUID  `gorm:"type:uuid;not null" json:"to_branch_id"`
	TransferDate     time.Time  `gorm:"not null;index" json:"transfer_date"`
	Reason           string     `gorm:"type:text;not null" json:"reason"`
	RequestedBy      uuid.UUID  `gorm:"type:uuid;not null" json:"requested_by"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	Status           string     `gorm:"size:20;not null;default:'pending'" json:"status"` // pending, approved, rejected, completed
	Notes            string     `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Asset            Asset      `gorm:"foreignKey:AssetID" json:"asset"`
	FromBranch       Branch     `gorm:"foreignKey:FromBranchID" json:"from_branch"`
	ToBranch         Branch     `gorm:"foreignKey:ToBranchID" json:"to_branch"`
	Requester        User       `gorm:"foreignKey:RequestedBy" json:"requester"`
	Approver         *User      `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

// TableName specifies table name
func (AssetTransfer) TableName() string {
	return "asset_transfers"
}

// Asset Status constants
const (
	AssetStatusActive   = "active"
	AssetStatusDisposed = "disposed"
	AssetStatusLost     = "lost"
	AssetStatusSold     = "sold"
)

// Asset Condition constants
const (
	AssetConditionExcellent = "excellent"
	AssetConditionGood      = "good"
	AssetConditionFair      = "fair"
	AssetConditionPoor      = "poor"
	AssetConditionDamaged   = "damaged"
)

// Depreciation Method constants
const (
	DepreciationMethodStraightLine    = "straight_line"
	DepreciationMethodDecliningBalance = "declining_balance"
)

// Maintenance Type constants
const (
	MaintenanceTypeRoutine   = "routine"
	MaintenanceTypeRepair    = "repair"
	MaintenanceTypeEmergency = "emergency"
)

// Maintenance Status constants
const (
	MaintenanceStatusScheduled = "scheduled"
	MaintenanceStatusCompleted = "completed"
	MaintenanceStatusCancelled = "cancelled"
)

// Transfer Status constants
const (
	TransferStatusPending   = "pending"
	TransferStatusApproved  = "approved"
	TransferStatusRejected  = "rejected"
	TransferStatusCompleted = "completed"
)

// CreateAssetRequest for creating asset
type CreateAssetRequest struct {
	Name               string    `json:"name" binding:"required"`
	CategoryID         uuid.UUID `json:"category_id" binding:"required"`
	BranchID           uuid.UUID `json:"branch_id" binding:"required"`
	Description        string    `json:"description"`
	Brand              string    `json:"brand"`
	SerialNumber       string    `json:"serial_number"`
	PurchaseDate       time.Time `json:"purchase_date" binding:"required"`
	PurchasePrice      float64   `json:"purchase_price" binding:"required,gt=0"`
	DepreciationMethod string    `json:"depreciation_method"`
	UsefulLife         int       `json:"useful_life"`
	SalvageValue       float64   `json:"salvage_value"`
	Location           string    `json:"location"`
	ResponsibleUser    *uuid.UUID `json:"responsible_user_id"`
}

// CreateMaintenanceRequest for creating maintenance
type CreateMaintenanceRequest struct {
	AssetID         uuid.UUID `json:"asset_id" binding:"required"`
	MaintenanceType string    `json:"maintenance_type" binding:"required"`
	MaintenanceDate time.Time `json:"maintenance_date" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	Cost            float64   `json:"cost"`
	Vendor          string    `json:"vendor"`
	NextScheduled   *time.Time `json:"next_scheduled"`
}

// CreateTransferRequest for creating transfer
type CreateTransferRequest struct {
	AssetID      uuid.UUID `json:"asset_id" binding:"required"`
	ToBranchID   uuid.UUID `json:"to_branch_id" binding:"required"`
	TransferDate time.Time `json:"transfer_date" binding:"required"`
	Reason       string    `json:"reason" binding:"required"`
}
