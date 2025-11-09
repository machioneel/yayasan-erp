package models

import (
	"time"

	"github.com/google/uuid"
)

// Budget represents a budget allocation
type Budget struct {
	BaseModel
	FiscalYearID uuid.UUID  `gorm:"type:uuid;not null;index" json:"fiscal_year_id"`
	BranchID     *uuid.UUID `gorm:"type:uuid;index" json:"branch_id,omitempty"`
	AccountID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"account_id"`
	FundID       *uuid.UUID `gorm:"type:uuid;index" json:"fund_id,omitempty"`
	ProgramID    *uuid.UUID `gorm:"type:uuid;index" json:"program_id,omitempty"`
	Period       string     `gorm:"size:7;not null" json:"period"` // YYYY-MM format
	Amount       float64    `gorm:"type:decimal(15,2);not null;default:0" json:"amount"`
	Description  string     `gorm:"type:text" json:"description"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	
	// Relationships
	FiscalYear FiscalYear `gorm:"foreignKey:FiscalYearID" json:"fiscal_year,omitempty"`
	Branch     *Branch    `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Account    Account    `gorm:"foreignKey:AccountID" json:"account"`
	Fund       *Fund      `gorm:"foreignKey:FundID" json:"fund,omitempty"`
	Program    *Program   `gorm:"foreignKey:ProgramID" json:"program,omitempty"`
}

// TableName specifies table name
func (Budget) TableName() string {
	return "budgets"
}

// FiscalYear represents a fiscal year
type FiscalYear struct {
	BaseModel
	Name        string     `gorm:"size:100;not null" json:"name"`
	StartDate   time.Time  `gorm:"not null" json:"start_date"`
	EndDate     time.Time  `gorm:"not null" json:"end_date"`
	IsCurrent   bool       `gorm:"default:false" json:"is_current"`
	IsClosed    bool       `gorm:"default:false" json:"is_closed"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	ClosedBy    *uuid.UUID `gorm:"type:uuid" json:"closed_by,omitempty"`
}

// TableName specifies table name
func (FiscalYear) TableName() string {
	return "fiscal_years"
}

// AccountBalance represents account balance at a point in time
type AccountBalance struct {
	BaseModel
	AccountID     uuid.UUID  `gorm:"type:uuid;not null;index:idx_account_balance" json:"account_id"`
	BranchID      *uuid.UUID `gorm:"type:uuid;index" json:"branch_id,omitempty"`
	FundID        *uuid.UUID `gorm:"type:uuid;index" json:"fund_id,omitempty"`
	ProgramID     *uuid.UUID `gorm:"type:uuid;index" json:"program_id,omitempty"`
	Period        string     `gorm:"size:7;not null;index:idx_account_balance" json:"period"` // YYYY-MM
	OpeningBalance float64   `gorm:"type:decimal(15,2);not null;default:0" json:"opening_balance"`
	Debit         float64    `gorm:"type:decimal(15,2);not null;default:0" json:"debit"`
	Credit        float64    `gorm:"type:decimal(15,2);not null;default:0" json:"credit"`
	ClosingBalance float64   `gorm:"type:decimal(15,2);not null;default:0" json:"closing_balance"`
	
	// Relationships
	Account Account  `gorm:"foreignKey:AccountID" json:"account"`
	Branch  *Branch  `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Fund    *Fund    `gorm:"foreignKey:FundID" json:"fund,omitempty"`
	Program *Program `gorm:"foreignKey:ProgramID" json:"program,omitempty"`
}

// TableName specifies table name
func (AccountBalance) TableName() string {
	return "account_balances"
}

// BudgetResponse for API responses
type BudgetResponse struct {
	ID           uuid.UUID  `json:"id"`
	FiscalYearID uuid.UUID  `json:"fiscal_year_id"`
	FiscalYear   string     `json:"fiscal_year"`
	BranchID     *uuid.UUID `json:"branch_id,omitempty"`
	BranchName   string     `json:"branch_name,omitempty"`
	AccountID    uuid.UUID  `json:"account_id"`
	AccountCode  string     `json:"account_code"`
	AccountName  string     `json:"account_name"`
	FundID       *uuid.UUID `json:"fund_id,omitempty"`
	FundName     string     `json:"fund_name,omitempty"`
	ProgramID    *uuid.UUID `json:"program_id,omitempty"`
	ProgramName  string     `json:"program_name,omitempty"`
	Period       string     `json:"period"`
	Amount       float64    `json:"amount"`
	Actual       float64    `json:"actual"`
	Variance     float64    `json:"variance"`
	VariancePct  float64    `json:"variance_pct"`
	Description  string     `json:"description,omitempty"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
}

