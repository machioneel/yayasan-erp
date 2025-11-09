package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type JournalRepository interface {
	GetAll(params *models.PaginationParams) ([]models.Journal, int64, error)
	GetByID(id uuid.UUID) (*models.Journal, error)
	GetByJournalNumber(number string) (*models.Journal, error)
	GetByBranch(branchID uuid.UUID, params *models.PaginationParams) ([]models.Journal, int64, error)
	GetByStatus(status string, params *models.PaginationParams) ([]models.Journal, int64, error)
	GetByDateRange(start, end time.Time) ([]models.Journal, error)
	Create(journal *models.Journal) error
	Update(journal *models.Journal) error
	Delete(id uuid.UUID) error
	GenerateJournalNumber(branchCode string, date time.Time) (string, error)
}

type journalRepository struct {
	db *gorm.DB
}

func NewJournalRepository(db *gorm.DB) JournalRepository {
	return &journalRepository{db: db}
}

func (r *journalRepository) GetAll(params *models.PaginationParams) ([]models.Journal, int64, error) {
	var journals []models.Journal
	var total int64

	query := r.db.Model(&models.Journal{})

	// Filter by status if provided
	if params.Search != "" {
		query = query.Where("status = ?", params.Search)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortOrder := "journal_date DESC, created_at DESC"
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
		Preload("Creator").
		Preload("JournalLines").
		Preload("JournalLines.Account").
		Preload("JournalLines.Fund").
		Preload("JournalLines.Program").
		Preload("JournalLines.Donor").
		Order(sortOrder).
		Limit(params.PageSize).
		Offset(offset).
		Find(&journals).Error

	return journals, total, err
}

func (r *journalRepository) GetByID(id uuid.UUID) (*models.Journal, error) {
	var journal models.Journal
	err := r.db.
		Preload("Branch").
		Preload("Creator").
		Preload("Reviewer").
		Preload("Approver").
		Preload("JournalLines").
		Preload("JournalLines.Account").
		Preload("JournalLines.Fund").
		Preload("JournalLines.Program").
		Preload("JournalLines.Donor").
		First(&journal, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("journal not found")
		}
		return nil, err
	}
	return &journal, nil
}

func (r *journalRepository) GetByJournalNumber(number string) (*models.Journal, error) {
	var journal models.Journal
	err := r.db.
		Preload("JournalLines").
		Preload("JournalLines.Account").
		First(&journal, "journal_number = ?", number).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("journal not found")
		}
		return nil, err
	}
	return &journal, nil
}

func (r *journalRepository) GetByBranch(branchID uuid.UUID, params *models.PaginationParams) ([]models.Journal, int64, error) {
	var journals []models.Journal
	var total int64

	query := r.db.Model(&models.Journal{}).Where("branch_id = ?", branchID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Creator").
		Preload("JournalLines").
		Preload("JournalLines.Account").
		Order("journal_date DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&journals).Error

	return journals, total, err
}

func (r *journalRepository) GetByStatus(status string, params *models.PaginationParams) ([]models.Journal, int64, error) {
	var journals []models.Journal
	var total int64

	query := r.db.Model(&models.Journal{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Branch").
		Preload("Creator").
		Preload("JournalLines").
		Order("journal_date DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&journals).Error

	return journals, total, err
}

func (r *journalRepository) GetByDateRange(start, end time.Time) ([]models.Journal, error) {
	var journals []models.Journal
	err := r.db.
		Where("journal_date BETWEEN ? AND ?", start, end).
		Preload("JournalLines").
		Preload("JournalLines.Account").
		Order("journal_date ASC").
		Find(&journals).Error
	return journals, err
}

func (r *journalRepository) Create(journal *models.Journal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create journal
		if err := tx.Create(journal).Error; err != nil {
			return err
		}

		// Calculate totals
		var totalDebit, totalCredit float64
		for _, line := range journal.JournalLines {
			totalDebit += line.Debit
			totalCredit += line.Credit
		}

		// Update journal totals
		journal.TotalDebit = totalDebit
		journal.TotalCredit = totalCredit

		return tx.Model(journal).Updates(map[string]interface{}{
			"total_debit":  totalDebit,
			"total_credit": totalCredit,
		}).Error
	})
}

func (r *journalRepository) Update(journal *models.Journal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing lines
		if err := tx.Where("journal_id = ?", journal.ID).Delete(&models.JournalLine{}).Error; err != nil {
			return err
		}

		// Update journal
		if err := tx.Save(journal).Error; err != nil {
			return err
		}

		// Calculate totals
		var totalDebit, totalCredit float64
		for _, line := range journal.JournalLines {
			totalDebit += line.Debit
			totalCredit += line.Credit
		}

		// Update totals
		journal.TotalDebit = totalDebit
		journal.TotalCredit = totalCredit

		return tx.Model(journal).Updates(map[string]interface{}{
			"total_debit":  totalDebit,
			"total_credit": totalCredit,
		}).Error
	})
}

func (r *journalRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete journal lines first
		if err := tx.Where("journal_id = ?", id).Delete(&models.JournalLine{}).Error; err != nil {
			return err
		}

		// Delete journal
		return tx.Delete(&models.Journal{}, "id = ?", id).Error
	})
}

func (r *journalRepository) GenerateJournalNumber(branchCode string, date time.Time) (string, error) {
	// Format: JE/BRANCH/YYYYMM/XXXX
	// Example: JE/YAY/202411/0001
	
	yearMonth := date.Format("200601")
	prefix := "JE/" + branchCode + "/" + yearMonth + "/"

	// Get last number for this month
	var lastJournal models.Journal
	err := r.db.
		Where("journal_number LIKE ?", prefix+"%").
		Order("journal_number DESC").
		First(&lastJournal).Error

	var sequence int
	if err == nil {
		// Extract sequence from last number
		// Parse the last 4 digits
		lastNum := lastJournal.JournalNumber[len(prefix):]
		// Convert to int and increment
		sequence = 1 // Simple increment for now
	} else {
		sequence = 1
	}

	// Generate number
	journalNumber := prefix + fmt.Sprintf("%04d", sequence)
	
	return journalNumber, nil
}
