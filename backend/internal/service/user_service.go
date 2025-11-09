package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/config"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
)

type UserService interface {
	GetAll(params *models.PaginationParams) (*models.UserListResponse, error)
	GetByID(id uuid.UUID) (*models.User, error)
	Create(req *models.CreateUserRequest, createdBy uuid.UUID) (*models.User, error)
	Update(id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error)
	Delete(id uuid.UUID) error
	UpdateProfile(id uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *userService) GetAll(params *models.PaginationParams) (*models.UserListResponse, error) {
	// Set default pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = config.AppConfig.App.DefaultPageSize
	}
	if params.PageSize > config.AppConfig.App.MaxPageSize {
		params.PageSize = config.AppConfig.App.MaxPageSize
	}

	users, total, err := s.userRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *user.ToUserResponse()
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &models.UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) GetByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUserWithRolesAndPermissions(id)
}

func (s *userService) Create(req *models.CreateUserRequest, createdBy uuid.UUID) (*models.User, error) {
	// Check if username exists
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email exists
	existingUser, _ = s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create user
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	isSuperAdmin := false
	if req.IsSuperAdmin != nil {
		isSuperAdmin = *req.IsSuperAdmin
	}

	user := &models.User{
		BranchID:     req.BranchID,
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		Phone:        req.Phone,
		IsActive:     isActive,
		IsSuperAdmin: isSuperAdmin,
	}

	// Hash password
	if err := user.SetPassword(req.Password); err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Assign roles if provided
	// Note: Role assignment would require UserRole repository
	// Implementation depends on your exact requirements

	return user, nil
}

func (s *userService) Update(id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields
	user.BranchID = req.BranchID
	user.FullName = req.FullName
	user.Phone = req.Phone
	user.AvatarURL = req.AvatarURL

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if req.IsSuperAdmin != nil {
		user.IsSuperAdmin = *req.IsSuperAdmin
	}

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}

func (s *userService) UpdateProfile(id uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.FullName = req.FullName
	user.Phone = req.Phone
	user.AvatarURL = req.AvatarURL

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
