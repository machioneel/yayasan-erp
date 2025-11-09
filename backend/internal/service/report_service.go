package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
	"gorm.io/gorm"
)

type ReportService interface {
	GetTrialBalance(req *models.TrialBalanceRequest) (*models.TrialBalanceResponse, error)
	GetBalanceSheet(req *models.BalanceSheetRequest) (*models.BalanceSheetResponse, error)
	GetIncomeStatement(req *models.IncomeStatementRequest) (*models.IncomeStatementResponse, error)
	GetGeneralLedger(req *models.GeneralLedgerRequest) (*models.GeneralLedgerResponse, error)
}

type reportService struct {
	db          *gorm.DB
	accountRepo repository.AccountRepository
	journalRepo repository.JournalRepository
}

func NewReportService(
	db *gorm.DB,
	accountRepo repository.AccountRepository,
	journalRepo repository.JournalRepository,
) ReportService {
	return &reportService{
		db:          db,
		accountRepo: accountRepo,
		journalRepo: journalRepo,
	}
}

func (s *reportService) GetTrialBalance(req *models.TrialBalanceRequest) (*models.TrialBalanceResponse, error) {
	// Get all detail accounts
	accounts, err := s.accountRepo.GetDetailAccounts()
	if err != nil {
		return nil, err
	}

	lines := make([]models.TrialBalanceLine, 0)
	var totalDebit, totalCredit float64

	// Calculate balance for each account
	for _, account := range accounts {
		balance, err := s.calculateAccountBalance(account.ID, req.AsOfDate, req.BranchID, req.FundID, req.ProgramID)
		if err != nil {
			continue
		}

		// Skip zero balances
		if balance == 0 {
			continue
		}

		line := models.TrialBalanceLine{
			AccountCode: account.Code,
			AccountName: account.Name,
			Category:    account.Category,
		}

		// Debit or credit based on normal balance
		if account.GetNormalBalance() == models.NormalBalanceDebit {
			if balance >= 0 {
				line.Debit = balance
			} else {
				line.Credit = -balance
			}
		} else {
			if balance >= 0 {
				line.Credit = balance
			} else {
				line.Debit = -balance
			}
		}

		totalDebit += line.Debit
		totalCredit += line.Credit

		lines = append(lines, line)
	}

	return &models.TrialBalanceResponse{
		AsOfDate: req.AsOfDate,
		Lines:    lines,
		Summary: models.TrialBalanceSummary{
			TotalDebit:  totalDebit,
			TotalCredit: totalCredit,
			Difference:  totalDebit - totalCredit,
			IsBalanced:  totalDebit == totalCredit,
		},
	}, nil
}

func (s *reportService) GetBalanceSheet(req *models.BalanceSheetRequest) (*models.BalanceSheetResponse, error) {
	// Get accounts by category
	assets, err := s.accountRepo.GetByCategory(models.AccountCategoryAsset)
	if err != nil {
		return nil, err
	}

	liabilities, err := s.accountRepo.GetByCategory(models.AccountCategoryLiability)
	if err != nil {
		return nil, err
	}

	equity, err := s.accountRepo.GetByCategory(models.AccountCategoryEquity)
	if err != nil {
		return nil, err
	}

	response := &models.BalanceSheetResponse{
		AsOfDate: req.AsOfDate,
	}

	// Assets section
	response.Assets = s.buildBalanceSheetSection(assets, req.AsOfDate, req.BranchID, req.FundID)
	response.TotalAssets = response.Assets.Total

	// Liabilities section
	response.Liabilities = s.buildBalanceSheetSection(liabilities, req.AsOfDate, req.BranchID, req.FundID)
	response.TotalLiabilities = response.Liabilities.Total

	// Equity section
	response.Equity = s.buildBalanceSheetSection(equity, req.AsOfDate, req.BranchID, req.FundID)
	
	// Add net income to equity
	netIncome, _ := s.calculateNetIncome(time.Date(req.AsOfDate.Year(), 1, 1, 0, 0, 0, 0, time.UTC), req.AsOfDate, req.BranchID, req.FundID)
	response.Equity.Lines = append(response.Equity.Lines, models.BalanceSheetLine{
		AccountCode: "NET_INCOME",
		AccountName: "Net Income (Current Year)",
		Amount:      netIncome,
		Level:       0,
		IsHeader:    false,
	})
	response.Equity.Total += netIncome
	response.TotalEquity = response.Equity.Total

	response.IsBalanced = (response.TotalAssets == (response.TotalLiabilities + response.TotalEquity))

	return response, nil
}

