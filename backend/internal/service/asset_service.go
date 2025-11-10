package service

import (
	"errors"
	"math"
	//"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type AssetService interface {
	GetAll(params *models.PaginationParams) (*models.PaginationResponse, error)
	GetByID(id uuid.UUID) (*models.Asset, error)
	GetByBranch(branchID uuid.UUID) ([]models.Asset, error)
	Create(req *models.CreateAssetRequest) (*models.Asset, error)
	Update(id uuid.UUID, req *models.CreateAssetRequest) (*models.Asset, error)
	Delete(id uuid.UUID) error
	
	
	// Depreciation
	CalculateDepreciation(assetID uuid.UUID, period string) error
	ProcessMonthlyDepreciation(period string) ([]models.AssetDepreciation, error)
	
	// Maintenance
	CreateMaintenance(req *models.CreateMaintenanceRequest) (*models.AssetMaintenance, error)
	GetUpcomingMaintenance() ([]models.AssetMaintenance, error)
	
	// Transfer
	CreateTransfer(req *models.CreateTransferRequest, userID uuid.UUID) (*models.AssetTransfer, error)
	ApproveTransfer(transferID, userID uuid.UUID) error
	
	// Statistics
	GetStatistics() (map[string]interface{}, error)
}

type assetService struct {
	assetRepo  repository.AssetRepository
	branchRepo repository.BranchRepository
}

func NewAssetService(assetRepo repository.AssetRepository, branchRepo repository.BranchRepository) AssetService {
	return &assetService{
		assetRepo:  assetRepo,
		branchRepo: branchRepo,
	}
}

