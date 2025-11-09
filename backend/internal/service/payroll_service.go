package service

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type PayrollService interface {
	GetAll(params *models.PaginationParams) (*models.PaginationResponse, error)
	GetByID(id uuid.UUID) (*models.Payroll, error)
	GetByPeriod(period string) ([]models.Payroll, error)
	GetByEmployee(employeeID uuid.UUID) ([]models.Payroll, error)
	Create(req *models.CreatePayrollRequest) (*models.Payroll, error)
	ProcessPayroll(payrollID uuid.UUID) (*models.Payroll, error)
	GenerateBulkPayroll(branchID uuid.UUID, period string, paymentDate time.Time) ([]models.Payroll, error)
	CalculateTax(grossSalary float64, maritalStatus string, dependents int) float64
	Delete(id uuid.UUID) error
}

type payrollService struct {
	payrollRepo  repository.PayrollRepository
	employeeRepo repository.EmployeeRepository
	branchRepo   repository.BranchRepository
}

func NewPayrollService(
	payrollRepo repository.PayrollRepository,
	employeeRepo repository.EmployeeRepository,
	branchRepo repository.BranchRepository,
) PayrollService {
	return &payrollService{
		payrollRepo:  payrollRepo,
		employeeRepo: employeeRepo,
		branchRepo:   branchRepo,
	}
}

