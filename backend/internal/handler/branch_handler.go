package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/service"
	"github.com/yayasan/erp-backend/internal/utils"
)

type BranchHandler struct {
	branchService service.BranchService
}

func NewBranchHandler(branchService service.BranchService) *BranchHandler {
	return &BranchHandler{branchService: branchService}
}

func (h *BranchHandler) GetAll(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	branches, total, err := h.branchService.GetAll(&params)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]models.BranchResponse, len(branches))
	for i, branch := range branches {
		responses[i] = *branch.ToBranchResponse()
	}

	utils.PaginatedResponse(c, responses, total, params.Page, params.PageSize)
}

func (h *BranchHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid branch ID")
		return
	}

	branch, err := h.branchService.GetByID(id)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Branch retrieved successfully", branch.ToBranchResponse())
}

func (h *BranchHandler) GetAllActive(c *gin.Context) {
	branches, err := h.branchService.GetAllActive()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	responses := make([]models.BranchResponse, len(branches))
	for i, branch := range branches {
		responses[i] = *branch.ToBranchResponse()
	}

	utils.SuccessResponse(c, http.StatusOK, "Active branches retrieved", responses)
}

func (h *BranchHandler) Create(c *gin.Context) {
	var req models.CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	createdBy, _ := c.Get("user_id")
	branch, err := h.branchService.Create(&req, createdBy.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Branch created successfully", branch.ToBranchResponse())
}

func (h *BranchHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid branch ID")
		return
	}

	var req models.UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	branch, err := h.branchService.Update(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Branch updated successfully", branch.ToBranchResponse())
}

func (h *BranchHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid branch ID")
		return
	}

	if err := h.branchService.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Branch deleted successfully", nil)
}
