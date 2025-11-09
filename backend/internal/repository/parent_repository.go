package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type ParentRepository interface {
	GetAll(params *models.PaginationParams) ([]models.Parent, int64, error)
	GetByID(id uuid.UUID) (*models.Parent, error)
	GetByNIK(nik string) (*models.Parent, error)
	GetByPhone(phone string) (*models.Parent, error)
	Search(keyword string) ([]models.Parent, error)
	Create(parent *models.Parent) error
	Update(parent *models.Parent) error
	Delete(id uuid.UUID) error
}

type parentRepository struct {
	db *gorm.DB
}

func NewParentRepository(db *gorm.DB) ParentRepository {
	return &parentRepository{db: db}
}

func (r *parentRepository) GetAll(params *models.PaginationParams) ([]models.Parent, int64, error) {
	var parents []models.Parent
	var total int64

	query := r.db.Model(&models.Parent{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("full_name LIKE ? OR nik LIKE ? OR phone LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Students").
		Preload("Students.Student").
		Order("full_name ASC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&parents).Error

	return parents, total, err
}

func (r *parentRepository) GetByID(id uuid.UUID) (*models.Parent, error) {
	var parent models.Parent
	err := r.db.
		Preload("Students").
		Preload("Students.Student").
		First(&parent, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("parent not found")
		}
		return nil, err
	}
	return &parent, nil
}

func (r *parentRepository) GetByNIK(nik string) (*models.Parent, error) {
	var parent models.Parent
	err := r.db.First(&parent, "nik = ?", nik).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &parent, nil
}

func (r *parentRepository) GetByPhone(phone string) (*models.Parent, error) {
	var parent models.Parent
	err := r.db.First(&parent, "phone = ?", phone).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &parent, nil
}

func (r *parentRepository) Search(keyword string) ([]models.Parent, error) {
	var parents []models.Parent
	searchPattern := "%" + keyword + "%"
	err := r.db.
		Where("full_name LIKE ? OR nik LIKE ? OR phone LIKE ? OR email LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Order("full_name ASC").
		Limit(20).
		Find(&parents).Error
	return parents, err
}

func (r *parentRepository) Create(parent *models.Parent) error {
	return r.db.Create(parent).Error
}

func (r *parentRepository) Update(parent *models.Parent) error {
	return r.db.Save(parent).Error
}

func (r *parentRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete student-parent relationships
		if err := tx.Where("parent_id = ?", id).Delete(&models.StudentParent{}).Error; err != nil {
			return err
		}
		// Delete parent
		return tx.Delete(&models.Parent{}, "id = ?", id).Error
	})
}