func (s *assetService) GetAll(params *models.PaginationParams) (*models.PaginationResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.GlobalConfig.App.DefaultPageSize
	}

	assets, total, err := s.assetRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(assets))
	for i := range assets {
		items[i] = assets[i]
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

func (s *assetService) GetByID(id uuid.UUID) (*models.Asset, error) {
	return s.assetRepo.GetByID(id)
}

func (s *assetService) GetByBranch(branchID uuid.UUID) ([]models.Asset, error) {
	return s.assetRepo.GetByBranch(branchID)
}

func (s *assetService) Create(req *models.CreateAssetRequest) (*models.Asset, error) {
	branch, err := s.branchRepo.GetByID(req.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	year := req.PurchaseDate.Year()
	assetNumber, err := s.assetRepo.GenerateAssetNumber(branch.Code, year)
	if err != nil {
		return nil, err
	}

	// Calculate initial book value
	bookValue := req.PurchasePrice

	asset := &models.Asset{
		AssetNumber:        assetNumber,
		Name:               req.Name,
		CategoryID:         req.CategoryID,
		BranchID:           req.BranchID,
		Description:        req.Description,
		Brand:              req.Brand,
		SerialNumber:       req.SerialNumber,
		PurchaseDate:       req.PurchaseDate,
		PurchasePrice:      req.PurchasePrice,
		DepreciationMethod: req.DepreciationMethod,
		UsefulLife:         req.UsefulLife,
		SalvageValue:       req.SalvageValue,
		AccumulatedDepreciation: 0,
		BookValue:          bookValue,
		Location:           req.Location,
		ResponsibleUser:    req.ResponsibleUser,
		Condition:          models.AssetConditionGood,
		Status:             models.AssetStatusActive,
	}

	if err := s.assetRepo.Create(asset); err != nil {
		return nil, err
	}

	return s.assetRepo.GetByID(asset.ID)
}

func (s *assetService) Update(id uuid.UUID, req *models.CreateAssetRequest) (*models.Asset, error) {
	asset, err := s.assetRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	asset.Name = req.Name
	asset.Description = req.Description
	asset.Brand = req.Brand
	asset.Location = req.Location
	asset.ResponsibleUser = req.ResponsibleUser

	if err := s.assetRepo.Update(asset); err != nil {
		return nil, err
	}

	return s.assetRepo.GetByID(asset.ID)
}

func (s *assetService) Delete(id uuid.UUID) error {
	return s.assetRepo.Delete(id)
}

func (s *assetService) CalculateDepreciation(assetID uuid.UUID, period string) error {
	asset, err := s.assetRepo.GetByID(assetID)
	if err != nil {
		return err
	}

	if asset.UsefulLife == 0 {
		return errors.New("useful life not set")
	}

	var depreciationAmount float64

	if asset.DepreciationMethod == models.DepreciationMethodStraightLine {
		// Straight-line: (Cost - Salvage) / Useful Life / 12
		depreciationAmount = (asset.PurchasePrice - asset.SalvageValue) / float64(asset.UsefulLife) / 12
	} else if asset.DepreciationMethod == models.DepreciationMethodDecliningBalance {
		// Declining balance: Book Value × (2 / Useful Life) / 12
		rate := 2.0 / float64(asset.UsefulLife)
		depreciationAmount = asset.BookValue * rate / 12
	}

	depreciationAmount = math.Round(depreciationAmount)

	newAccumulated := asset.AccumulatedDepreciation + depreciationAmount
	newBookValue := asset.PurchasePrice - newAccumulated

	// Create depreciation record
	depreciation := &models.AssetDepreciation{
		AssetID:            assetID,
		Period:             period,
		DepreciationAmount: depreciationAmount,
		AccumulatedAmount:  newAccumulated,
		BookValue:          newBookValue,
	}

	if err := s.assetRepo.CreateDepreciation(depreciation); err != nil {
		return err
	}

	// Update asset
	asset.AccumulatedDepreciation = newAccumulated
	asset.BookValue = newBookValue

	return s.assetRepo.Update(asset)
}

func (s *assetService) ProcessMonthlyDepreciation(period string) ([]models.AssetDepreciation, error) {
	assets, err := s.assetRepo.GetByStatus(models.AssetStatusActive)
	if err != nil {
		return nil, err
	}

	var depreciations []models.AssetDepreciation

	for _, asset := range assets {
		if asset.UsefulLife > 0 && asset.DepreciationMethod != "" {
			if err := s.CalculateDepreciation(asset.ID, period); err != nil {
				continue
			}
		}
	}

	return depreciations, nil
}

func (s *assetService) CreateMaintenance(req *models.CreateMaintenanceRequest) (*models.AssetMaintenance, error) {
	maintenance := &models.AssetMaintenance{
		AssetID:         req.AssetID,
		MaintenanceType: req.MaintenanceType,
		MaintenanceDate: req.MaintenanceDate,
		Description:     req.Description,
		Cost:            req.Cost,
		Vendor:          req.Vendor,
		NextScheduled:   req.NextScheduled,
		Status:          models.MaintenanceStatusCompleted,
	}

	if err := s.assetRepo.CreateMaintenance(maintenance); err != nil {
		return nil, err
	}

	return maintenance, nil
}

func (s *assetService) GetUpcomingMaintenance() ([]models.AssetMaintenance, error) {
	return s.assetRepo.GetUpcomingMaintenance()
}

func (s *assetService) CreateTransfer(req *models.CreateTransferRequest, userID uuid.UUID) (*models.AssetTransfer, error) {
	asset, err := s.assetRepo.GetByID(req.AssetID)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	transfer := &models.AssetTransfer{
		AssetID:      req.AssetID,
		FromBranchID: asset.BranchID,
		ToBranchID:   req.ToBranchID,
		TransferDate: req.TransferDate,
		Reason:       req.Reason,
		RequestedBy:  userID,
		Status:       models.TransferStatusPending,
	}

	if err := s.assetRepo.CreateTransfer(transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *assetService) ApproveTransfer(transferID, userID uuid.UUID) error {
	// Implementation would include updating transfer status and asset branch
	return s.assetRepo.UpdateTransferStatus(transferID, models.TransferStatusApproved)
}

func (s *assetService) GetStatistics() (map[string]interface{}, error) {
	totalActive, _ := s.assetRepo.CountByStatus(models.AssetStatusActive)
	totalValue, _ := s.assetRepo.GetTotalValue()
	totalBookValue, _ := s.assetRepo.GetTotalBookValue()

	return map[string]interface{}{
		"total_assets":    totalActive,
		"total_value":     totalValue,
		"total_book_value": totalBookValue,
		"depreciation":    totalValue - totalBookValue,
	}, nil
}
