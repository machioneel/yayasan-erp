package models

import (
	"time"

	"github.com/google/uuid"
)

// Student represents a student
type Student struct {
	BaseModel
	RegistrationNumber string     `gorm:"size:50;uniqueIndex;not null" json:"registration_number"`
	NISN              string     `gorm:"size:20;unique" json:"nisn"` // Nomor Induk Siswa Nasional
	FullName          string     `gorm:"size:200;not null;index" json:"full_name"`
	NickName          string     `gorm:"size:100" json:"nick_name"`
	Gender            string     `gorm:"size:10;not null" json:"gender"` // male, female
	BirthPlace        string     `gorm:"size:100" json:"birth_place"`
	BirthDate         *time.Time `json:"birth_date"`
	Religion          string     `gorm:"size:20" json:"religion"`
	Nationality       string     `gorm:"size:50;default:'Indonesia'" json:"nationality"`
	
	// Contact
	Email             string     `gorm:"size:100" json:"email"`
	Phone             string     `gorm:"size:20" json:"phone"`
	
	// Address
	Address           string     `gorm:"type:text" json:"address"`
	RT                string     `gorm:"size:5" json:"rt"`
	RW                string     `gorm:"size:5" json:"rw"`
	Village           string     `gorm:"size:100" json:"village"` // Kelurahan
	District          string     `gorm:"size:100" json:"district"` // Kecamatan
	City              string     `gorm:"size:100" json:"city"`
	Province          string     `gorm:"size:100" json:"province"`
	PostalCode        string     `gorm:"size:10" json:"postal_code"`
	
	// Education
	PreviousSchool    string     `gorm:"size:200" json:"previous_school"`
	PreviousSchoolAddress string `gorm:"type:text" json:"previous_school_address"`
	
	// Current Status
	BranchID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"branch_id"`
	CurrentClassID    *uuid.UUID `gorm:"type:uuid;index" json:"current_class_id,omitempty"`
	AcademicYearID    *uuid.UUID `gorm:"type:uuid;index" json:"academic_year_id,omitempty"`
	Status            string     `gorm:"size:20;not null;default:'active'" json:"status"` // active, graduated, dropped, transferred
	RegistrationDate  time.Time  `gorm:"not null" json:"registration_date"`
	
	// Documents
	PhotoURL          string     `gorm:"type:text" json:"photo_url,omitempty"`
	BirthCertificate  string     `gorm:"type:text" json:"birth_certificate,omitempty"`
	FamilyCardNumber  string     `gorm:"size:50" json:"family_card_number,omitempty"`
	
	// Health
	BloodType         string     `gorm:"size:5" json:"blood_type,omitempty"`
	Height            float64    `gorm:"type:decimal(5,2)" json:"height,omitempty"` // cm
	Weight            float64    `gorm:"type:decimal(5,2)" json:"weight,omitempty"` // kg
	SpecialNeeds      string     `gorm:"type:text" json:"special_needs,omitempty"`
	
	// Additional
	Notes             string     `gorm:"type:text" json:"notes,omitempty"`
	
	// Relationships
	Branch            Branch              `gorm:"foreignKey:BranchID" json:"branch"`
	CurrentClass      *Class              `gorm:"foreignKey:CurrentClassID" json:"current_class,omitempty"`
	AcademicYear      *AcademicYear       `gorm:"foreignKey:AcademicYearID" json:"academic_year,omitempty"`
	Parents           []StudentParent     `gorm:"foreignKey:StudentID" json:"parents,omitempty"`
	Payments          []Payment           `gorm:"foreignKey:StudentID" json:"payments,omitempty"`
}

// TableName specifies table name
func (Student) TableName() string {
	return "students"
}

// Parent represents a parent/guardian
type Parent struct {
	BaseModel
	FullName        string  `gorm:"size:200;not null;index" json:"full_name"`
	NIK             string  `gorm:"size:20;unique" json:"nik"` // Nomor Induk Kependudukan
	Gender          string  `gorm:"size:10;not null" json:"gender"`
	BirthPlace      string  `gorm:"size:100" json:"birth_place"`
	BirthDate       *time.Time `json:"birth_date,omitempty"`
	Religion        string  `gorm:"size:20" json:"religion"`
	Nationality     string  `gorm:"size:50;default:'Indonesia'" json:"nationality"`
	
	// Contact
	Email           string  `gorm:"size:100" json:"email"`
	Phone           string  `gorm:"size:20;not null" json:"phone"`
	WhatsApp        string  `gorm:"size:20" json:"whatsapp"`
	
	// Address
	Address         string  `gorm:"type:text" json:"address"`
	RT              string  `gorm:"size:5" json:"rt"`
	RW              string  `gorm:"size:5" json:"rw"`
	Village         string  `gorm:"size:100" json:"village"`
	District        string  `gorm:"size:100" json:"district"`
	City            string  `gorm:"size:100" json:"city"`
	Province        string  `gorm:"size:100" json:"province"`
	PostalCode      string  `gorm:"size:10" json:"postal_code"`
	
	// Employment
	Occupation      string  `gorm:"size:100" json:"occupation"`
	Company         string  `gorm:"size:200" json:"company,omitempty"`
	MonthlyIncome   float64 `gorm:"type:decimal(15,2)" json:"monthly_income,omitempty"`
	
	// Education
	Education       string  `gorm:"size:50" json:"education"` // SD, SMP, SMA, D3, S1, S2, S3
	
	// Relationships
	Students        []StudentParent `gorm:"foreignKey:ParentID" json:"students,omitempty"`
}

// TableName specifies table name
func (Parent) TableName() string {
	return "parents"
}

// StudentParent represents student-parent relationship
type StudentParent struct {
	BaseModel
	StudentID        uuid.UUID `gorm:"type:uuid;not null;index:idx_student_parent" json:"student_id"`
	ParentID         uuid.UUID `gorm:"type:uuid;not null;index:idx_student_parent" json:"parent_id"`
	Relationship     string    `gorm:"size:20;not null" json:"relationship"` // father, mother, guardian
	IsPrimaryContact bool      `gorm:"default:false" json:"is_primary_contact"`
	IsFinancial      bool      `gorm:"default:false" json:"is_financial"` // Responsible for payment
	
	// Relationships
	Student Student `gorm:"foreignKey:StudentID" json:"student"`
	Parent  Parent  `gorm:"foreignKey:ParentID" json:"parent"`
}

// TableName specifies table name
func (StudentParent) TableName() string {
	return "student_parents"
}

// Class represents a class/grade
type Class struct {
	BaseModel
	Code            string    `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Name            string    `gorm:"size:100;not null" json:"name"`
	Level           string    `gorm:"size:20;not null" json:"level"` // TK-A, TK-B, SD-1, SD-2, etc
	BranchID        uuid.UUID `gorm:"type:uuid;not null;index" json:"branch_id"`
	AcademicYearID  uuid.UUID `gorm:"type:uuid;not null;index" json:"academic_year_id"`
	HomeRoomTeacherID *uuid.UUID `gorm:"type:uuid" json:"homeroom_teacher_id,omitempty"`
	Capacity        int       `gorm:"default:30" json:"capacity"`
	CurrentStudents int       `gorm:"default:0" json:"current_students"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	
	// Relationships
	Branch          Branch        `gorm:"foreignKey:BranchID" json:"branch"`
	AcademicYear    AcademicYear  `gorm:"foreignKey:AcademicYearID" json:"academic_year"`
	Students        []Student     `gorm:"foreignKey:CurrentClassID" json:"students,omitempty"`
}

// TableName specifies table name
func (Class) TableName() string {
	return "classes"
}

// AcademicYear represents academic/school year
type AcademicYear struct {
	BaseModel
	Name        string     `gorm:"size:50;uniqueIndex;not null" json:"name"` // e.g. "2025/2026"
	StartDate   time.Time  `gorm:"not null" json:"start_date"`
	EndDate     time.Time  `gorm:"not null" json:"end_date"`
	IsCurrent   bool       `gorm:"default:false" json:"is_current"`
	IsClosed    bool       `gorm:"default:false" json:"is_closed"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	
	// Relationships
	Classes     []Class    `gorm:"foreignKey:AcademicYearID" json:"classes,omitempty"`
	Students    []Student  `gorm:"foreignKey:AcademicYearID" json:"students,omitempty"`
}

