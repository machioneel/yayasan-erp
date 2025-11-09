package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/service"
	"github.com/yayasan/erp-backend/internal/utils"
)

type JournalHandler struct {
	journalService service.JournalService
}

func NewJournalHandler(journalService service.JournalService) *JournalHandler {
	return &JournalHandler{journalService: journalService}
}

func (h *JournalHandler) GetAll(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.journalService.GetAll(&params)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journals retrieved successfully", result)
}

func (h *JournalHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID")
		return
	}

	journal, err := h.journalService.GetByID(id)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal retrieved successfully", journal.ToJournalResponse())
}

func (h *JournalHandler) GetByStatus(c *gin.Context) {
	status := c.Param("status")
	
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.journalService.GetByStatus(status, &params)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journals retrieved successfully", result)
}

func (h *JournalHandler) Create(c *gin.Context) {
	var req models.CreateJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, _ := c.Get("user_id")
	journal, err := h.journalService.Create(&req, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Journal created successfully", journal.ToJournalResponse())
}

func (h *JournalHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID")
		return
	}

	var req models.UpdateJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, _ := c.Get("user_id")
	journal, err := h.journalService.Update(id, &req, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal updated successfully", journal.ToJournalResponse())
}

func (h *JournalHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID")
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.journalService.Delete(id, userID.(uuid.UUID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal deleted successfully", nil)
}

func (h *JournalHandler) SubmitForReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID")
		return
	}

	userID, _ := c.Get("user_id")
	journal, err := h.journalService.SubmitForReview(id, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal submitted for review", journal.ToJournalResponse())
}

func (h *JournalHandler) Review(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID")
		return
	}

	var req models.ReviewJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, _ := c.Get("user_id")
	journal, err := h.journalService.Review(id, &req, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal reviewed successfully", journal.ToJournalResponse())
}

func (h *JournalHandler) Post(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID")
		return
	}

	var req models.PostJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, _ := c.Get("user_id")
	journal, err := h.journalService.Post(id, &req, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal posted successfully", journal.ToJournalResponse())
}

func (h *JournalHandler) Unpost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID")
		return
	}

	userID, _ := c.Get("user_id")
	journal, err := h.journalService.Unpost(id, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal unposted successfully", journal.ToJournalResponse())
}
