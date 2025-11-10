package service

import (
	"errors"
	//"time" // Sudah benar dikomentari karena tidak terpakai

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type InventoryService interface {
	// Items
	GetAllItems(params *models.PaginationParams) (*models.PaginationResponse, error)
	GetItemByID(id uuid.UUID) (*models.InventoryItem, error)
	GetLowStockItems(branchID uuid.UUID) ([]models.InventoryItem, error)
	SearchItems(keyword string) ([]models.InventoryItem, error)
	CreateItem(req *models.CreateInventoryItemRequest) (*models.InventoryItem, error)
	UpdateItem(id uuid.UUID, req *models.CreateInventoryItemRequest) (*models.InventoryItem, error)
	DeleteItem(id uuid.UUID) error
	
	// Transactions
	CreateStockIn(req *models.CreateStockTransactionRequest, userID uuid.UUID) (*models.StockTransaction, error)
	CreateStockOut(req *models.CreateStockTransactionRequest, userID uuid.UUID) (*models.StockTransaction, error)
	CreateAdjustment(req *models.CreateStockTransactionRequest, userID uuid.UUID) (*models.StockTransaction, error)
	GetTransactionHistory(itemID uuid.UUID) ([]models.StockTransaction, error)
	
	// Stock Opname
	CreateStockOpname(req *models.CreateStockOpnameRequest, userID uuid.UUID) (*models.StockOpname, error)
	ApproveStockOpname(id uuid.UUID, approverID uuid.UUID) (*models.StockOpname, error)
	ProcessStockOpname(id uuid.UUID) error
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
	branchRepo    repository.BranchRepository
}

func NewInventoryService(
	inventoryRepo repository.InventoryRepository,
	branchRepo repository.BranchRepository,
) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
		branchRepo:    branchRepo,
	}
}

func (s *inventoryService) GetAllItems(params *models.PaginationParams) (*models.PaginationResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.GlobalConfig.App.DefaultPageSize
	}

	items, total, err := s.inventoryRepo.GetAllItems(params)
	if err != nil {
		return nil, err
	}

	data := make([]interface{}, len(items))
	for i := range items {
		data[i] = items[i]
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &models.PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *inventoryService) GetItemByID(id uuid.UUID) (*models.InventoryItem, error) {
	return s.inventoryRepo.GetItemByID(id)
}

func (s *inventoryService) GetLowStockItems(branchID uuid.UUID) ([]models.InventoryItem, error) {
	return s.inventoryRepo.GetLowStockItems(branchID)
}

func (s *inventoryService) SearchItems(keyword string) ([]models.InventoryItem, error) {
	return s.inventoryRepo.SearchItems(keyword)
}

