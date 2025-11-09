package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type AccountRepository interface {
	GetAll(params *models.PaginationParams) ([]models.Account, int64, error)
	GetByID(id uuid.UUID) (*models.Account, error)
	GetByCode(code string) (*models.Account, error)
	GetByCategory(category string) ([]models.Account, error)
	GetActiveAccounts() ([]models.Account, error)
	GetDetailAccounts() ([]models.Account, error)
	GetTree() ([]models.Account, error)
	GetChildren(parentID uuid.UUID) ([]models.Account, error)
	Create(account *models.Account) error
	Update(account *models.Account) error
	Delete(id uuid.UUID) error
	BulkCreate(accounts []models.Account) error
	CodeExists(code string) (bool, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) GetAll(params *models.PaginationParams) ([]models.Account, int64, error) {
	var accounts []models.Account
	var total int64

	query := r.db.Model(&models.Account{})

	// Filter by category if provided
	if params.Search != "" {
		query = query.Where("category = ?", params.Search)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortOrder := "code ASC"
	if params.SortBy != "" {
		sortOrder = params.SortBy
		if params.SortOrder == "desc" {
			sortOrder += " DESC"
		} else {
			sortOrder += " ASC"
		}
	}

	// Get paginated results with parent
	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Parent").
		Order(sortOrder).
		Limit(params.PageSize).
		Offset(offset).
		Find(&accounts).Error

	return accounts, total, err
}

func (r *accountRepository) GetByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.
		Preload("Parent").
		Preload("Children").
		First(&account, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetByCode(code string) (*models.Account, error) {
	var account models.Account
	err := r.db.
		Preload("Parent").
		First(&account, "code = ?", code).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetByCategory(category string) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.
		Where("category = ? AND is_active = ?", category, true).
		Order("code ASC").
		Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) GetActiveAccounts() ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.
		Where("is_active = ?", true).
		Order("code ASC").
		Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) GetDetailAccounts() ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.
		Where("is_detail = ? AND is_active = ?", true, true).
		Order("code ASC").
		Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) GetTree() ([]models.Account, error) {
	var accounts []models.Account
	
	// Get all root accounts (no parent)
	err := r.db.
		Where("parent_id IS NULL").
		Order("code ASC").
		Find(&accounts).Error
	
	if err != nil {
		return nil, err
	}

	// Load children recursively
	for i := range accounts {
		if err := r.loadChildren(&accounts[i]); err != nil {
			return nil, err
		}
	}

	return accounts, nil
}

func (r *accountRepository) loadChildren(account *models.Account) error {
	var children []models.Account
	err := r.db.
		Where("parent_id = ?", account.ID).
		Order("code ASC").
		Find(&children).Error
	
	if err != nil {
		return err
	}

	account.Children = children

	// Recursively load children of children
	for i := range account.Children {
		if err := r.loadChildren(&account.Children[i]); err != nil {
			return err
		}
	}

	return nil
}

func (r *accountRepository) GetChildren(parentID uuid.UUID) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.
		Where("parent_id = ?", parentID).
		Order("code ASC").
		Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) Create(account *models.Account) error {
	// Set level based on parent
	if account.ParentID != nil {
		var parent models.Account
		if err := r.db.First(&parent, "id = ?", account.ParentID).Error; err != nil {
			return errors.New("parent account not found")
		}
		account.Level = parent.Level + 1
	} else {
		account.Level = 0
	}

	// Determine if detail account
	account.IsDetail = account.Type == models.AccountTypeDetail || 
					   account.Type == models.AccountTypeIncome ||
					   account.Type == models.AccountTypeRetained ||
					   account.Type == models.AccountTypeRetainedCurr

	// Set normal balance if not provided
	if account.NormalBalance == "" {
		account.NormalBalance = account.GetNormalBalance()
	}

	return r.db.Create(account).Error
}

func (r *accountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepository) Delete(id uuid.UUID) error {
	// Check if account has children
	var count int64
	if err := r.db.Model(&models.Account{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		return errors.New("cannot delete account with children")
	}

	// Check if account has transactions (will implement later)
	// For now, just soft delete
	return r.db.Delete(&models.Account{}, "id = ?", id).Error
}

func (r *accountRepository) BulkCreate(accounts []models.Account) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, account := range accounts {
			// Set level based on parent
			if account.ParentID != nil {
				var parent models.Account
				if err := tx.First(&parent, "id = ?", account.ParentID).Error; err != nil {
					return err
				}
				account.Level = parent.Level + 1
			} else {
				account.Level = 0
			}

			// Determine if detail account
			account.IsDetail = account.Type == models.AccountTypeDetail || 
							   account.Type == models.AccountTypeIncome ||
							   account.Type == models.AccountTypeRetained ||
							   account.Type == models.AccountTypeRetainedCurr

			// Set normal balance
			if account.NormalBalance == "" {
				account.NormalBalance = account.GetNormalBalance()
			}

			if err := tx.Create(&account).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *accountRepository) CodeExists(code string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Account{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}
