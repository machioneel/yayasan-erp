package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	// Items
	GetAllItems(params *models.PaginationParams) ([]models.InventoryItem, int64, error)
	GetItemByID(id uuid.UUID) (*models.InventoryItem, error)
	GetItemsByBranch(branchID uuid.UUID) ([]models.InventoryItem, error)
	GetItemsByCategory(categoryID uuid.UUID) ([]models.InventoryItem, error)
	GetLowStockItems(branchID uuid.UUID) ([]models.InventoryItem, error)
	SearchItems(keyword string) ([]models.InventoryItem, error)
	CreateItem(item *models.InventoryItem) error
	UpdateItem(item *models.InventoryItem) error
	DeleteItem(id uuid.UUID) error
	GenerateItemCode(categoryCode string) (string, error)
	
	// Stock Transactions
	GetAllTransactions(params *models.PaginationParams) ([]models.StockTransaction, int64, error)
	GetTransactionByID(id uuid.UUID) (*models.StockTransaction, error)
	GetTransactionsByItem(itemID uuid.UUID) ([]models.StockTransaction, error)
	GetTransactionsByDateRange(startDate, endDate time.Time) ([]models.StockTransaction, error)
	CreateTransaction(transaction *models.StockTransaction) error
	GenerateTransactionNumber(txType, branchCode string, date time.Time) (string, error)
	
	// Stock Opname
	GetAllOpnames(params *models.PaginationParams) ([]models.StockOpname, int64, error)
	GetOpnameByID(id uuid.UUID) (*models.StockOpname, error)
	CreateOpname(opname *models.StockOpname) error
	UpdateOpname(opname *models.StockOpname) error
	ApproveOpname(id uuid.UUID, approverID uuid.UUID) error
	GenerateOpnameNumber(branchCode string, date time.Time) (string, error)
	
	// Purchase Orders
	GetAllPOs(params *models.PaginationParams) ([]models.PurchaseOrder, int64, error)
	GetPOByID(id uuid.UUID) (*models.PurchaseOrder, error)
	CreatePO(po *models.PurchaseOrder) error
	UpdatePO(po *models.PurchaseOrder) error
	ApprovePO(id uuid.UUID, approverID uuid.UUID) error
	GeneratePONumber(branchCode string, date time.Time) (string, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// Items
func (r *inventoryRepository) GetAllItems(params *models.PaginationParams) ([]models.InventoryItem, int64, error) {
	var items []models.InventoryItem
	var total int64

	query := r.db.Model(&models.InventoryItem{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("name LIKE ? OR item_code LIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Category").
		Preload("Branch").
		Order("item_code ASC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&items).Error

	return items, total, err
}

func (r *inventoryRepository) GetItemByID(id uuid.UUID) (*models.InventoryItem, error) {
	var item models.InventoryItem
	err := r.db.
		Preload("Category").
		Preload("Branch").
		First(&item, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}
	return &item, nil
}

func (r *inventoryRepository) GetItemsByBranch(branchID uuid.UUID) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := r.db.
		Where("branch_id = ? AND is_active = ?", branchID, true).
		Preload("Category").
		Order("name ASC").
		Find(&items).Error
	return items, err
}

func (r *inventoryRepository) GetItemsByCategory(categoryID uuid.UUID) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := r.db.
		Where("category_id = ?", categoryID).
		Preload("Branch").
		Order("name ASC").
		Find(&items).Error
	return items, err
}

func (r *inventoryRepository) GetLowStockItems(branchID uuid.UUID) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	query := r.db.Where("current_stock <= minimum_stock AND is_active = ?", true)
	
	if branchID != uuid.Nil {
		query = query.Where("branch_id = ?", branchID)
	}
	
	err := query.
		Preload("Category").
		Preload("Branch").
		Order("current_stock ASC").
		Find(&items).Error
	return items, err
}

func (r *inventoryRepository) SearchItems(keyword string) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	searchPattern := "%" + keyword + "%"
	err := r.db.
		Where("name LIKE ? OR item_code LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern).
		Preload("Category").
		Limit(20).
		Find(&items).Error
	return items, err
}

func (r *inventoryRepository) CreateItem(item *models.InventoryItem) error {
	return r.db.Create(item).Error
}

func (r *inventoryRepository) UpdateItem(item *models.InventoryItem) error {
	return r.db.Save(item).Error
}

func (r *inventoryRepository) DeleteItem(id uuid.UUID) error {
	return r.db.Delete(&models.InventoryItem{}, "id = ?", id).Error
}

func (r *inventoryRepository) GenerateItemCode(categoryCode string) (string, error) {
	// Format: ITM/CATEGORY/XXXX
	prefix := fmt.Sprintf("ITM/%s/", categoryCode)

	var lastItem models.InventoryItem
	err := r.db.
		Where("item_code LIKE ?", prefix+"%").
		Order("item_code DESC").
		First(&lastItem).Error

	var sequence int
	if err == nil {
		var lastSeq int
		fmt.Sscanf(lastItem.ItemCode[len(prefix):], "%d", &lastSeq)
		sequence = lastSeq + 1
	} else {
		sequence = 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}

// Stock Transactions
func (r *inventoryRepository) GetAllTransactions(params *models.PaginationParams) ([]models.StockTransaction, int64, error) {
	var transactions []models.StockTransaction
	var total int64

	query := r.db.Model(&models.StockTransaction{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Item").
		Preload("Branch").
		Preload("User").
		Order("transaction_date DESC, created_at DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *inventoryRepository) GetTransactionByID(id uuid.UUID) (*models.StockTransaction, error) {
	var transaction models.StockTransaction
	err := r.db.
		Preload("Item").
		Preload("Branch").
		Preload("User").
		First(&transaction, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *inventoryRepository) GetTransactionsByItem(itemID uuid.UUID) ([]models.StockTransaction, error) {
	var transactions []models.StockTransaction
	err := r.db.
		Where("item_id = ?", itemID).
		Preload("User").
		Order("transaction_date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *inventoryRepository) GetTransactionsByDateRange(startDate, endDate time.Time) ([]models.StockTransaction, error) {
	var transactions []models.StockTransaction
	err := r.db.
		Where("transaction_date BETWEEN ? AND ?", startDate, endDate).
		Preload("Item").
		Preload("Branch").
		Order("transaction_date ASC").
		Find(&transactions).Error
	return transactions, err
}

func (r *inventoryRepository) CreateTransaction(transaction *models.StockTransaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get current item
		var item models.InventoryItem
		if err := tx.First(&item, "id = ?", transaction.ItemID).Error; err != nil {
			return err
		}

		// Set stock before
		transaction.StockBefore = item.CurrentStock

		// Calculate new stock
		var newStock float64
		switch transaction.TransactionType {
		case models.TransactionTypeIn:
			newStock = item.CurrentStock + transaction.Quantity
		case models.TransactionTypeOut:
			newStock = item.CurrentStock - transaction.Quantity
			if newStock < 0 {
				return errors.New("insufficient stock")
			}
		case models.TransactionTypeAdjustment, models.TransactionTypeOpname:
			// For adjustment, quantity can be positive or negative
			newStock = item.CurrentStock + transaction.Quantity
			if newStock < 0 {
				newStock = 0
			}
		}

		transaction.StockAfter = newStock

		// Create transaction
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// Update item stock
		return tx.Model(&models.InventoryItem{}).
			Where("id = ?", transaction.ItemID).
			Update("current_stock", newStock).
			Error
	})
}

func (r *inventoryRepository) GenerateTransactionNumber(txType, branchCode string, date time.Time) (string, error) {
	// Format: STK/IN|OUT|ADJ/BRANCH/YYYYMMDD/XXXX
	dateStr := date.Format("20060102")
	typeCode := "IN"
	switch txType {
	case models.TransactionTypeOut:
		typeCode = "OUT"
	case models.TransactionTypeAdjustment:
		typeCode = "ADJ"
	case models.TransactionTypeOpname:
		typeCode = "OPN"
	}
	
	prefix := fmt.Sprintf("STK/%s/%s/%s/", typeCode, branchCode, dateStr)

	var lastTxn models.StockTransaction
	err := r.db.
		Where("transaction_number LIKE ?", prefix+"%").
		Order("transaction_number DESC").
		First(&lastTxn).Error

	var sequence int
	if err == nil {
		var lastSeq int
		fmt.Sscanf(lastTxn.TransactionNumber[len(prefix):], "%d", &lastSeq)
		sequence = lastSeq + 1
	} else {
		sequence = 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}

// Stock Opname
func (r *inventoryRepository) GetAllOpnames(params *models.PaginationParams) ([]models.StockOpname, int64, error) {
	var opnames []models.StockOpname
	var total int64

	query := r.db.Model(&models.StockOpname{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Branch").
		Preload("Preparer").
		Preload("Approver").
		Order("opname_date DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&opnames).Error

	return opnames, total, err
}

func (r *inventoryRepository) GetOpnameByID(id uuid.UUID) (*models.StockOpname, error) {
	var opname models.StockOpname
	err := r.db.
		Preload("Branch").
		Preload("Preparer").
		Preload("Approver").
		Preload("Items").
		Preload("Items.Item").
		First(&opname, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("opname not found")
		}
		return nil, err
	}
	return &opname, nil
}

func (r *inventoryRepository) CreateOpname(opname *models.StockOpname) error {
	return r.db.Create(opname).Error
}

func (r *inventoryRepository) UpdateOpname(opname *models.StockOpname) error {
	return r.db.Save(opname).Error
}

func (r *inventoryRepository) ApproveOpname(id uuid.UUID, approverID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.StockOpname{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      models.OpnameStatusApproved,
			"approved_by": approverID,
			"approved_at": now,
		}).Error
}

func (r *inventoryRepository) GenerateOpnameNumber(branchCode string, date time.Time) (string, error) {
	// Format: OPN/BRANCH/YYYYMM/XXXX
	yearMonth := date.Format("200601")
	prefix := fmt.Sprintf("OPN/%s/%s/", branchCode, yearMonth)

	var lastOpname models.StockOpname
	err := r.db.
		Where("opname_number LIKE ?", prefix+"%").
		Order("opname_number DESC").
		First(&lastOpname).Error

	var sequence int
	if err == nil {
		var lastSeq int
		fmt.Sscanf(lastOpname.OpnameNumber[len(prefix):], "%d", &lastSeq)
		sequence = lastSeq + 1
	} else {
		sequence = 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}

// Purchase Orders
func (r *inventoryRepository) GetAllPOs(params *models.PaginationParams) ([]models.PurchaseOrder, int64, error) {
	var pos []models.PurchaseOrder
	var total int64

	query := r.db.Model(&models.PurchaseOrder{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Branch").
		Preload("Requester").
		Order("order_date DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&pos).Error

	return pos, total, err
}

func (r *inventoryRepository) GetPOByID(id uuid.UUID) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	err := r.db.
		Preload("Branch").
		Preload("Requester").
		Preload("Items").
		Preload("Items.Item").
		First(&po, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("purchase order not found")
		}
		return nil, err
	}
	return &po, nil
}

func (r *inventoryRepository) CreatePO(po *models.PurchaseOrder) error {
	return r.db.Create(po).Error
}

func (r *inventoryRepository) UpdatePO(po *models.PurchaseOrder) error {
	return r.db.Save(po).Error
}

func (r *inventoryRepository) ApprovePO(id uuid.UUID, approverID uuid.UUID) error {
	return r.db.Model(&models.PurchaseOrder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      models.POStatusSubmitted,
			"approved_by": approverID,
		}).Error
}

func (r *inventoryRepository) GeneratePONumber(branchCode string, date time.Time) (string, error) {
	// Format: PO/BRANCH/YYYYMM/XXXX
	yearMonth := date.Format("200601")
	prefix := fmt.Sprintf("PO/%s/%s/", branchCode, yearMonth)

	var lastPO models.PurchaseOrder
	err := r.db.
		Where("po_number LIKE ?", prefix+"%").
		Order("po_number DESC").
		First(&lastPO).Error

	var sequence int
	if err == nil {
		var lastSeq int
		fmt.Sscanf(lastPO.PONumber[len(prefix):], "%d", &lastSeq)
		sequence = lastSeq + 1
	} else {
		sequence = 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
