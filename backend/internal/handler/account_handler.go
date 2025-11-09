package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/service"
	"github.com/yayasan/erp-backend/internal/utils"
)

type AccountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) GetAll(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.accountService.GetAll(&params)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Accounts retrieved successfully", result)
}

func (h *AccountHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	account, err := h.accountService.GetByID(id)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account retrieved successfully", account.ToAccountResponse())
}

func (h *AccountHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	
	account, err := h.accountService.GetByCode(code)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account retrieved successfully", account.ToAccountResponse())
}

func (h *AccountHandler) GetTree(c *gin.Context) {
	tree, err := h.accountService.GetTree()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account tree retrieved successfully", tree)
}

func (h *AccountHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	
	accounts, err := h.accountService.GetByCategory(category)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]models.AccountResponse, len(accounts))
	for i, acc := range accounts {
		responses[i] = *acc.ToAccountResponse()
	}

	utils.SuccessResponse(c, http.StatusOK, "Accounts retrieved successfully", responses)
}

func (h *AccountHandler) GetDetailAccounts(c *gin.Context) {
	accounts, err := h.accountService.GetDetailAccounts()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]models.AccountResponse, len(accounts))
	for i, acc := range accounts {
		responses[i] = *acc.ToAccountResponse()
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail accounts retrieved successfully", responses)
}

func (h *AccountHandler) Create(c *gin.Context) {
	var req models.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	account, err := h.accountService.Create(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Account created successfully", account.ToAccountResponse())
}

func (h *AccountHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	var req models.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	account, err := h.accountService.Update(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account updated successfully", account.ToAccountResponse())
}

func (h *AccountHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	if err := h.accountService.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account deleted successfully", nil)
}

func (h *AccountHandler) BulkImport(c *gin.Context) {
	var req models.ImportAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	count, err := h.accountService.BulkImport(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Accounts imported successfully", gin.H{
		"imported_count": count,
	})
}