func (s *payrollService) GetAll(params *models.PaginationParams) (*models.PaginationResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.AppConfig.App.DefaultPageSize
	}

	payrolls, total, err := s.payrollRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(payrolls))
	for i := range payrolls {
		items[i] = payrolls[i]
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

func (s *payrollService) GetByID(id uuid.UUID) (*models.Payroll, error) {
	return s.payrollRepo.GetByID(id)
}

func (s *payrollService) GetByPeriod(period string) ([]models.Payroll, error) {
	return s.payrollRepo.GetByPeriod(period)
}

func (s *payrollService) GetByEmployee(employeeID uuid.UUID) ([]models.Payroll, error) {
	return s.payrollRepo.GetByEmployee(employeeID)
}

func (s *payrollService) Create(req *models.CreatePayrollRequest) (*models.Payroll, error) {
	// Validate employee
	employee, err := s.employeeRepo.GetByID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if employee.Status != models.EmployeeStatusActive {
		return nil, errors.New("employee is not active")
	}

	// Get branch
	branch, err := s.branchRepo.GetByID(employee.BranchID)
	if err != nil {
		return nil, err
	}

	// Generate payroll number
	payrollNumber, err := s.payrollRepo.GeneratePayrollNumber(branch.Code, req.PaymentDate)
	if err != nil {
		return nil, err
	}

	// Calculate gross salary
	grossSalary := req.BaseSalary
	allowances := 0.0
	overtime := 0.0
	bonus := 0.0

	for _, comp := range req.Components {
		// You would fetch component details here
		// For now, we'll just add the amounts
		if comp.Amount > 0 {
			allowances += comp.Amount
		}
	}

	grossSalary += allowances + overtime + bonus

	// Calculate tax (PPh 21)
	// Simplified calculation - in production, use actual marital status and dependents
	dependents := 0
	if employee.MaritalStatus == models.MaritalStatusMarried {
		dependents = 1 // Simplified
	}
	taxPPh21 := s.CalculateTax(grossSalary, employee.MaritalStatus, dependents)

	// Calculate BPJS (simplified - 4% of base salary)
	bpjs := req.BaseSalary * 0.04

	// Total deductions
	totalDeductions := taxPPh21 + bpjs

	// Net salary
	netSalary := grossSalary - totalDeductions

	// Create payroll
	payroll := &models.Payroll{
		EmployeeID:      req.EmployeeID,
		BranchID:        employee.BranchID,
		PayrollNumber:   payrollNumber,
		Period:          req.Period,
		PaymentDate:     req.PaymentDate,
		BaseSalary:      req.BaseSalary,
		Allowances:      allowances,
		Overtime:        overtime,
		Bonus:           bonus,
		TaxPPh21:        taxPPh21,
		BPJS:            bpjs,
		Loan:            0,
		OtherDeductions: 0,
		GrossSalary:     grossSalary,
		TotalDeductions: totalDeductions,
		NetSalary:       netSalary,
		PaymentMethod:   "transfer",
		PaymentStatus:   models.PayrollStatusPending,
		IsPosted:        false,
	}

	if err := s.payrollRepo.Create(payroll); err != nil {
		return nil, err
	}

	return s.payrollRepo.GetByID(payroll.ID)
}

func (s *payrollService) ProcessPayroll(payrollID uuid.UUID) (*models.Payroll, error) {
	payroll, err := s.payrollRepo.GetByID(payrollID)
	if err != nil {
		return nil, errors.New("payroll not found")
	}

	if payroll.PaymentStatus == models.PayrollStatusPaid {
		return nil, errors.New("payroll already paid")
	}

	// Mark as paid
	now := time.Now()
	payroll.PaymentStatus = models.PayrollStatusPaid
	payroll.PaidAt = &now

	// TODO: Create journal entry for payroll expense
	// This would be integrated with the journal system

	if err := s.payrollRepo.Update(payroll); err != nil {
		return nil, err
	}

	return s.payrollRepo.GetByID(payroll.ID)
}

func (s *payrollService) GenerateBulkPayroll(branchID uuid.UUID, period string, paymentDate time.Time) ([]models.Payroll, error) {
	// Get all active employees in branch
	employees, err := s.employeeRepo.GetByBranch(branchID)
	if err != nil {
		return nil, err
	}

	var payrolls []models.Payroll

	for _, employee := range employees {
		// Get latest contract to determine base salary
		// For now, we'll skip employees without contracts
		if len(employee.Contracts) == 0 {
			continue
		}

		// Get active contract
		var activeContract *models.EmployeeContract
		for i := range employee.Contracts {
			if employee.Contracts[i].Status == models.ContractStatusActive {
				activeContract = &employee.Contracts[i]
				break
			}
		}

		if activeContract == nil {
			continue
		}

		// Create payroll request
		req := &models.CreatePayrollRequest{
			EmployeeID:  employee.ID,
			Period:      period,
			PaymentDate: paymentDate,
			BaseSalary:  activeContract.BaseSalary,
			Components:  []models.PayrollComponentReq{},
		}

		payroll, err := s.Create(req)
		if err != nil {
			// Log error but continue with other employees
			continue
		}

		payrolls = append(payrolls, *payroll)
	}

	return payrolls, nil
}

func (s *payrollService) CalculateTax(grossSalary float64, maritalStatus string, dependents int) float64 {
	// PPh 21 Calculation (Indonesia Tax 2024 - Simplified)
	// PTKP (Penghasilan Tidak Kena Pajak) based on marital status
	
	// Annual gross salary
	annualGross := grossSalary * 12
	
	// Calculate PTKP
	ptkp := 54000000.0 // Base PTKP for single (TK/0)
	
	if maritalStatus == models.MaritalStatusMarried {
		ptkp = 58500000.0 // Married base (K/0)
		
		// Add PTKP for dependents (max 3)
		if dependents > 3 {
			dependents = 3
		}
		ptkp += float64(dependents) * 4500000.0
	}
	
	// Taxable income
	taxableIncome := annualGross - ptkp
	
	if taxableIncome <= 0 {
		return 0
	}
	
	// Progressive tax rates (2024)
	// 0-60jt: 5%
	// 60-250jt: 15%
	// 250-500jt: 25%
	// 500-5M: 30%
	// >5M: 35%
	
	var annualTax float64
	
	if taxableIncome <= 60000000 {
		annualTax = taxableIncome * 0.05
	} else if taxableIncome <= 250000000 {
		annualTax = 60000000*0.05 + (taxableIncome-60000000)*0.15
	} else if taxableIncome <= 500000000 {
		annualTax = 60000000*0.05 + 190000000*0.15 + (taxableIncome-250000000)*0.25
	} else if taxableIncome <= 5000000000 {
		annualTax = 60000000*0.05 + 190000000*0.15 + 250000000*0.25 + (taxableIncome-500000000)*0.30
	} else {
		annualTax = 60000000*0.05 + 190000000*0.15 + 250000000*0.25 + 4500000000*0.30 + (taxableIncome-5000000000)*0.35
	}
	
	// Monthly tax
	monthlyTax := annualTax / 12
	
	// Round to nearest rupiah
	return math.Round(monthlyTax)
}

func (s *payrollService) Delete(id uuid.UUID) error {
	payroll, err := s.payrollRepo.GetByID(id)
	if err != nil {
		return errors.New("payroll not found")
	}

	if payroll.PaymentStatus == models.PayrollStatusPaid {
		return errors.New("cannot delete paid payroll")
	}

	if payroll.IsPosted {
		return errors.New("cannot delete posted payroll")
	}

	return s.payrollRepo.Delete(id)
}
