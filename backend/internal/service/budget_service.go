package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
	"gorm.io/gorm"
)

type BudgetService interface {
	GetAll(params *models.PaginationParams) (*models.PaginationResponse, error)
	GetByID(id uuid.UUID) (*models.Budget, error)
	GetByFiscalYear(fiscalYearID uuid.UUID) ([]models.Budget, error)
	Create(req *models.CreateBudgetRequest) (*models.Budget, error)
	Update(id uuid.UUID, req *models.UpdateBudgetRequest) (*models.Budget, error)
	Delete(id uuid.UUID) error
	GetBudgetVsActual(req *models.BudgetVsActualRequest) (*models.BudgetVsActualResponse, error)
}

type budgetService struct {
	db             *gorm.DB
	budgetRepo     repository.BudgetRepository
	accountRepo    repository.AccountRepository
	fiscalYearRepo repository.FiscalYearRepository
}

func NewBudgetService(
	db *gorm.DB,
	budgetRepo repository.BudgetRepository,
	accountRepo repository.AccountRepository,
	fiscalYearRepo repository.FiscalYearRepository,
) BudgetService {
	return &budgetService{
		db:             db,
		budgetRepo:     budgetRepo,
		accountRepo:    accountRepo,
		fiscalYearRepo: fiscalYearRepo,
	}
}

func (s *budgetService) GetAll(params *models.PaginationParams) (*models.PaginationResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.AppConfig.App.DefaultPageSize
	}
	if params.PageSize > config.AppConfig.App.MaxPageSize {
		params.PageSize = config.AppConfig.App.MaxPageSize
	}

	budgets, total, err := s.budgetRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	// Convert to response
	items := make([]interface{}, len(budgets))
	for i, budget := range budgets {
		items[i] = s.toBudgetResponse(&budget)
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

func (s *budgetService) GetByID(id uuid.UUID) (*models.Budget, error) {
	return s.budgetRepo.GetByID(id)
}

func (s *budgetService) GetByFiscalYear(fiscalYearID uuid.UUID) ([]models.Budget, error) {
	return s.budgetRepo.GetByFiscalYear(fiscalYearID)
}

func (s *budgetService) Create(req *models.CreateBudgetRequest) (*models.Budget, error) {
	// Validate fiscal year exists
	fiscalYear, err := s.fiscalYearRepo.GetByID(req.FiscalYearID)
	if err != nil {
		return nil, errors.New("fiscal year not found")
	}

	if fiscalYear.IsClosed {
		return nil, errors.New("cannot create budget for closed fiscal year")
	}

	// Validate account exists
	account, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	// Only expense accounts can have budgets
	if account.Category != models.AccountCategoryExpense {
		return nil, errors.New("budgets can only be created for expense accounts")
	}

	// Check if budget already exists
	existing, _ := s.budgetRepo.GetByAccountPeriod(req.AccountID, req.Period, req.BranchID, req.FundID, req.ProgramID)
	if existing != nil {
		return nil, errors.New("budget already exists for this account and period")
	}

	// Create budget
	budget := &models.Budget{
		FiscalYearID: req.FiscalYearID,
		BranchID:     req.BranchID,
		AccountID:    req.AccountID,
		FundID:       req.FundID,
		ProgramID:    req.ProgramID,
		Period:       req.Period,
		Amount:       req.Amount,
		Description:  req.Description,
		IsActive:     true,
	}

	if err := s.budgetRepo.Create(budget); err != nil {
		return nil, err
	}

	return s.budgetRepo.GetByID(budget.ID)
}

func (s *budgetService) Update(id uuid.UUID, req *models.UpdateBudgetRequest) (*models.Budget, error) {
	budget, err := s.budgetRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("budget not found")
	}

	// Check if fiscal year is closed
	fiscalYear, err := s.fiscalYearRepo.GetByID(budget.FiscalYearID)
	if err != nil {
		return nil, err
	}

	if fiscalYear.IsClosed {
		return nil, errors.New("cannot update budget for closed fiscal year")
	}

	// Update fields
	budget.Amount = req.Amount
	budget.Description = req.Description

	if req.IsActive != nil {
		budget.IsActive = *req.IsActive
	}

	if err := s.budgetRepo.Update(budget); err != nil {
		return nil, err
	}

	return s.budgetRepo.GetByID(budget.ID)
}

func (s *budgetService) Delete(id uuid.UUID) error {
	budget, err := s.budgetRepo.GetByID(id)
	if err != nil {
		return errors.New("budget not found")
	}

	// Check if fiscal year is closed
	fiscalYear, err := s.fiscalYearRepo.GetByID(budget.FiscalYearID)
	if err != nil {
		return err
	}

	if fiscalYear.IsClosed {
		return errors.New("cannot delete budget for closed fiscal year")
	}

	return s.budgetRepo.Delete(id)
}