// TableName specifies table name
func (AcademicYear) TableName() string {
	return "academic_years"
}

// Student Status constants
const (
	StudentStatusActive      = "active"
	StudentStatusGraduated   = "graduated"
	StudentStatusDropped     = "dropped"
	StudentStatusTransferred = "transferred"
	StudentStatusSuspended   = "suspended"
)

// Gender constants
const (
	GenderMale   = "male"
	GenderFemale = "female"
)

// Parent Relationship constants
const (
	RelationshipFather   = "father"
	RelationshipMother   = "mother"
	RelationshipGuardian = "guardian"
)

// Education Level constants
const (
	EducationSD = "SD"
	EducationSMP = "SMP"
	EducationSMA = "SMA"
	EducationD3 = "D3"
	EducationS1 = "S1"
	EducationS2 = "S2"
	EducationS3 = "S3"
)

// Class Level constants
const (
	LevelTKA  = "TK-A"
	LevelTKB  = "TK-B"
	LevelSD1  = "SD-1"
	LevelSD2  = "SD-2"
	LevelSD3  = "SD-3"
	LevelSD4  = "SD-4"
	LevelSD5  = "SD-5"
	LevelSD6  = "SD-6"
	LevelSMP1 = "SMP-1"
	LevelSMP2 = "SMP-2"
	LevelSMP3 = "SMP-3"
)

