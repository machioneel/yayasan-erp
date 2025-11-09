package models

import (
	"time"

	"github.com/google/uuid"
)

// Attendance represents employee attendance
type Attendance struct {
	BaseModel
	EmployeeID   uuid.UUID  `gorm:"type:uuid;not null;index:idx_emp_date" json:"employee_id"`
	BranchID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	Date         time.Time  `gorm:"not null;index:idx_emp_date" json:"date"`
	CheckIn      *time.Time `json:"check_in,omitempty"`
	CheckOut     *time.Time `json:"check_out,omitempty"`
	WorkingHours float64    `gorm:"type:decimal(5,2);default:0" json:"working_hours"`
	OvertimeHours float64   `gorm:"type:decimal(5,2);default:0" json:"overtime_hours"`
	Status       string     `gorm:"size:20;not null" json:"status"` // present, absent, late, leave, sick, permit
	Notes        string     `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Employee     Employee   `gorm:"foreignKey:EmployeeID" json:"employee"`
	Branch       Branch     `gorm:"foreignKey:BranchID" json:"branch"`
}

// TableName specifies table name
func (Attendance) TableName() string {
	return "attendances"
}

// Leave represents leave request
type Leave struct {
	BaseModel
	EmployeeID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"employee_id"`
	LeaveType     string     `gorm:"size:50;not null" json:"leave_type"` // annual, sick, maternity, unpaid, permit
	StartDate     time.Time  `gorm:"not null;index" json:"start_date"`
	EndDate       time.Time  `gorm:"not null;index" json:"end_date"`
	TotalDays     int        `gorm:"not null" json:"total_days"`
	Reason        string     `gorm:"type:text;not null" json:"reason"`
	Status        string     `gorm:"size:20;not null;default:'pending'" json:"status"` // pending, approved, rejected, cancelled
	
	// Approval
	ApprovedBy    *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	RejectedBy    *uuid.UUID `gorm:"type:uuid" json:"rejected_by,omitempty"`
	RejectedAt    *time.Time `json:"rejected_at,omitempty"`
	RejectReason  string     `gorm:"type:text" json:"reject_reason,omitempty"`
	
	// Documents
	DocumentURL   string     `gorm:"type:text" json:"document_url,omitempty"` // Medical certificate, etc
	
	// Relationships
	Employee      Employee   `gorm:"foreignKey:EmployeeID" json:"employee"`
	Approver      *User      `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

// TableName specifies table name
func (Leave) TableName() string {
	return "leaves"
}

// Attendance Status constants
const (
	AttendanceStatusPresent = "present"
	AttendanceStatusAbsent  = "absent"
	AttendanceStatusLate    = "late"
	AttendanceStatusLeave   = "leave"
	AttendanceStatusSick    = "sick"
	AttendanceStatusPermit  = "permit"
)

// Leave Type constants
const (
	LeaveTypeAnnual    = "annual"
	LeaveTypeSick      = "sick"
	LeaveTypeMaternity = "maternity"
	LeaveTypeUnpaid    = "unpaid"
	LeaveTypePermit    = "permit"
)

// Leave Status constants
const (
	LeaveStatusPending   = "pending"
	LeaveStatusApproved  = "approved"
	LeaveStatusRejected  = "rejected"
	LeaveStatusCancelled = "cancelled"
)

// CreateLeaveRequest for creating leave
type CreateLeaveRequest struct {
	EmployeeID  uuid.UUID `json:"employee_id" binding:"required"`
	LeaveType   string    `json:"leave_type" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Reason      string    `json:"reason" binding:"required"`
}

// AttendanceReportRequest for attendance report
type AttendanceReportRequest struct {
	BranchID   *uuid.UUID `json:"branch_id"`
	EmployeeID *uuid.UUID `json:"employee_id"`
	StartDate  time.Time  `json:"start_date" binding:"required"`
	EndDate    time.Time  `json:"end_date" binding:"required"`
}
