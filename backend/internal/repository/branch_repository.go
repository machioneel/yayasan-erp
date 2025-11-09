package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type BranchRepository interface {
	Create(branch *models.Branch) error
	GetByID(id uuid.UUID) (*models.Branch, error)
	GetByCode(code string) (*models.Branch, error)
	GetAll(params *models.PaginationParams) ([]models.Branch, int64, error)
	GetAllActive() ([]models.Branch, error)
	Update(branch *models.Branch) error
	Delete(id uuid.UUID) error
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) Create(branch *models.Branch) error {
	return r.db.Create(branch).Error
}

func (r *branchRepository) GetByID(id uuid.UUID) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.First(&branch, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("branch not found")
		}
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) GetByCode(code string) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.First(&branch, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("branch not found")
		}
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) GetAll(params *models.PaginationParams) ([]models.Branch, int64, error) {
	var branches []models.Branch
	var total int64

	query := r.db.Model(&models.Branch{})

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where(
			"code LIKE ? OR name LIKE ? OR city LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if params.SortBy != "" {
		order := params.SortBy
		if params.SortDesc {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	// Execute query
	if err := query.Find(&branches).Error; err != nil {
		return nil, 0, err
	}

	return branches, total, nil
}

func (r *branchRepository) GetAllActive() ([]models.Branch, error) {
	var branches []models.Branch
	err := r.db.Where("is_active = ?", true).Order("name ASC").Find(&branches).Error
	return branches, err
}

func (r *branchRepository) Update(branch *models.Branch) error {
	return r.db.Save(branch).Error
}

func (r *branchRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Branch{}, "id = ?", id).Error
}
