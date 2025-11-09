package models

import (
	"time"

	"github.com/google/uuid"
)

// InventoryItem represents a stock item
type InventoryItem struct {
	BaseModel
	ItemCode         string     `gorm:"size:50;uniqueIndex;not null" json:"item_code"`
	Name             string     `gorm:"size:200;not null;index" json:"name"`
	CategoryID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"category_id"`
	BranchID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	
	// Item Details
	Description      string     `gorm:"type:text" json:"description,omitempty"`
	Unit             string     `gorm:"size:20;not null" json:"unit"` // pcs, box, rim, kg, liter
	
	// Stock
	CurrentStock     float64    `gorm:"type:decimal(15,3);default:0" json:"current_stock"`
	MinimumStock     float64    `gorm:"type:decimal(15,3);default:0" json:"minimum_stock"`
	MaximumStock     float64    `gorm:"type:decimal(15,3);default:0" json:"maximum_stock"`
	
	// Pricing
	PurchasePrice    float64    `gorm:"type:decimal(15,2);default:0" json:"purchase_price"`
	SellingPrice     float64    `gorm:"type:decimal(15,2);default:0" json:"selling_price"`
	
	// Accounting
	InventoryAccountID *uuid.UUID `gorm:"type:uuid" json:"inventory_account_id,omitempty"`
	COGSAccountID      *uuid.UUID `gorm:"type:uuid" json:"cogs_account_id,omitempty"`
	SalesAccountID     *uuid.UUID `gorm:"type:uuid" json:"sales_account_id,omitempty"`
	
	// Status
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	IsSaleable       bool       `gorm:"default:true" json:"is_saleable"` // Can be sold to students
	
	// Documents
	PhotoURL         string     `gorm:"type:text" json:"photo_url,omitempty"`
	
	Notes            string     `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Category         InventoryCategory     `gorm:"foreignKey:CategoryID" json:"category"`
	Branch           Branch                `gorm:"foreignKey:BranchID" json:"branch"`
	Transactions     []StockTransaction    `gorm:"foreignKey:ItemID" json:"transactions,omitempty"`
}

// TableName specifies table name
func (InventoryItem) TableName() string {
	return "inventory_items"
}

// InventoryCategory represents inventory category
type InventoryCategory struct {
	BaseModel
	Code             string     `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Name             string     `gorm:"size:100;not null" json:"name"`
	Description      string     `gorm:"type:text" json:"description,omitempty"`
	IsActive         bool       `gorm:"default:true" json:"is_active"`
}

// TableName specifies table name
func (InventoryCategory) TableName() string {
	return "inventory_categories"
}

// StockTransaction represents stock movement
type StockTransaction struct {
	BaseModel
	TransactionNumber string    `gorm:"size:50;uniqueIndex;not null" json:"transaction_number"`
	ItemID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"item_id"`
	BranchID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	TransactionType  string     `gorm:"size:20;not null;index" json:"transaction_type"` // in, out, adjustment, opname
	TransactionDate  time.Time  `gorm:"not null;index" json:"transaction_date"`
	
	// Quantity
	Quantity         float64    `gorm:"type:decimal(15,3);not null" json:"quantity"`
	UnitPrice        float64    `gorm:"type:decimal(15,2);default:0" json:"unit_price"`
	TotalValue       float64    `gorm:"type:decimal(15,2);default:0" json:"total_value"`
	
	// Stock Balance
	StockBefore      float64    `gorm:"type:decimal(15,3);not null" json:"stock_before"`
	StockAfter       float64    `gorm:"type:decimal(15,3);not null" json:"stock_after"`
	
	// Reference
	ReferenceType    string     `gorm:"size:50" json:"reference_type,omitempty"` // purchase, sale, adjustment
	ReferenceID      *uuid.UUID `gorm:"type:uuid" json:"reference_id,omitempty"`
	ReferenceNumber  string     `gorm:"size:100" json:"reference_number,omitempty"`
	
	// Details
	Supplier         string     `gorm:"size:200" json:"supplier,omitempty"`
	Customer         string     `gorm:"size:200" json:"customer,omitempty"`
	Reason           string     `gorm:"type:text" json:"reason,omitempty"`
	
	// User
	CreatedBy        uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	
	// Journal Integration
	IsPosted         bool       `gorm:"default:false" json:"is_posted"`
	JournalID        *uuid.UUID `gorm:"type:uuid" json:"journal_id,omitempty"`
	
	Notes            string     `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Item             InventoryItem `gorm:"foreignKey:ItemID" json:"item"`
	Branch           Branch        `gorm:"foreignKey:BranchID" json:"branch"`
	User             User          `gorm:"foreignKey:CreatedBy" json:"user"`
}

// TableName specifies table name
func (StockTransaction) TableName() string {
	return "stock_transactions"
}

// StockOpname represents stock taking
type StockOpname struct {
	BaseModel
	OpnameNumber     string     `gorm:"size:50;uniqueIndex;not null" json:"opname_number"`
	BranchID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	OpnameDate       time.Time  `gorm:"not null;index" json:"opname_date"`
	Status           string     `gorm:"size:20;not null;default:'draft'" json:"status"` // draft, completed, approved
	
	// Approval
	PreparedBy       uuid.UUID  `gorm:"type:uuid;not null" json:"prepared_by"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	
	Notes            string     `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Branch           Branch              `gorm:"foreignKey:BranchID" json:"branch"`
	Preparer         User                `gorm:"foreignKey:PreparedBy" json:"preparer"`
	Approver         *User               `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Items            []StockOpnameItem   `gorm:"foreignKey:OpnameID" json:"items"`
}

// TableName specifies table name
func (StockOpname) TableName() string {
	return "stock_opnames"
}

