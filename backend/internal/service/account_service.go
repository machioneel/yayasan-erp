package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type AccountService interface {
	GetAll(params *models.PaginationParams) (*models.AccountListResponse, error)
	GetByID(id uuid.UUID) (*models.Account, error)
	GetByCode(code string) (*models.Account, error)
	GetByCategory(category string) ([]models.Account, error)
	GetTree() ([]models.AccountTreeNode, error)
	GetActiveAccounts() ([]models.Account, error)
	GetDetailAccounts() ([]models.Account, error)
	Create(req *models.CreateAccountRequest) (*models.Account, error)
	Update(id uuid.UUID, req *models.UpdateAccountRequest) (*models.Account, error)
	Delete(id uuid.UUID) error
	ImportFromJSON(filepath string) (int, error)
	BulkImport(req *models.ImportAccountRequest) (int, error)
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

func (s *accountService) GetAll(params *models.PaginationParams) (*models.AccountListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.GlobalConfig.App.DefaultPageSize
	}
	if params.PageSize > config.GlobalConfig.App.MaxPageSize {
		params.PageSize = config.GlobalConfig.App.MaxPageSize
	}

	accounts, total, err := s.accountRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	// Convert to response
	accountResponses := make([]models.AccountResponse, len(accounts))
	for i, account := range accounts {
		accountResponses[i] = *account.ToAccountResponse()
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &models.AccountListResponse{
		Accounts:   accountResponses,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *accountService) GetByID(id uuid.UUID) (*models.Account, error) {
	return s.accountRepo.GetByID(id)
}

func (s *accountService) GetByCode(code string) (*models.Account, error) {
	return s.accountRepo.GetByCode(code)
}

func (s *accountService) GetByCategory(category string) ([]models.Account, error) {
	return s.accountRepo.GetByCategory(category)
}

func (s *accountService) GetTree() ([]models.AccountTreeNode, error) {
	accounts, err := s.accountRepo.GetTree()
	if err != nil {
		return nil, err
	}

	// Convert to tree nodes
	nodes := make([]models.AccountTreeNode, len(accounts))
	for i, account := range accounts {
		nodes[i] = s.convertToTreeNode(&account)
	}

	return nodes, nil
}

func (s *accountService) convertToTreeNode(account *models.Account) models.AccountTreeNode {
	node := models.AccountTreeNode{
		ID:       account.ID,
		Code:     account.Code,
		Name:     account.Name,
		Type:     account.Type,
		Category: account.Category,
		IsActive: account.IsActive,
		IsDetail: account.IsDetail,
		Level:    account.Level,
	}

	if len(account.Children) > 0 {
		node.Children = make([]models.AccountTreeNode, len(account.Children))
		for i, child := range account.Children {
			node.Children[i] = s.convertToTreeNode(&child)
		}
	}

	return node
}

func (s *accountService) GetActiveAccounts() ([]models.Account, error) {
	return s.accountRepo.GetActiveAccounts()
}

func (s *accountService) GetDetailAccounts() ([]models.Account, error) {
	return s.accountRepo.GetDetailAccounts()
}

func (s *accountService) Create(req *models.CreateAccountRequest) (*models.Account, error) {
	// Check if code exists
	exists, err := s.accountRepo.CodeExists(req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("account code already exists")
	}

	// Validate parent if provided
	if req.ParentID != nil {
		_, err := s.accountRepo.GetByID(*req.ParentID)
		if err != nil {
			return nil, errors.New("parent account not found")
		}
	}

	// Create account
	account := &models.Account{
		ParentID:      req.ParentID,
		Code:          req.Code,
		Name:          req.Name,
		NameEn:        req.NameEn,
		Type:          req.Type,
		Category:      req.Category,
		NormalBalance: req.NormalBalance,
		IsActive:      true,
		Description:   req.Description,
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountService) Update(id uuid.UUID, req *models.UpdateAccountRequest) (*models.Account, error) {
	account, err := s.accountRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("account not found")
	}

	// Update fields
	account.Name = req.Name
	account.NameEn = req.NameEn
	account.Description = req.Description

	if req.IsActive != nil {
		account.IsActive = *req.IsActive
	}

	if err := s.accountRepo.Update(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountService) Delete(id uuid.UUID) error {
	return s.accountRepo.Delete(id)
}

func (s *accountService) ImportFromJSON(filepath string) (int, error) {
	// Read JSON file
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return 0, errors.New("failed to read file: " + err.Error())
	}

	// Parse JSON
	var coaData []struct {
		No     int    `json:"no"`
		Code   string `json:"code"`
		Type   string `json:"type"`
		Name   string `json:"name"`
		Level1 string `json:"level1"`
		Level2 string `json:"level2"`
		Level3 string `json:"level3"`
		Level4 string `json:"level4"`
	}

	if err := json.Unmarshal(data, &coaData); err != nil {
		return 0, errors.New("failed to parse JSON: " + err.Error())
	}

	// Convert to accounts
	accounts := make([]models.Account, 0, len(coaData))
	codeToID := make(map[string]uuid.UUID)

	for _, item := range coaData {
		// Build name from levels
		name := item.Name
		if name == "" {
			nameParts := []string{}
			for _, level := range []string{item.Level1, item.Level2, item.Level3, item.Level4} {
				if level != "" {
					nameParts = append(nameParts, level)
				}
			}
			if len(nameParts) > 0 {
				name = nameParts[len(nameParts)-1]
			}
		}

		// Determine category from code
		category := s.getCategoryFromCode(item.Code)

		// Determine parent
		var parentID *uuid.UUID
		if len(item.Code) > 4 {
			// Has parent - get parent code
			parentCode := item.Code[:4]
			if pid, exists := codeToID[parentCode]; exists {
				parentID = &pid
			}
		}

		// Create account
		account := models.Account{
			BaseModel: models.BaseModel{
				ID: uuid.New(),
			},
			ParentID:  parentID,
			Code:      item.Code,
			Name:      name,
			Type:      item.Type,
			Category:  category,
			IsActive:  true,
		}

		codeToID[item.Code] = account.ID
		accounts = append(accounts, account)
	}

	// Bulk create
	if err := s.accountRepo.BulkCreate(accounts); err != nil {
		return 0, err
	}

	return len(accounts), nil
}

func (s *accountService) BulkImport(req *models.ImportAccountRequest) (int, error) {
	accounts := make([]models.Account, 0, len(req.Accounts))

	for _, reqAcc := range req.Accounts {
		// Check if code exists
		exists, err := s.accountRepo.CodeExists(reqAcc.Code)
		if err != nil {
			return 0, err
		}
		if exists {
			continue // Skip existing
		}

		account := models.Account{
			ParentID:      reqAcc.ParentID,
			Code:          reqAcc.Code,
			Name:          reqAcc.Name,
			NameEn:        reqAcc.NameEn,
			Type:          reqAcc.Type,
			Category:      reqAcc.Category,
			NormalBalance: reqAcc.NormalBalance,
			IsActive:      true,
			Description:   reqAcc.Description,
		}

		accounts = append(accounts, account)
	}

	if len(accounts) == 0 {
		return 0, errors.New("no new accounts to import")
	}

	if err := s.accountRepo.BulkCreate(accounts); err != nil {
		return 0, err
	}

	return len(accounts), nil
}

func (s *accountService) getCategoryFromCode(code string) string {
	if len(code) == 0 {
		return ""
	}

	firstDigit := code[0]
	switch firstDigit {
	case '1':
		return models.AccountCategoryAsset
	case '2':
		return models.AccountCategoryLiability
	case '3':
		return models.AccountCategoryEquity
	case '4':
		return models.AccountCategoryRevenue
	case '5', '6', '7', '8', '9':
		return models.AccountCategoryExpense
	default:
		return ""
	}
}
