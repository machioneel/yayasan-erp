package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/service"
	"github.com/yayasan/erp-backend/internal/utils"
)

type InventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

// Items
func (h *InventoryHandler) GetAllItems(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.inventoryService.GetAllItems(&params)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Items retrieved successfully", result)
}

func (h *InventoryHandler) GetItemByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID")
		return
	}

	item, err := h.inventoryService.GetItemByID(id)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item retrieved successfully", item)
}

func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	branchIDStr := c.Query("branch_id")
	var branchID uuid.UUID
	if branchIDStr != "" {
		branchID, _ = uuid.Parse(branchIDStr)
	}

	items, err := h.inventoryService.GetLowStockItems(branchID)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Low stock items retrieved successfully", items)
}

func (h *InventoryHandler) SearchItems(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Keyword is required")
		return
	}

	items, err := h.inventoryService.SearchItems(keyword)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Search results", items)
}

func (h *InventoryHandler) CreateItem(c *gin.Context) {
	var req models.CreateInventoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	item, err := h.inventoryService.CreateItem(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Item created successfully", item)
}

func (h *InventoryHandler) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req models.CreateInventoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	item, err := h.inventoryService.UpdateItem(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item updated successfully", item)
}

func (h *InventoryHandler) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID")
		return
	}

	if err := h.inventoryService.DeleteItem(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item deleted successfully", nil)
}

// Stock Transactions
func (h *InventoryHandler) CreateStockIn(c *gin.Context) {
	var req models.CreateStockTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, _ := c.Get("user_id")
	transaction, err := h.inventoryService.CreateStockIn(&req, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Stock IN created successfully", transaction)
}

func (h *InventoryHandler) CreateStockOut(c *gin.Context) {
	var req models.CreateStockTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, _ := c.Get("user_id")
	transaction, err := h.inventoryService.CreateStockOut(&req, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Stock OUT created successfully", transaction)
}

func (h *InventoryHandler) CreateAdjustment(c *gin.Context) {
	var req models.CreateStockTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, _ := c.Get("user_id")
	transaction, err := h.inventoryService.CreateAdjustment(&req, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Stock adjustment created successfully", transaction)
}

func (h *InventoryHandler) GetTransactionHistory(c *gin.Context) {
	idStr := c.Param("item_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID")
		return
	}

	transactions, err := h.inventoryService.GetTransactionHistory(id)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction history retrieved successfully", transactions)
}

// Stock Opname
func (h *InventoryHandler) CreateStockOpname(c *gin.Context) {
	var req models.CreateStockOpnameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, _ := c.Get("user_id")
	opname, err := h.inventoryService.CreateStockOpname(&req, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Stock opname created successfully", opname)
}

func (h *InventoryHandler) ApproveStockOpname(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid opname ID")
		return
	}

	userID, _ := c.Get("user_id")
	opname, err := h.inventoryService.ApproveStockOpname(id, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock opname approved successfully", opname)
}

func (h *InventoryHandler) ProcessStockOpname(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid opname ID")
		return
	}

	if err := h.inventoryService.ProcessStockOpname(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock opname processed successfully", nil)
}