// StockOpnameItem represents item in stock opname
type StockOpnameItem struct {
	BaseModel
	OpnameID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"opname_id"`
	ItemID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"item_id"`
	
	// Stock Count
	SystemStock      float64    `gorm:"type:decimal(15,3);not null" json:"system_stock"`
	PhysicalStock    float64    `gorm:"type:decimal(15,3);not null" json:"physical_stock"`
	Difference       float64    `gorm:"type:decimal(15,3);not null" json:"difference"`
	
	// Value
	UnitPrice        float64    `gorm:"type:decimal(15,2);default:0" json:"unit_price"`
	DifferenceValue  float64    `gorm:"type:decimal(15,2);default:0" json:"difference_value"`
	
	Remarks          string     `gorm:"type:text" json:"remarks,omitempty"`
	
	// Relationships
	Opname           StockOpname   `gorm:"foreignKey:OpnameID" json:"opname"`
	Item             InventoryItem `gorm:"foreignKey:ItemID" json:"item"`
}

// TableName specifies table name
func (StockOpnameItem) TableName() string {
	return "stock_opname_items"
}

// PurchaseOrder represents purchase order
type PurchaseOrder struct {
	BaseModel
	PONumber         string     `gorm:"size:50;uniqueIndex;not null" json:"po_number"`
	BranchID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	SupplierName     string     `gorm:"size:200;not null" json:"supplier_name"`
	SupplierContact  string     `gorm:"size:100" json:"supplier_contact,omitempty"`
	OrderDate        time.Time  `gorm:"not null;index" json:"order_date"`
	ExpectedDate     *time.Time `json:"expected_date,omitempty"`
	
	// Totals
	Subtotal         float64    `gorm:"type:decimal(15,2);default:0" json:"subtotal"`
	Tax              float64    `gorm:"type:decimal(15,2);default:0" json:"tax"`
	TotalAmount      float64    `gorm:"type:decimal(15,2);default:0" json:"total_amount"`
	
	// Status
	Status           string     `gorm:"size:20;not null;default:'draft'" json:"status"` // draft, submitted, received, cancelled
	
	// Approval
	RequestedBy      uuid.UUID  `gorm:"type:uuid;not null" json:"requested_by"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	
	Notes            string     `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Branch           Branch           `gorm:"foreignKey:BranchID" json:"branch"`
	Requester        User             `gorm:"foreignKey:RequestedBy" json:"requester"`
	Items            []PurchaseOrderItem `gorm:"foreignKey:POID" json:"items"`
}

// TableName specifies table name
func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

// PurchaseOrderItem represents item in purchase order
type PurchaseOrderItem struct {
	BaseModel
	POID             uuid.UUID `gorm:"type:uuid;not null;index" json:"po_id"`
	ItemID           uuid.UUID `gorm:"type:uuid;not null" json:"item_id"`
	Quantity         float64   `gorm:"type:decimal(15,3);not null" json:"quantity"`
	UnitPrice        float64   `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	Amount           float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	ReceivedQty      float64   `gorm:"type:decimal(15,3);default:0" json:"received_qty"`
	
	Notes            string    `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	PO               PurchaseOrder `gorm:"foreignKey:POID" json:"po"`
	Item             InventoryItem `gorm:"foreignKey:ItemID" json:"item"`
}

// TableName specifies table name
func (PurchaseOrderItem) TableName() string {
	return "purchase_order_items"
}

// Transaction Type constants
const (
	TransactionTypeIn         = "in"
	TransactionTypeOut        = "out"
	TransactionTypeAdjustment = "adjustment"
	TransactionTypeOpname     = "opname"
)

// Stock Opname Status constants
const (
	OpnameStatusDraft     = "draft"
	OpnameStatusCompleted = "completed"
	OpnameStatusApproved  = "approved"
)

// Purchase Order Status constants
const (
	POStatusDraft     = "draft"
	POStatusSubmitted = "submitted"
	POStatusReceived  = "received"
	POStatusCancelled = "cancelled"
)

// CreateInventoryItemRequest for creating item
type CreateInventoryItemRequest struct {
	Name          string    `json:"name" binding:"required"`
	CategoryID    uuid.UUID `json:"category_id" binding:"required"`
	BranchID      uuid.UUID `json:"branch_id" binding:"required"`
	Description   string    `json:"description"`
	Unit          string    `json:"unit" binding:"required"`
	MinimumStock  float64   `json:"minimum_stock"`
	PurchasePrice float64   `json:"purchase_price"`
	SellingPrice  float64   `json:"selling_price"`
	IsSaleable    bool      `json:"is_saleable"`
}

// CreateStockTransactionRequest for stock transaction
type CreateStockTransactionRequest struct {
	ItemID          uuid.UUID `json:"item_id" binding:"required"`
	TransactionType string    `json:"transaction_type" binding:"required"`
	TransactionDate time.Time `json:"transaction_date" binding:"required"`
	Quantity        float64   `json:"quantity" binding:"required,gt=0"`
	UnitPrice       float64   `json:"unit_price"`
	Supplier        string    `json:"supplier"`
	Customer        string    `json:"customer"`
	Reason          string    `json:"reason"`
}

// CreateStockOpnameRequest for stock opname
type CreateStockOpnameRequest struct {
	BranchID   uuid.UUID            `json:"branch_id" binding:"required"`
	OpnameDate time.Time            `json:"opname_date" binding:"required"`
	Items      []StockOpnameItemReq `json:"items" binding:"required,min=1"`
}

// StockOpnameItemReq for opname item
type StockOpnameItemReq struct {
	ItemID        uuid.UUID `json:"item_id" binding:"required"`
	PhysicalStock float64   `json:"physical_stock" binding:"required,gte=0"`
	Remarks       string    `json:"remarks"`
}