func (s *budgetService) GetBudgetVsActual(req *models.BudgetVsActualRequest) (*models.BudgetVsActualResponse, error) {
	// Get fiscal year
	fiscalYear, err := s.fiscalYearRepo.GetByID(req.FiscalYearID)
	if err != nil {
		return nil, errors.New("fiscal year not found")
	}

	// Get budgets
	var budgets []models.Budget
	query := s.db.Model(&models.Budget{}).
		Where("fiscal_year_id = ?", req.FiscalYearID).
		Preload("Account")

	if req.Period != "" {
		query = query.Where("period = ?", req.Period)
	}
	if req.BranchID != nil {
		query = query.Where("branch_id = ?", *req.BranchID)
	}
	if req.AccountID != nil {
		query = query.Where("account_id = ?", *req.AccountID)
	}

	if err := query.Find(&budgets).Error; err != nil {
		return nil, err
	}

	lines := make([]models.BudgetVsActualLine, 0, len(budgets))
	var totalBudget, totalActual float64

	for _, budget := range budgets {
		// Calculate actual spending
		actual, _ := s.calculateActualSpending(budget)

		variance := budget.Amount - actual
		variancePct := float64(0)
		if budget.Amount > 0 {
			variancePct = (variance / budget.Amount) * 100
		}

		lines = append(lines, models.BudgetVsActualLine{
			AccountCode: budget.Account.Code,
			AccountName: budget.Account.Name,
			Budget:      budget.Amount,
			Actual:      actual,
			Variance:    variance,
			VariancePct: variancePct,
		})

		totalBudget += budget.Amount
		totalActual += actual
	}

	totalVariance := totalBudget - totalActual
	totalVariancePct := float64(0)
	if totalBudget > 0 {
		totalVariancePct = (totalVariance / totalBudget) * 100
	}

	return &models.BudgetVsActualResponse{
		FiscalYear: fiscalYear.Name,
		Period:     req.Period,
		Lines:      lines,
		Summary: models.BudgetVsActualSummary{
			TotalBudget:      totalBudget,
			TotalActual:      totalActual,
			TotalVariance:    totalVariance,
			TotalVariancePct: totalVariancePct,
		},
	}, nil
}

// Helper functions

func (s *budgetService) calculateActualSpending(budget models.Budget) (float64, error) {
	var totalDebit, totalCredit float64

	query := s.db.Model(&models.JournalLine{}).
		Select("COALESCE(SUM(debit), 0) as total_debit, COALESCE(SUM(credit), 0) as total_credit").
		Joins("JOIN journals ON journals.id = journal_lines.journal_id").
		Where("journal_lines.account_id = ?", budget.AccountID).
		Where("DATE_FORMAT(journals.journal_date, '%Y-%m') = ?", budget.Period).
		Where("journals.is_posted = ?", true)

	if budget.BranchID != nil {
		query = query.Where("journals.branch_id = ?", *budget.BranchID)
	}
	if budget.FundID != nil {
		query = query.Where("journal_lines.fund_id = ?", *budget.FundID)
	}
	if budget.ProgramID != nil {
		query = query.Where("journal_lines.program_id = ?", *budget.ProgramID)
	}

	err := query.Row().Scan(&totalDebit, &totalCredit)
	if err != nil {
		return 0, err
	}

	// For expense accounts, debit increases the expense
	return totalDebit - totalCredit, nil
}

func (s *budgetService) toBudgetResponse(budget *models.Budget) *models.BudgetResponse {
	resp := &models.BudgetResponse{
		ID:           budget.ID,
		FiscalYearID: budget.FiscalYearID,
		BranchID:     budget.BranchID,
		AccountID:    budget.AccountID,
		AccountCode:  budget.Account.Code,
		AccountName:  budget.Account.Name,
		FundID:       budget.FundID,
		ProgramID:    budget.ProgramID,
		Period:       budget.Period,
		Amount:       budget.Amount,
		Description:  budget.Description,
		IsActive:     budget.IsActive,
		CreatedAt:    budget.CreatedAt,
	}

	if budget.FiscalYear.ID != uuid.Nil {
		resp.FiscalYear = budget.FiscalYear.Name
	}
	if budget.Branch != nil {
		resp.BranchName = budget.Branch.Name
	}
	if budget.Fund != nil {
		resp.FundName = budget.Fund.Name
	}
	if budget.Program != nil {
		resp.ProgramName = budget.Program.Name
	}

	// Calculate actual
	actual, _ := s.calculateActualSpending(*budget)
	resp.Actual = actual
	resp.Variance = budget.Amount - actual
	if budget.Amount > 0 {
		resp.VariancePct = (resp.Variance / budget.Amount) * 100
	}

	return resp
}