func (s *reportService) GetIncomeStatement(req *models.IncomeStatementRequest) (*models.IncomeStatementResponse, error) {
	// Get revenue accounts
	revenue, err := s.accountRepo.GetByCategory(models.AccountCategoryRevenue)
	if err != nil {
		return nil, err
	}

	// Get expense accounts
	expenses, err := s.accountRepo.GetByCategory(models.AccountCategoryExpense)
	if err != nil {
		return nil, err
	}

	response := &models.IncomeStatementResponse{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	// Revenue section
	response.Revenue = s.buildIncomeStatementSection(revenue, req.StartDate, req.EndDate, req.BranchID, req.FundID)
	response.TotalRevenue = response.Revenue.Total

	// Expenses section
	response.Expenses = s.buildIncomeStatementSection(expenses, req.StartDate, req.EndDate, req.BranchID, req.FundID)
	response.TotalExpenses = response.Expenses.Total

	// Net income
	response.NetIncome = response.TotalRevenue - response.TotalExpenses

	return response, nil
}

func (s *reportService) GetGeneralLedger(req *models.GeneralLedgerRequest) (*models.GeneralLedgerResponse, error) {
	// Get account
	account, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	// Calculate opening balance
	openingBalance, _ := s.calculateAccountBalance(req.AccountID, req.StartDate.AddDate(0, 0, -1), req.BranchID, nil, nil)

	// Get transactions
	var journalLines []models.JournalLine
	query := s.db.Model(&models.JournalLine{}).
		Joins("JOIN journals ON journals.id = journal_lines.journal_id").
		Where("journal_lines.account_id = ?", req.AccountID).
		Where("journals.journal_date BETWEEN ? AND ?", req.StartDate, req.EndDate).
		Where("journals.is_posted = ?", true).
		Order("journals.journal_date ASC, journals.created_at ASC")

	if req.BranchID != nil {
		query = query.Where("journals.branch_id = ?", *req.BranchID)
	}

	if err := query.Preload("Journal").Find(&journalLines).Error; err != nil {
		return nil, err
	}

	// Build response
	transactions := make([]models.GeneralLedgerLine, 0, len(journalLines))
	balance := openingBalance
	var totalDebit, totalCredit float64

	for _, line := range journalLines {
		// Calculate running balance
		if account.GetNormalBalance() == models.NormalBalanceDebit {
			balance += line.Debit - line.Credit
		} else {
			balance += line.Credit - line.Debit
		}

		totalDebit += line.Debit
		totalCredit += line.Credit

		transactions = append(transactions, models.GeneralLedgerLine{
			Date:          line.Journal.JournalDate,
			JournalNumber: line.Journal.JournalNumber,
			Description:   line.Description,
			Debit:         line.Debit,
			Credit:        line.Credit,
			Balance:       balance,
		})
	}

	return &models.GeneralLedgerResponse{
		Account:        *account.ToAccountResponse(),
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		OpeningBalance: openingBalance,
		Transactions:   transactions,
		TotalDebit:     totalDebit,
		TotalCredit:    totalCredit,
		ClosingBalance: balance,
	}, nil
}

// Helper functions

func (s *reportService) calculateAccountBalance(
	accountID uuid.UUID,
	asOfDate time.Time,
	branchID *uuid.UUID,
	fundID *uuid.UUID,
	programID *uuid.UUID,
) (float64, error) {
	var totalDebit, totalCredit float64

	query := s.db.Model(&models.JournalLine{}).
		Select("COALESCE(SUM(debit), 0) as total_debit, COALESCE(SUM(credit), 0) as total_credit").
		Joins("JOIN journals ON journals.id = journal_lines.journal_id").
		Where("journal_lines.account_id = ?", accountID).
		Where("journals.journal_date <= ?", asOfDate).
		Where("journals.is_posted = ?", true)

	if branchID != nil {
		query = query.Where("journals.branch_id = ?", *branchID)
	}
	if fundID != nil {
		query = query.Where("journal_lines.fund_id = ?", *fundID)
	}
	if programID != nil {
		query = query.Where("journal_lines.program_id = ?", *programID)
	}

	err := query.Row().Scan(&totalDebit, &totalCredit)
	if err != nil {
		return 0, err
	}

	// Get account to determine normal balance
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return 0, err
	}

	// Calculate balance based on normal balance
	if account.GetNormalBalance() == models.NormalBalanceDebit {
		return totalDebit - totalCredit, nil
	}
	return totalCredit - totalDebit, nil
}

