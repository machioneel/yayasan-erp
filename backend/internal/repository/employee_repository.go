package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	GetAll(params *models.PaginationParams) ([]models.Employee, int64, error)
	GetByID(id uuid.UUID) (*models.Employee, error)
	GetByNIK(nik string) (*models.Employee, error)
	GetByBranch(branchID uuid.UUID) ([]models.Employee, error)
	GetByStatus(status string) ([]models.Employee, error)
	GetTeachers() ([]models.Employee, error)
	Search(keyword string) ([]models.Employee, error)
	Create(employee *models.Employee) error
	Update(employee *models.Employee) error
	Delete(id uuid.UUID) error
	GenerateEmployeeNumber(branchCode string, year int) (string, error)
	CountByStatus(status string) (int64, error)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) GetAll(params *models.PaginationParams) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.Model(&models.Employee{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("full_name LIKE ? OR employee_number LIKE ? OR nik LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Branch").
		Order("full_name ASC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&employees).Error

	return employees, total, err
}

func (r *employeeRepository) GetByID(id uuid.UUID) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.
		Preload("Branch").
		Preload("Contracts").
		First(&employee, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) GetByNIK(nik string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.First(&employee, "nik = ?", nik).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) GetByBranch(branchID uuid.UUID) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.
		Where("branch_id = ? AND status = ?", branchID, models.EmployeeStatusActive).
		Order("full_name ASC").
		Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) GetByStatus(status string) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.
		Where("status = ?", status).
		Preload("Branch").
		Order("full_name ASC").
		Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) GetTeachers() ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.
		Where("is_teacher = ? AND status = ?", true, models.EmployeeStatusActive).
		Order("full_name ASC").
		Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) Search(keyword string) ([]models.Employee, error) {
	var employees []models.Employee
	searchPattern := "%" + keyword + "%"
	err := r.db.
		Where("full_name LIKE ? OR employee_number LIKE ? OR nik LIKE ? OR email LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Limit(20).
		Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

func (r *employeeRepository) Update(employee *models.Employee) error {
	return r.db.Save(employee).Error
}

func (r *employeeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Employee{}, "id = ?", id).Error
}

func (r *employeeRepository) GenerateEmployeeNumber(branchCode string, year int) (string, error) {
	// Format: EMP/BRANCH/YYYY/XXXX
	prefix := fmt.Sprintf("EMP/%s/%d/", branchCode, year)

	var lastEmployee models.Employee
	err := r.db.
		Where("employee_number LIKE ?", prefix+"%").
		Order("employee_number DESC").
		First(&lastEmployee).Error

	var sequence int
	if err == nil {
		var lastSeq int
		fmt.Sscanf(lastEmployee.EmployeeNumber[len(prefix):], "%d", &lastSeq)
		sequence = lastSeq + 1
	} else {
		sequence = 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}

func (r *employeeRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Employee{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}
