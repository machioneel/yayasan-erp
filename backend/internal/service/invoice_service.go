package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type InvoiceService interface {
	GetAll(params *models.PaginationParams) (*models.PaginationResponse, error)
	GetByID(id uuid.UUID) (*models.Invoice, error)
	GetByStudent(studentID uuid.UUID) ([]models.Invoice, error)
	Create(req *CreateInvoiceRequest) (*models.Invoice, error)
	Update(id uuid.UUID, req *CreateInvoiceRequest) (*models.Invoice, error)
	Delete(id uuid.UUID) error
	GetOverdue() ([]models.Invoice, error)
}

type invoiceService struct {
	invoiceRepo repository.InvoiceRepository
	studentRepo repository.StudentRepository
	branchRepo  repository.BranchRepository
}

func NewInvoiceService(
	invoiceRepo repository.InvoiceRepository,
	studentRepo repository.StudentRepository,
	branchRepo repository.BranchRepository,
) InvoiceService {
	return &invoiceService{
		invoiceRepo: invoiceRepo,
		studentRepo: studentRepo,
		branchRepo:  branchRepo,
	}
}

func (s *invoiceService) GetAll(params *models.PaginationParams) (*models.PaginationResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.GlobalConfig.App.DefaultPageSize
	}

	invoices, total, err := s.invoiceRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(invoices))
	for i := range invoices {
		items[i] = invoices[i]
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

func (s *invoiceService) GetByID(id uuid.UUID) (*models.Invoice, error) {
	return s.invoiceRepo.GetByID(id)
}

func (s *invoiceService) GetByStudent(studentID uuid.UUID) ([]models.Invoice, error) {
	return s.invoiceRepo.GetByStudent(studentID)
}

func (s *invoiceService) Create(req *CreateInvoiceRequest) (*models.Invoice, error) {
	// Validate student
	student, err := s.studentRepo.GetByID(req.StudentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	// Get branch
	branch, err := s.branchRepo.GetByID(student.BranchID)
	if err != nil {
		return nil, err
	}

	// Generate invoice number
	invoiceNumber, err := s.invoiceRepo.GenerateInvoiceNumber(branch.Code, req.InvoiceDate)
	if err != nil {
		return nil, err
	}

	// Create invoice
	invoice := &models.Invoice{
		StudentID:      req.StudentID,
		BranchID:       student.BranchID,
		AcademicYearID: req.AcademicYearID,
		InvoiceNumber:  invoiceNumber,
		InvoiceDate:    req.InvoiceDate,
		DueDate:        req.DueDate,
		Description:    req.Description,
	}

	// Add items
	for _, itemReq := range req.Items {
		invoice.Items = append(invoice.Items, models.InvoiceItem{
			FeeStructureID: itemReq.FeeStructureID,
			Description:    itemReq.Description,
			Quantity:       itemReq.Quantity,
			UnitPrice:      itemReq.UnitPrice,
			AccountID:      itemReq.AccountID,
		})
	}

	if err := s.invoiceRepo.Create(invoice); err != nil {
		return nil, err
	}

	return s.invoiceRepo.GetByID(invoice.ID)
}

func (s *invoiceService) Update(id uuid.UUID, req *CreateInvoiceRequest) (*models.Invoice, error) {
	invoice, err := s.invoiceRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	// Cannot update paid invoice
	if invoice.Status == models.InvoiceStatusPaid {
		return nil, errors.New("cannot update paid invoice")
	}

	// Update invoice
	invoice.InvoiceDate = req.InvoiceDate
	invoice.DueDate = req.DueDate
	invoice.Description = req.Description

	// Update items
	invoice.Items = []models.InvoiceItem{}
	for _, itemReq := range req.Items {
		invoice.Items = append(invoice.Items, models.InvoiceItem{
			InvoiceID:      invoice.ID,
			FeeStructureID: itemReq.FeeStructureID,
			Description:    itemReq.Description,
			Quantity:       itemReq.Quantity,
			UnitPrice:      itemReq.UnitPrice,
			AccountID:      itemReq.AccountID,
		})
	}

	if err := s.invoiceRepo.Update(invoice); err != nil {
		return nil, err
	}

	return s.invoiceRepo.GetByID(invoice.ID)
}

func (s *invoiceService) Delete(id uuid.UUID) error {
	return s.invoiceRepo.Delete(id)
}

func (s *invoiceService) GetOverdue() ([]models.Invoice, error) {
	return s.invoiceRepo.GetOverdue()
}

// CreateInvoiceRequest for creating invoice
type CreateInvoiceRequest struct {
	StudentID      uuid.UUID          `json:"student_id" binding:"required"`
	AcademicYearID uuid.UUID          `json:"academic_year_id" binding:"required"`
	InvoiceDate    time.Time          `json:"invoice_date" binding:"required"`
	DueDate        time.Time          `json:"due_date" binding:"required"`
	Description    string             `json:"description" binding:"required"`
	Items          []InvoiceItemRequest `json:"items" binding:"required,min=1"`
}

// InvoiceItemRequest for invoice item
type InvoiceItemRequest struct {
	FeeStructureID *uuid.UUID `json:"fee_structure_id"`
	Description    string     `json:"description" binding:"required"`
	Quantity       int        `json:"quantity" binding:"required,min=1"`
	UnitPrice      float64    `json:"unit_price" binding:"required,gt=0"`
	AccountID      *uuid.UUID `json:"account_id"`
}