// CreateBudgetRequest for creating budget
type CreateBudgetRequest struct {
	FiscalYearID uuid.UUID  `json:"fiscal_year_id" binding:"required"`
	BranchID     *uuid.UUID `json:"branch_id"`
	AccountID    uuid.UUID  `json:"account_id" binding:"required"`
	FundID       *uuid.UUID `json:"fund_id"`
	ProgramID    *uuid.UUID `json:"program_id"`
	Period       string     `json:"period" binding:"required"` // YYYY-MM
	Amount       float64    `json:"amount" binding:"required,min=0"`
	Description  string     `json:"description"`
}

// UpdateBudgetRequest for updating budget
type UpdateBudgetRequest struct {
	Amount      float64 `json:"amount" binding:"required,min=0"`
	Description string  `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

// TrialBalanceRequest for trial balance report
type TrialBalanceRequest struct {
	AsOfDate  time.Time  `json:"as_of_date" binding:"required"`
	BranchID  *uuid.UUID `json:"branch_id"`
	FundID    *uuid.UUID `json:"fund_id"`
	ProgramID *uuid.UUID `json:"program_id"`
}

// TrialBalanceResponse for trial balance report
type TrialBalanceResponse struct {
	AsOfDate time.Time            `json:"as_of_date"`
	Lines    []TrialBalanceLine   `json:"lines"`
	Summary  TrialBalanceSummary  `json:"summary"`
}

// TrialBalanceLine represents a line in trial balance
type TrialBalanceLine struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Category    string  `json:"category"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
}

// TrialBalanceSummary represents trial balance summary
type TrialBalanceSummary struct {
	TotalDebit  float64 `json:"total_debit"`
	TotalCredit float64 `json:"total_credit"`
	Difference  float64 `json:"difference"`
	IsBalanced  bool    `json:"is_balanced"`
}

// BalanceSheetRequest for balance sheet report
type BalanceSheetRequest struct {
	AsOfDate  time.Time  `json:"as_of_date" binding:"required"`
	BranchID  *uuid.UUID `json:"branch_id"`
	FundID    *uuid.UUID `json:"fund_id"`
}

// BalanceSheetResponse for balance sheet report
type BalanceSheetResponse struct {
	AsOfDate         time.Time             `json:"as_of_date"`
	Assets           BalanceSheetSection   `json:"assets"`
	Liabilities      BalanceSheetSection   `json:"liabilities"`
	Equity           BalanceSheetSection   `json:"equity"`
	TotalAssets      float64               `json:"total_assets"`
	TotalLiabilities float64               `json:"total_liabilities"`
	TotalEquity      float64               `json:"total_equity"`
	IsBalanced       bool                  `json:"is_balanced"`
}

// BalanceSheetSection represents a section in balance sheet
type BalanceSheetSection struct {
	Lines []BalanceSheetLine `json:"lines"`
	Total float64            `json:"total"`
}

// BalanceSheetLine represents a line in balance sheet
type BalanceSheetLine struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Amount      float64 `json:"amount"`
	Level       int     `json:"level"`
	IsHeader    bool    `json:"is_header"`
}

