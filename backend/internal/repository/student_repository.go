package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type StudentRepository interface {
	GetAll(params *models.PaginationParams) ([]models.Student, int64, error)
	GetByID(id uuid.UUID) (*models.Student, error)
	GetByRegistrationNumber(regNumber string) (*models.Student, error)
	GetByNISN(nisn string) (*models.Student, error)
	GetByBranch(branchID uuid.UUID, params *models.PaginationParams) ([]models.Student, int64, error)
	GetByClass(classID uuid.UUID) ([]models.Student, error)
	GetByStatus(status string) ([]models.Student, error)
	Search(keyword string, params *models.PaginationParams) ([]models.Student, int64, error)
	Create(student *models.Student) error
	Update(student *models.Student) error
	Delete(id uuid.UUID) error
	GenerateRegistrationNumber(branchCode string, year int) (string, error)
	CountByBranch(branchID uuid.UUID) (int64, error)
	CountByStatus(status string) (int64, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) GetAll(params *models.PaginationParams) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	query := r.db.Model(&models.Student{})

	// Search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where(
			"full_name LIKE ? OR registration_number LIKE ? OR nisn LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortOrder := "created_at DESC"
	if params.SortBy != "" {
		sortOrder = params.SortBy
		if params.SortOrder == "desc" {
			sortOrder += " DESC"
		} else {
			sortOrder += " ASC"
		}
	}

	// Get paginated results with relationships
	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Branch").
		Preload("CurrentClass").
		Preload("AcademicYear").
		Preload("Parents").
		Preload("Parents.Parent").
		Order(sortOrder).
		Limit(params.PageSize).
		Offset(offset).
		Find(&students).Error

	return students, total, err
}

func (r *studentRepository) GetByID(id uuid.UUID) (*models.Student, error) {
	var student models.Student
	err := r.db.
		Preload("Branch").
		Preload("CurrentClass").
		Preload("AcademicYear").
		Preload("Parents").
		Preload("Parents.Parent").
		First(&student, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByRegistrationNumber(regNumber string) (*models.Student, error) {
	var student models.Student
	err := r.db.
		Preload("Branch").
		Preload("CurrentClass").
		First(&student, "registration_number = ?", regNumber).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByNISN(nisn string) (*models.Student, error) {
	var student models.Student
	err := r.db.First(&student, "nisn = ?", nisn).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if not found (for checking existence)
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByBranch(branchID uuid.UUID, params *models.PaginationParams) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	query := r.db.Model(&models.Student{}).Where("branch_id = ?", branchID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("CurrentClass").
		Preload("Parents").
		Preload("Parents.Parent").
		Order("full_name ASC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&students).Error

	return students, total, err
}

func (r *studentRepository) GetByClass(classID uuid.UUID) ([]models.Student, error) {
	var students []models.Student
	err := r.db.
		Where("current_class_id = ?", classID).
		Order("full_name ASC").
		Find(&students).Error
	return students, err
}

func (r *studentRepository) GetByStatus(status string) ([]models.Student, error) {
	var students []models.Student
	err := r.db.
		Where("status = ?", status).
		Preload("Branch").
		Preload("CurrentClass").
		Order("full_name ASC").
		Find(&students).Error
	return students, err
}

func (r *studentRepository) Search(keyword string, params *models.PaginationParams) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	searchPattern := "%" + keyword + "%"
	query := r.db.Model(&models.Student{}).Where(
		"full_name LIKE ? OR registration_number LIKE ? OR nisn LIKE ? OR email LIKE ? OR phone LIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
	)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Branch").
		Preload("CurrentClass").
		Order("full_name ASC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&students).Error

	return students, total, err
}

func (r *studentRepository) Create(student *models.Student) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create student
		if err := tx.Create(student).Error; err != nil {
			return err
		}

		// Update class student count if assigned to class
		if student.CurrentClassID != nil {
			if err := tx.Model(&models.Class{}).
				Where("id = ?", student.CurrentClassID).
				UpdateColumn("current_students", gorm.Expr("current_students + ?", 1)).
				Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *studentRepository) Update(student *models.Student) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get old student data
		var oldStudent models.Student
		if err := tx.First(&oldStudent, "id = ?", student.ID).Error; err != nil {
			return err
		}

		// Update student
		if err := tx.Save(student).Error; err != nil {
			return err
		}

		// Update class counts if class changed
		if oldStudent.CurrentClassID != student.CurrentClassID {
			// Decrease old class count
			if oldStudent.CurrentClassID != nil {
				if err := tx.Model(&models.Class{}).
					Where("id = ?", oldStudent.CurrentClassID).
					UpdateColumn("current_students", gorm.Expr("current_students - ?", 1)).
					Error; err != nil {
					return err
				}
			}

			// Increase new class count
			if student.CurrentClassID != nil {
				if err := tx.Model(&models.Class{}).
					Where("id = ?", student.CurrentClassID).
					UpdateColumn("current_students", gorm.Expr("current_students + ?", 1)).
					Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *studentRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get student
		var student models.Student
		if err := tx.First(&student, "id = ?", id).Error; err != nil {
			return err
		}

		// Decrease class count if assigned
		if student.CurrentClassID != nil {
			if err := tx.Model(&models.Class{}).
				Where("id = ?", student.CurrentClassID).
				UpdateColumn("current_students", gorm.Expr("current_students - ?", 1)).
				Error; err != nil {
				return err
			}
		}

		// Delete student parents
		if err := tx.Where("student_id = ?", id).Delete(&models.StudentParent{}).Error; err != nil {
			return err
		}

		// Delete student
		return tx.Delete(&models.Student{}, "id = ?", id).Error
	})
}

func (r *studentRepository) GenerateRegistrationNumber(branchCode string, year int) (string, error) {
	// Format: BranchCode/YYYY/XXXX
	// Example: YAY/2025/0001
	
	prefix := fmt.Sprintf("%s/%d/", branchCode, year)

	// Get last number for this year and branch
	var lastStudent models.Student
	err := r.db.
		Where("registration_number LIKE ?", prefix+"%").
		Order("registration_number DESC").
		First(&lastStudent).Error

	var sequence int
	if err == nil {
		// Extract sequence from last number
		// Parse the last 4 digits
		var lastSeq int
		fmt.Sscanf(lastStudent.RegistrationNumber[len(prefix):], "%d", &lastSeq)
		sequence = lastSeq + 1
	} else {
		sequence = 1
	}

	// Generate number
	registrationNumber := fmt.Sprintf("%s%04d", prefix, sequence)
	
	return registrationNumber, nil
}

func (r *studentRepository) CountByBranch(branchID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Student{}).
		Where("branch_id = ?", branchID).
		Count(&count).Error
	return count, err
}

func (r *studentRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Student{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}
