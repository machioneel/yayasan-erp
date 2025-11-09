package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type PaymentService interface {
	GetAll(params *models.PaginationParams) (*models.PaginationResponse, error)
	GetByID(id uuid.UUID) (*models.Payment, error)
	GetByStudent(studentID uuid.UUID) ([]models.Payment, error)
	Create(req *CreatePaymentRequest, userID uuid.UUID) (*models.Payment, error)
	Post(id uuid.UUID, userID uuid.UUID) (*models.Payment, error)
	Delete(id uuid.UUID) error
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	invoiceRepo repository.InvoiceRepository
	branchRepo  repository.BranchRepository
	studentRepo repository.StudentRepository
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	invoiceRepo repository.InvoiceRepository,
	branchRepo repository.BranchRepository,
	studentRepo repository.StudentRepository,
) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
		branchRepo:  branchRepo,
		studentRepo: studentRepo,
	}
}

func (s *paymentService) GetAll(params *models.PaginationParams) (*models.PaginationResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.AppConfig.App.DefaultPageSize
	}

	payments, total, err := s.paymentRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(payments))
	for i := range payments {
		items[i] = payments[i]
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &models.PaginationResponse{
		Data:       items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *paymentService) GetByID(id uuid.UUID) (*models.Payment, error) {
	return s.paymentRepo.GetByID(id)
}

func (s *paymentService) GetByStudent(studentID uuid.UUID) ([]models.Payment, error) {
	return s.paymentRepo.GetByStudent(studentID)
}

func (s *paymentService) Create(req *CreatePaymentRequest, userID uuid.UUID) (*models.Payment, error) {
	// Validate invoice
	invoice, err := s.invoiceRepo.GetByID(req.InvoiceID)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	// Check if invoice is already fully paid
	if invoice.Status == models.InvoiceStatusPaid {
		return nil, errors.New("invoice is already fully paid")
	}

	// Check payment amount
	remaining := invoice.TotalAmount - invoice.PaidAmount
	if req.Amount > remaining {
		return nil, errors.New("payment amount exceeds remaining balance")
	}

	// Get branch
	branch, err := s.branchRepo.GetByID(invoice.BranchID)
	if err != nil {
		return nil, err
	}

	// Generate payment number
	paymentNumber, err := s.paymentRepo.GeneratePaymentNumber(branch.Code, req.PaymentDate)
	if err != nil {
		return nil, err
	}

	// Create payment
	payment := &models.Payment{
		InvoiceID:     req.InvoiceID,
		StudentID:     invoice.StudentID,
		BranchID:      invoice.BranchID,
		PaymentNumber: paymentNumber,
		PaymentDate:   req.PaymentDate,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		ReferenceNo:   req.ReferenceNo,
		Notes:         req.Notes,
		ReceivedBy:    userID,
		IsPosted:      false,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	return s.paymentRepo.GetByID(payment.ID)
}

func (s *paymentService) Post(id uuid.UUID, userID uuid.UUID) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	if payment.IsPosted {
		return nil, errors.New("payment is already posted")
	}

	// TODO: Create journal entry for payment
	// This will be implemented when journal integration is needed

	now := time.Now()
	payment.IsPosted = true
	payment.PostedAt = &now

	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, err
	}

	return s.paymentRepo.GetByID(payment.ID)
}

func (s *paymentService) Delete(id uuid.UUID) error {
	payment, err := s.paymentRepo.GetByID(id)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.IsPosted {
		return errors.New("cannot delete posted payment")
	}

	return s.paymentRepo.Delete(id)
}

// CreatePaymentRequest for creating payment
type CreatePaymentRequest struct {
	InvoiceID     uuid.UUID `json:"invoice_id" binding:"required"`
	PaymentDate   time.Time `json:"payment_date" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,gt=0"`
	PaymentMethod string    `json:"payment_method" binding:"required"`
	ReferenceNo   string    `json:"reference_no"`
	Notes         string    `json:"notes"`
}