// IncomeStatementRequest for income statement report
type IncomeStatementRequest struct {
	StartDate time.Time  `json:"start_date" binding:"required"`
	EndDate   time.Time  `json:"end_date" binding:"required"`
	BranchID  *uuid.UUID `json:"branch_id"`
	FundID    *uuid.UUID `json:"fund_id"`
}

// IncomeStatementResponse for income statement report
type IncomeStatementResponse struct {
	StartDate      time.Time              `json:"start_date"`
	EndDate        time.Time              `json:"end_date"`
	Revenue        IncomeStatementSection `json:"revenue"`
	Expenses       IncomeStatementSection `json:"expenses"`
	TotalRevenue   float64                `json:"total_revenue"`
	TotalExpenses  float64                `json:"total_expenses"`
	NetIncome      float64                `json:"net_income"`
}

// IncomeStatementSection represents a section in income statement
type IncomeStatementSection struct {
	Lines []IncomeStatementLine `json:"lines"`
	Total float64               `json:"total"`
}

// IncomeStatementLine represents a line in income statement
type IncomeStatementLine struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Amount      float64 `json:"amount"`
	Level       int     `json:"level"`
	IsHeader    bool    `json:"is_header"`
}

// GeneralLedgerRequest for general ledger report
type GeneralLedgerRequest struct {
	AccountID uuid.UUID  `json:"account_id" binding:"required"`
	StartDate time.Time  `json:"start_date" binding:"required"`
	EndDate   time.Time  `json:"end_date" binding:"required"`
	BranchID  *uuid.UUID `json:"branch_id"`
}

// GeneralLedgerResponse for general ledger report
type GeneralLedgerResponse struct {
	Account         AccountResponse      `json:"account"`
	StartDate       time.Time            `json:"start_date"`
	EndDate         time.Time            `json:"end_date"`
	OpeningBalance  float64              `json:"opening_balance"`
	Transactions    []GeneralLedgerLine  `json:"transactions"`
	TotalDebit      float64              `json:"total_debit"`
	TotalCredit     float64              `json:"total_credit"`
	ClosingBalance  float64              `json:"closing_balance"`
}

// GeneralLedgerLine represents a transaction in general ledger
type GeneralLedgerLine struct {
	Date          time.Time `json:"date"`
	JournalNumber string    `json:"journal_number"`
	Description   string    `json:"description"`
	Debit         float64   `json:"debit"`
	Credit        float64   `json:"credit"`
	Balance       float64   `json:"balance"`
}

// BudgetVsActualRequest for budget vs actual report
type BudgetVsActualRequest struct {
	FiscalYearID uuid.UUID  `json:"fiscal_year_id" binding:"required"`
	Period       string     `json:"period"` // Optional: specific period (YYYY-MM)
	BranchID     *uuid.UUID `json:"branch_id"`
	AccountID    *uuid.UUID `json:"account_id"` // Optional: specific account
}

// BudgetVsActualResponse for budget vs actual report
type BudgetVsActualResponse struct {
	FiscalYear string               `json:"fiscal_year"`
	Period     string               `json:"period,omitempty"`
	Lines      []BudgetVsActualLine `json:"lines"`
	Summary    BudgetVsActualSummary `json:"summary"`
}

// BudgetVsActualLine represents a line in budget vs actual report
type BudgetVsActualLine struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Budget      float64 `json:"budget"`
	Actual      float64 `json:"actual"`
	Variance    float64 `json:"variance"`
	VariancePct float64 `json:"variance_pct"`
}

// BudgetVsActualSummary represents summary of budget vs actual
type BudgetVsActualSummary struct {
	TotalBudget      float64 `json:"total_budget"`
	TotalActual      float64 `json:"total_actual"`
	TotalVariance    float64 `json:"total_variance"`
	TotalVariancePct float64 `json:"total_variance_pct"`
}
