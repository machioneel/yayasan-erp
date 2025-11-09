package models

import (
	"time"

	"github.com/google/uuid"
)

// Employee represents an employee/teacher
type Employee struct {
	BaseModel
	EmployeeNumber string     `gorm:"size:50;uniqueIndex;not null" json:"employee_number"`
	NIK            string     `gorm:"size:20;unique" json:"nik"` // Nomor Induk Kependudukan
	NPWP           string     `gorm:"size:30;unique" json:"npwp"` // Tax ID
	FullName       string     `gorm:"size:200;not null;index" json:"full_name"`
	Gender         string     `gorm:"size:10;not null" json:"gender"`
	BirthPlace     string     `gorm:"size:100" json:"birth_place"`
	BirthDate      *time.Time `json:"birth_date"`
	Religion       string     `gorm:"size:20" json:"religion"`
	MaritalStatus  string     `gorm:"size:20" json:"marital_status"` // single, married, divorced, widowed
	
	// Contact
	Email          string     `gorm:"size:100;unique" json:"email"`
	Phone          string     `gorm:"size:20;not null" json:"phone"`
	WhatsApp       string     `gorm:"size:20" json:"whatsapp"`
	
	// Address
	Address        string     `gorm:"type:text" json:"address"`
	RT             string     `gorm:"size:5" json:"rt"`
	RW             string     `gorm:"size:5" json:"rw"`
	Village        string     `gorm:"size:100" json:"village"`
	District       string     `gorm:"size:100" json:"district"`
	City           string     `gorm:"size:100" json:"city"`
	Province       string     `gorm:"size:100" json:"province"`
	PostalCode     string     `gorm:"size:10" json:"postal_code"`
	
	// Employment
	BranchID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	Department     string     `gorm:"size:100" json:"department"` // Teaching, Admin, Finance, etc
	Position       string     `gorm:"size:100;not null" json:"position"` // Teacher, Principal, Staff, etc
	EmploymentType string     `gorm:"size:50;not null" json:"employment_type"` // permanent, contract, temporary
	JoinDate       time.Time  `gorm:"not null" json:"join_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	Status         string     `gorm:"size:20;not null;default:'active'" json:"status"` // active, resigned, terminated, retired
	
	// Education
	Education      string     `gorm:"size:50" json:"education"` // SD, SMP, SMA, D3, S1, S2, S3
	Major          string     `gorm:"size:100" json:"major"`
	University     string     `gorm:"size:200" json:"university"`
	
	// Teaching (if teacher)
	IsTeacher      bool       `gorm:"default:false" json:"is_teacher"`
	TeachingSubjects string   `gorm:"type:text" json:"teaching_subjects,omitempty"` // JSON array
	
	// Bank Account
	BankName       string     `gorm:"size:100" json:"bank_name"`
	BankAccount    string     `gorm:"size:50" json:"bank_account"`
	AccountHolder  string     `gorm:"size:200" json:"account_holder"`
	
	// Emergency Contact
	EmergencyName  string     `gorm:"size:200" json:"emergency_name"`
	EmergencyPhone string     `gorm:"size:20" json:"emergency_phone"`
	EmergencyRelation string  `gorm:"size:50" json:"emergency_relation"`
	
	// Documents
	PhotoURL       string     `gorm:"type:text" json:"photo_url,omitempty"`
	CVUrl          string     `gorm:"type:text" json:"cv_url,omitempty"`
	CertificateURLs string    `gorm:"type:text" json:"certificate_urls,omitempty"` // JSON array
	
	// Relationships
	Branch         Branch             `gorm:"foreignKey:BranchID" json:"branch"`
	Contracts      []EmployeeContract `gorm:"foreignKey:EmployeeID" json:"contracts,omitempty"`
	Salaries       []Payroll          `gorm:"foreignKey:EmployeeID" json:"salaries,omitempty"`
}

// TableName specifies table name
func (Employee) TableName() string {
	return "employees"
}

// EmployeeContract represents an employment contract
type EmployeeContract struct {
	BaseModel
	EmployeeID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"employee_id"`
	ContractNumber string     `gorm:"size:50;uniqueIndex;not null" json:"contract_number"`
	ContractType   string     `gorm:"size:50;not null" json:"contract_type"` // permanent, pkwt, freelance
	StartDate      time.Time  `gorm:"not null" json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	Position       string     `gorm:"size:100;not null" json:"position"`
	BaseSalary     float64    `gorm:"type:decimal(15,2);not null" json:"base_salary"`
	
	// Contract Terms
	WorkingHours   int        `gorm:"default:40" json:"working_hours"` // per week
	LeaveDays      int        `gorm:"default:12" json:"leave_days"` // per year
	
	// Status
	Status         string     `gorm:"size:20;not null;default:'active'" json:"status"` // active, expired, terminated
	SignedDate     *time.Time `json:"signed_date,omitempty"`
	TerminationDate *time.Time `json:"termination_date,omitempty"`
	TerminationReason string  `gorm:"type:text" json:"termination_reason,omitempty"`
	
	Notes          string     `gorm:"type:text" json:"notes,omitempty"`
	DocumentURL    string     `gorm:"type:text" json:"document_url,omitempty"`
	
	// Relationships
	Employee       Employee   `gorm:"foreignKey:EmployeeID" json:"employee"`
}

// TableName specifies table name
func (EmployeeContract) TableName() string {
	return "employee_contracts"
}

// Payroll represents salary payment
type Payroll struct {
	BaseModel
	EmployeeID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"employee_id"`
	BranchID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	PayrollNumber   string     `gorm:"size:50;uniqueIndex;not null" json:"payroll_number"`
	Period          string     `gorm:"size:7;not null;index" json:"period"` // YYYY-MM
	PaymentDate     time.Time  `gorm:"not null" json:"payment_date"`
	
	// Salary Components
	BaseSalary      float64    `gorm:"type:decimal(15,2);not null" json:"base_salary"`
	Allowances      float64    `gorm:"type:decimal(15,2);default:0" json:"allowances"`
	Overtime        float64    `gorm:"type:decimal(15,2);default:0" json:"overtime"`
	Bonus           float64    `gorm:"type:decimal(15,2);default:0" json:"bonus"`
	
	// Deductions
	TaxPPh21        float64    `gorm:"type:decimal(15,2);default:0" json:"tax_pph21"`
	BPJS            float64    `gorm:"type:decimal(15,2);default:0" json:"bpjs"`
	Loan            float64    `gorm:"type:decimal(15,2);default:0" json:"loan"`
	OtherDeductions float64    `gorm:"type:decimal(15,2);default:0" json:"other_deductions"`
	
	// Totals
	GrossSalary     float64    `gorm:"type:decimal(15,2);not null" json:"gross_salary"`
	TotalDeductions float64    `gorm:"type:decimal(15,2);default:0" json:"total_deductions"`
	NetSalary       float64    `gorm:"type:decimal(15,2);not null" json:"net_salary"`
	
	// Payment
	PaymentMethod   string     `gorm:"size:50" json:"payment_method"` // transfer, cash
	PaymentStatus   string     `gorm:"size:20;not null;default:'pending'" json:"payment_status"` // pending, paid, cancelled
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	
	// Journal Integration
	IsPosted        bool       `gorm:"default:false" json:"is_posted"`
	JournalID       *uuid.UUID `gorm:"type:uuid" json:"journal_id,omitempty"`
	
	Notes           string     `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Employee        Employee        `gorm:"foreignKey:EmployeeID" json:"employee"`
	Branch          Branch          `gorm:"foreignKey:BranchID" json:"branch"`
	Components      []PayrollComponent `gorm:"foreignKey:PayrollID" json:"components,omitempty"`
}

// TableName specifies table name
func (Payroll) TableName() string {
	return "payrolls"
}

// PayrollComponent represents salary component detail
type PayrollComponent struct {
	BaseModel
	PayrollID   uuid.UUID `gorm:"type:uuid;not null;index" json:"payroll_id"`
	ComponentID uuid.UUID `gorm:"type:uuid;not null" json:"component_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Type        string    `gorm:"size:20;not null" json:"type"` // earning, deduction
	Amount      float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	IsTaxable   bool      `gorm:"default:true" json:"is_taxable"`
	
	// Relationships
	Payroll     Payroll          `gorm:"foreignKey:PayrollID" json:"payroll"`
	Component   SalaryComponent  `gorm:"foreignKey:ComponentID" json:"component"`
}

