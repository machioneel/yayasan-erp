package models

import (
	"time"

	"github.com/google/uuid"
)

// Account represents a Chart of Account entry
type Account struct {
	BaseModel
	ParentID    *uuid.UUID `gorm:"type:uuid;index" json:"parent_id"`
	Code        string     `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Name        string     `gorm:"size:200;not null" json:"name"`
	NameEn      string     `gorm:"size:200" json:"name_en"`
	Type        string     `gorm:"size:10;not null" json:"type"` // H, SH, B, I, R, R1
	Category    string     `gorm:"size:50;not null" json:"category"` // ASET, KEWAJIBAN, MODAL, PENDAPATAN, BIAYA
	NormalBalance string   `gorm:"size:10;not null" json:"normal_balance"` // debit, credit
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	IsDetail    bool       `gorm:"default:false" json:"is_detail"` // Can post transactions
	Level       int        `gorm:"default:0" json:"level"`
	Description string     `gorm:"type:text" json:"description"`
	
	// Relationships
	Parent   *Account   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Account  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName specifies table name
func (Account) TableName() string {
	return "accounts"
}

// AccountType constants
const (
	AccountTypeHeader       = "H"   // Header
	AccountTypeSubHeader    = "SH"  // Sub-Header
	AccountTypeDetail       = "B"   // Buku/Detail
	AccountTypeIncome       = "I"   // Income/Expense Detail
	AccountTypeRetained     = "R"   // Retained Earnings
	AccountTypeRetainedCurr = "R1"  // Current Year Retained
)

// AccountCategory constants
const (
	AccountCategoryAsset      = "ASET"
	AccountCategoryLiability  = "KEWAJIBAN"
	AccountCategoryEquity     = "MODAL"
	AccountCategoryRevenue    = "PENDAPATAN"
	AccountCategoryExpense    = "BIAYA"
)

// NormalBalance constants
const (
	NormalBalanceDebit  = "debit"
	NormalBalanceCredit = "credit"
)

// IsHeader checks if account is a header
func (a *Account) IsHeader() bool {
	return a.Type == AccountTypeHeader || a.Type == AccountTypeSubHeader
}

// CanPostTransaction checks if account can have transactions
func (a *Account) CanPostTransaction() bool {
	return a.IsDetail && a.IsActive && !a.IsHeader()
}

// GetNormalBalance returns normal balance based on category
func (a *Account) GetNormalBalance() string {
	if a.NormalBalance != "" {
		return a.NormalBalance
	}
	
	// Auto-determine based on category
	switch a.Category {
	case AccountCategoryAsset, AccountCategoryExpense:
		return NormalBalanceDebit
	case AccountCategoryLiability, AccountCategoryEquity, AccountCategoryRevenue:
		return NormalBalanceCredit
	default:
		return NormalBalanceDebit
	}
}

// AccountResponse for API responses
type AccountResponse struct {
	ID            uuid.UUID        `json:"id"`
	ParentID      *uuid.UUID       `json:"parent_id"`
	Code          string           `json:"code"`
	Name          string           `json:"name"`
	NameEn        string           `json:"name_en,omitempty"`
	Type          string           `json:"type"`
	TypeName      string           `json:"type_name"`
	Category      string           `json:"category"`
	NormalBalance string           `json:"normal_balance"`
	IsActive      bool             `json:"is_active"`
	IsDetail      bool             `json:"is_detail"`
	Level         int              `json:"level"`
	Description   string           `json:"description,omitempty"`
	Parent        *AccountResponse `json:"parent,omitempty"`
	Children      []AccountResponse `json:"children,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// ToAccountResponse converts Account to AccountResponse
func (a *Account) ToAccountResponse() *AccountResponse {
	resp := &AccountResponse{
		ID:            a.ID,
		ParentID:      a.ParentID,
		Code:          a.Code,
		Name:          a.Name,
		NameEn:        a.NameEn,
		Type:          a.Type,
		TypeName:      GetAccountTypeName(a.Type),
		Category:      a.Category,
		NormalBalance: a.GetNormalBalance(),
		IsActive:      a.IsActive,
		IsDetail:      a.IsDetail,
		Level:         a.Level,
		Description:   a.Description,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}

	if a.Parent != nil {
		resp.Parent = a.Parent.ToAccountResponse()
	}

	if len(a.Children) > 0 {
		resp.Children = make([]AccountResponse, len(a.Children))
		for i, child := range a.Children {
			resp.Children[i] = *child.ToAccountResponse()
		}
	}

	return resp
}

// GetAccountTypeName returns human-readable type name
func GetAccountTypeName(typeCode string) string {
	names := map[string]string{
		AccountTypeHeader:       "Header",
		AccountTypeSubHeader:    "Sub-Header",
		AccountTypeDetail:       "Detail",
		AccountTypeIncome:       "Income/Expense",
		AccountTypeRetained:     "Retained Earnings",
		AccountTypeRetainedCurr: "Current Year Retained",
	}
	return names[typeCode]
}

// CreateAccountRequest for creating new account
type CreateAccountRequest struct {
	ParentID      *uuid.UUID `json:"parent_id"`
	Code          string     `json:"code" binding:"required,max=20"`
	Name          string     `json:"name" binding:"required,max=200"`
	NameEn        string     `json:"name_en" binding:"max=200"`
	Type          string     `json:"type" binding:"required,oneof=H SH B I R R1"`
	Category      string     `json:"category" binding:"required,oneof=ASET KEWAJIBAN MODAL PENDAPATAN BIAYA"`
	NormalBalance string     `json:"normal_balance" binding:"omitempty,oneof=debit credit"`
	Description   string     `json:"description"`
}

// UpdateAccountRequest for updating account
type UpdateAccountRequest struct {
	Name          string  `json:"name" binding:"required,max=200"`
	NameEn        string  `json:"name_en" binding:"max=200"`
	IsActive      *bool   `json:"is_active"`
	Description   string  `json:"description"`
}

// AccountListResponse for paginated account list
type AccountListResponse struct {
	Accounts   []AccountResponse `json:"accounts"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// AccountTreeNode for hierarchical tree view
type AccountTreeNode struct {
	ID            uuid.UUID          `json:"id"`
	Code          string             `json:"code"`
	Name          string             `json:"name"`
	Type          string             `json:"type"`
	Category      string             `json:"category"`
	IsActive      bool               `json:"is_active"`
	IsDetail      bool               `json:"is_detail"`
	Level         int                `json:"level"`
	Children      []AccountTreeNode  `json:"children,omitempty"`
}

// ImportAccountRequest for bulk import
type ImportAccountRequest struct {
	Accounts []CreateAccountRequest `json:"accounts" binding:"required,min=1"`
}
