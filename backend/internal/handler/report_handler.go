package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/service"
	"github.com/yayasan/erp-backend/internal/utils"
)

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) GetTrialBalance(c *gin.Context) {
	var req models.TrialBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.reportService.GetTrialBalance(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Trial balance generated successfully", result)
}

func (h *ReportHandler) GetBalanceSheet(c *gin.Context) {
	var req models.BalanceSheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.reportService.GetBalanceSheet(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Balance sheet generated successfully", result)
}

func (h *ReportHandler) GetIncomeStatement(c *gin.Context) {
	var req models.IncomeStatementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.reportService.GetIncomeStatement(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Income statement generated successfully", result)
}

func (h *ReportHandler) GetGeneralLedger(c *gin.Context) {
	var req models.GeneralLedgerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.reportService.GetGeneralLedger(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "General ledger generated successfully", result)
}
