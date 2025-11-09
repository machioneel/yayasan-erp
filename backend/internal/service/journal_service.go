package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type JournalService interface {
	GetAll(params *models.PaginationParams) (*models.JournalListResponse, error)
	GetByID(id uuid.UUID) (*models.Journal, error)
	GetByStatus(status string, params *models.PaginationParams) (*models.JournalListResponse, error)
	Create(req *models.CreateJournalRequest, userID uuid.UUID) (*models.Journal, error)
	Update(id uuid.UUID, req *models.UpdateJournalRequest, userID uuid.UUID) (*models.Journal, error)
	Delete(id uuid.UUID, userID uuid.UUID) error
	SubmitForReview(id uuid.UUID, userID uuid.UUID) (*models.Journal, error)
	Review(id uuid.UUID, req *models.ReviewJournalRequest, userID uuid.UUID) (*models.Journal, error)
	Post(id uuid.UUID, req *models.PostJournalRequest, userID uuid.UUID) (*models.Journal, error)
	Unpost(id uuid.UUID, userID uuid.UUID) (*models.Journal, error)
}

type journalService struct {
	journalRepo repository.JournalRepository
	accountRepo repository.AccountRepository
	branchRepo  repository.BranchRepository
}

func NewJournalService(
	journalRepo repository.JournalRepository,
	accountRepo repository.AccountRepository,
	branchRepo repository.BranchRepository,
) JournalService {
	return &journalService{
		journalRepo: journalRepo,
		accountRepo: accountRepo,
		branchRepo:  branchRepo,
	}
}