// TableName specifies table name
func (PayrollComponent) TableName() string {
	return "payroll_components"
}

// SalaryComponent represents salary component configuration
type SalaryComponent struct {
	BaseModel
	Code        string  `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Name        string  `gorm:"size:100;not null" json:"name"`
	Type        string  `gorm:"size:20;not null" json:"type"` // earning, deduction
	Category    string  `gorm:"size:50;not null" json:"category"` // basic, allowance, overtime, tax, insurance, loan
	IsTaxable   bool    `gorm:"default:true" json:"is_taxable"`
	IsFixed     bool    `gorm:"default:false" json:"is_fixed"`
	Amount      float64 `gorm:"type:decimal(15,2)" json:"amount,omitempty"` // For fixed components
	Percentage  float64 `gorm:"type:decimal(5,2)" json:"percentage,omitempty"` // For percentage-based
	AccountID   *uuid.UUID `gorm:"type:uuid" json:"account_id,omitempty"` // For journal posting
	IsActive    bool    `gorm:"default:true" json:"is_active"`
	Description string  `gorm:"type:text" json:"description,omitempty"`
}

// TableName specifies table name
func (SalaryComponent) TableName() string {
	return "salary_components"
}

// Employee Status constants
const (
	EmployeeStatusActive     = "active"
	EmployeeStatusResigned   = "resigned"
	EmployeeStatusTerminated = "terminated"
	EmployeeStatusRetired    = "retired"
	EmployeeStatusSuspended  = "suspended"
)

// Employment Type constants
const (
	EmploymentTypePermanent  = "permanent"
	EmploymentTypeContract   = "contract"
	EmploymentTypeTemporary  = "temporary"
	EmploymentTypeFreelance  = "freelance"
)

// Marital Status constants
const (
	MaritalStatusSingle   = "single"
	MaritalStatusMarried  = "married"
	MaritalStatusDivorced = "divorced"
	MaritalStatusWidowed  = "widowed"
)

// Contract Type constants
const (
	ContractTypePermanent = "permanent"
	ContractTypePKWT      = "pkwt" // Perjanjian Kerja Waktu Tertentu
	ContractTypeFreelance = "freelance"
)

// Contract Status constants
const (
	ContractStatusActive     = "active"
	ContractStatusExpired    = "expired"
	ContractStatusTerminated = "terminated"
)

// Payroll Status constants
const (
	PayrollStatusPending   = "pending"
	PayrollStatusPaid      = "paid"
	PayrollStatusCancelled = "cancelled"
)

// Salary Component Type constants
const (
	ComponentTypeEarning   = "earning"
	ComponentTypeDeduction = "deduction"
)

// Salary Component Category constants
const (
	ComponentCategoryBasic     = "basic"
	ComponentCategoryAllowance = "allowance"
	ComponentCategoryOvertime  = "overtime"
	ComponentCategoryBonus     = "bonus"
	ComponentCategoryTax       = "tax"
	ComponentCategoryInsurance = "insurance"
	ComponentCategoryLoan      = "loan"
	ComponentCategoryOther     = "other"
)

// CreateEmployeeRequest for creating employee
type CreateEmployeeRequest struct {
	NIK              string     `json:"nik"`
	NPWP             string     `json:"npwp"`
	FullName         string     `json:"full_name" binding:"required"`
	Gender           string     `json:"gender" binding:"required,oneof=male female"`
	BirthPlace       string     `json:"birth_place"`
	BirthDate        *time.Time `json:"birth_date"`
	Religion         string     `json:"religion"`
	MaritalStatus    string     `json:"marital_status"`
	Email            string     `json:"email" binding:"required,email"`
	Phone            string     `json:"phone" binding:"required"`
	Address          string     `json:"address"`
	City             string     `json:"city"`
	BranchID         uuid.UUID  `json:"branch_id" binding:"required"`
	Department       string     `json:"department"`
	Position         string     `json:"position" binding:"required"`
	EmploymentType   string     `json:"employment_type" binding:"required"`
	JoinDate         time.Time  `json:"join_date" binding:"required"`
	IsTeacher        bool       `json:"is_teacher"`
	Education        string     `json:"education"`
	BankName         string     `json:"bank_name"`
	BankAccount      string     `json:"bank_account"`
}

// CreatePayrollRequest for creating payroll
type CreatePayrollRequest struct {
	EmployeeID    uuid.UUID              `json:"employee_id" binding:"required"`
	Period        string                 `json:"period" binding:"required"` // YYYY-MM
	PaymentDate   time.Time              `json:"payment_date" binding:"required"`
	BaseSalary    float64                `json:"base_salary" binding:"required,gt=0"`
	Components    []PayrollComponentReq  `json:"components"`
}

// PayrollComponentReq for payroll component
type PayrollComponentReq struct {
	ComponentID uuid.UUID `json:"component_id" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
}
