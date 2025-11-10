package service

import (
	"errors"
	//"time"

	"github.com/google/uuid"
	"github.com/yayasan/erp-backend/internal/models"
	"github.com/yayasan/erp-backend/internal/repository"
	"github.com/yayasan/erp-backend/internal/utils"
)

type AuthService interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	RefreshToken(refreshToken string) (string, error)
	ChangePassword(userID uuid.UUID, req *models.ChangePasswordRequest) error
	GetCurrentUser(userID uuid.UUID) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *authService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if username already exists
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, _ = s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create new user
	user := &models.User{
		BranchID: req.BranchID,
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
		IsActive: true,
		IsSuperAdmin: false,
	}

	// Hash password
	if err := user.SetPassword(req.Password); err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Find user by username or email
	user, err := s.userRepo.GetByUsernameOrEmail(req.Login)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Get user with roles and permissions
	userWithPerms, err := s.userRepo.GetUserWithRolesAndPermissions(user.ID)
	if err != nil {
		return nil, errors.New("failed to load user permissions")
	}

	// Collect permissions
	permissionSet := make(map[string]bool)
	for _, userRole := range userWithPerms.UserRoles {
		if userRole.Role != nil {
			for _, rolePerm := range userRole.Role.RolePermissions {
				if rolePerm.Permission != nil {
					permissionSet[rolePerm.Permission.Code] = true
				}
			}
		}
	}

	permissions := make([]string, 0, len(permissionSet))
	for perm := range permissionSet {
		permissions = append(permissions, perm)
	}

	// Generate access token
	accessToken, err := utils.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		user.BranchID,
		user.IsSuperAdmin,
		permissions,
	)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		// Log error but don't fail login
	}

	response := &models.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         userWithPerms.ToUserResponse(),
	}

	return response, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Check token type
	if claims.TokenType != "refresh" {
		return "", errors.New("invalid token type")
	}

	// Get user with permissions
	user, err := s.userRepo.GetUserWithRolesAndPermissions(claims.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return "", errors.New("user account is deactivated")
	}

	// Collect permissions
	permissionSet := make(map[string]bool)
	for _, userRole := range user.UserRoles {
		if userRole.Role != nil {
			for _, rolePerm := range userRole.Role.RolePermissions {
				if rolePerm.Permission != nil {
					permissionSet[rolePerm.Permission.Code] = true
				}
			}
		}
	}

	permissions := make([]string, 0, len(permissionSet))
	for perm := range permissionSet {
		permissions = append(permissions, perm)
	}

	// Generate new access token
	newAccessToken, err := utils.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		user.BranchID,
		user.IsSuperAdmin,
		permissions,
	)
	if err != nil {
		return "", errors.New("failed to generate new access token")
	}

	return newAccessToken, nil
}

func (s *authService) ChangePassword(userID uuid.UUID, req *models.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !user.CheckPassword(req.OldPassword) {
		return errors.New("old password is incorrect")
	}

	// Set new password
	if err := user.SetPassword(req.NewPassword); err != nil {
		return errors.New("failed to hash new password")
	}

	// Update user
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (s *authService) GetCurrentUser(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserWithRolesAndPermissions(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
