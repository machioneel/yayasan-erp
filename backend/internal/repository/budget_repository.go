package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type BudgetRepository interface {
	GetAll(params *models.PaginationParams) ([]models.Budget, int64, error)
	GetByID(id uuid.UUID) (*models.Budget, error)
	GetByFiscalYear(fiscalYearID uuid.UUID) ([]models.Budget, error)
	GetByAccountPeriod(accountID uuid.UUID, period string, branchID, fundID, programID *uuid.UUID) (*models.Budget, error)
	Create(budget *models.Budget) error
	Update(budget *models.Budget) error
	Delete(id uuid.UUID) error
}

type budgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) GetAll(params *models.PaginationParams) ([]models.Budget, int64, error) {
	var budgets []models.Budget
	var total int64

	query := r.db.Model(&models.Budget{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("FiscalYear").
		Preload("Branch").
		Preload("Account").
		Preload("Fund").
		Preload("Program").
		Order("period DESC, account_id ASC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&budgets).Error

	return budgets, total, err
}

func (r *budgetRepository) GetByID(id uuid.UUID) (*models.Budget, error) {
	var budget models.Budget
	err := r.db.
		Preload("FiscalYear").
		Preload("Branch").
		Preload("Account").
		Preload("Fund").
		Preload("Program").
		First(&budget, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("budget not found")
		}
		return nil, err
	}
	return &budget, nil
}

func (r *budgetRepository) GetByFiscalYear(fiscalYearID uuid.UUID) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.db.
		Where("fiscal_year_id = ?", fiscalYearID).
		Preload("Account").
		Preload("Branch").
		Preload("Fund").
		Preload("Program").
		Order("period ASC, account_id ASC").
		Find(&budgets).Error
	return budgets, err
}

func (r *budgetRepository) GetByAccountPeriod(
	accountID uuid.UUID,
	period string,
	branchID, fundID, programID *uuid.UUID,
) (*models.Budget, error) {
	var budget models.Budget
	query := r.db.Where("account_id = ? AND period = ?", accountID, period)

	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	} else {
		query = query.Where("branch_id IS NULL")
	}

	if fundID != nil {
		query = query.Where("fund_id = ?", *fundID)
	} else {
		query = query.Where("fund_id IS NULL")
	}

	if programID != nil {
		query = query.Where("program_id = ?", *programID)
	} else {
		query = query.Where("program_id IS NULL")
	}

	err := query.First(&budget).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &budget, nil
}

func (r *budgetRepository) Create(budget *models.Budget) error {
	return r.db.Create(budget).Error
}

func (r *budgetRepository) Update(budget *models.Budget) error {
	return r.db.Save(budget).Error
}

func (r *budgetRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Budget{}, "id = ?", id).Error
}