func (s *reportService) buildBalanceSheetSection(
	accounts []models.Account,
	asOfDate time.Time,
	branchID *uuid.UUID,
	fundID *uuid.UUID,
) models.BalanceSheetSection {
	lines := make([]models.BalanceSheetLine, 0)
	var total float64

	for _, account := range accounts {
		if !account.IsDetail {
			continue
		}

		balance, err := s.calculateAccountBalance(account.ID, asOfDate, branchID, fundID, nil)
		if err != nil || balance == 0 {
			continue
		}

		lines = append(lines, models.BalanceSheetLine{
			AccountCode: account.Code,
			AccountName: account.Name,
			Amount:      balance,
			Level:       account.Level,
			IsHeader:    false,
		})

		total += balance
	}

	return models.BalanceSheetSection{
		Lines: lines,
		Total: total,
	}
}

func (s *reportService) buildIncomeStatementSection(
	accounts []models.Account,
	startDate, endDate time.Time,
	branchID *uuid.UUID,
	fundID *uuid.UUID,
) models.IncomeStatementSection {
	lines := make([]models.IncomeStatementLine, 0)
	var total float64

	for _, account := range accounts {
		if !account.IsDetail {
			continue
		}

		// Calculate period balance
		amount, err := s.calculatePeriodBalance(account.ID, startDate, endDate, branchID, fundID)
		if err != nil || amount == 0 {
			continue
		}

		lines = append(lines, models.IncomeStatementLine{
			AccountCode: account.Code,
			AccountName: account.Name,
			Amount:      amount,
			Level:       account.Level,
			IsHeader:    false,
		})

		total += amount
	}

	return models.IncomeStatementSection{
		Lines: lines,
		Total: total,
	}
}

func (s *reportService) calculatePeriodBalance(
	accountID uuid.UUID,
	startDate, endDate time.Time,
	branchID *uuid.UUID,
	fundID *uuid.UUID,
) (float64, error) {
	var totalDebit, totalCredit float64

	query := s.db.Model(&models.JournalLine{}).
		Select("COALESCE(SUM(debit), 0) as total_debit, COALESCE(SUM(credit), 0) as total_credit").
		Joins("JOIN journals ON journals.id = journal_lines.journal_id").
		Where("journal_lines.account_id = ?", accountID).
		Where("journals.journal_date BETWEEN ? AND ?", startDate, endDate).
		Where("journals.is_posted = ?", true)

	if branchID != nil {
		query = query.Where("journals.branch_id = ?", *branchID)
	}
	if fundID != nil {
		query = query.Where("journal_lines.fund_id = ?", *fundID)
	}

	err := query.Row().Scan(&totalDebit, &totalCredit)
	if err != nil {
		return 0, err
	}

	return totalCredit - totalDebit, nil
}

func (s *reportService) calculateNetIncome(
	startDate, endDate time.Time,
	branchID *uuid.UUID,
	fundID *uuid.UUID,
) (float64, error) {
	// Get revenue
	revenue, _ := s.accountRepo.GetByCategory(models.AccountCategoryRevenue)
	var totalRevenue float64
	for _, account := range revenue {
		if account.IsDetail {
			amount, _ := s.calculatePeriodBalance(account.ID, startDate, endDate, branchID, fundID)
			totalRevenue += amount
		}
	}

	// Get expenses
	expenses, _ := s.accountRepo.GetByCategory(models.AccountCategoryExpense)
	var totalExpenses float64
	for _, account := range expenses {
		if account.IsDetail {
			amount, _ := s.calculatePeriodBalance(account.ID, startDate, endDate, branchID, fundID)
			totalExpenses += amount
		}
	}

	return totalRevenue - totalExpenses, nil
}
