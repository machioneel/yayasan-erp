package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type StudentService interface {
	GetAll(params *models.PaginationParams) (*models.PaginationResponse, error)
	GetByID(id uuid.UUID) (*models.Student, error)
	GetByBranch(branchID uuid.UUID, params *models.PaginationParams) (*models.PaginationResponse, error)
	Search(keyword string, params *models.PaginationParams) (*models.PaginationResponse, error)
	Create(req *models.CreateStudentRequest) (*models.Student, error)
	Update(id uuid.UUID, req *models.CreateStudentRequest) (*models.Student, error)
	Delete(id uuid.UUID) error
	GetStatistics() (map[string]interface{}, error)
}

type studentService struct {
	studentRepo repository.StudentRepository
	parentRepo  repository.ParentRepository
	branchRepo  repository.BranchRepository
}

func NewStudentService(
	studentRepo repository.StudentRepository,
	parentRepo repository.ParentRepository,
	branchRepo repository.BranchRepository,
) StudentService {
	return &studentService{
		studentRepo: studentRepo,
		parentRepo:  parentRepo,
		branchRepo:  branchRepo,
	}
}

func (s *studentService) GetAll(params *models.PaginationParams) (*models.PaginationResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.AppConfig.App.DefaultPageSize
	}

	students, total, err := s.studentRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(students))
	for i, student := range students {
		items[i] = s.toStudentResponse(&student)
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

func (s *studentService) GetByID(id uuid.UUID) (*models.Student, error) {
	return s.studentRepo.GetByID(id)
}

func (s *studentService) GetByBranch(branchID uuid.UUID, params *models.PaginationParams) (*models.PaginationResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.AppConfig.App.DefaultPageSize
	}

	students, total, err := s.studentRepo.GetByBranch(branchID, params)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(students))
	for i, student := range students {
		items[i] = s.toStudentResponse(&student)
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

func (s *studentService) Search(keyword string, params *models.PaginationParams) (*models.PaginationResponse, error) {
	students, total, err := s.studentRepo.Search(keyword, params)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(students))
	for i, student := range students {
		items[i] = s.toStudentResponse(&student)
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

func (s *studentService) Create(req *models.CreateStudentRequest) (*models.Student, error) {
	// Validate branch
	branch, err := s.branchRepo.GetByID(req.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	// Check NISN uniqueness if provided
	if req.NISN != "" {
		existing, _ := s.studentRepo.GetByNISN(req.NISN)
		if existing != nil {
			return nil, errors.New("NISN already exists")
		}
	}

	// Generate registration number if not provided
	registrationNumber := req.RegistrationNumber
	if registrationNumber == "" {
		year := time.Now().Year()
		if req.RegistrationDate.Year() > 0 {
			year = req.RegistrationDate.Year()
		}
		registrationNumber, err = s.studentRepo.GenerateRegistrationNumber(branch.Code, year)
		if err != nil {
			return nil, err
		}
	}

	// Create student
	student := &models.Student{
		RegistrationNumber: registrationNumber,
		NISN:              req.NISN,
		FullName:          req.FullName,
		NickName:          req.NickName,
		Gender:            req.Gender,
		BirthPlace:        req.BirthPlace,
		BirthDate:         req.BirthDate,
		Religion:          req.Religion,
		Email:             req.Email,
		Phone:             req.Phone,
		Address:           req.Address,
		City:              req.City,
		BranchID:          req.BranchID,
		CurrentClassID:    req.CurrentClassID,
		Status:            models.StudentStatusActive,
		RegistrationDate:  req.RegistrationDate,
	}

	// Handle parents
	for _, parentReq := range req.Parents {
		var parent *models.Parent
		
		if parentReq.ParentID != nil {
			// Use existing parent
			parent, err = s.parentRepo.GetByID(*parentReq.ParentID)
			if err != nil {
				return nil, errors.New("parent not found")
			}
		} else {
			// Create new parent
			parent = &models.Parent{
				FullName:   parentReq.FullName,
				NIK:        parentReq.NIK,
				Gender:     parentReq.Gender,
				Phone:      parentReq.Phone,
				Email:      parentReq.Email,
				Address:    parentReq.Address,
				Occupation: parentReq.Occupation,
			}
			
			if err := s.parentRepo.Create(parent); err != nil {
				return nil, err
			}
		}

		// Create student-parent relationship
		student.Parents = append(student.Parents, models.StudentParent{
			ParentID:         parent.ID,
			Relationship:     parentReq.Relationship,
			IsPrimaryContact: parentReq.IsPrimary,
			IsFinancial:      parentReq.IsFinancial,
		})
	}

	if err := s.studentRepo.Create(student); err != nil {
		return nil, err
	}

	return s.studentRepo.GetByID(student.ID)
}

func (s *studentService) Update(id uuid.UUID, req *models.CreateStudentRequest) (*models.Student, error) {
	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("student not found")
	}

	// Update basic info
	student.FullName = req.FullName
	student.NickName = req.NickName
	student.Gender = req.Gender
	student.BirthPlace = req.BirthPlace
	student.BirthDate = req.BirthDate
	student.Religion = req.Religion
	student.Email = req.Email
	student.Phone = req.Phone
	student.Address = req.Address
	student.City = req.City
	student.CurrentClassID = req.CurrentClassID

	if err := s.studentRepo.Update(student); err != nil {
		return nil, err
	}

	return s.studentRepo.GetByID(student.ID)
}

func (s *studentService) Delete(id uuid.UUID) error {
	return s.studentRepo.Delete(id)
}

func (s *studentService) GetStatistics() (map[string]interface{}, error) {
	totalActive, _ := s.studentRepo.CountByStatus(models.StudentStatusActive)
	totalGraduated, _ := s.studentRepo.CountByStatus(models.StudentStatusGraduated)
	totalDropped, _ := s.studentRepo.CountByStatus(models.StudentStatusDropped)

	return map[string]interface{}{
		"total_active":    totalActive,
		"total_graduated": totalGraduated,
		"total_dropped":   totalDropped,
		"total":           totalActive + totalGraduated + totalDropped,
	}, nil
}

func (s *studentService) toStudentResponse(student *models.Student) *models.StudentResponse {
	resp := &models.StudentResponse{
		ID:                 student.ID,
		RegistrationNumber: student.RegistrationNumber,
		NISN:               student.NISN,
		FullName:           student.FullName,
		NickName:           student.NickName,
		Gender:             student.Gender,
		BirthPlace:         student.BirthPlace,
		BirthDate:          student.BirthDate,
		Religion:           student.Religion,
		Email:              student.Email,
		Phone:              student.Phone,
		Address:            student.Address,
		City:               student.City,
		BranchID:           student.BranchID,
		CurrentClassID:     student.CurrentClassID,
		Status:             student.Status,
		StatusName:         getStatusName(student.Status),
		RegistrationDate:   student.RegistrationDate,
		PhotoURL:           student.PhotoURL,
		CreatedAt:          student.CreatedAt,
		UpdatedAt:          student.UpdatedAt,
	}

	if student.Branch.ID != uuid.Nil {
		resp.BranchName = student.Branch.Name
	}

	if student.CurrentClass != nil {
		resp.ClassName = student.CurrentClass.Name
		resp.ClassLevel = student.CurrentClass.Level
	}

	if len(student.Parents) > 0 {
		resp.Parents = make([]models.ParentSummary, len(student.Parents))
		for i, sp := range student.Parents {
			resp.Parents[i] = models.ParentSummary{
				ID:           sp.Parent.ID,
				FullName:     sp.Parent.FullName,
				Relationship: sp.Relationship,
				Phone:        sp.Parent.Phone,
				IsPrimary:    sp.IsPrimaryContact,
				IsFinancial:  sp.IsFinancial,
			}
		}
	}

	return resp
}

func getStatusName(status string) string {
	names := map[string]string{
		models.StudentStatusActive:      "Active",
		models.StudentStatusGraduated:   "Graduated",
		models.StudentStatusDropped:     "Dropped Out",
		models.StudentStatusTransferred: "Transferred",
		models.StudentStatusSuspended:   "Suspended",
	}
	return names[status]
}
