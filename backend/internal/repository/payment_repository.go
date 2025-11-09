package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	GetAll(params *models.PaginationParams) ([]models.Payment, int64, error)
	GetByID(id uuid.UUID) (*models.Payment, error)
	GetByStudent(studentID uuid.UUID) ([]models.Payment, error)
	GetByDateRange(start, end time.Time) ([]models.Payment, error)
	Create(payment *models.Payment) error
	Update(payment *models.Payment) error
	Delete(id uuid.UUID) error
	GeneratePaymentNumber(branchCode string, date time.Time) (string, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) GetAll(params *models.PaginationParams) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	query := r.db.Model(&models.Payment{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Invoice").
		Preload("Student").
		Preload("Branch").
		Order("payment_date DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&payments).Error

	return payments, total, err
}

func (r *paymentRepository) GetByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.
		Preload("Invoice").
		Preload("Invoice.Items").
		Preload("Student").
		Preload("Branch").
		First(&payment, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByStudent(studentID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.
		Where("student_id = ?", studentID).
		Preload("Invoice").
		Order("payment_date DESC").
		Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) GetByDateRange(start, end time.Time) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.
		Where("payment_date BETWEEN ? AND ?", start, end).
		Preload("Student").
		Preload("Invoice").
		Order("payment_date ASC").
		Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create payment
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// Update invoice paid amount
		var invoice models.Invoice
		if err := tx.First(&invoice, "id = ?", payment.InvoiceID).Error; err != nil {
			return err
		}

		newPaidAmount := invoice.PaidAmount + payment.Amount

		// Update status
		newStatus := models.InvoiceStatusPartial
		if newPaidAmount >= invoice.TotalAmount {
			newStatus = models.InvoiceStatusPaid
		}

		return tx.Model(&models.Invoice{}).
			Where("id = ?", payment.InvoiceID).
			Updates(map[string]interface{}{
				"paid_amount": newPaidAmount,
				"status":      newStatus,
			}).Error
	})
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get payment
		var payment models.Payment
		if err := tx.First(&payment, "id = ?", id).Error; err != nil {
			return err
		}

		// Update invoice paid amount
		var invoice models.Invoice
		if err := tx.First(&invoice, "id = ?", payment.InvoiceID).Error; err != nil {
			return err
		}

		newPaidAmount := invoice.PaidAmount - payment.Amount
		if newPaidAmount < 0 {
			newPaidAmount = 0
		}

		// Update status
		newStatus := models.InvoiceStatusUnpaid
		if newPaidAmount > 0 {
			newStatus = models.InvoiceStatusPartial
		}

		if err := tx.Model(&models.Invoice{}).
			Where("id = ?", payment.InvoiceID).
			Updates(map[string]interface{}{
				"paid_amount": newPaidAmount,
				"status":      newStatus,
			}).Error; err != nil {
			return err
		}

		// Delete payment
		return tx.Delete(&models.Payment{}, "id = ?", id).Error
	})
}

func (r *paymentRepository) GeneratePaymentNumber(branchCode string, date time.Time) (string, error) {
	// Format: PAY/BranchCode/YYYYMM/XXXX
	yearMonth := date.Format("200601")
	prefix := fmt.Sprintf("PAY/%s/%s/", branchCode, yearMonth)

	var lastPayment models.Payment
	err := r.db.
		Where("payment_number LIKE ?", prefix+"%").
		Order("payment_number DESC").
		First(&lastPayment).Error

	var sequence int
	if err == nil {
		var lastSeq int
		fmt.Sscanf(lastPayment.PaymentNumber[len(prefix):], "%d", &lastSeq)
		sequence = lastSeq + 1
	} else {
		sequence = 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