// StudentResponse for API responses
type StudentResponse struct {
	ID                 uuid.UUID         `json:"id"`
	RegistrationNumber string            `json:"registration_number"`
	NISN               string            `json:"nisn,omitempty"`
	FullName           string            `json:"full_name"`
	NickName           string            `json:"nick_name,omitempty"`
	Gender             string            `json:"gender"`
	BirthPlace         string            `json:"birth_place,omitempty"`
	BirthDate          *time.Time        `json:"birth_date,omitempty"`
	Religion           string            `json:"religion,omitempty"`
	Email              string            `json:"email,omitempty"`
	Phone              string            `json:"phone,omitempty"`
	Address            string            `json:"address,omitempty"`
	City               string            `json:"city,omitempty"`
	BranchID           uuid.UUID         `json:"branch_id"`
	BranchName         string            `json:"branch_name"`
	CurrentClassID     *uuid.UUID        `json:"current_class_id,omitempty"`
	ClassName          string            `json:"class_name,omitempty"`
	ClassLevel         string            `json:"class_level,omitempty"`
	Status             string            `json:"status"`
	StatusName         string            `json:"status_name"`
	RegistrationDate   time.Time         `json:"registration_date"`
	PhotoURL           string            `json:"photo_url,omitempty"`
	Parents            []ParentSummary   `json:"parents,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

// ParentSummary for student response
type ParentSummary struct {
	ID           uuid.UUID `json:"id"`
	FullName     string    `json:"full_name"`
	Relationship string    `json:"relationship"`
	Phone        string    `json:"phone"`
	IsPrimary    bool      `json:"is_primary_contact"`
	IsFinancial  bool      `json:"is_financial"`
}

// CreateStudentRequest for creating student
type CreateStudentRequest struct {
	RegistrationNumber string               `json:"registration_number"`
	NISN              string               `json:"nisn"`
	FullName          string               `json:"full_name" binding:"required"`
	NickName          string               `json:"nick_name"`
	Gender            string               `json:"gender" binding:"required,oneof=male female"`
	BirthPlace        string               `json:"birth_place"`
	BirthDate         *time.Time           `json:"birth_date"`
	Religion          string               `json:"religion"`
	Email             string               `json:"email"`
	Phone             string               `json:"phone"`
	Address           string               `json:"address"`
	City              string               `json:"city"`
	BranchID          uuid.UUID            `json:"branch_id" binding:"required"`
	CurrentClassID    *uuid.UUID           `json:"current_class_id"`
	RegistrationDate  time.Time            `json:"registration_date"`
	Parents           []CreateParentRequest `json:"parents" binding:"required,min=1"`
}

// CreateParentRequest for parent in student creation
type CreateParentRequest struct {
	ParentID     *uuid.UUID `json:"parent_id"` // If existing parent
	FullName     string     `json:"full_name" binding:"required"`
	NIK          string     `json:"nik"`
	Gender       string     `json:"gender" binding:"required,oneof=male female"`
	Phone        string     `json:"phone" binding:"required"`
	Email        string     `json:"email"`
	Address      string     `json:"address"`
	Occupation   string     `json:"occupation"`
	Relationship string     `json:"relationship" binding:"required,oneof=father mother guardian"`
	IsPrimary    bool       `json:"is_primary_contact"`
	IsFinancial  bool       `json:"is_financial"`
}
