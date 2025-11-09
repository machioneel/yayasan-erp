package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type AssetRepository interface {
	GetAll(params *models.PaginationParams) ([]models.Asset, int64, error)
	GetByID(id uuid.UUID) (*models.Asset, error)
	GetByBranch(branchID uuid.UUID) ([]models.Asset, error)
	GetByCategory(categoryID uuid.UUID) ([]models.Asset, error)
	GetByStatus(status string) ([]models.Asset, error)
	Search(keyword string) ([]models.Asset, error)
	Create(asset *models.Asset) error
	Update(asset *models.Asset) error
	Delete(id uuid.UUID) error
	GenerateAssetNumber(branchCode string, year int) (string, error)
	
	// Depreciation
	CreateDepreciation(depreciation *models.AssetDepreciation) error
	GetDepreciationsByPeriod(period string) ([]models.AssetDepreciation, error)
	GetDepreciationsByAsset(assetID uuid.UUID) ([]models.AssetDepreciation, error)
	
	// Maintenance
	CreateMaintenance(maintenance *models.AssetMaintenance) error
	GetMaintenancesByAsset(assetID uuid.UUID) ([]models.AssetMaintenance, error)
	GetUpcomingMaintenance() ([]models.AssetMaintenance, error)
	
	// Transfer
	CreateTransfer(transfer *models.AssetTransfer) error
	GetTransfersByAsset(assetID uuid.UUID) ([]models.AssetTransfer, error)
	UpdateTransferStatus(id uuid.UUID, status string) error
	
	// Statistics
	CountByStatus(status string) (int64, error)
	GetTotalValue() (float64, error)
	GetTotalBookValue() (float64, error)
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) GetAll(params *models.PaginationParams) ([]models.Asset, int64, error) {
	var assets []models.Asset
	var total int64

	query := r.db.Model(&models.Asset{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("name LIKE ? OR asset_number LIKE ? OR serial_number LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Category").
		Preload("Branch").
		Preload("ResponsiblePerson").
		Order("created_at DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&assets).Error

	return assets, total, err
}

func (r *assetRepository) GetByID(id uuid.UUID) (*models.Asset, error) {
	var asset models.Asset
	err := r.db.
		Preload("Category").
		Preload("Branch").
		Preload("ResponsiblePerson").
		Preload("Depreciations").
		Preload("Maintenances").
		Preload("Transfers").
		First(&asset, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("asset not found")
		}
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) GetByBranch(branchID uuid.UUID) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.
		Where("branch_id = ? AND status = ?", branchID, models.AssetStatusActive).
		Preload("Category").
		Order("name ASC").
		Find(&assets).Error
	return assets, err
}

func (r *assetRepository) GetByCategory(categoryID uuid.UUID) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.
		Where("category_id = ?", categoryID).
		Preload("Branch").
		Order("name ASC").
		Find(&assets).Error
	return assets, err
}

func (r *assetRepository) GetByStatus(status string) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.
		Where("status = ?", status).
		Preload("Category").
		Preload("Branch").
		Order("name ASC").
		Find(&assets).Error
	return assets, err
}

func (r *assetRepository) Search(keyword string) ([]models.Asset, error) {
	var assets []models.Asset
	searchPattern := "%" + keyword + "%"
	err := r.db.
		Where("name LIKE ? OR asset_number LIKE ? OR serial_number LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Limit(20).
		Find(&assets).Error
	return assets, err
}

func (r *assetRepository) Create(asset *models.Asset) error {
	return r.db.Create(asset).Error
}

func (r *assetRepository) Update(asset *models.Asset) error {
	return r.db.Save(asset).Error
}

func (r *assetRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Asset{}, "id = ?", id).Error
}

func (r *assetRepository) GenerateAssetNumber(branchCode string, year int) (string, error) {
	// Format: AST/BRANCH/YYYY/XXXX
	prefix := fmt.Sprintf("AST/%s/%d/", branchCode, year)

	var lastAsset models.Asset
	err := r.db.
		Where("asset_number LIKE ?", prefix+"%").
		Order("asset_number DESC").
		First(&lastAsset).Error

	var sequence int
	if err == nil {
		var lastSeq int
		fmt.Sscanf(lastAsset.AssetNumber[len(prefix):], "%d", &lastSeq)
		sequence = lastSeq + 1
	} else {
		sequence = 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}

// Depreciation methods
func (r *assetRepository) CreateDepreciation(depreciation *models.AssetDepreciation) error {
	return r.db.Create(depreciation).Error
}

func (r *assetRepository) GetDepreciationsByPeriod(period string) ([]models.AssetDepreciation, error) {
	var depreciations []models.AssetDepreciation
	err := r.db.
		Where("period = ?", period).
		Preload("Asset").
		Find(&depreciations).Error
	return depreciations, err
}

func (r *assetRepository) GetDepreciationsByAsset(assetID uuid.UUID) ([]models.AssetDepreciation, error) {
	var depreciations []models.AssetDepreciation
	err := r.db.
		Where("asset_id = ?", assetID).
		Order("period DESC").
		Find(&depreciations).Error
	return depreciations, err
}

// Maintenance methods
func (r *assetRepository) CreateMaintenance(maintenance *models.AssetMaintenance) error {
	return r.db.Create(maintenance).Error
}

func (r *assetRepository) GetMaintenancesByAsset(assetID uuid.UUID) ([]models.AssetMaintenance, error) {
	var maintenances []models.AssetMaintenance
	err := r.db.
		Where("asset_id = ?", assetID).
		Order("maintenance_date DESC").
		Find(&maintenances).Error
	return maintenances, err
}

func (r *assetRepository) GetUpcomingMaintenance() ([]models.AssetMaintenance, error) {
	var maintenances []models.AssetMaintenance
	now := time.Now()
	err := r.db.
		Where("next_scheduled <= ? AND status = ?", now.AddDate(0, 0, 30), models.MaintenanceStatusScheduled).
		Preload("Asset").
		Order("next_scheduled ASC").
		Find(&maintenances).Error
	return maintenances, err
}

// Transfer methods
func (r *assetRepository) CreateTransfer(transfer *models.AssetTransfer) error {
	return r.db.Create(transfer).Error
}

func (r *assetRepository) GetTransfersByAsset(assetID uuid.UUID) ([]models.AssetTransfer, error) {
	var transfers []models.AssetTransfer
	err := r.db.
		Where("asset_id = ?", assetID).
		Preload("FromBranch").
		Preload("ToBranch").
		Order("transfer_date DESC").
		Find(&transfers).Error
	return transfers, err
}

func (r *assetRepository) UpdateTransferStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.AssetTransfer{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Statistics methods
func (r *assetRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Asset{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

func (r *assetRepository) GetTotalValue() (float64, error) {
	var total float64
	err := r.db.Model(&models.Asset{}).
		Where("status = ?", models.AssetStatusActive).
		Select("COALESCE(SUM(purchase_price), 0)").
		Scan(&total).Error
	return total, err
}

func (r *assetRepository) GetTotalBookValue() (float64, error) {
	var total float64
	err := r.db.Model(&models.Asset{}).
		Where("status = ?", models.AssetStatusActive).
		Select("COALESCE(SUM(book_value), 0)").
		Scan(&total).Error
	return total, err
}
