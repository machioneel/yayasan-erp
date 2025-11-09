package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type BranchService interface {
	GetAll(params *models.PaginationParams) ([]models.Branch, int64, error)
	GetByID(id uuid.UUID) (*models.Branch, error)
	GetByCode(code string) (*models.Branch, error)
	GetAllActive() ([]models.Branch, error)
	Create(req *models.CreateBranchRequest, createdBy uuid.UUID) (*models.Branch, error)
	Update(id uuid.UUID, req *models.UpdateBranchRequest) (*models.Branch, error)
	Delete(id uuid.UUID) error
}

type branchService struct {
	branchRepo repository.BranchRepository
}

func NewBranchService(branchRepo repository.BranchRepository) BranchService {
	return &branchService{
		branchRepo: branchRepo,
	}
}

func (s *branchService) GetAll(params *models.PaginationParams) ([]models.Branch, int64, error) {
	// Set default pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.AppConfig.App.DefaultPageSize
	}
	if params.PageSize > config.AppConfig.App.MaxPageSize {
		params.PageSize = config.AppConfig.App.MaxPageSize
	}

	return s.branchRepo.GetAll(params)
}

func (s *branchService) GetByID(id uuid.UUID) (*models.Branch, error) {
	return s.branchRepo.GetByID(id)
}

func (s *branchService) GetByCode(code string) (*models.Branch, error) {
	return s.branchRepo.GetByCode(code)
}

func (s *branchService) GetAllActive() ([]models.Branch, error) {
	return s.branchRepo.GetAllActive()
}

func (s *branchService) Create(req *models.CreateBranchRequest, createdBy uuid.UUID) (*models.Branch, error) {
	// Check if code exists
	existing, _ := s.branchRepo.GetByCode(req.Code)
	if existing != nil {
		return nil, errors.New("branch code already exists")
	}

	branch := &models.Branch{
		Code:       req.Code,
		Name:       req.Name,
		Type:       req.Type,
		Address:    req.Address,
		City:       req.City,
		Province:   req.Province,
		PostalCode: req.PostalCode,
		Phone:      req.Phone,
		Email:      req.Email,
		IsActive:   true,
	}
	branch.CreatedBy = &createdBy

	if err := s.branchRepo.Create(branch); err != nil {
		return nil, err
	}

	return branch, nil
}

func (s *branchService) Update(id uuid.UUID, req *models.UpdateBranchRequest) (*models.Branch, error) {
	branch, err := s.branchRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	branch.Name = req.Name
	branch.Type = req.Type
	branch.Address = req.Address
	branch.City = req.City
	branch.Province = req.Province
	branch.PostalCode = req.PostalCode
	branch.Phone = req.Phone
	branch.Email = req.Email

	if req.IsActive != nil {
		branch.IsActive = *req.IsActive
	}

	if err := s.branchRepo.Update(branch); err != nil {
		return nil, err
	}

	return branch, nil
}

func (s *branchService) Delete(id uuid.UUID) error {
	return s.branchRepo.Delete(id)
}
