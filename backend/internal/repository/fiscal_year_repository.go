package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type FiscalYearRepository interface {
	GetAll() ([]models.FiscalYear, error)
	GetByID(id uuid.UUID) (*models.FiscalYear, error)
	GetCurrent() (*models.FiscalYear, error)
	Create(fiscalYear *models.FiscalYear) error
	Update(fiscalYear *models.FiscalYear) error
	Delete(id uuid.UUID) error
}

type fiscalYearRepository struct {
	db *gorm.DB
}

func NewFiscalYearRepository(db *gorm.DB) FiscalYearRepository {
	return &fiscalYearRepository{db: db}
}

func (r *fiscalYearRepository) GetAll() ([]models.FiscalYear, error) {
	var fiscalYears []models.FiscalYear
	err := r.db.Order("start_date DESC").Find(&fiscalYears).Error
	return fiscalYears, err
}

func (r *fiscalYearRepository) GetByID(id uuid.UUID) (*models.FiscalYear, error) {
	var fiscalYear models.FiscalYear
	err := r.db.First(&fiscalYear, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("fiscal year not found")
		}
		return nil, err
	}
	return &fiscalYear, nil
}

func (r *fiscalYearRepository) GetCurrent() (*models.FiscalYear, error) {
	var fiscalYear models.FiscalYear
	err := r.db.Where("is_current = ?", true).First(&fiscalYear).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no current fiscal year set")
		}
		return nil, err
	}
	return &fiscalYear, nil
}

func (r *fiscalYearRepository) Create(fiscalYear *models.FiscalYear) error {
	return r.db.Create(fiscalYear).Error
}

func (r *fiscalYearRepository) Update(fiscalYear *models.FiscalYear) error {
	return r.db.Save(fiscalYear).Error
}

func (r *fiscalYearRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.FiscalYear{}, "id = ?", id).Error
}
