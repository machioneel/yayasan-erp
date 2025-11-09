package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type EmployeeService interface {
	GetAll(params *models.PaginationParams) (*models.PaginationResponse, error)
	GetByID(id uuid.UUID) (*models.Employee, error)
	GetByBranch(branchID uuid.UUID) ([]models.Employee, error)
	GetTeachers() ([]models.Employee, error)
	Search(keyword string) ([]models.Employee, error)
	Create(req *models.CreateEmployeeRequest) (*models.Employee, error)
	Update(id uuid.UUID, req *models.CreateEmployeeRequest) (*models.Employee, error)
	Delete(id uuid.UUID) error
	GetStatistics() (map[string]interface{}, error)
}

type employeeService struct {
	employeeRepo repository.EmployeeRepository
	branchRepo   repository.BranchRepository
}

func NewEmployeeService(
	employeeRepo repository.EmployeeRepository,
	branchRepo repository.BranchRepository,
) EmployeeService {
	return &employeeService{
		employeeRepo: employeeRepo,
		branchRepo:   branchRepo,
	}
}

func (s *employeeService) GetAll(params *models.PaginationParams) (*models.PaginationResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.GlobalConfig.App.DefaultPageSize
	}

	employees, total, err := s.employeeRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(employees))
	for i := range employees {
		items[i] = employees[i]
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

func (s *employeeService) GetByID(id uuid.UUID) (*models.Employee, error) {
	return s.employeeRepo.GetByID(id)
}

func (s *employeeService) GetByBranch(branchID uuid.UUID) ([]models.Employee, error) {
	return s.employeeRepo.GetByBranch(branchID)
}

func (s *employeeService) GetTeachers() ([]models.Employee, error) {
	return s.employeeRepo.GetTeachers()
}

func (s *employeeService) Search(keyword string) ([]models.Employee, error) {
	return s.employeeRepo.Search(keyword)
}

func (s *employeeService) Create(req *models.CreateEmployeeRequest) (*models.Employee, error) {
	// Validate branch
	branch, err := s.branchRepo.GetByID(req.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	// Check NIK uniqueness if provided
	if req.NIK != "" {
		existing, _ := s.employeeRepo.GetByNIK(req.NIK)
		if existing != nil {
			return nil, errors.New("NIK already exists")
		}
	}

	// Generate employee number
	year := req.JoinDate.Year()
	employeeNumber, err := s.employeeRepo.GenerateEmployeeNumber(branch.Code, year)
	if err != nil {
		return nil, err
	}

	// Create employee
	employee := &models.Employee{
		EmployeeNumber: employeeNumber,
		NIK:            req.NIK,
		NPWP:           req.NPWP,
		FullName:       req.FullName,
		Gender:         req.Gender,
		BirthPlace:     req.BirthPlace,
		BirthDate:      req.BirthDate,
		Religion:       req.Religion,
		MaritalStatus:  req.MaritalStatus,
		Email:          req.Email,
		Phone:          req.Phone,
		Address:        req.Address,
		City:           req.City,
		BranchID:       req.BranchID,
		Department:     req.Department,
		Position:       req.Position,
		EmploymentType: req.EmploymentType,
		JoinDate:       req.JoinDate,
		Status:         models.EmployeeStatusActive,
		IsTeacher:      req.IsTeacher,
		Education:      req.Education,
		BankName:       req.BankName,
		BankAccount:    req.BankAccount,
	}

	if err := s.employeeRepo.Create(employee); err != nil {
		return nil, err
	}

	return s.employeeRepo.GetByID(employee.ID)
}

func (s *employeeService) Update(id uuid.UUID, req *models.CreateEmployeeRequest) (*models.Employee, error) {
	employee, err := s.employeeRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Update fields
	employee.FullName = req.FullName
	employee.Gender = req.Gender
	employee.BirthPlace = req.BirthPlace
	employee.BirthDate = req.BirthDate
	employee.Religion = req.Religion
	employee.MaritalStatus = req.MaritalStatus
	employee.Email = req.Email
	employee.Phone = req.Phone
	employee.Address = req.Address
	employee.City = req.City
	employee.Department = req.Department
	employee.Position = req.Position
	employee.EmploymentType = req.EmploymentType
	employee.IsTeacher = req.IsTeacher
	employee.Education = req.Education
	employee.BankName = req.BankName
	employee.BankAccount = req.BankAccount

	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, err
	}

	return s.employeeRepo.GetByID(employee.ID)
}

func (s *employeeService) Delete(id uuid.UUID) error {
	employee, err := s.employeeRepo.GetByID(id)
	if err != nil {
		return errors.New("employee not found")
	}

	// Soft delete - change status to resigned
	employee.Status = models.EmployeeStatusResigned
	now := time.Now()
	employee.EndDate = &now

	return s.employeeRepo.Update(employee)
}

func (s *employeeService) GetStatistics() (map[string]interface{}, error) {
	totalActive, _ := s.employeeRepo.CountByStatus(models.EmployeeStatusActive)
	totalResigned, _ := s.employeeRepo.CountByStatus(models.EmployeeStatusResigned)
	totalTerminated, _ := s.employeeRepo.CountByStatus(models.EmployeeStatusTerminated)

	teachers, _ := s.employeeRepo.GetTeachers()
	totalTeachers := int64(len(teachers))

	return map[string]interface{}{
		"total_active":     totalActive,
		"total_resigned":   totalResigned,
		"total_terminated": totalTerminated,
		"total_teachers":   totalTeachers,
		"total":            totalActive + totalResigned + totalTerminated,
	}, nil
}
