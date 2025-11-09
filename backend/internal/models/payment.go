package models

import (
	"time"

	"github.com/google/uuid"
)

// Payment represents a payment transaction
type Payment struct {
	BaseModel
	InvoiceID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"invoice_id"`
	StudentID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"student_id"`
	BranchID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	PaymentNumber  string     `gorm:"size:50;uniqueIndex;not null" json:"payment_number"`
	PaymentDate    time.Time  `gorm:"not null;index" json:"payment_date"`
	Amount         float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	PaymentMethod  string     `gorm:"size:50;not null" json:"payment_method"` // cash, transfer, card
	ReferenceNo    string     `gorm:"size:100" json:"reference_no,omitempty"`
	Notes          string     `gorm:"type:text" json:"notes,omitempty"`
	ReceivedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"received_by"`
	IsPosted       bool       `gorm:"default:false" json:"is_posted"`
	PostedAt       *time.Time `json:"posted_at,omitempty"`
	JournalID      *uuid.UUID `gorm:"type:uuid" json:"journal_id,omitempty"` // Link to journal entry
	
	// Relationships
	Invoice   Invoice  `gorm:"foreignKey:InvoiceID" json:"invoice"`
	Student   Student  `gorm:"foreignKey:StudentID" json:"student"`
	Branch    Branch   `gorm:"foreignKey:BranchID" json:"branch"`
	Receiver  User     `gorm:"foreignKey:ReceivedBy" json:"receiver"`
	Journal   *Journal `gorm:"foreignKey:JournalID" json:"journal,omitempty"`
}

// TableName specifies table name
func (Payment) TableName() string {
	return "payments"
}

// Invoice represents student invoice/billing
type Invoice struct {
	BaseModel
	StudentID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"student_id"`
	BranchID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	AcademicYearID uuid.UUID  `gorm:"type:uuid;not null;index" json:"academic_year_id"`
	InvoiceNumber  string     `gorm:"size:50;uniqueIndex;not null" json:"invoice_number"`
	InvoiceDate    time.Time  `gorm:"not null;index" json:"invoice_date"`
	DueDate        time.Time  `gorm:"not null;index" json:"due_date"`
	Description    string     `gorm:"type:text;not null" json:"description"`
	TotalAmount    float64    `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	PaidAmount     float64    `gorm:"type:decimal(15,2);default:0" json:"paid_amount"`
	Status         string     `gorm:"size:20;not null;default:'unpaid'" json:"status"` // unpaid, partial, paid, overdue
	
	// Relationships
	Student      Student       `gorm:"foreignKey:StudentID" json:"student"`
	Branch       Branch        `gorm:"foreignKey:BranchID" json:"branch"`
	AcademicYear AcademicYear  `gorm:"foreignKey:AcademicYearID" json:"academic_year"`
	Items        []InvoiceItem `gorm:"foreignKey:InvoiceID" json:"items,omitempty"`
	Payments     []Payment     `gorm:"foreignKey:InvoiceID" json:"payments,omitempty"`
}

// TableName specifies table name
func (Invoice) TableName() string {
	return "invoices"
}

// InvoiceItem represents invoice line item
type InvoiceItem struct {
	BaseModel
	InvoiceID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"invoice_id"`
	FeeStructureID  *uuid.UUID `gorm:"type:uuid" json:"fee_structure_id,omitempty"`
	Description     string     `gorm:"type:text;not null" json:"description"`
	Quantity        int        `gorm:"default:1" json:"quantity"`
	UnitPrice       float64    `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	Amount          float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	AccountID       *uuid.UUID `gorm:"type:uuid" json:"account_id,omitempty"` // Revenue account
	
	// Relationships
	Invoice      Invoice       `gorm:"foreignKey:InvoiceID" json:"invoice"`
	FeeStructure *FeeStructure `gorm:"foreignKey:FeeStructureID" json:"fee_structure,omitempty"`
	Account      *Account      `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

// TableName specifies table name
func (InvoiceItem) TableName() string {
	return "invoice_items"
}

// FeeStructure represents fee configuration
type FeeStructure struct {
	BaseModel
	BranchID      uuid.UUID `gorm:"type:uuid;not null;index" json:"branch_id"`
	Code          string    `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Name          string    `gorm:"size:200;not null" json:"name"`
	Description   string    `gorm:"type:text" json:"description,omitempty"`
	Amount        float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	FeeType       string    `gorm:"size:50;not null" json:"fee_type"` // monthly, annual, one_time
	Category      string    `gorm:"size:50;not null" json:"category"` // tuition, registration, uniform, book, etc
	AccountID     uuid.UUID `gorm:"type:uuid;not null" json:"account_id"` // Revenue account
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	
	// Relationships
	Branch  Branch  `gorm:"foreignKey:BranchID" json:"branch"`
	Account Account `gorm:"foreignKey:AccountID" json:"account"`
}

// TableName specifies table name
func (FeeStructure) TableName() string {
	return "fee_structures"
}

// Payment Method constants
const (
	PaymentMethodCash     = "cash"
	PaymentMethodTransfer = "transfer"
	PaymentMethodCard     = "card"
	PaymentMethodVA       = "virtual_account"
)

// Invoice Status constants
const (
	InvoiceStatusUnpaid  = "unpaid"
	InvoiceStatusPartial = "partial"
	InvoiceStatusPaid    = "paid"
	InvoiceStatusOverdue = "overdue"
)

// Fee Type constants
const (
	FeeTypeMonthly  = "monthly"
	FeeTypeAnnual   = "annual"
	FeeTypeOneTime  = "one_time"
)

// Fee Category constants
const (
	FeeCategoryTuition      = "tuition"
	FeeCategoryRegistration = "registration"
	FeeCategoryUniform      = "uniform"
	FeeCategoryBook         = "book"
	FeeCategoryActivity     = "activity"
	FeeCategoryOther        = "other"
)