func (s *journalService) GetAll(params *models.PaginationParams) (*models.JournalListResponse, error) {
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

	journals, total, err := s.journalRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	// Convert to response
	journalResponses := make([]models.JournalResponse, len(journals))
	for i, journal := range journals {
		journalResponses[i] = *journal.ToJournalResponse()
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &models.JournalListResponse{
		Journals:   journalResponses,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *journalService) GetByID(id uuid.UUID) (*models.Journal, error) {
	return s.journalRepo.GetByID(id)
}

func (s *journalService) GetByStatus(status string, params *models.PaginationParams) (*models.JournalListResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.AppConfig.App.DefaultPageSize
	}

	journals, total, err := s.journalRepo.GetByStatus(status, params)
	if err != nil {
		return nil, err
	}

	journalResponses := make([]models.JournalResponse, len(journals))
	for i, journal := range journals {
		journalResponses[i] = *journal.ToJournalResponse()
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &models.JournalListResponse{
		Journals:   journalResponses,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *journalService) Create(req *models.CreateJournalRequest, userID uuid.UUID) (*models.Journal, error) {
	// Validate branch exists
	branch, err := s.branchRepo.GetByID(req.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	// Validate journal lines
	if err := s.validateJournalLines(req.JournalLines); err != nil {
		return nil, err
	}

	// Check balance
	var totalDebit, totalCredit float64
	for _, line := range req.JournalLines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	if totalDebit != totalCredit {
		return nil, errors.New("journal is not balanced: debit != credit")
	}

	// Generate journal number
	journalNumber, err := s.journalRepo.GenerateJournalNumber(branch.Code, req.JournalDate)
	if err != nil {
		return nil, err
	}

	// Create journal
	journal := &models.Journal{
		BranchID:      req.BranchID,
		JournalNumber: journalNumber,
		JournalDate:   req.JournalDate,
		Description:   req.Description,
		ReferenceNo:   req.ReferenceNo,
		Status:        models.JournalStatusDraft,
		TotalDebit:    totalDebit,
		TotalCredit:   totalCredit,
		CreatedBy:     userID,
	}

	// Add journal lines
	journal.JournalLines = make([]models.JournalLine, len(req.JournalLines))
	for i, lineReq := range req.JournalLines {
		journal.JournalLines[i] = models.JournalLine{
			AccountID:   lineReq.AccountID,
			Description: lineReq.Description,
			Debit:       lineReq.Debit,
			Credit:      lineReq.Credit,
			FundID:      lineReq.FundID,
			ProgramID:   lineReq.ProgramID,
			DonorID:     lineReq.DonorID,
		}
	}

	if err := s.journalRepo.Create(journal); err != nil {
		return nil, err
	}

	return s.journalRepo.GetByID(journal.ID)
}

func (s *journalService) Update(id uuid.UUID, req *models.UpdateJournalRequest, userID uuid.UUID) (*models.Journal, error) {
	journal, err := s.journalRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("journal not found")
	}

	// Can only update draft journals
	if journal.Status != models.JournalStatusDraft {
		return nil, errors.New("can only update draft journals")
	}

	// Only creator can update
	if journal.CreatedBy != userID {
		return nil, errors.New("only creator can update this journal")
	}

	// Validate lines
	if err := s.validateJournalLines(req.JournalLines); err != nil {
		return nil, err
	}

	// Check balance
	var totalDebit, totalCredit float64
	for _, line := range req.JournalLines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	if totalDebit != totalCredit {
		return nil, errors.New("journal is not balanced")
	}

	// Update journal
	journal.JournalDate = req.JournalDate
	journal.Description = req.Description
	journal.ReferenceNo = req.ReferenceNo
	journal.TotalDebit = totalDebit
	journal.TotalCredit = totalCredit

	// Update lines
	journal.JournalLines = make([]models.JournalLine, len(req.JournalLines))
	for i, lineReq := range req.JournalLines {
		journal.JournalLines[i] = models.JournalLine{
			JournalID:   journal.ID,
			AccountID:   lineReq.AccountID,
			Description: lineReq.Description,
			Debit:       lineReq.Debit,
			Credit:      lineReq.Credit,
			FundID:      lineReq.FundID,
			ProgramID:   lineReq.ProgramID,
			DonorID:     lineReq.DonorID,
		}
	}

	if err := s.journalRepo.Update(journal); err != nil {
		return nil, err
	}

	return s.journalRepo.GetByID(journal.ID)
}

func (s *journalService) Delete(id uuid.UUID, userID uuid.UUID) error {
	journal, err := s.journalRepo.GetByID(id)
	if err != nil {
		return errors.New("journal not found")
	}

	// Can only delete draft journals
	if journal.Status != models.JournalStatusDraft {
		return errors.New("can only delete draft journals")
	}

	// Only creator can delete
	if journal.CreatedBy != userID {
		return errors.New("only creator can delete this journal")
	}

	return s.journalRepo.Delete(id)
}

func (s *journalService) SubmitForReview(id uuid.UUID, userID uuid.UUID) (*models.Journal, error) {
	journal, err := s.journalRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("journal not found")
	}

	// Must be draft
	if journal.Status != models.JournalStatusDraft {
		return nil, errors.New("can only submit draft journals")
	}

	// Must be creator
	if journal.CreatedBy != userID {
		return nil, errors.New("only creator can submit this journal")
	}

	// Must be balanced
	if !journal.IsBalanced() {
		return nil, errors.New("journal must be balanced before submission")
	}

	// Update status
	journal.Status = models.JournalStatusReview
	
	if err := s.journalRepo.Update(journal); err != nil {
		return nil, err
	}

	return s.journalRepo.GetByID(journal.ID)
}

func (s *journalService) Review(id uuid.UUID, req *models.ReviewJournalRequest, userID uuid.UUID) (*models.Journal, error) {
	journal, err := s.journalRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("journal not found")
	}

	// Must be in review status
	if journal.Status != models.JournalStatusReview {
		return nil, errors.New("journal is not in review status")
	}

	// Cannot review own journal
	if journal.CreatedBy == userID {
		return nil, errors.New("cannot review your own journal")
	}

	now := time.Now()

	if req.Action == "approve" {
		// Approve
		journal.Status = models.JournalStatusApproved
		journal.ApprovedBy = &userID
		journal.ApprovedAt = &now
	} else {
		// Reject
		journal.Status = models.JournalStatusRejected
		journal.RejectedBy = &userID
		journal.RejectedAt = &now
		journal.RejectReason = req.Notes
	}

	if err := s.journalRepo.Update(journal); err != nil {
		return nil, err
	}

	return s.journalRepo.GetByID(journal.ID)
}

func (s *journalService) Post(id uuid.UUID, req *models.PostJournalRequest, userID uuid.UUID) (*models.Journal, error) {
	journal, err := s.journalRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("journal not found")
	}

	// Must be approved
	if !journal.CanPost() {
		return nil, errors.New("journal cannot be posted")
	}

	// Post journal
	now := time.Now()
	journal.IsPosted = true
	journal.PostedAt = &now
	journal.PostedBy = &userID
	journal.Status = models.JournalStatusPosted

	if err := s.journalRepo.Update(journal); err != nil {
		return nil, err
	}

	// TODO: Update account balances (implement in next iteration)

	return s.journalRepo.GetByID(journal.ID)
}

func (s *journalService) Unpost(id uuid.UUID, userID uuid.UUID) (*models.Journal, error) {
	journal, err := s.journalRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("journal not found")
	}

	// Must be posted
	if !journal.IsPosted {
		return nil, errors.New("journal is not posted")
	}

	// Unpost journal
	journal.IsPosted = false
	journal.PostedAt = nil
	journal.PostedBy = nil
	journal.Status = models.JournalStatusApproved

	if err := s.journalRepo.Update(journal); err != nil {
		return nil, err
	}

	// TODO: Reverse account balances (implement in next iteration)

	return s.journalRepo.GetByID(journal.ID)
}

func (s *journalService) validateJournalLines(lines []models.CreateJournalLineReq) error {
	if len(lines) < 2 {
		return errors.New("journal must have at least 2 lines")
	}

	for i, line := range lines {
		// Validate account exists and can post
		account, err := s.accountRepo.GetByID(line.AccountID)
		if err != nil {
			return errors.New("account not found on line " + string(rune(i+1)))
		}

		if !account.CanPostTransaction() {
			return errors.New("account " + account.Code + " cannot have transactions")
		}

		// Either debit or credit, not both
		if line.Debit > 0 && line.Credit > 0 {
			return errors.New("line cannot have both debit and credit")
		}

		if line.Debit == 0 && line.Credit == 0 {
			return errors.New("line must have either debit or credit")
		}
	}

	return nil
}
