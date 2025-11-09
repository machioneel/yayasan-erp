package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	GetAll(params *models.PaginationParams) ([]models.Invoice, int64, error)
	GetByID(id uuid.UUID) (*models.Invoice, error)
	GetByStudent(studentID uuid.UUID) ([]models.Invoice, error)
	GetByStatus(status string) ([]models.Invoice, error)
	GetOverdue() ([]models.Invoice, error)
	Create(invoice *models.Invoice) error
	Update(invoice *models.Invoice) error
	Delete(id uuid.UUID) error
	GenerateInvoiceNumber(branchCode string, date time.Time) (string, error)
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) GetAll(params *models.PaginationParams) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	query := r.db.Model(&models.Invoice{})

	if params.Search != "" {
		query = query.Where("status = ?", params.Search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Student").
		Preload("Branch").
		Preload("Items").
		Order("invoice_date DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&invoices).Error

	return invoices, total, err
}

func (r *invoiceRepository) GetByID(id uuid.UUID) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.
		Preload("Student").
		Preload("Branch").
		Preload("AcademicYear").
		Preload("Items").
		Preload("Items.FeeStructure").
		Preload("Items.Account").
		Preload("Payments").
		First(&invoice, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) GetByStudent(studentID uuid.UUID) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := r.db.
		Where("student_id = ?", studentID).
		Preload("Items").
		Preload("Payments").
		Order("invoice_date DESC").
		Find(&invoices).Error
	return invoices, err
}

func (r *invoiceRepository) GetByStatus(status string) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := r.db.
		Where("status = ?", status).
		Preload("Student").
		Preload("Branch").
		Order("invoice_date DESC").
		Find(&invoices).Error
	return invoices, err
}

func (r *invoiceRepository) GetOverdue() ([]models.Invoice, error) {
	var invoices []models.Invoice
	now := time.Now()
	err := r.db.
		Where("due_date < ? AND status IN (?)", now, []string{models.InvoiceStatusUnpaid, models.InvoiceStatusPartial}).
		Preload("Student").
		Preload("Branch").
		Order("due_date ASC").
		Find(&invoices).Error

	// Update status to overdue
	for i := range invoices {
		invoices[i].Status = models.InvoiceStatusOverdue
		r.db.Model(&invoices[i]).Update("status", models.InvoiceStatusOverdue)
	}

	return invoices, err
}

func (r *invoiceRepository) Create(invoice *models.Invoice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Calculate total amount
		var total float64
		for _, item := range invoice.Items {
			item.Amount = item.UnitPrice * float64(item.Quantity)
			total += item.Amount
		}
		invoice.TotalAmount = total
		invoice.PaidAmount = 0
		invoice.Status = models.InvoiceStatusUnpaid

		// Create invoice
		return tx.Create(invoice).Error
	})
}

func (r *invoiceRepository) Update(invoice *models.Invoice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete old items
		if err := tx.Where("invoice_id = ?", invoice.ID).Delete(&models.InvoiceItem{}).Error; err != nil {
			return err
		}

		// Calculate total amount
		var total float64
		for i := range invoice.Items {
			invoice.Items[i].Amount = invoice.Items[i].UnitPrice * float64(invoice.Items[i].Quantity)
			total += invoice.Items[i].Amount
		}
		invoice.TotalAmount = total

		// Update status based on paid amount
		if invoice.PaidAmount >= invoice.TotalAmount {
			invoice.Status = models.InvoiceStatusPaid
		} else if invoice.PaidAmount > 0 {
			invoice.Status = models.InvoiceStatusPartial
		} else {
			invoice.Status = models.InvoiceStatusUnpaid
		}

		// Update invoice
		return tx.Save(invoice).Error
	})
}

func (r *invoiceRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if has payments
		var count int64
		if err := tx.Model(&models.Payment{}).Where("invoice_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("cannot delete invoice with payments")
		}

		// Delete invoice items
		if err := tx.Where("invoice_id = ?", id).Delete(&models.InvoiceItem{}).Error; err != nil {
			return err
		}

		// Delete invoice
		return tx.Delete(&models.Invoice{}, "id = ?", id).Error
	})
}

func (r *invoiceRepository) GenerateInvoiceNumber(branchCode string, date time.Time) (string, error) {
	// Format: INV/BranchCode/YYYYMM/XXXX
	yearMonth := date.Format("200601")
	prefix := fmt.Sprintf("INV/%s/%s/", branchCode, yearMonth)

	var lastInvoice models.Invoice
	err := r.db.
		Where("invoice_number LIKE ?", prefix+"%").
		Order("invoice_number DESC").
		First(&lastInvoice).Error

	var sequence int
	if err == nil {
		var lastSeq int
		fmt.Sscanf(lastInvoice.InvoiceNumber[len(prefix):], "%d", &lastSeq)
		sequence = lastSeq + 1
	} else {
		sequence = 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
