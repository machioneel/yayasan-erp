package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/service"
	"github.com/yayasan/erp-backend/internal/utils"
)

type EmployeeHandler struct {
	employeeService service.EmployeeService
	payrollService  service.PayrollService
}

func NewEmployeeHandler(employeeService service.EmployeeService, payrollService service.PayrollService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
		payrollService:  payrollService,
	}
}

// Employee endpoints
func (h *EmployeeHandler) GetAll(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.employeeService.GetAll(&params)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Employees retrieved successfully", result)
}

func (h *EmployeeHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	employee, err := h.employeeService.GetByID(id)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Employee retrieved successfully", employee)
}

func (h *EmployeeHandler) GetTeachers(c *gin.Context) {
	teachers, err := h.employeeService.GetTeachers()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Teachers retrieved successfully", teachers)
}

func (h *EmployeeHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Keyword is required")
		return
	}

	employees, err := h.employeeService.Search(keyword)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Search results", employees)
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var req models.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	employee, err := h.employeeService.Create(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Employee created successfully", employee)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	var req models.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	employee, err := h.employeeService.Update(id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Employee updated successfully", employee)
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	if err := h.employeeService.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Employee deleted successfully", nil)
}

func (h *EmployeeHandler) GetStatistics(c *gin.Context) {
	stats, err := h.employeeService.GetStatistics()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved successfully", stats)
}

// Payroll endpoints
func (h *EmployeeHandler) GetAllPayrolls(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.payrollService.GetAll(&params)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payrolls retrieved successfully", result)
}

func (h *EmployeeHandler) GetPayrollByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payroll ID")
		return
	}

	payroll, err := h.payrollService.GetByID(id)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payroll retrieved successfully", payroll)
}

func (h *EmployeeHandler) GetPayrollsByPeriod(c *gin.Context) {
	period := c.Param("period")
	if period == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Period is required")
		return
	}

	payrolls, err := h.payrollService.GetByPeriod(period)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payrolls retrieved successfully", payrolls)
}

func (h *EmployeeHandler) GetEmployeePayrolls(c *gin.Context) {
	idStr := c.Param("employee_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	payrolls, err := h.payrollService.GetByEmployee(id)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Employee payrolls retrieved successfully", payrolls)
}

func (h *EmployeeHandler) CreatePayroll(c *gin.Context) {
	var req models.CreatePayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	payroll, err := h.payrollService.Create(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Payroll created successfully", payroll)
}

func (h *EmployeeHandler) ProcessPayroll(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payroll ID")
		return
	}

	payroll, err := h.payrollService.ProcessPayroll(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payroll processed successfully", payroll)
}

func (h *EmployeeHandler) DeletePayroll(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payroll ID")
		return
	}

	if err := h.payrollService.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payroll deleted successfully", nil)
}

// Bulk payroll generation
type GenerateBulkPayrollRequest struct {
	BranchID    uuid.UUID `json:"branch_id" binding:"required"`
	Period      string    `json:"period" binding:"required"`
	PaymentDate string    `json:"payment_date" binding:"required"`
}

func (h *EmployeeHandler) GenerateBulkPayroll(c *gin.Context) {
	var req GenerateBulkPayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Parse payment date
	paymentDate, err := utils.ParseDate(req.PaymentDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment date format")
		return
	}

	payrolls, err := h.payrollService.GenerateBulkPayroll(req.BranchID, req.Period, paymentDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Bulk payroll generated successfully", map[string]interface{}{
		"total":    len(payrolls),
		"payrolls": payrolls,
	})
}
