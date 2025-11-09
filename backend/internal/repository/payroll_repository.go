package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type PayrollRepository interface {
	GetAll(params *models.PaginationParams) ([]models.Payroll, int64, error)
	GetByID(id uuid.UUID) (*models.Payroll, error)
	GetByPeriod(period string) ([]models.Payroll, error)
	GetByEmployee(employeeID uuid.UUID) ([]models.Payroll, error)
	Create(payroll *models.Payroll) error
	Update(payroll *models.Payroll) error
	Delete(id uuid.UUID) error
	GeneratePayrollNumber(branchCode string, date time.Time) (string, error)
}

type payrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) PayrollRepository {
	return &payrollRepository{db: db}
}

func (r *payrollRepository) GetAll(params *models.PaginationParams) ([]models.Payroll, int64, error) {
	var payrolls []models.Payroll
	var total int64

	query := r.db.Model(&models.Payroll{})

	if params.Search != "" {
		query = query.Where("period = ?", params.Search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Employee").
		Preload("Branch").
		Order("period DESC, created_at DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&payrolls).Error

	return payrolls, total, err
}

func (r *payrollRepository) GetByID(id uuid.UUID) (*models.Payroll, error) {
	var payroll models.Payroll
	err := r.db.
		Preload("Employee").
		Preload("Branch").
		Preload("Components").
		Preload("Components.Component").
		First(&payroll, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payroll not found")
		}
		return nil, err
	}
	return &payroll, nil
}

func (r *payrollRepository) GetByPeriod(period string) ([]models.Payroll, error) {
	var payrolls []models.Payroll
	err := r.db.
		Where("period = ?", period).
		Preload("Employee").
		Order("employee_id ASC").
		Find(&payrolls).Error
	return payrolls, err
}

func (r *payrollRepository) GetByEmployee(employeeID uuid.UUID) ([]models.Payroll, error) {
	var payrolls []models.Payroll
	err := r.db.
		Where("employee_id = ?", employeeID).
		Order("period DESC").
		Find(&payrolls).Error
	return payrolls, err
}

func (r *payrollRepository) Create(payroll *models.Payroll) error {
	return r.db.Create(payroll).Error
}

func (r *payrollRepository) Update(payroll *models.Payroll) error {
	return r.db.Save(payroll).Error
}

func (r *payrollRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("payroll_id = ?", id).Delete(&models.PayrollComponent{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Payroll{}, "id = ?", id).Error
	})
}

func (r *payrollRepository) GeneratePayrollNumber(branchCode string, date time.Time) (string, error) {
	// Format: PAY/BRANCH/YYYYMM/XXXX
	yearMonth := date.Format("200601")
	prefix := fmt.Sprintf("PAY/%s/%s/", branchCode, yearMonth)

	var lastPayroll models.Payroll
	err := r.db.
		Where("payroll_number LIKE ?", prefix+"%").
		Order("payroll_number DESC").
		First(&lastPayroll).Error

	var sequence int
	if err == nil {
		var lastSeq int
		fmt.Sscanf(lastPayroll.PayrollNumber[len(prefix):], "%d", &lastSeq)
		sequence = lastSeq + 1
	} else {
		sequence = 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