func (s *inventoryService) CreateItem(req *models.CreateInventoryItemRequest) (*models.InventoryItem, error) {
	// Validate branch
	// --- INI PERBAIKANNYA ---
	// Kita hanya perlu mengecek error-nya, bukan menggunakan variabel 'branch'
	_, err := s.branchRepo.GetByID(req.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	// Generate item code if needed (simple implementation)
	itemCode := req.Name // You can implement better code generation
	
	item := &models.InventoryItem{
		ItemCode:      itemCode,
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		BranchID:      req.BranchID,
		Description:   req.Description,
		Unit:          req.Unit,
		CurrentStock:  0,
		MinimumStock:  req.MinimumStock,
		PurchasePrice: req.PurchasePrice,
		SellingPrice:  req.SellingPrice,
		IsActive:      true,
		IsSaleable:    req.IsSaleable,
	}

	if err := s.inventoryRepo.CreateItem(item); err != nil {
		return nil, err
	}

	return s.inventoryRepo.GetItemByID(item.ID)
}

func (s *inventoryService) UpdateItem(id uuid.UUID, req *models.CreateInventoryItemRequest) (*models.InventoryItem, error) {
	item, err := s.inventoryRepo.GetItemByID(id)
	if err != nil {
		return nil, errors.New("item not found")
	}

	item.Name = req.Name
	item.Description = req.Description
	item.Unit = req.Unit
	item.MinimumStock = req.MinimumStock
	item.PurchasePrice = req.PurchasePrice
	item.SellingPrice = req.SellingPrice
	item.IsSaleable = req.IsSaleable

	if err := s.inventoryRepo.UpdateItem(item); err != nil {
		return nil, err
	}

	return s.inventoryRepo.GetItemByID(item.ID)
}

func (s *inventoryService) DeleteItem(id uuid.UUID) error {
	item, err := s.inventoryRepo.GetItemByID(id)
	if err != nil {
		return errors.New("item not found")
	}

	if item.CurrentStock > 0 {
		return errors.New("cannot delete item with stock")
	}

	return s.inventoryRepo.DeleteItem(id)
}

func (s *inventoryService) CreateStockIn(req *models.CreateStockTransactionRequest, userID uuid.UUID) (*models.StockTransaction, error) {
	return s.createTransaction(req, models.TransactionTypeIn, userID)
}

func (s *inventoryService) CreateStockOut(req *models.CreateStockTransactionRequest, userID uuid.UUID) (*models.StockTransaction, error) {
	return s.createTransaction(req, models.TransactionTypeOut, userID)
}

func (s *inventoryService) CreateAdjustment(req *models.CreateStockTransactionRequest, userID uuid.UUID) (*models.StockTransaction, error) {
	return s.createTransaction(req, models.TransactionTypeAdjustment, userID)
}

func (s *inventoryService) createTransaction(req *models.CreateStockTransactionRequest, txType string, userID uuid.UUID) (*models.StockTransaction, error) {
	// Validate item
	item, err := s.inventoryRepo.GetItemByID(req.ItemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Get branch
	branch, err := s.branchRepo.GetByID(item.BranchID)
	if err != nil {
		return nil, err
	}

	// Generate transaction number
	txNumber, err := s.inventoryRepo.GenerateTransactionNumber(txType, branch.Code, req.TransactionDate)
	if err != nil {
		return nil, err
	}

	// Create transaction
	transaction := &models.StockTransaction{
		TransactionNumber: txNumber,
		ItemID:           req.ItemID,
		BranchID:         item.BranchID,
		TransactionType:  txType,
		TransactionDate:  req.TransactionDate,
		Quantity:         req.Quantity,
		UnitPrice:        req.UnitPrice,
		TotalValue:       req.Quantity * req.UnitPrice,
		Supplier:         req.Supplier,
		Customer:         req.Customer,
		Reason:           req.Reason,
		CreatedBy:        userID,
	}

	if err := s.inventoryRepo.CreateTransaction(transaction); err != nil {
		return nil, err
	}

	return s.inventoryRepo.GetTransactionByID(transaction.ID)
}

func (s *inventoryService) GetTransactionHistory(itemID uuid.UUID) ([]models.StockTransaction, error) {
	return s.inventoryRepo.GetTransactionsByItem(itemID)
}

func (s *inventoryService) CreateStockOpname(req *models.CreateStockOpnameRequest, userID uuid.UUID) (*models.StockOpname, error) {
	// Get branch
	branch, err := s.branchRepo.GetByID(req.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	// Generate opname number
	opnameNumber, err := s.inventoryRepo.GenerateOpnameNumber(branch.Code, req.OpnameDate)
	if err != nil {
		return nil, err
	}

	// Create opname
	opname := &models.StockOpname{
		OpnameNumber: opnameNumber,
		BranchID:     req.BranchID,
		OpnameDate:   req.OpnameDate,
		Status:       models.OpnameStatusDraft,
		PreparedBy:   userID,
	}

	// Add items
	for _, itemReq := range req.Items {
		item, err := s.inventoryRepo.GetItemByID(itemReq.ItemID)
		if err != nil {
			continue
		}

		difference := itemReq.PhysicalStock - item.CurrentStock
		differenceValue := difference * item.PurchasePrice

		opname.Items = append(opname.Items, models.StockOpnameItem{
			ItemID:          itemReq.ItemID,
			SystemStock:     item.CurrentStock,
			PhysicalStock:   itemReq.PhysicalStock,
			Difference:      difference,
			UnitPrice:       item.PurchasePrice,
			DifferenceValue: differenceValue,
			Remarks:         itemReq.Remarks,
		})
	}

	if err := s.inventoryRepo.CreateOpname(opname); err != nil {
		return nil, err
	}

	return s.inventoryRepo.GetOpnameByID(opname.ID)
}

func (s *inventoryService) ApproveStockOpname(id uuid.UUID, approverID uuid.UUID) (*models.StockOpname, error) {
	opname, err := s.inventoryRepo.GetOpnameByID(id)
	if err != nil {
		return nil, errors.New("opname not found")
	}

	if opname.Status != models.OpnameStatusDraft {
		return nil, errors.New("opname is not in draft status")
	}

	if err := s.inventoryRepo.ApproveOpname(id, approverID); err != nil {
		return nil, err
	}

	return s.inventoryRepo.GetOpnameByID(id)
}

func (s *inventoryService) ProcessStockOpname(id uuid.UUID) error {
	opname, err := s.inventoryRepo.GetOpnameByID(id)
	if err != nil {
		return errors.New("opname not found")
	}

	if opname.Status != models.OpnameStatusApproved {
		return errors.New("opname is not approved")
	}

	// Create adjustment transactions for each item
	for _, item := range opname.Items {
		if item.Difference == 0 {
			continue
		}

		branch, _ := s.branchRepo.GetByID(opname.BranchID)
		txNumber, _ := s.inventoryRepo.GenerateTransactionNumber(
			models.TransactionTypeOpname,
			branch.Code,
			opname.OpnameDate,
		)

		transaction := &models.StockTransaction{
			TransactionNumber: txNumber,
			ItemID:           item.ItemID,
			BranchID:         opname.BranchID,
			TransactionType:  models.TransactionTypeOpname,
			TransactionDate:  opname.OpnameDate,
			Quantity:         item.Difference,
			UnitPrice:        item.UnitPrice,
			TotalValue:       item.DifferenceValue,
			Reason:           "Stock Opname: " + opname.OpnameNumber,
			CreatedBy:        opname.PreparedBy,
			ReferenceType:    "opname",
			ReferenceID:      &opname.ID,
			ReferenceNumber:  opname.OpnameNumber,
		}

		if err := s.inventoryRepo.CreateTransaction(transaction); err != nil {
			return err
		}
	}

	return nil
}