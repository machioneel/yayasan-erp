package models

import (
	"time"

	"github.com/google/uuid"
)

// Journal represents a journal entry
type Journal struct {
	BaseModel
	BranchID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	JournalNumber string     `gorm:"size:50;uniqueIndex;not null" json:"journal_number"`
	JournalDate   time.Time  `gorm:"not null;index" json:"journal_date"`
	Description   string     `gorm:"type:text;not null" json:"description"`
	ReferenceNo   string     `gorm:"size:100" json:"reference_no"`
	Status        string     `gorm:"size:20;not null;default:'draft'" json:"status"`
	TotalDebit    float64    `gorm:"type:decimal(15,2);not null;default:0" json:"total_debit"`
	TotalCredit   float64    `gorm:"type:decimal(15,2);not null;default:0" json:"total_credit"`
	IsPosted      bool       `gorm:"default:false" json:"is_posted"`
	PostedAt      *time.Time `gorm:"index" json:"posted_at,omitempty"`
	PostedBy      *uuid.UUID `gorm:"type:uuid" json:"posted_by,omitempty"`
	
	// Maker-Checker-Approver
	CreatedBy   uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	ReviewedBy  *uuid.UUID `gorm:"type:uuid" json:"reviewed_by,omitempty"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	ApprovedBy  *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	RejectedBy  *uuid.UUID `gorm:"type:uuid" json:"rejected_by,omitempty"`
	RejectedAt  *time.Time `json:"rejected_at,omitempty"`
	RejectReason string    `gorm:"type:text" json:"reject_reason,omitempty"`
	
	// Relationships
	Branch       Branch        `gorm:"foreignKey:BranchID" json:"branch"`
	JournalLines []JournalLine `gorm:"foreignKey:JournalID" json:"journal_lines,omitempty"`
	Creator      User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Reviewer     *User         `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
	Approver     *User         `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

// TableName specifies table name
func (Journal) TableName() string {
	return "journals"
}

// JournalLine represents a journal entry line
type JournalLine struct {
	BaseModel
	JournalID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"journal_id"`
	AccountID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"account_id"`
	Description string     `gorm:"type:text" json:"description"`
	Debit       float64    `gorm:"type:decimal(15,2);not null;default:0" json:"debit"`
	Credit      float64    `gorm:"type:decimal(15,2);not null;default:0" json:"credit"`
	
	// Multi-dimensional accounting
	FundID      *uuid.UUID `gorm:"type:uuid;index" json:"fund_id,omitempty"`
	ProgramID   *uuid.UUID `gorm:"type:uuid;index" json:"program_id,omitempty"`
	DonorID     *uuid.UUID `gorm:"type:uuid;index" json:"donor_id,omitempty"`
	ProjectID   *uuid.UUID `gorm:"type:uuid;index" json:"project_id,omitempty"`
	
	// Relationships
	Journal Journal  `gorm:"foreignKey:JournalID" json:"journal,omitempty"`
	Account Account  `gorm:"foreignKey:AccountID" json:"account"`
	Fund    *Fund    `gorm:"foreignKey:FundID" json:"fund,omitempty"`
	Program *Program `gorm:"foreignKey:ProgramID" json:"program,omitempty"`
	Donor   *Donor   `gorm:"foreignKey:DonorID" json:"donor,omitempty"`
}

// TableName specifies table name
func (JournalLine) TableName() string {
	return "journal_lines"
}

// Fund represents fund type (restricted/unrestricted)
type Fund struct {
	BaseModel
	Code        string `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Name        string `gorm:"size:200;not null" json:"name"`
	Type        string `gorm:"size:20;not null" json:"type"` // restricted, unrestricted
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

// TableName specifies table name
func (Fund) TableName() string {
	return "funds"
}

// Program represents program/project for tracking
type Program struct {
	BaseModel
	Code        string     `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Name        string     `gorm:"size:200;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
}

// TableName specifies table name
func (Program) TableName() string {
	return "programs"
}

// Journal Status constants
const (
	JournalStatusDraft    = "draft"
	JournalStatusReview   = "review"
	JournalStatusApproved = "approved"
	JournalStatusRejected = "rejected"
	JournalStatusPosted   = "posted"
)

// Fund Type constants
const (
	FundTypeRestricted   = "restricted"
	FundTypeUnrestricted = "unrestricted"
)

// IsBalanced checks if journal is balanced
func (j *Journal) IsBalanced() bool {
	return j.TotalDebit == j.TotalCredit
}

// CanPost checks if journal can be posted
func (j *Journal) CanPost() bool {
	return j.Status == JournalStatusApproved && 
		   !j.IsPosted && 
		   j.IsBalanced() &&
		   len(j.JournalLines) > 0
}

// CanApprove checks if journal can be approved
func (j *Journal) CanApprove() bool {
	return j.Status == JournalStatusReview && j.IsBalanced()
}

// JournalResponse for API responses
type JournalResponse struct {
	ID            uuid.UUID            `json:"id"`
	BranchID      uuid.UUID            `json:"branch_id"`
	JournalNumber string               `json:"journal_number"`
	JournalDate   time.Time            `json:"journal_date"`
	Description   string               `json:"description"`
	ReferenceNo   string               `json:"reference_no,omitempty"`
	Status        string               `json:"status"`
	StatusName    string               `json:"status_name"`
	TotalDebit    float64              `json:"total_debit"`
	TotalCredit   float64              `json:"total_credit"`
	IsPosted      bool                 `json:"is_posted"`
	IsBalanced    bool                 `json:"is_balanced"`
	PostedAt      *time.Time           `json:"posted_at,omitempty"`
	CreatedBy     uuid.UUID            `json:"created_by"`
	CreatorName   string               `json:"creator_name,omitempty"`
	ReviewedBy    *uuid.UUID           `json:"reviewed_by,omitempty"`
	ReviewedAt    *time.Time           `json:"reviewed_at,omitempty"`
	ApprovedBy    *uuid.UUID           `json:"approved_by,omitempty"`
	ApprovedAt    *time.Time           `json:"approved_at,omitempty"`
	RejectedBy    *uuid.UUID           `json:"rejected_by,omitempty"`
	RejectedAt    *time.Time           `json:"rejected_at,omitempty"`
	RejectReason  string               `json:"reject_reason,omitempty"`
	JournalLines  []JournalLineResponse `json:"journal_lines,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// JournalLineResponse for API responses
type JournalLineResponse struct {
	ID          uuid.UUID  `json:"id"`
	AccountID   uuid.UUID  `json:"account_id"`
	AccountCode string     `json:"account_code"`
	AccountName string     `json:"account_name"`
	Description string     `json:"description,omitempty"`
	Debit       float64    `json:"debit"`
	Credit      float64    `json:"credit"`
	FundID      *uuid.UUID `json:"fund_id,omitempty"`
	FundName    string     `json:"fund_name,omitempty"`
	ProgramID   *uuid.UUID `json:"program_id,omitempty"`
	ProgramName string     `json:"program_name,omitempty"`
	DonorID     *uuid.UUID `json:"donor_id,omitempty"`
	DonorName   string     `json:"donor_name,omitempty"`
}

// ToJournalResponse converts Journal to JournalResponse
func (j *Journal) ToJournalResponse() *JournalResponse {
	resp := &JournalResponse{
		ID:            j.ID,
		BranchID:      j.BranchID,
		JournalNumber: j.JournalNumber,
		JournalDate:   j.JournalDate,
		Description:   j.Description,
		ReferenceNo:   j.ReferenceNo,
		Status:        j.Status,
		StatusName:    GetJournalStatusName(j.Status),
		TotalDebit:    j.TotalDebit,
		TotalCredit:   j.TotalCredit,
		IsPosted:      j.IsPosted,
		IsBalanced:    j.IsBalanced(),
		PostedAt:      j.PostedAt,
		CreatedBy:     j.CreatedBy,
		ReviewedBy:    j.ReviewedBy,
		ReviewedAt:    j.ReviewedAt,
		ApprovedBy:    j.ApprovedBy,
		ApprovedAt:    j.ApprovedAt,
		RejectedBy:    j.RejectedBy,
		RejectedAt:    j.RejectedAt,
		RejectReason:  j.RejectReason,
		CreatedAt:     j.CreatedAt,
		UpdatedAt:     j.UpdatedAt,
	}

	if j.Creator.ID != uuid.Nil {
		resp.CreatorName = j.Creator.FullName
	}

	if len(j.JournalLines) > 0 {
		resp.JournalLines = make([]JournalLineResponse, len(j.JournalLines))
		for i, line := range j.JournalLines {
			resp.JournalLines[i] = JournalLineResponse{
				ID:          line.ID,
				AccountID:   line.AccountID,
				AccountCode: line.Account.Code,
				AccountName: line.Account.Name,
				Description: line.Description,
				Debit:       line.Debit,
				Credit:      line.Credit,
				FundID:      line.FundID,
				ProgramID:   line.ProgramID,
				DonorID:     line.DonorID,
			}
			
			if line.Fund != nil {
				resp.JournalLines[i].FundName = line.Fund.Name
			}
			if line.Program != nil {
				resp.JournalLines[i].ProgramName = line.Program.Name
			}
			if line.Donor != nil {
				resp.JournalLines[i].DonorName = line.Donor.Name
			}
		}
	}

	return resp
}

// GetJournalStatusName returns human-readable status name
func GetJournalStatusName(status string) string {
	names := map[string]string{
		JournalStatusDraft:    "Draft",
		JournalStatusReview:   "Under Review",
		JournalStatusApproved: "Approved",
		JournalStatusRejected: "Rejected",
		JournalStatusPosted:   "Posted",
	}
	return names[status]
}

// CreateJournalRequest for creating journal
type CreateJournalRequest struct {
	BranchID      uuid.UUID              `json:"branch_id" binding:"required"`
	JournalDate   time.Time              `json:"journal_date" binding:"required"`
	Description   string                 `json:"description" binding:"required"`
	ReferenceNo   string                 `json:"reference_no"`
	JournalLines  []CreateJournalLineReq `json:"journal_lines" binding:"required,min=2,dive"`
}

// CreateJournalLineReq for journal line
type CreateJournalLineReq struct {
	AccountID   uuid.UUID  `json:"account_id" binding:"required"`
	Description string     `json:"description"`
	Debit       float64    `json:"debit" binding:"min=0"`
	Credit      float64    `json:"credit" binding:"min=0"`
	FundID      *uuid.UUID `json:"fund_id"`
	ProgramID   *uuid.UUID `json:"program_id"`
	DonorID     *uuid.UUID `json:"donor_id"`
}

// UpdateJournalRequest for updating journal
type UpdateJournalRequest struct {
	JournalDate   time.Time              `json:"journal_date" binding:"required"`
	Description   string                 `json:"description" binding:"required"`
	ReferenceNo   string                 `json:"reference_no"`
	JournalLines  []CreateJournalLineReq `json:"journal_lines" binding:"required,min=2,dive"`
}

// SubmitForReviewRequest for submitting journal
type SubmitForReviewRequest struct {
	Notes string `json:"notes"`
}

// ReviewJournalRequest for reviewing journal
type ReviewJournalRequest struct {
	Action string `json:"action" binding:"required,oneof=approve reject"`
	Notes  string `json:"notes"`
}

// PostJournalRequest for posting journal
type PostJournalRequest struct {
	PostDate time.Time `json:"post_date" binding:"required"`
}

// JournalListResponse for paginated journal list
type JournalListResponse struct {
	Journals   []JournalResponse `json:"journals"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
